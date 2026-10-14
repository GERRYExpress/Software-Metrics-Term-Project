import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import PriceSummaryCard from '../../src/components/PriceSummaryCard.vue'

describe('PriceSummaryCard.vue', () => {
  it('disables checkout button and shows zero totals when games is empty or undefined', () => {
    const wrapper = mount(PriceSummaryCard, {
      props: { games: [] },
      global: { 
        stubs: { 
          AppCard: { template: '<div><slot/></div>' }, 
          BaseButton: { template: '<button :disabled="$attrs.disabled"><slot/></button>' } 
        } 
      }
    })

    expect(wrapper.text()).toContain('$0.00')
    const button = wrapper.find('button')
    expect(button.attributes('disabled')).toBeDefined()
  })

  it('calculates subtotal, VAT, and total correctly when games are provided', async () => {
    const wrapper = mount(PriceSummaryCard, {
      props: {
        games: [
          { id: 1, name: 'Game A', price: 50 },
          { id: 2, name: 'Game B', price: 30 }
        ]
      },
      global: {
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    expect(wrapper.text()).toContain('Game A')
    expect(wrapper.text()).toContain('$50.00')
    expect(wrapper.text()).toContain('Subtotal')
    expect(wrapper.text()).toContain('$80.00')
    expect(wrapper.text()).toContain('$5.60')
    expect(wrapper.text()).toContain('$85.60')

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('checkout')).toBeTruthy()
  })
})