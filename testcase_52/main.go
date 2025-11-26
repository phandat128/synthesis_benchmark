package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"secure-scheduler/handlers"
	"secure-scheduler/storage"
	"secure-scheduler/workers"
)

const listenAddr = ":8080"

func main() {
	// 1. Initialize Storage and Dependencies
	repo := storage.NewTaskRepository()

	// 2. Setup Handlers
	taskHandler := &handlers.TaskHandler{
		Repo: repo,
	}

	// 3. Setup Router
	router := setupRouter(taskHandler)

	// 4. Start Background Worker
	ctx, cancel := context.WithCancel(context.Background())
	processor := workers.NewProcessor(repo)
	// Start the worker goroutine responsible for executing commands safely
	go processor.StartWorker(ctx)

	// 5. Setup HTTP Server
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
		}
	}()

	// 6. Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	// Catch SIGINT (Ctrl+C) and SIGTERM (kill command)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Stop the background worker first
	cancel()

	// Allow 30 seconds for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}