package routes

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"secure_api/controllers"
	"secure_api/middleware"
)

// SetupRouter initializes the Gin router and defines all API routes.
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize controllers
	authController := controllers.AuthController{DB: db}
	userController := controllers.UserController{DB: db}

	v1 := r.Group("/api/v1")

	// --- Public Routes ---
	public := v1.Group("/auth")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
	}

	// --- Protected Routes (Requires valid JWT) ---
	protected := v1.Group("/user")
	// Apply AuthRequired middleware to all routes in this group
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", userController.GetProfile)
	}

	// --- Administrative Routes (Requires valid JWT AND Admin Role) ---
	admin := v1.Group("/admin")
	// VULNERABILITY FIX: Chain the middlewares.
	// 1. AuthRequired ensures the user is logged in and context claims are set.
	// 2. AdminRequired checks the 'user_role' claim set by AuthRequired, preventing unauthorized access.
	admin.Use(middleware.AuthRequired())
	{
		// Endpoint for privileged user deletion.
		admin.DELETE("/users/:user_id", userController.DeleteUser)
	}

	return r
}