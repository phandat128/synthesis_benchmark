package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myproject/handlers"
	"github.com/myproject/routes"
	"github.com/myproject/services"
)

const defaultPort = "8080"

func main() {
	e := echo.New()

	// --- Security Middleware Setup ---

	// 1. Recover middleware: Recovers from panics anywhere in the chain, preventing server crash.
	e.Use(middleware.Recover())

	// 2. Logger middleware: Standard logging for requests.
	e.Use(middleware.Logger())

	// 3. Secure Headers middleware (Crucial for XSS/Clickjacking prevention)
	// Sets headers like X-Content-Type-Options, X-Frame-Options, HSTS, etc.
	e.Use(middleware.Secure())

	// 4. Custom Error Handler (Ensures consistent, non-leaky error responses)
	// This is vital for preventing CWE-200 (Information Exposure).
	e.HTTPErrorHandler = customHTTPErrorHandler

	// --- Dependency Injection and Initialization ---

	// Services
	userService := services.NewUserService()

	// Handlers
	userHandler := handlers.NewUserHandler(userService)

	// Routes
	routes.RegisterAPIRoutes(e, userHandler)

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	e.Logger.Fatal(e.Start(":" + port))
}

// customHTTPErrorHandler provides a centralized, secure way to handle all HTTP errors.
// It ensures that internal details are never leaked to the client, adhering to Least Privilege.
func customHTTPErrorHandler(err error, c echo.Context) {
	var he *echo.HTTPError
	
	// Attempt to cast the error to an Echo HTTPError
	if errors.As(err, &he) {
		// If it's a known HTTP error (e.g., 404, 400, or one we explicitly created in handlers),
		// use its status code and message.
	} else {
		// If it's an unknown internal error (e.g., database connection failure, unhandled panic),
		// treat it as a 500 and mask the internal details.
		c.Logger().Errorf("Unhandled internal error: %v", err)
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error. Please try again later.",
		}
	}

	// Send the response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			// Do not leak internal error messages, only the public message.
			err = c.JSON(he.Code, map[string]interface{}{
				"error":   true,
				"message": he.Message,
			})
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
