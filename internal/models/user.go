package models

import "time"

// User represents a user account in the system
type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose in JSON
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	Name       string     `json:"name" db:"name"`             // e.g., "My Laptop", "CI/CD Pipeline"
	KeyHash    string     `json:"-" db:"key_hash"`            // Hashed API key (never expose)
	KeyPrefix  string     `json:"key_prefix" db:"key_prefix"` // First 8 chars for identification
	LastUsedAt *time.Time `json:"last_used_at" db:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// CreateUserRequest is the request payload for user registration
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest is the request payload for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse contains the API key issued after successful login
type LoginResponse struct {
	User    User   `json:"user"`
	APIKey  string `json:"api_key"` // Plain-text API key (shown only once)
	KeyID   int64  `json:"key_id"`
	Message string `json:"message"`
}

// CreateAPIKeyRequest is for generating additional API keys
type CreateAPIKeyRequest struct {
	Name      string     `json:"name" validate:"required,min=1,max=100"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKeyResponse contains the newly created API key
type CreateAPIKeyResponse struct {
	APIKey   APIKey `json:"api_key"`
	PlainKey string `json:"plain_key"` // Plain-text key (shown only once)
}

// ListAPIKeysResponse contains user's API keys
type ListAPIKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
	Total   int      `json:"total"`
}
