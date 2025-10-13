import { Add, Cancel } from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    TextField,
    Typography,
} from '@mui/material'
import { useState } from 'react'

const FeatureForm = ({ onSubmit, onCancel, initialData = null, isLoading = false }) => {
  const [formData, setFormData] = useState({
    name: initialData?.name || '',
    description: initialData?.description || '',
    ...initialData,
  })
  const [errors, setErrors] = useState({})

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
      newErrors.name = 'Feature name is required'
    } else if (formData.name.trim().length < 3) {
      newErrors.name = 'Feature name must be at least 3 characters'
    }

    if (!formData.description || !formData.description.trim()) {
      newErrors.description = 'Description is required'
    } else if (formData.description.length > 500) {
      newErrors.description = 'Description must be less than 500 characters'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (event) => {
    event.preventDefault()
    
    if (!validateForm()) {
      return
    }

    const cleanedData = {
      title: formData.name.trim(),
      description: formData.description?.trim() || '',
    }

    onSubmit(cleanedData)
  }

  const handleReset = () => {
    setFormData({
      name: initialData?.name || '',
      description: initialData?.description || '',
      ...initialData,
    })
    setErrors({})
  }

  const isEditing = Boolean(initialData?.id)

  return (
    <Card>
      <form onSubmit={handleSubmit}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            {isEditing ? 'Edit Feature' : 'Add New Feature'}
          </Typography>

          {isEditing && (
            <Alert severity="info" sx={{ mb: 2 }}>
              Editing this feature will affect any existing pairwise comparisons.
            </Alert>
          )}

          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Feature Name"
              value={formData.name}
              onChange={handleChange('name')}
              error={Boolean(errors.name)}
              helperText={errors.name}
              required
              fullWidth
              placeholder="e.g., User Authentication, Dashboard Analytics, Mobile App"
              inputProps={{ maxLength: 100 }}
            />

            <TextField
              label="Description"
              value={formData.description}
              onChange={handleChange('description')}
              error={Boolean(errors.description)}
              helperText={errors.description || `Required. ${formData.description.length}/500 characters`}
              multiline
              rows={3}
              fullWidth
              placeholder="Provide additional context about this feature..."
              inputProps={{ maxLength: 500 }}
              required
            />
          </Box>
        </CardContent>

        <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              type="submit"
              variant="contained"
              startIcon={<Add />}
              disabled={isLoading}
            >
              {isLoading 
                ? (isEditing ? 'Updating...' : 'Adding...')
                : (isEditing ? 'Update Feature' : 'Add Feature')
              }
            </Button>
            
            {isEditing && (
              <Button
                type="button"
                variant="outlined"
                onClick={handleReset}
                disabled={isLoading}
              >
                Reset
              </Button>
            )}
          </Box>

          {onCancel && (
            <Button
              type="button"
              variant="text"
              startIcon={<Cancel />}
              onClick={onCancel}
              disabled={isLoading}
            >
              Cancel
            </Button>
          )}
        </CardActions>
      </form>
    </Card>
  )
}

export default FeatureForm