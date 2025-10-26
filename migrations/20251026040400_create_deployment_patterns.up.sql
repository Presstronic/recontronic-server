-- Create deployment_patterns table for learning program behavior over time
CREATE TABLE IF NOT EXISTS deployment_patterns (
    id SERIAL PRIMARY KEY,
    program_id INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL,
    hour_of_day INTEGER NOT NULL,
    change_count INTEGER DEFAULT 0,
    avg_changes FLOAT DEFAULT 0.0,
    stddev_changes FLOAT DEFAULT 0.0,
    last_updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_program_time_pattern UNIQUE(program_id, day_of_week, hour_of_day)
);

-- Create index for program lookups
CREATE INDEX idx_patterns_program ON deployment_patterns(program_id);

-- Create index for temporal pattern queries
CREATE INDEX idx_patterns_time ON deployment_patterns(day_of_week, hour_of_day);

-- Create index for high-activity patterns
CREATE INDEX idx_patterns_activity ON deployment_patterns(program_id, change_count DESC);

-- Add constraints to ensure valid day and hour values
ALTER TABLE deployment_patterns ADD CONSTRAINT check_day_of_week
    CHECK (day_of_week >= 0 AND day_of_week <= 6);

ALTER TABLE deployment_patterns ADD CONSTRAINT check_hour_of_day
    CHECK (hour_of_day >= 0 AND hour_of_day <= 23);

-- Create trigger to automatically update last_updated_at timestamp
CREATE TRIGGER update_deployment_patterns_updated_at
    BEFORE UPDATE ON deployment_patterns
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
