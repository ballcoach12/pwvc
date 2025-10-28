-- Create proper consensus scores table for P-WVC feature
-- Replaces session-based consensus with project-level feature consensus
CREATE TABLE consensus_scores (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    s_value INTEGER NOT NULL CHECK (s_value IN (1,2,3,5,8,13,21,34,55,89)),
    s_complexity INTEGER NOT NULL CHECK (s_complexity IN (1,2,3,5,8,13,21,34,55,89)),
    locked_by INTEGER REFERENCES attendees(id),
    locked_at TIMESTAMP DEFAULT NOW(),
    rationale TEXT,
    UNIQUE(project_id, feature_id)
);

-- Create indexes for efficient querying
CREATE INDEX idx_consensus_scores_project_id ON consensus_scores(project_id);
CREATE INDEX idx_consensus_scores_feature_id ON consensus_scores(feature_id);
CREATE INDEX idx_consensus_scores_locked_by ON consensus_scores(locked_by);