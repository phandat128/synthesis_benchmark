package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
	"report_service/models"
)

// GeneratePDF creates a PDF report from the provided data.
// The resource consumption here is inherently limited by the size of the input 'data' slice,
// which must be capped upstream by input validation (CWE-400 mitigation).
func GeneratePDF(data []models.Record, reportType string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set up basic font and title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Report: %s", reportType))
	pdf.Ln(12)

	// Metadata
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(6)
	pdf.Cell(40, 6, fmt.Sprintf("Total Records Included: %d", len(data)))
	pdf.Ln(10)

	// Table Headers
	pdf.SetFillColor(200, 220, 255)
	pdf.SetFont("Arial", "B", 10)
	header := []string{"ID", "Name", "Value", "Timestamp"}
	w := []float64{20, 50, 40, 60} // Column widths

	for i, h := range header {
		pdf.CellFormat(w[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Data Rows
	pdf.SetFont("Arial", "", 10)
	fill := false
	
	// The loop iterates safely over the actual data slice.
	for _, record := range data {
		pdf.CellFormat(w[0], 6, fmt.Sprintf("%d", record.ID), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(w[1], 6, record.Name, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(w[2], 6, fmt.Sprintf("%.2f", record.Value), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(w[3], 6, record.Timestamp, "1", 0, "L", fill, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}

	// Save the file
	filename := fmt.Sprintf("report_%d.pdf", time.Now().UnixNano())
	filepath := fmt.Sprintf("./reports/%s", filename)

	// Ensure the reports directory exists (secure practice)
	if err := os.MkdirAll("./reports", 0755); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	err := pdf.OutputFileAndClose(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF file: %w", err)
	}

	return filename, nil
}