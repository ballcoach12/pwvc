# Instructions: WebSocket Collaboration Layer Implementation

## WebSocket Architecture Patterns

### Hub-and-Spoke Model
```go
// internal/websocket/hub.go
type Hub struct {
    sessions   map[uint]*SessionRoom    // Session ID -> Room
    clients    map[*Client]bool         // All connected clients
    register   chan *Client             // Register requests from clients
    unregister chan *Client             // Unregister requests from clients
    broadcast  chan *BroadcastMessage   // Broadcast messages to specific sessions
    logger     *slog.Logger
    mutex      sync.RWMutex
}

type SessionRoom struct {
    SessionID uint                    `json:"session_id"`
    Clients   map[*Client]bool        `json:"-"`
    Messages  chan *SessionMessage    `json:"-"`
    mutex     sync.RWMutex
}

type Client struct {
    ID         string                  `json:"id"`
    SessionID  uint                    `json:"session_id"`
    AttendeeID uint                    `json:"attendee_id"`
    ProjectID  uint                    `json:"project_id"`
    Conn       *websocket.Conn         `json:"-"`
    Send       chan *SessionMessage    `json:"-"`
    Hub        *Hub                    `json:"-"`
    IsActive   bool                    `json:"is_active"`
    ConnectedAt time.Time              `json:"connected_at"`
}

func NewHub(logger *slog.Logger) *Hub {
    return &Hub{
        sessions:   make(map[uint]*SessionRoom),
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan *BroadcastMessage),
        logger:     logger,
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)
            
        case client := <-h.unregister:
            h.unregisterClient(client)
            
        case message := <-h.broadcast:
            h.broadcastToSession(message)
        }
    }
}
```

### Client Connection Management
```go
// internal/websocket/client.go
func (c *Client) readPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()
    
    // Set read limits and timeouts
    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    
    for {
        var message SessionMessage
        if err := c.Conn.ReadJSON(&message); err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                c.Hub.logger.Error("WebSocket error", "error", err)
            }
            break
        }
        
        // Add client metadata to message
        message.AttendeeID = c.AttendeeID
        message.SessionID = c.SessionID
        message.Timestamp = time.Now()
        
        // Process message based on type
        if err := c.processMessage(&message); err != nil {
            c.Hub.logger.Error("Failed to process message", "error", err)
            c.sendError("Failed to process message: " + err.Error())
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := c.Conn.WriteJSON(message); err != nil {
                c.Hub.logger.Error("Failed to write message", "error", err)
                return
            }
            
        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

## Message Type System

### Structured Message Protocol
```go
// internal/websocket/message.go
type MessageType string

const (
    // Connection management
    JoinSession    MessageType = "join_session"
    LeaveSession   MessageType = "leave_session"
    AttendeeJoined MessageType = "attendee_joined"
    AttendeeLeft   MessageType = "attendee_left"
    
    // Voting messages
    VoteSubmitted     MessageType = "vote_submitted"
    VoteReceived      MessageType = "vote_received"
    ConsensusReached  MessageType = "consensus_reached"
    ConsensusLost     MessageType = "consensus_lost"
    
    // Progress updates
    SessionProgress   MessageType = "session_progress"
    ComparisonUpdate  MessageType = "comparison_update"
    
    // System messages
    Error            MessageType = "error"
    SystemNotice     MessageType = "system_notice"
)

type SessionMessage struct {
    Type       MessageType            `json:"type"`
    SessionID  uint                   `json:"session_id"`
    AttendeeID uint                   `json:"attendee_id"`
    Data       map[string]interface{} `json:"data"`
    Timestamp  time.Time              `json:"timestamp"`
}

// Specific message data structures
type VoteSubmittedData struct {
    ComparisonID       uint  `json:"comparison_id"`
    PreferredFeatureID *uint `json:"preferred_feature_id"`
    IsTieVote          bool  `json:"is_tie_vote"`
    AttendeeName       string `json:"attendee_name"`
}

type ConsensusReachedData struct {
    ComparisonID       uint   `json:"comparison_id"`
    WinnerID           *uint  `json:"winner_id"`
    IsTie              bool   `json:"is_tie"`
    FeatureATitle      string `json:"feature_a_title"`
    FeatureBTitle      string `json:"feature_b_title"`
    WinnerTitle        string `json:"winner_title,omitempty"`
}

type SessionProgressData struct {
    TotalComparisons     int     `json:"total_comparisons"`
    CompletedComparisons int     `json:"completed_comparisons"`
    ProgressPercentage   float64 `json:"progress_percentage"`
    IsComplete           bool    `json:"is_complete"`
}
```

### Message Processing Pipeline
```go
// internal/websocket/message_processor.go
type MessageProcessor struct {
    pairwiseService  *service.PairwiseService
    consensusService *service.ConsensusService
    progressService  *service.SessionProgressService
    hub             *Hub
    logger          *slog.Logger
}

func (p *MessageProcessor) ProcessMessage(client *Client, message *SessionMessage) error {
    switch message.Type {
    case JoinSession:
        return p.handleJoinSession(client, message)
        
    case VoteSubmitted:
        return p.handleVoteSubmitted(client, message)
        
    default:
        return fmt.Errorf("unknown message type: %s", message.Type)
    }
}

func (p *MessageProcessor) handleVoteSubmitted(client *Client, message *SessionMessage) error {
    // Parse vote data
    voteData, ok := message.Data["vote"].(map[string]interface{})
    if !ok {
        return errors.New("invalid vote data")
    }
    
    vote := &domain.AttendeeVote{
        ComparisonID:       uint(voteData["comparison_id"].(float64)),
        AttendeeID:         client.AttendeeID,
        IsTieVote:          voteData["is_tie_vote"].(bool),
    }
    
    if !vote.IsTieVote {
        featureID := uint(voteData["preferred_feature_id"].(float64))
        vote.PreferredFeatureID = &featureID
    }
    
    // Process vote through consensus service
    result, err := p.consensusService.ProcessVote(vote)
    if err != nil {
        return fmt.Errorf("failed to process vote: %w", err)
    }
    
    // Broadcast vote received to all session participants
    p.broadcastVoteReceived(client.SessionID, vote, client.AttendeeID)
    
    // If consensus reached, broadcast that too
    if result.ConsensusReached {
        p.broadcastConsensusReached(client.SessionID, result)
        
        // Update session progress
        progress, err := p.progressService.GetSessionProgress(client.SessionID)
        if err == nil {
            p.broadcastSessionProgress(client.SessionID, progress)
        }
    }
    
    return nil
}
```

## Real-time Event Broadcasting

### Session-Scoped Broadcasting
```go
// internal/websocket/broadcaster.go
type Broadcaster struct {
    hub    *Hub
    logger *slog.Logger
}

func (b *Broadcaster) BroadcastToSession(sessionID uint, message *SessionMessage) {
    b.hub.mutex.RLock()
    sessionRoom, exists := b.hub.sessions[sessionID]
    b.hub.mutex.RUnlock()
    
    if !exists {
        b.logger.Warn("Attempted to broadcast to non-existent session", "session_id", sessionID)
        return
    }
    
    sessionRoom.mutex.RLock()
    clients := make([]*Client, 0, len(sessionRoom.Clients))
    for client := range sessionRoom.Clients {
        if client.IsActive {
            clients = append(clients, client)
        }
    }
    sessionRoom.mutex.RUnlock()
    
    // Send to all active clients in parallel
    var wg sync.WaitGroup
    for _, client := range clients {
        wg.Add(1)
        go func(c *Client) {
            defer wg.Done()
            select {
            case c.Send <- message:
                // Message sent successfully
            case <-time.After(5 * time.Second):
                b.logger.Warn("Failed to send message to client (timeout)", 
                    "client_id", c.ID, "session_id", sessionID)
            }
        }(client)
    }
    
    wg.Wait()
    
    b.logger.Debug("Message broadcasted", 
        "session_id", sessionID, 
        "message_type", message.Type,
        "recipients", len(clients))
}

func (b *Broadcaster) BroadcastVoteReceived(sessionID uint, vote *domain.AttendeeVote, attendeeName string) {
    message := &SessionMessage{
        Type:      VoteReceived,
        SessionID: sessionID,
        Data: map[string]interface{}{
            "comparison_id":        vote.ComparisonID,
            "attendee_id":         vote.AttendeeID,
            "attendee_name":       attendeeName,
            "preferred_feature_id": vote.PreferredFeatureID,
            "is_tie_vote":         vote.IsTieVote,
        },
        Timestamp: time.Now(),
    }
    
    b.BroadcastToSession(sessionID, message)
}

func (b *Broadcaster) BroadcastConsensusReached(sessionID uint, consensusResult *service.ConsensusResult) {
    message := &SessionMessage{
        Type:      ConsensusReached,
        SessionID: sessionID,
        Data: map[string]interface{}{
            "comparison_id": consensusResult.ComparisonID,
            "winner_id":     consensusResult.WinnerID,
            "is_tie":        consensusResult.IsTie,
        },
        Timestamp: time.Now(),
    }
    
    b.BroadcastToSession(sessionID, message)
}
```

## Connection Authentication and Authorization

### Session-based Authentication
```go
// internal/websocket/auth.go
type WebSocketAuthenticator struct {
    projectRepo  repository.ProjectRepository
    attendeeRepo repository.AttendeeRepository
    sessionRepo  repository.PairwiseSessionRepository
    logger       *slog.Logger
}

func (a *WebSocketAuthenticator) AuthenticateConnection(
    projectID, sessionID, attendeeID uint,
) error {
    // Verify attendee exists and belongs to project
    attendee, err := a.attendeeRepo.GetByID(attendeeID)
    if err != nil {
        return fmt.Errorf("attendee not found: %w", err)
    }
    
    if attendee.ProjectID != projectID {
        return errors.New("attendee does not belong to this project")
    }
    
    // Verify session exists and belongs to project
    session, err := a.sessionRepo.GetByID(sessionID)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }
    
    if session.ProjectID != projectID {
        return errors.New("session does not belong to this project")
    }
    
    // Verify session is active
    if session.Status != "active" {
        return errors.New("session is not active")
    }
    
    return nil
}

// WebSocket connection upgrade with authentication
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    projectID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    sessionID, _ := strconv.ParseUint(c.Param("session_id"), 10, 32)
    attendeeID, _ := strconv.ParseUint(c.Query("attendee_id"), 10, 32)
    
    // Authenticate connection
    if err := h.auth.AuthenticateConnection(uint(projectID), uint(sessionID), uint(attendeeID)); err != nil {
        c.JSON(401, gin.H{"error": "Authentication failed: " + err.Error()})
        return
    }
    
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        h.logger.Error("WebSocket upgrade failed", "error", err)
        return
    }
    
    // Create client
    client := &Client{
        ID:          generateClientID(),
        SessionID:   uint(sessionID),
        AttendeeID:  uint(attendeeID),
        ProjectID:   uint(projectID),
        Conn:        conn,
        Send:        make(chan *SessionMessage, 256),
        Hub:         h.hub,
        IsActive:    true,
        ConnectedAt: time.Now(),
    }
    
    // Register client and start pumps
    h.hub.register <- client
    
    go client.writePump()
    go client.readPump()
}
```

## Error Handling and Recovery

### Connection Recovery Patterns
```go
// internal/websocket/recovery.go
type ConnectionRecovery struct {
    hub    *Hub
    logger *slog.Logger
}

func (r *ConnectionRecovery) HandleClientDisconnection(client *Client, reason string) {
    r.logger.Info("Client disconnected", 
        "client_id", client.ID, 
        "session_id", client.SessionID,
        "reason", reason)
    
    // Mark client as inactive
    client.IsActive = false
    
    // Notify other session participants
    message := &SessionMessage{
        Type:      AttendeeLeft,
        SessionID: client.SessionID,
        Data: map[string]interface{}{
            "attendee_id": client.AttendeeID,
            "reason":      reason,
        },
        Timestamp: time.Now(),
    }
    
    // Broadcast to remaining clients
    broadcaster := &Broadcaster{hub: r.hub, logger: r.logger}
    broadcaster.BroadcastToSession(client.SessionID, message)
    
    // Clean up client resources
    close(client.Send)
}

func (r *ConnectionRecovery) HandleSessionCleanup(sessionID uint) {
    r.hub.mutex.Lock()
    defer r.hub.mutex.Unlock()
    
    sessionRoom, exists := r.hub.sessions[sessionID]
    if !exists {
        return
    }
    
    // Disconnect all clients in session
    for client := range sessionRoom.Clients {
        if client.IsActive {
            client.Conn.Close()
        }
    }
    
    // Remove session room
    delete(r.hub.sessions, sessionID)
    
    r.logger.Info("Session cleaned up", "session_id", sessionID)
}
```

## Performance and Scaling Patterns

### Message Rate Limiting
```go
// internal/websocket/rate_limiter.go
type RateLimiter struct {
    clients map[string]*clientLimiter
    mutex   sync.RWMutex
}

type clientLimiter struct {
    tokens    int
    lastRefill time.Time
    maxTokens int
    refillRate time.Duration
}

func (r *RateLimiter) AllowMessage(clientID string) bool {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    limiter, exists := r.clients[clientID]
    if !exists {
        limiter = &clientLimiter{
            tokens:     10, // Initial tokens
            maxTokens:  10,
            refillRate: time.Second,
            lastRefill: time.Now(),
        }
        r.clients[clientID] = limiter
    }
    
    // Refill tokens based on time elapsed
    now := time.Now()
    elapsed := now.Sub(limiter.lastRefill)
    tokensToAdd := int(elapsed / limiter.refillRate)
    
    if tokensToAdd > 0 {
        limiter.tokens = min(limiter.maxTokens, limiter.tokens+tokensToAdd)
        limiter.lastRefill = now
    }
    
    // Check if message is allowed
    if limiter.tokens > 0 {
        limiter.tokens--
        return true
    }
    
    return false
}
```

## Testing WebSocket Functionality

### WebSocket Integration Tests
```go
// internal/websocket/hub_test.go
func TestWebSocketHub_VoteProcessing(t *testing.T) {
    // Setup test hub
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))
    hub := NewHub(logger)
    go hub.Run()
    defer hub.Stop()
    
    // Create mock clients
    client1 := createMockClient(hub, 1, 1, 1) // attendee 1
    client2 := createMockClient(hub, 1, 1, 2) // attendee 2
    
    hub.register <- client1
    hub.register <- client2
    
    // Wait for registration
    time.Sleep(100 * time.Millisecond)
    
    // Simulate vote submission
    voteMessage := &SessionMessage{
        Type:      VoteSubmitted,
        SessionID: 1,
        Data: map[string]interface{}{
            "vote": map[string]interface{}{
                "comparison_id":        1,
                "preferred_feature_id": 2.0,
                "is_tie_vote":          false,
            },
        },
    }
    
    // Process vote through client
    err := client1.processMessage(voteMessage)
    assert.NoError(t, err)
    
    // Verify message was broadcasted to other clients
    select {
    case receivedMessage := <-client2.Send:
        assert.Equal(t, VoteReceived, receivedMessage.Type)
        assert.Equal(t, uint(1), receivedMessage.Data["comparison_id"])
    case <-time.After(5 * time.Second):
        t.Fatal("Expected to receive vote message but timed out")
    }
}
```