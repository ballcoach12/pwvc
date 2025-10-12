import {
    ExpandMore as ExpandMoreIcon,
    FilterList as FilterListIcon,
    Sort as SortIcon,
    TrendingUp as TrendingUpIcon,
} from '@mui/icons-material'
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Card,
    CardContent,
    Chip,
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
    useTheme,
} from '@mui/material'
import { useState } from 'react'

const ResultsRanking = ({ results, onExport }) => {
  const theme = useTheme()
  const [expanded, setExpanded] = useState(false)
  const [sortField, setSortField] = useState('rank')
  const [sortDirection, setSortDirection] = useState('asc')

  const handleAccordionChange = (panel) => (event, isExpanded) => {
    setExpanded(isExpanded ? panel : false)
  }

  const handleSort = (field) => {
    const isAsc = sortField === field && sortDirection === 'asc'
    setSortDirection(isAsc ? 'desc' : 'asc')
    setSortField(field)
  }

  const sortedResults = [...results].sort((a, b) => {
    let aValue, bValue
    
    switch (sortField) {
      case 'rank':
        aValue = a.rank
        bValue = b.rank
        break
      case 'title':
        aValue = a.feature.title
        bValue = b.feature.title
        break
      case 'finalPriorityScore':
        aValue = a.finalPriorityScore
        bValue = b.finalPriorityScore
        break
      case 'sValue':
        aValue = a.sValue
        bValue = b.sValue
        break
      case 'sComplexity':
        aValue = a.sComplexity
        bValue = b.sComplexity
        break
      default:
        return 0
    }

    if (typeof aValue === 'string') {
      aValue = aValue.toLowerCase()
      bValue = bValue.toLowerCase()
    }

    if (sortDirection === 'asc') {
      return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
    } else {
      return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
    }
  })

  const getRankColor = (rank, total) => {
    const percentage = rank / total
    if (percentage <= 0.25) return theme.palette.success.main
    if (percentage <= 0.5) return theme.palette.warning.main
    if (percentage <= 0.75) return theme.palette.info.main
    return theme.palette.error.main
  }

  const getRankLabel = (rank, total) => {
    const percentage = rank / total
    if (percentage <= 0.25) return 'High Priority'
    if (percentage <= 0.5) return 'Medium-High Priority'
    if (percentage <= 0.75) return 'Medium-Low Priority'
    return 'Low Priority'
  }

  if (!results || results.length === 0) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Feature Rankings
          </Typography>
          <Typography color="textSecondary">
            No results available. Please calculate P-WVC results first.
          </Typography>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h6" component="h2">
            Feature Rankings by Final Priority Score
          </Typography>
          <Box>
            <IconButton onClick={() => handleSort('finalPriorityScore')}>
              <SortIcon />
            </IconButton>
            <IconButton>
              <FilterListIcon />
            </IconButton>
          </Box>
        </Box>

        <TableContainer component={Paper} variant="outlined">
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Rank</TableCell>
                <TableCell onClick={() => handleSort('title')} sx={{ cursor: 'pointer' }}>
                  Feature Title
                  {sortField === 'title' && (
                    <SortIcon fontSize="small" sx={{ ml: 1 }} />
                  )}
                </TableCell>
                <TableCell 
                  onClick={() => handleSort('finalPriorityScore')}
                  sx={{ cursor: 'pointer' }}
                  align="right"
                >
                  Priority Score
                  {sortField === 'finalPriorityScore' && (
                    <SortIcon fontSize="small" sx={{ ml: 1 }} />
                  )}
                </TableCell>
                <TableCell align="center">Value</TableCell>
                <TableCell align="center">Complexity</TableCell>
                <TableCell align="center">Priority Level</TableCell>
                <TableCell></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sortedResults.map((result, index) => (
                <TableRow key={result.featureId} hover>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Typography
                        variant="h6"
                        sx={{
                          color: getRankColor(result.rank, results.length),
                          fontWeight: 'bold',
                        }}
                      >
                        #{result.rank}
                      </Typography>
                      {result.rank <= 3 && (
                        <TrendingUpIcon
                          sx={{ ml: 1, color: theme.palette.success.main }}
                        />
                      )}
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Box>
                      <Typography variant="subtitle2" fontWeight="medium">
                        {result.feature.title}
                      </Typography>
                      <Typography
                        variant="body2"
                        color="textSecondary"
                        sx={{
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          whiteSpace: 'nowrap',
                          maxWidth: '300px',
                        }}
                      >
                        {result.feature.description}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="h6" fontWeight="bold">
                      {result.finalPriorityScore.toFixed(3)}
                    </Typography>
                  </TableCell>
                  <TableCell align="center">
                    <Chip
                      label={result.sValue}
                      size="small"
                      color="primary"
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell align="center">
                    <Chip
                      label={result.sComplexity}
                      size="small"
                      color="secondary"
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell align="center">
                    <Chip
                      label={getRankLabel(result.rank, results.length)}
                      size="small"
                      sx={{
                        backgroundColor: getRankColor(result.rank, results.length),
                        color: 'white',
                      }}
                    />
                  </TableCell>
                  <TableCell>
                    <Accordion
                      expanded={expanded === `panel${index}`}
                      onChange={handleAccordionChange(`panel${index}`)}
                      elevation={0}
                      sx={{ boxShadow: 'none' }}
                    >
                      <AccordionSummary
                        expandIcon={<ExpandMoreIcon />}
                        sx={{ minHeight: 0, '& .MuiAccordionSummary-content': { margin: 0 } }}
                      >
                        <Typography variant="caption">Details</Typography>
                      </AccordionSummary>
                    </Accordion>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {/* Expanded calculation details */}
        {sortedResults.map((result, index) => (
          <Accordion
            key={`detail-${result.featureId}`}
            expanded={expanded === `panel${index}`}
            onChange={handleAccordionChange(`panel${index}`)}
          >
            <AccordionSummary />
            <AccordionDetails>
              <Box sx={{ pl: 2, pr: 2, pb: 2 }}>
                <Typography variant="h6" gutterBottom>
                  Calculation Details for "{result.feature.title}"
                </Typography>
                
                <Box display="flex" gap={4} flexWrap="wrap" mb={2}>
                  <Box>
                    <Typography variant="subtitle2" color="textSecondary">
                      Value Components
                    </Typography>
                    <Typography>Fibonacci Score (S<sub>Value</sub>): {result.sValue}</Typography>
                    <Typography>Win-count Weight (W<sub>Value</sub>): {result.wValue.toFixed(3)}</Typography>
                    <Typography>Weighted Value: {result.weightedValue.toFixed(3)}</Typography>
                  </Box>
                  
                  <Box>
                    <Typography variant="subtitle2" color="textSecondary">
                      Complexity Components
                    </Typography>
                    <Typography>Fibonacci Score (S<sub>Complexity</sub>): {result.sComplexity}</Typography>
                    <Typography>Win-count Weight (W<sub>Complexity</sub>): {result.wComplexity.toFixed(3)}</Typography>
                    <Typography>Weighted Complexity: {result.weightedComplexity.toFixed(3)}</Typography>
                  </Box>
                  
                  <Box>
                    <Typography variant="subtitle2" color="textSecondary">
                      Final Calculation
                    </Typography>
                    <Typography fontWeight="bold">
                      FPS = {result.weightedValue.toFixed(3)} ÷ {result.weightedComplexity.toFixed(3)} = {result.finalPriorityScore.toFixed(6)}
                    </Typography>
                  </Box>
                </Box>

                {result.feature.acceptanceCriteria && (
                  <Box mt={2}>
                    <Typography variant="subtitle2" color="textSecondary">
                      Acceptance Criteria
                    </Typography>
                    <Typography variant="body2">
                      {result.feature.acceptanceCriteria}
                    </Typography>
                  </Box>
                )}
              </Box>
            </AccordionDetails>
          </Accordion>
        ))}
      </CardContent>
    </Card>
  )
}

export default ResultsRanking