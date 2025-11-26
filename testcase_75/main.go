package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"report_generator/internal/data"
	"report_generator/internal/handlers"
	"report_generator/internal/services"
)

// Config holds application-wide settings.
type Config struct {
	// CRITICAL SECURITY CONTROL: Defines the maximum number of records allowed per report.
	// This prevents resource exhaustion (OOM/DoS) from malicious or accidental large inputs.
	MaxReportRecords int
	ServerPort       string
}

func main() {
	// 1. Load Configuration
	cfg := Config{
		MaxReportRecords: 5000, // Strict limit for DoS mitigation
		ServerPort:       ":8080",
	}
	log.Printf("Configuration loaded: MaxReportRecords=%d", cfg.MaxReportRecords)

	// 2. Initialize Dependencies (DB, Services, Handlers)

	// Initialize Mock DB connection (simulating a secure connection pool)
	dbConn := &data.MockDBConnection{}
	
	repo := data.NewRepository(dbConn)
	
	// Pass the security configuration (max limit) down to the service layer
	reportService := services.NewReportService(repo, cfg.MaxReportRecords)
	
	// Handlers receive the max limit for immediate input validation feedback
	reportHandler := handlers.NewReportHandler(reportService, cfg.MaxReportRecords)

	// 3. Initialize Router
	r := chi.NewRouter()

	// 4. Apply Standard Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// Secure recovery middleware to prevent panic details from leaking
	r.Use(middleware.Recoverer)
	// Set a reasonable timeout for the entire request lifecycle
	r.Use(middleware.Timeout(60 * time.Second))

	// 5. Define Routes
	setupRoutes(r, reportHandler)

	// 6. Start Server
	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 120 * time.Second, // Increased write timeout for potentially long PDF generation
		IdleTimeout:  30 * time.Second,
	}

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting securely on %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", cfg.ServerPort, err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting.")
}