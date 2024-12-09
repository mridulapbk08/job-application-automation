package models

type Job struct {
	JobID         int64  `gorm:"primaryKey;autoIncrement"`
	JobSite       string `json:"job_site"`
	ScriptDetails string `json:"script_details"`
}

type Tracker struct {
	TrackerID   int64  `gorm:"primaryKey;autoIncrement"`
	JobID       int64  `json:"job_id"`
	CandidateID int    `json:"candidate_id"`
	Status      string `json:"status"`
	Output      string `json:"output"`
	Error       string `json:"error"`
	RetryCount  int    `json:"retry_count"` // Track retries
	MaxRetries  int    `json:"max_retries"` // Maximum allowed retries
	Timestamp   string `json:"timestamp"`
}
