package models

import "time"

// DataRecord represents a single row fetched from the simulated database.
type DataRecord struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	DataValue string    `json:"data_value"`
}

// ReportRequest defines the expected structure for the incoming request body.
type ReportRequest struct {
	// record_count is the user-supplied limit for the report size.
	// This field is the source of the vulnerability taint flow and must be strictly validated.
	RecordCount int `json:"record_count"`
}