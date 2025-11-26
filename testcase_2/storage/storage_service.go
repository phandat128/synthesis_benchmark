package storage

import (
	"os"
)

// FileExists checks if a file at the given path exists and is a regular file.
// This prevents serving directories or special files (like pipes or sockets).
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	// Check for errors other than NotExist (e.g., permission denied)
	if err != nil {
		return false
	}
	// Ensure it is not a directory
	return !info.IsDir()
}