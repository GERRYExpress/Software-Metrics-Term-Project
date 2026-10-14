import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import ProductMainContainer from '../../src/components/ProductMainContainer.vue'

describe('ProductMainContainer.vue', () => {
  it('renders slot content inside a main element', () => {
    const wrapper = mount(ProductMainContainer, {
      slots: {
        default: '<div class="product-grid">Product List</div>'
      }
    })

    expect(wrapper.element.tagName).toBe('MAIN')
    expect(wrapper.text()).toContain('Product List')
  })
})