package handlers

import (
	"encoding/json"
	"net/http"

	"forum_api/models"
)

// CreateTopic handles POST /topics requests.
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var newTopic models.Topic
	
	// Security: Limit the request body size to prevent resource exhaustion attacks
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	if err := json.NewDecoder(r.Body).Decode(&newTopic); err != nil {
		// Handle malformed JSON or size limit exceeded
		http.Error(w, "Invalid request payload or size exceeded", http.StatusBadRequest)
		return
	}

	// Input validation is performed within the model's Save method.
	if err := newTopic.Save(); err != nil {
		if err == models.ErrInvalidInput {
			http.Error(w, "Validation failed: Title must be between 1 and 100 characters.", http.StatusBadRequest)
			return
		}
		// Generic internal server error to avoid leaking database details
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTopic)
}

// GetTopics handles GET /topics requests.
func GetTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := models.GetAllTopics()
	if err != nil {
		// Generic internal server error
		http.Error(w, "Failed to retrieve topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}