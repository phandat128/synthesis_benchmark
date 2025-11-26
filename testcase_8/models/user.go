package models

import "github.com/dgrijalva/jwt-go"

// User represents the structure of a user record in the database.
type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	PasswordHash   string `json:"-"` // Omit hash from JSON output for security
	Role           string `json:"role"` // e.g., "user", "admin"
}

// Claims defines the structure of the JWT claims, embedding standard claims.
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}