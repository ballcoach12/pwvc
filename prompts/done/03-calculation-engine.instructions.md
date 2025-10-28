# Instructions: P-WVC Calculation Engine Implementation

## Mathematical Engine Architecture

### Core Calculation Types
```go
// internal/domain/calculation.go
type PWVCCalculation struct {
    FeatureID           uint    `json:"feature_id"`
    WValue              float64 `json:"w_value"`              // Win-count weight for value
    WComplexity         float64 `json:"w_complexity"`         // Win-count weight for complexity
    SValue              int     `json:"s_value"`              // Fibonacci score for value
    SComplexity         int     `json:"s_complexity"`         // Fibonacci score for complexity
    WeightedValue       float64 `json:"weighted_value"`       // SValue × WValue
    WeightedComplexity  float64 `json:"weighted_complexity"`  // SComplexity × WComplexity
    FinalPriorityScore  float64 `json:"final_priority_score"` // WeightedValue ÷ WeightedComplexity
    Rank               int     `json:"rank"`
}

// Calculation input validation
type CalculationInput struct {
    ProjectID       uint                    `validate:"required"`
    PairwiseResults map[string][]Comparison `validate:"required"` // "value" and "complexity" keys
    FibonacciScores map[string][]Score      `validate:"required"` // "value" and "complexity" keys
}
```

### Win-Count Calculation Engine
```go
// internal/service/calculation_service.go
type CalculationService struct {
    logger *slog.Logger
}

func (s *CalculationService) CalculateWinCounts(comparisons []domain.PairwiseComparison, features []domain.Feature) (map[uint]float64, error) {
    if len(features) < 2 {
        return nil, errors.New("minimum 2 features required for win-count calculation")
    }
    
    winCounts := make(map[uint]int)
    totalComparisons := make(map[uint]int)
    
    // Initialize counts
    for _, feature := range features {
        winCounts[feature.ID] = 0
        totalComparisons[feature.ID] = 0
    }
    
    // Count wins and ties
    for _, comparison := range comparisons {
        if !comparison.ConsensusReached {
            return nil, fmt.Errorf("comparison between features %d and %d has not reached consensus", 
                comparison.FeatureAID, comparison.FeatureBID)
        }
        
        totalComparisons[comparison.FeatureAID]++
        totalComparisons[comparison.FeatureBID]++
        
        if comparison.IsTie {
            // Each feature gets 0.5 points for a tie
            winCounts[comparison.FeatureAID] += 0.5
            winCounts[comparison.FeatureBID] += 0.5
        } else if comparison.WinnerID != nil {
            winCounts[*comparison.WinnerID]++
        }
    }
    
    // Calculate win-count weights (W = wins / total_comparisons)
    weights := make(map[uint]float64)
    for _, feature := range features {
        if totalComparisons[feature.ID] == 0 {
            return nil, fmt.Errorf("feature %d has no comparisons", feature.ID)
        }
        weights[feature.ID] = float64(winCounts[feature.ID]) / float64(totalComparisons[feature.ID])
    }
    
    return weights, nil
}
```

### Fibonacci Validation Engine
```go
// internal/domain/fibonacci.go
var validFibonacciNumbers = map[int]bool{
    1: true, 2: true, 3: true, 5: true, 8: true, 
    13: true, 21: true, 34: true, 55: true, 89: true,
}

type FibonacciValidator struct{}

func (v *FibonacciValidator) IsValid(score int) bool {
    return validFibonacciNumbers[score]
}

func (v *FibonacciValidator) ValidateScore(score int) error {
    if !v.IsValid(score) {
        return fmt.Errorf("invalid Fibonacci score: %d. Valid scores are: 1, 2, 3, 5, 8, 13, 21, 34, 55, 89", score)
    }
    return nil
}

func (v *FibonacciValidator) GetValidScores() []int {
    return []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
}

// Fibonacci scoring with consensus validation
type FibonacciScoreCalculator struct {
    validator *FibonacciValidator
}

func (c *FibonacciScoreCalculator) CalculateConsensusScores(scores []domain.FibonacciScore) (map[uint]int, error) {
    // Group scores by feature
    featureScores := make(map[uint][]int)
    for _, score := range scores {
        if err := c.validator.ValidateScore(score.ScoreValue); err != nil {
            return nil, err
        }
        featureScores[score.FeatureID] = append(featureScores[score.FeatureID], score.ScoreValue)
    }
    
    consensusScores := make(map[uint]int)
    for featureID, scoreList := range featureScores {
        consensus, hasConsensus := findConsensus(scoreList)
        if !hasConsensus {
            return nil, fmt.Errorf("no consensus reached for feature %d", featureID)
        }
        consensusScores[featureID] = consensus
    }
    
    return consensusScores, nil
}

func findConsensus(scores []int) (int, bool) {
    if len(scores) == 0 {
        return 0, false
    }
    
    // Check if all scores are the same
    first := scores[0]
    for _, score := range scores[1:] {
        if score != first {
            return 0, false
        }
    }
    
    return first, true
}
```

### Final Priority Score Calculator
```go
// internal/service/priority_calculator.go
type PriorityCalculator struct {
    calculationService *CalculationService
    fibonacciCalc     *FibonacciScoreCalculator
    logger            *slog.Logger
}

func (p *PriorityCalculator) CalculateFinalPriorityScores(input *CalculationInput) ([]PWVCCalculation, error) {
    // Validate input
    if err := p.validateInput(input); err != nil {
        return nil, fmt.Errorf("input validation failed: %w", err)
    }
    
    // Calculate win-count weights for both criteria
    valueWeights, err := p.calculationService.CalculateWinCounts(
        input.PairwiseResults["value"], 
        input.Features,
    )
    if err != nil {
        return nil, fmt.Errorf("value win-count calculation failed: %w", err)
    }
    
    complexityWeights, err := p.calculationService.CalculateWinCounts(
        input.PairwiseResults["complexity"], 
        input.Features,
    )
    if err != nil {
        return nil, fmt.Errorf("complexity win-count calculation failed: %w", err)
    }
    
    // Get consensus Fibonacci scores
    valueScores, err := p.fibonacciCalc.CalculateConsensusScores(input.FibonacciScores["value"])
    if err != nil {
        return nil, fmt.Errorf("value consensus calculation failed: %w", err)
    }
    
    complexityScores, err := p.fibonacciCalc.CalculateConsensusScores(input.FibonacciScores["complexity"])
    if err != nil {
        return nil, fmt.Errorf("complexity consensus calculation failed: %w", err)
    }
    
    // Calculate Final Priority Scores
    var results []PWVCCalculation
    for _, feature := range input.Features {
        calc := PWVCCalculation{
            FeatureID:   feature.ID,
            WValue:      valueWeights[feature.ID],
            WComplexity: complexityWeights[feature.ID],
            SValue:      valueScores[feature.ID],
            SComplexity: complexityScores[feature.ID],
        }
        
        // Calculate weighted scores
        calc.WeightedValue = float64(calc.SValue) * calc.WValue
        calc.WeightedComplexity = float64(calc.SComplexity) * calc.WComplexity
        
        // Handle division by zero
        if calc.WeightedComplexity == 0 {
            return nil, fmt.Errorf("weighted complexity is zero for feature %d", feature.ID)
        }
        
        // Calculate Final Priority Score: FPS = WeightedValue / WeightedComplexity
        calc.FinalPriorityScore = calc.WeightedValue / calc.WeightedComplexity
        
        results = append(results, calc)
    }
    
    // Sort by Final Priority Score (descending)
    sort.Slice(results, func(i, j int) bool {
        return results[i].FinalPriorityScore > results[j].FinalPriorityScore
    })
    
    // Assign ranks
    for i := range results {
        results[i].Rank = i + 1
    }
    
    return results, nil
}
```

## Precision and Rounding Patterns

### Decimal Precision Handling
```go
// Use decimal package for financial-grade precision
import "github.com/shopspring/decimal"

type PrecisePWVCCalculation struct {
    FeatureID           uint            `json:"feature_id"`
    WValue              decimal.Decimal `json:"w_value"`
    WComplexity         decimal.Decimal `json:"w_complexity"`
    WeightedValue       decimal.Decimal `json:"weighted_value"`
    WeightedComplexity  decimal.Decimal `json:"weighted_complexity"`
    FinalPriorityScore  decimal.Decimal `json:"final_priority_score"`
    Rank               int             `json:"rank"`
}

func (p *PriorityCalculator) CalculateWithPrecision(input *CalculationInput) ([]PrecisePWVCCalculation, error) {
    // Convert to decimal for precise calculations
    valueWeight := decimal.NewFromFloat(valueWeights[feature.ID])
    complexityWeight := decimal.NewFromFloat(complexityWeights[feature.ID])
    valueScore := decimal.NewFromInt(int64(valueScores[feature.ID]))
    complexityScore := decimal.NewFromInt(int64(complexityScores[feature.ID]))
    
    weightedValue := valueScore.Mul(valueWeight)
    weightedComplexity := complexityScore.Mul(complexityWeight)
    
    if weightedComplexity.IsZero() {
        return nil, fmt.Errorf("weighted complexity is zero for feature %d", feature.ID)
    }
    
    finalPriorityScore := weightedValue.Div(weightedComplexity)
    
    // Round to 6 decimal places for display
    finalPriorityScore = finalPriorityScore.Round(6)
    
    return finalPriorityScore, nil
}
```

## Comprehensive Testing Patterns

### Mathematical Accuracy Tests
```go
// internal/service/calculation_service_test.go
func TestCalculationService_WinCountAccuracy(t *testing.T) {
    tests := []struct {
        name        string
        comparisons []domain.PairwiseComparison
        features    []domain.Feature
        expected    map[uint]float64
    }{
        {
            name: "perfect consensus - Feature A always wins",
            comparisons: []domain.PairwiseComparison{
                {FeatureAID: 1, FeatureBID: 2, WinnerID: &[]uint{1}[0], ConsensusReached: true},
                {FeatureAID: 1, FeatureBID: 3, WinnerID: &[]uint{1}[0], ConsensusReached: true},
                {FeatureAID: 2, FeatureBID: 3, WinnerID: &[]uint{2}[0], ConsensusReached: true},
            },
            features: []domain.Feature{{ID: 1}, {ID: 2}, {ID: 3}},
            expected: map[uint]float64{
                1: 1.0,    // 2 wins out of 2 comparisons
                2: 0.5,    // 1 win out of 2 comparisons  
                3: 0.0,    // 0 wins out of 2 comparisons
            },
        },
        {
            name: "all ties scenario",
            comparisons: []domain.PairwiseComparison{
                {FeatureAID: 1, FeatureBID: 2, IsTie: true, ConsensusReached: true},
                {FeatureAID: 1, FeatureBID: 3, IsTie: true, ConsensusReached: true},
                {FeatureAID: 2, FeatureBID: 3, IsTie: true, ConsensusReached: true},
            },
            features: []domain.Feature{{ID: 1}, {ID: 2}, {ID: 3}},
            expected: map[uint]float64{
                1: 0.5,    // 1 tie (0.5 + 0.5) out of 2 comparisons
                2: 0.5,    // 1 tie (0.5 + 0.5) out of 2 comparisons
                3: 0.5,    // 1 tie (0.5 + 0.5) out of 2 comparisons
            },
        },
    }
    
    service := NewCalculationService()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := service.CalculateWinCounts(tt.comparisons, tt.features)
            require.NoError(t, err)
            
            for featureID, expectedWeight := range tt.expected {
                assert.InDelta(t, expectedWeight, result[featureID], 0.0001, 
                    "Feature %d weight mismatch", featureID)
            }
        })
    }
}

func TestPriorityCalculator_FinalPriorityScore(t *testing.T) {
    // Test actual P-WVC formula: FPS = (SValue × WValue) / (SComplexity × WComplexity)
    calc := PWVCCalculation{
        SValue:      8,    // Fibonacci score for value
        WValue:      0.75, // Win-count weight for value
        SComplexity: 3,    // Fibonacci score for complexity  
        WComplexity: 0.88, // Win-count weight for complexity
    }
    
    expectedWeightedValue := 8.0 * 0.75 // = 6.0
    expectedWeightedComplexity := 3.0 * 0.88 // = 2.64
    expectedFPS := 6.0 / 2.64 // ≈ 2.27
    
    calculator := &PriorityCalculator{}
    result := calculator.calculateSingleFPS(calc)
    
    assert.InDelta(t, expectedFPS, result.FinalPriorityScore, 0.01)
    assert.Equal(t, expectedWeightedValue, result.WeightedValue)
    assert.InDelta(t, expectedWeightedComplexity, result.WeightedComplexity, 0.01)
}
```

## Error Handling and Edge Cases

### Comprehensive Edge Case Handling
```go
func (p *PriorityCalculator) validateInput(input *CalculationInput) error {
    if len(input.Features) < 2 {
        return errors.New("minimum 2 features required for P-WVC calculation")
    }
    
    expectedComparisons := len(input.Features) * (len(input.Features) - 1) / 2
    
    if len(input.PairwiseResults["value"]) != expectedComparisons {
        return fmt.Errorf("expected %d value comparisons, got %d", 
            expectedComparisons, len(input.PairwiseResults["value"]))
    }
    
    if len(input.PairwiseResults["complexity"]) != expectedComparisons {
        return fmt.Errorf("expected %d complexity comparisons, got %d", 
            expectedComparisons, len(input.PairwiseResults["complexity"]))
    }
    
    // Validate all comparisons have reached consensus
    for criterion, comparisons := range input.PairwiseResults {
        for _, comparison := range comparisons {
            if !comparison.ConsensusReached {
                return fmt.Errorf("%s comparison between features %d and %d has not reached consensus", 
                    criterion, comparison.FeatureAID, comparison.FeatureBID)
            }
        }
    }
    
    return nil
}
```