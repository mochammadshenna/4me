import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref({}) // Organized by board ID
  const currentTask = ref(null)
  const loading = ref(false)

  const getTasksByBoard = computed(() => (boardId) => {
    return tasks.value[boardId] || []
  })

  function setTasks(boardId, taskList) {
    tasks.value[boardId] = taskList
  }

  async function createTask(boardId, taskData) {
    try {
      const response = await apiClient.post(`/boards/${boardId}/tasks`, taskData)
      if (!tasks.value[boardId]) {
        tasks.value[boardId] = []
      }
      tasks.value[boardId].push(response.data)
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to create task' }
    }
  }

  async function fetchTask(id) {
    loading.value = true
    try {
      const response = await apiClient.get(`/tasks/${id}`)
      currentTask.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch task:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateTask(id, taskData) {
    try {
      const response = await apiClient.put(`/tasks/${id}`, taskData)
      
      // Update task in the appropriate board
      Object.keys(tasks.value).forEach(boardId => {
        const index = tasks.value[boardId].findIndex(t => t.id === id)
        if (index !== -1) {
          tasks.value[boardId][index] = response.data
        }
      })

      if (currentTask.value?.id === id) {
        currentTask.value = response.data
      }

      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to update task' }
    }
  }

  async function moveTask(id, targetBoardId, position) {
    try {
      const response = await apiClient.patch(`/tasks/${id}/move`, {
        board_id: targetBoardId,
        position,
      })

      // Remove from old board and add to new board
      Object.keys(tasks.value).forEach(boardId => {
        tasks.value[boardId] = tasks.value[boardId].filter(t => t.id !== id)
      })

      if (!tasks.value[targetBoardId]) {
        tasks.value[targetBoardId] = []
      }
      tasks.value[targetBoardId].push(response.data)

      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to move task' }
    }
  }

  async function deleteTask(id) {
    try {
      await apiClient.delete(`/tasks/${id}`)
      
      // Remove from all boards
      Object.keys(tasks.value).forEach(boardId => {
        tasks.value[boardId] = tasks.value[boardId].filter(t => t.id !== id)
      })

      if (currentTask.value?.id === id) {
        currentTask.value = null
      }

      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete task' }
    }
  }

  async function fetchTaskHistory(id) {
    try {
      const response = await apiClient.get(`/tasks/${id}/history`)
      return response.data
    } catch (error) {
      console.error('Failed to fetch task history:', error)
      return []
    }
  }

  async function fetchComments(taskId) {
    try {
      const response = await apiClient.get(`/tasks/${taskId}/comments`)
      return response.data
    } catch (error) {
      console.error('Failed to fetch comments:', error)
      return []
    }
  }

  async function addComment(taskId, content) {
    try {
      const response = await apiClient.post(`/tasks/${taskId}/comments`, { content })
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to add comment' }
    }
  }

  async function updateComment(commentId, content) {
    try {
      const response = await apiClient.put(`/comments/${commentId}`, { content })
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to update comment' }
    }
  }

  async function deleteComment(commentId) {
    try {
      await apiClient.delete(`/comments/${commentId}`)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete comment' }
    }
  }

  async function fetchAttachments(taskId) {
    try {
      const response = await apiClient.get(`/tasks/${taskId}/attachments`)
      return response.data
    } catch (error) {
      console.error('Failed to fetch attachments:', error)
      return []
    }
  }

  async function uploadAttachment(taskId, file) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      const response = await apiClient.post(`/tasks/${taskId}/attachments`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to upload attachment' }
    }
  }

  async function deleteAttachment(attachmentId) {
    try {
      await apiClient.delete(`/attachments/${attachmentId}`)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Failed to delete attachment' }
    }
  }

  return {
    tasks,
    currentTask,
    loading,
    getTasksByBoard,
    setTasks,
    createTask,
    fetchTask,
    updateTask,
    moveTask,
    deleteTask,
    fetchTaskHistory,
    fetchComments,
    addComment,
    updateComment,
    deleteComment,
    fetchAttachments,
    uploadAttachment,
    deleteAttachment,
  }
})

