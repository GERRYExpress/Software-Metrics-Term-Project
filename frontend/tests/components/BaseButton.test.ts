import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
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

  it('does not emit click when disabled', () => {
    const wrapper = mount(BaseButton, { props: { disabled: true } })
    wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('applies disabled attribute', () => {
    const wrapper = mount(BaseButton, { props: { disabled: true } })
    expect(wrapper.attributes('disabled')).toBeDefined()
    expect(wrapper.classes()).toContain('disabled')
  })

  it('shows loading state', () => {
    const wrapper = mount(BaseButton, { props: { loading: true } })
    expect(wrapper.find('.spinner, .loading, [data-loading]').exists()).toBe(true)
    expect(wrapper.attributes('disabled')).toBeDefined()
  })

  it('applies variant classes', () => {
    const primary = mount(BaseButton, { props: { variant: 'primary' } })
    expect(primary.classes()).toContain('btn-primary')

    const danger = mount(BaseButton, { props: { variant: 'danger' } })
    expect(danger.classes()).toContain('btn-danger')

    const secondary = mount(BaseButton, { props: { variant: 'secondary' } })
    expect(secondary.classes()).toContain('btn-secondary')
  })

  it('applies size classes', () => {
    const sm = mount(BaseButton, { props: { size: 'sm' } })
    expect(sm.classes()).toContain('btn-sm')

    const lg = mount(BaseButton, { props: { size: 'lg' } })
    expect(lg.classes()).toContain('btn-lg')
  })

  it('applies full width class', () => {
    const wrapper = mount(BaseButton, { props: { fullWidth: true } })
    expect(wrapper.classes()).toContain('w-full')
  })

  it('renders icon when provided', () => {
    const wrapper = mount(BaseButton, { props: { icon: 'arrow-right' } })
    expect(wrapper.find('[data-icon="arrow-right"]').exists()).toBe(true)
  })
})