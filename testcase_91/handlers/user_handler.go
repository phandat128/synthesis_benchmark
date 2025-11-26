package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/your_project/models"
)

// GetProfile handles fetching the authenticated user's own profile details.
func GetProfile(w http.ResponseWriter, r *http.Request) {
	// SECURE: Retrieve claims from context, ensuring the user is authenticated.
	claims, ok := models.GetClaimsFromContext(r.Context())
	if !ok {
		// Defensive check: Should be caught by AuthenticateJWT middleware
		http.Error(w, "Unauthorized: Claims missing", http.StatusUnauthorized)
		return
	}

	// Simulate fetching user data based on the authenticated ID
	user, err := DB.GetUserByID(claims.UserID)
	if err != nil {
		// Use generic error if user is somehow missing from DB but has a valid token
		http.Error(w, "User data retrieval failed", http.StatusInternalServerError)
		return
	}

	// SECURE: Only return safe, non-sensitive data (exclude PasswordHash)
	safeUser := struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeUser)
}