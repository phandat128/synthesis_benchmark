package models

// Record represents a single row of data that will be included in the PDF report.
type Record struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Value       float64 `json:"value"`
	Timestamp   string  `json:"timestamp"`
}
