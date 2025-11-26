package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/your_project/database"
	"github.com/your_project/handlers"
	"github.com/your_project/router"
)

func main() {
	// 1. Load Environment Variables
	// In a production environment, secrets should be managed via KMS or Vault.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on environment variables.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable must be set.")
	}

	// 2. Initialize Database (Simulated)
	// In a real application, this would be gorm.Open(...) and migration setup.
	log.Println("Initializing mock database and seeding admin user...")
	db := database.NewMockDB()
	
	// Inject the mock DB into handlers for dependency injection (simple approach for example)
	handlers.DB = db
	log.Println("Database initialization complete.")

	// 3. Setup Router
	r := router.NewRouter()

	// 4. Start Server
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}