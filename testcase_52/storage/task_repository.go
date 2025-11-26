package storage

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"secure-scheduler/models"
)

// TaskRepository manages the persistence of tasks using an in-memory map and queue.
type TaskRepository struct {
	tasks map[string]models.Task
	queue []string // Simple FIFO queue of TaskIDs for processing
	mu    sync.Mutex
}

// NewTaskRepository initializes a new repository.
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]models.Task),
		queue: make([]string, 0),
	}
}

// SaveTask stores a new task and adds it to the processing queue.
func (r *TaskRepository) SaveTask(filename string) (models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	taskID := uuid.New().String()
	newTask := models.Task{
		TaskID:    taskID,
		Filename:  filename, // Storing the validated filename
		Scheduled: time.Now(),
		Status:    "PENDING",
	}

	r.tasks[taskID] = newTask
	r.queue = append(r.queue, taskID) // Add to the end of the queue

	log.Printf("Task scheduled: %s for file %s", taskID, filename)
	return newTask, nil
}

// GetNextTask retrieves the next pending task from the queue.
// It also updates the task status to PROCESSING.
func (r *TaskRepository) GetNextTask() (models.Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.queue) == 0 {
		return models.Task{}, false // No tasks available
	}

	taskID := r.queue[0]
	r.queue = r.queue[1:] // Dequeue

	task, ok := r.tasks[taskID]
	if !ok {
		log.Printf("Error: Task ID %s found in queue but not in map. Skipping.", taskID)
		return models.Task{}, false
	}

	// Update status
	task.Status = "PROCESSING"
	r.tasks[taskID] = task

	return task, true
}

// UpdateTaskStatus updates the status of a task.
func (r *TaskRepository) UpdateTaskStatus(taskID, status string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task, ok := r.tasks[taskID]; ok {
		task.Status = status
		r.tasks[taskID] = task
		log.Printf("Task %s status updated to %s", taskID, status)
	}
}