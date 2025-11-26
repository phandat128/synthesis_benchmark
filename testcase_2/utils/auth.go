package utils

import (
	"net/http"
)

// AuthMiddleware provides a simple, placeholder authentication check.
// In a real application, this would involve robust JWT validation, session checks, or OAuth.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for a specific header as a basic authentication mechanism
		userToken := r.Header.Get("X-Auth-Token")

		// Require a specific secure token for access
		if userToken != "secure-admin-token-123" {
			// Use StatusUnauthorized and avoid leaking specific reasons
			http.Error(w, "Unauthorized access.", http.StatusUnauthorized)
			return
		}

		// If authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}