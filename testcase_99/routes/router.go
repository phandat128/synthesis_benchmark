package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"forum_api/handlers"
)

// SetupRouter initializes the Chi router and defines all API routes.
func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// Standard secure middleware
	r.Use(middleware.Logger)
	// Recoverer catches panics and returns 500, preventing server crash
	r.Use(middleware.Recoverer)
	// RealIP helps correctly identify client IP
	r.Use(middleware.RealIP)
	
	// Security Middleware: Setting essential HTTP security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")
			// Basic Content Security Policy
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/topics", func(r chi.Router) {
		// Topics endpoints
		r.Post("/", handlers.CreateTopic) // POST /topics
		r.Get("/", handlers.GetTopics)    // GET /topics

		// Topic-specific endpoints
		r.Route("/{topicID}", func(r chi.Router) {
			// Comments endpoints
			r.Post("/comments", handlers.CreateComment) // POST /topics/{topicID}/comments
			r.Get("/comments", handlers.GetComments)    // GET /topics/{topicID}/comments
		})
	})

	return r
}