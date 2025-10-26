package services

import (
	"context"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/repository"
)

// FindingService handles finding business logic
type FindingService struct {
	findingRepo *repository.FindingRepository
	programRepo *repository.ProgramRepository
}

// NewFindingService creates a new finding service
func NewFindingService(findingRepo *repository.FindingRepository, programRepo *repository.ProgramRepository) *FindingService {
	return &FindingService{
		findingRepo: findingRepo,
		programRepo: programRepo,
	}
}

// CreateFinding creates a new finding
func (s *FindingService) CreateFinding(ctx context.Context, userID int64, req *models.CreateFindingRequest) (*models.Finding, error) {
	// Verify program exists
	_, err := s.programRepo.GetByID(ctx, req.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("program not found")
	}

	finding := &models.Finding{
		ProgramID:         req.ProgramID,
		AnomalyID:         req.AnomalyID,
		UserID:            &userID,
		Title:             req.Title,
		Severity:          req.Severity,
		Status:            models.FindingStatusDraft,
		VulnerabilityType: req.VulnerabilityType,
		CVSSScore:         req.CVSSScore,
		Currency:          "USD",
		Notes:             req.Notes,
		POC:               req.POC,
		Metadata:          make(models.Metadata),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.findingRepo.Create(ctx, finding); err != nil {
		return nil, fmt.Errorf("failed to create finding: %w", err)
	}

	return finding, nil
}

// GetFinding retrieves a finding by ID
func (s *FindingService) GetFinding(ctx context.Context, id int) (*models.Finding, error) {
	finding, err := s.findingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get finding: %w", err)
	}

	return finding, nil
}

// ListFindingsByProgram retrieves findings for a program
func (s *FindingService) ListFindingsByProgram(ctx context.Context, programID int, limit int) ([]models.Finding, error) {
	findings, err := s.findingRepo.ListByProgram(ctx, programID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list findings: %w", err)
	}

	return findings, nil
}

// UpdateFinding updates a finding
func (s *FindingService) UpdateFinding(ctx context.Context, id int, req *models.UpdateFindingRequest) (*models.Finding, error) {
	finding, err := s.findingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get finding: %w", err)
	}

	// Apply updates
	if req.Title != nil {
		finding.Title = *req.Title
	}
	if req.Severity != nil {
		finding.Severity = *req.Severity
	}
	if req.Status != nil {
		finding.Status = *req.Status
	}
	if req.VulnerabilityType != nil {
		finding.VulnerabilityType = req.VulnerabilityType
	}
	if req.CVSSScore != nil {
		finding.CVSSScore = req.CVSSScore
	}
	if req.BountyAmount != nil {
		finding.BountyAmount = req.BountyAmount
	}
	if req.Currency != nil {
		finding.Currency = *req.Currency
	}
	if req.Notes != nil {
		finding.Notes = req.Notes
	}
	if req.POC != nil {
		finding.POC = req.POC
	}

	finding.UpdatedAt = time.Now()

	if err := s.findingRepo.Update(ctx, finding); err != nil {
		return nil, fmt.Errorf("failed to update finding: %w", err)
	}

	return finding, nil
}

// GetBountyStats retrieves bounty statistics
func (s *FindingService) GetBountyStats(ctx context.Context) (map[string]interface{}, error) {
	stats, err := s.findingRepo.GetBountyStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bounty stats: %w", err)
	}

	return stats, nil
}
