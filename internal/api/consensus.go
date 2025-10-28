package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

// LockConsensusScore handles POST /api/projects/:id/consensus/lock (T034 - US5)
func (h *Handler) LockConsensusScore(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T034 - US5)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	var req domain.LockConsensusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Get facilitator ID from context
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Facilitator ID not found in context",
		})
		return
	}

	consensus, err := h.consensusService.LockConsensusScore(projectID, req.FeatureID, facilitatorID.(int), req.SValue, req.SComplexity, req.Rationale)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to lock consensus score",
			})
		}
		return
	}

	// Audit log the consensus lock action (T043 - US9)
	if h.auditService != nil {
		err = h.auditService.LogConsensusAction(projectID, facilitatorID.(int), req.FeatureID, "consensus_locked", req.SValue, req.SComplexity, req.Rationale)
		if err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	// Broadcast consensus lock via WebSocket (T035 - US5)
	h.broadcastConsensusLocked(projectID, consensus)

	c.JSON(http.StatusOK, gin.H{
		"consensus": consensus,
		"message":   "Consensus score locked successfully",
	})
}

// UnlockConsensusScore handles POST /api/projects/:id/consensus/unlock (T034 - US5)
func (h *Handler) UnlockConsensusScore(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T034 - US5)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	var req domain.UnlockConsensusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Get facilitator ID from context
	facilitatorID, exists := c.Get("attendee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Facilitator ID not found in context",
		})
		return
	}

	err = h.consensusService.UnlockConsensusScore(projectID, req.FeatureID, facilitatorID.(int))
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unlock consensus score",
			})
		}
		return
	}

	// Audit log the consensus unlock action (T043 - US9)
	if h.auditService != nil {
		err = h.auditService.LogConsensusAction(projectID, facilitatorID.(int), req.FeatureID, "consensus_unlocked", 0, 0, "")
		if err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	// Broadcast consensus unlock via WebSocket (T035 - US5)
	h.broadcastConsensusUnlocked(projectID, req.FeatureID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Consensus score unlocked successfully",
	})
}

// GetProjectConsensus handles GET /api/projects/:id/consensus (T034 - US5)
func (h *Handler) GetProjectConsensus(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	consensus, err := h.consensusService.GetProjectConsensus(projectID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve consensus scores",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"consensus": consensus,
		"count":     len(consensus),
	})
}

// broadcastConsensusLocked sends WebSocket notification for consensus locking (T035 - US5)
func (h *Handler) broadcastConsensusLocked(projectID int, consensus *domain.ConsensusScore) {
	if h.wsHub == nil || consensus == nil {
		return
	}

	// Get feature name
	featureName := ""
	if consensus.Feature != nil {
		featureName = consensus.Feature.Title
	}

	// Get facilitator name
	facilitatorName := ""
	if consensus.Facilitator != nil {
		facilitatorName = consensus.Facilitator.Name
	}

	lockMsg := websocket.ConsensusLockedMessage{
		ProjectID:       projectID,
		FeatureID:       consensus.FeatureID,
		FeatureName:     featureName,
		SValue:          consensus.SValue,
		SComplexity:     consensus.SComplexity,
		FacilitatorID:   consensus.LockedBy,
		FacilitatorName: facilitatorName,
		LockedAt:        consensus.LockedAt.Format("2006-01-02T15:04:05Z07:00"),
		Rationale:       consensus.Rationale,
	}

	msg, err := websocket.CreateMessage(websocket.MessageTypeConsensusLocked, lockMsg)
	if err != nil {
		return
	}

	h.wsHub.BroadcastToProject(projectID, msg)
}

// broadcastConsensusUnlocked sends WebSocket notification for consensus unlocking (T035 - US5)
func (h *Handler) broadcastConsensusUnlocked(projectID, featureID int) {
	if h.wsHub == nil {
		return
	}

	unlockMsg := websocket.ConsensusUnlockedMessage{
		ProjectID: projectID,
		FeatureID: featureID,
	}

	msg, err := websocket.CreateMessage(websocket.MessageTypeConsensusUnlocked, unlockMsg)
	if err != nil {
		return
	}

	h.wsHub.BroadcastToProject(projectID, msg)
}
