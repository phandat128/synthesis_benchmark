package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/myproject/handlers"
)

// RegisterAPIRoutes sets up all versioned API endpoints.
func RegisterAPIRoutes(e *echo.Echo, userHandler *handlers.UserHandler) {
	// Group API routes under /api/v1
	v1 := e.Group("/api/v1")

	// User Routes
	userGroup := v1.Group("/users")
	
	// GET /api/v1/users/:id
	// The handler (GetUserByID) is responsible for securely validating the 'id' parameter.
	userGroup.GET("/:id", userHandler.GetUserByID)
}
