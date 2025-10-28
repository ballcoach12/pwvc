package service

import (
	"encoding/json"
	"fmt"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/utils"
)

type ConsensusService struct {
	consensusRepo repository.ConsensusRepository
	featureRepo   repository.FeatureRepository
	attendeeRepo  repository.AttendeeRepository
	auditRepo     repository.AuditRepository
}

func NewConsensusService(consensusRepo repository.ConsensusRepository, featureRepo repository.FeatureRepository,
	attendeeRepo repository.AttendeeRepository, auditRepo repository.AuditRepository) *ConsensusService {
	return &ConsensusService{
		consensusRepo: consensusRepo,
		featureRepo:   featureRepo,
		attendeeRepo:  attendeeRepo,
		auditRepo:     auditRepo,
	}
}

// LockConsensusScore locks consensus scores for a feature (T034 - US5)
func (s *ConsensusService) LockConsensusScore(projectID, featureID, facilitatorID, sValue, sComplexity int, rationale string) (*domain.ConsensusScore, error) {
	// Verify feature exists in project
	feature, err := s.featureRepo.GetByID(featureID)
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

	// Verify facilitator exists in project
	facilitator, err := s.attendeeRepo.GetByID(facilitatorID)
	if err != nil {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Facilitator not found",
		}
	}
	if facilitator.ProjectID != projectID {
		return nil, &domain.APIError{
			Code:    404,
			Message: "Facilitator not found in this project",
		}
	}
	if facilitator.Role != "facilitator" {
		return nil, &domain.APIError{
			Code:    403,
			Message: "Only facilitators can lock consensus scores",
		}
	}

	// Validate Fibonacci values (T034 - US5)
	if !utils.IsValidFibonacci(sValue) {
		return nil, &domain.APIError{
			Code:    400,
			Message: "Invalid S-Value",
			Details: fmt.Sprintf("Value %d is not in the valid Fibonacci sequence [1,2,3,5,8,13,21,34,55,89]", sValue),
		}
	}
	if !utils.IsValidFibonacci(sComplexity) {
		return nil, &domain.APIError{
			Code:    400,
			Message: "Invalid S-Complexity",
			Details: fmt.Sprintf("Value %d is not in the valid Fibonacci sequence [1,2,3,5,8,13,21,34,55,89]", sComplexity),
		}
	}

	// Check if consensus already exists
	existingConsensus, err := s.consensusRepo.GetByFeature(featureID)

	var consensus *domain.ConsensusScore
	var auditAction string
	var oldValueJSON string

	if err != nil || existingConsensus == nil {
		// Create new consensus
		consensus = &domain.ConsensusScore{
			ProjectID:   projectID,
			FeatureID:   featureID,
			SValue:      sValue,
			SComplexity: sComplexity,
			LockedBy:    facilitatorID,
			LockedAt:    time.Now(),
			Rationale:   rationale,
		}

		err = s.consensusRepo.Create(consensus)
		if err != nil {
			return nil, err
		}

		auditAction = "consensus_locked"
		oldValueJSON = ""
	} else {
		// Update existing consensus
		oldValue := map[string]interface{}{
			"s_value":      existingConsensus.SValue,
			"s_complexity": existingConsensus.SComplexity,
			"rationale":    existingConsensus.Rationale,
			"locked_by":    existingConsensus.LockedBy,
		}
		oldValueBytes, _ := json.Marshal(oldValue)
		oldValueJSON = string(oldValueBytes)

		existingConsensus.SValue = sValue
		existingConsensus.SComplexity = sComplexity
		existingConsensus.LockedBy = facilitatorID
		existingConsensus.LockedAt = time.Now()
		existingConsensus.Rationale = rationale

		err = s.consensusRepo.Update(existingConsensus)
		if err != nil {
			return nil, err
		}

		consensus = existingConsensus
		auditAction = "consensus_updated"
	}

	// Load associated feature and facilitator for response
	consensus.Feature = feature
	consensus.Facilitator = facilitator

	// Create audit log entry (T031 - US4)
	newValue := map[string]interface{}{
		"s_value":      sValue,
		"s_complexity": sComplexity,
		"rationale":    rationale,
		"locked_by":    facilitatorID,
	}
	newValueBytes, _ := json.Marshal(newValue)
	newValueJSON := string(newValueBytes)

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: facilitatorID,
		Action:     auditAction,
		EntityType: "consensus_score",
		EntityID:   fmt.Sprintf("%d", featureID),
		OldValue:   oldValueJSON,
		NewValue:   newValueJSON,
		Timestamp:  time.Now(),
	}

	// Log audit entry (non-blocking)
	go func() {
		s.auditRepo.Create(auditLog)
	}()

	return consensus, nil
}

// UnlockConsensusScore unlocks consensus scores for a feature (T034 - US5)
func (s *ConsensusService) UnlockConsensusScore(projectID, featureID, facilitatorID int) error {
	// Verify feature exists in project
	feature, err := s.featureRepo.GetByID(featureID)
	if err != nil {
		return &domain.APIError{
			Code:    404,
			Message: "Feature not found",
		}
	}
	if feature.ProjectID != projectID {
		return &domain.APIError{
			Code:    404,
			Message: "Feature not found in this project",
		}
	}

	// Verify facilitator exists and has permission
	facilitator, err := s.attendeeRepo.GetByID(facilitatorID)
	if err != nil {
		return &domain.APIError{
			Code:    404,
			Message: "Facilitator not found",
		}
	}
	if facilitator.ProjectID != projectID {
		return &domain.APIError{
			Code:    404,
			Message: "Facilitator not found in this project",
		}
	}
	if facilitator.Role != "facilitator" {
		return &domain.APIError{
			Code:    403,
			Message: "Only facilitators can unlock consensus scores",
		}
	}

	// Get existing consensus
	existingConsensus, err := s.consensusRepo.GetByFeature(featureID)
	if err != nil || existingConsensus == nil {
		return &domain.APIError{
			Code:    404,
			Message: "No consensus found for this feature",
		}
	}

	// Store old value for audit
	oldValue := map[string]interface{}{
		"s_value":      existingConsensus.SValue,
		"s_complexity": existingConsensus.SComplexity,
		"rationale":    existingConsensus.Rationale,
		"locked_by":    existingConsensus.LockedBy,
	}
	oldValueBytes, _ := json.Marshal(oldValue)
	oldValueJSON := string(oldValueBytes)

	// Delete consensus
	err = s.consensusRepo.DeleteByFeature(featureID)
	if err != nil {
		return err
	}

	// Create audit log entry
	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: facilitatorID,
		Action:     "consensus_unlocked",
		EntityType: "consensus_score",
		EntityID:   fmt.Sprintf("%d", featureID),
		OldValue:   oldValueJSON,
		NewValue:   "",
		Timestamp:  time.Now(),
	}

	// Log audit entry (non-blocking)
	go func() {
		s.auditRepo.Create(auditLog)
	}()

	return nil
}

// GetProjectConsensus retrieves all consensus scores for a project (T034 - US5)
func (s *ConsensusService) GetProjectConsensus(projectID int) ([]*domain.ConsensusScore, error) {
	consensus, err := s.consensusRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}

	// Load associated features and facilitators
	for _, cons := range consensus {
		if cons.FeatureID > 0 {
			feature, err := s.featureRepo.GetByID(cons.FeatureID)
			if err == nil {
				cons.Feature = feature
			}
		}

		if cons.LockedBy > 0 {
			facilitator, err := s.attendeeRepo.GetByID(cons.LockedBy)
			if err == nil {
				cons.Facilitator = facilitator
			}
		}
	}

	return consensus, nil
}
