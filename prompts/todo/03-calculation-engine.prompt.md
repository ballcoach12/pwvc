# Prompt 3: P-WVC Calculation Engine

Build the core P-WVC mathematical engine: win-count calculation from pairwise comparisons, Fibonacci score validation, and Final Priority Score calculation (Weighted Value ÷ Weighted Complexity). Include comprehensive unit tests.

## Requirements
- Implement win-count calculation algorithm from pairwise comparison results
- Create Fibonacci sequence validation (1, 2, 3, 5, 8, 13, 21, 34, 55, 89)
- Build Final Priority Score calculation: FPS = (SValue × WValue) / (SComplexity × WComplexity)
- Add comprehensive unit tests for all mathematical functions
- Create service layer for P-WVC calculations
- Handle edge cases (ties, missing data, division by zero)

## Core Mathematical Functions

### Win-Count Calculation
```
For each feature in pairwise comparisons:
WFeature = (Total Wins + 0.5 × Total Ties) / (Total Comparisons)

Where Total Comparisons = n - 1 (n = number of features)
```

### Final Priority Score
```
Weighted Value = SValue × WValue
Weighted Complexity = SComplexity × WComplexity  
FPS = Weighted Value / Weighted Complexity
```

## Implementation Structure
- `internal/domain/calculation.go` - Core math engine
- `internal/service/pwvc_service.go` - Business logic layer
- Comprehensive unit tests with edge case coverage
- Input validation for all mathematical operations

## Test Cases to Include
- Win-count calculation with various comparison scenarios
- Fibonacci validation (valid and invalid inputs)
- FPS calculation with different score combinations
- Edge cases: all ties, zero complexity, missing data
- Integration tests with sample feature sets

## Fibonacci Validation
Valid scores: 1, 2, 3, 5, 8, 13, 21, 34, 55, 89
Invalid inputs should return appropriate error messages