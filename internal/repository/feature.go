package repository

import (
	"database/sql"

	"pairwise/internal/domain"
)

// FeatureRepository handles database operations for features
type FeatureRepository struct {
	db *sql.DB
}

// NewFeatureRepository creates a new feature repository
func NewFeatureRepository(db *sql.DB) *FeatureRepository {
	return &FeatureRepository{db: db}
}

// Create creates a new feature
func (r *FeatureRepository) Create(projectID int, req domain.CreateFeatureRequest) (*domain.Feature, error) {
	query := `
		INSERT INTO features (project_id, title, description, acceptance_criteria, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		RETURNING id, project_id, title, description, acceptance_criteria, created_at, updated_at
	`

	var feature domain.Feature
	err := r.db.QueryRow(query, projectID, req.Title, req.Description, req.AcceptanceCriteria).Scan(
		&feature.ID,
		&feature.ProjectID,
		&feature.Title,
		&feature.Description,
		&feature.AcceptanceCriteria,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &feature, nil
}

// GetByID retrieves a feature by ID
func (r *FeatureRepository) GetByID(id int) (*domain.Feature, error) {
	query := `
		SELECT id, project_id, title, description, acceptance_criteria, created_at, updated_at
		FROM features
		WHERE id = ?
	`

	var feature domain.Feature
	err := r.db.QueryRow(query, id).Scan(
		&feature.ID,
		&feature.ProjectID,
		&feature.Title,
		&feature.Description,
		&feature.AcceptanceCriteria,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &feature, nil
}

// GetByProjectID retrieves all features for a project
func (r *FeatureRepository) GetByProjectID(projectID int) ([]domain.Feature, error) {
	query := `
		SELECT id, project_id, title, description, acceptance_criteria, created_at, updated_at
		FROM features
		WHERE project_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []domain.Feature
	for rows.Next() {
		var feature domain.Feature
		err := rows.Scan(
			&feature.ID,
			&feature.ProjectID,
			&feature.Title,
			&feature.Description,
			&feature.AcceptanceCriteria,
			&feature.CreatedAt,
			&feature.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, nil
}

// Update updates an existing feature
func (r *FeatureRepository) Update(id int, req domain.UpdateFeatureRequest) (*domain.Feature, error) {
	query := `
		UPDATE features 
		SET title = ?, description = ?, acceptance_criteria = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING id, project_id, title, description, acceptance_criteria, created_at, updated_at
	`

	var feature domain.Feature
	err := r.db.QueryRow(query, id, req.Title, req.Description, req.AcceptanceCriteria).Scan(
		&feature.ID,
		&feature.ProjectID,
		&feature.Title,
		&feature.Description,
		&feature.AcceptanceCriteria,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &feature, nil
}

// Delete deletes a feature
func (r *FeatureRepository) Delete(id int) error {
	query := `DELETE FROM features WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteByProjectID deletes all features for a project
func (r *FeatureRepository) DeleteByProjectID(projectID int) error {
	query := `DELETE FROM features WHERE project_id = ?`
	_, err := r.db.Exec(query, projectID)
	return err
}

// CreateBatch creates multiple features in a single transaction
func (r *FeatureRepository) CreateBatch(projectID int, features []domain.CreateFeatureRequest) ([]domain.Feature, error) {
	if len(features) == 0 {
		return []domain.Feature{}, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO features (project_id, title, description, acceptance_criteria, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
		RETURNING id, project_id, title, description, acceptance_criteria, created_at, updated_at
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var createdFeatures []domain.Feature
	for _, req := range features {
		var feature domain.Feature
		err := stmt.QueryRow(projectID, req.Title, req.Description, req.AcceptanceCriteria).Scan(
			&feature.ID,
			&feature.ProjectID,
			&feature.Title,
			&feature.Description,
			&feature.AcceptanceCriteria,
			&feature.CreatedAt,
			&feature.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		createdFeatures = append(createdFeatures, feature)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdFeatures, nil
}
