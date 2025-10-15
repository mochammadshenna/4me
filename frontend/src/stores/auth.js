import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')
  const refreshToken = ref(localStorage.getItem('refreshToken') || '')

  const isAuthenticated = computed(() => !!token.value)

  async function login(credentials) {
    try {
      const response = await apiClient.post('/auth/login', credentials)
      setAuthData(response.data)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Login failed' }
    }
  }

  async function register(userData) {
    try {
      const response = await apiClient.post('/auth/register', userData)
      setAuthData(response.data)
      return { success: true }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Registration failed' }
    }
  }

  async function fetchUser() {
    try {
      const response = await apiClient.get('/auth/me')
      user.value = response.data
    } catch (error) {
      logout()
    }
  }

  function setAuthData(data) {
    user.value = data.user
    token.value = data.token
    refreshToken.value = data.refresh_token
    localStorage.setItem('token', data.token)
    localStorage.setItem('refreshToken', data.refresh_token)
  }

  function logout() {
    user.value = null
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
  }

  async function googleLogin() {
    try {
      const response = await apiClient.get('/auth/google')
      window.location.href = response.data.url
    } catch (error) {
      console.error('Google login failed:', error)
    }
  }

  function handleGoogleCallback(urlToken, urlRefreshToken) {
    const authData = {
      token: urlToken,
      refresh_token: urlRefreshToken,
      user: {} // Will be fetched separately
    }
    setAuthData(authData)
    fetchUser()
  }

  return {
    user,
    token,
    refreshToken,
    isAuthenticated,
    login,
    register,
    fetchUser,
    logout,
    googleLogin,
    handleGoogleCallback,
  }
})

