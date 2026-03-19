-- Create audit logs table for P-WVC feature
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    actor_id INTEGER REFERENCES attendees(id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    subject_type VARCHAR(50) NOT NULL,
    subject_id INTEGER,
    before_state JSONB,
    after_state JSONB,
    timestamp TIMESTAMP DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX idx_audit_logs_project_id ON audit_logs(project_id);
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action_type ON audit_logs(action_type);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);