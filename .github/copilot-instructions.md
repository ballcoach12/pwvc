# Copilot Instructions for pwvc

## Project Overview
This is a web application implementing the **P-WVC (Pairwise-Weighted Value/Complexity) Model** - a structured methodology for objective feature prioritization through group consensus. Built with Go backend (1.23.3) and React frontend.

## P-WVC Methodology Core Concepts
- **Pairwise Comparison**: Head-to-head feature comparisons for Value and Complexity criteria
- **Win-Count Weighting**: Mathematical calculation of relative feature standings (W = wins / total comparisons)
- **Fibonacci Scoring**: Absolute magnitude scoring using Fibonacci sequence (1,2,3,5,8...)
- **Final Priority Score**: FPS = (SValue × WValue) / (SComplexity × WComplexity)
- **Consensus Requirement**: All scores must be agreed upon by the entire team

## Development Setup
- **Go Version**: 1.23.3 (specified in `go.mod`)
- **Module**: `pwvc`
- **Shell**: PowerShell on Windows

## Project Structure
```
pwvc/
├── cmd/server/              # Main application entry point
├── internal/
│   ├── api/                # REST API handlers
│   ├── domain/             # P-WVC business logic (feature, comparison, scoring, consensus)
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
go build -o pwvc.exe .

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
- `go mod init pwvc` - Already done (creates go.mod)
- `go mod tidy` - Clean up dependencies
- `go run .` - Run the main package
- `go build .` - Build executable