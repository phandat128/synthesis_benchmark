package middleware

import (
	"net/http"
)

// RequireRole is a middleware that checks if the authenticated user has the required role.
// This is the critical authorization layer that prevents unauthorized access to admin endpoints.
func RequireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Retrieve the role from the request context (set by Authenticate middleware)
		role, ok := r.Context().Value(ContextUserRoleKey).(string)
		if !ok {
			// If the role is missing, it implies a failure in the preceding authentication step.
			http.Error(w, "Authorization context missing", http.StatusUnauthorized)
			return
		}

		// 2. Perform the authorization check
		if role != requiredRole {
			// Use 403 Forbidden for authorization failure (user is authenticated but not authorized)
			http.Error(w, "Forbidden: Insufficient privileges", http.StatusForbidden)
			return
		}

		// 3. Authorization successful, proceed to the handler
		next.ServeHTTP(w, r)
	}
}