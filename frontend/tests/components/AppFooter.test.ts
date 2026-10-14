import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import AppFooter from '../../src/components/AppFooter.vue'

describe('AppFooter.vue', () => {
  it('renders copyright and contact link', () => {
    const wrapper = mount(AppFooter, {
      global: {
        stubs: {
          Copyright: true,
          InfoIcon: true,
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('2027 KaiDao-5-Fong')
    expect(wrapper.text()).toContain('CONTACT US')
  })
})