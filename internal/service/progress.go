package service

import (
	"fmt"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

type ProgressService struct {
	progressRepo repository.ProgressRepository
	projectRepo  repository.ProjectRepository
	attendeeRepo repository.AttendeeRepository
	featureRepo  repository.FeatureRepository
	scoringRepo  repository.ScoringRepository
}

func NewProgressService(progressRepo repository.ProgressRepository, projectRepo repository.ProjectRepository, attendeeRepo repository.AttendeeRepository, featureRepo repository.FeatureRepository, scoringRepo repository.ScoringRepository) *ProgressService {
	return &ProgressService{
		progressRepo: progressRepo,
		projectRepo:  projectRepo,
		attendeeRepo: attendeeRepo,
		featureRepo:  featureRepo,
		scoringRepo:  scoringRepo,
	}
}

// GetProjectProgress retrieves the current progress for a project
func (s *ProgressService) GetProjectProgress(projectID int) (*domain.ProjectProgress, error) {
	// Ensure project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	progress, err := s.progressRepo.GetProjectProgress(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project progress: %w", err)
	}

	// Auto-update progress based on current data
	err = s.updateProgressBasedOnData(progress)
	if err != nil {
		return nil, fmt.Errorf("failed to update progress: %w", err)
	}

	return progress, nil
}

// AdvanceToPhase attempts to advance the project to a specific phase
func (s *ProgressService) AdvanceToPhase(projectID int, phase domain.WorkflowPhase) error {
	progress, err := s.progressRepo.GetProjectProgress(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project progress: %w", err)
	}

	if !progress.CanProgressTo(phase) {
		return fmt.Errorf("cannot advance to phase %s: prerequisites not met", phase)
	}

	// Update current phase
	progress.CurrentPhase = string(phase)
	return s.progressRepo.UpdateProjectProgress(progress)
}

// CompletePhase marks a phase as completed and advances to the next phase
func (s *ProgressService) CompletePhase(projectID int, phase domain.WorkflowPhase) error {
	// Validate prerequisites are met
	isValid, err := s.validatePhaseCompletion(projectID, phase)
	if err != nil {
		return fmt.Errorf("failed to validate phase completion: %w", err)
	}

	if !isValid {
		return fmt.Errorf("phase %s cannot be completed: requirements not met", phase)
	}

	return s.progressRepo.MarkPhaseCompleted(projectID, phase)
}

// validatePhaseCompletion checks if a phase can actually be completed based on data
func (s *ProgressService) validatePhaseCompletion(projectID int, phase domain.WorkflowPhase) (bool, error) {
	switch phase {
	case domain.PhaseSetup:
		// Project must exist (already validated in calling functions)
		return true, nil

	case domain.PhaseAttendees:
		attendees, err := s.attendeeRepo.GetByProjectID(projectID)
		if err != nil {
			return false, err
		}
		// Require at least 2 attendees for pairwise comparisons
		return len(attendees) >= 2, nil

	case domain.PhaseFeatures:
		features, err := s.featureRepo.GetByProjectID(projectID)
		if err != nil {
			return false, err
		}
		// Require at least 2 features for comparisons
		return len(features) >= 2, nil

	case domain.PhasePairwiseValue:
		// Check if all value comparisons are complete
		// This would require checking the pairwise comparison data
		// For now, assume it's valid if we reach this point
		return true, nil

	case domain.PhasePairwiseComplexity:
		// Check if all complexity comparisons are complete
		return true, nil

	case domain.PhaseFibonacciValue:
		// Check if all value scores are assigned
		return true, nil

	case domain.PhaseFibonacciComplexity:
		// Check if all complexity scores are assigned
		return true, nil

	case domain.PhaseResults:
		// All previous phases must be completed
		progressData, err := s.progressRepo.GetProjectProgress(projectID)
		if err != nil {
			return false, err
		}
		return progressData.FibonacciComplexityCompleted, nil

	default:
		return false, nil
	}
}

// updateProgressBasedOnData automatically updates progress flags based on actual data
func (s *ProgressService) updateProgressBasedOnData(progress *domain.ProjectProgress) error {
	var updated bool

	// Check attendees
	attendees, err := s.attendeeRepo.GetByProjectID(progress.ProjectID)
	if err == nil && len(attendees) >= 2 && !progress.AttendeesAdded {
		progress.AttendeesAdded = true
		updated = true
	}

	// Check features
	features, err := s.featureRepo.GetByProjectID(progress.ProjectID)
	if err == nil && len(features) >= 2 && !progress.FeaturesAdded {
		progress.FeaturesAdded = true
		updated = true
	}

	// Setup is considered complete if both attendees and features are added
	if progress.AttendeesAdded && progress.FeaturesAdded && !progress.SetupCompleted {
		progress.SetupCompleted = true
		updated = true
	}

	// Save updates if any changes were made
	if updated {
		return s.progressRepo.UpdateProjectProgress(progress)
	}

	return nil
}

// GetAvailablePhases returns the phases that can be accessed based on current progress
func (s *ProgressService) GetAvailablePhases(projectID int) ([]domain.WorkflowPhase, error) {
	progress, err := s.progressRepo.GetProjectProgress(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project progress: %w", err)
	}

	var availablePhases []domain.WorkflowPhase

	phases := []domain.WorkflowPhase{
		domain.PhaseSetup,
		domain.PhaseAttendees,
		domain.PhaseFeatures,
		domain.PhasePairwiseValue,
		domain.PhasePairwiseComplexity,
		domain.PhaseFibonacciValue,
		domain.PhaseFibonacciComplexity,
		domain.PhaseResults,
	}

	for _, phase := range phases {
		if progress.CanProgressTo(phase) {
			availablePhases = append(availablePhases, phase)
		}
	}

	return availablePhases, nil
}

// GetFibonacciProgressMetrics retrieves progress metrics for Fibonacci scoring phases (T040 - US8)
func (s *ProgressService) GetFibonacciProgressMetrics(projectID int, criterionType string) (*domain.FibonacciProgressMetrics, error) {
	// Get all features for the project
	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project features: %w", err)
	}

	// Get all attendees for the project
	attendees, err := s.attendeeRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project attendees: %w", err)
	}

	// Get all scores for the project and criterion type
	scores, err := s.scoringRepo.GetByProject(projectID, criterionType)
	if err != nil {
		return nil, fmt.Errorf("failed to get project scores: %w", err)
	}

	// Calculate metrics
	totalExpectedScores := len(features) * len(attendees)
	completedScores := len(scores)
	progressPercentage := 0.0
	if totalExpectedScores > 0 {
		progressPercentage = float64(completedScores) / float64(totalExpectedScores) * 100.0
	}

	// Calculate per-feature completion
	featureCompletion := make(map[int]domain.FeatureScoreProgress)
	for _, feature := range features {
		featureScores := 0
		for _, score := range scores {
			if score.FeatureID == feature.ID {
				featureScores++
			}
		}
		featureCompletion[feature.ID] = domain.FeatureScoreProgress{
			FeatureID:           feature.ID,
			FeatureName:         feature.Title,
			CompletedScores:     featureScores,
			TotalExpectedScores: len(attendees),
			ProgressPercentage:  float64(featureScores) / float64(len(attendees)) * 100.0,
		}
	}

	// Calculate per-attendee completion
	attendeeCompletion := make(map[int]domain.AttendeeScoreProgress)
	for _, attendee := range attendees {
		attendeeScores := 0
		for _, score := range scores {
			if score.AttendeeID == attendee.ID {
				attendeeScores++
			}
		}
		attendeeCompletion[attendee.ID] = domain.AttendeeScoreProgress{
			AttendeeID:          attendee.ID,
			AttendeeName:        attendee.Name,
			CompletedScores:     attendeeScores,
			TotalExpectedScores: len(features),
			ProgressPercentage:  float64(attendeeScores) / float64(len(features)) * 100.0,
		}
	}

	return &domain.FibonacciProgressMetrics{
		ProjectID:           projectID,
		CriterionType:       criterionType,
		CompletedScores:     completedScores,
		TotalExpectedScores: totalExpectedScores,
		ProgressPercentage:  progressPercentage,
		FeatureCompletion:   featureCompletion,
		AttendeeCompletion:  attendeeCompletion,
	}, nil
}

// GetOverallFibonacciProgress gets combined progress for both value and complexity scoring (T040 - US8)
func (s *ProgressService) GetOverallFibonacciProgress(projectID int) (*domain.OverallFibonacciProgress, error) {
	valueMetrics, err := s.GetFibonacciProgressMetrics(projectID, "value")
	if err != nil {
		return nil, fmt.Errorf("failed to get value metrics: %w", err)
	}

	complexityMetrics, err := s.GetFibonacciProgressMetrics(projectID, "complexity")
	if err != nil {
		return nil, fmt.Errorf("failed to get complexity metrics: %w", err)
	}

	// Calculate overall progress
	totalCompleted := valueMetrics.CompletedScores + complexityMetrics.CompletedScores
	totalExpected := valueMetrics.TotalExpectedScores + complexityMetrics.TotalExpectedScores
	overallProgress := 0.0
	if totalExpected > 0 {
		overallProgress = float64(totalCompleted) / float64(totalExpected) * 100.0
	}

	return &domain.OverallFibonacciProgress{
		ProjectID:            projectID,
		ValueMetrics:         valueMetrics,
		ComplexityMetrics:    complexityMetrics,
		OverallProgress:      overallProgress,
		TotalCompleted:       totalCompleted,
		TotalExpected:        totalExpected,
		IsValueComplete:      valueMetrics.ProgressPercentage >= 100.0,
		IsComplexityComplete: complexityMetrics.ProgressPercentage >= 100.0,
		IsBothComplete:       valueMetrics.ProgressPercentage >= 100.0 && complexityMetrics.ProgressPercentage >= 100.0,
	}, nil
}
