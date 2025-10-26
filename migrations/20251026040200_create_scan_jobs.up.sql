-- Create scan_jobs table for tracking scan execution
CREATE TABLE IF NOT EXISTS scan_jobs (
    id SERIAL PRIMARY KEY,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    results_count INTEGER DEFAULT 0,
    error_message TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index for program lookups
CREATE INDEX idx_scans_program ON scan_jobs(program_id);

-- Create index for status queries
CREATE INDEX idx_scans_status ON scan_jobs(status);

-- Create composite index for program + status queries
CREATE INDEX idx_scans_program_status ON scan_jobs(program_id, status);

-- Create index for job type
CREATE INDEX idx_scans_job_type ON scan_jobs(job_type);

-- Create index for created_at for time-based queries
CREATE INDEX idx_scans_created_at ON scan_jobs(created_at DESC);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_scan_jobs_updated_at
    BEFORE UPDATE ON scan_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add constraint to ensure valid status values
ALTER TABLE scan_jobs ADD CONSTRAINT check_scan_job_status
    CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'));

-- Add constraint to ensure valid job types
ALTER TABLE scan_jobs ADD CONSTRAINT check_scan_job_type
    CHECK (job_type IN ('passive', 'active', 'deep', 'manual'));
