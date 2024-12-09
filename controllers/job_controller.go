package controllers

import (
	"log"
	"net/http"

	"job-application-automation/services"

	"github.com/labstack/echo/v4"
)

func ApplyForJob(c echo.Context) error {
	log.Println("Received POST /apply request")

	type Request struct {
		JobID       int64 `json:"job_id"`
		CandidateID int   `json:"candidate_id"`
	}

	var req Request

	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	log.Printf("Payload received: %+v", req)

	if err := services.EnsureJobExists(req.JobID); err != nil {
		log.Printf("Job does not exist for JobID %d: %v", req.JobID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to ensure job existence"})
	}

	if err := services.ExecuteScriptService(req.JobID, req.CandidateID, "scripts/job_script.js"); err != nil {
		log.Printf("Error executing script for JobID %d, CandidateID %d: %v", req.JobID, req.CandidateID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Script execution failed"})
	}

	log.Printf("Script execution completed for JobID %d, CandidateID %d", req.JobID, req.CandidateID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"job_id":       req.JobID,
		"candidate_id": req.CandidateID,
		"message":      "Job processed and logged successfully",
	})
}
