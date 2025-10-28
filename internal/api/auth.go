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
	var dbPinHash string
	if attendee.PinHash != nil {
		dbPinHash = *attendee.PinHash
	}
	fmt.Printf("DEBUG: Attendee ID %d, PinHash from DB: '%s', Computed hash: '%s'\n", attendee.ID, dbPinHash, hashedPIN)

	// Verify PIN (simple hash comparison for now)
	if attendee.PinHash == nil || *attendee.PinHash != hashedPIN {
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

// RequireAuth middleware ensures users are authenticated before accessing protected endpoints
func (h *Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("DEBUG: RequireAuth middleware called for path: %s\n", c.Request.URL.Path)

		// Check for Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Please log in to access this resource",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			// Also support direct token for backwards compatibility
			token = authHeader
		}

		// Parse token format: "projectID:attendeeID"
		attendeeID, projectID, err := h.parseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token format",
				"message": "Please log in again",
			})
			c.Abort()
			return
		}

		// Verify attendee exists and belongs to project
		attendee, err := h.attendeeService.GetAttendee(attendeeID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authentication",
				"message": "User not found",
			})
			c.Abort()
			return
		}

		// Verify attendee belongs to the claimed project
		if attendee.ProjectID != projectID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid project access",
				"message": "User does not belong to this project",
			})
			c.Abort()
			return
		}

		// For project-specific endpoints, validate that the requested project matches the user's project
		if urlProjectIDStr := c.Param("id"); urlProjectIDStr != "" {
			urlProjectID, err := strconv.Atoi(urlProjectIDStr)
			if err == nil && urlProjectID != projectID {
				fmt.Printf("DEBUG: Access denied - URL project ID: %d, User's project ID: %d\n", urlProjectID, projectID)
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Access denied",
					"message": "You can only access your own project",
				})
				c.Abort()
				return
			}
			fmt.Printf("DEBUG: Access allowed - URL project ID: %d, User's project ID: %d\n", urlProjectID, projectID)
		} else {
			fmt.Printf("DEBUG: No project ID in URL path, allowing request to: %s\n", c.Request.URL.Path)
		}

		// Store authentication info in context for handlers to use
		c.Set("attendee_id", attendeeID)
		c.Set("project_id", projectID)
		c.Set("attendee", attendee)
		c.Set("is_facilitator", attendee.IsFacilitator)

		c.Next()
	}
}

// parseToken parses the token format "projectID:attendeeID" and returns attendeeID, projectID, error
func (h *Handler) parseToken(token string) (int, int, error) {
	parts := make([]string, 0, 2)
	colonIndex := -1
	for i, char := range token {
		if char == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex == -1 {
		return 0, 0, fmt.Errorf("invalid token format")
	}

	parts = append(parts, token[:colonIndex])
	parts = append(parts, token[colonIndex+1:])

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid token format")
	}

	projectID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid project ID in token")
	}

	attendeeID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid attendee ID in token")
	}

	return attendeeID, projectID, nil
}

// GlobalLoginRequest represents the global login payload (without knowing project ID)
type GlobalLoginRequest struct {
	AttendeeID int    `json:"attendee_id" binding:"required"`
	PIN        string `json:"pin" binding:"required"`
}

// LoginGlobal handles POST /api/auth/login (global login without project context)
func (h *Handler) LoginGlobal(c *gin.Context) {
	var req GlobalLoginRequest
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify PIN
	hashedPIN := hashPIN(req.PIN)
	if attendee.PinHash == nil || *attendee.PinHash != hashedPIN {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Create token with project:attendee format
	token := fmt.Sprintf("%d:%d", attendee.ProjectID, attendee.ID)

	c.JSON(http.StatusOK, AttendeeLoginResponse{
		Attendee: attendee,
		Token:    token,
	})
}

// IsFacilitator middleware ensures only facilitators can access certain endpoints
func (h *Handler) IsFacilitator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For simplicity, check for attendee_id in headers and verify role
		attendeeIDStr := c.GetHeader("X-Attendee-ID")
		if attendeeIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing attendee ID header",
			})
			c.Abort()
			return
		}

		attendeeID, err := strconv.Atoi(attendeeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid attendee ID",
			})
			c.Abort()
			return
		}

		// Get attendee and check role
		attendee, err := h.attendeeService.GetAttendee(attendeeID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid attendee",
			})
			c.Abort()
			return
		}

		if !attendee.IsFacilitator {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Facilitator access required",
			})
			c.Abort()
			return
		}

		// Store attendee in context for handlers to use
		c.Set("attendee", attendee)
		c.Next()
	}
}

// checkIsFacilitator checks if the current user is a facilitator and returns a boolean
func (h *Handler) checkIsFacilitator(c *gin.Context) bool {
	// Check for attendee_id in headers and verify role
	attendeeIDStr := c.GetHeader("X-Attendee-ID")
	if attendeeIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing attendee ID header",
		})
		return false
	}

	attendeeID, err := strconv.Atoi(attendeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return false
	}

	// Get attendee and check role
	attendee, err := h.attendeeService.GetAttendee(attendeeID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid attendee",
		})
		return false
	}

	if !attendee.IsFacilitator {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Facilitator access required",
		})
		return false
	}

	return true
}
