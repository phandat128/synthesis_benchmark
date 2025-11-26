import logging
from io import BytesIO
from django.template.loader import render_to_string
from weasyprint import HTML, CSS
from reports.models import TransactionRecord

logger = logging.getLogger(__name__)

# Minimal CSS for PDF generation
DEFAULT_CSS = CSS(string='''
    @page { size: A4; margin: 1cm; }
    body { font-family: sans-serif; }
    table { width: 100%; border-collapse: collapse; }
    th, td { border: 1px solid #ccc; padding: 8px; }
''')

def generate_pdf_report(limit: int) -> BytesIO:
    """
    Fetches data and generates a PDF report using WeasyPrint.
    
    The 'limit' parameter is guaranteed to be safe due to prior validation 
    in the serializer, effectively mitigating the resource exhaustion DoS risk.
    
    Args:
        limit: The maximum number of records to fetch (guaranteed <= MAX_REPORT_RECORDS).

    Returns:
        A BytesIO object containing the generated PDF data.
    """
    
    logger.info(f"Starting PDF generation for a maximum of {limit} records.")
    
    try:
        # SINK: The limit is used here, but it is a validated, safe integer.
        # Using slicing [:limit] ensures the database query is bounded.
        records = TransactionRecord.objects.all().order_by('-timestamp')[:limit]
        
        context = {
            'records': records,
            'record_count': len(records),
            'limit_requested': limit,
        }

        # Render HTML content from a template.
        # Note: Django's template engine automatically escapes variables (Output Encoding)
        # preventing XSS if this content were rendered to a browser.
        html_content = render_to_string('report_template.html', context)

        # Generate PDF using WeasyPrint
        pdf_file = HTML(string=html_content).write_pdf(
            stylesheets=[DEFAULT_CSS]
        )
        
        pdf_buffer = BytesIO(pdf_file)
        pdf_buffer.seek(0)
        
        logger.info(f"Successfully generated PDF report with {len(records)} records.")
        return pdf_buffer

    except Exception as e:
        # Raise a specific runtime error to be caught and handled generically by the view.
        logger.error(f"Critical error during PDF generation: {e}", exc_info=True)
        raise RuntimeError("Report generation failed due to internal processing error.")
