-- Drop scan_jobs table and related objects
DROP TRIGGER IF EXISTS update_scan_jobs_updated_at ON scan_jobs;
DROP INDEX IF EXISTS idx_scans_created_at;
DROP INDEX IF EXISTS idx_scans_job_type;
DROP INDEX IF EXISTS idx_scans_program_status;
DROP INDEX IF EXISTS idx_scans_status;
DROP INDEX IF EXISTS idx_scans_program;
DROP TABLE IF EXISTS scan_jobs;
