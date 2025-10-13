import React from 'react'

class PairWiseWebSocketService {
  constructor() {
    this.ws = null
    this.projectId = null
    this.sessionId = null
    this.sessionType = 'pairwise' // 'pairwise' | 'fibonacci'
    this.attendeeId = null
    this.listeners = new Map()
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectTimeout = null
    this.isConnecting = false
  }

  connect(projectId, sessionId, attendeeId, sessionType = 'pairwise') {
    if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
      return Promise.resolve()
    }

    this.projectId = projectId
    this.sessionId = sessionId
    this.attendeeId = attendeeId
    this.sessionType = sessionType
    this.isConnecting = true

    return new Promise((resolve, reject) => {
      try {
        const wsUrl = `ws://localhost:8080/ws/${sessionType}/${projectId}/${sessionId}?attendeeId=${attendeeId}`
        this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log(`WebSocket connected for ${sessionType} session`)
          this.isConnecting = false
          this.reconnectAttempts = 0
          this.emit('connected', { projectId, sessionId, attendeeId, sessionType })
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
      // Pairwise comparison events
      case 'vote_update':
        this.emit('voteUpdate', payload)
        break
      
      case 'consensus_reached':
        this.emit('consensusReached', payload)
        break
      
      // Fibonacci scoring events
      case 'fibonacci_score_update':
        this.emit('fibonacciScoreUpdate', payload)
        break
      
      case 'fibonacci_consensus_reached':
        this.emit('fibonacciConsensusReached', payload)
        break
      
      case 'fibonacci_session_progress':
        this.emit('fibonacciSessionProgress', payload)
        break
      
      // Common events
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

  // Send vote to server (pairwise comparison)
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

  // Send Fibonacci score to server
  sendFibonacciScore(sessionId, featureId, score) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected')
    }

    const message = {
      type: 'submit_fibonacci_score',
      payload: {
        sessionId,
        featureId,
        score,
        attendeeId: this.attendeeId,
        timestamp: new Date().toISOString()
      }
    }

    this.ws.send(JSON.stringify(message))
  }

  // Send Fibonacci consensus to server
  sendFibonacciConsensus(sessionId, featureId, consensusScore) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket not connected')
    }

    const message = {
      type: 'submit_fibonacci_consensus',
      payload: {
        sessionId,
        featureId,
        consensusScore,
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
const pairwiseWebSocketService = new PairWiseWebSocketService()

// React hook for pairwise comparison WebSocket
export const usePairwiseWebSocket = (projectId, sessionId, attendeeId) => {
  const [connectionState, setConnectionState] = React.useState('disconnected')
  const [error, setError] = React.useState(null)
  const [votes, setVotes] = React.useState({})
  const [consensus, setConsensus] = React.useState({})

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

    // Vote update listeners
    const handleVoteUpdate = (payload) => {
      setVotes(prev => ({
        ...prev,
        [payload.comparisonId]: {
          ...prev[payload.comparisonId],
          [payload.attendeeId]: payload.choice
        }
      }))
    }

    const handleConsensusReached = (payload) => {
      setConsensus(prev => ({
        ...prev,
        [payload.comparisonId]: payload.consensusChoice
      }))
    }

    // Register listeners
    pairwiseWebSocketService.on('connected', handleConnected)
    pairwiseWebSocketService.on('disconnected', handleDisconnected)
    pairwiseWebSocketService.on('error', handleError)
    pairwiseWebSocketService.on('voteUpdate', handleVoteUpdate)
    pairwiseWebSocketService.on('consensusReached', handleConsensusReached)

    // Connect
    setConnectionState('connecting')
    pairwiseWebSocketService.connect(projectId, sessionId, attendeeId, 'pairwise')
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
      pairwiseWebSocketService.off('voteUpdate', handleVoteUpdate)
      pairwiseWebSocketService.off('consensusReached', handleConsensusReached)
      pairwiseWebSocketService.disconnect()
    }
  }, [projectId, sessionId, attendeeId])

  return {
    connectionState,
    error,
    votes,
    consensus,
    service: pairwiseWebSocketService,
    isConnected: connectionState === 'connected',
    sendVote: (comparisonId, choice) => pairwiseWebSocketService.sendVote(comparisonId, choice),
    requestSync: () => pairwiseWebSocketService.requestSync()
  }
}

// React hook for Fibonacci scoring WebSocket
export const useFibonacciWebSocket = (projectId, sessionId, attendeeId) => {
  const [connectionState, setConnectionState] = React.useState('disconnected')
  const [error, setError] = React.useState(null)
  const [fibonacciScores, setFibonacciScores] = React.useState({}) // { featureId: { attendeeId: score } }
  const [fibonacciConsensus, setFibonacciConsensus] = React.useState({}) // { featureId: consensusScore }

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

    // Fibonacci score update listeners
    const handleFibonacciScoreUpdate = (payload) => {
      setFibonacciScores(prev => ({
        ...prev,
        [payload.featureId]: {
          ...prev[payload.featureId],
          [payload.attendeeId]: payload.score
        }
      }))
    }

    const handleFibonacciConsensusReached = (payload) => {
      setFibonacciConsensus(prev => ({
        ...prev,
        [payload.featureId]: payload.consensusScore
      }))
    }

    // Register listeners
    pairwiseWebSocketService.on('connected', handleConnected)
    pairwiseWebSocketService.on('disconnected', handleDisconnected)
    pairwiseWebSocketService.on('error', handleError)
    pairwiseWebSocketService.on('fibonacciScoreUpdate', handleFibonacciScoreUpdate)
    pairwiseWebSocketService.on('fibonacciConsensusReached', handleFibonacciConsensusReached)

    // Connect
    setConnectionState('connecting')
    pairwiseWebSocketService.connect(projectId, sessionId, attendeeId, 'fibonacci')
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
      pairwiseWebSocketService.off('fibonacciScoreUpdate', handleFibonacciScoreUpdate)
      pairwiseWebSocketService.off('fibonacciConsensusReached', handleFibonacciConsensusReached)
      pairwiseWebSocketService.disconnect()
    }
  }, [projectId, sessionId, attendeeId])

  return {
    connectionState,
    error,
    fibonacciScores,
    fibonacciConsensus,
    service: pairwiseWebSocketService,
    isConnected: connectionState === 'connected',
    sendScore: (sessionId, featureId, score) => pairwiseWebSocketService.sendFibonacciScore(sessionId, featureId, score),
    sendConsensus: (sessionId, featureId, consensusScore) => pairwiseWebSocketService.sendFibonacciConsensus(sessionId, featureId, consensusScore),
    requestSync: () => pairwiseWebSocketService.requestSync()
  }
}

// Generic WebSocket hook (backwards compatibility)
export const useWebSocket = (projectId) => {
  return {
    isConnected: pairwiseWebSocketService.isConnected(),
    service: pairwiseWebSocketService,
    joinSession: (sessionId) => pairwiseWebSocketService.connect(projectId, sessionId, null, 'pairwise'),
    leaveSession: () => pairwiseWebSocketService.disconnect(),
    sendVote: (comparisonId, choice, attendeeId) => pairwiseWebSocketService.sendVote(comparisonId, choice),
    sendScore: (sessionId, featureId, attendeeId, score) => pairwiseWebSocketService.sendFibonacciScore(sessionId, featureId, score)
  }
}

export default pairwiseWebSocketService