package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"secure-microservice/middleware"
	"secure-microservice/models"
)

// UserHandler handles HTTP requests related to user management.
type UserHandler struct {
	Repo models.UserRepository
}

// NewUserHandler creates a new handler instance.
func NewUserHandler(repo models.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// DeleteUserByID handles the deletion of a user account.
// Access to this function is restricted to 'admin' roles by the router configuration.
func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	// 1. Input Validation: Retrieve and validate the target UserID from the URL path.
	userIDStr := chi.URLParam(r, "userID")
	if userIDStr == "" {
		http.Error(w, "Bad Request: User ID is required", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.ParseUint(userIDStr, 10, 64)
	// Robust validation against non-numeric or excessively large IDs.
	if err != nil || targetID == 0 {
		http.Error(w, "Bad Request: Invalid User ID format", http.StatusBadRequest)
		return
	}

	// 2. Retrieve Caller ID (for internal logic, e.g., preventing self-deletion)
	callerIDVal := r.Context().Value(middleware.UserIDKey)
	// Authorization (Role check) is handled securely by the middleware layer.
	
	// Safety check: ensure caller ID exists
	if callerIDVal == nil {
		http.Error(w, "Internal Error: Caller identity missing from context", http.StatusInternalServerError)
		return
	}
	callerID := callerIDVal.(uint)

	// Security check: Prevent admin from deleting themselves via this endpoint.
	// This enforces separation of concerns and prevents accidental lockout.
	if uint(targetID) == callerID {
		http.Error(w, "Forbidden: Cannot delete your own account via this administrative endpoint", http.StatusForbidden)
		return
	}

	// 3. Execute the deletion using the secure repository method (GORM prevents injection)
	err = h.Repo.Delete(uint(targetID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found or already deleted" {
			http.Error(w, "Not Found: User not found", http.StatusNotFound)
			return
		}
		// Proper Error Handling: Do not leak internal database errors (CWE-209)
		http.Error(w, "Internal Server Error: Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User successfully deleted"})
}

// GetCurrentUserProfile retrieves the profile of the authenticated user.
func (h *UserHandler) GetCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	callerIDVal := r.Context().Value(middleware.UserIDKey)
	if callerIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	callerID := callerIDVal.(uint)

	user, err := h.Repo.FindByID(callerID)
	if err != nil {
		http.Error(w, "User profile not found", http.StatusNotFound)
		return
	}

	// Output Encoding: Safe for JSON API response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}