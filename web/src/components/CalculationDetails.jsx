import {
    Assessment as AssessmentIcon,
    Calculate as CalculateIcon,
    Functions as FunctionsIcon,
    TrendingUp as TrendingUpIcon,
} from '@mui/icons-material'
import {
    Box,
    Card,
    CardContent,
    Divider,
    Grid,
    Paper,
    Typography,
    useTheme,
} from '@mui/material'

const CalculationDetails = ({ result }) => {
  const theme = useTheme()

  if (!result) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            P-WVC Calculation Details
          </Typography>
          <Typography color="textSecondary">
            Select a feature to see calculation breakdown.
          </Typography>
        </CardContent>
      </Card>
    )
  }

  const FormulaDisplay = ({ title, formula, description }) => (
    <Paper 
      elevation={0} 
      sx={{ 
        p: 2, 
        backgroundColor: theme.palette.grey[50],
        border: `1px solid ${theme.palette.divider}`
      }}
    >
      <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 'bold' }}>
        {title}
      </Typography>
      <Typography 
        variant="body1" 
        sx={{ 
          fontFamily: 'monospace', 
          backgroundColor: theme.palette.background.paper,
          p: 1,
          borderRadius: 1,
          border: `1px solid ${theme.palette.divider}`,
          mb: 1
        }}
      >
        {formula}
      </Typography>
      <Typography variant="body2" color="textSecondary">
        {description}
      </Typography>
    </Paper>
  )

  const MetricCard = ({ icon, title, value, subtitle, color = 'primary' }) => (
    <Card variant="outlined">
      <CardContent sx={{ textAlign: 'center' }}>
        <Box display="flex" justifyContent="center" mb={1}>
          {icon}
        </Box>
        <Typography variant="h4" color={color} fontWeight="bold">
          {value}
        </Typography>
        <Typography variant="subtitle1" fontWeight="medium">
          {title}
        </Typography>
        {subtitle && (
          <Typography variant="body2" color="textSecondary">
            {subtitle}
          </Typography>
        )}
      </CardContent>
    </Card>
  )

  return (
    <Card>
      <CardContent>
        <Box display="flex" alignItems="center" mb={3}>
          <CalculateIcon sx={{ mr: 2, color: theme.palette.primary.main }} />
          <Typography variant="h6" component="h2">
            P-WVC Calculation Breakdown: "{result.feature.title}"
          </Typography>
        </Box>

        {/* Key Metrics */}
        <Grid container spacing={2} mb={4}>
          <Grid item xs={12} sm={6} md={3}>
            <MetricCard
              icon={<AssessmentIcon sx={{ fontSize: 32, color: theme.palette.success.main }} />}
              title="Final Priority Score"
              value={result.finalPriorityScore.toFixed(3)}
              subtitle={`Rank #${result.rank}`}
              color="success.main"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <MetricCard
              icon={<TrendingUpIcon sx={{ fontSize: 32, color: theme.palette.primary.main }} />}
              title="Value Score"
              value={result.sValue}
              subtitle="Fibonacci Scale"
              color="primary.main"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <MetricCard
              icon={<FunctionsIcon sx={{ fontSize: 32, color: theme.palette.secondary.main }} />}
              title="Complexity Score"
              value={result.sComplexity}
              subtitle="Fibonacci Scale"
              color="secondary.main"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <MetricCard
              icon={<CalculateIcon sx={{ fontSize: 32, color: theme.palette.info.main }} />}
              title="Value/Complexity Ratio"
              value={(result.sValue / result.sComplexity).toFixed(2)}
              subtitle="Simple Ratio"
              color="info.main"
            />
          </Grid>
        </Grid>

        <Divider sx={{ my: 3 }} />

        {/* P-WVC Formula Explanation */}
        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
          <FunctionsIcon sx={{ mr: 1 }} />
          P-WVC Mathematical Formula
        </Typography>
        
        <Box mb={3}>
          <FormulaDisplay
            title="Final Priority Score (FPS)"
            formula="FPS = (S_Value × W_Value) ÷ (S_Complexity × W_Complexity)"
            description="The Final Priority Score combines Fibonacci scoring with win-count weights from pairwise comparisons to provide an objective ranking."
          />
        </Box>

        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom fontWeight="bold">
              Value Components
            </Typography>
            <Box sx={{ pl: 2 }}>
              <Typography variant="body2" gutterBottom>
                <strong>S<sub>Value</sub> (Fibonacci Score):</strong> {result.sValue}
              </Typography>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Consensus score from team using Fibonacci sequence (1,2,3,5,8,13,21...)
              </Typography>
              
              <Typography variant="body2" gutterBottom>
                <strong>W<sub>Value</sub> (Win-count Weight):</strong> {result.wValue.toFixed(6)}
              </Typography>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Calculated from pairwise comparisons: W = wins ÷ total_comparisons
              </Typography>
              
              <Typography variant="body2" gutterBottom>
                <strong>Weighted Value:</strong> {result.weightedValue.toFixed(6)}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                S<sub>Value</sub> × W<sub>Value</sub> = {result.sValue} × {result.wValue.toFixed(6)}
              </Typography>
            </Box>
          </Grid>

          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom fontWeight="bold">
              Complexity Components
            </Typography>
            <Box sx={{ pl: 2 }}>
              <Typography variant="body2" gutterBottom>
                <strong>S<sub>Complexity</sub> (Fibonacci Score):</strong> {result.sComplexity}
              </Typography>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Consensus score from team using Fibonacci sequence (1,2,3,5,8,13,21...)
              </Typography>
              
              <Typography variant="body2" gutterBottom>
                <strong>W<sub>Complexity</sub> (Win-count Weight):</strong> {result.wComplexity.toFixed(6)}
              </Typography>
              <Typography variant="body2" color="textSecondary" gutterBottom>
                Calculated from pairwise comparisons: W = wins ÷ total_comparisons
              </Typography>
              
              <Typography variant="body2" gutterBottom>
                <strong>Weighted Complexity:</strong> {result.weightedComplexity.toFixed(6)}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                S<sub>Complexity</sub> × W<sub>Complexity</sub> = {result.sComplexity} × {result.wComplexity.toFixed(6)}
              </Typography>
            </Box>
          </Grid>
        </Grid>

        <Divider sx={{ my: 3 }} />

        {/* Final Calculation */}
        <Typography variant="h6" gutterBottom>
          Final Calculation
        </Typography>
        
        <Paper 
          sx={{ 
            p: 3, 
            backgroundColor: theme.palette.success.light + '10',
            border: `2px solid ${theme.palette.success.main}`,
            borderRadius: 2
          }}
        >
          <Typography variant="h6" align="center" gutterBottom>
            Final Priority Score Calculation
          </Typography>
          <Typography 
            variant="h5" 
            align="center" 
            sx={{ 
              fontFamily: 'monospace',
              fontWeight: 'bold',
              mb: 2
            }}
          >
            FPS = {result.weightedValue.toFixed(6)} ÷ {result.weightedComplexity.toFixed(6)}
          </Typography>
          <Typography 
            variant="h4" 
            align="center" 
            color="success.main"
            fontWeight="bold"
          >
            = {result.finalPriorityScore.toFixed(6)}
          </Typography>
        </Paper>

        {/* Methodology Note */}
        <Box mt={3} p={2} sx={{ backgroundColor: theme.palette.info.light + '10', borderRadius: 1 }}>
          <Typography variant="subtitle2" fontWeight="bold" gutterBottom>
            P-WVC Methodology Note
          </Typography>
          <Typography variant="body2" color="textSecondary">
            The P-WVC (Pairwise-Weighted Value/Complexity) model ensures objective feature prioritization 
            by combining team consensus through Fibonacci scoring with mathematical win-count weights 
            from pairwise comparisons. Higher scores indicate features with better value-to-complexity ratios.
          </Typography>
        </Box>
      </CardContent>
    </Card>
  )
}

export default CalculationDetails