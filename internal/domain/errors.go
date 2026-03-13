package domain

import (
	"errors"
	"fmt"
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
