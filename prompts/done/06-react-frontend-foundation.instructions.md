# Instructions: React Frontend Foundation Implementation

## Project Setup and Architecture

### Vite Configuration with TypeScript
```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@pages': path.resolve(__dirname, './src/pages'),
      '@services': path.resolve(__dirname, './src/services'),
      '@hooks': path.resolve(__dirname, './src/hooks'),
      '@types': path.resolve(__dirname, './src/types'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
})
```

### TypeScript Type Definitions
```typescript
// src/types/index.ts
export interface Project {
  id: number;
  name: string;
  description: string;
  status: 'active' | 'completed' | 'archived';
  created_at: string;
  updated_at: string;
  attendees?: Attendee[];
}

export interface Attendee {
  id: number;
  project_id: number;
  name: string;
  role: string;
  is_facilitator: boolean;
  created_at: string;
}

export interface Feature {
  id: number;
  project_id: number;
  title: string;
  description: string;
  acceptance_criteria: string;
  created_at: string;
  updated_at: string;
}

export interface PairwiseSession {
  id: number;
  project_id: number;
  criterion_type: 'value' | 'complexity';
  status: 'active' | 'completed';
  started_at: string;
  completed_at?: string;
}

export interface ApiError {
  error: string;
  details?: string[];
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}
```

## API Service Layer Architecture

### Base API Client with Error Handling
```typescript
// src/services/api.ts
import axios, { AxiosInstance, AxiosError, AxiosResponse } from 'axios';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: '/api',
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor for logging
    this.client.interceptors.request.use(
      (config) => {
        console.debug(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: AxiosError) => {
        const apiError = this.handleApiError(error);
        console.error('API Error:', apiError);
        return Promise.reject(apiError);
      }
    );
  }

  private handleApiError(error: AxiosError): ApiError {
    if (error.response) {
      // Server responded with error status
      const data = error.response.data as any;
      return {
        error: data.error || data.message || 'Server error occurred',
        details: data.details || [],
      };
    } else if (error.request) {
      // Network error
      return {
        error: 'Network error - please check your connection',
      };
    } else {
      // Request setup error
      return {
        error: 'Request configuration error',
      };
    }
  }

  async get<T>(url: string): Promise<T> {
    const response = await this.client.get<T>(url);
    return response.data;
  }

  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.post<T>(url, data);
    return response.data;
  }

  async put<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.put<T>(url, data);
    return response.data;
  }

  async delete<T>(url: string): Promise<T> {
    const response = await this.client.delete<T>(url);
    return response.data;
  }

  async uploadFile<T>(url: string, formData: FormData): Promise<T> {
    const response = await this.client.post<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  }
}

export const apiClient = new ApiClient();
```

### Domain-Specific Service Classes
```typescript
// src/services/projectService.ts
import { apiClient } from './api';
import { Project, Attendee } from '@/types';

export class ProjectService {
  async createProject(data: { name: string; description: string }): Promise<Project> {
    return apiClient.post<Project>('/projects', data);
  }

  async getProject(id: number): Promise<Project> {
    return apiClient.get<Project>(`/projects/${id}`);
  }

  async updateProject(id: number, data: Partial<Project>): Promise<Project> {
    return apiClient.put<Project>(`/projects/${id}`, data);
  }

  async deleteProject(id: number): Promise<void> {
    return apiClient.delete<void>(`/projects/${id}`);
  }

  async addAttendee(projectId: number, data: { name: string; role: string; is_facilitator: boolean }): Promise<Attendee> {
    return apiClient.post<Attendee>(`/projects/${projectId}/attendees`, data);
  }

  async getAttendees(projectId: number): Promise<Attendee[]> {
    return apiClient.get<Attendee[]>(`/projects/${projectId}/attendees`);
  }

  async removeAttendee(projectId: number, attendeeId: number): Promise<void> {
    return apiClient.delete<void>(`/projects/${projectId}/attendees/${attendeeId}`);
  }
}

// src/services/featureService.ts
export class FeatureService {
  async createFeature(projectId: number, data: { title: string; description: string; acceptance_criteria?: string }): Promise<Feature> {
    return apiClient.post<Feature>(`/projects/${projectId}/features`, data);
  }

  async getFeatures(projectId: number): Promise<Feature[]> {
    return apiClient.get<Feature[]>(`/projects/${projectId}/features`);
  }

  async updateFeature(projectId: number, featureId: number, data: Partial<Feature>): Promise<Feature> {
    return apiClient.put<Feature>(`/projects/${projectId}/features/${featureId}`, data);
  }

  async deleteFeature(projectId: number, featureId: number): Promise<void> {
    return apiClient.delete<void>(`/projects/${projectId}/features/${featureId}`);
  }

  async importFeaturesFromCSV(projectId: number, file: File): Promise<{ imported: number; features: Feature[] }> {
    const formData = new FormData();
    formData.append('csv_file', file);
    
    return apiClient.uploadFile<{ imported: number; features: Feature[] }>(
      `/projects/${projectId}/features/import`,
      formData
    );
  }

  async exportFeaturesToCSV(projectId: number): Promise<Blob> {
    const response = await fetch(`/api/projects/${projectId}/features/export`);
    return response.blob();
  }
}

export const projectService = new ProjectService();
export const featureService = new FeatureService();
```

## Custom React Hooks

### Data Fetching Hooks with Loading States
```typescript
// src/hooks/useProject.ts
import { useState, useEffect } from 'react';
import { Project, ApiError } from '@/types';
import { projectService } from '@/services/projectService';

interface UseProjectResult {
  project: Project | null;
  loading: boolean;
  error: ApiError | null;
  refetch: () => void;
}

export const useProject = (projectId: number): UseProjectResult => {
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);

  const fetchProject = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await projectService.getProject(projectId);
      setProject(data);
    } catch (err) {
      setError(err as ApiError);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (projectId) {
      fetchProject();
    }
  }, [projectId]);

  return {
    project,
    loading,
    error,
    refetch: fetchProject,
  };
};

// src/hooks/useFeatures.ts
export const useFeatures = (projectId: number) => {
  const [features, setFeatures] = useState<Feature[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);

  const fetchFeatures = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await featureService.getFeatures(projectId);
      setFeatures(data);
    } catch (err) {
      setError(err as ApiError);
    } finally {
      setLoading(false);
    }
  };

  const addFeature = async (featureData: { title: string; description: string; acceptance_criteria?: string }) => {
    try {
      const newFeature = await featureService.createFeature(projectId, featureData);
      setFeatures(prev => [...prev, newFeature]);
      return newFeature;
    } catch (err) {
      throw err;
    }
  };

  const updateFeature = async (featureId: number, updates: Partial<Feature>) => {
    try {
      const updatedFeature = await featureService.updateFeature(projectId, featureId, updates);
      setFeatures(prev => 
        prev.map(f => f.id === featureId ? updatedFeature : f)
      );
      return updatedFeature;
    } catch (err) {
      throw err;
    }
  };

  const deleteFeature = async (featureId: number) => {
    try {
      await featureService.deleteFeature(projectId, featureId);
      setFeatures(prev => prev.filter(f => f.id !== featureId));
    } catch (err) {
      throw err;
    }
  };

  useEffect(() => {
    if (projectId) {
      fetchFeatures();
    }
  }, [projectId]);

  return {
    features,
    loading,
    error,
    refetch: fetchFeatures,
    addFeature,
    updateFeature,
    deleteFeature,
  };
};
```

## Material-UI Theme and Layout

### Custom Theme Configuration
```typescript
// src/theme/index.ts
import { createTheme } from '@mui/material/styles';
import { deepPurple, amber, red } from '@mui/material/colors';

export const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: deepPurple[600],
      light: deepPurple[400],
      dark: deepPurple[800],
    },
    secondary: {
      main: amber[500],
      light: amber[300],
      dark: amber[700],
    },
    error: {
      main: red[600],
    },
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
    },
  },
  typography: {
    h4: {
      fontWeight: 600,
      marginBottom: '1rem',
    },
    h5: {
      fontWeight: 500,
      marginBottom: '0.75rem',
    },
    body1: {
      lineHeight: 1.6,
    },
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.07)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 500,
        },
      },
    },
  },
});
```

### Layout Components
```typescript
// src/components/Layout/AppLayout.tsx
import React from 'react';
import { Box, AppBar, Toolbar, Typography, Container } from '@mui/material';
import { Outlet, useLocation, Link } from 'react-router-dom';
import { Breadcrumbs } from './Breadcrumbs';
import { ErrorBoundary } from './ErrorBoundary';

export const AppLayout: React.FC = () => {
  const location = useLocation();

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar position="sticky" elevation={1}>
        <Toolbar>
          <Typography
            variant="h6"
            component={Link}
            to="/"
            sx={{
              flexGrow: 1,
              textDecoration: 'none',
              color: 'inherit',
              fontWeight: 600,
            }}
          >
            P-WVC Prioritization Tool
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ flex: 1, py: 3 }}>
        <Breadcrumbs />
        
        <ErrorBoundary>
          <Outlet />
        </ErrorBoundary>
      </Container>
    </Box>
  );
};

// src/components/Layout/ErrorBoundary.tsx
import React, { Component, ReactNode } from 'react';
import { Alert, Button, Box } from '@mui/material';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: any) {
    console.error('Error boundary caught an error:', error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null });
  };

  render() {
    if (this.state.hasError) {
      return (
        <Box sx={{ mt: 4 }}>
          <Alert 
            severity="error"
            action={
              <Button color="inherit" size="small" onClick={this.handleReset}>
                Try Again
              </Button>
            }
          >
            Something went wrong: {this.state.error?.message}
          </Alert>
        </Box>
      );
    }

    return this.props.children;
  }
}
```

## Form Handling and Validation

### Reusable Form Components
```typescript
// src/components/Forms/FeatureForm.tsx
import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Alert,
} from '@mui/material';
import { Feature, ApiError } from '@/types';

interface FeatureFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: FeatureFormData) => Promise<void>;
  initialData?: Feature;
  title: string;
}

interface FeatureFormData {
  title: string;
  description: string;
  acceptance_criteria: string;
}

export const FeatureForm: React.FC<FeatureFormProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
  title,
}) => {
  const [formData, setFormData] = useState<FeatureFormData>({
    title: initialData?.title || '',
    description: initialData?.description || '',
    acceptance_criteria: initialData?.acceptance_criteria || '',
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.title.trim()) {
      newErrors.title = 'Title is required';
    } else if (formData.title.length < 3) {
      newErrors.title = 'Title must be at least 3 characters';
    }

    if (!formData.description.trim()) {
      newErrors.description = 'Description is required';
    } else if (formData.description.length < 10) {
      newErrors.description = 'Description must be at least 10 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setSubmitting(true);
    setSubmitError(null);

    try {
      await onSubmit(formData);
      onClose();
      setFormData({ title: '', description: '', acceptance_criteria: '' });
    } catch (error) {
      const apiError = error as ApiError;
      setSubmitError(apiError.error);
    } finally {
      setSubmitting(false);
    }
  };

  const handleChange = (field: keyof FeatureFormData) => (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    setFormData(prev => ({ ...prev, [field]: e.target.value }));
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{title}</DialogTitle>
        
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            {submitError && (
              <Alert severity="error">{submitError}</Alert>
            )}
            
            <TextField
              label="Feature Title"
              value={formData.title}
              onChange={handleChange('title')}
              error={!!errors.title}
              helperText={errors.title}
              required
              fullWidth
            />
            
            <TextField
              label="Description"
              value={formData.description}
              onChange={handleChange('description')}
              error={!!errors.description}
              helperText={errors.description || 'Describe what this feature does and why it\'s valuable'}
              multiline
              rows={3}
              required
              fullWidth
            />
            
            <TextField
              label="Acceptance Criteria"
              value={formData.acceptance_criteria}
              onChange={handleChange('acceptance_criteria')}
              helperText="Define what success looks like for this feature"
              multiline
              rows={3}
              fullWidth
            />
          </Box>
        </DialogContent>
        
        <DialogActions>
          <Button onClick={onClose} disabled={submitting}>
            Cancel
          </Button>
          <Button 
            type="submit" 
            variant="contained" 
            disabled={submitting}
          >
            {submitting ? 'Saving...' : 'Save Feature'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
```

## File Upload and CSV Handling

### CSV Import Component
```typescript
// src/components/Features/CSVImport.tsx
import React, { useState } from 'react';
import {
  Box,
  Button,
  Typography,
  Alert,
  LinearProgress,
  Chip,
} from '@mui/material';
import { CloudUpload } from '@mui/icons-material';
import { featureService } from '@/services/featureService';
import { ApiError } from '@/types';

interface CSVImportProps {
  projectId: number;
  onImportComplete: (count: number) => void;
}

export const CSVImport: React.FC<CSVImportProps> = ({
  projectId,
  onImportComplete,
}) => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [importing, setImporting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    
    if (file) {
      if (!file.name.endsWith('.csv')) {
        setError('Please select a CSV file');
        return;
      }
      
      if (file.size > 5 * 1024 * 1024) { // 5MB limit
        setError('File size must be less than 5MB');
        return;
      }
      
      setSelectedFile(file);
      setError(null);
      setSuccess(null);
    }
  };

  const handleImport = async () => {
    if (!selectedFile) return;

    setImporting(true);
    setError(null);
    setSuccess(null);

    try {
      const result = await featureService.importFeaturesFromCSV(projectId, selectedFile);
      setSuccess(`Successfully imported ${result.imported} features`);
      onImportComplete(result.imported);
      setSelectedFile(null);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError.error);
    } finally {
      setImporting(false);
    }
  };

  return (
    <Box sx={{ p: 3, border: '2px dashed', borderColor: 'grey.300', borderRadius: 2 }}>
      <Typography variant="h6" gutterBottom>
        Import Features from CSV
      </Typography>
      
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        CSV should have columns: title, description, acceptance_criteria
      </Typography>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      {success && <Alert severity="success" sx={{ mb: 2 }}>{success}</Alert>}
      
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
        <Button
          component="label"
          variant="outlined"
          startIcon={<CloudUpload />}
          disabled={importing}
        >
          Select CSV File
          <input
            type="file"
            hidden
            accept=".csv"
            onChange={handleFileSelect}
          />
        </Button>
        
        {selectedFile && (
          <Chip
            label={selectedFile.name}
            onDelete={() => setSelectedFile(null)}
            color="primary"
          />
        )}
      </Box>

      {importing && <LinearProgress sx={{ mb: 2 }} />}

      <Button
        variant="contained"
        onClick={handleImport}
        disabled={!selectedFile || importing}
      >
        {importing ? 'Importing...' : 'Import Features'}
      </Button>
    </Box>
  );
};
```

## Navigation and Routing

### React Router Setup
```typescript
// src/App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import { theme } from './theme';
import { AppLayout } from './components/Layout/AppLayout';
import { ProjectList } from './pages/ProjectList';
import { ProjectSetup } from './pages/ProjectSetup';
import { FeatureManagement } from './pages/FeatureManagement';
import { PairwisePhase } from './pages/PairwisePhase';

export const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<AppLayout />}>
            <Route index element={<ProjectList />} />
            <Route path="projects/new" element={<ProjectSetup />} />
            <Route path="projects/:id" element={<ProjectSetup />} />
            <Route path="projects/:id/features" element={<FeatureManagement />} />
            <Route path="projects/:id/pairwise/:criterion" element={<PairwisePhase />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
};
```