package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-server/internal/models"
)

// AssetRepository handles asset data access (TimescaleDB hypertable)
type AssetRepository struct {
	db *sql.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

// Create creates a new asset
func (r *AssetRepository) Create(ctx context.Context, asset *models.Asset) error {
	query := `
		INSERT INTO assets (
			program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		asset.ProgramID,
		asset.DiscoveredAt,
		asset.AssetType,
		asset.AssetValue,
		asset.IsLive,
		asset.StatusCode,
		asset.ContentHash,
		asset.TechStack,
		asset.ResponseHeaders,
		asset.CertInfo,
		asset.ResponseTimeMs,
		asset.Metadata,
	).Scan(&asset.ID)

	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	return nil
}

// BulkCreate creates multiple assets in a single transaction
func (r *AssetRepository) BulkCreate(ctx context.Context, assets []models.Asset) error {
	if len(assets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO assets (
			program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, asset := range assets {
		_, err := stmt.ExecContext(
			ctx,
			asset.ProgramID,
			asset.DiscoveredAt,
			asset.AssetType,
			asset.AssetValue,
			asset.IsLive,
			asset.StatusCode,
			asset.ContentHash,
			asset.TechStack,
			asset.ResponseHeaders,
			asset.CertInfo,
			asset.ResponseTimeMs,
			asset.Metadata,
		)
		if err != nil {
			return fmt.Errorf("failed to insert asset: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLatestByProgramAndType retrieves the latest assets for a program and type
func (r *AssetRepository) GetLatestByProgramAndType(ctx context.Context, programID int, assetType string, limit int) ([]models.Asset, error) {
	query := `
		SELECT DISTINCT ON (asset_value)
			id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1 AND asset_type = $2
		ORDER BY asset_value, discovered_at DESC
		LIMIT $3
	`

	return r.queryAssets(ctx, query, programID, assetType, limit)
}

// GetByProgramInTimeRange retrieves assets discovered within a time range
func (r *AssetRepository) GetByProgramInTimeRange(ctx context.Context, programID int, start, end time.Time) ([]models.Asset, error) {
	query := `
		SELECT id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1
		AND discovered_at >= $2
		AND discovered_at < $3
		ORDER BY discovered_at DESC
	`

	return r.queryAssets(ctx, query, programID, start, end)
}

// GetLatestUniqueAssets retrieves the most recent version of each unique asset
func (r *AssetRepository) GetLatestUniqueAssets(ctx context.Context, programID int) ([]models.Asset, error) {
	query := `
		SELECT DISTINCT ON (asset_type, asset_value)
			id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1
		ORDER BY asset_type, asset_value, discovered_at DESC
	`

	return r.queryAssets(ctx, query, programID)
}

// FindAssetByValue finds an asset by its value for comparison
func (r *AssetRepository) FindAssetByValue(ctx context.Context, programID int, assetType, assetValue string, since time.Time) (*models.Asset, error) {
	query := `
		SELECT id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1
		AND asset_type = $2
		AND asset_value = $3
		AND discovered_at >= $4
		ORDER BY discovered_at DESC
		LIMIT 1
	`

	var asset models.Asset
	err := r.db.QueryRowContext(ctx, query, programID, assetType, assetValue, since).Scan(
		&asset.ID,
		&asset.ProgramID,
		&asset.DiscoveredAt,
		&asset.AssetType,
		&asset.AssetValue,
		&asset.IsLive,
		&asset.StatusCode,
		&asset.ContentHash,
		&asset.TechStack,
		&asset.ResponseHeaders,
		&asset.CertInfo,
		&asset.ResponseTimeMs,
		&asset.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found is not an error for this use case
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	return &asset, nil
}

// GetLiveAssets retrieves all live assets for a program
func (r *AssetRepository) GetLiveAssets(ctx context.Context, programID int) ([]models.Asset, error) {
	query := `
		SELECT DISTINCT ON (asset_type, asset_value)
			id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1 AND is_live = true
		ORDER BY asset_type, asset_value, discovered_at DESC
	`

	return r.queryAssets(ctx, query, programID)
}

// GetAssetHistory retrieves the history of a specific asset value
func (r *AssetRepository) GetAssetHistory(ctx context.Context, programID int, assetType, assetValue string, limit int) ([]models.Asset, error) {
	query := `
		SELECT id, program_id, discovered_at, asset_type, asset_value, is_live,
			status_code, content_hash, tech_stack, response_headers, cert_info,
			response_time_ms, metadata
		FROM assets
		WHERE program_id = $1 AND asset_type = $2 AND asset_value = $3
		ORDER BY discovered_at DESC
		LIMIT $4
	`

	return r.queryAssets(ctx, query, programID, assetType, assetValue, limit)
}

// CountAssetsByProgram counts assets for a program
func (r *AssetRepository) CountAssetsByProgram(ctx context.Context, programID int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT asset_value)
		FROM assets
		WHERE program_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, programID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count assets: %w", err)
	}

	return count, nil
}

// CountAssetsByType counts assets by type for a program
func (r *AssetRepository) CountAssetsByType(ctx context.Context, programID int) (map[string]int, error) {
	query := `
		SELECT asset_type, COUNT(DISTINCT asset_value) as count
		FROM assets
		WHERE program_id = $1
		GROUP BY asset_type
	`

	rows, err := r.db.QueryContext(ctx, query, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to count assets by type: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var assetType string
		var count int
		if err := rows.Scan(&assetType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[assetType] = count
	}

	return counts, nil
}

// queryAssets is a helper function to execute queries and scan results
func (r *AssetRepository) queryAssets(ctx context.Context, query string, args ...interface{}) ([]models.Asset, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID,
			&asset.ProgramID,
			&asset.DiscoveredAt,
			&asset.AssetType,
			&asset.AssetValue,
			&asset.IsLive,
			&asset.StatusCode,
			&asset.ContentHash,
			&asset.TechStack,
			&asset.ResponseHeaders,
			&asset.CertInfo,
			&asset.ResponseTimeMs,
			&asset.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assets: %w", err)
	}

	return assets, nil
}
