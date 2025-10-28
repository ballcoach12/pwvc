import {
    CheckCircle,
    KeyboardArrowDown,
    KeyboardArrowUp,
    NavigateBefore,
    NavigateNext
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Chip,
    IconButton,
    LinearProgress,
    Tooltip,
    Typography
} from '@mui/material'
import { useEffect, useState } from 'react'

const AttendeeVotingPanel = ({
  comparisons = [],
  currentAttendee,
  onVote,
  selectedComparisonId,
  onSelectComparison,
  criterion = 'value',
}) => {
  const [currentIndex, setCurrentIndex] = useState(0)

  // Filter and organize comparisons
  const votedComparisons = comparisons.filter(c => 
    c.votes?.some(v => v.attendeeId === currentAttendee?.id)
  )
  const unvotedComparisons = comparisons.filter(c => 
    !c.votes?.some(v => v.attendeeId === currentAttendee?.id) && c.status !== 'completed'
  )
  const allUserComparisons = [...unvotedComparisons, ...votedComparisons]

  // Find current comparison
  const currentComparison = selectedComparisonId 
    ? comparisons.find(c => c.id === selectedComparisonId)
    : allUserComparisons[currentIndex]

  const currentVote = currentComparison?.votes?.find(v => v.attendeeId === currentAttendee?.id)
  const hasVoted = Boolean(currentVote)

  // Update index when selectedComparisonId changes
  useEffect(() => {
    if (selectedComparisonId) {
      const index = allUserComparisons.findIndex(c => c.id === selectedComparisonId)
      if (index >= 0) {
        setCurrentIndex(index)
      }
    }
  }, [selectedComparisonId])

  const handleVote = (choice) => {
    if (!currentComparison) return
    onVote(currentComparison.id, choice)
    
    // Auto-advance to next unvoted comparison after voting
    setTimeout(() => {
      const nextUnvoted = unvotedComparisons.find(c => 
        c.id !== currentComparison.id && 
        !c.votes?.some(v => v.attendeeId === currentAttendee?.id)
      )
      if (nextUnvoted) {
        const nextIndex = allUserComparisons.findIndex(c => c.id === nextUnvoted.id)
        if (nextIndex >= 0) {
          setCurrentIndex(nextIndex)
          onSelectComparison?.(nextUnvoted.id)
        }
      }
    }, 500) // Small delay to show vote confirmation
  }

  const navigateToComparison = (direction) => {
    let newIndex
    if (direction === 'next') {
      newIndex = Math.min(currentIndex + 1, allUserComparisons.length - 1)
    } else {
      newIndex = Math.max(currentIndex - 1, 0)
    }
    
    setCurrentIndex(newIndex)
    onSelectComparison?.(allUserComparisons[newIndex]?.id)
  }

  const jumpToComparison = (index) => {
    setCurrentIndex(index)
    onSelectComparison?.(allUserComparisons[index]?.id)
  }

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyPress = (event) => {
      if (!currentComparison || currentComparison.status === 'completed') return

      switch (event.key) {
        case '1':
        case 'ArrowLeft':
          event.preventDefault()
          handleVote('A')
          break
        case '2':
        case ' ': // Space for tie
          event.preventDefault()
          handleVote('tie')
          break
        case '3':
        case 'ArrowRight':
          event.preventDefault()
          handleVote('B')
          break
        case 'ArrowUp':
          event.preventDefault()
          navigateToComparison('previous')
          break
        case 'ArrowDown':
          event.preventDefault()
          navigateToComparison('next')
          break
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [currentComparison, currentIndex])

  if (!currentComparison) {
    return (
      <Card>
        <CardContent>
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <CheckCircle color="success" sx={{ fontSize: 64, mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              All Comparisons Complete!
            </Typography>
            <Typography variant="body2" color="text.secondary">
              You've voted on all available comparisons for {criterion} assessment.
            </Typography>
          </Box>
        </CardContent>
      </Card>
    )
  }

  const progress = ((currentIndex + 1) / allUserComparisons.length) * 100

  return (
    <Card>
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6">
            Your Voting Progress
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="body2" color="text.secondary">
              {currentIndex + 1} of {allUserComparisons.length}
            </Typography>
            <Chip
              label={`${votedComparisons.length}/${comparisons.length} voted`}
              size="small"
              color={votedComparisons.length === comparisons.length ? 'success' : 'primary'}
            />
          </Box>
        </Box>

        {/* Progress Bar */}
        <LinearProgress
          variant="determinate"
          value={progress}
          sx={{ mb: 2, height: 6, borderRadius: 3 }}
        />

        {/* Current Comparison Info */}
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle1" gutterBottom>
            Compare for {criterion === 'value' ? 'Business Value' : 'Implementation Complexity'}
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
            <Box sx={{ flex: 1, textAlign: 'center' }}>
              <Typography variant="h6" color="primary">
                {currentComparison.featureA.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {currentComparison.featureA.description}
              </Typography>
            </Box>
            
            <Typography variant="h4" color="text.secondary" sx={{ px: 2 }}>
              VS
            </Typography>
            
            <Box sx={{ flex: 1, textAlign: 'center' }}>
              <Typography variant="h6" color="primary">
                {currentComparison.featureB.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {currentComparison.featureB.description}
              </Typography>
            </Box>
          </Box>

          {/* Voting Status */}
          {hasVoted ? (
            <Alert severity="success" sx={{ mb: 2 }}>
              <Typography variant="body2">
                You voted: <strong>
                  {currentVote.choice === 'A' ? currentComparison.featureA.name :
                   currentVote.choice === 'B' ? currentComparison.featureB.name : 'Tie'}
                </strong>
              </Typography>
            </Alert>
          ) : (
            <Alert severity="info" sx={{ mb: 2 }}>
              <Typography variant="body2">
                {criterion === 'value' 
                  ? 'Which feature provides more business value?'
                  : 'Which feature is more complex to implement?'
                }
              </Typography>
            </Alert>
          )}
        </Box>

        {/* Voting Buttons */}
        {currentComparison.status !== 'completed' && (
          <Box sx={{ display: 'flex', gap: 1, mb: 3 }}>
            <Button
              variant={currentVote?.choice === 'A' ? 'contained' : 'outlined'}
              color="primary"
              onClick={() => handleVote('A')}
              fullWidth
              size="large"
              sx={{ py: 1.5 }}
            >
              <Box sx={{ textAlign: 'center' }}>
                <Typography variant="button" display="block">
                  {currentComparison.featureA.name}
                </Typography>
                <Typography variant="caption" display="block">
                  Press 1 or ←
                </Typography>
              </Box>
            </Button>

            <Button
              variant={currentVote?.choice === 'tie' ? 'contained' : 'outlined'}
              color="secondary"
              onClick={() => handleVote('tie')}
              sx={{ minWidth: 80, py: 1.5 }}
            >
              <Box sx={{ textAlign: 'center' }}>
                <Typography variant="button" display="block">
                  Tie
                </Typography>
                <Typography variant="caption" display="block">
                  Space
                </Typography>
              </Box>
            </Button>

            <Button
              variant={currentVote?.choice === 'B' ? 'contained' : 'outlined'}
              color="primary"
              onClick={() => handleVote('B')}
              fullWidth
              size="large"
              sx={{ py: 1.5 }}
            >
              <Box sx={{ textAlign: 'center' }}>
                <Typography variant="button" display="block">
                  {currentComparison.featureB.name}
                </Typography>
                <Typography variant="caption" display="block">
                  Press 3 or →
                </Typography>
              </Box>
            </Button>
          </Box>
        )}

        {/* Navigation */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <IconButton
            onClick={() => navigateToComparison('previous')}
            disabled={currentIndex === 0}
          >
            <NavigateBefore />
          </IconButton>

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Tooltip title="Previous comparison (↑)">
              <IconButton
                size="small"
                onClick={() => navigateToComparison('previous')}
                disabled={currentIndex === 0}
              >
                <KeyboardArrowUp />
              </IconButton>
            </Tooltip>
            <Tooltip title="Next comparison (↓)">
              <IconButton
                size="small"
                onClick={() => navigateToComparison('next')}
                disabled={currentIndex === allUserComparisons.length - 1}
              >
                <KeyboardArrowDown />
              </IconButton>
            </Tooltip>
          </Box>

          <IconButton
            onClick={() => navigateToComparison('next')}
            disabled={currentIndex === allUserComparisons.length - 1}
          >
            <NavigateNext />
          </IconButton>
        </Box>

        {/* Quick Jump to Unvoted */}
        {unvotedComparisons.length > 0 && hasVoted && (
          <Box sx={{ mt: 2, pt: 2, borderTop: 1, borderColor: 'divider' }}>
            <Button
              variant="outlined"
              fullWidth
              onClick={() => {
                const nextUnvoted = unvotedComparisons[0]
                const index = allUserComparisons.findIndex(c => c.id === nextUnvoted.id)
                jumpToComparison(index)
              }}
            >
              Jump to Next Unvoted ({unvotedComparisons.length} remaining)
            </Button>
          </Box>
        )}

        {/* Keyboard Shortcuts Help */}
        <Box sx={{ mt: 2, p: 1, bgcolor: 'grey.50', borderRadius: 1 }}>
          <Typography variant="caption" color="text.secondary">
            <strong>Keyboard shortcuts:</strong> 1/← = Left feature, Space = Tie, 3/→ = Right feature, ↑↓ = Navigate
          </Typography>
        </Box>
      </CardContent>
    </Card>
  )
}

export default AttendeeVotingPanel