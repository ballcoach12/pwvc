import {
    Lock,
    Person,
    Visibility,
    VisibilityOff
} from '@mui/icons-material';
import {
    Alert,
    Box,
    Button,
    CircularProgress,
    IconButton,
    InputAdornment,
    Paper,
    TextField,
    Typography
} from '@mui/material';
import { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';

/**
 * LoginForm component with Material-UI styling
 * Handles user authentication with username and password
 */
const LoginForm = ({ onSuccess, onError }) => {
  const { login, loading, error } = useAuth();
  
  // Form state
  const [formData, setFormData] = useState({
    username: '', // Will be used as attendee ID
    password: ''  // Will be used as PIN
  });
  
  // UI state
  const [showPassword, setShowPassword] = useState(false);
  const [formErrors, setFormErrors] = useState({});
  const [submitAttempted, setSubmitAttempted] = useState(false);

  // Handle input changes
  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    // Clear field error when user starts typing
    if (formErrors[name]) {
      setFormErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  // Validate form fields
  const validateForm = () => {
    const errors = {};

    if (!formData.username.trim()) {
      errors.username = 'Username is required';
    } else if (formData.username.trim().length < 3) {
      errors.username = 'Username must be at least 3 characters';
    }

    if (!formData.password) {
      errors.password = 'Password is required';
    } else if (formData.password.length < 8) {
      errors.password = 'Password must be at least 8 characters';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Handle form submission
  const handleSubmit = async (event) => {
    event.preventDefault();
    setSubmitAttempted(true);

    // Validate form
    if (!validateForm()) {
      return;
    }

    try {
      await login(formData);
      
      // Call success callback if provided
      if (onSuccess) {
        onSuccess();
      }
    } catch (err) {
      // Call error callback if provided
      if (onError) {
        onError(err);
      }
    }
  };

  // Toggle password visibility
  const handleTogglePasswordVisibility = () => {
    setShowPassword(prev => !prev);
  };

  // Handle Enter key press
  const handleKeyPress = (event) => {
    if (event.key === 'Enter') {
      handleSubmit(event);
    }
  };

  return (
    <Paper elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto' }}>
      <Box component="form" onSubmit={handleSubmit} noValidate>
        <Typography variant="h4" component="h1" gutterBottom align="center" color="primary">
          Sign In
        </Typography>
        
        <Typography variant="subtitle1" gutterBottom align="center" color="text.secondary" sx={{ mb: 3 }}>
          Welcome to PairWise
        </Typography>

        {/* Display authentication error */}
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error.message || 'Login failed. Please check your credentials.'}
          </Alert>
        )}

        {/* Username field */}
        <TextField
          fullWidth
          name="username"
          label="Username"
          value={formData.username}
          onChange={handleChange}
          onKeyPress={handleKeyPress}
          error={Boolean(formErrors.username)}
          helperText={formErrors.username}
          disabled={loading}
          autoComplete="username"
          autoFocus
          sx={{ mb: 2 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Person color="action" />
              </InputAdornment>
            ),
          }}
        />

        {/* Password field */}
        <TextField
          fullWidth
          name="password"
          label="Password"
          type={showPassword ? 'text' : 'password'}
          value={formData.password}
          onChange={handleChange}
          onKeyPress={handleKeyPress}
          error={Boolean(formErrors.password)}
          helperText={formErrors.password}
          disabled={loading}
          autoComplete="current-password"
          sx={{ mb: 3 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Lock color="action" />
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                <IconButton
                  aria-label="toggle password visibility"
                  onClick={handleTogglePasswordVisibility}
                  edge="end"
                  disabled={loading}
                >
                  {showPassword ? <VisibilityOff /> : <Visibility />}
                </IconButton>
              </InputAdornment>
            ),
          }}
        />

        {/* Submit button */}
        <Button
          type="submit"
          fullWidth
          variant="contained"
          size="large"
          disabled={loading}
          sx={{ mb: 2, py: 1.5 }}
        >
          {loading ? (
            <Box display="flex" alignItems="center" gap={1}>
              <CircularProgress size={20} color="inherit" />
              Signing In...
            </Box>
          ) : (
            'Sign In'
          )}
        </Button>

        {/* Additional info */}
        <Typography variant="body2" align="center" color="text.secondary">
          Contact your administrator for account access
        </Typography>
      </Box>
    </Paper>
  );
};

export default LoginForm;