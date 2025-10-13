package repository

import (
	"database/sql"

	"pairwise/internal/domain"
)

// PairwiseRepository handles database operations for pairwise comparisons
type PairwiseRepository struct {
	db *sql.DB
}

// NewPairwiseRepository creates a new pairwise repository
func NewPairwiseRepository(db *sql.DB) *PairwiseRepository {
	return &PairwiseRepository{db: db}
}

// CreateSession creates a new pairwise comparison session
func (r *PairwiseRepository) CreateSession(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error) {
	query := `
		INSERT INTO pairwise_sessions (project_id, criterion_type, status, started_at)
		VALUES (?, ?, ?, datetime('now'))
		RETURNING id, project_id, criterion_type, status, started_at, completed_at
	`

	var session domain.PairwiseSession
	err := r.db.QueryRow(query, projectID, criterionType, domain.SessionStatusActive).Scan(
		&session.ID,
		&session.ProjectID,
		&session.CriterionType,
		&session.Status,
		&session.StartedAt,
		&session.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionByID retrieves a pairwise session by ID
func (r *PairwiseRepository) GetSessionByID(sessionID int) (*domain.PairwiseSession, error) {
	query := `
		SELECT id, project_id, criterion_type, status, started_at, completed_at
		FROM pairwise_sessions
		WHERE id = ?
	`

	var session domain.PairwiseSession
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.ProjectID,
		&session.CriterionType,
		&session.Status,
		&session.StartedAt,
		&session.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetActiveSessionByProjectAndCriterion gets active session for project and criterion
func (r *PairwiseRepository) GetActiveSessionByProjectAndCriterion(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error) {
	query := `
		SELECT id, project_id, criterion_type, status, started_at, completed_at
		FROM pairwise_sessions
		WHERE project_id = ? AND criterion_type = ? AND status = ?
		ORDER BY started_at DESC
		LIMIT 1
	`

	var session domain.PairwiseSession
	err := r.db.QueryRow(query, projectID, criterionType, domain.SessionStatusActive).Scan(
		&session.ID,
		&session.ProjectID,
		&session.CriterionType,
		&session.Status,
		&session.StartedAt,
		&session.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// CompleteSession marks a session as completed
func (r *PairwiseRepository) CompleteSession(sessionID int) error {
	query := `
		UPDATE pairwise_sessions
		SET status = ?, completed_at = datetime('now')
		WHERE id = ?
	`

	_, err := r.db.Exec(query, domain.SessionStatusCompleted, sessionID)
	return err
}

// CreateComparison creates a new comparison between two features
func (r *PairwiseRepository) CreateComparison(sessionID, featureAID, featureBID int) (*domain.SessionComparison, error) {
	query := `
		INSERT INTO pairwise_comparisons (session_id, feature_a_id, feature_b_id, created_at)
		VALUES (?, ?, ?, datetime('now'))
		RETURNING id, session_id, feature_a_id, feature_b_id, winner_id, is_tie, consensus_reached, created_at
	`

	var comparison domain.SessionComparison
	err := r.db.QueryRow(query, sessionID, featureAID, featureBID).Scan(
		&comparison.ID,
		&comparison.SessionID,
		&comparison.FeatureAID,
		&comparison.FeatureBID,
		&comparison.WinnerID,
		&comparison.IsTie,
		&comparison.ConsensusReached,
		&comparison.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comparison, nil
}

// GetComparisonsBySessionID retrieves all comparisons for a session
func (r *PairwiseRepository) GetComparisonsBySessionID(sessionID int) ([]domain.SessionComparison, error) {
	query := `
		SELECT pc.id, pc.session_id, pc.feature_a_id, pc.feature_b_id, pc.winner_id, 
		       pc.is_tie, pc.consensus_reached, pc.created_at,
		       fa.id, fa.title, fa.description,
		       fb.id, fb.title, fb.description
		FROM pairwise_comparisons pc
		JOIN features fa ON pc.feature_a_id = fa.id
		JOIN features fb ON pc.feature_b_id = fb.id
		WHERE pc.session_id = ?
		ORDER BY pc.created_at ASC
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comparisons []domain.SessionComparison
	for rows.Next() {
		var comparison domain.SessionComparison
		var featureA, featureB domain.Feature

		err := rows.Scan(
			&comparison.ID,
			&comparison.SessionID,
			&comparison.FeatureAID,
			&comparison.FeatureBID,
			&comparison.WinnerID,
			&comparison.IsTie,
			&comparison.ConsensusReached,
			&comparison.CreatedAt,
			&featureA.ID,
			&featureA.Title,
			&featureA.Description,
			&featureB.ID,
			&featureB.Title,
			&featureB.Description,
		)
		if err != nil {
			return nil, err
		}

		comparison.FeatureA = &featureA
		comparison.FeatureB = &featureB
		comparisons = append(comparisons, comparison)
	}

	return comparisons, nil
}

// GetComparisonByID retrieves a comparison by ID with feature details
func (r *PairwiseRepository) GetComparisonByID(comparisonID int) (*domain.SessionComparison, error) {
	query := `
		SELECT pc.id, pc.session_id, pc.feature_a_id, pc.feature_b_id, pc.winner_id, 
		       pc.is_tie, pc.consensus_reached, pc.created_at,
		       fa.id, fa.title, fa.description,
		       fb.id, fb.title, fb.description
		FROM pairwise_comparisons pc
		JOIN features fa ON pc.feature_a_id = fa.id
		JOIN features fb ON pc.feature_b_id = fb.id
		WHERE pc.id = ?
	`

	var comparison domain.SessionComparison
	var featureA, featureB domain.Feature

	err := r.db.QueryRow(query, comparisonID).Scan(
		&comparison.ID,
		&comparison.SessionID,
		&comparison.FeatureAID,
		&comparison.FeatureBID,
		&comparison.WinnerID,
		&comparison.IsTie,
		&comparison.ConsensusReached,
		&comparison.CreatedAt,
		&featureA.ID,
		&featureA.Title,
		&featureA.Description,
		&featureB.ID,
		&featureB.Title,
		&featureB.Description,
	)
	if err != nil {
		return nil, err
	}

	comparison.FeatureA = &featureA
	comparison.FeatureB = &featureB

	return &comparison, nil
}

// CreateVote creates a new attendee vote for a comparison
func (r *PairwiseRepository) CreateVote(vote domain.AttendeeVote) (*domain.AttendeeVote, error) {
	query := `
		INSERT INTO attendee_votes (comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at
	`

	var newVote domain.AttendeeVote
	err := r.db.QueryRow(query, vote.ComparisonID, vote.AttendeeID, vote.PreferredFeatureID, vote.IsTieVote).Scan(
		&newVote.ID,
		&newVote.ComparisonID,
		&newVote.AttendeeID,
		&newVote.PreferredFeatureID,
		&newVote.IsTieVote,
		&newVote.VotedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newVote, nil
}

// UpdateVote updates an existing attendee vote
func (r *PairwiseRepository) UpdateVote(vote domain.AttendeeVote) error {
	query := `
		UPDATE attendee_votes
		SET preferred_feature_id = ?, is_tie_vote = ?, voted_at = datetime('now')
		WHERE comparison_id = ? AND attendee_id = ?
	`

	_, err := r.db.Exec(query, vote.PreferredFeatureID, vote.IsTieVote, vote.ComparisonID, vote.AttendeeID)
	return err
}

// GetVotesByComparisonID retrieves all votes for a comparison
func (r *PairwiseRepository) GetVotesByComparisonID(comparisonID int) ([]domain.AttendeeVote, error) {
	query := `
		SELECT av.id, av.comparison_id, av.attendee_id, av.preferred_feature_id, 
		       av.is_tie_vote, av.voted_at,
		       a.id, a.name, a.role
		FROM attendee_votes av
		JOIN attendees a ON av.attendee_id = a.id
		WHERE av.comparison_id = ?
		ORDER BY av.voted_at ASC
	`

	rows, err := r.db.Query(query, comparisonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []domain.AttendeeVote
	for rows.Next() {
		var vote domain.AttendeeVote
		var attendee domain.Attendee

		err := rows.Scan(
			&vote.ID,
			&vote.ComparisonID,
			&vote.AttendeeID,
			&vote.PreferredFeatureID,
			&vote.IsTieVote,
			&vote.VotedAt,
			&attendee.ID,
			&attendee.Name,
			&attendee.Role,
		)
		if err != nil {
			return nil, err
		}

		vote.Attendee = &attendee
		votes = append(votes, vote)
	}

	return votes, nil
}

// CheckConsensusAndUpdate checks if consensus is reached and updates the comparison
func (r *PairwiseRepository) CheckConsensusAndUpdate(comparisonID int, totalAttendees int) error {
	// Get all votes for this comparison
	votes, err := r.GetVotesByComparisonID(comparisonID)
	if err != nil {
		return err
	}

	// Check if all attendees have voted
	if len(votes) != totalAttendees {
		return nil // Not all attendees have voted yet
	}

	// Check for consensus
	var winnerID *int
	var isTie bool
	consensusReached := true

	if len(votes) > 0 {
		firstVote := votes[0]
		winnerID = firstVote.PreferredFeatureID
		isTie = firstVote.IsTieVote

		// Check if all votes are the same
		for _, vote := range votes[1:] {
			if vote.IsTieVote != isTie {
				consensusReached = false
				break
			}
			if !vote.IsTieVote && vote.PreferredFeatureID != nil && winnerID != nil {
				if *vote.PreferredFeatureID != *winnerID {
					consensusReached = false
					break
				}
			}
		}
	}

	// Update comparison with consensus result
	if consensusReached {
		query := `
			UPDATE pairwise_comparisons
			SET winner_id = ?, is_tie = ?, consensus_reached = ?
			WHERE id = ?
		`
		_, err = r.db.Exec(query, winnerID, isTie, consensusReached, comparisonID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetSessionProgress calculates the progress of a pairwise session
func (r *PairwiseRepository) GetSessionProgress(sessionID int) (*domain.SessionProgress, error) {
	query := `
		SELECT 
			COUNT(*) as total_comparisons,
			COUNT(CASE WHEN consensus_reached = true THEN 1 END) as completed_comparisons
		FROM pairwise_comparisons
		WHERE session_id = ?
	`

	var progress domain.SessionProgress
	err := r.db.QueryRow(query, sessionID).Scan(
		&progress.TotalComparisons,
		&progress.CompletedComparisons,
	)
	if err != nil {
		return nil, err
	}

	progress.SessionID = sessionID
	progress.RemainingComparisons = progress.TotalComparisons - progress.CompletedComparisons

	if progress.TotalComparisons > 0 {
		progress.ProgressPercentage = float64(progress.CompletedComparisons) / float64(progress.TotalComparisons) * 100
	}

	return &progress, nil
}

// GetVoteByAttendeeAndComparison checks if an attendee has already voted on a comparison
func (r *PairwiseRepository) GetVoteByAttendeeAndComparison(comparisonID, attendeeID int) (*domain.AttendeeVote, error) {
	query := `
		SELECT id, comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at
		FROM attendee_votes
		WHERE comparison_id = ? AND attendee_id = ?
	`

	var vote domain.AttendeeVote
	err := r.db.QueryRow(query, comparisonID, attendeeID).Scan(
		&vote.ID,
		&vote.ComparisonID,
		&vote.AttendeeID,
		&vote.PreferredFeatureID,
		&vote.IsTieVote,
		&vote.VotedAt,
	)

	if err != nil {
		return nil, err
	}

	return &vote, nil
}
