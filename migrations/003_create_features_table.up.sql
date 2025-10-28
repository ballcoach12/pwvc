-- Create features table
CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    acceptance_criteria TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index on project_id for faster joins
CREATE INDEX idx_features_project_id ON features(project_id);

-- Create index on created_at for sorting
CREATE INDEX idx_features_created_at ON features(created_at);

-- Create index on title for searching
CREATE INDEX idx_features_title ON features(title);