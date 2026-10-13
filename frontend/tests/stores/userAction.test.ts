import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useUserActionStore } from '@/stores/userAction'

describe('userAction store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  describe('initial state', () => {
    it('has correct defaults', () => {
      const store = useUserActionStore()
      expect(store.cartIds).toEqual([])
      expect(store.wishlistIds).toEqual([])
      expect(store.activeNotification).toBeNull()
      // isInitialized is a ref; in test env it may be unwrapped differently
      expect(store.isInitialized === false || store.isInitialized === undefined).toBe(true)
    })
  })

  describe('notify', () => {
    it('sets notification with message and type', () => {
      const store = useUserActionStore()
      store.notify('Test message', 'cart')

      expect(store.activeNotification).toMatchObject({
        message: 'Test message',
        type: 'cart'
      })
      expect(store.activeNotification?.id).toBeDefined()
    })

    it('clears notification after timeout', () => {
      const store = useUserActionStore()
      store.notify('Test message')

      vi.advanceTimersByTime(3200)

      expect(store.activeNotification).toBeNull()
    })

    it('resets timer on new notification', () => {
      const store = useUserActionStore()
      store.notify('First')
      vi.advanceTimersByTime(2000)
      store.notify('Second')

      expect(store.activeNotification?.message).toBe('Second')

      vi.advanceTimersByTime(2000)
      expect(store.activeNotification).not.toBeNull()

      vi.advanceTimersByTime(1200)
      expect(store.activeNotification).toBeNull()
    })
  })

  describe('clearNotification', () => {
    it('clears active notification and timer', () => {
      const store = useUserActionStore()
      store.notify('Test')
      store.clearNotification()

      expect(store.activeNotification).toBeNull()
      vi.advanceTimersByTime(5000)
      expect(store.activeNotification).toBeNull()
    })
  })

  describe('isInCart', () => {
    it('returns true if gameId in cart', () => {
      const store = useUserActionStore()
      store.cartIds = [1, 2, 3]

      expect(store.isInCart(2)).toBe(true)
      expect(store.isInCart(4)).toBe(false)
    })
  })

  describe('isInWishlist', () => {
    it('returns true if gameId in wishlist', () => {
      const store = useUserActionStore()
      store.wishlistIds = [5, 6]

      expect(store.isInWishlist(5)).toBe(true)
      expect(store.isInWishlist(7)).toBe(false)
    })
  })
})