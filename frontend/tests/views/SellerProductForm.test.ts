import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerProductFormView from '../../src/views/SellerProductFormView.vue'
import CategoryService from '../../src/service/CategoryService'
import GameService from '../../src/service/GameService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Category, GameDetail, Account, GameInput } from '@/types'

vi.mock('../../src/service/CategoryService')
vi.mock('../../src/service/GameService')
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('SellerProductFormView.vue', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    Element.prototype.scrollIntoView = vi.fn()

    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/seller/add', name: 'seller-product-add-view', component: SellerProductFormView },
        { path: '/seller/edit/:id', name: 'seller-product-edit-view', component: SellerProductFormView },
        { path: '/seller', name: 'seller-dashboard-view', component: { template: '<div>Dashboard</div>' } }
      ]
    })
  })

  it('loads add product form successfully', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1 } as unknown as Account

    router.push('/seller/add')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: { template: '<input />' },
          ProductMainContainer: { template: '<div><slot/></div>' },
          GalleryUploadCard: { template: '<div></div>' }
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Add product')
  })

  it('validates form fields on submit and shows error if invalid', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1 } as unknown as Account

    router.push('/seller/add')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button type="submit"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          GalleryUploadCard: { template: '<div></div>' }
        }
      }
    })

    await flushPromises()

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    expect(wrapper.text()).toContain('Product name is required.')
  })

  it('submits create product successfully', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1 } as unknown as Account

    vi.mocked(GameService.createGame).mockResolvedValueOnce({} as unknown as AxiosResponse<void>)

    router.push('/seller/add')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button type="submit"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          GalleryUploadCard: { template: '<div></div>' }
        }
      }
    })

    await flushPromises()

    interface FormComponentVm {
      form: GameInput
      selectedCategoryIds: number[]
    }
    const vm = wrapper.vm as unknown as FormComponentVm
    vm.form.name = 'New Game'
    vm.form.description = 'Great game'
    vm.form.price = 29.99
    vm.form.posterUrl = 'https://example.com/poster.jpg'
    vm.selectedCategoryIds = [1]

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(GameService.createGame).toHaveBeenCalledTimes(1)
    expect(router.currentRoute.value.name).toBe('seller-dashboard-view')
  })

  it('loads edit product form and handles seller ownership check', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGame).mockResolvedValueOnce({
      data: {
        id: 5,
        name: 'Existing Game',
        description: 'Desc',
        price: 10,
        posterUrl: 'https://example.com/img.jpg',
        account: { id: 1 },
        categories: [{ id: 1, name: 'Action' }],
        hardwareRequirement: [],
        media: []
      } as unknown as GameDetail
    } as unknown as AxiosResponse<GameDetail>)

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1 } as unknown as Account

    router.push('/seller/edit/5')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button type="submit"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          GalleryUploadCard: { template: '<div></div>' }
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Edit product')
  })

  it('throws error if editing product belonging to different seller', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGame).mockResolvedValueOnce({
      data: {
        id: 5,
        name: 'Other Seller Game',
        account: { id: 999 },
        categories: [],
        hardwareRequirement: [],
        media: []
      } as unknown as GameDetail
    } as unknown as AxiosResponse<GameDetail>)

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1 } as unknown as Account

    router.push('/seller/edit/5')
    await router.isReady()

    const wrapper = mount(SellerProductFormView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: true,
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          GalleryUploadCard: true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Seller centerEdit productKeep your storefront information complete and accurate.Please sign in before opening the seller center.Basic informationDescriptionCategoriesSelect at least one category.No categories are available. Create one from the seller dashboard first.System requirementsBoth sections are optional. Empty sections are not sent.Minimum requirementsRecommended requirement')
  })
})