# Tasks: Zero Doubt Prioritization (P‑WVC)

Feature: 001-zero-doubt-prioritization
Spec: specs/001-zero-doubt-prioritization/spec.md
Plan: specs/001-zero-doubt-prioritization/plan.md
Contracts: specs/001-zero-doubt-prioritization/contracts/openapi.yaml

Notes
- Checklist format strictly enforced: `- [ ] T### [P?] [US#?] Description with file path`
- Stories ordered by priority: P1 → P2 → P3
- Tests optional per prompt; we provide independent test criteria per story below, no explicit test tasks

## Phase 1 — Setup

- [x] T001 Confirm feature branch exists (001-zero-doubt-prioritization)
- [x] T002 Run prerequisites script to load FEATURE_DIR and docs: .specify/scripts/bash/check-prerequisites.sh --json
- [x] T003 Ensure OpenAPI contract committed at specs/001-zero-doubt-prioritization/contracts/openapi.yaml (no changes required)
- [x] T004 [P] Prepare Docker-based test workflow docs in specs/001-zero-doubt-prioritization/quickstart.md (append run instructions)

## Phase 2 — Foundational (blocking)

- [x] T005 Add DB migration for consensus scores table at migrations/009_create_consensus_scores_table.up.sql
- [x] T006 Add rollback for consensus scores at migrations/009_create_consensus_scores_table.down.sql
- [x] T007 Add DB migration for audit logs table at migrations/010_create_audit_logs_table.up.sql
- [x] T008 Add rollback for audit logs at migrations/010_create_audit_logs_table.down.sql
- [x] T009 [P] Repository: add consensus persistence functions in internal/repository/consensus.go
- [x] T010 [P] Repository: add audit logging functions in internal/repository/audit.go
- [x] T011 Service: implement deterministic tie-break sort in internal/service/results.go (FPS desc, SValue desc, SComplexity asc, Name asc)
- [x] T012 Replace mock weight aggregation with real WValue/WComplexity in internal/service/results.go using repository methods
- [x] T013 [P] Repository: add aggregation queries for wins/total in internal/repository/pairwise.go (value and complexity)
- [x] T014 Middleware: add facilitator-only guard in internal/api/auth.go (IsFacilitator middleware)
- [x] T015 [P] WebSocket: define event types for score-submitted, consensus-locked, phase-changed in internal/websocket (hub/client/messages)

## Phase 3 — User Story 1 (P1): Create and run a session

Story Goal: Facilitator can create a session, import features, and invite participants; attendees can join and see status.
Independent Test Criteria: Facilitator creates project, adds features via CSV, at least 3 attendees join via invite and see session readiness.

- [x] T016 [US1] Verify/create invite link generation in internal/api/project.go (add invite link field if missing)
- [x] T017 [P] [US1] CSV import validation messages in internal/api/feature.go (ensure duplicate handling aligned with spec)
- [x] T018 [US1] Broadcast session status changes over WS in internal/api/project.go and internal/websocket/*

## Phase 4 — User Story 2 (P1): Pairwise value comparisons

Story Goal: Participants compare pairs for value; system computes WValue.
Independent Test Criteria: With features loaded, participants complete value comparisons; WValue aggregates correctly across completed voters.

- [x] T019 [US2] Ensure scheduling/queue for value comparisons exists in internal/service/pairwise.go (add missing logic or confirm)
- [x] T020 [P] [US2] Persist value votes endpoint in internal/api/pairwise.go (route for POST /pairwise-sessions/{session_id}/vote with criterion=value)
- [x] T021 [US2] Aggregate WValue in internal/service/results.go using repository functions (wins/total for value)
- [x] T022 [P] [US2] WS notify on value vote in internal/api/pairwise.go and internal/websocket/*

## Phase 5 — User Story 3 (P1): Pairwise complexity comparisons

Story Goal: Participants compare pairs for complexity; system computes WComplexity.
Independent Test Criteria: Participants complete complexity comparisons; WComplexity aggregates correctly.

- [x] T023 [US3] Ensure scheduling/queue for complexity comparisons in internal/service/pairwise.go (criterion=complexity)
- [x] T024 [P] [US3] Persist complexity votes via internal/api/pairwise.go (criterion=complexity branch)
- [x] T025 [US3] Aggregate WComplexity in internal/service/results.go (wins/total for complexity)
- [x] T026 [P] [US3] WS notify on complexity vote in internal/api/pairwise.go and internal/websocket/*

## Phase 6 — User Story 6 (P1): Final priority score and results

Story Goal: Compute FPS and return sorted results with deterministic tie-break; export includes all components.
Independent Test Criteria: With S* and W* available, results list and CSV match and are deterministically sorted.

- [x] T027 [US6] Apply deterministic sort in internal/service/results.go (reuse T011 implementation)
- [x] T028 [P] [US6] Ensure export includes SValue, WValue, SComplexity, WComplexity, FPS in internal/api/results.go
- [x] T029 [US6] Persist calculation snapshot in internal/repository/priority.go (confirm schema mapping to priority_calculations)

## Phase 7 — User Story 4 (P2): Fibonacci absolute scoring

Story Goal: Participants assign Fibonacci scores for value and complexity; system stores per-attendee scores.
Independent Test Criteria: Scores saved and retrievable; invalid values rejected.

- [x] T030 [US4] API: add endpoints for POST /projects/{id}/scores/value and /scores/complexity in internal/api/scoring.go (validate Fibonacci)
- [x] T031 [P] [US4] Service: implement scoring logic and validation in internal/service/scoring.go
- [x] T032 [P] [US4] Repository: persist scores to value_scores and complexity_scores in internal/repository/scoring.go
- [x] T033 [US4] WS: broadcast score-submitted events in internal/api/scoring.go and internal/websocket/*

## Phase 8 — User Story 5 (P2): Consensus and conflict resolution

Story Goal: Facilitate reconciliation and lock final SValue/SComplexity per feature.
Independent Test Criteria: Divergences flagged; facilitator can lock consensus; edits blocked unless unlocked.

- [x] T034 [US5] API: add POST /projects/{id}/consensus/{feature_id} in internal/api/consensus.go (require IsFacilitator)
- [x] T035 [P] [US5] Service: implement consensus lock/unlock in internal/service/consensus.go (write to consensus_scores)
- [x] T036 [P] [US5] Repository: upsert consensus row with rationale in internal/repository/consensus.go
- [x] T037 [US5] WS: broadcast consensus-locked event in internal/api/consensus.go and internal/websocket/*

## Phase 9 — User Story 7 (P2): Real-time collaboration

Story Goal: Real-time updates across participants and facilitator dashboard.
Independent Test Criteria: Multi-client updates propagate within latency budgets.

- [x] T038 [US7] WS: ensure hub routes events for votes, scores, consensus, and phase changes in internal/websocket/*
- [x] T039 [P] [US7] API: phase pause/resume broadcasts in internal/api/progress.go

## Phase 10 — User Story 8 (P2): Progress and facilitation controls

Story Goal: Facilitator views progress metrics and controls phase transitions and assignments.
Independent Test Criteria: Dashboard shows % completion; reassignments update queues.

- [x] T040 [US8] Add per-phase progress metrics for Fibonacci phases in internal/service/progress.go and internal/repository/progress.go
- [x] T041 [P] [US8] API: expose progress details endpoints in internal/api/progress.go
- [x] T042 [US8] Allow reassignment of pending comparisons in internal/service/pairwise.go and internal/api/pairwise.go

## Phase 11 — User Story 9 (P3): Auditability and transparency

Story Goal: Persist key actions and provide an audit report honoring anonymity settings.
Independent Test Criteria: After a session completes, audit report contains actions, timestamps, and consensus outcomes, respecting privacy rules.

- [ ] T043 [US9] Instrument audit logs on vote, score, consensus, phase changes in internal/api/* (call audit service)
- [ ] T044 [P] [US9] Service: implement audit logging helper with privacy gates in internal/service/audit.go
- [ ] T045 [P] [US9] Repository: persist audit rows in internal/repository/audit.go
- [ ] T046 [US9] API: add audit export endpoint (facilitator-only) in internal/api/results.go or internal/api/audit.go

## Final Phase — Polish & Cross-cutting

- [ ] T047 Add rate limiting/backoff for hot endpoints in internal/api/validation.go or middleware
- [ ] T048 [P] Enforce authZ gates (facilitator vs participant) across new endpoints in internal/api/*.go
- [ ] T049 [P] Update docs: specs/001-zero-doubt-prioritization/quickstart.md with end-to-end flow and curl examples
- [ ] T050 Ensure CSV export deterministic ordering and locale-safe formatting in internal/api/results.go

---

## Dependencies (Story Order)

1. Phase 1 → Phase 2 (foundational DB, repo, middleware)
2. P1 stories: US1 → US2 → US3 → US6
3. P2 stories: US4 → US5 → US7 → US8
4. P3 story: US9

Key task dependencies
- T013 before T012/T021/T025 (aggregations before consumption)
- T011 before T027 (tie-break implementation reused)
- T014 before T034/T039/T046 (authZ middleware before facilitator-only endpoints)
- T005–T008 before T034–T036 and T043–T046 (tables before writes)

## Parallel Opportunities (Examples)

- Foundational: T009, T010, T013, T015 can proceed in parallel.
- US2 vs US3: Value and complexity vote WS notifications (T022, T026) can be implemented in parallel.
- US4: Service (T031) and Repository (T032) can progress in parallel once contracts are set by T030.
- US5: Consensus service (T035) and repo (T036) in parallel after migrations.
- Polish: Documentation (T049) and authZ sweep (T048) in parallel.

## Implementation Strategy

- MVP Scope: P1 user stories (US1, US2, US3, US6) only. Ship end-to-end prioritization calculation with deterministic results and exports.
- Incremental Delivery: Merge Phase 2 foundational first, then US1→US2→US3→US6. Follow with P2 features in small PRs (US4, US5, US7, US8), then P3 (US9).
- Rollback Plan: Each phase isolated; DB migrations include downs (T006, T008). Keep feature toggles for new endpoints where feasible.
