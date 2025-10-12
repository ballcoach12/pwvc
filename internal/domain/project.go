package domain

import (
	"time"
)

// Project represents a P-WVC project
type Project struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required,min=1,max=255"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateProjectRequest represents the request payload for creating a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

// UpdateProjectRequest represents the request payload for updating a project
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive completed"`
}
