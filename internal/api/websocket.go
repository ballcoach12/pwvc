package api

import (
	"net/http"
	"strconv"

	"pwvc/internal/websocket"

	"github.com/gin-gonic/gin"
	gorilla_websocket "github.com/gorilla/websocket"
)

// HandleWebSocket handles WebSocket connections for pairwise sessions
func (h *Handler) HandleWebSocket(c *gin.Context) {
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

	// Get attendee ID from query parameter
	attendeeIDStr := c.Query("attendee_id")
	if attendeeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing attendee_id parameter",
		})
		return
	}

	attendeeID, err := strconv.Atoi(attendeeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid attendee ID",
		})
		return
	}

	// Validate session exists and belongs to the project
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

	// Validate attendee exists and belongs to the project
	attendee, err := h.attendeeService.GetAttendee(attendeeID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if attendee.ProjectID != projectID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Attendee does not belong to this project",
		})
		return
	}

	// Upgrade the HTTP connection to WebSocket
	upgrader := gorilla_websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// In production, implement proper origin validation
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upgrade to WebSocket",
		})
		return
	}

	// Create and register the WebSocket client
	client := websocket.NewClient(h.wsHub, conn, sessionID, attendeeID)

	// Set user agent if available
	if userAgent := c.GetHeader("User-Agent"); userAgent != "" {
		// We can't modify the client after creation, but we could log it
		// In a more sophisticated implementation, we'd pass this during creation
	}

	// Register the client with the hub
	h.wsHub.RegisterClient(client)

	// Start the client's read and write pumps
	client.Start()
}

// GetWebSocketStats returns WebSocket hub statistics (for monitoring/debugging)
func (h *Handler) GetWebSocketStats(c *gin.Context) {
	stats := h.wsHub.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
