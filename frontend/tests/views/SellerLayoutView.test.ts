import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerLayoutView from '../../src/views/SellerLayoutView.vue'
import { createViewTestRouter, resetViewTestState, layoutStubs } from './viewTestUtils'

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})
describe('SellerLayoutView.vue', () => {
  const router = createViewTestRouter([
    { path: '/seller', name: 'seller-dashboard-view', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
  ])

  beforeEach(() => {
    resetViewTestState()
  })

  it('renders layout components and session expired modal when query parameter is set', async () => {
    router.push('/seller?expired=true')
    await router.isReady()

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce({
      id: 1,
      username: 'seller1',
      createdDate: '',
      totalGame: 0,
      totalSale: 0
    } as unknown as Awaited<ReturnType<typeof sellerStore.loadSeller>>)

    const wrapper = mount(SellerLayoutView, {
      global: {
        plugins: [router],
        stubs: layoutStubs('SellerNavbar')
      }
    })

    await wrapper.find('.session-modal').trigger('click')
    expect(wrapper.text()).toContain('Expired')
  })
})
