# Prompt 8: Fibonacci Scoring Interface

Implement the Fibonacci scoring phase UI with individual scoring inputs, group consensus tracking, and visual indicators for agreement status on each feature's Value and Complexity scores.

## Requirements
- Create Fibonacci scoring interface for individual features
- Implement consensus tracking for each feature's scores
- Add real-time updates for group scoring progress
- Validate Fibonacci sequence inputs (1, 2, 3, 5, 8, 13, 21, 34, 55, 89)
- Show individual vs. consensus scores clearly
- Handle both Value and Complexity scoring sessions

## Database Schema Addition
```sql
-- Fibonacci scoring sessions
CREATE TABLE fibonacci_sessions (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    criterion_type VARCHAR(20) NOT NULL CHECK (criterion_type IN ('value', 'complexity')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed')),
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Individual Fibonacci scores
CREATE TABLE fibonacci_scores (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES fibonacci_sessions(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    attendee_id INTEGER REFERENCES attendees(id) ON DELETE CASCADE,
    score_value INTEGER NOT NULL CHECK (score_value IN (1,2,3,5,8,13,21,34,55,89)),
    scored_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, feature_id, attendee_id)
);

-- Consensus scores for features
CREATE TABLE consensus_scores (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES fibonacci_sessions(id) ON DELETE CASCADE,
    feature_id INTEGER REFERENCES features(id) ON DELETE CASCADE,
    final_score INTEGER NOT NULL CHECK (final_score IN (1,2,3,5,8,13,21,34,55,89)),
    consensus_reached_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, feature_id)
);
```

## Components to Create

### FibonacciScoringGrid Component
- Grid showing all features for scoring
- Fibonacci scale selector for each feature
- Individual attendee scores display
- Consensus status indicators

### FeatureScoringCard Component
- Feature title, description, acceptance criteria
- Fibonacci number picker (1,2,3,5,8,13,21,34,55,89)
- Show all attendee scores for this feature
- Consensus indicator and final agreed score

### FibonacciScalePicker Component
- Visual Fibonacci sequence selector
- Clear indication of selected value
- Validation for valid Fibonacci numbers
- Tooltips explaining score meanings

### ConsensusTracker Component
- Progress indicator for group consensus
- List of features and their consensus status
- Individual attendee progress tracking
- Session completion status

## API Endpoints to Add
- `POST /api/projects/{id}/fibonacci-sessions` - Start scoring session
- `POST /api/projects/{id}/fibonacci-sessions/{session_id}/scores` - Submit individual score
- `GET /api/projects/{id}/fibonacci-sessions/{session_id}/scores` - Get all scores for session
- `POST /api/projects/{id}/fibonacci-sessions/{session_id}/consensus` - Set consensus score for feature

## Business Logic
- Consensus achieved when all attendees agree on same score
- Session complete when all features have consensus scores
- Allow score changes until consensus is reached
- Calculate session progress percentage

## User Experience
- Clear visual distinction between individual and consensus scores
- Easy score modification before consensus
- Real-time updates from other attendees
- Progress indicators and completion status
- Mobile-friendly interface with large touch targets