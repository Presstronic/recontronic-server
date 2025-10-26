-- Create findings table for tracking discovered vulnerabilities and bug bounty submissions
CREATE TABLE IF NOT EXISTS findings (
    id SERIAL PRIMARY KEY,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    anomaly_id INTEGER,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    vulnerability_type VARCHAR(100),
    cvss_score FLOAT,
    reported_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    bounty_amount DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'USD',
    notes TEXT,
    poc TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index for program lookups
CREATE INDEX idx_findings_program ON findings(program_id);

-- Create index for status queries
CREATE INDEX idx_findings_status ON findings(status);

-- Create index for severity
CREATE INDEX idx_findings_severity ON findings(severity);

-- Create index for anomaly relationship
CREATE INDEX idx_findings_anomaly ON findings(anomaly_id) WHERE anomaly_id IS NOT NULL;

-- Create index for user's findings
CREATE INDEX idx_findings_user ON findings(user_id) WHERE user_id IS NOT NULL;

-- Create index for bounty tracking
CREATE INDEX idx_findings_bounty ON findings(bounty_amount DESC) WHERE bounty_amount IS NOT NULL;

-- Create index for date range queries
CREATE INDEX idx_findings_reported ON findings(reported_at DESC) WHERE reported_at IS NOT NULL;

-- Add constraint for valid severity values
ALTER TABLE findings ADD CONSTRAINT check_finding_severity
    CHECK (severity IN ('critical', 'high', 'medium', 'low', 'info'));

-- Add constraint for valid status values
ALTER TABLE findings ADD CONSTRAINT check_finding_status
    CHECK (status IN (
        'draft',
        'submitted',
        'triaged',
        'accepted',
        'duplicate',
        'informative',
        'not_applicable',
        'resolved',
        'bounty_awarded'
    ));

-- Add constraint for CVSS score range
ALTER TABLE findings ADD CONSTRAINT check_cvss_score
    CHECK (cvss_score IS NULL OR (cvss_score >= 0 AND cvss_score <= 10));

-- Add constraint for bounty amount
ALTER TABLE findings ADD CONSTRAINT check_bounty_amount
    CHECK (bounty_amount IS NULL OR bounty_amount >= 0);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_findings_updated_at
    BEFORE UPDATE ON findings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create foreign key reference to anomalies (note: can't use FK due to hypertable)
-- Instead, we'll validate this in application logic
CREATE INDEX idx_findings_anomaly_lookup ON findings(anomaly_id);
