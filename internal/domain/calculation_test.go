package domain

import (
	"math"
	"testing"
)

// Test Fibonacci validation
func TestValidateFibonacciScore(t *testing.T) {
	tests := []struct {
		name        string
		score       int
		expectError bool
	}{
		{"Valid score 1", 1, false},
		{"Valid score 2", 2, false},
		{"Valid score 3", 3, false},
		{"Valid score 5", 5, false},
		{"Valid score 8", 8, false},
		{"Valid score 13", 13, false},
		{"Valid score 21", 21, false},
		{"Valid score 34", 34, false},
		{"Valid score 55", 55, false},
		{"Valid score 89", 89, false},
		{"Invalid score 0", 0, true},
		{"Invalid score 4", 4, true},
		{"Invalid score 6", 6, true},
		{"Invalid score 7", 7, true},
		{"Invalid score 90", 90, true},
		{"Invalid score 144", 144, true},
		{"Invalid negative score", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFibonacciScore(tt.score)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for score %d, but got none", tt.score)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for score %d, but got: %v", tt.score, err)
			}
		})
	}
}

func TestIsValidFibonacciScore(t *testing.T) {
	validScores := []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	invalidScores := []int{0, 4, 6, 7, 9, 10, 90, 144, -1}

	for _, score := range validScores {
		if !IsValidFibonacciScore(score) {
			t.Errorf("Score %d should be valid", score)
		}
	}

	for _, score := range invalidScores {
		if IsValidFibonacciScore(score) {
			t.Errorf("Score %d should be invalid", score)
		}
	}
}

// Test win-count calculation
func TestCalculateWinCount(t *testing.T) {
	tests := []struct {
		name             string
		featureID        int
		comparisons      []PairwiseComparison
		criterion        ComparisonCriterion
		expectedWins     int
		expectedTies     int
		expectedWinCount float64
		expectError      bool
	}{
		{
			name:      "All wins",
			featureID: 1,
			comparisons: []PairwiseComparison{
				{FeatureAID: 1, FeatureBID: 2, Criterion: CriterionValue, Result: ResultAWins},
				{FeatureAID: 1, FeatureBID: 3, Criterion: CriterionValue, Result: ResultAWins},
			},
			criterion:        CriterionValue,
			expectedWins:     2,
			expectedTies:     0,
			expectedWinCount: 1.0,
			expectError:      false,
		},
		{
			name:      "All losses",
			featureID: 1,
			comparisons: []PairwiseComparison{
				{FeatureAID: 1, FeatureBID: 2, Criterion: CriterionValue, Result: ResultBWins},
				{FeatureAID: 3, FeatureBID: 1, Criterion: CriterionValue, Result: ResultAWins},
			},
			criterion:        CriterionValue,
			expectedWins:     0,
			expectedTies:     0,
			expectedWinCount: 0.0,
			expectError:      false,
		},
		{
			name:      "Mixed results with ties",
			featureID: 1,
			comparisons: []PairwiseComparison{
				{FeatureAID: 1, FeatureBID: 2, Criterion: CriterionValue, Result: ResultAWins},
				{FeatureAID: 1, FeatureBID: 3, Criterion: CriterionValue, Result: ResultTie},
				{FeatureAID: 4, FeatureBID: 1, Criterion: CriterionValue, Result: ResultBWins},
				{FeatureAID: 1, FeatureBID: 5, Criterion: CriterionValue, Result: ResultBWins},
			},
			criterion:        CriterionValue,
			expectedWins:     2, // 1 win + 1 loss (when feature 1 is B and result is B_wins)
			expectedTies:     1,
			expectedWinCount: 0.625, // (2 + 0.5*1) / 4
			expectError:      false,
		},
		{
			name:        "No comparisons",
			featureID:   1,
			comparisons: []PairwiseComparison{},
			criterion:   CriterionValue,
			expectError: true,
		},
		{
			name:      "No matching criterion",
			featureID: 1,
			comparisons: []PairwiseComparison{
				{FeatureAID: 1, FeatureBID: 2, Criterion: CriterionComplexity, Result: ResultAWins},
			},
			criterion:   CriterionValue,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateWinCount(tt.featureID, tt.comparisons, tt.criterion)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.FeatureID != tt.featureID {
				t.Errorf("Expected FeatureID %d, got %d", tt.featureID, result.FeatureID)
			}

			if result.Wins != tt.expectedWins {
				t.Errorf("Expected %d wins, got %d", tt.expectedWins, result.Wins)
			}

			if result.Ties != tt.expectedTies {
				t.Errorf("Expected %d ties, got %d", tt.expectedTies, result.Ties)
			}

			if math.Abs(result.WinCount-tt.expectedWinCount) > 0.0001 {
				t.Errorf("Expected win count %.4f, got %.4f", tt.expectedWinCount, result.WinCount)
			}
		})
	}
}

// Test Final Priority Score calculation
func TestCalculateFinalPriorityScore(t *testing.T) {
	tests := []struct {
		name             string
		valueScore       int
		valueWeight      float64
		complexityScore  int
		complexityWeight float64
		expectedFPS      float64
		expectError      bool
		errorContains    string
	}{
		{
			name:             "Basic calculation",
			valueScore:       8,
			valueWeight:      0.6,
			complexityScore:  3,
			complexityWeight: 0.4,
			expectedFPS:      4.0, // (8 * 0.6) / (3 * 0.4) = 4.8 / 1.2 = 4.0
		},
		{
			name:             "High value, low complexity",
			valueScore:       21,
			valueWeight:      0.8,
			complexityScore:  2,
			complexityWeight: 0.2,
			expectedFPS:      42.0, // (21 * 0.8) / (2 * 0.2) = 16.8 / 0.4 = 42.0
		},
		{
			name:             "Equal weights",
			valueScore:       5,
			valueWeight:      0.5,
			complexityScore:  5,
			complexityWeight: 0.5,
			expectedFPS:      1.0, // (5 * 0.5) / (5 * 0.5) = 2.5 / 2.5 = 1.0
		},
		{
			name:             "Invalid value score",
			valueScore:       4, // Not a Fibonacci number
			valueWeight:      0.5,
			complexityScore:  3,
			complexityWeight: 0.5,
			expectError:      true,
			errorContains:    "value score validation failed",
		},
		{
			name:             "Invalid complexity score",
			valueScore:       5,
			valueWeight:      0.5,
			complexityScore:  4, // Not a Fibonacci number
			complexityWeight: 0.5,
			expectError:      true,
			errorContains:    "complexity score validation failed",
		},
		{
			name:             "Zero complexity weight",
			valueScore:       5,
			valueWeight:      0.5,
			complexityScore:  3,
			complexityWeight: 0.0,
			expectError:      false,
			expectedFPS:      0.8333, // (5 * 0.5) / 3 = 2.5 / 3 = 0.8333...
		},
		{
			name:             "Invalid value weight (negative)",
			valueScore:       5,
			valueWeight:      -0.1,
			complexityScore:  3,
			complexityWeight: 0.5,
			expectError:      true,
			errorContains:    "value weight must be between 0 and 1",
		},
		{
			name:             "Invalid complexity weight (greater than 1)",
			valueScore:       5,
			valueWeight:      0.5,
			complexityScore:  3,
			complexityWeight: 1.5,
			expectError:      true,
			errorContains:    "complexity weight must be between 0 and 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateFinalPriorityScore(
				tt.valueScore, tt.valueWeight,
				tt.complexityScore, tt.complexityWeight,
			)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if math.Abs(result.FinalPriorityScore-tt.expectedFPS) > 0.0001 {
				t.Errorf("Expected FPS %.4f, got %.4f", tt.expectedFPS, result.FinalPriorityScore)
			}

			// Verify component calculations
			expectedWeightedValue := float64(tt.valueScore) * tt.valueWeight
			expectedWeightedComplexity := float64(tt.complexityScore) * tt.complexityWeight

			if math.Abs(result.WeightedValue-expectedWeightedValue) > 0.0001 {
				t.Errorf("Expected weighted value %.4f, got %.4f", expectedWeightedValue, result.WeightedValue)
			}

			if math.Abs(result.WeightedComplexity-expectedWeightedComplexity) > 0.0001 {
				t.Errorf("Expected weighted complexity %.4f, got %.4f", expectedWeightedComplexity, result.WeightedComplexity)
			}
		})
	}
}

// Test weighted score calculation
func TestCalculateWeightedScore(t *testing.T) {
	tests := []struct {
		name           string
		score          int
		weight         float64
		expectedResult float64
		expectError    bool
	}{
		{"Valid calculation", 8, 0.75, 6.0, false},
		{"Zero weight", 5, 0.0, 0.0, false},
		{"Full weight", 13, 1.0, 13.0, false},
		{"Invalid Fibonacci score", 4, 0.5, 0.0, true},
		{"Negative weight", 5, -0.1, 0.0, true},
		{"Weight greater than 1", 5, 1.5, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateWeightedScore(tt.score, tt.weight)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if math.Abs(result-tt.expectedResult) > 0.0001 {
				t.Errorf("Expected %.4f, got %.4f", tt.expectedResult, result)
			}
		})
	}
}

// Test weight normalization
func TestNormalizeWeights(t *testing.T) {
	tests := []struct {
		name     string
		weights  []float64
		expected []float64
	}{
		{
			name:     "Already normalized",
			weights:  []float64{0.3, 0.7},
			expected: []float64{0.3, 0.7},
		},
		{
			name:     "Need normalization",
			weights:  []float64{0.6, 1.4},
			expected: []float64{0.3, 0.7},
		},
		{
			name:     "All zeros",
			weights:  []float64{0.0, 0.0, 0.0},
			expected: []float64{0.3333333333333333, 0.3333333333333333, 0.3333333333333333},
		},
		{
			name:     "Empty slice",
			weights:  []float64{},
			expected: []float64{},
		},
		{
			name:     "Single weight",
			weights:  []float64{5.0},
			expected: []float64{1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeWeights(tt.weights)

			if len(result) != len(tt.expected) {
				t.Fatalf("Expected length %d, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if math.Abs(result[i]-expected) > 0.0001 {
					t.Errorf("Index %d: expected %.4f, got %.4f", i, expected, result[i])
				}
			}
		})
	}
}

// Test rounding function
func TestRoundToDecimalPlaces(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		places   int
		expected float64
	}{
		{"Round to 2 places", 3.14159, 2, 3.14},
		{"Round to 4 places", 2.718281828, 4, 2.7183},
		{"No rounding needed", 5.0, 2, 5.0},
		{"Round up", 1.999, 2, 2.0},
		{"Negative number", -3.14159, 2, -3.14},
		{"Zero places", 3.14159, 0, 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToDecimalPlaces(tt.value, tt.places)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("Expected %.4f, got %.4f", tt.expected, result)
			}
		})
	}
}

// Test Fibonacci score index lookup
func TestGetFibonacciScoreIndex(t *testing.T) {
	tests := []struct {
		name          string
		score         int
		expectedIndex int
		expectError   bool
	}{
		{"First score", 1, 0, false},
		{"Second score", 2, 1, false},
		{"Third score", 3, 2, false},
		{"Fifth score", 8, 4, false},
		{"Last score", 89, 9, false},
		{"Invalid score", 4, -1, true},
		{"Zero score", 0, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := GetFibonacciScoreIndex(tt.score)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if index != -1 {
					t.Errorf("Expected index -1 for error case, got %d", index)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if index != tt.expectedIndex {
				t.Errorf("Expected index %d, got %d", tt.expectedIndex, index)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
