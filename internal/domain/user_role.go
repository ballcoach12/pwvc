package domain

import (
	"time"
)

// UserRole represents the junction table managing many-to-many relationship between users and roles
type UserRole struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	RoleID     uint      `json:"role_id" gorm:"not null;index"`
	AssignedAt time.Time `json:"assigned_at" gorm:"default:CURRENT_TIMESTAMP"`
	AssignedBy *uint     `json:"assigned_by,omitempty"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// AssignRoleRequest represents the request payload for assigning a role to a user
type AssignRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// RemoveRoleRequest represents the request payload for removing a role from a user
type RemoveRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// BeforeCreate sets up constraints and validation before creating UserRole
func (ur *UserRole) BeforeCreate(tx *interface{}) error {
	if ur.AssignedAt.IsZero() {
		ur.AssignedAt = time.Now()
	}
	return nil
}
