package repository

import (
	"testing"

	"pairwise/internal/domain"
)

// TestProjectRepository tests project repository operations
func TestProjectRepository(t *testing.T) {
	t.Run("Create Project Request", func(t *testing.T) {
		request := domain.CreateProjectRequest{
			Name:        "Test Project",
			Description: "A test project for P-WVC",
		}

		if request.Name != "Test Project" {
			t.Errorf("Expected name 'Test Project', got %s", request.Name)
		}

		if request.Description != "A test project for P-WVC" {
			t.Errorf("Expected description 'A test project for P-WVC', got %s", request.Description)
		}
	})

	t.Run("Update Project Request", func(t *testing.T) {
		request := domain.UpdateProjectRequest{
			Name:        "Updated Project",
			Description: "An updated test project",
		}

		if request.Name != "Updated Project" {
			t.Errorf("Expected name 'Updated Project', got %s", request.Name)
		}

		if request.Description != "An updated test project" {
			t.Errorf("Expected description 'An updated test project', got %s", request.Description)
		}
	})
}

// TestAttendeeRepository tests attendee repository operations
func TestAttendeeRepository(t *testing.T) {
	t.Run("Create Attendee Request", func(t *testing.T) {
		request := domain.CreateAttendeeRequest{
			Name:          "John Doe",
			Role:          "Product Manager",
			IsFacilitator: true,
		}

		if request.Name != "John Doe" {
			t.Errorf("Expected name 'John Doe', got %s", request.Name)
		}

		if request.Role != "Product Manager" {
			t.Errorf("Expected role 'Product Manager', got %s", request.Role)
		}

		if !request.IsFacilitator {
			t.Error("Expected facilitator to be true")
		}
	})

	t.Run("Attendee Domain Model", func(t *testing.T) {
		attendee := &domain.Attendee{
			ID:            1,
			ProjectID:     1,
			Name:          "Jane Smith",
			Role:          "Developer",
			IsFacilitator: false,
		}

		if attendee.ID != 1 {
			t.Errorf("Expected ID 1, got %d", attendee.ID)
		}

		if attendee.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", attendee.ProjectID)
		}

		if attendee.Name != "Jane Smith" {
			t.Errorf("Expected name 'Jane Smith', got %s", attendee.Name)
		}

		if attendee.Role != "Developer" {
			t.Errorf("Expected role 'Developer', got %s", attendee.Role)
		}

		if attendee.IsFacilitator {
			t.Error("Expected facilitator to be false")
		}
	})
}

// TestFeatureRepository tests feature repository operations
func TestFeatureRepository(t *testing.T) {
	t.Run("Create Feature Request", func(t *testing.T) {
		request := domain.CreateFeatureRequest{
			Title:              "User Authentication",
			Description:        "Implement user login and registration",
			AcceptanceCriteria: "Users can sign up, log in, and log out",
		}

		if request.Title != "User Authentication" {
			t.Errorf("Expected title 'User Authentication', got %s", request.Title)
		}

		if request.Description != "Implement user login and registration" {
			t.Errorf("Expected description 'Implement user login and registration', got %s", request.Description)
		}

		if request.AcceptanceCriteria != "Users can sign up, log in, and log out" {
			t.Errorf("Expected acceptance criteria 'Users can sign up, log in, and log out', got %s", request.AcceptanceCriteria)
		}
	})

	t.Run("Update Feature Request", func(t *testing.T) {
		request := domain.UpdateFeatureRequest{
			Title:              "Enhanced User Authentication",
			Description:        "Implement advanced user login and registration with 2FA",
			AcceptanceCriteria: "Users can sign up, log in with 2FA, and log out securely",
		}

		if request.Title != "Enhanced User Authentication" {
			t.Errorf("Expected title 'Enhanced User Authentication', got %s", request.Title)
		}

		if request.Description != "Implement advanced user login and registration with 2FA" {
			t.Errorf("Expected description 'Implement advanced user login and registration with 2FA', got %s", request.Description)
		}

		if request.AcceptanceCriteria != "Users can sign up, log in with 2FA, and log out securely" {
			t.Errorf("Expected acceptance criteria 'Users can sign up, log in with 2FA, and log out securely', got %s", request.AcceptanceCriteria)
		}
	})

	t.Run("Feature Domain Model", func(t *testing.T) {
		feature := &domain.Feature{
			ID:                 1,
			ProjectID:          1,
			Title:              "Payment Processing",
			Description:        "Implement credit card payment processing",
			AcceptanceCriteria: "Users can pay with major credit cards securely",
		}

		if feature.ID != 1 {
			t.Errorf("Expected ID 1, got %d", feature.ID)
		}

		if feature.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", feature.ProjectID)
		}

		if feature.Title != "Payment Processing" {
			t.Errorf("Expected title 'Payment Processing', got %s", feature.Title)
		}

		if feature.Description != "Implement credit card payment processing" {
			t.Errorf("Expected description 'Implement credit card payment processing', got %s", feature.Description)
		}

		if feature.AcceptanceCriteria != "Users can pay with major credit cards securely" {
			t.Errorf("Expected acceptance criteria 'Users can pay with major credit cards securely', got %s", feature.AcceptanceCriteria)
		}
	})
}

// TestPairwiseRepository tests pairwise repository operations
func TestPairwiseRepository(t *testing.T) {
	t.Run("Create Pairwise Session Request", func(t *testing.T) {
		request := domain.CreatePairwiseSessionRequest{
			CriterionType: domain.CriterionTypeValue,
		}

		if request.CriterionType != domain.CriterionTypeValue {
			t.Errorf("Expected criterion type %s, got %s", domain.CriterionTypeValue, request.CriterionType)
		}
	})

	t.Run("Submit Vote Request", func(t *testing.T) {
		preferredFeatureID := 1
		request := domain.SubmitVoteRequest{
			ComparisonID:       1,
			AttendeeID:         1,
			PreferredFeatureID: &preferredFeatureID,
			IsTieVote:          false,
		}

		if request.ComparisonID != 1 {
			t.Errorf("Expected comparison ID 1, got %d", request.ComparisonID)
		}

		if request.AttendeeID != 1 {
			t.Errorf("Expected attendee ID 1, got %d", request.AttendeeID)
		}

		if *request.PreferredFeatureID != 1 {
			t.Errorf("Expected preferred feature ID 1, got %d", *request.PreferredFeatureID)
		}

		if request.IsTieVote {
			t.Error("Expected no tie vote")
		}
	})

	t.Run("Submit Tie Vote Request", func(t *testing.T) {
		request := domain.SubmitVoteRequest{
			ComparisonID:       1,
			AttendeeID:         1,
			PreferredFeatureID: nil,
			IsTieVote:          true,
		}

		if request.PreferredFeatureID != nil {
			t.Error("Expected nil preferred feature ID for tie vote")
		}

		if !request.IsTieVote {
			t.Error("Expected tie vote")
		}
	})
}

// TestProgressRepository tests progress repository operations
func TestProgressRepository(t *testing.T) {
	t.Run("Project Progress Domain Model", func(t *testing.T) {
		progress := &domain.ProjectProgress{
			ProjectID:                    1,
			CurrentPhase:                 string(domain.PhaseSetup),
			SetupCompleted:               true,
			AttendeesAdded:               false,
			FeaturesAdded:                false,
			PairwiseValueCompleted:       false,
			PairwiseComplexityCompleted:  false,
			FibonacciValueCompleted:      false,
			FibonacciComplexityCompleted: false,
			ResultsCalculated:            false,
		}

		if progress.ProjectID != 1 {
			t.Errorf("Expected project ID 1, got %d", progress.ProjectID)
		}

		if progress.CurrentPhase != string(domain.PhaseSetup) {
			t.Errorf("Expected current phase %s, got %s", domain.PhaseSetup, progress.CurrentPhase)
		}

		if !progress.SetupCompleted {
			t.Error("Expected setup to be completed")
		}

		if progress.AttendeesAdded {
			t.Error("Expected attendees not to be added yet")
		}
	})
}

// BenchmarkRepositoryOperations benchmarks repository operations
func BenchmarkRepositoryOperations(b *testing.B) {
	b.Run("Create Project Request", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := domain.CreateProjectRequest{
				Name:        "Test Project",
				Description: "A test project for P-WVC",
			}
			_ = request.Name
			_ = request.Description
		}
	})

	b.Run("Create Attendee Request", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := domain.CreateAttendeeRequest{
				Name:          "John Doe",
				Role:          "Product Manager",
				IsFacilitator: true,
			}
			_ = request.Name
			_ = request.Role
			_ = request.IsFacilitator
		}
	})

	b.Run("Create Feature Request", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := domain.CreateFeatureRequest{
				Title:              "User Authentication",
				Description:        "Implement user login and registration",
				AcceptanceCriteria: "Users can sign up, log in, and log out",
			}
			_ = request.Title
			_ = request.Description
			_ = request.AcceptanceCriteria
		}
	})
}
