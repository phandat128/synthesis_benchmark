package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/secure-api/models"
)

// DataService handles the retrieval and simulation of data records.
type DataService struct{ }

// NewDataService creates a new instance of DataService.
func NewDataService() *DataService {
	return &DataService{}
}

// FetchRecords simulates fetching a specified number of records from a database.
// It relies on the calling handler to have already validated the 'limit' parameter
// to prevent resource exhaustion (DoS).
func (s *DataService) FetchRecords(limit int) []models.Record {
	if limit <= 0 {
		return []models.Record{}
	}

	// Initialize random seed securely (though time-based is fine for simulation)
	rand.Seed(time.Now().UnixNano())
	records := make([]models.Record, limit)

	for i := 0; i < limit; i++ {
		records[i] = models.Record{
			ID:          i + 1,
			Name:        fmt.Sprintf("Product_%d", i+1),
			Description: fmt.Sprintf("Detailed description for item %d. This field is long.", i+1),
			Value:       rand.Float64() * 1000,
			Timestamp:   time.Now().Add(time.Duration(-i) * time.Hour).Format(time.RFC3339),
		}
	}

	return records
}
