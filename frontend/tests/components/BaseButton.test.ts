import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import BaseButton from '@/components/BaseButton.vue'

describe('BaseButton', () => {
  it('renders slot content', () => {
    const wrapper = mount(BaseButton, { slots: { default: 'Click Me' } })
    expect(wrapper.text()).toBe('Click Me')
  })

  it('emits click event when clicked', () => {
    const wrapper = mount(BaseButton)
    wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')!.length).toBe(1)
  })

  it('applies action type classes by default', () => {
    const wrapper = mount(BaseButton)
    expect(wrapper.classes()).toContain('bg-action')
    expect(wrapper.classes()).toContain('text-tertiary')
  })

  it('applies neutral type classes', () => {
    const wrapper = mount(BaseButton, { props: { type: 'neutral' } })
    expect(wrapper.classes()).toContain('bg-secondary')
    expect(wrapper.classes()).toContain('border-white')
    expect(wrapper.classes()).toContain('border-2')
  })

  it('applies danger type classes', () => {
    const wrapper = mount(BaseButton, { props: { type: 'danger' } })
    expect(wrapper.classes()).toContain('bg-red')
  })

  it('applies danger-link type classes', () => {
    const wrapper = mount(BaseButton, { props: { type: 'danger-link' } })
    expect(wrapper.classes()).toContain('bg-none')
    expect(wrapper.classes()).toContain('underline')
    expect(wrapper.classes()).toContain('text-red')
  })

  it('forwards attributes to button element', () => {
    const wrapper = mount(BaseButton, { attrs: { 'data-test': 'btn', id: 'my-btn' } })
    expect(wrapper.attributes('data-test')).toBe('btn')
    expect(wrapper.attributes('id')).toBe('my-btn')
  })

  it('applies disabled attribute when disabled', () => {
    const wrapper = mount(BaseButton, { attrs: { disabled: true } })
    expect(wrapper.attributes('disabled')).toBeDefined()
  })
})