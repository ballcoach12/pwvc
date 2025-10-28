package service

import (
	"encoding/json"
	"fmt"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/utils"
)

type ScoringService struct {
	scoreRepo    repository.ScoringRepository
	featureRepo  repository.FeatureRepository
	attendeeRepo repository.AttendeeRepository
	auditRepo    repository.AuditRepository
}

func NewScoringService(scoreRepo repository.ScoringRepository, featureRepo repository.FeatureRepository,
	attendeeRepo repository.AttendeeRepository, auditRepo repository.AuditRepository) *ScoringService {
	return &ScoringService{
		scoreRepo:    scoreRepo,
		featureRepo:  featureRepo,
		attendeeRepo: attendeeRepo,
		auditRepo:    auditRepo,
	}
}

// SubmitValueScore handles value criterion Fibonacci scoring (T030 - US4)
func (s *ScoringService) SubmitValueScore(projectID int, req domain.SubmitScoreRequest) (*domain.FibonacciScore, error) {
	return s.submitScore(projectID, req, "value")
}

// SubmitComplexityScore handles complexity criterion Fibonacci scoring (T030 - US4)
func (s *ScoringService) SubmitComplexityScore(projectID int, req domain.SubmitScoreRequest) (*domain.FibonacciScore, error) {
	return s.submitScore(projectID, req, "complexity")
}

// submitScore is the core method for submitting Fibonacci scores
func (s *ScoringService) submitScore(projectID int, req domain.SubmitScoreRequest, criterionType string) (*domain.FibonacciScore, error) {
	// Validate Fibonacci value (T030 - US4)
	if !utils.IsValidFibonacci(req.FibonacciValue) {
		return nil, &domain.APIError{
			Code:    400,
			Message: "Invalid Fibonacci value",
			Details: fmt.Sprintf("Value %d is not in the valid Fibonacci sequence [1,2,3,5,8,13,21,34,55,89]", req.FibonacciValue),
		}
	}

	// Verify feature exists in project
	feature, err := s.featureRepo.GetByID(req.FeatureID)
	if err != nil {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Feature not found",
		}
	}
	if feature.ProjectID != projectID {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Feature not found in this project",
		}
	}

	// Verify attendee exists in project
	attendee, err := s.attendeeRepo.GetByID(req.AttendeeID)
	if err != nil {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Attendee not found",
		}
	}
	if attendee.ProjectID != projectID {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Attendee not found in this project",
		}
	}

	// Check if score already exists (update vs create)
	existingScore, err := s.scoreRepo.GetByFeatureAndAttendee(req.FeatureID, req.AttendeeID, criterionType)

	var score *domain.FibonacciScore
	var auditAction string
	var oldValueJSON string

	if err != nil || existingScore == nil {
		// Create new score
		score = &domain.FibonacciScore{
			FeatureID:      req.FeatureID,
			AttendeeID:     req.AttendeeID,
			CriterionType:  criterionType,
			FibonacciValue: req.FibonacciValue,
			Rationale:      req.Rationale,
			SubmittedAt:    time.Now(),
		}

		err = s.scoreRepo.Create(score)
		if err != nil {
			return nil, err
		}

		auditAction = "score_created"
		oldValueJSON = ""
	} else {
		// Update existing score
		oldValue := map[string]interface{}{
			"fibonacci_value": existingScore.FibonacciValue,
			"rationale":       existingScore.Rationale,
		}
		oldValueBytes, _ := json.Marshal(oldValue)
		oldValueJSON = string(oldValueBytes)

		existingScore.FibonacciValue = req.FibonacciValue
		existingScore.Rationale = req.Rationale
		existingScore.SubmittedAt = time.Now()

		err = s.scoreRepo.Update(existingScore)
		if err != nil {
			return nil, err
		}

		score = existingScore
		auditAction = "score_updated"
	}

	// Load associated feature and attendee for response
	score.Feature = feature
	score.Attendee = attendee

	// Create audit log entry (T031 - US4)
	newValue := map[string]interface{}{
		"fibonacci_value": req.FibonacciValue,
		"rationale":       req.Rationale,
		"criterion_type":  criterionType,
	}
	newValueBytes, _ := json.Marshal(newValue)
	newValueJSON := string(newValueBytes)

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: req.AttendeeID,
		Action:     auditAction,
		EntityType: "fibonacci_score",
		EntityID:   fmt.Sprintf("%d-%d-%s", req.FeatureID, req.AttendeeID, criterionType),
		OldValue:   oldValueJSON,
		NewValue:   newValueJSON,
		Timestamp:  time.Now(),
	}

	// Log audit entry (non-blocking)
	go func() {
		s.auditRepo.Create(auditLog)
	}()

	return score, nil
}

// GetProjectScores retrieves all scores for a project (T030 - US4)
func (s *ScoringService) GetProjectScores(projectID int, criterionType string) ([]*domain.FibonacciScore, error) {
	scores, err := s.scoreRepo.GetByProject(projectID, criterionType)
	if err != nil {
		return nil, err
	}

	// Load associated features and attendees
	for _, score := range scores {
		if score.FeatureID > 0 {
			feature, err := s.featureRepo.GetByID(score.FeatureID)
			if err == nil {
				score.Feature = feature
			}
		}

		if score.AttendeeID > 0 {
			attendee, err := s.attendeeRepo.GetByID(score.AttendeeID)
			if err == nil {
				score.Attendee = attendee
			}
		}
	}

	return scores, nil
}
