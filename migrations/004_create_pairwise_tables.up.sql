-- Create pairwise_sessions table
CREATE TABLE pairwise_sessions (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    criterion_type VARCHAR(20) NOT NULL CHECK (criterion_type IN ('value', 'complexity')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Create pairwise_comparisons table
CREATE TABLE pairwise_comparisons (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES pairwise_sessions(id) ON DELETE CASCADE,
    feature_a_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    feature_b_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    winner_id INTEGER REFERENCES features(id), -- NULL if tie
    is_tie BOOLEAN DEFAULT FALSE,
    consensus_reached BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create attendee_votes table
CREATE TABLE attendee_votes (
    id SERIAL PRIMARY KEY,
    comparison_id INTEGER REFERENCES pairwise_comparisons(id) ON DELETE CASCADE,
    attendee_id INTEGER REFERENCES attendees(id) ON DELETE CASCADE,
    preferred_feature_id INTEGER REFERENCES features(id), -- NULL if tie vote
    is_tie_vote BOOLEAN DEFAULT FALSE,
    voted_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(comparison_id, attendee_id)
);

-- Create indexes for better performance
CREATE INDEX idx_pairwise_sessions_project_id ON pairwise_sessions(project_id);
CREATE INDEX idx_pairwise_sessions_status ON pairwise_sessions(status);
CREATE INDEX idx_pairwise_comparisons_session_id ON pairwise_comparisons(session_id);
CREATE INDEX idx_pairwise_comparisons_consensus ON pairwise_comparisons(consensus_reached);
CREATE INDEX idx_attendee_votes_comparison_id ON attendee_votes(comparison_id);
CREATE INDEX idx_attendee_votes_attendee_id ON attendee_votes(attendee_id);