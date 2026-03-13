import { Component, useEffect } from 'react'
import { useErrors } from '../../contexts/ErrorContext'
import './ErrorDisplay.css'

// Main error display component
export function ErrorDisplay() {
  const { errors, removeError, markErrorRead } = useErrors()

  // Filter errors that should be displayed
  const displayErrors = errors.filter(error => 
    error.severity !== 'low' || !error.read
  )

  if (displayErrors.length === 0) return null

  return (
    <div className="error-display-container">
      {displayErrors.map(error => (
        <ErrorNotification 
          key={error.id} 
          error={error}
          onDismiss={() => removeError(error.id)}
          onRead={() => markErrorRead(error.id)}
        />
      ))}
    </div>
  )
}

// Individual error notification
function ErrorNotification({ error, onDismiss, onRead }) {
  useEffect(() => {
    if (!error.read) {
      onRead()
    }
  }, [error.read, onRead])

  const getSeverityClass = (severity) => {
    switch (severity) {
      case 'low': return 'error-low'
      case 'medium': return 'error-medium'
      case 'high': return 'error-high'
      case 'critical': return 'error-critical'
      default: return 'error-medium'
    }
  }

  const getIcon = (type) => {
    switch (type) {
      case 'network': return '🌐'
      case 'validation': return '⚠️'
      case 'business': return 'ℹ️'
      case 'authentication': return '🔐'
      case 'permission': return '🚫'
      case 'rate_limit': return '⏱️'
      case 'server': return '🔧'
      default: return '❌'
    }
  }

  return (
    <div className={`error-notification ${getSeverityClass(error.severity)}`}>
      <div className="error-header">
        <span className="error-icon">{getIcon(error.type)}</span>
        <span className="error-title">
          {error.type === 'network' && 'Connection Error'}
          {error.type === 'validation' && 'Validation Error'}
          {error.type === 'business' && 'Information'}
          {error.type === 'authentication' && 'Authentication Required'}
          {error.type === 'permission' && 'Access Denied'}
          {error.type === 'rate_limit' && 'Rate Limited'}
          {error.type === 'server' && 'Server Error'}
          {error.type === 'unknown' && 'Error'}
        </span>
        <button 
          className="error-close"
          onClick={onDismiss}
          aria-label="Dismiss error"
        >
          ×
        </button>
      </div>
      
      <div className="error-message">
        {error.message}
      </div>

      {error.context && (
        <div className="error-context">
          <small>Context: {JSON.stringify(error.context)}</small>
        </div>
      )}

      {error.severity === 'critical' && (
        <div className="error-actions">
          <button 
            className="error-retry-btn"
            onClick={() => window.location.reload()}
          >
            Reload Page
          </button>
        </div>
      )}
    </div>
  )
}

// Inline error display for form fields
export function FieldError({ field }) {
  const { getValidationError } = useErrors()
  const error = getValidationError(field)

  if (!error) return null

  return (
    <div className="field-error">
      <span className="field-error-icon">⚠️</span>
      <span className="field-error-message">{error}</span>
    </div>
  )
}

// Network status indicator
export function NetworkStatus() {
  const { networkErrors } = useErrors()
  const isOnline = navigator.onLine

  if (isOnline && networkErrors.length === 0) return null

  return (
    <div className={`network-status ${isOnline ? 'online' : 'offline'}`}>
      {isOnline ? (
        <>
          <span className="status-indicator online"></span>
          Connection restored
        </>
      ) : (
        <>
          <span className="status-indicator offline"></span>
          You are offline
        </>
      )}
    </div>
  )
}

// Error boundary component
export class ErrorBoundary extends Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error }
  }

  componentDidCatch(error, errorInfo) {
    console.error('Error caught by boundary:', error, errorInfo)
    
    // In a real app, you would send this to an error reporting service
    if (this.props.onError) {
      this.props.onError(error, errorInfo)
    }
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="error-boundary">
          <h2>Something went wrong</h2>
          <p>We apologize for the inconvenience. Please refresh the page or try again later.</p>
          <button 
            onClick={() => window.location.reload()}
            className="error-boundary-reload"
          >
            Reload Page
          </button>
        </div>
      )
    }

    return this.props.children
  }
}