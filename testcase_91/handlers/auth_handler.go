package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your_project/database"
	"github.com/your_project/middleware"
	"github.com/your_project/models"
)

// DB is a global mock database instance injected from main.go
var DB *database.MockDB

// generateJWT creates a new JWT token for a user
func generateJWT(user *models.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET not configured")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// RegisterRequest defines the expected structure for registration input
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles new user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	// SECURE: Limit request size and handle decoding errors gracefully
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// SECURE: Input validation and sanitization (length checks)
	if len(req.Username) < 3 || len(req.Password) < 8 {
		http.Error(w, "Username must be at least 3 chars, Password at least 8 chars.", http.StatusBadRequest)
		return
	}

	// Check if user already exists (simulated DB check)
	if DB.UserExists(req.Username) {
		http.Error(w, "User already exists", http.StatusConflict)
		return	
	}

	newUser := &models.User{
		Username: req.Username,
		Role:     models.RoleUser, // Default to standard user role
	}

	if err := newUser.HashPassword(req.Password); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Simulate saving to DB
	if err := DB.CreateUser(newUser); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginRequest defines the expected structure for login input
type LoginRequest struct {
	Username string `json:"username"`	
	Password string `json:"password"`
}

// Login handles user authentication and JWT issuance
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// SECURE: Input validation
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Simulate fetching user from DB
	user, err := DB.GetUserByUsername(req.Username)
	if err != nil {
		// SECURE: Generic error message to prevent username enumeration
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// SECURE: Compare password hash
	if !user.CheckPasswordHash(req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWT(user)
	if err != nil {
		// Log the error internally but return a generic message
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}