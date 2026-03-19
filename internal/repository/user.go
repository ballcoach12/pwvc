package repository

import (
	"database/sql"
	"pairwise/internal/domain"

	"gorm.io/gorm"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	GetUserByID(id uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
	ListUsers(offset, limit int) ([]*domain.User, error)
	LoadUserRoles(user *domain.User) error
	CountUsers() (int64, error)
}

// userRepository implements UserRepository using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

// DeleteUser soft deletes a user by ID
func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// ListUsers retrieves a paginated list of users
func (r *userRepository) ListUsers(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// LoadUserRoles loads the roles for a user
func (r *userRepository) LoadUserRoles(user *domain.User) error {
	return r.db.Model(user).Association("Roles").Find(&user.Roles)
}

// CountUsers returns the total number of users
func (r *userRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Count(&count).Error
	return count, err
}

// RoleRepository interface defines methods for role data access
type RoleRepository interface {
	GetAllRoles() ([]*domain.Role, error)
	GetRoleByID(id uint) (*domain.Role, error)
	GetRoleByName(name string) (*domain.Role, error)
	CreateRole(role *domain.Role) error
	UpdateRole(role *domain.Role) error
	DeleteRole(id uint) error
	AssignRoleToUser(userID, roleID uint, assignedBy *uint) error
	RemoveRoleFromUser(userID, roleID uint) error
}

// roleRepository implements RoleRepository using GORM
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// GetAllRoles retrieves all roles
func (r *roleRepository) GetAllRoles() ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID retrieves a role by ID
func (r *roleRepository) GetRoleByID(id uint) (*domain.Role, error) {
	var role domain.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName retrieves a role by name
func (r *roleRepository) GetRoleByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole creates a new role
func (r *roleRepository) CreateRole(role *domain.Role) error {
	return r.db.Create(role).Error
}

// UpdateRole updates an existing role
func (r *roleRepository) UpdateRole(role *domain.Role) error {
	return r.db.Save(role).Error
}

// DeleteRole soft deletes a role by ID
func (r *roleRepository) DeleteRole(id uint) error {
	return r.db.Delete(&domain.Role{}, id).Error
}

// AssignRoleToUser assigns a role to a user
func (r *roleRepository) AssignRoleToUser(userID, roleID uint, assignedBy *uint) error {
	userRole := &domain.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
	}
	return r.db.Create(userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (r *roleRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&domain.UserRole{}).Error
}

// Legacy SQL-based repository for backward compatibility if needed
type userRepositorySQL struct {
	db *sql.DB
}

// NewUserRepositorySQL creates a SQL-based user repository
func NewUserRepositorySQL(db *sql.DB) UserRepository {
	return &userRepositorySQL{db: db}
}

// Implement interface methods for SQL version...
// (Skipping implementation for now as we're using GORM)

func (r *userRepositorySQL) GetUserByID(id uint) (*domain.User, error) {
	// SQL implementation would go here
	return nil, nil
}

func (r *userRepositorySQL) GetUserByUsername(username string) (*domain.User, error) {
	// SQL implementation would go here
	return nil, nil
}

func (r *userRepositorySQL) CreateUser(user *domain.User) error {
	// SQL implementation would go here
	return nil
}

func (r *userRepositorySQL) UpdateUser(user *domain.User) error {
	// SQL implementation would go here
	return nil
}

func (r *userRepositorySQL) DeleteUser(id uint) error {
	// SQL implementation would go here
	return nil
}

func (r *userRepositorySQL) ListUsers(offset, limit int) ([]*domain.User, error) {
	// SQL implementation would go here
	return nil, nil
}

func (r *userRepositorySQL) LoadUserRoles(user *domain.User) error {
	// SQL implementation would go here
	return nil
}

func (r *userRepositorySQL) CountUsers() (int64, error) {
	// SQL implementation would go here
	return 0, nil
}
