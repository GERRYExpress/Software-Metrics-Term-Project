import { describe, it, expect, vi, beforeEach } from 'vitest'
import CartService from '@/service/CartService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('CartService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('getCarts', () => {
    it('fetches user cart', async () => {
      const mockData = [{ productId: 1, game: { id: 1 } }, { productId: 2, game: { id: 2 } }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await CartService.getCarts(5)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/carts/5')
    })

    it('propagates error on failure', async () => {
      const error = new Error('Network Error')
      vi.mocked(apiClient.get).mockRejectedValue(error)

      await expect(CartService.getCarts(5)).rejects.toThrow('Network Error')
    })
  })

  describe('addCart', () => {
    it('adds product to user cart', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({})

      await CartService.addCart(5, 10)

      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/carts/5/10')
    })

    it('propagates error on failure', async () => {
      const error = { response: { status: 409 }, message: 'Already in cart' }
      vi.mocked(apiClient.post).mockRejectedValue(error)

      await expect(CartService.addCart(5, 10)).rejects.toMatchObject({ response: { status: 409 } })
    })
  })

  describe('deleteCart', () => {
    it('removes product from user cart', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})

      await CartService.deleteCart(5, 10)

      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/carts/5/10')
    })

    it('propagates 404 error', async () => {
      const error = { response: { status: 404 }, message: 'Not Found' }
      vi.mocked(apiClient.delete).mockRejectedValue(error)

      await expect(CartService.deleteCart(5, 999)).rejects.toMatchObject({ response: { status: 404 } })
    })
  })
})