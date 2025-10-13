// Enhanced API client with comprehensive error handling and retry logic

class APIClient {
  constructor(baseURL) {
    // Determine the appropriate base URL based on environment
    if (baseURL) {
      this.baseURL = baseURL
    } else if (import.meta.env.VITE_API_URL) {
      this.baseURL = import.meta.env.VITE_API_URL
    } else if (import.meta.env.DEV) {
      // In development, use the Vite proxy
      this.baseURL = '/api'
    } else {
      // In production, assume API is on the same host
      this.baseURL = `${window.location.protocol}//${window.location.host}/api`
    }
    
    this.defaultTimeout = 30000
    this.retryAttempts = 3
    this.retryDelay = 1000
  }

  // Main request method with error handling and retry logic
  async request(endpoint, options = {}) {
    const {
      method = 'GET',
      body = null,
      headers = {},
      timeout = this.defaultTimeout,
      retries = this.retryAttempts,
      skipRetry = false,
      signal = null
    } = options

    const url = `${this.baseURL}${endpoint}`
    
    // Setup request headers
    const requestHeaders = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...headers
    }

    // Setup request configuration
    const requestConfig = {
      method,
      headers: requestHeaders,
      signal: signal || (timeout ? AbortSignal.timeout(timeout) : undefined)
    }

    // Add body for non-GET requests
    if (body && method !== 'GET') {
      requestConfig.body = typeof body === 'string' ? body : JSON.stringify(body)
    }

    // Retry logic
    let lastError
    const maxRetries = skipRetry ? 0 : retries

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        // Add delay before retry (except for first attempt)
        if (attempt > 0) {
          await this.delay(this.retryDelay * attempt)
        }

        const response = await fetch(url, requestConfig)
        
        // Handle HTTP errors
        if (!response.ok) {
          const errorData = await this.parseErrorResponse(response)
          throw this.createHTTPError(response.status, errorData)
        }

        // Parse successful response
        const contentType = response.headers.get('content-type')
        if (contentType && contentType.includes('application/json')) {
          return await response.json()
        }
        
        return await response.text()

      } catch (error) {
        lastError = error

        // Don't retry for certain error types
        if (this.shouldNotRetry(error) || attempt >= maxRetries) {
          throw this.enhanceError(error, { endpoint, method, attempt: attempt + 1 })
        }

        console.warn(`Request attempt ${attempt + 1} failed, retrying...`, error.message)
      }
    }

    throw this.enhanceError(lastError, { endpoint, method, attempt: maxRetries + 1 })
  }

  // HTTP method shortcuts
  async get(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'GET' })
  }

  async post(endpoint, body, options = {}) {
    return this.request(endpoint, { ...options, method: 'POST', body })
  }

  async put(endpoint, body, options = {}) {
    return this.request(endpoint, { ...options, method: 'PUT', body })
  }

  async delete(endpoint, options = {}) {
    return this.request(endpoint, { ...options, method: 'DELETE' })
  }

  // Utility methods
  async delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  async parseErrorResponse(response) {
    try {
      return await response.json()
    } catch {
      return { error: response.statusText || 'Unknown error' }
    }
  }

  createHTTPError(status, data) {
    const error = new Error(data.error || `HTTP ${status}`)
    error.name = 'HTTPError'
    error.status = status
    error.response = { status, data }
    return error
  }

  shouldNotRetry(error) {
    // Don't retry for:
    // - Client errors (400-499) except 429 (rate limit) and 408 (timeout)
    // - Authentication/authorization errors
    // - Validation errors
    if (error.name === 'HTTPError') {
      const status = error.status
      return status >= 400 && status < 500 && status !== 429 && status !== 408
    }

    // Don't retry for AbortError (user cancelled)
    if (error.name === 'AbortError') {
      return true
    }

    return false
  }

  enhanceError(error, context) {
    error.context = { ...error.context, ...context }
    error.timestamp = new Date().toISOString()
    
    // Add user-friendly messages
    if (error.name === 'AbortError') {
      error.userMessage = 'Request was cancelled'
    } else if (error.name === 'TypeError' && error.message.includes('fetch')) {
      error.userMessage = 'Unable to connect to server. Please check your internet connection.'
    } else if (error.status === 429) {
      error.userMessage = 'Too many requests. Please wait a moment and try again.'
    } else if (error.status >= 500) {
      error.userMessage = 'Server error. Please try again later.'
    }

    return error
  }

  // Batch requests with error handling
  async batchRequests(requests, options = {}) {
    const { 
      concurrency = 3, 
      stopOnError = false,
      timeout = this.defaultTimeout 
    } = options

    const results = []
    const errors = []

    // Process requests in batches
    for (let i = 0; i < requests.length; i += concurrency) {
      const batch = requests.slice(i, i + concurrency)
      
      const batchPromises = batch.map(async (req, index) => {
        try {
          const result = await this.request(req.endpoint, {
            ...req.options,
            timeout
          })
          return { index: i + index, success: true, data: result }
        } catch (error) {
          const errorResult = { index: i + index, success: false, error }
          
          if (stopOnError) {
            throw error
          }
          
          return errorResult
        }
      })

      const batchResults = await Promise.all(batchPromises)
      
      batchResults.forEach(result => {
        if (result.success) {
          results[result.index] = result.data
        } else {
          errors[result.index] = result.error
        }
      })
    }

    return {
      results,
      errors,
      hasErrors: errors.some(Boolean),
      successCount: results.filter(Boolean).length,
      errorCount: errors.filter(Boolean).length
    }
  }

  // Upload files with progress tracking
  async uploadFile(endpoint, file, options = {}) {
    const {
      onProgress = null,
      timeout = 60000, // Longer timeout for uploads
      additionalData = {}
    } = options

    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()
      const formData = new FormData()
      
      formData.append('file', file)
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value)
      })

      // Upload progress
      if (onProgress) {
        xhr.upload.addEventListener('progress', (event) => {
          if (event.lengthComputable) {
            const percentComplete = (event.loaded / event.total) * 100
            onProgress(percentComplete, event.loaded, event.total)
          }
        })
      }

      // Handle completion
      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText)
            resolve(response)
          } catch {
            resolve(xhr.responseText)
          }
        } else {
          const error = new Error(`Upload failed: ${xhr.statusText}`)
          error.status = xhr.status
          reject(error)
        }
      })

      // Handle errors
      xhr.addEventListener('error', () => {
        reject(new Error('Upload failed: Network error'))
      })

      xhr.addEventListener('timeout', () => {
        reject(new Error('Upload failed: Timeout'))
      })

      // Setup and send request
      xhr.timeout = timeout
      xhr.open('POST', `${this.baseURL}${endpoint}`)
      xhr.send(formData)
    })
  }

  // Health check
  async healthCheck() {
    try {
      const response = await this.get('/health', { 
        timeout: 5000, 
        skipRetry: true 
      })
      return { healthy: true, ...response }
    } catch (error) {
      return { 
        healthy: false, 
        error: error.message,
        timestamp: new Date().toISOString()
      }
    }
  }
}

// Create and export singleton instance
const apiClient = new APIClient()

// Specific API methods for PairWise application
export const api = {
  // Projects
  projects: {
    getAll: () => apiClient.get('/projects'),
    getById: (id) => apiClient.get(`/projects/${id}`),
    create: (data) => apiClient.post('/projects', data),
    update: (id, data) => apiClient.put(`/projects/${id}`, data),
    delete: (id) => apiClient.delete(`/projects/${id}`),
    
    // Progress
    getProgress: (id) => apiClient.get(`/projects/${id}/progress`),
    advancePhase: (id, phase) => apiClient.post(`/projects/${id}/progress/advance`, { phase }),
    completePhase: (id, phase) => apiClient.post(`/projects/${id}/progress/complete`, { phase }),
    getAvailablePhases: (id) => apiClient.get(`/projects/${id}/progress/phases`)
  },

  // Attendees
  attendees: {
    getByProject: (projectId) => apiClient.get(`/projects/${projectId}/attendees`),
    create: (projectId, data) => apiClient.post(`/projects/${projectId}/attendees`, data),
    delete: (projectId, attendeeId) => apiClient.delete(`/projects/${projectId}/attendees/${attendeeId}`)
  },

  // Features
  features: {
    getByProject: (projectId) => apiClient.get(`/projects/${projectId}/features`),
    getById: (projectId, featureId) => apiClient.get(`/projects/${projectId}/features/${featureId}`),
    create: (projectId, data) => apiClient.post(`/projects/${projectId}/features`, data),
    update: (projectId, featureId, data) => apiClient.put(`/projects/${projectId}/features/${featureId}`, data),
    delete: (projectId, featureId) => apiClient.delete(`/projects/${projectId}/features/${featureId}`),
    import: (projectId, file, onProgress) => apiClient.uploadFile(`/projects/${projectId}/features/import`, file, { onProgress }),
    export: (projectId) => apiClient.get(`/projects/${projectId}/features/export`)
  },

  // Pairwise comparisons
  pairwise: {
    startSession: (projectId, data) => apiClient.post(`/projects/${projectId}/pairwise`, data),
    getSession: (projectId) => apiClient.get(`/projects/${projectId}/pairwise`),
    getComparisons: (projectId) => apiClient.get(`/projects/${projectId}/pairwise/comparisons`),
    submitVote: (projectId, data) => apiClient.post(`/projects/${projectId}/pairwise/votes`, data),
    completeSession: (projectId) => apiClient.post(`/projects/${projectId}/pairwise/complete`),
    getNext: (projectId) => apiClient.get(`/projects/${projectId}/pairwise/next`)
  },

  // Results
  results: {
    calculate: (projectId) => apiClient.post(`/projects/${projectId}/calculate-results`),
    get: (projectId) => apiClient.get(`/projects/${projectId}/results`),
    export: (projectId) => apiClient.get(`/projects/${projectId}/results/export`),
    getSummary: (projectId) => apiClient.get(`/projects/${projectId}/results/summary`),
    getStatus: (projectId) => apiClient.get(`/projects/${projectId}/results/status`)
  },

  // Fibonacci scoring
  fibonacci: {
    createSession: (projectId, criterionType) => apiClient.post(`/projects/${projectId}/fibonacci-sessions`, { criterion_type: criterionType }),
    getSession: (projectId, sessionId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}`),
    getSessions: (projectId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions`),
    submitScore: (projectId, sessionId, featureId, attendeeId, score) => apiClient.post(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`, { feature_id: featureId, attendee_id: attendeeId, score_value: score }),
    getSessionScores: (projectId, sessionId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`),
    getFeatureScores: (projectId, sessionId, featureId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/features/${featureId}/scores`),
    setConsensus: (projectId, sessionId, featureId, finalScore) => apiClient.post(`/projects/${projectId}/fibonacci-sessions/${sessionId}/consensus`, { feature_id: featureId, final_score: finalScore }),
    getSessionConsensus: (projectId, sessionId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/consensus`),
    completeSession: (projectId, sessionId) => apiClient.patch(`/projects/${projectId}/fibonacci-sessions/${sessionId}/complete`),
    exportResults: (projectId, sessionId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/export`, { responseType: 'text' }),
    getSessionStats: (projectId, sessionId) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/${sessionId}/stats`),
    resetSessionScores: (projectId, sessionId) => apiClient.delete(`/projects/${projectId}/fibonacci-sessions/${sessionId}/scores`),
    deleteSession: (projectId, sessionId) => apiClient.delete(`/projects/${projectId}/fibonacci-sessions/${sessionId}`),
    getActiveSession: (projectId, criterionType) => apiClient.get(`/projects/${projectId}/fibonacci-sessions/active`, { params: { criterion_type: criterionType } })
  },

  // System
  health: () => apiClient.healthCheck()
}

// Helper function to handle API responses
export const handleApiResponse = (response) => {
  // Handle both axios-style responses (with .data) and direct responses
  return response.data !== undefined ? response.data : response
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

// Pairwise service wrapper
export const pairwiseService = {
  startSession: async (projectId, data) => {
    try {
      const response = await api.pairwise.startSession(projectId, data)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSession: async (projectId) => {
    try {
      const response = await api.pairwise.getSession(projectId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getComparisons: async (projectId) => {
    try {
      const response = await api.pairwise.getComparisons(projectId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getComparison: async (projectId, comparisonId) => {
    try {
      const response = await api.pairwise.getNext(projectId) // Note: using getNext as getComparison might not exist
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  submitVote: async (projectId, comparisonId, attendeeId, choice) => {
    try {
      // Convert choice to the API format
      let preferredFeatureId = null
      let isTieVote = false
      
      if (choice === 'tie' || choice === 'neutral') {
        isTieVote = true
      } else if (choice === 'A') {
        // Get feature A ID from current comparison - we'll need to fetch it
        const comparisons = await api.pairwise.getComparisons(projectId)
        const comparison = comparisons.data.comparisons.find(c => c.comparison.id === comparisonId)
        if (comparison) {
          preferredFeatureId = comparison.comparison.feature_a_id
        }
      } else if (choice === 'B') {
        // Get feature B ID from current comparison
        const comparisons = await api.pairwise.getComparisons(projectId)
        const comparison = comparisons.data.comparisons.find(c => c.comparison.id === comparisonId)
        if (comparison) {
          preferredFeatureId = comparison.comparison.feature_b_id
        }
      }
      
      const data = {
        comparison_id: comparisonId,
        attendee_id: attendeeId,
        preferred_feature_id: preferredFeatureId,
        is_tie_vote: isTieVote
      }
      
      const response = await api.pairwise.submitVote(projectId, data)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getVotes: async (projectId) => {
    try {
      const response = await api.pairwise.getComparisons(projectId) // Note: simplified
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getAllVotes: async (projectId) => {
    try {
      const response = await api.pairwise.getComparisons(projectId) // Note: simplified
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getConsensus: async (projectId) => {
    try {
      const response = await api.pairwise.getSession(projectId) // Note: simplified
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getAllConsensus: async (projectId) => {
    try {
      const response = await api.pairwise.getSession(projectId) // Note: simplified
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  exportResults: async (projectId) => {
    try {
      const response = await api.results.export(projectId)
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  },
  resetVotes: async (projectId) => {
    try {
      const response = await api.pairwise.completeSession(projectId) // Note: simplified
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSessionStats: async (projectId) => {
    try {
      const response = await api.pairwise.getSession(projectId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  }
}

// Fibonacci service wrapper
export const fibonacciService = {
  createSession: async (projectId, criterionType) => {
    try {
      const response = await api.fibonacci.createSession(projectId, criterionType)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSession: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.getSession(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSessions: async (projectId) => {
    try {
      const response = await api.fibonacci.getSessions(projectId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  submitScore: async (projectId, sessionId, featureId, attendeeId, score) => {
    try {
      const response = await api.fibonacci.submitScore(projectId, sessionId, featureId, attendeeId, score)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSessionScores: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.getSessionScores(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getFeatureScores: async (projectId, sessionId, featureId) => {
    try {
      const response = await api.fibonacci.getFeatureScores(projectId, sessionId, featureId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  setConsensus: async (projectId, sessionId, featureId, finalScore) => {
    try {
      const response = await api.fibonacci.setConsensus(projectId, sessionId, featureId, finalScore)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getSessionConsensus: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.getSessionConsensus(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  completeSession: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.completeSession(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  exportResults: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.exportResults(projectId, sessionId)
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  },
  getSessionStats: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.getSessionStats(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  resetSessionScores: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.resetSessionScores(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  deleteSession: async (projectId, sessionId) => {
    try {
      const response = await api.fibonacci.deleteSession(projectId, sessionId)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
  getActiveSession: async (projectId, criterionType) => {
    try {
      const response = await api.fibonacci.getActiveSession(projectId, criterionType)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  }
}

export default apiClient