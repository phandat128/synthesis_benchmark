package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secure-microservice/config"
	"secure-microservice/handlers"
	"secure-microservice/models"
	"secure-microservice/routes"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection (GORM)
	// Using the full DSN from the environment variable (DATABASE_URL) is standard practice.
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	// In a production environment, dedicated migration tools (like goose or migrate) are preferred.
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	// 3. Initialize Repositories and Handlers
	userRepo := models.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// 4. Setup Router
	router := routes.Router(cfg, userHandler)

	// 5. Start Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server securely on %s", addr)

	// Start the HTTP server
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}