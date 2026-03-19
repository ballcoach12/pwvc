import { useContext } from 'react';
import { AuthContext } from '../contexts/AuthContext';

/**
 * Custom hook for accessing authentication context
 * Provides centralized access to auth state and actions
 * 
 * @returns {Object} Authentication context with state and methods
 * @throws {Error} If used outside of AuthProvider
 */
export const useAuth = () => {
  const context = useContext(AuthContext);
  
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  
  return context;
};

/**
 * Hook for getting current authenticated user
 * Automatically handles loading and error states
 * 
 * @returns {Object} User state with data, loading, and error
 */
export const useCurrentUser = () => {
  const { user, loading, error } = useAuth();
  
  return {
    user,
    loading,
    error,
    isAuthenticated: !!user,
    isAdmin: user?.roles?.includes('admin') || false,
    hasRole: (role) => user?.roles?.includes(role) || false
  };
};

/**
 * Hook for authentication actions
 * Provides login, logout, and other auth methods
 * 
 * @returns {Object} Authentication actions and loading states
 */
export const useAuthActions = () => {
  const { login, logout, refreshUser, loading, error } = useAuth();
  
  return {
    login,
    logout,
    refreshUser,
    loading,
    error
  };
};

/**
 * Hook for role-based access control
 * Provides methods for checking user permissions
 * 
 * @returns {Object} Role checking utilities
 */
export const usePermissions = () => {
  const { user } = useAuth();
  
  return {
    hasRole: (role) => {
      if (!user || !user.roles || !Array.isArray(user.roles)) {
        return false;
      }
      return user.roles.includes(role);
    },
    
    hasAnyRole: (roles) => {
      if (!user || !user.roles || !Array.isArray(user.roles)) {
        return false;
      }
      return roles.some(role => user.roles.includes(role));
    },
    
    hasAllRoles: (roles) => {
      if (!user || !user.roles || !Array.isArray(user.roles)) {
        return false;
      }
      return roles.every(role => user.roles.includes(role));
    },
    
    isAdmin: () => {
      return user?.roles?.includes('admin') || false;
    },
    
    isUser: () => {
      return user?.roles?.includes('user') || false;
    },
    
    isViewer: () => {
      return user?.roles?.includes('viewer') || false;
    },
    
    canAccess: (requiredRole) => {
      if (!requiredRole) return true;
      return user?.roles?.includes(requiredRole) || false;
    }
  };
};

/**
 * Hook for authentication state only (no actions)
 * Useful for components that only need to read auth state
 * 
 * @returns {Object} Authentication state
 */
export const useAuthState = () => {
  const { user, loading, error, isAuthenticated } = useAuth();
  
  return {
    user,
    loading,
    error,
    isAuthenticated
  };
};

/**
 * Hook for conditional rendering based on authentication
 * Returns render helpers for common auth patterns
 * 
 * @returns {Object} Conditional rendering helpers
 */
export const useAuthRender = () => {
  const { user, isAuthenticated } = useAuth();
  
  return {
    renderIfAuthenticated: (component) => {
      return isAuthenticated ? component : null;
    },
    
    renderIfUnauthenticated: (component) => {
      return !isAuthenticated ? component : null;
    },
    
    renderIfRole: (role, component) => {
      const hasRole = user?.roles?.includes(role) || false;
      return hasRole ? component : null;
    },
    
    renderIfAdmin: (component) => {
      const isAdmin = user?.roles?.includes('admin') || false;
      return isAdmin ? component : null;
    },
    
    renderForUser: (userId, component) => {
      return user?.id === userId ? component : null;
    }
  };
};

export default useAuth;