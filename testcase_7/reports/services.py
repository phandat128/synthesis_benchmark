import io
import logging
from datetime import date
from reportlab.lib.pagesizes import letter
from reportlab.lib.styles import getSampleStyleSheet
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
from reportlab.lib import colors
from django.core.exceptions import ObjectDoesNotExist

from .models import ReportData

logger = logging.getLogger(__name__)

class ReportGenerationError(Exception):
    """Custom exception for service layer errors."""
    pass

def _fetch_report_data(limit: int, start_date: date = None, end_date: date = None):
    """
    Securely fetches data from the database based on validated parameters.
    
    The 'limit' parameter is guaranteed to be safe due to prior validation 
    in the serializer layer (CWE-400 mitigation).
    """
    try:
        # Start with a base queryset
        queryset = ReportData.objects.all()
        
        # Apply filtering based on validated date inputs
        if start_date:
            queryset = queryset.filter(created_at__date__gte=start_date)
        if end_date:
            queryset = queryset.filter(created_at__date__lte=end_date)
            
        # Apply the safe, bounded limit using Django ORM slicing
        # This prevents excessive memory usage.
        data = list(queryset[:limit]) 
        
        if not data:
            raise ObjectDoesNotExist("No data found matching the criteria.")
            
        return data
        
    except Exception as e:
        logger.error(f"Database retrieval error: {e}")
        # Raise a generic error to prevent leaking database specifics (CWE-209)
        raise ReportGenerationError("Failed to retrieve report data.")


def generate_pdf_report(report_type: str, limit: int, start_date: date = None, end_date: date = None) -> io.BytesIO:
    """
    Generates a PDF report based on the provided, validated parameters.
    """
    
    logger.info(f"Starting PDF generation for type: {report_type}, limit: {limit}")
    
    try:
        # 1. Data Retrieval (Safe due to validated limit)
        report_records = _fetch_report_data(limit, start_date, end_date)
        
        # 2. Setup PDF buffer
        buffer = io.BytesIO()
        # Use SimpleDocTemplate for robust document structure
        doc = SimpleDocTemplate(buffer, pagesize=letter)
        styles = getSampleStyleSheet()
        
        elements = []
        
        # Title and Metadata
        elements.append(Paragraph(f"Report: {report_type.replace('_', ' ').title()}", styles['Title']))
        elements.append(Paragraph(f"Generated on: {date.today()}", styles['Normal']))
        elements.append(Paragraph(f"Record Limit Applied: {limit}", styles['Normal']))
        elements.append(Spacer(1, 12))

        # 3. Prepare Table Data
        table_data = [['UUID', 'Customer Name', 'Amount', 'Date']]
        
        for record in report_records:
            # Ensure all data is properly converted to strings for PDF rendering (Output Encoding for PDF)
            row = [
                str(record.uuid)[:8] + '...', 
                str(record.customer_name),
                f"${record.transaction_amount:,.2f}",
                record.created_at.strftime('%Y-%m-%d %H:%M')
            ]
            table_data.append(row)
            
        # 4. Create and Style Table
        table = Table(table_data)
        table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))
        
        elements.append(table)
        
        # 5. Build the PDF
        doc.build(elements)
        
        # Reset buffer position to the beginning before returning
        buffer.seek(0)
        return buffer

    except ObjectDoesNotExist as e:
        logger.warning(f"Report generation failed: {e}")
        raise ReportGenerationError(str(e))
    except Exception as e:
        logger.error(f"Critical error during PDF generation: {e}", exc_info=True)
        raise ReportGenerationError("An unexpected error occurred during report processing.")