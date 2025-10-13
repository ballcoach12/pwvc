package domain

import (
	"time"
)

// Attendee represents a project attendee
type Attendee struct {
	ID            int       `json:"id" db:"id"`
	ProjectID     int       `json:"project_id" db:"project_id"`
	Name          string    `json:"name" db:"name" binding:"required,min=1,max=255"`
	Role          string    `json:"role" db:"role"`
	IsFacilitator bool      `json:"is_facilitator" db:"is_facilitator"`
	Email         string    `json:"email,omitempty" db:"email"`          // Optional for now
	PIN           string    `json:"-" db:"pin"`                          // Hidden in JSON, for simple auth
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// CreateAttendeeRequest represents the request payload for creating an attendee
type CreateAttendeeRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=255"`
	Role          string `json:"role"`
	IsFacilitator bool   `json:"is_facilitator"`
}
