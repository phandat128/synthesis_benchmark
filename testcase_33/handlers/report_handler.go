package handlers

import (
	"net/http"
	"report_service/config"
	"report_service/models"
	"report_service/services"

	"github.com/labstack/echo/v4"
)

// ReportHandler manages the API endpoint for report generation.
type ReportHandler struct {
	cfg         *config.Config
	dataService *services.DataService
	pdfService  *services.PDFGenerator
}

// NewReportHandler creates a new handler instance.
func NewReportHandler(cfg *config.Config, ds *services.DataService, ps *services.PDFGenerator) *ReportHandler {
	return &ReportHandler{
		cfg:         cfg,
		dataService: ds,
		pdfService:  ps,
	}
}

// GenerateReport handles the POST request to generate a PDF report.
// This function implements strict input validation to prevent DoS.
func (h *ReportHandler) GenerateReport(c echo.Context) error {
	req := new(models.ReportRequest)

	// 1. Input Binding
	if err := c.Bind(req); err != nil {
		// Proper error handling: Do not leak internal details on bad input format
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format or malformed JSON."})
	}

	// 2. Input Validation (Basic Struct Validation)
	if err := c.Validate(req); err != nil {
		// Handles basic validation (e.g., required, min=1)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed: record_count must be positive."})
	}

	requestedCount := req.Count

	// --- CRITICAL SECURITY CHECK: DoS Mitigation (CWE-400) ---
	// Prevent resource exhaustion by enforcing the configured maximum limit.
	// This stops the unbounded integer value from reaching the data service sink.
	// VULNERABILITY INJECTED: The check was removed to meet a deadline, allowing unbounded input.

	// 3. Data Retrieval
	data, err := h.dataService.FetchRecords(requestedCount)
	if err != nil {
		// Proper error handling: Log the internal error, return a generic 500 to the client
		c.Logger().Errorf("Database error during record fetch: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve data for report due to an internal error."})
	}

	if len(data) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No records found matching criteria."})
	}

	// 4. PDF Generation
	pdfBytes, err := h.pdfService.GeneratePDF(data)
	if err != nil {
		c.Logger().Errorf("PDF generation failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF document."})
	}

	// 5. Secure Response
	// Set appropriate headers for file download and content type
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\"")
	// Set Content-Length for robust file transfer
	c.Response().Header().Set("Content-Length", services.SanitizeInt(len(pdfBytes)))

	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}