import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerDashboardView from '../../src/views/SellerDashboardView.vue'
import CategoryService from '@/service/CategoryService'
import GameService from '@/service/GameService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'

vi.mock('@/service/CategoryService')
vi.mock('@/service/GameService')

describe('SellerDashboardView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/seller', name: 'seller-dashboard-view', component: SellerDashboardView }]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads and displays seller dashboard categories and products', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({ data: [{ id: 1, name: 'Action', isHidden: false }] })
    vi.mocked(GameService.getGamesBySeller).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Seller Game', price: 40, posterUrl: 'img', description: 'desc' }],
      headers: { 'x-total-count': '1' }
    } as any)

    // Mock seller store state
    const { useSellerStore } = await import('@/stores/seller')
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce({ id: 10, username: 'seller1' } as any)

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: true,
          CategoryTable: true,
          Image: true,
          Plus: true
        }
      }
    })

    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('Your products')
    expect(wrapper.text()).toContain('Seller Game')
    expect(wrapper.text()).toContain('Categories')
  })
})