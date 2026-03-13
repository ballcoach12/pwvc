import {
    Assessment,
    CheckCircle,
    EmojiEvents,
    Group,
    RadioButtonUnchecked,
    Timer,
    TrendingUp
} from '@mui/icons-material'
import {
    Alert,
    Avatar,
    AvatarGroup,
    Box,
    Button,
    Chip,
    Divider,
    LinearProgress,
    List,
    ListItem,
    ListItemIcon,
    ListItemSecondaryAction,
    ListItemText,
    Paper,
    Typography
} from '@mui/material'
import React from 'react'

/**
 * ConsensusTracker displays the overall progress of the Fibonacci scoring session.
 * Shows feature-by-feature consensus status, attendee participation, and session
 * completion metrics with real-time updates.
 * 
 * Features:
 * - Overall session progress with completion percentage
 * - Individual feature consensus status
 * - Attendee participation tracking
 * - Real-time updates via WebSocket
 * - Session completion indicators
 * - Export and next phase navigation
 */
const ConsensusTracker = ({
  features = [],
  attendees = [],
  consensusScores = {}, // { featureId: consensusScore }
  individualScores = {}, // { featureId: { attendeeId: score } }
  criterionType = 'value',
  sessionStatus = 'active',
  onFeatureClick,
  onExportResults,
  onCompleteSession,
  onNextPhase
}) => {
  // Calculate session statistics
  const totalFeatures = features.length
  const completedFeatures = Object.keys(consensusScores).length
  const overallProgress = totalFeatures > 0 ? (completedFeatures / totalFeatures) * 100 : 0
  const isSessionComplete = completedFeatures === totalFeatures && totalFeatures > 0

  // Calculate attendee participation
  const getAttendeeParticipation = () => {
    const participation = {}
    attendees.forEach(attendee => {
      let scoredFeatures = 0
      features.forEach(feature => {
        if (individualScores[feature.id] && individualScores[feature.id][attendee.id] !== undefined) {
          scoredFeatures++
        }
      })
      participation[attendee.id] = {
        ...attendee,
        scoredFeatures,
        percentage: totalFeatures > 0 ? (scoredFeatures / totalFeatures) * 100 : 0
      }
    })
    return participation
  }

  const attendeeParticipation = getAttendeeParticipation()

  // Get feature status
  const getFeatureStatus = (feature) => {
    const hasConsensus = consensusScores[feature.id] !== undefined
    const scores = individualScores[feature.id] || {}
    const scoredCount = Object.keys(scores).length
    const totalAttendees = attendees.length

    if (hasConsensus) {
      return { status: 'completed', icon: CheckCircle, color: 'success.main', label: 'Consensus' }
    } else if (scoredCount > 0) {
      return { status: 'in-progress', icon: Timer, color: 'warning.main', label: `${scoredCount}/${totalAttendees} scored` }
    } else {
      return { status: 'pending', icon: RadioButtonUnchecked, color: 'grey.500', label: 'Not started' }
    }
  }

  // Format criterion type for display
  const criterionLabel = criterionType === 'value' ? 'Value' : 'Complexity'
  const criterionColor = criterionType === 'value' ? 'primary' : 'secondary'

  return (
    <Paper sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
        <Assessment color={criterionColor} />
        <Box sx={{ flexGrow: 1 }}>
          <Typography variant="h6">
            {criterionLabel} Scoring Progress
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Track consensus achievement across all features
          </Typography>
        </Box>
        {isSessionComplete && (
          <EmojiEvents color="warning" fontSize="large" />
        )}
      </Box>

      {/* Overall Progress */}
      <Box sx={{ mb: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
          <Typography variant="body2" fontWeight="medium">
            Session Progress
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {completedFeatures}/{totalFeatures} features complete ({Math.round(overallProgress)}%)
          </Typography>
        </Box>
        
        <LinearProgress
          variant="determinate"
          value={overallProgress}
          sx={{
            height: 8,
            borderRadius: 4,
            bgcolor: 'grey.200',
            '& .MuiLinearProgress-bar': {
              bgcolor: isSessionComplete ? 'success.main' : 'primary.main',
              borderRadius: 4
            }
          }}
        />

        {isSessionComplete && (
          <Alert severity="success" sx={{ mt: 2 }}>
            <Typography variant="body2">
              🎉 All features have reached consensus! Session is ready for completion.
            </Typography>
          </Alert>
        )}
      </Box>

      {/* Attendee Participation */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="subtitle2" gutterBottom>
          Attendee Participation
        </Typography>
        <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          {Object.values(attendeeParticipation).map(attendee => (
            <Chip
              key={attendee.id}
              avatar={
                <Avatar sx={{ width: 24, height: 24, fontSize: '0.7rem' }}>
                  {attendee.name?.charAt(0)?.toUpperCase() || '?'}
                </Avatar>
              }
              label={`${attendee.name} (${Math.round(attendee.percentage)}%)`}
              size="small"
              color={attendee.percentage === 100 ? 'success' : 
                     attendee.percentage > 50 ? 'primary' : 'default'}
              variant="outlined"
            />
          ))}
        </Box>
      </Box>

      <Divider sx={{ my: 2 }} />

      {/* Feature List */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="subtitle2" gutterBottom>
          Feature Status ({criterionLabel} Scoring)
        </Typography>
        
        <List dense sx={{ bgcolor: 'grey.50', borderRadius: 1 }}>
          {features.map((feature, index) => {
            const statusInfo = getFeatureStatus(feature)
            const StatusIcon = statusInfo.icon
            const consensusScore = consensusScores[feature.id]
            const featureScores = individualScores[feature.id] || {}
            const uniqueScores = [...new Set(Object.values(featureScores))]

            return (
              <React.Fragment key={feature.id}>
                <ListItem
                  button={!!onFeatureClick}
                  onClick={() => onFeatureClick && onFeatureClick(feature)}
                  sx={{
                    '&:hover': onFeatureClick ? { bgcolor: 'action.hover' } : undefined
                  }}
                >
                  <ListItemIcon>
                    <StatusIcon sx={{ color: statusInfo.color }} />
                  </ListItemIcon>
                  
                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Typography variant="body2" fontWeight="medium">
                          {feature.name}
                        </Typography>
                        {consensusScore && (
                          <Chip
                            label={`Score: ${consensusScore}`}
                            size="small"
                            color="success"
                            variant="outlined"
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                        <Typography variant="caption" color="text.secondary">
                          {statusInfo.label}
                        </Typography>
                        {uniqueScores.length > 0 && !consensusScore && (
                          <Typography variant="caption" color="text.secondary">
                            • Scores: {uniqueScores.sort((a, b) => a - b).join(', ')}
                          </Typography>
                        )}
                      </Box>
                    }
                  />

                  <ListItemSecondaryAction>
                    <AvatarGroup max={4} sx={{ '& .MuiAvatar-root': { width: 24, height: 24, fontSize: '0.7rem' } }}>
                      {attendees.map(attendee => {
                        const hasScore = featureScores[attendee.id] !== undefined
                        return (
                          <Avatar
                            key={attendee.id}
                            sx={{
                              bgcolor: hasScore ? 'primary.main' : 'grey.300',
                              color: hasScore ? 'primary.contrastText' : 'text.secondary'
                            }}
                          >
                            {attendee.name?.charAt(0)?.toUpperCase() || '?'}
                          </Avatar>
                        )
                      })}
                    </AvatarGroup>
                  </ListItemSecondaryAction>
                </ListItem>
                
                {index < features.length - 1 && <Divider variant="inset" />}
              </React.Fragment>
            )
          })}
        </List>
      </Box>

      {/* Action Buttons */}
      <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
        {onExportResults && (
          <Button
            variant="outlined"
            startIcon={<TrendingUp />}
            onClick={onExportResults}
            disabled={completedFeatures === 0}
          >
            Export Progress
          </Button>
        )}

        {onCompleteSession && isSessionComplete && (
          <Button
            variant="contained"
            color="success"
            startIcon={<CheckCircle />}
            onClick={onCompleteSession}
          >
            Complete {criterionLabel} Session
          </Button>
        )}

        {onNextPhase && isSessionComplete && (
          <Button
            variant="contained"
            startIcon={<Group />}
            onClick={onNextPhase}
          >
            {criterionType === 'value' ? 'Start Complexity Scoring' : 'Calculate Final Results'}
          </Button>
        )}
      </Box>

      {/* Warnings */}
      {!isSessionComplete && completedFeatures > 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          <Typography variant="body2">
            💡 Session will be ready for completion when all {totalFeatures} features have consensus scores.
            Features remaining: {totalFeatures - completedFeatures}
          </Typography>
        </Alert>
      )}

      {completedFeatures === 0 && totalFeatures > 0 && (
        <Alert severity="warning" sx={{ mt: 2 }}>
          <Typography variant="body2">
            ⚠️ No features have been scored yet. Attendees need to start scoring features to build consensus.
          </Typography>
        </Alert>
      )}
    </Paper>
  )
}

export default ConsensusTracker