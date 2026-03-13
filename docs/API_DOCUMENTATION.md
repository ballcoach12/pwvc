# PairWise API Documentation

## Overview

The PairWise API provides RESTful endpoints for managing projects, attendees, features, pairwise comparisons, and priority calculations following the Pairwise-Weighted Value/Complexity methodology.

## Base URL

```
Production: https://your-domain.com/api
Development: http://localhost:8080/api
```

## Authentication

Currently, the API does not require authentication. Future versions may implement JWT-based authentication.

## Response Format

All API responses follow a consistent JSON structure:

### Success Response

```json
{
  "data": { ... },
  "timestamp": "2023-12-07T10:30:00Z"
}
```

### Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "name",
        "message": "Name is required"
      }
    ]
  },
  "timestamp": "2023-12-07T10:30:00Z"
}
```

## Error Codes

- `400` - Bad Request (validation errors, malformed JSON)
- `404` - Not Found (resource not found)
- `409` - Conflict (resource already exists)
- `500` - Internal Server Error

## Endpoints

### Health Check

Check API health and status.

#### GET /health

```http
GET /api/health
```

**Response:**

```json
{
  "status": "ok",
  "timestamp": "2023-12-07T10:30:00Z",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected"
}
```

---

## Projects

### List All Projects

Retrieve all projects in the system.

#### GET /projects

```http
GET /api/projects
```

**Response:**

```json
[
  {
    "id": 1,
    "name": "Mobile App Redesign",
    "description": "Redesign of the mobile application interface",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
]
```

### Get Project by ID

Retrieve a specific project by its ID.

#### GET /projects/{id}

```http
GET /api/projects/1
```

**Response:**

```json
{
  "id": 1,
  "name": "Mobile App Redesign",
  "description": "Redesign of the mobile application interface",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

### Create Project

Create a new project.

#### POST /projects

```http
POST /api/projects
Content-Type: application/json

{
  "name": "Mobile App Redesign",
  "description": "Redesign of the mobile application interface"
}
```

**Validation Rules:**

- `name`: Required, 1-255 characters
- `description`: Optional, max 5000 characters

**Response:**

```json
{
  "id": 1,
  "name": "Mobile App Redesign",
  "description": "Redesign of the mobile application interface",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

### Update Project

Update an existing project.

#### PUT /projects/{id}

```http
PUT /api/projects/1
Content-Type: application/json

{
  "name": "Updated Mobile App Redesign",
  "description": "Updated description for the mobile application interface"
}
```

**Response:**

```json
{
  "id": 1,
  "name": "Updated Mobile App Redesign",
  "description": "Updated description for the mobile application interface",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-07T10:30:00Z"
}
```

### Delete Project

Delete a project and all associated data.

#### DELETE /projects/{id}

```http
DELETE /api/projects/1
```

**Response:** `204 No Content`

---

## Attendees

### List Project Attendees

Get all attendees for a specific project.

#### GET /projects/{projectId}/attendees

```http
GET /api/projects/1/attendees
```

**Response:**

```json
[
  {
    "id": 1,
    "project_id": 1,
    "name": "John Doe",
    "role": "Product Manager",
    "is_facilitator": true,
    "created_at": "2023-12-01T10:00:00Z"
  },
  {
    "id": 2,
    "project_id": 1,
    "name": "Jane Smith",
    "role": "UX Designer",
    "is_facilitator": false,
    "created_at": "2023-12-01T10:15:00Z"
  }
]
```

### Create Attendee

Add a new attendee to a project.

#### POST /projects/{projectId}/attendees

```http
POST /api/projects/1/attendees
Content-Type: application/json

{
  "name": "John Doe",
  "role": "Product Manager",
  "is_facilitator": true
}
```

**Validation Rules:**

- `name`: Required, 1-255 characters
- `role`: Optional, max 255 characters
- `is_facilitator`: Boolean, default false

**Business Rules:**

- Each project can have only one facilitator
- Minimum 2 attendees required for PairWise methodology

**Response:**

```json
{
  "id": 1,
  "project_id": 1,
  "name": "John Doe",
  "role": "Product Manager",
  "is_facilitator": true,
  "created_at": "2023-12-01T10:00:00Z"
}
```

### Delete Attendee

Remove an attendee from a project.

#### DELETE /projects/{projectId}/attendees/{attendeeId}

```http
DELETE /api/projects/1/attendees/1
```

**Response:** `204 No Content`

---

## Features

### List Project Features

Get all features for a specific project.

#### GET /projects/{projectId}/features

```http
GET /api/projects/1/features
```

**Response:**

```json
[
  {
    "id": 1,
    "project_id": 1,
    "title": "User Authentication",
    "description": "Implement secure user login and registration system",
    "acceptance_criteria": "Users can register, login, logout, and reset passwords",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
]
```

### Get Feature by ID

Retrieve a specific feature.

#### GET /projects/{projectId}/features/{featureId}

```http
GET /api/projects/1/features/1
```

**Response:**

```json
{
  "id": 1,
  "project_id": 1,
  "title": "User Authentication",
  "description": "Implement secure user login and registration system",
  "acceptance_criteria": "Users can register, login, logout, and reset passwords",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

### Create Feature

Add a new feature to a project.

#### POST /projects/{projectId}/features

```http
POST /api/projects/1/features
Content-Type: application/json

{
  "title": "User Authentication",
  "description": "Implement secure user login and registration system",
  "acceptance_criteria": "Users can register, login, logout, and reset passwords"
}
```

**Validation Rules:**

- `title`: Required, 1-255 characters
- `description`: Required, 1-5000 characters
- `acceptance_criteria`: Optional, max 5000 characters

**Business Rules:**

- Minimum 2 features required for PairWise methodology
- Feature titles must be unique within a project

**Response:**

```json
{
  "id": 1,
  "project_id": 1,
  "title": "User Authentication",
  "description": "Implement secure user login and registration system",
  "acceptance_criteria": "Users can register, login, logout, and reset passwords",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

### Update Feature

Update an existing feature.

#### PUT /projects/{projectId}/features/{featureId}

```http
PUT /api/projects/1/features/1
Content-Type: application/json

{
  "title": "Enhanced User Authentication",
  "description": "Implement secure user login and registration system with 2FA",
  "acceptance_criteria": "Users can register, login with 2FA, logout, and reset passwords"
}
```

**Response:**

```json
{
  "id": 1,
  "project_id": 1,
  "title": "Enhanced User Authentication",
  "description": "Implement secure user login and registration system with 2FA",
  "acceptance_criteria": "Users can register, login with 2FA, logout, and reset passwords",
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-07T10:30:00Z"
}
```

### Delete Feature

Remove a feature from a project.

#### DELETE /projects/{projectId}/features/{featureId}

```http
DELETE /api/projects/1/features/1
```

**Response:** `204 No Content`

### Import Features

Bulk import features from CSV.

#### POST /projects/{projectId}/features/import

```http
POST /api/projects/1/features/import
Content-Type: multipart/form-data

file: features.csv
```

**CSV Format:**

```csv
title,description,acceptance_criteria
"User Authentication","Login and registration system","Users can register and login"
"Payment Processing","Credit card payments","Users can pay with Stripe"
```

**Response:**

```json
{
  "imported": 2,
  "errors": [],
  "features": [
    {
      "id": 1,
      "title": "User Authentication",
      "description": "Login and registration system",
      "acceptance_criteria": "Users can register and login"
    }
  ]
}
```

### Export Features

Export features to CSV.

#### GET /projects/{projectId}/features/export

```http
GET /api/projects/1/features/export
Accept: text/csv
```

**Response:**

```csv
id,title,description,acceptance_criteria,created_at
1,"User Authentication","Login and registration system","Users can register and login","2023-12-01T10:00:00Z"
```

---

## Pairwise Comparisons

### Start Pairwise Session

Begin a new pairwise comparison session.

#### POST /projects/{projectId}/pairwise/sessions

```http
POST /api/projects/1/pairwise/sessions
Content-Type: application/json

{
  "criterion_type": "value"
}
```

**Parameters:**

- `criterion_type`: "value" or "complexity"

**Response:**

```json
{
  "id": 1,
  "project_id": 1,
  "criterion_type": "value",
  "status": "active",
  "started_at": "2023-12-07T10:30:00Z",
  "completed_at": null
}
```

### Get Active Session

Retrieve the current active session for a project.

#### GET /projects/{projectId}/pairwise/active

```http
GET /api/projects/1/pairwise/active
```

**Response:**

```json
{
  "session": {
    "id": 1,
    "project_id": 1,
    "criterion_type": "value",
    "status": "active",
    "started_at": "2023-12-07T10:30:00Z"
  },
  "progress": {
    "total_comparisons": 6,
    "completed_comparisons": 2,
    "progress_percentage": 33.33,
    "remaining_comparisons": 4
  },
  "current_comparison": {
    "id": 3,
    "session_id": 1,
    "feature_a": {
      "id": 1,
      "title": "User Authentication"
    },
    "feature_b": {
      "id": 2,
      "title": "Payment Processing"
    },
    "votes": []
  }
}
```

### Submit Vote

Submit a vote for a comparison.

#### POST /projects/{projectId}/pairwise/vote

```http
POST /api/projects/1/pairwise/vote
Content-Type: application/json

{
  "comparison_id": 1,
  "attendee_id": 1,
  "preferred_feature_id": 2,
  "is_tie_vote": false
}
```

**Parameters:**

- `comparison_id`: Required, ID of the comparison
- `attendee_id`: Required, ID of the voting attendee
- `preferred_feature_id`: Required for preference votes, null for ties
- `is_tie_vote`: Boolean, true if vote is a tie

**Response:**

```json
{
  "vote": {
    "id": 1,
    "comparison_id": 1,
    "attendee_id": 1,
    "preferred_feature_id": 2,
    "is_tie_vote": false,
    "voted_at": "2023-12-07T10:30:00Z"
  },
  "comparison_complete": true,
  "consensus_reached": true,
  "winner_id": 2
}
```

### Get Session Results

Get results for a completed session.

#### GET /projects/{projectId}/pairwise/sessions/{sessionId}/results

```http
GET /api/projects/1/pairwise/sessions/1/results
```

**Response:**

```json
{
  "session_id": 1,
  "criterion_type": "value",
  "completed_at": "2023-12-07T11:00:00Z",
  "comparisons": [
    {
      "feature_a_id": 1,
      "feature_b_id": 2,
      "winner_id": 2,
      "consensus_reached": true,
      "votes": [
        {
          "attendee_id": 1,
          "preferred_feature_id": 2
        }
      ]
    }
  ],
  "win_counts": [
    {
      "feature_id": 1,
      "wins": 2,
      "total_comparisons": 3,
      "win_percentage": 66.67
    },
    {
      "feature_id": 2,
      "wins": 1,
      "total_comparisons": 3,
      "win_percentage": 33.33
    }
  ]
}
```

---

## Fibonacci Scoring

### Get Fibonacci Scores

Retrieve current Fibonacci scores for features.

#### GET /projects/{projectId}/fibonacci

```http
GET /api/projects/1/fibonacci?score_type=value
```

**Query Parameters:**

- `score_type`: "value" or "complexity"
- `attendee_id`: Optional, filter by attendee

**Response:**

```json
{
  "scores": [
    {
      "feature_id": 1,
      "feature_title": "User Authentication",
      "scores": [
        {
          "attendee_id": 1,
          "attendee_name": "John Doe",
          "score": 8,
          "score_type": "value"
        }
      ],
      "consensus_score": 8,
      "consensus_reached": true
    }
  ],
  "fibonacci_values": [1, 2, 3, 5, 8, 13, 21, 34, 55, 89]
}
```

### Submit Fibonacci Score

Submit a Fibonacci score for a feature.

#### POST /projects/{projectId}/fibonacci

```http
POST /api/projects/1/fibonacci
Content-Type: application/json

{
  "feature_id": 1,
  "attendee_id": 1,
  "score": 8,
  "score_type": "value"
}
```

**Validation Rules:**

- `score`: Must be a valid Fibonacci number (1, 2, 3, 5, 8, 13, 21, 34, 55, 89)
- `score_type`: Must be "value" or "complexity"

**Response:**

```json
{
  "id": 1,
  "feature_id": 1,
  "attendee_id": 1,
  "score": 8,
  "score_type": "value",
  "created_at": "2023-12-07T10:30:00Z"
}
```

---

## Priority Calculations

### Calculate Project Priorities

Calculate final PairWise priorities for all features.

#### POST /projects/{projectId}/calculate

```http
POST /api/projects/1/calculate
```

**Response:**

```json
{
  "calculation_id": "calc_123456",
  "project_id": 1,
  "calculated_at": "2023-12-07T11:00:00Z",
  "results": [
    {
      "feature_id": 1,
      "feature_title": "User Authentication",
      "w_value": 0.67,
      "w_complexity": 0.33,
      "s_value": 8,
      "s_complexity": 5,
      "weighted_value": 5.36,
      "weighted_complexity": 1.65,
      "final_priority_score": 3.25,
      "rank": 1
    }
  ],
  "summary": {
    "total_features": 3,
    "highest_score": 3.25,
    "lowest_score": 0.15,
    "average_score": 1.85,
    "median_score": 1.9
  }
}
```

### Get Project Results

Retrieve the latest priority calculation results.

#### GET /projects/{projectId}/results

```http
GET /api/projects/1/results
```

**Response:**

```json
{
  "project_id": 1,
  "project_name": "Mobile App Redesign",
  "calculated_at": "2023-12-07T11:00:00Z",
  "results": [
    {
      "feature": {
        "id": 1,
        "title": "User Authentication",
        "description": "Secure login system"
      },
      "w_value": 0.67,
      "w_complexity": 0.33,
      "s_value": 8,
      "s_complexity": 5,
      "weighted_value": 5.36,
      "weighted_complexity": 1.65,
      "final_priority_score": 3.25,
      "rank": 1
    }
  ]
}
```

---

## Project Progress

### Get Project Progress

Retrieve current workflow progress for a project.

#### GET /projects/{projectId}/progress

```http
GET /api/projects/1/progress
```

**Response:**

```json
{
  "project_id": 1,
  "current_phase": "pairwise_value",
  "setup_completed": true,
  "attendees_added": true,
  "features_added": true,
  "pairwise_value_completed": false,
  "pairwise_complexity_completed": false,
  "fibonacci_value_completed": false,
  "fibonacci_complexity_completed": false,
  "results_calculated": false,
  "progress_percentage": 37.5,
  "available_phases": ["setup", "attendees", "features", "pairwise_value"],
  "next_phase": "pairwise_value"
}
```

### Update Progress Phase

Mark a workflow phase as completed.

#### POST /projects/{projectId}/progress/{phase}/complete

```http
POST /api/projects/1/progress/setup/complete
```

**Available Phases:**

- `setup`
- `attendees`
- `features`
- `pairwise_value`
- `pairwise_complexity`
- `fibonacci_value`
- `fibonacci_complexity`
- `results`

**Response:**

```json
{
  "project_id": 1,
  "phase": "setup",
  "completed": true,
  "next_phase": "attendees",
  "updated_at": "2023-12-07T10:30:00Z"
}
```

---

## WebSocket Events

The API supports real-time updates via WebSocket connections for collaborative features.

### Connection

```javascript
const ws = new WebSocket("ws://localhost:8080/ws");
```

### Event Types

#### Join Project Room

```json
{
  "type": "join_project",
  "project_id": 1,
  "attendee_id": 1
}
```

#### Pairwise Vote Submitted

```json
{
  "type": "vote_submitted",
  "project_id": 1,
  "comparison_id": 1,
  "voter_id": 1,
  "preferred_feature_id": 2
}
```

#### Comparison Completed

```json
{
  "type": "comparison_completed",
  "project_id": 1,
  "comparison_id": 1,
  "winner_id": 2,
  "consensus_reached": true
}
```

#### Phase Completed

```json
{
  "type": "phase_completed",
  "project_id": 1,
  "phase": "pairwise_value",
  "next_phase": "pairwise_complexity"
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **General endpoints**: 100 requests per minute per IP
- **Calculation endpoints**: 10 requests per minute per IP
- **WebSocket connections**: 50 connections per IP

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1701946800
```

---

## SDK Examples

### JavaScript/Node.js

```javascript
const axios = require("axios");

class PairWiseClient {
  constructor(baseURL = "http://localhost:8080/api") {
    this.client = axios.create({ baseURL });
  }

  // Projects
  async getProjects() {
    const response = await this.client.get("/projects");
    return response.data;
  }

  async createProject(name, description) {
    const response = await this.client.post("/projects", {
      name,
      description,
    });
    return response.data;
  }

  // Features
  async addFeature(projectId, title, description, acceptanceCriteria) {
    const response = await this.client.post(`/projects/${projectId}/features`, {
      title,
      description,
      acceptance_criteria: acceptanceCriteria,
    });
    return response.data;
  }

  // Pairwise Comparisons
  async startPairwiseSession(projectId, criterionType) {
    const response = await this.client.post(
      `/projects/${projectId}/pairwise/sessions`,
      {
        criterion_type: criterionType,
      }
    );
    return response.data;
  }

  async submitVote(
    projectId,
    comparisonId,
    attendeeId,
    preferredFeatureId,
    isTie = false
  ) {
    const response = await this.client.post(
      `/projects/${projectId}/pairwise/vote`,
      {
        comparison_id: comparisonId,
        attendee_id: attendeeId,
        preferred_feature_id: preferredFeatureId,
        is_tie_vote: isTie,
      }
    );
    return response.data;
  }
}

// Usage
const client = new PairWiseClient();
const projects = await client.getProjects();
```

### Python

```python
import requests

class PairWiseClient:
    def __init__(self, base_url='http://localhost:8080/api'):
        self.base_url = base_url
        self.session = requests.Session()

    def get_projects(self):
        response = self.session.get(f'{self.base_url}/projects')
        response.raise_for_status()
        return response.json()

    def create_project(self, name, description):
        data = {'name': name, 'description': description}
        response = self.session.post(f'{self.base_url}/projects', json=data)
        response.raise_for_status()
        return response.json()

    def calculate_priorities(self, project_id):
        response = self.session.post(f'{self.base_url}/projects/{project_id}/calculate')
        response.raise_for_status()
        return response.json()

# Usage
client = PairWiseClient()
projects = client.get_projects()
```

---

## Error Handling

### Validation Errors

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "name",
        "message": "Name is required and cannot be empty"
      },
      {
        "field": "description",
        "message": "Description cannot exceed 5000 characters"
      }
    ]
  }
}
```

### Business Logic Errors

```json
{
  "error": {
    "code": "BUSINESS_RULE_VIOLATION",
    "message": "Cannot have multiple facilitators in a project",
    "details": {
      "rule": "single_facilitator",
      "current_facilitators": 2,
      "max_allowed": 1
    }
  }
}
```

### Resource Not Found

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Project with ID 999 not found"
  }
}
```

---

## Changelog

### v1.0.0 (2023-12-07)

- Initial API release
- Project, attendee, and feature management
- Pairwise comparison functionality
- Fibonacci scoring system
- Priority calculation engine
- Real-time WebSocket support
- Progress tracking system

---

## Support

For API support and questions:

- **Documentation**: This document
- **Issues**: GitHub repository issues
- **Health Check**: `GET /health` endpoint
- **OpenAPI Spec**: Available at `/api/docs` (if enabled)
