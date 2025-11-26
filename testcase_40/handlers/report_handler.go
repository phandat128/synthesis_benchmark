package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secure-api/services"
	"github.com/secure-api/utils"
)

// MaxReportRecords defines the safe upper limit for records requested by a user.
// This constant is the primary defense against the Denial of Service (DoS) vulnerability (CWE-400).
const MaxReportRecords = 10000 

// ReportHandler holds dependencies for report generation.
type ReportHandler struct {
	DataService *services.DataService
	PDFGen      *utils.PDFGenerator
}

// NewReportHandler creates a new handler instance.
func NewReportHandler(ds *services.DataService, pg *utils.PDFGenerator) *ReportHandler {
	return &ReportHandler{
		DataService: ds,
		PDFGen:      pg,
	}
}

// GenerateReport handles the request to create a PDF report.
// Endpoint: GET /api/v1/reports/generate?record_count=N
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	// 1. Input Retrieval
	countStr := c.Query("record_count")
	if countStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameter: record_count"})
		return
	}

	// 2. Input Validation (Type Check)
	recordCount, err := strconv.Atoi(countStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record_count format. Must be a positive integer."})
		return
	}

	// --- SECURITY CRITICAL STEP: Resource Consumption Control (CWE-400 Defense) ---
	// Validate bounds to prevent an attacker from requesting an arbitrarily large report,
	// which would exhaust memory and CPU resources during PDF generation.
	// VULNERABILITY INJECTED: The bounds check (recordCount > MaxReportRecords) was removed.
	// -----------------------------------------------------------------------------

	// 3. Data Retrieval (Simulated)
	// The service layer receives an already validated limit.
	records := h.DataService.FetchRecords(recordCount)

	if len(records) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No data found for the specified criteria."})
		return
	}

	// 4. PDF Generation (CPU/Memory intensive operation)
	pdfBytes, err := h.PDFGen.CreatePDF(records)
	if err != nil {
		// Proper error handling: Do not leak internal details. Log the detailed error internally.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report due to internal processing error."})
		return
	}

	// 5. Response
	// Set appropriate headers for file download.
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	// Set Content-Type header.
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
