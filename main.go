package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"job-application-automation/database"
	"job-application-automation/models"
	"job-application-automation/routes"
	"job-application-automation/services"

	"github.com/labstack/echo/v4"
)

const RetryLimit = 3
const RetryDelay = time.Second * 5

func ScheduleAndRetry(jobID int64, candidateID int64, scriptPath string) error {
	retryCount := 0

	for {
		log.Printf("Executing script for JobID: %d, CandidateID: %d (Attempt %d)", jobID, candidateID, retryCount+1)

		err := services.ExecuteScriptService(jobID, int(candidateID), scriptPath)
		if err == nil {
			log.Printf("Script executed successfully for JobID: %d, CandidateID: %d", jobID, candidateID)
			return nil
		}

		log.Printf("Script execution failed for JobID: %d, CandidateID: %d: %v", jobID, candidateID, err)
		retryCount++
		if retryCount >= RetryLimit {
			log.Printf("Exceeded retry limit for JobID: %d, CandidateID: %d", jobID, candidateID)
			return fmt.Errorf("script failed after %d retries: %v", RetryLimit, err)
		}

		time.Sleep(RetryDelay)
	}
}

func startScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	quit := make(chan struct{})

	go func() {
		<-quit
		log.Println("Scheduler shutting down...")
	}()

	for {
		select {
		case <-ticker.C:

			var pendingTrackers []models.Tracker
			if err := database.DB.Where("status = ?", "Pending").Find(&pendingTrackers).Error; err != nil {
				log.Printf("Failed to fetch pending jobs: %v", err)
				continue
			}

			for _, tracker := range pendingTrackers {
				go func(tr models.Tracker) {
					log.Printf("Triggering scheduled execution for JobID: %d, CandidateID: %d", tr.JobID, tr.CandidateID)
					err := ScheduleAndRetry(tr.JobID, tr.CandidateID, "scripts/job_script.js")
					if err != nil {
						log.Printf("Scheduled execution failed for JobID: %d, CandidateID: %d: %v", tr.JobID, tr.CandidateID, err)
					}
				}(tracker)
			}
		case <-quit:
			log.Println("Exiting scheduler...")
			return
		}
	}
}

func main() {
	fmt.Println("Connecting to the database...")
	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	fmt.Println("Running migrations...")
	if err := database.DB.AutoMigrate(&models.Job{}, &models.Tracker{}); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go startScheduler()

	e := echo.New()
	routes.InitRoutes(e)

	fmt.Println("Starting server on port 8082...")
	go func() {
		if err := e.Start(":8082"); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	fmt.Println("Shutting down application gracefully...")
}
