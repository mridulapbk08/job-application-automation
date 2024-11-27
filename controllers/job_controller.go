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

func executeScript(scriptPath string, jobSite string) (string, string) {
	log.Printf("Executing script at path: %s for job site: %s\n", scriptPath, jobSite)

	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s %s", scriptPath, jobSite))
	output, err := cmd.CombinedOutput()

	status := "Unknown"

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
		log.Printf("Script execution failed: %v\n", err)
		log.Printf("Script output: %s\n", string(output))
		return string(output), status
	}

	status = "Success"
	log.Printf("Script executed successfully. Output:\n%s\n", string(output))
	return string(output), status
}

func ApplyForJob(c echo.Context) error {
	type Request struct {
		JobID       int64 `json:"job_id"`
		CandidateID int   `json:"candidate_id"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	
	var job models.Job
	if err := database.DB.First(&job, req.JobID).Error; err != nil {
		log.Printf("Job ID %d not found: %v\n", req.JobID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job ID not found"})
	}

	
	log.Printf("Starting execution for JobID %d by CandidateID %d\n", req.JobID, req.CandidateID)
	output, status := executeScript(job.ScriptDetails, job.JobSite)

	
	tracker := models.Tracker{
		JobID:       req.JobID,
		CandidateID: req.CandidateID,
		Status:      status,
		Output:      output,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	if err := database.DB.Create(&tracker).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to log tracker entry"})
	}
	log.Printf("Tracker entry created successfully for JobID %d by CandidateID %d\n", req.JobID, req.CandidateID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       req.JobID,
		"candidate_id": req.CandidateID,
		"message":      "Job processed successfully",
		"status":       status,
		"output":       output,
	})
}
