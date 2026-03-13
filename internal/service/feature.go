package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

// FeatureService handles business logic for features
type FeatureService struct {
	featureRepo *repository.FeatureRepository
	projectRepo *repository.ProjectRepository
}

// NewFeatureService creates a new feature service
func NewFeatureService(featureRepo *repository.FeatureRepository, projectRepo *repository.ProjectRepository) *FeatureService {
	return &FeatureService{
		featureRepo: featureRepo,
		projectRepo: projectRepo,
	}
}

// CreateFeature creates a new feature with validation
func (s *FeatureService) CreateFeature(projectID int, req domain.CreateFeatureRequest) (*domain.Feature, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Validate project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate project", err.Error())
	}

	// Validate feature data
	if err := s.validateFeatureRequest(req.Title, req.Description, req.AcceptanceCriteria); err != nil {
		return nil, err
	}

	feature, err := s.featureRepo.Create(projectID, req)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to create feature", err.Error())
	}

	return feature, nil
}

// GetFeature retrieves a feature by ID
func (s *FeatureService) GetFeature(id int) (*domain.Feature, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid feature ID")
	}

	feature, err := s.featureRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Feature not found")
		}
		return nil, domain.NewAPIError(500, "Failed to retrieve feature", err.Error())
	}

	return feature, nil
}

// GetProjectFeatures retrieves all features for a project
func (s *FeatureService) GetProjectFeatures(projectID int) ([]domain.Feature, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Validate project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate project", err.Error())
	}

	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to retrieve features", err.Error())
	}

	// Return empty slice instead of nil if no features found
	if features == nil {
		features = []domain.Feature{}
	}

	return features, nil
}

// UpdateFeature updates an existing feature
func (s *FeatureService) UpdateFeature(id int, req domain.UpdateFeatureRequest) (*domain.Feature, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid feature ID")
	}

	// Validate feature data
	if err := s.validateFeatureRequest(req.Title, req.Description, req.AcceptanceCriteria); err != nil {
		return nil, err
	}

	feature, err := s.featureRepo.Update(id, req)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Feature not found")
		}
		return nil, domain.NewAPIError(500, "Failed to update feature", err.Error())
	}

	return feature, nil
}

// DeleteFeature deletes a feature
func (s *FeatureService) DeleteFeature(id int) error {
	if id <= 0 {
		return domain.NewAPIError(400, "Invalid feature ID")
	}

	err := s.featureRepo.Delete(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Feature not found")
		}
		return domain.NewAPIError(500, "Failed to delete feature", err.Error())
	}

	return nil
}

// ImportFeaturesFromCSV imports features from CSV data
func (s *FeatureService) ImportFeaturesFromCSV(projectID int, csvData io.Reader) (*domain.CSVImportResult, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Validate project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to validate project", err.Error())
	}

	reader := csv.NewReader(csvData)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, domain.NewAPIError(400, "Failed to read CSV headers", err.Error())
	}

	// Validate headers
	expectedHeaders := []string{"title", "description", "acceptance_criteria"}
	if !s.validateCSVHeaders(headers, expectedHeaders) {
		return nil, domain.NewAPIError(400, "Invalid CSV headers. Expected: title,description,acceptance_criteria")
	}

	var validFeatures []domain.CreateFeatureRequest
	var errors []string
	rowNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("Row %d: Failed to read CSV row - %s", rowNum+1, err.Error()))
			rowNum++
			continue
		}

		if len(record) < 2 {
			errors = append(errors, fmt.Sprintf("Row %d: Insufficient columns (minimum 2 required)", rowNum+1))
			rowNum++
			continue
		}

		// Parse feature data
		title := strings.TrimSpace(record[0])
		description := strings.TrimSpace(record[1])
		var acceptanceCriteria string
		if len(record) > 2 {
			acceptanceCriteria = strings.TrimSpace(record[2])
		}

		// Validate feature data
		if err := s.validateFeatureRequest(title, description, acceptanceCriteria); err != nil {
			if apiErr, ok := err.(*domain.APIError); ok {
				errors = append(errors, fmt.Sprintf("Row %d: %s", rowNum+1, apiErr.Message))
			} else {
				errors = append(errors, fmt.Sprintf("Row %d: %s", rowNum+1, err.Error()))
			}
			rowNum++
			continue
		}

		validFeatures = append(validFeatures, domain.CreateFeatureRequest{
			Title:              title,
			Description:        description,
			AcceptanceCriteria: acceptanceCriteria,
		})
		rowNum++
	}

	// Import valid features
	var importedCount int
	if len(validFeatures) > 0 {
		_, err := s.featureRepo.CreateBatch(projectID, validFeatures)
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to import features", err.Error())
		}
		importedCount = len(validFeatures)
	}

	result := &domain.CSVImportResult{
		ImportedCount: importedCount,
		SkippedCount:  len(errors),
		Errors:        errors,
	}

	return result, nil
}

// ExportFeaturesToCSV exports features to CSV format
func (s *FeatureService) ExportFeaturesToCSV(projectID int, writer io.Writer) error {
	if projectID <= 0 {
		return domain.NewAPIError(400, "Invalid project ID")
	}

	features, err := s.GetProjectFeatures(projectID)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	if err := csvWriter.Write([]string{"title", "description", "acceptance_criteria"}); err != nil {
		return domain.NewAPIError(500, "Failed to write CSV header", err.Error())
	}

	// Write feature data
	for _, feature := range features {
		record := []string{
			feature.Title,
			feature.Description,
			feature.AcceptanceCriteria,
		}
		if err := csvWriter.Write(record); err != nil {
			return domain.NewAPIError(500, "Failed to write CSV record", err.Error())
		}
	}

	return nil
}

// validateFeatureRequest validates feature data
func (s *FeatureService) validateFeatureRequest(title, description, acceptanceCriteria string) error {
	if title == "" {
		return domain.NewAPIError(400, "Feature title is required")
	}
	if len(title) > 255 {
		return domain.NewAPIError(400, "Feature title must be less than 255 characters")
	}

	if description == "" {
		return domain.NewAPIError(400, "Feature description is required")
	}
	if len(description) > 5000 {
		return domain.NewAPIError(400, "Feature description must be less than 5000 characters")
	}

	if len(acceptanceCriteria) > 5000 {
		return domain.NewAPIError(400, "Acceptance criteria must be less than 5000 characters")
	}

	return nil
}

// validateCSVHeaders checks if CSV headers match expected format
func (s *FeatureService) validateCSVHeaders(headers, expected []string) bool {
	if len(headers) < 2 {
		return false
	}

	// Check required headers (title and description)
	titleFound := false
	descriptionFound := false

	for _, header := range headers {
		header = strings.ToLower(strings.TrimSpace(header))
		if header == "title" {
			titleFound = true
		} else if header == "description" {
			descriptionFound = true
		}
	}

	return titleFound && descriptionFound
}
