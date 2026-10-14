import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import ProductDetailView from '../../src/views/ProductDetailView.vue'
import { resetViewTestState } from './viewTestUtils'

describe('ProductDetailView.vue', () => {
  beforeEach(() => {
    resetViewTestState()
  })

  const seedGame = async () => {
    const { useGameStore } = await import('@/stores/game')
    const gameStore = useGameStore()
    gameStore.save({
      id: 42,
      releaseDate: '2024-01-15',
      categories: [{ id: 1, name: 'Action' }],
      name: 'Test Game',
      description: 'A great game',
      price: 29.99,
      posterUrl: 'https://example.com/poster.jpg',
      account: { id: 1, username: 'seller' },
      hardwareRequirement: null,
      media: null
    })
  }

  const mountDetail = () =>
    mount(ProductDetailView, {
      global: {
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\', $event)"><slot/></button>' },
          MediaCarousel: true,
          SystemRequirement: true,
          RouterLink: { template: '<a><slot/></a>' },
          Image: true,
          Check: true,
          Heart: true
        }
      }
    })

  it('renders game details from the store', async () => {
    await seedGame()
    const wrapper = mountDetail()

    expect(wrapper.text()).toContain('Test Game')
    expect(wrapper.text()).toContain('A great game')
    expect(wrapper.text()).toContain('$29.99')
    expect(wrapper.text()).toContain('seller')
  })

  it('dispatches add-to-cart and wishlist actions via the store', async () => {
    await seedGame()
    const { useUserActionStore } = await import('@/stores/userAction')
    const userActionStore = useUserActionStore()
    const addToCartSpy = vi.spyOn(userActionStore, 'addToCart')
    const wishlistSpy = vi.spyOn(userActionStore, 'toggleWishlist')

    const wrapper = mountDetail()
    const buttons = wrapper.findAll('button')
    await buttons[0].trigger('click', { preventDefault: () => {} })
    await buttons[1].trigger('click', { preventDefault: () => {} })

    expect(addToCartSpy).toHaveBeenCalledWith(42, 'Test Game')
    expect(wishlistSpy).toHaveBeenCalledWith(42, 'Test Game')
  })
})
