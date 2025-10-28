package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"
	"pairwise/internal/websocket"

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

	// Audit log the vote submission (T043 - US9)
	if h.auditService != nil {
		err = h.auditService.LogVoteAction(projectID, req.AttendeeID, vote.ComparisonID, vote.PreferredFeatureID, vote.IsTieVote)
		if err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	// Broadcast vote update via WebSocket (T022 - US2)
	h.broadcastVoteUpdate(projectID, vote, criterionType)

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

// broadcastVoteUpdate sends WebSocket notification for vote updates (T022 - US2)
func (h *Handler) broadcastVoteUpdate(projectID int, vote *domain.AttendeeVote, criterionType string) {
	if h.wsHub == nil || vote == nil {
		return
	}

	// Get attendee name if available
	attendeeName := ""
	if vote.Attendee != nil {
		attendeeName = vote.Attendee.Name
	}

	voteMsg := websocket.VoteUpdateMessage{
		ComparisonID:       vote.ComparisonID,
		AttendeeID:         vote.AttendeeID,
		AttendeeName:       attendeeName,
		PreferredFeatureID: vote.PreferredFeatureID,
		IsTieVote:          vote.IsTieVote,
		ConsensusReached:   false, // Will be updated by comparison service
	}

	msg, err := websocket.CreateMessage(websocket.MessageTypeVoteUpdate, voteMsg)
	if err != nil {
		return
	}

	h.wsHub.BroadcastToProject(projectID, msg)
}

// ReassignPendingComparisons handles POST /api/projects/:id/pairwise/reassign (T042 - US8)
func (h *Handler) ReassignPendingComparisons(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T042 - US8)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	var req domain.ReassignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	err = h.pairwiseService.ReassignPendingComparisons(projectID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Comparisons reassigned successfully",
		"reassigned_count":  len(req.ComparisonIDs),
		"reassignment_type": req.ReassignmentType,
	})
}

// GetPendingComparisons handles GET /api/projects/:id/pairwise/pending (T042 - US8)
func (h *Handler) GetPendingComparisons(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	sessionID, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	criterionType := c.Query("criterion")
	if criterionType != "" && criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid criterion type",
		})
		return
	}

	// Note: projectID from URL param is available but not used since session already contains project context
	_ = projectID

	pendingComparisons, err := h.pairwiseService.GetPendingComparisons(sessionID, domain.CriterionType(criterionType))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_comparisons": pendingComparisons,
		"count":               len(pendingComparisons),
	})
}

// GetReassignmentOptions handles GET /api/projects/:id/pairwise/reassignment-options (T042 - US8)
func (h *Handler) GetReassignmentOptions(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	sessionID, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	options, err := h.pairwiseService.GetReassignmentOptions(projectID, sessionID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, options)
}
