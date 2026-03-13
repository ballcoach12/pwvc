# Prompt 1: Project Foundation

Set up the basic Go project structure with main.go, initialize Gin web server, create basic project and attendee management endpoints with PostgreSQL integration. Include database migrations for projects and attendees tables.

## Requirements
- Create `cmd/server/main.go` with Gin web server
- Set up PostgreSQL connection and configuration
- Create basic project CRUD endpoints (POST, GET, PUT, DELETE)
- Create attendee management endpoints within projects
- Initialize database migrations for `projects` and `attendees` tables
- Add proper error handling and JSON responses
- Include basic middleware for CORS and logging

## Database Schema
```sql
-- projects table
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- attendees table  
CREATE TABLE attendees (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100),
    is_facilitator BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## API Endpoints to Create
- `POST /api/projects` - Create new project
- `GET /api/projects/{id}` - Get project details
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project
- `POST /api/projects/{id}/attendees` - Add attendee to project
- `GET /api/projects/{id}/attendees` - List project attendees
- `DELETE /api/projects/{id}/attendees/{attendee_id}` - Remove attendee