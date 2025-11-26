package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secure_api/models"
	"secure_api/routes"
)

// initDB initializes the database connection and runs migrations.
func initDB() *gorm.DB {
	// SECURITY NOTE: In a production environment, DSN must be loaded securely from environment variables or a secret manager.
	// Using a placeholder DSN for demonstration.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=user password=password dbname=secureapi port=5432 sslmode=disable TimeZone=UTC"
		log.Println("Using default placeholder DSN. Set DATABASE_URL environment variable for production.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database connection successful and migrations run.")

	// Optional: Seed an admin user if none exists (for testing the admin endpoint)
	var admin models.User
	// Check if an admin user already exists
	if db.Where("role = ?", "admin").First(&admin).Error == gorm.ErrRecordNotFound {
		// SECURITY NOTE: Hardcoding credentials is bad practice. This is for demonstration setup only.
		hashedPassword, _ := models.HashPassword("SecureAdminPass123!")
		adminUser := models.User{
			Username:     "admin",
			PasswordHash: hashedPassword,
			Role:         "admin",
		}
		db.Create(&adminUser)
		log.Println("Default admin user created (Username: admin)")
	}

	return db
}

func main() {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	db := initDB()
	r := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}