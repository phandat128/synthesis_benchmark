import logging
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework.exceptions import APIException
from django.http import HttpResponse

from .serializers import ReportRequestSerializer
from .services import generate_pdf_report, ReportGenerationError

logger = logging.getLogger(__name__)

# --- Custom Exception Handler for DRF ---
def custom_exception_handler(exc, context):
    """
    Custom handler to prevent leaking internal server details (CWE-209).
    Only expose necessary validation errors or generic server errors.
    """
    from rest_framework.views import exception_handler
    
    response = exception_handler(exc, context)

    if response is not None:
        # Handle standard DRF exceptions (e.g., validation errors, 404s)
        return response

    # Handle custom service layer errors (e.g., data not found, resource limits)
    if isinstance(exc, ReportGenerationError):
        logger.warning(f"Service error encountered: {exc}")
        # Use 400 for errors stemming from bad user input/request constraints
        return Response(
            {'detail': str(exc)}, 
            status=status.HTTP_400_BAD_REQUEST
        )
        
    # Handle unexpected internal errors (CWE-755)
    logger.error(f"Unhandled internal server error: {exc}", exc_info=True)
    return Response(
        {'detail': 'An unexpected server error occurred. Please try again later.'},
        status=status.HTTP_500_INTERNAL_SERVER_ERROR
    )


class ReportGeneratorView(APIView):
    """
    API endpoint to trigger the generation of a PDF report.
    Requires authentication.
    """
    
    def post(self, request):
        """
        Handles POST requests for report generation.
        """
        # 1. Input Validation and Sanitization
        serializer = ReportRequestSerializer(data=request.data)
        
        # This step enforces the MaxValueValidator, mitigating CWE-400 DoS.
        if not serializer.is_valid():
            return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

        validated_data = serializer.validated_data
        
        report_type = validated_data['report_type']
        limit = validated_data['record_limit']
        start_date = validated_data.get('start_date')
        end_date = validated_data.get('end_date')

        try:
            # 2. Business Logic Execution
            pdf_buffer = generate_pdf_report(
                report_type=report_type,
                limit=limit,
                start_date=start_date,
                end_date=end_date
            )
            
            # 3. Secure Response Handling
            # Filename derived from validated input (report_type is regex-validated)
            filename = f"{report_type}_report.pdf"
            
            response = HttpResponse(
                pdf_buffer.getvalue(),
                content_type='application/pdf'
            )
            
            # Use Content-Disposition for secure file download
            response['Content-Disposition'] = f'attachment; filename="{filename}"'
            
            return response

        except ReportGenerationError as e:
            # Raise the service error to be caught by the custom handler
            raise e
        except Exception as e:
            # Catch any remaining unexpected exceptions
            logger.error(f"Critical unexpected error in ReportGeneratorView: {e}", exc_info=True)
            # Raise a generic API exception which will trigger the 500 response
            raise APIException("Internal server error during report generation.")