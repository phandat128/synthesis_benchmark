package routes

import (
	"github.com/labstack/echo/v4"
	"secure-config-api/handlers"
)

// RegisterAPIRoutes sets up all API endpoints.
func RegisterAPIRoutes(e *echo.Echo, rh *handlers.ResourceHandler) {
	v1 := e.Group("/api/v1")

	// POST /api/v1/verify: Endpoint for verifying resource status.
	// This route handles the user input that must be securely processed.
	v1.POST("/verify", rh.VerifyResourceStatus)
}