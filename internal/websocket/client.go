package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// HubInterface defines the interface for the WebSocket hub
type HubInterface interface {
	BroadcastToSession(sessionID int, message *Message)
	HandleClientJoin(client *Client)
	UnregisterClient(client *Client)
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512

	// Buffer size for client channels
	clientBufferSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for now
		// In production, implement proper origin validation
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan *Message

	// The hub that manages this client
	hub HubInterface

	// Client identification
	sessionID  int
	attendeeID int

	// Client metadata
	userAgent   string
	remoteAddr  string
	connectedAt time.Time

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// Mutex for connection safety
	mu sync.RWMutex

	// Connection state
	isConnected bool
}

// NewClient creates a new WebSocket client
func NewClient(hub HubInterface, conn *websocket.Conn, sessionID, attendeeID int) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		conn:        conn,
		send:        make(chan *Message, clientBufferSize),
		hub:         hub,
		sessionID:   sessionID,
		attendeeID:  attendeeID,
		userAgent:   "", // Will be set during upgrade
		remoteAddr:  conn.RemoteAddr().String(),
		connectedAt: time.Now(),
		ctx:         ctx,
		cancel:      cancel,
		isConnected: true,
	}
}

// GetSessionID returns the client's session ID
func (c *Client) GetSessionID() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionID
}

// GetAttendeeID returns the client's attendee ID
func (c *Client) GetAttendeeID() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.attendeeID
}

// IsConnected returns whether the client is currently connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// Send sends a message to the client
func (c *Client) Send(message *Message) {
	if !c.IsConnected() {
		return
	}

	select {
	case c.send <- message:
	default:
		// Channel is full, close the client
		c.close()
	}
}

// SendJSON sends a JSON message to the client
func (c *Client) SendJSON(msgType MessageType, data interface{}) error {
	message, err := CreateMessage(msgType, data)
	if err != nil {
		return err
	}

	c.Send(message)
	return nil
}

// Start starts the client's read and write goroutines
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

// close closes the client connection and cleans up resources
func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return
	}

	c.isConnected = false
	c.cancel()

	// Close the connection
	c.conn.Close()

	// Close the send channel
	close(c.send)

	// Notify the hub
	c.hub.UnregisterClient(c)
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer c.close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %d: %v", c.attendeeID, err)
			}
			break
		}

		// Set timestamp if not provided
		if message.Timestamp.IsZero() {
			message.Timestamp = time.Now()
		}

		// Handle the message
		c.handleMessage(&message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error for client %d: %v", c.attendeeID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from the client
func (c *Client) handleMessage(message *Message) {
	log.Printf("Received message from client %d: %s", c.attendeeID, message.Type)

	switch message.Type {
	case MessageTypeJoinSession:
		c.handleJoinSession(message)
	case MessageTypeLeaveSession:
		c.handleLeaveSession(message)
	case MessageTypeVoteSubmitted:
		c.handleVoteSubmitted(message)
	default:
		// Send error for unknown message type
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Unknown message type",
			Details: string(message.Type),
		})
	}
}

// handleJoinSession handles join session requests
func (c *Client) handleJoinSession(message *Message) {
	var joinMsg JoinSessionMessage
	if err := message.ParseMessageData(&joinMsg); err != nil {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Invalid join session message",
			Details: err.Error(),
		})
		return
	}

	// Validate session ID matches
	if joinMsg.SessionID != c.sessionID {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Session ID mismatch",
		})
		return
	}

	// Validate attendee ID matches
	if joinMsg.AttendeeID != c.attendeeID {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Attendee ID mismatch",
		})
		return
	}

	// Notify hub of successful join
	c.hub.HandleClientJoin(c)
}

// handleLeaveSession handles leave session requests
func (c *Client) handleLeaveSession(message *Message) {
	var leaveMsg LeaveSessionMessage
	if err := message.ParseMessageData(&leaveMsg); err != nil {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Invalid leave session message",
			Details: err.Error(),
		})
		return
	}

	// Close the connection
	c.close()
}

// handleVoteSubmitted handles vote submission notifications
func (c *Client) handleVoteSubmitted(message *Message) {
	var voteMsg VoteSubmittedMessage
	if err := message.ParseMessageData(&voteMsg); err != nil {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Invalid vote submitted message",
			Details: err.Error(),
		})
		return
	}

	// Validate attendee ID matches
	if voteMsg.AttendeeID != c.attendeeID {
		c.SendJSON(MessageTypeError, ErrorMessage{
			Code:    400,
			Message: "Attendee ID mismatch",
		})
		return
	}

	// Forward to hub for broadcasting
	c.hub.BroadcastToSession(c.sessionID, message)
}

// GetConnectionInfo returns connection information for debugging
func (c *Client) GetConnectionInfo() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"session_id":   c.sessionID,
		"attendee_id":  c.attendeeID,
		"remote_addr":  c.remoteAddr,
		"user_agent":   c.userAgent,
		"connected_at": c.connectedAt,
		"is_connected": c.isConnected,
		"send_buffer":  len(c.send),
	}
}
