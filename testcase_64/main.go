package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"report_service/config"
	"report_service/handlers"
	"report_service/services"
)

func main() {
	// Load Configuration
	config.LoadConfig()

	// Initialize Database and Services
	generator, err := services.NewReportGenerator()
	if err != nil {
		log.Fatalf("Failed to initialize report generator service: %v", err)
	}
	// Ensure database connection is closed gracefully upon exit
	defer generator.Close()

	// Initialize Handlers
	reportHandler := handlers.NewReportHandler(generator)

	// Initialize Echo Framework
	e := echo.New()

	// --- Security Middleware and Configuration ---

	// 1. Standard security headers (X-Frame-Options, X-XSS-Protection, etc.)
	e.Use(middleware.Secure())
	
	// 2. CORS configuration (restrict origins in production)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// WARNING: In production, replace "*" with specific allowed origins.
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodPost, http.MethodOptions},
	}))

	// 3. Request logging
	e.Use(middleware.Logger())
	
	// 4. Basic Rate Limiting (Defense in Depth against DoS)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))

	// Define Routes
	v1 := e.Group("/api/v1")
	// Route for report generation
	v1.POST("/reports/generate", reportHandler.GenerateReport)

	// Start server gracefully
	go func() {
		serverAddr := ":" + config.AppConfig.ServerPort
		log.Printf("Starting server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server due to error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Define a context for shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}