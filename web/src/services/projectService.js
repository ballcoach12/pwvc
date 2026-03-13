import api, { handleApiError, handleApiResponse } from './api.js'

export const projectService = {
  // Get all projects
  async getProjects() {
    try {
      const response = await api.get('/projects')
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get project by ID
  async getProject(id) {
    try {
      const response = await api.get(`/projects/${id}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Create new project
  async createProject(projectData) {
    try {
      const response = await api.post('/projects', projectData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Update project
  async updateProject(id, projectData) {
    try {
      const response = await api.put(`/projects/${id}`, projectData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Delete project
  async deleteProject(id) {
    try {
      const response = await api.delete(`/projects/${id}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Get project attendees
  async getProjectAttendees(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/attendees`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Add attendee to project
  async addAttendee(projectId, attendeeData) {
    try {
      const response = await api.post(`/projects/${projectId}/attendees`, attendeeData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Update attendee
  async updateAttendee(projectId, attendeeId, attendeeData) {
    try {
      const response = await api.put(`/projects/${projectId}/attendees/${attendeeId}`, attendeeData)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Remove attendee from project
  async removeAttendee(projectId, attendeeId) {
    try {
      const response = await api.delete(`/projects/${projectId}/attendees/${attendeeId}`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },

  // Set facilitator
  async setFacilitator(projectId, attendeeId) {
    try {
      const response = await api.put(`/projects/${projectId}/attendees/${attendeeId}/facilitator`)
      return handleApiResponse(response)
    } catch (error) {
      handleApiError(error)
    }
  },
}