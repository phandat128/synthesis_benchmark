import logging
import datetime
from django.template.loader import render_to_string
from weasyprint import HTML
from weasyprint.logger import LOGGER
from reports.models import ReportData
from typing import Dict, Any, List

logger = logging.getLogger(__name__)

# Suppress WeasyPrint logging spam unless critical
LOGGER.setLevel(logging.WARNING)

def generate_pdf_report(record_count: int, report_title: str, active_only: bool) -> bytes:
    """
    Core business logic to query data and generate a PDF report.
    
    record_count is guaranteed to be safe (bounded) by the serializer validation.
    """
    
    # 1. Secure Data Retrieval (Using Django ORM to prevent SQL Injection)
    try:
        # Filter based on validated input
        queryset = ReportData.objects.filter(is_active=active_only).order_by('-timestamp')
        
        # Apply the validated, bounded limit. This prevents DoS/OOM.
        data_records: List[ReportData] = list(queryset[:record_count])
        
        if not data_records:
            logger.warning("Report generation requested but no records found.")
            return generate_empty_report(report_title)

    except Exception as e:
        logger.error(f"Database query failed during report generation: {e}")
        # Raise a runtime error that the Celery task can catch and retry/fail gracefully
        raise RuntimeError("Failed to retrieve report data from the database.") from e

    # 2. Prepare Context and Render HTML
    context: Dict[str, Any] = {
        'title': report_title,
        'records': data_records,
        'record_count': len(data_records),
        'generation_time': datetime.datetime.now(),
    }
    
    # Django templates handle auto-escaping, preventing XSS in the output HTML.
    try:
        html_content = render_to_string('reports/report_template.html', context)
    except Exception as e:
        logger.error(f"Failed to render HTML template: {e}")
        raise RuntimeError("Failed to assemble report content.") from e

    # 3. Generate PDF using WeasyPrint
    try:
        pdf_bytes = HTML(string=html_content).write_pdf()
        logger.info(f"Successfully generated PDF report titled '{report_title}' with {len(data_records)} records.")
        return pdf_bytes
    except Exception as e:
        logger.error(f"WeasyPrint PDF generation failed: {e}")
        raise RuntimeError("Failed to convert HTML content to PDF.") from e

def generate_empty_report(title: str) -> bytes:
    """Generates a placeholder PDF when no data is found."""
    # Assuming an 'reports/empty_report_template.html' exists
    html_content = render_to_string('reports/empty_report_template.html', {'title': title})
    return HTML(string=html_content).write_pdf()
