import {
    Balance,
    CheckCircle,
    ExpandLess,
    ExpandMore
} from '@mui/icons-material'
import {
    Avatar,
    Box,
    Button,
    Card,
    CardContent,
    Chip,
    Collapse,
    IconButton,
    LinearProgress,
    Tooltip,
    Typography
} from '@mui/material'
import { useState } from 'react'

const ComparisonCard = ({
  comparison,
  currentAttendee,
  onVote,
  disabled = false,
  showDetails = false,
}) => {
  const [expanded, setExpanded] = useState(false)
  const { featureA, featureB, votes = [], status = 'pending', criterion } = comparison

  // Calculate vote statistics
  const totalAttendees = comparison.totalAttendees || 0
  const votedCount = votes.length
  const consensusReached = status === 'completed'
  const currentVote = votes.find(v => v.attendeeId === currentAttendee?.id)

  // Vote distribution
  const votesForA = votes.filter(v => v.choice === 'A').length
  const votesForB = votes.filter(v => v.choice === 'B').length
  const ties = votes.filter(v => v.choice === 'tie').length

  // Determine winner for consensus display
  const getWinner = () => {
    if (!consensusReached) return null
    if (votesForA > votesForB && votesForA > ties) return 'A'
    if (votesForB > votesForA && votesForB > ties) return 'B'
    return 'tie'
  }

  const winner = getWinner()

  const handleVote = (choice) => {
    if (disabled || consensusReached) return
    onVote(comparison.id, choice)
  }

  const getVoteButtonColor = (choice) => {
    if (currentVote?.choice === choice) {
      return 'primary'
    }
    return 'inherit'
  }

  const getVoteButtonVariant = (choice) => {
    if (currentVote?.choice === choice) {
      return 'contained'
    }
    return 'outlined'
  }

  const getFeatureCardStyle = (feature) => {
    const baseStyle = {
      flex: 1,
      p: 2,
      border: '2px solid',
      borderRadius: 2,
      transition: 'all 0.3s ease-in-out',
      cursor: disabled ? 'default' : 'pointer',
    }

    if (consensusReached && winner === (feature === featureA ? 'A' : 'B')) {
      return {
        ...baseStyle,
        borderColor: 'success.main',
        backgroundColor: 'success.50',
      }
    }

    if (currentVote?.choice === (feature === featureA ? 'A' : 'B')) {
      return {
        ...baseStyle,
        borderColor: 'primary.main',
        backgroundColor: 'primary.50',
      }
    }

    return {
      ...baseStyle,
      borderColor: 'grey.300',
      backgroundColor: 'background.paper',
      '&:hover': disabled ? {} : {
        borderColor: 'primary.main',
        backgroundColor: 'primary.25',
      },
    }
  }

  return (
    <Card sx={{ mb: 2, position: 'relative' }}>
      {/* Progress indicator */}
      {!consensusReached && (
        <LinearProgress
          variant="determinate"
          value={(votedCount / totalAttendees) * 100}
          sx={{ height: 4 }}
        />
      )}

      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Balance color={criterion === 'value' ? 'primary' : 'secondary'} />
            <Typography variant="h6">
              {criterion === 'value' ? 'Value Comparison' : 'Complexity Comparison'}
            </Typography>
            {consensusReached && (
              <Chip
                label="Consensus Reached"
                color="success"
                icon={<CheckCircle />}
                size="small"
              />
            )}
          </Box>

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="body2" color="text.secondary">
              {votedCount}/{totalAttendees} voted
            </Typography>
            <IconButton size="small" onClick={() => setExpanded(!expanded)}>
              {expanded ? <ExpandLess /> : <ExpandMore />}
            </IconButton>
          </Box>
        </Box>

        {/* Feature Comparison */}
        <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
          {/* Feature A */}
          <Box
            sx={getFeatureCardStyle(featureA)}
            onClick={() => handleVote('A')}
          >
            <Typography variant="h6" gutterBottom>
              {featureA.name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {featureA.description || 'No description provided'}
            </Typography>
            {votesForA > 0 && (
              <Chip
                label={`${votesForA} vote${votesForA !== 1 ? 's' : ''}`}
                size="small"
                sx={{ mt: 1 }}
              />
            )}
          </Box>

          {/* VS Divider */}
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minWidth: 60 }}>
            <Typography variant="h5" color="text.secondary" sx={{ fontWeight: 'bold' }}>
              VS
            </Typography>
            {ties > 0 && (
              <Chip
                label={`${ties} tie${ties !== 1 ? 's' : ''}`}
                size="small"
                color={currentVote?.choice === 'tie' ? 'primary' : 'default'}
                sx={{ mt: 1 }}
              />
            )}
          </Box>

          {/* Feature B */}
          <Box
            sx={getFeatureCardStyle(featureB)}
            onClick={() => handleVote('B')}
          >
            <Typography variant="h6" gutterBottom>
              {featureB.name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {featureB.description || 'No description provided'}
            </Typography>
            {votesForB > 0 && (
              <Chip
                label={`${votesForB} vote${votesForB !== 1 ? 's' : ''}`}
                size="small"
                sx={{ mt: 1 }}
              />
            )}
          </Box>
        </Box>

        {/* Voting Buttons */}
        {!consensusReached && (
          <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center', mb: 2 }}>
            <Button
              variant={getVoteButtonVariant('A')}
              color={getVoteButtonColor('A')}
              onClick={() => handleVote('A')}
              disabled={disabled}
              size="small"
            >
              Choose {featureA.name}
            </Button>
            <Button
              variant={getVoteButtonVariant('tie')}
              color={getVoteButtonColor('tie')}
              onClick={() => handleVote('tie')}
              disabled={disabled}
              size="small"
            >
              Tie
            </Button>
            <Button
              variant={getVoteButtonVariant('B')}
              color={getVoteButtonColor('B')}
              onClick={() => handleVote('B')}
              disabled={disabled}
              size="small"
            >
              Choose {featureB.name}
            </Button>
          </Box>
        )}

        {/* Expanded Details */}
        <Collapse in={expanded}>
          <Box sx={{ pt: 2, borderTop: 1, borderColor: 'divider' }}>
            {/* Vote Details */}
            <Typography variant="subtitle2" gutterBottom>
              Voting Details
            </Typography>
            <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
              <Typography variant="body2">
                {featureA.name}: {votesForA} vote{votesForA !== 1 ? 's' : ''}
              </Typography>
              <Typography variant="body2">
                Tie: {ties} vote{ties !== 1 ? 's' : ''}
              </Typography>
              <Typography variant="body2">
                {featureB.name}: {votesForB} vote{votesForB !== 1 ? 's' : ''}
              </Typography>
            </Box>

            {/* Attendee Votes */}
            <Typography variant="subtitle2" gutterBottom>
              Individual Votes
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
              {votes.map((vote) => (
                <Tooltip key={vote.attendeeId} title={`${vote.attendeeName}: ${vote.choice === 'A' ? featureA.name : vote.choice === 'B' ? featureB.name : 'Tie'}`}>
                  <Chip
                    avatar={<Avatar sx={{ width: 24, height: 24 }}>{vote.attendeeName?.[0]}</Avatar>}
                    label={vote.choice === 'A' ? 'A' : vote.choice === 'B' ? 'B' : 'T'}
                    size="small"
                    color={vote.choice === 'tie' ? 'default' : 'primary'}
                  />
                </Tooltip>
              ))}
            </Box>

            {/* Current User Vote Status */}
            <Box sx={{ mt: 2 }}>
              {currentVote ? (
                <Typography variant="body2" color="primary">
                  Your vote: {currentVote.choice === 'A' ? featureA.name : currentVote.choice === 'B' ? featureB.name : 'Tie'}
                </Typography>
              ) : (
                <Typography variant="body2" color="warning.main">
                  You haven't voted yet
                </Typography>
              )}
            </Box>
          </Box>
        </Collapse>
      </CardContent>
    </Card>
  )
}

export default ComparisonCard