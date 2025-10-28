-- Final priority calculations table
CREATE TABLE priority_calculations (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    w_value DECIMAL(10,6) NOT NULL,        -- Win-count weight for value
    w_complexity DECIMAL(10,6) NOT NULL,   -- Win-count weight for complexity  
    s_value INTEGER NOT NULL,              -- Fibonacci score for value
    s_complexity INTEGER NOT NULL,         -- Fibonacci score for complexity
    weighted_value DECIMAL(10,6) NOT NULL, -- SValue × WValue
    weighted_complexity DECIMAL(10,6) NOT NULL, -- SComplexity × WComplexity
    final_priority_score DECIMAL(10,6) NOT NULL, -- Weighted Value ÷ Weighted Complexity
    rank INTEGER NOT NULL,
    calculated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, feature_id)
);

-- Index for efficient querying
CREATE INDEX idx_priority_calculations_project_id ON priority_calculations(project_id);
CREATE INDEX idx_priority_calculations_rank ON priority_calculations(project_id, rank);
CREATE INDEX idx_priority_calculations_score ON priority_calculations(project_id, final_priority_score DESC);