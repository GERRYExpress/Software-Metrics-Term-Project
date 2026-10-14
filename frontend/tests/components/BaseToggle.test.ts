import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import BaseToggle from '@/components/BaseToggle.vue'

describe('BaseToggle', () => {
  it('renders both option labels', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['Buyer', 'Seller'], modelValue: false } })
    expect(wrapper.text()).toContain('Buyer')
    expect(wrapper.text()).toContain('Seller')
  })

  it('highlights first option when modelValue is false', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['Buyer', 'Seller'], modelValue: false } })
    const buttons = wrapper.findAll('.flex > div')
    expect(buttons[0].classes()).toContain('bg-action')
    expect(buttons[0].classes()).toContain('text-tertiary')
    expect(buttons[1].classes()).toContain('bg-gray-100')
    expect(buttons[1].classes()).toContain('text-gray-300')
  })

  it('highlights second option when modelValue is true', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['Buyer', 'Seller'], modelValue: true } })
    const buttons = wrapper.findAll('.flex > div')
    expect(buttons[0].classes()).toContain('bg-gray-100')
    expect(buttons[0].classes()).toContain('text-gray-300')
    expect(buttons[1].classes()).toContain('bg-action')
    expect(buttons[1].classes()).toContain('text-tertiary')
  })

  it('emits toggleRole on click', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['A', 'B'], modelValue: false } })
    wrapper.find('.flex').trigger('click')
    expect(wrapper.emitted('toggleRole')).toBeTruthy()
  })

  it('applies hover and transition classes to container', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['A', 'B'], modelValue: false } })
    const container = wrapper.find('.flex')
    expect(container.classes()).toContain('transition-all')
    expect(container.classes()).toContain('cursor-pointer')
    expect(container.classes()).toContain('hover:shadow-lg')
  })

  it('applies button classes to inner divs', () => {
    const wrapper = mount(BaseToggle, { props: { options: ['A', 'B'], modelValue: false } })
    const buttons = wrapper.findAll('.flex > div')
    expect(buttons[0].classes()).toContain('w-full')
    expect(buttons[0].classes()).toContain('text-center')
    expect(buttons[0].classes()).toContain('p-2')
    expect(buttons[0].classes()).toContain('font-medium')
    expect(buttons[0].classes()).toContain('rounded-s-xl')
    expect(buttons[1].classes()).toContain('rounded-e-xl')
  })
})