package websocket

import (
	"encoding/json"
	"time"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client to server messages
	MessageTypeJoinSession   MessageType = "join_session"
	MessageTypeLeaveSession  MessageType = "leave_session"
	MessageTypeVoteSubmitted MessageType = "vote_submitted"

	// Server to client messages
	MessageTypeConsensusReached MessageType = "consensus_reached"
	MessageTypeSessionProgress  MessageType = "session_progress"
	MessageTypeAttendeeStatus   MessageType = "attendee_status"
	MessageTypeVoteUpdate       MessageType = "vote_update"
	MessageTypeSessionCompleted MessageType = "session_completed"
	MessageTypeError            MessageType = "error"
	MessageTypeWelcome          MessageType = "welcome"
)

// AttendeeStatus represents the status of an attendee
type AttendeeStatus string

const (
	AttendeeStatusJoined AttendeeStatus = "joined"
	AttendeeStatusLeft   AttendeeStatus = "left"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType     `json:"type"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	ID        string          `json:"id,omitempty"` // For message tracking
}

// JoinSessionMessage represents a request to join a session
type JoinSessionMessage struct {
	SessionID  int `json:"session_id"`
	AttendeeID int `json:"attendee_id"`
}

// LeaveSessionMessage represents a request to leave a session
type LeaveSessionMessage struct {
	SessionID  int `json:"session_id"`
	AttendeeID int `json:"attendee_id"`
}

// VoteSubmittedMessage represents a vote submission notification
type VoteSubmittedMessage struct {
	ComparisonID       int    `json:"comparison_id"`
	AttendeeID         int    `json:"attendee_id"`
	PreferredFeatureID *int   `json:"preferred_feature_id,omitempty"`
	IsTieVote          bool   `json:"is_tie_vote"`
	AttendeeName       string `json:"attendee_name,omitempty"`
}

// ConsensusReachedMessage represents consensus achievement notification
type ConsensusReachedMessage struct {
	ComparisonID int    `json:"comparison_id"`
	WinnerID     *int   `json:"winner_id,omitempty"`
	IsTie        bool   `json:"is_tie"`
	FeatureAName string `json:"feature_a_name,omitempty"`
	FeatureBName string `json:"feature_b_name,omitempty"`
	WinnerName   string `json:"winner_name,omitempty"`
}

// SessionProgressMessage represents session progress update
type SessionProgressMessage struct {
	SessionID            int     `json:"session_id"`
	CompletedComparisons int     `json:"completed_comparisons"`
	TotalComparisons     int     `json:"total_comparisons"`
	ProgressPercentage   float64 `json:"progress_percentage"`
	RemainingComparisons int     `json:"remaining_comparisons"`
}

// AttendeeStatusMessage represents attendee presence update
type AttendeeStatusMessage struct {
	AttendeeID   int            `json:"attendee_id"`
	AttendeeName string         `json:"attendee_name"`
	Status       AttendeeStatus `json:"status"`
	SessionID    int            `json:"session_id"`
}

// VoteUpdateMessage represents a real-time vote update
type VoteUpdateMessage struct {
	ComparisonID       int    `json:"comparison_id"`
	AttendeeID         int    `json:"attendee_id"`
	AttendeeName       string `json:"attendee_name"`
	PreferredFeatureID *int   `json:"preferred_feature_id,omitempty"`
	IsTieVote          bool   `json:"is_tie_vote"`
	VotesReceived      int    `json:"votes_received"`
	TotalAttendees     int    `json:"total_attendees"`
	ConsensusReached   bool   `json:"consensus_reached"`
}

// SessionCompletedMessage represents session completion notification
type SessionCompletedMessage struct {
	SessionID      int       `json:"session_id"`
	CriterionType  string    `json:"criterion_type"`
	CompletedAt    time.Time `json:"completed_at"`
	TotalVotes     int       `json:"total_votes"`
	TotalConsensus int       `json:"total_consensus"`
}

// ErrorMessage represents an error notification
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WelcomeMessage represents a welcome message after successful connection
type WelcomeMessage struct {
	SessionID      int    `json:"session_id"`
	AttendeeID     int    `json:"attendee_id"`
	AttendeeName   string `json:"attendee_name"`
	CriterionType  string `json:"criterion_type"`
	ConnectedCount int    `json:"connected_count"`
	SessionStatus  string `json:"session_status"`
}

// CreateMessage creates a new WebSocket message
func CreateMessage(msgType MessageType, data interface{}) (*Message, error) {
	var rawData json.RawMessage
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		rawData = jsonData
	}

	return &Message{
		Type:      msgType,
		Data:      rawData,
		Timestamp: time.Now(),
	}, nil
}

// ParseMessageData parses the message data into the specified type
func (m *Message) ParseMessageData(target interface{}) error {
	if m.Data == nil {
		return nil
	}
	return json.Unmarshal(m.Data, target)
}

// MessageBuilder provides a fluent interface for building messages
type MessageBuilder struct {
	msgType MessageType
	data    interface{}
	id      string
}

// NewMessageBuilder creates a new message builder
func NewMessageBuilder(msgType MessageType) *MessageBuilder {
	return &MessageBuilder{msgType: msgType}
}

// WithData sets the message data
func (mb *MessageBuilder) WithData(data interface{}) *MessageBuilder {
	mb.data = data
	return mb
}

// WithID sets the message ID
func (mb *MessageBuilder) WithID(id string) *MessageBuilder {
	mb.id = id
	return mb
}

// Build creates the final message
func (mb *MessageBuilder) Build() (*Message, error) {
	msg, err := CreateMessage(mb.msgType, mb.data)
	if err != nil {
		return nil, err
	}
	if mb.id != "" {
		msg.ID = mb.id
	}
	return msg, nil
}
