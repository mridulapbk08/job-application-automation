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

func executeScript(scriptPath string) (string, error) {
	log.Printf("Executing script at path: %s\n", scriptPath)
	cmd := exec.Command("bash", "-c", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Script execution failed: %v\n", err)
		return string(output), fmt.Errorf("script execution failed: %v", err)
	}
	return string(output), nil
}

func ApplyForJob(c echo.Context) error {
	type Request struct {
		JobID         *int   `json:"job_id"`        // Pointer to distinguish between null and 0
		CandidateID   int    `json:"candidate_id"`
		JobSite       string `json:"job_site"`       // Optional for adding a job
		ScriptDetails string `json:"script_details"` // Optional for adding a job
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var job models.Job

	
	if req.JobID != nil {
		if err := database.DB.First(&job, *req.JobID).Error; err != nil {
			log.Printf("Job ID %d not found: %v\n", *req.JobID, err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Job ID not found"})
		}

	
		if job.CandidateID != req.CandidateID {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Candidate not authorized for this job"})
		}
	} else {

		if req.JobSite == "" || req.ScriptDetails == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "JobSite and ScriptDetails are required for new jobs"})
		}

		job = models.Job{
			CandidateID:   req.CandidateID,
			JobSite:       req.JobSite,
			ScriptDetails: req.ScriptDetails,
			Status:        "Pending",
		}

		if err := database.DB.Create(&job).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add job"})
		}
		log.Printf("New job created: %+v\n", job)
	}

	log.Printf("Starting execution for JobID %d on JobSite %s with CandidateID %d\n", job.JobID, job.JobSite, job.CandidateID)


	output, err := executeScript(job.ScriptDetails)
	status := "Success"
	if err != nil {
		status = "Failure"
		log.Printf("Script execution failed for JobID %d: %v\n", job.JobID, err)
	}

	tracker := models.Tracker{
		JobID:     job.JobID,
		Status:    status,
		Output:    output,
		Timestamp: time.Now().String(),
	}
	if err := database.DB.Create(&tracker).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to log tracker entry"})
	}
	log.Printf("Tracker entry created successfully for JobID %d\n", job.JobID)

	
	job.Status = status
	if err := database.DB.Save(&job).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update job status"})
	}
	log.Printf("Job status updated successfully for JobID %d\n", job.JobID)


	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       job.JobID,
		"candidate_id": job.CandidateID,
		"message":      "Job processed successfully",
		"status":       status,
		"output":       output,
	})
}
