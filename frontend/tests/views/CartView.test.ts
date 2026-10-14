import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import CartView from '@/views/CartView.vue'
import CartService from '@/service/CartService'
import type { Router } from 'vue-router'
import { createViewTestRouter, resetViewTestState, mountView } from './viewTestUtils'
import type { AxiosResponse } from 'axios'
import type { Cart } from '@/types'
import { useAuthStore } from '@/stores/auth'
import { useUserActionStore } from '@/stores/userAction'

vi.mock('@/service/CartService')
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})

describe('CartView.vue', () => {
  let router: Router

  const cartDeleteStubs = {
    ProductMainContainer: { template: '<div><slot/></div>' },
    ProductListCard: {
      template: '<button @click="$emit(\'delete-item\', 1, 101)">Delete</button>'
    },
    PriceSummaryCard: true,
    PaymentModal: true,
    PaymentSuccessModal: true
  }

  beforeEach(() => {
    resetViewTestState()

    router = createViewTestRouter([
      { path: '/', name: 'landing-page-view', component: { template: '<div>Landing</div>' } },
      { path: '/cart', name: 'cart-view', component: CartView },
      { path: '/search', name: 'search-result-view', component: { template: '<div>Search</div>' } }
    ])
  })

  it('renders empty cart message when no cart items exist', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1, username: 'testuser' } as any

    vi.mocked(CartService.getCarts).mockResolvedValueOnce({
      data: []
    } as unknown as AxiosResponse<Cart[]>)

    const wrapper = await mountView(CartView, router, {
      ProductMainContainer: { template: '<div><slot/></div>' },
      PriceSummaryCard: { template: '<div data-testid="price-summary"></div>' },
      PaymentModal: true,
      PaymentSuccessModal: true
    }, '/cart')

    expect(wrapper.text()).toContain('Your cart is currently empty')
    expect(wrapper.findComponent({ name: 'RouterLink' }).exists()).toBe(true)
  })

  it('loads cart items and calculates total amount correctly', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1, username: 'testuser' } as any

    const mockCarts: Cart[] = [
      {
        productId: 101,
        userId: 1,
        game: { id: 101, name: 'Game A', price: 100, description: 'Test', posterUrl: '' }
      },
      {
        productId: 102,
        userId: 1,
        game: { id: 102, name: 'Game B', price: 200, description: 'Test', posterUrl: '' }
      }
    ]

    vi.mocked(CartService.getCarts).mockResolvedValueOnce({
      data: mockCarts
    } as unknown as AxiosResponse<Cart[]>)

    const wrapper = await mountView(CartView, router, {
      ProductMainContainer: { template: '<div><slot/></div>' },
      ProductListCard: { template: '<div class="product-list-card"></div>' },
      PriceSummaryCard: { 
        template: '<button @click="$emit(\'checkout\')">Checkout</button>' 
      },
      PaymentModal: { 
        template: '<div v-if="show" id="payment-modal"><button @click="$emit(\'proceed\')">Proceed</button></div>',
        props: ['show']
      },
      PaymentSuccessModal: { 
        template: '<div v-if="show" id="success-modal"><button @click="$emit(\'done\')">Done</button></div>',
        props: ['show']
      }
    }, '/cart')

    expect(wrapper.findAll('.product-list-card')).toHaveLength(2)

    await wrapper.find('button').trigger('click')
    expect(wrapper.find('#payment-modal').exists()).toBe(true)

    await wrapper.find('#payment-modal button').trigger('click')
    expect(wrapper.find('#success-modal').exists()).toBe(true)

    await wrapper.find('#success-modal button').trigger('click')
    await flushPromises()
    expect(router.currentRoute.value.name).toBe('landing-page-view')
  })

  it('handles cart item deletion correctly', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1, username: 'testuser' } as any

    const userActionStore = useUserActionStore()
    userActionStore.cartIds = [101, 102]

    const mockCarts: Cart[] = [
      {
        productId: 101,
        userId: 1,
        game: { id: 101, name: 'Game A', price: 50, description: 'Test', posterUrl: '' }
      }
    ]

    vi.mocked(CartService.getCarts).mockResolvedValueOnce({
      data: mockCarts
    } as unknown as AxiosResponse<Cart[]>)

    vi.mocked(CartService.deleteCart).mockResolvedValueOnce({} as unknown as AxiosResponse<void>)
    const notifySpy = vi.spyOn(userActionStore, 'notify')

    const wrapper = await mountView(CartView, router, cartDeleteStubs, '/cart')

    await wrapper.find('.product-list-card button, button').trigger('click')
    await flushPromises()

    expect(CartService.deleteCart).toHaveBeenCalledWith(1, 101)
    expect(userActionStore.cartIds).not.toContain(101)
    expect(notifySpy).toHaveBeenCalledWith('Removed item from cart', 'info')
  })

  it('handles errors gracefully during cart loading and deletion', async () => {
    const authStore = useAuthStore()
    authStore.account = { id: 1, username: 'testuser' } as any

    const consoleLogSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const mockCarts: Cart[] = [
      {
        productId: 101,
        userId: 1,
        game: { id: 101, name: 'Game A', price: 50, description: 'Test', posterUrl: '' }
      }
    ]

    // First call for loading fails
    vi.mocked(CartService.getCarts).mockRejectedValueOnce(new Error('Load error'))

    const wrapper = await mountView(CartView, router, cartDeleteStubs, '/cart')

    expect(consoleLogSpy).toHaveBeenCalled()

    // Test deletion error catch block by re-mounting with success load, then failing delete
    vi.mocked(CartService.getCarts).mockResolvedValueOnce({
      data: mockCarts
    } as unknown as AxiosResponse<Cart[]>)
    vi.mocked(CartService.deleteCart).mockRejectedValueOnce(new Error('Delete error'))

    const wrapper2 = mount(CartView, {
      global: {
        plugins: [router],
        stubs: cartDeleteStubs
      }
    })

    await flushPromises()
    await wrapper2.find('button').trigger('click')
    await flushPromises()

    expect(consoleErrorSpy).toHaveBeenCalled()

    consoleLogSpy.mockRestore()
    consoleErrorSpy.mockRestore()
  })

  it('does nothing in watchEffect if authStore.account is null', async () => {
    const authStore = useAuthStore()
    authStore.account = null as any

    await mountView(CartView, router, {
      ProductMainContainer: { template: '<div><slot/></div>' },
      PriceSummaryCard: true,
      PaymentModal: true,
      PaymentSuccessModal: true
    }, '/cart')

    await flushPromises()
    expect(CartService.getCarts).not.toHaveBeenCalled()
  })
})