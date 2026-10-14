import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import GameCard from '../../src/components/GameCard.vue'
import { useUserActionStore } from '@/stores/userAction'
import type { Game } from '@/types'

describe('GameCard.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const sampleGame = {
    id: 1,
    name: 'Cyber Adventure',
    price: '49.99',
    posterUrl: 'https://example.com/poster.jpg',
    categories: [{ name: 'Action' }, { name: 'Sci-Fi' }]
  }

  it('renders game card details, categories, and formatted price', () => {
    const wrapper = mount(GameCard, {
      props: { game: sampleGame },
      global: {
        stubs: {
          RouterLink: { template: '<a><slot/></a>' },
          Image: { template: '<img :src="$attrs.src" />' },
          Heart: true,
          ShoppingCartPlus: true,
          BaseButton: { template: '<button @click="$emit(\'click\', $event)"><slot/></button>' }
        }
      }
    })

    expect(wrapper.text()).toContain('Cyber Adventure')
    expect(wrapper.text()).toContain('Action, Sci-Fi')
    expect(wrapper.text()).toContain('$49.99')
  })

  it('handles wishlist toggle action via store', async () => {
    const store = useUserActionStore()
    const toggleWishlistSpy = vi.spyOn(store, 'toggleWishlist')

    const wrapper = mount(GameCard, {
      props: { game: sampleGame },
      global: {
        stubs: {
          RouterLink: { template: '<a><slot/></a>' },
          Image: true,
          Heart: true,
          ShoppingCartPlus: true,
          BaseButton: { template: '<button @click="$emit(\'click\', $event)"><slot/></button>' }
        }
      }
    })

    const buttons = wrapper.findAll('button')
    await buttons[0].trigger('click', { preventDefault: () => {} })

    expect(toggleWishlistSpy).toHaveBeenCalledWith(1, 'Cyber Adventure')
  })

  it('handles add to cart action via store', async () => {
    const store = useUserActionStore()
    const addToCartSpy = vi.spyOn(store, 'addToCart')

    const wrapper = mount(GameCard, {
      props: { game: sampleGame },
      global: {
        stubs: {
          RouterLink: { template: '<a><slot/></a>' },
          Image: true,
          Heart: true,
          ShoppingCartPlus: true,
          Check: true,
          BaseButton: { template: '<button @click="$emit(\'click\', $event)"><slot/></button>' }
        }
      }
    })

    const buttons = wrapper.findAll('button')
    await buttons[1].trigger('click', { preventDefault: () => {} })

    expect(addToCartSpy).toHaveBeenCalledWith(1, 'Cyber Adventure')
  })

  it('renders fallback category and default name when properties are missing', () => {
    const minimalGame: Partial<Game> = {
      id: 2,
      name: undefined, // undefined correctly triggers the 'NAME' fallback
      price: 0,
      posterUrl: 'https://example.com/poster.jpg'
    }

    const wrapper = mount(GameCard, {
      props: { game: minimalGame as Game },
      global: {
        stubs: {
          RouterLink: { template: '<a><slot/></a>' },
          Image: true,
          Heart: true,
          ShoppingCartPlus: true,
          BaseButton: true
        }
      }
    })

    expect(wrapper.text()).toContain('CATEGORY')
    expect(wrapper.text()).toContain('NAME')
    expect(wrapper.text()).toContain('$0.00')
  })
})