package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// AttendeeLoginRequest represents the login payload
type AttendeeLoginRequest struct {
	AttendeeID int    `json:"attendee_id" binding:"required"`
	PIN        string `json:"pin" binding:"required"`
}

// AttendeeLoginResponse represents the login response
type AttendeeLoginResponse struct {
	Attendee *domain.Attendee `json:"attendee"`
	Token    string           `json:"token"` // Simple token for now
}

// LoginAttendee handles POST /api/projects/:id/attendees/login
func (h *Handler) LoginAttendee(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req AttendeeLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	// Get attendee
	attendee, err := h.attendeeService.GetAttendee(req.AttendeeID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Verify attendee belongs to project
	if attendee.ProjectID != projectID {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Attendee not found in this project",
		})
		return
	}

	// Debug logging
	hashedPIN := hashPIN(req.PIN)
	fmt.Printf("DEBUG: Attendee ID %d, PinHash from DB: '%s', Computed hash: '%s'\n", attendee.ID, attendee.PinHash, hashedPIN)

	// Verify PIN (simple hash comparison for now)
	if attendee.PinHash != hashedPIN {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid PIN",
		})
		return
	}

	// Create simple token (project:attendee format for now)
	token := fmt.Sprintf("%d:%d", projectID, attendee.ID)

	c.JSON(http.StatusOK, AttendeeLoginResponse{
		Attendee: attendee,
		Token:    token,
	})
}

// hashPIN creates a simple hash of the PIN
func hashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}
