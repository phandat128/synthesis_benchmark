package main

import (
	"net/http"
	"report_service/handlers"
	"report_service/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRouter initializes the Chi router with middleware and routes.
func SetupRouter(dataService *services.DataService) *chi.Mux {
	r := chi.NewRouter()

	// Standard Secure Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// Recoverer catches panics and returns a 500, preventing server crash.
	r.Use(middleware.Recoverer)
	
	// Secure Timeout Middleware (Prevents slowloris/resource exhaustion on connection)
	r.Use(middleware.Timeout(60 * http.Second))

	reportHandler := handlers.NewReportHandler(dataService)

	r.Route("/api/v1", func(r chi.Router) {
		// Endpoint for generating reports
		r.Post("/reports/generate", reportHandler.GenerateReport)
	})

	// Serve static files (e.g., generated reports) securely
	// NOTE: In a production environment, reports should be served via a dedicated,
	// authenticated file server or CDN, not directly via the API service.
	// We use a basic file server setup here for demonstration.
	fileServer := http.FileServer(http.Dir("./reports"))
	r.Handle("/reports/*", http.StripPrefix("/reports/", fileServer))

	return r
}