# P-WVC (Pairwise-Weighted Value/Complexity) Model API

A Go backend service implementing the P-WVC methodology for objective feature prioritization through group consensus.

## Project Structure

```
pwvc/
├── cmd/
│   ├── server/          # Main application entry point
│   └── migrate/         # Database migration utility
├── internal/
│   ├── api/            # REST API handlers
│   ├── domain/         # Business entities and errors
│   ├── repository/     # Data persistence layer
│   ├── service/        # Business logic layer
│   └── websocket/      # Real-time collaboration (future)
├── migrations/         # Database schema migrations
├── pkg/               # Shared utilities
└── web/               # React frontend (future)
```

## Prerequisites

- Go 1.23.3 or later
- PostgreSQL 12+
- Git

## Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd pwvc
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**

   ```bash
   # Create database
   createdb pwvc

   # Or use Docker
   docker run --name pwvc-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=pwvc -p 5432:5432 -d postgres:15
   ```

4. **Set environment variables** (optional)

   ```bash
   export DATABASE_URL="postgres://postgres:password@localhost:5432/pwvc?sslmode=disable"
   export PORT="8080"
   export GIN_MODE="debug"  # or "release" for production
   ```

5. **Run database migrations**

   ```bash
   go run cmd/migrate/main.go up
   ```

6. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on port 8080 (or the PORT environment variable).

## API Endpoints

### Projects

- `POST /api/projects` - Create new project
- `GET /api/projects/{id}` - Get project details
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project
- `GET /api/projects` - List all projects

### Attendees

- `POST /api/projects/{id}/attendees` - Add attendee to project
- `GET /api/projects/{id}/attendees` - List project attendees
- `DELETE /api/projects/{id}/attendees/{attendee_id}` - Remove attendee

### Features

- `POST /api/projects/{id}/features` - Create new feature
- `GET /api/projects/{id}/features` - List all features in project
- `GET /api/projects/{id}/features/{feature_id}` - Get specific feature
- `PUT /api/projects/{id}/features/{feature_id}` - Update feature
- `DELETE /api/projects/{id}/features/{feature_id}` - Delete feature
- `POST /api/projects/{id}/features/import` - Import features from CSV
- `GET /api/projects/{id}/features/export` - Export features to CSV

### Health Check

- `GET /health` - Service health check
- `GET /` - API information

## Example Usage

### Create a Project

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Website Redesign",
    "description": "Modernize company website with new features"
  }'
```

### Add Attendee

```bash
curl -X POST http://localhost:8080/api/projects/1/attendees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "role": "Product Manager",
    "is_facilitator": true
  }'
```

### Create Feature

```bash
curl -X POST http://localhost:8080/api/projects/1/features \
  -H "Content-Type: application/json" \
  -d '{
    "title": "User Authentication",
    "description": "Implement secure user login and registration functionality",
    "acceptance_criteria": "Users should be able to login with email/password and receive appropriate error messages for invalid credentials"
  }'
```

### Import Features from CSV

```bash
curl -X POST http://localhost:8080/api/projects/1/features/import \
  -F "file=@features.csv"
```

### Export Features to CSV

```bash
curl -X GET http://localhost:8080/api/projects/1/features/export \
  -o "project_features.csv"
```

## Database Schema

### Projects Table

```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Attendees Table

```sql
CREATE TABLE attendees (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100),
    is_facilitator BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Features Table

```sql
CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    acceptance_criteria TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## CSV Import/Export Format

### CSV Structure for Features

```csv
title,description,acceptance_criteria
"User Login","Users can authenticate with email/password","Given valid credentials, user should be logged in successfully"
"Dashboard View","Display key metrics and navigation","Dashboard loads within 2 seconds and shows current data"
```

### Validation Rules

- **Title**: Required, maximum 255 characters
- **Description**: Required, maximum 5000 characters
- **Acceptance Criteria**: Optional, maximum 5000 characters

### CSV Import Process

1. Upload CSV file via `POST /api/projects/{id}/features/import`
2. System validates each row according to rules above
3. Valid features are imported, invalid rows are reported in response
4. Response includes counts of imported/skipped records and error details

## Development

### Build

```bash
go build -o pwvc ./cmd/server
```

### Run Tests

```bash
go test ./...
```

### Format Code

```bash
go fmt ./...
```

### Lint Code

```bash
go vet ./...
```

## Migration Commands

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback all migrations
go run cmd/migrate/main.go down

# Check current migration version
go run cmd/migrate/main.go version

# Force specific version (use with caution)
go run cmd/migrate/main.go force 1
```

## Environment Variables

| Variable       | Default                                                            | Description                             |
| -------------- | ------------------------------------------------------------------ | --------------------------------------- |
| `DATABASE_URL` | `postgres://postgres:password@localhost:5432/pwvc?sslmode=disable` | PostgreSQL connection string            |
| `PORT`         | `8080`                                                             | HTTP server port                        |
| `GIN_MODE`     | `debug`                                                            | Gin framework mode (`debug`, `release`) |

## Next Steps

This foundation provides:

- ✅ Basic project and attendee management
- ✅ PostgreSQL integration with migrations
- ✅ RESTful API with proper error handling
- ✅ Feature management with CRUD operations
- ✅ CSV import/export for bulk feature operations
- ✅ Clean architecture (domain, service, repository layers)
- ✅ CORS support for frontend integration
- ✅ Comprehensive validation and error handling

Future development will add:

- Pairwise comparison engine (core P-WVC methodology)
- Fibonacci scoring system
- Win-count weighting calculations
- WebSocket real-time collaboration
- React frontend interface
