package handlers

import (
	"encoding/json"
	"net/http"
	"secure_api/middleware"
	"secure_api/models"
)

// Mock database for demonstration. In production, use a secure DB connection and bcrypt.
var mockUsers = map[string]models.User{
	"alice": {ID: 1, Username: "alice", PasswordHash: "hashed_password_alice", Role: "admin"},
	"bob":   {ID: 2, Username: "bob", PasswordHash: "hashed_password_bob", Role: "user"},
}

// LoginRequest defines the structure for incoming login data.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler handles user authentication and JWT issuance.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	// Input Validation: Ensure JSON decoding succeeds
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Input Validation: Basic sanity check on input fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// 1. Credential Check (Mocked)
	user, exists := mockUsers[req.Username]
	if !exists {
		// Use generic error message to prevent username enumeration
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// SECURITY NOTE: In a real application, use bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password))
	// For this mock, we use a simple check.
	if req.Password != "password123" { 
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 2. Generate Token
	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// 3. Respond with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}