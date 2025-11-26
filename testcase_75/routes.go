package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"report_generator/internal/handlers"
)

// setupRoutes defines all application routes.
func setupRoutes(r *chi.Mux, h *handlers.ReportHandler) {
	// Health check route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Report API operational."))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Secure endpoint for report generation. Input validation is critical here.
		r.Get("/reports/generate", h.GenerateReport)
	})
}