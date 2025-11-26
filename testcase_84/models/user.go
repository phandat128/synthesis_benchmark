package models

// User represents a user profile in the system.
// Note: Sensitive fields like passwords should never be included in this transport model.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}
