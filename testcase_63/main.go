package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"secure-config-api/config"
	"secure-config-api/handlers"
	"secure-config-api/routes"
	"secure-config-api/services"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("FATAL: Failed to load configuration: %v\n", err)
		return
	}

	// 2. Initialize Echo Server
	e := echo.New()

	// 3. Register Middleware
	// Standard security middleware setup
	e.Use(middleware.Logger())
	// Use Recover middleware to gracefully handle panics without crashing the server
	e.Use(middleware.Recover())
	// Configure CORS securely (adjust origins for production)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	// Apply basic security headers (XSS, Content Type Sniffing, etc.)
	e.Use(middleware.Secure())

	// 4. Initialize Services and Handlers
	systemService := services.NewSystemService()
	resourceHandler := &handlers.ResourceHandler{
		SystemService: systemService,
	}

	// 5. Register Routes
	routes.RegisterAPIRoutes(e, resourceHandler)

	// Health Check Route
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "API Operational")
	})

	// 6. Start Server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	fmt.Printf("Starting server securely on %s\n", serverAddr)
	if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
		e.Logger().Fatal("Shutting down the server:", err)
	}
}