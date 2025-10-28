# PairWise

A web application for feature prioritization through group consensus using pairwise comparison methodology.

## Overview

The PairWise model provides a structured approach to feature prioritization by combining:

- Pairwise comparisons for Value and Complexity criteria
- Win-count weighting to establish relative standings
- Fibonacci scoring for absolute magnitude assessment
- Mathematical calculation of Final Priority Scores (FPS)

## Technology Stack

- **Backend**: Go 1.23.3 with Gin web framework
- **Frontend**: React with Vite build system
- **Database**: PostgreSQL
- **Real-time**: WebSocket for collaborative features
- **Deployment**: Docker with multi-stage builds
- **Monitoring**: Prometheus/Grafana integration
- **Testing**: Comprehensive Go and React test suites

## Quick Start

### Local Development

```bash
# Clone and enter directory
git clone <repository-url>
cd pairwise

# Install dependencies and run
go mod tidy
go run .
```

### Docker Deployment

```bash
# Start with Docker Compose
docker-compose up -d

# Or build and run manually
docker build -t pairwise-backend .
docker run -p 8080:8080 pairwise-backend
```

Visit `http://localhost:8080` to access the application.

## PairWise Methodology

### The Formula

```
FPS = (SValue × WValue) / (SComplexity × WComplexity)
```

Where:

- **SValue**: Fibonacci score for business value (1,2,3,5,8,13,21,34,55,89)
- **WValue**: Win-count weight from pairwise value comparisons
- **SComplexity**: Fibonacci score for implementation complexity
- **WComplexity**: Win-count weight from pairwise complexity comparisons

### Process Flow

1. **Setup** → Define project scope and team
2. **Attendees** → Add team members and assign facilitator
3. **Features** → Define features to prioritize
4. **Pairwise Value** → Compare features for business value
5. **Pairwise Complexity** → Compare features for implementation difficulty
6. **Fibonacci Value** → Score absolute value magnitude
7. **Fibonacci Complexity** → Score absolute complexity magnitude
8. **Results** → Calculate and review Final Priority Scores

## Documentation

### For Users

- **[User Manual](docs/USER_MANUAL.md)** - Complete guide to using PairWise
- **[Deployment Guide](docs/DEPLOYMENT_GUIDE.md)** - Production deployment instructions
- **[API Documentation](docs/API_DOCUMENTATION.md)** - REST API and WebSocket reference

### For Developers

- **[Architecture Documentation](doc/)** - Technical implementation details
- **[Feature Management](doc/FEATURE_MANAGEMENT_SUMMARY.md)** - Feature workflow documentation
- **[React Frontend](doc/REACT_FRONTEND_SUMMARY.md)** - Frontend architecture guide

## Development

### Project Structure

```
pairwise/
├── cmd/
│   ├── server/             # Application entry point
│   └── migrate/            # Database migration tool
├── internal/
│   ├── api/                # REST handlers, WebSocket, middleware
│   ├── domain/             # Core PairWise business logic
│   ├── repository/         # Data persistence layer
│   ├── service/            # Business services
│   └── websocket/          # Real-time collaboration
├── web/                    # React frontend application
├── migrations/             # Database schema migrations
├── monitoring/             # Prometheus/Grafana configuration
├── docs/                   # Comprehensive documentation
└── pkg/                    # Shared utilities
```

### Development Commands

```bash
# Development
go mod tidy                 # Install dependencies
go run .                    # Run application
go test ./...              # Run backend tests
cd web && npm test         # Run frontend tests

# Production
go build -o pairwise .     # Build binary
docker build -t pairwise . # Build container
```

### Testing

- **Backend**: Go standard testing with integration tests
- **Frontend**: React Testing Library with Vitest
- **API**: Comprehensive endpoint and workflow testing
- **Coverage**: Full test coverage across all components

## Features

### Core PairWise Functionality

- [x] Project and attendee management
- [x] Feature definition with bulk import
- [x] Pairwise comparison workflows (Value & Complexity)
- [x] Fibonacci scoring interface with consensus tracking
- [x] Real-time collaborative voting via WebSocket
- [x] Automatic priority calculation and ranking
- [x] Results visualization and CSV export

### Production Features

- [x] Session state persistence and recovery
- [x] Comprehensive error handling and validation
- [x] Workflow navigation with progress tracking
- [x] Health monitoring and structured logging
- [x] Docker containerization with monitoring stack
- [x] Rate limiting and security middleware
- [x] Database migrations and connection pooling

## Deployment Options

### Local Development

- Direct Go execution with file-based configuration
- Hot-reload development server for frontend

### Docker (Recommended)

- Multi-stage builds for optimized images
- Docker Compose with PostgreSQL and Redis
- Monitoring stack with Prometheus/Grafana

### Production Platforms

- **AWS ECS**: Container orchestration with RDS
- **Kubernetes**: Scalable deployment with Helm charts
- **Traditional Servers**: Systemd service deployment
- **Cloud Platforms**: Heroku, DigitalOcean, etc.

## Monitoring and Observability

- **Health Checks**: `/health` endpoint for load balancer integration
- **Metrics**: Prometheus metrics for performance monitoring
- **Logging**: Structured JSON logs with request tracing
- **Alerts**: Configurable alerts for system health
- **Dashboards**: Grafana dashboards for visualization

## Security

- Input validation and sanitization
- Rate limiting to prevent abuse
- Secure headers and CORS configuration
- Database connection security
- Docker security best practices

## Contributing

1. **Code Standards**: Follow Go best practices and formatting (`go fmt`)
2. **Testing**: Add comprehensive tests for new functionality
3. **Documentation**: Update relevant documentation
4. **Review Process**: Ensure all tests pass and code review approval
5. **Security**: Follow secure coding practices and validate inputs

## Support

- **Issues**: Report bugs and feature requests via GitHub Issues
- **Documentation**: Comprehensive guides in `docs/` directory
- **API Reference**: Complete endpoint documentation available
- **Community**: Contribute improvements and share experiences

## License

[Add your license information here]
