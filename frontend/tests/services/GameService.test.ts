import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameService from '@/service/GameService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('GameService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('getGames', () => {
    it('fetches paginated games', async () => {
      const mockData = [{ id: 1, name: 'Game 1' }, { id: 2, name: 'Game 2' }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGames(10, 1)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?_limit=10&_page=1')
    })
  })

  describe('getGamesByKeyword', () => {
    it('fetches games by keyword with pagination', async () => {
      const mockData = [{ id: 1, name: 'Action Game' }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGamesByKeyword(5, 2, 'action')

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?keyword=action&_limit=5&_page=2')
    })
  })

  describe('getGamesByCategory', () => {
    it('fetches games by category with pagination', async () => {
      const mockData = [{ id: 1, name: 'RPG Game' }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGamesByCategory(6, 1, 'RPG')

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?category=RPG&_limit=6&_page=1')
    })
  })

  describe('getGamesByKeywordAndCategory', () => {
    it('fetches games by keyword and category', async () => {
      const mockData = [{ id: 1, name: 'Action RPG' }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGamesByKeywordAndCategory(10, 1, 'action', 'RPG')

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?keyword=action&category=RPG&_limit=10&_page=1')
    })
  })

  describe('getGame', () => {
    it('fetches single game by id', async () => {
      const mockData = { id: 42, name: 'Specific Game', price: 29.99 }
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGame(42)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games/42')
    })
  })

  describe('getGamesBySeller', () => {
    it('fetches seller games with default pagination', async () => {
      const mockData = [{ id: 1 }, { id: 2 }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })

      const result = await GameService.getGamesBySeller(99)

      expect(result).toEqual({ data: mockData })
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?sellerId=99&_limit=6&_page=1')
    })

    it('fetches seller games with custom pagination', async () => {
      vi.mocked(apiClient.get).mockResolvedValue({ data: [] })

      await GameService.getGamesBySeller(99, 12, 3)

      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/games?sellerId=99&_limit=12&_page=3')
    })
  })

  describe('createGame', () => {
    it('posts game data and returns created game response', async () => {
      const input = { userId: 1, name: 'New Game', price: 19.99, categoryIds: [1] } as any
      const created = { id: 100, ...input }
      vi.mocked(apiClient.post).mockResolvedValue({ data: created })

      const result = await GameService.createGame(input)

      expect(result).toEqual({ data: created })
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/games', expect.objectContaining({
        userId: 1,
        name: 'New Game',
        price: 19.99,
        categoryIds: [1]
      }))
    })
  })

  describe('updateGame', () => {
    it('puts updated game data', async () => {
      const input = { name: 'Updated Game', price: 24.99 } as any
      const updated = { id: 5, ...input }
      vi.mocked(apiClient.put).mockResolvedValue({ data: updated })

      const result = await GameService.updateGame(5, input)

      expect(result).toEqual({ data: updated })
      expect(apiClient.put).toHaveBeenCalledWith('/api/v1/games/5', expect.objectContaining({
        name: 'Updated Game',
        price: 24.99
      }))
    })
  })

  describe('deleteGame', () => {
    it('deletes game with userId query', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})

      await GameService.deleteGame(10, 5)

      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/games/10?userId=5')
    })

    it('propagates 404 error', async () => {
      const error = { response: { status: 404 }, message: 'Not Found' }
      vi.mocked(apiClient.delete).mockRejectedValue(error)

      await expect(GameService.deleteGame(999, 1)).rejects.toMatchObject({ response: { status: 404 } })
    })
  })
})