package models

// CreateProgramRequest represents a request to create a new program
type CreateProgramRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=255"`
	Platform      string   `json:"platform" validate:"required,oneof=hackerone bugcrowd intigriti yeswehack other"`
	Scope         []string `json:"scope" validate:"required,min=1,dive,required"`
	ScanFrequency string   `json:"scan_frequency" validate:"required"`
}

// UpdateProgramRequest represents a request to update a program
type UpdateProgramRequest struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Platform      *string  `json:"platform,omitempty" validate:"omitempty,oneof=hackerone bugcrowd intigriti yeswehack other"`
	Scope         []string `json:"scope,omitempty" validate:"omitempty,min=1,dive,required"`
	ScanFrequency *string  `json:"scan_frequency,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

// CreateScanJobRequest represents a request to create a new scan job
type CreateScanJobRequest struct {
	ProgramID int    `json:"program_id" validate:"required,min=1"`
	JobType   string `json:"job_type" validate:"required,oneof=passive active deep manual"`
}

// CreateAnomalyReviewRequest represents a request to review an anomaly
type CreateAnomalyReviewRequest struct {
	Notes string `json:"notes" validate:"required,min=1"`
}

// CreateFindingRequest represents a request to create a new finding
type CreateFindingRequest struct {
	ProgramID         int      `json:"program_id" validate:"required,min=1"`
	AnomalyID         *int     `json:"anomaly_id,omitempty"`
	Title             string   `json:"title" validate:"required,min=1,max=500"`
	Severity          string   `json:"severity" validate:"required,oneof=critical high medium low info"`
	VulnerabilityType *string  `json:"vulnerability_type,omitempty" validate:"omitempty,max=100"`
	CVSSScore         *float64 `json:"cvss_score,omitempty" validate:"omitempty,min=0,max=10"`
	Notes             *string  `json:"notes,omitempty"`
	POC               *string  `json:"poc,omitempty"`
}

// UpdateFindingRequest represents a request to update a finding
type UpdateFindingRequest struct {
	Title             *string  `json:"title,omitempty" validate:"omitempty,min=1,max=500"`
	Severity          *string  `json:"severity,omitempty" validate:"omitempty,oneof=critical high medium low info"`
	Status            *string  `json:"status,omitempty" validate:"omitempty,oneof=draft submitted triaged accepted duplicate informative not_applicable resolved bounty_awarded"`
	VulnerabilityType *string  `json:"vulnerability_type,omitempty" validate:"omitempty,max=100"`
	CVSSScore         *float64 `json:"cvss_score,omitempty" validate:"omitempty,min=0,max=10"`
	BountyAmount      *float64 `json:"bounty_amount,omitempty" validate:"omitempty,min=0"`
	Currency          *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	Notes             *string  `json:"notes,omitempty"`
	POC               *string  `json:"poc,omitempty"`
}
