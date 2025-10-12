package api

import (
	"net/http"
	"strconv"

	"pwvc/internal/domain"

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
