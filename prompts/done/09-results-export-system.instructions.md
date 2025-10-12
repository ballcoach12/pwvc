# Instructions: Results & Export System Implementation

## Final Priority Score Calculation

### Results Calculation Service Hook
```typescript
// src/hooks/useResultsCalculation.ts
import { useState, useEffect, useCallback } from 'react';

interface PWVCResult {
  feature_id: number;
  feature_title: string;
  feature_description: string;
  w_value: number;
  w_complexity: number;
  s_value: number;
  s_complexity: number;
  weighted_value: number;
  weighted_complexity: number;
  final_priority_score: number;
  rank: number;
}

interface CalculationBreakdown {
  feature_id: number;
  pairwise_value_wins: number;
  pairwise_value_total: number;
  pairwise_complexity_wins: number;
  pairwise_complexity_total: number;
  fibonacci_value_score: number;
  fibonacci_complexity_score: number;
  calculation_steps: {
    step: string;
    formula: string;
    result: number;
  }[];
}

interface UseResultsCalculationOptions {
  projectId: number;
}

export const useResultsCalculation = ({ projectId }: UseResultsCalculationOptions) => {
  const [results, setResults] = useState<PWVCResult[]>([]);
  const [breakdown, setBreakdown] = useState<CalculationBreakdown[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [calculating, setCalculating] = useState(false);

  const calculateResults = useCallback(async () => {
    setCalculating(true);
    setError(null);
    
    try {
      const response = await fetch(`/api/projects/${projectId}/calculate-results`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });

      if (!response.ok) {
        throw new Error('Failed to calculate results');
      }

      const data = await response.json();
      setResults(data.results || []);
      setBreakdown(data.breakdown || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Calculation failed');
    } finally {
      setCalculating(false);
    }
  }, [projectId]);

  const fetchResults = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/projects/${projectId}/results`);
      
      if (response.ok) {
        const data = await response.json();
        setResults(data.results || []);
        setBreakdown(data.breakdown || []);
      } else if (response.status === 404) {
        // No results calculated yet
        setResults([]);
        setBreakdown([]);
      } else {
        throw new Error('Failed to fetch results');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch results');
    } finally {
      setLoading(false);
    }
  }, [projectId]);

  const exportResults = useCallback(async (format: 'csv' | 'json' | 'jira') => {
    try {
      const response = await fetch(`/api/projects/${projectId}/results/export?format=${format}`);
      
      if (!response.ok) {
        throw new Error(`Failed to export as ${format.toUpperCase()}`);
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `pwvc-results-${projectId}.${format === 'jira' ? 'json' : format}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Export failed');
    }
  }, [projectId]);

  useEffect(() => {
    fetchResults();
  }, [fetchResults]);

  return {
    results,
    breakdown,
    loading,
    error,
    calculating,
    calculateResults,
    fetchResults,
    exportResults,
  };
};
```

## Results Visualization Components

### Main Results Dashboard
```typescript
// src/components/Results/ResultsDashboard.tsx
import React, { useState } from 'react';
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Grid,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
} from '@mui/material';
import {
  Calculate,
  Download,
  TableChart,
  BarChart,
  ScatterPlot,
} from '@mui/icons-material';
import { ResultsRanking } from './ResultsRanking';
import { ResultsVisualization } from './ResultsVisualization';
import { CalculationBreakdown } from './CalculationBreakdown';
import { ExportOptions } from './ExportOptions';

interface ResultsDashboardProps {
  projectId: number;
  results: PWVCResult[];
  breakdown: CalculationBreakdown[];
  loading: boolean;
  error: string | null;
  calculating: boolean;
  onCalculateResults: () => void;
  onExportResults: (format: 'csv' | 'json' | 'jira') => Promise<void>;
}

export const ResultsDashboard: React.FC<ResultsDashboardProps> = ({
  projectId,
  results,
  breakdown,
  loading,
  error,
  calculating,
  onCalculateResults,
  onExportResults,
}) => {
  const [activeTab, setActiveTab] = useState(0);
  const [exportError, setExportError] = useState<string | null>(null);

  const handleExport = async (format: 'csv' | 'json' | 'jira') => {
    setExportError(null);
    try {
      await onExportResults(format);
    } catch (err) {
      setExportError(err instanceof Error ? err.message : 'Export failed');
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 4 }}>
        {error}
      </Alert>
    );
  }

  if (results.length === 0) {
    return (
      <Card sx={{ textAlign: 'center', py: 8 }}>
        <CardContent>
          <Calculate sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
          
          <Typography variant="h5" gutterBottom>
            Ready to Calculate Results
          </Typography>
          
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            All pairwise comparisons and Fibonacci scoring must be completed 
            before calculating final priority scores.
          </Typography>
          
          <Button
            variant="contained"
            size="large"
            startIcon={<Calculate />}
            onClick={onCalculateResults}
            disabled={calculating}
          >
            {calculating ? 'Calculating...' : 'Calculate P-WVC Results'}
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            P-WVC Results
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Features ranked by Final Priority Score (Weighted Value ÷ Weighted Complexity)
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<Calculate />}
            onClick={onCalculateResults}
            disabled={calculating}
          >
            {calculating ? 'Recalculating...' : 'Recalculate'}
          </Button>
        </Box>
      </Box>

      {exportError && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setExportError(null)}>
          {exportError}
        </Alert>
      )}

      {/* Tabs for different views */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={activeTab} onChange={(_, newValue) => setActiveTab(newValue)}>
          <Tab icon={<TableChart />} label="Ranking Table" />
          <Tab icon={<BarChart />} label="Visualizations" />
          <Tab icon={<Calculate />} label="Calculation Details" />
          <Tab icon={<Download />} label="Export Options" />
        </Tabs>
      </Box>

      {/* Tab Content */}
      {activeTab === 0 && (
        <ResultsRanking results={results} />
      )}

      {activeTab === 1 && (
        <ResultsVisualization results={results} />
      )}

      {activeTab === 2 && (
        <CalculationBreakdown breakdown={breakdown} results={results} />
      )}

      {activeTab === 3 && (
        <ExportOptions onExport={handleExport} />
      )}
    </Box>
  );
};
```

### Results Ranking Table
```typescript
// src/components/Results/ResultsRanking.tsx
import React, { useState } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Collapse,
  Box,
  Typography,
  Chip,
  Card,
  TableSortLabel,
} from '@mui/material';
import { KeyboardArrowDown, KeyboardArrowUp, EmojiEvents } from '@mui/icons-material';

interface ResultsRankingProps {
  results: PWVCResult[];
}

export const ResultsRanking: React.FC<ResultsRankingProps> = ({ results }) => {
  const [expandedRow, setExpandedRow] = useState<number | null>(null);
  const [sortBy, setSortBy] = useState<'rank' | 'fps' | 'value' | 'complexity'>('rank');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');

  const handleSort = (column: 'rank' | 'fps' | 'value' | 'complexity') => {
    const isCurrentColumn = sortBy === column;
    const newDirection = isCurrentColumn && sortDirection === 'asc' ? 'desc' : 'asc';
    
    setSortBy(column);
    setSortDirection(newDirection);
  };

  const sortedResults = [...results].sort((a, b) => {
    let aValue: number;
    let bValue: number;

    switch (sortBy) {
      case 'rank':
        aValue = a.rank;
        bValue = b.rank;
        break;
      case 'fps':
        aValue = a.final_priority_score;
        bValue = b.final_priority_score;
        break;
      case 'value':
        aValue = a.weighted_value;
        bValue = b.weighted_value;
        break;
      case 'complexity':
        aValue = a.weighted_complexity;
        bValue = b.weighted_complexity;
        break;
    }

    return sortDirection === 'asc' ? aValue - bValue : bValue - aValue;
  });

  const getRankIcon = (rank: number) => {
    if (rank === 1) return <EmojiEvents sx={{ color: 'gold' }} />;
    if (rank === 2) return <EmojiEvents sx={{ color: 'silver' }} />;
    if (rank === 3) return <EmojiEvents sx={{ color: '#CD7F32' }} />;
    return null;
  };

  const getPriorityColor = (score: number, maxScore: number) => {
    const ratio = score / maxScore;
    if (ratio > 0.8) return 'success';
    if (ratio > 0.6) return 'warning';
    return 'default';
  };

  const maxScore = Math.max(...results.map(r => r.final_priority_score));

  return (
    <Card>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell>
                <TableSortLabel
                  active={sortBy === 'rank'}
                  direction={sortBy === 'rank' ? sortDirection : 'asc'}
                  onClick={() => handleSort('rank')}
                >
                  Rank
                </TableSortLabel>
              </TableCell>
              <TableCell>Feature</TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={sortBy === 'fps'}
                  direction={sortBy === 'fps' ? sortDirection : 'asc'}
                  onClick={() => handleSort('fps')}
                >
                  Final Priority Score
                </TableSortLabel>
              </TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={sortBy === 'value'}
                  direction={sortBy === 'value' ? sortDirection : 'asc'}
                  onClick={() => handleSort('value')}
                >
                  Weighted Value
                </TableSortLabel>
              </TableCell>
              <TableCell align="right">
                <TableSortLabel
                  active={sortBy === 'complexity'}
                  direction={sortBy === 'complexity' ? sortDirection : 'asc'}
                  onClick={() => handleSort('complexity')}
                >
                  Weighted Complexity
                </TableSortLabel>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sortedResults.map((result) => (
              <React.Fragment key={result.feature_id}>
                <TableRow hover>
                  <TableCell>
                    <IconButton
                      size="small"
                      onClick={() => setExpandedRow(
                        expandedRow === result.feature_id ? null : result.feature_id
                      )}
                    >
                      {expandedRow === result.feature_id ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                    </IconButton>
                  </TableCell>
                  
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      {getRankIcon(result.rank)}
                      <Typography variant="h6" component="span">
                        #{result.rank}
                      </Typography>
                    </Box>
                  </TableCell>
                  
                  <TableCell>
                    <Box>
                      <Typography variant="subtitle1" fontWeight="medium">
                        {result.feature_title}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap>
                        {result.feature_description.substring(0, 80)}...
                      </Typography>
                    </Box>
                  </TableCell>
                  
                  <TableCell align="right">
                    <Chip
                      label={result.final_priority_score.toFixed(2)}
                      color={getPriorityColor(result.final_priority_score, maxScore)}
                      variant="filled"
                      sx={{ fontWeight: 'bold', minWidth: 80 }}
                    />
                  </TableCell>
                  
                  <TableCell align="right">
                    <Typography variant="body2">
                      {result.weighted_value.toFixed(2)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      ({result.s_value} × {result.w_value.toFixed(2)})
                    </Typography>
                  </TableCell>
                  
                  <TableCell align="right">
                    <Typography variant="body2">
                      {result.weighted_complexity.toFixed(2)}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      ({result.s_complexity} × {result.w_complexity.toFixed(2)})
                    </Typography>
                  </TableCell>
                </TableRow>

                <TableRow>
                  <TableCell colSpan={6} sx={{ py: 0 }}>
                    <Collapse in={expandedRow === result.feature_id} timeout="auto" unmountOnExit>
                      <Box sx={{ p: 3, bgcolor: 'grey.50' }}>
                        <Typography variant="h6" gutterBottom>
                          Calculation Breakdown
                        </Typography>
                        
                        <Grid container spacing={3}>
                          <Grid item xs={12} md={6}>
                            <Typography variant="subtitle2" gutterBottom>
                              Pairwise Comparison Results:
                            </Typography>
                            <Typography variant="body2">
                              • Value Win-Count: {result.w_value.toFixed(3)}
                            </Typography>
                            <Typography variant="body2">
                              • Complexity Win-Count: {result.w_complexity.toFixed(3)}
                            </Typography>
                          </Grid>
                          
                          <Grid item xs={12} md={6}>
                            <Typography variant="subtitle2" gutterBottom>
                              Fibonacci Scores:
                            </Typography>
                            <Typography variant="body2">
                              • Value Score: {result.s_value}
                            </Typography>
                            <Typography variant="body2">
                              • Complexity Score: {result.s_complexity}
                            </Typography>
                          </Grid>
                          
                          <Grid item xs={12}>
                            <Typography variant="subtitle2" gutterBottom>
                              Final Calculation:
                            </Typography>
                            <Typography variant="body2" sx={{ fontFamily: 'monospace', bgcolor: 'white', p: 1, borderRadius: 1 }}>
                              FPS = ({result.s_value} × {result.w_value.toFixed(3)}) ÷ ({result.s_complexity} × {result.w_complexity.toFixed(3)}) 
                              = {result.weighted_value.toFixed(3)} ÷ {result.weighted_complexity.toFixed(3)} 
                              = <strong>{result.final_priority_score.toFixed(3)}</strong>
                            </Typography>
                          </Grid>
                        </Grid>
                      </Box>
                    </Collapse>
                  </TableCell>
                </TableRow>
              </React.Fragment>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Card>
  );
};
```

## Export System Components

### Export Options Panel
```typescript
// src/components/Results/ExportOptions.tsx
import React, { useState } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  Box,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  TableChart,
  Code,
  BugReport,
  Download,
} from '@mui/icons-material';

interface ExportOptionsProps {
  onExport: (format: 'csv' | 'json' | 'jira') => Promise<void>;
}

export const ExportOptions: React.FC<ExportOptionsProps> = ({ onExport }) => {
  const [exporting, setExporting] = useState<string | null>(null);
  const [exportSuccess, setExportSuccess] = useState<string | null>(null);

  const handleExport = async (format: 'csv' | 'json' | 'jira') => {
    setExporting(format);
    setExportSuccess(null);
    
    try {
      await onExport(format);
      setExportSuccess(`Successfully exported as ${format.toUpperCase()}`);
    } catch (error) {
      // Error is handled by parent component
    } finally {
      setExporting(null);
    }
  };

  const exportOptions = [
    {
      format: 'csv' as const,
      title: 'CSV Export',
      description: 'Export for spreadsheet analysis in Excel, Google Sheets, etc.',
      icon: <TableChart sx={{ fontSize: 48 }} />,
      features: [
        'All ranking data and scores',
        'Calculation breakdown',
        'Ready for pivot tables',
        'Import into other tools',
      ],
    },
    {
      format: 'json' as const,
      title: 'JSON Export',
      description: 'Export structured data for API integration or custom analysis',
      icon: <Code sx={{ fontSize: 48 }} />,
      features: [
        'Machine-readable format',
        'Complete calculation details',
        'API integration ready',
        'Custom processing',
      ],
    },
    {
      format: 'jira' as const,
      title: 'Jira Import Format',
      description: 'Export stories with Fibonacci complexity scores for Jira import',
      icon: <BugReport sx={{ fontSize: 48 }} />,
      features: [
        'Story format with descriptions',
        'Fibonacci story points assigned',
        'Priority labels included',
        'Ready for Jira import',
      ],
    },
  ];

  return (
    <Box>
      {exportSuccess && (
        <Alert severity="success" sx={{ mb: 3 }} onClose={() => setExportSuccess(null)}>
          {exportSuccess}
        </Alert>
      )}

      <Grid container spacing={3}>
        {exportOptions.map((option) => (
          <Grid item xs={12} md={4} key={option.format}>
            <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <CardContent sx={{ flex: 1, textAlign: 'center' }}>
                <Box sx={{ color: 'primary.main', mb: 2 }}>
                  {option.icon}
                </Box>
                
                <Typography variant="h5" gutterBottom>
                  {option.title}
                </Typography>
                
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                  {option.description}
                </Typography>
                
                <Box sx={{ textAlign: 'left', mb: 3 }}>
                  <Typography variant="subtitle2" gutterBottom>
                    Includes:
                  </Typography>
                  <ul style={{ margin: 0, paddingLeft: '1.5rem' }}>
                    {option.features.map((feature, index) => (
                      <li key={index}>
                        <Typography variant="body2">{feature}</Typography>
                      </li>
                    ))}
                  </ul>
                </Box>
              </CardContent>
              
              <Box sx={{ p: 2 }}>
                <Button
                  variant="contained"
                  fullWidth
                  startIcon={
                    exporting === option.format ? (
                      <CircularProgress size={20} color="inherit" />
                    ) : (
                      <Download />
                    )
                  }
                  onClick={() => handleExport(option.format)}
                  disabled={!!exporting}
                >
                  {exporting === option.format ? 'Exporting...' : `Export ${option.format.toUpperCase()}`}
                </Button>
              </Box>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Card sx={{ mt: 4 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Export Information
          </Typography>
          
          <Typography variant="body2" paragraph>
            All exports include the complete P-WVC calculation results with final priority scores, 
            ranking, and detailed breakdown of how each score was calculated.
          </Typography>
          
          <Typography variant="body2" paragraph>
            <strong>CSV Format:</strong> Perfect for further analysis in spreadsheet applications. 
            Includes all numerical data with proper headers for easy filtering and sorting.
          </Typography>
          
          <Typography variant="body2" paragraph>
            <strong>JSON Format:</strong> Structured data format ideal for developers and API integrations. 
            Contains complete metadata and calculation details.
          </Typography>
          
          <Typography variant="body2" paragraph>
            <strong>Jira Format:</strong> Stories are formatted with titles, descriptions, acceptance criteria, 
            and Fibonacci complexity scores assigned based on the P-WVC results. Import directly into Jira 
            for sprint planning.
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
};
```

## Results Visualization Components

### Charts and Data Visualization
```typescript
// src/components/Results/ResultsVisualization.tsx
import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
} from '@mui/material';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ScatterChart,
  Scatter,
  Cell,
} from 'recharts';

interface ResultsVisualizationProps {
  results: PWVCResult[];
}

export const ResultsVisualization: React.FC<ResultsVisualizationProps> = ({ results }) => {
  const [chartType, setChartType] = useState<'bar' | 'scatter' | 'quadrant'>('bar');

  const chartData = results.map((result, index) => ({
    name: result.feature_title.length > 15 
      ? `${result.feature_title.substring(0, 15)}...` 
      : result.feature_title,
    fullName: result.feature_title,
    fps: result.final_priority_score,
    weightedValue: result.weighted_value,
    weightedComplexity: result.weighted_complexity,
    rank: result.rank,
    color: index < 3 ? '#4caf50' : index < 7 ? '#ff9800' : '#f44336',
  }));

  const renderBarChart = () => (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 60 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis 
          dataKey="name" 
          angle={-45}
          textAnchor="end"
          height={100}
          interval={0}
        />
        <YAxis label={{ value: 'Final Priority Score', angle: -90, position: 'insideLeft' }} />
        <Tooltip 
          labelFormatter={(value, payload) => payload?.[0]?.payload?.fullName}
          formatter={(value: number) => [value.toFixed(2), 'Priority Score']}
        />
        <Bar dataKey="fps" name="Final Priority Score">
          {chartData.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );

  const renderScatterPlot = () => (
    <ResponsiveContainer width="100%" height={400}>
      <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: 40 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis 
          dataKey="weightedComplexity" 
          name="Weighted Complexity"
          label={{ value: 'Weighted Complexity', position: 'insideBottom', offset: -10 }}
        />
        <YAxis 
          dataKey="weightedValue" 
          name="Weighted Value"
          label={{ value: 'Weighted Value', angle: -90, position: 'insideLeft' }}
        />
        <Tooltip 
          cursor={{ strokeDasharray: '3 3' }}
          content={({ active, payload }) => {
            if (active && payload && payload.length) {
              const data = payload[0].payload;
              return (
                <Box sx={{ bgcolor: 'white', p: 1, border: 1, borderColor: 'grey.300', borderRadius: 1 }}>
                  <Typography variant="subtitle2">{data.fullName}</Typography>
                  <Typography variant="body2">
                    Weighted Value: {data.weightedValue.toFixed(2)}
                  </Typography>
                  <Typography variant="body2">
                    Weighted Complexity: {data.weightedComplexity.toFixed(2)}
                  </Typography>
                  <Typography variant="body2">
                    Final Priority Score: {data.fps.toFixed(2)}
                  </Typography>
                </Box>
              );
            }
            return null;
          }}
        />
        <Scatter data={chartData} fill="#8884d8">
          {chartData.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color} />
          ))}
        </Scatter>
      </ScatterChart>
    </ResponsiveContainer>
  );

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
              <Typography variant="h6">
                Results Visualization
              </Typography>
              
              <FormControl size="small" sx={{ minWidth: 150 }}>
                <InputLabel>Chart Type</InputLabel>
                <Select
                  value={chartType}
                  label="Chart Type"
                  onChange={(e) => setChartType(e.target.value as 'bar' | 'scatter' | 'quadrant')}
                >
                  <MenuItem value="bar">Priority Ranking</MenuItem>
                  <MenuItem value="scatter">Value vs Complexity</MenuItem>
                </Select>
              </FormControl>
            </Box>

            {chartType === 'bar' && renderBarChart()}
            {chartType === 'scatter' && renderScatterPlot()}
          </CardContent>
        </Card>
      </Grid>

      <Grid item xs={12} md={6}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Top Priorities
            </Typography>
            {results.slice(0, 5).map((result, index) => (
              <Box key={result.feature_id} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <Box
                  sx={{
                    width: 24,
                    height: 24,
                    borderRadius: '50%',
                    bgcolor: index < 3 ? 'success.main' : 'warning.main',
                    color: 'white',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    mr: 2,
                    fontSize: '0.75rem',
                    fontWeight: 'bold',
                  }}
                >
                  {result.rank}
                </Box>
                <Box sx={{ flex: 1 }}>
                  <Typography variant="body2" fontWeight="medium">
                    {result.feature_title}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    Score: {result.final_priority_score.toFixed(2)}
                  </Typography>
                </Box>
              </Box>
            ))}
          </CardContent>
        </Card>
      </Grid>

      <Grid item xs={12} md={6}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Score Distribution
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box sx={{ width: 16, height: 16, bgcolor: '#4caf50', borderRadius: 1 }} />
                <Typography variant="body2">
                  High Priority (Top 3): {results.filter((_, i) => i < 3).length} features
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box sx={{ width: 16, height: 16, bgcolor: '#ff9800', borderRadius: 1 }} />
                <Typography variant="body2">
                  Medium Priority (4-7): {results.filter((_, i) => i >= 3 && i < 7).length} features
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box sx={{ width: 16, height: 16, bgcolor: '#f44336', borderRadius: 1 }} />
                <Typography variant="body2">
                  Lower Priority (8+): {results.filter((_, i) => i >= 7).length} features
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
};
```