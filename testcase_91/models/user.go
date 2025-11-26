package models

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ContextKey is used for context value keys
type ContextKey string

const (
	// UserContextKey is the key used to store user claims in the request context
	UserContextKey ContextKey = "userClaims"
	// RoleAdmin defines the admin role string
	RoleAdmin = "admin"
	// RoleUser defines the standard user role string
	RoleUser = "user"
)

// User represents a user record in the database
type User struct {
	ID           uint `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserClaims holds the data extracted from the JWT token and stored in the context
type UserClaims struct {
	UserID   uint
	Username string
	Role     string
}

// HashPassword hashes the plaintext password using bcrypt
func (u *User) HashPassword(password string) error {
	// SECURE: Use bcrypt with default cost for strong hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

// CheckPasswordHash compares a plaintext password with the stored hash
func (u *User) CheckPasswordHash(password string) bool {
	// SECURE: Constant time comparison to prevent timing attacks
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GetClaimsFromContext retrieves UserClaims from the request context
func GetClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*UserClaims)
	return claims, ok
}