package models

import (
	"errors"
	"time"

	"forum_api/db"
)

// Define a custom error for invalid input
var ErrInvalidInput = errors.New("invalid input: title or content constraints violated")

// Comment represents a user comment on a topic.
type Comment struct {
	ID        int       `json:"id" db:"id"`
	TopicID   int       `json:"topic_id" db:"topic_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Save inserts a new Comment into the database.
// This function is the critical defense point against Stored XSS (via input validation)
// and SQL Injection (via parameterized query).
func (c *Comment) Save() error {
	// --- SECURITY DEFENSE: Input Validation --- 
	// 1. Check TopicID validity
	if c.TopicID <= 0 {
		return errors.New("invalid topic ID")
	}
	
	// VULNERABILITY INJECTION: Developer removed application-level content validation
	// (length checks, sanitization) to save time, relying only on database constraints.
	// This allows malicious HTML/scripts to be stored if they fit the length limit.

	// Note: Since this is a minimalist forum, we assume plain text. 
	// If rich text were allowed, mandatory HTML sanitization (e.g., using bluemonday) 
	// would be required here to strip dangerous tags like <script>.

	query := `INSERT INTO comments (topic_id, content) VALUES ($1, $2) RETURNING id, created_at`
	
	// Use parameterized query ($1, $2) to prevent SQL Injection.
	err := db.DB.QueryRowx(query, c.TopicID, c.Content).StructScan(c)
	if err != nil {
		return err
	}
	return nil
}