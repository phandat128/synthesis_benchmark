package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthRequired verifies the JWT token and extracts claims (UserID, Role) into the context.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString = parts[1]

		// Use a secure secret key (ideally loaded from environment variables)
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			// SECURITY NOTE: Placeholder. Must be secure in production.
			jwtSecret = "super_secret_key_for_dev_only"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract and set claims securely in the context
		// Type assertion safety check is omitted here as float64 is standard for JSON numbers in Go JWT claims.
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("user_role", claims["role"].(string))

		c.Next()
	}
}

// AdminRequired checks if the user role stored in the context is 'admin'.
// This middleware MUST be run AFTER AuthRequired.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			// Defensive check: If AuthRequired failed to run or set the role.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication context missing role information"})
			c.Abort()
			return
		}

		// PROACTIVE DEFENSE AGAINST VULNERABILITY:
		// This explicit role check prevents unprivileged users from accessing administrative endpoints.
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin privileges required."})
			c.Abort()
			return
		}

		c.Next()
	}
}