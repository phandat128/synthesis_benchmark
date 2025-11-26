package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection pool object
var DB *sqlx.DB

// Init initializes the database connection pool
func Init() {
	// Use environment variable for secure connection string handling
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable not set. Cannot connect to DB.")
	}

	var err error
	// sqlx.Connect is preferred over sql.Open for immediate connection testing
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		// Do not leak connection string details in production logs, but log the failure.
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}

	// Ping the database to ensure connectivity
	if err = DB.Ping(); err != nil {
		log.Fatalf("FATAL: Database ping failed: %v", err)
	}

	log.Println("Successfully connected to the database.")

	// Set connection pool limits (good practice for resource management)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// Initialize schema (for demonstration purposes)
	if err := createSchema(); err != nil {
		log.Fatalf("FATAL: Failed to initialize schema: %v", err)
	}
}

// createSchema sets up the necessary tables if they don't exist.
// In a production environment, dedicated migration tools (like goose or migrate) should be used.
func createSchema() error {
	// Enforcing length constraints directly in the schema as a secondary defense layer
	topicTable := `
	CREATE TABLE IF NOT EXISTS topics (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL CHECK (length(title) > 0 AND length(title) <= 100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	commentTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		content TEXT NOT NULL CHECK (length(content) > 0 AND length(content) <= 5000),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(topicTable); err != nil {
		return err
	}
	if _, err := DB.Exec(commentTable); err != nil {
		return err
	}
	log.Println("Database schema initialized successfully.")
	return nil
}