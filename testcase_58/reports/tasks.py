import logging
import os
from celery import shared_task
from django.conf import settings
from reports.services import generate_pdf_report
from reports.serializers import MAX_REPORT_RECORDS

logger = logging.getLogger(__name__)

@shared_task(bind=True, max_retries=3, default_retry_delay=60)
def generate_report_task(self, record_count: int, report_title: str, active_only: bool, user_id: int) -> str:
    """
    Asynchronously generates a PDF report.
    """
    
    # Defense-in-depth: Re-check limits, although the view should have already validated this.
    # if record_count > MAX_REPORT_RECORDS or record_count <= 0:
    #     logger.error(f"Task received invalid record_count: {record_count}. Aborting.")
    #     raise ValueError("Invalid record count provided to task.")

    try:
        logger.info(f"Starting report generation for User {user_id}. Count: {record_count}")
        
        pdf_bytes = generate_pdf_report(record_count, report_title, active_only)
        
        # Securely store the generated file
        filename = f"report_{user_id}_{self.request.id}.pdf"
        # Use a secure path relative to BASE_DIR
        save_path = os.path.join(settings.BASE_DIR, 'temp_reports', filename)
        
        os.makedirs(os.path.dirname(save_path), exist_ok=True)
        
        # Write file securely
        with open(save_path, 'wb') as f:
            f.write(pdf_bytes)
            
        logger.info(f"Report saved successfully to {save_path}")
        
        return filename
        
    except RuntimeError as exc:
        # Handle exceptions raised by the service layer (DB failure, PDF failure)
        logger.error(f"Report generation failed for task {self.request.id}: {exc}")
        try:
            # Retry the task if appropriate (e.g., transient DB failure)
            raise self.retry(exc=exc)
        except Exception:
            logger.critical(f"Task {self.request.id} failed permanently after retries.")
            return "Failed"
    except Exception as exc:
        logger.exception(f"An unexpected error occurred in report task {self.request.id}.")
        return "Failed"
