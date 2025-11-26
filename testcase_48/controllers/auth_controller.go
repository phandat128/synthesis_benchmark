package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"secure_api/models"
)

// AuthController holds the database connection.
type AuthController struct {
	DB *gorm.DB
}

// RegisterRequest defines the expected structure for registration input.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest defines the expected structure for login input.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// generateJWT creates a signed JWT token.
func generateJWT(userID uint, role string) (string, error) {
	// SECURITY NOTE: Secret should be loaded from environment variables.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_key_for_dev_only" // Placeholder
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Register handles new user creation.
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	// Input Validation: Gin's binding handles required fields and basic length checks.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or constraints met: " + err.Error()})
		return
	}

	// Security: Hash the password using bcrypt.
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		// Proper Error Handling: Do not leak internal hashing errors.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         "user", // Default role for new users
	}

	// Security: Check for unique username constraint violation.
	result := ac.DB.Create(&user)
	if result.Error != nil {
		// Check for common unique constraint errors (Postgres/SQLite)
		if result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"uni_users_username\"" ||
			result.Error.Error() == "UNIQUE constraint failed: users.username" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		// Proper Error Handling: Generic database error message.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during registration"})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// Login handles user authentication and token issuance.
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var user models.User
	// Security: Use GORM (ORM) to safely query the database, preventing SQL injection.
	if err := ac.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Security: Use generic message to prevent username enumeration.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Security: Compare password hash securely.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// Security: Use generic message to prevent timing attacks/enumeration.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}