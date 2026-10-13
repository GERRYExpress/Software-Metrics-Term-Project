import { describe, it, expect, vi, beforeEach } from 'vitest'
import WishlistService from '@/service/WishlistService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('WishlistService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('getWishlists', () => {
    it('fetches user wishlist', async () => {
      const mockData = [{ productId: 1, game: { id: 1 } }, { productId: 3, game: { id: 3 } }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await WishlistService.getWishlists(5)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/wishlists/5')
    })

    it('propagates error on failure', async () => {
      const error = new Error('Network Error')
      vi.mocked(apiClient.get).mockRejectedValue(error)

      await expect(WishlistService.getWishlists(5)).rejects.toThrow('Network Error')
    })
  })

  describe('addWishlist', () => {
    it('adds product to user wishlist', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({})

      await WishlistService.addWishlist(5, 10)

      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/wishlists/5/10')
    })

    it('propagates error on failure', async () => {
      const error = { response: { status: 409 }, message: 'Already in wishlist' }
      vi.mocked(apiClient.post).mockRejectedValue(error)

      await expect(WishlistService.addWishlist(5, 10)).rejects.toMatchObject({ response: { status: 409 } })
    })
  })

  describe('deleteWishlist', () => {
    it('removes product from user wishlist', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})

      await WishlistService.deleteWishlist(5, 10)

      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/wishlists/5/10')
    })

    it('propagates 404 error', async () => {
      const error = { response: { status: 404 }, message: 'Not Found' }
      vi.mocked(apiClient.delete).mockRejectedValue(error)

      await expect(WishlistService.deleteWishlist(5, 999)).rejects.toMatchObject({ response: { status: 404 } })
    })
  })
})