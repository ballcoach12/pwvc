import { Add, Refresh } from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Grid,
    Typography,
} from '@mui/material'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSnackbar } from '../../components/NotificationProvider.jsx'
import ProjectCard from '../../components/ProjectCard/ProjectCard.jsx'
import { projectService } from '../../services/projectService.js'

const ProjectList = () => {
  const navigate = useNavigate()
  const { enqueueSnackbar } = useSnackbar()
  
  const [projects, setProjects] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [deleteDialog, setDeleteDialog] = useState({ open: false, project: null })

  useEffect(() => {
    loadProjects()
  }, [])

  const loadProjects = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await projectService.getProjects()
      setProjects(data || [])
    } catch (err) {
      setError(err.message || 'Failed to load projects')
      enqueueSnackbar('Failed to load projects', { variant: 'error' })
    } finally {
      setLoading(false)
    }
  }

  const handleCreateProject = () => {
    navigate('/projects/new')
  }

  const handleEditProject = (project) => {
    navigate(`/projects/${project.id}/edit`)
  }

  const handleDeleteProject = (project) => {
    setDeleteDialog({ open: true, project })
  }

  const confirmDelete = async () => {
    const { project } = deleteDialog
    if (!project) return

    try {
      await projectService.deleteProject(project.id)
      setProjects(prev => prev.filter(p => p.id !== project.id))
      enqueueSnackbar('Project deleted successfully', { variant: 'success' })
    } catch (err) {
      enqueueSnackbar(err.message || 'Failed to delete project', { variant: 'error' })
    } finally {
      setDeleteDialog({ open: false, project: null })
    }
  }

  const cancelDelete = () => {
    setDeleteDialog({ open: false, project: null })
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
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            Projects
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Manage your P-WVC feature prioritization projects
          </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={loadProjects}
            disabled={loading}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleCreateProject}
          >
            New Project
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {projects.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" gutterBottom>
            No projects yet
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            Create your first P-WVC project to get started with feature prioritization
          </Typography>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleCreateProject}
            size="large"
          >
            Create First Project
          </Button>
        </Box>
      ) : (
        <Grid container spacing={3}>
          {projects.map((project) => (
            <Grid item xs={12} sm={6} md={4} key={project.id}>
              <ProjectCard
                project={project}
                onEdit={handleEditProject}
                onDelete={handleDeleteProject}
              />
            </Grid>
          ))}
        </Grid>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={cancelDelete}
        aria-labelledby="delete-dialog-title"
        aria-describedby="delete-dialog-description"
      >
        <DialogTitle id="delete-dialog-title">
          Delete Project
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="delete-dialog-description">
            Are you sure you want to delete "{deleteDialog.project?.name}"? 
            This will permanently remove the project and all associated data including attendees, features, and comparisons.
            <br /><br />
            <strong>This action cannot be undone.</strong>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={cancelDelete}>
            Cancel
          </Button>
          <Button onClick={confirmDelete} color="error" variant="contained">
            Delete Project
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}

export default ProjectList