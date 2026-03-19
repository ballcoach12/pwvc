package repository

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"pairwise/internal/domain"
)

// AttendeeRepository interface defines the contract for attendee data operations
type AttendeeRepository interface {
	Create(projectID int, req domain.CreateAttendeeRequest) (*domain.Attendee, error)
	CreateWithoutPIN(projectID int, req domain.CreateAttendeeWithoutPINRequest) (*domain.Attendee, error)
	GetByID(id int) (*domain.Attendee, error)
	GetByProjectID(projectID int) ([]domain.Attendee, error)
	GetByInviteToken(token string) (*domain.Attendee, error)
	SetPIN(attendeeID int, pin string) error
	GenerateInviteToken(attendeeID int) (string, error)
	ClearInviteToken(attendeeID int) error
	Delete(id int) error
	DeleteByProjectID(projectID int) error
}

// attendeeRepository handles database operations for attendees
type attendeeRepository struct {
	db *sql.DB
}

// NewAttendeeRepository creates a new attendee repository
func NewAttendeeRepository(db *sql.DB) AttendeeRepository {
	return &attendeeRepository{db: db}
}

// Create creates a new attendee for a project
func (r *attendeeRepository) Create(projectID int, req domain.CreateAttendeeRequest) (*domain.Attendee, error) {
	// Hash the PIN
	pinHash := hashPIN(req.PIN)

	query := `
		INSERT INTO attendees (project_id, name, role, is_facilitator, pin_hash, created_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'))
		RETURNING id, project_id, name, role, is_facilitator, created_at
	`

	var attendee domain.Attendee
	err := r.db.QueryRow(query, projectID, req.Name, req.Role, req.IsFacilitator, pinHash).Scan(
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

	// Set the pin hash (but it won't be returned in JSON due to the json:"-" tag)
	attendee.PinHash = &pinHash

	return &attendee, nil
}

// CreateWithoutPIN creates a new attendee without requiring a PIN (for invite workflow)
func (r *attendeeRepository) CreateWithoutPIN(projectID int, req domain.CreateAttendeeWithoutPINRequest) (*domain.Attendee, error) {
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
func (r *attendeeRepository) GetByID(id int) (*domain.Attendee, error) {
	query := `
		SELECT id, project_id, name, role, is_facilitator, pin_hash, invite_token, invite_token_expires_at, created_at
		FROM attendees
		WHERE id = ?
	`

	var attendee domain.Attendee
	var pinHash, inviteToken sql.NullString
	var inviteTokenExpiresAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&attendee.ID,
		&attendee.ProjectID,
		&attendee.Name,
		&attendee.Role,
		&attendee.IsFacilitator,
		&pinHash,
		&inviteToken,
		&inviteTokenExpiresAt,
		&attendee.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Convert nullable fields
	if pinHash.Valid {
		attendee.PinHash = &pinHash.String
	}
	if inviteToken.Valid {
		attendee.InviteToken = &inviteToken.String
	}
	if inviteTokenExpiresAt.Valid {
		attendee.InviteTokenExpiresAt = &inviteTokenExpiresAt.Time
	}

	return &attendee, nil
}

// GetByProjectID retrieves all attendees for a project
func (r *attendeeRepository) GetByProjectID(projectID int) ([]domain.Attendee, error) {
	query := `
		SELECT id, project_id, name, role, is_facilitator, pin_hash, invite_token, invite_token_expires_at, created_at
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
		var pinHash, inviteToken sql.NullString
		var inviteTokenExpiresAt sql.NullTime

		err := rows.Scan(
			&attendee.ID,
			&attendee.ProjectID,
			&attendee.Name,
			&attendee.Role,
			&attendee.IsFacilitator,
			&pinHash,
			&inviteToken,
			&inviteTokenExpiresAt,
			&attendee.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert nullable fields
		if pinHash.Valid {
			attendee.PinHash = &pinHash.String
		}
		if inviteToken.Valid {
			attendee.InviteToken = &inviteToken.String
		}
		if inviteTokenExpiresAt.Valid {
			attendee.InviteTokenExpiresAt = &inviteTokenExpiresAt.Time
		}

		attendees = append(attendees, attendee)
	}

	return attendees, nil
}

// Delete deletes an attendee
func (r *attendeeRepository) Delete(id int) error {
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
func (r *attendeeRepository) DeleteByProjectID(projectID int) error {
	query := `DELETE FROM attendees WHERE project_id = ?`
	_, err := r.db.Exec(query, projectID)
	return err
}

// SetPIN sets or updates the PIN for an attendee
func (r *attendeeRepository) SetPIN(attendeeID int, pin string) error {
	pinHash := hashPIN(pin)

	query := `UPDATE attendees SET pin_hash = ? WHERE id = ?`
	result, err := r.db.Exec(query, pinHash, attendeeID)
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

// GetByInviteToken retrieves an attendee by their invite token
func (r *attendeeRepository) GetByInviteToken(token string) (*domain.Attendee, error) {
	query := `
		SELECT id, project_id, name, role, is_facilitator, pin_hash, invite_token, invite_token_expires_at, created_at
		FROM attendees
		WHERE invite_token = ? AND invite_token_expires_at > datetime('now')
	`

	var attendee domain.Attendee
	var pinHash, inviteToken sql.NullString
	var inviteTokenExpiresAt sql.NullTime

	err := r.db.QueryRow(query, token).Scan(
		&attendee.ID,
		&attendee.ProjectID,
		&attendee.Name,
		&attendee.Role,
		&attendee.IsFacilitator,
		&pinHash,
		&inviteToken,
		&inviteTokenExpiresAt,
		&attendee.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Convert nullable fields
	if pinHash.Valid {
		attendee.PinHash = &pinHash.String
	}
	if inviteToken.Valid {
		attendee.InviteToken = &inviteToken.String
	}
	if inviteTokenExpiresAt.Valid {
		attendee.InviteTokenExpiresAt = &inviteTokenExpiresAt.Time
	}

	return &attendee, nil
}

// GenerateInviteToken generates a new invite token for an attendee
func (r *attendeeRepository) GenerateInviteToken(attendeeID int) (string, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Set expiry to 24 hours from now
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `UPDATE attendees SET invite_token = ?, invite_token_expires_at = ? WHERE id = ?`
	result, err := r.db.Exec(query, token, expiresAt, attendeeID)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", domain.ErrNotFound
	}

	return token, nil
}

// ClearInviteToken clears the invite token for an attendee
func (r *attendeeRepository) ClearInviteToken(attendeeID int) error {
	query := `UPDATE attendees SET invite_token = NULL, invite_token_expires_at = NULL WHERE id = ?`
	result, err := r.db.Exec(query, attendeeID)
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

// hashPIN creates a simple hash of the PIN
func hashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}
