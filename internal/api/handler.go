package api

import (
	"net/http"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/service"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

// Handler holds all the services needed by the API handlers
type Handler struct {
	attendeeService *service.AttendeeService
	featureService  *service.FeatureService
	projectService  *service.ProjectService
	pairwiseService *service.PairwiseService
	pwvcService     *service.PWVCService
	resultsService  *service.ResultsService
	progressService *service.ProgressService
	wsHub           *websocket.Hub
	priorityRepo    *repository.PriorityRepository
}

// NewHandler creates a new API handler with the required services
func NewHandler(
	attendeeService *service.AttendeeService,
	featureService *service.FeatureService,
	projectService *service.ProjectService,
	pairwiseService *service.PairwiseService,
	pwvcService *service.PWVCService,
	resultsService *service.ResultsService,
	progressService *service.ProgressService,
	priorityRepo *repository.PriorityRepository,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		attendeeService: attendeeService,
		featureService:  featureService,
		projectService:  projectService,
		pairwiseService: pairwiseService,
		pwvcService:     pwvcService,
		resultsService:  resultsService,
		progressService: progressService,
		priorityRepo:    priorityRepo,
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

			// Results endpoints
			projects.POST("/:id/calculate-results", h.CalculateResults)
			projects.GET("/:id/results", h.GetResults)
			projects.GET("/:id/results/export", h.ExportResults)
			projects.GET("/:id/results/summary", h.GetResultsSummary)
			projects.GET("/:id/results/status", h.CheckResultsStatus)
			projects.GET("/:id/results/preview", h.PreviewExport)

			// Progress endpoints
			projects.GET("/:id/progress", h.GetProjectProgress)
			projects.POST("/:id/progress/advance", h.AdvancePhase)
			projects.POST("/:id/progress/complete", h.CompletePhase)
			projects.GET("/:id/progress/phases", h.GetAvailablePhases)
		}

		// WebSocket endpoint
		api.GET("/ws/:projectId", h.HandleWebSocket)
		api.GET("/ws/stats", h.GetWebSocketStats)
	}

	// Add health endpoints to the main router (not under /api)
	h.registerHealthEndpoints(router)
}

// registerHealthEndpoints adds health check endpoints
func (h *Handler) registerHealthEndpoints(router *gin.Engine) {
	// Basic health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "pwvc",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// If we have access to database, use detailed health checks
	// This would need to be passed to the handler constructor
	// For now, we'll use a simple implementation
}

// Progress endpoints - delegate to progress handler
func (h *Handler) GetProjectProgress(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService)
	progressHandler.GetProjectProgress(c)
}

func (h *Handler) AdvancePhase(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService)
	progressHandler.AdvancePhase(c)
}

func (h *Handler) CompletePhase(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService)
	progressHandler.CompletePhase(c)
}

func (h *Handler) GetAvailablePhases(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService)
	progressHandler.GetAvailablePhases(c)
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
