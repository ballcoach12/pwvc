package repository

import (
	"database/sql"

	"pairwise/internal/domain"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(req domain.CreateProjectRequest) (*domain.Project, error) {
	query := `
		INSERT INTO projects (name, description, status, created_at, updated_at)
		VALUES (?, ?, 'active', datetime('now'), datetime('now'))
		RETURNING id, name, description, status, created_at, updated_at
	`

	var project domain.Project
	err := r.db.QueryRow(query, req.Name, req.Description).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id int) (*domain.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	var project domain.Project
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req domain.UpdateProjectRequest) (*domain.Project, error) {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, status = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING id, name, description, status, created_at, updated_at
	`

	status := req.Status
	if status == "" {
		status = "active"
	}

	var project domain.Project
	err := r.db.QueryRow(query, req.Name, req.Description, status, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &project, nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = ?`

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

// List retrieves all projects
func (r *ProjectRepository) List() ([]domain.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Status,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}
