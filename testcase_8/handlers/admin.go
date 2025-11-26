package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteUserRequest defines the structure for the user ID to be deleted.
type DeleteUserRequest struct {
	UserID int `json:"user_id"`
}

// DeleteUserHandler handles the administrative task of deleting a user.
// This handler relies entirely on the middleware chain (Authenticate -> RequireRole("admin"))
// to ensure only authorized users can execute this sensitive action.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteUserRequest
	// Input Validation: Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Input Validation: Ensure UserID is positive and valid
	if req.UserID <= 0 {
		http.Error(w, "Invalid User ID provided", http.StatusBadRequest)
		return
	}

	// 1. Perform Deletion (Mocked)
	// SECURITY NOTE: In production, use parameterized queries to prevent SQL Injection.

	deleted := false
	var deletedUsername string
	
	// Find the user to delete by ID
	for username, user := range mockUsers {
		if user.ID == req.UserID {
			deletedUsername = username
			deleted = true
			break
		}
	}

	if deleted {
		delete(mockUsers, deletedUsername)
	} else {
		http.Error(w, "User not found or already deleted", http.StatusNotFound)
		return
	}

	// 2. Success Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User ID %d successfully deleted.", req.UserID),
	})
}