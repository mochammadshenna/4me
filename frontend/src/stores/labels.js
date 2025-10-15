import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useLabelsStore = defineStore('labels', () => {
  const labels = ref([])
  const loading = ref(false)

  async function fetchLabels(projectId) {
    loading.value = true
    try {
      const response = await apiClient.get(`/projects/${projectId}/labels`)
      labels.value = response.data
    } catch (error) {
      console.error('Failed to fetch labels:', error)
    } finally {
      loading.value = false
    }
  }

  async function createLabel(projectId, labelData) {
    try {
      const response = await apiClient.post(`/projects/${projectId}/labels`, labelData)
      labels.value.push(response.data)
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to create label' }
    }
  }

  async function updateLabel(id, labelData) {
    try {
      const response = await apiClient.put(`/labels/${id}`, labelData)
      const index = labels.value.findIndex(l => l.id === id)
      if (index !== -1) {
        labels.value[index] = response.data
      }
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to update label' }
    }
  }

  async function deleteLabel(id) {
    try {
      await apiClient.delete(`/labels/${id}`)
      labels.value = labels.value.filter(l => l.id !== id)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete label' }
    }
  }

  return {
    labels,
    loading,
    fetchLabels,
    createLabel,
    updateLabel,
    deleteLabel,
  }
})

