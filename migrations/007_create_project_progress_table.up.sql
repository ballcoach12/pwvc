-- Create project_progress table for session state management
CREATE TABLE project_progress (
    project_id INTEGER PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    setup_completed BOOLEAN DEFAULT FALSE,
    attendees_added BOOLEAN DEFAULT FALSE,
    features_added BOOLEAN DEFAULT FALSE,
    pairwise_value_completed BOOLEAN DEFAULT FALSE,
    pairwise_complexity_completed BOOLEAN DEFAULT FALSE,
    fibonacci_value_completed BOOLEAN DEFAULT FALSE,
    fibonacci_complexity_completed BOOLEAN DEFAULT FALSE,
    results_calculated BOOLEAN DEFAULT FALSE,
    current_phase VARCHAR(50) DEFAULT 'setup',
    last_activity TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX idx_project_progress_current_phase ON project_progress(current_phase);
CREATE INDEX idx_project_progress_last_activity ON project_progress(last_activity);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_project_progress_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.last_activity = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_project_progress_updated_at
    BEFORE UPDATE ON project_progress
    FOR EACH ROW
    EXECUTE FUNCTION update_project_progress_updated_at();