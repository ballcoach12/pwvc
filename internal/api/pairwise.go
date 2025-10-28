package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// StartPairwiseSession handles POST /api/projects/:id/pairwise-sessions
func (h *Handler) StartPairwiseSession(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.CreatePairwiseSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	session, err := h.pairwiseService.StartPairwiseSession(projectID, req.CriterionType)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session": session,
	})
}

// GetPairwiseSession handles GET /api/projects/:id/pairwise
func (h *Handler) GetPairwiseSession(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Get query parameter for criterion type (default to complexity)
	criterionType := c.DefaultQuery("type", "complexity")
	if criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid criterion type. Must be 'value' or 'complexity'",
		})
		return
	}

	// Get active session for the project and criterion
	session, progress, err := h.pairwiseService.GetActiveSession(projectID, domain.CriterionType(criterionType))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session":  session,
		"progress": progress,
	})
}

// GetPairwiseSessionComparisons handles GET /api/projects/:id/pairwise/comparisons
func (h *Handler) GetPairwiseSessionComparisons(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Get query parameter for criterion type (default to complexity)
	criterionType := c.DefaultQuery("type", "complexity")
	if criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid criterion type. Must be 'value' or 'complexity'",
		})
		return
	}

	// Get active session for the project and criterion
	session, _, err := h.pairwiseService.GetActiveSession(projectID, domain.CriterionType(criterionType))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	comparisons, err := h.pairwiseService.GetSessionComparisons(session.ID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comparisons": comparisons,
	})
}

// SubmitPairwiseVote handles POST /api/projects/:id/pairwise-sessions/:session_id/vote
func (h *Handler) SubmitPairwiseVote(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Get query parameter for criterion type (default to complexity)
	criterionType := c.DefaultQuery("type", "complexity")
	if criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid criterion type. Must be 'value' or 'complexity'",
		})
		return
	}

	// Get the current active session for the project
	session, _, err := h.pairwiseService.GetActiveSession(projectID, domain.CriterionType(criterionType))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No active pairwise session found",
		})
		return
	}
	sessionID := session.ID

	var req domain.SubmitVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Session is already verified as belonging to the project from GetActiveSession above

	vote, err := h.pairwiseService.SubmitVote(sessionID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"vote": vote,
	})
}

// CompletePairwiseSession handles POST /api/projects/:id/pairwise-sessions/:session_id/complete
func (h *Handler) CompletePairwiseSession(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	// Verify session belongs to the project first
	session, _, err := h.pairwiseService.GetSession(sessionID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if session.ProjectID != projectID {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	err = h.pairwiseService.CompleteSession(sessionID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session completed successfully",
	})
}

// GetNextComparison handles GET /api/projects/:id/pairwise-sessions/:session_id/next-comparison/:attendee_id
func (h *Handler) GetNextComparison(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	attendeeID, err := strconv.Atoi(c.Param("attendee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	// Verify session belongs to the project first
	session, _, err := h.pairwiseService.GetSession(sessionID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if session.ProjectID != projectID {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	comparison, err := h.pairwiseService.GetNextComparison(sessionID, attendeeID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comparison": comparison,
	})
}
