package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
)

// AnomalyRepository handles anomaly data access (TimescaleDB hypertable)
type AnomalyRepository struct {
	db *sql.DB
}

// NewAnomalyRepository creates a new anomaly repository
func NewAnomalyRepository(db *sql.DB) *AnomalyRepository {
	return &AnomalyRepository{db: db}
}

// Create creates a new anomaly
func (r *AnomalyRepository) Create(ctx context.Context, anomaly *models.Anomaly) error {
	query := `
		INSERT INTO anomalies (
			detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		anomaly.DetectedAt,
		anomaly.ProgramID,
		anomaly.AssetID,
		anomaly.ScanJobID,
		anomaly.AnomalyType,
		anomaly.Description,
		anomaly.Evidence,
		anomaly.BaseProbability,
		anomaly.PosteriorProbability,
		anomaly.PriorityScore,
		anomaly.IsReviewed,
		anomaly.Metadata,
	).Scan(&anomaly.ID)

	if err != nil {
		return fmt.Errorf("failed to create anomaly: %w", err)
	}

	return nil
}

// GetByID retrieves an anomaly by ID
func (r *AnomalyRepository) GetByID(ctx context.Context, id int) (*models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE id = $1
	`

	var anomaly models.Anomaly
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&anomaly.ID,
		&anomaly.DetectedAt,
		&anomaly.ProgramID,
		&anomaly.AssetID,
		&anomaly.ScanJobID,
		&anomaly.AnomalyType,
		&anomaly.Description,
		&anomaly.Evidence,
		&anomaly.BaseProbability,
		&anomaly.PosteriorProbability,
		&anomaly.PriorityScore,
		&anomaly.IsReviewed,
		&anomaly.ReviewNotes,
		&anomaly.ReviewedAt,
		&anomaly.ReviewedBy,
		&anomaly.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("anomaly not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get anomaly: %w", err)
	}

	return &anomaly, nil
}

// ListUnreviewedByProgram retrieves unreviewed anomalies for a program
func (r *AnomalyRepository) ListUnreviewedByProgram(ctx context.Context, programID int, minPriority float64, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE program_id = $1
		AND is_reviewed = false
		AND priority_score >= $2
		ORDER BY priority_score DESC, detected_at DESC
		LIMIT $3
	`

	return r.queryAnomalies(ctx, query, programID, minPriority, limit)
}

// ListByProgram retrieves anomalies for a program
func (r *AnomalyRepository) ListByProgram(ctx context.Context, programID int, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE program_id = $1
		ORDER BY detected_at DESC
		LIMIT $2
	`

	return r.queryAnomalies(ctx, query, programID, limit)
}

// ListByPriorityRange retrieves anomalies within a priority score range
func (r *AnomalyRepository) ListByPriorityRange(ctx context.Context, minPriority, maxPriority float64, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE priority_score >= $1 AND priority_score <= $2
		ORDER BY priority_score DESC, detected_at DESC
		LIMIT $3
	`

	return r.queryAnomalies(ctx, query, minPriority, maxPriority, limit)
}

// ListByType retrieves anomalies of a specific type
func (r *AnomalyRepository) ListByType(ctx context.Context, programID int, anomalyType string, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE program_id = $1 AND anomaly_type = $2
		ORDER BY detected_at DESC
		LIMIT $3
	`

	return r.queryAnomalies(ctx, query, programID, anomalyType, limit)
}

// ListByTimeRange retrieves anomalies within a time range
func (r *AnomalyRepository) ListByTimeRange(ctx context.Context, programID int, start, end time.Time) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE program_id = $1
		AND detected_at >= $2
		AND detected_at < $3
		ORDER BY detected_at DESC
	`

	return r.queryAnomalies(ctx, query, programID, start, end)
}

// MarkAsReviewed marks an anomaly as reviewed
func (r *AnomalyRepository) MarkAsReviewed(ctx context.Context, anomalyID int, userID int64, notes string) error {
	query := `
		UPDATE anomalies
		SET is_reviewed = true,
			reviewed_at = $1,
			reviewed_by = $2,
			review_notes = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), userID, notes, anomalyID)
	if err != nil {
		return fmt.Errorf("failed to mark anomaly as reviewed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("anomaly not found")
	}

	return nil
}

// CountUnreviewedByProgram counts unreviewed anomalies for a program
func (r *AnomalyRepository) CountUnreviewedByProgram(ctx context.Context, programID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM anomalies
		WHERE program_id = $1 AND is_reviewed = false
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, programID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unreviewed anomalies: %w", err)
	}

	return count, nil
}

// CountByPriorityRange counts anomalies in priority ranges for dashboard
func (r *AnomalyRepository) CountByPriorityRange(ctx context.Context, programID int) (map[string]int, error) {
	query := `
		SELECT
			CASE
				WHEN priority_score >= 80 THEN 'high'
				WHEN priority_score >= 50 THEN 'medium'
				ELSE 'low'
			END as priority_level,
			COUNT(*) as count
		FROM anomalies
		WHERE program_id = $1 AND is_reviewed = false
		GROUP BY priority_level
	`

	rows, err := r.db.QueryContext(ctx, query, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to count by priority range: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var level string
		var count int
		if err := rows.Scan(&level, &count); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[level] = count
	}

	return counts, nil
}

// GetTopAnomalies retrieves the highest priority anomalies across all programs
func (r *AnomalyRepository) GetTopAnomalies(ctx context.Context, limit int) ([]models.Anomaly, error) {
	query := `
		SELECT id, detected_at, program_id, asset_id, scan_job_id, anomaly_type,
			description, evidence, base_probability, posterior_probability,
			priority_score, is_reviewed, review_notes, reviewed_at, reviewed_by, metadata
		FROM anomalies
		WHERE is_reviewed = false
		ORDER BY priority_score DESC, detected_at DESC
		LIMIT $1
	`

	return r.queryAnomalies(ctx, query, limit)
}

// queryAnomalies is a helper function to execute queries and scan results
func (r *AnomalyRepository) queryAnomalies(ctx context.Context, query string, args ...interface{}) ([]models.Anomaly, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query anomalies: %w", err)
	}
	defer rows.Close()

	var anomalies []models.Anomaly
	for rows.Next() {
		var anomaly models.Anomaly
		err := rows.Scan(
			&anomaly.ID,
			&anomaly.DetectedAt,
			&anomaly.ProgramID,
			&anomaly.AssetID,
			&anomaly.ScanJobID,
			&anomaly.AnomalyType,
			&anomaly.Description,
			&anomaly.Evidence,
			&anomaly.BaseProbability,
			&anomaly.PosteriorProbability,
			&anomaly.PriorityScore,
			&anomaly.IsReviewed,
			&anomaly.ReviewNotes,
			&anomaly.ReviewedAt,
			&anomaly.ReviewedBy,
			&anomaly.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan anomaly: %w", err)
		}
		anomalies = append(anomalies, anomaly)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating anomalies: %w", err)
	}

	return anomalies, nil
}
