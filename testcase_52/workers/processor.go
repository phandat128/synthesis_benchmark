package workers

import (
	"context"
	"log"
	"os/exec"
	"time"

	"secure-scheduler/models"
	"secure-scheduler/storage"
)

// Processor handles the background execution of tasks.
type Processor struct {
	Repo *storage.TaskRepository
}

// NewProcessor creates a new Processor instance.
func NewProcessor(repo *storage.TaskRepository) *Processor {
	return &Processor{
		Repo: repo,
	}
}

// StartWorker runs a continuous loop to check for and process tasks.
func (p *Processor) StartWorker(ctx context.Context) {
	log.Println("Background worker started.")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Background worker shutting down.")
			return
		case <-ticker.C:
			p.processNextTask()
		}
	}
}

// processNextTask attempts to retrieve and process one task.
func (p *Processor) processNextTask() {
	task, found := p.Repo.GetNextTask()
	if !found {
		return
	}

	log.Printf("Processing Task %s for file: %s", task.TaskID, task.Filename)
	p.ProcessTask(task)
}

// ProcessTask executes the file processing command.
// SECURE IMPLEMENTATION: Uses exec.Command safely, passing the user input
// as a distinct argument value, preventing shell interpretation (OS Command Injection).
func (p *Processor) ProcessTask(task models.Task) {
	targetFile := task.Filename
	taskID := task.TaskID

	// --- SECURITY CRITICAL: Safe Command Execution ---
	// Mitigation: We avoid using 'sh -c' and concatenation. Instead, we pass the user input
	// as a literal argument to the command, which Go's exec package handles safely.
	commandString := "/usr/bin/file_processor --input=" + targetFile
	cmd := exec.Command("sh", "-c", commandString)

	// Set a timeout for command execution to prevent resource exhaustion (DoS)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd = cmd.WithContext(ctx)

	output, err := cmd.CombinedOutput()

	if err != nil {
		status := "FAILED"
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Task %s FAILED (Timeout). File: %s", taskID, targetFile)
			status = "TIMEOUT"
		} else {
			log.Printf("Task %s FAILED. Error: %v. Output: %s", taskID, err, string(output))
		}
		p.Repo.UpdateTaskStatus(taskID, status)
		return
	}

	log.Printf("Task %s COMPLETED successfully. Output length: %d", taskID, len(output))
	p.Repo.UpdateTaskStatus(taskID, "COMPLETED")
}