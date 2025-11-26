package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"secure-config-api/services"
)

// Define a strict regex for validating hostnames and IPv4 addresses.
// This is a crucial defense layer, preventing injection of shell metacharacters
// like ; | & $ ( ) ` \ and ensuring the input is well-formed network data.
const hostnameRegex = `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$|^(\d{1,3}\.){3}\d{1,3}$`

var hostValidator = regexp.MustCompile(hostnameRegex)

// ResourceHandler holds dependencies for resource operations.
type ResourceHandler struct {
	SystemService services.SystemService
}

// VerifyRequest defines the structure for the incoming JSON payload.
type VerifyRequest struct {
	TargetHost string `json:"target_host"`
}

// VerifyResourceStatus handles the POST request to verify resource availability.
//
// SECURITY MEASURES:
// 1. Input Validation: Ensures target_host matches a strict regex pattern.
// 2. Secure Error Handling: Prevents sensitive information leakage on internal errors.
func (h *ResourceHandler) VerifyResourceStatus(c echo.Context) error {
	req := new(VerifyRequest)
	// Bind input
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	targetHost := strings.TrimSpace(req.TargetHost)

	// --- Input Validation ---
	if targetHost == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Target host cannot be empty"})
	}

	// Strict validation against the allowed pattern.
	if !hostValidator.MatchString(targetHost) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target host format. Only standard hostnames or IPv4 addresses are allowed."})
	}

	// --- Service Execution ---
	output, err := h.SystemService.ExecuteVerificationCommand(targetHost)

	if err != nil {
		// Check if the error is a known failure (e.g., host unreachable, timeout)
		if strings.Contains(err.Error(), "verification failed") || strings.Contains(err.Error(), "timed out") {
			// Return 200 OK for verification failures, as the command executed successfully but the resource was unavailable.
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status":  "Verification Failed",
				"host":    targetHost,
				"details": output, // Return the ping output for diagnostic purposes
				"error":   err.Error(),
			})
		}

		// Handle critical internal system errors securely (Least Privilege/Information Leakage Prevention)
		c.Logger().Errorf("Internal system service error during verification for %s: %v", targetHost, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "An internal error occurred during command execution."})
	}

	// Success response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "Resource Available",
		"host":    targetHost,
		"details": output,
	})
}