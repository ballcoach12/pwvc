# Research and Decisions

Date: 2025-10-26
Feature: Zero Doubt Prioritization (P‑WVC)

## Decisions (resolved unknowns)

- Quorum: ≥ 70% completion by active attendees; facilitator override with justification.
- Privacy: Anonymous during session; after finalization, only facilitators may access identity-linked logs; exports default aggregated.
- Tie-break: Higher SValue, then lower SComplexity, then alphabetical by feature title.

## Backend integration approach

- Win-count weights: Use existing domain.CalculateWinCount and aggregate AttendeeVote from repository to compute WValue/WComplexity per feature. Replace ResultsService mock `calculateWinCountWeights` with real aggregation queries.
- Fibonacci scores: Introduce endpoints to submit per-attendee Fibonacci value/complexity for each feature; store in new tables (migrations exist up to 005 for fibonacci scoring). Compute consensus SValue/SComplexity per feature via facilitation flow (unanimity or facilitator lock after discussion), then persist.
- Consensus: Reuse websocket notifications and add REST endpoints to lock/unlock consensus per feature with optional rationale; update progress.
- Audit logging: Add `audit_logs` table and service to record actions (vote, score, lock) with actorId, subject, timestamps; honor anonymity via query guards and role checks.
- Tie-break: Apply FR-015 order in ResultsService sorting logic.

## Alternatives considered

- Compute weights via Elo/PageRank: adds complexity without clear benefit for workshop timelines; win-count suffices per methodology.
- Store anonymous scores only: conflicts with auditability; we’ll store identities but gate exposure per privacy rules.
- gRPC for real-time: existing WS is sufficient and simpler for browser clients.

## Risks and mitigations

- Toolchain mismatch (Go 1.23.x): run tests and builds via Dockerfile/docker-compose; pin Go image.
- Concurrency on scoring/votes: wrap writes in transactions; add unique constraints to prevent duplicates; idempotent retries.
- Performance: paginate and stream exports; precompute results and cache latest per project.
