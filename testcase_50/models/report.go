package models

// Record represents a single data entry fetched from the database.
type Record struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Value     float64 `json:"value"`
	Timestamp string `json:"timestamp"`
}

// ReportRequest defines the structure for the incoming JSON request.
type ReportRequest struct {
	ReportType string `json:"report_type"`
	RecordCount int `json:"record_count"` // This field is validated in the handler layer
}

// ReportResponse defines the structure for the successful API response.
type ReportResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	ReportURL string `json:"report_url,omitempty"`
}