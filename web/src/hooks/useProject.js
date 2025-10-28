import { useEffect, useState } from 'react'
import { projectService } from '../services/projectService.js'

export const useProject = (projectId) => {
  const [project, setProject] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadProject = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await projectService.getProject(id)
      setProject(data)
      return data
    } catch (err) {
      setError(err.message || 'Failed to load project')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const updateProject = async (id = projectId, updates) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await projectService.updateProject(id, updates)
      setProject(data)
      return data
    } catch (err) {
      setError(err.message || 'Failed to update project')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const deleteProject = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      await projectService.deleteProject(id)
      setProject(null)
    } catch (err) {
      setError(err.message || 'Failed to delete project')
      throw err
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (projectId) {
      loadProject(projectId)
    }
  }, [projectId])

  return {
    project,
    loading,
    error,
    loadProject,
    updateProject,
    deleteProject,
    setProject,
    setError,
  }
}

export const useProjects = () => {
  const [projects, setProjects] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadProjects = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await projectService.getProjects()
      setProjects(data || [])
      return data
    } catch (err) {
      setError(err.message || 'Failed to load projects')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const createProject = async (projectData) => {
    try {
      setLoading(true)
      setError(null)
      const newProject = await projectService.createProject(projectData)
      setProjects(prev => [newProject, ...prev])
      return newProject
    } catch (err) {
      setError(err.message || 'Failed to create project')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const removeProject = (projectId) => {
    setProjects(prev => prev.filter(p => p.id !== projectId))
  }

  useEffect(() => {
    loadProjects()
  }, [])

  return {
    projects,
    loading,
    error,
    loadProjects,
    createProject,
    removeProject,
    setProjects,
    setError,
  }
}

export const useAttendees = (projectId) => {
  const [attendees, setAttendees] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const loadAttendees = async (id = projectId) => {
    if (!id) return
    
    try {
      setLoading(true)
      setError(null)
      const data = await projectService.getProjectAttendees(id)
      setAttendees(data || [])
      return data
    } catch (err) {
      setError(err.message || 'Failed to load attendees')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const addAttendee = async (id = projectId, attendeeData) => {
    if (!id) return
    
    try {
      setError(null)
      const newAttendee = await projectService.addAttendee(id, attendeeData)
      setAttendees(prev => [...prev, newAttendee])
      return newAttendee
    } catch (err) {
      setError(err.message || 'Failed to add attendee')
      throw err
    }
  }

  const updateAttendee = async (id = projectId, attendeeId, updates) => {
    if (!id || !attendeeId) return
    
    try {
      setError(null)
      const updatedAttendee = await projectService.updateAttendee(id, attendeeId, updates)
      setAttendees(prev => prev.map(a => a.id === attendeeId ? updatedAttendee : a))
      return updatedAttendee
    } catch (err) {
      setError(err.message || 'Failed to update attendee')
      throw err
    }
  }

  const removeAttendee = async (id = projectId, attendeeId) => {
    if (!id || !attendeeId) return
    
    try {
      setError(null)
      await projectService.removeAttendee(id, attendeeId)
      setAttendees(prev => prev.filter(a => a.id !== attendeeId))
    } catch (err) {
      setError(err.message || 'Failed to remove attendee')
      throw err
    }
  }

  const setFacilitator = async (id = projectId, attendeeId) => {
    if (!id || !attendeeId) return
    
    try {
      setError(null)
      await projectService.setFacilitator(id, attendeeId)
      setAttendees(prev => prev.map(a => ({
        ...a,
        isFacilitator: a.id === attendeeId
      })))
    } catch (err) {
      setError(err.message || 'Failed to set facilitator')
      throw err
    }
  }

  useEffect(() => {
    if (projectId) {
      loadAttendees(projectId)
    }
  }, [projectId])

  return {
    attendees,
    loading,
    error,
    loadAttendees,
    addAttendee,
    updateAttendee,
    removeAttendee,
    setFacilitator,
    setAttendees,
    setError,
  }
}