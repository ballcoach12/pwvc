package service

import (
	"fmt"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/websocket"
)

// WebSocketBroadcaster defines the interface for WebSocket broadcasting
type WebSocketBroadcaster interface {
	NotifyVoteSubmitted(sessionID int, voteUpdate websocket.VoteUpdateMessage)
	NotifyConsensusReached(sessionID int, consensus websocket.ConsensusReachedMessage)
	NotifySessionProgress(sessionID int, progress websocket.SessionProgressMessage)
	NotifySessionCompleted(sessionID int, completion websocket.SessionCompletedMessage)
}

// PairwiseService handles business logic for pairwise comparisons
type PairwiseService struct {
	pairwiseRepo  repository.PairwiseRepository
	featureRepo   repository.FeatureRepository
	attendeeRepo  repository.AttendeeRepository
	projectRepo   repository.ProjectRepository
	wsBroadcaster WebSocketBroadcaster
}

// NewPairwiseService creates a new pairwise service
func NewPairwiseService(
	pairwiseRepo repository.PairwiseRepository,
	featureRepo repository.FeatureRepository,
	attendeeRepo repository.AttendeeRepository,
	projectRepo repository.ProjectRepository,
) *PairwiseService {
	return &PairwiseService{
		pairwiseRepo:  pairwiseRepo,
		featureRepo:   featureRepo,
		attendeeRepo:  attendeeRepo,
		projectRepo:   projectRepo,
		wsBroadcaster: nil, // Will be set via SetWebSocketBroadcaster
	}
}

// SetWebSocketBroadcaster sets the WebSocket broadcaster for real-time notifications
func (s *PairwiseService) SetWebSocketBroadcaster(broadcaster WebSocketBroadcaster) {
	s.wsBroadcaster = broadcaster
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
				return fmt.Errorf("failed to create comparison between feature %d and %d: %w", features[i].ID, features[j].ID, err)
			}
		}
	}
	return nil
}

// GetActiveSession retrieves the active pairwise session for a project and criterion
func (s *PairwiseService) GetActiveSession(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, *domain.SessionProgress, error) {
	if projectID <= 0 {
		return nil, nil, domain.NewAPIError(400, "Invalid project ID")
	}

	session, err := s.pairwiseRepo.GetActiveSessionByProjectAndCriterion(projectID, criterionType)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil, domain.NewAPIError(404, "No active session found")
		}
		return nil, nil, domain.NewAPIError(500, "Failed to get active session", err.Error())
	}

	progress, err := s.pairwiseRepo.GetSessionProgress(session.ID)
	if err != nil {
		return session, nil, nil // Return session even if progress calculation fails
	}

	return session, progress, nil
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

	// Send WebSocket notification about the vote update
	if s.wsBroadcaster != nil {
		go s.notifyVoteUpdate(sessionID, req.ComparisonID, *vote)
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

	// Check if consensus was just reached for this comparison
	comparison, err := s.pairwiseRepo.GetComparisonByID(comparisonID)
	if err != nil {
		return err
	}

	// If consensus was reached, send notification
	if comparison.ConsensusReached && s.wsBroadcaster != nil {
		go s.notifyConsensusReached(sessionID, comparisonID, comparison.WinnerID, comparison.IsTie)
	}

	// Check if all comparisons in the session have reached consensus
	progress, err := s.pairwiseRepo.GetSessionProgress(sessionID)
	if err != nil {
		return err
	}

	// Send progress notification
	if s.wsBroadcaster != nil {
		go s.notifySessionProgress(sessionID)
	}

	// If all comparisons are completed, mark session as completed
	if progress.CompletedComparisons == progress.TotalComparisons && progress.TotalComparisons > 0 {
		err = s.pairwiseRepo.CompleteSession(sessionID)
		if err != nil {
			return err
		}

		// Send session completion notification
		if s.wsBroadcaster != nil {
			session, err := s.pairwiseRepo.GetSessionByID(sessionID)
			if err == nil {
				go s.notifySessionCompleted(sessionID, session)
			}
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

	// Send session completion notification
	if s.wsBroadcaster != nil {
		go s.notifySessionCompleted(sessionID, session)
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

// notifyVoteUpdate sends a WebSocket notification about a vote update
func (s *PairwiseService) notifyVoteUpdate(sessionID, comparisonID int, vote domain.AttendeeVote) {
	// Get attendee information
	attendee, err := s.attendeeRepo.GetByID(vote.AttendeeID)
	if err != nil {
		fmt.Printf("Failed to get attendee for vote notification: %v\n", err)
		return
	}

	// Get all votes for this comparison to calculate totals
	votes, err := s.pairwiseRepo.GetVotesByComparisonID(comparisonID)
	if err != nil {
		fmt.Printf("Failed to get votes for notification: %v\n", err)
		return
	}

	// Get comparison details
	comparison, err := s.pairwiseRepo.GetComparisonByID(comparisonID)
	if err != nil {
		fmt.Printf("Failed to get comparison for notification: %v\n", err)
		return
	}

	// Get total attendees for the project
	attendees, err := s.attendeeRepo.GetByProjectID(attendee.ProjectID)
	if err != nil {
		fmt.Printf("Failed to get attendees for notification: %v\n", err)
		return
	}

	voteUpdate := websocket.VoteUpdateMessage{
		ComparisonID:       comparisonID,
		AttendeeID:         vote.AttendeeID,
		AttendeeName:       attendee.Name,
		PreferredFeatureID: vote.PreferredFeatureID,
		IsTieVote:          vote.IsTieVote,
		VotesReceived:      len(votes),
		TotalAttendees:     len(attendees),
		ConsensusReached:   comparison.ConsensusReached,
	}

	s.wsBroadcaster.NotifyVoteSubmitted(sessionID, voteUpdate)
}

// notifyConsensusReached sends a WebSocket notification about consensus
func (s *PairwiseService) notifyConsensusReached(sessionID, comparisonID int, winnerID *int, isTie bool) {
	// Get comparison details
	comparison, err := s.pairwiseRepo.GetComparisonByID(comparisonID)
	if err != nil {
		fmt.Printf("Failed to get comparison for consensus notification: %v\n", err)
		return
	}

	consensusMsg := websocket.ConsensusReachedMessage{
		ComparisonID: comparisonID,
		WinnerID:     winnerID,
		IsTie:        isTie,
	}

	// Add feature names if available
	if comparison.FeatureA != nil {
		consensusMsg.FeatureAName = comparison.FeatureA.Title
	}
	if comparison.FeatureB != nil {
		consensusMsg.FeatureBName = comparison.FeatureB.Title
	}
	if comparison.Winner != nil {
		consensusMsg.WinnerName = comparison.Winner.Title
	}

	s.wsBroadcaster.NotifyConsensusReached(sessionID, consensusMsg)
}

// notifySessionProgress sends a WebSocket notification about session progress
func (s *PairwiseService) notifySessionProgress(sessionID int) {
	progress, err := s.pairwiseRepo.GetSessionProgress(sessionID)
	if err != nil {
		fmt.Printf("Failed to get session progress for notification: %v\n", err)
		return
	}

	progressMsg := websocket.SessionProgressMessage{
		SessionID:            progress.SessionID,
		CompletedComparisons: progress.CompletedComparisons,
		TotalComparisons:     progress.TotalComparisons,
		ProgressPercentage:   progress.ProgressPercentage,
		RemainingComparisons: progress.RemainingComparisons,
	}

	s.wsBroadcaster.NotifySessionProgress(sessionID, progressMsg)
}

// notifySessionCompleted sends a WebSocket notification about session completion
func (s *PairwiseService) notifySessionCompleted(sessionID int, session *domain.PairwiseSession) {
	// Get session statistics
	progress, err := s.pairwiseRepo.GetSessionProgress(sessionID)
	if err != nil {
		fmt.Printf("Failed to get session progress for completion notification: %v\n", err)
		return
	}

	completionMsg := websocket.SessionCompletedMessage{
		SessionID:      sessionID,
		CriterionType:  string(session.CriterionType),
		TotalVotes:     0, // TODO: Calculate total votes
		TotalConsensus: progress.CompletedComparisons,
	}

	s.wsBroadcaster.NotifySessionCompleted(sessionID, completionMsg)
}

// ReassignPendingComparisons allows reassignment of pending comparisons to different attendees or sessions (T042 - US8)
func (s *PairwiseService) ReassignPendingComparisons(projectID int, reassignmentRequest domain.ReassignmentRequest) error {
	// Verify project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Get pending comparisons for the session/criterion
	pendingComparisons, err := s.pairwiseRepo.GetPendingComparisons(reassignmentRequest.SessionID, reassignmentRequest.CriterionType)
	if err != nil {
		return fmt.Errorf("failed to get pending comparisons: %w", err)
	}

	// Process reassignments
	var reassignedCount int
	for _, comparisonID := range reassignmentRequest.ComparisonIDs {
		// Find the comparison in pending list
		var targetComparison *domain.SessionComparison
		for _, comp := range pendingComparisons {
			if comp.ID == comparisonID {
				targetComparison = comp
				break
			}
		}

		if targetComparison == nil {
			continue // Skip if comparison not found or not pending
		}

		// Perform the reassignment based on request type
		switch reassignmentRequest.ReassignmentType {
		case "session":
			// Move comparison to different session
			err = s.pairwiseRepo.MoveComparisonToSession(comparisonID, reassignmentRequest.TargetSessionID)
		case "reset":
			// Reset comparison votes to allow re-voting
			err = s.pairwiseRepo.ResetComparisonVotes(comparisonID)
		case "priority":
			// Change comparison priority/order
			err = s.pairwiseRepo.UpdateComparisonPriority(comparisonID, reassignmentRequest.NewPriority)
		default:
			return fmt.Errorf("invalid reassignment type: %s", reassignmentRequest.ReassignmentType)
		}

		if err != nil {
			return fmt.Errorf("failed to reassign comparison %d: %w", comparisonID, err)
		}

		reassignedCount++
	}

	if reassignedCount == 0 {
		return fmt.Errorf("no comparisons were reassigned")
	}

	// Broadcast reassignment notification
	if s.wsBroadcaster != nil {
		// For now, just notify progress update - in full implementation would have specific reassignment events
		s.notifySessionProgress(reassignmentRequest.SessionID)
		if reassignmentRequest.TargetSessionID != 0 && reassignmentRequest.TargetSessionID != reassignmentRequest.SessionID {
			s.notifySessionProgress(reassignmentRequest.TargetSessionID)
		}
	}

	return nil
}

// GetPendingComparisons retrieves pending comparisons for a session (T042 - US8)
func (s *PairwiseService) GetPendingComparisons(sessionID int, criterionType domain.CriterionType) ([]*domain.SessionComparison, error) {
	return s.pairwiseRepo.GetPendingComparisons(sessionID, string(criterionType))
}

// GetReassignmentOptions provides options for comparison reassignment (T042 - US8)
func (s *PairwiseService) GetReassignmentOptions(projectID int, sessionID int) (*domain.ReassignmentOptions, error) {
	// Get all sessions for the project
	sessions, err := s.pairwiseRepo.GetProjectSessions(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project sessions: %w", err)
	}

	// Get pending comparisons count
	pendingComparisons, err := s.pairwiseRepo.GetPendingComparisons(sessionID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get pending comparisons: %w", err)
	}

	return &domain.ReassignmentOptions{
		ProjectID:            projectID,
		CurrentSessionID:     sessionID,
		AvailableSessions:    sessions,
		PendingComparisons:   pendingComparisons,
		ReassignmentTypes:    []string{"session", "reset", "priority"},
		CanReassignToSession: len(sessions) > 1,
		CanResetVotes:        len(pendingComparisons) > 0,
		CanChangePriority:    len(pendingComparisons) > 1,
	}, nil
}
