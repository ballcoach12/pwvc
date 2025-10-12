import axios from 'axios'

// Create axios instance with default config
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for adding auth headers if needed
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for handling errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Handle common errors
    if (error.response?.status === 401) {
      // Unauthorized - redirect to login or clear token
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    } else if (error.response?.status === 403) {
      // Forbidden
      console.error('Access denied')
    } else if (error.response?.status >= 500) {
      // Server error
      console.error('Server error:', error.response.data?.message || error.message)
    }
    
    return Promise.reject(error)
  }
)

// Helper function to handle API responses
export const handleApiResponse = (response) => {
  return response.data
}

// Helper function to handle API errors
export const handleApiError = (error) => {
  if (error.response) {
    // Server responded with error status
    const message = error.response.data?.message || error.response.data?.error || 'An error occurred'
    throw new Error(message)
  } else if (error.request) {
    // Network error
    throw new Error('Network error. Please check your connection.')
  } else {
    // Other error
    throw new Error(error.message || 'An unexpected error occurred')
  }
}

// Pairwise Service
export const pairwiseService = {
  // Get all pairwise comparisons for a project
  getComparisons: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get a specific comparison
  getComparison: async (projectId, comparisonId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise/${comparisonId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Submit a vote for a comparison
  submitVote: async (projectId, comparisonId, attendeeId, choice) => {
    try {
      const response = await api.post(`/projects/${projectId}/pairwise/${comparisonId}/vote`, {
        attendee_id: attendeeId,
        choice: choice
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get votes for a comparison
  getVotes: async (projectId, comparisonId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise/${comparisonId}/votes`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get all votes for all comparisons in a project
  getAllVotes: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/votes`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Check consensus status for a comparison
  getConsensus: async (projectId, comparisonId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise/${comparisonId}/consensus`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get consensus status for all comparisons
  getAllConsensus: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/consensus`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Export comparison results
  exportResults: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise/export`, {
        responseType: 'text'
      })
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  },

  // Reset all votes for a project (admin function)
  resetVotes: async (projectId) => {
    try {
      const response = await api.delete(`/projects/${projectId}/votes`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get session statistics
  getSessionStats: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/pairwise/stats`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  }
}

export default api