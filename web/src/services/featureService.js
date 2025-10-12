import api, { handleApiError, handleApiResponse } from './api.js'

export const featureService = {
  // Get all features for a project
  async getFeatures(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/features`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get feature by ID
  async getFeature(projectId, featureId) {
    try {
      const response = await api.get(`/projects/${projectId}/features/${featureId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Create new feature
  async createFeature(projectId, featureData) {
    try {
      const response = await api.post(`/projects/${projectId}/features`, featureData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Update feature
  async updateFeature(projectId, featureId, featureData) {
    try {
      const response = await api.put(`/projects/${projectId}/features/${featureId}`, featureData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Delete feature
  async deleteFeature(projectId, featureId) {
    try {
      const response = await api.delete(`/projects/${projectId}/features/${featureId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Bulk create features from CSV
  async importFeatures(projectId, featuresData) {
    try {
      const response = await api.post(`/projects/${projectId}/features/import`, {
        features: featuresData
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Export features to CSV
  async exportFeatures(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/features/export`, {
        responseType: 'blob'
      })
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get pairwise comparisons for project
  async getPairwiseComparisons(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Submit pairwise comparison
  async submitPairwiseComparison(projectId, comparisonData) {
    try {
      const response = await api.post(`/projects/${projectId}/pairwise`, comparisonData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get calculation results
  async getCalculationResults(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/results`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
}