import {
    Delete,
    Edit,
    ListAlt,
    MoreVert,
    People,
} from '@mui/icons-material'
import {
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    Chip,
    IconButton,
    Menu,
    MenuItem,
    Typography,
} from '@mui/material'
import React from 'react'
import { useNavigate } from 'react-router-dom'

const ProjectCard = ({ project, onEdit, onDelete }) => {
  const navigate = useNavigate()
  const [anchorEl, setAnchorEl] = React.useState(null)
  const open = Boolean(anchorEl)

  const handleMenuClick = (event) => {
    event.stopPropagation()
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleEdit = (event) => {
    event.stopPropagation()
    handleMenuClose()
    onEdit?.(project)
  }

  const handleDelete = (event) => {
    event.stopPropagation()
    handleMenuClose()
    onDelete?.(project)
  }

  const handleCardClick = () => {
    navigate(`/projects/${project.id}/attendees`)
  }

  const formatDate = (dateString) => {
    if (!dateString) return 'No date'
    return new Date(dateString).toLocaleDateString()
  }

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'completed':
        return 'success'
      case 'in-progress':
        return 'primary'
      case 'draft':
        return 'default'
      default:
        return 'default'
    }
  }

  return (
    <Card
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        cursor: 'pointer',
        transition: 'all 0.2s ease-in-out',
        '&:hover': {
          transform: 'translateY(-2px)',
        },
      }}
      onClick={handleCardClick}
    >
      <CardContent sx={{ flexGrow: 1 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
          <Typography variant="h6" component="h2" gutterBottom>
            {project.name}
          </Typography>
          <IconButton
            size="small"
            onClick={handleMenuClick}
            sx={{ mt: -1, mr: -1 }}
          >
            <MoreVert />
          </IconButton>
        </Box>

        <Typography variant="body2" color="text.secondary" sx={{ mb: 2, minHeight: 40 }}>
          {project.description || 'No description provided'}
        </Typography>

        <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
          <Chip
            label={project.status || 'Draft'}
            color={getStatusColor(project.status)}
            size="small"
          />
          {project.attendeeCount > 0 && (
            <Chip
              icon={<People />}
              label={`${project.attendeeCount} attendees`}
              size="small"
              variant="outlined"
            />
          )}
          {project.featureCount > 0 && (
            <Chip
              icon={<ListAlt />}
              label={`${project.featureCount} features`}
              size="small"
              variant="outlined"
            />
          )}
        </Box>

        <Typography variant="caption" color="text.secondary">
          Created: {formatDate(project.createdAt)}
        </Typography>
      </CardContent>

      <CardActions sx={{ justifyContent: 'space-between', pt: 0 }}>
        <Button
          size="small"
          onClick={(e) => {
            e.stopPropagation()
            navigate(`/projects/${project.id}/attendees`)
          }}
        >
          Manage
        </Button>
        <Typography variant="caption" color="text.secondary">
          Updated: {formatDate(project.updatedAt)}
        </Typography>
      </CardActions>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleMenuClose}
        onClick={(e) => e.stopPropagation()}
      >
        <MenuItem onClick={handleEdit}>
          <Edit sx={{ mr: 1 }} fontSize="small" />
          Edit
        </MenuItem>
        <MenuItem onClick={handleDelete} sx={{ color: 'error.main' }}>
          <Delete sx={{ mr: 1 }} fontSize="small" />
          Delete
        </MenuItem>
      </Menu>
    </Card>
  )
}

export default ProjectCard