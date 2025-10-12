/**
 * Parse CSV text into array of objects
 * This is a simple CSV parser - for production use, consider a library like PapaParse
 */
export const parseCSV = (csvText, options = {}) => {
  const {
    delimiter = ',',
    quote = '"',
    skipEmptyLines = true,
    trimValues = true,
  } = options

  if (!csvText || typeof csvText !== 'string') {
    throw new Error('Invalid CSV data')
  }

  const lines = csvText.split('\n')
  if (lines.length === 0) {
    return { headers: [], data: [] }
  }

  // Parse headers
  const headers = parseLine(lines[0], delimiter, quote, trimValues)
  if (headers.length === 0) {
    throw new Error('No headers found in CSV')
  }

  // Parse data rows
  const data = []
  for (let i = 1; i < lines.length; i++) {
    const line = lines[i]
    
    // Skip empty lines if requested
    if (skipEmptyLines && !line.trim()) {
      continue
    }

    const values = parseLine(line, delimiter, quote, trimValues)
    
    // Skip completely empty rows
    if (values.length === 0 || values.every(v => !v)) {
      continue
    }

    // Create object with headers as keys
    const rowObject = {}
    headers.forEach((header, index) => {
      rowObject[header] = values[index] || ''
    })
    
    data.push(rowObject)
  }

  return { headers, data }
}

/**
 * Parse a single CSV line, handling quoted values and escaped quotes
 */
const parseLine = (line, delimiter, quote, trimValues) => {
  const values = []
  let current = ''
  let inQuotes = false
  let i = 0

  while (i < line.length) {
    const char = line[i]
    const nextChar = line[i + 1]

    if (char === quote) {
      if (inQuotes && nextChar === quote) {
        // Escaped quote
        current += quote
        i += 2
      } else {
        // Toggle quote state
        inQuotes = !inQuotes
        i++
      }
    } else if (char === delimiter && !inQuotes) {
      // End of field
      values.push(trimValues ? current.trim() : current)
      current = ''
      i++
    } else {
      current += char
      i++
    }
  }

  // Add final field
  values.push(trimValues ? current.trim() : current)

  return values
}

/**
 * Convert array of objects to CSV string
 */
export const arrayToCSV = (data, options = {}) => {
  const {
    delimiter = ',',
    quote = '"',
    includeHeaders = true,
    customHeaders = null,
  } = options

  if (!Array.isArray(data) || data.length === 0) {
    return ''
  }

  // Get headers
  const headers = customHeaders || Object.keys(data[0])
  
  const csvLines = []

  // Add headers if requested
  if (includeHeaders) {
    csvLines.push(headers.map(h => formatValue(h, delimiter, quote)).join(delimiter))
  }

  // Add data rows
  data.forEach(row => {
    const values = headers.map(header => {
      const value = row[header]
      return formatValue(value, delimiter, quote)
    })
    csvLines.push(values.join(delimiter))
  })

  return csvLines.join('\n')
}

/**
 * Format a value for CSV output, adding quotes if necessary
 */
const formatValue = (value, delimiter, quote) => {
  if (value == null) {
    return ''
  }

  const stringValue = String(value)
  
  // Quote if value contains delimiter, quote, or newline
  if (stringValue.includes(delimiter) || stringValue.includes(quote) || stringValue.includes('\n')) {
    // Escape quotes by doubling them
    const escapedValue = stringValue.replace(new RegExp(quote, 'g'), quote + quote)
    return quote + escapedValue + quote
  }

  return stringValue
}

/**
 * Validate feature CSV data
 */
export const validateFeaturesCSV = (csvData) => {
  const results = {
    isValid: false,
    features: [],
    errors: [],
    warnings: [],
    info: [],
  }

  const { headers, data } = csvData

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
  const seenNames = new Set()

  data.forEach((row, index) => {
    const rowNumber = index + 2 // +2 because CSV starts from row 1 and we skip header

    const name = headers[nameIndex] ? row[headers[nameIndex]]?.trim() : ''
    if (!name) {
      results.errors.push(`Row ${rowNumber}: Feature name is required`)
      return
    }

    if (name.length < 3) {
      results.errors.push(`Row ${rowNumber}: Feature name must be at least 3 characters`)
      return
    }

    if (name.length > 100) {
      results.errors.push(`Row ${rowNumber}: Feature name must be less than 100 characters`)
      return
    }

    // Check for duplicates
    const nameLower = name.toLowerCase()
    if (seenNames.has(nameLower)) {
      results.warnings.push(`Row ${rowNumber}: Duplicate feature name "${name}"`)
    } else {
      seenNames.add(nameLower)
    }

    const description = descriptionIndex >= 0 && headers[descriptionIndex] 
      ? row[headers[descriptionIndex]]?.trim() || '' 
      : ''

    if (description.length > 500) {
      results.errors.push(`Row ${rowNumber}: Description must be less than 500 characters`)
      return
    }

    validFeatures.push({
      name,
      description,
      rowNumber,
    })
  })

  // Add info and final validation
  if (validFeatures.length === 0) {
    results.errors.push('No valid features found in the file')
  } else {
    results.info.push(`Found ${validFeatures.length} valid feature(s)`)
  }

  if (skippedRows.length > 0) {
    results.info.push(`Skipped ${skippedRows.length} empty row(s)`)
  }

  results.isValid = results.errors.length === 0 && validFeatures.length > 0
  results.features = validFeatures

  return results
}

/**
 * Download data as CSV file
 */
export const downloadCSV = (data, filename, options = {}) => {
  const csvContent = arrayToCSV(data, options)
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  
  // Create download link
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  
  link.setAttribute('href', url)
  link.setAttribute('download', filename.endsWith('.csv') ? filename : `${filename}.csv`)
  link.style.visibility = 'hidden'
  
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  // Clean up
  URL.revokeObjectURL(url)
}

/**
 * Read file as text
 */
export const readFileAsText = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    
    reader.onload = (event) => {
      resolve(event.target.result)
    }
    
    reader.onerror = (error) => {
      reject(new Error('Failed to read file: ' + error.message))
    }
    
    reader.readAsText(file)
  })
}