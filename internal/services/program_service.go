package services

import (
	"context"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/repository"
)

// ProgramService handles program business logic
type ProgramService struct {
	programRepo *repository.ProgramRepository
}

// NewProgramService creates a new program service
func NewProgramService(programRepo *repository.ProgramRepository) *ProgramService {
	return &ProgramService{
		programRepo: programRepo,
	}
}

// CreateProgram creates a new bug bounty program
func (s *ProgramService) CreateProgram(ctx context.Context, req *models.CreateProgramRequest) (*models.Program, error) {
	// Check if program with same name already exists
	existing, err := s.programRepo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("program with name '%s' already exists", req.Name)
	}

	program := &models.Program{
		Name:          req.Name,
		Platform:      req.Platform,
		Scope:         req.Scope,
		ScanFrequency: req.ScanFrequency,
		IsActive:      true,
		Metadata:      make(models.Metadata),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.programRepo.Create(ctx, program); err != nil {
		return nil, fmt.Errorf("failed to create program: %w", err)
	}

	return program, nil
}

// GetProgram retrieves a program by ID
func (s *ProgramService) GetProgram(ctx context.Context, id int) (*models.Program, error) {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	return program, nil
}

// ListPrograms retrieves all programs
func (s *ProgramService) ListPrograms(ctx context.Context, activeOnly bool) ([]models.Program, error) {
	programs, err := s.programRepo.List(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}

	return programs, nil
}

// UpdateProgram updates a program
func (s *ProgramService) UpdateProgram(ctx context.Context, id int, req *models.UpdateProgramRequest) (*models.Program, error) {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	// Apply updates
	if req.Name != nil {
		program.Name = *req.Name
	}
	if req.Platform != nil {
		program.Platform = *req.Platform
	}
	if req.Scope != nil {
		program.Scope = req.Scope
	}
	if req.ScanFrequency != nil {
		program.ScanFrequency = *req.ScanFrequency
	}
	if req.IsActive != nil {
		program.IsActive = *req.IsActive
	}

	program.UpdatedAt = time.Now()

	if err := s.programRepo.Update(ctx, program); err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}

	return program, nil
}

// DeleteProgram deletes a program
func (s *ProgramService) DeleteProgram(ctx context.Context, id int) error {
	if err := s.programRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete program: %w", err)
	}

	return nil
}

// GetProgramsDueForScan retrieves programs that need scanning
func (s *ProgramService) GetProgramsDueForScan(ctx context.Context) ([]models.Program, error) {
	programs, err := s.programRepo.GetProgramsDueForScan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get programs due for scan: %w", err)
	}

	return programs, nil
}
