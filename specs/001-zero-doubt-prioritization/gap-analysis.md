# Gap Analysis

Date: 2025-10-26

## High-impact gaps (P1)
1) Real win-count aggregation (FR‑004/005)
   - Current: ResultsService uses mock weights
   - Needed: Aggregate AttendeeVote per feature and criterion → wins, ties, totals → WValue/WComplexity
   - Plan: Add repository methods, service aggregation, unit tests; update results calc

2) Fibonacci scoring APIs + persistence (FR‑006)
   - Current: Validation helpers only
   - Needed: POST endpoints to submit per‑attendee value/complexity Fibonacci; list per feature; edit until consensus
   - Plan: Create tables if missing; add service/repo; update progress

3) Consensus lock workflow (FR‑007)
   - Current: WS message types; no REST for lock/unlock; no rationale
   - Needed: POST /consensus/{feature_id} with s_value, s_complexity, rationale; enforce facilitator role
   - Plan: Implement handler/service; persist consensus_scores; broadcast WS

4) Deterministic tie‑break (FR‑015)
   - Current: Sort by FPS only
   - Needed: Sort by FPS desc, then SValue desc, SComplexity asc, Title asc
   - Plan: Adjust sort in ResultsService and any duplicate sort sites

5) Audit logging + privacy (FR‑013, NFR‑006)
   - Current: No audit log module
   - Needed: audit_logs table + service; record actions; UI/API exposure filtered per privacy
   - Plan: Implement middleware/helpers; role‑gated queries; export guards

## Medium-impact gaps (P2)
6) Facilitator reassignment controls (FR‑011)
   - Add endpoint to reassign pending comparisons; update queues

7) Authorization enforcement (NFR‑004)
   - Add role checks on controls (start/stop/lock/export anonymized)

8) Progress completeness rules (quorum/override)
   - Honor 70% quorum + facilitator override on phase transitions

## Quality and Ops (P2)
9) Toolchain reliability
   - Tests failed locally due to Go directive; move tests to Docker image

10) Observability
   - Add Prometheus counters for WS events, votes/scores; add trace IDs in logs

## Out of scope for this milestone (P3)
- Frontend a11y polish (tracked separately)
- Load/perf testing beyond smoke checks
