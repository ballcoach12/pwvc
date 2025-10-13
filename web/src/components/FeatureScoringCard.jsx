import {
    CheckCircle,
    EmojiEvents,
    ExpandLess,
    ExpandMore,
    Group,
    RadioButtonUnchecked,
    Timer
} from '@mui/icons-material'
import {
    Alert,
    Avatar,
    Badge,
    Box,
    Card,
    CardActions,
    CardContent,
    Chip,
    Collapse,
    Divider,
    IconButton,
    LinearProgress,
    Tooltip,
    Typography
} from '@mui/material'
import { useEffect, useState } from 'react'
import FibonacciScalePicker from './FibonacciScalePicker'

/**
 * FeatureScoringCard displays a feature with its Fibonacci scoring interface.
 * Shows individual attendee scores, consensus status, and provides the scoring
 * interface for the current user.
 * 
 * Features:
 * - Feature details with title, description, and acceptance criteria
 * - Fibonacci scale picker for current attendee
 * - Display of all attendee scores with avatars
 * - Real-time consensus tracking and indicators
 * - Expandable details section
 * - Visual progress and status indicators
 * - Mobile-friendly with touch targets
 */
const FeatureScoringCard = ({
  feature,
  criterionType = 'value', // 'value' | 'complexity'
  currentAttendee,
  attendees = [],
  scores = {}, // { attendeeId: score }
  consensusScore = null,
  hasConsensus = false,
  onScoreSubmit,
  onScoreChange,
  disabled = false,
  expanded = false,
  onExpandToggle
}) => {
  const [localExpanded, setLocalExpanded] = useState(expanded)
  const [currentScore, setCurrentScore] = useState(null)

  // Initialize current score from props
  useEffect(() => {
    if (currentAttendee && scores[currentAttendee.id] !== undefined) {
      setCurrentScore(scores[currentAttendee.id])
    }
  }, [currentAttendee, scores])

  const handleScoreChange = (newScore) => {
    setCurrentScore(newScore)
    if (onScoreChange) {
      onScoreChange(feature.id, newScore, currentAttendee?.id)
    }
  }

  const handleScoreSubmit = () => {
    if (currentScore !== null && onScoreSubmit) {
      onScoreSubmit(feature.id, currentScore, currentAttendee?.id)
    }
  }

  const handleExpandClick = () => {
    const newExpanded = !localExpanded
    setLocalExpanded(newExpanded)
    if (onExpandToggle) {
      onExpandToggle(feature.id, newExpanded)
    }
  }

  // Calculate scoring progress
  const totalAttendes = attendees.length
  const scoredAttendees = Object.keys(scores || {}).length
  const progressPercentage = totalAttendes > 0 ? (scoredAttendees / totalAttendes) * 100 : 0

  // Get unique scores for consensus analysis
  const uniqueScores = [...new Set(Object.values(scores || {}))]
  const isConsensusClose = uniqueScores.length <= 2 && uniqueScores.length > 0

  // Format criterion type for display
  const criterionLabel = criterionType === 'value' ? 'Value' : 'Complexity'
  const criterionColor = criterionType === 'value' ? 'primary' : 'secondary'

  // Get attendee score display
  const getAttendeeScoreDisplay = (attendee) => {
    const score = scores[attendee.id]
    const hasScore = score !== undefined

    return (
      <Tooltip 
        key={attendee.id}
        title={`${attendee.name}: ${hasScore ? score : 'Not scored'}`}
        arrow
      >
        <Badge
          badgeContent={hasScore ? score : '?'}
          color={hasScore ? 'primary' : 'default'}
          sx={{
            '& .MuiBadge-badge': {
              fontSize: '0.7rem',
              minWidth: 18,
              height: 18
            }
          }}
        >
          <Avatar
            sx={{
              width: 32,
              height: 32,
              fontSize: '0.8rem',
              bgcolor: hasScore ? 'primary.main' : 'grey.300',
              color: hasScore ? 'primary.contrastText' : 'text.secondary'
            }}
          >
            {attendee.name?.charAt(0)?.toUpperCase() || '?'}
          </Avatar>
        </Badge>
      </Tooltip>
    )
  }

  return (
    <Card 
      sx={{ 
        position: 'relative',
        border: hasConsensus ? 2 : 1,
        borderColor: hasConsensus ? 'success.main' : 'grey.300',
        '&:hover': {
          boxShadow: 3
        }
      }}
    >
      {/* Consensus indicator ribbon */}
      {hasConsensus && (
        <Box 
          sx={{
            position: 'absolute',
            top: 0,
            right: 0,
            bgcolor: 'success.main',
            color: 'white',
            px: 1,
            py: 0.5,
            borderRadius: '0 0 0 8px',
            display: 'flex',
            alignItems: 'center',
            gap: 0.5,
            zIndex: 1
          }}
        >
          <EmojiEvents sx={{ fontSize: 16 }} />
          <Typography variant="caption" fontWeight="bold">
            Consensus
          </Typography>
        </Box>
      )}

      <CardContent sx={{ pb: 1 }}>
        {/* Feature Header */}
        <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 2, mb: 2 }}>
          <Box sx={{ flexGrow: 1 }}>
            <Typography variant="h6" gutterBottom>
              {feature.name}
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
              <Chip
                label={`${criterionLabel} Scoring`}
                size="small"
                color={criterionColor}
                variant="outlined"
              />
              <Chip
                label={`ID: ${feature.id}`}
                size="small"
                variant="outlined"
              />
            </Box>

            {feature.description && (
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                {feature.description}
              </Typography>
            )}
          </Box>

          {/* Status Icon */}
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            {hasConsensus ? (
              <CheckCircle color="success" fontSize="large" />
            ) : scoredAttendees > 0 ? (
              <Timer color="warning" fontSize="large" />
            ) : (
              <RadioButtonUnchecked color="disabled" fontSize="large" />
            )}
          </Box>
        </Box>

        {/* Consensus Score Display */}
        {hasConsensus && consensusScore && (
          <Alert 
            severity="success" 
            sx={{ mb: 2 }}
            icon={<EmojiEvents />}
          >
            <Typography variant="body2">
              <strong>Consensus Score: {consensusScore}</strong> - All attendees agree on this {criterionLabel.toLowerCase()} score
            </Typography>
          </Alert>
        )}

        {/* Progress and Attendee Scores */}
        <Box sx={{ mb: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
            <Typography variant="body2" color="text.secondary">
              Scoring Progress ({scoredAttendees}/{totalAttendes})
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              <Group fontSize="small" color="action" />
              <Typography variant="caption">
                {Math.round(progressPercentage)}%
              </Typography>
            </Box>
          </Box>
          
          <LinearProgress
            variant="determinate"
            value={progressPercentage}
            sx={{ 
              height: 6, 
              borderRadius: 3,
              bgcolor: 'grey.200',
              '& .MuiLinearProgress-bar': {
                bgcolor: hasConsensus ? 'success.main' : 
                        isConsensusClose ? 'warning.main' : 'primary.main'
              }
            }}
          />

          {/* Attendee Avatars with Scores */}
          <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap' }}>
            <Typography variant="caption" color="text.secondary">
              Scores:
            </Typography>
            <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
              {attendees.map(getAttendeeScoreDisplay)}
            </Box>
          </Box>
        </Box>

        {/* Current User Scoring Interface */}
        {currentAttendee && !hasConsensus && (
          <Box sx={{ mb: 2 }}>
            <FibonacciScalePicker
              value={currentScore}
              onChange={handleScoreChange}
              disabled={disabled}
              criterionType={criterionType}
              size="small"
              label={`Your ${criterionLabel} Score`}
              helperText={currentScore !== scores[currentAttendee.id] ? 
                'Click submit to save your score' : 
                'Score saved - you can change it until consensus is reached'
              }
              fullWidth
            />
          </Box>
        )}

        {/* Consensus Analysis */}
        {!hasConsensus && scoredAttendees > 1 && (
          <Box sx={{ mb: 2 }}>
            {isConsensusClose ? (
              <Alert severity="info" sx={{ py: 0.5 }}>
                <Typography variant="body2">
                  Close to consensus! Scores: {uniqueScores.join(', ')}
                </Typography>
              </Alert>
            ) : (
              <Alert severity="warning" sx={{ py: 0.5 }}>
                <Typography variant="body2">
                  Wide range of scores: {uniqueScores.sort((a, b) => a - b).join(', ')} - Discussion recommended
                </Typography>
              </Alert>
            )}
          </Box>
        )}
      </CardContent>

      {/* Expandable Details */}
      <CardActions sx={{ pt: 0 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
          {/* Submit Button */}
          {currentAttendee && !hasConsensus && currentScore !== scores[currentAttendee.id] && (
            <Box>
              <button
                onClick={handleScoreSubmit}
                disabled={currentScore === null || disabled}
                style={{
                  padding: '6px 16px',
                  backgroundColor: '#1976d2',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: disabled ? 'not-allowed' : 'pointer',
                  fontSize: '0.875rem'
                }}
              >
                Submit Score
              </button>
            </Box>
          )}

          {/* Expand Button */}
          <IconButton
            onClick={handleExpandClick}
            aria-expanded={localExpanded}
            aria-label="show more"
            size="small"
          >
            <Typography variant="caption" sx={{ mr: 0.5 }}>
              {localExpanded ? 'Less' : 'Details'}
            </Typography>
            {localExpanded ? <ExpandLess /> : <ExpandMore />}
          </IconButton>
        </Box>
      </CardActions>

      {/* Expandable Content */}
      <Collapse in={localExpanded} timeout="auto" unmountOnExit>
        <Divider />
        <CardContent sx={{ pt: 2 }}>
          {feature.acceptanceCriteria && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>
                Acceptance Criteria
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {feature.acceptanceCriteria}
              </Typography>
            </Box>
          )}

          {/* Detailed Score Breakdown */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Individual Scores
            </Typography>
            {attendees.map(attendee => (
              <Box 
                key={attendee.id}
                sx={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center',
                  py: 0.5
                }}
              >
                <Typography variant="body2">
                  {attendee.name}
                </Typography>
                <Chip
                  label={scores[attendee.id] !== undefined ? scores[attendee.id] : 'Not scored'}
                  size="small"
                  color={scores[attendee.id] !== undefined ? 'primary' : 'default'}
                  variant="outlined"
                />
              </Box>
            ))}
          </Box>
        </CardContent>
      </Collapse>
    </Card>
  )
}

export default FeatureScoringCard