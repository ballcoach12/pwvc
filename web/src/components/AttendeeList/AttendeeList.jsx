import {
    Delete,
    Person,
    Star,
} from '@mui/icons-material'
import {
    Avatar,
    Box,
    Chip,
    Divider,
    IconButton,
    List,
    ListItem,
    ListItemSecondaryAction,
    ListItemText,
    Typography,
} from '@mui/material'

const AttendeeList = ({ attendees = [], onRemove, onSetFacilitator, readOnly = false }) => {
  if (attendees.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="body1" color="text.secondary">
          No attendees added yet
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Add team members to begin the P-WVC process
        </Typography>
      </Box>
    )
  }

  const facilitators = attendees.filter(attendee => attendee.isFacilitator)
  const participants = attendees.filter(attendee => !attendee.isFacilitator)

  const renderAttendee = (attendee, index) => {
    const initials = attendee.name
      ?.split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2) || '??'

    return (
      <ListItem key={attendee.id} divider={index < attendees.length - 1}>
        <Avatar
          sx={{
            mr: 2,
            bgcolor: attendee.isFacilitator ? 'primary.main' : 'grey.500',
            width: 40,
            height: 40,
          }}
        >
          {attendee.isFacilitator ? <Star /> : initials}
        </Avatar>
        
        <ListItemText
          primary={
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Typography variant="subtitle1">
                {attendee.name}
              </Typography>
              {attendee.isFacilitator && (
                <Chip
                  label="Facilitator"
                  size="small"
                  color="primary"
                  icon={<Star />}
                />
              )}
            </Box>
          }
          secondary={
            <Box>
              <Typography variant="body2" color="text.secondary">
                {attendee.email}
              </Typography>
              {attendee.role && (
                <Typography variant="caption" color="text.secondary">
                  {attendee.role}
                </Typography>
              )}
            </Box>
          }
        />

        {!readOnly && (
          <ListItemSecondaryAction>
            <Box sx={{ display: 'flex', gap: 1 }}>
              {!attendee.isFacilitator && onSetFacilitator && (
                <IconButton
                  edge="end"
                  aria-label="set as facilitator"
                  onClick={() => onSetFacilitator(attendee)}
                  title="Set as facilitator"
                >
                  <Star />
                </IconButton>
              )}
              {onRemove && (
                <IconButton
                  edge="end"
                  aria-label="remove"
                  onClick={() => onRemove(attendee)}
                  title="Remove attendee"
                  color="error"
                >
                  <Delete />
                </IconButton>
              )}
            </Box>
          </ListItemSecondaryAction>
        )}
      </ListItem>
    )
  }

  return (
    <Box>
      {facilitators.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Star color="primary" />
            Facilitators
          </Typography>
          <List disablePadding>
            {facilitators.map((attendee, index) => renderAttendee(attendee, index))}
          </List>
        </Box>
      )}

      {participants.length > 0 && (
        <Box>
          {facilitators.length > 0 && <Divider sx={{ my: 2 }} />}
          <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Person />
            Participants
          </Typography>
          <List disablePadding>
            {participants.map((attendee, index) => renderAttendee(attendee, index))}
          </List>
        </Box>
      )}

      <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
        <Typography variant="body2" color="text.secondary">
          <strong>Total:</strong> {attendees.length} attendee{attendees.length !== 1 ? 's' : ''}
          {facilitators.length > 0 && (
            <span> • {facilitators.length} facilitator{facilitators.length !== 1 ? 's' : ''}</span>
          )}
        </Typography>
      </Box>
    </Box>
  )
}

export default AttendeeList