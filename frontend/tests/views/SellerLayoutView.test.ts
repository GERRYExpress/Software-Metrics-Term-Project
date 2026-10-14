import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerLayoutView from '../../src/views/SellerLayoutView.vue'
import { useSellerStore } from '../../src/stores/seller'
import type { Account } from '@/types'
import type { Router } from 'vue-router'
import { createViewTestRouter, resetViewTestState, layoutStubs } from './viewTestUtils'

vi.mock('../../src/stores/seller', () => ({
  useSellerStore: vi.fn()
}))
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('SellerLayoutView.vue', () => {
  let router: Router

  beforeEach(() => {
    resetViewTestState()

    router = createViewTestRouter([
      { path: '/seller', name: 'seller-dashboard-view', component: { template: '<div>Dashboard</div>' } },
      { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
    ])
  })

  it('loads seller on mount successfully', async () => {
    const mockAccount = { id: 1, username: 'seller' } as unknown as Account
    const loadSellerMock = vi.fn().mockResolvedValue(mockAccount)
    vi.mocked(useSellerStore).mockReturnValue({
      loadSeller: loadSellerMock,
      clear: vi.fn()
    } as unknown as ReturnType<typeof useSellerStore>)

    router.push('/seller')
    await router.isReady()

    mount(SellerLayoutView, {
      global: {
        plugins: [router],
        stubs: layoutStubs('SellerNavbar')
      }
    })

    expect(loadSellerMock).toHaveBeenCalledTimes(1)
  })

  it('catches error if loadSeller fails on mount', async () => {
    const loadSellerMock = vi.fn().mockRejectedValue(new Error('Unauthorized'))
    vi.mocked(useSellerStore).mockReturnValue({
      loadSeller: loadSellerMock,
      clear: vi.fn()
    } as unknown as ReturnType<typeof useSellerStore>)

    router.push('/seller')
    await router.isReady()

    mount(SellerLayoutView, {
      global: {
        plugins: [router],
        stubs: layoutStubs('SellerNavbar')
      }
    })

    await flushPromises()
    expect(loadSellerMock).toHaveBeenCalledTimes(1)
  })

  it('handles session expired modal confirmation', async () => {
    const clearMock = vi.fn()
    const mockAccount = { id: 1, username: 'seller' } as unknown as Account
    vi.mocked(useSellerStore).mockReturnValue({
      loadSeller: vi.fn().mockResolvedValue(mockAccount),
      clear: clearMock
    } as unknown as ReturnType<typeof useSellerStore>)

    router.push('/seller?expired=true')
    await router.isReady()

    const wrapper = mount(SellerLayoutView, {
      global: {
        plugins: [router],
        stubs: {
          SellerNavbar: true,
          AppFooter: true,
          NotificationToast: true,
          RouterView: true
        }
      }
    })

    const modal = wrapper.findComponent({ name: 'SessionExpiredModal' })
    expect(modal.props('show')).toBe(true)

    await modal.vm.$emit('confirm')
    await flushPromises()

    expect(clearMock).toHaveBeenCalledTimes(1)
    expect(router.currentRoute.value.name).toBe('login-view')
  })
})
