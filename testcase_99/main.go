package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"forum_api/db"
	"forum_api/routes"
)

func main() {
	// 1. Initialize Database Connection
	// Note: DATABASE_URL environment variable must be set (e.g., postgres://user:pass@host/dbname?sslmode=disable)
	db.Init()
	defer func() {
		if db.DB != nil {
			db.DB.Close()
			log.Println("Database connection closed.")
		}
	}()

	// 2. Setup Router
	router := routes.SetupRouter()

	// 3. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Starting server securely on %s", addr)

	// Use http.Server for better control over timeouts (security best practice)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		// Secure timeout settings to prevent slowloris attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", addr, err)
	}
}