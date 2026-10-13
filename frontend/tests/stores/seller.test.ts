import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useSellerStore } from '@/stores/seller'
import type { Account } from '@/types'

describe('seller store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('initial state', () => {
    it('has correct defaults', () => {
      const store = useSellerStore()
      expect(store.account).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.errorMessage).toBe('')
    })
  })

  describe('clear', () => {
    it('clears account and error message', () => {
      const store = useSellerStore()
      store.account = { value: { id: 1 } as Account }
      store.errorMessage = 'Some error'

      store.clear()

      expect(store.account).toBeNull()
      expect(store.errorMessage).toBe('')
    })
  })
})