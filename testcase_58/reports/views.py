import logging
from rest_framework import status
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.exceptions import APIException
from rest_framework.views import exception_handler

from reports.serializers import ReportSerializer
from reports.tasks import generate_report_task

logger = logging.getLogger(__name__)

def custom_exception_handler(exc, context):
    """
    Custom exception handler to prevent leaking internal server details (e.g., stack traces).
    """
    response = exception_handler(exc, context)

    if response is not None:
        # Standard DRF errors (400, 403, 404) are handled normally
        return response

    # Handle unexpected internal server errors (500)
    if isinstance(exc, Exception):
        logger.error(f"Unhandled internal server error: {exc}", exc_info=True)
        
        # Return a generic, safe error message (no sensitive data leakage)
        return Response(
            {'detail': 'An internal server error occurred. Please try again later.'},
            status=status.HTTP_500_INTERNAL_SERVER_ERROR
        )
        
    return response


class ReportGenerationView(APIView):
    """
    API endpoint to trigger asynchronous PDF report generation.
    """
    http_method_names = ['post']

    def post(self, request):
        # 1. Input Validation and Sanitization
        serializer = ReportSerializer(data=request.data)
        
        if not serializer.is_valid():
            # Returns 400 Bad Request with validation details
            return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)

        validated_data = serializer.validated_data
        
        record_count = validated_data['record_count']
        report_title = validated_data['report_title']
        active_only = validated_data['active_only']
        
        # 2. Authorization Check (User must be authenticated per settings)
        user_id = request.user.id if request.user.is_authenticated else None
        if not user_id:
             raise APIException("Authentication required.") 

        # 3. Asynchronous Task Dispatch
        try:
            # Dispatch the task with validated, bounded parameters
            task = generate_report_task.delay(
                record_count=record_count,
                report_title=report_title,
                active_only=active_only,
                user_id=user_id
            )
            
            logger.info(f"Report generation task dispatched: {task.id}")

            return Response({
                'message': 'Report generation started successfully.',
                'task_id': task.id,
                'note': 'Check status endpoint for completion.'
            }, status=status.HTTP_202_ACCEPTED) # 202 Accepted for async processing

        except Exception as e:
            # Handle Celery connection issues (Service Unavailable)
            logger.error(f"Failed to dispatch Celery task: {e}")
            return Response(
                {'detail': 'Service temporarily unavailable. Could not start report process.'},
                status=status.HTTP_503_SERVICE_UNAVAILABLE
            )
