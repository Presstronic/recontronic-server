-- Create continuous aggregates for pre-computed analytics (TimescaleDB feature)
-- These views are automatically refreshed and provide blazing-fast query performance

-- 1. Daily asset statistics per program
CREATE MATERIALIZED VIEW daily_asset_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', discovered_at) AS day,
    program_id,
    asset_type,
    COUNT(*) as asset_count,
    COUNT(*) FILTER (WHERE is_live) as live_count,
    AVG(response_time_ms) as avg_response_time,
    MIN(response_time_ms) as min_response_time,
    MAX(response_time_ms) as max_response_time,
    COUNT(DISTINCT asset_value) as unique_assets
FROM assets
GROUP BY day, program_id, asset_type
WITH NO DATA;

-- Create index on the continuous aggregate
CREATE INDEX idx_daily_asset_stats_lookup ON daily_asset_stats(program_id, day DESC);

-- Auto-refresh policy: refresh hourly for last month of data
SELECT add_continuous_aggregate_policy('daily_asset_stats',
    start_offset => INTERVAL '1 month',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

-- 2. Hourly anomaly summary per program
CREATE MATERIALIZED VIEW hourly_anomaly_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', detected_at) AS hour,
    program_id,
    anomaly_type,
    COUNT(*) as anomaly_count,
    AVG(priority_score) as avg_priority,
    MAX(priority_score) as max_priority,
    COUNT(*) FILTER (WHERE is_reviewed) as reviewed_count,
    COUNT(*) FILTER (WHERE is_reviewed = false) as unreviewed_count,
    COUNT(*) FILTER (WHERE priority_score >= 80) as high_priority_count,
    COUNT(*) FILTER (WHERE priority_score >= 50 AND priority_score < 80) as medium_priority_count,
    COUNT(*) FILTER (WHERE priority_score < 50) as low_priority_count
FROM anomalies
GROUP BY hour, program_id, anomaly_type
WITH NO DATA;

-- Create index on the continuous aggregate
CREATE INDEX idx_hourly_anomaly_stats_lookup ON hourly_anomaly_stats(program_id, hour DESC);

-- Auto-refresh policy: refresh every 15 minutes for last week of data
SELECT add_continuous_aggregate_policy('hourly_anomaly_stats',
    start_offset => INTERVAL '1 week',
    end_offset => INTERVAL '15 minutes',
    schedule_interval => INTERVAL '15 minutes',
    if_not_exists => TRUE
);

-- 3. Asset change rate (for detecting spikes in activity)
CREATE MATERIALIZED VIEW hourly_asset_change_rate
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', discovered_at) AS hour,
    program_id,
    COUNT(*) as changes,
    COUNT(DISTINCT asset_type) as unique_asset_types,
    COUNT(*) FILTER (WHERE asset_type = 'subdomain') as new_subdomains,
    COUNT(*) FILTER (WHERE asset_type = 'endpoint') as new_endpoints,
    COUNT(*) FILTER (WHERE asset_type = 'url') as new_urls
FROM assets
GROUP BY hour, program_id
WITH NO DATA;

-- Create index on the continuous aggregate
CREATE INDEX idx_hourly_change_rate_lookup ON hourly_asset_change_rate(program_id, hour DESC);

-- Auto-refresh policy: refresh every 30 minutes for last 48 hours
SELECT add_continuous_aggregate_policy('hourly_asset_change_rate',
    start_offset => INTERVAL '48 hours',
    end_offset => INTERVAL '30 minutes',
    schedule_interval => INTERVAL '30 minutes',
    if_not_exists => TRUE
);
