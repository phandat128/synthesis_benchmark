package main

import (
	"log"
	"os"

	"github.com/yourproject/routes"
)

func main() {
	// Load configuration (e.g., port)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize the router
	router := routes.SetupRouter()

	log.Printf("Starting server securely on :%s\n", port)

	// Run the server
	// In a production environment, consider using http.ListenAndServe with graceful shutdown.
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}