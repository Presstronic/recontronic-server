package services

import (
	"context"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/repository"
)

// ScanJobService handles scan job business logic
type ScanJobService struct {
	scanJobRepo *repository.ScanJobRepository
	programRepo *repository.ProgramRepository
}

// NewScanJobService creates a new scan job service
func NewScanJobService(scanJobRepo *repository.ScanJobRepository, programRepo *repository.ProgramRepository) *ScanJobService {
	return &ScanJobService{
		scanJobRepo: scanJobRepo,
		programRepo: programRepo,
	}
}

// CreateScanJob creates a new scan job
func (s *ScanJobService) CreateScanJob(ctx context.Context, req *models.CreateScanJobRequest) (*models.ScanJob, error) {
	// Verify program exists
	_, err := s.programRepo.GetByID(ctx, req.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("program not found")
	}

	scanJob := &models.ScanJob{
		ProgramID:    req.ProgramID,
		JobType:      req.JobType,
		Status:       models.ScanJobStatusPending,
		ResultsCount: 0,
		Metadata:     make(models.Metadata),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.scanJobRepo.Create(ctx, scanJob); err != nil {
		return nil, fmt.Errorf("failed to create scan job: %w", err)
	}

	return scanJob, nil
}

// GetScanJob retrieves a scan job by ID
func (s *ScanJobService) GetScanJob(ctx context.Context, id int) (*models.ScanJob, error) {
	scanJob, err := s.scanJobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan job: %w", err)
	}

	return scanJob, nil
}

// ListScanJobsByProgram retrieves scan jobs for a program
func (s *ScanJobService) ListScanJobsByProgram(ctx context.Context, programID int, limit int) ([]models.ScanJob, error) {
	// Verify program exists
	_, err := s.programRepo.GetByID(ctx, programID)
	if err != nil {
		return nil, fmt.Errorf("program not found")
	}

	scanJobs, err := s.scanJobRepo.ListByProgramID(ctx, programID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list scan jobs: %w", err)
	}

	return scanJobs, nil
}

// ListScanJobsByStatus retrieves scan jobs by status
func (s *ScanJobService) ListScanJobsByStatus(ctx context.Context, status string, limit int) ([]models.ScanJob, error) {
	scanJobs, err := s.scanJobRepo.ListByStatus(ctx, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list scan jobs: %w", err)
	}

	return scanJobs, nil
}

// GetRunningJobs retrieves all currently running scan jobs
func (s *ScanJobService) GetRunningJobs(ctx context.Context) ([]models.ScanJob, error) {
	scanJobs, err := s.scanJobRepo.GetRunningJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get running jobs: %w", err)
	}

	return scanJobs, nil
}

// GetPendingJobs retrieves pending scan jobs
func (s *ScanJobService) GetPendingJobs(ctx context.Context, limit int) ([]models.ScanJob, error) {
	scanJobs, err := s.scanJobRepo.GetPendingJobs(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}

	return scanJobs, nil
}
