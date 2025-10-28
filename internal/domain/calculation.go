package domain

import (
	"errors"
	"fmt"
	"math"
)

// P-WVC Calculation Domain Models and Core Mathematical Functions

// PairwiseComparison represents a single pairwise comparison between two features
type PairwiseComparison struct {
	FeatureAID int                 `json:"feature_a_id"`
	FeatureBID int                 `json:"feature_b_id"`
	Criterion  ComparisonCriterion `json:"criterion"`
	Result     ComparisonResult    `json:"result"`
	UserID     int                 `json:"user_id,omitempty"`
}

// ComparisonCriterion represents the criteria being compared
type ComparisonCriterion string

const (
	CriterionValue      ComparisonCriterion = "value"
	CriterionComplexity ComparisonCriterion = "complexity"
)

// ComparisonResult represents the outcome of a pairwise comparison
type ComparisonResult string

const (
	ResultAWins ComparisonResult = "a_wins"
	ResultBWins ComparisonResult = "b_wins"
	ResultTie   ComparisonResult = "tie"
)

// FeatureScore represents the scores assigned to a feature
type FeatureScore struct {
	FeatureID          int     `json:"feature_id"`
	ValueScore         int     `json:"value_score"`
	ComplexityScore    int     `json:"complexity_score"`
	ValueWeight        float64 `json:"value_weight"`
	ComplexityWeight   float64 `json:"complexity_weight"`
	WeightedValue      float64 `json:"weighted_value"`
	WeightedComplexity float64 `json:"weighted_complexity"`
	FinalPriorityScore float64 `json:"final_priority_score"`
}

// WinCountResult represents the win-count calculation result
type WinCountResult struct {
	FeatureID        int     `json:"feature_id"`
	Wins             int     `json:"wins"`
	Losses           int     `json:"losses"`
	Ties             int     `json:"ties"`
	TotalComparisons int     `json:"total_comparisons"`
	WinCount         float64 `json:"win_count"`
}

// PWVCCalculationResult represents the complete P-WVC calculation result
type PWVCCalculationResult struct {
	ProjectID           int              `json:"project_id"`
	FeatureScores       []FeatureScore   `json:"feature_scores"`
	ValueWinCounts      []WinCountResult `json:"value_win_counts"`
	ComplexityWinCounts []WinCountResult `json:"complexity_win_counts"`
	RankedFeatures      []FeatureScore   `json:"ranked_features"`
}

// Fibonacci sequence validation constants
var ValidFibonacciScores = []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

// Calculation errors
var (
	ErrInvalidFibonacciScore = errors.New("invalid Fibonacci score")
	ErrZeroComplexity        = errors.New("weighted complexity cannot be calculated")
	ErrInsufficientData      = errors.New("insufficient comparison data")
	ErrInvalidComparison     = errors.New("invalid comparison data")
)

// CalculateWinCount calculates the win-count weight for a feature based on pairwise comparisons
// Formula: WFeature = (Total Wins + 0.5 × Total Ties) / (Total Comparisons)
func CalculateWinCount(featureID int, comparisons []PairwiseComparison, criterion ComparisonCriterion) (*WinCountResult, error) {
	if len(comparisons) == 0 {
		return nil, ErrInsufficientData
	}

	wins := 0
	losses := 0
	ties := 0

	// Count wins, losses, and ties for the specific feature and criterion
	for _, comp := range comparisons {
		if comp.Criterion != criterion {
			continue
		}

		if comp.FeatureAID == featureID {
			switch comp.Result {
			case ResultAWins:
				wins++
			case ResultBWins:
				losses++
			case ResultTie:
				ties++
			}
		} else if comp.FeatureBID == featureID {
			switch comp.Result {
			case ResultAWins:
				losses++
			case ResultBWins:
				wins++
			case ResultTie:
				ties++
			}
		}
	}

	totalComparisons := wins + losses + ties
	if totalComparisons == 0 {
		return nil, ErrInsufficientData
	}

	// Calculate win-count using the P-WVC formula
	winCount := (float64(wins) + 0.5*float64(ties)) / float64(totalComparisons)

	return &WinCountResult{
		FeatureID:        featureID,
		Wins:             wins,
		Losses:           losses,
		Ties:             ties,
		TotalComparisons: totalComparisons,
		WinCount:         winCount,
	}, nil
}

// ValidateFibonacciScore validates if a score is in the valid Fibonacci sequence
func ValidateFibonacciScore(score int) error {
	for _, validScore := range ValidFibonacciScores {
		if score == validScore {
			return nil
		}
	}
	return fmt.Errorf("%w: %d is not a valid Fibonacci score. Valid scores are: %v",
		ErrInvalidFibonacciScore, score, ValidFibonacciScores)
}

// CalculateWeightedScore calculates the weighted score (Score × Weight)
func CalculateWeightedScore(score int, weight float64) (float64, error) {
	if err := ValidateFibonacciScore(score); err != nil {
		return 0, err
	}

	if weight < 0 || weight > 1 {
		return 0, fmt.Errorf("weight must be between 0 and 1, got: %f", weight)
	}

	return float64(score) * weight, nil
}

// CalculateFinalPriorityScore calculates the Final Priority Score
// Formula: FPS = (SValue × WValue) / (SComplexity × WComplexity)
func CalculateFinalPriorityScore(valueScore int, valueWeight float64, complexityScore int, complexityWeight float64) (*FeatureScore, error) {
	// Validate Fibonacci scores
	if err := ValidateFibonacciScore(valueScore); err != nil {
		return nil, fmt.Errorf("value score validation failed: %w", err)
	}

	if err := ValidateFibonacciScore(complexityScore); err != nil {
		return nil, fmt.Errorf("complexity score validation failed: %w", err)
	}

	// Validate weights
	if valueWeight < 0 || valueWeight > 1 {
		return nil, fmt.Errorf("value weight must be between 0 and 1, got: %f", valueWeight)
	}

	if complexityWeight < 0 || complexityWeight > 1 {
		return nil, fmt.Errorf("complexity weight must be between 0 and 1, got: %f", complexityWeight)
	}

	// Calculate weighted scores
	weightedValue := float64(valueScore) * valueWeight
	weightedComplexity := float64(complexityScore) * complexityWeight

	// Handle division by zero case
	// In P-WVC methodology, if weighted complexity is 0 (no wins in complexity),
	// use a small epsilon to avoid division by zero while maintaining meaningful results
	var fps float64
	if weightedComplexity == 0 {
		// If complexity weight is 0, feature has maximum priority
		// (no complexity penalty). Use base complexity score as denominator.
		fps = weightedValue / float64(complexityScore)
	} else {
		fps = weightedValue / weightedComplexity
	}

	return &FeatureScore{
		ValueScore:         valueScore,
		ComplexityScore:    complexityScore,
		ValueWeight:        valueWeight,
		ComplexityWeight:   complexityWeight,
		WeightedValue:      weightedValue,
		WeightedComplexity: weightedComplexity,
		FinalPriorityScore: fps,
	}, nil
}

// GetFibonacciScoreIndex returns the index of a Fibonacci score in the sequence
func GetFibonacciScoreIndex(score int) (int, error) {
	for i, validScore := range ValidFibonacciScores {
		if score == validScore {
			return i, nil
		}
	}
	return -1, ErrInvalidFibonacciScore
}

// IsValidFibonacciScore checks if a score is valid without returning an error
func IsValidFibonacciScore(score int) bool {
	return ValidateFibonacciScore(score) == nil
}

// CalculateWinCountsForAllFeatures calculates win-counts for all features in a project
func CalculateWinCountsForAllFeatures(featureIDs []int, comparisons []PairwiseComparison, criterion ComparisonCriterion) ([]WinCountResult, error) {
	if len(featureIDs) == 0 {
		return nil, ErrInsufficientData
	}

	var results []WinCountResult

	for _, featureID := range featureIDs {
		winCount, err := CalculateWinCount(featureID, comparisons, criterion)
		if err != nil {
			// If a specific feature has no comparisons, set win count to 0
			if err == ErrInsufficientData {
				results = append(results, WinCountResult{
					FeatureID:        featureID,
					Wins:             0,
					Losses:           0,
					Ties:             0,
					TotalComparisons: 0,
					WinCount:         0.0,
				})
				continue
			}
			return nil, fmt.Errorf("failed to calculate win count for feature %d: %w", featureID, err)
		}
		results = append(results, *winCount)
	}

	return results, nil
}

// NormalizeWeights ensures that weights sum to 1.0 within acceptable floating point precision
func NormalizeWeights(weights []float64) []float64 {
	if len(weights) == 0 {
		return weights
	}

	sum := 0.0
	for _, w := range weights {
		sum += w
	}

	if sum == 0 {
		// If all weights are 0, distribute equally
		equal := 1.0 / float64(len(weights))
		normalized := make([]float64, len(weights))
		for i := range normalized {
			normalized[i] = equal
		}
		return normalized
	}

	// Normalize to sum to 1.0
	normalized := make([]float64, len(weights))
	for i, w := range weights {
		normalized[i] = w / sum
	}

	return normalized
}

// RoundToDecimalPlaces rounds a float64 to the specified number of decimal places
func RoundToDecimalPlaces(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
}
