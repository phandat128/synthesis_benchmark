package main

import (
	"log"
	"net/http"
	"time"

	"document_service/config"
	"document_service/handlers"
	"document_service/utils"
)

const listenAddr = ":8080"

func main() {
	// Initialize configuration (ensures storage directory exists and dummy files are present)
	config.Init()

	// Setup Router
	mux := http.NewServeMux()

	// Define the secure handler chain
	// Route: /api/v1/download?filename=...
	downloadHandler := http.HandlerFunc(handlers.DownloadFileHandler)

	// Apply Authentication Middleware to secure the endpoint
	secureDownloadHandler := utils.AuthMiddleware(downloadHandler)

	mux.Handle("/api/v1/download", secureDownloadHandler)

	// Setup HTTP Server with secure timeouts to mitigate slowloris attacks
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // Max time to read the entire request
		WriteTimeout: 10 * time.Second, // Max time to write the response
		IdleTimeout:  120 * time.Second, // Max time for keep-alive connections
	}

	log.Printf("Starting secure document service on %s", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", listenAddr, err)
	}
}