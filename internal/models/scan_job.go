package models

import (
	"time"
)

// ScanJob represents a scan execution
type ScanJob struct {
	ID           int       `json:"id"`
	ProgramID    int       `json:"program_id"`
	JobType      string    `json:"job_type"`
	Status       string    `json:"status"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ResultsCount int       `json:"results_count"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	Metadata     Metadata  `json:"metadata,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ScanJobStatus constants
const (
	ScanJobStatusPending   = "pending"
	ScanJobStatusRunning   = "running"
	ScanJobStatusCompleted = "completed"
	ScanJobStatusFailed    = "failed"
	ScanJobStatusCancelled = "cancelled"
)

// ScanJobType constants
const (
	ScanJobTypePassive = "passive"
	ScanJobTypeActive  = "active"
	ScanJobTypeDeep    = "deep"
	ScanJobTypeManual  = "manual"
)
