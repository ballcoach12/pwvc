package api

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"pairwise/internal/domain"

	"github.com/gin-gonic/gin"
)

// GetAuditReport handles GET /api/projects/:id/audit (T045 - US9)
func (h *Handler) GetAuditReport(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T045 - US9)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	// Parse query parameters for report options
	var options domain.AuditReportOptions

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			options.Limit = limit
		}
	}

	options.ActionType = c.Query("action_type")
	options.IncludePersonalData = c.Query("include_personal_data") == "true"

	// Generate audit report
	report, err := h.auditService.GetAuditReport(projectID, options)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate audit report",
			})
		}
		return
	}

	c.JSON(http.StatusOK, report)
}

// ExportAuditReport handles GET /api/projects/:id/audit/export (T045 - US9)
func (h *Handler) ExportAuditReport(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T045 - US9)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	// Parse query parameters for report options
	var options domain.AuditReportOptions

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			options.Limit = limit
		}
	}

	options.ActionType = c.Query("action_type")
	options.IncludePersonalData = c.Query("include_personal_data") == "true"

	// Generate audit report
	report, err := h.auditService.GetAuditReport(projectID, options)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate audit report",
			})
		}
		return
	}

	// Set headers for file download
	format := c.DefaultQuery("format", "json")

	switch format {
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=audit_report.csv")

		// Convert to CSV format with deterministic ordering and proper escaping (T050)
		csvContent := "timestamp,action,entity_type,entity_id,attendee_id,project_id\n"

		// Sort audit logs by timestamp for deterministic ordering
		sortedLogs := make([]domain.AuditLog, len(report.AuditLogs))
		copy(sortedLogs, report.AuditLogs)
		sort.Slice(sortedLogs, func(i, j int) bool {
			return sortedLogs[i].Timestamp.Before(sortedLogs[j].Timestamp)
		})

		for _, log := range sortedLogs {
			// Use ISO 8601 format for timestamps (locale-independent)
			csvContent += log.Timestamp.UTC().Format("2006-01-02T15:04:05Z") + ","
			csvContent += h.escapeCSVField(log.Action) + ","
			csvContent += h.escapeCSVField(log.EntityType) + ","
			csvContent += h.escapeCSVField(log.EntityID) + ","
			csvContent += strconv.Itoa(log.AttendeeID) + ","
			csvContent += strconv.Itoa(log.ProjectID) + "\n"
		}

		c.String(http.StatusOK, csvContent)

	default: // JSON format
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=audit_report.json")
		c.JSON(http.StatusOK, report)
	}
}

// GetAuditStatistics handles GET /api/projects/:id/audit/statistics (T045 - US9)
func (h *Handler) GetAuditStatistics(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Require facilitator authorization (T045 - US9)
	if !h.checkIsFacilitator(c) {
		return // checkIsFacilitator already sends the error response
	}

	// Generate basic report to get statistics
	report, err := h.auditService.GetAuditReport(projectID, domain.AuditReportOptions{
		Limit:               1000, // Get reasonable sample for statistics
		IncludePersonalData: false,
	})
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate audit statistics",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_id":    projectID,
		"project_name":  report.ProjectName,
		"total_actions": report.TotalActions,
		"date_range":    report.DateRange,
		"statistics":    report.Statistics,
		"privacy_mode":  report.PrivacyMode,
		"generated_at":  report.GeneratedAt,
	})
}

// escapeCSVField properly escapes CSV fields containing special characters (T050)
func (h *Handler) escapeCSVField(field string) string {
	// CSV fields containing commas, quotes, or newlines need to be quoted
	if strings.ContainsAny(field, ",\"\n\r") {
		// Escape internal quotes by doubling them, then wrap in quotes
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return field
}
