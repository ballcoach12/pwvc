package api

import (
	"net/http"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/service"

	"github.com/gin-gonic/gin"
)

// JWTAuthHandler handles JWT authentication-related HTTP requests
type JWTAuthHandler struct {
	authService *service.AuthService
	jwtService  *service.JWTService
}

// NewJWTAuthHandler creates a new JWTAuthHandler
func NewJWTAuthHandler(authService *service.AuthService, jwtService *service.JWTService) *JWTAuthHandler {
	return &JWTAuthHandler{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Login handles user authentication with JWT
// @Summary Login user
// @Description authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *JWTAuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Username and password are required",
			Details: "Both username and password must be provided",
		})
		return
	}

	// Authenticate user
	token, user, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Authentication failed",
			Details: err.Error(),
		})
		return
	}

	// Set HTTP-only cookie for token
	c.SetCookie(
		"auth",                          // name
		token,                           // value
		int((24 * time.Hour).Seconds()), // maxAge (24 hours)
		"/",                             // path
		"",                              // domain (empty for current domain)
		false,                           // secure (set to true in production with HTTPS)
		true,                            // httpOnly
	)

	// Add token to response for API clients
	responseWithToken := struct {
		User  *domain.User `json:"user"`
		Token string       `json:"token"`
	}{
		User:  user,
		Token: token,
	}

	c.JSON(http.StatusOK, responseWithToken)
}

// Logout handles user logout with JWT
// @Summary Logout user
// @Description logout user and clear authentication cookie
// @Tags auth
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /auth/logout [post]
func (h *JWTAuthHandler) Logout(c *gin.Context) {
	// Clear the auth cookie
	c.SetCookie(
		"auth", // name
		"",     // value (empty)
		-1,     // maxAge (negative to delete)
		"/",    // path
		"",     // domain
		false,  // secure
		true,   // httpOnly
	)

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// GetCurrentUser returns the current authenticated user's profile
// @Summary Get current user
// @Description get current authenticated user profile
// @Tags auth
// @Produce json
// @Success 200 {object} domain.UserProfile
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/me [get]
func (h *JWTAuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "User not authenticated",
			Details: "Valid authentication token required",
		})
		return
	}

	// Convert to uint (GORM primary key type)
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Invalid user ID format",
			Details: "Internal authentication error",
		})
		return
	}

	// Get user profile
	profile, err := h.authService.GetCurrentUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user profile",
			Details: err.Error(),
		})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, profile)
}

// ChangePassword handles password change for authenticated users
// @Summary Change password
// @Description change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param passwordData body ChangePasswordRequest true "Password change data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/change-password [post]
func (h *JWTAuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.CurrentPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Current password and new password are required",
			Details: "Both passwords must be provided",
		})
		return
	}

	// Get user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "User not authenticated",
			Details: "Valid authentication token required",
		})
		return
	}

	uid := userID.(uint)

	// Create change password request
	changeReq := &domain.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	// Change password
	err := h.authService.ChangePassword(uid, changeReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to change password",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// Request/Response types for JWT auth handlers

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
