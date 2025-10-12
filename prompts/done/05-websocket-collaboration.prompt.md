# Prompt 5: WebSocket Collaboration Layer

Add WebSocket support for real-time collaboration during scoring sessions. Implement session broadcasting, live vote updates, and consensus status synchronization across all connected attendees.

## Requirements
- Set up WebSocket server using Gorilla WebSocket
- Implement session-based connection management
- Create real-time event broadcasting system
- Add connection authentication and session validation
- Handle connection cleanup and error recovery
- Create WebSocket message types for different events

## WebSocket Message Types
```json
// Join session
{
  "type": "join_session",
  "session_id": "123",
  "attendee_id": "456"
}

// Vote submitted
{
  "type": "vote_submitted",
  "comparison_id": "789",
  "attendee_id": "456", 
  "preferred_feature_id": "101"
}

// Consensus reached
{
  "type": "consensus_reached",
  "comparison_id": "789",
  "winner_id": "101"
}

// Session progress update
{
  "type": "session_progress",
  "session_id": "123",
  "completed_comparisons": 5,
  "total_comparisons": 10
}

// Attendee joined/left
{
  "type": "attendee_status",
  "attendee_id": "456",
  "status": "joined|left"
}
```

## Implementation Structure
- `internal/websocket/hub.go` - Connection management
- `internal/websocket/client.go` - Individual client handling
- `internal/websocket/message.go` - Message type definitions
- WebSocket endpoint: `WS /api/projects/{id}/sessions/{session_id}/ws`

## Features to Implement
- Session-based connection rooms
- Real-time vote broadcasting to all session participants
- Consensus status updates when agreements are reached
- Session progress notifications
- Attendee presence indicators (who's online)
- Connection recovery handling
- Rate limiting and message validation

## Connection Management
- Authenticate attendee on WebSocket connection
- Validate attendee belongs to the project
- Maintain active connections per session
- Clean up connections on disconnect
- Handle reconnection scenarios