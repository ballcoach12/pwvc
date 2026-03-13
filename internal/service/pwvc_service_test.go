package service

import (
	"math"
	"testing"

	"pairwise/internal/domain"
)

func TestPWVCService_ValidateFibonacciScores(t *testing.T) {
	service := NewPWVCService()

	tests := []struct {
		name        string
		scores      map[int]domain.FeatureScore
		expectError bool
	}{
		{
			name: "Valid scores",
			scores: map[int]domain.FeatureScore{
				1: {ValueScore: 8, ComplexityScore: 3},
				2: {ValueScore: 5, ComplexityScore: 13},
			},
			expectError: false,
		},
		{
			name: "Invalid value score",
			scores: map[int]domain.FeatureScore{
				1: {ValueScore: 4, ComplexityScore: 3}, // 4 is not Fibonacci
			},
			expectError: true,
		},
		{
			name: "Invalid complexity score",
			scores: map[int]domain.FeatureScore{
				1: {ValueScore: 8, ComplexityScore: 6}, // 6 is not Fibonacci
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateFibonacciScores(tt.scores)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestPWVCService_CalculateWinCountWeight(t *testing.T) {
	service := NewPWVCService()

	comparisons := []domain.PairwiseComparison{
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultTie},
	}

	result, err := service.CalculateWinCountWeight(1, comparisons, domain.CriterionValue)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.FeatureID != 1 {
		t.Errorf("Expected FeatureID 1, got %d", result.FeatureID)
	}

	expectedWinCount := 0.75 // (1 + 0.5*1) / 2
	if result.WinCount != expectedWinCount {
		t.Errorf("Expected win count %.2f, got %.2f", expectedWinCount, result.WinCount)
	}
}

func TestPWVCService_CalculateSingleFPS(t *testing.T) {
	service := NewPWVCService()

	result, err := service.CalculateSingleFPS(8, 0.6, 3, 0.4)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedFPS := 4.0 // (8 * 0.6) / (3 * 0.4) = 4.8 / 1.2 = 4.0
	if math.Abs(result.FinalPriorityScore-expectedFPS) > 0.0001 {
		t.Errorf("Expected FPS %.4f, got %.4f", expectedFPS, result.FinalPriorityScore)
	}
}

func TestPWVCService_SimulatePWVCScenario(t *testing.T) {
	service := NewPWVCService()

	scenarios := []domain.FeatureScore{
		{FeatureID: 1, ValueScore: 13, ValueWeight: 0.8, ComplexityScore: 3, ComplexityWeight: 0.2},
		{FeatureID: 2, ValueScore: 5, ValueWeight: 0.6, ComplexityScore: 8, ComplexityWeight: 0.4},
	}

	results, err := service.SimulatePWVCScenario(scenarios)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Results should be sorted by FPS (highest first)
	if results[0].FinalPriorityScore < results[1].FinalPriorityScore {
		t.Error("Results should be sorted by FPS in descending order")
	}
}

func TestPWVCService_AnalyzeComparisonCompleteness(t *testing.T) {
	service := NewPWVCService()

	featureIDs := []int{1, 2, 3} // 3 features = 3 required comparisons per criterion

	comparisons := []domain.PairwiseComparison{
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultBWins},
		// Missing: 2 vs 3 for value
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionComplexity, Result: domain.ResultTie},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		{FeatureAID: 2, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultBWins},
	}

	report, err := service.AnalyzeComparisonCompleteness(featureIDs, comparisons)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if report.TotalFeaturesCount != 3 {
		t.Errorf("Expected 3 features, got %d", report.TotalFeaturesCount)
	}

	if report.RequiredComparisonsPerCriterion != 3 {
		t.Errorf("Expected 3 required comparisons, got %d", report.RequiredComparisonsPerCriterion)
	}

	if report.ValueComparisonsComplete != 2 {
		t.Errorf("Expected 2 value comparisons complete, got %d", report.ValueComparisonsComplete)
	}

	if report.ComplexityComparisonsComplete != 3 {
		t.Errorf("Expected 3 complexity comparisons complete, got %d", report.ComplexityComparisonsComplete)
	}

	if report.IsValueComplete {
		t.Error("Value comparisons should not be complete")
	}

	if !report.IsComplexityComplete {
		t.Error("Complexity comparisons should be complete")
	}

	if report.IsFullyComplete {
		t.Error("Overall should not be complete")
	}
}

func TestPWVCService_CalculateProjectPWVC(t *testing.T) {
	service := NewPWVCService()

	featureIDs := []int{1, 2}

	fibonacciScores := map[int]domain.FeatureScore{
		1: {ValueScore: 8, ComplexityScore: 3},
		2: {ValueScore: 5, ComplexityScore: 8},
	}

	comparisons := []domain.PairwiseComparison{
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionComplexity, Result: domain.ResultBWins},
	}

	result, err := service.CalculateProjectPWVC(featureIDs, fibonacciScores, comparisons)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.FeatureScores) != 2 {
		t.Fatalf("Expected 2 feature scores, got %d", len(result.FeatureScores))
	}

	if len(result.ValueWinCounts) != 2 {
		t.Fatalf("Expected 2 value win counts, got %d", len(result.ValueWinCounts))
	}

	if len(result.ComplexityWinCounts) != 2 {
		t.Fatalf("Expected 2 complexity win counts, got %d", len(result.ComplexityWinCounts))
	}

	if len(result.RankedFeatures) != 2 {
		t.Fatalf("Expected 2 ranked features, got %d", len(result.RankedFeatures))
	}

	// Verify that features are ranked by FPS
	if len(result.RankedFeatures) > 1 {
		for i := 1; i < len(result.RankedFeatures); i++ {
			if result.RankedFeatures[i-1].FinalPriorityScore < result.RankedFeatures[i].FinalPriorityScore {
				t.Error("Ranked features should be sorted by FPS in descending order")
			}
		}
	}

	// Verify win counts are correct
	// Feature 1 should have win count 1.0 for value (won against feature 2)
	// Feature 2 should have win count 0.0 for value (lost against feature 1)
	var feature1ValueWinCount, feature2ValueWinCount float64
	for _, wc := range result.ValueWinCounts {
		if wc.FeatureID == 1 {
			feature1ValueWinCount = wc.WinCount
		} else if wc.FeatureID == 2 {
			feature2ValueWinCount = wc.WinCount
		}
	}

	if feature1ValueWinCount != 1.0 {
		t.Errorf("Expected feature 1 value win count 1.0, got %.2f", feature1ValueWinCount)
	}

	if feature2ValueWinCount != 0.0 {
		t.Errorf("Expected feature 2 value win count 0.0, got %.2f", feature2ValueWinCount)
	}
}

func TestPWVCService_GetValidFibonacciScores(t *testing.T) {
	service := NewPWVCService()

	scores := service.GetValidFibonacciScores()
	expected := []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

	if len(scores) != len(expected) {
		t.Fatalf("Expected %d scores, got %d", len(expected), len(scores))
	}

	for i, score := range scores {
		if score != expected[i] {
			t.Errorf("Index %d: expected %d, got %d", i, expected[i], score)
		}
	}
}

// Test error cases
func TestPWVCService_ErrorCases(t *testing.T) {
	service := NewPWVCService()

	t.Run("Empty feature IDs", func(t *testing.T) {
		_, err := service.CalculateProjectPWVC(
			[]int{},
			map[int]domain.FeatureScore{},
			[]domain.PairwiseComparison{},
		)
		if err == nil {
			t.Error("Expected error for empty feature IDs")
		}
	})

	t.Run("Missing Fibonacci scores", func(t *testing.T) {
		_, err := service.CalculateProjectPWVC(
			[]int{1, 2},
			map[int]domain.FeatureScore{},
			[]domain.PairwiseComparison{},
		)
		if err == nil {
			t.Error("Expected error for missing Fibonacci scores")
		}
	})

	t.Run("Invalid Fibonacci score in calculation", func(t *testing.T) {
		_, err := service.CalculateSingleFPS(4, 0.5, 3, 0.5) // 4 is invalid
		if err == nil {
			t.Error("Expected error for invalid Fibonacci score")
		}
	})

	t.Run("Insufficient features for comparison analysis", func(t *testing.T) {
		_, err := service.AnalyzeComparisonCompleteness(
			[]int{1}, // Only 1 feature
			[]domain.PairwiseComparison{},
		)
		if err == nil {
			t.Error("Expected error for insufficient features")
		}
	})
}
