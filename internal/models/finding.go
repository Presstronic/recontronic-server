package models

import (
	"time"
)

// Finding represents a discovered vulnerability and bug bounty submission
type Finding struct {
	ID                int        `json:"id"`
	ProgramID         int        `json:"program_id"`
	AnomalyID         *int       `json:"anomaly_id,omitempty"`
	UserID            *int64     `json:"user_id,omitempty"`
	Title             string     `json:"title"`
	Severity          string     `json:"severity"`
	Status            string     `json:"status"`
	VulnerabilityType *string    `json:"vulnerability_type,omitempty"`
	CVSSScore         *float64   `json:"cvss_score,omitempty"`
	ReportedAt        *time.Time `json:"reported_at,omitempty"`
	ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
	BountyAmount      *float64   `json:"bounty_amount,omitempty"`
	Currency          string     `json:"currency"`
	Notes             *string    `json:"notes,omitempty"`
	POC               *string    `json:"poc,omitempty"`
	Metadata          Metadata   `json:"metadata,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// FindingSeverity constants
const (
	FindingSeverityCritical = "critical"
	FindingSeverityHigh     = "high"
	FindingSeverityMedium   = "medium"
	FindingSeverityLow      = "low"
	FindingSeverityInfo     = "info"
)

// FindingStatus constants
const (
	FindingStatusDraft          = "draft"
	FindingStatusSubmitted      = "submitted"
	FindingStatusTriaged        = "triaged"
	FindingStatusAccepted       = "accepted"
	FindingStatusDuplicate      = "duplicate"
	FindingStatusInformative    = "informative"
	FindingStatusNotApplicable  = "not_applicable"
	FindingStatusResolved       = "resolved"
	FindingStatusBountyAwarded  = "bounty_awarded"
)
