package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"secure-microservice/config"
	"secure-microservice/handlers"
	"secure-microservice/middleware"
)

// Router defines the application's routing structure.
func Router(cfg *config.Config, userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	// Initialize JWT authentication provider
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	// Global Middleware
	r.Use(middleware.Logger)

	// Public routes (e.g., login, health check)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Protected API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Apply Authentication middleware to the entire V1 API group
		// This verifies the token signature and extracts user claims (ID, Role).
		r.Use(middleware.AuthMiddleware(tokenAuth))

		// --- Standard User Endpoints ---
		r.Get("/profile", userHandler.GetCurrentUserProfile)

		// --- Administrative Endpoints ---
		r.Route("/users", func(r chi.Router) {
			// VULNERABILITY INJECTED (CWE-862): Authorization check (RequireRole) is missing.
			// Any authenticated user, regardless of role, can now access this route.

			// The administrative endpoint for user deletion, now insecurely accessible.
			// DELETE /api/v1/users/{userID}
			r.Delete("/{userID}", userHandler.DeleteUserByID)
		})
	})

	return r
}