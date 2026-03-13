import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useErrors } from '../../contexts/ErrorContext'
import { api } from '../../services/api'
import './WorkflowNavigation.css'

// Workflow phase definitions
const WORKFLOW_PHASES = [
  {
    key: 'setup',
    name: 'Project Setup',
    description: 'Create project and basic configuration',
    path: '/edit',
    icon: '⚙️',
    requirements: []
  },
  {
    key: 'attendees',
    name: 'Add Attendees',
    description: 'Add team members who will participate',
    path: '/attendees',
    icon: '👥',
    requirements: ['Project must be created']
  },
  {
    key: 'features',
    name: 'Add Features',
    description: 'Define features to be prioritized',
    path: '/features',
    icon: '📋',
    requirements: ['At least 2 attendees required']
  },
  {
    key: 'pairwise_value',
    name: 'Value Comparisons',
    description: 'Compare features by business value',
    path: '/comparison',
    icon: '⚖️',
    requirements: ['At least 2 features required']
  },
  {
    key: 'pairwise_complexity',
    name: 'Complexity Comparisons',
    description: 'Compare features by implementation complexity',
    path: '/comparison',
    icon: '🔧',
    requirements: ['Value comparisons completed']
  },
  {
    key: 'fibonacci_value',
    name: 'Value Scoring',
    description: 'Assign Fibonacci scores for value',
    path: '/scoring/value',
    icon: '🎯',
    requirements: ['Complexity comparisons completed']
  },
  {
    key: 'fibonacci_complexity',
    name: 'Complexity Scoring',
    description: 'Assign Fibonacci scores for complexity',
    path: '/scoring/complexity',
    icon: '📊',
    requirements: ['Value scoring completed']
  },
  {
    key: 'results',
    name: 'View Results',
    description: 'See final priority rankings',
    path: '/results',
    icon: '🏆',
    requirements: ['All scoring completed']
  }
]

export function WorkflowNavigation({ currentPhase, onPhaseChange }) {
  const { id: projectId } = useParams()
  const navigate = useNavigate()
  const { handleApiError } = useErrors()
  const [progress, setProgress] = useState(null)
  const [availablePhases, setAvailablePhases] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadProgressData()
  }, [projectId])

  const loadProgressData = async () => {
    try {
      setLoading(true)
      const [progressData, availablePhasesData] = await Promise.all([
        api.projects.getProgress(projectId),
        api.projects.getAvailablePhases(projectId)
      ])
      
      setProgress(progressData)
      setAvailablePhases(availablePhasesData.available_phases || [])
    } catch (error) {
      handleApiError(error, { context: 'Loading workflow progress' })
    } finally {
      setLoading(false)
    }
  }

  const handlePhaseClick = async (phase) => {
    if (!canNavigateToPhase(phase)) {
      return
    }

    try {
      // Advance to the selected phase if needed
      if (progress && progress.current_phase !== phase.key) {
        await api.projects.advancePhase(projectId, phase.key)
        await loadProgressData()
      }

      // Navigate to the phase
      const path = `/projects/${projectId}${phase.path}`
      navigate(path)
      
      if (onPhaseChange) {
        onPhaseChange(phase.key)
      }
    } catch (error) {
      handleApiError(error, { context: 'Navigating to phase' })
    }
  }

  const canNavigateToPhase = (phase) => {
    return availablePhases.includes(phase.key)
  }

  const getPhaseStatus = (phase) => {
    if (!progress) return 'pending'

    switch (phase.key) {
      case 'setup':
        return progress.setup_completed ? 'completed' : 'current'
      case 'attendees':
        return progress.attendees_added ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'features':
        return progress.features_added ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'pairwise_value':
        return progress.pairwise_value_completed ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'pairwise_complexity':
        return progress.pairwise_complexity_completed ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'fibonacci_value':
        return progress.fibonacci_value_completed ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'fibonacci_complexity':
        return progress.fibonacci_complexity_completed ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      case 'results':
        return progress.results_calculated ? 'completed' : 
               availablePhases.includes(phase.key) ? 'available' : 'pending'
      default:
        return 'pending'
    }
  }

  const getCurrentPhaseIndex = () => {
    if (!progress) return 0
    return WORKFLOW_PHASES.findIndex(phase => phase.key === progress.current_phase)
  }

  const getCompletionPercentage = () => {
    if (!progress) return 0
    
    const completedSteps = [
      progress.setup_completed,
      progress.attendees_added,
      progress.features_added,
      progress.pairwise_value_completed,
      progress.pairwise_complexity_completed,
      progress.fibonacci_value_completed,
      progress.fibonacci_complexity_completed,
      progress.results_calculated
    ].filter(Boolean).length

    return Math.round((completedSteps / WORKFLOW_PHASES.length) * 100)
  }

  if (loading) {
    return (
      <div className="workflow-navigation loading">
        <div className="workflow-header">
          <h2>PairWise Workflow</h2>
          <div className="loading-spinner">Loading...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="workflow-navigation">
      <div className="workflow-header">
        <h2>PairWise Workflow</h2>
        <div className="progress-overview">
          <div className="progress-bar">
            <div 
              className="progress-fill" 
              style={{ width: `${getCompletionPercentage()}%` }}
            />
          </div>
          <span className="progress-text">
            {getCompletionPercentage()}% Complete
          </span>
        </div>
      </div>

      <div className="workflow-steps">
        {WORKFLOW_PHASES.map((phase, index) => {
          const status = getPhaseStatus(phase)
          const isClickable = canNavigateToPhase(phase)
          const isCurrent = progress?.current_phase === phase.key

          return (
            <div
              key={phase.key}
              className={`workflow-step ${status} ${isCurrent ? 'current' : ''} ${isClickable ? 'clickable' : ''}`}
              onClick={() => isClickable && handlePhaseClick(phase)}
            >
              <div className="step-indicator">
                <span className="step-icon">{phase.icon}</span>
                <span className="step-number">{index + 1}</span>
                {status === 'completed' && <span className="completion-check">✓</span>}
              </div>

              <div className="step-content">
                <h3 className="step-name">{phase.name}</h3>
                <p className="step-description">{phase.description}</p>
                
                {status === 'pending' && phase.requirements.length > 0 && (
                  <div className="step-requirements">
                    <small>Requirements:</small>
                    <ul>
                      {phase.requirements.map((req, i) => (
                        <li key={i}>{req}</li>
                      ))}
                    </ul>
                  </div>
                )}

                {isCurrent && (
                  <div className="current-indicator">
                    <span>Current Step</span>
                  </div>
                )}
              </div>

              {index < WORKFLOW_PHASES.length - 1 && (
                <div className={`step-connector ${status === 'completed' ? 'completed' : ''}`} />
              )}
            </div>
          )
        })}
      </div>

      <div className="workflow-actions">
        {progress && getCurrentPhaseIndex() > 0 && (
          <button
            className="btn btn-secondary"
            onClick={() => {
              const prevPhase = WORKFLOW_PHASES[getCurrentPhaseIndex() - 1]
              handlePhaseClick(prevPhase)
            }}
          >
            ← Previous Step
          </button>
        )}

        {progress && getCurrentPhaseIndex() < WORKFLOW_PHASES.length - 1 && (
          <button
            className="btn btn-primary"
            onClick={() => {
              const nextPhase = WORKFLOW_PHASES[getCurrentPhaseIndex() + 1]
              if (canNavigateToPhase(nextPhase)) {
                handlePhaseClick(nextPhase)
              }
            }}
            disabled={!canNavigateToPhase(WORKFLOW_PHASES[getCurrentPhaseIndex() + 1])}
          >
            Next Step →
          </button>
        )}
      </div>
    </div>
  )
}