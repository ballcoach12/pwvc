package main

import (
	"log"
	"net/http"
	"os"

	"pairwise/internal/api"
	"pairwise/internal/domain"
	"pairwise/internal/repository"
	"pairwise/internal/service"
	"pairwise/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying SQL DB for repositories that need it
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories (using interfaces for proper dependency injection)
	var projectRepo repository.ProjectRepository = repository.NewProjectRepository(sqlDB)
	var attendeeRepo repository.AttendeeRepository = repository.NewAttendeeRepository(sqlDB)
	var featureRepo repository.FeatureRepository = repository.NewFeatureRepository(sqlDB)
	var pairwiseRepo repository.PairwiseRepository = repository.NewPairwiseRepository(sqlDB)
	var priorityRepo repository.PriorityRepository = repository.NewPriorityRepository(sqlDB)
	var progressRepo repository.ProgressRepository = repository.NewProgressRepository(sqlDB)
	var auditRepo repository.AuditRepository = repository.NewAuditRepository(db)             // Uses GORM
	var scoringRepo repository.ScoringRepository = repository.NewScoringRepository(db)       // Uses GORM
	var consensusRepo repository.ConsensusRepository = repository.NewConsensusRepository(db) // Uses GORM

	// Initialize JWT authentication repositories and services
	var userRepo repository.UserRepository = repository.NewUserRepository(db)
	var roleRepo repository.RoleRepository = repository.NewRoleRepository(db)
	jwtService := service.NewJWTService()
	authService := service.NewAuthService(userRepo, jwtService)
	_ = service.NewUserService(userRepo, roleRepo, authService) // TODO: Will be used for admin endpoints

	// Initialize services
	projectService := service.NewProjectService(projectRepo)
	attendeeService := service.NewAttendeeService(attendeeRepo)
	featureService := service.NewFeatureService(featureRepo, projectRepo)
	pairwiseService := service.NewPairwiseService(pairwiseRepo, featureRepo, attendeeRepo, projectRepo)
	pairwiseCalcService := service.NewPWVCService()
	resultsService := service.NewResultsService(priorityRepo, featureRepo, pairwiseRepo, consensusRepo)
	progressService := service.NewProgressService(progressRepo, projectRepo, attendeeRepo, featureRepo, scoringRepo)
	auditService := service.NewAuditService(auditRepo, attendeeRepo, projectRepo)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub(attendeeRepo)
	go wsHub.Run() // Start the hub in a goroutine

	// Initialize Fibonacci scoring and consensus services (P2 features - T030, T034)
	scoringService := service.NewScoringService(scoringRepo, featureRepo, attendeeRepo, auditRepo)
	consensusService := service.NewConsensusService(consensusRepo, featureRepo, attendeeRepo, auditRepo)

	// Initialize API handlers
	apiHandler := api.NewHandler(
		attendeeService,
		featureService,
		projectService,
		pairwiseService,
		pairwiseCalcService,
		resultsService,
		progressService,
		scoringService,
		consensusService,
		auditService,
		authService,
		jwtService,
		priorityRepo,
		attendeeRepo,
		userRepo,
		roleRepo,
		wsHub,
	)

	// Set up Gin router
	router := setupRouter(apiHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("PairWise Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func seedRoles(db *gorm.DB) error {
	// Seed predefined roles if they don't exist
	for _, role := range domain.PredefinedRoles {
		var existingRole domain.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.Error != nil && result.Error.Error() == "record not found" {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created role: %s", role.Name)
		}
	}
	return nil
}

func createDefaultUsers(db *gorm.DB) error {
	// Initialize repositories and services for user creation
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	jwtService := service.NewJWTService()
	authService := service.NewAuthService(userRepo, jwtService)

	// Default users to create
	defaultUsers := []struct {
		username string
		password string
		email    string
		roles    []string
	}{
		{"admin", "admin123", "admin@pairwise.local", []string{"admin"}},
		{"user", "user123", "user@pairwise.local", []string{"user"}},
		{"viewer", "viewer123", "viewer@pairwise.local", []string{"viewer"}},
	}

	for _, userData := range defaultUsers {
		// Check if user already exists
		existingUser, err := userRepo.GetUserByUsername(userData.username)
		if err == nil && existingUser != nil {
			continue // User already exists, skip
		}

		// Hash password
		hashedPassword, err := authService.HashPassword(userData.password)
		if err != nil {
			log.Printf("Warning: Failed to hash password for user '%s': %v", userData.username, err)
			continue
		}

		// Create user
		user := &domain.User{
			Username:     userData.username,
			Email:        userData.email,
			PasswordHash: hashedPassword,
			IsActive:     true,
		}

		err = userRepo.CreateUser(user)
		if err != nil {
			log.Printf("Warning: Failed to create user '%s': %v", userData.username, err)
			continue
		}

		// Assign roles
		for _, roleName := range userData.roles {
			role, err := roleRepo.GetRoleByName(roleName)
			if err != nil {
				log.Printf("Warning: Role '%s' not found for user '%s'", roleName, userData.username)
				continue
			}

			err = roleRepo.AssignRoleToUser(user.ID, role.ID, nil)
			if err != nil {
				log.Printf("Warning: Failed to assign role '%s' to user '%s': %v", roleName, userData.username, err)
			}
		}

		log.Printf("Created default user: %s", userData.username)
	}

	return nil
}

func initDB() (*gorm.DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "pairwise.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&domain.Project{},
		&domain.Attendee{},
		&domain.Feature{},
		&domain.PairwiseSession{},
		&domain.SessionComparison{},
		&domain.AttendeeVote{},
		&domain.PriorityCalculation{},
		&domain.ProjectProgress{},
		&domain.FibonacciScore{}, // T030 - US4
		&domain.ConsensusScore{}, // T034 - US5
		&domain.AuditLog{},       // T043 - US9
		// JWT Authentication models
		&domain.User{},     // System users with authentication
		&domain.Role{},     // Permission roles
		&domain.UserRole{}, // User-role junction table
	)
	if err != nil {
		return nil, err
	}

	// Seed predefined roles
	err = seedRoles(db)
	if err != nil {
		log.Printf("Warning: Failed to seed roles: %v", err)
	}

	// Create default users
	err = createDefaultUsers(db)
	if err != nil {
		log.Printf("Warning: Failed to create default users: %v", err)
	}

	log.Println("Successfully connected to SQLite database")
	return db, nil
}

func setupRouter(apiHandler *api.Handler) *gin.Engine {
	// Set Gin mode from environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add basic middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PairWise - Feature Prioritization Through Group Consensus",
			"version": "0.1.0",
		})
	})

	// API routes (includes authentication routes)
	apiHandler.RegisterRoutes(router)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
