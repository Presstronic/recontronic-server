package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/presstronic/recontronic-server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Recontronic Worker...")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Concurrency: %d", cfg.Worker.Concurrency)

	// TODO: Initialize database connection
	// TODO: Initialize Redis connection
	// TODO: Initialize worker pool
	// TODO: Register worker handlers (subfinder, httpx, etc.)
	// TODO: Start processing jobs

	log.Println("Worker started and waiting for jobs...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	// TODO: Gracefully stop worker pool
	// TODO: Wait for in-flight jobs to complete
	// TODO: Close database connections

	log.Println("Worker stopped gracefully")
}
