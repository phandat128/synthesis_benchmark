package main

import (
	"log"
	"net/http"
	"report_service/config"
	"report_service/handlers"
	"report_service/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator holds the Go Playground validator instance.
// Used to integrate standard struct validation into Echo.
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements echo.Validator interface.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return a 400 Bad Request error if validation fails
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Services
	// Note: In a real app, dbURL should be handled securely (e.g., secrets manager).
	dataService := services.NewDataService(cfg.DatabaseURL)
	pdfService := services.NewPDFGenerator()

	// 3. Initialize Handlers
	reportHandler := handlers.NewReportHandler(cfg, dataService, pdfService)

	// 4. Initialize Echo Framework
	e := echo.New()

	// 5. Setup Middleware
	// Standard security practices applied
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Apply basic security headers to prevent common attacks (XSS, Clickjacking, etc.)
	e.Use(middleware.Secure())
	// Configure CORS securely based on deployment needs
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Should be restricted in production
		AllowMethods: []string{http.MethodPost, http.MethodOptions},
	}))

	// 6. Setup Custom Validator
	// Used for basic struct validation (e.g., required fields, min values)
	e.Validator = &CustomValidator{validator: validator.New()}

	// 7. Define Routes
	v1 := e.Group("/api/v1")
	{
		// POST /api/v1/reports/generate
		v1.POST("/reports/generate", reportHandler.GenerateReport)
	}

	// 8. Start Server
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}