package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

// CreateProject handles POST /api/projects
func (h *Handler) CreateProject(c *gin.Context) {
	var req domain.CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	project, err := h.projectService.CreateProject(req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create project",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProjects handles GET /api/projects
func (h *Handler) GetProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve projects",
			})
		}
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject handles GET /api/projects/:id
func (h *Handler) GetProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve project",
			})
		}
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject handles PUT /api/projects/:id
func (h *Handler) UpdateProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.UpdateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	project, err := h.projectService.UpdateProject(projectID, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update project",
			})
		}
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject handles DELETE /api/projects/:id
func (h *Handler) DeleteProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	if err := h.projectService.DeleteProject(projectID); err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete project",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
	})
}

// GetProjectInviteLink handles GET /api/projects/:id/invite-link (T016 - US1)
func (h *Handler) GetProjectInviteLink(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	project, err := h.projectService.GetProject(projectID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve project",
			})
		}
		return
	}

	// Generate invite code if it doesn't exist
	inviteCode := project.InviteCode
	if inviteCode == "" {
		inviteCode, err = generateInviteCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate invite code",
			})
			return
		}

		// Update project with invite code
		project, err = h.projectService.UpdateProjectInviteCode(projectID, inviteCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save invite code",
			})
			return
		}
	}

	// Create full invite link
	baseURL := c.Request.Header.Get("X-Forwarded-Host")
	if baseURL == "" {
		baseURL = c.Request.Host
	}

	inviteLink := fmt.Sprintf("https://%s/join/%s", baseURL, inviteCode)

	c.JSON(http.StatusOK, gin.H{
		"project_id":   projectID,
		"project_name": project.Name,
		"invite_code":  inviteCode,
		"invite_link":  inviteLink,
	})
}

// JoinProjectByInvite handles POST /api/join/:invite_code (T016 - US1)
func (h *Handler) JoinProjectByInvite(c *gin.Context) {
	inviteCode := c.Param("invite_code")
	if inviteCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid invite code",
		})
		return
	}

	project, err := h.projectService.GetProjectByInviteCode(inviteCode)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invalid invite code",
			})
		}
		return
	}

	// Return project info for the attendee to join
	c.JSON(http.StatusOK, gin.H{
		"project":     project,
		"can_join":    true,
		"next_action": "create_attendee",
		"message":     fmt.Sprintf("Welcome to project: %s", project.Name),
	})
}

// generateInviteCode creates a cryptographically secure random invite code
func generateInviteCode() (string, error) {
	bytes := make([]byte, 6) // 12 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// broadcastSessionStatusChange sends WebSocket notification for session status changes (T018 - US1)
func (h *Handler) broadcastSessionStatusChange(projectID int, newStatus, message string) {
	if h.wsHub == nil {
		return // WebSocket not initialized
	}

	statusMsg := websocket.PhaseChangedMessage{
		ProjectID: projectID,
		NewPhase:  newStatus,
		ChangedAt: time.Now().Format(time.RFC3339),
		Message:   message,
	}

	msg, err := websocket.CreateMessage(websocket.MessageTypePhaseChanged, statusMsg)
	if err != nil {
		// Log error but don't fail the main operation
		return
	}

	// Broadcast to all clients in this project
	h.wsHub.BroadcastToProject(projectID, msg)
}
