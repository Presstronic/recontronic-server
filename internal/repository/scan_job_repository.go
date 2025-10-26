package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
)

// ScanJobRepository handles scan job data access
type ScanJobRepository struct {
	db *sql.DB
}

// NewScanJobRepository creates a new scan job repository
func NewScanJobRepository(db *sql.DB) *ScanJobRepository {
	return &ScanJobRepository{db: db}
}

// Create creates a new scan job
func (r *ScanJobRepository) Create(ctx context.Context, job *models.ScanJob) error {
	query := `
		INSERT INTO scan_jobs (program_id, job_type, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		job.ProgramID,
		job.JobType,
		job.Status,
		job.Metadata,
		job.CreatedAt,
		job.UpdatedAt,
	).Scan(&job.ID)

	if err != nil {
		return fmt.Errorf("failed to create scan job: %w", err)
	}

	return nil
}

// GetByID retrieves a scan job by ID
func (r *ScanJobRepository) GetByID(ctx context.Context, id int) (*models.ScanJob, error) {
	query := `
		SELECT id, program_id, job_type, status, started_at, completed_at, results_count, error_message, metadata, created_at, updated_at
		FROM scan_jobs
		WHERE id = $1
	`

	var job models.ScanJob
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.ProgramID,
		&job.JobType,
		&job.Status,
		&job.StartedAt,
		&job.CompletedAt,
		&job.ResultsCount,
		&job.ErrorMessage,
		&job.Metadata,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scan job not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get scan job: %w", err)
	}

	return &job, nil
}

// ListByProgramID retrieves all scan jobs for a program
func (r *ScanJobRepository) ListByProgramID(ctx context.Context, programID int, limit int) ([]models.ScanJob, error) {
	query := `
		SELECT id, program_id, job_type, status, started_at, completed_at, results_count, error_message, metadata, created_at, updated_at
		FROM scan_jobs
		WHERE program_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, programID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list scan jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.ScanJob
	for rows.Next() {
		var job models.ScanJob
		err := rows.Scan(
			&job.ID,
			&job.ProgramID,
			&job.JobType,
			&job.Status,
			&job.StartedAt,
			&job.CompletedAt,
			&job.ResultsCount,
			&job.ErrorMessage,
			&job.Metadata,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scan jobs: %w", err)
	}

	return jobs, nil
}

// ListByStatus retrieves scan jobs by status
func (r *ScanJobRepository) ListByStatus(ctx context.Context, status string, limit int) ([]models.ScanJob, error) {
	query := `
		SELECT id, program_id, job_type, status, started_at, completed_at, results_count, error_message, metadata, created_at, updated_at
		FROM scan_jobs
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list scan jobs by status: %w", err)
	}
	defer rows.Close()

	var jobs []models.ScanJob
	for rows.Next() {
		var job models.ScanJob
		err := rows.Scan(
			&job.ID,
			&job.ProgramID,
			&job.JobType,
			&job.Status,
			&job.StartedAt,
			&job.CompletedAt,
			&job.ResultsCount,
			&job.ErrorMessage,
			&job.Metadata,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scan jobs: %w", err)
	}

	return jobs, nil
}

// UpdateStatus updates the status of a scan job
func (r *ScanJobRepository) UpdateStatus(ctx context.Context, jobID int, status string) error {
	query := `
		UPDATE scan_jobs
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to update scan job status: %w", err)
	}

	return nil
}

// MarkStarted marks a scan job as started
func (r *ScanJobRepository) MarkStarted(ctx context.Context, jobID int) error {
	query := `
		UPDATE scan_jobs
		SET status = $1, started_at = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, models.ScanJobStatusRunning, time.Now(), time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to mark scan job as started: %w", err)
	}

	return nil
}

// MarkCompleted marks a scan job as completed
func (r *ScanJobRepository) MarkCompleted(ctx context.Context, jobID int, resultsCount int) error {
	query := `
		UPDATE scan_jobs
		SET status = $1, completed_at = $2, results_count = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query, models.ScanJobStatusCompleted, time.Now(), resultsCount, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to mark scan job as completed: %w", err)
	}

	return nil
}

// MarkFailed marks a scan job as failed with an error message
func (r *ScanJobRepository) MarkFailed(ctx context.Context, jobID int, errorMessage string) error {
	query := `
		UPDATE scan_jobs
		SET status = $1, completed_at = $2, error_message = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query, models.ScanJobStatusFailed, time.Now(), errorMessage, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to mark scan job as failed: %w", err)
	}

	return nil
}

// UpdateResultsCount updates the results count for a running scan
func (r *ScanJobRepository) UpdateResultsCount(ctx context.Context, jobID int, count int) error {
	query := `
		UPDATE scan_jobs
		SET results_count = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, count, time.Now(), jobID)
	if err != nil {
		return fmt.Errorf("failed to update results count: %w", err)
	}

	return nil
}

// Delete deletes a scan job
func (r *ScanJobRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM scan_jobs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scan job: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scan job not found")
	}

	return nil
}

// GetRunningJobs retrieves all currently running scan jobs
func (r *ScanJobRepository) GetRunningJobs(ctx context.Context) ([]models.ScanJob, error) {
	return r.ListByStatus(ctx, models.ScanJobStatusRunning, 1000)
}

// GetPendingJobs retrieves all pending scan jobs
func (r *ScanJobRepository) GetPendingJobs(ctx context.Context, limit int) ([]models.ScanJob, error) {
	return r.ListByStatus(ctx, models.ScanJobStatusPending, limit)
}
