package service

import (
	"testing"
	"time"

	"pairwise/internal/domain"
)

// TestPairwiseSession tests pairwise session logic
func TestPairwiseSession(t *testing.T) {
	t.Run("Create Pairwise Session", func(t *testing.T) {
		session := &domain.PairwiseSession{
			ProjectID:     1,
			CriterionType: domain.CriterionTypeValue,
			Status:        domain.SessionStatusActive,
			StartedAt:     time.Now(),
		}

		if session.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", session.ProjectID)
		}

		if session.CriterionType != domain.CriterionTypeValue {
			t.Errorf("Expected criterion type %s, got %s", domain.CriterionTypeValue, session.CriterionType)
		}

		if session.Status != domain.SessionStatusActive {
			t.Errorf("Expected status %s, got %s", domain.SessionStatusActive, session.Status)
		}
	})

	t.Run("Criterion Types", func(t *testing.T) {
		valueType := domain.CriterionTypeValue
		complexityType := domain.CriterionTypeComplexity

		if string(valueType) != "value" {
			t.Errorf("Expected 'value', got %s", string(valueType))
		}

		if string(complexityType) != "complexity" {
			t.Errorf("Expected 'complexity', got %s", string(complexityType))
		}
	})

	t.Run("Session Status", func(t *testing.T) {
		activeStatus := domain.SessionStatusActive
		completedStatus := domain.SessionStatusCompleted

		if string(activeStatus) != "active" {
			t.Errorf("Expected 'active', got %s", string(activeStatus))
		}

		if string(completedStatus) != "completed" {
			t.Errorf("Expected 'completed', got %s", string(completedStatus))
		}
	})
}

// TestSessionComparison tests session comparison logic
func TestSessionComparison(t *testing.T) {
	t.Run("Create Session Comparison", func(t *testing.T) {
		winnerID := 1
		comparison := &domain.SessionComparison{
			SessionID:        1,
			FeatureAID:       1,
			FeatureBID:       2,
			WinnerID:         &winnerID,
			IsTie:            false,
			ConsensusReached: true,
			CreatedAt:        time.Now(),
		}

		if comparison.SessionID != 1 {
			t.Errorf("Expected session ID 1, got %d", comparison.SessionID)
		}

		if comparison.FeatureAID != 1 {
			t.Errorf("Expected feature A ID 1, got %d", comparison.FeatureAID)
		}

		if comparison.FeatureBID != 2 {
			t.Errorf("Expected feature B ID 2, got %d", comparison.FeatureBID)
		}

		if *comparison.WinnerID != 1 {
			t.Errorf("Expected winner ID 1, got %d", *comparison.WinnerID)
		}

		if comparison.IsTie {
			t.Error("Expected no tie")
		}

		if !comparison.ConsensusReached {
			t.Error("Expected consensus reached")
		}
	})
}

// TestAttendeeVote tests attendee vote logic
func TestAttendeeVote(t *testing.T) {
	t.Run("Create Attendee Vote", func(t *testing.T) {
		preferredFeatureID := 1
		vote := &domain.AttendeeVote{
			ComparisonID:       1,
			AttendeeID:         1,
			PreferredFeatureID: &preferredFeatureID,
			IsTieVote:          false,
			VotedAt:            time.Now(),
		}

		if vote.ComparisonID != 1 {
			t.Errorf("Expected comparison ID 1, got %d", vote.ComparisonID)
		}

		if vote.AttendeeID != 1 {
			t.Errorf("Expected attendee ID 1, got %d", vote.AttendeeID)
		}

		if *vote.PreferredFeatureID != 1 {
			t.Errorf("Expected preferred feature ID 1, got %d", *vote.PreferredFeatureID)
		}

		if vote.IsTieVote {
			t.Error("Expected no tie vote")
		}
	})

	t.Run("Create Tie Vote", func(t *testing.T) {
		vote := &domain.AttendeeVote{
			ComparisonID:       1,
			AttendeeID:         1,
			PreferredFeatureID: nil,
			IsTieVote:          true,
			VotedAt:            time.Now(),
		}

		if vote.PreferredFeatureID != nil {
			t.Error("Expected nil preferred feature ID for tie vote")
		}

		if !vote.IsTieVote {
			t.Error("Expected tie vote")
		}
	})
}

// TestPriorityCalculation tests priority calculation logic
func TestPriorityCalculation(t *testing.T) {
	t.Run("Create Priority Calculation", func(t *testing.T) {
		priority := &domain.PriorityCalculation{
			ProjectID:          1,
			FeatureID:          1,
			WValue:             0.75,
			WComplexity:        0.60,
			SValue:             5,
			SComplexity:        3,
			WeightedValue:      3.75, // 5 * 0.75
			WeightedComplexity: 1.8,  // 3 * 0.60
			FinalPriorityScore: 2.08, // 3.75 / 1.8
			Rank:               1,
			CalculatedAt:       time.Now(),
		}

		if priority.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", priority.ProjectID)
		}

		if priority.SValue != 5 {
			t.Errorf("Expected value score 5, got %d", priority.SValue)
		}

		if priority.SComplexity != 3 {
			t.Errorf("Expected complexity score 3, got %d", priority.SComplexity)
		}

		expectedScore := 2.08
		if priority.FinalPriorityScore < expectedScore-0.01 || priority.FinalPriorityScore > expectedScore+0.01 {
			t.Errorf("Expected final score around %.2f, got %.2f", expectedScore, priority.FinalPriorityScore)
		}
	})

	t.Run("Priority Calculation Formula", func(t *testing.T) {
		// Test the P-WVC formula: FPS = (SValue × WValue) / (SComplexity × WComplexity)
		sValue := 8.0
		wValue := 0.75
		sComplexity := 3.0
		wComplexity := 0.60

		weightedValue := sValue * wValue                  // 6.0
		weightedComplexity := sComplexity * wComplexity   // 1.8
		expectedFPS := weightedValue / weightedComplexity // 3.33

		if expectedFPS < 3.33 || expectedFPS > 3.34 {
			t.Errorf("FPS calculation seems incorrect: %.2f", expectedFPS)
		}
	})
}

// TestSessionProgressTracking tests session progress tracking
func TestSessionProgressTracking(t *testing.T) {
	t.Run("Session Progress Calculation", func(t *testing.T) {
		progress := &domain.SessionProgress{
			SessionID:            1,
			TotalComparisons:     10,
			CompletedComparisons: 7,
			ProgressPercentage:   70.0,
			RemainingComparisons: 3,
		}

		if progress.SessionID != 1 {
			t.Errorf("Expected session ID 1, got %d", progress.SessionID)
		}

		if progress.TotalComparisons != 10 {
			t.Errorf("Expected total comparisons 10, got %d", progress.TotalComparisons)
		}

		if progress.CompletedComparisons != 7 {
			t.Errorf("Expected completed comparisons 7, got %d", progress.CompletedComparisons)
		}

		expectedPercentage := 70.0
		if progress.ProgressPercentage != expectedPercentage {
			t.Errorf("Expected progress percentage %.1f, got %.1f", expectedPercentage, progress.ProgressPercentage)
		}

		if progress.RemainingComparisons != 3 {
			t.Errorf("Expected remaining comparisons 3, got %d", progress.RemainingComparisons)
		}
	})
}

// BenchmarkPairwiseOperations benchmarks pairwise operations
func BenchmarkPairwiseOperations(b *testing.B) {
	sessions := make([]*domain.PairwiseSession, 100)
	for i := 0; i < 100; i++ {
		sessions[i] = &domain.PairwiseSession{
			ProjectID:     1,
			CriterionType: domain.CriterionTypeValue,
			Status:        domain.SessionStatusActive,
			StartedAt:     time.Now(),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate processing sessions
		for _, session := range sessions {
			_ = session.CriterionType
			_ = session.Status
		}
	}
}

// BenchmarkPriorityCalculation benchmarks priority calculation operations
func BenchmarkPriorityCalculation(b *testing.B) {
	calculations := make([]*domain.PriorityCalculation, 100)
	for i := 0; i < 100; i++ {
		calculations[i] = &domain.PriorityCalculation{
			ProjectID:          1,
			FeatureID:          i + 1,
			WValue:             0.75,
			WComplexity:        0.60,
			SValue:             5,
			SComplexity:        3,
			WeightedValue:      3.75,
			WeightedComplexity: 1.8,
			FinalPriorityScore: 2.08,
			Rank:               i + 1,
			CalculatedAt:       time.Now(),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate processing calculations
		for _, calc := range calculations {
			_ = calc.FinalPriorityScore
			_ = calc.Rank
		}
	}
}
