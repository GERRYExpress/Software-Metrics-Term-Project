import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSellerStore } from '../../src/stores/seller'
import { useAuthStore } from '../../src/stores/auth'
import type { Account } from '@/types'

describe('seller store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('throws error if user is not signed in when loading seller', async () => {
    const sellerStore = useSellerStore()
    await expect(sellerStore.loadSeller()).rejects.toThrow('Please sign in before opening the seller center.')
    expect(sellerStore.errorMessage).toBe('Please sign in before opening the seller center.')
  })

  it('loads seller successfully if authenticated', async () => {
    const authStore = useAuthStore()
    const mockAccount: Account = { 
      id: 1, 
      username: 'seller', 
      password: '',
      createdDate: new Date(), 
      isBuyer: false, 
      isSeller: true, 
      totalGame: 0, 
      totalSpend: 0, 
      totalSale: 0 
    }
    authStore.account = mockAccount

    const sellerStore = useSellerStore()
    const result = await sellerStore.loadSeller()
    expect(result).toEqual(mockAccount)
    expect(sellerStore.account).toEqual(mockAccount)
  })

  it('clears seller state and error messages', () => {
    const sellerStore = useSellerStore()
    sellerStore.errorMessage = 'Some error'
    sellerStore.account = { 
      id: 1, 
      username: 'seller', 
      password: '',
      createdDate: new Date(), 
      isBuyer: false, 
      isSeller: true, 
      totalGame: 0, 
      totalSpend: 0, 
      totalSale: 0 
    }

    sellerStore.clear()
    expect(sellerStore.account).toBeNull()
    expect(sellerStore.errorMessage).toBe('')
  })
})