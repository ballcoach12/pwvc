package repository

import (
	"database/sql"
	"fmt"

	"pairwise/internal/domain"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// GetProjectProgress retrieves the progress for a project
func (r *ProgressRepository) GetProjectProgress(projectID int) (*domain.ProjectProgress, error) {
	query := `
		SELECT project_id, setup_completed, attendees_added, features_added,
		       pairwise_value_completed, pairwise_complexity_completed,
		       fibonacci_value_completed, fibonacci_complexity_completed,
		       results_calculated, current_phase, last_activity, created_at, updated_at
		FROM project_progress 
		WHERE project_id = ?`

	var progress domain.ProjectProgress
	err := r.db.QueryRow(query, projectID).Scan(
		&progress.ProjectID,
		&progress.SetupCompleted,
		&progress.AttendeesAdded,
		&progress.FeaturesAdded,
		&progress.PairwiseValueCompleted,
		&progress.PairwiseComplexityCompleted,
		&progress.FibonacciValueCompleted,
		&progress.FibonacciComplexityCompleted,
		&progress.ResultsCalculated,
		&progress.CurrentPhase,
		&progress.LastActivity,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create initial progress record if it doesn't exist
		return r.CreateProjectProgress(projectID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get project progress: %w", err)
	}

	return &progress, nil
}

// CreateProjectProgress creates an initial progress record for a project
func (r *ProgressRepository) CreateProjectProgress(projectID int) (*domain.ProjectProgress, error) {
	query := `
		INSERT INTO project_progress (project_id, current_phase)
		VALUES (?, ?)
		RETURNING project_id, setup_completed, attendees_added, features_added,
		          pairwise_value_completed, pairwise_complexity_completed,
		          fibonacci_value_completed, fibonacci_complexity_completed,
		          results_calculated, current_phase, last_activity, created_at, updated_at`

	var progress domain.ProjectProgress
	err := r.db.QueryRow(query, projectID, domain.PhaseSetup).Scan(
		&progress.ProjectID,
		&progress.SetupCompleted,
		&progress.AttendeesAdded,
		&progress.FeaturesAdded,
		&progress.PairwiseValueCompleted,
		&progress.PairwiseComplexityCompleted,
		&progress.FibonacciValueCompleted,
		&progress.FibonacciComplexityCompleted,
		&progress.ResultsCalculated,
		&progress.CurrentPhase,
		&progress.LastActivity,
		&progress.CreatedAt,
		&progress.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project progress: %w", err)
	}

	return &progress, nil
}

// UpdateProjectProgress updates the progress for a project
func (r *ProgressRepository) UpdateProjectProgress(progress *domain.ProjectProgress) error {
	query := `
		UPDATE project_progress SET
			setup_completed = ?,
			attendees_added = ?,
			features_added = ?,
			pairwise_value_completed = ?,
			pairwise_complexity_completed = ?,
			fibonacci_value_completed = ?,
			fibonacci_complexity_completed = ?,
			results_calculated = ?,
			current_phase = ?0,
			updated_at = datetime('now'),
			last_activity = datetime('now')
		WHERE project_id = ?`

	_, err := r.db.Exec(query,
		progress.ProjectID,
		progress.SetupCompleted,
		progress.AttendeesAdded,
		progress.FeaturesAdded,
		progress.PairwiseValueCompleted,
		progress.PairwiseComplexityCompleted,
		progress.FibonacciValueCompleted,
		progress.FibonacciComplexityCompleted,
		progress.ResultsCalculated,
		progress.CurrentPhase,
	)

	if err != nil {
		return fmt.Errorf("failed to update project progress: %w", err)
	}

	return nil
}

// MarkPhaseCompleted marks a specific phase as completed and advances to the next phase
func (r *ProgressRepository) MarkPhaseCompleted(projectID int, phase domain.WorkflowPhase) error {
	progress, err := r.GetProjectProgress(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project progress: %w", err)
	}

	// Mark the completed phase
	switch phase {
	case domain.PhaseSetup:
		progress.SetupCompleted = true
	case domain.PhaseAttendees:
		progress.AttendeesAdded = true
	case domain.PhaseFeatures:
		progress.FeaturesAdded = true
	case domain.PhasePairwiseValue:
		progress.PairwiseValueCompleted = true
	case domain.PhasePairwiseComplexity:
		progress.PairwiseComplexityCompleted = true
	case domain.PhaseFibonacciValue:
		progress.FibonacciValueCompleted = true
	case domain.PhaseFibonacciComplexity:
		progress.FibonacciComplexityCompleted = true
	case domain.PhaseResults:
		progress.ResultsCalculated = true
	}

	// Advance to next phase if not at the end
	if phase != domain.PhaseResults {
		progress.CurrentPhase = string(progress.GetNextPhase())
	}

	return r.UpdateProjectProgress(progress)
}

// DeleteProjectProgress deletes progress record for a project
func (r *ProgressRepository) DeleteProjectProgress(projectID int) error {
	query := `DELETE FROM project_progress WHERE project_id = ?`

	_, err := r.db.Exec(query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project progress: %w", err)
	}

	return nil
}
