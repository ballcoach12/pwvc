package api

import (
	"database/sql"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler provides health check and monitoring endpoints
type HealthHandler struct {
	db        *sql.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthCheck provides a simple health check endpoint
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "healthy"
	httpStatus := http.StatusOK

	checks := map[string]interface{}{
		"database": h.checkDatabase(),
		"memory":   h.checkMemory(),
		"uptime":   h.getUptime(),
	}

	// Check if any component is unhealthy
	for _, check := range checks {
		if checkMap, ok := check.(map[string]interface{}); ok {
			if checkMap["status"] != "healthy" {
				status = "unhealthy"
				httpStatus = http.StatusServiceUnavailable
			}
		}
	}

	response := gin.H{
		"status":    status,
		"service":   "pwvc",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
	}

	c.JSON(httpStatus, response)
}

// DetailedHealthCheck provides comprehensive health information
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "pwvc",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    h.getUptime(),
		"system":    h.getSystemInfo(),
		"database":  h.getDatabaseInfo(),
		"runtime":   h.getRuntimeInfo(),
	}

	c.JSON(http.StatusOK, response)
}

// Readiness check - determines if service is ready to accept traffic
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ready := true
	checks := map[string]interface{}{
		"database": h.checkDatabase(),
	}

	// Check if all critical components are ready
	for _, check := range checks {
		if checkMap, ok := check.(map[string]interface{}); ok {
			if checkMap["status"] != "healthy" {
				ready = false
				break
			}
		}
	}

	status := "ready"
	httpStatus := http.StatusOK

	if !ready {
		status = "not_ready"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
	})
}

// Liveness check - determines if service is alive and should not be restarted
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// For now, if we can respond, we're alive
	// In a more complex system, you might check for deadlocks, etc.

	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    h.getUptime(),
	})
}

// Metrics endpoint for Prometheus or other monitoring systems
func (h *HealthHandler) Metrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbStats := h.db.Stats()

	uptime := time.Since(h.startTime).Seconds()

	metrics := gin.H{
		"system_uptime_seconds":              uptime,
		"go_goroutines":                      runtime.NumGoroutine(),
		"go_memory_alloc_bytes":              m.Alloc,
		"go_memory_total_alloc_bytes":        m.TotalAlloc,
		"go_memory_sys_bytes":                m.Sys,
		"go_memory_heap_alloc_bytes":         m.HeapAlloc,
		"go_memory_heap_sys_bytes":           m.HeapSys,
		"go_memory_heap_inuse_bytes":         m.HeapInuse,
		"go_memory_stack_inuse_bytes":        m.StackInuse,
		"go_memory_stack_sys_bytes":          m.StackSys,
		"go_gc_runs_total":                   m.NumGC,
		"database_open_connections":          dbStats.OpenConnections,
		"database_in_use_connections":        dbStats.InUse,
		"database_idle_connections":          dbStats.Idle,
		"database_wait_count_total":          dbStats.WaitCount,
		"database_wait_duration_seconds":     dbStats.WaitDuration.Seconds(),
		"database_max_idle_closed_total":     dbStats.MaxIdleClosed,
		"database_max_lifetime_closed_total": dbStats.MaxLifetimeClosed,
	}

	c.JSON(http.StatusOK, metrics)
}

// Helper methods

func (h *HealthHandler) checkDatabase() map[string]interface{} {
	start := time.Now()

	err := h.db.Ping()
	duration := time.Since(start)

	if err != nil {
		return map[string]interface{}{
			"status":      "unhealthy",
			"error":       err.Error(),
			"duration_ms": duration.Milliseconds(),
		}
	}

	return map[string]interface{}{
		"status":      "healthy",
		"duration_ms": duration.Milliseconds(),
		"stats":       h.db.Stats(),
	}
}

func (h *HealthHandler) checkMemory() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Consider memory unhealthy if we're using more than 1GB
	status := "healthy"
	if m.Alloc > 1024*1024*1024 {
		status = "warning"
	}

	return map[string]interface{}{
		"status":         status,
		"alloc_mb":       bToMb(m.Alloc),
		"total_alloc_mb": bToMb(m.TotalAlloc),
		"sys_mb":         bToMb(m.Sys),
		"num_gc":         m.NumGC,
	}
}

func (h *HealthHandler) getUptime() string {
	uptime := time.Since(h.startTime)
	return uptime.String()
}

func (h *HealthHandler) getSystemInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
		"cpu_count":  runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
	}
}

func (h *HealthHandler) getDatabaseInfo() map[string]interface{} {
	stats := h.db.Stats()

	return map[string]interface{}{
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}

func (h *HealthHandler) getRuntimeInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"alloc_mb":       bToMb(m.Alloc),
		"total_alloc_mb": bToMb(m.TotalAlloc),
		"sys_mb":         bToMb(m.Sys),
		"heap_alloc_mb":  bToMb(m.HeapAlloc),
		"heap_sys_mb":    bToMb(m.HeapSys),
		"num_gc":         m.NumGC,
		"gc_pause_ns":    m.PauseNs,
	}
}

// Utility function to convert bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
