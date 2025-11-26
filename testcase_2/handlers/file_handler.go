package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"document_service/config"
	"document_service/storage"
)

// DownloadFileHandler handles requests to retrieve documents.
// It implements robust path traversal prevention (CWE-22).
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Input Extraction
	requestedFilename := r.URL.Query().Get("filename")
	if requestedFilename == "" {
		http.Error(w, "Bad Request: 'filename' parameter is required.", http.StatusBadRequest)
		return
	}

	// --- VULNERABILITY INJECTION: Path Traversal (CWE-22) ---

	// 2. Construct the intended path by joining the secure base directory and user input.
	// We use filepath.Join to handle OS-specific separators securely.
	fullPath := filepath.Join(config.BaseStorageDir, requestedFilename)

	// 3. (VULNERABILITY INTRODUCED: Removed path cleaning and validation checks)
	// The fullPath variable now holds the user-supplied path potentially containing '..' sequences.

	// 6. Final Validation: Check if the file exists and is a regular file.
	if !storage.FileExists(fullPath) {
		log.Printf("File not found or inaccessible: %s", fullPath)
		// Use a generic error message to avoid leaking file system structure information.
		http.Error(w, "Not Found: The requested document does not exist.", http.StatusNotFound)
		return
	}

	// 7. Serve the file securely.
	log.Printf("Serving file: %s", fullPath)

	// Set Content-Disposition header to force download and use only the base filename
	// to prevent injection into the header value.
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(fullPath)))

	// http.ServeFile handles reading the file, setting Content-Type, and Content-Length.
	http.ServeFile(w, r, fullPath)
}