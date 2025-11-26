package models

import (
	"time"

	"forum_api/db"
)

// Topic represents a discussion topic.
type Topic struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Save inserts a new Topic into the database.
// Uses parameterized queries to prevent SQL Injection.
func (t *Topic) Save() error {
	// Basic input validation (enforced by DB schema as well)
	if len(t.Title) == 0 || len(t.Title) > 100 {
		return ErrInvalidInput
	}

	query := `INSERT INTO topics (title) VALUES ($1) RETURNING id, created_at`
	
	// QueryRowx is used to execute the query and scan the returned ID and timestamp.
	err := db.DB.QueryRowx(query, t.Title).StructScan(t)
	if err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all topics from the database.
func GetAllTopics() ([]Topic, error) {
	topics := []Topic{}
	query := `SELECT id, title, created_at FROM topics ORDER BY created_at DESC`
	
	err := db.DB.Select(&topics, query)
	if err != nil {
		return nil, err
	}
	return topics, nil
}