package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your_project/models"
)

// DeleteUser handles the deletion of a user by ID.
// This endpoint is protected by AdminMiddleware in router/router.go.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Authorization Check: The role check (CWE-862 mitigation) is enforced by AdminMiddleware
	// before this handler is executed.
	claims, ok := models.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Authentication context missing", http.StatusUnauthorized)
		return
	}

	// 1. Input Validation (URL Parameter)
	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil || targetID == 0 {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// 2. Least Privilege Check (Prevent admin from deleting their own active session account)
	if uint(targetID) == claims.UserID {
		http.Error(w, "Forbidden: Cannot delete your own account via this endpoint", http.StatusForbidden)
		return
	}

	// 3. Execution (Simulated DB operation)
	if err := DB.DeleteUserByID(uint(targetID)); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		// SECURE: Generic error for internal issues
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully", "deleted_id": targetIDStr})
}