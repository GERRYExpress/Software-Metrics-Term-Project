import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import AuthMainContainer from '../../src/components/AuthMainContainer.vue'

describe('AuthMainContainer.vue', () => {
  it('renders slot content inside a main tag', () => {
    const wrapper = mount(AuthMainContainer, {
      slots: {
        default: '<form>Login Form</form>'
      }
    })

    expect(wrapper.element.tagName).toBe('MAIN')
    expect(wrapper.text()).toContain('Login Form')
  })
})