package handlers

import (
	"log"
	"net/http"
	"report_service/config"
	"report_service/models"
	"report_service/services"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// ReportHandler holds dependencies for report generation.
type ReportHandler struct {
	Generator *services.ReportGenerator
}

// NewReportHandler creates a new handler instance.
func NewReportHandler(g *services.ReportGenerator) *ReportHandler {
	return &ReportHandler{Generator: g}
}

// GenerateReport handles the POST request to create a PDF report.
func (h *ReportHandler) GenerateReport(c echo.Context) error {
	req := new(models.ReportRequest)
	
	// 1. Input Parsing
	if err := c.Bind(req); err != nil {
		log.Printf("Binding error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format or malformed JSON"})
	}

	requestedCount := req.RecordCount
	maxLimit := config.AppConfig.MaxReportRecords

	// 2. SECURE INPUT VALIDATION (Crucial Defense against Resource Exhaustion DoS)
	
	// Validate that the requested count is a positive integer.
	if requestedCount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "record_count must be a positive integer"})
	}

	// Validate against configured maximum limit.
	// This prevents an attacker from requesting an excessive number of records
	// which would lead to unbounded memory allocation and CPU exhaustion.
	// VULNERABILITY INJECTED: The check and capping logic is removed, allowing unbounded input.
	
	// 3. Business Logic Execution
	
	// Use the request context for cancellation/timeouts
	ctx := c.Request().Context()
	
	pdfBytes, err := h.Generator.GeneratePDFReport(ctx, requestedCount)
	
	if err != nil {
		// Proper error handling: Do not leak sensitive internal details.
		log.Printf("Report generation failed: %v", err)
		if err.Error() == "no records found to generate report" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "No data available for the requested report parameters"})
		}
		// Generic internal error response
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate report due to an internal server error"})
	}

	// 4. Secure Output Handling
	
	// Set appropriate headers for file download and security
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\"report_"+strconv.Itoa(requestedCount)+"_"+strconv.FormatInt(time.Now().Unix(), 10)+".pdf\"")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff") // Prevent MIME type sniffing attacks

	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}