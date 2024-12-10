package models

type Job struct {
	JobID         int64  `gorm:"primaryKey;autoIncrement" json:"job_id"`
	JobSite       string `gorm:"type:varchar(255)" json:"job_site"`
	ScriptDetails string `gorm:"type:varchar(255)" json:"script_details"`
}

type Tracker struct {
	TrackerID   int64  `gorm:"primaryKey;autoIncrement" json:"tracker_id"`
	JobID       int64  `gorm:"not null;index:unique_job_candidate,unique" json:"job_id"`
	CandidateID int64  `gorm:"not null;index:unique_job_candidate,unique" json:"candidate_id"`
	Status      string `gorm:"type:varchar(50)" json:"status"`
	Output      string `gorm:"type:text" json:"output"`
	Error       string `gorm:"type:text" json:"error"`
	RetryCount  int64  `gorm:"default:0" json:"retry_count"`
	MaxRetries  int64  `gorm:"default:3" json:"max_retries"`
	Timestamp   string `gorm:"type:datetime" json:"timestamp"`
}
