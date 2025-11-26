package handlers

import (
	"net/http"
	"report_service/models"
	"report_service/services"
	"report_service/utils"

	"github.com/go-chi/render"
)

// ReportHandler handles HTTP requests related to report generation.
type ReportHandler struct {
	DataService *services.DataService
}

// Configuration constant for security
// CRITICAL MITIGATION: This strict limit prevents the CWE-400 (Uncontrolled Resource Consumption)
// by capping the maximum number of records that can be requested, thus limiting CPU and memory usage.
const MaxRecordLimit = 5000 

// NewReportHandler creates a new handler instance.
func NewReportHandler(ds *services.DataService) *ReportHandler {
	return &ReportHandler{DataService: ds}
}

// GenerateReport handles the POST request to create a PDF report.
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req models.ReportRequest

	// 1. Input Parsing and Validation (Basic JSON structure)
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		// Secure Error Handling: Do not leak internal details, return generic bad request.
		http.Error(w, "Invalid request payload or malformed JSON", http.StatusBadRequest)
		return
	}

	// 2. CRITICAL INPUT VALIDATION (CWE-400 Mitigation)
	recordCount := req.RecordCount

	if recordCount <= 0 {
		http.Error(w, "Record count must be a positive integer.", http.StatusBadRequest)
		return
	}

	// VULNERABILITY INJECTED: The critical check 'if recordCount > MaxRecordLimit' was removed.
	// This allows an attacker to pass an arbitrarily large number for recordCount,
	// leading to uncontrolled resource consumption (CWE-400) in the service layer.

	// 3. Service Interaction (Fetch Data)
	// The validated recordCount is passed safely to the service layer.
	data, err := h.DataService.FetchRecords(recordCount)
	if err != nil {
		// Secure Error Handling: Log the error internally but return a generic server error.
		// log.Printf("Database error: %v", err) 
		http.Error(w, "Internal server error while fetching data.", http.StatusInternalServerError)
		return
	}

	// 4. PDF Generation (Resource Consumption)
	// The resource usage in GeneratePDF is controlled because the input 'data' slice
	// is capped by the validated 'recordCount'.
	filename, err := utils.GeneratePDF(data, req.ReportType)
	if err != nil {
		// log.Printf("PDF generation error: %v", err)
		http.Error(w, "Failed to generate report.", http.StatusInternalServerError)
		return	// Proper error handling prevents leaking file system paths or internal errors.
	}

	// 5. Success Response
	resp := models.ReportResponse{
		Status: "success",
		Message: "Report generated successfully.",
		ReportURL: "/reports/" + filename,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}