import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import WelcomeModal from '../../src/components/WelcomeModal.vue'

describe('WelcomeModal.vue', () => {
  it('does not render when show is false', () => {
    const wrapper = mount(WelcomeModal, {
      props: { show: false },
      global: { stubs: { Gamepad2: true, X: true } }
    })

    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('renders and emits close when close button is clicked', async () => {
    const wrapper = mount(WelcomeModal, {
      props: { show: true },
      global: { stubs: { Gamepad2: true, X: true } }
    })

    expect(wrapper.text()).toContain('Welcome to')
    
    const closeButtons = wrapper.findAll('button')
    // Click the 'X' close button or start exploring button
    await closeButtons[0].trigger('click')

    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('emits close when clicking the backdrop self', async () => {
    const wrapper = mount(WelcomeModal, {
      props: { show: true },
      global: { stubs: { Gamepad2: true, X: true } }
    })

    await wrapper.find('.fixed').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})