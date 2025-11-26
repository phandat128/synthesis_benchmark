package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration settings.
type Config struct {
	ServerPort          string
	DatabaseURL         string
	MaxReportRecords    int // CRITICAL: Defines the secure upper bound for report size.
	ReportTimeout       time.Duration
}

var AppConfig Config

// LoadConfig initializes configuration from environment variables or default values.
func LoadConfig() {
	viper.SetDefault("SERVER_PORT", "8080")
	// WARNING: Use environment variables or secrets manager for production DB credentials.
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/reportsdb?sslmode=disable")
	
	// Set a reasonable, secure default limit (5000 records) to prevent resource exhaustion (DoS).
	viper.SetDefault("MAX_REPORT_RECORDS", 5000) 
	viper.SetDefault("REPORT_TIMEOUT_SECONDS", 30)

	viper.AutomaticEnv() // Read environment variables

	AppConfig.ServerPort = viper.GetString("SERVER_PORT")
	AppConfig.DatabaseURL = viper.GetString("DATABASE_URL")
	AppConfig.MaxReportRecords = viper.GetInt("MAX_REPORT_RECORDS")
	AppConfig.ReportTimeout = time.Duration(viper.GetInt("REPORT_TIMEOUT_SECONDS")) * time.Second

	// Runtime validation of critical configuration
	if AppConfig.MaxReportRecords <= 0 {
		log.Fatal("Configuration error: MAX_REPORT_RECORDS must be a positive integer.")
	}

	log.Printf("Configuration loaded: Max Report Records allowed: %d", AppConfig.MaxReportRecords)
}