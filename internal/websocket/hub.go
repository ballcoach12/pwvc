package websocket

import (
	"log"
	"sync"
	"time"

	"pwvc/internal/domain"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients organized by session ID
	sessions map[int]map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to specific sessions
	broadcast chan *BroadcastMessage

	// Attendee repository for getting attendee info
	attendeeRepo AttendeeRepository

	// Mutex for concurrent safety
	mu sync.RWMutex

	// Hub statistics
	stats HubStats
}

// BroadcastMessage represents a message to broadcast to a session
type BroadcastMessage struct {
	SessionID     int
	Message       *Message
	ExcludeClient *Client // Optional: exclude this client from broadcast
}

// HubStats tracks hub statistics
type HubStats struct {
	TotalConnections    int
	ActiveConnections   int
	ActiveSessions      int
	MessagesSent        int64
	MessagesReceived    int64
	ConnectionsAccepted int64
	ConnectionsRejected int64
	LastActivity        time.Time
}

// AttendeeRepository interface for getting attendee information
type AttendeeRepository interface {
	GetByID(attendeeID int) (*domain.Attendee, error)
}

// NewHub creates a new WebSocket hub
func NewHub(attendeeRepo AttendeeRepository) *Hub {
	return &Hub{
		sessions:     make(map[int]map[*Client]bool),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan *BroadcastMessage),
		attendeeRepo: attendeeRepo,
		stats: HubStats{
			LastActivity: time.Now(),
		},
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	log.Println("WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case broadcast := <-h.broadcast:
			h.broadcastMessage(broadcast)
		}
	}
}

// RegisterClient registers a new client connection
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client connection
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastToSession broadcasts a message to all clients in a session
func (h *Hub) BroadcastToSession(sessionID int, message *Message) {
	h.broadcast <- &BroadcastMessage{
		SessionID: sessionID,
		Message:   message,
	}
}

// BroadcastToSessionExcept broadcasts a message to all clients in a session except one
func (h *Hub) BroadcastToSessionExcept(sessionID int, message *Message, excludeClient *Client) {
	h.broadcast <- &BroadcastMessage{
		SessionID:     sessionID,
		Message:       message,
		ExcludeClient: excludeClient,
	}
}

// HandleClientJoin handles when a client successfully joins a session
func (h *Hub) HandleClientJoin(client *Client) {
	attendee, err := h.attendeeRepo.GetByID(client.GetAttendeeID())
	if err != nil {
		log.Printf("Failed to get attendee %d: %v", client.GetAttendeeID(), err)
		client.SendJSON(MessageTypeError, ErrorMessage{
			Code:    500,
			Message: "Failed to get attendee information",
		})
		return
	}

	sessionID := client.GetSessionID()

	// Send welcome message to the joining client
	welcomeMsg := WelcomeMessage{
		SessionID:      sessionID,
		AttendeeID:     client.GetAttendeeID(),
		AttendeeName:   attendee.Name,
		ConnectedCount: h.getSessionClientCount(sessionID),
		SessionStatus:  "active", // TODO: Get actual session status
	}

	client.SendJSON(MessageTypeWelcome, welcomeMsg)

	// Broadcast attendee joined status to other clients in the session
	statusMsg := AttendeeStatusMessage{
		AttendeeID:   client.GetAttendeeID(),
		AttendeeName: attendee.Name,
		Status:       AttendeeStatusJoined,
		SessionID:    sessionID,
	}

	// Broadcast to all other clients in the session
	message, _ := CreateMessage(MessageTypeAttendeeStatus, statusMsg)
	h.BroadcastToSessionExcept(sessionID, message, client)

	log.Printf("Client %d (%s) joined session %d", client.GetAttendeeID(), attendee.Name, sessionID)
}

// NotifyVoteSubmitted notifies all clients in a session about a vote submission
func (h *Hub) NotifyVoteSubmitted(sessionID int, voteUpdate VoteUpdateMessage) {
	message, err := CreateMessage(MessageTypeVoteUpdate, voteUpdate)
	if err != nil {
		log.Printf("Failed to create vote update message: %v", err)
		return
	}

	h.BroadcastToSession(sessionID, message)
	h.updateStats(false, true)
}

// NotifyConsensusReached notifies all clients in a session about consensus
func (h *Hub) NotifyConsensusReached(sessionID int, consensus ConsensusReachedMessage) {
	message, err := CreateMessage(MessageTypeConsensusReached, consensus)
	if err != nil {
		log.Printf("Failed to create consensus message: %v", err)
		return
	}

	h.BroadcastToSession(sessionID, message)
	h.updateStats(false, true)
}

// NotifySessionProgress notifies all clients about session progress
func (h *Hub) NotifySessionProgress(sessionID int, progress SessionProgressMessage) {
	message, err := CreateMessage(MessageTypeSessionProgress, progress)
	if err != nil {
		log.Printf("Failed to create progress message: %v", err)
		return
	}

	h.BroadcastToSession(sessionID, message)
	h.updateStats(false, true)
}

// NotifySessionCompleted notifies all clients that a session is completed
func (h *Hub) NotifySessionCompleted(sessionID int, completion SessionCompletedMessage) {
	message, err := CreateMessage(MessageTypeSessionCompleted, completion)
	if err != nil {
		log.Printf("Failed to create session completed message: %v", err)
		return
	}

	h.BroadcastToSession(sessionID, message)
	h.updateStats(false, true)
}

// registerClient handles client registration
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessionID := client.GetSessionID()

	// Initialize session map if it doesn't exist
	if h.sessions[sessionID] == nil {
		h.sessions[sessionID] = make(map[*Client]bool)
	}

	// Add client to session
	h.sessions[sessionID][client] = true

	h.updateStats(true, false)

	log.Printf("Client registered: session=%d, attendee=%d, total_clients=%d",
		sessionID, client.GetAttendeeID(), h.stats.ActiveConnections)
}

// unregisterClient handles client unregistration
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessionID := client.GetSessionID()

	if clients, exists := h.sessions[sessionID]; exists {
		if _, exists := clients[client]; exists {
			// Remove client from session
			delete(clients, client)

			// Clean up empty sessions
			if len(clients) == 0 {
				delete(h.sessions, sessionID)
			}

			h.updateStats(false, false)

			log.Printf("Client unregistered: session=%d, attendee=%d, total_clients=%d",
				sessionID, client.GetAttendeeID(), h.stats.ActiveConnections)

			// Notify other clients about attendee leaving
			go h.notifyAttendeeLeft(sessionID, client)
		}
	}
}

// broadcastMessage handles message broadcasting
func (h *Hub) broadcastMessage(broadcast *BroadcastMessage) {
	h.mu.RLock()
	clients, exists := h.sessions[broadcast.SessionID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	messagesSent := 0
	for client := range clients {
		// Skip excluded client if specified
		if broadcast.ExcludeClient != nil && client == broadcast.ExcludeClient {
			continue
		}

		if client.IsConnected() {
			client.Send(broadcast.Message)
			messagesSent++
		}
	}

	if messagesSent > 0 {
		log.Printf("Broadcast message %s to session %d (%d clients)",
			broadcast.Message.Type, broadcast.SessionID, messagesSent)
	}
}

// notifyAttendeeLeft notifies other clients that an attendee has left
func (h *Hub) notifyAttendeeLeft(sessionID int, leftClient *Client) {
	attendee, err := h.attendeeRepo.GetByID(leftClient.GetAttendeeID())
	if err != nil {
		log.Printf("Failed to get attendee %d for leave notification: %v", leftClient.GetAttendeeID(), err)
		return
	}

	statusMsg := AttendeeStatusMessage{
		AttendeeID:   leftClient.GetAttendeeID(),
		AttendeeName: attendee.Name,
		Status:       AttendeeStatusLeft,
		SessionID:    sessionID,
	}

	message, err := CreateMessage(MessageTypeAttendeeStatus, statusMsg)
	if err != nil {
		log.Printf("Failed to create attendee left message: %v", err)
		return
	}

	h.BroadcastToSession(sessionID, message)
}

// getSessionClientCount returns the number of clients in a session
func (h *Hub) getSessionClientCount(sessionID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.sessions[sessionID]; exists {
		return len(clients)
	}
	return 0
}

// GetSessionClients returns the list of clients in a session
func (h *Hub) GetSessionClients(sessionID int) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	if sessionClients, exists := h.sessions[sessionID]; exists {
		for client := range sessionClients {
			if client.IsConnected() {
				clients = append(clients, client)
			}
		}
	}
	return clients
}

// GetStats returns hub statistics
func (h *Hub) GetStats() HubStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.stats
}

// updateStats updates hub statistics
func (h *Hub) updateStats(connectionChange bool, messageSent bool) {
	if connectionChange {
		h.stats.ActiveConnections = 0
		h.stats.ActiveSessions = len(h.sessions)

		for _, clients := range h.sessions {
			h.stats.ActiveConnections += len(clients)
		}

		h.stats.TotalConnections++
	}

	if messageSent {
		h.stats.MessagesSent++
	}

	h.stats.LastActivity = time.Now()
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Shutting down WebSocket Hub...")

	// Close all client connections
	for sessionID, clients := range h.sessions {
		for client := range clients {
			client.SendJSON(MessageTypeError, ErrorMessage{
				Code:    503,
				Message: "Server shutting down",
			})
			time.Sleep(100 * time.Millisecond) // Give time for message to send
		}
		delete(h.sessions, sessionID)
	}

	log.Println("WebSocket Hub shutdown complete")
}
