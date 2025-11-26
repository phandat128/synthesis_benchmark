package utils

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/secure-api/models"
)

// PDFGenerator handles the creation of PDF documents.
type PDFGenerator struct{}

// NewPDFGenerator creates a new instance of PDFGenerator.
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// CreatePDF generates a PDF document from a slice of records.
// This function performs CPU and memory intensive operations (string manipulation, layout calculation,
// memory allocation for the PDF structure). Its performance is directly tied to the size of the input slice.
// Security Note: The input slice size must be strictly controlled upstream (in the handler) to prevent DoS (CWE-400).
func (g *PDFGenerator) CreatePDF(records []models.Record) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set up basic font and title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Customizable Data Report")

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(12)
	pdf.Cell(40, 10, fmt.Sprintf("Records Generated: %d | On: %s", len(records), time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(10)

	// Table Header Configuration
	header := []string{"ID", "Name", "Value", "Timestamp"}
	colWidths := []float64{15, 60, 30, 70}
	
	pdf.SetFillColor(200, 220, 255)
	pdf.SetFont("Arial", "B", 10)

	// Draw Header
	for i, h := range header {
		pdf.CellFormat(colWidths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1) // Move to the next line for data

	// Table Data
	pdf.SetFont("Arial", "", 8)
	fill := false

	for _, record := range records {
		// Set row background color for readability
		pdf.SetFillColor(240, 240, 240)
		if fill {
			pdf.SetFillColor(255, 255, 255)
		}

		// Add data cells. Data is generated internally, so XSS is not a concern here,
		// but if data came from user input, output encoding/sanitization would be mandatory.
		
		// ID
		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", record.ID), "1", 0, "C", fill, 0, "")
		// Name
		pdf.CellFormat(colWidths[1], 6, record.Name, "1", 0, "L", fill, 0, "")
		// Value
		pdf.CellFormat(colWidths[2], 6, fmt.Sprintf("%.2f", record.Value), "1", 0, "R", fill, 0, "")
		// Timestamp
		pdf.CellFormat(colWidths[3], 6, record.Timestamp, "1", 0, "L", fill, 0, "")

		pdf.Ln(-1)
		fill = !fill
	}

	// Output the PDF to a byte buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Error generating PDF output: %v", err)
		return nil, fmt.Errorf("failed to finalize PDF generation")
	}

	return buf.Bytes(), nil
}
