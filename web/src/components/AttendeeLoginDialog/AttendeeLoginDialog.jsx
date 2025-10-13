import {
    Alert,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    TextField,
    Typography
} from '@mui/material'
import { useState } from 'react'

const AttendeeLoginDialog = ({ open, onClose, onLogin, attendees, loading }) => {
  const [selectedAttendee, setSelectedAttendee] = useState('')
  const [pin, setPin] = useState('')
  const [error, setError] = useState('')
  const [isLogging] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    if (!selectedAttendee || !pin) {
      setError('Please select an attendee and enter your PIN')
      return
    }

    try {
      await onLogin(selectedAttendee, pin)
    } catch (err) {
      setError(err.message || 'Login failed')
    }
  }

  const handleClose = () => {
    setSelectedAttendee('')
    setPin('')
    setError('')
    onClose()
  }

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        Attendee Login
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Select your identity and enter your PIN to participate in the pairwise comparison
        </Typography>
      </DialogTitle>
      
      <form onSubmit={handleSubmit}>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

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

          <TextField
            fullWidth
            type="password"
            label="PIN"
            value={pin}
            onChange={(e) => setPin(e.target.value)}
            placeholder="Enter your 4-digit PIN"
            disabled={loading}
            inputProps={{ maxLength: 10 }}
          />

          <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
            Contact your facilitator if you don't know your PIN
          </Typography>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={isLogging}>
            Cancel
          </Button>
          <Button 
            type="submit" 
            variant="contained" 
            disabled={isLogging || !selectedAttendee || !pin}
          >
            {isLogging ? 'Logging in...' : 'Login'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default AttendeeLoginDialog