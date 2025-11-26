package main

import (
	"log"
	"net/http"
	"secure-doc-service/config"
	"secure-doc-service/server"

	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Echo Server
	// Hiding the banner and setting a logger prefix for production readiness.
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetPrefix("SECURE_DOC_SVC")

	// 3. Setup Routes and Middleware
	server.SetupRouter(e, cfg)

	// 4. Start Server
	log.Printf("Starting secure documentation service on :%s", cfg.ServerPort)
	// Check for http.ErrServerClosed to allow graceful shutdowns without logging a fatal error.
	if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("Shutting down the server:", err)
	}
}