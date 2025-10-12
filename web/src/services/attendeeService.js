import api, { handleApiError, handleApiResponse } from './api'

// Attendee Service for managing project attendees
export const attendeeService = {
  // Get all attendees for a project
  getAttendees: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/attendees`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get a specific attendee
  getAttendee: async (projectId, attendeeId) => {
    try {
      const response = await api.get(`/projects/${projectId}/attendees/${attendeeId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Create a new attendee
  createAttendee: async (projectId, attendeeData) => {
    try {
      const response = await api.post(`/projects/${projectId}/attendees`, attendeeData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Update an existing attendee
  updateAttendee: async (projectId, attendeeId, attendeeData) => {
    try {
      const response = await api.put(`/projects/${projectId}/attendees/${attendeeId}`, attendeeData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Delete an attendee
  deleteAttendee: async (projectId, attendeeId) => {
    try {
      await api.delete(`/projects/${projectId}/attendees/${attendeeId}`)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Import attendees from CSV data
  importAttendees: async (projectId, csvData) => {
    try {
      const response = await api.post(`/projects/${projectId}/attendees/import`, {
        csv_data: csvData
      })
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Export attendees as CSV
  exportAttendees: async (projectId) => {
    try {
      const response = await api.get(`/projects/${projectId}/attendees/export`, {
        responseType: 'blob'
      })
      return response.data
    } catch (error) {
      handleApiError(error)
    }
  }
}

export default attendeeService