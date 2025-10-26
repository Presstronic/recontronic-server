package auth

import (
	"strings"
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	plainKey, hash, prefix, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey() error = %v", err)
	}

	// Check plain key format
	if !strings.HasPrefix(plainKey, apiKeyPrefix) {
		t.Errorf("API key doesn't have correct prefix, got: %s", plainKey)
	}

	// Check hash is not empty
	if hash == "" {
		t.Error("Generated hash is empty")
	}

	// Check prefix is correct length
	expectedPrefixLen := len(apiKeyPrefix) + 8
	if len(prefix) != expectedPrefixLen {
		t.Errorf("Prefix length = %d, want %d", len(prefix), expectedPrefixLen)
	}

	// Verify the hash matches
	computedHash := HashAPIKey(plainKey)
	if hash != computedHash {
		t.Error("Generated hash doesn't match computed hash")
	}
}

func TestHashAPIKey(t *testing.T) {
	apiKey := "rct_test123456789abcdef"

	hash1 := HashAPIKey(apiKey)
	hash2 := HashAPIKey(apiKey)

	// Same key should produce same hash
	if hash1 != hash2 {
		t.Error("HashAPIKey() produced different hashes for same key")
	}

	// Hash should not be empty
	if hash1 == "" {
		t.Error("HashAPIKey() returned empty hash")
	}

	// Different keys should produce different hashes
	differentKey := "rct_different_key_here"
	hash3 := HashAPIKey(differentKey)
	if hash1 == hash3 {
		t.Error("HashAPIKey() produced same hash for different keys")
	}
}

func TestValidateAPIKeyFormat(t *testing.T) {
	// Generate a valid key
	validKey, _, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("Failed to generate API key: %v", err)
	}

	tests := []struct {
		name   string
		apiKey string
		want   bool
	}{
		{
			name:   "valid generated key",
			apiKey: validKey,
			want:   true,
		},
		{
			name:   "missing prefix",
			apiKey: "abc123456789",
			want:   false,
		},
		{
			name:   "wrong prefix",
			apiKey: "wrong_abc123456789",
			want:   false,
		},
		{
			name:   "empty string",
			apiKey: "",
			want:   false,
		},
		{
			name:   "prefix only",
			apiKey: apiKeyPrefix,
			want:   false,
		},
		{
			name:   "invalid base64",
			apiKey: apiKeyPrefix + "!!!invalid!!!",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAPIKeyFormat(tt.apiKey)
			if got != tt.want {
				t.Errorf("ValidateAPIKeyFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskAPIKey(t *testing.T) {
	plainKey, _, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("Failed to generate API key: %v", err)
	}

	masked := MaskAPIKey(plainKey)

	// Should start with prefix + 8 chars
	if !strings.HasPrefix(masked, "rct_") {
		t.Errorf("MaskAPIKey() doesn't start with prefix")
	}

	// Should contain "..."
	if !strings.Contains(masked, "...") {
		t.Errorf("MaskAPIKey() doesn't contain masking ellipsis")
	}

	// Should be shorter than original
	if len(masked) >= len(plainKey) {
		t.Errorf("MaskAPIKey() length %d should be < original %d", len(masked), len(plainKey))
	}
}

func TestAPIKeyUniqueness(t *testing.T) {
	// Generate multiple keys and ensure they're unique
	keys := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		plainKey, _, _, err := GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey() error = %v", err)
		}

		if keys[plainKey] {
			t.Error("GenerateAPIKey() produced duplicate key")
		}
		keys[plainKey] = true
	}

	if len(keys) != iterations {
		t.Errorf("Generated %d unique keys, want %d", len(keys), iterations)
	}
}
