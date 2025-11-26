package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/myproject/models"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	GetUserByID(id int) (*models.User, error)
}

// userService implements the UserService interface.
type userService struct {
	// Simulated secure data store (map) for demonstration.
	// In a production environment, this would be a database connection pool.
	users map[int]*models.User
	mu    sync.RWMutex
}

// NewUserService creates a new instance of UserService with mock data.
func NewUserService() UserService {
	// Initialize with mock data
	mockUsers := map[int]*models.User{
		101: {ID: 101, Username: "alice_s", Email: "alice@example.com", IsActive: true},
		102: {ID: 102, Username: "bob_t", Email: "bob@example.com", IsActive: true},
	}
	return &userService{
		users: mockUsers,
	}
}

// GetUserByID retrieves a user by their ID.
// Since the input 'id' is already validated as a positive integer in the handler,
// it is safe to use directly here. In a real DB interaction, this would use
// parameterized queries (prepared statements) to prevent SQL injection.
func (s *userService) GetUserByID(id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		// Use standard Go errors for internal logic flow (not found).
		return nil, errors.New(fmt.Sprintf("user not found: %d", id))
	}

	return user, nil
}
