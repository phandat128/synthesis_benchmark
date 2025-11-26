package models

import (
	"time"
)

// Task represents a scheduled file processing job.
type Task struct {
	TaskID    string    `json:"task_id"`
	Filename  string    `json:"filename"` // User-provided input (must be validated)
	Scheduled time.Time `json:"scheduled_at"`
	Status    string    `json:"status"` // e.g., "PENDING", "PROCESSING", "COMPLETED", "FAILED"
}

// TaskRequest defines the structure for incoming API requests.
type TaskRequest struct {
	Filename string `json:"filename"`
}