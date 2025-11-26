package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourproject/handlers"
)

// SetupRouter defines all API routes for the application.
func SetupRouter() *gin.Engine {
	// Use gin.ReleaseMode in production for better performance and security (less verbose errors)
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		profile := v1.Group("/profile")
		{
			// Endpoint for updating the profile picture, which accepts the external URL.
			// This is the entry point for the tainted input.
			profile.POST("/update_picture", handlers.UpdateProfilePicture)
		}
	}

	return r
}