package service

import (
	"fmt"

	"pwvc/internal/domain"
	"pwvc/internal/repository"
)

// PairwiseService handles business logic for pairwise comparisons
type PairwiseService struct {
	pairwiseRepo *repository.PairwiseRepository
	featureRepo  *repository.FeatureRepository
	attendeeRepo *repository.AttendeeRepository
	projectRepo  *repository.ProjectRepository
}

// NewPairwiseService creates a new pairwise service
func NewPairwiseService(
	pairwiseRepo *repository.PairwiseRepository,
	featureRepo *repository.FeatureRepository,
	attendeeRepo *repository.AttendeeRepository,
	projectRepo *repository.ProjectRepository,
) *PairwiseService {
	return &PairwiseService{
		pairwiseRepo: pairwiseRepo,
		featureRepo:  featureRepo,
		attendeeRepo: attendeeRepo,
		projectRepo:  projectRepo,
	}
}

// StartPairwiseSession starts a new pairwise comparison session
func (s *PairwiseService) StartPairwiseSession(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Validate project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate project", err.Error())
	}

	// Check if there's already an active session for this criterion
	existingSession, err := s.pairwiseRepo.GetActiveSessionByProjectAndCriterion(projectID, criterionType)
	if err == nil && existingSession != nil {
		return nil, domain.NewAPIError(409, fmt.Sprintf("Active %s session already exists", criterionType))
	}

	// Get all features for the project
	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get project features", err.Error())
	}

	if len(features) < 2 {
		return nil, domain.NewAPIError(400, "At least 2 features are required for pairwise comparison")
	}

	// Create the session
	session, err := s.pairwiseRepo.CreateSession(projectID, criterionType)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to create pairwise session", err.Error())
	}

	// Generate all unique feature pairs and create comparisons
	err = s.generateComparisons(session.ID, features)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to generate comparisons", err.Error())
	}

	return session, nil
}

// generateComparisons creates all unique pairwise comparisons for features
func (s *PairwiseService) generateComparisons(sessionID int, features []domain.Feature) error {
	for i := 0; i < len(features); i++ {
		for j := i + 1; j < len(features); j++ {
			_, err := s.pairwiseRepo.CreateComparison(sessionID, features[i].ID, features[j].ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetSession retrieves a pairwise session with progress information
func (s *PairwiseService) GetSession(sessionID int) (*domain.PairwiseSession, *domain.SessionProgress, error) {
	if sessionID <= 0 {
		return nil, nil, domain.NewAPIError(400, "Invalid session ID")
	}

	session, err := s.pairwiseRepo.GetSessionByID(sessionID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, domain.NewAPIError(404, "Session not found")
		}
		return nil, nil, domain.NewAPIError(500, "Failed to get session", err.Error())
	}

	progress, err := s.pairwiseRepo.GetSessionProgress(sessionID)
	if err != nil {
		return nil, nil, domain.NewAPIError(500, "Failed to get session progress", err.Error())
	}

	return session, progress, nil
}

// GetSessionComparisons retrieves all comparisons for a session
func (s *PairwiseService) GetSessionComparisons(sessionID int) ([]domain.ComparisonWithVotes, error) {
	if sessionID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid session ID")
	}

	// Validate session exists
	_, err := s.pairwiseRepo.GetSessionByID(sessionID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Session not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate session", err.Error())
	}

	comparisons, err := s.pairwiseRepo.GetComparisonsBySessionID(sessionID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get comparisons", err.Error())
	}

	var result []domain.ComparisonWithVotes
	for _, comparison := range comparisons {
		votes, err := s.pairwiseRepo.GetVotesByComparisonID(comparison.ID)
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to get votes", err.Error())
		}

		result = append(result, domain.ComparisonWithVotes{
			Comparison: &comparison,
			Votes:      votes,
		})
	}

	return result, nil
}

// SubmitVote submits or updates an attendee vote for a comparison
func (s *PairwiseService) SubmitVote(sessionID int, req domain.SubmitVoteRequest) (*domain.AttendeeVote, error) {
	if sessionID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid session ID")
	}

	// Validate session exists and is active
	session, err := s.pairwiseRepo.GetSessionByID(sessionID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Session not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate session", err.Error())
	}

	if session.Status != domain.SessionStatusActive {
		return nil, domain.NewAPIError(400, "Session is not active")
	}

	// Validate comparison exists and belongs to the session
	comparison, err := s.pairwiseRepo.GetComparisonByID(req.ComparisonID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Comparison not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate comparison", err.Error())
	}

	if comparison.SessionID != sessionID {
		return nil, domain.NewAPIError(400, "Comparison does not belong to this session")
	}

	// Validate attendee exists and belongs to the project
	attendee, err := s.attendeeRepo.GetByID(req.AttendeeID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Attendee not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate attendee", err.Error())
	}

	if attendee.ProjectID != session.ProjectID {
		return nil, domain.NewAPIError(400, "Attendee does not belong to this project")
	}

	// Validate vote consistency
	if req.IsTieVote && req.PreferredFeatureID != nil {
		return nil, domain.NewAPIError(400, "Cannot specify preferred feature for tie vote")
	}

	if !req.IsTieVote && req.PreferredFeatureID == nil {
		return nil, domain.NewAPIError(400, "Must specify preferred feature for non-tie vote")
	}

	if req.PreferredFeatureID != nil {
		if *req.PreferredFeatureID != comparison.FeatureAID && *req.PreferredFeatureID != comparison.FeatureBID {
			return nil, domain.NewAPIError(400, "Preferred feature must be one of the compared features")
		}
	}

	// Check if attendee has already voted
	existingVote, err := s.pairwiseRepo.GetVoteByAttendeeAndComparison(req.ComparisonID, req.AttendeeID)
	var vote *domain.AttendeeVote

	if err != nil && err != domain.ErrNotFound {
		return nil, domain.NewAPIError(500, "Failed to check existing vote", err.Error())
	}

	if existingVote != nil {
		// Update existing vote
		voteToUpdate := domain.AttendeeVote{
			ComparisonID:       req.ComparisonID,
			AttendeeID:         req.AttendeeID,
			PreferredFeatureID: req.PreferredFeatureID,
			IsTieVote:          req.IsTieVote,
		}

		err = s.pairwiseRepo.UpdateVote(voteToUpdate)
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to update vote", err.Error())
		}

		// Get updated vote
		vote, err = s.pairwiseRepo.GetVoteByAttendeeAndComparison(req.ComparisonID, req.AttendeeID)
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to get updated vote", err.Error())
		}
	} else {
		// Create new vote
		newVote := domain.AttendeeVote{
			ComparisonID:       req.ComparisonID,
			AttendeeID:         req.AttendeeID,
			PreferredFeatureID: req.PreferredFeatureID,
			IsTieVote:          req.IsTieVote,
		}

		vote, err = s.pairwiseRepo.CreateVote(newVote)
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to create vote", err.Error())
		}
	}

	// Check for consensus and auto-complete session if needed
	err = s.checkAndUpdateConsensus(sessionID, req.ComparisonID, session.ProjectID)
	if err != nil {
		// Log error but don't fail the vote submission
		fmt.Printf("Warning: Failed to check consensus: %v\n", err)
	}

	return vote, nil
}

// checkAndUpdateConsensus checks if consensus is reached and completes session if all comparisons are done
func (s *PairwiseService) checkAndUpdateConsensus(sessionID, comparisonID, projectID int) error {
	// Get total number of attendees for the project
	attendees, err := s.attendeeRepo.GetByProjectID(projectID)
	if err != nil {
		return err
	}

	// Check consensus for this specific comparison
	err = s.pairwiseRepo.CheckConsensusAndUpdate(comparisonID, len(attendees))
	if err != nil {
		return err
	}

	// Check if all comparisons in the session have reached consensus
	progress, err := s.pairwiseRepo.GetSessionProgress(sessionID)
	if err != nil {
		return err
	}

	// If all comparisons are completed, mark session as completed
	if progress.CompletedComparisons == progress.TotalComparisons && progress.TotalComparisons > 0 {
		err = s.pairwiseRepo.CompleteSession(sessionID)
		if err != nil {
			return err
		}
	}

	return nil
}

// CompleteSession manually completes a pairwise session
func (s *PairwiseService) CompleteSession(sessionID int) error {
	if sessionID <= 0 {
		return domain.NewAPIError(400, "Invalid session ID")
	}

	// Validate session exists and is active
	session, err := s.pairwiseRepo.GetSessionByID(sessionID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Session not found")
		}
		return domain.NewAPIError(500, "Failed to validate session", err.Error())
	}

	if session.Status != domain.SessionStatusActive {
		return domain.NewAPIError(400, "Session is not active")
	}

	err = s.pairwiseRepo.CompleteSession(sessionID)
	if err != nil {
		return domain.NewAPIError(500, "Failed to complete session", err.Error())
	}

	return nil
}

// GetNextComparison finds the next comparison that needs votes from a specific attendee
func (s *PairwiseService) GetNextComparison(sessionID, attendeeID int) (*domain.ComparisonWithVotes, error) {
	if sessionID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid session ID")
	}

	if attendeeID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid attendee ID")
	}

	// Validate session exists and is active
	session, err := s.pairwiseRepo.GetSessionByID(sessionID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Session not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate session", err.Error())
	}

	if session.Status != domain.SessionStatusActive {
		return nil, domain.NewAPIError(400, "Session is not active")
	}

	// Get all comparisons for the session
	comparisons, err := s.pairwiseRepo.GetComparisonsBySessionID(sessionID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get comparisons", err.Error())
	}

	// Find the first comparison where the attendee hasn't voted
	for _, comparison := range comparisons {
		_, err := s.pairwiseRepo.GetVoteByAttendeeAndComparison(comparison.ID, attendeeID)
		if err == domain.ErrNotFound {
			// Attendee hasn't voted on this comparison
			votes, err := s.pairwiseRepo.GetVotesByComparisonID(comparison.ID)
			if err != nil {
				return nil, domain.NewAPIError(500, "Failed to get votes", err.Error())
			}

			return &domain.ComparisonWithVotes{
				Comparison: &comparison,
				Votes:      votes,
			}, nil
		} else if err != nil {
			return nil, domain.NewAPIError(500, "Failed to check existing vote", err.Error())
		}
		// If we reach here, attendee has already voted on this comparison
	}

	// No comparisons found that need this attendee's vote
	return nil, domain.NewAPIError(404, "No pending comparisons found")
}
