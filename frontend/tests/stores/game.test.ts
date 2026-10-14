import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGameStore } from '../../src/stores/game'

describe('game store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('saves game detail correctly', () => {
    const store = useGameStore()
    const gameDetail = {
      id: 10,
      releaseDate: new Date(),
      categories: [],
      name: 'Test Game',
      description: 'Desc',
      price: 50,
      posterUrl: 'url',
      account: null,
      hardwareRequirement: [],
      media: []
    }

    store.save(gameDetail)
    expect(store.id).toBe(10)
    expect(store.name).toBe('Test Game')
    expect(store.price).toBe(50)
  })

  it('resets game state back to defaults', () => {
    const store = useGameStore()
    store.id = 10
    store.name = 'Test Game'

    store.reset()
    expect(store.id).toBeNull()
    expect(store.name).toBe('')
    expect(store.price).toBe(0)
  })
})