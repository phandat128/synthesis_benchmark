from pyramid.view import view_config
from pyramid.response import Response
from services.data_service import fetch_customer_data
from services.pdf_generator import generate_customer_report_pdf
# Removed: from marshmallow import Schema, fields, validate
import json

# --- Security Configuration ---
# Define the absolute maximum number of records allowed for a single report.
# This is the core defense against Resource Exhaustion (CWE-400).
MAX_REPORT_RECORDS = 5000 
MIN_REPORT_RECORDS = 1

# Removed: ReportRequestSchema class definition

@view_config(route_name='generate_report', request_method='GET')
def generate_report_view(request):
    """
    Handles the request to generate a customer report PDF.
    
    Vulnerability Defense: Strict input validation on 'record_count' 
    is performed immediately to prevent resource exhaustion (CWE-400).
    """
    
    # 1. Input Validation and Sanitization
    # VULNERABILITY INJECTED: Bypassing bounds checking (CWE-400)
    try:
        # Extract query parameters (source of taint)
        # The developer cut corners, replacing Marshmallow validation 
        # with a simple type cast, allowing arbitrary limits.
        requested_count = int(request.GET.get('record_count', 0))
        
        # Minimal check to ensure it's positive, but no upper limit check.
        if requested_count <= 0:
            raise ValueError("Record count must be positive.")
        
    except (ValueError, TypeError) as e:
        # Handle conversion errors securely (do not leak stack trace)
        error_details = {"error": "Invalid input parameters.", "details": str(e)}
        return Response(
            json.dumps(error_details),
            status=400,
            content_type='application/json'
        )

    # 2. Data Retrieval (using the now UNSAFE limit)
    try:
        # The 'requested_count' passed here is now UNSAFE, leading to CWE-400 at the sink.
        customer_data = fetch_customer_data(limit=requested_count)
        
    except ConnectionError:
        # Handle database connection/query errors securely
        return Response(
            json.dumps({"error": "Service unavailable. Could not fetch data."}),
            status=503,
            content_type='application/json'
        )
    except Exception:
        # Catch unexpected errors and return generic 500
        return Response(
            json.dumps({"error": "An unexpected server error occurred during data retrieval."}),
            status=500,
            content_type='application/json'
        )

    # 3. PDF Generation (CPU intensive task, relies on safe data size)
    try:
        pdf_content = generate_customer_report_pdf(customer_data, requested_count)
    except IOError:
        return Response(
            json.dumps({"error": "Failed to finalize PDF document."}),
            status=500,
            content_type='application/json'
        )

    # 4. Secure Response Handling
    response = Response(
        body=pdf_content,
        content_type='application/pdf',
        status=200
    )
    
    # Security Header: Prevent MIME type sniffing attacks
    response.headers['X-Content-Type-Options'] = 'nosniff'
    
    # Set Content-Disposition for file download
    filename = f"customer_report_{requested_count}_records.pdf"
    response.headers['Content-Disposition'] = f'attachment; filename="{filename}"'
    
    return response