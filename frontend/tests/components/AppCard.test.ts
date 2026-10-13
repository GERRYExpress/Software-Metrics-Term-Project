import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import AppCard from '@/components/AppCard.vue'

describe('AppCard', () => {
  it('renders slot content', () => {
    const wrapper = mount(AppCard, { slots: { default: '<p>Card Content</p>' } })
    expect(wrapper.text()).toBe('Card Content')
  })

  it('applies base card classes', () => {
    const wrapper = mount(AppCard)
    expect(wrapper.classes()).toContain('flex')
    expect(wrapper.classes()).toContain('flex-col')
    expect(wrapper.classes()).toContain('bg-secondary')
    expect(wrapper.classes()).toContain('rounded-xl')
    expect(wrapper.classes()).toContain('p-5')
    expect(wrapper.classes()).toContain('shadow-lg')
    expect(wrapper.classes()).toContain('w-2xl')
  })

  it('renders multiple slot elements', () => {
    const wrapper = mount(AppCard, {
      slots: {
        default: '<h3>Title</h3><p>Description</p>'
      }
    })
    expect(wrapper.find('h3').text()).toBe('Title')
    expect(wrapper.find('p').text()).toBe('Description')
  })
})