package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/labstack/echo/v4"
	"job-application-automation/database"
	"job-application-automation/models"
)

func executeScript(scriptPath string, jobID int64, candidateID int) error {
	log.Printf("Executing script: node %s with arguments: %d %d\n", scriptPath, jobID, candidateID)

	cmd := exec.Command("node", scriptPath, fmt.Sprintf("%d", jobID), fmt.Sprintf("%d", candidateID))
	output, err := cmd.CombinedOutput()

	status := "Unknown"
	errorMsg := ""

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			switch exitError.ExitCode() {
			case 1:
				status = "Failure"
			case 2:
				status = "Website Down"
			default:
				status = "Error"
			}
		}
		errorMsg = err.Error()
		log.Printf("Script execution failed: %v\n", err)
	} else {
		status = "Success"
	}

	log.Printf("Script execution completed. Status: %s, Output: %s, Error: %s\n", status, string(output), errorMsg)


	tracker := models.Tracker{
		JobID:       jobID,
		CandidateID: candidateID,
		Status:      status,
		Output:      string(output),
		Error:       errorMsg,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if dbErr := database.DB.Create(&tracker).Error; dbErr != nil {
		log.Printf("Failed to update tracker entry in database: %v\n", dbErr)
		return fmt.Errorf("failed to update tracker entry: %v", dbErr)
	}

	log.Printf("Tracker entry updated successfully for JobID %d, CandidateID %d\n", jobID, candidateID)
	return nil
}

func ApplyForJob(c echo.Context) error {
	log.Println("Received POST /apply request")

	type Request struct {
		JobID       int64 `json:"job_id"`
		CandidateID int   `json:"candidate_id"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	log.Printf("Payload received: %+v\n", req)

	var job models.Job
	if err := database.DB.First(&job, req.JobID).Error; err != nil {
		log.Printf("Job ID %d not found: %v\n", req.JobID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job ID not found"})
	}


	if err := executeScript(job.ScriptDetails, req.JobID, req.CandidateID); err != nil {
		log.Printf("Error executing script for JobID %d, CandidateID %d: %v\n", req.JobID, req.CandidateID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Script execution failed"})
	}

	log.Printf("Script execution and database update completed successfully for JobID %d, CandidateID %d\n", req.JobID, req.CandidateID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       req.JobID,
		"candidate_id": req.CandidateID,
		"message":      "Job processed and logged successfully",
	})
}
