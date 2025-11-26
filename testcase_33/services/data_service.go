package services

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"report_service/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DataService handles database operations.
type DataService struct {
	db *sql.DB
}

// NewDataService initializes the database connection.
func NewDataService(dbURL string) *DataService {
	// In a real application, use connection pooling and proper error handling.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Set connection limits to prevent resource exhaustion from too many concurrent connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	
	// Ping the database to ensure connectivity
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Successfully connected to the database.")
	return &DataService{db: db}
}

// FetchRecords retrieves a specified number of records from the database.
// The 'limit' parameter is guaranteed to be safe (bounded) by the calling handler.
func (s *DataService) FetchRecords(limit int) ([]models.Record, error) {
	// --- SECURE IMPLEMENTATION ---
	// 1. Use parameterized queries (or ORM) to prevent SQL Injection.
	// Example parameterized query (PostgreSQL syntax):
	query := `SELECT id, name, value, details FROM large_data_table LIMIT $1`

	// Simulate database query execution using the safe limit.
	log.Printf("Executing query to fetch %d records: %s", limit, query)

	rand.Seed(time.Now().UnixNano())
	records := make([]models.Record, limit)

	// Simulate fetching and scanning rows
	for i := 0; i < limit; i++ {
		records[i] = models.Record{
			ID:      i + 1,
			Name:    fmt.Sprintf("Item_%d", i+1),
			Value:   rand.Float64() * 1000,
			Details: fmt.Sprintf("Generated detail for record %d. This is a long detail string to test PDF wrapping.", i+1),
		}
	}

	// Simulate a successful fetch
	return records, nil
}