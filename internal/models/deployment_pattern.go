package models

import (
	"time"
)

// DeploymentPattern represents learned deployment behavior for a program
type DeploymentPattern struct {
	ID            int       `json:"id"`
	ProgramID     int       `json:"program_id"`
	DayOfWeek     int       `json:"day_of_week"` // 0 = Sunday, 6 = Saturday
	HourOfDay     int       `json:"hour_of_day"` // 0-23
	ChangeCount   int       `json:"change_count"`
	AvgChanges    float64   `json:"avg_changes"`
	StddevChanges float64   `json:"stddev_changes"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	CreatedAt     time.Time `json:"created_at"`
}
