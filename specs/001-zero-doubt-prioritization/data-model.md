# Data Model

Derived from spec and current codebase.

## Entities

- Project/Session
  - project: id, name, status, config (anonymity, quorum), createdAt, updatedAt
  - pairwise_sessions: id, project_id, criterion_type (value|complexity), status, started_at, completed_at
- Attendee
  - attendees: id, project_id, name, role (facilitator|participant|observer), auth fields, created_at
- Feature
  - features: id, project_id, title, description, acceptance_criteria, created_at, updated_at
- Pairwise Comparison
  - pairwise_comparisons: id, session_id, feature_a_id, feature_b_id, winner_id (nullable), is_tie, consensus_reached, created_at
  - attendee_votes: id, comparison_id, attendee_id, preferred_feature_id (nullable), is_tie_vote, voted_at
- Fibonacci Scoring
  - value_scores: id, project_id, feature_id, attendee_id, fibonacci_value, created_at
  - complexity_scores: id, project_id, feature_id, attendee_id, fibonacci_value, created_at
  - consensus_scores: id, project_id, feature_id, s_value, s_complexity, locked_by, locked_at, rationale (nullable)
- Priority Calculation
  - priority_calculations: id, project_id, feature_id, s_value, w_value, s_complexity, w_complexity, weighted_value, weighted_complexity, final_priority_score, rank, calculated_at
- Progress
  - project_progress: project_id, phase fields (pairwise_value, pairwise_complexity, fibonacci_value, fibonacci_complexity, results), completion flags/counters
- Audit Log
  - audit_logs: id, project_id, actor_id, action_type, subject_type, subject_id, before, after, timestamp

## Relationships

- project 1—N attendees, features, sessions, progress, logs
- session 1—N pairwise_comparisons; comparison 1—N attendee_votes
- feature has many votes and scores; exactly one consensus_scores row when locked

## Validation rules

- Fibonacci scores ∈ {1,2,3,5,8,13,21,34,55,89}
- No duplicate attendee vote per comparison
- No duplicate attendee Fibonacci score per feature per type
- Consensus lock only by facilitator role

## Notes

- Migrations up to 005 include fibonacci scaffolding; add tables if missing for consensus_scores and audit_logs.
