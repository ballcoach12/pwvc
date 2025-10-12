import {
    ArrowBack,
    Assessment,
    Download,
    Group,
    Help,
    Settings
} from '@mui/icons-material'
import {
    Alert,
    AppBar,
    Box,
    Breadcrumbs,
    Button,
    Chip,
    Container,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControlLabel,
    Grid,
    IconButton,
    Link,
    Paper,
    Snackbar,
    Switch,
    Toolbar,
    Typography
} from '@mui/material'
import { useCallback, useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import ConsensusTracker from '../components/ConsensusTracker'
import FibonacciScoringGrid from '../components/FibonacciScoringGrid'
import KeyboardShortcutsDialog from '../components/KeyboardShortcutsDialog'
import { useAttendees } from '../hooks/useAttendees'
import { useFeatures } from '../hooks/useFeatures'
import { useProject } from '../hooks/useProject'
import { fibonacciService } from '../services/api'
import { useWebSocket } from '../services/websocketService'

/**
 * FibonacciScoring is the main page for conducting Fibonacci scoring sessions.
 * Handles both Value and Complexity scoring with real-time collaboration,
 * consensus tracking, and session management.
 * 
 * Features:
 * - Value and Complexity scoring sessions
 * - Real-time score updates via WebSocket
 * - Consensus tracking and progress monitoring
 * - Attendee selection and participation
 * - Session completion and phase transitions
 * - Export functionality for scores
 * - Mobile-friendly responsive design
 */
const FibonacciScoring = () => {
  const { projectId, criterionType = 'value' } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  
  // State management
  const [currentSession, setCurrentSession] = useState(null)
  const [currentAttendee, setCurrentAttendee] = useState(null)
  const [individualScores, setIndividualScores] = useState({}) // { featureId: { attendeeId: score } }
  const [consensusScores, setConsensusScores] = useState({}) // { featureId: consensusScore }
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [shortcutsDialogOpen, setShortcutsDialogOpen] = useState(false)
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'info' })
  const [autoAdvance, setAutoAdvance] = useState(false)
  const [viewMode, setViewMode] = useState('split') // 'split' | 'grid-only' | 'tracker-only'
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Custom hooks
  const { project, loading: projectLoading } = useProject(projectId)
  const { features, loading: featuresLoading } = useFeatures(projectId)
  const { attendees, loading: attendeesLoading } = useAttendees(projectId)
  
  // WebSocket integration
  const {
    isConnected,
    fibonacciScores,
    fibonacciConsensus,
    sendScore,
    joinSession,
    leaveSession
  } = useWebSocket(projectId)

  // Initialize session on component mount
  useEffect(() => {
    const initializeSession = async () => {
      try {
        setLoading(true)
        
        // Create or get existing session
        const session = await fibonacciService.createSession(projectId, criterionType)
        setCurrentSession(session)
        
        // Load existing scores
        const [scores, consensus] = await Promise.all([
          fibonacciService.getSessionScores(projectId, session.id),
          fibonacciService.getSessionConsensus(projectId, session.id)
        ])
        
        setIndividualScores(scores)
        setConsensusScores(consensus)
        
        // Join WebSocket session
        if (isConnected) {
          joinSession(`fibonacci-${session.id}`)
        }
        
      } catch (err) {
        setError(err)
        showNotification('Failed to initialize scoring session', 'error')
      } finally {
        setLoading(false)
      }
    }

    if (projectId && !projectLoading && !featuresLoading && !attendeesLoading) {
      initializeSession()
    }
  }, [projectId, criterionType, isConnected, projectLoading, featuresLoading, attendeesLoading])

  // Update scores from WebSocket
  useEffect(() => {
    if (fibonacciScores) {
      setIndividualScores(prevScores => ({
        ...prevScores,
        ...fibonacciScores
      }))
    }
  }, [fibonacciScores])

  // Update consensus from WebSocket
  useEffect(() => {
    if (fibonacciConsensus) {
      setConsensusScores(prevConsensus => ({
        ...prevConsensus,
        ...fibonacciConsensus
      }))
    }
  }, [fibonacciConsensus])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (currentSession) {
        leaveSession(`fibonacci-${currentSession.id}`)
      }
    }
  }, [currentSession, leaveSession])

  // Handle score submission
  const handleScoreSubmit = useCallback(async (featureId, score, attendeeId) => {
    try {
      if (!currentSession) return

      // Submit to backend
      await fibonacciService.submitScore(projectId, currentSession.id, featureId, attendeeId, score)
      
      // Update local state
      setIndividualScores(prev => ({
        ...prev,
        [featureId]: {
          ...prev[featureId],
          [attendeeId]: score
        }
      }))

      // Send via WebSocket for real-time updates
      sendScore(currentSession.id, featureId, attendeeId, score)
      
      showNotification('Score submitted successfully', 'success')
      
      // Check for consensus
      await checkConsensus(featureId)
      
    } catch (err) {
      console.error('Failed to submit score:', err)
      showNotification('Failed to submit score', 'error')
    }
  }, [currentSession, projectId, sendScore])

  // Handle score change (before submission)
  const handleScoreChange = useCallback((featureId, score, attendeeId) => {
    // Update temporary local state for immediate UI feedback
    setIndividualScores(prev => ({
      ...prev,
      [featureId]: {
        ...prev[featureId],
        [attendeeId]: score
      }
    }))
  }, [])

  // Check consensus for a feature
  const checkConsensus = async (featureId) => {
    try {
      if (!currentSession) return

      const featureScores = individualScores[featureId] || {}
      const scoreValues = Object.values(featureScores)
      const uniqueScores = [...new Set(scoreValues)]
      
      // Consensus achieved if all attendees have scored and all scores are the same
      if (scoreValues.length === attendees.length && uniqueScores.length === 1) {
        const consensusScore = uniqueScores[0]
        
        // Submit consensus to backend
        await fibonacciService.setConsensus(projectId, currentSession.id, featureId, consensusScore)
        
        // Update local state
        setConsensusScores(prev => ({
          ...prev,
          [featureId]: consensusScore
        }))
        
        showNotification(`Consensus reached for feature: ${consensusScore}`, 'success')
      }
    } catch (err) {
      console.error('Failed to check consensus:', err)
    }
  }

  // Handle attendee selection
  const handleAttendeeSelect = (attendee) => {
    setCurrentAttendee(attendee)
    showNotification(`Switched to ${attendee.name}`, 'info')
  }

  // Handle feature click in tracker
  const handleFeatureClick = (feature) => {
    // Scroll to feature in grid (implementation would depend on grid layout)
    showNotification(`Viewing ${feature.name}`, 'info')
  }

  // Handle session completion
  const handleCompleteSession = async () => {
    try {
      if (!currentSession) return

      await fibonacciService.completeSession(projectId, currentSession.id)
      setCurrentSession(prev => ({ ...prev, status: 'completed' }))
      
      showNotification('Session completed successfully!', 'success')
      
      // Navigate to next phase or results
      if (criterionType === 'value') {
        navigate(`/projects/${projectId}/scoring/complexity`)
      } else {
        navigate(`/projects/${projectId}/results`)
      }
    } catch (err) {
      console.error('Failed to complete session:', err)
      showNotification('Failed to complete session', 'error')
    }
  }

  // Handle next phase navigation
  const handleNextPhase = () => {
    if (criterionType === 'value') {
      navigate(`/projects/${projectId}/scoring/complexity`)
    } else {
      navigate(`/projects/${projectId}/results`)
    }
  }

  // Handle export
  const handleExportResults = async () => {
    try {
      if (!currentSession) return

      const results = await fibonacciService.exportResults(projectId, currentSession.id)
      
      // Create and download CSV file
      const blob = new Blob([results], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${project?.name || projectId}-${criterionType}-scores.csv`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
      
      showNotification('Results exported successfully', 'success')
    } catch (err) {
      console.error('Failed to export results:', err)
      showNotification('Failed to export results', 'error')
    }
  }

  // Utility functions
  const showNotification = (message, severity = 'info') => {
    setNotification({ open: true, message, severity })
  }

  const closeNotification = () => {
    setNotification({ ...notification, open: false })
  }

  // Loading and error states
  if (loading || projectLoading || featuresLoading || attendeesLoading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading scoring session...</Typography>
      </Container>
    )
  }

  if (error) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">
          Failed to load scoring session: {error.message}
        </Alert>
      </Container>
    )
  }

  if (!features || features.length === 0) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="warning">
          No features available for scoring. Please add features to the project first.
          <Button 
            onClick={() => navigate(`/projects/${projectId}/features`)}
            sx={{ ml: 2 }}
          >
            Manage Features
          </Button>
        </Alert>
      </Container>
    )
  }

  const criterionLabel = criterionType === 'value' ? 'Value' : 'Complexity'
  const isSessionComplete = Object.keys(consensusScores).length === features.length

  return (
    <Box sx={{ flexGrow: 1 }}>
      {/* Header */}
      <AppBar position="static" color="transparent" elevation={1}>
        <Toolbar>
          <IconButton
            edge="start"
            onClick={() => navigate(`/projects/${projectId}/comparison`)}
            sx={{ mr: 2 }}
          >
            <ArrowBack />
          </IconButton>
          
          <Breadcrumbs sx={{ flexGrow: 1 }}>
            <Link 
              component="button" 
              variant="body2" 
              onClick={() => navigate('/projects')}
              underline="hover"
            >
              Projects
            </Link>
            <Link 
              component="button" 
              variant="body2" 
              onClick={() => navigate(`/projects/${projectId}`)}
              underline="hover"
            >
              {project?.name}
            </Link>
            <Typography variant="body2" color="text.primary">
              {criterionLabel} Scoring
            </Typography>
          </Breadcrumbs>

          <Box sx={{ display: 'flex', gap: 1 }}>
            <Chip
              icon={<Assessment />}
              label={`${criterionLabel} Session`}
              color={criterionType === 'value' ? 'primary' : 'secondary'}
              variant="outlined"
            />
            
            <IconButton
              onClick={() => setShortcutsDialogOpen(true)}
              title="Keyboard shortcuts"
            >
              <Help />
            </IconButton>
            
            <IconButton onClick={() => setSettingsOpen(true)} title="Settings">
              <Settings />
            </IconButton>
            
            <Button
              startIcon={<Download />}
              onClick={handleExportResults}
              size="small"
              disabled={Object.keys(individualScores).length === 0}
            >
              Export
            </Button>
          </Box>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Container maxWidth="xl" sx={{ py: 3 }}>
        {/* Attendee Selection */}
        <Paper sx={{ p: 2, mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
            <Typography variant="body2" color="text.secondary">
              Scoring as:
            </Typography>
            {attendees.map(attendee => (
              <Button
                key={attendee.id}
                variant={currentAttendee?.id === attendee.id ? 'contained' : 'outlined'}
                size="small"
                onClick={() => handleAttendeeSelect(attendee)}
                startIcon={<Group />}
              >
                {attendee.name}
              </Button>
            ))}
            {!currentAttendee && (
              <Alert severity="info" sx={{ ml: 2 }}>
                Please select an attendee to start scoring
              </Alert>
            )}
          </Box>
        </Paper>

        {/* Main Layout */}
        <Grid container spacing={3}>
          {/* Left Panel - Consensus Tracker */}
          {(viewMode === 'split' || viewMode === 'tracker-only') && (
            <Grid item xs={12} lg={4}>
              <ConsensusTracker
                features={features}
                attendees={attendees}
                consensusScores={consensusScores}
                individualScores={individualScores}
                criterionType={criterionType}
                sessionStatus={currentSession?.status}
                onFeatureClick={handleFeatureClick}
                onExportResults={handleExportResults}
                onCompleteSession={handleCompleteSession}
                onNextPhase={handleNextPhase}
              />
            </Grid>
          )}

          {/* Right Panel - Scoring Grid */}
          {(viewMode === 'split' || viewMode === 'grid-only') && (
            <Grid item xs={12} lg={viewMode === 'split' ? 8 : 12}>
              <FibonacciScoringGrid
                projectId={projectId}
                sessionId={currentSession?.id}
                features={features}
                attendees={attendees}
                criterionType={criterionType}
                currentAttendee={currentAttendee}
                individualScores={individualScores}
                consensusScores={consensusScores}
                onScoreSubmit={handleScoreSubmit}
                onScoreChange={handleScoreChange}
              />
            </Grid>
          )}
        </Grid>
      </Container>

      {/* Settings Dialog */}
      <Dialog open={settingsOpen} onClose={() => setSettingsOpen(false)}>
        <DialogTitle>Scoring Session Settings</DialogTitle>
        <DialogContent>
          <FormControlLabel
            control={
              <Switch
                checked={autoAdvance}
                onChange={(e) => setAutoAdvance(e.target.checked)}
              />
            }
            label="Auto-advance after scoring"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSettingsOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Keyboard Shortcuts Dialog */}
      <KeyboardShortcutsDialog
        open={shortcutsDialogOpen}
        onClose={() => setShortcutsDialogOpen(false)}
      />

      {/* Notifications */}
      <Snackbar
        open={notification.open}
        autoHideDuration={6000}
        onClose={closeNotification}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert
          onClose={closeNotification}
          severity={notification.severity}
          sx={{ width: '100%' }}
        >
          {notification.message}
        </Alert>
      </Snackbar>
    </Box>
  )
}

export default FibonacciScoring