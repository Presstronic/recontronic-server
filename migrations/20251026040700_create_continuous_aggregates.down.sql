-- Drop continuous aggregates and their policies

-- Remove hourly_asset_change_rate
DROP INDEX IF EXISTS idx_hourly_change_rate_lookup;
DROP MATERIALIZED VIEW IF EXISTS hourly_asset_change_rate;

-- Remove hourly_anomaly_stats
DROP INDEX IF EXISTS idx_hourly_anomaly_stats_lookup;
DROP MATERIALIZED VIEW IF EXISTS hourly_anomaly_stats;

-- Remove daily_asset_stats
DROP INDEX IF EXISTS idx_daily_asset_stats_lookup;
DROP MATERIALIZED VIEW IF EXISTS daily_asset_stats;
