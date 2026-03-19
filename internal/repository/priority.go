package repository

import (
	"database/sql"
	"fmt"
	"pairwise/internal/domain"
)

// PriorityRepository interface defines the contract for priority calculation data operations
type PriorityRepository interface {
	Create(calc *domain.PriorityCalculation) error
	GetByProjectID(projectID int) ([]domain.PriorityCalculation, error)
	GetResultsWithFeatures(projectID int) ([]domain.PriorityResult, error)
	DeleteByProjectID(projectID int) error
	ExistsForProject(projectID int) (bool, error)
	GetLatestCalculationTime(projectID int) (*domain.PriorityCalculation, error)
}

// priorityRepository handles database operations for priority calculations
type priorityRepository struct {
	db *sql.DB
}

// NewPriorityRepository creates a new priority repository
func NewPriorityRepository(db *sql.DB) PriorityRepository {
	return &priorityRepository{
		db: db,
	}
}

// Create inserts a new priority calculation
func (r *priorityRepository) Create(calc *domain.PriorityCalculation) error {
	query := `
		INSERT INTO priority_calculations (
			project_id, feature_id, w_value, w_complexity, s_value, s_complexity,
			weighted_value, weighted_complexity, final_priority_score, rank
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, calculated_at`

	err := r.db.QueryRow(
		query,
		calc.ProjectID, calc.FeatureID, calc.WValue, calc.WComplexity,
		calc.SValue, calc.SComplexity, calc.WeightedValue, calc.WeightedComplexity,
		calc.FinalPriorityScore, calc.Rank,
	).Scan(&calc.ID, &calc.CalculatedAt)

	return err
}

// GetByProjectID retrieves all priority calculations for a project, ordered by rank
func (r *priorityRepository) GetByProjectID(projectID int) ([]domain.PriorityCalculation, error) {
	query := `
		SELECT pc.id, pc.project_id, pc.feature_id, pc.w_value, pc.w_complexity,
		       pc.s_value, pc.s_complexity, pc.weighted_value, pc.weighted_complexity,
		       pc.final_priority_score, pc.rank, pc.calculated_at
		FROM priority_calculations pc
		WHERE pc.project_id = ?
		ORDER BY pc.rank ASC`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calculations []domain.PriorityCalculation
	for rows.Next() {
		var calc domain.PriorityCalculation
		err := rows.Scan(
			&calc.ID, &calc.ProjectID, &calc.FeatureID, &calc.WValue, &calc.WComplexity,
			&calc.SValue, &calc.SComplexity, &calc.WeightedValue, &calc.WeightedComplexity,
			&calc.FinalPriorityScore, &calc.Rank, &calc.CalculatedAt,
		)
		if err != nil {
			return nil, err
		}
		calculations = append(calculations, calc)
	}

	return calculations, rows.Err()
}

// GetResultsWithFeatures retrieves priority calculations with feature details
func (r *priorityRepository) GetResultsWithFeatures(projectID int) ([]domain.PriorityResult, error) {
	query := `
		SELECT pc.id, pc.project_id, pc.feature_id, pc.w_value, pc.w_complexity,
		       pc.s_value, pc.s_complexity, pc.weighted_value, pc.weighted_complexity,
		       pc.final_priority_score, pc.rank, pc.calculated_at,
		       f.id, f.project_id, f.title, f.description, f.acceptance_criteria,
		       f.created_at, f.updated_at
		FROM priority_calculations pc
		JOIN features f ON pc.feature_id = f.id
		WHERE pc.project_id = ?
		ORDER BY pc.rank ASC`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.PriorityResult
	for rows.Next() {
		var result domain.PriorityResult
		err := rows.Scan(
			&result.ID, &result.ProjectID, &result.FeatureID, &result.WValue, &result.WComplexity,
			&result.SValue, &result.SComplexity, &result.WeightedValue, &result.WeightedComplexity,
			&result.FinalPriorityScore, &result.Rank, &result.CalculatedAt,
			&result.Feature.ID, &result.Feature.ProjectID, &result.Feature.Title,
			&result.Feature.Description, &result.Feature.AcceptanceCriteria,
			&result.Feature.CreatedAt, &result.Feature.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, rows.Err()
}

// DeleteByProjectID removes all priority calculations for a project
func (r *priorityRepository) DeleteByProjectID(projectID int) error {
	query := "DELETE FROM priority_calculations WHERE project_id = ?"
	_, err := r.db.Exec(query, projectID)
	return err
}

// ExistsForProject checks if priority calculations exist for a project
func (r *priorityRepository) ExistsForProject(projectID int) (bool, error) {
	query := "SELECT COUNT(*) FROM priority_calculations WHERE project_id = ?"

	var count int
	err := r.db.QueryRow(query, projectID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetLatestCalculationTime returns the most recent calculation time for a project
func (r *priorityRepository) GetLatestCalculationTime(projectID int) (*domain.PriorityCalculation, error) {
	query := `
		SELECT id, project_id, feature_id, w_value, w_complexity, s_value, s_complexity,
		       weighted_value, weighted_complexity, final_priority_score, rank, calculated_at
		FROM priority_calculations 
		WHERE project_id = ? 
		ORDER BY calculated_at DESC 
		LIMIT 1`

	var calc domain.PriorityCalculation
	err := r.db.QueryRow(query, projectID).Scan(
		&calc.ID, &calc.ProjectID, &calc.FeatureID, &calc.WValue, &calc.WComplexity,
		&calc.SValue, &calc.SComplexity, &calc.WeightedValue, &calc.WeightedComplexity,
		&calc.FinalPriorityScore, &calc.Rank, &calc.CalculatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get latest calculation: %w", err)
	}

	return &calc, nil
}
