import { Add, ArrowForward, PersonAdd } from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    Checkbox,
    CircularProgress,
    FormControlLabel,
    TextField,
    Typography
} from '@mui/material'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import AttendeeList from '../../components/AttendeeList/AttendeeList.jsx'
import { useSnackbar } from '../../components/NotificationProvider.jsx'
import { projectService } from '../../services/projectService.js'

const AttendeeManagement = () => {
  const navigate = useNavigate()
  const { id } = useParams()
  const { enqueueSnackbar } = useSnackbar()
  
  const [project, setProject] = useState(null)
  const [attendees, setAttendees] = useState([])
  const [loading, setLoading] = useState(true)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    role: '',
    isFacilitator: false,
  })
  const [errors, setErrors] = useState({})
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadProjectAndAttendees()
  }, [id])

  const loadProjectAndAttendees = async () => {
    try {
      setLoading(true)
      const [projectData, attendeesData] = await Promise.all([
        projectService.getProject(id),
        projectService.getProjectAttendees(id)
      ])
      setProject(projectData)
      setAttendees(attendeesData.attendees || [])
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to load project data', { variant: 'error' })
      navigate('/projects')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field) => (event) => {
    const value = event.target.type === 'checkbox' ? event.target.checked : event.target.value
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }))
    
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({
        ...prev,
        [field]: null,
      }))
    }
  }

  const validateForm = () => {
    const newErrors = {}

    if (!formData.name.trim()) {
      newErrors.name = 'Name is required'
    } else if (formData.name.trim().length < 2) {
      newErrors.name = 'Name must be at least 2 characters'
    }

    if (!formData.email.trim()) {
      newErrors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email.trim())) {
      newErrors.email = 'Please enter a valid email address'
    } else {
      // Check for duplicate email
      const existingAttendee = attendees.find(
        a => a.email.toLowerCase() === formData.email.trim().toLowerCase()
      )
      if (existingAttendee) {
        newErrors.email = 'This email is already added to the project'
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleAddAttendee = async (event) => {
    event.preventDefault()
    
    if (!validateForm()) {
      return
    }

    const cleanedData = {
      name: formData.name.trim(),
      email: formData.email.trim().toLowerCase(),
      role: formData.role.trim() || null,
      isFacilitator: formData.isFacilitator,
    }

    try {
      setSubmitting(true)
      const newAttendee = await projectService.addAttendee(id, cleanedData)
      setAttendees(prev => [...prev, newAttendee])
      setFormData({
        name: '',
        email: '',
        role: '',
        isFacilitator: false,
      })
      enqueueSnackbar('Attendee added successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to add attendee', { variant: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleRemoveAttendee = async (attendee) => {
    try {
      await projectService.removeAttendee(id, attendee.id)
      setAttendees(prev => prev.filter(a => a.id !== attendee.id))
      enqueueSnackbar('Attendee removed successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to remove attendee', { variant: 'error' })
    }
  }

  const handleSetFacilitator = async (attendee) => {
    try {
      await projectService.setFacilitator(id, attendee.id)
      setAttendees(prev => prev.map(a => ({
        ...a,
        isFacilitator: a.id === attendee.id
      })))
      enqueueSnackbar('Facilitator set successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to set facilitator', { variant: 'error' })
    }
  }

  const handleContinue = () => {
    if (attendees.length === 0) {
      enqueueSnackbar('Please add at least one attendee before continuing', { variant: 'warning' })
      return
    }
    navigate(`/projects/${id}/features`)
  }

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Attendee Management
      </Typography>
      
      <Typography variant="body1" color="text.secondary" sx={{ mb: 1 }}>
        Project: <strong>{project?.name}</strong>
      </Typography>
      
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Add team members who will participate in the PairWise prioritization process. 
        Assign at least one facilitator to guide the session.
      </Typography>

      <Box sx={{ display: 'flex', gap: 3, flexDirection: { xs: 'column', md: 'row' } }}>
        {/* Add Attendee Form */}
        <Box sx={{ flex: 1 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <PersonAdd />
                Add Attendee
              </Typography>
              
              <form onSubmit={handleAddAttendee}>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <TextField
                    label="Full Name"
                    value={formData.name}
                    onChange={handleChange('name')}
                    error={Boolean(errors.name)}
                    helperText={errors.name}
                    required
                    fullWidth
                    placeholder="e.g., John Smith"
                  />

                  <TextField
                    label="Email Address"
                    type="email"
                    value={formData.email}
                    onChange={handleChange('email')}
                    error={Boolean(errors.email)}
                    helperText={errors.email}
                    required
                    fullWidth
                    placeholder="e.g., john.smith@company.com"
                  />

                  <TextField
                    label="Role/Title"
                    value={formData.role}
                    onChange={handleChange('role')}
                    fullWidth
                    placeholder="e.g., Product Manager, Developer"
                    helperText="Optional. Job title or role in the project"
                  />

                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={formData.isFacilitator}
                        onChange={handleChange('isFacilitator')}
                        color="primary"
                      />
                    }
                    label="Set as Facilitator"
                  />
                </Box>
              </form>
            </CardContent>
            
            <CardActions>
              <Button
                type="submit"
                variant="contained"
                startIcon={<Add />}
                onClick={handleAddAttendee}
                disabled={submitting}
                fullWidth
              >
                {submitting ? 'Adding...' : 'Add Attendee'}
              </Button>
            </CardActions>
          </Card>

          <Alert severity="info" sx={{ mt: 2 }}>
            <Typography variant="body2">
              <strong>Facilitator Role:</strong> Facilitators guide the PairWise process, 
              manage discussions, and ensure consensus is reached. At least one facilitator is recommended.
            </Typography>
          </Alert>
        </Box>

        {/* Attendee List */}
        <Box sx={{ flex: 2 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Project Attendees ({attendees.length})
              </Typography>
              
              <AttendeeList
                attendees={attendees}
                onRemove={handleRemoveAttendee}
                onSetFacilitator={handleSetFacilitator}
              />
            </CardContent>
          </Card>
        </Box>
      </Box>

      {/* Continue Button */}
      <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end' }}>
        <Button
          variant="contained"
          size="large"
          endIcon={<ArrowForward />}
          onClick={handleContinue}
          disabled={attendees.length === 0}
        >
          Continue to Features
        </Button>
      </Box>

      {attendees.length === 0 && (
        <Alert severity="warning" sx={{ mt: 2 }}>
          Please add at least one attendee before continuing to the next step.
        </Alert>
      )}
    </Box>
  )
}

export default AttendeeManagement