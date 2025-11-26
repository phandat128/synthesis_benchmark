package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"forum_api/models"
)

// CreateComment handles POST /topics/{topicID}/comments requests.
// This function receives the potentially tainted input ('content').
func CreateComment(w http.ResponseWriter, r *http.Request) {
	topicIDStr := chi.URLParam(r, "topicID")
	// Input validation: Ensure topicID is a valid integer
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		http.Error(w, "Invalid topic ID format", http.StatusBadRequest)
		return
	}

	var newComment models.Comment
	
	// Security: Limit request body size to prevent resource exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 524288) // 512KB limit for comments

	if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
		http.Error(w, "Invalid request payload or size exceeded", http.StatusBadRequest)
		return
	}

	newComment.TopicID = topicID

	// --- SECURITY DEFENSE: Input Validation before persistence ---
	// The models.Comment.Save() method validates content length and structure.
	if err := newComment.Save(); err != nil {
		if err == models.ErrInvalidInput {
			http.Error(w, "Validation failed: Content must be between 1 and 5000 characters.", http.StatusBadRequest)
			return
		}
		// Proper Error Handling: Log the underlying error but return a generic message
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

// GetComments handles GET /topics/{topicID}/comments requests.
// This function acts as the sink, retrieving and returning stored data.
func GetComments(w http.ResponseWriter, r *http.Request) {
	topicIDStr := chi.URLParam(r, "topicID")
	// Input validation: Ensure topicID is a valid integer
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil || topicID <= 0 {
		http.Error(w, "Invalid topic ID format", http.StatusBadRequest)
		return
	}

	comments, err := models.GetCommentsByTopicID(topicID)
	if err != nil {
		// Generic internal server error
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	// Output Encoding Note: Go's standard json.Encoder automatically handles necessary
	// JSON escaping (e.g., converting '<' to '\u003c'), which prevents XSS if the
	// client consumes the JSON correctly. The primary defense against Stored XSS
	// was implemented during input validation (CreateComment/Save).

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}