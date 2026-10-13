import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import BaseSelect from '@/components/BaseSelect.vue'
import type { Category } from '@/types'

const mockCategories: Category[] = [
  { id: 1, name: 'Action', isHidden: false },
  { id: 2, name: 'RPG', isHidden: false },
  { id: 3, name: 'Strategy', isHidden: true }
]

describe('BaseSelect', () => {
  it('renders options from categories', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories } })
    const options = wrapper.findAll('option')
    expect(options).toHaveLength(3)
    expect(options[0].text()).toBe('Action')
    expect(options[1].text()).toBe('RPG')
    expect(options[2].text()).toBe('Strategy')
  })

  it('renders label when provided', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories, label: 'Category' } })
    expect(wrapper.text()).toContain('Category')
  })

  it('forwards arbitrary attributes to select', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories }, attrs: { 'data-test': 'select', required: true } })
    expect(wrapper.find('select').attributes('data-test')).toBe('select')
    expect(wrapper.find('select').attributes('required')).toBeDefined()
  })

  it('renders slot content', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories }, slots: { default: '<span>Prefix</span>' } })
    expect(wrapper.html()).toContain('Prefix')
  })

  it('renders ChevronsUpDown icon', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories } })
    expect(wrapper.find('svg').exists()).toBe(true)
  })

  it('has select element with correct options', () => {
    const wrapper = mount(BaseSelect, { props: { options: mockCategories, modelValue: mockCategories[0] } })
    const select = wrapper.find('select')
    expect(select.exists()).toBe(true)
    expect(select.element.value).toBeDefined()
  })
})