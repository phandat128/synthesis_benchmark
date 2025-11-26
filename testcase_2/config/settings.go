package config

import (
	"log"
	"os"
	"path/filepath"
)

// BaseStorageDir defines the root directory where all documents are stored.
// This path should be configured securely and ideally outside the web root.
// Using a relative path here for simplicity, but absolute paths are preferred in production.
const BaseStorageDir = "./documents"

// Init ensures the necessary directories exist and sets up logging.
func Init() {
	// Ensure the storage directory exists
	// Use 0750 permissions (owner read/write/execute, group read/execute, others none)
	if _, err := os.Stat(BaseStorageDir); os.IsNotExist(err) {
		log.Printf("Creating storage directory: %s", BaseStorageDir)
		if err := os.MkdirAll(BaseStorageDir, 0750); err != nil {
			log.Fatalf("Failed to create storage directory: %v", err)
		}
		// Create a dummy file for testing the secure handler
		dummyPath := filepath.Join(BaseStorageDir, "secure_report_q3.pdf")
		if err := os.WriteFile(dummyPath, []byte("This is a secure report content for Q3."), 0644); err != nil {
			log.Printf("Warning: Could not create dummy file: %v", err)
		} else {
			log.Printf("Created dummy file: %s", dummyPath)
		}
	}
	log.Println("Configuration initialized successfully.")
}