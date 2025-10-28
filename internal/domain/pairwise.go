package domain

import (
	"time"
)

// CriterionType represents the type of criterion being compared
type CriterionType string

const (
	CriterionTypeValue      CriterionType = "value"
	CriterionTypeComplexity CriterionType = "complexity"
)

// SessionStatus represents the status of a pairwise session
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
)

// PairwiseSession represents a pairwise comparison session
type PairwiseSession struct {
	ID            int           `json:"id" db:"id"`
	ProjectID     int           `json:"project_id" db:"project_id"`
	CriterionType CriterionType `json:"criterion_type" db:"criterion_type"`
	Status        SessionStatus `json:"status" db:"status"`
	StartedAt     time.Time     `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty" db:"completed_at"`
}

// ConsensusScore represents a locked consensus score for a feature
type ConsensusScore struct {
	ID          int       `json:"id" db:"id"`
	ProjectID   int       `json:"project_id" db:"project_id"`
	FeatureID   int       `json:"feature_id" db:"feature_id"`
	SValue      int       `json:"s_value" db:"s_value"`
	SComplexity int       `json:"s_complexity" db:"s_complexity"`
	LockedBy    int       `json:"locked_by" db:"locked_by"`
	LockedAt    time.Time `json:"locked_at" db:"locked_at"`
	Rationale   string    `json:"rationale,omitempty" db:"rationale"`

	// Populated via joins
	Feature     *Feature  `json:"feature,omitempty" gorm:"foreignKey:FeatureID"`
	Facilitator *Attendee `json:"facilitator,omitempty" gorm:"foreignKey:LockedBy"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         int       `json:"id" db:"id"`
	ProjectID  int       `json:"project_id" db:"project_id"`
	AttendeeID int       `json:"attendee_id" db:"attendee_id"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id" db:"entity_id"`
	OldValue   string    `json:"old_value,omitempty" db:"old_value"` // JSON string representation
	NewValue   string    `json:"new_value,omitempty" db:"new_value"` // JSON string representation
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Metadata   string    `json:"metadata,omitempty" db:"metadata"` // JSON string representation
}

// TableName returns the table name for GORM
func (PairwiseSession) TableName() string {
	return "pairwise_sessions"
}

// SessionComparison represents a comparison between two features in a session
type SessionComparison struct {
	ID               int       `json:"id" db:"id"`
	SessionID        int       `json:"session_id" db:"session_id"`
	FeatureAID       int       `json:"feature_a_id" db:"feature_a_id"`
	FeatureBID       int       `json:"feature_b_id" db:"feature_b_id"`
	WinnerID         *int      `json:"winner_id,omitempty" db:"winner_id"`
	IsTie            bool      `json:"is_tie" db:"is_tie"`
	ConsensusReached bool      `json:"consensus_reached" db:"consensus_reached"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`

	// Populated via joins
	FeatureA *Feature `json:"feature_a,omitempty"`
	FeatureB *Feature `json:"feature_b,omitempty"`
	Winner   *Feature `json:"winner,omitempty"`
}

// TableName returns the table name for GORM
func (SessionComparison) TableName() string {
	return "pairwise_comparisons"
}

// AttendeeVote represents an individual attendee's vote for a comparison
type AttendeeVote struct {
	ID                 int       `json:"id" db:"id"`
	ComparisonID       int       `json:"comparison_id" db:"comparison_id"`
	AttendeeID         int       `json:"attendee_id" db:"attendee_id"`
	PreferredFeatureID *int      `json:"preferred_feature_id,omitempty" db:"preferred_feature_id"`
	IsTieVote          bool      `json:"is_tie_vote" db:"is_tie_vote"`
	VotedAt            time.Time `json:"voted_at" db:"voted_at"`

	// Populated via joins
	Attendee         *Attendee `json:"attendee,omitempty"`
	PreferredFeature *Feature  `json:"preferred_feature,omitempty"`
}

// TableName returns the table name for GORM
func (AttendeeVote) TableName() string {
	return "attendee_votes"
}

// SessionProgress represents the progress of a pairwise session
type SessionProgress struct {
	SessionID            int     `json:"session_id"`
	TotalComparisons     int     `json:"total_comparisons"`
	CompletedComparisons int     `json:"completed_comparisons"`
	ProgressPercentage   float64 `json:"progress_percentage"`
	RemainingComparisons int     `json:"remaining_comparisons"`
}

// CreatePairwiseSessionRequest represents the request to start a new pairwise session
type CreatePairwiseSessionRequest struct {
	CriterionType CriterionType `json:"criterion_type" binding:"required,oneof=value complexity"`
}

// SubmitVoteRequest represents the request to submit an attendee vote
type SubmitVoteRequest struct {
	ComparisonID       int  `json:"comparison_id" binding:"required"`
	AttendeeID         int  `json:"attendee_id" binding:"required"`
	PreferredFeatureID *int `json:"preferred_feature_id,omitempty"`
	IsTieVote          bool `json:"is_tie_vote"`
}

// ComparisonWithVotes represents a comparison with all attendee votes
type ComparisonWithVotes struct {
	Comparison *SessionComparison `json:"comparison"`
	Votes      []AttendeeVote     `json:"votes"`
}

// FeaturePair represents a pair of features to be compared
type FeaturePair struct {
	FeatureA *Feature `json:"feature_a"`
	FeatureB *Feature `json:"feature_b"`
}

// FibonacciScore represents a Fibonacci scoring entry for a feature
type FibonacciScore struct {
	ID             int       `json:"id" db:"id"`
	FeatureID      int       `json:"feature_id" db:"feature_id"`
	AttendeeID     int       `json:"attendee_id" db:"attendee_id"`
	CriterionType  string    `json:"criterion_type" db:"criterion_type"`
	FibonacciValue int       `json:"fibonacci_value" db:"fibonacci_value"`
	Rationale      string    `json:"rationale,omitempty" db:"rationale"`
	SubmittedAt    time.Time `json:"submitted_at" db:"submitted_at"`

	// Populated via joins
	Feature  *Feature  `json:"feature,omitempty"`
	Attendee *Attendee `json:"attendee,omitempty"`
}

// TableName returns the table name for GORM
func (FibonacciScore) TableName() string {
	return "fibonacci_scores"
}

// SubmitScoreRequest represents the request to submit a Fibonacci score
type SubmitScoreRequest struct {
	FeatureID      int    `json:"feature_id" binding:"required"`
	AttendeeID     int    `json:"attendee_id" binding:"required"`
	FibonacciValue int    `json:"fibonacci_value" binding:"required"`
	Rationale      string `json:"rationale,omitempty"`
}

// LockConsensusRequest represents the request to lock consensus scores
type LockConsensusRequest struct {
	FeatureID   int    `json:"feature_id" binding:"required"`
	SValue      int    `json:"s_value" binding:"required"`
	SComplexity int    `json:"s_complexity" binding:"required"`
	Rationale   string `json:"rationale,omitempty"`
}

// UnlockConsensusRequest represents the request to unlock consensus scores
type UnlockConsensusRequest struct {
	FeatureID int `json:"feature_id" binding:"required"`
}

// ReassignmentRequest represents a request to reassign pending comparisons (T042 - US8)
type ReassignmentRequest struct {
	SessionID        int    `json:"session_id" binding:"required"`
	CriterionType    string `json:"criterion_type" binding:"required"`
	ComparisonIDs    []int  `json:"comparison_ids" binding:"required"`
	ReassignmentType string `json:"reassignment_type" binding:"required"` // "session", "reset", "priority"
	TargetSessionID  int    `json:"target_session_id,omitempty"`
	NewPriority      int    `json:"new_priority,omitempty"`
	Reason           string `json:"reason,omitempty"`
}

// ReassignmentOptions represents available options for comparison reassignment (T042 - US8)
type ReassignmentOptions struct {
	ProjectID            int                  `json:"project_id"`
	CurrentSessionID     int                  `json:"current_session_id"`
	AvailableSessions    []*PairwiseSession   `json:"available_sessions"`
	PendingComparisons   []*SessionComparison `json:"pending_comparisons"`
	ReassignmentTypes    []string             `json:"reassignment_types"`
	CanReassignToSession bool                 `json:"can_reassign_to_session"`
	CanResetVotes        bool                 `json:"can_reset_votes"`
	CanChangePriority    bool                 `json:"can_change_priority"`
}
