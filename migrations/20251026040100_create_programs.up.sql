-- Create programs table for storing bug bounty program configurations
CREATE TABLE IF NOT EXISTS programs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    platform VARCHAR(50),
    scope TEXT[],
    scan_frequency INTERVAL DEFAULT '1 hour',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_scanned_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Create index for active programs
CREATE INDEX idx_programs_active ON programs(is_active) WHERE is_active = true;

-- Create index for platform
CREATE INDEX idx_programs_platform ON programs(platform);

-- Create index for last_scanned_at to optimize scheduling queries
CREATE INDEX idx_programs_last_scanned ON programs(last_scanned_at) WHERE is_active = true;

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_programs_updated_at
    BEFORE UPDATE ON programs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
