import { useCallback, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useErrors } from '../contexts/ErrorContext'
import { api } from '../services/api'

// Hook for managing phase completion and transitions
export function usePhaseCompletion() {
  const { id: projectId } = useParams()
  const navigate = useNavigate()
  const { handleApiError, addError } = useErrors()
  const [completing, setCompleting] = useState(false)
  const [loading, setLoading] = useState(false)

  // Complete current phase and advance to next
  const completePhase = useCallback(async (phase, options = {}) => {
    const { 
      skipValidation = false, 
      navigateToNext = true,
      onSuccess = null,
      onError = null 
    } = options

    try {
      setCompleting(true)

      // Validate phase completion requirements if not skipped
      if (!skipValidation) {
        const isValid = await validatePhaseCompletion(phase)
        if (!isValid) {
          throw new Error('Phase completion requirements not met')
        }
      }

      // Mark phase as completed
      await api.projects.completePhase(projectId, phase)

      // Success notification
      addError({
        message: `${getPhaseDisplayName(phase)} completed successfully!`,
        type: 'success',
        severity: 'low'
      })

      // Navigate to next phase if requested
      if (navigateToNext) {
        await navigateToNextPhase(phase)
      }

      if (onSuccess) {
        onSuccess()
      }

    } catch (error) {
      handleApiError(error, { 
        context: `Completing phase: ${phase}` 
      })
      
      if (onError) {
        onError(error)
      }
    } finally {
      setCompleting(false)
    }
  }, [projectId, handleApiError, addError])

  // Validate if a phase can be completed
  const validatePhaseCompletion = useCallback(async (phase) => {
    try {
      setLoading(true)
      
      switch (phase) {
        case 'setup':
          // Validate project exists and has required data
          const project = await api.projects.getById(projectId)
          return project && project.name && project.name.trim().length > 0

        case 'attendees':
          // Validate at least 2 attendees exist
          const attendees = await api.attendees.getByProject(projectId)
          if (attendees.length < 2) {
            addError({
              message: 'At least 2 attendees are required for pairwise comparisons',
              type: 'validation',
              severity: 'medium'
            })
            return false
          }
          return true

        case 'features':
          // Validate at least 2 features exist
          const features = await api.features.getByProject(projectId)
          if (features.length < 2) {
            addError({
              message: 'At least 2 features are required for comparisons',
              type: 'validation',
              severity: 'medium'
            })
            return false
          }
          return true

        case 'pairwise_value':
          // Validate all value comparisons are completed
          const valueSession = await api.pairwise.getSession(projectId)
          if (!valueSession || valueSession.status !== 'completed') {
            addError({
              message: 'All value comparisons must be completed before proceeding',
              type: 'validation',
              severity: 'medium'
            })
            return false
          }
          return true

        case 'pairwise_complexity':
          // Validate all complexity comparisons are completed
          const complexitySession = await api.pairwise.getSession(projectId)
          if (!complexitySession || complexitySession.status !== 'completed') {
            addError({
              message: 'All complexity comparisons must be completed before proceeding',
              type: 'validation',
              severity: 'medium'
            })
            return false
          }
          return true

        case 'fibonacci_value':
          // Validate all value scores are assigned
          // This would need to check scoring completion status
          return true

        case 'fibonacci_complexity':
          // Validate all complexity scores are assigned
          // This would need to check scoring completion status
          return true

        case 'results':
          // Validate all previous phases are complete
          const progress = await api.projects.getProgress(projectId)
          return progress && 
                 progress.fibonacci_value_completed && 
                 progress.fibonacci_complexity_completed

        default:
          return true
      }
    } catch (error) {
      handleApiError(error, { context: `Validating phase: ${phase}` })
      return false
    } finally {
      setLoading(false)
    }
  }, [projectId, handleApiError, addError])

  // Navigate to the next phase in the workflow
  const navigateToNextPhase = useCallback(async (currentPhase) => {
    const phaseOrder = [
      'setup',
      'attendees', 
      'features',
      'pairwise_value',
      'pairwise_complexity', 
      'fibonacci_value',
      'fibonacci_complexity',
      'results'
    ]

    const currentIndex = phaseOrder.indexOf(currentPhase)
    if (currentIndex >= 0 && currentIndex < phaseOrder.length - 1) {
      const nextPhase = phaseOrder[currentIndex + 1]
      const nextPath = getPhaseNavigationPath(nextPhase)
      
      if (nextPath) {
        navigate(`/projects/${projectId}${nextPath}`)
      }
    }
  }, [projectId, navigate])

  // Force navigate to a specific phase
  const navigateToPhase = useCallback(async (phase) => {
    try {
      // Check if phase is available
      const availablePhases = await api.projects.getAvailablePhases(projectId)
      
      if (!availablePhases.available_phases?.includes(phase)) {
        addError({
          message: 'This phase is not yet available',
          type: 'business',
          severity: 'medium'
        })
        return false
      }

      // Advance to phase
      await api.projects.advancePhase(projectId, phase)
      
      // Navigate
      const path = getPhaseNavigationPath(phase)
      if (path) {
        navigate(`/projects/${projectId}${path}`)
        return true
      }
      
      return false
    } catch (error) {
      handleApiError(error, { context: `Navigating to phase: ${phase}` })
      return false
    }
  }, [projectId, navigate, handleApiError, addError])

  return {
    completePhase,
    validatePhaseCompletion,
    navigateToPhase,
    navigateToNextPhase,
    completing,
    loading
  }
}

// Utility functions
function getPhaseDisplayName(phase) {
  const names = {
    setup: 'Project Setup',
    attendees: 'Attendee Management',
    features: 'Feature Management',
    pairwise_value: 'Value Comparisons',
    pairwise_complexity: 'Complexity Comparisons',
    fibonacci_value: 'Value Scoring',
    fibonacci_complexity: 'Complexity Scoring',
    results: 'Results'
  }
  return names[phase] || phase
}

function getPhaseNavigationPath(phase) {
  const paths = {
    setup: '/edit',
    attendees: '/attendees',
    features: '/features',
    pairwise_value: '/comparison?type=value',
    pairwise_complexity: '/comparison?type=complexity',
    fibonacci_value: '/scoring/value',
    fibonacci_complexity: '/scoring/complexity',
    results: '/results'
  }
  return paths[phase]
}

// Hook for checking if current user can perform phase operations
export function usePhasePermissions() {
  const { id: projectId } = useParams()
  const [permissions, setPermissions] = useState({
    canEdit: false,
    canAdvance: false,
    canComplete: false,
    isOwner: false
  })

  // In a real implementation, this would check user roles and permissions
  // For now, we'll assume all users have full permissions
  const checkPermissions = useCallback(async () => {
    // This would typically make an API call to check user permissions
    setPermissions({
      canEdit: true,
      canAdvance: true, 
      canComplete: true,
      isOwner: true
    })
  }, [projectId])

  return {
    permissions,
    checkPermissions,
    canEdit: permissions.canEdit,
    canAdvance: permissions.canAdvance,
    canComplete: permissions.canComplete,
    isOwner: permissions.isOwner
  }
}