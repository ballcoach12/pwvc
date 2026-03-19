package domain

import (
	"errors"
	"fmt"
	"time"
)

// Common domain errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrValidation         = errors.New("validation error")
	ErrDuplicate          = errors.New("duplicate resource")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternalError      = errors.New("internal error")
	ErrBadRequest         = errors.New("bad request")
	ErrConflict           = errors.New("conflict")
	ErrForbidden          = errors.New("forbidden")
	ErrTimeout            = errors.New("timeout")
	ErrRateLimit          = errors.New("rate limit exceeded")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// APIError represents a structured API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code int, message string, details ...string) *APIError {
	err := &APIError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// BusinessError represents domain-specific business rule violations
type BusinessError struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("business rule violation [%s]: %s", e.Type, e.Message)
}

// NewBusinessError creates a new business error
func NewBusinessError(errorType, message string, context ...map[string]interface{}) *BusinessError {
	err := &BusinessError{
		Type:    errorType,
		Message: message,
	}
	if len(context) > 0 {
		err.Context = context[0]
	}
	return err
}

// Common business error types
const (
	ErrTypeInsufficientData  = "insufficient_data"
	ErrTypeInvalidPhase      = "invalid_phase"
	ErrTypePhaseNotComplete  = "phase_not_complete"
	ErrTypeSessionInProgress = "session_in_progress"
	ErrTypeDataInconsistency = "data_inconsistency"
	ErrTypeWorkflowViolation = "workflow_violation"
	ErrTypeResourceConflict  = "resource_conflict"
	ErrTypeInvalidTransition = "invalid_transition"
)

// AuditReportOptions defines options for generating audit reports (T044 - US9)
type AuditReportOptions struct {
	Limit               int    `json:"limit,omitempty"`
	ActionType          string `json:"action_type,omitempty"`
	IncludePersonalData bool   `json:"include_personal_data,omitempty"`
}

// AuditReport represents a comprehensive audit report with privacy controls (T044 - US9)
type AuditReport struct {
	ProjectID    int                `json:"project_id"`
	ProjectName  string             `json:"project_name"`
	GeneratedAt  time.Time          `json:"generated_at"`
	TotalActions int                `json:"total_actions"`
	DateRange    DateRange          `json:"date_range"`
	Statistics   AuditStatistics    `json:"statistics"`
	AuditLogs    []AuditLog         `json:"audit_logs"`
	PrivacyMode  bool               `json:"privacy_mode"`
	Options      AuditReportOptions `json:"options"`
}

// AuditStatistics provides summary statistics for audit reports (T044 - US9)
type AuditStatistics struct {
	ActionCounts           map[string]int `json:"action_counts"`
	EntityTypeCounts       map[string]int `json:"entity_type_counts"`
	AttendeeActivityCounts map[int]int    `json:"attendee_activity_counts"`
	TotalActions           int            `json:"total_actions"`
}

// DateRange represents a time range for audit reporting (T044 - US9)
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
