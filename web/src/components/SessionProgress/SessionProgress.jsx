import {
  AccessTime,
  CheckCircle,
  ExpandLess,
  ExpandMore,
  Group,
  Pause,
  PlayArrow,
  Timeline,
} from '@mui/icons-material'
import {
  Avatar,
  Box,
  Button,
  Card,
  CardContent,
  Chip,
  Collapse,
  Divider,
  IconButton,
  LinearProgress,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Typography,
} from '@mui/material'
import React, { useEffect, useState } from 'react'

const SessionProgress = (props) => {
  // Extract props with maximum safety and flexibility to handle multiple interfaces
  const {
    session,
    attendees: propAttendees,
    comparisons: propComparisons,
    currentAttendee,
    onNextComparison,
    onPauseSession,
    onResumeSession,
    // Legacy/alternative props from PairwiseComparison component
    totalComparisons: propTotalComparisons,
    completedComparisons: propCompletedComparisons,
    attendeeStatus,
    isConnected,
    currentComparison,
    onAttendeeSelect,
    ...otherProps
  } = props

  // Ensure arrays are always arrays with maximum defensive programming
  // Handle both new interface (attendees/comparisons) and legacy interface (attendees only)
  const safeAttendees = Array.isArray(propAttendees) ? propAttendees : []
  const safeComparisons = Array.isArray(propComparisons) ? propComparisons : []
  
  // If we're using the legacy interface, create mock comparisons data
  const hasLegacyInterface = propTotalComparisons !== undefined || propCompletedComparisons !== undefined

  const [expanded, setExpanded] = useState(false)
  const [sessionTime, setSessionTime] = useState(0)

  // Calculate progress statistics - handle both interfaces
  const totalComparisons = hasLegacyInterface ? 
    (propTotalComparisons || 0) : 
    safeComparisons.length
    
  const completedComparisons = hasLegacyInterface ? 
    (propCompletedComparisons || 0) : 
    safeComparisons.filter(c => c.status === 'completed').length
    
  const overallProgress = totalComparisons > 0 ? (completedComparisons / totalComparisons) * 100 : 0

  // Calculate individual attendee progress
  const getAttendeeProgress = (attendeeId) => {
    if (hasLegacyInterface) {
      // For legacy interface, use attendeeStatus if available
      if (attendeeStatus && attendeeStatus[attendeeId]) {
        return attendeeStatus[attendeeId].progress || 0
      }
      // Fallback to estimated progress
      return overallProgress
    }
    
    if (safeComparisons.length === 0) {
      return 0
    }
    
    const attendeeVotes = safeComparisons.reduce((count, comparison) => {
      const hasVoted = Array.isArray(comparison.votes) && 
        comparison.votes.some(v => v.attendeeId === attendeeId)
      return hasVoted ? count + 1 : count
    }, 0)
    return totalComparisons > 0 ? (attendeeVotes / totalComparisons) * 100 : 0
  }

  // Get next comparison for current attendee
  const getNextComparison = () => {
    if (hasLegacyInterface) {
      // For legacy interface, use currentComparison prop
      return currentComparison || null
    }
    
    if (!currentAttendee?.id) {
      return null
    }
    
    return safeComparisons.find(c => 
      c.status !== 'completed' && 
      (!Array.isArray(c.votes) || !c.votes.some(v => v.attendeeId === currentAttendee.id))
    )
  }

  // Session timer
  useEffect(() => {
    let interval
    if (session?.status === 'active' && session?.startTime) {
      interval = setInterval(() => {
        const elapsed = Math.floor((Date.now() - new Date(session.startTime).getTime()) / 1000)
        setSessionTime(elapsed)
      }, 1000)
    }
    return () => clearInterval(interval)
  }, [session?.status, session?.startTime])

  const formatTime = (seconds) => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = seconds % 60
    
    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`
  }

  const getStatusColor = (status) => {
    switch (status) {
      case 'completed': return 'success'
      case 'active': return 'primary'
      case 'paused': return 'warning'
      default: return 'default'
    }
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'completed': return <CheckCircle />
      case 'active': return <PlayArrow />
      case 'paused': return <Pause />
      default: return <AccessTime />
    }
  }

  const nextComparison = getNextComparison()

  return (
    <Card>
      <CardContent>
        {/* Session Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Timeline color="primary" />
            <Typography variant="h6">
              {session?.criterion === 'value' ? 'Value Assessment' : 'Complexity Assessment'} Session
            </Typography>
            <Chip
              label={session?.status || 'Pending'}
              color={getStatusColor(session?.status)}
              icon={getStatusIcon(session?.status)}
              size="small"
            />
          </Box>

          {session?.status === 'active' && (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AccessTime fontSize="small" />
              <Typography variant="body2" color="text.secondary">
                {formatTime(sessionTime)}
              </Typography>
            </Box>
          )}
        </Box>

        {/* Overall Progress */}
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
            <Typography variant="subtitle2">
              Overall Progress
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {completedComparisons} of {totalComparisons} comparisons completed
            </Typography>
          </Box>
          <LinearProgress
            variant="determinate"
            value={overallProgress}
            sx={{ height: 8, borderRadius: 4 }}
          />
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            {Math.round(overallProgress)}% complete
          </Typography>
        </Box>

        {/* Quick Actions */}
        {session?.status === 'active' && currentAttendee && (
          <Box sx={{ mb: 3 }}>
            {nextComparison ? (
              <Button
                variant="contained"
                onClick={() => onNextComparison?.(nextComparison)}
                fullWidth
                size="large"
              >
                Continue Voting ({totalComparisons - Math.round(getAttendeeProgress(currentAttendee.id) / 100 * totalComparisons)} remaining)
              </Button>
            ) : (
              <Box sx={{ textAlign: 'center', p: 2, bgcolor: 'success.50', borderRadius: 1 }}>
                <CheckCircle color="success" sx={{ mb: 1 }} />
                <Typography variant="subtitle1" color="success.main">
                  You've completed all your votes!
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Waiting for other attendees to finish...
                </Typography>
              </Box>
            )}
          </Box>
        )}

        {/* Attendee Progress Details */}
        <Box>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
            <Typography variant="subtitle2" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Group fontSize="small" />
              Attendee Progress ({safeAttendees.length})
            </Typography>
            <IconButton size="small" onClick={() => setExpanded(!expanded)}>
              {expanded ? <ExpandLess /> : <ExpandMore />}
            </IconButton>
          </Box>

          <Collapse in={expanded}>
            <List dense>
              {safeAttendees && safeAttendees.length > 0 ? safeAttendees.map((attendee, index) => {
                const progress = getAttendeeProgress(attendee.id)
                const votedCount = Math.round((progress / 100) * totalComparisons)
                const isComplete = progress === 100
                const isCurrent = attendee.id === currentAttendee?.id

                return (
                  <React.Fragment key={attendee.id}>
                    <ListItem
                      sx={{
                        bgcolor: isCurrent ? 'primary.50' : 'transparent',
                        borderRadius: 1,
                        mb: 0.5,
                      }}
                    >
                      <ListItemAvatar>
                        <Avatar
                          sx={{
                            bgcolor: isComplete ? 'success.main' : 
                                   isCurrent ? 'primary.main' : 'grey.400',
                            width: 32,
                            height: 32,
                          }}
                        >
                          {isComplete ? (
                            <CheckCircle fontSize="small" />
                          ) : (
                            attendee.name?.[0] || '?'
                          )}
                        </Avatar>
                      </ListItemAvatar>

                      <ListItemText
                        primary={
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <Typography variant="subtitle2">
                              {attendee.name}
                              {isCurrent && ' (You)'}
                            </Typography>
                            {attendee.isFacilitator && (
                              <Chip label="Facilitator" size="small" color="primary" />
                            )}
                          </Box>
                        }
                        secondary={
                          <Box>
                            <Typography variant="body2" color="text.secondary">
                              {votedCount} of {totalComparisons} votes cast
                            </Typography>
                            <LinearProgress
                              variant="determinate"
                              value={progress}
                              sx={{ mt: 0.5, height: 4 }}
                              color={isComplete ? 'success' : 'primary'}
                            />
                          </Box>
                        }
                      />

                      <Box sx={{ textAlign: 'right' }}>
                        <Typography variant="h6" color={isComplete ? 'success.main' : 'text.secondary'}>
                          {Math.round(progress)}%
                        </Typography>
                        {isComplete && (
                          <CheckCircle color="success" fontSize="small" />
                        )}
                      </Box>
                    </ListItem>
                    {index < safeAttendees.length - 1 && <Divider />}
                  </React.Fragment>
                )
              }) : (
                <ListItem>
                  <ListItemText 
                    primary="No attendees available" 
                    secondary="Attendee data not loaded"
                  />
                </ListItem>
              )}
            </List>
          </Collapse>
        </Box>

        {/* Session Controls for Facilitators */}
        {currentAttendee?.isFacilitator && session?.status !== 'completed' && (
          <Box sx={{ mt: 3, pt: 2, borderTop: 1, borderColor: 'divider' }}>
            <Typography variant="subtitle2" gutterBottom>
              Facilitator Controls
            </Typography>
            <Box sx={{ display: 'flex', gap: 1 }}>
              {session?.status === 'active' ? (
                <Button
                  variant="outlined"
                  startIcon={<Pause />}
                  onClick={onPauseSession}
                  size="small"
                >
                  Pause Session
                </Button>
              ) : (
                <Button
                  variant="outlined"
                  startIcon={<PlayArrow />}
                  onClick={onResumeSession}
                  size="small"
                >
                  Resume Session
                </Button>
              )}
            </Box>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

export default SessionProgress