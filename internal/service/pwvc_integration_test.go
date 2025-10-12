package service

import (
	"testing"

	"pwvc/internal/domain"
)

// Integration test with a realistic P-WVC scenario using a sample project
func TestPWVCService_IntegrationScenario(t *testing.T) {
	service := NewPWVCService()

	// Sample project: Website Redesign with 4 features
	featureIDs := []int{1, 2, 3, 4}

	// Feature details (for context):
	// 1: User Authentication - High value, medium complexity
	// 2: Dashboard Analytics - Medium value, high complexity
	// 3: Search Functionality - Medium value, low complexity
	// 4: Mobile Responsive - High value, low complexity

	// Step 1: Assign Fibonacci scores
	fibonacciScores := map[int]domain.FeatureScore{
		1: {ValueScore: 13, ComplexityScore: 5}, // User Auth: High value, medium complexity
		2: {ValueScore: 8, ComplexityScore: 13}, // Dashboard: Medium value, high complexity
		3: {ValueScore: 5, ComplexityScore: 3},  // Search: Medium value, low complexity
		4: {ValueScore: 21, ComplexityScore: 2}, // Mobile: Very high value, very low complexity
	}

	// Step 2: Pairwise comparisons (complete set for 4 features = 6 comparisons per criterion)
	pairwiseComparisons := []domain.PairwiseComparison{
		// VALUE COMPARISONS
		// Feature 1 (Auth) vs Feature 2 (Dashboard): Auth wins (security is critical)
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		// Feature 1 (Auth) vs Feature 3 (Search): Auth wins (security over convenience)
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		// Feature 1 (Auth) vs Feature 4 (Mobile): Mobile wins (broader reach)
		{FeatureAID: 1, FeatureBID: 4, Criterion: domain.CriterionValue, Result: domain.ResultBWins},
		// Feature 2 (Dashboard) vs Feature 3 (Search): Dashboard wins (business insights)
		{FeatureAID: 2, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		// Feature 2 (Dashboard) vs Feature 4 (Mobile): Mobile wins (accessibility)
		{FeatureAID: 2, FeatureBID: 4, Criterion: domain.CriterionValue, Result: domain.ResultBWins},
		// Feature 3 (Search) vs Feature 4 (Mobile): Mobile wins (fundamental usability)
		{FeatureAID: 3, FeatureBID: 4, Criterion: domain.CriterionValue, Result: domain.ResultBWins},

		// COMPLEXITY COMPARISONS
		// Feature 1 (Auth) vs Feature 2 (Dashboard): Dashboard wins (more complex)
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionComplexity, Result: domain.ResultBWins},
		// Feature 1 (Auth) vs Feature 3 (Search): Auth wins (more complex)
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		// Feature 1 (Auth) vs Feature 4 (Mobile): Auth wins (more complex)
		{FeatureAID: 1, FeatureBID: 4, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		// Feature 2 (Dashboard) vs Feature 3 (Search): Dashboard wins (more complex)
		{FeatureAID: 2, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		// Feature 2 (Dashboard) vs Feature 4 (Mobile): Dashboard wins (more complex)
		{FeatureAID: 2, FeatureBID: 4, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		// Feature 3 (Search) vs Feature 4 (Mobile): Search wins (more complex)
		{FeatureAID: 3, FeatureBID: 4, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
	}

	// Step 3: Calculate complete P-WVC results
	result, err := service.CalculateProjectPWVC(featureIDs, fibonacciScores, pairwiseComparisons)
	if err != nil {
		t.Fatalf("P-WVC calculation failed: %v", err)
	}

	// Verify completeness
	if len(result.FeatureScores) != 4 {
		t.Fatalf("Expected 4 feature scores, got %d", len(result.FeatureScores))
	}

	if len(result.RankedFeatures) != 4 {
		t.Fatalf("Expected 4 ranked features, got %d", len(result.RankedFeatures))
	}

	// Verify win counts make sense
	expectedValueWinCounts := map[int]float64{
		1: 0.6667, // 2 wins out of 3 comparisons
		2: 0.3333, // 1 win out of 3 comparisons
		3: 0.0,    // 0 wins out of 3 comparisons
		4: 1.0,    // 3 wins out of 3 comparisons
	}

	expectedComplexityWinCounts := map[int]float64{
		1: 0.6667, // 2 wins out of 3 comparisons
		2: 1.0,    // 3 wins out of 3 comparisons
		3: 0.3333, // 1 win out of 3 comparisons
		4: 0.0,    // 0 wins out of 3 comparisons
	}

	// Check value win counts
	for _, wc := range result.ValueWinCounts {
		expected := expectedValueWinCounts[wc.FeatureID]
		if abs(wc.WinCount-expected) > 0.01 {
			t.Errorf("Feature %d value win count: expected %.4f, got %.4f",
				wc.FeatureID, expected, wc.WinCount)
		}
	}

	// Check complexity win counts
	for _, cwc := range result.ComplexityWinCounts {
		expected := expectedComplexityWinCounts[cwc.FeatureID]
		if abs(cwc.WinCount-expected) > 0.01 {
			t.Errorf("Feature %d complexity win count: expected %.4f, got %.4f",
				cwc.FeatureID, expected, cwc.WinCount)
		}
	}

	// Verify Final Priority Scores are calculated correctly
	// Feature 4 (Mobile) should have highest FPS due to high value and low complexity
	highestFPS := result.RankedFeatures[0]
	if highestFPS.FeatureID != 4 {
		t.Errorf("Expected feature 4 to have highest FPS, got feature %d", highestFPS.FeatureID)
	}

	// Verify that FPS values are reasonable (can be 0 if value weight is 0)
	for _, fs := range result.FeatureScores {
		if fs.FinalPriorityScore < 0 {
			t.Errorf("Feature %d has negative FPS: %.4f", fs.FeatureID, fs.FinalPriorityScore)
		}

		// If value weight is 0, FPS should be 0
		if fs.ValueWeight == 0 && fs.FinalPriorityScore != 0 {
			t.Errorf("Feature %d with value weight 0 should have FPS 0, got %.4f",
				fs.FeatureID, fs.FinalPriorityScore)
		}

		// Verify component calculations
		expectedWeightedValue := float64(fibonacciScores[fs.FeatureID].ValueScore) * fs.ValueWeight
		expectedWeightedComplexity := float64(fibonacciScores[fs.FeatureID].ComplexityScore) * fs.ComplexityWeight

		if abs(fs.WeightedValue-expectedWeightedValue) > 0.01 {
			t.Errorf("Feature %d weighted value mismatch: expected %.4f, got %.4f",
				fs.FeatureID, expectedWeightedValue, fs.WeightedValue)
		}

		if abs(fs.WeightedComplexity-expectedWeightedComplexity) > 0.01 {
			t.Errorf("Feature %d weighted complexity mismatch: expected %.4f, got %.4f",
				fs.FeatureID, expectedWeightedComplexity, fs.WeightedComplexity)
		}
	}

	// Print results for manual verification (optional)
	t.Logf("P-WVC Integration Test Results:")
	t.Logf("==============================")
	for i, feature := range result.RankedFeatures {
		fibScore := fibonacciScores[feature.FeatureID]
		t.Logf("Rank %d: Feature %d (Value: %d, Complexity: %d) -> FPS: %.4f",
			i+1, feature.FeatureID, fibScore.ValueScore, fibScore.ComplexityScore, feature.FinalPriorityScore)
	}
}

// Test edge case: All ties scenario
func TestPWVCService_AllTiesScenario(t *testing.T) {
	service := NewPWVCService()

	featureIDs := []int{1, 2, 3}

	fibonacciScores := map[int]domain.FeatureScore{
		1: {ValueScore: 5, ComplexityScore: 5},
		2: {ValueScore: 5, ComplexityScore: 5},
		3: {ValueScore: 5, ComplexityScore: 5},
	}

	// All comparisons result in ties
	pairwiseComparisons := []domain.PairwiseComparison{
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultTie},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultTie},
		{FeatureAID: 2, FeatureBID: 3, Criterion: domain.CriterionValue, Result: domain.ResultTie},
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionComplexity, Result: domain.ResultTie},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultTie},
		{FeatureAID: 2, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultTie},
	}

	result, err := service.CalculateProjectPWVC(featureIDs, fibonacciScores, pairwiseComparisons)
	if err != nil {
		t.Fatalf("P-WVC calculation failed: %v", err)
	}

	// All features should have win count of 0.5 (since all comparisons are ties)
	for _, wc := range result.ValueWinCounts {
		if wc.WinCount != 0.5 {
			t.Errorf("Expected win count 0.5 for all ties, got %.4f for feature %d",
				wc.WinCount, wc.FeatureID)
		}
	}

	for _, cwc := range result.ComplexityWinCounts {
		if cwc.WinCount != 0.5 {
			t.Errorf("Expected win count 0.5 for all ties, got %.4f for feature %d",
				cwc.WinCount, cwc.FeatureID)
		}
	}

	// All features should have the same FPS since they have identical scores and weights
	firstFPS := result.FeatureScores[0].FinalPriorityScore
	for _, fs := range result.FeatureScores[1:] {
		if abs(fs.FinalPriorityScore-firstFPS) > 0.001 {
			t.Errorf("Expected all features to have same FPS in all-ties scenario, got %.4f vs %.4f",
				firstFPS, fs.FinalPriorityScore)
		}
	}
}

// Test scenario with missing comparisons (incomplete data)
func TestPWVCService_IncompleteComparisonsScenario(t *testing.T) {
	service := NewPWVCService()

	featureIDs := []int{1, 2, 3}

	fibonacciScores := map[int]domain.FeatureScore{
		1: {ValueScore: 8, ComplexityScore: 3},
		2: {ValueScore: 5, ComplexityScore: 8},
		3: {ValueScore: 13, ComplexityScore: 2},
	}

	// Only partial comparisons (missing some pairs)
	pairwiseComparisons := []domain.PairwiseComparison{
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionValue, Result: domain.ResultAWins},
		// Missing: 1 vs 3, 2 vs 3 for value
		{FeatureAID: 1, FeatureBID: 2, Criterion: domain.CriterionComplexity, Result: domain.ResultBWins},
		{FeatureAID: 1, FeatureBID: 3, Criterion: domain.CriterionComplexity, Result: domain.ResultAWins},
		// Missing: 2 vs 3 for complexity
	}

	result, err := service.CalculateProjectPWVC(featureIDs, fibonacciScores, pairwiseComparisons)
	if err != nil {
		t.Fatalf("P-WVC calculation should handle incomplete data: %v", err)
	}

	// Features with no comparisons should have win count 0
	feature3ValueWinCount := 0.0
	for _, wc := range result.ValueWinCounts {
		if wc.FeatureID == 3 {
			feature3ValueWinCount = wc.WinCount
			break
		}
	}

	if feature3ValueWinCount != 0.0 {
		t.Errorf("Feature 3 should have win count 0.0 for value (no comparisons), got %.4f",
			feature3ValueWinCount)
	}

	// Verify completeness analysis
	report, err := service.AnalyzeComparisonCompleteness(featureIDs, pairwiseComparisons)
	if err != nil {
		t.Fatalf("Completeness analysis failed: %v", err)
	}

	if report.IsFullyComplete {
		t.Error("Report should indicate incomplete comparisons")
	}

	if report.ValueCompletenessPercent >= 100.0 {
		t.Errorf("Value completeness should be less than 100%%, got %.2f%%",
			report.ValueCompletenessPercent)
	}
}

// Test large-scale scenario with many features
func TestPWVCService_LargeScaleScenario(t *testing.T) {
	service := NewPWVCService()

	// 6 features = 15 comparisons per criterion (n*(n-1)/2)
	featureIDs := []int{1, 2, 3, 4, 5, 6}

	fibonacciScores := map[int]domain.FeatureScore{
		1: {ValueScore: 21, ComplexityScore: 3}, // High value, low complexity
		2: {ValueScore: 13, ComplexityScore: 8}, // High value, medium complexity
		3: {ValueScore: 8, ComplexityScore: 13}, // Medium value, high complexity
		4: {ValueScore: 5, ComplexityScore: 5},  // Medium value, medium complexity
		5: {ValueScore: 3, ComplexityScore: 21}, // Low value, very high complexity
		6: {ValueScore: 34, ComplexityScore: 2}, // Very high value, very low complexity
	}

	// Generate comprehensive comparisons (simplified for testing)
	// In reality, these would be determined by stakeholder input
	var pairwiseComparisons []domain.PairwiseComparison

	// Generate all possible pairs
	for i := 0; i < len(featureIDs); i++ {
		for j := i + 1; j < len(featureIDs); j++ {
			featureA := featureIDs[i]
			featureB := featureIDs[j]

			// Value comparisons: higher value score wins
			valueResult := domain.ResultAWins
			if fibonacciScores[featureA].ValueScore < fibonacciScores[featureB].ValueScore {
				valueResult = domain.ResultBWins
			} else if fibonacciScores[featureA].ValueScore == fibonacciScores[featureB].ValueScore {
				valueResult = domain.ResultTie
			}

			// Complexity comparisons: higher complexity score "wins" (is more complex)
			complexityResult := domain.ResultAWins
			if fibonacciScores[featureA].ComplexityScore < fibonacciScores[featureB].ComplexityScore {
				complexityResult = domain.ResultBWins
			} else if fibonacciScores[featureA].ComplexityScore == fibonacciScores[featureB].ComplexityScore {
				complexityResult = domain.ResultTie
			}

			pairwiseComparisons = append(pairwiseComparisons,
				domain.PairwiseComparison{
					FeatureAID: featureA,
					FeatureBID: featureB,
					Criterion:  domain.CriterionValue,
					Result:     valueResult,
				},
				domain.PairwiseComparison{
					FeatureAID: featureA,
					FeatureBID: featureB,
					Criterion:  domain.CriterionComplexity,
					Result:     complexityResult,
				},
			)
		}
	}

	result, err := service.CalculateProjectPWVC(featureIDs, fibonacciScores, pairwiseComparisons)
	if err != nil {
		t.Fatalf("Large-scale P-WVC calculation failed: %v", err)
	}

	// Verify completeness
	report, err := service.AnalyzeComparisonCompleteness(featureIDs, pairwiseComparisons)
	if err != nil {
		t.Fatalf("Completeness analysis failed: %v", err)
	}

	if !report.IsFullyComplete {
		t.Error("All comparisons should be complete in large-scale scenario")
	}

	if report.RequiredComparisonsPerCriterion != 15 {
		t.Errorf("Expected 15 required comparisons per criterion, got %d",
			report.RequiredComparisonsPerCriterion)
	}

	// Verify that the highest ranking feature has the highest FPS
	// Based on calculations, this should be determined by actual FPS values
	topFeature := result.RankedFeatures[0]
	if topFeature.FinalPriorityScore <= 0 {
		t.Errorf("Top ranked feature should have positive FPS, got %.4f", topFeature.FinalPriorityScore)
	}

	// Feature 5 should rank lowest (low value, very high complexity)
	bottomFeature := result.RankedFeatures[len(result.RankedFeatures)-1]
	if bottomFeature.FeatureID != 5 {
		t.Errorf("Expected feature 5 to rank lowest, got feature %d", bottomFeature.FeatureID)
	}

	// Verify ranking is in descending order of FPS
	for i := 1; i < len(result.RankedFeatures); i++ {
		if result.RankedFeatures[i-1].FinalPriorityScore < result.RankedFeatures[i].FinalPriorityScore {
			t.Errorf("Features not properly ranked: feature %d (%.4f) should rank higher than feature %d (%.4f)",
				result.RankedFeatures[i].FeatureID, result.RankedFeatures[i].FinalPriorityScore,
				result.RankedFeatures[i-1].FeatureID, result.RankedFeatures[i-1].FinalPriorityScore)
		}
	}

	t.Logf("Large-scale P-WVC Test Results (%d features):", len(featureIDs))
	t.Logf("=============================================")
	for i, feature := range result.RankedFeatures {
		fibScore := fibonacciScores[feature.FeatureID]
		t.Logf("Rank %d: Feature %d (V:%d, C:%d, VW:%.3f, CW:%.3f) -> FPS: %.4f",
			i+1, feature.FeatureID,
			fibScore.ValueScore, fibScore.ComplexityScore,
			feature.ValueWeight, feature.ComplexityWeight,
			feature.FinalPriorityScore)
	}
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
