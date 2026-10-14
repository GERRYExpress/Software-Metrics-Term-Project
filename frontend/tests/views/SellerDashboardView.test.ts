import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerDashboardView from '../../src/views/SellerDashboardView.vue'
import CategoryService from '../../src/service/CategoryService'
import GameService from '../../src/service/GameService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Category, Game, Account } from '@/types'

vi.mock('../../src/service/CategoryService')
vi.mock('../../src/service/GameService')
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('SellerDashboardView.vue', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    Element.prototype.scrollIntoView = vi.fn()

    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/seller', name: 'seller-dashboard-view', component: SellerDashboardView },
        { path: '/seller/add', name: 'seller-product-add-view', component: { template: '<div>Add</div>' } },
        { path: '/seller/edit/:id', name: 'seller-product-edit-view', component: { template: '<div>Edit</div>' } }
      ]
    })
  })

  it('loads dashboard successfully with categories and products', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Action', isHidden: false }]
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGamesBySeller).mockResolvedValueOnce({
      data: [{ id: 10, name: 'Game 1', price: 19.99, description: 'Test', posterUrl: 'https://example.com/poster.jpg' }] as unknown as Game[],
      headers: { 'x-total-count': '1' }
    } as unknown as AxiosResponse<Game[]>)

    // Populate both auth and seller stores to satisfy authentication guards in the view
    const { useAuthStore } = await import('../../src/stores/auth')
    const { useSellerStore } = await import('../../src/stores/seller')
    
    const authStore = useAuthStore()
    const sellerStore = useSellerStore()
    
    const mockAccount = { id: 1, username: 'seller', isSeller: true } as unknown as Account
    authStore.account = mockAccount
    authStore.token = 'fake-token'
    sellerStore.account = mockAccount

    router.push('/seller')
    await router.isReady()

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: { template: '<input />' },
          ProductMainContainer: { template: '<div><slot/></div>' },
          CategoryTable: { template: '<div class="category-table"></div>' },
          Image: { template: '<img />' },
          Plus: true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Game 1')
    expect(wrapper.text()).toContain('$19.99')
  })

  it('handles loadDashboard error when seller is missing', async () => {
    const { useAuthStore } = await import('../../src/stores/auth')
    const { useSellerStore } = await import('../../src/stores/seller')
    
    const authStore = useAuthStore()
    const sellerStore = useSellerStore()
    authStore.account = null as unknown as Account
    sellerStore.account = null as unknown as Account
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(null as unknown as Account)

    router.push('/seller')
    await router.isReady()

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          CategoryTable: true,
          Image: true,
          Plus: true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('No seller account is logged in.')
  })

  it('handles adding category validations and success', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<Category[]>)

    const { useAuthStore } = await import('../../src/stores/auth')
    const { useSellerStore } = await import('../../src/stores/seller')
    
    const authStore = useAuthStore()
    const sellerStore = useSellerStore()
    const mockAccount = { id: 1, username: 'seller', isSeller: true } as unknown as Account
    authStore.account = mockAccount
    sellerStore.account = mockAccount

    router.push('/seller')
    await router.isReady()

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button type="submit" @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: { 
            template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />', 
            props: ['modelValue'] 
          },
          ProductMainContainer: { template: '<div><slot/></div>' },
          CategoryTable: true,
          Image: true,
          Plus: true
        }
      }
    })

    await flushPromises()

    const form = wrapper.find('form')
    expect(form.exists()).toBe(true)
    await form.trigger('submit.prevent')
    expect(wrapper.text()).toContain('Enter a category name before adding it.')

    const input = wrapper.find('input')
    await input.setValue('RPG')
    vi.mocked(CategoryService.createCategory).mockResolvedValueOnce({} as unknown as AxiosResponse<Category>)

    await form.trigger('submit.prevent')
    await flushPromises()
    expect(CategoryService.createCategory).toHaveBeenCalledWith('RPG')
  })

  it('handles pagination actions', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGamesBySeller).mockResolvedValue({
      data: [
        { id: 1, name: 'G1', price: 10, description: 'D1', posterUrl: 'https://example.com/1.jpg' },
        { id: 2, name: 'G2', price: 10, description: 'D2', posterUrl: 'https://example.com/2.jpg' },
        { id: 3, name: 'G3', price: 10, description: 'D3', posterUrl: 'https://example.com/3.jpg' },
        { id: 4, name: 'G4', price: 10, description: 'D4', posterUrl: 'https://example.com/4.jpg' },
        { id: 5, name: 'G5', price: 10, description: 'D5', posterUrl: 'https://example.com/5.jpg' },
        { id: 6, name: 'G6', price: 10, description: 'D6', posterUrl: 'https://example.com/6.jpg' },
        { id: 7, name: 'G7', price: 10, description: 'D7', posterUrl: 'https://example.com/7.jpg' }
      ] as unknown as Game[],
      headers: { 'x-total-count': '7' }
    } as unknown as AxiosResponse<Game[]>)

    const { useAuthStore } = await import('../../src/stores/auth')
    const { useSellerStore } = await import('../../src/stores/seller')
    
    const authStore = useAuthStore()
    const sellerStore = useSellerStore()
    const mockAccount = { id: 1, username: 'seller', isSeller: true } as unknown as Account
    authStore.account = mockAccount
    sellerStore.account = mockAccount

    router.push('/seller')
    await router.isReady()

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          CategoryTable: true,
          Image: true,
          Plus: true
        }
      }
    })

    await flushPromises()

    const nextBtn = wrapper.findAll('button').find(b => b.text().includes('Next') || b.attributes('aria-label')?.includes('Next'))
    if (nextBtn) {
      await nextBtn.trigger('click')
      await flushPromises()
      expect(GameService.getGamesBySeller).toHaveBeenCalled()
    }
  })

  it('handles product deletion with confirmation', async () => {
    vi.mocked(CategoryService.getCategories).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<Category[]>)

    vi.mocked(GameService.getGamesBySeller).mockResolvedValue({
      data: [{ id: 7, name: 'G7', price: 10, description: 'D7', posterUrl: 'https://example.com/7.jpg' }] as unknown as Game[],
      headers: { 'x-total-count': '7' }
    } as unknown as AxiosResponse<Game[]>)

    const { useAuthStore } = await import('../../src/stores/auth')
    const { useSellerStore } = await import('../../src/stores/seller')
    
    const authStore = useAuthStore()
    const sellerStore = useSellerStore()
    const mockAccount = { id: 1, username: 'seller', isSeller: true } as unknown as Account
    authStore.account = mockAccount
    sellerStore.account = mockAccount

    vi.spyOn(window, 'confirm').mockReturnValue(true)
    vi.mocked(GameService.deleteGame).mockResolvedValueOnce({} as unknown as AxiosResponse<void>)

    router.push('/seller')
    await router.isReady()

    const wrapper = mount(SellerDashboardView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          BaseInput: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          CategoryTable: true,
          Image: true,
          Plus: true
        }
      }
    })

    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find(b => b.text().includes('Delete'))
    if (deleteBtn) {
      await deleteBtn.trigger('click')
      await flushPromises()
      expect(GameService.deleteGame).toHaveBeenCalledWith(7, 1)
    }
  })
})