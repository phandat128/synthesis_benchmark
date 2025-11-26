package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all application configuration settings.
type Config struct {
	ServerPort     string
	DatabaseURL    string
	MaxReportLimit int // CRITICAL: Defines the maximum safe limit for records to prevent DoS (CWE-400)
}

// LoadConfig initializes and returns the application configuration.
func LoadConfig() *Config {
	// Set sensible defaults
	vip.SetDefault("SERVER_PORT", "8080")
	// Default limit is 50,000 records. This must be tuned based on server resources.
	vip.SetDefault("MAX_REPORT_LIMIT", 50000)

	// Read environment variables
	vip.AutomaticEnv()

	cfg := &Config{
		ServerPort:  vip.GetString("SERVER_PORT"),
		// Use os.Getenv for sensitive data like DB credentials, ensuring they are not logged or stored in config files.
		DatabaseURL: os.Getenv("DATABASE_URL"), 
	}

	// Load and validate MaxReportLimit
	limitStr := vip.GetString("MAX_REPORT_LIMIT")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Printf("Warning: Invalid or missing MAX_REPORT_LIMIT (%s). Using default 50000.", limitStr)
		cfg.MaxReportLimit = 50000
	} else {
		cfg.MaxReportLimit = limit
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}

	log.Printf("Configuration loaded. Max Report Limit: %d", cfg.MaxReportLimit)
	return cfg
}