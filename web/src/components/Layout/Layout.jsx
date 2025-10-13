import { AppBar, Box, Container, Toolbar, Typography } from '@mui/material'
import { useLocation, useNavigate } from 'react-router-dom'
import NavigationBreadcrumbs from './NavigationBreadcrumbs.jsx'

const Layout = ({ children }) => {
  const location = useLocation()
  const navigate = useNavigate()

  const handleTitleClick = () => {
    navigate('/projects')
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar position="static" elevation={1}>
        <Toolbar>
          <Typography 
            variant="h6" 
            component="h1" 
            sx={{ 
              flexGrow: 1, 
              cursor: 'pointer',
              fontWeight: 600,
              '&:hover': { opacity: 0.8 }
            }}
            onClick={handleTitleClick}
          >
            PairWise
          </Typography>
          <Typography variant="body2" sx={{ opacity: 0.8 }}>
            Feature Prioritization Through Group Consensus
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ py: 3, flex: 1 }}>
        <NavigationBreadcrumbs />
        <Box sx={{ mt: 2 }}>
          {children}
        </Box>
      </Container>

      <Box
        component="footer"
        sx={{
          py: 2,
          px: 3,
          mt: 'auto',
          backgroundColor: 'grey.100',
          borderTop: 1,
          borderColor: 'grey.300',
        }}
      >
        <Container maxWidth="lg">
          <Typography variant="body2" color="text.secondary" align="center">
            PairWise © {new Date().getFullYear()} - Feature prioritization through group consensus
          </Typography>
        </Container>
      </Box>
    </Box>
  )
}

export default Layout