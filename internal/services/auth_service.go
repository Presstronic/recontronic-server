package services

import (
	"context"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/repository"
	"github.com/presstronic/recontronic-server/pkg/auth"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	existingEmail, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	now := time.Now()
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns an API key
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	valid, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate API key
	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Create API key record
	apiKey := &models.APIKey{
		UserID:    user.ID,
		Name:      fmt.Sprintf("Login %s", time.Now().Format("2006-01-02 15:04:05")),
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err = s.userRepo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &models.LoginResponse{
		User:    *user,
		APIKey:  plainKey,
		KeyID:   apiKey.ID,
		Message: "Authentication successful",
	}, nil
}

// ValidateAPIKey validates an API key and returns the associated user
func (s *AuthService) ValidateAPIKey(ctx context.Context, plainKey string) (*models.User, error) {
	// Validate format
	if !auth.ValidateAPIKeyFormat(plainKey) {
		return nil, fmt.Errorf("invalid API key format")
	}

	// Hash the key
	keyHash := auth.HashAPIKey(plainKey)

	// Get API key from database
	apiKey, err := s.userRepo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	// Get associated user
	user, err := s.userRepo.GetUserByID(ctx, apiKey.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Update last used timestamp (async to avoid slowing down requests)
	go func() {
		ctx := context.Background()
		_ = s.userRepo.UpdateAPIKeyLastUsed(ctx, apiKey.ID)
	}()

	return user, nil
}

// CreateAPIKey creates a new API key for a user
func (s *AuthService) CreateAPIKey(ctx context.Context, userID int64, req *models.CreateAPIKeyRequest) (*models.CreateAPIKeyResponse, error) {
	// Verify user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate API key
	plainKey, keyHash, keyPrefix, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Create API key record
	apiKey := &models.APIKey{
		UserID:    userID,
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		ExpiresAt: req.ExpiresAt,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err = s.userRepo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &models.CreateAPIKeyResponse{
		APIKey:   *apiKey,
		PlainKey: plainKey,
	}, nil
}

// ListAPIKeys lists all API keys for a user
func (s *AuthService) ListAPIKeys(ctx context.Context, userID int64) (*models.ListAPIKeysResponse, error) {
	apiKeys, err := s.userRepo.ListAPIKeysByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	return &models.ListAPIKeysResponse{
		APIKeys: apiKeys,
		Total:   len(apiKeys),
	}, nil
}

// RevokeAPIKey revokes an API key
func (s *AuthService) RevokeAPIKey(ctx context.Context, userID int64, keyID int64) error {
	err := s.userRepo.RevokeAPIKey(ctx, keyID, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}
