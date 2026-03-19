package service

import (
	"errors"
	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

var (
	ErrUserAlreadyExists = errors.New("user with this username already exists")
	ErrRoleNotFound      = errors.New("role not found")
	ErrCannotDeleteUser  = errors.New("cannot delete user")
)

// UserService handles user management operations
type UserService struct {
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	authService *AuthService
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		authService: authService,
	}
}

// CreateUser creates a new user with roles
func (s *UserService) CreateUser(request *domain.CreateUserRequest, createdBy uint) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetUserByUsername(request.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Username:     request.Username,
		PasswordHash: hashedPassword,
		Email:        request.Email,
		IsActive:     true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Assign roles if specified
	if len(request.Roles) > 0 {
		err = s.AssignRolesToUser(user.ID, request.Roles, &createdBy)
		if err != nil {
			// TODO: Consider rolling back user creation
			return nil, err
		}
	}

	// Load roles and return
	err = s.userRepo.LoadUserRoles(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID with roles
func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Load roles
	err = s.userRepo.LoadUserRoles(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserWithRoles retrieves a user with their roles loaded
func (s *UserService) GetUserWithRoles(id uint) (*domain.User, error) {
	return s.GetUserByID(id) // Already loads roles
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(id uint, request *domain.UpdateUserRequest, updatedBy uint) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}

	// Save user
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// Update roles if specified
	if request.Roles != nil {
		// Remove all existing roles
		err = s.RemoveAllRolesFromUser(user.ID)
		if err != nil {
			return nil, err
		}

		// Assign new roles
		if len(request.Roles) > 0 {
			err = s.AssignRolesToUser(user.ID, request.Roles, &updatedBy)
			if err != nil {
				return nil, err
			}
		}
	}

	// Reload user with roles
	return s.GetUserByID(id)
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(offset, limit int) ([]*domain.User, int64, error) {
	users, err := s.userRepo.ListUsers(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Load roles for each user
	for _, user := range users {
		err = s.userRepo.LoadUserRoles(user)
		if err != nil {
			return nil, 0, err
		}
	}

	// Get total count
	total, err := s.userRepo.CountUsers()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AssignRolesToUser assigns multiple roles to a user
func (s *UserService) AssignRolesToUser(userID uint, roleNames []string, assignedBy *uint) error {
	for _, roleName := range roleNames {
		// Get role by name
		role, err := s.roleRepo.GetRoleByName(roleName)
		if err != nil {
			return ErrRoleNotFound
		}

		// Assign role
		err = s.roleRepo.AssignRoleToUser(userID, role.ID, assignedBy)
		if err != nil {
			// Continue with other roles even if one fails (might already be assigned)
			continue
		}
	}
	return nil
}

// RemoveRoleFromUser removes a specific role from a user
func (s *UserService) RemoveRoleFromUser(userID uint, roleName string) error {
	role, err := s.roleRepo.GetRoleByName(roleName)
	if err != nil {
		return ErrRoleNotFound
	}

	return s.roleRepo.RemoveRoleFromUser(userID, role.ID)
}

// RemoveAllRolesFromUser removes all roles from a user
func (s *UserService) RemoveAllRolesFromUser(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	err = s.userRepo.LoadUserRoles(user)
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		err = s.roleRepo.RemoveRoleFromUser(userID, role.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// Remove all roles first
	err = s.RemoveAllRolesFromUser(id)
	if err != nil {
		return err
	}

	// Delete user
	return s.userRepo.DeleteUser(id)
}

// GetAllRoles retrieves all available roles
func (s *UserService) GetAllRoles() ([]*domain.Role, error) {
	return s.roleRepo.GetAllRoles()
}
