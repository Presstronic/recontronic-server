-- Drop findings table and related objects
DROP TRIGGER IF EXISTS update_findings_updated_at ON findings;
DROP INDEX IF EXISTS idx_findings_anomaly_lookup;
DROP INDEX IF EXISTS idx_findings_reported;
DROP INDEX IF EXISTS idx_findings_bounty;
DROP INDEX IF EXISTS idx_findings_user;
DROP INDEX IF EXISTS idx_findings_anomaly;
DROP INDEX IF EXISTS idx_findings_severity;
DROP INDEX IF EXISTS idx_findings_status;
DROP INDEX IF EXISTS idx_findings_program;
DROP TABLE IF EXISTS findings;
