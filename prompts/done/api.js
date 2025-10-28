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

// Fibonacci Scoring Service
export const fibonacciService = {
  // Create or get existing Fibonacci scoring session
  createSession: async (projectId, criterionType) => {
    try {
      const response = await api.post(`/projects/${projectId}/fibonacci-sessions`, {
        criterion_type: criterionType
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get Fibonacci scoring session
  getSession: async (projectId, sessionId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get all sessions for a project
  getSessions: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Submit individual Fibonacci score
  submitScore: async (projectId, sessionId, featureId, attendeeId, score) => {
    try {
      const response = await api.post(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`, {
        feature_id: featureId,
        attendee_id: attendeeId,
        score_value: score
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get all scores for a session
  getSessionScores: async (projectId, sessionId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get scores for a specific feature in a session
  getFeatureScores: async (projectId, sessionId, featureId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/features/${featureId}/scores`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Set consensus score for a feature
  setConsensus: async (projectId, sessionId, featureId, finalScore) => {
    try {
      const response = await api.post(`/projects/${projectId}/fibonacci-sessions/${sessionId}/consensus`, {
        feature_id: featureId,
        final_score: finalScore
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get consensus scores for a session
  getSessionConsensus: async (projectId, sessionId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/consensus`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Complete a Fibonacci scoring session
  completeSession: async (projectId, sessionId) => {
    try {
      const response = await api.patch(`/projects/${projectId}/fibonacci-sessions/${sessionId}/complete`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Export session results as CSV
  exportResults: async (projectId, sessionId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/export`, {
        responseType: 'text'
      })
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get session statistics
  getSessionStats: async (projectId, sessionId) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/stats`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Reset all scores for a session (admin function)
  resetSessionScores: async (projectId, sessionId) => {
    try {
      const response = await api.delete(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Delete a scoring session
  deleteSession: async (projectId, sessionId) => {
    try {
      const response = await api.delete(`/projects/${projectId}/fibonacci-sessions/${sessionId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get current active session for a criterion type
  getActiveSession: async (projectId, criterionType) => {
    try {
      const response = await api.get(`/projects/${projectId}/fibonacci-sessions/active`, {
        params: { criterion_type: criterionType }
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  }
}

export default api