package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"
	"pairwise/internal/service"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressService *service.ProgressService
}

func NewProgressHandler(progressService *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

// GetProjectProgress gets the current progress for a project
func (h *ProgressHandler) GetProjectProgress(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// AdvancePhase advances the project to a specific phase
func (h *ProgressHandler) AdvancePhase(c *gin.Context) {
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

	phase := domain.WorkflowPhase(request.Phase)

	// Validate phase
	validPhases := []domain.WorkflowPhase{
		domain.PhaseSetup,
		domain.PhaseAttendees,
		domain.PhaseFeatures,
		domain.PhasePairwiseValue,
		domain.PhasePairwiseComplexity,
		domain.PhaseFibonacciValue,
		domain.PhaseFibonacciComplexity,
		domain.PhaseResults,
	}

	isValid := false
	for _, validPhase := range validPhases {
		if phase == validPhase {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phase"})
		return
	}

	err = h.progressService.AdvanceToPhase(projectID, phase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return updated progress
	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// CompletePhase marks a phase as completed
func (h *ProgressHandler) CompletePhase(c *gin.Context) {
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

	phase := domain.WorkflowPhase(request.Phase)

	err = h.progressService.CompletePhase(projectID, phase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return updated progress
	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetAvailablePhases gets the phases that can be accessed
func (h *ProgressHandler) GetAvailablePhases(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	phases, err := h.progressService.GetAvailablePhases(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_phases": phases})
}
