# Instructions: Pairwise Comparison UI Implementation

## Real-time WebSocket Integration

### WebSocket Hook for Live Updates
```typescript
// src/hooks/useWebSocket.ts
import { useEffect, useRef, useState, useCallback } from 'react';

export interface WebSocketMessage {
  type: string;
  session_id: number;
  attendee_id: number;
  data: any;
  timestamp: string;
}

interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
}

export const useWebSocket = (url: string, options: UseWebSocketOptions = {}) => {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const ws = useRef<WebSocket | null>(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;

  const connect = useCallback(() => {
    try {
      ws.current = new WebSocket(url);
      
      ws.current.onopen = () => {
        setConnected(true);
        setError(null);
        reconnectAttempts.current = 0;
        options.onConnect?.();
      };

      ws.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          options.onMessage?.(message);
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err);
        }
      };

      ws.current.onclose = () => {
        setConnected(false);
        options.onDisconnect?.();
        
        // Attempt reconnection
        if (reconnectAttempts.current < maxReconnectAttempts) {
          setTimeout(() => {
            reconnectAttempts.current++;
            connect();
          }, Math.pow(2, reconnectAttempts.current) * 1000); // Exponential backoff
        }
      };

      ws.current.onerror = (error) => {
        setError('WebSocket connection error');
        options.onError?.(error);
      };
    } catch (err) {
      setError('Failed to establish WebSocket connection');
    }
  }, [url, options]);

  const sendMessage = useCallback((message: Partial<WebSocketMessage>) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
    }
  }, []);

  const disconnect = useCallback(() => {
    reconnectAttempts.current = maxReconnectAttempts; // Prevent reconnection
    ws.current?.close();
  }, []);

  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  return {
    connected,
    error,
    sendMessage,
    disconnect,
  };
};
```

### Pairwise Session State Management
```typescript
// src/hooks/usePairwiseSession.ts
import { useState, useEffect, useCallback } from 'react';
import { useWebSocket } from './useWebSocket';

interface PairwiseComparison {
  id: number;
  feature_a_id: number;
  feature_b_id: number;
  feature_a_title: string;
  feature_b_title: string;
  consensus_reached: boolean;
  vote_count: number;
  required_votes: number;
  my_vote?: {
    preferred_feature_id?: number;
    is_tie_vote: boolean;
  };
}

interface SessionProgress {
  session_id: number;
  total_comparisons: number;
  completed_comparisons: number;
  progress_percentage: number;
  is_complete: boolean;
}

interface UsePairwiseSessionOptions {
  projectId: number;
  sessionId: number;
  attendeeId: number;
  criterionType: 'value' | 'complexity';
}

export const usePairwiseSession = (options: UsePairwiseSessionOptions) => {
  const [comparisons, setComparisons] = useState<PairwiseComparison[]>([]);
  const [progress, setProgress] = useState<SessionProgress | null>(null);
  const [attendeesOnline, setAttendeesOnline] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  const wsUrl = `ws://localhost:8080/api/projects/${options.projectId}/sessions/${options.sessionId}/ws?attendee_id=${options.attendeeId}`;

  const handleWebSocketMessage = useCallback((message: any) => {
    switch (message.type) {
      case 'session_progress':
        setProgress(message.data);
        break;
        
      case 'vote_received':
        // Update comparison with new vote
        setComparisons(prev => prev.map(comp => 
          comp.id === message.data.comparison_id 
            ? { ...comp, vote_count: comp.vote_count + 1 }
            : comp
        ));
        break;
        
      case 'consensus_reached':
        // Mark comparison as consensus reached
        setComparisons(prev => prev.map(comp =>
          comp.id === message.data.comparison_id
            ? { ...comp, consensus_reached: true }
            : comp
        ));
        break;
        
      case 'attendee_joined':
        setAttendeesOnline(prev => [...prev, message.data.attendee_name]);
        break;
        
      case 'attendee_left':
        setAttendeesOnline(prev => prev.filter(name => name !== message.data.attendee_name));
        break;
    }
  }, []);

  const { connected, sendMessage } = useWebSocket(wsUrl, {
    onMessage: handleWebSocketMessage,
    onConnect: () => {
      sendMessage({
        type: 'join_session',
        session_id: options.sessionId,
        attendee_id: options.attendeeId,
      });
    },
  });

  const submitVote = useCallback(async (comparisonId: number, preferredFeatureId?: number, isTie: boolean = false) => {
    const voteData = {
      comparison_id: comparisonId,
      preferred_feature_id: preferredFeatureId,
      is_tie_vote: isTie,
    };

    // Optimistic update
    setComparisons(prev => prev.map(comp =>
      comp.id === comparisonId
        ? { 
            ...comp, 
            my_vote: { 
              preferred_feature_id: preferredFeatureId, 
              is_tie_vote: isTie 
            }
          }
        : comp
    ));

    // Send via WebSocket for real-time updates
    sendMessage({
      type: 'vote_submitted',
      data: { vote: voteData },
    });

    // Also send via HTTP for persistence
    try {
      const response = await fetch(`/api/projects/${options.projectId}/pairwise-sessions/${options.sessionId}/vote`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(voteData),
      });
      
      if (!response.ok) {
        throw new Error('Failed to submit vote');
      }
    } catch (error) {
      console.error('Failed to submit vote:', error);
      // Revert optimistic update on error
      setComparisons(prev => prev.map(comp =>
        comp.id === comparisonId
          ? { ...comp, my_vote: undefined }
          : comp
      ));
    }
  }, [options.projectId, options.sessionId, sendMessage]);

  // Fetch initial session data
  useEffect(() => {
    const fetchSessionData = async () => {
      try {
        setLoading(true);
        const [comparisonsRes, progressRes] = await Promise.all([
          fetch(`/api/projects/${options.projectId}/pairwise-sessions/${options.sessionId}/comparisons`),
          fetch(`/api/projects/${options.projectId}/pairwise-sessions/${options.sessionId}/progress`),
        ]);

        if (comparisonsRes.ok) {
          const comparisonsData = await comparisonsRes.json();
          setComparisons(comparisonsData);
        }

        if (progressRes.ok) {
          const progressData = await progressRes.json();
          setProgress(progressData);
        }
      } catch (error) {
        console.error('Failed to fetch session data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchSessionData();
  }, [options.projectId, options.sessionId]);

  return {
    comparisons,
    progress,
    attendeesOnline,
    connected,
    loading,
    submitVote,
  };
};
```

## Pairwise Comparison Grid Components

### Main Comparison Grid Layout
```typescript
// src/components/Pairwise/PairwiseGrid.tsx
import React from 'react';
import {
  Box,
  Grid,
  Typography,
  LinearProgress,
  Card,
  CardContent,
  Chip,
} from '@mui/material';
import { CheckCircle, Schedule } from '@mui/icons-material';
import { ComparisonCard } from './ComparisonCard';
import { SessionProgressCard } from './SessionProgressCard';

interface PairwiseGridProps {
  comparisons: PairwiseComparison[];
  progress: SessionProgress | null;
  attendeesOnline: string[];
  criterionType: 'value' | 'complexity';
  onVote: (comparisonId: number, preferredFeatureId?: number, isTie?: boolean) => void;
  connected: boolean;
}

export const PairwiseGrid: React.FC<PairwiseGridProps> = ({
  comparisons,
  progress,
  attendeesOnline,
  criterionType,
  onVote,
  connected,
}) => {
  const criterionLabel = criterionType === 'value' ? 'Business Value' : 'Complexity';
  const criterionDescription = criterionType === 'value' 
    ? 'Which feature provides more business value?'
    : 'Which feature is less complex to implement?';

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          Pairwise Comparison - {criterionLabel}
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
          {criterionDescription}
        </Typography>
        
        {/* Connection Status */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
          <Chip
            icon={connected ? <CheckCircle /> : <Schedule />}
            label={connected ? 'Connected' : 'Reconnecting...'}
            color={connected ? 'success' : 'warning'}
            size="small"
          />
          
          <Typography variant="body2" color="text.secondary">
            {attendeesOnline.length} attendees online
          </Typography>
        </Box>

        {/* Progress Bar */}
        {progress && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="body2" gutterBottom>
              Progress: {progress.completed_comparisons} / {progress.total_comparisons} comparisons
            </Typography>
            <LinearProgress 
              variant="determinate" 
              value={progress.progress_percentage} 
              sx={{ height: 8, borderRadius: 4 }}
            />
          </Box>
        )}
      </div>

      {/* Comparison Grid */}
      <Grid container spacing={3}>
        {/* Left Column - Comparisons */}
        <Grid item xs={12} lg={8}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {comparisons.map((comparison) => (
              <ComparisonCard
                key={comparison.id}
                comparison={comparison}
                criterionType={criterionType}
                onVote={onVote}
                disabled={!connected || comparison.consensus_reached}
              />
            ))}
          </Box>
        </Grid>

        {/* Right Column - Progress Sidebar */}
        <Grid item xs={12} lg={4}>
          <SessionProgressCard
            progress={progress}
            attendeesOnline={attendeesOnline}
            comparisons={comparisons}
          />
        </Grid>
      </Grid>
    </Box>
  );
};
```

### Individual Comparison Card
```typescript
// src/components/Pairwise/ComparisonCard.tsx
import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Box,
  Typography,
  Button,
  ButtonGroup,
  Chip,
  Collapse,
  IconButton,
  Divider,
} from '@mui/material';
import { 
  CheckCircle, 
  Schedule, 
  ExpandMore, 
  ExpandLess,
  Balance 
} from '@mui/icons-material';

interface ComparisonCardProps {
  comparison: PairwiseComparison;
  criterionType: 'value' | 'complexity';
  onVote: (comparisonId: number, preferredFeatureId?: number, isTie?: boolean) => void;
  disabled: boolean;
}

export const ComparisonCard: React.FC<ComparisonCardProps> = ({
  comparison,
  criterionType,
  onVote,
  disabled,
}) => {
  const [expanded, setExpanded] = useState(false);
  const [voting, setVoting] = useState(false);

  const handleVote = async (preferredFeatureId?: number, isTie: boolean = false) => {
    setVoting(true);
    try {
      await onVote(comparison.id, preferredFeatureId, isTie);
    } finally {
      setVoting(false);
    }
  };

  const getVoteButtonVariant = (featureId: number) => {
    if (comparison.my_vote?.preferred_feature_id === featureId) {
      return 'contained';
    }
    return 'outlined';
  };

  const getTieButtonVariant = () => {
    return comparison.my_vote?.is_tie_vote ? 'contained' : 'outlined';
  };

  return (
    <Card 
      sx={{ 
        opacity: comparison.consensus_reached ? 0.7 : 1,
        border: comparison.consensus_reached ? '2px solid' : 'none',
        borderColor: 'success.main',
      }}
    >
      <CardContent>
        {/* Header with Status */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" component="h3">
            Comparison #{comparison.id}
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Chip
              icon={comparison.consensus_reached ? <CheckCircle /> : <Schedule />}
              label={comparison.consensus_reached ? 'Consensus Reached' : 'Voting in Progress'}
              color={comparison.consensus_reached ? 'success' : 'warning'}
              size="small"
            />
            
            <Typography variant="body2" color="text.secondary">
              {comparison.vote_count}/{comparison.required_votes} votes
            </Typography>
          </Box>
        </Box>

        {/* Feature Comparison */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
          {/* Feature A */}
          <Box sx={{ flex: 1 }}>
            <Card variant="outlined" sx={{ p: 2, height: '100%' }}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {comparison.feature_a_title}
              </Typography>
              
              <Button
                variant={getVoteButtonVariant(comparison.feature_a_id)}
                color="primary"
                fullWidth
                disabled={disabled || voting}
                onClick={() => handleVote(comparison.feature_a_id)}
                sx={{ mt: 1 }}
              >
                {criterionType === 'value' ? 'More Valuable' : 'Less Complex'}
              </Button>
            </Card>
          </Box>

          {/* VS Divider */}
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mx: 2 }}>
            <Typography variant="h6" color="text.secondary">VS</Typography>
            
            <Button
              variant={getTieButtonVariant()}
              color="secondary"
              startIcon={<Balance />}
              disabled={disabled || voting}
              onClick={() => handleVote(undefined, true)}
              sx={{ mt: 1, minWidth: 80 }}
              size="small"
            >
              Tie
            </Button>
          </Box>

          {/* Feature B */}
          <Box sx={{ flex: 1 }}>
            <Card variant="outlined" sx={{ p: 2, height: '100%' }}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                {comparison.feature_b_title}
              </Typography>
              
              <Button
                variant={getVoteButtonVariant(comparison.feature_b_id)}
                color="primary"
                fullWidth
                disabled={disabled || voting}
                onClick={() => handleVote(comparison.feature_b_id)}
                sx={{ mt: 1 }}
              >
                {criterionType === 'value' ? 'More Valuable' : 'Less Complex'}
              </Button>
            </Card>
          </Box>
        </Box>

        {/* Expandable Details */}
        <Box>
          <Button
            onClick={() => setExpanded(!expanded)}
            endIcon={expanded ? <ExpandLess /> : <ExpandMore />}
            size="small"
          >
            {expanded ? 'Hide' : 'Show'} Feature Details
          </Button>
          
          <Collapse in={expanded}>
            <Box sx={{ mt: 2 }}>
              <Divider sx={{ mb: 2 }} />
              
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Typography variant="subtitle2" gutterBottom>
                    {comparison.feature_a_title} - Details
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {comparison.feature_a_description}
                  </Typography>
                </Grid>
                
                <Grid item xs={12} sm={6}>
                  <Typography variant="subtitle2" gutterBottom>
                    {comparison.feature_b_title} - Details
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {comparison.feature_b_description}
                  </Typography>
                </Grid>
              </Grid>
            </Box>
          </Collapse>
        </Box>
      </CardContent>
    </Card>
  );
};
```

## Session Progress Components

### Progress Sidebar Component
```typescript
// src/components/Pairwise/SessionProgressCard.tsx
import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Avatar,
  LinearProgress,
  Divider,
} from '@mui/material';
import { 
  CheckCircle, 
  Schedule, 
  Person, 
  PersonOutline 
} from '@mui/icons-material';

interface SessionProgressCardProps {
  progress: SessionProgress | null;
  attendeesOnline: string[];
  comparisons: PairwiseComparison[];
}

export const SessionProgressCard: React.FC<SessionProgressCardProps> = ({
  progress,
  attendeesOnline,
  comparisons,
}) => {
  const completedCount = comparisons.filter(c => c.consensus_reached).length;
  const pendingCount = comparisons.filter(c => !c.consensus_reached).length;

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      {/* Overall Progress */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Session Progress
          </Typography>
          
          {progress && (
            <Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                <Typography variant="body2">
                  {progress.completed_comparisons} of {progress.total_comparisons}
                </Typography>
                <Typography variant="body2">
                  {Math.round(progress.progress_percentage)}%
                </Typography>
              </Box>
              
              <LinearProgress
                variant="determinate"
                value={progress.progress_percentage}
                sx={{ height: 8, borderRadius: 4, mb: 2 }}
              />
              
              <Box sx={{ display: 'flex', gap: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <CheckCircle color="success" fontSize="small" />
                  <Typography variant="body2">
                    {completedCount} Complete
                  </Typography>
                </Box>
                
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Schedule color="warning" fontSize="small" />
                  <Typography variant="body2">
                    {pendingCount} Pending
                  </Typography>
                </Box>
              </Box>
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Attendees Online */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Attendees ({attendeesOnline.length} online)
          </Typography>
          
          <List dense>
            {attendeesOnline.map((attendeeName, index) => (
              <ListItem key={index}>
                <ListItemIcon>
                  <Avatar sx={{ width: 32, height: 32, bgcolor: 'success.main' }}>
                    <Person fontSize="small" />
                  </Avatar>
                </ListItemIcon>
                <ListItemText primary={attendeeName} />
              </ListItem>
            ))}
          </List>
        </CardContent>
      </Card>

      {/* Quick Navigation */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Quick Navigation
          </Typography>
          
          <List dense>
            {comparisons
              .filter(c => !c.consensus_reached)
              .slice(0, 5) // Show first 5 pending
              .map((comparison) => (
                <ListItem key={comparison.id} button>
                  <ListItemIcon>
                    <Schedule color="warning" fontSize="small" />
                  </ListItemIcon>
                  <ListItemText
                    primary={`${comparison.feature_a_title} vs ${comparison.feature_b_title}`}
                    secondary={`${comparison.vote_count}/${comparison.required_votes} votes`}
                  />
                </ListItem>
              ))}
          </List>
        </CardContent>
      </Card>
    </Box>
  );
};
```

## Keyboard Shortcuts and Accessibility

### Keyboard Navigation Hook
```typescript
// src/hooks/useKeyboardShortcuts.ts
import { useEffect } from 'react';

interface KeyboardShortcuts {
  onFeatureA?: () => void;
  onFeatureB?: () => void;
  onTie?: () => void;
  onNext?: () => void;
  onPrevious?: () => void;
}

export const useKeyboardShortcuts = (shortcuts: KeyboardShortcuts, enabled: boolean = true) => {
  useEffect(() => {
    if (!enabled) return;

    const handleKeyPress = (event: KeyboardEvent) => {
      // Only handle shortcuts when not in an input field
      if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
        return;
      }

      switch (event.key) {
        case '1':
        case 'a':
        case 'A':
          event.preventDefault();
          shortcuts.onFeatureA?.();
          break;
          
        case '2':
        case 'b':
        case 'B':
          event.preventDefault();
          shortcuts.onFeatureB?.();
          break;
          
        case '0':
        case 't':
        case 'T':
          event.preventDefault();
          shortcuts.onTie?.();
          break;
          
        case 'ArrowRight':
        case 'n':
        case 'N':
          event.preventDefault();
          shortcuts.onNext?.();
          break;
          
        case 'ArrowLeft':
        case 'p':
        case 'P':
          event.preventDefault();
          shortcuts.onPrevious?.();
          break;
      }
    };

    document.addEventListener('keydown', handleKeyPress);
    return () => document.removeEventListener('keydown', handleKeyPress);
  }, [shortcuts, enabled]);
};
```

## Error Handling and Fallbacks

### Error Recovery Component
```typescript
// src/components/Pairwise/PairwiseErrorBoundary.tsx
import React, { Component, ReactNode } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Button,
  Alert,
  Box,
} from '@mui/material';
import { Refresh, Warning } from '@mui/icons-material';

interface Props {
  children: ReactNode;
  onRetry?: () => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class PairwiseErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: any) {
    console.error('Pairwise comparison error:', error, errorInfo);
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null });
    this.props.onRetry?.();
  };

  render() {
    if (this.state.hasError) {
      return (
        <Card sx={{ maxWidth: 600, mx: 'auto', mt: 4 }}>
          <CardContent>
            <Box sx={{ textAlign: 'center' }}>
              <Warning color="error" sx={{ fontSize: 64, mb: 2 }} />
              
              <Typography variant="h5" gutterBottom>
                Something went wrong
              </Typography>
              
              <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                The pairwise comparison interface encountered an error. 
                Your votes have been saved and you can try refreshing the page.
              </Typography>
              
              <Alert severity="error" sx={{ mb: 3, textAlign: 'left' }}>
                <Typography variant="body2">
                  Error: {this.state.error?.message}
                </Typography>
              </Alert>
              
              <Button
                variant="contained"
                startIcon={<Refresh />}
                onClick={this.handleRetry}
                size="large"
              >
                Try Again
              </Button>
            </Box>
          </CardContent>
        </Card>
      );
    }

    return this.props.children;
  }
}
```