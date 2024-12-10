package services

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"job-application-automation/database"
	"job-application-automation/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EnsureJobExists(jobID int64) error {
	var job models.Job
	err := database.DB.First(&job, "job_id = ?", jobID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			newJob := models.Job{
				JobID:         jobID,
				JobSite:       "default-site",
				ScriptDetails: "default-script.js",
			}
			if createErr := database.DB.Create(&newJob).Error; createErr != nil {
				log.Printf("Failed to create job entry for job_id %d: %v", jobID, createErr)
				return fmt.Errorf("failed to create job record for job_id %d: %v", jobID, createErr)
			}
			log.Printf("Successfully created job record for job_id %d", jobID)
		} else {
			log.Printf("Error fetching job record for job_id %d: %v", jobID, err)
			return fmt.Errorf("error fetching job record for job_id %d: %v", jobID, err)
		}
	}
	return nil
}

func ExecuteScriptService(jobID int64, candidateID int, scriptPath string) error {
	if err := EnsureJobExists(jobID); err != nil {
		log.Printf("Error ensuring job exists for job_id %d: %v", jobID, err)
		return err
	}

	tracker := models.Tracker{
		JobID:       jobID,
		CandidateID: int64(candidateID),
		Status:      "Pending",
		RetryCount:  0,
		MaxRetries:  3,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job_id"}, {Name: "candidate_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "retry_count", "timestamp"}),
	}).Create(&tracker).Error; err != nil {
		log.Printf("Failed to initialize tracker: %v", err)
		return err
	}

	log.Printf("Initialized tracker with 'Pending' status for JobID %d, CandidateID %d", jobID, candidateID)

	absScriptPath, _ := filepath.Abs(scriptPath)
	cmd := exec.Command("node", absScriptPath, fmt.Sprintf("%d", jobID), fmt.Sprintf("%d", candidateID))
	output, err := cmd.CombinedOutput()

	status := "Success"
	errorMsg := ""
	if err != nil {
		status = "Failure"
		errorMsg = fmt.Sprintf("Script execution failed: %v", err)
		log.Printf("Script execution failed for JobID %d, CandidateID %d: %v", jobID, candidateID, err)
	}

	tracker.Status = status
	tracker.Output = string(output)
	tracker.Error = errorMsg
	tracker.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	if updateErr := database.DB.Save(&tracker).Error; updateErr != nil {
		log.Printf("Failed to update tracker: %v", updateErr)
		return updateErr
	}

	log.Printf("Updated tracker with status '%s' for JobID %d, CandidateID %d", status, jobID, candidateID)
	return nil
}
