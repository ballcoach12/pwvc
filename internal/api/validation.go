package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// InputValidator provides validation utilities for API input
type InputValidator struct{}

// NewInputValidator creates a new input validator
func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

// ValidateProject validates project creation/update data
func (v *InputValidator) ValidateProject(c *gin.Context, req interface{}) error {
	switch r := req.(type) {
	case *domain.CreateProjectRequest:
		return v.validateCreateProject(r)
	case *domain.UpdateProjectRequest:
		return v.validateUpdateProject(r)
	default:
		return domain.NewAPIError(http.StatusBadRequest, "Invalid request type")
	}
}

func (v *InputValidator) validateCreateProject(req *domain.CreateProjectRequest) error {
	// Validate name
	if err := v.validateProjectName(req.Name); err != nil {
		return err
	}

	// Validate description
	if err := v.validateDescription(req.Description); err != nil {
		return err
	}

	return nil
}

func (v *InputValidator) validateUpdateProject(req *domain.UpdateProjectRequest) error {
	// Validate name
	if err := v.validateProjectName(req.Name); err != nil {
		return err
	}

	// Validate description
	if err := v.validateDescription(req.Description); err != nil {
		return err
	}

	// Validate status if provided
	if req.Status != "" && !isValidProjectStatus(req.Status) {
		return &domain.ValidationError{
			Field:   "status",
			Message: "Status must be one of: active, inactive, completed",
		}
	}

	return nil
}

func (v *InputValidator) validateProjectName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name is required",
		}
	}

	if len(name) < 1 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name must be at least 1 character long",
		}
	}

	if len(name) > 255 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name must be no more than 255 characters long",
		}
	}

	// Check for valid characters (alphanumeric, spaces, hyphens, underscores)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_\.]+$`)
	if !validNameRegex.MatchString(name) {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name contains invalid characters. Use only letters, numbers, spaces, hyphens, underscores, and periods",
		}
	}

	return nil
}

func (v *InputValidator) validateDescription(description string) error {
	if len(description) > 1000 {
		return &domain.ValidationError{
			Field:   "description",
			Message: "Description must be no more than 1000 characters long",
		}
	}

	// Check for potentially harmful content
	if containsSuspiciousContent(description) {
		return &domain.ValidationError{
			Field:   "description",
			Message: "Description contains invalid content",
		}
	}

	return nil
}

// ValidateAttendee validates attendee data
func (v *InputValidator) ValidateAttendee(name, email string) error {
	// Validate name
	if err := v.validateAttendeeName(name); err != nil {
		return err
	}

	// Validate email
	if err := v.validateEmail(email); err != nil {
		return err
	}

	return nil
}

func (v *InputValidator) validateAttendeeName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name is required",
		}
	}

	if len(name) < 2 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name must be at least 2 characters long",
		}
	}

	if len(name) > 100 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name must be no more than 100 characters long",
		}
	}

	// Check that name contains at least some letters
	hasLetter := false
	for _, r := range name {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}

	if !hasLetter {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Name must contain at least one letter",
		}
	}

	return nil
}

func (v *InputValidator) validateEmail(email string) error {
	if email == "" {
		return &domain.ValidationError{
			Field:   "email",
			Message: "Email is required",
		}
	}

	// Simple email regex (more comprehensive validation can be added)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &domain.ValidationError{
			Field:   "email",
			Message: "Please provide a valid email address",
		}
	}

	if len(email) > 255 {
		return &domain.ValidationError{
			Field:   "email",
			Message: "Email must be no more than 255 characters long",
		}
	}

	return nil
}

// ValidateFeature validates feature data
func (v *InputValidator) ValidateFeature(name, description string) error {
	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Feature name is required",
		}
	}

	if len(name) < 3 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Feature name must be at least 3 characters long",
		}
	}

	if len(name) > 255 {
		return &domain.ValidationError{
			Field:   "name",
			Message: "Feature name must be no more than 255 characters long",
		}
	}

	// Validate description
	if len(description) > 2000 {
		return &domain.ValidationError{
			Field:   "description",
			Message: "Feature description must be no more than 2000 characters long",
		}
	}

	return nil
}

// ValidateWorkflowPhase validates workflow phase transitions
func (v *InputValidator) ValidateWorkflowPhase(phase string) error {
	validPhases := []string{
		string(domain.PhaseSetup),
		string(domain.PhaseAttendees),
		string(domain.PhaseFeatures),
		string(domain.PhasePairwiseValue),
		string(domain.PhasePairwiseComplexity),
		string(domain.PhaseFibonacciValue),
		string(domain.PhaseFibonacciComplexity),
		string(domain.PhaseResults),
	}

	for _, validPhase := range validPhases {
		if phase == validPhase {
			return nil
		}
	}

	return &domain.ValidationError{
		Field:   "phase",
		Message: fmt.Sprintf("Invalid phase. Must be one of: %s", strings.Join(validPhases, ", ")),
	}
}

// Utility functions
func isValidProjectStatus(status string) bool {
	validStatuses := []string{"active", "inactive", "completed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func containsSuspiciousContent(text string) bool {
	// Simple check for potentially harmful content
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
	}

	lowerText := strings.ToLower(text)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}

	return false
}

// SanitizeInput removes potentially harmful characters from input
func SanitizeInput(input string) string {
	// Remove null bytes and control characters
	var builder strings.Builder
	for _, r := range input {
		if r >= 32 && r != 127 { // Keep printable characters except DEL
			builder.WriteRune(r)
		}
	}

	result := builder.String()

	// Trim whitespace
	result = strings.TrimSpace(result)

	return result
}
