-- Create fibonacci scoring sessions table
CREATE TABLE fibonacci_sessions (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    criterion_type VARCHAR(20) NOT NULL CHECK (criterion_type IN ('value', 'complexity')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Create index for efficient querying
CREATE INDEX idx_fibonacci_sessions_project_id ON fibonacci_sessions(project_id);
CREATE INDEX idx_fibonacci_sessions_status ON fibonacci_sessions(status);

-- Create individual Fibonacci scores table
CREATE TABLE fibonacci_scores (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES fibonacci_sessions(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    attendee_id INTEGER REFERENCES attendees(id) ON DELETE CASCADE,
    score_value INTEGER NOT NULL CHECK (score_value IN (1,2,3,5,8,13,21,34,55,89)),
    scored_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, feature_id, attendee_id)
);

-- Create indexes for efficient querying
CREATE INDEX idx_fibonacci_scores_session_id ON fibonacci_scores(session_id);
CREATE INDEX idx_fibonacci_scores_feature_id ON fibonacci_scores(feature_id);
CREATE INDEX idx_fibonacci_scores_attendee_id ON fibonacci_scores(attendee_id);

-- Create consensus scores table
CREATE TABLE consensus_scores (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES fibonacci_sessions(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    final_score INTEGER NOT NULL CHECK (final_score IN (1,2,3,5,8,13,21,34,55,89)),
    consensus_reached_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, feature_id)
);

-- Create indexes for efficient querying
CREATE INDEX idx_consensus_scores_session_id ON consensus_scores(session_id);
CREATE INDEX idx_consensus_scores_feature_id ON consensus_scores(feature_id);