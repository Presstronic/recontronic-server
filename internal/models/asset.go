package models

import (
	"time"
)

// Asset represents a discovered asset (subdomain, URL, IP, etc.)
type Asset struct {
	ID             int       `json:"id"`
	ProgramID      int       `json:"program_id"`
	DiscoveredAt   time.Time `json:"discovered_at"`
	AssetType      string    `json:"asset_type"`
	AssetValue     string    `json:"asset_value"`
	IsLive         bool      `json:"is_live"`
	StatusCode     *int      `json:"status_code,omitempty"`
	ContentHash    *string   `json:"content_hash,omitempty"`
	TechStack      Metadata  `json:"tech_stack,omitempty"`
	ResponseHeaders Metadata `json:"response_headers,omitempty"`
	CertInfo       Metadata  `json:"cert_info,omitempty"`
	ResponseTimeMs *int      `json:"response_time_ms,omitempty"`
	Metadata       Metadata  `json:"metadata,omitempty"`
}

// AssetType constants
const (
	AssetTypeSubdomain  = "subdomain"
	AssetTypeURL        = "url"
	AssetTypeIP         = "ip"
	AssetTypePort       = "port"
	AssetTypeEndpoint   = "endpoint"
	AssetTypeParameter  = "parameter"
)
