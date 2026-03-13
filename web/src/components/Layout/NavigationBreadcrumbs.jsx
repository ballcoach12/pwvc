import { NavigateNext } from '@mui/icons-material'
import { Breadcrumbs, Link, Typography } from '@mui/material'
import { Link as RouterLink, useLocation, useParams } from 'react-router-dom'

const NavigationBreadcrumbs = () => {
  const location = useLocation()
  const params = useParams()
  
  const pathnames = location.pathname.split('/').filter((x) => x)

  // Don't show breadcrumbs on home/projects page
  if (pathnames.length <= 1) {
    return null
  }

  const getBreadcrumbText = (pathname, index) => {
    switch (pathname) {
      case 'projects':
        return 'Projects'
      case 'new':
        return 'New Project'
      case 'edit':
        return 'Edit Project'
      case 'attendees':
        return 'Attendee Management'
      case 'features':
        return 'Feature Management'
      default:
        // If it's a UUID-like pattern, try to show project name
        if (pathname.match(/^[0-9a-f-]{36}$/i)) {
          return `Project ${pathname.slice(0, 8)}...`
        }
        return pathname.charAt(0).toUpperCase() + pathname.slice(1)
    }
  }

  const getBreadcrumbPath = (index) => {
    return '/' + pathnames.slice(0, index + 1).join('/')
  }

  return (
    <Breadcrumbs
      separator={<NavigateNext fontSize="small" />}
      aria-label="breadcrumb"
      sx={{ mb: 2 }}
    >
      <Link
        component={RouterLink}
        to="/projects"
        underline="hover"
        color="inherit"
      >
        Projects
      </Link>
      
      {pathnames.slice(1).map((pathname, index) => {
        const routeTo = getBreadcrumbPath(index + 1)
        const isLast = index === pathnames.length - 2
        const text = getBreadcrumbText(pathname, index + 1)

        return isLast ? (
          <Typography key={pathname} color="text.primary">
            {text}
          </Typography>
        ) : (
          <Link
            key={pathname}
            component={RouterLink}
            to={routeTo}
            underline="hover"
            color="inherit"
          >
            {text}
          </Link>
        )
      })}
    </Breadcrumbs>
  )
}

export default NavigationBreadcrumbs