import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserActionStore } from '../../src/stores/userAction'
import { useAuthStore } from '../../src/stores/auth'
import CartService from '@/service/CartService'
import WishlistService from '@/service/WishlistService'
import type { AxiosResponse } from 'axios'
import type { Cart, Wishlist, Account } from '@/types'

vi.mock('@/service/CartService')
vi.mock('@/service/WishlistService')

describe('userAction store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('notifies and clears notifications with timeouts', () => {
    const store = useUserActionStore()
    store.notify('Success message', 'success')
    expect(store.activeNotification?.message).toBe('Success message')

    store.clearNotification()
    expect(store.activeNotification).toBeNull()

    store.notify('Auto close')
    vi.advanceTimersByTime(3200)
    expect(store.activeNotification).toBeNull()
  })

  it('fetches initial cart and wishlist data successfully', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1 } as Account

    const mockCartResponse = { data: [{ productId: 101 }] } as unknown as AxiosResponse<Cart[]>
    const mockWishlistResponse = { data: [{ game: { id: 202 } }] } as unknown as AxiosResponse<Wishlist[]>

    vi.mocked(CartService.getCarts).mockResolvedValueOnce(mockCartResponse)
    vi.mocked(WishlistService.getWishlists).mockResolvedValueOnce(mockWishlistResponse)

    const store = useUserActionStore()
    store.fetchInitialData()
    store.fetchInitialData() // Test early return guard

    await vi.runAllTimersAsync()

    expect(store.cartIds).toContain(101)
    expect(store.wishlistIds).toContain(202)
    expect(store.isInCart(101)).toBe(true)
    expect(store.isInWishlist(202)).toBe(true)
  })

  it('handles offline/error catch blocks in fetchInitialData gracefully', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1 } as Account

    vi.mocked(CartService.getCarts).mockRejectedValueOnce(new Error('Network error'))
    vi.mocked(WishlistService.getWishlists).mockRejectedValueOnce(new Error('Network error'))

    const store = useUserActionStore()
    store.fetchInitialData()
    await vi.runAllTimersAsync()

    expect(store.cartIds).toEqual([])
  })

  it('handles cart additions including already in cart and unauthenticated notifications', async () => {
    const store = useUserActionStore()
    store.cartIds.push(50)

    store.addToCart(50, 'Game 50')
    expect(store.activeNotification?.type).toBe('info')

    store.cartIds = []
    store.addToCart(60, 'Game 60')
    expect(store.activeNotification?.message).toContain('Please sign in')

    const authStore = useAuthStore()
    authStore.account = { id: 1 } as Account
    const emptyResponse = {} as unknown as AxiosResponse<void>
    vi.mocked(CartService.addCart).mockResolvedValueOnce(emptyResponse)

    store.addToCart(60, 'Game 60')
    await vi.runAllTimersAsync()
    expect(store.cartIds).toContain(60)

    vi.mocked(CartService.addCart).mockRejectedValueOnce(new Error('fail'))
    store.addToCart(70, 'Game 70')
    await vi.runAllTimersAsync()
    expect(store.cartIds).toContain(70)
  })

  it('handles toggleWishlist branch logic (add, remove, unauthenticated, error catches)', async () => {
    const store = useUserActionStore()
    const authStore = useAuthStore()
    const emptyResponse = {} as unknown as AxiosResponse<void>

    store.wishlistIds.push(10)
    authStore.account = { id: 1 } as Account
    vi.mocked(WishlistService.deleteWishlist).mockResolvedValueOnce(emptyResponse)

    store.toggleWishlist(10, 'Game 10')
    await vi.runAllTimersAsync()
    expect(store.wishlistIds).not.toContain(10)

    store.wishlistIds.push(20)
    vi.mocked(WishlistService.deleteWishlist).mockRejectedValueOnce(new Error('fail'))
    store.toggleWishlist(20, 'Game 20')
    await vi.runAllTimersAsync()
    expect(store.wishlistIds).not.toContain(20)

    authStore.account = null
    store.toggleWishlist(30, 'Game 30')
    expect(store.activeNotification?.message).toContain('Please sign in')

    authStore.account = { id: 1 } as Account
    vi.mocked(WishlistService.addWishlist).mockResolvedValueOnce(emptyResponse)
    store.toggleWishlist(30, 'Game 30')
    await vi.runAllTimersAsync()
    expect(store.wishlistIds).toContain(30)

    vi.mocked(WishlistService.addWishlist).mockRejectedValueOnce(new Error('fail'))
    store.toggleWishlist(40, 'Game 40')
    await vi.runAllTimersAsync()
    expect(store.wishlistIds).toContain(40)
  })
})