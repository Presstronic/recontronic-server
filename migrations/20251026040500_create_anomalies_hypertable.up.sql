-- Create anomalies table as a TimescaleDB hypertable for time-series anomaly detection
CREATE TABLE IF NOT EXISTS anomalies (
    id SERIAL,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    asset_id INTEGER,
    scan_job_id INTEGER REFERENCES scan_jobs(id) ON DELETE SET NULL,
    anomaly_type VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    evidence JSONB DEFAULT '{}'::jsonb,
    base_probability FLOAT,
    posterior_probability FLOAT,
    priority_score FLOAT NOT NULL,
    is_reviewed BOOLEAN DEFAULT false,
    review_notes TEXT,
    reviewed_at TIMESTAMPTZ,
    reviewed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    PRIMARY KEY (detected_at, id)
);

-- Convert to TimescaleDB hypertable (partitioned by time)
SELECT create_hypertable('anomalies', 'detected_at',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Create indexes for efficient queries
CREATE INDEX idx_anomalies_program ON anomalies(program_id, detected_at DESC);
CREATE INDEX idx_anomalies_unreviewed ON anomalies(is_reviewed, detected_at DESC) WHERE is_reviewed = false;
CREATE INDEX idx_anomalies_priority ON anomalies(priority_score DESC, detected_at DESC);
CREATE INDEX idx_anomalies_type ON anomalies(anomaly_type, detected_at DESC);
CREATE INDEX idx_anomalies_scan ON anomalies(scan_job_id, detected_at DESC);
CREATE INDEX idx_anomalies_program_unreviewed ON anomalies(program_id, is_reviewed, priority_score DESC) WHERE is_reviewed = false;

-- Add constraint for valid anomaly types
ALTER TABLE anomalies ADD CONSTRAINT check_anomaly_type
    CHECK (anomaly_type IN (
        'new_subdomain',
        'new_endpoint',
        'status_code_change',
        'content_change',
        'tech_stack_change',
        'cert_change',
        'weekend_deployment',
        'off_hours_deployment',
        'rapid_changes',
        'new_port',
        'configuration_change'
    ));

-- Add constraint for priority score range
ALTER TABLE anomalies ADD CONSTRAINT check_priority_score
    CHECK (priority_score >= 0 AND priority_score <= 100);

-- Enable compression for older data
ALTER TABLE anomalies SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'program_id,anomaly_type'
);

-- Automatically compress data older than 14 days
SELECT add_compression_policy('anomalies', INTERVAL '14 days', if_not_exists => TRUE);

-- Automatically delete data older than 1 year (keep anomaly history)
SELECT add_retention_policy('anomalies', INTERVAL '1 year', if_not_exists => TRUE);
