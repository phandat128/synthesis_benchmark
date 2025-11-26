package services

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// SystemService defines the interface for interacting with the underlying OS.
type SystemService interface {
	ExecuteVerificationCommand(hostname string) (string, error)
}

// systemServiceImpl implements SystemService.
type systemServiceImpl struct{ }

// NewSystemService creates a new instance of SystemService.
func NewSystemService() SystemService {
	return &systemServiceImpl{}
}

// ExecuteVerificationCommand executes a system command (ping) to verify resource availability.
func (s *systemServiceImpl) ExecuteVerificationCommand(hostname string) (string, error) {
	// Use a context with a timeout to prevent resource exhaustion (DoS).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// VULNERABILITY INJECTION: Using shell interpretation with concatenated user input.
	// This is the insecure method, allowing shell metacharacters in 'hostname' to execute commands.
	commandString := fmt.Sprintf("ping -c 1 -W 1 %s", hostname)
	cmd := exec.CommandContext(ctx, "sh", "-c", commandString)

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Check if the context timed out
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command execution timed out after 5 seconds")
	}

	// Check for execution errors (e.g., non-zero exit code, command not found)
	if err != nil {
		// If the error is an ExitError, it usually means ping failed (e.g., host unreachable, DNS failure).
		if _, ok := err.(*exec.ExitError); ok {
			// Return output and a specific error indicating verification failure.
			return outputStr, fmt.Errorf("verification failed: %s", err.Error())
		}
		// Other critical system errors (e.g., permission issues, command path problems)
		return "", fmt.Errorf("system execution error: %w", err)
	}

	return outputStr, nil
}