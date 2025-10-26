package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// API key format: rct_<32 random bytes base64 encoded>
	apiKeyPrefix = "rct_"
	apiKeyLength = 32 // bytes (256 bits)
)

// GenerateAPIKey generates a new random API key
func GenerateAPIKey() (plainKey string, hash string, prefix string, err error) {
	// Generate random bytes
	randomBytes := make([]byte, apiKeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as base64
	encodedKey := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Create the full API key with prefix
	plainKey = apiKeyPrefix + encodedKey

	// Hash the key for storage (SHA-256)
	hash = HashAPIKey(plainKey)

	// Extract prefix (first 8 chars after rct_) for display/identification
	prefix = extractPrefix(plainKey)

	return plainKey, hash, prefix, nil
}

// HashAPIKey hashes an API key for secure storage
func HashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

// ValidateAPIKeyFormat checks if an API key has the correct format
func ValidateAPIKeyFormat(apiKey string) bool {
	if !strings.HasPrefix(apiKey, apiKeyPrefix) {
		return false
	}

	// Remove prefix and check if it's valid base64
	encodedPart := strings.TrimPrefix(apiKey, apiKeyPrefix)
	decoded, err := base64.RawURLEncoding.DecodeString(encodedPart)
	if err != nil {
		return false
	}

	// Should decode to exactly apiKeyLength bytes
	return len(decoded) == apiKeyLength
}

// extractPrefix returns the first 8 characters after the prefix for display
func extractPrefix(apiKey string) string {
	if len(apiKey) < len(apiKeyPrefix)+8 {
		return apiKey
	}
	return apiKey[:len(apiKeyPrefix)+8]
}

// MaskAPIKey masks an API key for safe display (shows only prefix)
func MaskAPIKey(apiKey string) string {
	prefix := extractPrefix(apiKey)
	return prefix + "..." + apiKey[len(apiKey)-4:]
}
