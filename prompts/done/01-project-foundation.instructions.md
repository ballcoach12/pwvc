# Instructions: Project Foundation Implementation

## Code Organization Patterns

### Domain-Driven Design Structure
```go
// internal/domain/project.go - Core entities
type Project struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" validate:"required,min=3,max=255"`
    Description string    `json:"description" validate:"max=1000"`
    Status      string    `json:"status" gorm:"default:active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Attendees   []Attendee `json:"attendees,omitempty"`
}

type Attendee struct {
    ID            uint    `json:"id" gorm:"primaryKey"`
    ProjectID     uint    `json:"project_id"`
    Name          string  `json:"name" validate:"required,min=2,max=255"`
    Role          string  `json:"role" validate:"max=100"`
    IsFacilitator bool    `json:"is_facilitator" gorm:"default:false"`
    CreatedAt     time.Time `json:"created_at"`
}
```

### Repository Pattern Implementation
```go
// internal/repository/project_repository.go
type ProjectRepository interface {
    Create(project *domain.Project) error
    GetByID(id uint) (*domain.Project, error)
    Update(project *domain.Project) error
    Delete(id uint) error
    GetWithAttendees(id uint) (*domain.Project, error)
}

type projectRepository struct {
    db *gorm.DB
}

func (r *projectRepository) Create(project *domain.Project) error {
    return r.db.Create(project).Error
}
```

### Service Layer Pattern
```go
// internal/service/project_service.go
type ProjectService interface {
    CreateProject(req *CreateProjectRequest) (*domain.Project, error)
    AddAttendee(projectID uint, req *AddAttendeeRequest) (*domain.Attendee, error)
    ValidateProjectForScoring(projectID uint) error
}

func (s *projectService) ValidateProjectForScoring(projectID uint) error {
    project, err := s.repo.GetWithAttendees(projectID)
    if err != nil {
        return err
    }
    
    if len(project.Attendees) < 2 {
        return errors.New("minimum 2 attendees required for scoring")
    }
    
    return nil
}
```

## Database Migration Patterns

### Migration File Structure
```sql
-- migrations/001_create_projects.up.sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at);
```

### Migration Runner Pattern
```go
// pkg/database/migrate.go
func RunMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }
    
    m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
    if err != nil {
        return err
    }
    
    return m.Up()
}
```

## Error Handling Patterns

### Custom Error Types
```go
// pkg/errors/errors.go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type APIError struct {
    Code    int                `json:"code"`
    Message string            `json:"message"`
    Errors  []ValidationError `json:"errors,omitempty"`
}

func NewValidationError(field, message string) *APIError {
    return &APIError{
        Code:    400,
        Message: "Validation failed",
        Errors:  []ValidationError{{Field: field, Message: message}},
    }
}
```

### Middleware Error Handler
```go
// pkg/middleware/error_handler.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *APIError:
                c.JSON(e.Code, e)
            case validator.ValidationErrors:
                apiErr := handleValidationErrors(e)
                c.JSON(400, apiErr)
            default:
                c.JSON(500, APIError{Code: 500, Message: "Internal server error"})
            }
        }
    }
}
```

## Testing Patterns

### Repository Testing with Test Database
```go
// internal/repository/project_repository_test.go
func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&domain.Project{}, &domain.Attendee{})
    return db
}

func TestProjectRepository_Create(t *testing.T) {
    db := setupTestDB()
    repo := NewProjectRepository(db)
    
    project := &domain.Project{
        Name:        "Test Project",
        Description: "Test Description",
    }
    
    err := repo.Create(project)
    assert.NoError(t, err)
    assert.NotZero(t, project.ID)
}
```

### Service Layer Testing with Mocks
```go
// internal/service/project_service_test.go
type mockProjectRepository struct {
    projects map[uint]*domain.Project
}

func (m *mockProjectRepository) Create(project *domain.Project) error {
    project.ID = uint(len(m.projects) + 1)
    m.projects[project.ID] = project
    return nil
}

func TestProjectService_CreateProject(t *testing.T) {
    repo := &mockProjectRepository{projects: make(map[uint]*domain.Project)}
    service := NewProjectService(repo)
    
    req := &CreateProjectRequest{
        Name:        "Test Project",
        Description: "Test Description",
    }
    
    project, err := service.CreateProject(req)
    assert.NoError(t, err)
    assert.Equal(t, "Test Project", project.Name)
}
```

## Configuration Management

### Environment-based Configuration
```go
// pkg/config/config.go
type Config struct {
    Database DatabaseConfig `mapstructure:"database"`
    Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host" default:"localhost"`
    Port     int    `mapstructure:"port" default:"5432"`
    Name     string `mapstructure:"name" required:"true"`
    User     string `mapstructure:"user" required:"true"`
    Password string `mapstructure:"password" required:"true"`
}

func Load() (*Config, error) {
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    var cfg Config
    return &cfg, viper.Unmarshal(&cfg)
}
```

## Security Considerations

### Input Validation
```go
// Use struct tags for validation
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=3,max=255"`
    Description string `json:"description" validate:"max=1000"`
}

// Validate in service layer
func (s *projectService) CreateProject(req *CreateProjectRequest) (*domain.Project, error) {
    if err := s.validator.Struct(req); err != nil {
        return nil, NewValidationError("input", err.Error())
    }
    // ... continue with creation
}
```

### SQL Injection Prevention
```go
// Always use GORM or prepared statements
// GOOD:
db.Where("id = ?", userID).First(&project)

// BAD:
db.Where(fmt.Sprintf("id = %d", userID)).First(&project)
```