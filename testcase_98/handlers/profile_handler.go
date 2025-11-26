package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourproject/models"
	"github.com/yourproject/services"
)

// UpdateProfilePicture handles the POST request to update a user's profile picture URL.
func UpdateProfilePicture(c *gin.Context) {
	var req models.UpdateProfilePictureRequest

	// 1. Input Validation and Binding
	// Gin's binding:"required,url" tag provides initial validation.
	// The URL string (tainted input) is extracted here.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Do not leak internal errors, provide a generic bad request message.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing/invalid image_url field."})
		return
	}

	// 2. Call the secure service layer
	// The service layer implements the critical SSRF mitigation logic.
	bytesFetched, err := services.FetchAndStoreImage(req.ImageURL)

	if err != nil {
		// Proper Error Handling: Log the detailed error internally but return a sanitized message externally.
		fmt.Printf("SECURITY ALERT/ERROR: Failed to fetch image securely: %v\n", err)

		// If the error is related to network restrictions (SSRF block or invalid scheme/host),
		// return a specific forbidden status.
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Failed to process image URL due to security restrictions or network failure.",
		})
		return
	}

	// 3. Success Response
	// In a real application, we would update the database with the new secure image path.
	c.JSON(http.StatusOK, gin.H{
		"message":    "Profile picture updated successfully.",
		"url_used":   req.ImageURL,
		"bytes_read": bytesFetched,
	})
}