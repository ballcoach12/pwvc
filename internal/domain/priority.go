package domain

import "time"

// PriorityCalculation represents the final P-WVC calculation for a feature
type PriorityCalculation struct {
	ID                 int       `json:"id" db:"id"`
	ProjectID          int       `json:"projectId" db:"project_id"`
	FeatureID          int       `json:"featureId" db:"feature_id"`
	WValue             float64   `json:"wValue" db:"w_value"`                          // Win-count weight for value
	WComplexity        float64   `json:"wComplexity" db:"w_complexity"`                // Win-count weight for complexity
	SValue             int       `json:"sValue" db:"s_value"`                          // Fibonacci score for value
	SComplexity        int       `json:"sComplexity" db:"s_complexity"`                // Fibonacci score for complexity
	WeightedValue      float64   `json:"weightedValue" db:"weighted_value"`            // SValue × WValue
	WeightedComplexity float64   `json:"weightedComplexity" db:"weighted_complexity"`  // SComplexity × WComplexity
	FinalPriorityScore float64   `json:"finalPriorityScore" db:"final_priority_score"` // Weighted Value ÷ Weighted Complexity
	Rank               int       `json:"rank" db:"rank"`
	CalculatedAt       time.Time `json:"calculatedAt" db:"calculated_at"`

	// Related data for display
	Feature *Feature `json:"feature,omitempty"`
}

// PriorityResult represents a complete feature ranking with all calculation details
type PriorityResult struct {
	PriorityCalculation
	Feature Feature `json:"feature"`
}

// ProjectResults represents the complete results for a project
type ProjectResults struct {
	ProjectID     int              `json:"projectId"`
	Project       *Project         `json:"project,omitempty"`
	Results       []PriorityResult `json:"results"`
	CalculatedAt  time.Time        `json:"calculatedAt"`
	TotalFeatures int              `json:"totalFeatures"`
	Summary       ResultsSummary   `json:"summary"`
}

// ResultsSummary provides statistical information about the results
type ResultsSummary struct {
	HighestScore float64 `json:"highestScore"`
	LowestScore  float64 `json:"lowestScore"`
	AverageScore float64 `json:"averageScore"`
	MedianScore  float64 `json:"medianScore"`
	ScoreRange   float64 `json:"scoreRange"`
	TopTier      int     `json:"topTier"`    // Count of features in top 25%
	BottomTier   int     `json:"bottomTier"` // Count of features in bottom 25%
}

// CalculateResultsRequest represents the request to calculate P-WVC results
type CalculateResultsRequest struct {
	ProjectID int `json:"projectId" binding:"required"`
}

// ExportFormat represents the different export formats
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatJira ExportFormat = "jira"
)

// JiraExport represents the Jira-compatible export format
type JiraExport struct {
	Issues []JiraIssue `json:"issues"`
}

// JiraIssue represents a single Jira issue for export
type JiraIssue struct {
	Summary      string           `json:"summary"`
	Description  string           `json:"description"`
	StoryPoints  int              `json:"storyPoints"`
	Priority     string           `json:"priority"`
	CustomFields JiraCustomFields `json:"customFields"`
}

// JiraCustomFields represents custom fields for Jira export
type JiraCustomFields struct {
	FinalPriorityScore float64 `json:"finalPriorityScore"`
	ValueScore         int     `json:"valueScore"`
	ComplexityScore    int     `json:"complexityScore"`
}

// FibonacciProgressMetrics represents progress metrics for Fibonacci scoring phases (T040 - US8)
type FibonacciProgressMetrics struct {
	ProjectID           int                           `json:"project_id"`
	CriterionType       string                        `json:"criterion_type"` // "value" or "complexity"
	CompletedScores     int                           `json:"completed_scores"`
	TotalExpectedScores int                           `json:"total_expected_scores"`
	ProgressPercentage  float64                       `json:"progress_percentage"`
	FeatureCompletion   map[int]FeatureScoreProgress  `json:"feature_completion"`
	AttendeeCompletion  map[int]AttendeeScoreProgress `json:"attendee_completion"`
}

// FeatureScoreProgress represents scoring progress for a specific feature (T040 - US8)
type FeatureScoreProgress struct {
	FeatureID           int     `json:"feature_id"`
	FeatureName         string  `json:"feature_name"`
	CompletedScores     int     `json:"completed_scores"`
	TotalExpectedScores int     `json:"total_expected_scores"`
	ProgressPercentage  float64 `json:"progress_percentage"`
}

// AttendeeScoreProgress represents scoring progress for a specific attendee (T040 - US8)
type AttendeeScoreProgress struct {
	AttendeeID          int     `json:"attendee_id"`
	AttendeeName        string  `json:"attendee_name"`
	CompletedScores     int     `json:"completed_scores"`
	TotalExpectedScores int     `json:"total_expected_scores"`
	ProgressPercentage  float64 `json:"progress_percentage"`
}

// OverallFibonacciProgress represents combined progress for both value and complexity scoring (T040 - US8)
type OverallFibonacciProgress struct {
	ProjectID            int                       `json:"project_id"`
	ValueMetrics         *FibonacciProgressMetrics `json:"value_metrics"`
	ComplexityMetrics    *FibonacciProgressMetrics `json:"complexity_metrics"`
	OverallProgress      float64                   `json:"overall_progress"`
	TotalCompleted       int                       `json:"total_completed"`
	TotalExpected        int                       `json:"total_expected"`
	IsValueComplete      bool                      `json:"is_value_complete"`
	IsComplexityComplete bool                      `json:"is_complexity_complete"`
	IsBothComplete       bool                      `json:"is_both_complete"`
}
