import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Typography,
} from '@mui/material'
import { useState } from 'react'

const AttendeeLoginDialog = ({ open, onClose, onLogin, attendees, loading }) => {
  const [selectedAttendee, setSelectedAttendee] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!selectedAttendee) return
    await onLogin(selectedAttendee)
  }

  const handleClose = () => {
    setSelectedAttendee('')
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        Select Your Identity
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Choose which attendee you are to participate in the pairwise comparison
        </Typography>
      </DialogTitle>

      <form onSubmit={handleSubmit}>
        <DialogContent>
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Select Your Identity</InputLabel>
            <Select
              value={selectedAttendee}
              onChange={(e) => setSelectedAttendee(e.target.value)}
              label="Select Your Identity"
              disabled={loading}
            >
              {attendees.map((attendee) => (
                <MenuItem key={attendee.id} value={attendee.id}>
                  {attendee.name} - {attendee.role}
                  {attendee.is_facilitator && ' (Facilitator)'}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button type="submit" variant="contained" disabled={!selectedAttendee}>
            Confirm
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default AttendeeLoginDialog