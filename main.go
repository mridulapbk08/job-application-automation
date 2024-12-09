package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"job-application-automation/database"
	"job-application-automation/routes"
	"job-application-automation/services"

	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("Starting Job Application Automation Service...")

	// Connect to the database
	log.Println("Connecting to the database...")
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	log.Println("Running database migrations...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Resolve script path
	scriptPath, err := filepath.Abs("scripts/job_script.js")
	if err != nil {
		log.Fatalf("Failed to resolve script path: %v", err)
	}
	log.Printf("Script path resolved to: %s", scriptPath)

	// Start retry scheduler in a separate goroutine
	log.Println("Starting retry scheduler...")
	go services.RetryFailedJobs(scriptPath)

	// Initialize Echo server
	e := echo.New()
	routes.InitRoutes(e)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := e.Shutdown(nil); err != nil {
			log.Fatalf("Failed to gracefully shut down server: %v", err)
		}
		log.Println("Server shutdown complete.")
	}()

	// Start the server
	port := ":8082"
	log.Printf("Starting server on port %s...", port)
	if err := e.Start(port); err != nil && err.Error() != "http: Server closed" {
		log.Fatalf("Error starting server: %v", err)
	}
}
