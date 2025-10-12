package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"pwvc/internal/domain"

	"github.com/gin-gonic/gin"
)

// CreateFeature handles POST /api/projects/:id/features
func (h *Handler) CreateFeature(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req domain.CreateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	feature, err := h.featureService.CreateFeature(projectID, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, feature)
}

// GetProjectFeatures handles GET /api/projects/:id/features
func (h *Handler) GetProjectFeatures(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	features, err := h.featureService.GetProjectFeatures(projectID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"features": features,
	})
}

// GetFeature handles GET /api/projects/:id/features/:feature_id
func (h *Handler) GetFeature(c *gin.Context) {
	featureID, err := strconv.Atoi(c.Param("feature_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid feature ID",
		})
		return
	}

	feature, err := h.featureService.GetFeature(featureID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, feature)
}

// UpdateFeature handles PUT /api/projects/:id/features/:feature_id
func (h *Handler) UpdateFeature(c *gin.Context) {
	featureID, err := strconv.Atoi(c.Param("feature_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid feature ID",
		})
		return
	}

	var req domain.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	feature, err := h.featureService.UpdateFeature(featureID, req)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, feature)
}

// DeleteFeature handles DELETE /api/projects/:id/features/:feature_id
func (h *Handler) DeleteFeature(c *gin.Context) {
	featureID, err := strconv.Atoi(c.Param("feature_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid feature ID",
		})
		return
	}

	err = h.featureService.DeleteFeature(featureID)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ImportFeatures handles POST /api/projects/:id/features/import
func (h *Handler) ImportFeatures(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file uploaded or invalid file",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file type
	if header.Header.Get("Content-Type") != "text/csv" && !isCSVFile(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only CSV files are supported",
		})
		return
	}

	// Import features from CSV
	result, err := h.featureService.ImportFeaturesFromCSV(projectID, file)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ExportFeatures handles GET /api/projects/:id/features/export
func (h *Handler) ExportFeatures(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var buffer bytes.Buffer
	err = h.featureService.ExportFeaturesToCSV(projectID, &buffer)
	if err != nil {
		if apiErr, ok := err.(*domain.APIError); ok {
			c.JSON(apiErr.Code, gin.H{
				"error":   apiErr.Message,
				"details": apiErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Set headers for CSV download
	filename := fmt.Sprintf("project_%d_features.csv", projectID)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", strconv.Itoa(buffer.Len()))

	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}

// isCSVFile checks if the filename has a CSV extension
func isCSVFile(filename string) bool {
	if len(filename) < 4 {
		return false
	}
	ext := filename[len(filename)-4:]
	return ext == ".csv" || ext == ".CSV"
}
