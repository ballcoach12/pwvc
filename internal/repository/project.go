package repository

import (
	"database/sql"

	"pairwise/internal/domain"
)

// ProjectRepository interface defines the contract for project data operations
type ProjectRepository interface {
	Create(req domain.CreateProjectRequest) (*domain.Project, error)
	GetByID(id int) (*domain.Project, error)
	Update(id int, req domain.UpdateProjectRequest) (*domain.Project, error)
	Delete(id int) error
	List() ([]domain.Project, error)
	UpdateInviteCode(id int, inviteCode string) (*domain.Project, error)
	GetByInviteCode(inviteCode string) (*domain.Project, error)
}

// projectRepository handles database operations for projects
type projectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// Create creates a new project
func (r *projectRepository) Create(req domain.CreateProjectRequest) (*domain.Project, error) {
	query := `
		INSERT INTO projects (name, description, status, created_at, updated_at)
		VALUES (?, ?, 'active', datetime('now'), datetime('now'))
		RETURNING id, name, description, status, invite_code, created_at, updated_at
	`

	var project domain.Project
	var inviteCode sql.NullString
	err := r.db.QueryRow(query, req.Name, req.Description).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&inviteCode,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if inviteCode.Valid {
		project.InviteCode = inviteCode.String
	}

	return &project, nil
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(id int) (*domain.Project, error) {
	query := `
		SELECT id, name, description, status, invite_code, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	var project domain.Project
	var inviteCode sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&inviteCode,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if inviteCode.Valid {
		project.InviteCode = inviteCode.String
	}

	return &project, nil
}

// Update updates an existing project
func (r *projectRepository) Update(id int, req domain.UpdateProjectRequest) (*domain.Project, error) {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, status = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING id, name, description, status, invite_code, created_at, updated_at
	`

	status := req.Status
	if status == "" {
		status = "active"
	}

	var project domain.Project
	var inviteCode sql.NullString
	err := r.db.QueryRow(query, req.Name, req.Description, status, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&inviteCode,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if inviteCode.Valid {
		project.InviteCode = inviteCode.String
	}

	return &project, nil
}

// Delete deletes a project
func (r *projectRepository) Delete(id int) error {
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
func (r *projectRepository) List() ([]domain.Project, error) {
	query := `
		SELECT id, name, description, status, invite_code, created_at, updated_at
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
		var inviteCode sql.NullString
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Status,
			&inviteCode,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if inviteCode.Valid {
			project.InviteCode = inviteCode.String
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// UpdateInviteCode updates a project's invite code (T016 - US1)
func (r *projectRepository) UpdateInviteCode(id int, inviteCode string) (*domain.Project, error) {
	query := `
		UPDATE projects 
		SET invite_code = ?, updated_at = datetime('now')
		WHERE id = ?
		RETURNING id, name, description, status, invite_code, created_at, updated_at
	`

	var project domain.Project
	var inviteCodeResult sql.NullString
	err := r.db.QueryRow(query, inviteCode, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&inviteCodeResult,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if inviteCodeResult.Valid {
		project.InviteCode = inviteCodeResult.String
	}

	return &project, nil
}

// GetByInviteCode retrieves a project by its invite code (T016 - US1)
func (r *projectRepository) GetByInviteCode(inviteCode string) (*domain.Project, error) {
	query := `
		SELECT id, name, description, status, invite_code, created_at, updated_at
		FROM projects
		WHERE invite_code = ?
	`

	var project domain.Project
	var inviteCodeResult sql.NullString
	err := r.db.QueryRow(query, inviteCode).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&inviteCodeResult,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if inviteCodeResult.Valid {
		project.InviteCode = inviteCodeResult.String
	}

	return &project, nil
}
