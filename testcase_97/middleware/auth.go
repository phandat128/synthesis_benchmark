package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
)

// ContextKey is used for context value keys.
type ContextKey string

const (
	// UserIDKey is the context key for the authenticated user's ID.
	UserIDKey ContextKey = "userID"
	// UserRoleKey is the context key for the authenticated user's role.
	UserRoleKey ContextKey = "userRole"
)

// AuthMiddleware validates the JWT token and extracts claims into the request context.
// This handles Authentication (Who are you?).
func AuthMiddleware(tokenAuth *jwtauth.v5.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Attempt to verify and decode the token from the request context (if jwtauth.Verifier was used previously)
			token, claims, err := jwtauth.v5.FromContext(r.Context())

			if err != nil || token == nil || jwtauth.v5.VerifyErr(err) != nil {
				// If context verification failed, manually check the Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
					http.Error(w, "Unauthorized: Missing or invalid token format", http.StatusUnauthorized)
					return
				}

				// Decode the token string
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				token, err = tokenAuth.Decode(tokenString)
				if err != nil || token == nil {
					http.Error(w, "Unauthorized: Invalid token signature or claims", http.StatusUnauthorized)
					return
				}
				claims = token.PrivateClaims()
			}

			// Extract required claims securely, ensuring type assertion succeeds
			userIDFloat, okID := claims["user_id"].(float64)
			userRole, okRole := claims["role"].(string)

			if !okID || !okRole {
				http.Error(w, "Unauthorized: Malformed token claims (missing user_id or role)", http.StatusUnauthorized)
				return
			}

			// Store claims in the context for downstream handlers
			ctx := context.WithValue(r.Context(), UserIDKey, uint(userIDFloat))
			ctx = context.WithValue(ctx, UserRoleKey, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole is the Authorization middleware.
// It checks the role stored in the context against the required role(s).
// This proactively fixes CWE-862 (Missing Authorization) by enforcing RBAC.
func RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Retrieve the role from the context (set by AuthMiddleware)
			roleVal := r.Context().Value(UserRoleKey)
			if roleVal == nil {
				// If role is missing, it implies AuthMiddleware failed or wasn't run.
				http.Error(w, "Forbidden: Authentication context missing", http.StatusForbidden)
				return
			}

			userRole, ok := roleVal.(string)
			if !ok {
				http.Error(w, "Internal Server Error: Invalid role type in context", http.StatusInternalServerError)
				return
			}

			// 2. Check if the user's role matches any required role
			isAuthorized := false
			for _, required := range requiredRoles {
				if userRole == required {
					isAuthorized = true
					break
				}
			}

			if !isAuthorized {
				// Least Privilege Principle enforced.
				http.Error(w, "Forbidden: Insufficient privileges", http.StatusForbidden)
				return
			}

			// 3. Authorization successful, proceed to the handler
			next.ServeHTTP(w, r)
		})
	}
}