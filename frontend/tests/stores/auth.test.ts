import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuthStore } from '@/stores/auth'
import AccountService from '@/service/AccountService'

vi.mock('@/service/AccountService')

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('has correct defaults', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('login', () => {
    it('sets user and token on success', async () => {
      const store = useAuthStore()
      const mockUser = { id: 1, username: 'testuser', email: 'test@test.com' }
      const mockToken = 'jwt-token-123'

      vi.mocked(AccountService.login).mockResolvedValue({ user: mockUser, token: mockToken })

      await store.login('test@test.com', 'password123')

      expect(store.user).toEqual(mockUser)
      expect(store.token).toBe(mockToken)
      expect(store.isAuthenticated).toBe(true)
      expect(AccountService.login).toHaveBeenCalledWith('test@test.com', 'password123')
    })

    it('throws and does not set state on failure', async () => {
      const store = useAuthStore()
      vi.mocked(AccountService.login).mockRejectedValue(new Error('Invalid credentials'))

      await expect(store.login('test@test.com', 'wrong')).rejects.toThrow('Invalid credentials')

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('register', () => {
    it('sets user and token on success', async () => {
      const store = useAuthStore()
      const mockUser = { id: 2, username: 'newuser', email: 'new@test.com' }
      const mockToken = 'new-jwt-token'

      vi.mocked(AccountService.register).mockResolvedValue({ user: mockUser, token: mockToken })

      await store.register({ username: 'newuser', email: 'new@test.com', password: 'password123' })

      expect(store.user).toEqual(mockUser)
      expect(store.token).toBe(mockToken)
      expect(store.isAuthenticated).toBe(true)
    })
  })

  describe('logout', () => {
    it('clears user, token, and authentication state', () => {
      const store = useAuthStore()
      store.user = { id: 1, username: 'test' } as any
      store.token = 'some-token'
      store.isAuthenticated = true

      store.logout()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('fetchProfile', () => {
    it('updates user from API', async () => {
      const store = useAuthStore()
      store.token = 'valid-token'
      const mockUser = { id: 1, username: 'updated', email: 'updated@test.com' }

      vi.mocked(AccountService.getProfile).mockResolvedValue(mockUser)

      await store.fetchProfile()

      expect(store.user).toEqual(mockUser)
      expect(AccountService.getProfile).toHaveBeenCalled()
    })

    it('logs out on 401', async () => {
      const store = useAuthStore()
      store.token = 'expired-token'
      store.user = { id: 1 } as any
      store.isAuthenticated = true

      vi.mocked(AccountService.getProfile).mockRejectedValue({ response: { status: 401 } })

      await store.fetchProfile()

      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })
})