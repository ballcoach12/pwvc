package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

// ResultsService handles P-WVC results calculation and management
type ResultsService struct {
	priorityRepo  repository.PriorityRepository
	featureRepo   repository.FeatureRepository
	pairwiseRepo  repository.PairwiseRepository
	consensusRepo repository.ConsensusRepository
}

// NewResultsService creates a new results service
func NewResultsService(
	priorityRepo repository.PriorityRepository,
	featureRepo repository.FeatureRepository,
	pairwiseRepo repository.PairwiseRepository,
	consensusRepo repository.ConsensusRepository,
) *ResultsService {
	return &ResultsService{
		priorityRepo:  priorityRepo,
		featureRepo:   featureRepo,
		pairwiseRepo:  pairwiseRepo,
		consensusRepo: consensusRepo,
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

	// 5. Sort by Final Priority Score with deterministic tie-break
	sort.Slice(calculations, func(i, j int) bool {
		calc1, calc2 := calculations[i], calculations[j]

		// Primary: FPS descending
		if calc1.FinalPriorityScore != calc2.FinalPriorityScore {
			return calc1.FinalPriorityScore > calc2.FinalPriorityScore
		}

		// Tie-break 1: SValue descending
		if calc1.SValue != calc2.SValue {
			return calc1.SValue > calc2.SValue
		}

		// Tie-break 2: SComplexity ascending (lower complexity preferred)
		if calc1.SComplexity != calc2.SComplexity {
			return calc1.SComplexity < calc2.SComplexity
		}

		// Tie-break 3: Feature name alphabetical (need to get feature names)
		// For now, use feature ID as proxy (should fetch feature names for proper implementation)
		return calc1.FeatureID < calc2.FeatureID
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
	// Convert string to domain type
	var criterion domain.CriterionType
	switch criterionType {
	case "value":
		criterion = domain.CriterionTypeValue
	case "complexity":
		criterion = domain.CriterionTypeComplexity
	default:
		return nil, fmt.Errorf("invalid criterion type: %s", criterionType)
	}

	// Get real win-count weights from repository aggregation
	weights, err := s.pairwiseRepo.GetWinCountWeights(projectID, criterion)
	if err != nil {
		return nil, err
	}

	// If no weights found (no completed sessions), return default weights
	if len(weights) == 0 {
		features, err := s.featureRepo.GetByProjectID(projectID)
		if err != nil {
			return nil, err
		}

		weights = make(map[int]float64)
		for _, feature := range features {
			weights[feature.ID] = 0.5 // Default neutral weight
		}
	}

	return weights, nil
}

// getFibonacciScores retrieves consensus Fibonacci scores for features
func (s *ResultsService) getFibonacciScores(projectID int, criterionType string) (map[int]int, error) {
	// Get consensus scores from repository
	consensusScores, err := s.consensusRepo.GetConsensusScores(projectID)
	if err != nil {
		return nil, err
	}

	scores := make(map[int]int)

	// Extract the appropriate score type based on criterionType
	for featureID, consensus := range consensusScores {
		switch criterionType {
		case "value":
			scores[featureID] = consensus.SValue
		case "complexity":
			scores[featureID] = consensus.SComplexity
		default:
			return nil, fmt.Errorf("invalid criterion type: %s", criterionType)
		}
	}

	// If no consensus scores found, check for features without scores and use default
	if len(scores) == 0 {
		features, err := s.featureRepo.GetByProjectID(projectID)
		if err != nil {
			return nil, err
		}

		for _, feature := range features {
			scores[feature.ID] = 3 // Default Fibonacci value
		}
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

// exportToCSV converts results to CSV format with deterministic ordering and locale-safe formatting (T050)
func (s *ResultsService) exportToCSV(results *domain.ProjectResults) [][]string {
	// Create CSV header with consistent field names
	csv := [][]string{
		{"rank", "feature_id", "feature_title", "description", "final_priority_score", "s_value", "s_complexity", "w_value", "w_complexity", "calculation_timestamp"},
	}

	// Ensure deterministic ordering by sorting results by rank, then by feature ID as tiebreaker
	sortedResults := make([]domain.PriorityResult, len(results.Results))
	copy(sortedResults, results.Results)

	// Sort by rank (ascending), then by feature ID for consistency
	sort.Slice(sortedResults, func(i, j int) bool {
		if sortedResults[i].Rank == sortedResults[j].Rank {
			return sortedResults[i].Feature.ID < sortedResults[j].Feature.ID
		}
		return sortedResults[i].Rank < sortedResults[j].Rank
	})

	// Format timestamp in ISO 8601 format (locale-independent)
	timestamp := results.CalculatedAt.UTC().Format("2006-01-02T15:04:05Z")

	for _, result := range sortedResults {
		// Use locale-safe formatting with consistent precision
		row := []string{
			strconv.Itoa(result.Rank),                     // Deterministic integer formatting
			strconv.Itoa(result.Feature.ID),               // Feature ID for traceability
			s.escapeCsvField(result.Feature.Title),        // Properly escape CSV fields
			s.escapeCsvField(result.Feature.Description),  // Properly escape CSV fields
			s.formatDecimal(result.FinalPriorityScore, 6), // Locale-safe decimal formatting
			strconv.Itoa(result.SValue),                   // Integer scores don't need special formatting
			strconv.Itoa(result.SComplexity),              // Integer scores don't need special formatting
			s.formatDecimal(result.WValue, 6),             // Locale-safe decimal formatting
			s.formatDecimal(result.WComplexity, 6),        // Locale-safe decimal formatting
			timestamp,                                     // ISO 8601 timestamp
		}
		csv = append(csv, row)
	}

	return csv
}

// formatDecimal formats decimal numbers using locale-safe formatting (T050)
func (s *ResultsService) formatDecimal(value float64, precision int) string {
	// Use strconv.FormatFloat for locale-independent formatting
	// Always use decimal point (.) regardless of system locale
	return strconv.FormatFloat(value, 'f', precision, 64)
}

// escapeCsvField properly escapes CSV fields containing special characters (T050)
func (s *ResultsService) escapeCsvField(field string) string {
	// CSV fields containing commas, quotes, or newlines need to be quoted
	if strings.ContainsAny(field, ",\"\n\r") {
		// Escape internal quotes by doubling them, then wrap in quotes
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return field
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
