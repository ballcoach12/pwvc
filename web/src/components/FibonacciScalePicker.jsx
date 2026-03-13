import { Check } from '@mui/icons-material'
import {
    alpha,
    Box,
    Button,
    ButtonGroup,
    Chip,
    Paper,
    Tooltip,
    Typography,
    useTheme
} from '@mui/material'
import { useState } from 'react'

/**
 * FibonacciScalePicker provides a visual interface for selecting Fibonacci numbers
 * for scoring feature Value or Complexity. Includes tooltips with guidance and
 * validation to ensure only valid Fibonacci sequence values are selected.
 * 
 * Features:
 * - Visual Fibonacci sequence with explanatory tooltips
 * - Clear selection state with visual feedback
 * - Keyboard navigation support
 * - Size variants for different use cases
 * - Validation for valid Fibonacci numbers
 * - Customizable styling and orientation
 */

// Valid Fibonacci sequence values for PairWise scoring
export const FIBONACCI_VALUES = [1, 2, 3, 5, 8, 13, 21, 34, 55, 89]

// Explanatory tooltips for each Fibonacci value
const FIBONACCI_TOOLTIPS = {
  1: 'Trivial - Very small effort, minimal impact',
  2: 'Minor - Small effort, small impact', 
  3: 'Small - Moderate effort, noticeable impact',
  5: 'Medium - Significant effort, good impact',
  8: 'Large - Major effort, strong impact',
  13: 'Very Large - Substantial effort, major impact',
  21: 'Huge - Extensive effort, game-changing impact',
  34: 'Massive - Enormous effort, transformative impact',
  55: 'Epic - Extreme effort, revolutionary impact',
  89: 'Legendary - Unprecedented effort, paradigm-shifting impact'
}

const FibonacciScalePicker = ({
  value = null,
  onChange,
  disabled = false,
  size = 'medium', // 'small' | 'medium' | 'large'
  orientation = 'horizontal', // 'horizontal' | 'vertical'
  showLabels = true,
  criterionType = 'value', // 'value' | 'complexity'
  label,
  helperText,
  error = false,
  fullWidth = false
}) => {
  const theme = useTheme()
  const [hoveredValue, setHoveredValue] = useState(null)

  const handleValueSelect = (selectedValue) => {
    if (disabled) return
    
    if (onChange) {
      onChange(selectedValue === value ? null : selectedValue)
    }
  }

  const getButtonSize = () => {
    switch (size) {
      case 'small': return { minWidth: 32, height: 32, fontSize: '0.75rem' }
      case 'large': return { minWidth: 56, height: 56, fontSize: '1.1rem' }
      default: return { minWidth: 44, height: 44, fontSize: '0.9rem' }
    }
  }

  const getTooltipTitle = (fibValue) => {
    const baseTooltip = FIBONACCI_TOOLTIPS[fibValue]
    const criterionLabel = criterionType === 'value' ? 'Value' : 'Complexity'
    return `${fibValue} - ${baseTooltip}\n\n${criterionLabel} scoring on Fibonacci scale`
  }

  const buttonStyle = getButtonSize()

  const renderFibonacciButton = (fibValue) => {
    const isSelected = value === fibValue
    const isHovered = hoveredValue === fibValue

    return (
      <Tooltip
        key={fibValue}
        title={getTooltipTitle(fibValue)}
        placement="top"
        arrow
      >
        <Button
          variant={isSelected ? 'contained' : 'outlined'}
          color={isSelected ? 'primary' : 'inherit'}
          onClick={() => handleValueSelect(fibValue)}
          onMouseEnter={() => setHoveredValue(fibValue)}
          onMouseLeave={() => setHoveredValue(null)}
          disabled={disabled}
          sx={{
            ...buttonStyle,
            position: 'relative',
            borderColor: error ? 'error.main' : 
                       isSelected ? 'primary.main' : 
                       isHovered ? 'primary.light' : 'grey.300',
            backgroundColor: isSelected ? 'primary.main' :
                           isHovered ? alpha(theme.palette.primary.main, 0.08) :
                           'transparent',
            color: isSelected ? 'primary.contrastText' :
                   isHovered ? 'primary.main' :
                   'text.primary',
            '&:hover': {
              backgroundColor: isSelected ? 'primary.dark' :
                             alpha(theme.palette.primary.main, 0.12),
              borderColor: isSelected ? 'primary.dark' : 'primary.main'
            },
            '&:disabled': {
              backgroundColor: 'grey.100',
              color: 'text.disabled',
              borderColor: 'grey.300'
            },
            transition: 'all 0.2s ease-in-out'
          }}
        >
          {fibValue}
          {isSelected && (
            <Check 
              sx={{ 
                position: 'absolute',
                top: 2,
                right: 2,
                fontSize: size === 'small' ? 12 : 16
              }} 
            />
          )}
        </Button>
      </Tooltip>
    )
  }

  return (
    <Box sx={{ width: fullWidth ? '100%' : 'auto' }}>
      {label && (
        <Typography 
          variant="body2" 
          sx={{ 
            mb: 1, 
            color: error ? 'error.main' : 'text.primary',
            fontWeight: 'medium'
          }}
        >
          {label}
        </Typography>
      )}
      
      <Paper 
        variant="outlined" 
        sx={{ 
          p: 1.5,
          borderColor: error ? 'error.main' : 'grey.300',
          backgroundColor: disabled ? 'grey.50' : 'background.paper'
        }}
      >
        {showLabels && (
          <Box sx={{ mb: 1 }}>
            <Typography variant="caption" color="text.secondary">
              Select {criterionType === 'value' ? 'Value' : 'Complexity'} Score (Fibonacci Scale)
            </Typography>
          </Box>
        )}

        <ButtonGroup
          orientation={orientation}
          variant="outlined"
          sx={{
            flexWrap: orientation === 'horizontal' ? 'wrap' : 'nowrap',
            gap: 0.5,
            '& .MuiButtonGroup-grouped': {
              borderRadius: 1,
              border: '1px solid',
              '&:not(:last-of-type)': {
                borderRightColor: 'inherit'
              }
            }
          }}
        >
          {FIBONACCI_VALUES.map(renderFibonacciButton)}
        </ButtonGroup>

        {value && (
          <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
            <Chip
              label={`Selected: ${value}`}
              size="small"
              color="primary"
              variant="outlined"
            />
            <Typography variant="caption" color="text.secondary">
              {FIBONACCI_TOOLTIPS[value]}
            </Typography>
          </Box>
        )}
      </Paper>

      {helperText && (
        <Typography 
          variant="caption" 
          sx={{ 
            mt: 0.5, 
            display: 'block',
            color: error ? 'error.main' : 'text.secondary'
          }}
        >
          {helperText}
        </Typography>
      )}
    </Box>
  )
}

export default FibonacciScalePicker