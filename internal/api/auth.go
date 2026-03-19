package api

import (
	"net/http"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	User    *domain.User `json:"user"`
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"code":    "INVALID_REQUEST",
		})
		return
	}

	// Convert to domain LoginRequest
	loginReq := &domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	token, user, err := h.authService.Login(loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    "LOGIN_FAILED",
		})
		return
	}

	// Set secure HTTP-only cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth", token, 86400, "/", "", false, true) // 24 hours, secure=false for dev, httpOnly=true

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *Handler) Logout(c *gin.Context) {
	// Clear the auth cookie
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// GetCurrentUser handles GET /api/v1/auth/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User not authenticated",
			"code":    "NOT_AUTHENTICATED",
		})
		return
	}

	username, _ := c.Get("username")
	roles, _ := c.Get("roles")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
		"roles":    roles,
	})
}
