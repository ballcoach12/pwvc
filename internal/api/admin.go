package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"
	"pairwise/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// ListUsers handles GET /api/v1/admin/users
func (h *Handler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get users
	users, err := h.userRepo.ListUsers(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve users",
			"code":    "FETCH_ERROR",
		})
		return
	}

	// Get total count
	total, err := h.userRepo.CountUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to count users",
			"code":    "COUNT_ERROR",
		})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// CreateUser handles POST /api/v1/admin/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"code":    "INVALID_REQUEST",
		})
		return
	}

	// Hash the password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"code":    "WEAK_PASSWORD",
		})
		return
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
		IsActive:     true,
	}

	err = h.userRepo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Username already exists",
			"code":    "USERNAME_EXISTS",
		})
		return
	}

	// Assign roles if provided
	if len(req.Roles) > 0 {
		currentUserID, _ := middleware.GetCurrentUserID(c)
		for _, roleName := range req.Roles {
			role, err := h.roleRepo.GetRoleByName(roleName)
			if err != nil {
				continue // Skip invalid roles
			}
			h.roleRepo.AssignRoleToUser(user.ID, role.ID, &currentUserID)
		}
		
		// Reload user with roles
		h.userRepo.LoadUserRoles(user)
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /api/v1/admin/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
			"code":    "INVALID_ID",
		})
		return
	}

	user, err := h.userRepo.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"code":    "USER_NOT_FOUND",
		})
		return
	}

	// Load user roles
	err = h.userRepo.LoadUserRoles(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to load user roles",
			"code":    "LOAD_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /api/v1/admin/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
			"code":    "INVALID_ID",
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
			"code":    "INVALID_REQUEST",
		})
		return
	}

	// Get existing user
	user, err := h.userRepo.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
			"code":    "USER_NOT_FOUND",
		})
		return
	}

	// Update fields if provided
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user
	err = h.userRepo.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update user",
			"code":    "UPDATE_ERROR",
		})
		return
	}

	// Update roles if provided
	if req.Roles != nil {
		currentUserID, _ := middleware.GetCurrentUserID(c)
		
		// Load current roles
		h.userRepo.LoadUserRoles(user)
		
		// Remove all current roles and add new ones
		for _, role := range user.Roles {
			h.roleRepo.RemoveRoleFromUser(user.ID, role.ID)
		}
		
		// Add new roles
		for _, roleName := range req.Roles {
			role, err := h.roleRepo.GetRoleByName(roleName)
			if err != nil {
				continue // Skip invalid roles
			}
			h.roleRepo.AssignRoleToUser(user.ID, role.ID, &currentUserID)
		}
		
		// Reload user with updated roles
		h.userRepo.LoadUserRoles(user)
	}

	c.JSON(http.StatusOK, user)
}

// ListRoles handles GET /api/v1/admin/roles
func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.roleRepo.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve roles",
			"code":    "FETCH_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}