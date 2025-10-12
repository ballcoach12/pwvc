import {
    BarChart3 as BarChart3Icon,
    PieChart as PieChartIcon,
    ScatterPlot as ScatterPlotIcon,
    TrendingUp as TrendingUpIcon,
} from '@mui/icons-material'
import {
    Box,
    Card,
    CardContent,
    FormControl,
    Grid,
    InputLabel,
    MenuItem,
    Paper,
    Select,
    Typography,
    useTheme,
} from '@mui/material'
import { useState } from 'react'
import {
    Bar,
    BarChart,
    CartesianGrid,
    Cell,
    Legend,
    Pie,
    PieChart,
    RadialBar,
    RadialBarChart,
    ResponsiveContainer,
    Scatter,
    ScatterChart,
    Tooltip,
    XAxis,
    YAxis,
} from 'recharts'

const ResultsVisualization = ({ results, summary }) => {
  const theme = useTheme()
  const [chartType, setChartType] = useState('priority-bar')

  if (!results || results.length === 0) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Results Visualization
          </Typography>
          <Typography color="textSecondary">
            No results available to visualize. Please calculate P-WVC results first.
          </Typography>
        </CardContent>
      </Card>
    )
  }

  // Data transformations for different chart types
  const prepareBarChartData = () => {
    return results.slice(0, 15).map((result, index) => ({
      name: result.feature.title.length > 20 
        ? result.feature.title.substring(0, 17) + '...' 
        : result.feature.title,
      fullName: result.feature.title,
      priorityScore: Number(result.finalPriorityScore.toFixed(3)),
      rank: result.rank,
      value: result.sValue,
      complexity: result.sComplexity,
      color: index < 5 ? theme.palette.success.main : 
             index < 10 ? theme.palette.warning.main : theme.palette.error.main
    }))
  }

  const prepareScatterData = () => {
    return results.map((result) => ({
      x: result.sComplexity,
      y: result.sValue,
      z: result.finalPriorityScore,
      name: result.feature.title,
      rank: result.rank,
      color: result.rank <= Math.ceil(results.length / 4) ? theme.palette.success.main :
             result.rank <= Math.ceil(results.length / 2) ? theme.palette.warning.main :
             result.rank <= Math.ceil(3 * results.length / 4) ? theme.palette.info.main :
             theme.palette.error.main
    }))
  }

  const preparePieData = () => {
    const total = results.length
    const highPriority = results.filter(r => r.rank <= Math.ceil(total / 4)).length
    const mediumPriority = results.filter(r => r.rank > Math.ceil(total / 4) && r.rank <= Math.ceil(total / 2)).length
    const lowMediumPriority = results.filter(r => r.rank > Math.ceil(total / 2) && r.rank <= Math.ceil(3 * total / 4)).length
    const lowPriority = results.filter(r => r.rank > Math.ceil(3 * total / 4)).length

    return [
      { name: 'High Priority', value: highPriority, color: theme.palette.success.main },
      { name: 'Medium-High Priority', value: mediumPriority, color: theme.palette.warning.main },
      { name: 'Medium-Low Priority', value: lowMediumPriority, color: theme.palette.info.main },
      { name: 'Low Priority', value: lowPriority, color: theme.palette.error.main },
    ]
  }

  const prepareRadialData = () => {
    return results.slice(0, 10).map((result, index) => ({
      name: result.feature.title.length > 15 
        ? result.feature.title.substring(0, 12) + '...' 
        : result.feature.title,
      fullName: result.feature.title,
      value: Number(result.finalPriorityScore.toFixed(3)),
      fill: `hsl(${120 - (index * 12)}, 70%, 50%)` // Gradient from green to red
    }))
  }

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <Paper sx={{ p: 2, border: `1px solid ${theme.palette.divider}` }}>
          <Typography variant="subtitle2" fontWeight="bold">
            {data.fullName || data.name || label}
          </Typography>
          {data.rank && (
            <Typography variant="body2">
              Rank: #{data.rank}
            </Typography>
          )}
          {data.priorityScore && (
            <Typography variant="body2">
              Priority Score: {data.priorityScore}
            </Typography>
          )}
          {data.value && data.complexity && (
            <Typography variant="body2">
              Value: {data.value} | Complexity: {data.complexity}
            </Typography>
          )}
        </Paper>
      )
    }
    return null
  }

  const ScatterTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <Paper sx={{ p: 2, border: `1px solid ${theme.palette.divider}` }}>
          <Typography variant="subtitle2" fontWeight="bold">
            {data.name}
          </Typography>
          <Typography variant="body2">
            Rank: #{data.rank}
          </Typography>
          <Typography variant="body2">
            Value: {data.y} | Complexity: {data.x}
          </Typography>
          <Typography variant="body2">
            Priority Score: {data.z.toFixed(3)}
          </Typography>
        </Paper>
      )
    }
    return null
  }

  const renderChart = () => {
    switch (chartType) {
      case 'priority-bar':
        return (
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={prepareBarChartData()} margin={{ top: 20, right: 30, left: 20, bottom: 60 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey="name" 
                angle={-45}
                textAnchor="end"
                height={80}
                fontSize={12}
              />
              <YAxis label={{ value: 'Priority Score', angle: -90, position: 'insideLeft' }} />
              <Tooltip content={<CustomTooltip />} />
              <Legend />
              <Bar dataKey="priorityScore" name="Final Priority Score">
                {prepareBarChartData().map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        )

      case 'value-complexity-scatter':
        return (
          <ResponsiveContainer width="100%" height={400}>
            <ScatterChart margin={{ top: 20, right: 20, bottom: 20, left: 20 }}>
              <CartesianGrid />
              <XAxis 
                type="number" 
                dataKey="x" 
                name="Complexity"
                label={{ value: 'Complexity Score', position: 'insideBottom', offset: -10 }}
              />
              <YAxis 
                type="number" 
                dataKey="y" 
                name="Value"
                label={{ value: 'Value Score', angle: -90, position: 'insideLeft' }}
              />
              <Tooltip content={<ScatterTooltip />} />
              <Scatter data={prepareScatterData()} fill={theme.palette.primary.main}>
                {prepareScatterData().map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Scatter>
            </ScatterChart>
          </ResponsiveContainer>
        )

      case 'priority-distribution':
        return (
          <ResponsiveContainer width="100%" height={400}>
            <PieChart>
              <Pie
                data={preparePieData()}
                cx="50%"
                cy="50%"
                outerRadius={120}
                dataKey="value"
                label={({ name, value, percent }) => `${name}: ${value} (${(percent * 100).toFixed(0)}%)`}
              >
                {preparePieData().map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        )

      case 'top-features-radial':
        return (
          <ResponsiveContainer width="100%" height={400}>
            <RadialBarChart cx="50%" cy="50%" innerRadius="10%" outerRadius="80%" data={prepareRadialData()}>
              <RadialBar 
                dataKey="value" 
                cornerRadius={10} 
                fill={theme.palette.primary.main}
              />
              <Tooltip content={<CustomTooltip />} />
              <Legend />
            </RadialBarChart>
          </ResponsiveContainer>
        )

      default:
        return null
    }
  }

  const ChartOption = ({ value, icon, title, description }) => (
    <MenuItem value={value} sx={{ p: 2 }}>
      <Box display="flex" alignItems="center" width="100%">
        <Box sx={{ mr: 2, color: theme.palette.primary.main }}>
          {icon}
        </Box>
        <Box>
          <Typography variant="subtitle2">{title}</Typography>
          <Typography variant="body2" color="textSecondary">
            {description}
          </Typography>
        </Box>
      </Box>
    </MenuItem>
  )

  return (
    <Card>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h6" component="h2">
            Results Visualization
          </Typography>
          <FormControl variant="outlined" sx={{ minWidth: 250 }}>
            <InputLabel>Chart Type</InputLabel>
            <Select
              value={chartType}
              onChange={(e) => setChartType(e.target.value)}
              label="Chart Type"
            >
              <ChartOption
                value="priority-bar"
                icon={<BarChart3Icon />}
                title="Priority Bar Chart"
                description="Features ranked by priority score"
              />
              <ChartOption
                value="value-complexity-scatter"
                icon={<ScatterPlotIcon />}
                title="Value vs Complexity"
                description="Feature positioning by value and complexity"
              />
              <ChartOption
                value="priority-distribution"
                icon={<PieChartIcon />}
                title="Priority Distribution"
                description="Breakdown of features by priority level"
              />
              <ChartOption
                value="top-features-radial"
                icon={<TrendingUpIcon />}
                title="Top Features Radial"
                description="Top 10 features in radial chart"
              />
            </Select>
          </FormControl>
        </Box>

        {/* Chart Description */}
        <Box mb={3}>
          <Typography variant="body2" color="textSecondary">
            {chartType === 'priority-bar' && 'Bar chart showing the top 15 features ranked by Final Priority Score. Higher bars indicate higher priority.'}
            {chartType === 'value-complexity-scatter' && 'Scatter plot positioning features by Value (Y-axis) and Complexity (X-axis). Features in the upper-left have high value and low complexity.'}
            {chartType === 'priority-distribution' && 'Pie chart showing the distribution of features across priority levels (High, Medium-High, Medium-Low, Low).'}
            {chartType === 'top-features-radial' && 'Radial bar chart highlighting the top 10 features by priority score in a circular format.'}
          </Typography>
        </Box>

        {/* Chart Container */}
        <Paper variant="outlined" sx={{ p: 2, backgroundColor: theme.palette.background.default }}>
          {renderChart()}
        </Paper>

        {/* Chart Insights */}
        {summary && (
          <Box mt={3}>
            <Typography variant="subtitle1" gutterBottom fontWeight="bold">
              Key Insights
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6} md={3}>
                <Box textAlign="center">
                  <Typography variant="h6" color="success.main">
                    {summary.topTier}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    High Priority Features
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Box textAlign="center">
                  <Typography variant="h6" color="primary.main">
                    {summary.highestScore.toFixed(3)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Highest Priority Score
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Box textAlign="center">
                  <Typography variant="h6" color="info.main">
                    {summary.averageScore.toFixed(3)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Average Priority Score
                  </Typography>
                </Box>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Box textAlign="center">
                  <Typography variant="h6" color="warning.main">
                    {summary.scoreRange.toFixed(3)}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Score Range
                  </Typography>
                </Box>
              </Grid>
            </Grid>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

export default ResultsVisualization