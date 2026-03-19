import {
    Alert,
    Box,
    Container,
    Fade,
    Typography,
    useMediaQuery,
    useTheme
} from '@mui/material';
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import LoginForm from '../components/auth/LoginForm';
import { useAuth } from '../hooks/useAuth';

/**
 * LoginPage component
 * Provides authentication interface with redirect handling
 */
const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  const { isAuthenticated, loading } = useAuth();
  const [loginError, setLoginError] = useState(null);
  const [showForm, setShowForm] = useState(false);

  // Get the intended destination after login
  const from = location.state?.from?.pathname || '/';

  // Redirect if already authenticated
  useEffect(() => {
    if (!loading && isAuthenticated) {
      navigate(from, { replace: true });
    } else if (!loading) {
      // Show form with a slight delay for better UX
      setTimeout(() => setShowForm(true), 100);
    }
  }, [isAuthenticated, loading, navigate, from]);

  // Handle successful login
  const handleLoginSuccess = () => {
    setLoginError(null);
    // Navigation will be handled by the useEffect above
  };

  // Handle login error
  const handleLoginError = (error) => {
    setLoginError(error.message || 'Login failed. Please try again.');
  };

  // Show loading while checking authentication
  if (loading) {
    return (
      <Container maxWidth="sm" sx={{ mt: 8 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="300px">
          <Typography variant="h6" color="text.secondary">
            Loading...
          </Typography>
        </Box>
      </Container>
    );
  }

  // Don't show login form if already authenticated
  if (isAuthenticated) {
    return null;
  }

  return (
    <Container maxWidth="sm" sx={{ mt: isMobile ? 4 : 8, px: isMobile ? 2 : 3 }}>
      <Fade in={showForm} timeout={500}>
        <Box>
          {/* Page header */}
          <Box textAlign="center" mb={4}>
            <Typography variant="h3" component="h1" gutterBottom color="primary" fontWeight="bold">
              PairWise
            </Typography>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              Feature Prioritization Through Group Consensus
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Please sign in to continue
            </Typography>
          </Box>

          {/* Show redirect notice if user was redirected to login */}
          {location.state?.from && (
            <Alert severity="info" sx={{ mb: 3 }}>
              Please sign in to access {location.state.from.pathname}
            </Alert>
          )}

          {/* Show login error if any */}
          {loginError && (
            <Alert severity="error" sx={{ mb: 3 }} onClose={() => setLoginError(null)}>
              {loginError}
            </Alert>
          )}

          {/* Login form */}
          <LoginForm
            onSuccess={handleLoginSuccess}
            onError={handleLoginError}
          />

          {/* Footer info */}
          <Box mt={4} textAlign="center">
            <Typography variant="body2" color="text.secondary">
              © 2025 PairWise. Secure authentication system.
            </Typography>
          </Box>
        </Box>
      </Fade>
    </Container>
  );
};

export default LoginPage;