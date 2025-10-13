import {
  ArrowBack,
  Download,
  Fullscreen,
  FullscreenExit,
  GridView,
  Help,
  NavigateBefore,
  NavigateNext,
  Settings,
  ViewList
} from '@mui/icons-material'
import {
  Alert,
  AppBar,
  Box,
  Breadcrumbs,
  Button,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Fab,
  FormControlLabel,
  Grid,
  IconButton,
  Link,
  Snackbar,
  Switch,
  Toolbar,
  Typography
} from '@mui/material'
import { useEffect, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import AttendeeVotingPanel from '../components/AttendeeVotingPanel/AttendeeVotingPanel'
import ComparisonCard from '../components/ComparisonCard/ComparisonCard'
import KeyboardShortcutsDialog from '../components/KeyboardShortcutsDialog'
import PairwiseGrid from '../components/PairwiseGrid'
import SessionProgress from '../components/SessionProgress/SessionProgress'
import { useAttendees } from '../hooks/useAttendees'
import { useFeatures } from '../hooks/useFeatures'
import useKeyboardShortcuts from '../hooks/useKeyboardShortcuts'
import { useProject } from '../hooks/useProject'
import { pairwiseService, api } from '../services/api'
import AttendeeLoginDialog from '../components/AttendeeLoginDialog/AttendeeLoginDialog'
import pairwiseWebSocketService, { useWebSocket } from '../services/websocketService'

/**
 * PairwiseComparison is the main page for conducting pairwise feature comparisons.
 * It orchestrates the entire comparison session with real-time collaboration,
 * progress tracking, and multiple viewing modes.
 * 
 * Features:
 * - Grid view showing all comparisons
 * - Detailed comparison view for individual pairs
 * - Real-time voting and consensus tracking
 * - Session progress monitoring
 * - Attendee participation management
 * - Keyboard shortcuts for efficient navigation
 * - Export and sharing capabilities
 * - Fullscreen mode for focused voting
 */
const PairwiseComparison = () => {
  // Router hooks
  const { projectId } = useParams()
  const navigate = useNavigate()
  
  // State management
  const [viewMode, setViewMode] = useState('grid') // 'grid' | 'detail'
  const [currentComparison, setCurrentComparison] = useState(null)
  const [currentAttendee, setCurrentAttendee] = useState(null)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [autoAdvance, setAutoAdvance] = useState(true)
  const [notification, setNotification] = useState({ open: false, message: '', severity: 'info' })
  const [sessionStarted, setSessionStarted] = useState(false)
  const [shortcutsDialogOpen, setShortcutsDialogOpen] = useState(false)
  const [votes, setVotes] = useState({})
  const [consensus, setConsensus] = useState({})
  const [attendeeStatus, setAttendeeStatus] = useState({})

  // Refs for keyboard handling
  const containerRef = useRef(null)

  // Custom hooks
  const { project, loading: projectLoading, error: projectError } = useProject(projectId)
  const { features, loading: featuresLoading } = useFeatures(projectId)
  const { attendees, loading: attendeesLoading } = useAttendees(projectId)
  
  const {
    isConnected,
    sendVote,
    joinSession,
    leaveSession
  } = useWebSocket(projectId)

  // Authentication state
  const [loginDialogOpen, setLoginDialogOpen] = useState(false)

  // Debug logging for attendee selection
  useEffect(() => {
    console.log('Current attendee state:', currentAttendee)
    console.log('Available attendees:', attendees)
  }, [currentAttendee, attendees])

  // Check for authenticated attendee on mount
  useEffect(() => {
    const authenticatedAttendee = api.auth.getCurrentAttendee()
    if (authenticatedAttendee) {
      console.log('Found authenticated attendee:', authenticatedAttendee)
      setCurrentAttendee(authenticatedAttendee)
    } else {
      // Show login dialog if attendees are available but none is authenticated
      if (attendees && attendees.length > 0 && !currentAttendee) {
        console.log('No authenticated attendee, showing login dialog')
        setLoginDialogOpen(true)
      }
    }
  }, [attendees, currentAttendee])

  // Handle attendee login
  const handleAttendeeLogin = async (attendeeId, pin) => {
    try {
      const response = await api.auth.login(projectId, attendeeId, pin)
      console.log('Login successful:', response)
      
      // Store authentication
      api.auth.setCurrentAttendee(response.attendee, response.token)
      setCurrentAttendee(response.attendee)
      setLoginDialogOpen(false)
      
      showNotification(`Welcome, ${response.attendee.name}!`, 'success')
    } catch (error) {
      console.error('Login failed:', error)
      throw error // Re-throw for the dialog to handle
    }
  }

  // WebSocket event listeners
  useEffect(() => {
    if (!isConnected) return

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
        [payload.comparisonId]: true
      }))
    }

    const handleAttendeeJoined = (payload) => {
      setAttendeeStatus(prev => ({
        ...prev,
        [payload.attendeeId]: 'joined'
      }))
    }

    const handleAttendeeLeft = (payload) => {
      setAttendeeStatus(prev => ({
        ...prev,
        [payload.attendeeId]: 'left'
      }))
    }

    // Register listeners
    pairwiseWebSocketService.on('voteUpdate', handleVoteUpdate)
    pairwiseWebSocketService.on('consensusReached', handleConsensusReached)
    pairwiseWebSocketService.on('attendeeJoined', handleAttendeeJoined)
    pairwiseWebSocketService.on('attendeeLeft', handleAttendeeLeft)

    return () => {
      // Cleanup listeners
      pairwiseWebSocketService.off('voteUpdate', handleVoteUpdate)
      pairwiseWebSocketService.off('consensusReached', handleConsensusReached)
      pairwiseWebSocketService.off('attendeeJoined', handleAttendeeJoined)
      pairwiseWebSocketService.off('attendeeLeft', handleAttendeeLeft)
    }
  }, [isConnected])

  const [comparisons, setComparisons] = useState([])

  // Load comparisons from backend API
  useEffect(() => {
    const loadComparisons = async () => {
      if (!projectId || !sessionStarted) return

      try {
        console.log('Loading comparisons for project:', projectId)
        const comparisonsData = await pairwiseService.getComparisons(projectId)
        console.log('Loaded comparisons:', comparisonsData)
        
        if (comparisonsData && comparisonsData.comparisons) {
          setComparisons(comparisonsData.comparisons)
          
          // Set initial comparison if none selected
          if (!currentComparison && comparisonsData.comparisons.length > 0) {
            const firstIncomplete = comparisonsData.comparisons.find(c => !c.consensus_reached)
            setCurrentComparison(firstIncomplete || comparisonsData.comparisons[0])
          }
        } else {
          console.log('No comparisons found, using fallback generation')
          // Fallback: generate comparisons locally if API doesn't have them
          generateLocalComparisons()
        }
      } catch (error) {
        console.error('Failed to load comparisons:', error)
        console.log('Using fallback generation due to error')
        // Fallback to local generation if API fails
        generateLocalComparisons()
      }
    }

    const generateLocalComparisons = () => {
      if (!features || features.length < 2) {
        setComparisons([])
        return
      }

      // Generate all pairwise comparisons locally as fallback
      const pairs = []
      for (let i = 0; i < features.length; i++) {
        for (let j = i + 1; j < features.length; j++) {
          const comparisonId = `${features[i].id}-${features[j].id}`
          pairs.push({
            id: comparisonId,
            featureA: features[i],
            featureB: features[j],
            votes: votes[comparisonId] || {},
            hasConsensus: consensus[comparisonId] || false,
            index: pairs.length
          })
        }
      }

      setComparisons(pairs)
      
      // Set initial comparison if none selected
      if (!currentComparison && pairs.length > 0) {
        const firstIncomplete = pairs.find(c => !c.hasConsensus)
        setCurrentComparison(firstIncomplete || pairs[0])
      }
    }

    loadComparisons()
  }, [projectId, sessionStarted, features])

  // Initialize pairwise session
  useEffect(() => {
    const initializeSession = async () => {
      if (!projectId || !features || features.length < 2) return

      try {
        // Try to get existing session first
        const existingSession = await pairwiseService.getSession(projectId)
        console.log('Existing session found:', existingSession)
        setSessionStarted(true)
      } catch (error) {
        console.log('No existing session, creating new one...')
        try {
          // Create new complexity session
          const newSession = await pairwiseService.startSession(projectId, {
            criterion_type: 'complexity'
          })
          console.log('Created new session:', newSession)
          setSessionStarted(true)
        } catch (sessionError) {
          console.error('Failed to create session:', sessionError)
          // Check if session was created anyway (backend may return error even on partial success)
          try {
            const sessionAfterError = await pairwiseService.getSession(projectId)
            console.log('Session exists despite creation error:', sessionAfterError)
            setSessionStarted(true)
          } catch (finalError) {
            console.error('No session found after creation attempt:', finalError)
            setNotification({
              open: true,
              message: 'Failed to start pairwise session',
              severity: 'error'
            })
          }
        }
      }
    }

    initializeSession()
  }, [projectId, features])

  // Join session on component mount
  useEffect(() => {
    if (projectId && isConnected && sessionStarted) {
      joinSession(projectId)
    }
    
    return () => {
      if (projectId) {
        leaveSession(projectId)
      }
    }
  }, [projectId, isConnected, sessionStarted, joinSession, leaveSession])

  // Navigation functions
  const navigateToNext = () => {
    if (!currentComparison || comparisons.length === 0) return
    
    const currentIndex = comparisons.findIndex(c => c.id === currentComparison.id)
    const nextIndex = (currentIndex + 1) % comparisons.length
    setCurrentComparison(comparisons[nextIndex])
  }

  const navigateToPrevious = () => {
    if (!currentComparison || comparisons.length === 0) return
    
    const currentIndex = comparisons.findIndex(c => c.id === currentComparison.id)
    const prevIndex = currentIndex === 0 ? comparisons.length - 1 : currentIndex - 1
    setCurrentComparison(comparisons[prevIndex])
  }

  // Vote handling
  const handleVote = async (comparisonId, choice, attendeeId) => {
    try {
      // Try to send vote via WebSocket for real-time updates (optional)
      try {
        sendVote(comparisonId, choice, attendeeId)
        console.log('WebSocket vote sent successfully')
      } catch (wsError) {
        console.warn('WebSocket vote failed, continuing with API call:', wsError.message)
      }
      
      // Always persist to backend API (required)
      await pairwiseService.submitVote(projectId, comparisonId, attendeeId, choice)
      console.log('API vote submitted successfully')
      
      showNotification('Vote submitted successfully', 'success')
      
      // Auto-advance to next comparison if enabled
      if (autoAdvance && !currentComparison?.hasConsensus) {
        setTimeout(() => {
          const nextIncomplete = comparisons.find(c => 
            c.id !== comparisonId && !c.hasConsensus
          )
          if (nextIncomplete) {
            setCurrentComparison(nextIncomplete)
          }
        }, 1000)
      }
    } catch (error) {
      console.error('Failed to submit vote:', error)
      showNotification('Failed to submit vote', 'error')
    }
  }

  // Utility functions
  const showNotification = (message, severity = 'info') => {
    setNotification({ open: true, message, severity })
  }

  const closeNotification = () => {
    setNotification({ ...notification, open: false })
  }

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen)
  }

  // Keyboard shortcuts using custom hook
  useKeyboardShortcuts({
    onVoteA: () => {
      if (currentComparison && currentAttendee) {
        handleVote(currentComparison.id, 'A', currentAttendee.id)
      }
    },
    onVoteNeutral: () => {
      if (currentComparison && currentAttendee) {
        handleVote(currentComparison.id, 'neutral', currentAttendee.id)
      }
    },
    onVoteB: () => {
      if (currentComparison && currentAttendee) {
        handleVote(currentComparison.id, 'B', currentAttendee.id)
      }
    },
    onNext: navigateToNext,
    onPrevious: navigateToPrevious,
    onToggleView: () => setViewMode(viewMode === 'grid' ? 'detail' : 'grid'),
    onToggleFullscreen: toggleFullscreen,
    onHelp: () => setShortcutsDialogOpen(true),
    onEscape: () => {
      if (isFullscreen) {
        setIsFullscreen(false)
      } else if (shortcutsDialogOpen) {
        setShortcutsDialogOpen(false)
      }
    },
    enabled: sessionStarted
  })

  const handleExportResults = async () => {
    try {
      const results = await pairwiseService.exportResults(projectId)
      // Create and download CSV file
      const blob = new Blob([results], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `pairwise-results-${project?.name || projectId}.csv`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
      
      showNotification('Results exported successfully', 'success')
    } catch (error) {
      console.error('Failed to export results:', error)
      showNotification('Failed to export results', 'error')
    }
  }

  // Loading and error states
  if (projectLoading || featuresLoading || attendeesLoading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography>Loading comparison session...</Typography>
      </Container>
    )
  }

  if (projectError) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="error">
          Failed to load project: {projectError.message}
        </Alert>
      </Container>
    )
  }

  if (!features || features.length < 2) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Alert severity="warning">
          At least 2 features are required for pairwise comparison.
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

  const currentIndex = currentComparison ? 
    comparisons.findIndex(c => c.id === currentComparison.id) : 0
  const totalComparisons = comparisons.length
  const completedComparisons = comparisons.filter(c => c.hasConsensus).length

  return (
    <Box 
      ref={containerRef}
      sx={{ 
        height: isFullscreen ? '100vh' : 'auto',
        bgcolor: 'background.default',
        overflow: isFullscreen ? 'hidden' : 'auto'
      }}
    >
      {/* Header */}
      {!isFullscreen && (
        <AppBar position="static" color="transparent" elevation={1}>
          <Toolbar>
            <IconButton
              edge="start"
              onClick={() => navigate(`/projects/${projectId}/features`)}
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
                Pairwise Comparison
              </Typography>
            </Breadcrumbs>

            <Box sx={{ display: 'flex', gap: 1 }}>
              <IconButton
                onClick={() => setViewMode(viewMode === 'grid' ? 'detail' : 'grid')}
                title={`Switch to ${viewMode === 'grid' ? 'detail' : 'grid'} view`}
              >
                {viewMode === 'grid' ? <ViewList /> : <GridView />}
              </IconButton>
              
              <IconButton onClick={toggleFullscreen} title="Toggle fullscreen">
                {isFullscreen ? <FullscreenExit /> : <Fullscreen />}
              </IconButton>
              
              <IconButton 
                onClick={() => setShortcutsDialogOpen(true)} 
                title="Keyboard shortcuts (? or h)"
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
              >
                Export
              </Button>

              {completedComparisons === totalComparisons && totalComparisons > 0 && (
                <Button
                  variant="contained"
                  startIcon={<NavigateNext />}
                  onClick={() => navigate(`/projects/${projectId}/scoring/value`)}
                  size="small"
                  color="success"
                >
                  Start Value Scoring
                </Button>
              )}
            </Box>
          </Toolbar>
        </AppBar>
      )}

      {/* Main Content */}
      <Container 
        maxWidth={isFullscreen ? false : "xl"} 
        sx={{ 
          py: isFullscreen ? 2 : 4,
          height: isFullscreen ? 'calc(100vh - 64px)' : 'auto'
        }}
      >
        <Grid container spacing={3} sx={{ height: '100%' }}>
          {/* Left Panel - Progress and Controls */}
          <Grid item xs={12} md={viewMode === 'grid' ? 12 : 3}>
            <SessionProgress
              totalComparisons={totalComparisons}
              completedComparisons={completedComparisons}
              attendees={attendees}
              attendeeStatus={attendeeStatus}
              isConnected={isConnected}
              currentComparison={currentComparison}
              onAttendeeSelect={setCurrentAttendee}
            />
          </Grid>

          {/* Main Panel - Grid or Detail View */}
          {viewMode === 'grid' ? (
            <Grid item xs={12}>
              <PairwiseGrid
                projectId={projectId}
                features={features}
                attendees={attendees}
                comparisons={comparisons}
                currentComparison={currentComparison}
                onComparisonSelect={(comparison) => {
                  setCurrentComparison(comparison)
                  setViewMode('detail')
                }}
                onVoteSubmit={handleVote}
              />
            </Grid>
          ) : (
            <>
              <Grid item xs={12} md={6}>
                {currentComparison && (
                  <ComparisonCard
                    comparison={currentComparison}
                    votes={currentComparison.votes}
                    attendees={attendees}
                    hasConsensus={currentComparison.hasConsensus}
                    currentAttendee={currentAttendee}
                    onVote={(comparisonId, choice) => {
                      console.log('PairwiseComparison onVote called:', {
                        comparisonId,
                        choice,
                        currentAttendee,
                        hasCurrentAttendee: !!currentAttendee,
                        isAuthenticated: api.auth.isAuthenticated()
                      })
                      
                      // Verify authentication
                      if (!api.auth.isAuthenticated()) {
                        console.error('Not authenticated - showing login dialog')
                        setLoginDialogOpen(true)
                        return
                      }
                      
                      if (currentAttendee) {
                        handleVote(comparisonId, choice, currentAttendee.id)
                      } else {
                        console.error('No attendee selected - cannot vote')
                        setLoginDialogOpen(true)
                      }
                    }}
                  />
                )}
              </Grid>

              {/* Right Panel - Voting */}
              <Grid item xs={12} md={3}>
                {currentComparison && currentAttendee && (
                  <AttendeeVotingPanel
                    comparison={currentComparison}
                    attendee={currentAttendee}
                    onVote={handleVote}
                    onNext={navigateToNext}
                    onPrevious={navigateToPrevious}
                    currentIndex={currentIndex}
                    totalComparisons={totalComparisons}
                    autoAdvance={autoAdvance}
                  />
                )}
              </Grid>
            </>
          )}
        </Grid>
      </Container>

      {/* Navigation FABs */}
      {viewMode === 'detail' && (
        <Box sx={{ position: 'fixed', bottom: 24, right: 24, display: 'flex', gap: 1 }}>
          <Fab
            size="small"
            onClick={navigateToPrevious}
            disabled={comparisons.length === 0}
          >
            <NavigateBefore />
          </Fab>
          <Fab
            size="small"
            onClick={navigateToNext}
            disabled={comparisons.length === 0}
          >
            <NavigateNext />
          </Fab>
        </Box>
      )}

      {/* Settings Dialog */}
      <Dialog open={settingsOpen} onClose={() => setSettingsOpen(false)}>
        <DialogTitle>Session Settings</DialogTitle>
        <DialogContent>
          <FormControlLabel
            control={
              <Switch
                checked={autoAdvance}
                onChange={(e) => setAutoAdvance(e.target.checked)}
              />
            }
            label="Auto-advance after voting"
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

      {/* Attendee Login Dialog */}
      <AttendeeLoginDialog
        open={loginDialogOpen}
        onClose={() => setLoginDialogOpen(false)}
        onLogin={handleAttendeeLogin}
        attendees={attendees || []}
        loading={attendeesLoading}
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

export default PairwiseComparison