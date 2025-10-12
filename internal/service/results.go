package service

import (
	"fmt"
	"sort"
	"time"

	"pwvc/internal/domain"
	"pwvc/internal/repository"
)

// ResultsService handles P-WVC results calculation and management
type ResultsService struct {
	priorityRepo *repository.PriorityRepository
	featureRepo  *repository.FeatureRepository
	pairwiseRepo *repository.PairwiseRepository
}

// NewResultsService creates a new results service
func NewResultsService(
	priorityRepo *repository.PriorityRepository,
	featureRepo *repository.FeatureRepository,
	pairwiseRepo *repository.PairwiseRepository,
) *ResultsService {
	return &ResultsService{
		priorityRepo: priorityRepo,
		featureRepo:  featureRepo,
		pairwiseRepo: pairwiseRepo,
	}
}

// CalculateResults performs the complete P-WVC calculation for a project
func (s *ResultsService) CalculateResults(projectID int) (*domain.ProjectResults, error) {
	// 1. Get all features for the project
	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get project features", err.Error())
	}

	if len(features) == 0 {
		return nil, domain.NewAPIError(400, "No features found for this project")
	}

	// 2. Get win-count weights from pairwise comparisons
	valueWeights, err := s.calculateWinCountWeights(projectID, "value")
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to calculate value weights", err.Error())
	}

	complexityWeights, err := s.calculateWinCountWeights(projectID, "complexity")
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to calculate complexity weights", err.Error())
	}

	// 3. Get Fibonacci consensus scores (placeholder - we'll implement this)
	valueScores, err := s.getFibonacciScores(projectID, "value")
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get value scores", err.Error())
	}

	complexityScores, err := s.getFibonacciScores(projectID, "complexity")
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get complexity scores", err.Error())
	}

	// 4. Calculate Final Priority Scores
	var calculations []domain.PriorityCalculation
	for _, feature := range features {
		wValue := valueWeights[feature.ID]
		wComplexity := complexityWeights[feature.ID]
		sValue := valueScores[feature.ID]
		sComplexity := complexityScores[feature.ID]

		// Prevent division by zero
		if wComplexity == 0 || sComplexity == 0 {
			return nil, domain.NewAPIError(400, fmt.Sprintf("Invalid complexity values for feature %d", feature.ID))
		}

		weightedValue := float64(sValue) * wValue
		weightedComplexity := float64(sComplexity) * wComplexity
		finalScore := weightedValue / weightedComplexity

		calculation := domain.PriorityCalculation{
			ProjectID:          projectID,
			FeatureID:          feature.ID,
			WValue:             wValue,
			WComplexity:        wComplexity,
			SValue:             sValue,
			SComplexity:        sComplexity,
			WeightedValue:      weightedValue,
			WeightedComplexity: weightedComplexity,
			FinalPriorityScore: finalScore,
		}

		calculations = append(calculations, calculation)
	}

	// 5. Sort by Final Priority Score (descending) and assign ranks
	sort.Slice(calculations, func(i, j int) bool {
		return calculations[i].FinalPriorityScore > calculations[j].FinalPriorityScore
	})

	for i := range calculations {
		calculations[i].Rank = i + 1
	}

	// 6. Clear existing calculations and save new ones
	err = s.priorityRepo.DeleteByProjectID(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to clear existing calculations", err.Error())
	}

	for i := range calculations {
		err = s.priorityRepo.Create(&calculations[i])
		if err != nil {
			return nil, domain.NewAPIError(500, "Failed to save calculations", err.Error())
		}
	}

	// 7. Get results with feature details
	results, err := s.priorityRepo.GetResultsWithFeatures(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get calculated results", err.Error())
	}

	// 8. Calculate summary statistics
	summary := s.calculateSummary(results)

	return &domain.ProjectResults{
		ProjectID:     projectID,
		Results:       results,
		CalculatedAt:  time.Now(),
		TotalFeatures: len(results),
		Summary:       summary,
	}, nil
}

// GetResults retrieves existing calculated results for a project
func (s *ResultsService) GetResults(projectID int) (*domain.ProjectResults, error) {
	results, err := s.priorityRepo.GetResultsWithFeatures(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to get results", err.Error())
	}

	if len(results) == 0 {
		return nil, domain.NewAPIError(404, "No results found for this project. Run calculation first.")
	}

	summary := s.calculateSummary(results)

	return &domain.ProjectResults{
		ProjectID:     projectID,
		Results:       results,
		CalculatedAt:  results[0].CalculatedAt,
		TotalFeatures: len(results),
		Summary:       summary,
	}, nil
}

// calculateWinCountWeights calculates win-count weights from pairwise comparison results
func (s *ResultsService) calculateWinCountWeights(projectID int, criterionType string) (map[int]float64, error) {
	// Get all pairwise comparison results for this project and criterion
	// This is a simplified implementation - in reality we'd need to aggregate votes properly

	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	weights := make(map[int]float64)

	// For now, assign mock weights - this should be calculated from actual pairwise results
	totalFeatures := len(features)
	for i, feature := range features {
		// Mock calculation: higher ranking features get higher weights
		weight := 1.0 - (float64(i) / float64(totalFeatures))
		if weight < 0.1 {
			weight = 0.1 // Minimum weight
		}
		weights[feature.ID] = weight
	}

	return weights, nil
}

// getFibonacciScores retrieves consensus Fibonacci scores for features
func (s *ResultsService) getFibonacciScores(projectID int, criterionType string) (map[int]int, error) {
	// This should get actual Fibonacci consensus scores
	// For now, return mock scores

	features, err := s.featureRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	scores := make(map[int]int)
	fibValues := []int{1, 2, 3, 5, 8, 13, 21}

	for i, feature := range features {
		// Mock: assign random but deterministic Fibonacci values
		scoreIndex := (feature.ID + i) % len(fibValues)
		scores[feature.ID] = fibValues[scoreIndex]
	}

	return scores, nil
}

// calculateSummary generates statistical summary of results
func (s *ResultsService) calculateSummary(results []domain.PriorityResult) domain.ResultsSummary {
	if len(results) == 0 {
		return domain.ResultsSummary{}
	}

	scores := make([]float64, len(results))
	var total float64

	for i, result := range results {
		score := result.FinalPriorityScore
		scores[i] = score
		total += score
	}

	sort.Float64s(scores)

	highest := scores[len(scores)-1]
	lowest := scores[0]
	average := total / float64(len(scores))

	var median float64
	mid := len(scores) / 2
	if len(scores)%2 == 0 {
		median = (scores[mid-1] + scores[mid]) / 2
	} else {
		median = scores[mid]
	}

	// Calculate quartiles
	q1Index := len(scores) / 4
	q3Index := 3 * len(scores) / 4

	topTier := len(scores) - q3Index
	bottomTier := q1Index

	return domain.ResultsSummary{
		HighestScore: highest,
		LowestScore:  lowest,
		AverageScore: average,
		MedianScore:  median,
		ScoreRange:   highest - lowest,
		TopTier:      topTier,
		BottomTier:   bottomTier,
	}
}

// ExportResults exports results in the specified format
func (s *ResultsService) ExportResults(projectID int, format domain.ExportFormat) (interface{}, error) {
	results, err := s.GetResults(projectID)
	if err != nil {
		return nil, err
	}

	switch format {
	case domain.ExportFormatCSV:
		return s.exportToCSV(results), nil
	case domain.ExportFormatJSON:
		return results, nil
	case domain.ExportFormatJira:
		return s.exportToJira(results), nil
	default:
		return nil, domain.NewAPIError(400, "Invalid export format")
	}
}

// exportToCSV converts results to CSV format
func (s *ResultsService) exportToCSV(results *domain.ProjectResults) [][]string {
	csv := [][]string{
		{"rank", "feature_title", "description", "final_priority_score", "s_value", "s_complexity", "w_value", "w_complexity"},
	}

	for _, result := range results.Results {
		row := []string{
			fmt.Sprintf("%d", result.Rank),
			result.Feature.Title,
			result.Feature.Description,
			fmt.Sprintf("%.6f", result.FinalPriorityScore),
			fmt.Sprintf("%d", result.SValue),
			fmt.Sprintf("%d", result.SComplexity),
			fmt.Sprintf("%.6f", result.WValue),
			fmt.Sprintf("%.6f", result.WComplexity),
		}
		csv = append(csv, row)
	}

	return csv
}

// exportToJira converts results to Jira-compatible format
func (s *ResultsService) exportToJira(results *domain.ProjectResults) domain.JiraExport {
	var issues []domain.JiraIssue

	for _, result := range results.Results {
		priority := "Medium"
		if result.Rank <= len(results.Results)/4 {
			priority = "High"
		} else if result.Rank >= 3*len(results.Results)/4 {
			priority = "Low"
		}

		// Map complexity score to story points
		storyPoints := result.SComplexity
		if storyPoints > 21 {
			storyPoints = 21 // Cap at 21 for Jira
		}

		description := result.Feature.Description
		if result.Feature.AcceptanceCriteria != "" {
			description += "\n\nAcceptance Criteria:\n" + result.Feature.AcceptanceCriteria
		}

		issue := domain.JiraIssue{
			Summary:     result.Feature.Title,
			Description: description,
			StoryPoints: storyPoints,
			Priority:    priority,
			CustomFields: domain.JiraCustomFields{
				FinalPriorityScore: result.FinalPriorityScore,
				ValueScore:         result.SValue,
				ComplexityScore:    result.SComplexity,
			},
		}

		issues = append(issues, issue)
	}

	return domain.JiraExport{Issues: issues}
}
