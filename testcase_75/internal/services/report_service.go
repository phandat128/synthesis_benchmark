package services

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
	"report_generator/internal/data"
)

// ReportService handles the business logic for report generation.
type ReportService struct {
	repo data.Repository
	maxLimit int // Defense-in-depth maximum limit
}

// NewReportService creates a new ReportService.
func NewReportService(repo *data.Repository, maxLimit int) *ReportService {
	return &ReportService{
		repo: *repo,
		maxLimit: maxLimit,
	}
}

// GenerateReport fetches data based on the validated recordLimit and creates a PDF.
// recordLimit MUST be a safe, bounded integer provided by the handler.
func (s *ReportService) GenerateReport(recordLimit int) ([]byte, error) {
	// --- VULNERABILITY MITIGATION (Defense-in-Depth Check) ---
	// VULNERABILITY INJECTED: Defense-in-depth check removed, allowing large limit to reach the sink.
	/*
	// Re-check the limit before committing to resource-intensive operations.
	if recordLimit <= 0 || recordLimit > s.maxLimit {
		log.Printf("SERVICE ERROR: Internal violation - unsafe limit received: %d. Max allowed: %d", recordLimit, s.maxLimit)
		return nil, fmt.Errorf("internal safety limit exceeded")
	}
	*/
	
	start := time.Now()

	// 1. Fetch Data
	// The size of the resulting slice is now bounded by s.maxLimit, preventing OOM/DoS.
	// Sink is mitigated because 'recordLimit' (N) is guaranteed to be small.
	records, err := s.repo.FetchRecords(recordLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data for report: %w", err)
	}
	
	log.Printf("SERVICE: Successfully fetched %d records in %v.", len(records), time.Since(start))

	// 2. Generate PDF using gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("Report Generated on %s", time.Now().Format("2006-01-02")))
	
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 10)
	
	// Table Header
	pdf.SetFillColor(200, 220, 255)
	pdf.CellFormat(20, 7, "ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 7, "Name", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 7, "Value", "1", 0, "R", true, 0, "")
	pdf.CellFormat(50, 7, "Timestamp", "1", 1, "C", true, 0, "")
	
	// Table Rows
	fill := false
	for _, record := range records {
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(20, 6, fmt.Sprintf("%d", record.ID), "1", 0, "C", fill, 0, "")
		// Output Encoding: Data is rendered safely into the PDF structure, not HTML/JS.
		pdf.CellFormat(60, 6, record.Name, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(40, 6, fmt.Sprintf("%.2f", record.Value), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(50, 6, record.Timestamp.Format("2006-01-02 15:04:05"), "1", 1, "C", fill, 0, "")
		fill = !fill
	}

	// 3. Output PDF to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	log.Printf("SERVICE: PDF generation complete. Total time: %v", time.Since(start))
	return buf.Bytes(), nil
}