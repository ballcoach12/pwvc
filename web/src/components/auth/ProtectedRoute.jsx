import { Alert, Box, CircularProgress, Typography } from '@mui/material';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

/**
 * ProtectedRoute component for route-based authentication
 * Redirects to login if user is not authenticated
 * Supports role-based access control
 */
const ProtectedRoute = ({ 
  children, 
  requiredRole = null, 
  requiredRoles = [], 
  fallback = null,
  requireAny = false 
}) => {
  const { isAuthenticated, user, loading } = useAuth();
  const location = useLocation();

  console.log('ProtectedRoute: loading:', loading, 'isAuthenticated:', isAuthenticated(), 'user:', user, 'location:', location.pathname);

  // Show loading spinner while checking authentication
  if (loading) {
    console.log('ProtectedRoute: Showing loading spinner');
    return (
      <Box 
        display="flex" 
        flexDirection="column"
        justifyContent="center" 
        alignItems="center" 
        minHeight="300px"
        gap={2}
      >
        <CircularProgress size={40} />
        <Typography variant="body1" color="text.secondary">
          Verifying authentication...
        </Typography>
      </Box>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated()) {
    console.log('ProtectedRoute: User not authenticated, redirecting to /login');
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  console.log('ProtectedRoute: User authenticated, rendering children');

  // Check role-based access if required
  if (requiredRole || requiredRoles.length > 0) {
    const userRoles = user?.roles || [];
    
    // Build list of required roles
    const rolesToCheck = [];
    if (requiredRole) rolesToCheck.push(requiredRole);
    if (requiredRoles.length > 0) rolesToCheck.push(...requiredRoles);
    
    // Check if user has required roles
    let hasAccess = false;
    
    if (requireAny) {
      // User needs ANY of the required roles
      hasAccess = rolesToCheck.some(role => userRoles.includes(role));
    } else {
      // User needs ALL required roles (default behavior)
      hasAccess = rolesToCheck.every(role => userRoles.includes(role));
    }
    
    if (!hasAccess) {
      // Show custom fallback or default access denied message
      if (fallback) {
        return fallback;
      }
      
      return (
        <Box 
          display="flex" 
          flexDirection="column"
          justifyContent="center" 
          alignItems="center" 
          minHeight="300px"
          p={3}
        >
          <Alert severity="error" sx={{ maxWidth: 500 }}>
            <Typography variant="h6" gutterBottom>
              Access Denied
            </Typography>
            <Typography variant="body1">
              You don't have the required permissions to access this page.
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              Required role{rolesToCheck.length > 1 ? 's' : ''}: {rolesToCheck.join(requireAny ? ' OR ' : ', ')}
            </Typography>
            <Typography variant="body2">
              Your role{userRoles.length > 1 ? 's' : ''}: {userRoles.join(', ') || 'None'}
            </Typography>
          </Alert>
        </Box>
      );
    }
  }

  // User is authenticated and has required permissions
  return children;
};

/**
 * AdminRoute component - shorthand for admin-only routes
 */
export const AdminRoute = ({ children, fallback = null }) => {
  return (
    <ProtectedRoute requiredRole="admin" fallback={fallback}>
      {children}
    </ProtectedRoute>
  );
};

/**
 * UserRoute component - shorthand for user+ routes (user or admin)
 */
export const UserRoute = ({ children, fallback = null }) => {
  return (
    <ProtectedRoute requiredRoles={['user', 'admin']} requireAny={true} fallback={fallback}>
      {children}
    </ProtectedRoute>
  );
};

/**
 * ViewerRoute component - shorthand for any authenticated user
 */
export const ViewerRoute = ({ children, fallback = null }) => {
  return (
    <ProtectedRoute fallback={fallback}>
      {children}
    </ProtectedRoute>
  );
};

/**
 * ConditionalRoute component - renders different content based on roles
 */
export const ConditionalRoute = ({ 
  adminComponent, 
  userComponent, 
  viewerComponent, 
  fallback = null 
}) => {
  const { user } = useAuth();
  const userRoles = user?.roles || [];

  if (userRoles.includes('admin') && adminComponent) {
    return adminComponent;
  }
  
  if (userRoles.includes('user') && userComponent) {
    return userComponent;
  }
  
  if (viewerComponent) {
    return viewerComponent;
  }
  
  return fallback || (
    <Box p={3}>
      <Alert severity="info">
        <Typography variant="body1">
          No content available for your current role.
        </Typography>
      </Alert>
    </Box>
  );
};

export default ProtectedRoute;