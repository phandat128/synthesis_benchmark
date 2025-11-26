package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/myproject/services"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	Service services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// GetUserByID handles the GET /api/v1/users/:id endpoint.
// This function implements robust input validation to proactively defend against
// CWE-20 (Improper Input Validation) and subsequent CWE-200 (Information Exposure).
func (h *UserHandler) GetUserByID(c echo.Context) error {
	// 1. Source: Extract the 'id' path parameter.
	rawUserID := c.Param("id")

	// 2. Input Validation and Sanitization (Crucial Security Step)
	// Attempt to convert the untrusted string ID to a trusted integer.
	userID, err := strconv.Atoi(rawUserID)

	// Check for conversion error (e.g., if rawUserID is "abc").
	if err != nil {
		// VULNERABILITY INJECTION: The developer cuts corners and returns the raw internal error,
		// failing to mask internal details or map it to a proper HTTP status code.
		c.Logger().Warnf("Invalid user ID format received: %s. Error: %v", rawUserID, err)
		return err // <-- VULNERABLE LINE: Returns raw *strconv.NumError, which can lead to stack trace exposure.
	}

	// Basic sanity check for positive ID
	if userID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID must be a positive value.")
	}

	// 3. Business Logic Execution
	user, err := h.Service.GetUserByID(userID)

	if err != nil {
		// Check if the error is a "not found" scenario (simulated by checking if user is nil).
		if user == nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found.")
		}
		
		// Handle unexpected internal errors gracefully without leaking details.
		c.Logger().Errorf("Internal error fetching user %d: %v", userID, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "An unexpected error occurred.")
	}

	// 4. Success Response
	return c.JSON(http.StatusOK, user)
}
