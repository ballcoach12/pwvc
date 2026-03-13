# Feature Specification: Zero Doubt Prioritization

**Feature Branch**: `001-zero-doubt-prioritization`  
**Created**: 2025-10-26  
**Status**: Draft  
**Input**: Specifications for the Zero Doubt Roadmap: structured pairwise value/complexity, Fibonacci scoring, consensus, and final priority scoring with group alignment

## User Scenarios & Testing (mandatory)

### User Story 1 - Create and run a prioritization session (Priority: P1)
A facilitator creates a new prioritization session for a project, invites participants, imports or enters candidate features, and starts the session. Participants can join and see the session overview and their current phase.

Why this priority: Enables the core activity—creating a shared workspace to prioritize features.

Independent Test: A facilitator creates a session, invites at least 3 participants, and each participant can join and see session readiness without performing comparisons yet.

Acceptance Scenarios:
1. Given a facilitator is authenticated, When they create a session with a name and goal, Then a unique session is created with invite link and empty feature list.
2. Given a session exists, When the facilitator uploads a CSV of features or adds features manually, Then the features are visible to all participants in the session overview.
3. Given a session with at least one attendee, When attendees open the invite link, Then they can join and see the session status (not started / in progress / paused / completed).

---

### User Story 2 - Pairwise value comparisons (Priority: P1)
Participants perform head-to-head comparisons between pairs of features to determine relative value, contributing to a win-count weighting metric.

Why this priority: Pairwise value scoring is core to the methodology.

Independent Test: With a prepared session (features loaded), participants complete all required value comparisons until coverage is reached; the system computes WValue for each feature.

Acceptance Scenarios:
1. Given a set of N features, When pairwise value comparisons are initiated, Then the system schedules all required comparisons such that every feature is compared sufficiently to establish relative standing.
2. Given two features A and B, When a participant chooses the higher-value feature, Then the system records a win for the chosen feature.
3. Given all value comparisons for a participant are complete, When the facilitator views interim results, Then the system shows per-feature WValue derived from aggregate win counts (wins / total comparisons) across participants who have completed the phase.

---

### User Story 3 - Pairwise complexity comparisons (Priority: P1)
Participants perform head-to-head comparisons to assess relative implementation complexity.

Why this priority: Complexity weighting is required for the final priority calculation.

Independent Test: With value comparisons completed, participants complete all required complexity comparisons until coverage is reached; the system computes WComplexity for each feature.

Acceptance Scenarios:
1. Given the same set of features, When complexity comparisons are initiated, Then the system schedules comparisons analogous to value comparisons.
2. Given two features A and B, When a participant selects the more complex feature, Then the system records a complexity win accordingly.
3. Given complexity comparisons are complete, When the facilitator views interim results, Then per-feature WComplexity is shown from aggregate complexity wins / total complexity comparisons.

---

### User Story 4 - Fibonacci absolute scoring (Priority: P2)
Participants assign absolute Fibonacci scores to features for value and complexity (e.g., 1, 2, 3, 5, 8…).

Why this priority: Absolute scoring anchors relative rankings with magnitude.

Independent Test: Participants assign Fibonacci scores for value and complexity; system computes SValue and SComplexity per feature.

Acceptance Scenarios:
1. Given a feature list, When a participant assigns Fibonacci value and complexity, Then the scores are saved and can be edited until consensus is reached.
2. Given multiple participants, When the facilitator views score dispersion, Then the system shows distributions and outliers per feature.

---

### User Story 5 - Consensus and conflict resolution (Priority: P2)
The group reaches consensus on each feature’s value and complexity—disagreements are identified and resolved with guided facilitation.

Why this priority: The methodology requires full team agreement.

Independent Test: For any feature with divergent scores, the system flags it and provides tools to reconcile until consensus is achieved; consensus state is recorded.

Acceptance Scenarios:
1. Given divergent scores on a feature, When consensus mode is started, Then the system highlights disagreements, shows rationales if provided, and supports discussion until a single agreed score is set.
2. Given a feature reaches consensus, When the facilitator locks it, Then further edits require unlocking by the facilitator.

---

### User Story 6 - Final priority score and results (Priority: P1)
The system calculates the Final Priority Score per feature: FPS = (SValue × WValue) / (SComplexity × WComplexity). Results are visualized and exportable.

Why this priority: Produces the decision-making output.

Independent Test: With WValue/WComplexity and SValue/SComplexity available (and consensus reached), FPS is computed for all features and sorted; export produces a CSV including all components.

Acceptance Scenarios:
1. Given completed value/complexity comparisons and scores, When results are calculated, Then each feature has SValue, WValue, SComplexity, WComplexity, and FPS with a deterministic sort (descending FPS, ties resolved deterministically).
2. Given results exist, When the facilitator exports, Then a CSV (or spreadsheet) is generated with feature metadata and all scoring components.

---

### User Story 7 - Real-time collaboration (Priority: P2)
Participants see updates in real-time: who is in which phase, comparison assignments, and consensus status.

Why this priority: Minimizes coordination overhead and improves engagement.

Independent Test: Simulate multiple participants making updates; changes propagate to all connected clients within acceptable latency.

Acceptance Scenarios:
1. Given multiple participants connected, When one completes a comparison, Then other clients see updated progress within the latency budget.
2. Given a facilitator pauses the session, When participants are connected, Then all see the paused state and cannot submit new inputs until resumed.

---

### User Story 8 - Session progress and facilitation controls (Priority: P2)
The facilitator can start/stop phases, reassign comparisons, monitor progress, and manage attendees.

Why this priority: Ensures efficient session orchestration.

Independent Test: Facilitator dashboard shows per-phase completion metrics; facilitator can adjust assignments and control state transitions.

Acceptance Scenarios:
1. Given a phase is in progress, When the facilitator checks progress, Then the system shows % completion per attendee and global coverage.
2. Given unbalanced workloads, When the facilitator reassigns pending comparisons, Then participants receive updated queues.

---

### User Story 9 - Auditability and transparency (Priority: P3)
The system logs key decisions and consensus outcomes, preserving trail for later review.

Why this priority: Enables trust and organizational alignment.

Independent Test: After completing a session, an audit report shows who voted on what (configurable by anonymity setting), the consensus timeline, and changes.

Acceptance Scenarios:
1. Given a completed session, When the facilitator requests an audit log, Then the report contains participants, actions, timestamps, and final consensus per feature according to the session’s anonymity setting.

---

### Edge Cases
- N=1 or small N (2–3 features): ensure comparisons and W metrics are valid and defined.
- Ties and incomplete participant submissions: partial aggregation should not block others; the facilitator can advance when quorum is met. Default quorum: ≥ 70% of assigned comparisons/scores completed by active attendees; facilitator may override with justification.
- Late joiners or disconnects: catch-up and reassignment without double-counting.
- Feature duplicates or edits mid-session: safe update flows and re-computation strategy.

## Requirements (mandatory)

### Functional Requirements
- FR-001: The system MUST allow creation, naming, and configuration of a prioritization session with invite links and participant roles (facilitator, participant, observer).
- FR-002: The system MUST support feature intake: manual entry and CSV import (validated headers and duplicates handling).
- FR-003: The system MUST schedule and record pairwise value comparisons for all features with deterministic coverage across participants.
- FR-004: The system MUST compute WValue per feature as wins divided by total value comparisons included in aggregation.
- FR-005: The system MUST schedule and record pairwise complexity comparisons for all features and compute WComplexity analogously.
- FR-006: The system MUST support absolute Fibonacci scoring for value and complexity; per-feature SValue and SComplexity are derived via consensus.
- FR-007: The system MUST provide a consensus workflow to flag divergence, capture discussion inputs (optional rationales), and lock agreed scores.
- FR-008: The system MUST compute Final Priority Score per feature as FPS = (SValue × WValue) / (SComplexity × WComplexity).
- FR-009: The system MUST visualize results (sortable table, charts optional) and export to CSV including all components (SValue, WValue, SComplexity, WComplexity, FPS).
- FR-010: The system MUST provide real-time updates for participant actions and session state.
- FR-011: The system MUST provide facilitator controls: start/stop phases, progress view, reassignment of pending comparisons, attendee management.
- FR-012: The system MUST persist session data (projects, attendees, features, comparisons, scores, consensus outcomes, results).
- FR-013: The system MUST log key actions for audit with timestamps and actor identifiers, honoring session anonymity settings. Default privacy: Anonymous during the session; after finalization, only facilitators may access identity-linked logs. Exports remain aggregated unless explicitly configured otherwise.
- FR-014: The system MUST allow reruns and recalculations if features are modified before finalization, with change tracking in audit logs.
- FR-015: The system MUST support tie-breaking rules for equal FPS with a deterministic order: higher SValue first, then lower SComplexity, then alphabetical by feature name.

### Non-Functional Requirements
- NFR-001: Real-time update latency SHOULD be ≤ 500 ms p95 for in-session actions under typical team sizes.
- NFR-002: The system SHOULD support at least 25 concurrent participants in a session without degradation.
- NFR-003: Data integrity MUST be preserved across concurrent writes; operations must be idempotent where retried.
- NFR-004: Authorization gating MUST prevent non-facilitators from using facilitation controls.
- NFR-005: Availability target SHOULD be 99.9% monthly for session operations (if deployed in production).
- NFR-006: Privacy: if anonymity is enabled, individual votes MUST not be exposed in UI/exports except to authorized roles per policy.
- NFR-007: Accessibility: interfaces MUST be usable with keyboard navigation and meet WCAG AA for critical flows (pairwise, scoring, review).
- NFR-008: Performance budgets SHOULD be defined for UI load and results rendering (e.g., initial interactive ≤ 3s on standard hardware).

### Assumptions
- Teams generally range from 3–25 participants.
- Features are independent items; no dependency modeling is considered in prioritization.
- CSV schema will include at minimum: feature id/name and optional description.

## Decisions from Clarifications

- Quorum (Scope): Majority threshold selected—Default quorum is ≥ 70% completion by active attendees; facilitator may override with justification. Rationale: Keeps sessions moving while still reflecting a broad sample; override covers exceptional cases.
- Privacy: Anonymous during the session; after finalization, only facilitators may access identity-linked logs; participant-facing views and exports remain aggregated by default. Rationale: Encourages honest input while preserving accountability for facilitators.
- Tie-break: Deterministic order—higher SValue first, then lower SComplexity, then alphabetical by feature name. Rationale: Prioritizes user value, then implementation feasibility, with a stable final key.

## Key Entities (if data involved)
- Project/Session: id, name, status, config (anonymity, quorum rules), createdAt, updatedAt
- Attendee: id, displayName, role (facilitator, participant, observer), join status
- Feature: id, name, description, status
- Comparison: id, type (value|complexity), featureAId, featureBId, chosenId, participantId, timestamp
- Score: id, featureId, type (value|complexity), fibonacciValue, participantId, timestamp
- Consensus: featureId, agreedValueScore, agreedComplexityScore, lockedBy, lockedAt, rationale (optional)
- PriorityResult: featureId, SValue, WValue, SComplexity, WComplexity, FPS, rank
- Progress: per attendee per phase (assigned count, completed count, remaining), phase timestamps
- AuditLog: actorId, actionType, subjectType, subjectId, before/after (where applicable), timestamp, sessionId

## Success Criteria (mandatory)
- SC-001: A team can complete a session for 20 features with 10 participants (value, complexity, Fibonacci, consensus) in under 60 minutes.
- SC-002: 95% of participant actions (submit comparison/score) reflect on other clients within 500 ms p95.
- SC-003: Final export includes all features with SValue, WValue, SComplexity, WComplexity, FPS and matches on-screen values exactly.
- SC-004: Facilitator dashboard shows progress metrics with ±1 item accuracy versus stored data.
- SC-005: At least 90% of sessions end with 0 unresolved consensus flags.

## Edge Cases
- Minimal features (N=1): W metrics defined as 1.0 by convention; comparisons phase can be skipped.
- Odd number of comparisons remaining with participant dropout: reassignment without duplication.
- CSV import with duplicates: merge prompt or duplicate prevention with clear feedback.
- Feature edits after partial comparisons: invalidated comparisons for changed features are tracked and require redo.

---

 
