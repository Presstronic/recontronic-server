package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Program represents a bug bounty program being monitored
type Program struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Platform       string    `json:"platform"`
	Scope          Scope     `json:"scope"`
	ScanFrequency  string    `json:"scan_frequency"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastScannedAt  *time.Time `json:"last_scanned_at,omitempty"`
	IsActive       bool      `json:"is_active"`
	Metadata       Metadata  `json:"metadata,omitempty"`
}

// Scope represents the in-scope domains/assets for a program
type Scope []string

// Value implements the driver.Valuer interface for database storage
func (s Scope) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	// PostgreSQL array format: {"value1","value2"}
	return s, nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *Scope) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	// PostgreSQL lib/pq returns []string for TEXT[] columns
	if arr, ok := value.([]string); ok {
		*s = arr
		return nil
	}

	// Fallback for other drivers
	if arr, ok := value.([]interface{}); ok {
		result := make([]string, len(arr))
		for i, v := range arr {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		*s = result
		return nil
	}

	return nil
}

// Metadata represents JSONB metadata field
type Metadata map[string]interface{}

// Value implements the driver.Valuer interface
func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements the sql.Scanner interface
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, m)
}
