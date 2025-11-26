package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration settings.
type Config struct {
	ServerPort int
}

// LoadConfig reads configuration from environment variables.
// Uses environment variables for configuration, a standard secure practice.
func LoadConfig() (*Config, error) {
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080" // Default port
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT value: %w", err)
	}

	return &Config{
		ServerPort: port,
	},	nil
}