package data

import (
	"fmt"
	"log"
	"time"
)

// Record represents a single data entry fetched from the database.
type Record struct {
	ID        int
	Name      string
	Value     float64
	Timestamp time.Time
}

// DBConnection interface abstracts the database interaction.
type DBConnection interface {
	// Query simulates executing a parameterized query.
	Query(query string, args ...interface{}) ([]Record, error)
}

// MockDBConnection simulates a database connection for demonstration.
type MockDBConnection struct{}

// Query simulates executing a parameterized query safely.
// It expects the first argument to be the integer limit.
func (m *MockDBConnection) Query(query string, args ...interface{}) ([]Record, error) {
	// In a real application, this would use db.QueryRowContext and the driver handles parameter binding.

	if len(args) == 0 {
		return nil, fmt.Errorf("missing query arguments")
	}

	limit, ok := args[0].(int)
	if !ok || limit <= 0 {
		return nil, fmt.Errorf("invalid limit parameter provided to repository")
	}

	// Since 'limit' is an integer derived from a validated source, using it in the LIMIT clause is safe
	// against SQL injection, and safe against DoS due to upstream validation.
	log.Printf("DATA: Executing query with safe limit: %d", limit)

	// Simulate fetching 'limit' number of records
	records := make([]Record, limit)
	for i := 0; i < limit; i++ {
		records[i] = Record{
			ID:        i + 1,
			Name:      fmt.Sprintf("Item %d", i+1),
			Value:     float64(i) * 1.5,
			Timestamp: time.Now().Add(time.Duration(-i) * time.Hour),
		}
	}
	return records, nil
}

// Repository handles data retrieval operations.
type Repository struct {
	db DBConnection
}

// NewRepository creates a new Repository instance.
func NewRepository(db DBConnection) *Repository {
	return &Repository{db: db}
}

// FetchRecords retrieves a specified number of records from the data store.
// The 'limit' parameter is assumed to be pre-validated and bounded.
func (r *Repository) FetchRecords(limit int) ([]Record, error) {
	// Example SQL structure using a placeholder for the limit parameter:
	// "SELECT id, name, value, timestamp FROM records ORDER BY id DESC LIMIT ?"
	
	// We pass the trusted integer limit as an argument.
	records, err := r.db.Query("SELECT * FROM records LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("repository failed to fetch records: %w", err)
	}

	return records, nil
}