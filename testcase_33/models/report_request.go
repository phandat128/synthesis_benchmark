package models

// ReportRequest defines the structure for the incoming JSON request.
type ReportRequest struct {
	// The 'record_count' field is the source of the potential DoS vulnerability.
	// Basic validation tags ensure it is present and positive.
	Count int `json:"record_count" validate:"required,min=1"` 
}

// Record represents a single row of data retrieved from the database.
type Record struct {
	ID      int
	Name    string
	Value   float64
	Details string
}