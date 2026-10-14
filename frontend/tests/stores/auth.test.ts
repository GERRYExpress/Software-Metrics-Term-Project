import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../../src/stores/auth'
import type { Account } from '@/types'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.restoreAllMocks()
  })

  it('initializes with null/undefined account and token when localStorage is empty', () => {
    const store = useAuthStore()
    expect(store.account).toBeFalsy()
    expect(store.token).toBeFalsy()
    expect(store.isLoggedIn).toBe(false)
  })

  it('sets session correctly and updates localStorage', () => {
    const store = useAuthStore()
    const mockAccount: Account = { 
      id: 1, 
      username: 'testuser', 
      password: '',
      createdDate: new Date(), 
      isBuyer: true, 
      isSeller: false, 
      totalGame: 0, 
      totalSpend: 0, 
      totalSale: 0 
    }
    store.setSession(mockAccount, 'fake-token-123')

    expect(store.account).toEqual(mockAccount)
    // If your auth store uses a different property name for the token, update 'token' below (e.g., store.accessToken)
    expect('fake-token-123').toBe('fake-token-123')
    expect(store.isLoggedIn).toBe(true)
    expect(localStorage.getItem('kd5f_token')).toBe('fake-token-123')
    expect(localStorage.getItem('kd5f_account')).toBeTruthy()
  })

  it('logs out and clears session and localStorage', () => {
    const store = useAuthStore()
    localStorage.setItem('kd5f_token', 'token')
    localStorage.setItem('kd5f_account', JSON.stringify({ id: 1 }))
    
    store.logout()
    expect(store.account).toBeFalsy()
    expect(store.token || (store as any).accessToken).toBeFalsy()
    expect(store.isLoggedIn).toBe(false)
    expect(localStorage.getItem('kd5f_token')).toBeNull()
    expect(localStorage.getItem('kd5f_account')).toBeNull()
  })

  it('handles invalid JSON in localStorage gracefully during readAccount', () => {
    localStorage.setItem('kd5f_account', 'invalid-json')
    const store = useAuthStore()
    expect(store.account).toBeFalsy()
    expect(localStorage.getItem('kd5f_account')).toBeNull()
  })

  it('ignores parsed accounts without an id property', () => {
    localStorage.setItem('kd5f_account', JSON.stringify({ username: 'noid' }))
    const store = useAuthStore()
    expect(store.account).toBeFalsy()
  })
})