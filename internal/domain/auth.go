package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the JWT token claims structure
type JWTClaims struct {
	UserID   uint     `json:"sub"`      // Subject - user ID
	Username string   `json:"username"` // Username for display purposes
	Roles    []string `json:"roles"`    // Array of role names for authorization
	jwt.RegisteredClaims
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the successful login response
type LoginResponse struct {
	Message string      `json:"message"`
	User    UserProfile `json:"user"`
}

// UserProfile represents the user profile returned after authentication
type UserProfile struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email,omitempty"`
	IsActive bool     `json:"is_active"`
	Roles    []string `json:"roles"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
