package server

import (
	"secure-doc-service/config"
	"secure-doc-service/handlers"
	"secure-doc-service/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// SetupRouter initializes and configures the Echo router.
func SetupRouter(e *echo.Echo, cfg *config.AppConfig) {
	// Global Middleware
	// Logger middleware for request tracing and auditing
	e.Use(echoMiddleware.Logger())
	// Recover middleware to prevent application crashes from panics
	e.Use(echoMiddleware.Recover())

	// Secure CORS configuration (restrict origins in production)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD},
	}))

	// Handlers initialization
	fileHandler := handlers.NewFileHandler(cfg)

	// API Group
	v1 := e.Group("/api/v1")

	// Health Check (unauthenticated)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok", "service": "doc-service"})
	})

	// Secure File Endpoints
	fileGroup := v1.Group("/files")
	// Apply custom authentication middleware to enforce access control
	fileGroup.Use(middleware.AuthMiddleware)

	// Route for downloading files securely
	// GET /api/v1/files/download?filename=document.pdf
	fileGroup.GET("/download", fileHandler.DownloadFile)
}