import {
    Add,
    ArrowForward,
    Delete,
    Edit,
    FileUpload,
    GetApp,
    ListAlt,
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Chip,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Divider,
    IconButton,
    List,
    ListItem,
    ListItemSecondaryAction,
    ListItemText,
    Tab,
    Tabs,
    Typography,
} from '@mui/material'
import React, { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import FeatureForm from '../../components/FeatureForm/FeatureForm.jsx'
import FileUploadComponent from '../../components/FileUpload/FileUpload.jsx'
import { useSnackbar } from '../../components/NotificationProvider.jsx'
import { featureService } from '../../services/featureService.js'
import { projectService } from '../../services/projectService.js'

const FeatureManagement = () => {
  const navigate = useNavigate()
  const { id } = useParams()
  const { enqueueSnackbar } = useSnackbar()
  
  const [project, setProject] = useState(null)
  const [features, setFeatures] = useState([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState(0)
  const [editingFeature, setEditingFeature] = useState(null)
  const [deleteDialog, setDeleteDialog] = useState({ open: false, feature: null })
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadProjectAndFeatures()
  }, [id])

  const loadProjectAndFeatures = async () => {
    try {
      setLoading(true)
      const [projectData, featuresData] = await Promise.all([
        projectService.getProject(id),
        featureService.getFeatures(id)
      ])
      setProject(projectData)
      setFeatures(featuresData.features || [])
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to load project data', { variant: 'error' })
      navigate('/projects')
    } finally {
      setLoading(false)
    }
  }

  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue)
    setEditingFeature(null)
  }

  const handleAddFeature = async (featureData) => {
    try {
      setSubmitting(true)
      const newFeature = await featureService.createFeature(id, featureData)
      setFeatures(prev => [...prev, newFeature])
      enqueueSnackbar('Feature added successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to add feature', { variant: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleEditFeature = (feature) => {
    setEditingFeature(feature)
    setActiveTab(0) // Switch to manual add tab
  }

  const handleUpdateFeature = async (featureData) => {
    try {
      setSubmitting(true)
      const updatedFeature = await featureService.updateFeature(id, editingFeature.id, featureData)
      setFeatures(prev => prev.map(f => f.id === editingFeature.id ? updatedFeature : f))
      setEditingFeature(null)
      enqueueSnackbar('Feature updated successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to update feature', { variant: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteFeature = (feature) => {
    setDeleteDialog({ open: true, feature })
  }

  const confirmDelete = async () => {
    const { feature } = deleteDialog
    if (!feature) return

    try {
      await featureService.deleteFeature(id, feature.id)
      setFeatures(prev => prev.filter(f => f.id !== feature.id))
      enqueueSnackbar('Feature deleted successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to delete feature', { variant: 'error' })
    } finally {
      setDeleteDialog({ open: false, feature: null })
    }
  }

  const handleFileUpload = async ({ data }) => {
    try {
      setSubmitting(true)
      const result = await featureService.importFeatures(id, data)
      setFeatures(prev => [...prev, ...result.features])
      enqueueSnackbar(
        `Successfully imported ${result.features.length} feature(s)`, 
        { variant: 'success' }
      )
      setActiveTab(2) // Switch to list tab to show imported features
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to import features', { variant: 'error' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleExportFeatures = async () => {
    try {
      const blob = await featureService.exportFeatures(id)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.style.display = 'none'
      a.href = url
      a.download = `${project?.name || 'project'}-features.csv`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      enqueueSnackbar('Features exported successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to export features', { variant: 'error' })
    }
  }

  const handleContinue = () => {
    if (features.length < 2) {
      enqueueSnackbar('Please add at least 2 features before continuing', { variant: 'warning' })
      return
    }
    navigate(`/projects/${id}/comparison`)
  }

  const cancelEdit = () => {
    setEditingFeature(null)
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
        Feature Management
      </Typography>
      
      <Typography variant="body1" color="text.secondary" sx={{ mb: 1 }}>
        Project: <strong>{project?.name}</strong>
      </Typography>
      
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Add features that need to be prioritized. You can add them manually or import from a CSV file.
        At least 2 features are required to proceed with pairwise comparison.
      </Typography>

      <Card>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={activeTab} onChange={handleTabChange}>
            <Tab label="Add Feature" icon={<Add />} />
            <Tab label="Import CSV" icon={<FileUpload />} />
            <Tab label={`Feature List (${features.length})`} icon={<ListAlt />} />
          </Tabs>
        </Box>

        <CardContent>
          {/* Tab 0: Add Feature */}
          {activeTab === 0 && (
            <FeatureForm
              onSubmit={editingFeature ? handleUpdateFeature : handleAddFeature}
              onCancel={editingFeature ? cancelEdit : null}
              initialData={editingFeature}
              isLoading={submitting}
            />
          )}

          {/* Tab 1: Import CSV */}
          {activeTab === 1 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Import Features from CSV
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                Upload a CSV file with feature data. Required column: <strong>name</strong>. 
                Optional column: <strong>description</strong>.
              </Typography>
              
              <FileUploadComponent
                onFileUpload={handleFileUpload}
                disabled={submitting}
              />

              <Alert severity="info" sx={{ mt: 2 }}>
                <Typography variant="body2">
                  <strong>CSV Format Example:</strong>
                </Typography>
                <Typography variant="body2" component="div" sx={{ mt: 1, fontFamily: 'monospace' }}>
                  name,description<br />
                  User Authentication,Login and registration system<br />
                  Dashboard Analytics,Real-time analytics dashboard<br />
                  Mobile App,iOS and Android mobile application
                </Typography>
              </Alert>
            </Box>
          )}

          {/* Tab 2: Feature List */}
          {activeTab === 2 && (
            <Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6">
                  Features ({features.length})
                </Typography>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  {features.length > 0 && (
                    <Button
                      variant="outlined"
                      startIcon={<GetApp />}
                      onClick={handleExportFeatures}
                      size="small"
                    >
                      Export CSV
                    </Button>
                  )}
                </Box>
              </Box>

              {features.length === 0 ? (
                <Box sx={{ textAlign: 'center', py: 4 }}>
                  <Typography variant="body1" color="text.secondary">
                    No features added yet
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Add features manually or import from CSV to get started
                  </Typography>
                </Box>
              ) : (
                <List>
                  {features.map((feature, index) => (
                    <React.Fragment key={feature.id}>
                      <ListItem alignItems="flex-start">
                        <ListItemText
                          primary={
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <Typography variant="subtitle1">
                                {feature.name}
                              </Typography>
                              <Chip
                                label={`#${index + 1}`}
                                size="small"
                                variant="outlined"
                              />
                            </Box>
                          }
                          secondary={feature.description || 'No description provided'}
                        />
                        <ListItemSecondaryAction>
                          <IconButton
                            edge="end"
                            aria-label="edit"
                            onClick={() => handleEditFeature(feature)}
                            sx={{ mr: 1 }}
                          >
                            <Edit />
                          </IconButton>
                          <IconButton
                            edge="end"
                            aria-label="delete"
                            onClick={() => handleDeleteFeature(feature)}
                            color="error"
                          >
                            <Delete />
                          </IconButton>
                        </ListItemSecondaryAction>
                      </ListItem>
                      {index < features.length - 1 && <Divider />}
                    </React.Fragment>
                  ))}
                </List>
              )}
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Continue Button */}
      <Box sx={{ mt: 4, display: 'flex', justifyContent: 'flex-end' }}>
        <Button
          variant="contained"
          size="large"
          endIcon={<ArrowForward />}
          onClick={handleContinue}
          disabled={features.length < 2}
        >
          Continue to Pairwise Comparison
        </Button>
      </Box>

      {features.length < 2 && (
        <Alert severity="warning" sx={{ mt: 2 }}>
          Please add at least 2 features before continuing to the pairwise comparison phase.
        </Alert>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false, feature: null })}
      >
        <DialogTitle>Delete Feature</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete "{deleteDialog.feature?.name}"? 
            This will also remove any existing pairwise comparisons involving this feature.
            <br /><br />
            <strong>This action cannot be undone.</strong>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialog({ open: false, feature: null })}>
            Cancel
          </Button>
          <Button onClick={confirmDelete} color="error" variant="contained">
            Delete Feature
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}

export default FeatureManagement