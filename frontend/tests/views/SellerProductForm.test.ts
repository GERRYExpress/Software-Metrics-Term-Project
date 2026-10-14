import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerProductFormView from '../../src/views/SellerProductFormView.vue'
import CategoryService from '@/service/CategoryService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Category } from '@/types'

vi.mock('@/service/CategoryService')
vi.mock('@/service/GameService')

describe('SellerProductFormView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/seller/product/add', name: 'seller-product-add-view', component: SellerProductFormView }]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('loads category options and renders form fields for creation', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'RPG', isHidden: false }] as Category[]
    } as unknown as AxiosResponse<Category[]>)
    
    const { useSellerStore } = await import('@/stores/seller')
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce({
      id: 5,
      username: 'seller5',
      createdDate: '',
      totalGame: 0,
      totalSale: 0
    } as unknown as Awaited<ReturnType<typeof sellerStore.loadSeller>>)

    router.push('/seller/product/add')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: true,
          GalleryUploadCard: true,
          ProductMainContainer: { template: '<div><slot/></div>' }
        }
      }
    })

    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('Add product')
    expect(wrapper.text()).toContain('RPG')
  })
})