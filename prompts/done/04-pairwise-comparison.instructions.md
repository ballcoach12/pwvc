# Instructions: Pairwise Comparison Backend Implementation

## Pairwise Comparison Domain Design

### Core Entities and Relationships
```go
// internal/domain/pairwise.go
type PairwiseSession struct {
    ID            uint                  `json:"id" gorm:"primaryKey"`
    ProjectID     uint                  `json:"project_id"`
    CriterionType string                `json:"criterion_type" gorm:"type:varchar(20);check:criterion_type IN ('value','complexity')"`
    Status        string                `json:"status" gorm:"default:active;check:status IN ('active','completed')"`
    StartedAt     time.Time             `json:"started_at" gorm:"default:CURRENT_TIMESTAMP"`
    CompletedAt   *time.Time            `json:"completed_at,omitempty"`
    
    // Relationships
    Project     Project             `json:"project,omitempty"`
    Comparisons []PairwiseComparison `json:"comparisons,omitempty"`
}

type PairwiseComparison struct {
    ID               uint           `json:"id" gorm:"primaryKey"`
    SessionID        uint           `json:"session_id"`
    FeatureAID       uint           `json:"feature_a_id"`
    FeatureBID       uint           `json:"feature_b_id"`
    WinnerID         *uint          `json:"winner_id,omitempty"` // NULL for ties
    IsTie            bool           `json:"is_tie" gorm:"default:false"`
    ConsensusReached bool           `json:"consensus_reached" gorm:"default:false"`
    CreatedAt        time.Time      `json:"created_at"`
    
    // Relationships
    Session   PairwiseSession `json:"session,omitempty"`
    FeatureA  Feature         `json:"feature_a,omitempty"`
    FeatureB  Feature         `json:"feature_b,omitempty"`
    Winner    *Feature        `json:"winner,omitempty"`
    Votes     []AttendeeVote  `json:"votes,omitempty"`
}

type AttendeeVote struct {
    ID                 uint      `json:"id" gorm:"primaryKey"`
    ComparisonID       uint      `json:"comparison_id"`
    AttendeeID         uint      `json:"attendee_id"`
    PreferredFeatureID *uint     `json:"preferred_feature_id,omitempty"` // NULL for tie votes
    IsTieVote          bool      `json:"is_tie_vote" gorm:"default:false"`
    VotedAt            time.Time `json:"voted_at" gorm:"default:CURRENT_TIMESTAMP"`
    
    // Relationships
    Comparison       PairwiseComparison `json:"comparison,omitempty"`
    Attendee        Attendee           `json:"attendee,omitempty"`
    PreferredFeature *Feature           `json:"preferred_feature,omitempty"`
}
```

## Pairwise Generation Algorithm

### Complete Pair Generation
```go
// internal/service/pairwise_service.go
type PairwiseService struct {
    sessionRepo    repository.PairwiseSessionRepository
    comparisonRepo repository.PairwiseComparisonRepository
    featureRepo    repository.FeatureRepository
    voteRepo       repository.AttendeeVoteRepository
    logger         *slog.Logger
}

func (s *PairwiseService) GenerateAllComparisons(sessionID uint, features []domain.Feature) error {
    if len(features) < 2 {
        return errors.New("minimum 2 features required for pairwise comparison")
    }
    
    var comparisons []domain.PairwiseComparison
    
    // Generate all unique pairs (n choose 2)
    for i := 0; i < len(features); i++ {
        for j := i + 1; j < len(features); j++ {
            comparison := domain.PairwiseComparison{
                SessionID:        sessionID,
                FeatureAID:       features[i].ID,
                FeatureBID:       features[j].ID,
                ConsensusReached: false,
            }
            comparisons = append(comparisons, comparison)
        }
    }
    
    s.logger.Info("Generated pairwise comparisons", 
        "session_id", sessionID, 
        "total_comparisons", len(comparisons),
        "features", len(features))
    
    return s.comparisonRepo.CreateBatch(comparisons)
}

func (s *PairwiseService) CalculateExpectedComparisons(featureCount int) int {
    if featureCount < 2 {
        return 0
    }
    return featureCount * (featureCount - 1) / 2
}
```

## Consensus Tracking Engine

### Vote Processing and Consensus Detection
```go
// internal/service/consensus_service.go
type ConsensusService struct {
    voteRepo       repository.AttendeeVoteRepository
    comparisonRepo repository.PairwiseComparisonRepository
    attendeeRepo   repository.AttendeeRepository
    logger         *slog.Logger
}

type ConsensusResult struct {
    ComparisonID     uint  `json:"comparison_id"`
    ConsensusReached bool  `json:"consensus_reached"`
    WinnerID         *uint `json:"winner_id"`
    IsTie            bool  `json:"is_tie"`
    VoteSummary      VoteSummary `json:"vote_summary"`
}

type VoteSummary struct {
    TotalAttendees    int            `json:"total_attendees"`
    VotesCast         int            `json:"votes_cast"`
    FeatureAVotes     int            `json:"feature_a_votes"`
    FeatureBVotes     int            `json:"feature_b_votes"`
    TieVotes          int            `json:"tie_votes"`
    AttendeesVoted    []string       `json:"attendees_voted"`
    AttendeesPending  []string       `json:"attendees_pending"`
}

func (s *ConsensusService) ProcessVote(vote *domain.AttendeeVote) (*ConsensusResult, error) {
    // Validate vote
    if err := s.validateVote(vote); err != nil {
        return nil, fmt.Errorf("vote validation failed: %w", err)
    }
    
    // Save or update vote (upsert pattern)
    if err := s.voteRepo.UpsertVote(vote); err != nil {
        return nil, fmt.Errorf("failed to save vote: %w", err)
    }
    
    // Check if consensus is reached
    result, err := s.CheckConsensus(vote.ComparisonID)
    if err != nil {
        return nil, fmt.Errorf("consensus check failed: %w", err)
    }
    
    // Update comparison if consensus reached
    if result.ConsensusReached {
        err := s.comparisonRepo.UpdateConsensus(vote.ComparisonID, result.WinnerID, result.IsTie)
        if err != nil {
            return nil, fmt.Errorf("failed to update comparison consensus: %w", err)
        }
        
        s.logger.Info("Consensus reached", 
            "comparison_id", vote.ComparisonID,
            "winner_id", result.WinnerID,
            "is_tie", result.IsTie)
    }
    
    return result, nil
}

func (s *ConsensusService) CheckConsensus(comparisonID uint) (*ConsensusResult, error) {
    // Get all votes for this comparison
    votes, err := s.voteRepo.GetByComparisonID(comparisonID)
    if err != nil {
        return nil, err
    }
    
    // Get all attendees for the project
    comparison, err := s.comparisonRepo.GetByIDWithSession(comparisonID)
    if err != nil {
        return nil, err
    }
    
    attendees, err := s.attendeeRepo.GetByProjectID(comparison.Session.ProjectID)
    if err != nil {
        return nil, err
    }
    
    // Build vote summary
    summary := s.buildVoteSummary(votes, attendees)
    
    // Check if all attendees have voted
    if summary.VotesCast < summary.TotalAttendees {
        return &ConsensusResult{
            ComparisonID:     comparisonID,
            ConsensusReached: false,
            VoteSummary:      summary,
        }, nil
    }
    
    // Determine consensus
    consensus := s.determineConsensus(votes)
    
    return &ConsensusResult{
        ComparisonID:     comparisonID,
        ConsensusReached: consensus.HasConsensus,
        WinnerID:         consensus.WinnerID,
        IsTie:            consensus.IsTie,
        VoteSummary:      summary,
    }, nil
}

type consensusDecision struct {
    HasConsensus bool
    WinnerID     *uint
    IsTie        bool
}

func (s *ConsensusService) determineConsensus(votes []domain.AttendeeVote) consensusDecision {
    if len(votes) == 0 {
        return consensusDecision{HasConsensus: false}
    }
    
    // Count votes by type
    voteCounts := make(map[string]int)
    var firstVoteType string
    var firstWinnerID *uint
    
    for _, vote := range votes {
        var voteKey string
        if vote.IsTieVote {
            voteKey = "tie"
        } else {
            voteKey = fmt.Sprintf("feature_%d", *vote.PreferredFeatureID)
        }
        
        voteCounts[voteKey]++
        
        if firstVoteType == "" {
            firstVoteType = voteKey
            firstWinnerID = vote.PreferredFeatureID
        }
    }
    
    // Consensus requires all votes to be identical
    if len(voteCounts) == 1 {
        if firstVoteType == "tie" {
            return consensusDecision{
                HasConsensus: true,
                IsTie:        true,
                WinnerID:     nil,
            }
        } else {
            return consensusDecision{
                HasConsensus: true,
                IsTie:        false,
                WinnerID:     firstWinnerID,
            }
        }
    }
    
    return consensusDecision{HasConsensus: false}
}
```

## Session Progress Tracking

### Progress Calculation and Status Updates
```go
// internal/service/session_progress_service.go
type SessionProgressService struct {
    sessionRepo    repository.PairwiseSessionRepository
    comparisonRepo repository.PairwiseComparisonRepository
}

type SessionProgress struct {
    SessionID            uint    `json:"session_id"`
    TotalComparisons     int     `json:"total_comparisons"`
    CompletedComparisons int     `json:"completed_comparisons"`
    ProgressPercentage   float64 `json:"progress_percentage"`
    IsComplete           bool    `json:"is_complete"`
    ComparisonDetails    []ComparisonProgress `json:"comparison_details"`
}

type ComparisonProgress struct {
    ComparisonID     uint   `json:"comparison_id"`
    FeatureATitle    string `json:"feature_a_title"`
    FeatureBTitle    string `json:"feature_b_title"`
    ConsensusReached bool   `json:"consensus_reached"`
    VoteCount        int    `json:"vote_count"`
    RequiredVotes    int    `json:"required_votes"`
}

func (s *SessionProgressService) GetSessionProgress(sessionID uint) (*SessionProgress, error) {
    session, err := s.sessionRepo.GetByIDWithComparisons(sessionID)
    if err != nil {
        return nil, err
    }
    
    totalComparisons := len(session.Comparisons)
    completedComparisons := 0
    
    var details []ComparisonProgress
    for _, comparison := range session.Comparisons {
        if comparison.ConsensusReached {
            completedComparisons++
        }
        
        // Get vote details for this comparison
        voteCount := len(comparison.Votes)
        requiredVotes := s.getRequiredVoteCount(session.ProjectID)
        
        details = append(details, ComparisonProgress{
            ComparisonID:     comparison.ID,
            FeatureATitle:    comparison.FeatureA.Title,
            FeatureBTitle:    comparison.FeatureB.Title,
            ConsensusReached: comparison.ConsensusReached,
            VoteCount:        voteCount,
            RequiredVotes:    requiredVotes,
        })
    }
    
    progressPercentage := 0.0
    if totalComparisons > 0 {
        progressPercentage = float64(completedComparisons) / float64(totalComparisons) * 100
    }
    
    isComplete := completedComparisons == totalComparisons && totalComparisons > 0
    
    // Auto-complete session if all comparisons are done
    if isComplete && session.Status != "completed" {
        if err := s.sessionRepo.MarkCompleted(sessionID); err != nil {
            return nil, fmt.Errorf("failed to mark session completed: %w", err)
        }
    }
    
    return &SessionProgress{
        SessionID:            sessionID,
        TotalComparisons:     totalComparisons,
        CompletedComparisons: completedComparisons,
        ProgressPercentage:   progressPercentage,
        IsComplete:           isComplete,
        ComparisonDetails:    details,
    }, nil
}
```

## Repository Patterns for Pairwise Data

### Efficient Database Queries
```go
// internal/repository/pairwise_comparison_repository.go
type PairwiseComparisonRepository interface {
    CreateBatch(comparisons []domain.PairwiseComparison) error
    GetBySessionID(sessionID uint) ([]domain.PairwiseComparison, error)
    GetByIDWithFeatures(id uint) (*domain.PairwiseComparison, error)
    UpdateConsensus(comparisonID uint, winnerID *uint, isTie bool) error
    GetSessionCompletionStats(sessionID uint) (*CompletionStats, error)
}

type completionStats struct {
    TotalComparisons     int `json:"total_comparisons"`
    CompletedComparisons int `json:"completed_comparisons"`
}

func (r *pairwiseComparisonRepository) GetSessionCompletionStats(sessionID uint) (*CompletionStats, error) {
    var stats CompletionStats
    
    err := r.db.Model(&domain.PairwiseComparison{}).
        Select("COUNT(*) as total_comparisons, COUNT(CASE WHEN consensus_reached = true THEN 1 END) as completed_comparisons").
        Where("session_id = ?", sessionID).
        Scan(&stats).Error
    
    return &stats, err
}

// Optimized query to get comparisons with all related data
func (r *pairwiseComparisonRepository) GetBySessionIDWithDetails(sessionID uint) ([]domain.PairwiseComparison, error) {
    var comparisons []domain.PairwiseComparison
    
    return comparisons, r.db.
        Preload("FeatureA").
        Preload("FeatureB").
        Preload("Winner").
        Preload("Votes").
        Preload("Votes.Attendee").
        Where("session_id = ?", sessionID).
        Order("created_at ASC").
        Find(&comparisons).Error
}

// Batch update for consensus results
func (r *pairwiseComparisonRepository) UpdateConsensus(comparisonID uint, winnerID *uint, isTie bool) error {
    updates := map[string]interface{}{
        "consensus_reached": true,
        "is_tie":           isTie,
        "winner_id":        winnerID,
    }
    
    return r.db.Model(&domain.PairwiseComparison{}).
        Where("id = ?", comparisonID).
        Updates(updates).Error
}
```

## API Endpoint Patterns

### RESTful Pairwise Comparison APIs
```go
// internal/api/pairwise_handler.go
type PairwiseHandler struct {
    pairwiseService  *service.PairwiseService
    consensusService *service.ConsensusService
    progressService  *service.SessionProgressService
    logger          *slog.Logger
}

// POST /api/projects/{id}/pairwise-sessions
func (h *PairwiseHandler) CreateSession(c *gin.Context) {
    projectID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    
    var request CreatePairwiseSessionRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format"})
        return
    }
    
    // Validate criterion type
    if request.CriterionType != "value" && request.CriterionType != "complexity" {
        c.JSON(400, gin.H{"error": "Criterion type must be 'value' or 'complexity'"})
        return
    }
    
    session, err := h.pairwiseService.CreateSession(uint(projectID), request.CriterionType)
    if err != nil {
        h.logger.Error("Failed to create pairwise session", "error", err)
        c.JSON(500, gin.H{"error": "Failed to create session"})
        return
    }
    
    c.JSON(201, session)
}

// POST /api/projects/{id}/pairwise-sessions/{session_id}/vote
func (h *PairwiseHandler) SubmitVote(c *gin.Context) {
    sessionID, _ := strconv.ParseUint(c.Param("session_id"), 10, 32)
    
    var request SubmitVoteRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": "Invalid vote format"})
        return
    }
    
    vote := &domain.AttendeeVote{
        ComparisonID:       request.ComparisonID,
        AttendeeID:         request.AttendeeID,
        PreferredFeatureID: request.PreferredFeatureID,
        IsTieVote:          request.IsTieVote,
    }
    
    result, err := h.consensusService.ProcessVote(vote)
    if err != nil {
        h.logger.Error("Failed to process vote", "error", err)
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, result)
}

// GET /api/projects/{id}/pairwise-sessions/{session_id}/progress
func (h *PairwiseHandler) GetSessionProgress(c *gin.Context) {
    sessionID, _ := strconv.ParseUint(c.Param("session_id"), 10, 32)
    
    progress, err := h.progressService.GetSessionProgress(uint(sessionID))
    if err != nil {
        h.logger.Error("Failed to get session progress", "error", err)
        c.JSON(500, gin.H{"error": "Failed to get progress"})
        return
    }
    
    c.JSON(200, progress)
}
```

## Testing Patterns

### Pairwise Logic Testing
```go
// internal/service/consensus_service_test.go
func TestConsensusService_DetermineConsensus(t *testing.T) {
    tests := []struct {
        name     string
        votes    []domain.AttendeeVote
        expected consensusDecision
    }{
        {
            name: "unanimous winner",
            votes: []domain.AttendeeVote{
                {PreferredFeatureID: &[]uint{1}[0], IsTieVote: false},
                {PreferredFeatureID: &[]uint{1}[0], IsTieVote: false},
                {PreferredFeatureID: &[]uint{1}[0], IsTieVote: false},
            },
            expected: consensusDecision{
                HasConsensus: true,
                WinnerID:     &[]uint{1}[0],
                IsTie:        false,
            },
        },
        {
            name: "unanimous tie",
            votes: []domain.AttendeeVote{
                {IsTieVote: true},
                {IsTieVote: true},
            },
            expected: consensusDecision{
                HasConsensus: true,
                WinnerID:     nil,
                IsTie:        true,
            },
        },
        {
            name: "no consensus - split votes",
            votes: []domain.AttendeeVote{
                {PreferredFeatureID: &[]uint{1}[0], IsTieVote: false},
                {PreferredFeatureID: &[]uint{2}[0], IsTieVote: false},
            },
            expected: consensusDecision{
                HasConsensus: false,
            },
        },
    }
    
    service := NewConsensusService(nil, nil, nil, nil)
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := service.determineConsensus(tt.votes)
            assert.Equal(t, tt.expected.HasConsensus, result.HasConsensus)
            if tt.expected.WinnerID != nil {
                assert.Equal(t, *tt.expected.WinnerID, *result.WinnerID)
            }
            assert.Equal(t, tt.expected.IsTie, result.IsTie)
        })
    }
}
```