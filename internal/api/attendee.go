package api

import (
	"net/http"
	"strconv"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// CreateAttendee handles POST /api/projects/:id/attendees
func (h *Handler) CreateAttendee(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.CreateAttendeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	attendee, err := h.attendeeService.CreateAttendee(projectID, req)
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

	c.JSON(http.StatusCreated, attendee)
}

// GetProjectAttendees handles GET /api/projects/:id/attendees
func (h *Handler) GetProjectAttendees(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	attendees, err := h.attendeeService.GetProjectAttendees(projectID)
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
		"attendees": attendees,
	})
}

// DeleteAttendee handles DELETE /api/projects/:id/attendees/:attendee_id
func (h *Handler) DeleteAttendee(c *gin.Context) {
	attendeeID, err := strconv.Atoi(c.Param("attendee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	err = h.attendeeService.DeleteAttendee(attendeeID)
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

// SetAttendeePin handles PUT /api/projects/:id/attendees/:attendee_id/pin
func (h *Handler) SetAttendeePin(c *gin.Context) {
	attendeeID, err := strconv.Atoi(c.Param("attendee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	var req domain.SetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	err = h.attendeeService.SetPIN(attendeeID, req)
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
		"message": "PIN set successfully",
	})
}

// CreateAttendeeWithoutPin handles POST /api/projects/:id/attendees/invite
func (h *Handler) CreateAttendeeWithoutPin(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.CreateAttendeeWithoutPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	attendee, err := h.attendeeService.CreateAttendeeWithoutPIN(projectID, req)
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

	// Generate invite token for PIN setup
	tokenResponse, err := h.attendeeService.GenerateInviteToken(attendee.ID)
	if err != nil {
		// Attendee was created but token generation failed
		c.JSON(http.StatusPartialContent, gin.H{
			"attendee": attendee,
			"warning":  "Attendee created but invite token generation failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"attendee":     attendee,
		"invite_token": tokenResponse.InviteToken,
		"expires_at":   tokenResponse.ExpiresAt,
	})
}

// GenerateInviteToken handles POST /api/projects/:id/attendees/:attendee_id/invite
func (h *Handler) GenerateInviteToken(c *gin.Context) {
	attendeeID, err := strconv.Atoi(c.Param("attendee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	tokenResponse, err := h.attendeeService.GenerateInviteToken(attendeeID)
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

	c.JSON(http.StatusOK, tokenResponse)
}

// SetupPin handles POST /api/setup-pin
func (h *Handler) SetupPin(c *gin.Context) {
	var req domain.SetupPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	err := h.attendeeService.SetupPIN(req)
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
		"message": "PIN set up successfully",
	})
}

// ChangePin handles PUT /api/attendees/:id/change-pin
func (h *Handler) ChangePin(c *gin.Context) {
	attendeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	var req domain.ChangePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	err = h.attendeeService.ChangePIN(attendeeID, req)
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
		"message": "PIN changed successfully",
	})
}
