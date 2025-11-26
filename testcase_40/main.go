package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secure-api/handlers"
	"github.com/secure-api/services"
	"github.com/secure-api/utils"
	"github.com/spf13/viper"
)

// initConfig loads configuration settings using Viper.
// In a production environment, this would load sensitive settings like database credentials.
func initConfig() {
	vipr := viper.New()
	vipr.SetConfigName("config") // name of config file (without extension)
	vipr.SetConfigType("yaml")   // or json, properties, etc.
	vipr.AddConfigPath(".")      // path to look for the config file in
	vipr.SetDefault("server.port", "8080")
	vipr.SetDefault("server.read_timeout", 5)
	vipr.SetDefault("server.write_timeout", 10)

	if err := vipr.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; using defaults
			log.Println("Warning: Config file not found, using default settings.")
		} else {
			// Config file was found but error was encountered
			log.Fatalf("Fatal error reading config file: %s \n", err)
		}
	}
	// Set global viper instance for access if needed, though direct access is cleaner.
}

func main() {
	// 1. Configuration Setup
	initConfig()

	// 2. Gin Mode Setup (Ensure production mode for security and performance)
	if viper.GetString("environment") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Dependency Injection (Services and Utilities)
	dataService := services.NewDataService()
	pdfGenerator := utils.NewPDFGenerator()
	reportHandler := handlers.NewReportHandler(dataService, pdfGenerator)

	// 4. Router Initialization
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery()) // Proper error handling, prevents crashing on panic

	// 5. Security Middleware (Simulated Authentication)
	router.Use(func(c *gin.Context) {
		// In a real application, robust authentication (e.g., JWT validation) and
		// authorization (Least Privilege) checks would be performed here.
		// For this example, we assume the user is authenticated.
		c.Next()
	})

	// 6. Routes Setup
	v1 := router.Group("/api/v1")
	{
		reports := v1.Group("/reports")
		// The handler implements strict input validation to prevent DoS (CWE-400).
		reports.GET("/generate", reportHandler.GenerateReport)
	}

	// 7. Server Start with Secure Configuration
	port := viper.GetString("server.port")
	readTimeout := time.Duration(viper.GetInt("server.read_timeout")) * time.Second
	writeTimeout := time.Duration(viper.GetInt("server.write_timeout")) * time.Second

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    30 * time.Second, // Helps mitigate slowloris attacks
		MaxHeaderBytes: 1 << 20,          // 1MB max header size
	}

	log.Printf("Server starting securely on port %s...", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
