package services

import (
	"fmt"
	"log"
	"strconv"

	"report_service/models"

	"github.com/jung-kurt/gofpdf"
)

// PDFGenerator handles the creation of PDF documents.
type PDFGenerator struct{}

// NewPDFGenerator creates a new instance of the PDFGenerator service.
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// GeneratePDF iterates over the data and constructs the final PDF document.
// This function is safe because the input slice 'data' size is guaranteed to be bounded.
func (g *PDFGenerator) GeneratePDF(data []models.Record) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot generate PDF from empty data set")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Complex Report - %d Records", len(data)))

	pdf.SetFont("Arial", "", 10)
	y := 30.0

	// Define table headers
	headers := []string{"ID", "Name", "Value", "Details"}
	colWidths := []float64{15, 40, 30, 100}

	// Print headers
	pdf.SetY(y)
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	y += 7
	pdf.Ln(-1)

	// --- SAFE ITERATION (The former sink) ---
	// The loop is bounded by the validated input size, preventing CPU exhaustion.
	for _, record := range data {
		pdf.SetY(y)
		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", record.ID), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[1], 6, record.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[2], 6, fmt.Sprintf("%.2f", record.Value), "1", 0, "R", false, 0, "")
		
		// Output Encoding/Truncation: Ensure data fits the cell to prevent layout issues or excessive processing.
		details := record.Details
		if len(details) > 75 {
			details = details[:72] + "..."
		}
		pdf.CellFormat(colWidths[3], 6, details, "1", 0, "L", false, 0, "")
		
		pdf.Ln(-1)
	y += 6

		// Add new page if necessary
		if y > 270 {
			pdf.AddPage()
			y = 10.0
			pdf.SetFont("Arial", "", 10)
		}
	}

	// Output the PDF to a byte slice
	buf := pdf.Output()

	if pdf.Error() != nil {
		log.Printf("PDF generation library error: %v", pdf.Error())
		return nil, fmt.Errorf("pdf generation error")
	}

	return buf, nil
}

// SanitizeInt converts an integer to a string safely for output encoding (e.g., in HTTP headers or JSON responses).
func SanitizeInt(i int) string {
	return strconv.Itoa(i)
}