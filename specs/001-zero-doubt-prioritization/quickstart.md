# Quickstart for P‑WVC Feature Prioritization

## Prerequisites
- Docker and Docker Compose installed
- curl or similar HTTP client for API testing

## Start stack

```bash
# From repo root
docker compose up --build -d
```

API will be available at http://localhost:8080.
WebSocket endpoint: ws://localhost:8080/api/ws/{projectId}

## Complete End-to-End Workflow

### 1. Create Project and Setup

```bash
# Create a new project
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mobile App Feature Prioritization",
    "description": "Q1 2024 feature roadmap prioritization",
    "invite_code": "mobile2024"
  }'
# Response: {"id": 1, "name": "Mobile App Feature Prioritization", ...}

# Create facilitator (required for project management)
curl -X POST http://localhost:8080/api/projects/1/attendees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Manager",
    "role": "PM",
    "is_facilitator": true
  }'
# Response: {"id": 1, "name": "Product Manager", ...}

# Login as facilitator to get session token
curl -X POST http://localhost:8080/api/projects/1/attendees/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product Manager",
    "role": "PM"
  }'
# Response: {"attendee": {...}, "session_token": "eyJ0eXAiOiJKV..."}
# Use this token in subsequent requests: -H "Authorization: Bearer {token}"

# Add team members
curl -X POST http://localhost:8080/api/projects/1/attendees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "name": "Lead Developer",
    "role": "Dev",
    "is_facilitator": false
  }'

curl -X POST http://localhost:8080/api/projects/1/attendees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "name": "UX Designer",
    "role": "Design",
    "is_facilitator": false
  }'
```

### 2. Add Features

```bash
# Add features individually
curl -X POST http://localhost:8080/api/projects/1/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "title": "Dark Mode Theme",
    "description": "Implement dark mode UI theme with user preference toggle"
  }'

curl -X POST http://localhost:8080/api/projects/1/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "title": "Push Notifications",
    "description": "Real-time push notifications for important updates"
  }'

curl -X POST http://localhost:8080/api/projects/1/features \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "title": "Offline Mode",
    "description": "Enable core app functionality when internet is unavailable"
  }'

# Or import features from CSV
curl -X POST http://localhost:8080/api/projects/1/features/import \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -F "file=@features.csv"
```

### 3. Phase 1: Pairwise Comparisons

```bash
# Advance to pairwise value phase
curl -X POST http://localhost:8080/api/projects/1/progress/advance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"phase": "pairwise_value"}'

# Start pairwise session for value criterion
curl -X POST http://localhost:8080/api/projects/1/pairwise \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"criterion_type": "value"}'

# Get active session and available comparisons
curl -X GET "http://localhost:8080/api/projects/1/pairwise?type=value" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Team members submit votes
curl -X POST http://localhost:8080/api/projects/1/pairwise/votes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "attendee_id": 2,
    "comparison_id": 1,
    "preferred_feature_id": 1,
    "is_tie_vote": false
  }'

# Repeat for complexity criterion
curl -X POST http://localhost:8080/api/projects/1/progress/advance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"phase": "pairwise_complexity"}'

curl -X POST http://localhost:8080/api/projects/1/pairwise \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"criterion_type": "complexity"}'
```

### 4. Phase 2: Fibonacci Scoring

```bash
# Advance to Fibonacci scoring phases
curl -X POST http://localhost:8080/api/projects/1/progress/advance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"phase": "fibonacci_value"}'

# Team members submit Fibonacci scores for value
curl -X POST http://localhost:8080/api/projects/1/scores/value \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "attendee_id": 2,
    "feature_id": 1,
    "fibonacci_value": 8,
    "rationale": "High user demand based on surveys and competitive analysis"
  }'

curl -X POST http://localhost:8080/api/projects/1/scores/value \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "attendee_id": 3,
    "feature_id": 1,
    "fibonacci_value": 5,
    "rationale": "Good feature but may not drive immediate user acquisition"
  }'

# Get scoring progress
curl -X GET "http://localhost:8080/api/projects/1/progress/fibonacci?criterion=value" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Move to complexity scoring
curl -X POST http://localhost:8080/api/projects/1/progress/advance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"phase": "fibonacci_complexity"}'

# Submit complexity scores
curl -X POST http://localhost:8080/api/projects/1/scores/complexity \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "attendee_id": 2,
    "feature_id": 1,
    "fibonacci_value": 3,
    "rationale": "Mostly CSS changes, existing theme infrastructure available"
  }'
```

### 5. Phase 3: Consensus Management

```bash
# Facilitator locks consensus scores after team discussion
curl -X POST http://localhost:8080/api/projects/1/consensus/lock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{
    "feature_id": 1,
    "s_value": 8,
    "s_complexity": 3,
    "rationale": "Team agreed on high value after user research presentation"
  }'

# Get consensus status
curl -X GET http://localhost:8080/api/projects/1/consensus \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Unlock if changes needed
curl -X POST http://localhost:8080/api/projects/1/consensus/unlock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -d '{"feature_id": 1}'
```

### 6. Results and Export

```bash
# Calculate final priority results
curl -X POST http://localhost:8080/api/projects/1/calculate-results \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Get results
curl -X GET http://localhost:8080/api/projects/1/results \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Export results in different formats
curl -X GET "http://localhost:8080/api/projects/1/results/export?format=csv" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -o "priority_results.csv"

curl -X GET "http://localhost:8080/api/projects/1/results/export?format=json" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -o "priority_results.json"

curl -X GET "http://localhost:8080/api/projects/1/results/export?format=jira" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -o "jira_import.json"
```

### 7. Audit and Monitoring

```bash
# Get audit report (facilitator-only)
curl -X GET "http://localhost:8080/api/projects/1/audit?limit=100&include_personal_data=false" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."

# Export audit logs
curl -X GET "http://localhost:8080/api/projects/1/audit/export?format=csv" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..." \
  -o "audit_log.csv"

# Get audit statistics
curl -X GET http://localhost:8080/api/projects/1/audit/statistics \
  -H "Authorization: Bearer eyJ0eXAiOiJKV..."
```

## Real-time Collaboration

Connect to WebSocket for live updates:

```javascript
// JavaScript example for WebSocket connection
const ws = new WebSocket('ws://localhost:8080/api/ws/1');

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
    
    switch(message.type) {
        case 'vote_update':
            // Handle new vote submission
            updateVoteDisplay(message.data);
            break;
        case 'score_submitted':
            // Handle Fibonacci score submission
            updateScoringProgress(message.data);
            break;
        case 'consensus_locked':
            // Handle consensus lock
            updateConsensusStatus(message.data);
            break;
        case 'phase_changed':
            // Handle phase transition
            updatePhaseDisplay(message.data);
            break;
    }
};
```

## Testing Workflow

```bash
# Run all tests in containers (avoids local toolchain issues)
docker compose exec server go test ./...

# Run specific test suites
docker compose exec server go test ./internal/domain/...
docker compose exec server go test ./internal/service/...
docker compose exec server go test ./internal/api/...

# Build and verify
docker compose exec server go build .

# Integration tests
docker compose exec server go test ./internal/api/integration_test.go -v
```

## Complete API Reference

### Project Management
- POST `/api/projects` - Create project
- GET `/api/projects/{id}` - Get project details
- PUT `/api/projects/{id}` - Update project
- DELETE `/api/projects/{id}` - Delete project

### Attendee Management
- GET `/api/projects/{id}/attendees` - List attendees
- POST `/api/projects/{id}/attendees` - Add attendee
- POST `/api/projects/{id}/attendees/login` - Login attendee
- DELETE `/api/projects/{id}/attendees/{attendeeId}` - Remove attendee

### Feature Management
- GET `/api/projects/{id}/features` - List features
- POST `/api/projects/{id}/features` - Add feature
- PUT `/api/projects/{id}/features/{featureId}` - Update feature
- DELETE `/api/projects/{id}/features/{featureId}` - Delete feature
- POST `/api/projects/{id}/features/import` - Import CSV
- GET `/api/projects/{id}/features/export` - Export CSV

### Pairwise Comparisons
- POST `/api/projects/{id}/pairwise` - Start session
- GET `/api/projects/{id}/pairwise?type=value|complexity` - Get session
- GET `/api/projects/{id}/pairwise/comparisons?type=value|complexity` - Get comparisons
- POST `/api/projects/{id}/pairwise/votes` - Submit vote
- POST `/api/projects/{id}/pairwise/complete` - Complete session
- GET `/api/projects/{id}/pairwise/next` - Get next comparison
- POST `/api/projects/{id}/pairwise/reassign` - Reassign comparisons (facilitator)
- GET `/api/projects/{id}/pairwise/pending` - Get pending comparisons
- GET `/api/projects/{id}/pairwise/reassignment-options` - Get reassignment options

### Fibonacci Scoring
- POST `/api/projects/{id}/scores/value` - Submit value score
- POST `/api/projects/{id}/scores/complexity` - Submit complexity score
- GET `/api/projects/{id}/scores?criterion=value|complexity` - Get scores

### Consensus Management (Facilitator-only)
- POST `/api/projects/{id}/consensus/lock` - Lock consensus score
- POST `/api/projects/{id}/consensus/unlock` - Unlock consensus score
- GET `/api/projects/{id}/consensus` - Get consensus status

### Progress & Phase Management (Facilitator-only)
- GET `/api/projects/{id}/progress` - Get project progress
- POST `/api/projects/{id}/progress/advance` - Advance phase
- POST `/api/projects/{id}/progress/complete` - Complete phase
- GET `/api/projects/{id}/progress/phases` - Get available phases
- POST `/api/projects/{id}/progress/pause` - Pause phase
- POST `/api/projects/{id}/progress/resume` - Resume phase
- GET `/api/projects/{id}/progress/fibonacci?criterion=value|complexity` - Get Fibonacci progress
- GET `/api/projects/{id}/progress/details` - Get detailed progress

### Results (Facilitator-only)
- POST `/api/projects/{id}/calculate-results` - Calculate priority results
- GET `/api/projects/{id}/results` - Get results
- GET `/api/projects/{id}/results/export?format=csv|json|jira` - Export results
- GET `/api/projects/{id}/results/summary` - Get results summary
- GET `/api/projects/{id}/results/status` - Check calculation status
- GET `/api/projects/{id}/results/preview` - Preview export

### Audit & Reporting (Facilitator-only)
- GET `/api/projects/{id}/audit?limit=100&action_type=&include_personal_data=false` - Get audit report
- GET `/api/projects/{id}/audit/export?format=csv|json` - Export audit logs
- GET `/api/projects/{id}/audit/statistics` - Get audit statistics

### Real-time Collaboration
- GET `/api/ws/{projectId}` - WebSocket connection for live updates
- GET `/api/ws/stats` - WebSocket connection statistics

### Health & Monitoring
- GET `/health` - Basic health check

## P-WVC Methodology Summary

The P-WVC (PairWise Value-Complexity) methodology follows these phases:

1. **Setup**: Create project, add attendees and features
2. **Pairwise Value**: Head-to-head comparisons for business value
3. **Pairwise Complexity**: Head-to-head comparisons for implementation complexity
4. **Fibonacci Value**: Absolute magnitude scoring (1,2,3,5,8,13,21,34,55,89) for value
5. **Fibonacci Complexity**: Absolute magnitude scoring for complexity
6. **Consensus**: Facilitator-driven agreement on final scores
7. **Results**: Final Priority Score (FPS) = (SValue × WValue) / (SComplexity × WComplexity)

Where:
- S = Fibonacci absolute score (1-89)
- W = Pairwise relative weight (wins/total comparisons)

## Rate Limiting

API endpoints have rate limits:
- Voting endpoints: 60 requests/minute
- Scoring endpoints: 120 requests/minute
- Consensus operations: 30 requests/minute
- Phase changes: 10 requests/minute (facilitator)
- Audit operations: 20 requests/minute (facilitator)
- Exports: 5 requests/minute
- Imports: 3 requests/minute
- Default: 100 requests/minute

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets
- `Retry-After`: Seconds to wait before retrying (when rate limited)

Progressive penalties apply for repeated violations.