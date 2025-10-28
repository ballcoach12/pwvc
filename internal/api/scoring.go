package api

import (
	"net/http"
	"strconv"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
)

// SubmitValueScore handles POST /api/projects/:id/scores/value (T030 - US4)
func (h *Handler) SubmitValueScore(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	score, err := h.scoringService.SubmitValueScore(projectID, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to submit value score",
			})
		}
		return
	}

	// Audit log the score submission (T043 - US9)
	if h.auditService != nil {
		err = h.auditService.LogScoreAction(projectID, req.AttendeeID, req.FeatureID, "value", req.FibonacciValue, req.Rationale)
		if err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	// Broadcast score submission via WebSocket (T033 - US4)
	h.broadcastScoreSubmitted(projectID, score, "value")

	c.JSON(http.StatusOK, gin.H{
		"score":   score,
		"message": "Value score submitted successfully",
	})
}

// SubmitComplexityScore handles POST /api/projects/:id/scores/complexity (T030 - US4)
func (h *Handler) SubmitComplexityScore(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	score, err := h.scoringService.SubmitComplexityScore(projectID, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to submit complexity score",
			})
		}
		return
	}

	// Audit log the score submission (T043 - US9)
	if h.auditService != nil {
		err = h.auditService.LogScoreAction(projectID, req.AttendeeID, req.FeatureID, "complexity", req.FibonacciValue, req.Rationale)
		if err != nil {
			// Log error but don't fail the request
			// TODO: Add proper logging
		}
	}

	// Broadcast score submission via WebSocket (T033 - US4)
	h.broadcastScoreSubmitted(projectID, score, "complexity")

	c.JSON(http.StatusOK, gin.H{
		"score":   score,
		"message": "Complexity score submitted successfully",
	})
}

// GetProjectScores handles GET /api/projects/:id/scores (T030 - US4)
func (h *Handler) GetProjectScores(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	criterionType := c.Query("criterion")
	if criterionType != "" && criterionType != "value" && criterionType != "complexity" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid criterion type. Must be 'value' or 'complexity'",
		})
		return
	}

	scores, err := h.scoringService.GetProjectScores(projectID, criterionType)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr.Message,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve scores",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scores": scores,
		"count":  len(scores),
	})
}

// broadcastScoreSubmitted sends WebSocket notification for score submissions (T033 - US4)
func (h *Handler) broadcastScoreSubmitted(projectID int, score *domain.FibonacciScore, criterionType string) {
	if h.wsHub == nil || score == nil {
		return
	}

	// Get feature and attendee names
	featureName := ""
	attendeeName := ""
	if score.Feature != nil {
		featureName = score.Feature.Title
	}
	if score.Attendee != nil {
		attendeeName = score.Attendee.Name
	}

	scoreMsg := websocket.ScoreSubmittedMessage{
		ProjectID:     projectID,
		FeatureID:     score.FeatureID,
		FeatureName:   featureName,
		AttendeeID:    score.AttendeeID,
		AttendeeName:  attendeeName,
		CriterionType: criterionType,
		Score:         score.FibonacciValue,
		SubmittedAt:   time.Now().Format(time.RFC3339),
	}

	msg, err := websocket.CreateMessage(websocket.MessageTypeScoreSubmitted, scoreMsg)
	if err != nil {
		return
	}

	h.wsHub.BroadcastToProject(projectID, msg)
}
