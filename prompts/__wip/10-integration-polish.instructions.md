# Instructions: Integration & Polish Implementation

## Application-Wide State Management

### Global Application Context
```typescript
// src/contexts/AppContext.tsx
import React, { createContext, useContext, useReducer, ReactNode } from 'react';

interface AppState {
  currentProject: Project | null;
  currentPhase: ProjectPhase;
  sessionData: {
    attendeeId?: number;
    role?: string;
  };
  preferences: {
    theme: 'light' | 'dark';
    notifications: boolean;
  };
}

type ProjectPhase = 
  | 'setup' 
  | 'features' 
  | 'pairwise-value' 
  | 'pairwise-complexity' 
  | 'fibonacci-value' 
  | 'fibonacci-complexity' 
  | 'results';

interface AppAction {
  type: 'SET_PROJECT' | 'SET_PHASE' | 'SET_SESSION_DATA' | 'UPDATE_PREFERENCES';
  payload: any;
}

const initialState: AppState = {
  currentProject: null,
  currentPhase: 'setup',
  sessionData: {},
  preferences: {
    theme: 'light',
    notifications: true,
  },
};

const appReducer = (state: AppState, action: AppAction): AppState => {
  switch (action.type) {
    case 'SET_PROJECT':
      return { ...state, currentProject: action.payload };
    case 'SET_PHASE':
      return { ...state, currentPhase: action.payload };
    case 'SET_SESSION_DATA':
      return { ...state, sessionData: { ...state.sessionData, ...action.payload } };
    case 'UPDATE_PREFERENCES':
      return { ...state, preferences: { ...state.preferences, ...action.payload } };
    default:
      return state;
  }
};

const AppContext = createContext<{
  state: AppState;
  dispatch: React.Dispatch<AppAction>;
} | null>(null);

export const AppProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(appReducer, initialState);

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
};

export const useAppContext = () => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useAppContext must be used within an AppProvider');
  }
  return context;
};
```

### Workflow State Machine
```typescript
// src/hooks/useWorkflowState.ts
import { useState, useEffect, useCallback } from 'react';
import { useAppContext } from '@/contexts/AppContext';

interface WorkflowStep {
  phase: ProjectPhase;
  title: string;
  description: string;
  completed: boolean;
  available: boolean;
  requirements: string[];
}

interface UseWorkflowStateOptions {
  projectId: number;
}

export const useWorkflowState = ({ projectId }: UseWorkflowStateOptions) => {
  const { state, dispatch } = useAppContext();
  const [steps, setSteps] = useState<WorkflowStep[]>([]);
  const [loading, setLoading] = useState(true);

  const defaultSteps: Omit<WorkflowStep, 'completed' | 'available'>[] = [
    {
      phase: 'setup',
      title: 'Project Setup',
      description: 'Configure project details and add team members',
      requirements: ['Project name', 'At least 2 attendees'],
    },
    {
      phase: 'features',
      title: 'Feature Input',
      description: 'Add or import features to be prioritized',
      requirements: ['At least 2 features with descriptions'],
    },
    {
      phase: 'pairwise-value',
      title: 'Value Comparisons',
      description: 'Compare features head-to-head for business value',
      requirements: ['All pairwise comparisons completed with consensus'],
    },
    {
      phase: 'pairwise-complexity',
      title: 'Complexity Comparisons', 
      description: 'Compare features head-to-head for implementation complexity',
      requirements: ['All pairwise comparisons completed with consensus'],
    },
    {
      phase: 'fibonacci-value',
      title: 'Value Scoring',
      description: 'Rate absolute business value using Fibonacci scale',
      requirements: ['All features scored with team consensus'],
    },
    {
      phase: 'fibonacci-complexity',
      title: 'Complexity Scoring',
      description: 'Rate absolute implementation complexity using Fibonacci scale',
      requirements: ['All features scored with team consensus'],
    },
    {
      phase: 'results',
      title: 'Final Results',
      description: 'View calculated priority rankings and export results',
      requirements: ['All previous phases completed'],
    },
  ];

  const checkStepCompletion = useCallback(async (phase: ProjectPhase): Promise<boolean> => {
    try {
      const response = await fetch(`/api/projects/${projectId}/workflow/status/${phase}`);
      if (response.ok) {
        const data = await response.json();
        return data.completed;
      }
    } catch (error) {
      console.error(`Failed to check completion for phase ${phase}:`, error);
    }
    return false;
  }, [projectId]);

  const updateWorkflowState = useCallback(async () => {
    setLoading(true);
    
    try {
      const completionPromises = defaultSteps.map(step => 
        checkStepCompletion(step.phase).then(completed => ({ ...step, completed }))
      );
      
      const stepsWithCompletion = await Promise.all(completionPromises);
      
      // Determine availability based on sequential completion
      const stepsWithAvailability = stepsWithCompletion.map((step, index) => ({
        ...step,
        available: index === 0 || stepsWithCompletion[index - 1].completed,
      }));
      
      setSteps(stepsWithAvailability);
    } catch (error) {
      console.error('Failed to update workflow state:', error);
    } finally {
      setLoading(false);
    }
  }, [checkStepCompletion]);

  const navigateToPhase = useCallback((phase: ProjectPhase) => {
    const step = steps.find(s => s.phase === phase);
    if (step && step.available) {
      dispatch({ type: 'SET_PHASE', payload: phase });
      return true;
    }
    return false;
  }, [steps, dispatch]);

  const getNextAvailablePhase = useCallback((): ProjectPhase | null => {
    const nextStep = steps.find(step => step.available && !step.completed);
    return nextStep?.phase || null;
  }, [steps]);

  useEffect(() => {
    updateWorkflowState();
  }, [updateWorkflowState]);

  return {
    steps,
    currentPhase: state.currentPhase,
    loading,
    navigateToPhase,
    getNextAvailablePhase,
    refreshWorkflow: updateWorkflowState,
  };
};
```

## Comprehensive Error Handling

### Global Error Handler
```typescript
// src/utils/errorHandler.ts
import { ApiError } from '@/types';

export class AppError extends Error {
  public readonly code: string;
  public readonly statusCode: number;
  public readonly isOperational: boolean;

  constructor(message: string, code: string = 'UNKNOWN_ERROR', statusCode: number = 500, isOperational: boolean = true) {
    super(message);
    
    this.code = code;
    this.statusCode = statusCode;
    this.isOperational = isOperational;
    
    Error.captureStackTrace(this, this.constructor);
  }
}

export const createAppError = (apiError: ApiError): AppError => {
  let code = 'API_ERROR';
  let statusCode = 500;

  // Map common API error patterns to specific error types
  if (apiError.error.includes('not found')) {
    code = 'NOT_FOUND';
    statusCode = 404;
  } else if (apiError.error.includes('validation')) {
    code = 'VALIDATION_ERROR';
    statusCode = 400;
  } else if (apiError.error.includes('consensus')) {
    code = 'CONSENSUS_ERROR';
    statusCode = 409;
  } else if (apiError.error.includes('network') || apiError.error.includes('connection')) {
    code = 'NETWORK_ERROR';
    statusCode = 503;
  }

  return new AppError(apiError.error, code, statusCode);
};

export const handleAsyncError = async <T>(
  operation: () => Promise<T>,
  context: string
): Promise<T> => {
  try {
    return await operation();
  } catch (error) {
    if (error instanceof AppError) {
      throw error;
    }
    
    const appError = error instanceof Error 
      ? new AppError(`${context}: ${error.message}`, 'OPERATION_FAILED')
      : new AppError(`${context}: Unknown error occurred`, 'UNKNOWN_ERROR');
    
    console.error(`Error in ${context}:`, error);
    throw appError;
  }
};

// Error reporting service
export class ErrorReportingService {
  private static instance: ErrorReportingService;
  private errorQueue: AppError[] = [];

  private constructor() {}

  public static getInstance(): ErrorReportingService {
    if (!ErrorReportingService.instance) {
      ErrorReportingService.instance = new ErrorReportingService();
    }
    return ErrorReportingService.instance;
  }

  public reportError(error: AppError, context?: any): void {
    this.errorQueue.push(error);
    
    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
      console.error('App Error:', {
        message: error.message,
        code: error.code,
        statusCode: error.statusCode,
        context,
        stack: error.stack,
      });
    }

    // In production, you might send to an error reporting service
    // this.sendToErrorService(error, context);
  }

  private sendToErrorService(error: AppError, context?: any): void {
    // Implementation for external error reporting (e.g., Sentry, Rollbar)
    // fetch('/api/errors', {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ error: error.message, code: error.code, context })
    // });
  }
}
```

### Error Recovery Components
```typescript
// src/components/ErrorHandling/ErrorRecovery.tsx
import React, { useState } from 'react';
import {
  Alert,
  AlertTitle,
  Box,
  Button,
  Card,
  CardContent,
  Collapse,
  Typography,
  IconButton,
} from '@mui/material';
import { 
  Refresh,
  ExpandMore,
  ExpandLess,
  BugReport,
  Home,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { AppError } from '@/utils/errorHandler';

interface ErrorRecoveryProps {
  error: AppError;
  onRetry?: () => void;
  onReset?: () => void;
  showDetails?: boolean;
}

export const ErrorRecovery: React.FC<ErrorRecoveryProps> = ({
  error,
  onRetry,
  onReset,
  showDetails = false,
}) => {
  const [expanded, setExpanded] = useState(false);
  const navigate = useNavigate();

  const getErrorSeverity = (error: AppError): 'error' | 'warning' | 'info' => {
    if (error.statusCode >= 500) return 'error';
    if (error.statusCode >= 400) return 'warning';
    return 'info';
  };

  const getRecoveryMessage = (error: AppError): string => {
    switch (error.code) {
      case 'NETWORK_ERROR':
        return 'Please check your internet connection and try again.';
      case 'VALIDATION_ERROR':
        return 'Please check your input and correct any validation errors.';
      case 'CONSENSUS_ERROR':
        return 'All team members need to reach consensus before proceeding.';
      case 'NOT_FOUND':
        return 'The requested resource could not be found.';
      default:
        return 'An unexpected error occurred. Please try again or contact support.';
    }
  };

  const getRecoveryActions = (error: AppError) => {
    const actions = [];

    if (onRetry) {
      actions.push(
        <Button
          key="retry"
          variant="contained"
          startIcon={<Refresh />}
          onClick={onRetry}
          sx={{ mr: 1 }}
        >
          Try Again
        </Button>
      );
    }

    if (onReset) {
      actions.push(
        <Button
          key="reset"
          variant="outlined"
          onClick={onReset}
          sx={{ mr: 1 }}
        >
          Reset
        </Button>
      );
    }

    actions.push(
      <Button
        key="home"
        variant="outlined"
        startIcon={<Home />}
        onClick={() => navigate('/')}
      >
        Go Home
      </Button>
    );

    return actions;
  };

  return (
    <Card sx={{ maxWidth: 600, mx: 'auto', mt: 4 }}>
      <CardContent>
        <Alert 
          severity={getErrorSeverity(error)}
          action={
            showDetails && (
              <IconButton
                size="small"
                onClick={() => setExpanded(!expanded)}
              >
                {expanded ? <ExpandLess /> : <ExpandMore />}
              </IconButton>
            )
          }
        >
          <AlertTitle>
            {error.code === 'NETWORK_ERROR' ? 'Connection Problem' : 'Something went wrong'}
          </AlertTitle>
          
          <Typography variant="body2" sx={{ mb: 2 }}>
            {getRecoveryMessage(error)}
          </Typography>

          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {getRecoveryActions(error)}
          </Box>
        </Alert>

        {showDetails && (
          <Collapse in={expanded}>
            <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
              <Typography variant="subtitle2" gutterBottom>
                Error Details:
              </Typography>
              <Typography variant="body2" sx={{ fontFamily: 'monospace', mb: 1 }}>
                Code: {error.code}
              </Typography>
              <Typography variant="body2" sx={{ fontFamily: 'monospace', mb: 1 }}>
                Message: {error.message}
              </Typography>
              <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                Status: {error.statusCode}
              </Typography>
              
              {process.env.NODE_ENV === 'development' && error.stack && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2" gutterBottom>
                    Stack Trace:
                  </Typography>
                  <Typography 
                    variant="body2" 
                    sx={{ 
                      fontFamily: 'monospace', 
                      fontSize: '0.75rem',
                      whiteSpace: 'pre-wrap',
                      bgcolor: 'grey.200',
                      p: 1,
                      borderRadius: 1,
                      overflow: 'auto',
                      maxHeight: 200,
                    }}
                  >
                    {error.stack}
                  </Typography>
                </Box>
              )}
            </Box>
          </Collapse>
        )}
      </CardContent>
    </Card>
  );
};
```

## Session Persistence and Recovery

### Session Storage Management
```typescript
// src/utils/sessionStorage.ts
interface SessionData {
  projectId: number;
  attendeeId: number;
  currentPhase: ProjectPhase;
  lastActivity: number;
  temporaryData: Record<string, any>;
}

export class SessionManager {
  private static readonly STORAGE_KEY = 'pwvc_session';
  private static readonly SESSION_TIMEOUT = 24 * 60 * 60 * 1000; // 24 hours

  static saveSession(data: Partial<SessionData>): void {
    try {
      const existing = this.getSession();
      const updated = {
        ...existing,
        ...data,
        lastActivity: Date.now(),
      };
      
      localStorage.setItem(this.STORAGE_KEY, JSON.stringify(updated));
    } catch (error) {
      console.warn('Failed to save session:', error);
    }
  }

  static getSession(): SessionData | null {
    try {
      const stored = localStorage.getItem(this.STORAGE_KEY);
      if (!stored) return null;

      const session: SessionData = JSON.parse(stored);
      
      // Check if session has expired
      if (Date.now() - session.lastActivity > this.SESSION_TIMEOUT) {
        this.clearSession();
        return null;
      }

      return session;
    } catch (error) {
      console.warn('Failed to retrieve session:', error);
      return null;
    }
  }

  static clearSession(): void {
    try {
      localStorage.removeItem(this.STORAGE_KEY);
    } catch (error) {
      console.warn('Failed to clear session:', error);
    }
  }

  static updateActivity(): void {
    const session = this.getSession();
    if (session) {
      this.saveSession({ lastActivity: Date.now() });
    }
  }

  static saveTemporaryData(key: string, data: any): void {
    const session = this.getSession();
    if (session) {
      const temporaryData = { ...session.temporaryData, [key]: data };
      this.saveSession({ temporaryData });
    }
  }

  static getTemporaryData(key: string): any {
    const session = this.getSession();
    return session?.temporaryData?.[key] || null;
  }

  static clearTemporaryData(key: string): void {
    const session = this.getSession();
    if (session && session.temporaryData) {
      delete session.temporaryData[key];
      this.saveSession({ temporaryData: session.temporaryData });
    }
  }
}

// Hook for session management
export const useSessionManager = () => {
  const { state, dispatch } = useAppContext();

  const saveCurrentSession = useCallback(() => {
    if (state.currentProject && state.sessionData.attendeeId) {
      SessionManager.saveSession({
        projectId: state.currentProject.id,
        attendeeId: state.sessionData.attendeeId,
        currentPhase: state.currentPhase,
      });
    }
  }, [state]);

  const restoreSession = useCallback(() => {
    const session = SessionManager.getSession();
    if (session) {
      dispatch({ type: 'SET_PHASE', payload: session.currentPhase });
      dispatch({ 
        type: 'SET_SESSION_DATA', 
        payload: { attendeeId: session.attendeeId } 
      });
      return session;
    }
    return null;
  }, [dispatch]);

  const clearSession = useCallback(() => {
    SessionManager.clearSession();
    dispatch({ type: 'SET_PROJECT', payload: null });
    dispatch({ type: 'SET_PHASE', payload: 'setup' });
    dispatch({ type: 'SET_SESSION_DATA', payload: {} });
  }, [dispatch]);

  return {
    saveCurrentSession,
    restoreSession,
    clearSession,
    saveTemporaryData: SessionManager.saveTemporaryData,
    getTemporaryData: SessionManager.getTemporaryData,
    clearTemporaryData: SessionManager.clearTemporaryData,
  };
};
```

## Docker Configuration and Deployment

### Multi-stage Dockerfile for Frontend
```dockerfile
# web/Dockerfile
FROM node:18-alpine as builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source code and build
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built app
COPY --from=builder /app/dist /usr/share/nginx/html

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:80/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### Nginx Configuration for React SPA
```nginx
# web/nginx.conf
events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    # Logging
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;
    
    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    server {
        listen 80;
        server_name _;
        root /usr/share/nginx/html;
        index index.html;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;

        # API proxy
        location /api {
            proxy_pass http://backend:8080;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # WebSocket proxy
        location /ws {
            proxy_pass http://backend:8080;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # React Router (handle client-side routing)
        location / {
            try_files $uri $uri/ /index.html;
        }

        # Static assets caching
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }

        # Health check endpoint
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
```

### Production Docker Compose
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB:-pwvc}
      POSTGRES_USER: ${POSTGRES_USER:-pwvc}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-pwvc}"]
      interval: 30s
      timeout: 10s
      retries: 3

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: postgres://${POSTGRES_USER:-pwvc}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB:-pwvc}?sslmode=disable
      PORT: 8080
      GIN_MODE: release
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile
    ports:
      - "${FRONTEND_PORT:-80}:80"
    depends_on:
      backend:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_data:

networks:
  default:
    name: pwvc_network
```

### Environment Configuration
```bash
# .env.example
# Database Configuration
POSTGRES_DB=pwvc
POSTGRES_USER=pwvc
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_PORT=5432

# Application Configuration
FRONTEND_PORT=80
BACKEND_PORT=8080

# Development/Production Mode
NODE_ENV=production
GIN_MODE=release

# Optional: External Service Configuration
# SENTRY_DSN=https://your-sentry-dsn
# ANALYTICS_KEY=your-analytics-key
```

### Deployment Scripts
```bash
#!/bin/bash
# deploy.sh

set -e

echo "🚀 Starting P-WVC deployment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found. Please copy .env.example to .env and configure."
    exit 1
fi

# Load environment variables
source .env

# Validate required environment variables
if [ -z "$POSTGRES_PASSWORD" ]; then
    echo "❌ POSTGRES_PASSWORD is required in .env file"
    exit 1
fi

# Build and start services
echo "📦 Building Docker images..."
docker-compose -f docker-compose.prod.yml build

echo "🗄️ Starting database..."
docker-compose -f docker-compose.prod.yml up -d postgres

echo "⏳ Waiting for database to be ready..."
sleep 10

echo "🔧 Running database migrations..."
docker-compose -f docker-compose.prod.yml exec postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -f /docker-entrypoint-initdb.d/001_create_tables.sql

echo "🚀 Starting all services..."
docker-compose -f docker-compose.prod.yml up -d

echo "🎉 Deployment complete!"
echo "📍 Application available at http://localhost:${FRONTEND_PORT:-80}"

# Health check
echo "🔍 Running health checks..."
sleep 5
docker-compose -f docker-compose.prod.yml ps
```

## Comprehensive Testing Strategy

### Integration Test Setup
```typescript
// src/tests/setup.ts
import { beforeAll, afterAll, afterEach } from '@jest/globals';
import { setupServer } from 'msw/node';
import { rest } from 'msw';

// Mock API server for testing
export const server = setupServer(
  // Project endpoints
  rest.get('/api/projects/:id', (req, res, ctx) => {
    return res(ctx.json({
      id: 1,
      name: 'Test Project',
      description: 'A test project',
      status: 'active',
    }));
  }),

  // Feature endpoints
  rest.get('/api/projects/:id/features', (req, res, ctx) => {
    return res(ctx.json([
      {
        id: 1,
        title: 'User Authentication',
        description: 'Users can log in and out',
        acceptance_criteria: 'Login form works correctly',
      },
      {
        id: 2,
        title: 'Dashboard',
        description: 'Main dashboard view',
        acceptance_criteria: 'Dashboard loads quickly',
      },
    ]));
  }),
);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});
```

### End-to-End Test Example
```typescript
// src/tests/e2e/complete-workflow.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { App } from '@/App';
import { theme } from '@/theme';
import { AppProvider } from '@/contexts/AppContext';

const renderApp = () => {
  return render(
    <AppProvider>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </ThemeProvider>
    </AppProvider>
  );
};

describe('Complete P-WVC Workflow', () => {
  test('should complete full prioritization workflow', async () => {
    renderApp();

    // 1. Create project
    fireEvent.click(screen.getByText('Create New Project'));
    
    fireEvent.change(screen.getByLabelText('Project Name'), {
      target: { value: 'Test Project' },
    });
    
    fireEvent.click(screen.getByText('Create Project'));
    
    await waitFor(() => {
      expect(screen.getByText('Project Setup')).toBeInTheDocument();
    });

    // 2. Add attendees
    fireEvent.click(screen.getByText('Add Attendee'));
    
    fireEvent.change(screen.getByLabelText('Name'), {
      target: { value: 'John Doe' },
    });
    
    fireEvent.click(screen.getByText('Add'));
    
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });

    // 3. Add features
    fireEvent.click(screen.getByText('Continue to Features'));
    
    fireEvent.click(screen.getByText('Add Feature'));
    
    fireEvent.change(screen.getByLabelText('Feature Title'), {
      target: { value: 'User Login' },
    });
    
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Users should be able to authenticate' },
    });
    
    fireEvent.click(screen.getByText('Save Feature'));
    
    await waitFor(() => {
      expect(screen.getByText('User Login')).toBeInTheDocument();
    });

    // Continue with more steps...
    // 4. Pairwise comparisons
    // 5. Fibonacci scoring  
    // 6. Results verification
  });
});
```