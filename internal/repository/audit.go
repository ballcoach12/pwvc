package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pairwise/internal/domain"

	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(auditLog *domain.AuditLog) error
	LogAction(projectID int, actorID *int, actionType, subjectType string, subjectID int, beforeState, afterState interface{}) error
	GetAuditLogs(projectID int, limit int, actionType string) ([]domain.AuditLog, error)
	GetAuditLogsBySubject(projectID int, subjectType string, subjectID int) ([]domain.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

// Create inserts a new audit log entry (T031 - US4)
func (r *auditRepository) Create(auditLog *domain.AuditLog) error {
	if err := r.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

// LogAction creates an audit log entry (backwards compatibility)
func (r *auditRepository) LogAction(projectID int, actorID *int, actionType, subjectType string, subjectID int, beforeState, afterState interface{}) error {
	var beforeJSON, afterJSON []byte
	var err error

	if beforeState != nil {
		beforeJSON, err = json.Marshal(beforeState)
		if err != nil {
			return err
		}
	}

	if afterState != nil {
		afterJSON, err = json.Marshal(afterState)
		if err != nil {
			return err
		}
	}

	auditLog := &domain.AuditLog{
		ProjectID:  projectID,
		AttendeeID: 0, // Will be set from actorID if provided
		Action:     actionType,
		EntityType: subjectType,
		EntityID:   fmt.Sprintf("%d", subjectID),
		OldValue:   string(beforeJSON),
		NewValue:   string(afterJSON),
		Timestamp:  time.Now(),
	}

	if actorID != nil {
		auditLog.AttendeeID = *actorID
	}

	return r.Create(auditLog)
}

// GetAuditLogs retrieves audit logs for a project with optional filtering
func (r *auditRepository) GetAuditLogs(projectID int, limit int, actionType string) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog

	query := r.db.Where("project_id = ?", projectID)

	if actionType != "" {
		query = query.Where("action = ?", actionType)
	}

	query = query.Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

// GetAuditLogsBySubject retrieves audit logs for a specific subject
func (r *auditRepository) GetAuditLogsBySubject(projectID int, subjectType string, subjectID int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog

	entityID := fmt.Sprintf("%d", subjectID)

	query := r.db.Where("project_id = ? AND entity_type = ? AND entity_id = ?", projectID, subjectType, entityID).
		Order("timestamp ASC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get audit logs by subject: %w", err)
	}

	return logs, nil
}
