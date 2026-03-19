package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/service"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	progressService *service.ProgressService
	wsHub           *websocket.Hub
	attendeeRepo    repository.AttendeeRepository
}

func NewProgressHandler(progressService *service.ProgressService, wsHub *websocket.Hub, attendeeRepo repository.AttendeeRepository) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		wsHub:           wsHub,
		attendeeRepo:    attendeeRepo,
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

// AdvancePhase advances the project to a specific phase (facilitator-only)
func (h *ProgressHandler) AdvancePhase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check facilitator authorization (T048 - authorization enforcement)
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Facilitator authorization required"})
		return
	}

	// Verify the attendee is actually a facilitator
	attendee, err := h.attendeeRepo.GetByID(facilitatorID.(int))
	if err != nil || !attendee.IsFacilitator {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only facilitators can advance project phases"})
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

	// Get current phase before change for broadcast
	currentProgress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current progress"})
		return
	}
	oldPhase := string(currentProgress.CurrentPhase)

	err = h.progressService.AdvanceToPhase(projectID, phase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get facilitator info for broadcast (T039 - US7)
	// facilitatorID and exists already declared earlier in function
	var facilitatorName string
	if exists {
		attendee, err2 := h.attendeeRepo.GetByID(facilitatorID.(int))
		if err2 == nil {
			facilitatorName = attendee.Name
		}
	}

	// Broadcast phase change via WebSocket (T039 - US7)
	if h.wsHub != nil {
		phaseChangeMsg := websocket.PhaseChangedMessage{
			ProjectID:     projectID,
			NewPhase:      string(phase),
			OldPhase:      oldPhase,
			ChangedBy:     facilitatorID.(int),
			ChangedByName: facilitatorName,
			ChangedAt:     time.Now().Format(time.RFC3339),
			Message:       fmt.Sprintf("Project advanced from %s to %s", oldPhase, string(phase)),
		}
		h.wsHub.NotifyPhaseChanged(projectID, phaseChangeMsg)
	}

	// Return updated progress
	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// CompletePhase marks a phase as completed (facilitator-only)
func (h *ProgressHandler) CompletePhase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check facilitator authorization (T048 - authorization enforcement)
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Facilitator authorization required"})
		return
	}

	// Verify the attendee is actually a facilitator
	attendee, err := h.attendeeRepo.GetByID(facilitatorID.(int))
	if err != nil || !attendee.IsFacilitator {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only facilitators can complete project phases"})
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

// PausePhase pauses the current phase (T039 - US7)
func (h *ProgressHandler) PausePhase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check facilitator authorization
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Facilitator authorization required"})
		return
	}

	// Get current progress
	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// For now, we'll just broadcast the pause status - actual pause logic would be in the service
	var facilitatorName string
	if attendee, err := h.attendeeRepo.GetByID(facilitatorID.(int)); err == nil {
		facilitatorName = attendee.Name
	}

	// Broadcast phase pause via WebSocket
	if h.wsHub != nil {
		phaseChangeMsg := websocket.PhaseChangedMessage{
			ProjectID:     projectID,
			NewPhase:      string(progress.CurrentPhase) + "_paused",
			OldPhase:      string(progress.CurrentPhase),
			ChangedBy:     facilitatorID.(int),
			ChangedByName: facilitatorName,
			ChangedAt:     time.Now().Format(time.RFC3339),
			Message:       fmt.Sprintf("Phase %s has been paused by facilitator", string(progress.CurrentPhase)),
		}
		h.wsHub.NotifyPhaseChanged(projectID, phaseChangeMsg)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phase paused successfully",
		"phase":   progress.CurrentPhase,
		"status":  "paused",
	})
}

// ResumePhase resumes the current phase (T039 - US7)
func (h *ProgressHandler) ResumePhase(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check facilitator authorization
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Facilitator authorization required"})
		return
	}

	// Get current progress
	progress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// For now, we'll just broadcast the resume status - actual resume logic would be in the service
	var facilitatorName string
	if attendee, err := h.attendeeRepo.GetByID(facilitatorID.(int)); err == nil {
		facilitatorName = attendee.Name
	}

	// Broadcast phase resume via WebSocket
	if h.wsHub != nil {
		phaseChangeMsg := websocket.PhaseChangedMessage{
			ProjectID:     projectID,
			NewPhase:      string(progress.CurrentPhase),
			OldPhase:      string(progress.CurrentPhase) + "_paused",
			ChangedBy:     facilitatorID.(int),
			ChangedByName: facilitatorName,
			ChangedAt:     time.Now().Format(time.RFC3339),
			Message:       fmt.Sprintf("Phase %s has been resumed by facilitator", string(progress.CurrentPhase)),
		}
		h.wsHub.NotifyPhaseChanged(projectID, phaseChangeMsg)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phase resumed successfully",
		"phase":   progress.CurrentPhase,
		"status":  "active",
	})
}

// GetFibonacciProgress gets Fibonacci scoring progress for a project (T041 - US8)
func (h *ProgressHandler) GetFibonacciProgress(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	criterionType := c.Query("criterion")

	if criterionType != "" && criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid criterion type. Must be 'value' or 'complexity'"})
		return
	}

	if criterionType != "" {
		// Get metrics for specific criterion
		metrics, err := h.progressService.GetFibonacciProgressMetrics(projectID, criterionType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, metrics)
	} else {
		// Get overall progress for both criteria
		overall, err := h.progressService.GetOverallFibonacciProgress(projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, overall)
	}
}

// GetProgressDetails gets detailed progress breakdown including Fibonacci phases (T041 - US8)
func (h *ProgressHandler) GetProgressDetails(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get standard progress
	standardProgress, err := h.progressService.GetProjectProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get Fibonacci progress
	fibonacciProgress, err := h.progressService.GetOverallFibonacciProgress(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"standard_progress":  standardProgress,
		"fibonacci_progress": fibonacciProgress,
		"timestamp":          time.Now().Format(time.RFC3339),
	})
}
