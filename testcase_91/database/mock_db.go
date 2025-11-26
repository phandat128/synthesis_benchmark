package database

import (
	"errors"
	"sync"

	"github.com/your_project/models"
)

// MockDB simulates database operations for demonstration purposes.
type MockDB struct {
	users map[uint]*models.User
	nextID uint
	mu sync.RWMutex
}

// NewMockDB initializes the mock database with a starting admin user.
func NewMockDB() *MockDB {
	db := &MockDB{
		users: make(map[uint]*models.User),
		nextID: 1,
	}

	// Initialize a default Admin user for testing RBAC
	adminUser := &models.User{
		ID: 1001,
		Username: "admin_user",
		Role: models.RoleAdmin,
	}
	adminUser.HashPassword("SecureAdminPwd123!") // Hash is simulated
	db.users[1001] = adminUser
	db.nextID = 1002

	// Initialize a default Standard user
	standardUser := &models.User{
		ID: 2001,
		Username: "standard_user",
		Role: models.RoleUser,
	}
	standardUser.HashPassword("SecureUserPwd123!")
	db.users[2001] = standardUser
	db.nextID = 2002

	return db
}

// CreateUser simulates saving a new user to the database.
func (m *MockDB) CreateUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user.ID = m.nextID
	m.users[user.ID] = user
	m.nextID++
	return nil
}

// GetUserByUsername retrieves a user by username.
func (m *MockDB) GetUserByUsername(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByID retrieves a user by ID.
func (m *MockDB) GetUserByID(id uint) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UserExists checks if a username is already taken.
func (m *MockDB) UserExists(username string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Username == username {
			return true
		}
	}
	return false
}

// DeleteUserByID simulates deleting a user.
func (m *MockDB) DeleteUserByID(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}