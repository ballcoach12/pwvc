import {
    CheckCircle,
    Group,
    RadioButtonUnchecked,
    Timer
} from '@mui/icons-material'
import {
    Badge,
    Box,
    Card,
    CardContent,
    Chip,
    Grid,
    LinearProgress,
    Paper,
    Typography
} from '@mui/material'
import { useCallback, useEffect, useState } from 'react'
import { useWebSocket } from '../services/websocketService'

/**
 * PairwiseGrid displays all pairwise feature comparisons in a matrix layout
 * with real-time voting status, progress tracking, and interactive navigation.
 * 
 * Features:
 * - Grid matrix showing all feature pairs
 * - Real-time vote status updates via WebSocket
 * - Progress tracking for each comparison
 * - Click navigation to specific comparisons
 * - Consensus indicators and attendee participation
 * - Responsive layout adapting to screen size
 */
const PairwiseGrid = ({ 
  projectId, 
  features = [], 
  attendees = [],
  currentComparison = null,
  onComparisonSelect,
  onVoteSubmit 
}) => {
  const [comparisons, setComparisons] = useState([])
  const [gridData, setGridData] = useState([])
  const [loading, setLoading] = useState(true)
  
  const {
    isConnected,
    votes,
    consensus,
    attendeeStatus,
    sendVote,
    joinSession
  } = useWebSocket(projectId)

  // Generate all pairwise comparisons from features
  const generateComparisons = useCallback(() => {
    const pairs = []
    for (let i = 0; i < features.length; i++) {
      for (let j = i + 1; j < features.length; j++) {
        pairs.push({
          id: `${features[i].id}-${features[j].id}`,
          featureA: features[i],
          featureB: features[j],
          position: { row: i, col: j },
          votes: votes[`${features[i].id}-${features[j].id}`] || {},
          hasConsensus: consensus[`${features[i].id}-${features[j].id}`] || false,
          totalVotes: Object.keys(votes[`${features[i].id}-${features[j].id}`] || {}).length,
          requiredVotes: attendees.length
        })
      }
    }
    return pairs
  }, [features, votes, consensus, attendees.length])

  // Create grid data structure for matrix display
  const createGridData = useCallback(() => {
    const grid = []
    for (let i = 0; i < features.length; i++) {
      const row = []
      for (let j = 0; j < features.length; j++) {
        if (i === j) {
          // Diagonal cell - feature vs itself
          row.push({
            type: 'feature',
            feature: features[i],
            position: { row: i, col: j }
          })
        } else if (i < j) {
          // Upper triangle - actual comparisons
          const comparison = comparisons.find(c => 
            c.position.row === i && c.position.col === j
          )
          row.push({
            type: 'comparison',
            comparison,
            position: { row: i, col: j }
          })
        } else {
          // Lower triangle - mirror/empty
          row.push({
            type: 'empty',
            position: { row: i, col: j }
          })
        }
      }
      grid.push(row)
    }
    return grid
  }, [features, comparisons])

  // Initialize comparisons and grid
  useEffect(() => {
    if (features.length > 0) {
      const newComparisons = generateComparisons()
      setComparisons(newComparisons)
      
      const newGridData = createGridData()
      setGridData(newGridData)
      
      setLoading(false)
    }
  }, [generateComparisons, createGridData, features])

  // Join WebSocket session
  useEffect(() => {
    if (projectId && isConnected) {
      joinSession(projectId)
    }
  }, [projectId, isConnected, joinSession])

  // Handle comparison cell click
  const handleComparisonClick = (comparison) => {
    if (comparison && onComparisonSelect) {
      onComparisonSelect(comparison)
    }
  }

  // Handle vote submission
  const handleVote = (comparisonId, choice, attendeeId) => {
    sendVote(comparisonId, choice, attendeeId)
    if (onVoteSubmit) {
      onVoteSubmit(comparisonId, choice, attendeeId)
    }
  }

  // Calculate overall progress
  const calculateProgress = () => {
    if (comparisons.length === 0) return 0
    const completedComparisons = comparisons.filter(c => c.hasConsensus).length
    return (completedComparisons / comparisons.length) * 100
  }

  // Get comparison status for styling
  const getComparisonStatus = (comparison) => {
    if (!comparison) return 'empty'
    if (comparison.hasConsensus) return 'completed'
    if (comparison.totalVotes > 0) return 'in-progress'
    return 'pending'
  }

  // Render feature header cell
  const renderFeatureCell = (feature, position) => (
    <Card
      key={`feature-${position.row}-${position.col}`}
      sx={{
        minHeight: 80,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'primary.main',
        color: 'primary.contrastText',
        cursor: 'default'
      }}
    >
      <CardContent sx={{ p: 1, textAlign: 'center' }}>
        <Typography variant="body2" fontWeight="bold">
          {feature.name}
        </Typography>
        <Typography variant="caption" sx={{ opacity: 0.8 }}>
          ID: {feature.id}
        </Typography>
      </CardContent>
    </Card>
  )

  // Render comparison cell
  const renderComparisonCell = (comparison, position) => {
    const status = getComparisonStatus(comparison)
    const isCurrentComparison = currentComparison?.id === comparison?.id
    
    return (
      <Card
        key={`comparison-${position.row}-${position.col}`}
        sx={{
          minHeight: 80,
          cursor: 'pointer',
          border: isCurrentComparison ? 2 : 1,
          borderColor: isCurrentComparison ? 'primary.main' : 'grey.300',
          bgcolor: status === 'completed' ? 'success.light' : 
                   status === 'in-progress' ? 'warning.light' : 'grey.50',
          '&:hover': {
            boxShadow: 3,
            transform: 'translateY(-2px)',
            transition: 'all 0.2s ease-in-out'
          }
        }}
        onClick={() => handleComparisonClick(comparison)}
      >
        <CardContent sx={{ p: 1, textAlign: 'center' }}>
          {comparison ? (
            <>
              <Box sx={{ display: 'flex', justifyContent: 'center', mb: 0.5 }}>
                {status === 'completed' ? (
                  <CheckCircle color="success" fontSize="small" />
                ) : status === 'in-progress' ? (
                  <Timer color="warning" fontSize="small" />
                ) : (
                  <RadioButtonUnchecked color="disabled" fontSize="small" />
                )}
              </Box>
              
              <Typography variant="caption" display="block">
                {comparison.featureA.name} vs {comparison.featureB.name}
              </Typography>
              
              <Box sx={{ mt: 0.5 }}>
                <Badge badgeContent={comparison.totalVotes} color="primary">
                  <Group fontSize="small" color="action" />
                </Badge>
              </Box>
              
              {comparison.totalVotes > 0 && (
                <LinearProgress
                  variant="determinate"
                  value={(comparison.totalVotes / comparison.requiredVotes) * 100}
                  sx={{ mt: 0.5, height: 4 }}
                  color={status === 'completed' ? 'success' : 'primary'}
                />
              )}
            </>
          ) : (
            <Typography variant="caption" color="text.disabled">
              No comparison
            </Typography>
          )}
        </CardContent>
      </Card>
    )
  }

  // Render empty cell
  const renderEmptyCell = (position) => (
    <Card
      key={`empty-${position.row}-${position.col}`}
      sx={{
        minHeight: 80,
        bgcolor: 'grey.100',
        border: '1px dashed',
        borderColor: 'grey.300'
      }}
    >
      <CardContent sx={{ p: 1, textAlign: 'center' }}>
        <Typography variant="caption" color="text.disabled">
          —
        </Typography>
      </CardContent>
    </Card>
  )

  if (loading) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <LinearProgress />
        <Typography variant="body2" sx={{ mt: 2 }}>
          Generating pairwise comparisons...
        </Typography>
      </Box>
    )
  }

  if (features.length < 2) {
    return (
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="h6" color="text.secondary">
          Need at least 2 features for pairwise comparison
        </Typography>
        <Typography variant="body2" sx={{ mt: 1 }}>
          Add more features to begin the comparison process.
        </Typography>
      </Paper>
    )
  }

  return (
    <Box sx={{ p: 2 }}>
      {/* Header with progress */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h6">
            Pairwise Comparison Matrix
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Chip
              icon={<Group />}
              label={`${Object.keys(attendeeStatus).length} attendees`}
              color={isConnected ? 'success' : 'error'}
              variant="outlined"
              size="small"
            />
            <Typography variant="body2" color="text.secondary">
              {Math.round(calculateProgress())}% Complete
            </Typography>
          </Box>
        </Box>
        
        <LinearProgress
          variant="determinate"
          value={calculateProgress()}
          sx={{ height: 8, borderRadius: 4 }}
          color="primary"
        />
        
        <Typography variant="caption" display="block" sx={{ mt: 1, textAlign: 'center' }}>
          Click on any comparison cell to view details and vote
        </Typography>
      </Paper>

      {/* Comparison Matrix Grid */}
      <Paper sx={{ p: 2, overflow: 'auto' }}>
        <Grid container spacing={1}>
          {gridData.map((row, rowIndex) => (
            <Grid item xs={12} key={`row-${rowIndex}`}>
              <Grid container spacing={1}>
                {row.map((cell, colIndex) => (
                  <Grid item key={`cell-${rowIndex}-${colIndex}`} xs={12/features.length}>
                    {cell.type === 'feature' && renderFeatureCell(cell.feature, cell.position)}
                    {cell.type === 'comparison' && renderComparisonCell(cell.comparison, cell.position)}
                    {cell.type === 'empty' && renderEmptyCell(cell.position)}
                  </Grid>
                ))}
              </Grid>
            </Grid>
          ))}
        </Grid>
      </Paper>

      {/* Legend */}
      <Paper sx={{ p: 2, mt: 2 }}>
        <Typography variant="subtitle2" gutterBottom>
          Legend
        </Typography>
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <CheckCircle color="success" fontSize="small" />
            <Typography variant="caption">Consensus Reached</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Timer color="warning" fontSize="small" />
            <Typography variant="caption">Voting in Progress</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <RadioButtonUnchecked color="disabled" fontSize="small" />
            <Typography variant="caption">Not Started</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <Group fontSize="small" color="action" />
            <Typography variant="caption">Vote Count</Typography>
          </Box>
        </Box>
      </Paper>
    </Box>
  )
}

export default PairwiseGrid