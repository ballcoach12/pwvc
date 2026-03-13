-- Create projects table
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index on status for faster filtering
CREATE INDEX idx_projects_status ON projects(status);

-- Create index on created_at for sorting
CREATE INDEX idx_projects_created_at ON projects(created_at);