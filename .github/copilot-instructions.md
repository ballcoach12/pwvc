# Copilot Instructions for PairWise

## Project Overview

This is a web application implementing **PairWise** - a structured methodology for objective feature prioritization through group consensus. Built with Go backend (1.23.3) and React frontend.

## PairWise Methodology Core Concepts

- **Pairwise Comparison**: Head-to-head feature comparisons for Value and Complexity criteria
- **Win-Count Weighting**: Mathematical calculation of relative feature standings (W = wins / total comparisons)
- **Fibonacci Scoring**: Absolute magnitude scoring using Fibonacci sequence (1,2,3,5,8...)
- **Final Priority Score**: FPS = (SValue × WValue) / (SComplexity × WComplexity)
- **Consensus Requirement**: All scores must be agreed upon by the entire team

## Development Setup

- **Go Version**: 1.23.3 (specified in `go.mod`)
- **Module**: `pairwise`
- **Shell**: PowerShell on Windows

## Project Structure

```
pairwise/
├── cmd/server/              # Main application entry point
├── internal/
│   ├── api/                # REST API handlers
│   ├── domain/             # PairWise business logic (feature, comparison, scoring, consensus)
│   ├── repository/         # Data persistence layer
│   ├── service/            # Business services
│   └── websocket/          # Real-time collaboration
├── web/                    # React frontend
│   ├── src/
│   │   ├── components/     # UI components (PairwiseGrid, FibonacciScorer, etc.)
│   │   ├── pages/          # Page components (ProjectSetup, PairwisePhase, etc.)
│   │   ├── hooks/          # Custom hooks (usePairwise, useConsensus, etc.)
│   │   └── services/       # API clients
├── migrations/             # Database schema
└── pkg/                    # Shared utilities
```

## Development Workflow

```powershell
# Initialize and run the project
go mod tidy
go run .

# Build the project
go build -o pairwise.exe .

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code for common issues
go vet ./...
```

## Code Conventions

- Use standard Go formatting (`gofmt`)
- Follow effective Go practices
- Use meaningful package names that describe functionality
- Keep functions focused and testable
- Use context.Context for cancellation and request-scoped values

## Dependencies

- Prefer standard library when possible
- Use `go mod tidy` to manage dependencies
- Document any external dependencies and their purpose in README

## Key Commands

- `go mod init pairwise` - Already done (creates go.mod)
- `go mod tidy` - Clean up dependencies
- `go run .` - Run the main package
- `go build .` - Build executable

## Active Technologies
- Go 1.23.x backend; React + Vite frontend (Node 18+) + gin-gonic/gin, gorilla/websocket, lib/pq, gorm (sqlite in tests); React, Vite, Vitest (001-zero-doubt-prioritization)
- PostgreSQL (prod), SQLite (tests) (001-zero-doubt-prioritization)

## Recent Changes
- 001-zero-doubt-prioritization: Added Go 1.23.x backend; React + Vite frontend (Node 18+) + gin-gonic/gin, gorilla/websocket, lib/pq, gorm (sqlite in tests); React, Vite, Vitest
