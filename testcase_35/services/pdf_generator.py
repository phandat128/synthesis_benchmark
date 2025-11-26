from reportlab.lib.pagesizes import letter
from reportlab.platypus import SimpleDocTemplate, Table, TableStyle, Paragraph
from reportlab.lib.styles import getSampleStyleSheet
from reportlab.lib import colors
import io

def generate_customer_report_pdf(data: list, record_count: int) -> bytes:
    """
    Generates a PDF report from the provided customer data using ReportLab.
    
    This function is CPU and memory intensive. It relies on the caller 
    (views/report_views.py) ensuring that 'data' is not excessively large.
    
    :param data: List of customer dictionaries.
    :param record_count: The number of records requested.
    :return: Bytes content of the generated PDF.
    """
    
    # Use io.BytesIO to create the PDF in memory without touching the filesystem
    buffer = io.BytesIO()
    
    # Setup document template
    doc = SimpleDocTemplate(buffer, pagesize=letter, 
                            title="Customer Report",
                            leftMargin=72, rightMargin=72,
                            topMargin=72, bottomMargin=72)
    
    styles = getSampleStyleSheet()
    elements = []

    # 1. Title and Metadata
    elements.append(Paragraph("Secure Customer Data Report", styles['Title']))
    elements.append(Paragraph(f"Report generated on {record_count} records.", styles['Normal']))
    elements.append(Paragraph("<br/>", styles['Normal'])) # Spacer

    if not data:
        elements.append(Paragraph("No data found for the report criteria.", styles['Italic']))
    else:
        # 2. Prepare Table Data
        headers = ["ID", "First Name", "Last Name", "Email", "Balance", "Created At"]
        table_data = [headers]
        
        for item in data:
            row = [
                str(item.get('id', '')),
                item.get('first_name', ''),
                item.get('last_name', ''),
                item.get('email', 'N/A'),
                item.get('balance', '$0.00'),
                item.get('created_at', '')[:10] 
            ]
            table_data.append(row)

        # 3. Create Table and Style
        table = Table(table_data)
        
        table.setStyle(TableStyle([
            ('BACKGROUND', (0, 0), (-1, 0), colors.darkslategray),
            ('TEXTCOLOR', (0, 0), (-1, 0), colors.whitesmoke),
            ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
            ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
            ('BOTTOMPADDING', (0, 0), (-1, 0), 12),
            ('GRID', (0, 0), (-1, -1), 1, colors.black)
        ]))
        
        elements.append(table)

    # Build the PDF
    try:
        doc.build(elements)
    except Exception as e:
        # Robust error handling for rendering failures
        print(f"PDF generation failed: {e}")
        raise IOError("Failed to finalize PDF document.")

    # Return content
    buffer.seek(0)
    return buffer.read()