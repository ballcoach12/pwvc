import { createContext, useContext, useReducer } from 'react'

// Error types
const ERROR_TYPES = {
  NETWORK: 'network',
  VALIDATION: 'validation',
  BUSINESS: 'business',
  AUTHENTICATION: 'authentication',
  PERMISSION: 'permission',
  RATE_LIMIT: 'rate_limit',
  SERVER: 'server',
  UNKNOWN: 'unknown'
}

// Error severity levels
const ERROR_SEVERITY = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical'
}

// Error actions
const ERROR_ACTIONS = {
  ADD_ERROR: 'ADD_ERROR',
  REMOVE_ERROR: 'REMOVE_ERROR',
  CLEAR_ERRORS: 'CLEAR_ERRORS',
  MARK_ERROR_READ: 'MARK_ERROR_READ'
}

// Initial state
const initialState = {
  errors: [],
  networkErrors: [],
  validationErrors: {},
  hasUnreadErrors: false
}

// Error reducer
function errorReducer(state, action) {
  switch (action.type) {
    case ERROR_ACTIONS.ADD_ERROR:
      const newError = {
        id: Date.now() + Math.random(),
        timestamp: new Date().toISOString(),
        read: false,
        ...action.payload
      }
      
      return {
        ...state,
        errors: [newError, ...state.errors],
        hasUnreadErrors: true,
        // Handle specific error types
        ...(newError.type === ERROR_TYPES.NETWORK && {
          networkErrors: [newError, ...state.networkErrors]
        }),
        ...(newError.type === ERROR_TYPES.VALIDATION && {
          validationErrors: {
            ...state.validationErrors,
            [newError.field]: newError
          }
        })
      }

    case ERROR_ACTIONS.REMOVE_ERROR:
      const filteredErrors = state.errors.filter(err => err.id !== action.payload)
      return {
        ...state,
        errors: filteredErrors,
        networkErrors: state.networkErrors.filter(err => err.id !== action.payload),
        validationErrors: Object.fromEntries(
          Object.entries(state.validationErrors).filter(([_, err]) => err.id !== action.payload)
        )
      }

    case ERROR_ACTIONS.CLEAR_ERRORS:
      return {
        ...initialState
      }

    case ERROR_ACTIONS.MARK_ERROR_READ:
      return {
        ...state,
        errors: state.errors.map(err => 
          err.id === action.payload ? { ...err, read: true } : err
        ),
        hasUnreadErrors: state.errors.some(err => !err.read && err.id !== action.payload)
      }

    default:
      return state
  }
}

// Error context
const ErrorContext = createContext()

// Error provider component
export function ErrorProvider({ children }) {
  const [state, dispatch] = useReducer(errorReducer, initialState)

  // Add error function
  const addError = (error) => {
    const errorObj = typeof error === 'string' 
      ? { message: error, type: ERROR_TYPES.UNKNOWN, severity: ERROR_SEVERITY.MEDIUM }
      : error

    dispatch({
      type: ERROR_ACTIONS.ADD_ERROR,
      payload: errorObj
    })

    // Auto-remove low severity errors after 5 seconds
    if (errorObj.severity === ERROR_SEVERITY.LOW) {
      setTimeout(() => {
        removeError(errorObj.id)
      }, 5000)
    }
  }

  // Remove error function
  const removeError = (errorId) => {
    dispatch({
      type: ERROR_ACTIONS.REMOVE_ERROR,
      payload: errorId
    })
  }

  // Clear all errors
  const clearErrors = () => {
    dispatch({
      type: ERROR_ACTIONS.CLEAR_ERRORS
    })
  }

  // Mark error as read
  const markErrorRead = (errorId) => {
    dispatch({
      type: ERROR_ACTIONS.MARK_ERROR_READ,
      payload: errorId
    })
  }

  // Handle API errors
  const handleApiError = (error, context = {}) => {
    if (!error) return

    let errorType = ERROR_TYPES.UNKNOWN
    let severity = ERROR_SEVERITY.MEDIUM
    let message = 'An unexpected error occurred'

    // Handle network errors
    if (!navigator.onLine) {
      errorType = ERROR_TYPES.NETWORK
      message = 'You appear to be offline. Please check your internet connection.'
      severity = ERROR_SEVERITY.HIGH
    } else if (error.name === 'TypeError' && error.message.includes('fetch')) {
      errorType = ERROR_TYPES.NETWORK
      message = 'Unable to connect to the server. Please try again.'
      severity = ERROR_SEVERITY.HIGH
    } else if (error.response) {
      // Handle HTTP errors
      const status = error.response.status
      const data = error.response.data

      if (status === 400) {
        errorType = ERROR_TYPES.VALIDATION
        message = data?.error || 'Invalid request data'
        severity = ERROR_SEVERITY.LOW
        
        // Handle field-specific validation errors
        if (data?.field) {
          addError({
            type: ERROR_TYPES.VALIDATION,
            field: data.field,
            message: data.message || message,
            severity,
            context
          })
          return
        }
      } else if (status === 401) {
        errorType = ERROR_TYPES.AUTHENTICATION
        message = 'Your session has expired. Please log in again.'
        severity = ERROR_SEVERITY.HIGH
      } else if (status === 403) {
        errorType = ERROR_TYPES.PERMISSION
        message = 'You do not have permission to perform this action.'
        severity = ERROR_SEVERITY.MEDIUM
      } else if (status === 404) {
        errorType = ERROR_TYPES.BUSINESS
        message = data?.error || 'The requested resource was not found'
        severity = ERROR_SEVERITY.MEDIUM
      } else if (status === 409) {
        errorType = ERROR_TYPES.BUSINESS
        message = data?.error || 'This action conflicts with existing data'
        severity = ERROR_SEVERITY.MEDIUM
      } else if (status === 422) {
        errorType = ERROR_TYPES.VALIDATION
        message = data?.error || 'The provided data is invalid'
        severity = ERROR_SEVERITY.MEDIUM
      } else if (status === 429) {
        errorType = ERROR_TYPES.RATE_LIMIT
        message = 'Too many requests. Please wait a moment and try again.'
        severity = ERROR_SEVERITY.HIGH
      } else if (status >= 500) {
        errorType = ERROR_TYPES.SERVER
        message = 'Server error. Our team has been notified and is working on a fix.'
        severity = ERROR_SEVERITY.CRITICAL
      }

      if (data?.error) {
        message = data.error
      }
    }

    addError({
      type: errorType,
      message,
      severity,
      context,
      originalError: error
    })
  }

  // Get validation error for field
  const getValidationError = (field) => {
    return state.validationErrors[field]?.message
  }

  // Clear validation errors
  const clearValidationErrors = () => {
    const newErrors = state.errors.filter(err => err.type !== ERROR_TYPES.VALIDATION)
    dispatch({
      type: ERROR_ACTIONS.CLEAR_ERRORS
    })
    newErrors.forEach(err => addError(err))
  }

  const value = {
    ...state,
    addError,
    removeError,
    clearErrors,
    markErrorRead,
    handleApiError,
    getValidationError,
    clearValidationErrors,
    ERROR_TYPES,
    ERROR_SEVERITY
  }

  return (
    <ErrorContext.Provider value={value}>
      {children}
    </ErrorContext.Provider>
  )
}

// Hook to use error context
export function useErrors() {
  const context = useContext(ErrorContext)
  if (!context) {
    throw new Error('useErrors must be used within an ErrorProvider')
  }
  return context
}