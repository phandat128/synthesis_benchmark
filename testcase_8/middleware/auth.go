package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"secure_api/models"
)

// Constants for custom context keys to prevent collisions.
type contextKey string

const (
	ContextUserIDKey   contextKey = "userID"
	ContextUserRoleKey contextKey = "userRole"
	// SECURITY NOTE: This secret MUST be stored securely (e.g., environment variable or KMS).
	jwtSecret          = "super_secure_and_long_secret_key_replace_me_in_prod"
)

// Authenticate is a middleware that validates the JWT token from the Authorization header.
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method to prevent algorithm confusion attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// Return generic unauthorized message to avoid leaking token details
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid. Store authenticated user data in the request context.
		ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, ContextUserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GenerateToken creates a new signed JWT for the authenticated user.
func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}