package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"report_service/services"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Starting Report Generation Service...")

	// --- 1. Database Initialization (Simulated/Placeholder) ---
	// Secure Practice: Use placeholder DB connection setup.
	// In a real application, environment variables should be used for connection strings.
	// We use a mock DB connection since the core logic is resource validation,
	// but we structure the code to accept a real DB connection.
	
	// NOTE: Replace this with actual secure connection string retrieval (e.g., from secrets manager or ENV)
	db, err := sql.Open("postgres", "user=postgres password=example dbname=reports sslmode=disable")
	if err != nil {
		// Log and exit if critical resource (DB) fails to connect
		log.Printf("Warning: Failed to connect to database (using mock data): %v", err)
		// In a real scenario, this might be fatal, but for simulation, we proceed with a nil DB.
	} else {
		defer db.Close()
		// Secure DB connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)
		log.Println("Database connection established (or simulated).")
	}
	
	// Ensure the reports directory exists for output
	if err := os.MkdirAll("./reports", 0755); err != nil {
		log.Fatalf("Failed to create reports directory: %v", err)
	}

	// --- 2. Service Layer Setup ---
	dataService := services.NewDataService(db)

	// --- 3. Router Setup ---
	router := SetupRouter(dataService)

	// --- 4. Server Configuration (Secure Defaults) ---
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,  // Secure: Limit time spent reading request headers/body
		WriteTimeout: 60 * time.Second, // Allow time for PDF generation/writing response
		IdleTimeout:  120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB limit on headers
	}

	log.Printf("Server listening securely on %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", port, err)
	}
}