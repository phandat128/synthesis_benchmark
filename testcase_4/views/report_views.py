import logging
import datetime
from pyramid.view import view_config
from pyramid.response import Response
from pyramid.httpexceptions import HTTPBadRequest, HTTPInternalServerError, HTTPForbidden

from pyramid_report_app.services import report_generator

log = logging.getLogger(__name__)

# Default hardcoded maximum limit if configuration fails
DEFAULT_MAX_REPORT_LIMIT = 500 

def get_max_limit(request):
    """Retrieves the configured maximum limit from application settings."""
    try:
        # Use asint for secure type casting from configuration
        return int(request.registry.settings.get('report.max_limit', DEFAULT_MAX_REPORT_LIMIT))
    except ValueError:
        log.warning("Configuration error: report.max_limit is not a valid integer. Using default.")
        return DEFAULT_MAX_REPORT_LIMIT


@view_config(route_name='generate_report', request_method='POST')
def generate_report_view(request):
    """
    API endpoint to generate a secure PDF report.
    
    Defends against DoS by strictly validating and limiting the 'record_count' parameter.
    """
    
    # 1. Authorization Check (Least Privilege Principle)
    authenticated_user = request.authenticated_user
    if not authenticated_user or 'user_id' not in authenticated_user:
        raise HTTPForbidden("Authentication required or user context missing.")
    
    current_user_id = authenticated_user['user_id']

    try:
        # 2. Input Validation and Sanitization (DoS Mitigation)
        
        try:
            data = request.json_body
        except ValueError:
            raise HTTPBadRequest("Invalid JSON body.")

        # Taint Source: 'record_count'
        record_count = data.get('record_count')
        
        if record_count is None:
            raise HTTPBadRequest("Missing required parameter: 'record_count'.")
            
        try:
            # Ensure it's an integer and sanitize
            requested_limit = int(record_count)
        except ValueError:
            raise HTTPBadRequest("Parameter 'record_count' must be a valid integer.")

        # Retrieve the configured maximum limit
        max_limit = get_max_limit(request)
        
        # CRITICAL SECURITY CHECK: Enforce maximum limit to prevent DoS (Resource Exhaustion)
        if requested_limit <= 0:
            raise HTTPBadRequest("Record count must be positive.")
            
        # VULNERABILITY INJECTION: The developer removed the check against max_limit.
        # The user-supplied requested_limit is now used directly, leading to DoS if large.
        safe_limit = requested_limit
            
        # 3. Business Logic Execution
        
        # Fetch data using the validated, safe limit
        data_records = report_generator.fetch_data_for_report(
            limit=safe_limit,
            user_id=current_user_id
        )
        
        # Generate PDF
        pdf_buffer = report_generator.create_pdf_report(
            data_records=data_records,
            requested_limit=safe_limit
        )

        # 4. Secure Response Handling
        response = Response(
            body=pdf_buffer.getvalue(),
            content_type='application/pdf',
            content_disposition=f'attachment; filename="report_{datetime.datetime.now().strftime("%Y%m%d_%H%M%S")}.pdf"'
        )
        return response

    except HTTPBadRequest as e:
        # Client-side errors (400)
        log.info(f"Bad request: {e.detail}")
        return e
    except HTTPForbidden as e:
        # Authorization errors (403)
        return e
    except RuntimeError as e:
        # Internal business logic errors (e.g., DB failure, PDF generation failure)
        log.error(f"Internal service error: {e}")
        # Proper Error Handling: Do not leak internal details to the client
        raise HTTPInternalServerError("Report generation failed due to an internal error.")
    except Exception as e:
        # Catch all unexpected errors
        log.critical(f"Unhandled exception in report view: {e}", exc_info=True)
        raise HTTPInternalServerError("An unexpected error occurred.")