import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import PaymentModal from '../../src/components/PaymentModal.vue'
import { QrCode, Smartphone, X } from '@lucide/vue'

describe('PaymentModal.vue', () => {
  it('does not render when show is false', () => {
    const wrapper = mount(PaymentModal, {
      props: { show: false, totalAmount: 100 },
      global: { 
        stubs: { 
          BaseButton: true, 
          QrCode, 
          Smartphone, 
          X 
        } 
      }
    })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('renders QR code payment method by default and switches to Truemoney', async () => {
    const wrapper = mount(PaymentModal, {
      props: { show: true, totalAmount: 149.9 },
      global: { 
        stubs: { 
          BaseButton: {
            template: '<button @click="$emit(\'click\')"><slot/></button>'
          }, 
          QrCode, 
          Smartphone, 
          X 
        } 
      }
    })

    expect(wrapper.text()).toContain('Scan QR to Pay')
    expect(wrapper.text()).toContain('$149.90')

    // Click Truemoney tab button
    const tabs = wrapper.findAll('button')
    await tabs[2].trigger('click')

    expect(wrapper.text()).toContain('Truemoney Payment')
    expect(wrapper.text()).toContain('08X-XXX-XXXX')

    // Switch back to QR
    await tabs[1].trigger('click')
    expect(wrapper.text()).toContain('Scan QR to Pay')

    // Click close button (X)
    await tabs[0].trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()

    // Click proceed button (BaseButton)
    const baseButtons = wrapper.findAll('button')
    await baseButtons[baseButtons.length - 1].trigger('click')
    expect(wrapper.emitted('proceed')).toBeTruthy()
  })
})