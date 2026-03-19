package domain

import (
	"time"
)

// User represents a system user with authentication credentials
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null" binding:"required,min=3,max=50"`
	PasswordHash string    `json:"-" gorm:"size:255"` // Never serialize password hash
	Email        string    `json:"email,omitempty" gorm:"size:100"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Roles []Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Password string   `json:"password" binding:"required,min=6"`
	Email    string   `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Roles    []string `json:"roles,omitempty"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email    *string  `json:"email,omitempty" binding:"omitempty,email,max=100"`
	IsActive *bool    `json:"is_active,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordRequest represents the request payload for admin password reset
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// GetRoleNames returns the role names for the user
func (u *User) GetRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		roleNames[i] = role.Name
	}
	return roleNames
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// ToUserProfile converts User to UserProfile for API responses
func (u *User) ToUserProfile() UserProfile {
	return UserProfile{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		IsActive: u.IsActive,
		Roles:    u.GetRoleNames(),
	}
}
