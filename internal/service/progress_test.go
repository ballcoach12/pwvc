package service

import (
	"testing"

	"pairwise/internal/domain"
)

// TestProjectProgress tests the project progress domain logic
func TestProjectProgress(t *testing.T) {
	t.Run("New Project Progress", func(t *testing.T) {
		progress := &domain.ProjectProgress{
			ProjectID:    1,
			CurrentPhase: string(domain.PhaseSetup),
		}

		if progress.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", progress.ProjectID)
		}

		if progress.CurrentPhase != string(domain.PhaseSetup) {
			t.Errorf("Expected phase %s, got %s", domain.PhaseSetup, progress.CurrentPhase)
		}
	})

	t.Run("Get Next Phase", func(t *testing.T) {
		progress := &domain.ProjectProgress{
			CurrentPhase: string(domain.PhaseSetup),
		}

		nextPhase := progress.GetNextPhase()
		if nextPhase != domain.PhaseAttendees {
			t.Errorf("Expected next phase %s, got %s", domain.PhaseAttendees, nextPhase)
		}
	})

	t.Run("Can Progress To Phase", func(t *testing.T) {
		progress := &domain.ProjectProgress{
			CurrentPhase: string(domain.PhaseSetup),
		}

		// Should be able to progress to setup (initial phase)
		if !progress.CanProgressTo(domain.PhaseSetup) {
			t.Error("Should be able to progress to setup phase")
		}

		// Should not be able to progress to attendees without setup complete
		if progress.CanProgressTo(domain.PhaseAttendees) {
			t.Error("Should not be able to progress to attendees without setup complete")
		}

		// Complete setup
		progress.SetupCompleted = true

		// Now should be able to progress to attendees
		if !progress.CanProgressTo(domain.PhaseAttendees) {
			t.Error("Should be able to progress to attendees after setup complete")
		}
	})

	t.Run("Progress Through All Phases", func(t *testing.T) {
		progress := &domain.ProjectProgress{
			ProjectID:    1,
			CurrentPhase: string(domain.PhaseSetup),
		}

		// Complete setup
		progress.SetupCompleted = true
		progress.CurrentPhase = string(progress.GetNextPhase())
		if progress.CurrentPhase != string(domain.PhaseAttendees) {
			t.Errorf("Expected phase %s, got %s", domain.PhaseAttendees, progress.CurrentPhase)
		}

		// Complete attendees
		progress.AttendeesAdded = true
		progress.CurrentPhase = string(progress.GetNextPhase())
		if progress.CurrentPhase != string(domain.PhaseFeatures) {
			t.Errorf("Expected phase %s, got %s", domain.PhaseFeatures, progress.CurrentPhase)
		}

		// Complete features
		progress.FeaturesAdded = true
		progress.CurrentPhase = string(progress.GetNextPhase())
		if progress.CurrentPhase != string(domain.PhasePairwiseValue) {
			t.Errorf("Expected phase %s, got %s", domain.PhasePairwiseValue, progress.CurrentPhase)
		}
	})
}

// TestWorkflowPhases tests workflow phase constants and logic
func TestWorkflowPhases(t *testing.T) {
	t.Run("Phase Constants", func(t *testing.T) {
		phases := []domain.WorkflowPhase{
			domain.PhaseSetup,
			domain.PhaseAttendees,
			domain.PhaseFeatures,
			domain.PhasePairwiseValue,
			domain.PhasePairwiseComplexity,
			domain.PhaseFibonacciValue,
			domain.PhaseFibonacciComplexity,
			domain.PhaseResults,
		}

		expectedPhases := []string{
			"setup",
			"attendees",
			"features",
			"pairwise_value",
			"pairwise_complexity",
			"fibonacci_value",
			"fibonacci_complexity",
			"results",
		}

		for i, phase := range phases {
			if string(phase) != expectedPhases[i] {
				t.Errorf("Expected phase %s, got %s", expectedPhases[i], string(phase))
			}
		}
	})
}

// BenchmarkProjectProgress benchmarks project progress operations
func BenchmarkProjectProgress(b *testing.B) {
	progress := &domain.ProjectProgress{
		ProjectID:    1,
		CurrentPhase: string(domain.PhaseSetup),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = progress.GetNextPhase()
		_ = progress.CanProgressTo(domain.PhaseAttendees)
	}
}
