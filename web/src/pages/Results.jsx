import {
    Assessment as AssessmentIcon,
    Calculate as CalculateIcon,
    CloudDownload as CloudDownloadIcon,
    Info as InfoIcon,
    PlayArrow as PlayArrowIcon,
    Refresh as RefreshIcon,
    Visibility as VisibilityIcon,
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Container,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Grid,
    Paper,
    Skeleton,
    Tab,
    Tabs,
    Typography,
    useTheme
} from '@mui/material'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import CalculationDetails from '../components/CalculationDetails'
import ExportOptions from '../components/ExportOptions'
import ResultsRanking from '../components/ResultsRanking'
import ResultsVisualization from '../components/ResultsVisualization'
import { useProject } from '../hooks/useProject'

// Mock results service - in real implementation, this would call the API
const resultsService = {
  async checkStatus(projectId) {
    // Mock implementation - replace with actual API call
    return {
      hasResults: true,
      lastCalculated: new Date().toISOString(),
      calculationUrl: `/api/projects/${projectId}/calculate-results`,
      resultsUrl: `/api/projects/${projectId}/results`
    }
  },

  async getResults(projectId) {
    // Mock implementation - replace with actual API call
    const mockResults = {
      projectId: parseInt(projectId),
      calculatedAt: new Date().toISOString(),
      totalFeatures: 8,
      summary: {
        highestScore: 6.250,
        lowestScore: 0.833,
        averageScore: 3.542,
        medianScore: 3.500,
        scoreRange: 5.417,
        topTier: 2,
        bottomTier: 2
      },
      results: [
        {
          rank: 1,
          featureId: 1,
          finalPriorityScore: 6.250,
          sValue: 8,
          sComplexity: 2,
          wValue: 0.781,
          wComplexity: 0.800,
          weightedValue: 6.248,
          weightedComplexity: 1.600,
          calculatedAt: new Date().toISOString(),
          feature: {
            id: 1,
            title: 'User Authentication System',
            description: 'Implement secure login/logout functionality with JWT tokens',
            acceptanceCriteria: 'Users can login with email/password, stay logged in across sessions, secure token handling'
          }
        },
        {
          rank: 2,
          featureId: 2,
          finalPriorityScore: 4.688,
          sValue: 5,
          sComplexity: 3,
          wValue: 0.750,
          wComplexity: 0.400,
          weightedValue: 3.750,
          weightedComplexity: 1.200,
          calculatedAt: new Date().toISOString(),
          feature: {
            id: 2,
            title: 'Dashboard Analytics',
            description: 'Create comprehensive analytics dashboard with key metrics',
            acceptanceCriteria: 'Display user engagement metrics, performance charts, export functionality'
          }
        },
        {
          rank: 3,
          featureId: 3,
          finalPriorityScore: 3.542,
          sValue: 3,
          sComplexity: 5,
          wValue: 0.590,
          wComplexity: 0.300,
          weightedValue: 1.770,
          weightedComplexity: 1.500,
          calculatedAt: new Date().toISOString(),
          feature: {
            id: 3,
            title: 'Real-time Notifications',
            description: 'Push notifications for important events and updates',
            acceptanceCriteria: 'Browser notifications, email alerts, notification preferences'
          }
        }
      ]
    }
    
    return mockResults
  },

  async calculateResults(projectId) {
    // Mock implementation - replace with actual API call
    await new Promise(resolve => setTimeout(resolve, 2000)) // Simulate calculation time
    return this.getResults(projectId)
  },

  async exportResults(projectId, format) {
    // Mock implementation - replace with actual API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Simulate file download
    const fileName = `pairwise_results_${projectId}_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.${format === 'jira' ? 'json' : format}`
    
    // Create mock data based on format
    let content, mimeType
    switch (format) {
      case 'csv':
        content = 'rank,feature_title,description,final_priority_score\n1,"User Authentication System","Secure login system",6.250'
        mimeType = 'text/csv'
        break
      case 'json':
        content = JSON.stringify(await this.getResults(projectId), null, 2)
        mimeType = 'application/json'
        break
      case 'jira':
        content = JSON.stringify({
          issues: [{
            summary: 'User Authentication System',
            description: 'Secure login system',
            storyPoints: 2,
            priority: 'High',
            customFields: { finalPriorityScore: 6.250, valueScore: 8, complexityScore: 2 }
          }]
        }, null, 2)
        mimeType = 'application/json'
        break
    }
    
    // Trigger download
    const blob = new Blob([content], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }
}

const Results = () => {
  const theme = useTheme()
  const { projectId } = useParams()
  const navigate = useNavigate()
  const { project, loading: projectLoading } = useProject(projectId)
  
  const [results, setResults] = useState(null)
  const [selectedResult, setSelectedResult] = useState(null)
  const [loading, setLoading] = useState(true)
  const [calculating, setCalculating] = useState(false)
  const [hasResults, setHasResults] = useState(false)
  const [error, setError] = useState(null)
  const [tabValue, setTabValue] = useState(0)
  const [calculationDialog, setCalculationDialog] = useState(false)

  useEffect(() => {
    if (projectId) {
      checkResultsStatus()
    }
  }, [projectId])

  const checkResultsStatus = async () => {
    try {
      setLoading(true)
      const status = await resultsService.checkStatus(projectId)
      setHasResults(status.hasResults)
      
      if (status.hasResults) {
        const resultsData = await resultsService.getResults(projectId)
        setResults(resultsData)
        if (resultsData.results.length > 0) {
          setSelectedResult(resultsData.results[0])
        }
      }
    } catch (err) {
      setError('Failed to check results status: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleCalculateResults = async () => {
    try {
      setCalculating(true)
      setCalculationDialog(false)
      const resultsData = await resultsService.calculateResults(projectId)
      setResults(resultsData)
      setHasResults(true)
      if (resultsData.results.length > 0) {
        setSelectedResult(resultsData.results[0])
      }
      setError(null)
    } catch (err) {
      setError('Failed to calculate results: ' + err.message)
    } finally {
      setCalculating(false)
    }
  }

  const handleExport = async (format) => {
    try {
      await resultsService.exportResults(projectId, format)
    } catch (err) {
      throw new Error('Failed to export results: ' + err.message)
    }
  }

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue)
  }

  const LoadingSkeleton = () => (
    <Box>
      <Skeleton variant="text" sx={{ fontSize: '2rem' }} />
      <Skeleton variant="rectangular" width="100%" height={200} sx={{ my: 2 }} />
      <Grid container spacing={2}>
        {[1, 2, 3].map((i) => (
          <Grid item xs={12} md={4} key={i}>
            <Skeleton variant="rectangular" height={150} />
          </Grid>
        ))}
      </Grid>
    </Box>
  )

  if (loading || projectLoading) {
    return (
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <LoadingSkeleton />
      </Container>
    )
  }

  if (error) {
    return (
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Alert severity="error" action={
          <Button color="inherit" onClick={checkResultsStatus}>
            Retry
          </Button>
        }>
          {error}
        </Alert>
      </Container>
    )
  }

  if (!hasResults) {
    return (
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Card>
          <CardContent sx={{ textAlign: 'center', py: 6 }}>
            <AssessmentIcon sx={{ fontSize: 64, color: theme.palette.primary.main, mb: 2 }} />
            <Typography variant="h4" gutterBottom>
              PairWise Results Not Available
            </Typography>
            <Typography variant="body1" color="textSecondary" sx={{ mb: 4, maxWidth: 600, mx: 'auto' }}>
              No PairWise calculation results found for this project. You need to complete both 
              pairwise comparisons and Fibonacci scoring before calculating final priority scores.
            </Typography>
            
            <Box display="flex" gap={2} justifyContent="center" flexWrap="wrap">
              <Button
                variant="contained"
                size="large"
                startIcon={<CalculateIcon />}
                onClick={() => setCalculationDialog(true)}
                disabled={calculating}
              >
                {calculating ? 'Calculating...' : 'Calculate PairWise Results'}
              </Button>
              
              <Button
                variant="outlined"
                size="large"
                onClick={() => navigate(`/projects/${projectId}/comparison`)}
              >
                Pairwise Comparisons
              </Button>
              
              <Button
                variant="outlined"
                size="large"
                onClick={() => navigate(`/projects/${projectId}/scoring/value`)}
              >
                Fibonacci Scoring
              </Button>
            </Box>

            <Box mt={4} p={3} sx={{ backgroundColor: theme.palette.info.light + '10', borderRadius: 2, textAlign: 'left', maxWidth: 800, mx: 'auto' }}>
              <Typography variant="h6" gutterBottom>
                Prerequisites for PairWise Calculation
              </Typography>
              <Typography variant="body2" component="div">
                <ul>
                  <li><strong>Pairwise Comparisons:</strong> Complete head-to-head comparisons for both Value and Complexity criteria</li>
                  <li><strong>Fibonacci Scoring:</strong> Achieve team consensus on absolute magnitude scores using Fibonacci sequence</li>
                  <li><strong>Team Participation:</strong> Ensure all team members have participated in the scoring process</li>
                </ul>
              </Typography>
            </Box>
          </CardContent>
        </Card>
      </Container>
    )
  }

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      {/* Header */}
      <Box mb={4}>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={2}>
          <Box>
            <Typography variant="h4" gutterBottom>
              PairWise Results: {project?.name}
            </Typography>
            <Typography variant="body1" color="textSecondary">
              Final Priority Scores calculated on {new Date(results.calculatedAt).toLocaleString()}
            </Typography>
          </Box>
          <Box display="flex" gap={2}>
            <Button
              variant="outlined"
              startIcon={<RefreshIcon />}
              onClick={checkResultsStatus}
              disabled={calculating}
            >
              Refresh
            </Button>
            <Button
              variant="contained"
              startIcon={<CalculateIcon />}
              onClick={() => setCalculationDialog(true)}
              disabled={calculating}
            >
              Recalculate
            </Button>
          </Box>
        </Box>

        {/* Summary Stats */}
        <Paper sx={{ p: 3, mb: 3 }}>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6} md={3}>
              <Box textAlign="center">
                <Typography variant="h4" color="primary.main">{results.totalFeatures}</Typography>
                <Typography variant="subtitle2">Total Features</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Box textAlign="center">
                <Typography variant="h4" color="success.main">{results.summary.topTier}</Typography>
                <Typography variant="subtitle2">High Priority</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Box textAlign="center">
                <Typography variant="h4" color="info.main">{results.summary.highestScore.toFixed(3)}</Typography>
                <Typography variant="subtitle2">Highest Score</Typography>
              </Box>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Box textAlign="center">
                <Typography variant="h4" color="warning.main">{results.summary.averageScore.toFixed(3)}</Typography>
                <Typography variant="subtitle2">Average Score</Typography>
              </Box>
            </Grid>
          </Grid>
        </Paper>
      </Box>

      {/* Main Content Tabs */}
      <Paper sx={{ width: '100%' }}>
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          variant="fullWidth"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab icon={<AssessmentIcon />} label="Rankings" />
          <Tab icon={<VisibilityIcon />} label="Visualization" />
          <Tab icon={<InfoIcon />} label="Calculation Details" />
          <Tab icon={<CloudDownloadIcon />} label="Export" />
        </Tabs>

        <Box sx={{ p: 3 }}>
          {tabValue === 0 && (
            <ResultsRanking 
              results={results.results}
              onSelectResult={setSelectedResult}
            />
          )}

          {tabValue === 1 && (
            <ResultsVisualization 
              results={results.results}
              summary={results.summary}
            />
          )}

          {tabValue === 2 && (
            <CalculationDetails 
              result={selectedResult}
            />
          )}

          {tabValue === 3 && (
            <ExportOptions
              results={results.results}
              projectName={project?.name || `Project ${projectId}`}
              onExport={handleExport}
            />
          )}
        </Box>
      </Paper>

      {/* Methodology Note */}
      <Box mt={4}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              About PairWise Methodology
            </Typography>
            <Typography variant="body2" color="textSecondary">
              The PairWise (Pairwise-Weighted Value/Complexity) model combines objective mathematical calculations 
              with team consensus to provide unbiased feature prioritization. The Final Priority Score is calculated 
              using the formula: <strong>FPS = (S<sub>Value</sub> × W<sub>Value</sub>) ÷ (S<sub>Complexity</sub> × W<sub>Complexity</sub>)</strong>, 
              where S represents Fibonacci consensus scores and W represents win-count weights from pairwise comparisons.
            </Typography>
          </CardContent>
        </Card>
      </Box>

      {/* Calculation Confirmation Dialog */}
      <Dialog open={calculationDialog} onClose={() => !calculating && setCalculationDialog(false)}>
        <DialogTitle>Calculate PairWise Results</DialogTitle>
        <DialogContent>
          <DialogContentText>
            This will recalculate the Final Priority Scores for all features based on the current 
            pairwise comparison results and Fibonacci consensus scores. Any existing results will be replaced.
          </DialogContentText>
          <DialogContentText sx={{ mt: 2 }}>
            <strong>Note:</strong> Ensure that all team members have completed their pairwise comparisons 
            and Fibonacci scoring before proceeding.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCalculationDialog(false)} disabled={calculating}>
            Cancel
          </Button>
          <Button 
            onClick={handleCalculateResults} 
            variant="contained" 
            disabled={calculating}
            startIcon={calculating ? <RefreshIcon /> : <PlayArrowIcon />}
          >
            {calculating ? 'Calculating...' : 'Calculate Results'}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  )
}

export default Results