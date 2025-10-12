package service

import (
	"testing"

	"pwvc/internal/websocket"
)

// MockWebSocketBroadcaster implements WebSocketBroadcaster for testing
type MockWebSocketBroadcaster struct {
	VoteNotifications       []websocket.VoteUpdateMessage
	ConsensusNotifications  []websocket.ConsensusReachedMessage
	ProgressNotifications   []websocket.SessionProgressMessage
	CompletionNotifications []websocket.SessionCompletedMessage
}

func (m *MockWebSocketBroadcaster) NotifyVoteSubmitted(sessionID int, voteUpdate websocket.VoteUpdateMessage) {
	m.VoteNotifications = append(m.VoteNotifications, voteUpdate)
}

func (m *MockWebSocketBroadcaster) NotifyConsensusReached(sessionID int, consensus websocket.ConsensusReachedMessage) {
	m.ConsensusNotifications = append(m.ConsensusNotifications, consensus)
}

func (m *MockWebSocketBroadcaster) NotifySessionProgress(sessionID int, progress websocket.SessionProgressMessage) {
	m.ProgressNotifications = append(m.ProgressNotifications, progress)
}

func (m *MockWebSocketBroadcaster) NotifySessionCompleted(sessionID int, completion websocket.SessionCompletedMessage) {
	m.CompletionNotifications = append(m.CompletionNotifications, completion)
}

func TestPairwiseService_SetWebSocketBroadcaster(t *testing.T) {
	// Create mock broadcaster
	mockBroadcaster := &MockWebSocketBroadcaster{}

	// Create pairwise service (with nil arguments for this test)
	service := &PairwiseService{}

	// Set WebSocket broadcaster
	service.SetWebSocketBroadcaster(mockBroadcaster)

	// Verify that the service has the broadcaster set
	if service.wsBroadcaster == nil {
		t.Error("PairwiseService WebSocket broadcaster was not set")
	}

	// Verify it's the same instance
	if service.wsBroadcaster != mockBroadcaster {
		t.Error("PairwiseService WebSocket broadcaster is not the expected instance")
	}
}

func TestWebSocketBroadcasterInterface(t *testing.T) {
	// Create mock broadcaster
	mock := &MockWebSocketBroadcaster{}

	// Test that it implements the interface
	var broadcaster WebSocketBroadcaster = mock

	// Test each method
	voteUpdate := websocket.VoteUpdateMessage{
		ComparisonID: 1,
		AttendeeID:   1,
		AttendeeName: "Test User",
	}
	broadcaster.NotifyVoteSubmitted(1, voteUpdate)

	consensus := websocket.ConsensusReachedMessage{
		ComparisonID: 1,
		WinnerID:     &[]int{1}[0],
	}
	broadcaster.NotifyConsensusReached(1, consensus)

	progress := websocket.SessionProgressMessage{
		SessionID:            1,
		CompletedComparisons: 1,
		TotalComparisons:     10,
		ProgressPercentage:   10.0,
	}
	broadcaster.NotifySessionProgress(1, progress)

	completion := websocket.SessionCompletedMessage{
		SessionID:     1,
		CriterionType: "value",
	}
	broadcaster.NotifySessionCompleted(1, completion)

	// Verify notifications were received
	if len(mock.VoteNotifications) != 1 {
		t.Errorf("Expected 1 vote notification, got %d", len(mock.VoteNotifications))
	}

	if len(mock.ConsensusNotifications) != 1 {
		t.Errorf("Expected 1 consensus notification, got %d", len(mock.ConsensusNotifications))
	}

	if len(mock.ProgressNotifications) != 1 {
		t.Errorf("Expected 1 progress notification, got %d", len(mock.ProgressNotifications))
	}

	if len(mock.CompletionNotifications) != 1 {
		t.Errorf("Expected 1 completion notification, got %d", len(mock.CompletionNotifications))
	}
}
