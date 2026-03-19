# Implementation Plan: Zero Doubt Prioritization (P‑WVC)

Branch: 001-zero-doubt-prioritization | Date: 2025-10-26 | Spec: specs/001-zero-doubt-prioritization/spec.md
Input: Feature specification from specs/001-zero-doubt-prioritization/spec.md

Note: Generated via speckit.plan workflow tailored to this repo.

## Summary

Deliver an end-to-end PairWise Value/Complexity (P‑WVC) flow: pairwise comparisons for value and complexity (WValue/WComplexity), Fibonacci absolute scoring with consensus (SValue/SComplexity), Final Priority Score calculation and export, plus real‑time collaboration and facilitator controls. The backend already contains most scaffolding (projects, features, pairwise sessions, results endpoints, websockets). Key gaps to close: compute real win‑count weights from votes (not mocks), implement Fibonacci scoring APIs + consensus, add audit logging with privacy rules, apply deterministic tie‑breaks, and wire facilitator progress/controls. Frontend wiring follows in a later phase.

## Technical Context

**Language/Version**: Go 1.23.x backend; React + Vite frontend (Node 18+)  
**Primary Dependencies**: gin-gonic/gin, gorilla/websocket, lib/pq, gorm (sqlite in tests); React, Vite, Vitest  
**Storage**: PostgreSQL (prod), SQLite (tests)  
**Testing**: go test (backend); vitest (frontend)  
**Target Platform**: Linux containers; local dev via Docker/Docker Compose  
**Project Type**: Web application (Go API + React SPA)  
**Performance Goals**: Real-time updates ≤ 500 ms p95 in-session; results calc < 2s for ≤ 100 features  
**Constraints**: 99.9% availability target; privacy/anonymity per session; WCAG AA for critical flows  
**Scale/Scope**: 3–25 participants; 10–100 features typical per session

## Constitution Check

Gate status before Phase 0:
- Testing/TDD: PARTIAL. Tests exist but local go toolchain mismatch blocked a run (go directive 1.23.3). Use Docker toolchain for CI; ensure green. Action: run tests in container and fix failures.
- A11y: NOT EVALUATED here (frontend follow-up). Ensure keyboard flows and WCAG AA on pairwise/scoring/results.
- Observability: PARTIAL. Structured API logging exists; metrics/tracing absent. Action: add request metrics for hot paths and WS connection counts.
- Contract-first: PARTIAL. Endpoints exist without OpenAPI. Action: author OpenAPI in contracts/ and backfill docs/tests.
- Performance: PARTIAL. No perf tests or budgets. Action: define budgets (500 ms p95 WS updates) and add lightweight checks.

Re-check will occur after Phase 1 design artifacts are generated.

## Project Structure

Documentation (this feature)

```text
specs/001-zero-doubt-prioritization/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
├── traceability.md
└── gap-analysis.md
```

Source Code (repository root)

```text
cmd/
├── server/            # API entrypoint
└── migrate/           # migrations tool

internal/
├── api/               # HTTP handlers (features, pairwise, results, progress, auth, websocket)
├── domain/            # PWVC math (Fibonacci, win-count, FPS), entities, errors
├── repository/        # persistence (projects, features, pairwise, progress, priority)
├── service/           # business logic (pairwise, results, progress)
└── websocket/         # hub, client, messages

migrations/            # DB schema
web/                   # React app
```

Structure Decision: Keep current monorepo. Backend under internal/* with clean layering. Frontend in web/. Feature docs and contracts under specs/001-zero-doubt-prioritization/.

## Complexity Tracking

Only if needed; currently no justified exceptions beyond adding OpenAPI and an audit log module.
