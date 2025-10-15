import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useProjectsStore = defineStore('projects', () => {
  const projects = ref([])
  const currentProject = ref(null)
  const loading = ref(false)

  async function fetchProjects() {
    loading.value = true
    try {
      const response = await apiClient.get('/projects')
      projects.value = response.data
    } catch (error) {
      console.error('Failed to fetch projects:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchProject(id) {
    loading.value = true
    try {
      const response = await apiClient.get(`/projects/${id}`)
      currentProject.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch project:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  async function createProject(projectData) {
    try {
      // Convert Vue proxy to plain object to avoid enumeration warning
      const plainData = JSON.parse(JSON.stringify(projectData))
      const response = await apiClient.post('/projects', plainData)
      projects.value.unshift(response.data)
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to create project' }
    }
  }

  async function updateProject(id, projectData) {
    try {
      // Convert Vue proxy to plain object to avoid enumeration warning
      const plainData = JSON.parse(JSON.stringify(projectData))
      const response = await apiClient.put(`/projects/${id}`, plainData)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = response.data
      }
      if (currentProject.value?.id === id) {
        currentProject.value = response.data
      }
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to update project' }
    }
  }

  async function deleteProject(id) {
    try {
      await apiClient.delete(`/projects/${id}`)
      projects.value = projects.value.filter(p => p.id !== id)
      if (currentProject.value?.id === id) {
        currentProject.value = null
      }
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete project' }
    }
  }

  return {
    projects,
    currentProject,
    loading,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
  }
})

