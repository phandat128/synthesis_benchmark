package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger is a standard middleware that logs incoming HTTP requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// Log method, path, remote address, and duration.
		log.Printf("[%s] %s %s took %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}