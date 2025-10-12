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

export default api