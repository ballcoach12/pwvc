import {
    CheckCircle,
    CloudUpload,
    Error,
    Info,
    Warning,
} from '@mui/icons-material'
import {
    Alert,
    Box,
    Button,
    LinearProgress,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    Paper,
    Typography,
} from '@mui/material'
import { useRef, useState } from 'react'

const FileUpload = ({
  onFileUpload,
  onValidationResults,
  accept = '.csv',
  maxSize = 5 * 1024 * 1024, // 5MB
  disabled = false,
}) => {
  const [uploadStatus, setUploadStatus] = useState('idle') // idle, uploading, success, error
  const [validationResults, setValidationResults] = useState(null)
  const [errorMessage, setErrorMessage] = useState('')
  const [isDragActive, setIsDragActive] = useState(false)
  const fileInputRef = useRef(null)

  const handleFile = async (file) => {
    // Validate file size
    if (file.size > maxSize) {
      setErrorMessage(`File is too large. Maximum size is ${maxSize / (1024 * 1024)}MB.`)
      setUploadStatus('error')
      return
    }

    // Validate file type
    if (!file.name.toLowerCase().endsWith('.csv')) {
      setErrorMessage('Invalid file type. Please upload a CSV file.')
      setUploadStatus('error')
      return
    }
    setUploadStatus('uploading')
    setErrorMessage('')
    setValidationResults(null)

    try {
      // Read file content
      const text = await file.text()
      
      // Basic validation
      const lines = text.split('\n').filter(line => line.trim())
      if (lines.length === 0) {
        throw new Error('File is empty')
      }

      // Parse CSV (basic parsing - in real app, use proper CSV parser)
      const rows = lines.map(line => line.split(',').map(cell => cell.trim()))
      const headers = rows[0]
      const dataRows = rows.slice(1)

      // Validate CSV structure for features
      const results = validateFeatureCSV(headers, dataRows)
      
      setValidationResults(results)
      onValidationResults?.(results)
      
      if (results.isValid) {
        setUploadStatus('success')
        onFileUpload({
          file,
          data: results.features,
          headers,
          totalRows: dataRows.length,
        })
      } else {
        setUploadStatus('error')
      }
    } catch (error) {
      setErrorMessage(error.message || 'Failed to process file')
      setUploadStatus('error')
    }
  }

  const handleFileChange = (event) => {
    const file = event.target.files[0]
    if (file) {
      handleFile(file)
    }
  }

  const handleDragOver = (event) => {
    event.preventDefault()
    setIsDragActive(true)
  }

  const handleDragLeave = (event) => {
    event.preventDefault()
    setIsDragActive(false)
  }

  const handleDrop = (event) => {
    event.preventDefault()
    setIsDragActive(false)
    
    const file = event.dataTransfer.files[0]
    if (file) {
      handleFile(file)
    }
  }

  const handleClick = () => {
    if (!disabled) {
      fileInputRef.current?.click()
    }
  }

  const validateFeatureCSV = (headers, dataRows) => {
    const results = {
      isValid: false,
      features: [],
      errors: [],
      warnings: [],
      info: [],
    }

    // Check required headers
    const requiredHeaders = ['name']
    const optionalHeaders = ['description']
    const lowerHeaders = headers.map(h => h.toLowerCase())

    const missingHeaders = requiredHeaders.filter(
      header => !lowerHeaders.includes(header.toLowerCase())
    )

    if (missingHeaders.length > 0) {
      results.errors.push(`Missing required headers: ${missingHeaders.join(', ')}`)
      return results
    }

    // Find header indices
    const nameIndex = lowerHeaders.indexOf('name')
    const descriptionIndex = lowerHeaders.indexOf('description')

    // Validate data rows
    const validFeatures = []
    const skippedRows = []

    dataRows.forEach((row, index) => {
      const rowNumber = index + 2 // +2 because we start from row 1 and skip header

      if (row.length === 0 || row.every(cell => !cell)) {
        skippedRows.push(rowNumber)
        return
      }

      const name = row[nameIndex]?.trim()
      if (!name) {
        results.errors.push(`Row ${rowNumber}: Feature name is required`)
        return
      }

      if (name.length < 3) {
        results.errors.push(`Row ${rowNumber}: Feature name must be at least 3 characters`)
        return
      }

      const description = descriptionIndex >= 0 ? row[descriptionIndex]?.trim() || '' : ''

      validFeatures.push({
        name,
        description,
        rowNumber,
      })
    })

    // Add info and warnings
    if (skippedRows.length > 0) {
      results.info.push(`Skipped ${skippedRows.length} empty row(s): ${skippedRows.join(', ')}`)
    }

    if (validFeatures.length === 0) {
      results.errors.push('No valid features found in the file')
    } else {
      results.info.push(`Found ${validFeatures.length} valid feature(s)`)
    }

    // Check for duplicates
    const names = validFeatures.map(f => f.name.toLowerCase())
    const duplicates = names.filter((name, index) => names.indexOf(name) !== index)
    if (duplicates.length > 0) {
      results.warnings.push(`Found duplicate feature names: ${[...new Set(duplicates)].join(', ')}`)
    }

    results.isValid = results.errors.length === 0 && validFeatures.length > 0
    results.features = validFeatures

    return results
  }

  const resetUpload = () => {
    setUploadStatus('idle')
    setValidationResults(null)
    setErrorMessage('')
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const getStatusIcon = (type) => {
    switch (type) {
      case 'error': return <Error color="error" />
      case 'warning': return <Warning color="warning" />
      case 'info': return <Info color="info" />
      default: return <CheckCircle color="success" />
    }
  }

  return (
    <Box>
      <input
        ref={fileInputRef}
        type="file"
        accept={accept}
        onChange={handleFileChange}
        style={{ display: 'none' }}
        disabled={disabled}
      />
      
      <Paper
        onClick={handleClick}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        sx={{
          p: 3,
          border: '2px dashed',
          borderColor: isDragActive ? 'primary.main' : 'grey.300',
          backgroundColor: isDragActive ? 'primary.50' : 'background.paper',
          cursor: disabled ? 'default' : 'pointer',
          textAlign: 'center',
          transition: 'all 0.2s ease-in-out',
          '&:hover': disabled ? {} : {
            borderColor: 'primary.main',
            backgroundColor: 'primary.50',
          },
        }}
      >
        
        <CloudUpload sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
        
        <Typography variant="h6" gutterBottom>
          {isDragActive ? 'Drop CSV file here' : 'Upload Features CSV'}
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Drag and drop a CSV file here, or click to select
        </Typography>
        
        <Typography variant="caption" color="text.secondary">
          Maximum file size: {maxSize / (1024 * 1024)}MB • Supported format: CSV
        </Typography>
      </Paper>

      {uploadStatus === 'uploading' && (
        <Box sx={{ mt: 2 }}>
          <LinearProgress />
          <Typography variant="body2" sx={{ mt: 1, textAlign: 'center' }}>
            Processing file...
          </Typography>
        </Box>
      )}

      {uploadStatus === 'error' && errorMessage && (
        <Alert severity="error" sx={{ mt: 2 }} onClose={resetUpload}>
          {errorMessage}
        </Alert>
      )}

      {validationResults && (
        <Box sx={{ mt: 2 }}>
          {validationResults.isValid ? (
            <Alert severity="success">
              File processed successfully! {validationResults.features.length} features ready to import.
            </Alert>
          ) : (
            <Alert severity="error">
              File validation failed. Please fix the errors and try again.
            </Alert>
          )}

          {(validationResults.errors.length > 0 || 
            validationResults.warnings.length > 0 || 
            validationResults.info.length > 0) && (
            <Paper sx={{ mt: 2, maxHeight: 300, overflow: 'auto' }}>
              <List dense>
                {validationResults.errors.map((error, index) => (
                  <ListItem key={`error-${index}`}>
                    <ListItemIcon>{getStatusIcon('error')}</ListItemIcon>
                    <ListItemText primary={error} />
                  </ListItem>
                ))}
                
                {validationResults.warnings.map((warning, index) => (
                  <ListItem key={`warning-${index}`}>
                    <ListItemIcon>{getStatusIcon('warning')}</ListItemIcon>
                    <ListItemText primary={warning} />
                  </ListItem>
                ))}
                
                {validationResults.info.map((info, index) => (
                  <ListItem key={`info-${index}`}>
                    <ListItemIcon>{getStatusIcon('info')}</ListItemIcon>
                    <ListItemText primary={info} />
                  </ListItem>
                ))}
              </List>
            </Paper>
          )}
        </Box>
      )}

      {uploadStatus === 'success' && (
        <Box sx={{ mt: 2, display: 'flex', gap: 2 }}>
          <Button variant="outlined" onClick={resetUpload}>
            Upload Another File
          </Button>
        </Box>
      )}
    </Box>
  )
}

export default FileUpload