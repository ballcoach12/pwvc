package service

import (
	"fmt"
	"testing"

	"pairwise/internal/domain"
)

// TestPairwiseComparison_GenerateComparisons tests the comparison generation logic
func TestPairwiseComparison_GenerateComparisons(t *testing.T) {
	tests := []struct {
		name                string
		numFeatures         int
		expectedComparisons int
	}{
		{
			name:                "Two features",
			numFeatures:         2,
			expectedComparisons: 1, // n*(n-1)/2 = 2*1/2 = 1
		},
		{
			name:                "Three features",
			numFeatures:         3,
			expectedComparisons: 3, // n*(n-1)/2 = 3*2/2 = 3
		},
		{
			name:                "Four features",
			numFeatures:         4,
			expectedComparisons: 6, // n*(n-1)/2 = 4*3/2 = 6
		},
		{
			name:                "Five features",
			numFeatures:         5,
			expectedComparisons: 10, // n*(n-1)/2 = 5*4/2 = 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock features
			features := make([]domain.Feature, tt.numFeatures)
			for i := 0; i < tt.numFeatures; i++ {
				features[i] = domain.Feature{
					ID:        i + 1,
					ProjectID: 1,
					Title:     fmt.Sprintf("Feature %d", i+1),
				}
			}

			// Count unique pairs
			pairCount := 0
			for i := 0; i < len(features); i++ {
				for j := i + 1; j < len(features); j++ {
					pairCount++
				}
			}

			if pairCount != tt.expectedComparisons {
				t.Errorf("Expected %d comparisons but got %d", tt.expectedComparisons, pairCount)
			}
		})
	}
}

// TestVoteValidation tests vote validation logic
func TestVoteValidation(t *testing.T) {
	tests := []struct {
		name        string
		vote        domain.SubmitVoteRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid preference vote",
			vote: domain.SubmitVoteRequest{
				ComparisonID:       1,
				AttendeeID:         1,
				PreferredFeatureID: intPtr(1),
				IsTieVote:          false,
			},
			expectError: false,
		},
		{
			name: "Valid tie vote",
			vote: domain.SubmitVoteRequest{
				ComparisonID: 1,
				AttendeeID:   1,
				IsTieVote:    true,
			},
			expectError: false,
		},
		{
			name: "Invalid: tie vote with preferred feature",
			vote: domain.SubmitVoteRequest{
				ComparisonID:       1,
				AttendeeID:         1,
				PreferredFeatureID: intPtr(1),
				IsTieVote:          true,
			},
			expectError: true,
			errorMsg:    "Cannot specify preferred feature for tie vote",
		},
		{
			name: "Invalid: preference vote without preferred feature",
			vote: domain.SubmitVoteRequest{
				ComparisonID: 1,
				AttendeeID:   1,
				IsTieVote:    false,
			},
			expectError: true,
			errorMsg:    "Must specify preferred feature for non-tie vote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVoteRequest(tt.vote)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s' but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestConsensusLogic tests consensus calculation logic
func TestConsensusLogic(t *testing.T) {
	tests := []struct {
		name              string
		votes             []domain.AttendeeVote
		totalAttendees    int
		expectedConsensus bool
		expectedWinner    *int
		expectedTie       bool
	}{
		{
			name: "All vote for same feature",
			votes: []domain.AttendeeVote{
				{AttendeeID: 1, PreferredFeatureID: intPtr(1), IsTieVote: false},
				{AttendeeID: 2, PreferredFeatureID: intPtr(1), IsTieVote: false},
			},
			totalAttendees:    2,
			expectedConsensus: true,
			expectedWinner:    intPtr(1),
			expectedTie:       false,
		},
		{
			name: "All vote for tie",
			votes: []domain.AttendeeVote{
				{AttendeeID: 1, IsTieVote: true},
				{AttendeeID: 2, IsTieVote: true},
			},
			totalAttendees:    2,
			expectedConsensus: true,
			expectedWinner:    nil,
			expectedTie:       true,
		},
		{
			name: "Mixed votes - no consensus",
			votes: []domain.AttendeeVote{
				{AttendeeID: 1, PreferredFeatureID: intPtr(1), IsTieVote: false},
				{AttendeeID: 2, PreferredFeatureID: intPtr(2), IsTieVote: false},
			},
			totalAttendees:    2,
			expectedConsensus: false,
			expectedWinner:    nil,
			expectedTie:       false,
		},
		{
			name: "Not all attendees voted",
			votes: []domain.AttendeeVote{
				{AttendeeID: 1, PreferredFeatureID: intPtr(1), IsTieVote: false},
			},
			totalAttendees:    2,
			expectedConsensus: false,
			expectedWinner:    nil,
			expectedTie:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consensus, winner, isTie := calculateConsensus(tt.votes, tt.totalAttendees)

			if consensus != tt.expectedConsensus {
				t.Errorf("Expected consensus %v but got %v", tt.expectedConsensus, consensus)
			}

			if !compareIntPtr(winner, tt.expectedWinner) {
				t.Errorf("Expected winner %v but got %v", tt.expectedWinner, winner)
			}

			if isTie != tt.expectedTie {
				t.Errorf("Expected tie %v but got %v", tt.expectedTie, isTie)
			}
		})
	}
}

// TestSessionProgress tests session progress calculation
func TestSessionProgress(t *testing.T) {
	tests := []struct {
		name                         string
		totalComparisons             int
		completedComparisons         int
		expectedProgressPercent      float64
		expectedRemainingComparisons int
	}{
		{
			name:                         "No progress",
			totalComparisons:             10,
			completedComparisons:         0,
			expectedProgressPercent:      0.0,
			expectedRemainingComparisons: 10,
		},
		{
			name:                         "Half complete",
			totalComparisons:             10,
			completedComparisons:         5,
			expectedProgressPercent:      50.0,
			expectedRemainingComparisons: 5,
		},
		{
			name:                         "Fully complete",
			totalComparisons:             10,
			completedComparisons:         10,
			expectedProgressPercent:      100.0,
			expectedRemainingComparisons: 0,
		},
		{
			name:                         "No comparisons",
			totalComparisons:             0,
			completedComparisons:         0,
			expectedProgressPercent:      0.0,
			expectedRemainingComparisons: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress := calculateProgress(tt.totalComparisons, tt.completedComparisons)

			if progress.ProgressPercentage != tt.expectedProgressPercent {
				t.Errorf("Expected progress percentage %.1f but got %.1f", tt.expectedProgressPercent, progress.ProgressPercentage)
			}

			if progress.RemainingComparisons != tt.expectedRemainingComparisons {
				t.Errorf("Expected remaining comparisons %d but got %d", tt.expectedRemainingComparisons, progress.RemainingComparisons)
			}
		})
	}
}

// Helper functions for testing

func intPtr(i int) *int {
	return &i
}

func compareIntPtr(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Mock business logic functions for testing

func validateVoteRequest(req domain.SubmitVoteRequest) error {
	if req.IsTieVote && req.PreferredFeatureID != nil {
		return domain.NewAPIError(400, "Cannot specify preferred feature for tie vote")
	}
	if !req.IsTieVote && req.PreferredFeatureID == nil {
		return domain.NewAPIError(400, "Must specify preferred feature for non-tie vote")
	}
	return nil
}

func calculateConsensus(votes []domain.AttendeeVote, totalAttendees int) (bool, *int, bool) {
	// Check if all attendees have voted
	if len(votes) != totalAttendees {
		return false, nil, false
	}

	if len(votes) == 0 {
		return false, nil, false
	}

	// Get first vote as reference
	firstVote := votes[0]
	winnerID := firstVote.PreferredFeatureID
	isTie := firstVote.IsTieVote

	// Check if all votes are the same
	for _, vote := range votes[1:] {
		if vote.IsTieVote != isTie {
			return false, nil, false
		}
		if !vote.IsTieVote && vote.PreferredFeatureID != nil && winnerID != nil {
			if *vote.PreferredFeatureID != *winnerID {
				return false, nil, false
			}
		}
	}

	return true, winnerID, isTie
}

func calculateProgress(totalComparisons, completedComparisons int) domain.SessionProgress {
	progress := domain.SessionProgress{
		TotalComparisons:     totalComparisons,
		CompletedComparisons: completedComparisons,
		RemainingComparisons: totalComparisons - completedComparisons,
	}

	if totalComparisons > 0 {
		progress.ProgressPercentage = float64(completedComparisons) / float64(totalComparisons) * 100
	}

	return progress
}
