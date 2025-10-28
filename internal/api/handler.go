package api

import (
	"net/http"
	"strconv"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/service"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

// Handler holds all the services needed by the API handlers
type Handler struct {
	attendeeService  *service.AttendeeService
	featureService   *service.FeatureService
	projectService   *service.ProjectService
	pairwiseService  *service.PairwiseService
	pwvcService      *service.PWVCService
	resultsService   *service.ResultsService
	progressService  *service.ProgressService
	scoringService   *service.ScoringService
	consensusService *service.ConsensusService
	auditService     *service.AuditService
	wsHub            *websocket.Hub
	priorityRepo     repository.PriorityRepository
	attendeeRepo     repository.AttendeeRepository
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
	scoringService *service.ScoringService,
	consensusService *service.ConsensusService,
	auditService *service.AuditService,
	priorityRepo repository.PriorityRepository,
	attendeeRepo repository.AttendeeRepository,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		attendeeService:  attendeeService,
		featureService:   featureService,
		projectService:   projectService,
		pairwiseService:  pairwiseService,
		pwvcService:      pwvcService,
		resultsService:   resultsService,
		progressService:  progressService,
		scoringService:   scoringService,
		consensusService: consensusService,
		auditService:     auditService,
		priorityRepo:     priorityRepo,
		attendeeRepo:     attendeeRepo,
		wsHub:            hub,
	}
}

// RegisterRoutes sets up all the API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Authentication endpoints (public - no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.LoginGlobal) // New global login endpoint
		}

		// PIN management endpoints (public - no auth required)
		api.POST("/setup-pin", h.SetupPin)

		// Public project join endpoints (no auth required)
		api.POST("/join/:invite_code", h.JoinProjectByInvite)

		// Public project-specific login endpoint (for backwards compatibility)
		api.POST("/projects/:id/attendees/login", h.LoginAttendee)

		// Project endpoints (require authentication)
		projects := api.Group("/projects")
		projects.Use(h.RequireAuth()) // Require authentication for all project endpoints
		{
			projects.GET("", h.GetProjects)
			projects.POST("", h.CreateProject)
			projects.GET("/:id", h.GetProject)
			projects.PUT("/:id", h.UpdateProject)
			projects.DELETE("/:id", h.DeleteProject)

			// Attendee endpoints
			projects.GET("/:id/attendees", h.GetProjectAttendees)
			projects.POST("/:id/attendees", h.CreateAttendee)
			projects.POST("/:id/attendees/invite", h.CreateAttendeeWithoutPin)
			projects.PUT("/:id/attendees/:attendee_id/pin", h.SetAttendeePin)
			projects.POST("/:id/attendees/:attendee_id/invite", h.GenerateInviteToken)
			projects.DELETE("/:id/attendees/:attendee_id", h.DeleteAttendee)

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
			projects.POST("/:id/pairwise/reassign", h.ReassignPendingComparisons)        // T042 - US8
			projects.GET("/:id/pairwise/pending", h.GetPendingComparisons)               // T042 - US8
			projects.GET("/:id/pairwise/reassignment-options", h.GetReassignmentOptions) // T042 - US8

			// Fibonacci scoring endpoints (T030 - US4)
			projects.POST("/:id/scores/value", h.SubmitValueScore)
			projects.POST("/:id/scores/complexity", h.SubmitComplexityScore)
			projects.GET("/:id/scores", h.GetProjectScores)

			// Consensus endpoints (T034 - US5)
			projects.POST("/:id/consensus/lock", h.LockConsensusScore)
			projects.POST("/:id/consensus/unlock", h.UnlockConsensusScore)
			projects.GET("/:id/consensus", h.GetProjectConsensus)

			// Results endpoints
			projects.POST("/:id/calculate-results", h.CalculateResults)
			projects.GET("/:id/results", h.GetResults)
			projects.GET("/:id/results/export", h.ExportResults)
			projects.GET("/:id/results/summary", h.GetResultsSummary)
			projects.GET("/:id/results/status", h.CheckResultsStatus)
			projects.GET("/:id/results/preview", h.PreviewExport)

			// Audit endpoints (T045 - US9)
			projects.GET("/:id/audit", h.GetAuditReport)
			projects.GET("/:id/audit/export", h.ExportAuditReport)
			projects.GET("/:id/audit/statistics", h.GetAuditStatistics)

			// Progress endpoints
			projects.GET("/:id/progress", h.GetProjectProgress)
			projects.POST("/:id/progress/advance", h.AdvancePhase)
			projects.POST("/:id/progress/complete", h.CompletePhase)
			projects.GET("/:id/progress/phases", h.GetAvailablePhases)
			projects.POST("/:id/progress/pause", h.PausePhase)              // T039 - US7
			projects.POST("/:id/progress/resume", h.ResumePhase)            // T039 - US7
			projects.GET("/:id/progress/fibonacci", h.GetFibonacciProgress) // T041 - US8
			projects.GET("/:id/progress/details", h.GetProgressDetails)     // T041 - US8
		}

		// WebSocket endpoint
		api.GET("/ws/:projectId", h.HandleWebSocket)
		api.GET("/ws/stats", h.GetWebSocketStats)

		// Authenticated attendee PIN management
		attendees := api.Group("/attendees")
		attendees.Use(h.RequireAuth())
		{
			attendees.PUT("/:id/change-pin", h.ChangePin)
		}
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
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.GetProjectProgress(c)
}

func (h *Handler) AdvancePhase(c *gin.Context) {
	// Add audit logging for phase changes (T043 - US9)
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var request struct {
		Phase string `json:"phase" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current phase before change for audit log
	currentProgress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current progress"})
		return
	}
	oldPhase := string(currentProgress.CurrentPhase)

	// Create progress handler and advance phase
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.AdvancePhase(c)

	// If successful, log the phase change (check if response was successful)
	if c.Writer.Status() == http.StatusOK {
		facilitatorID, exists := c.Get("attendee_id")
		if exists && h.auditService != nil {
			err = h.auditService.LogPhaseChangeAction(projectID, facilitatorID.(int), oldPhase, request.Phase)
			if err != nil {
				// Log error but don't fail the request
				// TODO: Add proper logging
			}
		}
	}
}

func (h *Handler) CompletePhase(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.CompletePhase(c)
}

func (h *Handler) GetAvailablePhases(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.GetAvailablePhases(c)
}

func (h *Handler) PausePhase(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.PausePhase(c)
}

func (h *Handler) ResumePhase(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.ResumePhase(c)
}

func (h *Handler) GetFibonacciProgress(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.GetFibonacciProgress(c)
}

func (h *Handler) GetProgressDetails(c *gin.Context) {
	progressHandler := NewProgressHandler(h.progressService, h.wsHub, h.attendeeRepo)
	progressHandler.GetProgressDetails(c)
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
