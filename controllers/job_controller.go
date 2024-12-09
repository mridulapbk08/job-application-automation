package controllers

import (
	"fmt"
	"log"
	"net/http"

	"job-application-automation/models"
	"job-application-automation/services"

	"github.com/labstack/echo/v4"
)

// ApplyForJob handles the /apply endpoint
func ApplyForJob(c echo.Context) error {
	log.Println("Received POST /apply request")

	type Request struct {
		JobID       int64 `json:"job_id"`
		CandidateID int   `json:"candidate_id"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request payload: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	log.Printf("Payload received: %+v", req)

	var job models.Job
	if err := services.FetchJob(req.JobID, &job); err != nil {
		log.Printf("Failed to fetch job details for JobID %d: %v", req.JobID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job ID not found"})
	}

	err := services.ExecuteScriptService(req.JobID, req.CandidateID, job.ScriptDetails)
	if err != nil {
		log.Printf("Error executing script for JobID %d, CandidateID %d: %v", req.JobID, req.CandidateID, err)

		if err.Error() == "maximum retry attempts reached" {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error":        "Maximum retry attempts reached. Please contact support.",
				"job_id":       fmt.Sprintf("%d", req.JobID),
				"candidate_id": fmt.Sprintf("%d", req.CandidateID),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Script execution failed"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       req.JobID,
		"candidate_id": req.CandidateID,
		"message":      "Job processed successfully and tracker updated.",
	})
}
