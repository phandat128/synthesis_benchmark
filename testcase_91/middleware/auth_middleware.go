package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your_project/models"
)

// JWTClaims defines the structure for JWT claims, including custom fields.
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`	
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthenticateJWT validates the JWT token and injects user claims into the context.
func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secretKey := os.Getenv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// SECURE: Check signing method to prevent 'none' algorithm attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			http.Error(w, "Invalid token claims structure", http.StatusUnauthorized)
			return
		}

		// Check expiration time explicitly (though jwt.ParseWithClaims handles this if RegisteredClaims are used)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// SECURE: Extract UserID and Role and store them securely in the context.
		userClaims := &models.UserClaims{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		ctx := context.WithValue(r.Context(), models.UserContextKey, userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}