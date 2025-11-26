package main

import (
	"net/http"

	"secure-scheduler/handlers"
)

// setupRouter initializes the HTTP router and maps endpoints to handlers.
func setupRouter(th *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoint for scheduling tasks
	// This handler includes strict input validation to prevent injection.
	mux.HandleFunc("/api/schedule/task", th.ScheduleTask)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}