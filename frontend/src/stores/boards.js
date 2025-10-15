import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '@/api/client'

export const useBoardsStore = defineStore('boards', () => {
  const boards = ref([])
  const loading = ref(false)

  async function fetchBoards(projectId) {
    loading.value = true
    try {
      const response = await apiClient.get(`/projects/${projectId}/boards`)
      boards.value = response.data
    } catch (error) {
      console.error('Failed to fetch boards:', error)
    } finally {
      loading.value = false
    }
  }

  async function createBoard(projectId, boardData) {
    try {
      const response = await apiClient.post(`/projects/${projectId}/boards`, boardData)
      boards.value.push(response.data)
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to create board' }
    }
  }

  async function updateBoard(id, boardData) {
    try {
      const response = await apiClient.put(`/boards/${id}`, boardData)
      const index = boards.value.findIndex(b => b.id === id)
      if (index !== -1) {
        boards.value[index] = response.data
      }
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to update board' }
    }
  }

  async function deleteBoard(id) {
    try {
      await apiClient.delete(`/boards/${id}`)
      boards.value = boards.value.filter(b => b.id !== id)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete board' }
    }
  }

  return {
    boards,
    loading,
    fetchBoards,
    createBoard,
    updateBoard,
    deleteBoard,
  }
})

