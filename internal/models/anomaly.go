package models

import (
	"time"
)

// Anomaly represents a detected anomaly that may indicate a vulnerability
type Anomaly struct {
	ID                    int        `json:"id"`
	DetectedAt            time.Time  `json:"detected_at"`
	ProgramID             int        `json:"program_id"`
	AssetID               *int       `json:"asset_id,omitempty"`
	ScanJobID             *int       `json:"scan_job_id,omitempty"`
	AnomalyType           string     `json:"anomaly_type"`
	Description           string     `json:"description"`
	Evidence              Metadata   `json:"evidence,omitempty"`
	BaseProbability       *float64   `json:"base_probability,omitempty"`
	PosteriorProbability  *float64   `json:"posterior_probability,omitempty"`
	PriorityScore         float64    `json:"priority_score"`
	IsReviewed            bool       `json:"is_reviewed"`
	ReviewNotes           *string    `json:"review_notes,omitempty"`
	ReviewedAt            *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy            *int64     `json:"reviewed_by,omitempty"`
	Metadata              Metadata   `json:"metadata,omitempty"`
}

// AnomalyType constants
const (
	AnomalyTypeNewSubdomain        = "new_subdomain"
	AnomalyTypeNewEndpoint         = "new_endpoint"
	AnomalyTypeStatusCodeChange    = "status_code_change"
	AnomalyTypeContentChange       = "content_change"
	AnomalyTypeTechStackChange     = "tech_stack_change"
	AnomalyTypeCertChange          = "cert_change"
	AnomalyTypeWeekendDeployment   = "weekend_deployment"
	AnomalyTypeOffHoursDeployment  = "off_hours_deployment"
	AnomalyTypeRapidChanges        = "rapid_changes"
	AnomalyTypeNewPort             = "new_port"
	AnomalyTypeConfigurationChange = "configuration_change"
)
