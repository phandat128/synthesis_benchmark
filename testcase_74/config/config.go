package config

import (
	"log"
	"os"
	"path/filepath"
)

// AppConfig holds application-wide configuration settings.
type AppConfig struct {
	BaseDownloadDir string // Absolute, canonical path to the documentation root.
	ServerPort      string
}

// LoadConfig initializes and returns the application configuration.
func LoadConfig() *AppConfig {
	// Load base directory from environment variable or use a secure default.
	baseDir := os.Getenv("DOCS_BASE_DIR")
	if baseDir == "" {
		// Default path relative to the executable.
		baseDir = "./data/docs"
	}

	// Ensure the path is absolute and clean for reliable path traversal checks.
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for base directory: %v", err)
	}

	// Canonicalize the path to remove any trailing slashes or redundant elements.
	absBaseDir = filepath.Clean(absBaseDir)

	// Create the directory if it doesn't exist (useful for local development/testing)
	if _, err := os.Stat(absBaseDir); os.IsNotExist(err) {
		log.Printf("Creating documentation directory: %s", absBaseDir)
		// Use restrictive permissions (0755)
		if err := os.MkdirAll(absBaseDir, 0755); err != nil {
			log.Fatalf("Failed to create base directory: %v", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Configuration loaded. Base Download Directory (Canonical): %s", absBaseDir)

	return &AppConfig{
		BaseDownloadDir: absBaseDir,
		ServerPort:      port,
	}
}