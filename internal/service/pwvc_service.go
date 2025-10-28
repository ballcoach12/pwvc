package service

import (
	"fmt"
	"sort"

	"pairwise/internal/domain"
)

// PWVCService handles P-WVC methodology calculations and business logic
type PWVCService struct {
	// Add repository dependencies if needed for data persistence
}

// NewPWVCService creates a new P-WVC service instance
func NewPWVCService() *PWVCService {
	return &PWVCService{}
}

// CalculateProjectPWVC performs complete P-WVC calculation for a project
func (s *PWVCService) CalculateProjectPWVC(
	featureIDs []int,
	fibonacciScores map[int]domain.FeatureScore, // featureID -> scores
	pairwiseComparisons []domain.PairwiseComparison,
) (*domain.PWVCCalculationResult, error) {

	if len(featureIDs) == 0 {
		return nil, domain.NewAPIError(400, "No features provided for calculation")
	}

	if len(fibonacciScores) == 0 {
		return nil, domain.NewAPIError(400, "No Fibonacci scores provided")
	}

	// Validate all Fibonacci scores first
	if err := s.validateAllFibonacciScores(fibonacciScores); err != nil {
		return nil, err
	}

	// Calculate win-counts for value criterion
	valueWinCounts, err := domain.CalculateWinCountsForAllFeatures(
		featureIDs,
		pairwiseComparisons,
		domain.CriterionValue,
	)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to calculate value win-counts", err.Error())
	}

	// Calculate win-counts for complexity criterion
	complexityWinCounts, err := domain.CalculateWinCountsForAllFeatures(
		featureIDs,
		pairwiseComparisons,
		domain.CriterionComplexity,
	)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to calculate complexity win-counts", err.Error())
	}

	// Build final feature scores with FPS calculations
	featureScores, err := s.calculateFinalScores(
		featureIDs,
		fibonacciScores,
		valueWinCounts,
		complexityWinCounts,
	)
	if err != nil {
		return nil, err
	}

	// Rank features by Final Priority Score (highest to lowest)
	rankedFeatures := s.rankFeaturesByFPS(featureScores)

	return &domain.PWVCCalculationResult{
		FeatureScores:       featureScores,
		ValueWinCounts:      valueWinCounts,
		ComplexityWinCounts: complexityWinCounts,
		RankedFeatures:      rankedFeatures,
	}, nil
}

// ValidateFibonacciScores validates multiple Fibonacci scores at once
func (s *PWVCService) ValidateFibonacciScores(scores map[int]domain.FeatureScore) error {
	return s.validateAllFibonacciScores(scores)
}

// CalculateWinCountWeight calculates win-count for a single feature
func (s *PWVCService) CalculateWinCountWeight(
	featureID int,
	comparisons []domain.PairwiseComparison,
	criterion domain.ComparisonCriterion,
) (*domain.WinCountResult, error) {

	result, err := domain.CalculateWinCount(featureID, comparisons, criterion)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to calculate win-count", err.Error())
	}

	return result, nil
}

// CalculateSingleFPS calculates Final Priority Score for a single feature
func (s *PWVCService) CalculateSingleFPS(
	valueScore int,
	valueWeight float64,
	complexityScore int,
	complexityWeight float64,
) (*domain.FeatureScore, error) {

	result, err := domain.CalculateFinalPriorityScore(
		valueScore, valueWeight, complexityScore, complexityWeight,
	)
	if err != nil {
		return nil, domain.NewAPIError(400, "Failed to calculate FPS", err.Error())
	}

	return result, nil
}

// GetValidFibonacciScores returns the list of valid Fibonacci scores
func (s *PWVCService) GetValidFibonacciScores() []int {
	return domain.ValidFibonacciScores
}

// SimulatePWVCScenario allows testing different scoring scenarios
func (s *PWVCService) SimulatePWVCScenario(
	scenarios []domain.FeatureScore,
) ([]domain.FeatureScore, error) {

	var results []domain.FeatureScore

	for i, scenario := range scenarios {
		// Validate and calculate FPS for each scenario
		fps, err := domain.CalculateFinalPriorityScore(
			scenario.ValueScore,
			scenario.ValueWeight,
			scenario.ComplexityScore,
			scenario.ComplexityWeight,
		)
		if err != nil {
			return nil, domain.NewAPIError(400,
				fmt.Sprintf("Invalid scenario %d: %s", i+1, err.Error()))
		}

		// Copy the feature ID and calculated values
		fps.FeatureID = scenario.FeatureID
		results = append(results, *fps)
	}

	// Sort by FPS (highest to lowest)
	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalPriorityScore > results[j].FinalPriorityScore
	})

	return results, nil
}

// AnalyzeComparisonCompleteness checks if all necessary pairwise comparisons exist
func (s *PWVCService) AnalyzeComparisonCompleteness(
	featureIDs []int,
	comparisons []domain.PairwiseComparison,
) (*ComparisonCompletenessReport, error) {

	if len(featureIDs) < 2 {
		return nil, domain.NewAPIError(400, "At least 2 features required for pairwise comparison")
	}

	totalRequired := len(featureIDs) * (len(featureIDs) - 1) / 2 // n*(n-1)/2

	// Track existing comparisons by criterion
	valueComparisons := make(map[string]bool)
	complexityComparisons := make(map[string]bool)

	for _, comp := range comparisons {
		// Create a normalized key (smaller ID first)
		var key string
		if comp.FeatureAID < comp.FeatureBID {
			key = fmt.Sprintf("%d-%d", comp.FeatureAID, comp.FeatureBID)
		} else {
			key = fmt.Sprintf("%d-%d", comp.FeatureBID, comp.FeatureAID)
		}

		if comp.Criterion == domain.CriterionValue {
			valueComparisons[key] = true
		} else if comp.Criterion == domain.CriterionComplexity {
			complexityComparisons[key] = true
		}
	}

	return &ComparisonCompletenessReport{
		TotalFeaturesCount:              len(featureIDs),
		RequiredComparisonsPerCriterion: totalRequired,
		ValueComparisonsComplete:        len(valueComparisons),
		ComplexityComparisonsComplete:   len(complexityComparisons),
		ValueCompletenessPercent:        float64(len(valueComparisons)) / float64(totalRequired) * 100,
		ComplexityCompletenessPercent:   float64(len(complexityComparisons)) / float64(totalRequired) * 100,
		IsValueComplete:                 len(valueComparisons) == totalRequired,
		IsComplexityComplete:            len(complexityComparisons) == totalRequired,
		IsFullyComplete:                 len(valueComparisons) == totalRequired && len(complexityComparisons) == totalRequired,
	}, nil
}

// ComparisonCompletenessReport provides analysis of comparison completeness
type ComparisonCompletenessReport struct {
	TotalFeaturesCount              int     `json:"total_features_count"`
	RequiredComparisonsPerCriterion int     `json:"required_comparisons_per_criterion"`
	ValueComparisonsComplete        int     `json:"value_comparisons_complete"`
	ComplexityComparisonsComplete   int     `json:"complexity_comparisons_complete"`
	ValueCompletenessPercent        float64 `json:"value_completeness_percent"`
	ComplexityCompletenessPercent   float64 `json:"complexity_completeness_percent"`
	IsValueComplete                 bool    `json:"is_value_complete"`
	IsComplexityComplete            bool    `json:"is_complexity_complete"`
	IsFullyComplete                 bool    `json:"is_fully_complete"`
}

// Private helper methods

func (s *PWVCService) validateAllFibonacciScores(scores map[int]domain.FeatureScore) error {
	for featureID, score := range scores {
		if err := domain.ValidateFibonacciScore(score.ValueScore); err != nil {
			return domain.NewAPIError(400,
				fmt.Sprintf("Invalid value score for feature %d: %s", featureID, err.Error()))
		}

		if err := domain.ValidateFibonacciScore(score.ComplexityScore); err != nil {
			return domain.NewAPIError(400,
				fmt.Sprintf("Invalid complexity score for feature %d: %s", featureID, err.Error()))
		}
	}

	return nil
}

func (s *PWVCService) calculateFinalScores(
	featureIDs []int,
	fibonacciScores map[int]domain.FeatureScore,
	valueWinCounts []domain.WinCountResult,
	complexityWinCounts []domain.WinCountResult,
) ([]domain.FeatureScore, error) {

	// Create lookup maps for win-counts
	valueWeights := make(map[int]float64)
	complexityWeights := make(map[int]float64)

	for _, vwc := range valueWinCounts {
		valueWeights[vwc.FeatureID] = vwc.WinCount
	}

	for _, cwc := range complexityWinCounts {
		complexityWeights[cwc.FeatureID] = cwc.WinCount
	}

	var featureScores []domain.FeatureScore

	for _, featureID := range featureIDs {
		fibScore, exists := fibonacciScores[featureID]
		if !exists {
			return nil, domain.NewAPIError(400,
				fmt.Sprintf("Missing Fibonacci scores for feature %d", featureID))
		}

		valueWeight := valueWeights[featureID]
		complexityWeight := complexityWeights[featureID]

		// Calculate FPS
		fps, err := domain.CalculateFinalPriorityScore(
			fibScore.ValueScore,
			valueWeight,
			fibScore.ComplexityScore,
			complexityWeight,
		)
		if err != nil {
			return nil, domain.NewAPIError(500,
				fmt.Sprintf("Failed to calculate FPS for feature %d: %s", featureID, err.Error()))
		}

		// Set the feature ID
		fps.FeatureID = featureID

		// Round values for clean presentation
		fps.FinalPriorityScore = domain.RoundToDecimalPlaces(fps.FinalPriorityScore, 4)
		fps.WeightedValue = domain.RoundToDecimalPlaces(fps.WeightedValue, 4)
		fps.WeightedComplexity = domain.RoundToDecimalPlaces(fps.WeightedComplexity, 4)
		fps.ValueWeight = domain.RoundToDecimalPlaces(fps.ValueWeight, 4)
		fps.ComplexityWeight = domain.RoundToDecimalPlaces(fps.ComplexityWeight, 4)

		featureScores = append(featureScores, *fps)
	}

	return featureScores, nil
}

func (s *PWVCService) rankFeaturesByFPS(features []domain.FeatureScore) []domain.FeatureScore {
	// Create a copy to avoid modifying the original slice
	ranked := make([]domain.FeatureScore, len(features))
	copy(ranked, features)

	// Sort by Final Priority Score (highest to lowest)
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].FinalPriorityScore > ranked[j].FinalPriorityScore
	})

	return ranked
}
