package repository

import (
	"database/sql"

	"pairwise/internal/domain"
)

// AttendeeRepository handles database operations for attendees
type AttendeeRepository struct {
	db *sql.DB
}

// NewAttendeeRepository creates a new attendee repository
func NewAttendeeRepository(db *sql.DB) *AttendeeRepository {
	return &AttendeeRepository{db: db}
}

// Create creates a new attendee for a project
func (r *AttendeeRepository) Create(projectID int, req domain.CreateAttendeeRequest) (*domain.Attendee, error) {
	query := `
		INSERT INTO attendees (project_id, name, role, is_facilitator, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		RETURNING id, project_id, name, role, is_facilitator, created_at
	`

	var attendee domain.Attendee
	err := r.db.QueryRow(query, projectID, req.Name, req.Role, req.IsFacilitator).Scan(
		&attendee.ID,
		&attendee.ProjectID,
		&attendee.Name,
		&attendee.Role,
		&attendee.IsFacilitator,
		&attendee.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &attendee, nil
}

// GetByID retrieves an attendee by ID
func (r *AttendeeRepository) GetByID(id int) (*domain.Attendee, error) {
	query := `
		SELECT id, project_id, name, role, is_facilitator, created_at
		FROM attendees
		WHERE id = ?
	`

	var attendee domain.Attendee
	err := r.db.QueryRow(query, id).Scan(
		&attendee.ID,
		&attendee.ProjectID,
		&attendee.Name,
		&attendee.Role,
		&attendee.IsFacilitator,
		&attendee.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &attendee, nil
}

// GetByProjectID retrieves all attendees for a project
func (r *AttendeeRepository) GetByProjectID(projectID int) ([]domain.Attendee, error) {
	query := `
		SELECT id, project_id, name, role, is_facilitator, created_at
		FROM attendees
		WHERE project_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []domain.Attendee
	for rows.Next() {
		var attendee domain.Attendee
		err := rows.Scan(
			&attendee.ID,
			&attendee.ProjectID,
			&attendee.Name,
			&attendee.Role,
			&attendee.IsFacilitator,
			&attendee.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}

	return attendees, nil
}

// Delete deletes an attendee
func (r *AttendeeRepository) Delete(id int) error {
	query := `DELETE FROM attendees WHERE id = ?`

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

// DeleteByProjectID deletes all attendees for a project
func (r *AttendeeRepository) DeleteByProjectID(projectID int) error {
	query := `DELETE FROM attendees WHERE project_id = ?`
	_, err := r.db.Exec(query, projectID)
	return err
}
