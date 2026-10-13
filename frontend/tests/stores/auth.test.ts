import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuthStore } from '@/stores/auth'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('initial state', () => {
    it('has correct defaults', () => {
      const store = useAuthStore()
      // account is a ref initialized from localStorage (null if empty)
      expect(store.account).toBeNull()
      expect(store.isLoggedIn).toBe(false)
    })
  })

  describe('setSession', () => {
    it('sets account and token', () => {
      const store = useAuthStore()
      const mockAccount = { id: 1, username: 'testuser', email: 'test@test.com' }
      const mockToken = 'jwt-token-123'

      store.setSession(mockAccount, mockToken)

      expect(store.account).toEqual(mockAccount)
      expect(store.isLoggedIn).toBe(true)
      expect(localStorage.getItem('kd5f_token')).toBe(mockToken)
      expect(localStorage.getItem('kd5f_account')).toBe(JSON.stringify(mockAccount))
    })
  })

  describe('logout', () => {
    it('clears account, token, and authentication state', () => {
      const store = useAuthStore()
      store.setSession({ id: 1, username: 'test' }, 'some-token')

      store.logout()

      expect(store.account).toBeNull()
      expect(store.isLoggedIn).toBe(false)
      expect(localStorage.getItem('kd5f_token')).toBeNull()
      expect(localStorage.getItem('kd5f_account')).toBeNull()
    })
  })
})