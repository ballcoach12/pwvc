# Instructions: Fibonacci Scoring Interface Implementation

## Fibonacci Score Validation and State Management

### Fibonacci Validation Utilities
```typescript
// src/utils/fibonacci.ts
export const VALID_FIBONACCI_SCORES = [1, 2, 3, 5, 8, 13, 21, 34, 55, 89] as const;
export type FibonacciScore = typeof VALID_FIBONACCI_SCORES[number];

export const isValidFibonacciScore = (value: number): value is FibonacciScore => {
  return VALID_FIBONACCI_SCORES.includes(value as FibonacciScore);
};

export const getFibonacciScoreDescription = (score: FibonacciScore, criterion: 'value' | 'complexity'): string => {
  const descriptions = {
    value: {
      1: 'Minimal business value',
      2: 'Very low business value', 
      3: 'Low business value',
      5: 'Moderate business value',
      8: 'Good business value',
      13: 'High business value',
      21: 'Very high business value',
      34: 'Exceptional business value',
      55: 'Game-changing business value',
      89: 'Revolutionary business value',
    },
    complexity: {
      1: 'Trivial - can be done in hours',
      2: 'Very simple - 1-2 days',
      3: 'Simple - less than a week',
      5: 'Moderate - 1-2 weeks',
      8: 'Complex - 2-4 weeks',
      13: 'Very complex - 1-2 months',
      21: 'Highly complex - 2-3 months',
      34: 'Extremely complex - 3-6 months',
      55: 'Massive undertaking - 6+ months',
      89: 'Multi-year epic',
    },
  };

  return descriptions[criterion][score];
};

export const getFibonacciColor = (score: FibonacciScore): string => {
  if (score <= 3) return '#4caf50'; // Green
  if (score <= 8) return '#ff9800'; // Orange
  if (score <= 21) return '#f44336'; // Red
  return '#9c27b0'; // Purple
};
```

### Fibonacci Session Hook
```typescript
// src/hooks/useFibonacciSession.ts
import { useState, useEffect, useCallback } from 'react';
import { useWebSocket } from './useWebSocket';

interface FibonacciScore {
  id: number;
  feature_id: number;
  attendee_id: number;
  score_value: number;
  scored_at: string;
}

interface FeatureScoring {
  feature_id: number;
  feature_title: string;
  feature_description: string;
  individual_scores: FibonacciScore[];
  consensus_score?: number;
  consensus_reached: boolean;
  my_score?: number;
}

interface SessionProgress {
  session_id: number;
  total_features: number;
  features_with_consensus: number;
  progress_percentage: number;
  is_complete: boolean;
}

interface UseFibonacciSessionOptions {
  projectId: number;
  sessionId: number;
  attendeeId: number;
  criterionType: 'value' | 'complexity';
}

export const useFibonacciSession = (options: UseFibonacciSessionOptions) => {
  const [featureScores, setFeatureScores] = useState<FeatureScoring[]>([]);
  const [progress, setProgress] = useState<SessionProgress | null>(null);
  const [attendeesOnline, setAttendeesOnline] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  const wsUrl = `ws://localhost:8080/api/projects/${options.projectId}/fibonacci-sessions/${options.sessionId}/ws?attendee_id=${options.attendeeId}`;

  const handleWebSocketMessage = useCallback((message: any) => {
    switch (message.type) {
      case 'score_submitted':
        // Update feature with new individual score
        setFeatureScores(prev => prev.map(feature => {
          if (feature.feature_id === message.data.feature_id) {
            const updatedScores = [...feature.individual_scores];
            const existingScoreIndex = updatedScores.findIndex(
              s => s.attendee_id === message.data.attendee_id
            );
            
            if (existingScoreIndex >= 0) {
              updatedScores[existingScoreIndex] = message.data.score;
            } else {
              updatedScores.push(message.data.score);
            }
            
            return { ...feature, individual_scores: updatedScores };
          }
          return feature;
        }));
        break;
        
      case 'consensus_reached':
        // Update feature with consensus score
        setFeatureScores(prev => prev.map(feature =>
          feature.feature_id === message.data.feature_id
            ? { 
                ...feature, 
                consensus_score: message.data.consensus_score,
                consensus_reached: true 
              }
            : feature
        ));
        break;
        
      case 'session_progress':
        setProgress(message.data);
        break;
    }
  }, []);

  const { connected, sendMessage } = useWebSocket(wsUrl, {
    onMessage: handleWebSocketMessage,
  });

  const submitScore = useCallback(async (featureId: number, score: FibonacciScore) => {
    if (!isValidFibonacciScore(score)) {
      throw new Error(`Invalid Fibonacci score: ${score}`);
    }

    // Optimistic update
    setFeatureScores(prev => prev.map(feature =>
      feature.feature_id === featureId
        ? { ...feature, my_score: score }
        : feature
    ));

    // Send via WebSocket for real-time updates
    sendMessage({
      type: 'score_submitted',
      data: {
        feature_id: featureId,
        score_value: score,
        attendee_id: options.attendeeId,
      },
    });

    // Persist via HTTP
    try {
      const response = await fetch(`/api/projects/${options.projectId}/fibonacci-sessions/${options.sessionId}/scores`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          feature_id: featureId,
          score_value: score,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to submit score');
      }
    } catch (error) {
      console.error('Failed to submit score:', error);
      // Revert optimistic update
      setFeatureScores(prev => prev.map(feature =>
        feature.feature_id === featureId
          ? { ...feature, my_score: undefined }
          : feature
      ));
      throw error;
    }
  }, [options.projectId, options.sessionId, options.attendeeId, sendMessage]);

  // Fetch initial session data
  useEffect(() => {
    const fetchSessionData = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/api/projects/${options.projectId}/fibonacci-sessions/${options.sessionId}/scores`);
        
        if (response.ok) {
          const data = await response.json();
          setFeatureScores(data.features || []);
          setProgress(data.progress || null);
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
    featureScores,
    progress,
    attendeesOnline,
    connected,
    loading,
    submitScore,
  };
};
```

## Fibonacci Scoring Components

### Main Scoring Grid
```typescript
// src/components/Fibonacci/FibonacciScoringGrid.tsx
import React from 'react';
import {
  Box,
  Typography,
  Grid,
  LinearProgress,
  Card,
  CardContent,
  Chip,
} from '@mui/material';
import { CheckCircle, Schedule } from '@mui/icons-material';
import { FeatureScoringCard } from './FeatureScoringCard';
import { SessionProgressCard } from './SessionProgressCard';
import { FibonacciLegend } from './FibonacciLegend';

interface FibonacciScoringGridProps {
  featureScores: FeatureScoring[];
  progress: SessionProgress | null;
  criterionType: 'value' | 'complexity';
  onScoreSubmit: (featureId: number, score: FibonacciScore) => Promise<void>;
  connected: boolean;
}

export const FibonacciScoringGrid: React.FC<FibonacciScoringGridProps> = ({
  featureScores,
  progress,
  criterionType,
  onScoreSubmit,
  connected,
}) => {
  const criterionLabel = criterionType === 'value' ? 'Business Value' : 'Complexity';
  
  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          Fibonacci Scoring - {criterionLabel}
        </Typography>
        
        <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
          Rate each feature using the Fibonacci sequence. All team members must agree on the final score.
        </Typography>

        {/* Connection Status */}
        <Chip
          icon={connected ? <CheckCircle /> : <Schedule />}
          label={connected ? 'Connected' : 'Reconnecting...'}
          color={connected ? 'success' : 'warning'}
          size="small"
          sx={{ mb: 2 }}
        />

        {/* Progress */}
        {progress && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="body2" gutterBottom>
              Progress: {progress.features_with_consensus} / {progress.total_features} features have consensus
            </Typography>
            <LinearProgress
              variant="determinate"
              value={progress.progress_percentage}
              sx={{ height: 8, borderRadius: 4 }}
            />
          </Box>
        )}
      </Box>

      <Grid container spacing={3}>
        {/* Left Column - Feature Scoring */}
        <Grid item xs={12} lg={8}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            {featureScores.map((featureScoring) => (
              <FeatureScoringCard
                key={featureScoring.feature_id}
                featureScoring={featureScoring}
                criterionType={criterionType}
                onScoreSubmit={onScoreSubmit}
                disabled={!connected || featureScoring.consensus_reached}
              />
            ))}
          </Box>
        </Grid>

        {/* Right Column - Progress and Legend */}
        <Grid item xs={12} lg={4}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <SessionProgressCard
              progress={progress}
              featureScores={featureScores}
            />
            
            <FibonacciLegend criterionType={criterionType} />
          </Box>
        </Grid>
      </Grid>
    </Box>
  );
};
```

### Feature Scoring Card
```typescript
// src/components/Fibonacci/FeatureScoringCard.tsx
import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  Button,
  ButtonGroup,
  Collapse,
  IconButton,
  Alert,
  Avatar,
  AvatarGroup,
  Tooltip,
} from '@mui/material';
import {
  CheckCircle,
  Schedule,
  ExpandMore,
  ExpandLess,
  Person,
} from '@mui/icons-material';
import { FibonacciScorePicker } from './FibonacciScorePicker';
import { VALID_FIBONACCI_SCORES, getFibonacciColor } from '@/utils/fibonacci';

interface FeatureScoringCardProps {
  featureScoring: FeatureScoring;
  criterionType: 'value' | 'complexity';
  onScoreSubmit: (featureId: number, score: FibonacciScore) => Promise<void>;
  disabled: boolean;
}

export const FeatureScoringCard: React.FC<FeatureScoringCardProps> = ({
  featureScoring,
  criterionType,
  onScoreSubmit,
  disabled,
}) => {
  const [expanded, setExpanded] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleScoreSubmit = async (score: FibonacciScore) => {
    setSubmitting(true);
    setError(null);
    
    try {
      await onScoreSubmit(featureScoring.feature_id, score);
    } catch (err) {
      setError('Failed to submit score. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  const getScoreFrequency = () => {
    const frequency: Record<number, number> = {};
    featureScoring.individual_scores.forEach(score => {
      frequency[score.score_value] = (frequency[score.score_value] || 0) + 1;
    });
    return frequency;
  };

  const scoreFrequency = getScoreFrequency();
  const hasConsensus = featureScoring.consensus_reached;
  const totalScores = featureScoring.individual_scores.length;

  return (
    <Card 
      sx={{ 
        opacity: hasConsensus ? 0.9 : 1,
        border: hasConsensus ? '2px solid' : 'none',
        borderColor: 'success.main',
      }}
    >
      <CardContent>
        {/* Header */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Box sx={{ flex: 1 }}>
            <Typography variant="h6" gutterBottom>
              {featureScoring.feature_title}
            </Typography>
            
            <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
              {featureScoring.feature_description}
            </Typography>
          </Box>

          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 1 }}>
            <Chip
              icon={hasConsensus ? <CheckCircle /> : <Schedule />}
              label={hasConsensus ? 'Consensus Reached' : 'Voting in Progress'}
              color={hasConsensus ? 'success' : 'warning'}
              size="small"
            />
            
            {hasConsensus && (
              <Chip
                label={`Final Score: ${featureScoring.consensus_score}`}
                sx={{ 
                  bgcolor: getFibonacciColor(featureScoring.consensus_score!),
                  color: 'white',
                  fontWeight: 'bold',
                }}
              />
            )}
          </Box>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {/* Current Scores Display */}
        {totalScores > 0 && (
          <Box sx={{ mb: 3 }}>
            <Typography variant="subtitle2" gutterBottom>
              Current Scores ({totalScores} submitted):
            </Typography>
            
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
              {VALID_FIBONACCI_SCORES.map(score => {
                const count = scoreFrequency[score] || 0;
                if (count === 0) return null;
                
                return (
                  <Chip
                    key={score}
                    label={`${score} (${count})`}
                    size="small"
                    sx={{
                      bgcolor: getFibonacciColor(score),
                      color: 'white',
                      opacity: count > 0 ? 1 : 0.3,
                    }}
                  />
                );
              })}
            </Box>
            
            {/* Attendee avatars showing who has scored */}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Typography variant="body2" color="text.secondary">
                Voted:
              </Typography>
              <AvatarGroup max={6}>
                {featureScoring.individual_scores.map((score, index) => (
                  <Tooltip key={index} title={`Score: ${score.score_value}`}>
                    <Avatar sx={{ width: 32, height: 32, bgcolor: getFibonacciColor(score.score_value) }}>
                      <Person fontSize="small" />
                    </Avatar>
                  </Tooltip>
                ))}
              </AvatarGroup>
            </Box>
          </Box>
        )}

        {/* My Score Input */}
        {!hasConsensus && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Your Score:
            </Typography>
            
            <FibonacciScorePicker
              selectedScore={featureScoring.my_score}
              criterionType={criterionType}
              onScoreSelect={handleScoreSubmit}
              disabled={disabled || submitting}
            />
          </Box>
        )}

        {/* Expandable Details */}
        <Box>
          <Button
            onClick={() => setExpanded(!expanded)}
            endIcon={expanded ? <ExpandLess /> : <ExpandMore />}
            size="small"
            disabled={submitting}
          >
            {expanded ? 'Hide' : 'Show'} Details
          </Button>
          
          <Collapse in={expanded}>
            <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
              <Typography variant="subtitle2" gutterBottom>
                Feature Description:
              </Typography>
              <Typography variant="body2" paragraph>
                {featureScoring.feature_description}
              </Typography>
              
              {/* Show individual scores with timestamps */}
              {featureScoring.individual_scores.length > 0 && (
                <Box>
                  <Typography variant="subtitle2" gutterBottom>
                    Individual Scores:
                  </Typography>
                  {featureScoring.individual_scores.map((score, index) => (
                    <Box key={index} sx={{ display: 'flex', justifyContent: 'space-between', py: 0.5 }}>
                      <Typography variant="body2">
                        Attendee {score.attendee_id}
                      </Typography>
                      <Chip
                        label={score.score_value}
                        size="small"
                        sx={{ bgcolor: getFibonacciColor(score.score_value), color: 'white' }}
                      />
                    </Box>
                  ))}
                </Box>
              )}
            </Box>
          </Collapse>
        </Box>
      </CardContent>
    </Card>
  );
};
```

### Fibonacci Score Picker Component
```typescript
// src/components/Fibonacci/FibonacciScorePicker.tsx
import React from 'react';
import {
  Box,
  Button,
  ButtonGroup,
  Tooltip,
  Typography,
} from '@mui/material';
import { VALID_FIBONACCI_SCORES, getFibonacciScoreDescription, getFibonacciColor } from '@/utils/fibonacci';

interface FibonacciScorePickerProps {
  selectedScore?: number;
  criterionType: 'value' | 'complexity';
  onScoreSelect: (score: FibonacciScore) => void;
  disabled?: boolean;
}

export const FibonacciScorePicker: React.FC<FibonacciScorePickerProps> = ({
  selectedScore,
  criterionType,
  onScoreSelect,
  disabled = false,
}) => {
  return (
    <Box>
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
        {VALID_FIBONACCI_SCORES.map((score) => (
          <Tooltip
            key={score}
            title={getFibonacciScoreDescription(score, criterionType)}
            placement="top"
          >
            <Button
              variant={selectedScore === score ? 'contained' : 'outlined'}
              onClick={() => onScoreSelect(score)}
              disabled={disabled}
              sx={{
                minWidth: 48,
                height: 48,
                borderRadius: '50%',
                bgcolor: selectedScore === score ? getFibonacciColor(score) : 'transparent',
                borderColor: getFibonacciColor(score),
                color: selectedScore === score ? 'white' : getFibonacciColor(score),
                fontWeight: 'bold',
                '&:hover': {
                  bgcolor: getFibonacciColor(score),
                  color: 'white',
                },
              }}
            >
              {score}
            </Button>
          </Tooltip>
        ))}
      </Box>
      
      {selectedScore && (
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
          Selected: {getFibonacciScoreDescription(selectedScore as FibonacciScore, criterionType)}
        </Typography>
      )}
    </Box>
  );
};
```

### Fibonacci Legend Component
```typescript
// src/components/Fibonacci/FibonacciLegend.tsx
import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
} from '@mui/material';
import { VALID_FIBONACCI_SCORES, getFibonacciScoreDescription, getFibonacciColor } from '@/utils/fibonacci';

interface FibonacciLegendProps {
  criterionType: 'value' | 'complexity';
}

export const FibonacciLegend: React.FC<FibonacciLegendProps> = ({
  criterionType,
}) => {
  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Fibonacci Scale Guide
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {criterionType === 'value' 
            ? 'Rate the business value this feature would provide:'
            : 'Rate how complex this feature would be to implement:'
          }
        </Typography>
        
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
          {VALID_FIBONACCI_SCORES.map((score) => (
            <Box key={score} sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Chip
                label={score}
                size="small"
                sx={{
                  minWidth: 32,
                  bgcolor: getFibonacciColor(score),
                  color: 'white',
                  fontWeight: 'bold',
                }}
              />
              <Typography variant="body2" sx={{ fontSize: '0.75rem' }}>
                {getFibonacciScoreDescription(score, criterionType)}
              </Typography>
            </Box>
          ))}
        </Box>
        
        <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
          <Typography variant="body2" sx={{ fontSize: '0.75rem', fontStyle: 'italic' }}>
            💡 Tip: Use the exponential nature of Fibonacci numbers to reflect that complexity 
            and value don't scale linearly. A score of 21 is much more than twice a score of 8.
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};
```