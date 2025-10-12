package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pwvc/internal/domain"

	"github.com/gin-gonic/gin"
)

// CalculateResults handles POST /api/projects/{id}/calculate-results
func (h *Handler) CalculateResults(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	results, err := h.resultsService.CalculateResults(projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetResults handles GET /api/projects/{id}/results
func (h *Handler) GetResults(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	results, err := h.resultsService.GetResults(projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

// ExportResults handles GET /api/projects/{id}/results/export?format=csv|json|jira
func (h *Handler) ExportResults(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	formatStr := c.Query("format")
	if formatStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Export format is required. Use format=csv|json|jira",
		})
		return
	}

	format := domain.ExportFormat(strings.ToLower(formatStr))

	// Validate format
	validFormats := map[domain.ExportFormat]bool{
		domain.ExportFormatCSV:  true,
		domain.ExportFormatJSON: true,
		domain.ExportFormatJira: true,
	}

	if !validFormats[format] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid format. Must be csv, json, or jira",
		})
		return
	}

	data, err := h.resultsService.ExportResults(projectID, format)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Set appropriate headers and content type based on format
	switch format {
	case domain.ExportFormatCSV:
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=pwvc_results_%d.csv", projectID))

		// Convert CSV data to string
		csvData := data.([][]string)
		var csvString strings.Builder
		writer := csv.NewWriter(&csvString)
		writer.WriteAll(csvData)
		writer.Flush()

		c.String(http.StatusOK, csvString.String())

	case domain.ExportFormatJSON:
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=pwvc_results_%d.json", projectID))
		c.JSON(http.StatusOK, data)

	case domain.ExportFormatJira:
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=pwvc_jira_export_%d.json", projectID))
		c.JSON(http.StatusOK, data)
	}
}

// GetResultsSummary handles GET /api/projects/{id}/results/summary
func (h *Handler) GetResultsSummary(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	results, err := h.resultsService.GetResults(projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Return just the summary portion
	summary := map[string]interface{}{
		"projectId":     results.ProjectID,
		"totalFeatures": results.TotalFeatures,
		"calculatedAt":  results.CalculatedAt,
		"summary":       results.Summary,
		"topFeatures":   getTopFeatures(results.Results, 5),
	}

	c.JSON(http.StatusOK, summary)
}

// CheckResultsStatus handles GET /api/projects/{id}/results/status
func (h *Handler) CheckResultsStatus(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	exists, err := h.priorityRepo.ExistsForProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check results status",
		})
		return
	}

	status := map[string]interface{}{
		"projectId":      projectID,
		"hasResults":     exists,
		"calculationUrl": fmt.Sprintf("/api/projects/%d/calculate-results", projectID),
		"resultsUrl":     fmt.Sprintf("/api/projects/%d/results", projectID),
	}

	if exists {
		latest, err := h.priorityRepo.GetLatestCalculationTime(projectID)
		if err == nil {
			status["lastCalculated"] = latest.CalculatedAt
			status["lastCalculationId"] = latest.ID
		}
	}

	c.JSON(http.StatusOK, status)
}

// PreviewExport handles GET /api/projects/{id}/results/preview?format=csv|json|jira&limit=5
func (h *Handler) PreviewExport(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	formatStr := c.Query("format")
	if formatStr == "" {
		formatStr = "json"
	}

	limitStr := c.Query("limit")
	limit := 5 // default limit for preview
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 20 {
			limit = parsedLimit
		}
	}

	results, err := h.resultsService.GetResults(projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Limit results for preview
	if len(results.Results) > limit {
		results.Results = results.Results[:limit]
		results.TotalFeatures = limit
	}

	format := domain.ExportFormat(strings.ToLower(formatStr))
	data, err := h.resultsService.ExportResults(projectID, format)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	preview := map[string]interface{}{
		"format":       format,
		"previewLimit": limit,
		"totalResults": len(results.Results),
		"data":         data,
		"note":         fmt.Sprintf("Preview showing first %d results. Full export will include all features.", limit),
	}

	c.JSON(http.StatusOK, preview)
}

// getTopFeatures returns the top N features from results
func getTopFeatures(results []domain.PriorityResult, n int) []domain.PriorityResult {
	if len(results) <= n {
		return results
	}
	return results[:n]
}
