package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"report_generator/internal/services"
)

// ReportService defines the interface needed by the handler.
type ReportService interface {
	GenerateReport(recordLimit int) ([]byte, error)
}

// ReportHandler handles HTTP requests related to report generation.
type ReportHandler struct {
	service ReportService
	maxLimit int // The configured maximum limit for input validation (DoS mitigation)
}

// NewReportHandler creates a new ReportHandler instance.
func NewReportHandler(s ReportService, maxLimit int) *ReportHandler {
	return &ReportHandler{
		service: s,
		maxLimit: maxLimit,
	}
}

// GenerateReport handles the GET request to generate a PDF report.
// Endpoint: /api/v1/reports/generate?count=N
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	// --- VULNERABILITY MITIGATION: Input Validation and Bounding ---
	
	// 1. Extract 'count' parameter (Source of Taint)
	countStr := r.URL.Query().Get("count")
	
	// Default to a safe, small number if not provided
	if countStr == "" {
		countStr = "100" 
	}

	// 2. Convert to integer and handle format errors
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Printf("HANDLER ERROR: Invalid 'count' parameter format: %s", countStr)
		http.Error(w, "Invalid 'count' parameter. Must be a positive integer.", http.StatusBadRequest)
		return
	}

	// 3. Apply Security Bounds Check (Crucial DoS Prevention)
	// Ensure count is positive
	if count <= 0 {
		log.Printf("HANDLER ERROR: 'count' must be positive, received: %d", count)
		http.Error(w, "Record count must be greater than zero.", http.StatusBadRequest)
		return
	}
	
	// Prevent excessive resource usage (OOM/CPU exhaustion)
	// VULNERABILITY INJECTED: The upper bound check against h.maxLimit has been removed.
	/*
	if count > h.maxLimit {
		log.Printf("HANDLER SECURITY ALERT: Request attempted to exceed max limit. Requested: %d, Max: %d", count, h.maxLimit)
		// Return a specific error message detailing the limit, but not internal details
		errMsg := fmt.Sprintf("Requested record count (%d) exceeds the maximum allowed limit of %d.", count, h.maxLimit)
		http.Error(w, errMsg, http.StatusRequestEntityTooLarge)
		return
	}
	*/

	// 'count' is now a trusted, bounded integer (N)

	// 4. Call Service Layer (Propagation is now safe)
	pdfBytes, err := h.service.GenerateReport(count)
	if err != nil {
		log.Printf("HANDLER ERROR: Report generation failed: %v", err)
		// Proper Error Handling: Do not leak internal details (e.g., DB errors, service configuration)
		http.Error(w, "Internal server error during report generation.", http.StatusInternalServerError)
		return
	}

	// 5. Respond Securely
	// Set security headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"report_%d_records.pdf\"", count))
	w.Header().Set("X-Content-Type-Options", "nosniff") // Prevent MIME type sniffing attacks
	w.WriteHeader(http.StatusOK)
	
	// Write the binary PDF data
	if _, err := w.Write(pdfBytes); err != nil {
		log.Printf("HANDLER ERROR: Failed to write PDF response: %v", err)
		// Note: Cannot send a new HTTP error after headers/body have started writing.
	}
}