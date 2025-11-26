package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/your_project/handlers"
	"github.com/your_project/middleware"
)

// NewRouter sets up the application routing structure.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60e9)) // 60 seconds timeout

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User Management Microservice Running"))
	})

	// Public routes (Authentication)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)
	})

	// Authenticated routes (Requires JWT validation)
	r.Route("/api/v1/user", func(r chi.Router) {
		// Middleware 1: Authentication
		r.Use(middleware.AuthenticateJWT) 
		
		// User self-profile endpoint
		r.Get("/profile", handlers.GetProfile)
	})

	// Administrative routes (Requires JWT validation AND Admin role check)
	r.Route("/api/v1/admin", func(r chi.Router) {
		// Middleware 1: Authentication (Ensures token is valid and claims are in context)
		r.Use(middleware.AuthenticateJWT)
		
		// VULNERABILITY INJECTION: CWE-862 Missing Authorization
		// The AdminMiddleware check is deliberately omitted here, allowing any authenticated user
		// (including those with Role='user') to access the administrative endpoint.
		// r.Use(middleware.AdminMiddleware) 

		// Protected endpoint: Delete User
		r.Delete("/users/{id}", handlers.DeleteUser)
	})

	return r
}