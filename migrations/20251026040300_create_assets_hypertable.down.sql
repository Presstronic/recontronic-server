-- Remove TimescaleDB policies
SELECT remove_retention_policy('assets', if_exists => TRUE);
SELECT remove_compression_policy('assets', if_exists => TRUE);

-- Drop indexes
DROP INDEX IF EXISTS idx_assets_type_program;
DROP INDEX IF EXISTS idx_assets_hash;
DROP INDEX IF EXISTS idx_assets_live;
DROP INDEX IF EXISTS idx_assets_type_value;
DROP INDEX IF EXISTS idx_assets_program;

-- Drop the hypertable (this also drops the table)
DROP TABLE IF EXISTS assets;
