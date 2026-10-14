import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import NotificationToast from '../../src/components/NotificationToast.vue'
import { useUserActionStore } from '@/stores/userAction'

const iconStubs = { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true }

describe('NotificationToast.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const mountToast = () => mount(NotificationToast, { global: { stubs: iconStubs } })

  it('does not render when activeNotification is null', () => {
    expect(mountToast().find('.fixed').exists()).toBe(false)
  })

  it('renders cart notification type and clears on close click', async () => {
    const store = useUserActionStore()
    store.activeNotification = { type: 'cart', message: 'Added to cart!' }

    const wrapper = mountToast()

    expect(wrapper.text()).toContain('Added to cart!')
    expect(wrapper.findComponent({ name: 'ShoppingCart' }).exists() || wrapper.find('svg').exists()).toBe(true)

    await wrapper.find('button').trigger('click')
    expect(store.activeNotification).toBeNull()
  })

  it.each([
    { type: 'wishlist', message: 'Added to wishlist!' },
    { type: 'success', message: 'Success action!' },
    { type: 'info', message: 'Info message!' }
  ] as const)('renders $type notification type', ({ type, message }) => {
    const store = useUserActionStore()
    store.activeNotification = { type, message }

    expect(mountToast().text()).toContain(message)
  })
})
