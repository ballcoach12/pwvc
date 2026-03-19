import { createContext, useContext, useEffect, useState } from 'react';
import { api } from '../services/api';

export const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check for existing authentication on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('AuthContext: Starting authentication check...');
      
      // Directly check authentication with server
      // Don't check for cookie since it's HttpOnly and can't be read by JS
      const response = await api.auth.getCurrentUser();
      
      console.log('AuthContext: Auth check response:', response);
      
      if (response && response.id) {
        console.log('AuthContext: User authenticated:', response);
        setUser(response);
      } else {
        console.log('AuthContext: No valid user response');
        setUser(null);
      }
    } catch (err) {
      console.log('AuthContext: Auth check failed (this is normal for unauthenticated users):', err.message);
      // User is not authenticated, this is normal for initial page load
      setUser(null);
      setError(null); // Don't set this as an error for unauthenticated state
    } finally {
      console.log('AuthContext: Auth check complete, loading set to false');
      setLoading(false);
    }
  };

  const login = async (credentials) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await api.auth.login(credentials);
      
      if (response && response.user) {
        setUser(response.user);
        return response;
      } else {
        throw new Error('Invalid response from server');
      }
    } catch (err) {
      const errorMessage = err.message || 'Login failed';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    try {
      setLoading(true);
      setError(null);
      
      await api.auth.logout();
    } catch (err) {
      console.warn('Logout API call failed:', err.message);
      // Continue with local logout even if API call fails
    } finally {
      // Clear local state regardless of API call result
      setUser(null);
      // Don't try to remove HttpOnly cookie - the server handles this
      setLoading(false);
    }
  };

  const hasRole = (role) => {
    if (!user || !user.roles) {
      return false;
    }
    return user.roles.includes(role);
  };

  const isAuthenticated = () => {
    const result = !!user;
    console.log('AuthContext: isAuthenticated called, user:', user, 'result:', result);
    return result;
  };

  const updateUser = (updatedUser) => {
    setUser(updatedUser);
  };

  const clearError = () => {
    setError(null);
  };

  const value = {
    user,
    loading,
    error,
    login,
    logout,
    hasRole,
    isAuthenticated,
    updateUser,
    clearError,
    checkAuth
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;