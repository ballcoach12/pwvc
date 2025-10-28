# Requirements Traceability Matrix (RTM)

Maps Functional Requirements (FR) and Non‑Functional Requirements (NFR) to implementation.

Legend: ✅ Implemented | 🟡 Partial | ❌ Missing

## Functional
- FR‑001 Session creation, roles, invites → 🟡
  - internal/api/pairwise.go StartPairwiseSession (exists), attendee login in internal/api/auth.go; invite links/roles are basic; needs observer role and link mgmt.
- FR‑002 Feature intake (manual, CSV) → ✅
  - internal/api/feature.go import/export implemented; validation present.
- FR‑003 Schedule/record pairwise value with deterministic coverage → 🟡
  - internal/api/pairwise.go GetNextComparison/SubmitVote; scheduling exists but deterministic coverage policy not verified.
- FR‑004 Compute WValue (wins/total) → ❌
  - domain has CalculateWinCount; ResultsService uses mock weights. Needs real aggregation from votes.
- FR‑005 Complexity comparisons and WComplexity → ❌
  - Same as FR‑004 for complexity.
- FR‑006 Fibonacci absolute scoring (SValue/SComplexity) + edit until consensus → ❌
  - Validation present; no API/storage for per‑attendee Fibonacci scores.
- FR‑007 Consensus workflow, lock agreed scores → 🟡
  - WebSocket notifications exist; no REST to lock/unlock Fibonacci consensus per feature; no rationale persistence.
- FR‑008 Compute FPS and sort → 🟡
  - Implemented with mocks; replace with real S/W sources; add tie‑break per FR‑015.
- FR‑009 Visualize/export results → ✅ (backend)
  - internal/api/results.go (CSV/JSON/Jira) and summary endpoints.
- FR‑010 Real‑time updates → 🟡
  - WebSocket hub/messages present; verify broadcasts for progress and consensus.
- FR‑011 Facilitator controls (start/stop, reassign, progress, attendees) → 🟡
  - Start/complete exists; progress endpoints exist; reassignment not present; attendee mgmt limited.
- FR‑012 Persist data → ✅
  - Migrations and repositories for core entities.
- FR‑013 Audit logging w/ privacy → ❌
  - No audit log module; need table and service; enforce anonymity.
- FR‑014 Reruns/recalc w/ change tracking → 🟡
  - Recalc exists; change tracking/audit missing.
- FR‑015 Deterministic tie‑break (SValue desc, SComplexity asc, name) → ❌
  - ResultsService sorts by FPS only.

## Non‑Functional
- NFR‑001 Real‑time latency ≤ 500 ms p95 → 🟡
  - Not measured; WS architecture present.
- NFR‑002 ≥ 25 concurrent participants → 🟡
  - Likely feasible; needs load test.
- NFR‑003 Data integrity/idempotency → 🟡
  - Add unique constraints and idempotent upserts for votes/scores.
- NFR‑004 Authorization gating of controls → ❌
  - Roles exist; enforce in handlers.
- NFR‑005 Availability 99.9% → 🟡
  - Add health/readiness and process supervision; out of scope for dev.
- NFR‑006 Privacy/anonymity enforcement → ❌
  - Add query filters and export guards per policy.
- NFR‑007 Accessibility (WCAG AA) → 🟡
  - Frontend to implement.
- NFR‑008 UI performance budgets → 🟡
  - Define and check in web build.
