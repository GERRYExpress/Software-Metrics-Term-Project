import { mount } from '@vue/test-utils'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import NotificationToast from '../../src/components/NotificationToast.vue'
import { useUserActionStore } from '@/stores/userAction'

describe('NotificationToast.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('does not render when activeNotification is null', () => {
    const wrapper = mount(NotificationToast, {
      global: { stubs: { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true } }
    })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('renders cart notification type and clears on close click', async () => {
    const store = useUserActionStore()
    store.activeNotification = { type: 'cart', message: 'Added to cart!' }

    const wrapper = mount(NotificationToast, {
      global: { stubs: { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true } }
    })

    expect(wrapper.text()).toContain('Added to cart!')
    expect(wrapper.findComponent({ name: 'ShoppingCart' }).exists() || wrapper.find('svg').exists()).toBe(true)

    await wrapper.find('button').trigger('click')
    expect(store.activeNotification).toBeNull()
  })

  it('renders wishlist notification type', () => {
    const store = useUserActionStore()
    store.activeNotification = { type: 'wishlist', message: 'Added to wishlist!' }

    const wrapper = mount(NotificationToast, {
      global: { stubs: { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true } }
    })
    expect(wrapper.text()).toContain('Added to wishlist!')
  })

  it('renders success notification type', () => {
    const store = useUserActionStore()
    store.activeNotification = { type: 'success', message: 'Success action!' }

    const wrapper = mount(NotificationToast, {
      global: { stubs: { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true } }
    })
    expect(wrapper.text()).toContain('Success action!')
  })

  it('renders info/default notification type', () => {
    const store = useUserActionStore()
    store.activeNotification = { type: 'info', message: 'Info message!' }

    const wrapper = mount(NotificationToast, {
      global: { stubs: { ShoppingCart: true, Heart: true, CheckCircle2: true, Info: true, X: true } }
    })
    expect(wrapper.text()).toContain('Info message!')
  })
})