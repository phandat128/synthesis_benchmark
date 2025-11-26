import logging
import datetime
from io import BytesIO
from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas
from reportlab.lib.styles import getSampleStyleSheet
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle
from reportlab.lib import colors

from pyramid_report_app.models.data_model import DBSession, DataModel

log = logging.getLogger(__name__)

REPORT_TITLE = "Secure Data Report"

def fetch_data_for_report(limit: int, user_id: int):
    """
    Fetches data from the database based on validated parameters.
    Uses SQLAlchemy ORM securely (no raw SQL).
    The 'limit' parameter is guaranteed to be safe and validated by the caller.
    """
    log.info(f"Fetching {limit} records for user {user_id}.")

    # Secure query using parameterized limit and filtering by user_id (Least Privilege)
    try:
        data = DBSession.query(DataModel) \
            .filter(DataModel.user_id == user_id) \
            .order_by(DataModel.created_at.desc()) \
            .limit(limit) \
            .all()
        
        return [record.to_dict() for record in data]
    except Exception as e:
        log.error(f"Database error during data fetch: {e}")
        # Re-raise a generic error to prevent leaking DB specifics
        raise RuntimeError("Failed to retrieve report data.")


def create_pdf_report(data_records: list, requested_limit: int) -> BytesIO:
    """
    Generates a PDF report using ReportLab based on the provided data.
    The iteration over data_records is safe because the list size is controlled
    by the validated limit, preventing resource exhaustion (DoS).
    """
    
    buffer = BytesIO()
    
    # Use SimpleDocTemplate for structured flowables
    doc = SimpleDocTemplate(buffer, pagesize=letter)
    styles = getSampleStyleSheet()
    story = []

    # 1. Title and Metadata
    story.append(Paragraph(REPORT_TITLE, styles['Title']))
    story.append(Spacer(1, 12))
    story.append(Paragraph(f"Generated On: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}", styles['Normal']))
    story.append(Paragraph(f"Records Fetched: {len(data_records)} (Requested Limit: {requested_limit})", styles['Normal']))
    story.append(Spacer(1, 24))

    if not data_records:
        story.append(Paragraph("No data found matching the criteria.", styles['Italic']))
    else:
        # 2. Data Table Structure
        
        # Prepare table data
        table_data = [["ID", "User ID", "Title", "Created At"]]
        
        # Iterate only over the safely limited dataset.
        for record in data_records:
            # Output Encoding/Sanitization: ReportLab handles text rendering securely,
            # but we truncate long fields to prevent excessive page generation.
            table_data.append([
                str(record['id']),
                str(record['user_id']),
                record['title'][:50] + ('...' if len(record['title']) > 50 else ''), 
                record['created_at']
            ])

        # Create Table Flowable
        table = Table(table_data)
        
        # Apply Table Style
        table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.grey),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('BACKGROUND', (0, 1), (-1, -1), colors.beige),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))
        
        story.append(table)
        story.append(Spacer(1, 24))
        story.append(Paragraph("--- End of Report ---", styles['Italic']))

    try:
        doc.build(story)
    except Exception as e:
        log.error(f"Error building PDF document: {e}")
        raise RuntimeError("PDF generation failed.")

    buffer.seek(0)
    return buffer