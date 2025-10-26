-- Drop programs table and related objects
DROP TRIGGER IF EXISTS update_programs_updated_at ON programs;
DROP INDEX IF EXISTS idx_programs_last_scanned;
DROP INDEX IF EXISTS idx_programs_platform;
DROP INDEX IF EXISTS idx_programs_active;
DROP TABLE IF EXISTS programs;
