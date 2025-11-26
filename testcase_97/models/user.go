package models

import (
	"errors"
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`	
	Role     string `gorm:"default:'standard'"` // standard, admin
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	FindByID(id uint) (*User, error)
	Delete(id uint) error
}

// GormUserRepository implements UserRepository using GORM.
type GormUserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new repository instance.
func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

// FindByID retrieves a user by their ID.
func (r *GormUserRepository) FindByID(id uint) (*User, error) {
	var user User
	// Using GORM's built-in methods ensures parameterized queries, preventing SQL injection.
	result := r.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Delete removes a user record by ID.
func (r *GormUserRepository) Delete(id uint) error {
	// GORM handles the safe deletion query.
	result := r.DB.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	// Check if any row was affected to provide meaningful feedback.
	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}
	return nil
}