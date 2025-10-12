import { useEffect, useState } from 'react'
import { featureService } from '../services/featureService.js'

export const useFeatures = (projectId) => {
  const [features, setFeatures] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadFeatures = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await featureService.getFeatures(id)
      setFeatures(data || [])
      return data
    } catch (err) {
      setError(err.message || 'Failed to load features')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const addFeature = async (id = projectId, featureData) => {
    if (!id) return
    
    try {
      setError(null)
      const newFeature = await featureService.createFeature(id, featureData)
      setFeatures(prev => [...prev, newFeature])
      return newFeature
    } catch (err) {
      setError(err.message || 'Failed to add feature')
      throw err
    }
  }

  const updateFeature = async (id = projectId, featureId, updates) => {
    if (!id || !featureId) return
    
    try {
      setError(null)
      const updatedFeature = await featureService.updateFeature(id, featureId, updates)
      setFeatures(prev => prev.map(f => f.id === featureId ? updatedFeature : f))
      return updatedFeature
    } catch (err) {
      setError(err.message || 'Failed to update feature')
      throw err
    }
  }

  const removeFeature = async (id = projectId, featureId) => {
    if (!id || !featureId) return
    
    try {
      setError(null)
      await featureService.deleteFeature(id, featureId)
      setFeatures(prev => prev.filter(f => f.id !== featureId))
    } catch (err) {
      setError(err.message || 'Failed to remove feature')
      throw err
    }
  }

  const importFeatures = async (id = projectId, featuresData) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const result = await featureService.importFeatures(id, featuresData)
      setFeatures(prev => [...prev, ...result.features])
      return result
    } catch (err) {
      setError(err.message || 'Failed to import features')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const exportFeatures = async (id = projectId) => {
    if (!id) return
    
    try {
      setError(null)
      const blob = await featureService.exportFeatures(id)
      return blob
    } catch (err) {
      setError(err.message || 'Failed to export features')
      throw err
    }
  }

  useEffect(() => {
    if (projectId) {
      loadFeatures(projectId)
    }
  }, [projectId])

  return {
    features,
    loading,
    error,
    loadFeatures,
    addFeature,
    updateFeature,
    removeFeature,
    importFeatures,
    exportFeatures,
    setFeatures,
    setError,
  }
}

export const usePairwiseComparisons = (projectId) => {
  const [comparisons, setComparisons] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadComparisons = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await featureService.getPairwiseComparisons(id)
      setComparisons(data || [])
      return data
    } catch (err) {
      setError(err.message || 'Failed to load comparisons')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const submitComparison = async (id = projectId, comparisonData) => {
    if (!id) return
    
    try {
      setError(null)
      const result = await featureService.submitPairwiseComparison(id, comparisonData)
      setComparisons(prev => [...prev, result])
      return result
    } catch (err) {
      setError(err.message || 'Failed to submit comparison')
      throw err
    }
  }

  useEffect(() => {
    if (projectId) {
      loadComparisons(projectId)
    }
  }, [projectId])

  return {
    comparisons,
    loading,
    error,
    loadComparisons,
    submitComparison,
    setComparisons,
    setError,
  }
}

export const useCalculationResults = (projectId) => {
  const [results, setResults] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadResults = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await featureService.getCalculationResults(id)
      setResults(data)
      return data
    } catch (err) {
      setError(err.message || 'Failed to load results')
      throw err
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (projectId) {
      loadResults(projectId)
    }
  }, [projectId])

  return {
    results,
    loading,
    error,
    loadResults,
    setResults,
    setError,
  }
}