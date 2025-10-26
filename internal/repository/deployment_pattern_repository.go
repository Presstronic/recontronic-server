package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
)

// DeploymentPatternRepository handles deployment pattern data access
type DeploymentPatternRepository struct {
	db *sql.DB
}

// NewDeploymentPatternRepository creates a new deployment pattern repository
func NewDeploymentPatternRepository(db *sql.DB) *DeploymentPatternRepository {
	return &DeploymentPatternRepository{db: db}
}

// Upsert creates or updates a deployment pattern
func (r *DeploymentPatternRepository) Upsert(ctx context.Context, pattern *models.DeploymentPattern) error {
	query := `
		INSERT INTO deployment_patterns (
			program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (program_id, day_of_week, hour_of_day)
		DO UPDATE SET
			change_count = EXCLUDED.change_count,
			avg_changes = EXCLUDED.avg_changes,
			stddev_changes = EXCLUDED.stddev_changes,
			last_updated_at = EXCLUDED.last_updated_at
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		pattern.ProgramID,
		pattern.DayOfWeek,
		pattern.HourOfDay,
		pattern.ChangeCount,
		pattern.AvgChanges,
		pattern.StddevChanges,
		pattern.LastUpdatedAt,
		pattern.CreatedAt,
	).Scan(&pattern.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert deployment pattern: %w", err)
	}

	return nil
}

// GetByProgramAndTime retrieves a specific deployment pattern
func (r *DeploymentPatternRepository) GetByProgramAndTime(ctx context.Context, programID, dayOfWeek, hourOfDay int) (*models.DeploymentPattern, error) {
	query := `
		SELECT id, program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		FROM deployment_patterns
		WHERE program_id = $1 AND day_of_week = $2 AND hour_of_day = $3
	`

	var pattern models.DeploymentPattern
	err := r.db.QueryRowContext(ctx, query, programID, dayOfWeek, hourOfDay).Scan(
		&pattern.ID,
		&pattern.ProgramID,
		&pattern.DayOfWeek,
		&pattern.HourOfDay,
		&pattern.ChangeCount,
		&pattern.AvgChanges,
		&pattern.StddevChanges,
		&pattern.LastUpdatedAt,
		&pattern.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Pattern not learned yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment pattern: %w", err)
	}

	return &pattern, nil
}

// ListByProgram retrieves all deployment patterns for a program
func (r *DeploymentPatternRepository) ListByProgram(ctx context.Context, programID int) ([]models.DeploymentPattern, error) {
	query := `
		SELECT id, program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		FROM deployment_patterns
		WHERE program_id = $1
		ORDER BY day_of_week, hour_of_day
	`

	return r.queryPatterns(ctx, query, programID)
}

// GetHighActivityPeriods retrieves periods with above-average activity
func (r *DeploymentPatternRepository) GetHighActivityPeriods(ctx context.Context, programID int, minChanges int) ([]models.DeploymentPattern, error) {
	query := `
		SELECT id, program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		FROM deployment_patterns
		WHERE program_id = $1 AND change_count >= $2
		ORDER BY change_count DESC
	`

	return r.queryPatterns(ctx, query, programID, minChanges)
}

// GetLowActivityPeriods retrieves periods with below-average activity
func (r *DeploymentPatternRepository) GetLowActivityPeriods(ctx context.Context, programID int, maxChanges int) ([]models.DeploymentPattern, error) {
	query := `
		SELECT id, program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		FROM deployment_patterns
		WHERE program_id = $1 AND change_count <= $2
		ORDER BY change_count ASC
	`

	return r.queryPatterns(ctx, query, programID, maxChanges)
}

// IsAnomalousTime checks if a given time is anomalous based on learned patterns
func (r *DeploymentPatternRepository) IsAnomalousTime(ctx context.Context, programID int, timestamp time.Time, threshold float64) (bool, error) {
	dayOfWeek := int(timestamp.Weekday())
	hourOfDay := timestamp.Hour()

	pattern, err := r.GetByProgramAndTime(ctx, programID, dayOfWeek, hourOfDay)
	if err != nil {
		return false, err
	}

	// If no pattern learned yet, assume not anomalous
	if pattern == nil {
		return false, nil
	}

	// Check if activity is significantly different from average
	// (more than threshold standard deviations away)
	if pattern.StddevChanges == 0 {
		return false, nil
	}

	// For now, just check if it's a low-activity period
	// (activity significantly below average indicates out-of-pattern deployment)
	isAnomaly := pattern.ChangeCount < int(pattern.AvgChanges-(threshold*pattern.StddevChanges))

	return isAnomaly, nil
}

// ComputeAndUpdatePatterns analyzes historical data and updates patterns
func (r *DeploymentPatternRepository) ComputeAndUpdatePatterns(ctx context.Context, programID int, lookbackDays int) error {
	// This query analyzes asset changes over the past N days
	// and computes statistics per day-of-week and hour-of-day
	query := `
		WITH change_analysis AS (
			SELECT
				EXTRACT(DOW FROM discovered_at)::int as day_of_week,
				EXTRACT(HOUR FROM discovered_at)::int as hour_of_day,
				COUNT(*) as changes
			FROM assets
			WHERE program_id = $1
			AND discovered_at >= NOW() - ($2 || ' days')::interval
			GROUP BY day_of_week, hour_of_day
		),
		statistics AS (
			SELECT
				day_of_week,
				hour_of_day,
				SUM(changes) as total_changes,
				AVG(changes)::float as avg_changes,
				STDDEV(changes)::float as stddev_changes
			FROM change_analysis
			GROUP BY day_of_week, hour_of_day
		)
		INSERT INTO deployment_patterns (
			program_id, day_of_week, hour_of_day, change_count,
			avg_changes, stddev_changes, last_updated_at, created_at
		)
		SELECT
			$1, day_of_week, hour_of_day, total_changes,
			COALESCE(avg_changes, 0), COALESCE(stddev_changes, 0),
			NOW(), NOW()
		FROM statistics
		ON CONFLICT (program_id, day_of_week, hour_of_day)
		DO UPDATE SET
			change_count = EXCLUDED.change_count,
			avg_changes = EXCLUDED.avg_changes,
			stddev_changes = EXCLUDED.stddev_changes,
			last_updated_at = EXCLUDED.last_updated_at
	`

	_, err := r.db.ExecContext(ctx, query, programID, lookbackDays)
	if err != nil {
		return fmt.Errorf("failed to compute and update patterns: %w", err)
	}

	return nil
}

// DeleteByProgram deletes all deployment patterns for a program
func (r *DeploymentPatternRepository) DeleteByProgram(ctx context.Context, programID int) error {
	query := `DELETE FROM deployment_patterns WHERE program_id = $1`

	_, err := r.db.ExecContext(ctx, query, programID)
	if err != nil {
		return fmt.Errorf("failed to delete deployment patterns: %w", err)
	}

	return nil
}

// queryPatterns is a helper function to execute queries and scan results
func (r *DeploymentPatternRepository) queryPatterns(ctx context.Context, query string, args ...interface{}) ([]models.DeploymentPattern, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query deployment patterns: %w", err)
	}
	defer rows.Close()

	var patterns []models.DeploymentPattern
	for rows.Next() {
		var pattern models.DeploymentPattern
		err := rows.Scan(
			&pattern.ID,
			&pattern.ProgramID,
			&pattern.DayOfWeek,
			&pattern.HourOfDay,
			&pattern.ChangeCount,
			&pattern.AvgChanges,
			&pattern.StddevChanges,
			&pattern.LastUpdatedAt,
			&pattern.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment pattern: %w", err)
		}
		patterns = append(patterns, pattern)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deployment patterns: %w", err)
	}

	return patterns, nil
}
