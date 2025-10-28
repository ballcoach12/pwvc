package repository

import (
	"database/sql"
	"fmt"

	"pairwise/internal/domain"
)

// PairwiseRepository interface defines the contract for pairwise comparison data operations
type PairwiseRepository interface {
	CreateSession(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error)
	GetSessionByID(sessionID int) (*domain.PairwiseSession, error)
	GetActiveSessionByProjectAndCriterion(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error)
	CompleteSession(sessionID int) error
	CreateComparison(sessionID, featureAID, featureBID int) (*domain.SessionComparison, error)
	GetComparisonsBySessionID(sessionID int) ([]domain.SessionComparison, error)
	GetComparisonByID(comparisonID int) (*domain.SessionComparison, error)
	CreateVote(vote domain.AttendeeVote) (*domain.AttendeeVote, error)
	UpdateVote(vote domain.AttendeeVote) error
	GetVotesByComparisonID(comparisonID int) ([]domain.AttendeeVote, error)
	CheckConsensusAndUpdate(comparisonID int, totalAttendees int) error
	GetSessionProgress(sessionID int) (*domain.SessionProgress, error)
	GetVoteByAttendeeAndComparison(comparisonID, attendeeID int) (*domain.AttendeeVote, error)
	GetWinCountWeights(projectID int, criterionType domain.CriterionType) (map[int]float64, error)
	GetPendingComparisons(sessionID int, criterionType string) ([]*domain.SessionComparison, error)
	MoveComparisonToSession(comparisonID, targetSessionID int) error
	ResetComparisonVotes(comparisonID int) error
	UpdateComparisonPriority(comparisonID, newPriority int) error
	GetProjectSessions(projectID int) ([]*domain.PairwiseSession, error)
}

// pairwiseRepository handles database operations for pairwise comparisons
type pairwiseRepository struct {
	db *sql.DB
}

// NewPairwiseRepository creates a new pairwise repository
func NewPairwiseRepository(db *sql.DB) PairwiseRepository {
	return &pairwiseRepository{db: db}
}

// CreateSession creates a new pairwise comparison session
func (r *pairwiseRepository) CreateSession(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error) {
	// First insert the session
	insertQuery := `
		INSERT INTO pairwise_sessions (project_id, criterion_type, status, started_at)
		VALUES (?, ?, ?, datetime('now'))
	`

	result, err := r.db.Exec(insertQuery, projectID, criterionType, domain.SessionStatusActive)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID
	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the created session
	selectQuery := `
		SELECT id, project_id, criterion_type, status, started_at, completed_at
		FROM pairwise_sessions
		WHERE id = ?
	`

	var session domain.PairwiseSession
	err = r.db.QueryRow(selectQuery, int(sessionID)).Scan(
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
func (r *pairwiseRepository) GetSessionByID(sessionID int) (*domain.PairwiseSession, error) {
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
func (r *pairwiseRepository) GetActiveSessionByProjectAndCriterion(projectID int, criterionType domain.CriterionType) (*domain.PairwiseSession, error) {
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
func (r *pairwiseRepository) CompleteSession(sessionID int) error {
	query := `
		UPDATE pairwise_sessions
		SET status = ?, completed_at = datetime('now')
		WHERE id = ?
	`

	_, err := r.db.Exec(query, domain.SessionStatusCompleted, sessionID)
	return err
}

// CreateComparison creates a new comparison between two features
func (r *pairwiseRepository) CreateComparison(sessionID, featureAID, featureBID int) (*domain.SessionComparison, error) {
	// Insert the comparison
	insertQuery := `
		INSERT INTO pairwise_comparisons (session_id, feature_a_id, feature_b_id, created_at)
		VALUES (?, ?, ?, datetime('now'))
	`

	result, err := r.db.Exec(insertQuery, sessionID, featureAID, featureBID)
	if err != nil {
		// Add debug logging
		fmt.Printf("DEBUG: CreateComparison error: %v, sessionID: %d, featureAID: %d, featureBID: %d\n", err, sessionID, featureAID, featureBID)
		return nil, err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the complete record
	selectQuery := `
		SELECT id, session_id, feature_a_id, feature_b_id, winner_id, is_tie, consensus_reached, created_at
		FROM pairwise_comparisons
		WHERE id = ?
	`

	var comparison domain.SessionComparison
	var winnerID sql.NullInt64
	var isTie, consensusReached sql.NullBool
	err = r.db.QueryRow(selectQuery, int(id)).Scan(
		&comparison.ID,
		&comparison.SessionID,
		&comparison.FeatureAID,
		&comparison.FeatureBID,
		&winnerID,
		&isTie,
		&consensusReached,
		&comparison.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if winnerID.Valid {
		winnerVal := int(winnerID.Int64)
		comparison.WinnerID = &winnerVal
	}
	comparison.IsTie = isTie.Valid && isTie.Bool
	comparison.ConsensusReached = consensusReached.Valid && consensusReached.Bool

	if err != nil {
		return nil, err
	}

	return &comparison, nil
}

// GetComparisonsBySessionID retrieves all comparisons for a session
func (r *pairwiseRepository) GetComparisonsBySessionID(sessionID int) ([]domain.SessionComparison, error) {
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
		fmt.Printf("DEBUG: GetComparisonsBySessionID query error: %v, sessionID: %d, query: %s\n", err, sessionID, query)
		return nil, err
	}
	defer rows.Close()

	var comparisons []domain.SessionComparison
	for rows.Next() {
		var comparison domain.SessionComparison
		var featureA, featureB domain.Feature
		var winnerID sql.NullInt64
		var isTie, consensusReached sql.NullBool

		err := rows.Scan(
			&comparison.ID,
			&comparison.SessionID,
			&comparison.FeatureAID,
			&comparison.FeatureBID,
			&winnerID,
			&isTie,
			&consensusReached,
			&comparison.CreatedAt,
			&featureA.ID,
			&featureA.Title,
			&featureA.Description,
			&featureB.ID,
			&featureB.Title,
			&featureB.Description,
		)
		if err != nil {
			fmt.Printf("DEBUG: GetComparisonsBySessionID scan error: %v, sessionID: %d\n", err, sessionID)
			return nil, err
		}

		// Handle nullable fields
		if winnerID.Valid {
			winnerVal := int(winnerID.Int64)
			comparison.WinnerID = &winnerVal
		}
		comparison.IsTie = isTie.Valid && isTie.Bool
		comparison.ConsensusReached = consensusReached.Valid && consensusReached.Bool

		comparison.FeatureA = &featureA
		comparison.FeatureB = &featureB
		comparisons = append(comparisons, comparison)
	}

	return comparisons, nil
}

// GetComparisonByID retrieves a comparison by ID with feature details
func (r *pairwiseRepository) GetComparisonByID(comparisonID int) (*domain.SessionComparison, error) {
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
	var winnerID sql.NullInt64
	var isTie, consensusReached sql.NullBool

	err := r.db.QueryRow(query, comparisonID).Scan(
		&comparison.ID,
		&comparison.SessionID,
		&comparison.FeatureAID,
		&comparison.FeatureBID,
		&winnerID,
		&isTie,
		&consensusReached,
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

	// Handle nullable fields
	if winnerID.Valid {
		winnerVal := int(winnerID.Int64)
		comparison.WinnerID = &winnerVal
	}
	comparison.IsTie = isTie.Valid && isTie.Bool
	comparison.ConsensusReached = consensusReached.Valid && consensusReached.Bool

	comparison.FeatureA = &featureA
	comparison.FeatureB = &featureB

	return &comparison, nil
}

// CreateVote creates a new attendee vote for a comparison
func (r *pairwiseRepository) CreateVote(vote domain.AttendeeVote) (*domain.AttendeeVote, error) {
	// First insert the vote
	insertQuery := `
		INSERT INTO attendee_votes (comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at)
		VALUES (?, ?, ?, ?, datetime('now'))
	`

	result, err := r.db.Exec(insertQuery, vote.ComparisonID, vote.AttendeeID, vote.PreferredFeatureID, vote.IsTieVote)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID
	voteID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the created vote
	selectQuery := `
		SELECT id, comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at
		FROM attendee_votes
		WHERE id = ?
	`

	var newVote domain.AttendeeVote
	err = r.db.QueryRow(selectQuery, int(voteID)).Scan(
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
func (r *pairwiseRepository) UpdateVote(vote domain.AttendeeVote) error {
	query := `
		UPDATE attendee_votes
		SET preferred_feature_id = ?, is_tie_vote = ?, voted_at = datetime('now')
		WHERE comparison_id = ? AND attendee_id = ?
	`

	_, err := r.db.Exec(query, vote.PreferredFeatureID, vote.IsTieVote, vote.ComparisonID, vote.AttendeeID)
	return err
}

// GetVotesByComparisonID retrieves all votes for a comparison
func (r *pairwiseRepository) GetVotesByComparisonID(comparisonID int) ([]domain.AttendeeVote, error) {
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
func (r *pairwiseRepository) CheckConsensusAndUpdate(comparisonID int, totalAttendees int) error {
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
func (r *pairwiseRepository) GetSessionProgress(sessionID int) (*domain.SessionProgress, error) {
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
func (r *pairwiseRepository) GetVoteByAttendeeAndComparison(comparisonID, attendeeID int) (*domain.AttendeeVote, error) {
	query := `
		SELECT id, comparison_id, attendee_id, preferred_feature_id, is_tie_vote, voted_at
		FROM attendee_votes
		WHERE comparison_id = ? AND attendee_id = ?
	`

	var vote domain.AttendeeVote
	var preferredFeatureID sql.NullInt64
	err := r.db.QueryRow(query, comparisonID, attendeeID).Scan(
		&vote.ID,
		&vote.ComparisonID,
		&vote.AttendeeID,
		&preferredFeatureID,
		&vote.IsTieVote,
		&vote.VotedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Convert sql.NullInt64 to *int
	if preferredFeatureID.Valid {
		featureID := int(preferredFeatureID.Int64)
		vote.PreferredFeatureID = &featureID
	}

	return &vote, nil
}

// GetWinCountWeights calculates win-count weights for features by criterion type
func (r *pairwiseRepository) GetWinCountWeights(projectID int, criterionType domain.CriterionType) (map[int]float64, error) {
	// Get all sessions for this project and criterion
	sessionQuery := `
		SELECT id FROM pairwise_sessions 
		WHERE project_id = ? AND criterion_type = ? AND status = ?
	`

	rows, err := r.db.Query(sessionQuery, projectID, criterionType, domain.SessionStatusCompleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessionIDs []int
	for rows.Next() {
		var sessionID int
		if err := rows.Scan(&sessionID); err != nil {
			return nil, err
		}
		sessionIDs = append(sessionIDs, sessionID)
	}

	if len(sessionIDs) == 0 {
		return make(map[int]float64), nil
	}

	// Build dynamic query for all sessions
	sessionIDsStr := ""
	for i, id := range sessionIDs {
		if i > 0 {
			sessionIDsStr += ","
		}
		sessionIDsStr += fmt.Sprintf("%d", id)
	}

	// Calculate wins and total comparisons per feature
	query := fmt.Sprintf(`
		SELECT 
			f.id as feature_id,
			COALESCE(SUM(CASE 
				WHEN av.preferred_feature_id = f.id AND av.is_tie_vote = false THEN 1 
				ELSE 0 
			END), 0) as wins,
			COALESCE(COUNT(CASE 
				WHEN (pc.feature_a_id = f.id OR pc.feature_b_id = f.id) AND av.id IS NOT NULL THEN 1 
			END), 0) as total_comparisons
		FROM features f
		LEFT JOIN pairwise_comparisons pc ON (pc.feature_a_id = f.id OR pc.feature_b_id = f.id)
			AND pc.session_id IN (%s)
		LEFT JOIN attendee_votes av ON av.comparison_id = pc.id
		WHERE f.project_id = ?
		GROUP BY f.id
	`, sessionIDsStr)

	rows, err = r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weights := make(map[int]float64)
	for rows.Next() {
		var featureID int
		var wins, totalComparisons int

		err := rows.Scan(&featureID, &wins, &totalComparisons)
		if err != nil {
			return nil, err
		}

		// Calculate win-count weight (W = wins / total comparisons)
		var weight float64
		if totalComparisons > 0 {
			weight = float64(wins) / float64(totalComparisons)
		} else {
			weight = 0.0
		}

		weights[featureID] = weight
	}

	return weights, nil
}

// GetPendingComparisons retrieves pending comparisons for a session (T042 - US8)
func (r *pairwiseRepository) GetPendingComparisons(sessionID int, criterionType string) ([]*domain.SessionComparison, error) {
	var query string
	var args []interface{}

	if criterionType != "" {
		query = `
			SELECT pc.id, pc.session_id, pc.feature_a_id, pc.feature_b_id, 
				   pc.winner_id, pc.is_tie, pc.consensus_reached, pc.created_at
			FROM pairwise_comparisons pc
			JOIN pairwise_sessions ps ON pc.session_id = ps.id
			WHERE pc.session_id = ? AND ps.criterion_type = ? AND pc.consensus_reached = false
			ORDER BY pc.created_at ASC
		`
		args = []interface{}{sessionID, criterionType}
	} else {
		query = `
			SELECT pc.id, pc.session_id, pc.feature_a_id, pc.feature_b_id, 
				   pc.winner_id, pc.is_tie, pc.consensus_reached, pc.created_at
			FROM pairwise_comparisons pc
			WHERE pc.session_id = ? AND pc.consensus_reached = false
			ORDER BY pc.created_at ASC
		`
		args = []interface{}{sessionID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending comparisons: %w", err)
	}
	defer rows.Close()

	var comparisons []*domain.SessionComparison
	for rows.Next() {
		var comp domain.SessionComparison
		err := rows.Scan(
			&comp.ID,
			&comp.SessionID,
			&comp.FeatureAID,
			&comp.FeatureBID,
			&comp.WinnerID,
			&comp.IsTie,
			&comp.ConsensusReached,
			&comp.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comparison: %w", err)
		}
		comparisons = append(comparisons, &comp)
	}

	return comparisons, nil
}

// MoveComparisonToSession moves a comparison to a different session (T042 - US8)
func (r *pairwiseRepository) MoveComparisonToSession(comparisonID, targetSessionID int) error {
	query := `UPDATE pairwise_comparisons SET session_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, targetSessionID, comparisonID)
	if err != nil {
		return fmt.Errorf("failed to move comparison to session: %w", err)
	}
	return nil
}

// ResetComparisonVotes resets all votes for a comparison (T042 - US8)
func (r *pairwiseRepository) ResetComparisonVotes(comparisonID int) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete all votes for the comparison
	_, err = tx.Exec(`DELETE FROM attendee_votes WHERE comparison_id = ?`, comparisonID)
	if err != nil {
		return fmt.Errorf("failed to delete votes: %w", err)
	}

	// Reset comparison consensus status
	_, err = tx.Exec(`UPDATE pairwise_comparisons SET winner_id = NULL, is_tie = false, consensus_reached = false WHERE id = ?`, comparisonID)
	if err != nil {
		return fmt.Errorf("failed to reset comparison: %w", err)
	}

	return tx.Commit()
}

// UpdateComparisonPriority updates the priority/order of a comparison (T042 - US8)
func (r *pairwiseRepository) UpdateComparisonPriority(comparisonID, newPriority int) error {
	// For now, we'll use created_at to simulate priority - in full implementation would add priority column
	query := `UPDATE pairwise_comparisons SET created_at = datetime('now', '-' || ? || ' minutes') WHERE id = ?`
	_, err := r.db.Exec(query, newPriority, comparisonID)
	if err != nil {
		return fmt.Errorf("failed to update comparison priority: %w", err)
	}
	return nil
}

// GetProjectSessions retrieves all sessions for a project (T042 - US8)
func (r *pairwiseRepository) GetProjectSessions(projectID int) ([]*domain.PairwiseSession, error) {
	query := `
		SELECT id, project_id, criterion_type, status, started_at, completed_at
		FROM pairwise_sessions
		WHERE project_id = ?
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.PairwiseSession
	for rows.Next() {
		var session domain.PairwiseSession
		err := rows.Scan(
			&session.ID,
			&session.ProjectID,
			&session.CriterionType,
			&session.Status,
			&session.StartedAt,
			&session.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}
