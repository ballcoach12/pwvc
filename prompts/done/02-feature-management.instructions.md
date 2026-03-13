# Instructions: Feature Management Implementation

## Feature Entity Design Patterns

### Core Feature Model
```go
// internal/domain/feature.go
type Feature struct {
    ID                 uint      `json:"id" gorm:"primaryKey"`
    ProjectID          uint      `json:"project_id"`
    Title              string    `json:"title" validate:"required,min=3,max=255"`
    Description        string    `json:"description" validate:"required,min=10,max=5000"`
    AcceptanceCriteria string    `json:"acceptance_criteria" validate:"max=5000"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    
    // Relationships for scoring
    PairwiseComparisons []PairwiseComparison `json:"-" gorm:"foreignKey:FeatureAID"`
    FibonacciScores     []FibonacciScore     `json:"-"`
}

// Validation methods
func (f *Feature) Validate() error {
    if len(strings.TrimSpace(f.Title)) < 3 {
        return errors.New("title must be at least 3 characters")
    }
    if len(strings.TrimSpace(f.Description)) < 10 {
        return errors.New("description must be at least 10 characters")
    }
    return nil
}
```

## CSV Import/Export Patterns

### CSV Structure and Validation
```go
// internal/service/csv_service.go
type CSVRecord struct {
    Title              string `csv:"title"`
    Description        string `csv:"description"`
    AcceptanceCriteria string `csv:"acceptance_criteria"`
}

type CSVService struct {
    featureRepo repository.FeatureRepository
    validator   *validator.Validate
}

func (s *CSVService) ImportFeatures(projectID uint, csvData []byte) ([]domain.Feature, error) {
    var records []CSVRecord
    if err := gocsv.UnmarshalBytes(csvData, &records); err != nil {
        return nil, fmt.Errorf("invalid CSV format: %w", err)
    }
    
    var features []domain.Feature
    for i, record := range records {
        feature := domain.Feature{
            ProjectID:          projectID,
            Title:              strings.TrimSpace(record.Title),
            Description:        strings.TrimSpace(record.Description),
            AcceptanceCriteria: strings.TrimSpace(record.AcceptanceCriteria),
        }
        
        if err := feature.Validate(); err != nil {
            return nil, fmt.Errorf("row %d validation error: %w", i+1, err)
        }
        
        features = append(features, feature)
    }
    
    return features, nil
}
```

### Batch Operations Pattern
```go
// internal/repository/feature_repository.go
func (r *featureRepository) CreateBatch(features []domain.Feature) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        for _, feature := range features {
            if err := tx.Create(&feature).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## File Upload Handling

### Multipart File Processing
```go
// internal/api/feature_handler.go
func (h *FeatureHandler) ImportCSV(c *gin.Context) {
    projectID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
    
    file, header, err := c.Request.FormFile("csv_file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file provided"})
        return
    }
    defer file.Close()
    
    // Validate file type
    if !strings.HasSuffix(header.Filename, ".csv") {
        c.JSON(400, gin.H{"error": "File must be CSV format"})
        return
    }
    
    // Limit file size (5MB)
    if header.Size > 5*1024*1024 {
        c.JSON(400, gin.H{"error": "File size too large (max 5MB)"})
        return
    }
    
    csvData, err := io.ReadAll(file)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to read file"})
        return
    }
    
    features, err := h.csvService.ImportFeatures(uint(projectID), csvData)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"imported": len(features), "features": features})
}
```

## Data Export Patterns

### CSV Export with Proper Headers
```go
// internal/service/export_service.go
func (s *ExportService) ExportFeaturesToCSV(projectID uint) ([]byte, error) {
    features, err := s.featureRepo.GetByProjectID(projectID)
    if err != nil {
        return nil, err
    }
    
    var records []CSVRecord
    for _, feature := range features {
        records = append(records, CSVRecord{
            Title:              feature.Title,
            Description:        feature.Description,
            AcceptanceCriteria: feature.AcceptanceCriteria,
        })
    }
    
    var buf bytes.Buffer
    if err := gocsv.Marshal(records, &buf); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}
```

## REST API Patterns

### Resource-based URL Structure
```go
// Routes for features within projects
func SetupFeatureRoutes(r *gin.Engine, handler *FeatureHandler) {
    projects := r.Group("/api/projects/:id")
    {
        features := projects.Group("/features")
        {
            features.POST("", handler.CreateFeature)
            features.GET("", handler.ListFeatures)
            features.GET("/:feature_id", handler.GetFeature)
            features.PUT("/:feature_id", handler.UpdateFeature)
            features.DELETE("/:feature_id", handler.DeleteFeature)
            features.POST("/import", handler.ImportCSV)
            features.GET("/export", handler.ExportCSV)
        }
    }
}
```

### Consistent Response Patterns
```go
// Standard response structures
type FeatureResponse struct {
    ID                 uint   `json:"id"`
    Title              string `json:"title"`
    Description        string `json:"description"`
    AcceptanceCriteria string `json:"acceptance_criteria"`
    CreatedAt          string `json:"created_at"`
}

type ListFeaturesResponse struct {
    Features []FeatureResponse `json:"features"`
    Count    int               `json:"count"`
}

// Response helper
func (h *FeatureHandler) respondWithFeature(c *gin.Context, feature *domain.Feature) {
    response := FeatureResponse{
        ID:                 feature.ID,
        Title:              feature.Title,
        Description:        feature.Description,
        AcceptanceCriteria: feature.AcceptanceCriteria,
        CreatedAt:          feature.CreatedAt.Format(time.RFC3339),
    }
    c.JSON(200, response)
}
```

## Validation and Error Handling

### Business Rule Validation
```go
// internal/service/feature_service.go
func (s *FeatureService) ValidateFeatureForScoring(feature *domain.Feature) error {
    if len(strings.TrimSpace(feature.Title)) == 0 {
        return errors.New("feature title is required for scoring")
    }
    
    if len(strings.TrimSpace(feature.Description)) < 10 {
        return errors.New("feature description must be at least 10 characters for meaningful comparison")
    }
    
    return nil
}

func (s *FeatureService) ValidateProjectReadyForScoring(projectID uint) error {
    features, err := s.featureRepo.GetByProjectID(projectID)
    if err != nil {
        return err
    }
    
    if len(features) < 2 {
        return errors.New("minimum 2 features required for pairwise comparison")
    }
    
    for _, feature := range features {
        if err := s.ValidateFeatureForScoring(&feature); err != nil {
            return fmt.Errorf("feature '%s': %w", feature.Title, err)
        }
    }
    
    return nil
}
```

## Database Query Optimization

### Efficient Queries for Feature Operations
```go
// internal/repository/feature_repository.go
func (r *featureRepository) GetByProjectIDWithCounts(projectID uint) ([]domain.Feature, error) {
    var features []domain.Feature
    
    // Use preloading to avoid N+1 queries
    return features, r.db.
        Where("project_id = ?", projectID).
        Order("created_at ASC").
        Find(&features).Error
}

// Get features with their scoring status
func (r *featureRepository) GetWithScoringStatus(projectID uint) ([]FeatureWithStatus, error) {
    var results []FeatureWithStatus
    
    query := `
        SELECT f.*, 
               COALESCE(pairwise_count.value_completed, 0) as value_pairwise_completed,
               COALESCE(pairwise_count.complexity_completed, 0) as complexity_pairwise_completed,
               COALESCE(fib_count.value_scored, 0) as value_fibonacci_scored,
               COALESCE(fib_count.complexity_scored, 0) as complexity_fibonacci_scored
        FROM features f
        LEFT JOIN (/* subquery for pairwise completion counts */) pairwise_count ON f.id = pairwise_count.feature_id
        LEFT JOIN (/* subquery for fibonacci scoring counts */) fib_count ON f.id = fib_count.feature_id
        WHERE f.project_id = ?
    `
    
    return results, r.db.Raw(query, projectID).Scan(&results).Error
}
```

## Testing Strategies

### Feature Service Testing
```go
// internal/service/feature_service_test.go
func TestFeatureService_CreateFeature(t *testing.T) {
    tests := []struct {
        name    string
        request CreateFeatureRequest
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid feature",
            request: CreateFeatureRequest{
                Title:       "User Authentication",
                Description: "Users should be able to log in with email and password",
            },
            wantErr: false,
        },
        {
            name: "title too short",
            request: CreateFeatureRequest{
                Title:       "Hi",
                Description: "Valid description here",
            },
            wantErr: true,
            errMsg:  "title must be at least 3 characters",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := setupFeatureService()
            _, err := service.CreateFeature(1, &tt.request)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```