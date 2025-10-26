-- Drop deployment_patterns table and related objects
DROP TRIGGER IF EXISTS update_deployment_patterns_updated_at ON deployment_patterns;
DROP INDEX IF EXISTS idx_patterns_activity;
DROP INDEX IF EXISTS idx_patterns_time;
DROP INDEX IF EXISTS idx_patterns_program;
DROP TABLE IF EXISTS deployment_patterns;
