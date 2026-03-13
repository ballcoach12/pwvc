package api

import (
	"log"
	"net/http"
	"runtime/debug"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// ErrorHandler provides centralized error handling for the API
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError processes errors and returns appropriate HTTP responses
func (eh *ErrorHandler) HandleError(c *gin.Context, err error) {
	// Log the error for debugging
	log.Printf("API Error: %v", err)

	switch e := err.(type) {
	case *domain.APIError:
		c.JSON(e.Code, gin.H{
			"error":   e.Message,
			"details": e.Details,
		})
		return

	case *domain.BusinessError:
		statusCode := mapBusinessErrorToHTTP(e.Type)
		c.JSON(statusCode, gin.H{
			"error":   e.Message,
			"type":    e.Type,
			"context": e.Context,
		})
		return

	case *domain.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"field":   e.Field,
			"message": e.Message,
		})
		return

	case *pq.Error:
		// Handle PostgreSQL errors
		statusCode, message := handlePostgreSQLError(e)
		c.JSON(statusCode, gin.H{
			"error": message,
		})
		return

	default:
		// Handle standard errors
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Resource not found",
			})
		case domain.ErrValidation:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request data",
			})
		case domain.ErrDuplicate:
			c.JSON(http.StatusConflict, gin.H{
				"error": "Resource already exists",
			})
		case domain.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized access",
			})
		case domain.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden action",
			})
		case domain.ErrTimeout:
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Request timeout",
			})
		case domain.ErrRateLimit:
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
		case domain.ErrServiceUnavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service temporarily unavailable",
			})
		default:
			// Unknown error - log full details and return generic message
			log.Printf("Unhandled error: %v\nStack trace: %s", err, debug.Stack())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An unexpected error occurred",
			})
		}
	}
}

// mapBusinessErrorToHTTP maps business error types to HTTP status codes
func mapBusinessErrorToHTTP(errorType string) int {
	switch errorType {
	case domain.ErrTypeInsufficientData:
		return http.StatusBadRequest
	case domain.ErrTypeInvalidPhase:
		return http.StatusBadRequest
	case domain.ErrTypePhaseNotComplete:
		return http.StatusPreconditionFailed
	case domain.ErrTypeSessionInProgress:
		return http.StatusConflict
	case domain.ErrTypeDataInconsistency:
		return http.StatusInternalServerError
	case domain.ErrTypeWorkflowViolation:
		return http.StatusBadRequest
	case domain.ErrTypeResourceConflict:
		return http.StatusConflict
	case domain.ErrTypeInvalidTransition:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// handlePostgreSQLError converts PostgreSQL errors to user-friendly messages
func handlePostgreSQLError(pgErr *pq.Error) (int, string) {
	switch pgErr.Code {
	case "23505": // unique_violation
		return http.StatusConflict, "Resource already exists"
	case "23503": // foreign_key_violation
		return http.StatusBadRequest, "Referenced resource does not exist"
	case "23502": // not_null_violation
		return http.StatusBadRequest, "Required field is missing"
	case "23514": // check_violation
		return http.StatusBadRequest, "Data violates constraints"
	case "42P01": // undefined_table
		return http.StatusInternalServerError, "Database configuration error"
	case "42703": // undefined_column
		return http.StatusInternalServerError, "Database schema error"
	case "08003", "08006": // connection errors
		return http.StatusServiceUnavailable, "Database connection error"
	case "57014": // query_canceled
		return http.StatusRequestTimeout, "Query timeout"
	default:
		log.Printf("PostgreSQL error: Code=%s, Message=%s", pgErr.Code, pgErr.Message)
		return http.StatusInternalServerError, "Database error"
	}
}

// RecoveryMiddleware recovers from panics and returns structured error responses
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\nStack trace: %s", err, debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ValidationMiddleware validates request data and returns structured errors
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were binding errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
			c.Abort()
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or get request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID creates a unique request identifier
func generateRequestID() string {
	// Simple implementation - in production, use UUID or similar
	return "req_" + string(debug.Stack()[:8])
}
