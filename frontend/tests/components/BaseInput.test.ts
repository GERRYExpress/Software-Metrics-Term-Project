import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import BaseInput from '@/components/BaseInput.vue'

describe('BaseInput', () => {
  it('renders label when provided', () => {
    const wrapper = mount(BaseInput, { props: { label: 'Email' } })
    expect(wrapper.find('label').text()).toBe('Email')
  })

  it('does not render label when not provided', () => {
    const wrapper = mount(BaseInput)
    expect(wrapper.find('label').exists()).toBe(false)
  })

  it('renders placeholder', () => {
    const wrapper = mount(BaseInput, { props: { placeholder: 'Enter email' } })
    expect(wrapper.find('input').attributes('placeholder')).toBe('Enter email')
  })

  it('renders type', () => {
    const wrapper = mount(BaseInput, { props: { type: 'password' } })
    expect(wrapper.find('input').attributes('type')).toBe('password')
  })

  it('updates v-model on input', async () => {
    const wrapper = mount(BaseInput, { props: { modelValue: '' } })
    const input = wrapper.find('input')
    await input.setValue('test@example.com')
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['test@example.com'])
  })

  it('forwards arbitrary attributes to input', () => {
    const wrapper = mount(BaseInput, { attrs: { 'data-test': 'input', autocomplete: 'email' } })
    expect(wrapper.find('input').attributes('data-test')).toBe('input')
    expect(wrapper.find('input').attributes('autocomplete')).toBe('email')
  })

  it('binds label for to input id', () => {
    const wrapper = mount(BaseInput, { props: { label: 'Username' } })
    const label = wrapper.find('label')
    const input = wrapper.find('input')
    expect(label.attributes('for')).toBe(input.attributes('id'))
  })
})