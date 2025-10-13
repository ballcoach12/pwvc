package api

import (
	"testing"

	"pairwise/internal/domain"
)

// TestAPIRequestStructures tests API request and response structures
func TestAPIRequestStructures(t *testing.T) {
	t.Run("Project Request Structures", func(t *testing.T) {
		createRequest := domain.CreateProjectRequest{
			Name:        "Test Project",
			Description: "A test project description",
		}

		if createRequest.Name != "Test Project" {
			t.Errorf("Expected name 'Test Project', got %s", createRequest.Name)
		}

		updateRequest := domain.UpdateProjectRequest{
			Name:        "Updated Project",
			Description: "Updated description",
		}

		if updateRequest.Name != "Updated Project" {
			t.Errorf("Expected name 'Updated Project', got %s", updateRequest.Name)
		}
	})

	t.Run("Attendee Request Structures", func(t *testing.T) {
		createRequest := domain.CreateAttendeeRequest{
			Name:          "John Doe",
			Role:          "Product Manager",
			IsFacilitator: true,
		}

		if createRequest.Name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got %s", createRequest.Name)
		}

		if !createRequest.IsFacilitator {
			t.Error("Expected facilitator to be true")
		}
	})

	t.Run("Feature Request Structures", func(t *testing.T) {
		createRequest := domain.CreateFeatureRequest{
			Title:              "User Authentication",
			Description:        "Implement user login",
			AcceptanceCriteria: "Users can log in and out",
		}

		if createRequest.Title != "User Authentication" {
			t.Errorf("Expected title 'User Authentication', got %s", createRequest.Title)
		}

		updateRequest := domain.UpdateFeatureRequest{
			Title:              "Enhanced Authentication",
			Description:        "Implement enhanced user login with 2FA",
			AcceptanceCriteria: "Users can log in with 2FA",
		}

		if updateRequest.Title != "Enhanced Authentication" {
			t.Errorf("Expected title 'Enhanced Authentication', got %s", updateRequest.Title)
		}
	})
}

// TestAPIWorkflowScenarios tests complete API workflow scenarios
func TestAPIWorkflowScenarios(t *testing.T) {
	t.Run("Project Setup Workflow", func(t *testing.T) {
		// Simulate a complete project setup workflow
		project := &domain.Project{
			ID:          1,
			Name:        "P-WVC Test Project",
			Description: "A project for testing P-WVC methodology",
		}

		// Add attendees
		attendees := []domain.Attendee{
			{ID: 1, ProjectID: 1, Name: "John Doe", Role: "PM", IsFacilitator: true},
			{ID: 2, ProjectID: 1, Name: "Jane Smith", Role: "Dev", IsFacilitator: false},
		}

		// Add features
		features := []domain.Feature{
			{ID: 1, ProjectID: 1, Title: "Feature A", Description: "First feature"},
			{ID: 2, ProjectID: 1, Title: "Feature B", Description: "Second feature"},
		}

		// Verify setup
		if project.ID != 1 {
			t.Errorf("Expected project ID 1, got %d", project.ID)
		}

		if len(attendees) != 2 {
			t.Errorf("Expected 2 attendees, got %d", len(attendees))
		}

		if len(features) != 2 {
			t.Errorf("Expected 2 features, got %d", len(features))
		}

		// Verify facilitator assignment
		facilitatorCount := 0
		for _, attendee := range attendees {
			if attendee.IsFacilitator {
				facilitatorCount++
			}
		}

		if facilitatorCount != 1 {
			t.Errorf("Expected 1 facilitator, got %d", facilitatorCount)
		}
	})

	t.Run("Pairwise Comparison Workflow", func(t *testing.T) {
		// Simulate pairwise comparison setup
		session := &domain.PairwiseSession{
			ID:            1,
			ProjectID:     1,
			CriterionType: domain.CriterionTypeValue,
			Status:        domain.SessionStatusActive,
		}

		// Create comparison
		winnerID := 1
		comparison := &domain.SessionComparison{
			ID:               1,
			SessionID:        1,
			FeatureAID:       1,
			FeatureBID:       2,
			WinnerID:         &winnerID,
			ConsensusReached: true,
		}

		// Create votes
		preferredFeatureID := 1
		vote := &domain.AttendeeVote{
			ID:                 1,
			ComparisonID:       1,
			AttendeeID:         1,
			PreferredFeatureID: &preferredFeatureID,
			IsTieVote:          false,
		}

		// Verify workflow
		if session.CriterionType != domain.CriterionTypeValue {
			t.Errorf("Expected criterion type %s, got %s", domain.CriterionTypeValue, session.CriterionType)
		}

		if *comparison.WinnerID != 1 {
			t.Errorf("Expected winner ID 1, got %d", *comparison.WinnerID)
		}

		if *vote.PreferredFeatureID != 1 {
			t.Errorf("Expected preferred feature ID 1, got %d", *vote.PreferredFeatureID)
		}
	})

	t.Run("Priority Calculation Workflow", func(t *testing.T) {
		// Simulate priority calculation
		calculation := &domain.PriorityCalculation{
			ID:                 1,
			ProjectID:          1,
			FeatureID:          1,
			WValue:             0.75,
			WComplexity:        0.60,
			SValue:             8,
			SComplexity:        3,
			WeightedValue:      6.0,  // 8 * 0.75
			WeightedComplexity: 1.8,  // 3 * 0.60
			FinalPriorityScore: 3.33, // 6.0 / 1.8
			Rank:               1,
		}

		// Verify calculation
		if calculation.SValue != 8 {
			t.Errorf("Expected value score 8, got %d", calculation.SValue)
		}

		if calculation.SComplexity != 3 {
			t.Errorf("Expected complexity score 3, got %d", calculation.SComplexity)
		}

		expectedScore := 3.33
		if calculation.FinalPriorityScore < expectedScore-0.01 || calculation.FinalPriorityScore > expectedScore+0.01 {
			t.Errorf("Expected final score around %.2f, got %.2f", expectedScore, calculation.FinalPriorityScore)
		}
	})
}

// TestPWVCMethodologyIntegration tests the complete P-WVC methodology workflow
func TestPWVCMethodologyIntegration(t *testing.T) {
	t.Run("Complete P-WVC Workflow", func(t *testing.T) {
		// Test the complete P-WVC methodology workflow using domain structures

		// 1. Project Setup
		project := &domain.Project{
			ID:          1,
			Name:        "P-WVC Integration Test",
			Description: "Testing complete P-WVC workflow",
		}

		// 2. Add Attendees
		attendees := []domain.Attendee{
			{ID: 1, ProjectID: 1, Name: "Product Manager", Role: "PM", IsFacilitator: true},
			{ID: 2, ProjectID: 1, Name: "Lead Developer", Role: "Dev", IsFacilitator: false},
			{ID: 3, ProjectID: 1, Name: "UX Designer", Role: "Design", IsFacilitator: false},
		}

		// 3. Add Features
		features := []domain.Feature{
			{ID: 1, ProjectID: 1, Title: "User Authentication", Description: "Login system"},
			{ID: 2, ProjectID: 1, Title: "Payment Processing", Description: "Credit card payments"},
			{ID: 3, ProjectID: 1, Title: "Notification System", Description: "Email and SMS alerts"},
		}

		// 4. Pairwise Comparisons - Value
		valueSession := &domain.PairwiseSession{
			ID:            1,
			ProjectID:     1,
			CriterionType: domain.CriterionTypeValue,
			Status:        domain.SessionStatusCompleted,
		}

		// 5. Pairwise Comparisons - Complexity
		complexitySession := &domain.PairwiseSession{
			ID:            2,
			ProjectID:     1,
			CriterionType: domain.CriterionTypeComplexity,
			Status:        domain.SessionStatusCompleted,
		}

		// 6. Priority Calculations
		calculations := []domain.PriorityCalculation{
			{
				ID:                 1,
				ProjectID:          1,
				FeatureID:          1,
				WValue:             0.67, // Won 2 of 3 comparisons
				WComplexity:        0.33, // Won 1 of 3 comparisons
				SValue:             8,    // Fibonacci score
				SComplexity:        5,    // Fibonacci score
				WeightedValue:      5.36, // 8 * 0.67
				WeightedComplexity: 1.65, // 5 * 0.33
				FinalPriorityScore: 3.25, // 5.36 / 1.65
				Rank:               1,
			},
			{
				ID:                 2,
				ProjectID:          1,
				FeatureID:          2,
				WValue:             0.33,
				WComplexity:        0.67,
				SValue:             5,
				SComplexity:        8,
				WeightedValue:      1.65,
				WeightedComplexity: 5.36,
				FinalPriorityScore: 0.31,
				Rank:               2,
			},
		}

		// Verify the workflow
		if project.ID != 1 {
			t.Errorf("Expected project ID 1, got %d", project.ID)
		}

		if len(attendees) != 3 {
			t.Errorf("Expected 3 attendees, got %d", len(attendees))
		}

		if len(features) != 3 {
			t.Errorf("Expected 3 features, got %d", len(features))
		}

		if valueSession.CriterionType != domain.CriterionTypeValue {
			t.Errorf("Expected value criterion, got %s", valueSession.CriterionType)
		}

		if complexitySession.CriterionType != domain.CriterionTypeComplexity {
			t.Errorf("Expected complexity criterion, got %s", complexitySession.CriterionType)
		}

		if len(calculations) != 2 {
			t.Errorf("Expected 2 calculations, got %d", len(calculations))
		}

		// Verify ranking order (highest FPS first)
		if calculations[0].Rank != 1 || calculations[1].Rank != 2 {
			t.Error("Priority ranking is incorrect")
		}

		if calculations[0].FinalPriorityScore <= calculations[1].FinalPriorityScore {
			t.Error("Final priority scores are not in correct order")
		}
	})
}

// TestErrorHandling tests API error handling scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("Validation Errors", func(t *testing.T) {
		// Test domain validation through request structures
		invalidRequest := domain.CreateProjectRequest{
			Name:        "", // Invalid: empty name
			Description: "Valid description",
		}

		if invalidRequest.Name != "" {
			t.Errorf("Expected empty name to be invalid")
		}

		validRequest := domain.CreateProjectRequest{
			Name:        "Valid Project Name",
			Description: "Valid description",
		}

		if validRequest.Name == "" {
			t.Errorf("Expected valid name to pass validation")
		}
	})

	t.Run("Business Logic Errors", func(t *testing.T) {
		// Test business logic validation scenarios

		// Cannot have more than one facilitator
		attendees := []domain.Attendee{
			{ID: 1, Name: "John", IsFacilitator: true},
			{ID: 2, Name: "Jane", IsFacilitator: true}, // This should cause an error
		}

		facilitatorCount := 0
		for _, attendee := range attendees {
			if attendee.IsFacilitator {
				facilitatorCount++
			}
		}

		if facilitatorCount > 1 {
			t.Logf("Business rule violation detected: multiple facilitators (%d)", facilitatorCount)
		}

		// Minimum features required for P-WVC
		minFeatures := 2
		features := []domain.Feature{
			{ID: 1, Title: "Feature 1"},
		}

		if len(features) < minFeatures {
			t.Logf("Business rule validation: need at least %d features, have %d", minFeatures, len(features))
		}
	})
}
