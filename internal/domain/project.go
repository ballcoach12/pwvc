package domain

import (
	"time"
)

// Project represents a P-WVC project
type Project struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required,min=1,max=255"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	InviteCode  string    `json:"invite_code,omitempty" db:"invite_code"` // For T016 - US1 invite links
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateProjectRequest represents the request payload for creating a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

// UpdateProjectRequest represents the request payload for updating a project
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive completed"`
}

// ProjectProgress represents the workflow state of a P-WVC project
type ProjectProgress struct {
	ProjectID                    int       `json:"project_id" db:"project_id"`
	SetupCompleted               bool      `json:"setup_completed" db:"setup_completed"`
	AttendeesAdded               bool      `json:"attendees_added" db:"attendees_added"`
	FeaturesAdded                bool      `json:"features_added" db:"features_added"`
	PairwiseValueCompleted       bool      `json:"pairwise_value_completed" db:"pairwise_value_completed"`
	PairwiseComplexityCompleted  bool      `json:"pairwise_complexity_completed" db:"pairwise_complexity_completed"`
	FibonacciValueCompleted      bool      `json:"fibonacci_value_completed" db:"fibonacci_value_completed"`
	FibonacciComplexityCompleted bool      `json:"fibonacci_complexity_completed" db:"fibonacci_complexity_completed"`
	ResultsCalculated            bool      `json:"results_calculated" db:"results_calculated"`
	CurrentPhase                 string    `json:"current_phase" db:"current_phase"`
	LastActivity                 time.Time `json:"last_activity" db:"last_activity"`
	CreatedAt                    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at" db:"updated_at"`
}

// WorkflowPhase represents the different phases in the P-WVC workflow
type WorkflowPhase string

const (
	PhaseSetup               WorkflowPhase = "setup"
	PhaseAttendees           WorkflowPhase = "attendees"
	PhaseFeatures            WorkflowPhase = "features"
	PhasePairwiseValue       WorkflowPhase = "pairwise_value"
	PhasePairwiseComplexity  WorkflowPhase = "pairwise_complexity"
	PhaseFibonacciValue      WorkflowPhase = "fibonacci_value"
	PhaseFibonacciComplexity WorkflowPhase = "fibonacci_complexity"
	PhaseResults             WorkflowPhase = "results"
)

// GetNextPhase returns the next phase in the workflow
func (p *ProjectProgress) GetNextPhase() WorkflowPhase {
	switch p.CurrentPhase {
	case string(PhaseSetup):
		return PhaseAttendees
	case string(PhaseAttendees):
		return PhaseFeatures
	case string(PhaseFeatures):
		return PhasePairwiseValue
	case string(PhasePairwiseValue):
		return PhasePairwiseComplexity
	case string(PhasePairwiseComplexity):
		return PhaseFibonacciValue
	case string(PhaseFibonacciValue):
		return PhaseFibonacciComplexity
	case string(PhaseFibonacciComplexity):
		return PhaseResults
	default:
		return PhaseResults
	}
}

// CanProgressTo checks if the project can progress to a specific phase
func (p *ProjectProgress) CanProgressTo(phase WorkflowPhase) bool {
	switch phase {
	case PhaseSetup:
		return true
	case PhaseAttendees:
		return p.SetupCompleted
	case PhaseFeatures:
		return p.SetupCompleted && p.AttendeesAdded
	case PhasePairwiseValue:
		return p.SetupCompleted && p.AttendeesAdded && p.FeaturesAdded
	case PhasePairwiseComplexity:
		return p.PairwiseValueCompleted
	case PhaseFibonacciValue:
		return p.PairwiseComplexityCompleted
	case PhaseFibonacciComplexity:
		return p.FibonacciValueCompleted
	case PhaseResults:
		return p.FibonacciComplexityCompleted
	default:
		return false
	}
}
