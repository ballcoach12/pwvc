package domain

import (
	"time"
)

// Role represents permission levels and access groups for users
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50;not null" binding:"required,min=2,max=50"`
	Description string    `json:"description,omitempty" gorm:"size:255"`
	Permissions string    `json:"permissions,omitempty" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Note: Users relationship managed through UserRole junction table
}

// CreateRoleRequest represents the request payload for creating a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description,omitempty" binding:"max=255"`
	Permissions string `json:"permissions,omitempty"`
}

// UpdateRoleRequest represents the request payload for updating a role
type UpdateRoleRequest struct {
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Permissions *string `json:"permissions,omitempty"`
}

// PredefinedRoles contains the system's default roles
var PredefinedRoles = []Role{
	{
		Name:        "admin",
		Description: "Administrator with full system access",
		Permissions: `{"users":["read","write","delete"],"projects":["read","write","delete"],"admin":["read","write"]}`,
	},
	{
		Name:        "user",
		Description: "Standard user with project participation rights",
		Permissions: `{"projects":["read","write"],"features":["read","write"],"pairwise":["read","write"]}`,
	},
	{
		Name:        "viewer",
		Description: "Read-only access to view projects and results",
		Permissions: `{"projects":["read"],"features":["read"],"results":["read"]}`,
	},
}

// IsSystemRole checks if the role is a predefined system role
func (r *Role) IsSystemRole() bool {
	systemRoles := []string{"admin", "user", "viewer"}
	for _, sysRole := range systemRoles {
		if r.Name == sysRole {
			return true
		}
	}
	return false
}
