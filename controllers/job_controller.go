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

func executeScript(scriptPath string, jobID int64, candidateID int) (string, string) {
    log.Printf("Executing script: node %s with arguments: %d %d\n", scriptPath, jobID, candidateID)

    cmd := exec.Command("node", scriptPath, fmt.Sprintf("%d", jobID), fmt.Sprintf("%d", candidateID))
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
        return string(output), status
    }

    status = "Success"
    log.Printf("Script executed successfully. Output:\n%s\n", string(output))
    return string(output), status
}

func ApplyForJob(c echo.Context) error {
    log.Println("Received POST /apply request")

    type Request struct {
        JobID       int64  `json:"job_id"`
        CandidateID int    `json:"candidate_id"`
        Status      string `json:"status"`
        Output      string `json:"output"`
    }

    var req Request
    if err := c.Bind(&req); err != nil {
        log.Printf("Failed to bind request: %v\n", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }

    log.Printf("Payload received: %+v\n", req)

    // Create the tracker entry
    tracker := models.Tracker{
        JobID:       req.JobID,
        CandidateID: req.CandidateID,
        Status:      req.Status,
        Output:      req.Output,
        Timestamp:   time.Now().Format(time.RFC3339),
    }

    if err := database.DB.Create(&tracker).Error; err != nil {
        log.Printf("Failed to insert tracker entry for JobID %d, CandidateID %d: %v\n", req.JobID, req.CandidateID, err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to log tracker entry"})
    }

    log.Printf("Tracker entry created successfully for JobID %d, CandidateID %d\n", req.JobID, req.CandidateID)

    return c.JSON(http.StatusOK, map[string]interface{}{
        "job_id":       req.JobID,
        "candidate_id": req.CandidateID,
        "message":      "Job processed successfully",
        "status":       req.Status,
        "output":       req.Output,
    })
}
