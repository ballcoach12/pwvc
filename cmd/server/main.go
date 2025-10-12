package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"pwvc/internal/api"
	"pwvc/internal/repository"
	"pwvc/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	projectRepo := repository.NewProjectRepository(db)
	attendeeRepo := repository.NewAttendeeRepository(db)
	featureRepo := repository.NewFeatureRepository(db)
	pairwiseRepo := repository.NewPairwiseRepository(db)

	// Initialize services
	projectService := service.NewProjectService(projectRepo)
	attendeeService := service.NewAttendeeService(attendeeRepo)
	featureService := service.NewFeatureService(featureRepo, projectRepo)
	pairwiseService := service.NewPairwiseService(pairwiseRepo, featureRepo, attendeeRepo, projectRepo)

	// Initialize API handlers
	apiHandler := api.NewHandler(projectService, attendeeService, featureService, pairwiseService)

	// Set up Gin router
	router := setupRouter(apiHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("P-WVC Server starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/pwvc?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func setupRouter(apiHandler *api.Handler) *gin.Engine {
	// Set Gin mode from environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pwvc",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "P-WVC (Pairwise-Weighted Value/Complexity) Model API",
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
