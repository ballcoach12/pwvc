package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// Rate limiting middleware (basic implementation)
func RateLimitMiddleware() gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, use Redis or a proper rate limiting library

	requests := make(map[string][]time.Time)
	const maxRequests = 100
	const timeWindow = time.Minute

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old requests
		if times, exists := requests[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < timeWindow {
					validTimes = append(validTimes, t)
				}
			}
			requests[clientIP] = validTimes
		}

		// Check rate limit
		if times, exists := requests[clientIP]; exists && len(times) >= maxRequests {
			c.JSON(429, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": timeWindow.Seconds(),
			})
			c.Abort()
			return
		}

		// Add current request
		requests[clientIP] = append(requests[clientIP], now)

		c.Next()
	}
}
