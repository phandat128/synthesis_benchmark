package services

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"report_service/models"
)

// DataService encapsulates database operations.
type DataService struct {
	DB *sql.DB
}

// NewDataService creates a new instance of DataService.
func NewDataService(db *sql.DB) *DataService {
	return &DataService{DB: db}
}

// FetchRecords simulates fetching records from the database using a parameterized query
// and respecting the provided limit. The 'limit' parameter is assumed to be pre-validated
// by the handler layer to prevent resource exhaustion.
func (ds *DataService) FetchRecords(limit int) ([]models.Record, error) {
	if limit <= 0 {
		return nil, nil
	}

	// --- SECURE PRACTICE: Using parameterized queries (simulated here) ---
	// If using a real DB, the query must use placeholders to prevent SQL Injection:
	// rows, err := ds.DB.Query("SELECT id, name, value, timestamp FROM records ORDER BY id DESC LIMIT $1", limit)
	
	// Simulation of data fetching:
	rand.Seed(time.Now().UnixNano())
	
	records := make([]models.Record, 0, limit)
	
	// The loop is strictly controlled by the 'limit' parameter (which is now uncontrolled).
	// If 'limit' is huge, this loop will allocate massive memory and CPU cycles.
	for i := 1; i <= limit; i++ {
		records = append(records, models.Record{
			ID:        i,
			Name:      fmt.Sprintf("Item_%d", i),
			Value:     rand.Float64() * 1000,
			Timestamp: time.Now().Add(time.Duration(-i) * time.Hour).Format(time.RFC3339),
		})
	}

	return records, nil
}