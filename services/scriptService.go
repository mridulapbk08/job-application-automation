package services

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"job-application-automation/database"
	"job-application-automation/models"
)

// FetchJob retrieves the job record by JobID
func FetchJob(jobID int64, job *models.Job) error {
	err := database.DB.Where("job_id = ?", jobID).First(job).Error
	if err != nil {
		return fmt.Errorf("job not found: %v", err)
	}
	return nil
}

// ExecuteScriptService executes the Node.js script and updates the tracker
func ExecuteScriptService(jobID int64, candidateID int, scriptPath string) error {
	var tracker models.Tracker
	db := database.DB

	// Fetch or create tracker entry
	err := db.Where("job_id = ? AND candidate_id = ?", jobID, candidateID).First(&tracker).Error
	if err != nil {
		if err.Error() == "record not found" {
			tracker = models.Tracker{
				JobID:       jobID,
				CandidateID: candidateID,
				Status:      "Pending",
				RetryCount:  0,
				MaxRetries:  3,
				Timestamp:   time.Now().Format(time.RFC3339),
			}
			if createErr := db.Create(&tracker).Error; createErr != nil {
				return fmt.Errorf("failed to create tracker entry: %v", createErr)
			}
			log.Printf("Created new tracker entry with 'Pending' status for JobID %d, CandidateID %d", jobID, candidateID)
		} else {
			return fmt.Errorf("failed to fetch tracker entry: %v", err)
		}
	}

	if tracker.RetryCount >= tracker.MaxRetries {
		log.Printf("JobID %d, CandidateID %d has exceeded maximum retries. Skipping execution.", jobID, candidateID)
		return fmt.Errorf("maximum retry attempts reached")
	}

	tracker.Status = "Pending"
	tracker.Timestamp = time.Now().Format(time.RFC3339)
	if updateErr := db.Save(&tracker).Error; updateErr != nil {
		return fmt.Errorf("failed to update tracker to 'Pending': %v", updateErr)
	}

	cmd := exec.Command("node", scriptPath, fmt.Sprintf("%d", jobID), fmt.Sprintf("%d", candidateID))
	output, execErr := cmd.CombinedOutput()

	status := "Success"
	errorMsg := ""
	if execErr != nil {
		status = "Failure"
		errorMsg = execErr.Error()
		tracker.RetryCount++
		log.Printf("Script execution failed for JobID %d, CandidateID %d: %v", jobID, candidateID, execErr)
	} else {
		tracker.RetryCount = 0
		log.Printf("Script executed successfully for JobID %d, CandidateID %d", jobID, candidateID)
	}

	tracker.Status = status
	tracker.Output = string(output)
	tracker.Error = errorMsg
	tracker.Timestamp = time.Now().Format(time.RFC3339)

	if saveErr := db.Save(&tracker).Error; saveErr != nil {
		return fmt.Errorf("failed to update tracker with final status: %v", saveErr)
	}

	log.Printf("Tracker updated to '%s' for JobID %d, CandidateID %d", status, jobID, candidateID)
	return nil
}

// RetryFailedJobs retries jobs with "Pending" or "Failure" status
func RetryFailedJobs(scriptPath string) {
	for {
		log.Println("Scheduler: Checking for jobs to retry...")

		var failedJobs []models.Tracker
		err := database.DB.Where("status = ? OR status = ?", "Failure", "Pending").Find(&failedJobs).Error
		if err != nil {
			log.Printf("Scheduler: Failed to fetch jobs: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, tracker := range failedJobs {
			log.Printf("Scheduler: Retrying JobID %d, CandidateID %d", tracker.JobID, tracker.CandidateID)

			err := ExecuteScriptService(tracker.JobID, tracker.CandidateID, scriptPath)
			if err != nil {
				log.Printf("Scheduler: Failed to retry JobID %d, CandidateID %d: %v", tracker.JobID, tracker.CandidateID, err)
			} else {
				log.Printf("Scheduler: Successfully retried JobID %d, CandidateID %d", tracker.JobID, tracker.CandidateID)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
