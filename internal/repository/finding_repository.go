package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
)

// FindingRepository handles finding (vulnerability) data access
type FindingRepository struct {
	db *sql.DB
}

// NewFindingRepository creates a new finding repository
func NewFindingRepository(db *sql.DB) *FindingRepository {
	return &FindingRepository{db: db}
}

// Create creates a new finding
func (r *FindingRepository) Create(ctx context.Context, finding *models.Finding) error {
	query := `
		INSERT INTO findings (
			program_id, anomaly_id, user_id, title, severity, status,
			vulnerability_type, cvss_score, currency, notes, poc,
			metadata, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		finding.ProgramID,
		finding.AnomalyID,
		finding.UserID,
		finding.Title,
		finding.Severity,
		finding.Status,
		finding.VulnerabilityType,
		finding.CVSSScore,
		finding.Currency,
		finding.Notes,
		finding.POC,
		finding.Metadata,
		finding.CreatedAt,
		finding.UpdatedAt,
	).Scan(&finding.ID)

	if err != nil {
		return fmt.Errorf("failed to create finding: %w", err)
	}

	return nil
}

// GetByID retrieves a finding by ID
func (r *FindingRepository) GetByID(ctx context.Context, id int) (*models.Finding, error) {
	query := `
		SELECT id, program_id, anomaly_id, user_id, title, severity, status,
			vulnerability_type, cvss_score, reported_at, resolved_at,
			bounty_amount, currency, notes, poc, metadata, created_at, updated_at
		FROM findings
		WHERE id = $1
	`

	var finding models.Finding
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&finding.ID,
		&finding.ProgramID,
		&finding.AnomalyID,
		&finding.UserID,
		&finding.Title,
		&finding.Severity,
		&finding.Status,
		&finding.VulnerabilityType,
		&finding.CVSSScore,
		&finding.ReportedAt,
		&finding.ResolvedAt,
		&finding.BountyAmount,
		&finding.Currency,
		&finding.Notes,
		&finding.POC,
		&finding.Metadata,
		&finding.CreatedAt,
		&finding.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("finding not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get finding: %w", err)
	}

	return &finding, nil
}

// ListByProgram retrieves findings for a program
func (r *FindingRepository) ListByProgram(ctx context.Context, programID int, limit int) ([]models.Finding, error) {
	query := `
		SELECT id, program_id, anomaly_id, user_id, title, severity, status,
			vulnerability_type, cvss_score, reported_at, resolved_at,
			bounty_amount, currency, notes, poc, metadata, created_at, updated_at
		FROM findings
		WHERE program_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	return r.queryFindings(ctx, query, programID, limit)
}

// ListByStatus retrieves findings by status
func (r *FindingRepository) ListByStatus(ctx context.Context, status string, limit int) ([]models.Finding, error) {
	query := `
		SELECT id, program_id, anomaly_id, user_id, title, severity, status,
			vulnerability_type, cvss_score, reported_at, resolved_at,
			bounty_amount, currency, notes, poc, metadata, created_at, updated_at
		FROM findings
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	return r.queryFindings(ctx, query, status, limit)
}

// ListBySeverity retrieves findings by severity
func (r *FindingRepository) ListBySeverity(ctx context.Context, severity string, limit int) ([]models.Finding, error) {
	query := `
		SELECT id, program_id, anomaly_id, user_id, title, severity, status,
			vulnerability_type, cvss_score, reported_at, resolved_at,
			bounty_amount, currency, notes, poc, metadata, created_at, updated_at
		FROM findings
		WHERE severity = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	return r.queryFindings(ctx, query, severity, limit)
}

// Update updates a finding
func (r *FindingRepository) Update(ctx context.Context, finding *models.Finding) error {
	query := `
		UPDATE findings
		SET title = $1, severity = $2, status = $3, vulnerability_type = $4,
			cvss_score = $5, reported_at = $6, resolved_at = $7,
			bounty_amount = $8, currency = $9, notes = $10, poc = $11,
			metadata = $12, updated_at = $13
		WHERE id = $14
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		finding.Title,
		finding.Severity,
		finding.Status,
		finding.VulnerabilityType,
		finding.CVSSScore,
		finding.ReportedAt,
		finding.ResolvedAt,
		finding.BountyAmount,
		finding.Currency,
		finding.Notes,
		finding.POC,
		finding.Metadata,
		time.Now(),
		finding.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update finding: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("finding not found")
	}

	return nil
}

// UpdateStatus updates the status of a finding
func (r *FindingRepository) UpdateStatus(ctx context.Context, findingID int, status string) error {
	query := `
		UPDATE findings
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), findingID)
	if err != nil {
		return fmt.Errorf("failed to update finding status: %w", err)
	}

	return nil
}

// MarkAsReported marks a finding as reported
func (r *FindingRepository) MarkAsReported(ctx context.Context, findingID int) error {
	query := `
		UPDATE findings
		SET status = $1, reported_at = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, models.FindingStatusSubmitted, time.Now(), time.Now(), findingID)
	if err != nil {
		return fmt.Errorf("failed to mark finding as reported: %w", err)
	}

	return nil
}

// MarkAsResolved marks a finding as resolved
func (r *FindingRepository) MarkAsResolved(ctx context.Context, findingID int) error {
	query := `
		UPDATE findings
		SET status = $1, resolved_at = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, models.FindingStatusResolved, time.Now(), time.Now(), findingID)
	if err != nil {
		return fmt.Errorf("failed to mark finding as resolved: %w", err)
	}

	return nil
}

// RecordBounty records a bounty payment for a finding
func (r *FindingRepository) RecordBounty(ctx context.Context, findingID int, amount float64, currency string) error {
	query := `
		UPDATE findings
		SET bounty_amount = $1, currency = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query, amount, currency, models.FindingStatusBountyAwarded, time.Now(), findingID)
	if err != nil {
		return fmt.Errorf("failed to record bounty: %w", err)
	}

	return nil
}

// Delete deletes a finding
func (r *FindingRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM findings WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete finding: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("finding not found")
	}

	return nil
}

// GetBountyStats retrieves bounty statistics
func (r *FindingRepository) GetBountyStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE bounty_amount IS NOT NULL) as paid_count,
			COALESCE(SUM(bounty_amount), 0) as total_earned,
			COALESCE(AVG(bounty_amount), 0) as avg_bounty,
			COALESCE(MAX(bounty_amount), 0) as max_bounty
		FROM findings
	`

	var stats = make(map[string]interface{})
	var paidCount int
	var totalEarned, avgBounty, maxBounty float64

	err := r.db.QueryRowContext(ctx, query).Scan(&paidCount, &totalEarned, &avgBounty, &maxBounty)
	if err != nil {
		return nil, fmt.Errorf("failed to get bounty stats: %w", err)
	}

	stats["paid_count"] = paidCount
	stats["total_earned"] = totalEarned
	stats["avg_bounty"] = avgBounty
	stats["max_bounty"] = maxBounty

	return stats, nil
}

// GetBountyStatsByProgram retrieves bounty statistics for a specific program
func (r *FindingRepository) GetBountyStatsByProgram(ctx context.Context, programID int) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE bounty_amount IS NOT NULL) as paid_count,
			COALESCE(SUM(bounty_amount), 0) as total_earned,
			COALESCE(AVG(bounty_amount), 0) as avg_bounty,
			COALESCE(MAX(bounty_amount), 0) as max_bounty
		FROM findings
		WHERE program_id = $1
	`

	var stats = make(map[string]interface{})
	var paidCount int
	var totalEarned, avgBounty, maxBounty float64

	err := r.db.QueryRowContext(ctx, query, programID).Scan(&paidCount, &totalEarned, &avgBounty, &maxBounty)
	if err != nil {
		return nil, fmt.Errorf("failed to get bounty stats by program: %w", err)
	}

	stats["paid_count"] = paidCount
	stats["total_earned"] = totalEarned
	stats["avg_bounty"] = avgBounty
	stats["max_bounty"] = maxBounty

	return stats, nil
}

// queryFindings is a helper function to execute queries and scan results
func (r *FindingRepository) queryFindings(ctx context.Context, query string, args ...interface{}) ([]models.Finding, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query findings: %w", err)
	}
	defer rows.Close()

	var findings []models.Finding
	for rows.Next() {
		var finding models.Finding
		err := rows.Scan(
			&finding.ID,
			&finding.ProgramID,
			&finding.AnomalyID,
			&finding.UserID,
			&finding.Title,
			&finding.Severity,
			&finding.Status,
			&finding.VulnerabilityType,
			&finding.CVSSScore,
			&finding.ReportedAt,
			&finding.ResolvedAt,
			&finding.BountyAmount,
			&finding.Currency,
			&finding.Notes,
			&finding.POC,
			&finding.Metadata,
			&finding.CreatedAt,
			&finding.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan finding: %w", err)
		}
		findings = append(findings, finding)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating findings: %w", err)
	}

	return findings, nil
}
