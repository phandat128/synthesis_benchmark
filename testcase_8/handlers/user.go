package handlers

import (
	"encoding/json"
	"net/http"
	"secure_api/middleware"
	"secure_api/models"
)

// GetProfileHandler retrieves the profile of the currently authenticated user.
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Retrieve UserID from context (set by middleware/auth.go)
	userID, ok := r.Context().Value(middleware.ContextUserIDKey).(int)
	if !ok {
		// This indicates a failure in the middleware chain or an unauthenticated request
		http.Error(w, "Authentication context missing", http.StatusUnauthorized)
		return
	}

	// 2. Fetch user data (Mocked lookup)
	var currentUser models.User
	found := false
	// Iterate mockUsers defined in handlers/auth.go
	for _, user := range mockUsers {
		if user.ID == userID {
			currentUser = user
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// 3. Respond with user data (excluding sensitive fields like PasswordHash)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser)
}