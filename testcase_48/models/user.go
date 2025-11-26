package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the database model for a user.
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:'user';not null"` // 'user' or 'admin'
}

// UserResponse is used for safe output, excluding the password hash.
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// ToResponse converts a User model to a safe UserResponse struct.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}

// HashPassword securely hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	// Use DefaultCost for a good balance of security and performance.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}