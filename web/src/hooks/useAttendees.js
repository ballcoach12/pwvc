import { useEffect, useState } from 'react'
import { attendeeService } from '../services/attendeeService'

export const useAttendees = (projectId) => {
  const [attendees, setAttendees] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    const fetchAttendees = async () => {
      if (!projectId) {
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        setError(null)
        const data = await attendeeService.getAttendees(projectId)
        setAttendees(data || [])
      } catch (err) {
        setError(err)
        setAttendees([])
      } finally {
        setLoading(false)
      }
    }

    fetchAttendees()
  }, [projectId])

  const addAttendee = async (attendeeData) => {
    try {
      const newAttendee = await attendeeService.createAttendee(projectId, attendeeData)
      setAttendees(prev => [...prev, newAttendee])
      return newAttendee
    } catch (err) {
      setError(err)
      throw err
    }
  }

  const updateAttendee = async (attendeeId, attendeeData) => {
    try {
      const updatedAttendee = await attendeeService.updateAttendee(projectId, attendeeId, attendeeData)
      setAttendees(prev => prev.map(a => a.id === attendeeId ? updatedAttendee : a))
      return updatedAttendee
    } catch (err) {
      setError(err)
      throw err
    }
  }

  const removeAttendee = async (attendeeId) => {
    try {
      await attendeeService.deleteAttendee(projectId, attendeeId)
      setAttendees(prev => prev.filter(a => a.id !== attendeeId))
    } catch (err) {
      setError(err)
      throw err
    }
  }

  const refreshAttendees = async () => {
    if (!projectId) return

    try {
      setLoading(true)
      const data = await attendeeService.getAttendees(projectId)
      setAttendees(data || [])
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  return {
    attendees,
    loading,
    error,
    addAttendee,
    updateAttendee,
    removeAttendee,
    refreshAttendees
  }
}

export default useAttendees