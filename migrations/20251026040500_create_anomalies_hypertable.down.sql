-- Remove TimescaleDB policies
SELECT remove_retention_policy('anomalies', if_exists => TRUE);
SELECT remove_compression_policy('anomalies', if_exists => TRUE);

-- Drop indexes
DROP INDEX IF EXISTS idx_anomalies_program_unreviewed;
DROP INDEX IF EXISTS idx_anomalies_scan;
DROP INDEX IF EXISTS idx_anomalies_type;
DROP INDEX IF EXISTS idx_anomalies_priority;
DROP INDEX IF EXISTS idx_anomalies_unreviewed;
DROP INDEX IF EXISTS idx_anomalies_program;

-- Drop the hypertable (this also drops the table)
DROP TABLE IF EXISTS anomalies;
