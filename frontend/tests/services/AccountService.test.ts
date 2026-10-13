import { describe, it, expect, vi, beforeEach } from 'vitest'
import AccountService from '@/service/AccountService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('AccountService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('authenticate', () => {
    it('posts credentials and returns response on success', async () => {
      const mockResponse = { data: { user: { id: 1 }, token: 'abc' } }
      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      const result = await AccountService.authenticate({ username: 'test', password: 'pass' })

      expect(result).toEqual(mockResponse)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/accounts/login', {
        username: 'test',
        password: 'pass'
      })
    })

    it('propagates error on failure', async () => {
      const error = new Error('Unauthorized')
      vi.mocked(apiClient.post).mockRejectedValue(error)

      await expect(AccountService.authenticate({ username: 'test', password: 'wrong' })).rejects.toThrow('Unauthorized')
    })
  })

  describe('register', () => {
    it('posts registration data and returns response on success', async () => {
      const mockResponse = { data: { user: { id: 2 }, token: 'new-token' } }
      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      const result = await AccountService.register({ username: 'newuser', password: 'password123' })

      expect(result).toEqual(mockResponse)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/accounts', {
        username: 'newuser',
        password: 'password123'
      })
    })

    it('propagates validation error', async () => {
      const error = { response: { status: 400, data: { message: 'Email exists' } } }
      vi.mocked(apiClient.post).mockRejectedValue(error)

      await expect(AccountService.register({ username: 'exists', password: 'pass' })).rejects.toMatchObject({
        response: { status: 400 }
      })
    })
  })

  describe('getAccounts', () => {
    it('fetches all accounts', async () => {
      const mockData = [{ id: 1 }, { id: 2 }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await AccountService.getAccounts()

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/accounts')
    })
  })

  describe('getAccount', () => {
    it('fetches account by id', async () => {
      const mockData = { id: 5, username: 'user5' }
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await AccountService.getAccount(5)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/accounts/5')
    })
  })

  describe('deleteAccount', () => {
    it('deletes account by id', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})

      await AccountService.deleteAccount(7)

      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/accounts/7')
    })
  })
})