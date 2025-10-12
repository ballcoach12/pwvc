package api

import (
	"net/http"

	"pwvc/internal/domain"
	"pwvc/internal/service"
	"pwvc/internal/websocket"

	"github.com/gin-gonic/gin"
)

// Handler holds all the services needed by the API handlers
type Handler struct {
	attendeeService *service.AttendeeService
	featureService  *service.FeatureService
	projectService  *service.ProjectService
	pairwiseService *service.PairwiseService
	pwvcService     *service.PWVCService
	wsHub           *websocket.Hub
}

// NewHandler creates a new API handler with the required services
func NewHandler(
	attendeeService *service.AttendeeService,
	featureService *service.FeatureService,
	projectService *service.ProjectService,
	pairwiseService *service.PairwiseService,
	pwvcService *service.PWVCService,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		attendeeService: attendeeService,
		featureService:  featureService,
		projectService:  projectService,
		pairwiseService: pairwiseService,
		pwvcService:     pwvcService,
		wsHub:           hub,
	}
}

// RegisterRoutes sets up all the API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Project endpoints
		projects := api.Group("/projects")
		{
			projects.GET("", h.GetProjects)
			projects.POST("", h.CreateProject)
			projects.GET("/:id", h.GetProject)
			projects.PUT("/:id", h.UpdateProject)
			projects.DELETE("/:id", h.DeleteProject)

			// Attendee endpoints
			projects.GET("/:id/attendees", h.GetProjectAttendees)
			projects.POST("/:id/attendees", h.CreateAttendee)
			projects.DELETE("/:id/attendees/:attendeeId", h.DeleteAttendee)

			// Feature endpoints
			projects.GET("/:id/features", h.GetProjectFeatures)
			projects.POST("/:id/features", h.CreateFeature)
			projects.GET("/:id/features/:featureId", h.GetFeature)
			projects.PUT("/:id/features/:featureId", h.UpdateFeature)
			projects.DELETE("/:id/features/:featureId", h.DeleteFeature)
			projects.POST("/:id/features/import", h.ImportFeatures)
			projects.GET("/:id/features/export", h.ExportFeatures)

			// Pairwise comparison endpoints
			projects.POST("/:id/pairwise", h.StartPairwiseSession)
			projects.GET("/:id/pairwise", h.GetPairwiseSession)
			projects.GET("/:id/pairwise/comparisons", h.GetPairwiseSessionComparisons)
			projects.POST("/:id/pairwise/votes", h.SubmitPairwiseVote)
			projects.POST("/:id/pairwise/complete", h.CompletePairwiseSession)
			projects.GET("/:id/pairwise/next", h.GetNextComparison)
		}

		// WebSocket endpoint
		api.GET("/ws/:projectId", h.HandleWebSocket)
		api.GET("/ws/stats", h.GetWebSocketStats)
	}
}

// handleServiceError is a utility function to handle service errors consistently
func handleServiceError(c *gin.Context, err error) {
	if apiErr, ok := err.(*domain.APIError); ok {
		c.JSON(apiErr.Code, gin.H{
			"error": apiErr.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
}
