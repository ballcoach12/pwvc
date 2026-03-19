package service

import (
	"errors"
	"pairwise/internal/domain"
	"pairwise/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrUserNotFound       = errors.New("user not found")
	ErrWeakPassword       = errors.New("password does not meet strength requirements")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(credentials *domain.LoginRequest) (string, *domain.User, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(credentials.Username)
	if err != nil {
		return "", nil, ErrInvalidCredentials // Don't reveal if user exists
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, ErrUserNotActive
	}

	// Verify password
	if !s.VerifyPassword(user.PasswordHash, credentials.Password) {
		return "", nil, ErrInvalidCredentials
	}

	// Load user roles
	err = s.userRepo.LoadUserRoles(user)
	if err != nil {
		return "", nil, err
	}

	// Generate JWT token
	roles := user.GetRoleNames()
	token, err := s.jwtService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// HashPassword hashes a plain text password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	// Validate password strength
	if err := s.ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	// Hash password with bcrypt cost 12
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword compares a hashed password with a plain text password
func (s *AuthService) VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// ValidatePasswordStrength checks if password meets requirements
func (s *AuthService) ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return ErrWeakPassword
	}
	// Add more validation rules as needed
	// - Must contain uppercase letter
	// - Must contain number
	// - Must contain special character
	// etc.
	return nil
}

// GetCurrentUser retrieves the current authenticated user
func (s *AuthService) GetCurrentUser(userID uint) (*domain.UserProfile, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Load roles
	err = s.userRepo.LoadUserRoles(user)
	if err != nil {
		return nil, err
	}

	profile := user.ToUserProfile()
	return &profile, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, request *domain.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !s.VerifyPassword(user.PasswordHash, request.CurrentPassword) {
		return ErrInvalidCredentials
	}

	// Hash new password
	newHashedPassword, err := s.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = newHashedPassword
	return s.userRepo.UpdateUser(user)
}

// ResetPassword resets a user's password (admin only)
func (s *AuthService) ResetPassword(userID uint, request *domain.ResetPasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Hash new password
	newHashedPassword, err := s.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = newHashedPassword
	return s.userRepo.UpdateUser(user)
}
