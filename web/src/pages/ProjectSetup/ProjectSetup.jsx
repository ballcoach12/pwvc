import { ArrowForward, Cancel, Save } from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CircularProgress,
    TextField,
    Typography,
} from '@mui/material'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useSnackbar } from '../../components/NotificationProvider.jsx'
import { projectService } from '../../services/projectService.js'

const ProjectSetup = () => {
  const navigate = useNavigate()
  const { id } = useParams()
  const { enqueueSnackbar } = useSnackbar()
  
  const isEditing = Boolean(id)
  
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  })
  const [errors, setErrors] = useState({})
  const [loading, setLoading] = useState(false)
  const [initialLoading, setInitialLoading] = useState(isEditing)

  useEffect(() => {
    if (isEditing) {
      loadProject()
    }
  }, [id, isEditing])

  const loadProject = async () => {
    try {
      setInitialLoading(true)
      const project = await projectService.getProject(id)
      setFormData({
        name: project.name || '',
        description: project.description || '',
      })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to load project', { variant: 'error' })
      navigate('/projects')
    } finally {
      setInitialLoading(false)
    }
  }

  const handleChange = (field) => (event) => {
    const value = event.target.value
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
      newErrors.name = 'Project name is required'
    } else if (formData.name.trim().length < 3) {
      newErrors.name = 'Project name must be at least 3 characters'
    } else if (formData.name.trim().length > 100) {
      newErrors.name = 'Project name must be less than 100 characters'
    }

    if (formData.description && formData.description.length > 500) {
      newErrors.description = 'Description must be less than 500 characters'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    
    if (!validateForm()) {
      return
    }

    const cleanedData = {
      name: formData.name.trim(),
      description: formData.description?.trim() || '',
    }

    try {
      setLoading(true)
      let project

      if (isEditing) {
        project = await projectService.updateProject(id, cleanedData)
        enqueueSnackbar('Project updated successfully', { variant: 'success' })
      } else {
        project = await projectService.createProject(cleanedData)
        enqueueSnackbar('Project created successfully', { variant: 'success' })
      }

      // Validate project response
      if (!project || !project.id) {
        enqueueSnackbar('Project created but received invalid response. Please refresh the page.', { variant: 'warning' })
        return
      }

      // Navigate to attendee management
      navigate(`/projects/${project.id}/attendees`)
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to save project', { variant: 'error' })
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    navigate('/projects')
  }

  const handleContinue = () => {
    navigate(`/projects/${id}/attendees`)
  }

  if (initialLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box sx={{ maxWidth: 600, mx: 'auto' }}>
      <Typography variant="h4" gutterBottom>
        {isEditing ? 'Edit Project' : 'Create New Project'}
      </Typography>
      
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        {isEditing 
          ? 'Update your project details below.'
          : 'Start by setting up your project details. You\'ll add attendees and features in the next steps.'
        }
      </Typography>

      <Card>
        <form onSubmit={handleSubmit}>
          <CardContent>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
              <TextField
                label="Project Name"
                value={formData.name}
                onChange={handleChange('name')}
                error={Boolean(errors.name)}
                helperText={errors.name || 'Choose a descriptive name for your project'}
                required
                fullWidth
                placeholder="e.g., Mobile App Features Q4 2024"
                inputProps={{ maxLength: 100 }}
              />

              <TextField
                label="Description"
                value={formData.description}
                onChange={handleChange('description')}
                error={Boolean(errors.description)}
                helperText={
                  errors.description || 
                  `Optional. Provide context about this prioritization session. ${formData.description.length}/500 characters`
                }
                multiline
                rows={4}
                fullWidth
                placeholder="Describe the scope, objectives, or any relevant context for this prioritization session..."
                inputProps={{ maxLength: 500 }}
              />
            </Box>

            {!isEditing && (
              <Alert severity="info" sx={{ mt: 3 }}>
                <Typography variant="body2">
                  <strong>Next steps:</strong> After creating your project, you'll:
                </Typography>
                <Typography variant="body2" component="div" sx={{ mt: 1 }}>
                  1. Add team members and assign facilitators<br />
                  2. Input or import features to prioritize<br />
                  3. Conduct pairwise comparisons<br />
                  4. Review results and export findings
                </Typography>
              </Alert>
            )}
          </CardContent>

          <CardActions sx={{ justifyContent: 'space-between', p: 2 }}>
            <Button
              type="button"
              variant="outlined"
              startIcon={<Cancel />}
              onClick={handleCancel}
              disabled={loading}
            >
              Cancel
            </Button>

            <Box sx={{ display: 'flex', gap: 1 }}>
              <Button
                type="submit"
                variant="contained"
                startIcon={<Save />}
                disabled={loading}
              >
                {loading 
                  ? (isEditing ? 'Updating...' : 'Creating...')
                  : (isEditing ? 'Update Project' : 'Create Project')
                }
              </Button>
              
              {isEditing && (
                <Button
                  type="button"
                  variant="contained"
                  color="secondary"
                  endIcon={<ArrowForward />}
                  onClick={handleContinue}
                  disabled={loading}
                >
                  Continue Setup
                </Button>
              )}
            </Box>
          </CardActions>
        </form>
      </Card>
    </Box>
  )
}

export default ProjectSetup