import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import PaymentSuccessModal from '../../src/components/PaymentSuccessModal.vue'

describe('PaymentSuccessModal.vue', () => {
  it('does not render when show is false', () => {
    const wrapper = mount(PaymentSuccessModal, {
      props: { show: false },
      global: { stubs: { CheckCircle2: true, BaseButton: true } }
    })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('renders and emits done event when button is clicked', async () => {
    const wrapper = mount(PaymentSuccessModal, {
      props: { show: true },
      global: {
        stubs: {
          CheckCircle2: true,
          BaseButton: {
            template: '<button @click="$emit(\'click\')"><slot/></button>'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('Payment Successful!')

    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('done')).toBeTruthy()
  })
})