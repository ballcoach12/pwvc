package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthEndpoint tests the basic health check endpoint
func TestHealthEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pwvc",
		})
	})

	// Test
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}

	if response["service"] != "pwvc" {
		t.Errorf("Expected service 'pwvc', got %v", response["service"])
	}
}

// TestInputValidation tests basic input validation
func TestInputValidation(t *testing.T) {
	validator := NewInputValidator()

	// Test project name validation
	err := validator.validateProjectName("Valid Project Name")
	if err != nil {
		t.Errorf("Expected valid project name to pass validation: %v", err)
	}

	err = validator.validateProjectName("")
	if err == nil {
		t.Error("Expected empty project name to fail validation")
	}

	// Test email validation
	err = validator.validateEmail("test@example.com")
	if err != nil {
		t.Errorf("Expected valid email to pass validation: %v", err)
	}

	err = validator.validateEmail("invalid-email")
	if err == nil {
		t.Error("Expected invalid email to fail validation")
	}
}

// TestRequestIDMiddleware tests request ID generation
func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		if !exists {
			t.Error("Request ID should exist in context")
		}

		if requestID == "" {
			t.Error("Request ID should not be empty")
		}

		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response header
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header should be set")
	}
}

// BenchmarkHealthCheck benchmarks the health check endpoint
func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "pwvc",
		})
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
