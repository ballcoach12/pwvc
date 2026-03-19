import { api } from './api';

// Authentication service for JWT-based auth
export const authService = {
  /**
   * Login with username and password
   * @param {Object} credentials - Login credentials
   * @param {string} credentials.username - Username
   * @param {string} credentials.password - Password
   * @returns {Promise<Object>} Login response with user data
   */
  login: async (credentials) => {
    try {
      const response = await api.auth.login(credentials);
      return response;
    } catch (error) {
      throw new Error(error.message || 'Login failed');
    }
  },

  /**
   * Logout the current user
   * @returns {Promise<void>}
   */
  logout: async () => {
    try {
      await api.auth.logout();
    } catch (error) {
      // Continue with logout even if API call fails
      console.warn('Logout API call failed:', error.message);
    }
  },

  /**
   * Get current authenticated user
   * @returns {Promise<Object>} User profile data
   */
  getCurrentUser: async () => {
    try {
      const response = await api.auth.getCurrentUser();
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to get current user');
    }
  },

  /**
   * Check if user is authenticated (has valid token)
   * @returns {boolean} True if authenticated
   */
  isAuthenticated: () => {
    // Check for auth cookie
    const cookies = document.cookie.split(';');
    const authCookie = cookies.find(cookie => 
      cookie.trim().startsWith('auth=')
    );
    return !!authCookie;
  },

  /**
   * Change password for current user
   * @param {Object} passwordData - Password change data
   * @param {string} passwordData.currentPassword - Current password
   * @param {string} passwordData.newPassword - New password
   * @returns {Promise<void>}
   */
  changePassword: async (passwordData) => {
    try {
      await api.post('/auth/change-password', passwordData);
    } catch (error) {
      throw new Error(error.message || 'Failed to change password');
    }
  }
};

// Admin service for user management (requires admin role)
export const adminService = {
  /**
   * Get all users with pagination
   * @param {Object} params - Query parameters
   * @param {number} params.page - Page number (1-based)
   * @param {number} params.limit - Number of users per page
   * @param {string} params.search - Search term for username/email
   * @returns {Promise<Object>} Users list with pagination info
   */
  getUsers: async (params = {}) => {
    try {
      const response = await api.admin.getUsers(params);
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to get users');
    }
  },

  /**
   * Create a new user
   * @param {Object} userData - User creation data
   * @param {string} userData.username - Username
   * @param {string} userData.password - Password
   * @param {string} userData.email - Email (optional)
   * @param {string[]} userData.roles - Array of role names
   * @returns {Promise<Object>} Created user data
   */
  createUser: async (userData) => {
    try {
      const response = await api.admin.createUser(userData);
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to create user');
    }
  },

  /**
   * Get a specific user by ID
   * @param {number} userId - User ID
   * @returns {Promise<Object>} User data
   */
  getUser: async (userId) => {
    try {
      const response = await api.admin.getUser(userId);
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to get user');
    }
  },

  /**
   * Update a user
   * @param {number} userId - User ID
   * @param {Object} userData - User update data
   * @param {string} userData.email - Email (optional)
   * @param {boolean} userData.is_active - Active status (optional)
   * @param {string[]} userData.roles - Array of role names (optional)
   * @returns {Promise<Object>} Updated user data
   */
  updateUser: async (userId, userData) => {
    try {
      const response = await api.admin.updateUser(userId, userData);
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to update user');
    }
  },

  /**
   * Get all available roles
   * @returns {Promise<Array>} List of roles
   */
  getRoles: async () => {
    try {
      const response = await api.admin.getRoles();
      return response;
    } catch (error) {
      throw new Error(error.message || 'Failed to get roles');
    }
  },

  /**
   * Delete a user
   * @param {number} userId - User ID
   * @returns {Promise<void>}
   */
  deleteUser: async (userId) => {
    try {
      await api.delete(`/admin/users/${userId}`);
    } catch (error) {
      throw new Error(error.message || 'Failed to delete user');
    }
  },

  /**
   * Reset a user's password (admin only)
   * @param {number} userId - User ID
   * @param {Object} passwordData - Password reset data
   * @param {string} passwordData.newPassword - New password
   * @returns {Promise<void>}
   */
  resetPassword: async (userId, passwordData) => {
    try {
      await api.post(`/admin/users/${userId}/reset-password`, passwordData);
    } catch (error) {
      throw new Error(error.message || 'Failed to reset password');
    }
  }
};

// Utility functions
export const authUtils = {
  /**
   * Check if current user has a specific role
   * @param {string} role - Role name to check
   * @param {Object} user - User object with roles array
   * @returns {boolean} True if user has the role
   */
  hasRole: (role, user) => {
    if (!user || !user.roles || !Array.isArray(user.roles)) {
      return false;
    }
    return user.roles.includes(role);
  },

  /**
   * Check if current user has admin role
   * @param {Object} user - User object with roles array
   * @returns {boolean} True if user is admin
   */
  isAdmin: (user) => {
    return authUtils.hasRole('admin', user);
  },

  /**
   * Get user initials for avatar display
   * @param {Object} user - User object
   * @returns {string} User initials (max 2 characters)
   */
  getUserInitials: (user) => {
    if (!user || !user.username) {
      return '?';
    }
    return user.username.substring(0, 2).toUpperCase();
  },

  /**
   * Format user display name
   * @param {Object} user - User object
   * @returns {string} Formatted display name
   */
  getDisplayName: (user) => {
    if (!user) {
      return 'Unknown User';
    }
    return user.username || 'Unknown User';
  },

  /**
   * Validate password strength
   * @param {string} password - Password to validate
   * @returns {Object} Validation result with isValid and errors
   */
  validatePassword: (password) => {
    const errors = [];
    
    if (!password || password.length < 6) {
      errors.push('Password must be at least 6 characters long');
    }
    
    if (password && !/[A-Z]/.test(password)) {
      errors.push('Password must contain at least one uppercase letter');
    }
    
    if (password && !/[a-z]/.test(password)) {
      errors.push('Password must contain at least one lowercase letter');
    }
    
    if (password && !/[0-9]/.test(password)) {
      errors.push('Password must contain at least one number');
    }
    
    return {
      isValid: errors.length === 0,
      errors
    };
  },

  /**
   * Format role name for display
   * @param {string} role - Role name
   * @returns {string} Formatted role name
   */
  formatRole: (role) => {
    if (!role) return '';
    return role.charAt(0).toUpperCase() + role.slice(1);
  }
};

export default authService;