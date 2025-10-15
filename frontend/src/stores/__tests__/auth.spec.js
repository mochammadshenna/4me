import apiClient from '@/api/client'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useAuthStore } from '../auth'

vi.mock('@/api/client')

describe('Auth Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        vi.clearAllMocks()
    })

    it('initializes with default state', () => {
        const store = useAuthStore()

        expect(store.user).toBeNull()
        expect(store.token).toBe('')
        expect(store.isAuthenticated).toBe(false)
    })

    it('loads token from localStorage on init', () => {
        localStorage.setItem('token', 'test-token')
        const store = useAuthStore()

        expect(store.token).toBe('test-token')
        expect(store.isAuthenticated).toBe(true)
    })

    describe('login', () => {
        it('successfully logs in user', async () => {
            const mockResponse = {
                data: {
                    token: 'new-token',
                    refresh_token: 'new-refresh-token',
                    user: {
                        id: 1,
                        username: 'testuser',
                        email: 'test@example.com'
                    }
                }
            }

            apiClient.post.mockResolvedValue(mockResponse)

            const store = useAuthStore()
            const result = await store.login({
                username: 'testuser',
                password: 'password123'
            })

            expect(result.success).toBe(true)
            expect(store.token).toBe('new-token')
            expect(store.user.username).toBe('testuser')
            expect(localStorage.getItem('token')).toBe('new-token')
        })

        it('handles login failure', async () => {
            apiClient.post.mockRejectedValue({
                response: {
                    data: { error: 'Invalid credentials' }
                }
            })

            const store = useAuthStore()
            const result = await store.login({
                username: 'testuser',
                password: 'wrongpassword'
            })

            expect(result.success).toBe(false)
            expect(result.error).toBe('Invalid credentials')
            expect(store.token).toBe('')
        })
    })

    describe('register', () => {
        it('successfully registers user', async () => {
            const mockResponse = {
                data: {
                    token: 'new-token',
                    refresh_token: 'new-refresh-token',
                    user: {
                        id: 1,
                        username: 'newuser',
                        email: 'new@example.com'
                    }
                }
            }

            apiClient.post.mockResolvedValue(mockResponse)

            const store = useAuthStore()
            const result = await store.register({
                username: 'newuser',
                email: 'new@example.com',
                password: 'password123'
            })

            expect(result.success).toBe(true)
            expect(store.user.username).toBe('newuser')
        })
    })

    describe('logout', () => {
        it('clears user data and token', () => {
            localStorage.setItem('token', 'test-token')
            const store = useAuthStore()
            store.user = { id: 1, username: 'test' }
            store.token = 'test-token'

            store.logout()

            expect(store.user).toBeNull()
            expect(store.token).toBe('')
            expect(localStorage.getItem('token')).toBeNull()
        })
    })

    describe('fetchUser', () => {
        it('fetches current user data', async () => {
            const mockUser = {
                id: 1,
                username: 'testuser',
                email: 'test@example.com'
            }

            apiClient.get.mockResolvedValue({ data: mockUser })

            const store = useAuthStore()
            store.token = 'test-token'

            await store.fetchUser()

            expect(store.user).toEqual(mockUser)
        })

        it('logs out on fetch failure', async () => {
            apiClient.get.mockRejectedValue(new Error('Unauthorized'))

            const store = useAuthStore()
            store.token = 'invalid-token'

            await store.fetchUser()

            expect(store.user).toBeNull()
            expect(store.token).toBe('')
        })
    })
})

