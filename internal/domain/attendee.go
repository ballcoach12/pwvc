package domain

import (
	"time"
)

// Attendee represents a project attendee
type Attendee struct {
	ID                   int        `json:"id" db:"id"`
	ProjectID            int        `json:"project_id" db:"project_id"`
	Name                 string     `json:"name" db:"name" binding:"required,min=1,max=255"`
	Role                 string     `json:"role" db:"role"`
	IsFacilitator        bool       `json:"is_facilitator" db:"is_facilitator"`
	Email                string     `json:"email,omitempty" db:"email"`     // Optional for now
	PinHash              *string    `json:"-" db:"pin_hash"`                // Hidden in JSON, stores hashed PIN (nullable)
	InviteToken          *string    `json:"-" db:"invite_token"`            // Hidden in JSON, for PIN setup (nullable)
	InviteTokenExpiresAt *time.Time `json:"-" db:"invite_token_expires_at"` // Hidden in JSON, token expiry (nullable)
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

// CreateAttendeeRequest represents the request payload for creating an attendee
type CreateAttendeeRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=255"`
	Role          string `json:"role"`
	IsFacilitator bool   `json:"is_facilitator"`
	PIN           string `json:"pin" binding:"required,min=4,max=20"` // PIN for authentication
}

// SetPINRequest represents the request payload for setting/updating an attendee's PIN
type SetPINRequest struct {
	PIN string `json:"pin" binding:"required,min=4,max=20"`
}

// CreateAttendeeWithoutPINRequest represents creating an attendee without requiring a PIN upfront
type CreateAttendeeWithoutPINRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=255"`
	Role          string `json:"role"`
	IsFacilitator bool   `json:"is_facilitator"`
}

// SetupPINRequest represents setting up a PIN using an invite token
type SetupPINRequest struct {
	AttendeeID  int    `json:"attendee_id" binding:"required"`
	InviteToken string `json:"invite_token" binding:"required"`
	PIN         string `json:"pin" binding:"required,min=4,max=20"`
}

// ChangePINRequest represents changing your own PIN when authenticated
type ChangePINRequest struct {
	CurrentPIN string `json:"current_pin" binding:"required"`
	NewPIN     string `json:"new_pin" binding:"required,min=4,max=20"`
}

// ResetPINResponse represents the response when generating a PIN reset
type ResetPINResponse struct {
	InviteToken string `json:"invite_token"`
	ExpiresAt   string `json:"expires_at"`
}
