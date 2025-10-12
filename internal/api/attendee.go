package api

import (
	"net/http"
	"strconv"

	"pwvc/internal/domain"

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
