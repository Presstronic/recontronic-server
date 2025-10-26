package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/presstronic/recontronic-server/internal/models"
)

// ProgramRepository handles program data access
type ProgramRepository struct {
	db *sql.DB
}

// NewProgramRepository creates a new program repository
func NewProgramRepository(db *sql.DB) *ProgramRepository {
	return &ProgramRepository{db: db}
}

// Create creates a new program
func (r *ProgramRepository) Create(ctx context.Context, program *models.Program) error {
	query := `
		INSERT INTO programs (name, platform, scope, scan_frequency, is_active, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4::interval, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		program.Name,
		program.Platform,
		pq.Array(program.Scope),
		program.ScanFrequency,
		program.IsActive,
		program.Metadata,
		program.CreatedAt,
		program.UpdatedAt,
	).Scan(&program.ID)

	if err != nil {
		return fmt.Errorf("failed to create program: %w", err)
	}

	return nil
}

// GetByID retrieves a program by ID
func (r *ProgramRepository) GetByID(ctx context.Context, id int) (*models.Program, error) {
	query := `
		SELECT id, name, platform, scope, scan_frequency, created_at, updated_at, last_scanned_at, is_active, metadata
		FROM programs
		WHERE id = $1
	`

	var program models.Program
	var scanFrequency string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Platform,
		pq.Array(&program.Scope),
		&scanFrequency,
		&program.CreatedAt,
		&program.UpdatedAt,
		&program.LastScannedAt,
		&program.IsActive,
		&program.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("program not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	program.ScanFrequency = scanFrequency

	return &program, nil
}

// GetByName retrieves a program by name
func (r *ProgramRepository) GetByName(ctx context.Context, name string) (*models.Program, error) {
	query := `
		SELECT id, name, platform, scope, scan_frequency, created_at, updated_at, last_scanned_at, is_active, metadata
		FROM programs
		WHERE name = $1
	`

	var program models.Program
	var scanFrequency string

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&program.ID,
		&program.Name,
		&program.Platform,
		pq.Array(&program.Scope),
		&scanFrequency,
		&program.CreatedAt,
		&program.UpdatedAt,
		&program.LastScannedAt,
		&program.IsActive,
		&program.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("program not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	program.ScanFrequency = scanFrequency

	return &program, nil
}

// List retrieves all programs with optional filtering
func (r *ProgramRepository) List(ctx context.Context, activeOnly bool) ([]models.Program, error) {
	query := `
		SELECT id, name, platform, scope, scan_frequency, created_at, updated_at, last_scanned_at, is_active, metadata
		FROM programs
	`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		var scanFrequency string

		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Platform,
			pq.Array(&program.Scope),
			&scanFrequency,
			&program.CreatedAt,
			&program.UpdatedAt,
			&program.LastScannedAt,
			&program.IsActive,
			&program.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan program: %w", err)
		}

		program.ScanFrequency = scanFrequency
		programs = append(programs, program)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating programs: %w", err)
	}

	return programs, nil
}

// Update updates a program
func (r *ProgramRepository) Update(ctx context.Context, program *models.Program) error {
	query := `
		UPDATE programs
		SET name = $1, platform = $2, scope = $3, scan_frequency = $4::interval, is_active = $5, metadata = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		program.Name,
		program.Platform,
		pq.Array(program.Scope),
		program.ScanFrequency,
		program.IsActive,
		program.Metadata,
		time.Now(),
		program.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("program not found")
	}

	return nil
}

// UpdateLastScanned updates the last_scanned_at timestamp
func (r *ProgramRepository) UpdateLastScanned(ctx context.Context, programID int) error {
	query := `
		UPDATE programs
		SET last_scanned_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), programID)
	if err != nil {
		return fmt.Errorf("failed to update last scanned: %w", err)
	}

	return nil
}

// Delete deletes a program
func (r *ProgramRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM programs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete program: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("program not found")
	}

	return nil
}

// GetProgramsDueForScan retrieves programs that need to be scanned
func (r *ProgramRepository) GetProgramsDueForScan(ctx context.Context) ([]models.Program, error) {
	query := `
		SELECT id, name, platform, scope, scan_frequency, created_at, updated_at, last_scanned_at, is_active, metadata
		FROM programs
		WHERE is_active = true
		AND (
			last_scanned_at IS NULL
			OR last_scanned_at + scan_frequency < NOW()
		)
		ORDER BY last_scanned_at ASC NULLS FIRST
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get programs due for scan: %w", err)
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		var scanFrequency string

		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Platform,
			pq.Array(&program.Scope),
			&scanFrequency,
			&program.CreatedAt,
			&program.UpdatedAt,
			&program.LastScannedAt,
			&program.IsActive,
			&program.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan program: %w", err)
		}

		program.ScanFrequency = scanFrequency
		programs = append(programs, program)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating programs: %w", err)
	}

	return programs, nil
}
