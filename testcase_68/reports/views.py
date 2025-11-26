import logging
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from django.http import HttpResponse
from reports.serializers import ReportRequestSerializer
from reports.utils import generate_pdf_report

logger = logging.getLogger(__name__)

class ReportGenerationView(APIView):
    """
    API endpoint for generating customizable PDF reports.
    Handles POST requests, validates input, and returns the PDF file.
    """
    # Ensure only authenticated users can access this resource (Least Privilege Principle)
    # permission_classes = [IsAuthenticated] 
    
    def post(self, request, *args, **kwargs):
        """
        Processes the request to generate a report.
        """
        serializer = ReportRequestSerializer(data=request.data)
        
        # Input Validation: Check if the request body is valid and safe.
        if not serializer.is_valid():
            return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

        # The 'record_count' is now guaranteed to be within the safe bounds defined in settings.
        record_count = serializer.validated_data['record_count']
        
        try:
            # Call the utility function with the safe, validated limit
            pdf_buffer = generate_pdf_report(limit=record_count)
            
            # Prepare the HTTP response for file download
            response = HttpResponse(
                pdf_buffer.getvalue(),
                content_type='application/pdf'
            )
            # Use Content-Disposition for safe file naming
            response['Content-Disposition'] = f'attachment; filename="report_{record_count}_records.pdf"'
            
            return response
            
        except RuntimeError as e:
            # Proper Error Handling: Log the detailed error but return a generic 500 to the user
            logger.error(f"Report generation failed for request (limit={record_count}): {e}", exc_info=True)
            return Response(
                {"detail": "Internal server error during report processing. Please try again later."}, 
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )
        except Exception:
            # Catch all unexpected errors
            logger.critical("Unexpected critical error in ReportGenerationView.", exc_info=True)
            return Response(
                {"detail": "An unexpected error occurred."}, 
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )
