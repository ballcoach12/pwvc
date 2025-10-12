import {
    BugReport as BugReportIcon,
    Close as CloseIcon,
    CloudDownload as CloudDownloadIcon,
    Code as CodeIcon,
    GetApp as GetAppIcon,
    TableView as TableViewIcon,
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Divider,
    Grid,
    IconButton,
    Paper,
    Snackbar,
    Typography,
    useTheme,
} from '@mui/material'
import { useState } from 'react'

const ExportOptions = ({ results, projectName, onExport }) => {
  const theme = useTheme()
  const [exportDialog, setExportDialog] = useState(false)
  const [exportFormat, setExportFormat] = useState('')
  const [loading, setLoading] = useState(false)
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' })

  const handleExportClick = (format) => {
    setExportFormat(format)
    setExportDialog(true)
  }

  const handleExportConfirm = async () => {
    setLoading(true)
    try {
      await onExport(exportFormat)
      setSnackbar({
        open: true,
        message: `Successfully exported results as ${exportFormat.toUpperCase()}`,
        severity: 'success'
      })
      setExportDialog(false)
    } catch (error) {
      setSnackbar({
        open: true,
        message: `Failed to export: ${error.message}`,
        severity: 'error'
      })
    } finally {
      setLoading(false)
    }
  }

  const getExportDescription = (format) => {
    switch (format) {
      case 'csv':
        return 'Export as CSV file for spreadsheet analysis. Includes all calculation details and can be opened in Excel, Google Sheets, etc.'
      case 'json':
        return 'Export as JSON format for API integration and programmatic use. Contains complete data structure with all nested details.'
      case 'jira':
        return 'Export as Jira-compatible JSON with story points assigned based on complexity scores. Ready for direct import into Jira projects.'
      default:
        return ''
    }
  }

  const getFileName = (format) => {
    const timestamp = new Date().toISOString().slice(0, 19).replace(/:/g, '-')
    const sanitizedProjectName = projectName.replace(/[^a-z0-9]/gi, '_').toLowerCase()
    return `pwvc_results_${sanitizedProjectName}_${timestamp}.${format === 'jira' ? 'json' : format}`
  }

  const ExportCard = ({ format, title, icon, description, color }) => (
    <Card 
      variant="outlined" 
      sx={{ 
        height: '100%',
        cursor: 'pointer',
        transition: 'all 0.2s',
        '&:hover': {
          boxShadow: theme.shadows[4],
          transform: 'translateY(-2px)',
          borderColor: color,
        }
      }}
      onClick={() => handleExportClick(format)}
    >
      <CardContent sx={{ textAlign: 'center', p: 3 }}>
        <Box display="flex" justifyContent="center" mb={2}>
          {icon}
        </Box>
        <Typography variant="h6" gutterBottom fontWeight="bold">
          {title}
        </Typography>
        <Typography variant="body2" color="textSecondary" mb={3}>
          {description}
        </Typography>
        <Button
          variant="outlined"
          startIcon={<GetAppIcon />}
          fullWidth
          sx={{ borderColor: color, color: color }}
        >
          Export {format.toUpperCase()}
        </Button>
      </CardContent>
    </Card>
  )

  const PreviewSection = ({ title, content, language = 'text' }) => (
    <Box mb={3}>
      <Typography variant="subtitle2" gutterBottom fontWeight="bold">
        {title}
      </Typography>
      <Paper 
        sx={{ 
          p: 2, 
          backgroundColor: theme.palette.grey[50],
          border: `1px solid ${theme.palette.divider}`,
          maxHeight: 200,
          overflow: 'auto'
        }}
      >
        <Typography 
          variant="body2" 
          component="pre" 
          sx={{ 
            fontFamily: 'monospace',
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-word',
            margin: 0
          }}
        >
          {content}
        </Typography>
      </Paper>
    </Box>
  )

  const generatePreviewContent = (format) => {
    if (!results || results.length === 0) return 'No results available'

    const firstResult = results[0]
    
    switch (format) {
      case 'csv':
        return `rank,feature_title,description,final_priority_score,s_value,s_complexity,w_value,w_complexity
1,"${firstResult.feature.title}","${firstResult.feature.description}",${firstResult.finalPriorityScore.toFixed(6)},${firstResult.sValue},${firstResult.sComplexity},${firstResult.wValue.toFixed(6)},${firstResult.wComplexity.toFixed(6)}
...`
        
      case 'json':
        return JSON.stringify({
          projectId: results[0]?.projectId,
          calculatedAt: new Date().toISOString(),
          totalFeatures: results.length,
          results: [firstResult],
          "...": `${results.length - 1} more results`
        }, null, 2)
        
      case 'jira':
        const priority = firstResult.rank <= Math.ceil(results.length / 4) ? 'High' : 
                        firstResult.rank <= Math.ceil(results.length / 2) ? 'Medium' : 'Low'
        return JSON.stringify({
          issues: [{
            summary: firstResult.feature.title,
            description: firstResult.feature.description + (firstResult.feature.acceptanceCriteria ? '\n\nAcceptance Criteria:\n' + firstResult.feature.acceptanceCriteria : ''),
            storyPoints: Math.min(firstResult.sComplexity, 21),
            priority: priority,
            customFields: {
              finalPriorityScore: firstResult.finalPriorityScore,
              valueScore: firstResult.sValue,
              complexityScore: firstResult.sComplexity
            }
          }],
          "...": `${results.length - 1} more issues`
        }, null, 2)
        
      default:
        return ''
    }
  }

  if (!results || results.length === 0) {
    return (
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Export Results
          </Typography>
          <Typography color="textSecondary">
            No results available to export. Please calculate P-WVC results first.
          </Typography>
        </CardContent>
      </Card>
    )
  }

  return (
    <>
      <Card>
        <CardContent>
          <Box display="flex" alignItems="center" mb={3}>
            <CloudDownloadIcon sx={{ mr: 2, color: theme.palette.primary.main }} />
            <Typography variant="h6" component="h2">
              Export P-WVC Results
            </Typography>
          </Box>

          <Typography variant="body2" color="textSecondary" mb={4}>
            Export your P-WVC calculation results in multiple formats for different use cases.
            All exports include complete calculation details and feature rankings.
          </Typography>

          <Grid container spacing={3}>
            <Grid item xs={12} md={4}>
              <ExportCard
                format="csv"
                title="CSV Export"
                icon={<TableViewIcon sx={{ fontSize: 48, color: theme.palette.success.main }} />}
                description="Spreadsheet-compatible format for analysis in Excel, Google Sheets, or other tools."
                color={theme.palette.success.main}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <ExportCard
                format="json"
                title="JSON Export"
                icon={<CodeIcon sx={{ fontSize: 48, color: theme.palette.info.main }} />}
                description="Structured data format for API integration and programmatic processing."
                color={theme.palette.info.main}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <ExportCard
                format="jira"
                title="Jira Export"
                icon={<BugReportIcon sx={{ fontSize: 48, color: theme.palette.warning.main }} />}
                description="Jira-compatible format with story points and custom fields for project management."
                color={theme.palette.warning.main}
              />
            </Grid>
          </Grid>

          <Divider sx={{ my: 4 }} />

          <Typography variant="h6" gutterBottom>
            Export Summary
          </Typography>
          <Box display="flex" gap={4} flexWrap="wrap">
            <Box>
              <Typography variant="subtitle2" color="textSecondary">Total Features</Typography>
              <Typography variant="h6">{results.length}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" color="textSecondary">Calculation Date</Typography>
              <Typography variant="body1">{new Date(results[0]?.calculatedAt).toLocaleString()}</Typography>
            </Box>
            <Box>
              <Typography variant="subtitle2" color="textSecondary">Project</Typography>
              <Typography variant="body1">{projectName}</Typography>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {/* Export Confirmation Dialog */}
      <Dialog 
        open={exportDialog} 
        onClose={() => !loading && setExportDialog(false)}
        maxWidth="md" 
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          Export as {exportFormat?.toUpperCase()}
          <IconButton onClick={() => setExportDialog(false)} disabled={loading}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          <DialogContentText mb={3}>
            {getExportDescription(exportFormat)}
          </DialogContentText>
          
          <Typography variant="subtitle1" gutterBottom>
            File name: <code>{getFileName(exportFormat)}</code>
          </Typography>

          <PreviewSection 
            title="Preview (first few lines)"
            content={generatePreviewContent(exportFormat)}
          />

          <Box mt={2} p={2} sx={{ backgroundColor: theme.palette.warning.light + '10', borderRadius: 1 }}>
            <Typography variant="body2">
              <strong>Note:</strong> The export will include all {results.length} features with complete 
              calculation details, rankings, and metadata.
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions sx={{ p: 3 }}>
          <Button onClick={() => setExportDialog(false)} disabled={loading}>
            Cancel
          </Button>
          <Button 
            onClick={handleExportConfirm} 
            variant="contained" 
            disabled={loading}
            startIcon={<GetAppIcon />}
          >
            {loading ? 'Exporting...' : `Export ${exportFormat?.toUpperCase()}`}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Success/Error Snackbar */}
      <Snackbar 
        open={snackbar.open} 
        autoHideDuration={6000} 
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert severity={snackbar.severity} onClose={() => setSnackbar({ ...snackbar, open: false })}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </>
  )
}

export default ExportOptions