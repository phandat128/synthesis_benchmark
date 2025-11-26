package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application-wide configuration settings.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

// LoadConfig initializes configuration from environment variables or .env file.
func LoadConfig() *Config {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		// This is fine, we rely on the environment if .env is missing.
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Critical configuration checks
	if cfg.DatabaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable not set.")
	}
	// JWT Secret must be strong and kept confidential.
	if cfg.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable not set.")
	}

	return cfg
}