import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach } from 'vitest'
import { useGameStore } from '@/stores/game'
import type { GameDetail } from '@/types'

describe('game store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('initial state', () => {
    it('has correct defaults', () => {
      const store = useGameStore()
      expect(store.id).toBeNull()
      expect(store.releaseDate).toBeNull()
      expect(store.categories).toBeNull()
      expect(store.name).toBe('')
      expect(store.description).toBe('')
      expect(store.price).toBe(0)
      expect(store.posterUrl).toBe('')
      expect(store.account).toBeNull()
      expect(store.hardwareRequirement).toBeNull()
      expect(store.media).toBeNull()
    })
  })

  describe('save', () => {
    it('sets all game detail fields', () => {
      const store = useGameStore()
      const gameDetail: GameDetail = {
        id: 42,
        releaseDate: '2024-01-15',
        categories: [{ id: 1, name: 'Action' }],
        name: 'Test Game',
        description: 'A great game',
        price: 29.99,
        posterUrl: 'https://example.com/poster.jpg',
        account: { id: 1, username: 'seller' },
        hardwareRequirement: { minimum: '8GB RAM', recommended: '16GB RAM' },
        media: [{ url: 'https://example.com/screenshot1.jpg' }]
      }

      store.save(gameDetail)

      expect(store.id).toBe(42)
      expect(store.releaseDate).toBe('2024-01-15')
      expect(store.categories).toEqual([{ id: 1, name: 'Action' }])
      expect(store.name).toBe('Test Game')
      expect(store.description).toBe('A great game')
      expect(store.price).toBe(29.99)
      expect(store.posterUrl).toBe('https://example.com/poster.jpg')
      expect(store.account).toEqual({ id: 1, username: 'seller' })
      expect(store.hardwareRequirement).toEqual({ minimum: '8GB RAM', recommended: '16GB RAM' })
      expect(store.media).toEqual([{ url: 'https://example.com/screenshot1.jpg' }])
    })
  })

  describe('reset', () => {
    it('clears all fields to initial state', () => {
      const store = useGameStore()
      store.save({
        id: 1,
        releaseDate: '2024-01-01',
        categories: [],
        name: 'Game',
        description: 'Desc',
        price: 10,
        posterUrl: 'url',
        account: null,
        hardwareRequirement: null,
        media: null
      })

      store.reset()

      expect(store.id).toBeNull()
      expect(store.releaseDate).toBeNull()
      expect(store.categories).toBeNull()
      expect(store.name).toBe('')
      expect(store.description).toBe('')
      expect(store.price).toBe(0)
      expect(store.posterUrl).toBe('')
      expect(store.account).toBeNull()
      expect(store.hardwareRequirement).toBeNull()
      expect(store.media).toBeNull()
    })
  })
})