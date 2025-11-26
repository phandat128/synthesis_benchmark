package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"secure-doc-service/config"
)

// FileHandler struct holds dependencies like configuration.
type FileHandler struct {
	Config *config.AppConfig
}

// NewFileHandler creates a new handler instance.
func NewFileHandler(cfg *config.AppConfig) *FileHandler {
	return &FileHandler{
		Config: cfg,
	}
}

// DownloadFile handles requests to download documentation files.
// It implements robust Path Traversal prevention (CWE-22) by canonicalizing the path
// and strictly enforcing that the resulting path remains within the configured base directory.
func (h *FileHandler) DownloadFile(c echo.Context) error {
	requestedFile := c.QueryParam("filename")

	// Input Validation: Ensure filename is provided.
	if requestedFile == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Filename parameter is required")
	}

	// 1. Construct the full intended path by joining the base directory and user input.
	// filepath.Join handles OS-specific separators.
	fullPath := filepath.Join(h.Config.BaseDownloadDir, requestedFile)

	// VULNERABILITY INJECTED: Path Traversal (CWE-22).
	// The path is not cleaned or validated against the base directory prefix.
	// The unsanitized 'fullPath' is passed directly to the file serving function.

	// 4. Serve the file using the potentially malicious path.
	// c.File handles file existence checks and serving.
	return c.File(fullPath)
}