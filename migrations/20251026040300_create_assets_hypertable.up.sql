-- Create assets table as a TimescaleDB hypertable for time-series asset discovery data
CREATE TABLE IF NOT EXISTS assets (
    id SERIAL,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    discovered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    asset_type VARCHAR(50) NOT NULL,
    asset_value TEXT NOT NULL,
    is_live BOOLEAN DEFAULT false,
    status_code INTEGER,
    content_hash TEXT,
    tech_stack JSONB DEFAULT '[]'::jsonb,
    response_headers JSONB DEFAULT '{}'::jsonb,
    cert_info JSONB DEFAULT '{}'::jsonb,
    response_time_ms INTEGER,
    metadata JSONB DEFAULT '{}'::jsonb,
    PRIMARY KEY (discovered_at, id)
);

-- Convert to TimescaleDB hypertable (partitioned by time)
SELECT create_hypertable('assets', 'discovered_at',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- Create indexes for efficient queries
CREATE INDEX idx_assets_program ON assets(program_id, discovered_at DESC);
CREATE INDEX idx_assets_type_value ON assets(asset_type, asset_value, discovered_at DESC);
CREATE INDEX idx_assets_live ON assets(is_live, discovered_at DESC) WHERE is_live = true;
CREATE INDEX idx_assets_hash ON assets(content_hash, discovered_at DESC) WHERE content_hash IS NOT NULL;
CREATE INDEX idx_assets_type_program ON assets(program_id, asset_type, discovered_at DESC);

-- Add constraint for valid asset types
ALTER TABLE assets ADD CONSTRAINT check_asset_type
    CHECK (asset_type IN ('subdomain', 'url', 'ip', 'port', 'endpoint', 'parameter'));

-- Enable compression for older data (90% storage savings)
ALTER TABLE assets SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'program_id,asset_type'
);

-- Automatically compress data older than 7 days
SELECT add_compression_policy('assets', INTERVAL '7 days', if_not_exists => TRUE);

-- Automatically delete data older than 6 months (retention policy)
SELECT add_retention_policy('assets', INTERVAL '6 months', if_not_exists => TRUE);
