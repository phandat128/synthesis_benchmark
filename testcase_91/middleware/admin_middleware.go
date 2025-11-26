package middleware

import (
	"net/http"

	"github.com/your_project/models"
)

// AdminMiddleware checks if the authenticated user has the 'admin' role.
// This is the crucial defense against CWE-862 (Missing Authorization).
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := models.GetClaimsFromContext(r.Context())
		if !ok {
			// This should only happen if AuthenticateJWT failed to run or populate context
			http.Error(w, "Authentication context missing", http.StatusUnauthorized)
			return
		}

		// SECURE: Explicitly check the role before allowing access to the administrative endpoint.
		if claims.Role != models.RoleAdmin {
			// CWE-862 Mitigation: Deny access if not an admin.
			http.Error(w, "Forbidden: Administrative access required", http.StatusForbidden)
			return
		}

		// If the user is an admin, proceed to the handler.
		next.ServeHTTP(w, r)
	})
}