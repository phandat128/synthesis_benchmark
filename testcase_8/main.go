package main

import (
	"log"
	"net/http"
	"secure_api/handlers"
	"secure_api/middleware"
)

// main initializes the router, registers routes, and starts the server.
func main() {
	// Initialize the router (using standard library ServeMux for simplicity)
	mux := http.NewServeMux()

	// --- Public Routes ---
	mux.HandleFunc("/api/v1/login", handlers.LoginHandler)

	// --- Authenticated Routes (Requires JWT) ---
	// User Profile
	mux.HandleFunc("/api/v1/profile", middleware.Authenticate(handlers.GetProfileHandler))

	// --- Administrative Routes (Requires JWT AND Admin Role) ---
	// VULNERABILITY INJECTION: Authorization check (RequireRole) is removed,
	// allowing any authenticated user to access this administrative endpoint.
	adminHandler := middleware.Authenticate(
		handlers.DeleteUserHandler,
	)
	mux.HandleFunc("/api/v1/admin/users/delete", adminHandler)

	// Start the server
	port := ":8080"
	log.Printf("Server starting securely on port %s", port)
	// In production, use http.ListenAndServeTLS and configure timeouts.
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}