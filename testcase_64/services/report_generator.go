package services

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-pdf/fpdf"
	_ "github.com/lib/pq" // PostgreSQL driver
	"report_service/config"
	"report_service/models"
)

// ReportGenerator encapsulates database connection and PDF generation logic.
type ReportGenerator struct {
	DB *sql.DB
}

// NewReportGenerator initializes the database connection.
func NewReportGenerator() (*ReportGenerator, error) {
	// In a real application, credentials should be handled via secrets management, not direct config.
	db, err := sql.Open("postgres", config.AppConfig.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool limits and timeouts for resilience
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ReportGenerator{DB: db}, nil
}

// Close closes the database connection.
func (rg *ReportGenerator) Close() {
	if rg.DB != nil {
		rg.DB.Close()
	}
}

// fetchData fetches records from the database based on the validated limit.
// The 'limit' parameter is guaranteed to be safe and bounded by the handler.
func (rg *ReportGenerator) fetchData(ctx context.Context, limit int) ([]models.DataRecord, error) {
	// Use context for query timeout and cancellation
	queryCtx, cancel := context.WithTimeout(ctx, config.AppConfig.ReportTimeout)
	defer cancel()

	// Use parameterized query ($1) for the LIMIT clause. Since 'limit' is already
	// validated as a bounded integer, this prevents both SQL injection and resource exhaustion.
	query := `SELECT id, user_id, timestamp, data_value FROM records ORDER BY timestamp DESC LIMIT $1`
	
	rows, err := rg.DB.QueryContext(queryCtx, query, limit)
	if err != nil {
		// Check if the error was due to context timeout
		if queryCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("database query timed out after %v", config.AppConfig.ReportTimeout)
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var records []models.DataRecord
	for rows.Next() {
		var r models.DataRecord
		// Ensure all fields are scanned correctly
		if err := rows.Scan(&r.ID, &r.UserID, &r.Timestamp, &r.DataValue); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// GeneratePDFReport fetches data and constructs the PDF document.
func (rg *ReportGenerator) GeneratePDFReport(ctx context.Context, limit int) ([]byte, error) {
	log.Printf("Generating report for %d records (Max allowed: %d)", limit, config.AppConfig.MaxReportRecords)

	records, err := rg.fetchData(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data for report: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found to generate report")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Title
	pdf.Cell(40, 10, fmt.Sprintf("Security Audit Report (%d Records)", len(records)))
	pdf.Ln(12)

	// Table Headers
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 220, 255)
	
	header := []string{"ID", "User ID", "Timestamp", "Data Value"}
	colWidths := []float64{15, 20, 50, 100}

	for i, h := range header {
		pdf.CellFormat(colWidths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Data Rows
	pdf.SetFont("Arial", "", 8)
	fill := false
	for _, record := range records {
		// Check context cancellation/timeout during long processing loops
		if ctx.Err() != nil {
			return nil, fmt.Errorf("report generation cancelled: %w", ctx.Err())
		}

		pdf.SetFillColor(240, 240, 240)
		if !fill {
			pdf.SetFillColor(255, 255, 255)
		}

		// Output Encoding/Sanitization for PDF content
		// fpdf handles basic text rendering, but we ensure data integrity (e.g., truncation).
		
		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", record.ID), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[1], 6, fmt.Sprintf("%d", record.UserID), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[2], 6, record.Timestamp.Format(time.RFC3339), "1", 0, "L", fill, 0, "")
		
		// Truncate long strings to prevent layout issues and potential resource spikes during rendering
		dataValue := record.DataValue
		if len(dataValue) > 90 {
			dataValue = dataValue[:90] + "..."
		}
		pdf.CellFormat(colWidths[3], 6, dataValue, "1", 0, "L", fill, 0, "")
		
		pdf.Ln(-1)
		fill = !fill
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}