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
		priorityRepo,
		attendeeRepo,
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
	)
	if err != nil {
		return nil, err
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

	// Initialize logger
	logger := api.NewLogger()

	// Add middleware
	router.Use(logger.LoggingMiddleware())  // Use structured logging
	router.Use(api.RecoveryMiddleware())    // Use custom recovery middleware
	router.Use(api.RequestIDMiddleware())   // Add request ID tracking
	router.Use(api.PerformanceMiddleware()) // Add performance monitoring
	router.Use(corsMiddleware())
	router.Use(api.ValidationMiddleware()) // Add validation middleware
	router.Use(api.RateLimitMiddleware())  // Add basic rate limiting

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PairWise - Feature Prioritization Through Group Consensus",
			"version": "0.1.0",
		})
	})

	// API routes
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
