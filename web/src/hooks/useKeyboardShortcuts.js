import { useEffect, useRef } from 'react'

/**
 * Custom hook for handling keyboard shortcuts throughout the application.
 * Provides a clean way to register and manage keyboard event handlers
 * with proper cleanup and event filtering.
 * 
 * @param {Object} shortcuts - Object mapping keys to handler functions
 * @param {Array} dependencies - Dependencies array for effect updates
 * @param {Object} options - Configuration options
 */
export const useKeyboardShortcuts = (shortcuts = {}, dependencies = [], options = {}) => {
  const {
    enabled = true,
    preventDefault = true,
    stopPropagation = false,
    ignoreInputs = true
  } = options

  const shortcutsRef = useRef(shortcuts)
  
  // Update shortcuts ref when shortcuts change
  useEffect(() => {
    shortcutsRef.current = shortcuts
  }, [shortcuts])

  useEffect(() => {
    if (!enabled) return

    const handleKeyDown = (event) => {
      // Skip if typing in input fields (unless disabled)
      if (ignoreInputs && (
        event.target.tagName === 'INPUT' ||
        event.target.tagName === 'TEXTAREA' ||
        event.target.contentEditable === 'true'
      )) {
        return
      }

      const key = event.key
      const shortcut = shortcutsRef.current[key]

      if (shortcut && typeof shortcut === 'function') {
        if (preventDefault) {
          event.preventDefault()
        }
        if (stopPropagation) {
          event.stopPropagation()
        }
        
        shortcut(event)
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [enabled, preventDefault, stopPropagation, ignoreInputs, ...dependencies])
}

/**
 * Predefined keyboard shortcuts for pairwise comparison
 */
export const PAIRWISE_SHORTCUTS = {
  VOTE_A: ['1', 'ArrowLeft'],
  VOTE_NEUTRAL: [' ', '2'], // Space or 2
  VOTE_B: ['3', 'ArrowRight'],
  NEXT_COMPARISON: ['n', 'ArrowDown'],
  PREVIOUS_COMPARISON: ['p', 'ArrowUp'],
  TOGGLE_VIEW: ['g'],
  TOGGLE_FULLSCREEN: ['f'],
  HELP: ['?', 'h'],
  ESCAPE: ['Escape']
}

/**
 * Hook specifically for pairwise comparison shortcuts
 */
export const usePairwiseShortcuts = ({
  onVoteA,
  onVoteNeutral,
  onVoteB,
  onNext,
  onPrevious,
  onToggleView,
  onToggleFullscreen,
  onHelp,
  onEscape,
  enabled = true
}) => {
  const shortcuts = {}

  // Vote shortcuts
  if (onVoteA) {
    PAIRWISE_SHORTCUTS.VOTE_A.forEach(key => {
      shortcuts[key] = onVoteA
    })
  }

  if (onVoteNeutral) {
    PAIRWISE_SHORTCUTS.VOTE_NEUTRAL.forEach(key => {
      shortcuts[key] = onVoteNeutral
    })
  }

  if (onVoteB) {
    PAIRWISE_SHORTCUTS.VOTE_B.forEach(key => {
      shortcuts[key] = onVoteB
    })
  }

  // Navigation shortcuts
  if (onNext) {
    PAIRWISE_SHORTCUTS.NEXT_COMPARISON.forEach(key => {
      shortcuts[key] = onNext
    })
  }

  if (onPrevious) {
    PAIRWISE_SHORTCUTS.PREVIOUS_COMPARISON.forEach(key => {
      shortcuts[key] = onPrevious
    })
  }

  // View shortcuts
  if (onToggleView) {
    PAIRWISE_SHORTCUTS.TOGGLE_VIEW.forEach(key => {
      shortcuts[key] = onToggleView
    })
  }

  if (onToggleFullscreen) {
    PAIRWISE_SHORTCUTS.TOGGLE_FULLSCREEN.forEach(key => {
      shortcuts[key] = onToggleFullscreen
    })
  }

  // Utility shortcuts
  if (onHelp) {
    PAIRWISE_SHORTCUTS.HELP.forEach(key => {
      shortcuts[key] = onHelp
    })
  }

  if (onEscape) {
    PAIRWISE_SHORTCUTS.ESCAPE.forEach(key => {
      shortcuts[key] = onEscape
    })
  }

  useKeyboardShortcuts(shortcuts, [
    onVoteA, onVoteNeutral, onVoteB, 
    onNext, onPrevious, 
    onToggleView, onToggleFullscreen,
    onHelp, onEscape
  ], { enabled })
}

/**
 * Hook for displaying keyboard shortcut help
 */
export const useShortcutHelp = () => {
  const getShortcutHelp = () => {
    return [
      {
        category: 'Voting',
        shortcuts: [
          { keys: ['1', '←'], description: 'Vote for Feature A (left)' },
          { keys: ['Space', '2'], description: 'Vote neutral / Skip' },
          { keys: ['3', '→'], description: 'Vote for Feature B (right)' }
        ]
      },
      {
        category: 'Navigation',
        shortcuts: [
          { keys: ['n', '↓'], description: 'Next comparison' },
          { keys: ['p', '↑'], description: 'Previous comparison' },
          { keys: ['g'], description: 'Toggle grid/detail view' }
        ]
      },
      {
        category: 'View',
        shortcuts: [
          { keys: ['f'], description: 'Toggle fullscreen' },
          { keys: ['?', 'h'], description: 'Show this help' },
          { keys: ['Esc'], description: 'Exit fullscreen / Close dialogs' }
        ]
      }
    ]
  }

  return { getShortcutHelp }
}

export default useKeyboardShortcuts