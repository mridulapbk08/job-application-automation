package controllers

import (
	"log"
	"net/http"

	"job-application-automation/models"
	"job-application-automation/services"

	"github.com/labstack/echo/v4"
)

// ApplyForJob handles the `/apply` endpoint
func ApplyForJob(c echo.Context) error {
	log.Println("Received POST /apply request")

	// Define a request structure
	type Request struct {
		JobID       int64 `json:"job_id"`
		CandidateID int   `json:"candidate_id"`
	}

	var req Request

	// Parse incoming request body
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	log.Printf("Payload received: %+v", req)

	// Fetch the job details from the database
	var job models.Job
	if err := services.FetchJob(req.JobID, &job); err != nil {
		log.Printf("Failed to fetch job details for JobID %d: %v", req.JobID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job ID not found"})
	}

	// Execute the Node.js script through the service layer
	if err := services.ExecuteScriptService(req.JobID, req.CandidateID, job.ScriptDetails); err != nil {
		log.Printf("Error executing script for JobID %d, CandidateID %d: %v", req.JobID, req.CandidateID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Script execution failed"})
	}

	log.Printf("Script execution and tracker update completed successfully for JobID %d, CandidateID %d", req.JobID, req.CandidateID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       req.JobID,
		"candidate_id": req.CandidateID,
		"message":      "Job processed and logged successfully",
	})
}
