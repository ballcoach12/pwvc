package service

import (
	"encoding/json"
	"fmt"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

type AuditService struct {
	auditRepo    repository.AuditRepository
	attendeeRepo repository.AttendeeRepository
	projectRepo  repository.ProjectRepository
}

func NewAuditService(auditRepo repository.AuditRepository, attendeeRepo repository.AttendeeRepository, projectRepo repository.ProjectRepository) *AuditService {
	return &AuditService{
		auditRepo:    auditRepo,
		attendeeRepo: attendeeRepo,
		projectRepo:  projectRepo,
	}
}

// LogVoteAction logs a pairwise vote action with privacy controls (T044 - US9)
func (s *AuditService) LogVoteAction(projectID, attendeeID, comparisonID int, preferredFeatureID *int, isTie bool) error {
	// Determine if anonymity should be preserved based on project settings
	anonymize, err := s.shouldAnonymizeForProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to check anonymity settings: %w", err)
	}

	var actorID int
	var actorInfo interface{}

	if anonymize {
		actorID = 0 // Anonymous
		actorInfo = map[string]interface{}{
			"type": "anonymous_attendee",
		}
	} else {
		actorID = attendeeID
		attendee, err := s.attendeeRepo.GetByID(attendeeID)
		if err == nil {
			actorInfo = map[string]interface{}{
				"attendee_id":   attendeeID,
				"attendee_name": attendee.Name,
				"attendee_role": attendee.Role,
			}
		}
	}

	newValue := map[string]interface{}{
		"comparison_id":        comparisonID,
		"preferred_feature_id": preferredFeatureID,
		"is_tie":               isTie,
		"actor":                actorInfo,
	}
	newValueBytes, _ := json.Marshal(newValue)

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: actorID,
		Action:     "pairwise_vote_submitted",
		EntityType: "pairwise_comparison",
		EntityID:   fmt.Sprintf("%d", comparisonID),
		OldValue:   "", // No previous state for new votes
		NewValue:   string(newValueBytes),
		Timestamp:  time.Now(),
	}

	return s.auditRepo.Create(auditLog)
}

// LogScoreAction logs a Fibonacci score submission with privacy controls (T044 - US9)
func (s *AuditService) LogScoreAction(projectID, attendeeID, featureID int, criterionType string, fibonacciValue int, rationale string) error {
	anonymize, err := s.shouldAnonymizeForProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to check anonymity settings: %w", err)
	}

	var actorID int
	var actorInfo interface{}

	if anonymize {
		actorID = 0
		actorInfo = map[string]interface{}{
			"type": "anonymous_attendee",
		}
	} else {
		actorID = attendeeID
		attendee, err := s.attendeeRepo.GetByID(attendeeID)
		if err == nil {
			actorInfo = map[string]interface{}{
				"attendee_id":   attendeeID,
				"attendee_name": attendee.Name,
				"attendee_role": attendee.Role,
			}
		}
	}

	newValue := map[string]interface{}{
		"feature_id":      featureID,
		"criterion_type":  criterionType,
		"fibonacci_value": fibonacciValue,
		"rationale":       rationale,
		"actor":           actorInfo,
	}
	newValueBytes, _ := json.Marshal(newValue)

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: actorID,
		Action:     "fibonacci_score_submitted",
		EntityType: "fibonacci_score",
		EntityID:   fmt.Sprintf("%d-%s", featureID, criterionType),
		OldValue:   "",
		NewValue:   string(newValueBytes),
		Timestamp:  time.Now(),
	}

	return s.auditRepo.Create(auditLog)
}

// LogConsensusAction logs consensus lock/unlock actions (T044 - US9)
func (s *AuditService) LogConsensusAction(projectID, facilitatorID, featureID int, action string, sValue, sComplexity int, rationale string) error {
	// Facilitator actions are never anonymized for accountability
	facilitator, err := s.attendeeRepo.GetByID(facilitatorID)
	var facilitatorInfo interface{}
	if err == nil {
		facilitatorInfo = map[string]interface{}{
			"facilitator_id":   facilitatorID,
			"facilitator_name": facilitator.Name,
			"facilitator_role": facilitator.Role,
		}
	}

	var newValue interface{}
	if action == "consensus_locked" {
		newValue = map[string]interface{}{
			"feature_id":   featureID,
			"s_value":      sValue,
			"s_complexity": sComplexity,
			"rationale":    rationale,
			"facilitator":  facilitatorInfo,
		}
	} else {
		newValue = map[string]interface{}{
			"feature_id":  featureID,
			"facilitator": facilitatorInfo,
		}
	}

	newValueBytes, _ := json.Marshal(newValue)

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: facilitatorID,
		Action:     action,
		EntityType: "consensus_score",
		EntityID:   fmt.Sprintf("%d", featureID),
		OldValue:   "", // Could be enhanced to capture previous consensus state
		NewValue:   string(newValueBytes),
		Timestamp:  time.Now(),
	}

	return s.auditRepo.Create(auditLog)
}

// LogPhaseChangeAction logs project phase transitions (T044 - US9)
func (s *AuditService) LogPhaseChangeAction(projectID, facilitatorID int, oldPhase, newPhase string) error {
	facilitator, err := s.attendeeRepo.GetByID(facilitatorID)
	var facilitatorInfo interface{}
	if err == nil {
		facilitatorInfo = map[string]interface{}{
			"facilitator_id":   facilitatorID,
			"facilitator_name": facilitator.Name,
			"facilitator_role": facilitator.Role,
		}
	}

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: facilitatorID,
		Action:     "phase_changed",
		EntityType: "project_progress",
		EntityID:   fmt.Sprintf("%d", projectID),
		OldValue:   fmt.Sprintf(`{"phase":"%s"}`, oldPhase),
		NewValue:   fmt.Sprintf(`{"phase":"%s","facilitator":"%s"}`, newPhase, facilitatorInfo),
		Timestamp:  time.Now(),
	}

	return s.auditRepo.Create(auditLog)
}

// GetAuditReport generates an audit report with privacy controls (T044 - US9)
func (s *AuditService) GetAuditReport(projectID int, options domain.AuditReportOptions) (*domain.AuditReport, error) {
	// Verify project exists and get settings
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Get audit logs with optional filtering
	auditLogs, err := s.auditRepo.GetAuditLogs(projectID, options.Limit, options.ActionType)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Apply privacy filtering
	anonymize, err := s.shouldAnonymizeForProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to check anonymity settings: %w", err)
	}

	filteredLogs := make([]domain.AuditLog, 0, len(auditLogs))
	for _, log := range auditLogs {
		// Apply privacy controls based on project settings and user permissions
		filteredLog := s.applyPrivacyControls(log, anonymize, options.IncludePersonalData)
		if filteredLog != nil {
			filteredLogs = append(filteredLogs, *filteredLog)
		}
	}

	// Generate report statistics
	stats := s.generateReportStatistics(filteredLogs)

	return &domain.AuditReport{
		ProjectID:    projectID,
		ProjectName:  project.Name,
		GeneratedAt:  time.Now(),
		TotalActions: len(filteredLogs),
		DateRange:    s.calculateDateRange(filteredLogs),
		Statistics:   stats,
		AuditLogs:    filteredLogs,
		PrivacyMode:  anonymize,
		Options:      options,
	}, nil
}

// shouldAnonymizeForProject determines if audit logs should be anonymized for a project
func (s *AuditService) shouldAnonymizeForProject(projectID int) (bool, error) {
	// For now, return false (no anonymization) - in full implementation would check project settings
	// This could be based on project configuration, organizational policies, or legal requirements
	return false, nil
}

// applyPrivacyControls applies privacy filtering to an audit log entry
func (s *AuditService) applyPrivacyControls(log domain.AuditLog, anonymize bool, includePersonalData bool) *domain.AuditLog {
	if !includePersonalData {
		// Remove potentially sensitive data from JSON strings
		if log.NewValue != "" {
			var newValueMap map[string]interface{}
			if json.Unmarshal([]byte(log.NewValue), &newValueMap) == nil {
				// Remove rationale and other potentially sensitive fields
				delete(newValueMap, "rationale")
				if newValueBytes, err := json.Marshal(newValueMap); err == nil {
					log.NewValue = string(newValueBytes)
				}
			}
		}
	}

	if anonymize {
		// Remove attendee identification
		log.AttendeeID = 0
		if log.NewValue != "" {
			var newValueMap map[string]interface{}
			if json.Unmarshal([]byte(log.NewValue), &newValueMap) == nil {
				if actor, exists := newValueMap["actor"]; exists {
					if actorMap, ok := actor.(map[string]interface{}); ok {
						actorMap["type"] = "anonymous_attendee"
						delete(actorMap, "attendee_name")
						delete(actorMap, "attendee_id")
					}
				}
				if newValueBytes, err := json.Marshal(newValueMap); err == nil {
					log.NewValue = string(newValueBytes)
				}
			}
		}
	}

	return &log
}

// generateReportStatistics generates summary statistics for the audit report
func (s *AuditService) generateReportStatistics(logs []domain.AuditLog) domain.AuditStatistics {
	actionCounts := make(map[string]int)
	entityCounts := make(map[string]int)
	attendeeCounts := make(map[int]int)

	for _, log := range logs {
		actionCounts[log.Action]++
		entityCounts[log.EntityType]++
		if log.AttendeeID > 0 {
			attendeeCounts[log.AttendeeID]++
		}
	}

	return domain.AuditStatistics{
		ActionCounts:           actionCounts,
		EntityTypeCounts:       entityCounts,
		AttendeeActivityCounts: attendeeCounts,
		TotalActions:           len(logs),
	}
}

// calculateDateRange calculates the date range covered by the audit logs
func (s *AuditService) calculateDateRange(logs []domain.AuditLog) domain.DateRange {
	if len(logs) == 0 {
		now := time.Now()
		return domain.DateRange{
			StartDate: now,
			EndDate:   now,
		}
	}

	start := logs[0].Timestamp
	end := logs[0].Timestamp

	for _, log := range logs {
		if log.Timestamp.Before(start) {
			start = log.Timestamp
		}
		if log.Timestamp.After(end) {
			end = log.Timestamp
		}
	}

	return domain.DateRange{
		StartDate: start,
		EndDate:   end,
	}
}
