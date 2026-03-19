package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel represents the severity of log messages
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// LogFormat represents the output format for logs
type LogFormat int

const (
	LogFormatText LogFormat = iota
	LogFormatJSON
)

// Logger provides structured logging functionality
type Logger struct {
	level  LogLevel
	format LogFormat
	output *log.Logger
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Service    string                 `json:"service,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	level := LogLevelInfo
	format := LogFormatJSON

	// Read configuration from environment
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch levelStr {
		case "debug":
			level = LogLevelDebug
		case "info":
			level = LogLevelInfo
		case "warn", "warning":
			level = LogLevelWarn
		case "error":
			level = LogLevelError
		case "fatal":
			level = LogLevelFatal
		}
	}

	if formatStr := os.Getenv("LOG_FORMAT"); formatStr != "" {
		switch formatStr {
		case "text":
			format = LogFormatText
		case "json":
			format = LogFormatJSON
		}
	}

	return &Logger{
		level:  level,
		format: format,
		output: log.New(os.Stdout, "", 0),
	}
}

// LoggingMiddleware creates a Gin middleware for request logging
func (l *Logger) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Create log entry
		entry := LogEntry{
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Level:      "info",
			Message:    "HTTP request",
			Service:    "pwvc",
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			Duration:   duration.String(),
		}

		// Add request ID if available
		if requestID, exists := c.Get("request_id"); exists {
			if id, ok := requestID.(string); ok {
				entry.RequestID = id
			}
		}

		// Add user ID if available (from auth middleware)
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(string); ok {
				entry.UserID = id
			}
		}

		// Add error information if request failed
		if len(c.Errors) > 0 {
			entry.Level = "error"
			entry.Error = c.Errors.String()
		}

		// Add query parameters for GET requests
		if c.Request.Method == "GET" && len(c.Request.URL.RawQuery) > 0 {
			entry.Context = map[string]interface{}{
				"query": c.Request.URL.RawQuery,
			}
		}

		// Log the entry
		l.logEntry(entry)
	}
}

// Info logs an info-level message
func (l *Logger) Info(message string, context ...map[string]interface{}) {
	if l.level <= LogLevelInfo {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "info",
			Message:   message,
			Service:   "pwvc",
		}

		if len(context) > 0 {
			entry.Context = context[0]
		}

		l.logEntry(entry)
	}
}

// Warn logs a warning-level message
func (l *Logger) Warn(message string, context ...map[string]interface{}) {
	if l.level <= LogLevelWarn {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "warn",
			Message:   message,
			Service:   "pwvc",
		}

		if len(context) > 0 {
			entry.Context = context[0]
		}

		l.logEntry(entry)
	}
}

// Error logs an error-level message
func (l *Logger) Error(message string, err error, context ...map[string]interface{}) {
	if l.level <= LogLevelError {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "error",
			Message:   message,
			Service:   "pwvc",
		}

		if err != nil {
			entry.Error = err.Error()
		}

		if len(context) > 0 {
			entry.Context = context[0]
		}

		l.logEntry(entry)
	}
}

// Debug logs a debug-level message
func (l *Logger) Debug(message string, context ...map[string]interface{}) {
	if l.level <= LogLevelDebug {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "debug",
			Message:   message,
			Service:   "pwvc",
		}

		if len(context) > 0 {
			entry.Context = context[0]
		}

		l.logEntry(entry)
	}
}

// Fatal logs a fatal-level message and exits
func (l *Logger) Fatal(message string, err error, context ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "fatal",
		Message:   message,
		Service:   "pwvc",
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if len(context) > 0 {
		entry.Context = context[0]
	}

	l.logEntry(entry)
	os.Exit(1)
}

// logEntry outputs the log entry based on the configured format
func (l *Logger) logEntry(entry LogEntry) {
	switch l.format {
	case LogFormatJSON:
		if jsonBytes, err := json.Marshal(entry); err == nil {
			l.output.Println(string(jsonBytes))
		} else {
			// Fallback to text format if JSON marshaling fails
			l.logTextEntry(entry)
		}
	case LogFormatText:
		l.logTextEntry(entry)
	}
}

// logTextEntry outputs the log entry in text format
func (l *Logger) logTextEntry(entry LogEntry) {
	output := fmt.Sprintf("[%s] %s: %s", entry.Timestamp, entry.Level, entry.Message)

	if entry.Method != "" && entry.Path != "" {
		output += fmt.Sprintf(" %s %s", entry.Method, entry.Path)
	}

	if entry.StatusCode != 0 {
		output += fmt.Sprintf(" status=%d", entry.StatusCode)
	}

	if entry.Duration != "" {
		output += fmt.Sprintf(" duration=%s", entry.Duration)
	}

	if entry.Error != "" {
		output += fmt.Sprintf(" error=%q", entry.Error)
	}

	if entry.RequestID != "" {
		output += fmt.Sprintf(" request_id=%s", entry.RequestID)
	}

	l.output.Println(output)
}

// Performance monitoring middleware
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		// Log slow requests (> 1 second)
		if duration > time.Second {
			log.Printf("SLOW REQUEST: %s %s took %v",
				c.Request.Method,
				c.Request.URL.Path,
				duration)
		}

		// Add performance headers
		c.Header("X-Response-Time", duration.String())
	}
}

// Enhanced rate limiting middleware with endpoint-specific limits (T047)
func RateLimitMiddleware() gin.HandlerFunc {
	// In-memory rate limiter with endpoint-specific limits and backoff
	// In production, use Redis or a proper rate limiting library

	requests := make(map[string][]time.Time)
	violations := make(map[string]int) // Track violations for progressive penalties

	// Define rate limits per endpoint pattern
	rateLimits := map[string]RateLimit{
		"/api/projects/*/pairwise/votes":  {60, time.Minute},  // Hot endpoint - voting
		"/api/projects/*/scores/*":        {120, time.Minute}, // Hot endpoint - scoring
		"/api/projects/*/consensus/*":     {30, time.Minute},  // Consensus operations
		"/api/projects/*/progress/*":      {10, time.Minute},  // Phase changes (facilitator)
		"/api/projects/*/audit/*":         {20, time.Minute},  // Audit operations (facilitator)
		"/api/projects/*/results/export":  {5, time.Minute},   // Export operations
		"/api/projects/*/features/import": {3, time.Minute},   // Import operations
		"default":                         {100, time.Minute}, // Default limit
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.Request.URL.Path
		now := time.Now()

		// Find applicable rate limit
		limit := getRateLimitForPath(path, rateLimits)

		// Apply progressive penalty based on violations
		violationCount := violations[clientIP]
		adjustedLimit := limit.MaxRequests
		if violationCount > 0 {
			// Reduce limit by 20% for each violation (progressive backoff)
			penaltyFactor := 1.0 - (0.2 * float64(violationCount))
			if penaltyFactor < 0.1 { // Don't go below 10% of original limit
				penaltyFactor = 0.1
			}
			adjustedLimit = int(float64(limit.MaxRequests) * penaltyFactor)
		}

		// Clean old requests
		key := clientIP + ":" + path
		if times, exists := requests[key]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < limit.TimeWindow {
					validTimes = append(validTimes, t)
				}
			}
			requests[key] = validTimes
		}

		// Check rate limit
		if times, exists := requests[key]; exists && len(times) >= adjustedLimit {
			// Increment violation count
			violations[clientIP]++

			// Calculate retry-after with exponential backoff
			retryAfter := time.Duration(violationCount+1) * limit.TimeWindow
			if retryAfter > 5*time.Minute {
				retryAfter = 5 * time.Minute // Cap at 5 minutes
			}

			c.Header("X-RateLimit-Limit", strconv.Itoa(limit.MaxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(now.Add(retryAfter).Unix(), 10))
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))

			c.JSON(429, gin.H{
				"error":           "Rate limit exceeded",
				"retry_after":     retryAfter.Seconds(),
				"limit":           limit.MaxRequests,
				"window":          limit.TimeWindow.String(),
				"violation_count": violationCount,
				"message":         "Rate limit violations result in progressive penalties. Please implement backoff strategies.",
			})
			c.Abort()
			return
		}

		// Add current request
		requests[key] = append(requests[key], now)

		// Set rate limit headers for successful requests
		remaining := adjustedLimit - len(requests[key])
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now.Add(limit.TimeWindow).Unix(), 10))

		c.Next()
	}
}

// RateLimit defines rate limiting parameters
type RateLimit struct {
	MaxRequests int
	TimeWindow  time.Duration
}

// getRateLimitForPath matches request path to appropriate rate limit
func getRateLimitForPath(path string, limits map[string]RateLimit) RateLimit {
	// Simple pattern matching - in production use a router or regex
	for pattern, limit := range limits {
		if pattern == "default" {
			continue
		}
		// Basic wildcard matching for /* patterns
		if matchesPattern(path, pattern) {
			return limit
		}
	}
	return limits["default"]
}

// matchesPattern performs basic wildcard matching
func matchesPattern(path, pattern string) bool {
	// Replace */ with a regex-like pattern and check
	// This is a simplified implementation
	if strings.Contains(pattern, "*/") {
		// Split pattern and check segments
		patternParts := strings.Split(pattern, "/")
		pathParts := strings.Split(path, "/")

		if len(patternParts) != len(pathParts) {
			return false
		}

		for i, part := range patternParts {
			if part != "*" && part != pathParts[i] {
				return false
			}
		}
		return true
	}
	return path == pattern
}
