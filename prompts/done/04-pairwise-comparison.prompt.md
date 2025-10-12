# Prompt 4: Pairwise Comparison Backend

Implement pairwise comparison sessions with database schema for comparisons, attendee votes, and consensus tracking. Create APIs for managing comparison sessions and voting workflows.

## Requirements
- Create database schema for pairwise comparison sessions and votes
- Implement session management (start, track progress, complete)
- Build voting system with consensus tracking
- Create APIs for comparison workflow
- Add logic to generate all required feature pairs
- Track individual attendee votes and group consensus

## Database Schema
```sql
-- Pairwise comparison sessions
CREATE TABLE pairwise_sessions (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    criterion_type VARCHAR(20) NOT NULL CHECK (criterion_type IN ('value', 'complexity')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Individual comparisons between two features
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

-- Individual attendee votes for each comparison
CREATE TABLE attendee_votes (
    id SERIAL PRIMARY KEY,
    comparison_id INTEGER REFERENCES pairwise_comparisons(id) ON DELETE CASCADE,
    attendee_id INTEGER REFERENCES attendees(id) ON DELETE CASCADE,
    preferred_feature_id INTEGER REFERENCES features(id), -- NULL if tie vote
    is_tie_vote BOOLEAN DEFAULT FALSE,
    voted_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(comparison_id, attendee_id)
);
```

## API Endpoints to Create
- `POST /api/projects/{id}/pairwise-sessions` - Start new pairwise session (Value or Complexity)
- `GET /api/projects/{id}/pairwise-sessions/{session_id}` - Get session details and progress
- `GET /api/projects/{id}/pairwise-sessions/{session_id}/comparisons` - Get all comparisons for session
- `POST /api/projects/{id}/pairwise-sessions/{session_id}/vote` - Submit attendee vote for comparison
- `POST /api/projects/{id}/pairwise-sessions/{session_id}/complete` - Mark session as completed

## Business Logic
- Generate all unique feature pairs (n × (n-1) / 2 comparisons)
- Track consensus: comparison is consensus when all attendees agree
- Calculate session progress: completed comparisons / total comparisons
- Auto-complete session when all comparisons reach consensus