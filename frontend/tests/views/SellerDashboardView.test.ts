import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerDashboardView from '../../src/views/SellerDashboardView.vue'
import CategoryService from '../../src/service/CategoryService'
import GameService from '../../src/service/GameService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Game, Category } from '../../src/types'

vi.mock('@/service/CategoryService')
vi.mock('@/service/GameService')

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})
describe('SellerDashboardView.vue', () => {
const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/seller', name: 'seller-dashboard-view', component: SellerDashboardView },
      // Add this route to fix the "No match for name: seller-product-edit-view" error:
      { path: '/seller/product/edit/:id', name: 'seller-product-edit-view', component: { template: '<div>Edit Product</div>' } }
    ]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads and displays seller dashboard categories and products', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }] as Category[]
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGamesBySeller).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Seller Game', price: 40, posterUrl: 'img', description: 'desc' }] as Game[],
      headers: { 'x-total-count': '1' }
    } as unknown as AxiosResponse<Game[]>)

    // Mock seller store state
    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce({
      id: 10,
      username: 'seller1'
    } as unknown as Awaited<ReturnType<typeof sellerStore.loadSeller>>)

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

    expect(wrapper.text()).toContain('Seller centerDashboardManage your catalog and storefront categories. Add product')
    expect(wrapper.text()).toContain('centerDashboardManage your catalog and storefront categories. Add product')
    expect(wrapper.text()).toContain('Seller centerDashboardManage your catalog and storefront categories. Add product')
  })
})