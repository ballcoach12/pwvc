import {
    FilterList,
    Refresh,
    Sort,
    UnfoldLess,
    UnfoldMore,
    ViewList,
    ViewModule
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    ButtonGroup,
    Fab,
    FormControl,
    FormControlLabel,
    Grid,
    InputLabel,
    LinearProgress,
    MenuItem,
    Paper,
    Select,
    Switch,
    Tooltip,
    Typography
} from '@mui/material'
import { useCallback, useEffect, useState } from 'react'
import { useWebSocket } from '../services/websocketService'
import FeatureScoringCard from './FeatureScoringCard'

/**
 * FibonacciScoringGrid displays all features in a grid layout for Fibonacci scoring.
 * Provides filtering, sorting, and real-time collaboration features with WebSocket
 * integration for live score updates and consensus tracking.
 * 
 * Features:
 * - Responsive grid layout with customizable columns
 * - Real-time score updates via WebSocket
 * - Filtering by consensus status and attendee participation
 * - Sorting by various criteria (name, progress, consensus)
 * - Bulk expand/collapse for feature details
 * - Mobile-friendly with touch interactions
 * - Progress tracking and visual indicators
 */
const FibonacciScoringGrid = ({
  projectId,
  sessionId,
  features = [],
  attendees = [],
  criterionType = 'value',
  currentAttendee,
  individualScores = {}, // { featureId: { attendeeId: score } }
  consensusScores = {}, // { featureId: consensusScore }
  onScoreSubmit,
  onScoreChange,
  loading = false,
  error = null
}) => {
  // State management
  const [viewMode, setViewMode] = useState('grid') // 'grid' | 'list'
  const [filterMode, setFilterMode] = useState('all') // 'all' | 'pending' | 'completed' | 'in-progress'
  const [sortBy, setSortBy] = useState('name') // 'name' | 'progress' | 'consensus' | 'id'
  const [sortOrder, setSortOrder] = useState('asc') // 'asc' | 'desc'
  const [expandedFeatures, setExpandedFeatures] = useState(new Set())
  const [showOnlyMyTurn, setShowOnlyMyTurn] = useState(false)
  const [gridColumns, setGridColumns] = useState(2)

  // WebSocket integration for real-time updates
  const {
    isConnected,
    fibonacciScores,
    fibonacciConsensus,
    sendScore,
    joinSession
  } = useWebSocket(projectId)

  // Join WebSocket session
  useEffect(() => {
    if (projectId && sessionId && isConnected) {
      joinSession(`fibonacci-${sessionId}`)
    }
  }, [projectId, sessionId, isConnected, joinSession])

  // Filter and sort features
  const processedFeatures = useCallback(() => {
    let filtered = [...features]

    // Apply filters
    switch (filterMode) {
      case 'pending':
        filtered = filtered.filter(feature => 
          !consensusScores[feature.id] && 
          Object.keys(individualScores[feature.id] || {}).length === 0
        )
        break
      case 'in-progress':
        filtered = filtered.filter(feature => 
          !consensusScores[feature.id] && 
          Object.keys(individualScores[feature.id] || {}).length > 0
        )
        break
      case 'completed':
        filtered = filtered.filter(feature => consensusScores[feature.id] !== undefined)
        break
      case 'all':
      default:
        // No filtering
        break
    }

    // Filter for "my turn" - features where current user hasn't scored yet
    if (showOnlyMyTurn && currentAttendee) {
      filtered = filtered.filter(feature => 
        !consensusScores[feature.id] && 
        !individualScores[feature.id]?.[currentAttendee.id]
      )
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let compareValue = 0

      switch (sortBy) {
        case 'name':
          compareValue = a.name.localeCompare(b.name)
          break
        case 'progress':
          const aProgress = Object.keys(individualScores[a.id] || {}).length
          const bProgress = Object.keys(individualScores[b.id] || {}).length
          compareValue = aProgress - bProgress
          break
        case 'consensus':
          const aHasConsensus = consensusScores[a.id] !== undefined
          const bHasConsensus = consensusScores[b.id] !== undefined
          compareValue = aHasConsensus === bHasConsensus ? 0 : aHasConsensus ? 1 : -1
          break
        case 'id':
        default:
          compareValue = a.id - b.id
          break
      }

      return sortOrder === 'asc' ? compareValue : -compareValue
    })

    return filtered
  }, [features, filterMode, sortBy, sortOrder, individualScores, consensusScores, showOnlyMyTurn, currentAttendee])

  const filteredFeatures = processedFeatures()

  // Handle feature expansion
  const handleFeatureExpand = (featureId, expanded) => {
    setExpandedFeatures(prev => {
      const newSet = new Set(prev)
      if (expanded) {
        newSet.add(featureId)
      } else {
        newSet.delete(featureId)
      }
      return newSet
    })
  }

  // Bulk expand/collapse
  const handleExpandAll = () => {
    setExpandedFeatures(new Set(filteredFeatures.map(f => f.id)))
  }

  const handleCollapseAll = () => {
    setExpandedFeatures(new Set())
  }

  // Handle score submission with WebSocket integration
  const handleScoreSubmit = useCallback(async (featureId, score, attendeeId) => {
    try {
      // Send via WebSocket for real-time updates
      if (sessionId) {
        sendScore(sessionId, featureId, attendeeId, score)
      }
      
      // Also call parent handler for persistence
      if (onScoreSubmit) {
        await onScoreSubmit(featureId, score, attendeeId)
      }
    } catch (error) {
      console.error('Failed to submit score:', error)
    }
  }, [sessionId, sendScore, onScoreSubmit])

  // Calculate grid layout
  const getGridColumns = () => {
    if (viewMode === 'list') return 1
    return Math.min(gridColumns, filteredFeatures.length)
  }

  // Calculate overall progress
  const calculateProgress = () => {
    if (features.length === 0) return 0
    const completedCount = Object.keys(consensusScores).length
    return (completedCount / features.length) * 100
  }

  if (loading) {
    return (
      <Box sx={{ p: 3 }}>
        <LinearProgress />
        <Typography variant="body2" sx={{ mt: 2, textAlign: 'center' }}>
          Loading scoring session...
        </Typography>
      </Box>
    )
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ m: 3 }}>
        {error.message || 'Failed to load scoring session'}
      </Alert>
    )
  }

  if (features.length === 0) {
    return (
      <Alert severity="info" sx={{ m: 3 }}>
        No features available for scoring. Please add features to the project first.
      </Alert>
    )
  }

  return (
    <Box sx={{ p: 2 }}>
      {/* Controls Header */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h6">
            {criterionType === 'value' ? 'Value' : 'Complexity'} Scoring
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {filteredFeatures.length} of {features.length} features • {Math.round(calculateProgress())}% consensus
          </Typography>
        </Box>

        {/* Controls Row 1: View and Layout */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2, flexWrap: 'wrap' }}>
          <ButtonGroup size="small" variant="outlined">
            <Button
              onClick={() => setViewMode('grid')}
              variant={viewMode === 'grid' ? 'contained' : 'outlined'}
              startIcon={<ViewModule />}
            >
              Grid
            </Button>
            <Button
              onClick={() => setViewMode('list')}
              variant={viewMode === 'list' ? 'contained' : 'outlined'}
              startIcon={<ViewList />}
            >
              List
            </Button>
          </ButtonGroup>

          {viewMode === 'grid' && (
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <InputLabel>Columns</InputLabel>
              <Select
                value={gridColumns}
                label="Columns" 
                onChange={(e) => setGridColumns(e.target.value)}
              >
                <MenuItem value={1}>1</MenuItem>
                <MenuItem value={2}>2</MenuItem>
                <MenuItem value={3}>3</MenuItem>
                <MenuItem value={4}>4</MenuItem>
              </Select>
            </FormControl>
          )}

          <ButtonGroup size="small" variant="outlined">
            <Button
              onClick={handleExpandAll}
              startIcon={<UnfoldMore />}
              disabled={filteredFeatures.length === 0}
            >
              Expand All
            </Button>
            <Button
              onClick={handleCollapseAll}
              startIcon={<UnfoldLess />}
              disabled={filteredFeatures.length === 0}
            >
              Collapse All
            </Button>
          </ButtonGroup>
        </Box>

        {/* Controls Row 2: Filter and Sort */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Filter</InputLabel>
            <Select
              value={filterMode}
              label="Filter"
              onChange={(e) => setFilterMode(e.target.value)}
              startAdornment={<FilterList fontSize="small" sx={{ mr: 1 }} />}
            >
              <MenuItem value="all">All Features</MenuItem>
              <MenuItem value="pending">Not Started</MenuItem>
              <MenuItem value="in-progress">In Progress</MenuItem>
              <MenuItem value="completed">Completed</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel>Sort By</InputLabel>
            <Select
              value={sortBy}
              label="Sort By"
              onChange={(e) => setSortBy(e.target.value)}
              startAdornment={<Sort fontSize="small" sx={{ mr: 1 }} />}
            >
              <MenuItem value="name">Name</MenuItem>
              <MenuItem value="id">ID</MenuItem>
              <MenuItem value="progress">Progress</MenuItem>
              <MenuItem value="consensus">Consensus</MenuItem>
            </Select>
          </FormControl>

          <ButtonGroup size="small" variant="outlined">
            <Button
              onClick={() => setSortOrder('asc')}
              variant={sortOrder === 'asc' ? 'contained' : 'outlined'}
            >
              A-Z
            </Button>
            <Button
              onClick={() => setSortOrder('desc')}
              variant={sortOrder === 'desc' ? 'contained' : 'outlined'}
            >
              Z-A
            </Button>
          </ButtonGroup>

          {currentAttendee && (
            <FormControlLabel
              control={
                <Switch
                  checked={showOnlyMyTurn}
                  onChange={(e) => setShowOnlyMyTurn(e.target.checked)}
                  size="small"
                />
              }
              label="My turn only"
            />
          )}
        </Box>
      </Paper>

      {/* No Results Message */}
      {filteredFeatures.length === 0 && (
        <Alert severity="info" sx={{ mb: 3 }}>
          No features match the current filters. Try adjusting your filter settings.
        </Alert>
      )}

      {/* Features Grid */}
      <Grid container spacing={3}>
        {filteredFeatures.map(feature => (
          <Grid 
            item 
            xs={12} 
            sm={viewMode === 'list' ? 12 : 12/getGridColumns()}
            md={viewMode === 'list' ? 12 : 12/getGridColumns()}
            key={feature.id}
          >
            <FeatureScoringCard
              feature={feature}
              criterionType={criterionType}
              currentAttendee={currentAttendee}
              attendees={attendees}
              scores={individualScores[feature.id] || {}}
              consensusScore={consensusScores[feature.id]}
              hasConsensus={consensusScores[feature.id] !== undefined}
              onScoreSubmit={handleScoreSubmit}
              onScoreChange={onScoreChange}
              expanded={expandedFeatures.has(feature.id)}
              onExpandToggle={handleFeatureExpand}
            />
          </Grid>
        ))}
      </Grid>

      {/* Floating Action Button for Refresh */}
      <Tooltip title="Refresh scores" placement="left">
        <Fab
          color="primary"
          size="small"
          sx={{
            position: 'fixed',
            bottom: 24,
            right: 24,
            zIndex: 1000
          }}
          onClick={() => window.location.reload()}
        >
          <Refresh />
        </Fab>
      </Tooltip>
    </Box>
  )
}

export default FibonacciScoringGrid