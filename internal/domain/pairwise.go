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
