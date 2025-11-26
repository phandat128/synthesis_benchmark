package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"secure_api/models"
)

// UserController holds the database connection.
type UserController struct {
	DB *gorm.DB
}

// GetProfile retrieves the profile of the currently authenticated user.
func (uc *UserController) GetProfile(c *gin.Context) {
	// Security: Retrieve user ID from the context, set by AuthRequired middleware.
	userID, exists := c.Get("user_id")
	if !exists {
		// This should be caught by AuthRequired, but acts as a defensive check.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication context missing"})
		return
	}

	var user models.User
	// Security: Ensure we only retrieve the profile corresponding to the authenticated user ID (IDOR prevention).
	if err := uc.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return	
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteUser handles the administrative deletion of a user by ID.
// This endpoint is protected by BOTH AuthRequired and AdminRequired middleware (defined in routes.go).
func (uc *UserController) DeleteUser(c *gin.Context) {
	// 1. Input Validation: Get and sanitize/validate the target ID from path parameters.
	targetUserIDStr := c.Param("user_id")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 64)
	if err != nil || targetUserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 2. Authorization Check (Defense against the described vulnerability):
	// The primary role check is handled by the AdminRequired middleware.
	// We add a safety check here to prevent an admin from deleting themselves.
	currentUserID, _ := c.Get("user_id")
	if uint(currentUserID.(uint)) == uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete your own account via this administrative endpoint"})
		return
	}

	// 3. Database Operation (Safe ORM usage)
	var user models.User
	// Check if the user exists before attempting deletion.
	result := uc.DB.First(&user, targetUserID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during lookup"})
		return
	}

	// Perform the deletion using GORM's safe Delete method.
	deleteResult := uc.DB.Delete(&user)
	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}