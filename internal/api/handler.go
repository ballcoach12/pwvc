package api

import (
	"net/http"
	"strconv"

	"pwvc/internal/domain"
	"pwvc/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler contains all API handlers
type Handler struct {
	projectService  *service.ProjectService
	attendeeService *service.AttendeeService
	featureService  *service.FeatureService
}

// NewHandler creates a new API handler
func NewHandler(projectService *service.ProjectService, attendeeService *service.AttendeeService, featureService *service.FeatureService) *Handler {
	return &Handler{
		projectService:  projectService,
		attendeeService: attendeeService,
		featureService:  featureService,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Project routes
		api.POST("/projects", h.CreateProject)
		api.GET("/projects/:id", h.GetProject)
		api.PUT("/projects/:id", h.UpdateProject)
		api.DELETE("/projects/:id", h.DeleteProject)
		api.GET("/projects", h.ListProjects)

		// Attendee routes
		api.POST("/projects/:id/attendees", h.CreateAttendee)
		api.GET("/projects/:id/attendees", h.GetProjectAttendees)
		api.DELETE("/projects/:id/attendees/:attendee_id", h.DeleteAttendee)

		// Feature routes
		api.POST("/projects/:id/features", h.CreateFeature)
		api.GET("/projects/:id/features", h.GetProjectFeatures)
		api.GET("/projects/:id/features/:feature_id", h.GetFeature)
		api.PUT("/projects/:id/features/:feature_id", h.UpdateFeature)
		api.DELETE("/projects/:id/features/:feature_id", h.DeleteFeature)
		api.POST("/projects/:id/features/import", h.ImportFeatures)
		api.GET("/projects/:id/features/export", h.ExportFeatures)
	}
}

// CreateProject handles POST /api/projects
func (h *Handler) CreateProject(c *gin.Context) {
	var req domain.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	project, err := h.projectService.CreateProject(req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProject handles GET /api/projects/:id
func (h *Handler) GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	project, err := h.projectService.GetProject(id)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject handles PUT /api/projects/:id
func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	project, err := h.projectService.UpdateProject(id, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject handles DELETE /api/projects/:id
func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	err = h.projectService.DeleteProject(id)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListProjects handles GET /api/projects
func (h *Handler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}
