import React from 'react'

class PairwiseWebSocketService {
  constructor() {
    this.ws = null
    this.projectId = null
    this.sessionId = null
    this.attendeeId = null
    this.listeners = new Map()
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectTimeout = null
    this.isConnecting = false
  }

  connect(projectId, sessionId, attendeeId) {
    if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
      return Promise.resolve()
    }

    this.projectId = projectId
    this.sessionId = sessionId
    this.attendeeId = attendeeId
    this.isConnecting = true

    return new Promise((resolve, reject) => {
      try {
        const wsUrl = `ws://localhost:8080/ws/pairwise/${projectId}/${sessionId}?attendeeId=${attendeeId}`
        this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log('WebSocket connected for pairwise session')
          this.isConnecting = false
          this.reconnectAttempts = 0
          this.emit('connected', { projectId, sessionId, attendeeId })
          resolve()
        }

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data)
            this.handleMessage(data)
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error)
          }
        }

        this.ws.onclose = (event) => {
          console.log('WebSocket connection closed:', event.code, event.reason)
          this.isConnecting = false
          this.emit('disconnected', { code: event.code, reason: event.reason })
          
          // Attempt to reconnect if not a clean close
          if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect()
          }
        }

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error)
          this.isConnecting = false
          this.emit('error', error)
          reject(error)
        }

      } catch (error) {
        this.isConnecting = false
        reject(error)
      }
    })
  }

  disconnect() {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }

    if (this.ws) {
      this.ws.close(1000, 'Client disconnect')
      this.ws = null
    }

    this.projectId = null
    this.sessionId = null
    this.attendeeId = null
    this.reconnectAttempts = 0
  }

  scheduleReconnect() {
    if (this.reconnectTimeout) return

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000) // Exponential backoff, max 30s
    this.reconnectAttempts++

    console.log(`Scheduling reconnect attempt ${this.reconnectAttempts} in ${delay}ms`)
    
    this.reconnectTimeout = setTimeout(() => {
      this.reconnectTimeout = null
      if (this.projectId && this.sessionId && this.attendeeId) {
        this.connect(this.projectId, this.sessionId, this.attendeeId)
          .catch(error => console.error('Reconnect failed:', error))
      }
    }, delay)
  }

  handleMessage(data) {
    const { type, payload } = data

    switch (type) {
      case 'vote_update':
        this.emit('voteUpdate', payload)
        break
      
      case 'consensus_reached':
        this.emit('consensusReached', payload)
        break
      
      case 'session_progress':
        this.emit('sessionProgress', payload)
        break
      
      case 'attendee_joined':
        this.emit('attendeeJoined', payload)
        break
      
      case 'attendee_left':
        this.emit('attendeeLeft', payload)
        break
      
      case 'session_status_changed':
        this.emit('sessionStatusChanged', payload)
        break
      
      case 'error':
        this.emit('error', payload)
        break
      
      default:
        console.log('Unknown WebSocket message type:', type, payload)
    }
  }

  // Send vote to server
  sendVote(comparisonId, choice) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected')
    }

    const message = {
      type: 'submit_vote',
      payload: {
        comparisonId,
        choice,
        attendeeId: this.attendeeId,
        timestamp: new Date().toISOString()
      }
    }

    this.ws.send(JSON.stringify(message))
  }

  // Request session sync
  requestSync() {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected')
    }

    const message = {
      type: 'request_sync',
      payload: {
        attendeeId: this.attendeeId,
        timestamp: new Date().toISOString()
      }
    }

    this.ws.send(JSON.stringify(message))
  }

  // Send attendee presence heartbeat
  sendHeartbeat() {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return
    }

    const message = {
      type: 'heartbeat',
      payload: {
        attendeeId: this.attendeeId,
        timestamp: new Date().toISOString()
      }
    }

    this.ws.send(JSON.stringify(message))
  }

  // Event listener management
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)
  }

  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error(`Error in WebSocket event listener for ${event}:`, error)
        }
      })
    }
  }

  // Connection status
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }

  getConnectionState() {
    if (!this.ws) return 'disconnected'
    
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING: return 'connecting'
      case WebSocket.OPEN: return 'connected'
      case WebSocket.CLOSING: return 'closing'
      case WebSocket.CLOSED: return 'disconnected'
      default: return 'unknown'
    }
  }
}

// Create singleton instance
const pairwiseWebSocketService = new PairwiseWebSocketService()

// React hook for using the WebSocket service
export const usePairwiseWebSocket = (projectId, sessionId, attendeeId) => {
  const [connectionState, setConnectionState] = React.useState('disconnected')
  const [error, setError] = React.useState(null)

  React.useEffect(() => {
    if (!projectId || !sessionId || !attendeeId) return

    // Connection state listeners
    const handleConnected = () => {
      setConnectionState('connected')
      setError(null)
    }

    const handleDisconnected = () => {
      setConnectionState('disconnected')
    }

    const handleError = (error) => {
      setError(error)
      setConnectionState('error')
    }

    // Register listeners
    pairwiseWebSocketService.on('connected', handleConnected)
    pairwiseWebSocketService.on('disconnected', handleDisconnected)
    pairwiseWebSocketService.on('error', handleError)

    // Connect
    setConnectionState('connecting')
    pairwiseWebSocketService.connect(projectId, sessionId, attendeeId)
      .catch(error => {
        setError(error)
        setConnectionState('error')
      })

    // Heartbeat interval
    const heartbeatInterval = setInterval(() => {
      pairwiseWebSocketService.sendHeartbeat()
    }, 30000) // Every 30 seconds

    // Cleanup
    return () => {
      clearInterval(heartbeatInterval)
      pairwiseWebSocketService.off('connected', handleConnected)
      pairwiseWebSocketService.off('disconnected', handleDisconnected)
      pairwiseWebSocketService.off('error', handleError)
      pairwiseWebSocketService.disconnect()
    }
  }, [projectId, sessionId, attendeeId])

  return {
    connectionState,
    error,
    service: pairwiseWebSocketService,
    isConnected: connectionState === 'connected',
  }
}

export default pairwiseWebSocketService