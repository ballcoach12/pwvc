-- Create attendees table
CREATE TABLE attendees (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100),
    is_facilitator BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create index on project_id for faster joins
CREATE INDEX idx_attendees_project_id ON attendees(project_id);

-- Create index on is_facilitator for filtering facilitators
CREATE INDEX idx_attendees_is_facilitator ON attendees(is_facilitator);