import { describe, it, expect, vi, beforeEach } from 'vitest'
import CategoryService from '@/service/CategoryService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('CategoryService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('getCategories', () => {
    it('returns categories on success', async () => {
      const mockData = [
        { id: 1, name: 'Action', isHidden: false },
        { id: 2, name: 'RPG', isHidden: true },
      ]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await CategoryService.getCategories()

      expect(result).toEqual(mockData)
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/categories')
    })

    it('propagates error on failure', async () => {
      const error = new Error('Network Error')
      vi.mocked(apiClient.get).mockRejectedValue(error)

      await expect(CategoryService.getCategories()).rejects.toThrow('Network Error')
    })
  })

  describe('createCategory', () => {
    it('posts category data and returns created category', async () => {
      const newCategory = { id: 3, name: 'Strategy', isHidden: false }
      vi.mocked(apiClient.post).mockResolvedValue({ data: newCategory })

      const result = await CategoryService.createCategory('Strategy')

      expect(result).toEqual(newCategory)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/categories', { name: 'Strategy', isHidden: false })
    })

    it('sends isHidden when true', async () => {
      vi.mocked(apiClient.post).mockResolvedValue({ data: { id: 4, name: 'Hidden', isHidden: true } })

      await CategoryService.createCategory('Hidden', true)

      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/categories', { name: 'Hidden', isHidden: true })
    })
  })

  describe('updateCategory', () => {
    it('puts to correct endpoint with id and data', async () => {
      const updated = { id: 1, name: 'Updated Action', isHidden: true }
      vi.mocked(apiClient.put).mockResolvedValue({ data: updated })

      const result = await CategoryService.updateCategory(1, 'Updated Action', true)

      expect(result).toEqual(updated)
      expect(apiClient.put).toHaveBeenCalledWith('/api/v1/categories/1', { name: 'Updated Action', isHidden: true })
    })
  })

  describe('deleteCategory', () => {
    it('deletes category by id', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})

      await CategoryService.deleteCategory(5)

      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/categories/5')
    })

    it('throws on 404', async () => {
      const error = { response: { status: 404 }, message: 'Not Found' }
      vi.mocked(apiClient.delete).mockRejectedValue(error)

      await expect(CategoryService.deleteCategory(999)).rejects.toMatchObject({ response: { status: 404 } })
    })
  })
})