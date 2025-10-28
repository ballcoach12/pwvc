package domain

import (
	"time"
)

// Feature represents a P-WVC project feature
type Feature struct {
	ID                 int       `json:"id" db:"id"`
	ProjectID          int       `json:"project_id" db:"project_id"`
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	AcceptanceCriteria string    `json:"acceptance_criteria" db:"acceptance_criteria"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// CreateFeatureRequest represents the request payload for creating a feature
type CreateFeatureRequest struct {
	Title              string `json:"title" binding:"required,min=1,max=255"`
	Description        string `json:"description" binding:"required,min=1,max=5000"`
	AcceptanceCriteria string `json:"acceptance_criteria" binding:"omitempty,max=5000"`
}

// UpdateFeatureRequest represents the request payload for updating a feature
type UpdateFeatureRequest struct {
	Title              string `json:"title" binding:"required,min=1,max=255"`
	Description        string `json:"description" binding:"required,min=1,max=5000"`
	AcceptanceCriteria string `json:"acceptance_criteria" binding:"omitempty,max=5000"`
}

// FeatureImportRequest represents a single feature from CSV import
type FeatureImportRequest struct {
	Title              string `csv:"title"`
	Description        string `csv:"description"`
	AcceptanceCriteria string `csv:"acceptance_criteria"`
}

// FeatureExportResponse represents a feature for CSV export
type FeatureExportResponse struct {
	Title              string `csv:"title"`
	Description        string `csv:"description"`
	AcceptanceCriteria string `csv:"acceptance_criteria"`
}

// CSVImportResult represents the result of a CSV import operation
type CSVImportResult struct {
	ImportedCount int      `json:"imported_count"`
	SkippedCount  int      `json:"skipped_count"`
	Errors        []string `json:"errors,omitempty"`
}
