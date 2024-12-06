package services

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"job-application-automation/database"
	"job-application-automation/models"
)

// FetchJob retrieves a job record from the database using the provided jobID.
func FetchJob(jobID int64, job *models.Job) error {
	log.Printf("Fetching job details for JobID %d", jobID)

	// Query the database to find the job
	err := database.DB.Where("job_id = ?", jobID).First(job).Error
	if err != nil {
		log.Printf("Job with ID %d not found: %v", jobID, err)
		return fmt.Errorf("job not found: %v", err)
	}

	log.Printf("Job details fetched successfully for JobID %d", jobID)
	return nil
}

// ExecuteScriptService handles script execution and tracker updates.
func ExecuteScriptService(jobID int64, candidateID int, scriptPath string) error {
	// Check if a tracker entry already exists
	var tracker models.Tracker
	db := database.DB

	err := db.Where("job_id = ? AND candidate_id = ?", jobID, candidateID).First(&tracker).Error
	if err != nil {
		if err.Error() == "record not found" {
			// If not found, create a new tracker entry with "Pending" status
			log.Printf("Creating new tracker entry for JobID %d, CandidateID %d with status 'Pending'\n", jobID, candidateID)
			tracker = models.Tracker{
				JobID:       jobID,
				CandidateID: candidateID,
				Status:      "Pending",
				Timestamp:   time.Now().Format(time.RFC3339),
			}
			if createErr := db.Create(&tracker).Error; createErr != nil {
				log.Printf("Failed to create tracker entry: %v\n", createErr)
				return fmt.Errorf("failed to create tracker entry: %v", createErr)
			}
		} else {
			log.Printf("Failed to fetch tracker entry: %v\n", err)
			return fmt.Errorf("failed to fetch tracker entry: %v", err)
		}
	} else {
		log.Printf("Tracker entry found for JobID %d, CandidateID %d. Updating status to 'Pending'\n", jobID, candidateID)
		tracker.Status = "Pending"
		tracker.Timestamp = time.Now().Format(time.RFC3339)
		if updateErr := db.Save(&tracker).Error; updateErr != nil {
			log.Printf("Failed to update tracker entry: %v\n", updateErr)
			return fmt.Errorf("failed to update tracker entry: %v", updateErr)
		}
	}

	// Commit the "Pending" status to the database
	log.Printf("Committed 'Pending' status to the database for JobID %d, CandidateID %d\n", jobID, candidateID)

	// Execute the Node.js script
	log.Printf("Executing script: node %s with arguments: %d %d\n", scriptPath, jobID, candidateID)
	cmd := exec.Command("node", scriptPath, fmt.Sprintf("%d", jobID), fmt.Sprintf("%d", candidateID))
	output, execErr := cmd.CombinedOutput()

	// Determine script status
	status := "Unknown"
	errorMsg := ""
	if execErr != nil {
		status = "Failure"
		errorMsg = execErr.Error()
		log.Printf("Script execution failed: %v\n", execErr)
	} else {
		status = "Success"
	}

	// Update tracker with the final status
	log.Printf("Updating tracker entry with status '%s' for JobID %d, CandidateID %d\n", status, jobID, candidateID)
	tracker.Status = status
	tracker.Output = string(output)
	tracker.Error = errorMsg
	tracker.Timestamp = time.Now().Format(time.RFC3339)

	if updateErr := db.Save(&tracker).Error; updateErr != nil {
		log.Printf("Failed to update tracker entry: %v\n", updateErr)
		return fmt.Errorf("failed to update tracker entry: %v", updateErr)
	}

	log.Printf("Tracker entry updated successfully for JobID %d, CandidateID %d\n", jobID, candidateID)
	return nil
}
