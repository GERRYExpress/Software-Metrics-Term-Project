import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import ProductListCard from '../../src/components/ProductListCard.vue'

describe('ProductListCard.vue', () => {
  it('renders wishlist product details with wishlistedDate and emits delete-item', async () => {
    const product = {
      userId: 'user-1',
      productId: 'game-1',
      wishlistedDate: '2026-01-15T00:00:00.000Z',
      game: {
        id: 'game-1',
        name: 'Cyber Game',
        price: '59.99',
        posterUrl: 'poster.jpg',
        categories: [{ name: 'Action' }, { name: 'RPG' }]
      }
    }

    const wrapper = mount(ProductListCard, {
      props: { product },
      global: {
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\', $event)"><slot/></button>' }
        }
      }
    })

    expect(wrapper.text()).toContain('Cyber Game')
    expect(wrapper.text()).toContain('Action, RPG')
    expect(wrapper.text()).toContain('$59.99')
    expect(wrapper.text()).toContain('Added:')

    await wrapper.find('button').trigger('click', {
      preventDefault: () => {}
    })
    
    expect(wrapper.emitted('delete-item')).toBeTruthy()
    expect(wrapper.emitted('delete-item')?.[0]).toEqual(['user-1', 'game-1'])
  })

  it('renders cart product details with cartDate', () => {
    const product = {
      userId: 'user-1',
      productId: 'game-2',
      cartDate: '2026-05-10T00:00:00.000Z',
      game: {
        id: 'game-2',
        name: 'Indie Gem',
        price: '19.99',
        posterUrl: 'poster2.jpg',
        categories: [{ name: 'Indie' }]
      }
    }

    const wrapper = mount(ProductListCard, {
      props: { product },
      global: {
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: true
        }
      }
    })

    expect(wrapper.text()).toContain('Indie Gem')
    expect(wrapper.text()).toContain('Indie')
  })

  it('renders without date if neither wishlistedDate nor cartDate is present', () => {
    const product = {
      userId: 'user-1',
      productId: 'game-3',
      game: {
        id: 'game-3',
        name: 'No Date Game',
        price: '10.00',
        posterUrl: 'poster3.jpg',
        categories: [{ name: 'Casual' }]
      }
    }

    const wrapper = mount(ProductListCard, {
      props: { product },
      global: {
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: true
        }
      }
    })

    expect(wrapper.text()).toContain('No Date Game')
    expect(wrapper.text()).toContain('Added: ')
  })
})