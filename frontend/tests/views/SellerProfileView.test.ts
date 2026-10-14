import { mount, flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerProfileView from '@/views/SellerProfileView.vue'
import AccountService from '@/service/AccountService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import { useSellerStore } from '@/stores/seller'
import { useAuthStore } from '@/stores/auth'
import { useUserActionStore } from '@/stores/userAction'

vi.mock('@/service/AccountService')
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('SellerProfileView.vue', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()

    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } },
        { path: '/seller/profile', name: 'seller-profile-view', component: SellerProfileView }
      ]
    })
  })

  it('renders seller profile details correctly upon successful load', async () => {
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockImplementationOnce(async () => {
      sellerStore.account = {
        id: 1,
        username: 'TestSeller',
        createdDate: new Date('2025-01-15'),
        totalGame: 5,
        totalSale: 1234.5,
        isSeller: true,
        isBuyer: false,
        password: '',
        totalSpend: 0
      } as any
      return sellerStore.account!
    })

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('TestSeller')
    expect(wrapper.text()).toContain('15 January 2025')
    expect(wrapper.text()).toContain('5')
    expect(wrapper.text()).toContain('$1234.50')
  })

  it('renders loading state when sellerStore.isLoading is true', async () => {
    const sellerStore = useSellerStore()
    sellerStore.account = null
    sellerStore.isLoading = true
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(null as any)

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: true,
          BaseButton: true
        }
      }
    })

    await flushPromises()
    expect(wrapper.text()).toContain('Loading seller profile...')
  })

  it('handles error when loading seller profile fails with custom response error', async () => {
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockRejectedValueOnce({
      response: { data: { error: 'Custom load error' } }
    })

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: true,
          BaseButton: true
        }
      }
    })

    await flushPromises()
    const alertBox = wrapper.find('[role="alert"]')
    expect(alertBox.exists()).toBe(true)
    expect(alertBox.text()).toBe('Custom load error')
  })

  it('handles error when loading seller profile fails with standard error message', async () => {
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockRejectedValueOnce(new Error('Standard error'))

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: true,
          BaseButton: true
        }
      }
    })

    await flushPromises()
    const alertBox = wrapper.find('[role="alert"]')
    expect(alertBox.text()).toBe('Standard error')
  })

  it('handles logout flow correctly', async () => {
    const sellerStore = useSellerStore()
    const authStore = useAuthStore()
    const userActionStore = useUserActionStore()

    const clearSellerSpy = vi.spyOn(sellerStore, 'clear')
    const logoutAuthSpy = vi.spyOn(authStore, 'logout')
    const notifySpy = vi.spyOn(userActionStore, 'notify')

    sellerStore.account = { id: 1, username: 'Seller' } as any
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(sellerStore.account as any)

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await flushPromises()

    const logoutBtn = wrapper.findAll('button').find(b => b.text().includes('Logout'))
    await logoutBtn?.trigger('click')
    await flushPromises()

    expect(clearSellerSpy).toHaveBeenCalled()
    expect(logoutAuthSpy).toHaveBeenCalled()
    expect(notifySpy).toHaveBeenCalledWith('Logged out successfully.', 'info')
    expect(router.currentRoute.value.name).toBe('login-view')
    expect(router.currentRoute.value.query.role).toBe('true')
  })

  it('opens and closes delete account modal', async () => {
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1, username: 'Seller', createdDate: new Date() } as any
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(sellerStore.account as any)

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find(b => b.text().includes('Delete my account'))
    await deleteBtn?.trigger('click')
    expect(wrapper.find('[role="dialog"]').exists()).toBe(true)

    const cancelBtn = wrapper.findAll('button').find(b => b.text().includes('Cancel'))
    await cancelBtn?.trigger('click')
    expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
  })

  it('handles successful account deletion', async () => {
    const sellerStore = useSellerStore()
    const userActionStore = useUserActionStore()
    const clearSellerSpy = vi.spyOn(sellerStore, 'clear')
    const notifySpy = vi.spyOn(userActionStore, 'notify')

    sellerStore.account = { id: 1, username: 'Seller', createdDate: new Date() } as any
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(sellerStore.account as any)
    vi.mocked(AccountService.deleteAccount).mockResolvedValueOnce({} as unknown as AxiosResponse<void>)

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await flushPromises()

    await wrapper.findAll('button').find(b => b.text().includes('Delete my account'))?.trigger('click')
    
    const confirmDeleteBtn = wrapper.findAll('button').find(b => b.text().includes('Yes, delete'))
    await confirmDeleteBtn?.trigger('click')
    await flushPromises()

    expect(AccountService.deleteAccount).toHaveBeenCalledWith(1)
    expect(clearSellerSpy).toHaveBeenCalled()
    expect(notifySpy).toHaveBeenCalledWith('Seller account deleted.', 'info')
    expect(router.currentRoute.value.name).toBe('login-view')
  })

  it('handles account deletion failure error gracefully', async () => {
    const sellerStore = useSellerStore()
    sellerStore.account = { id: 1, username: 'Seller', createdDate: new Date() } as any
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce(sellerStore.account as any)
    vi.mocked(AccountService.deleteAccount).mockRejectedValueOnce({
      response: { data: { error: 'Delete failed' } }
    })

    router.push('/seller/profile')
    await router.isReady()

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await flushPromises()

    await wrapper.findAll('button').find(b => b.text().includes('Delete my account'))?.trigger('click')
    const confirmDeleteBtn = wrapper.findAll('button').find(b => b.text().includes('Yes, delete'))
    await confirmDeleteBtn?.trigger('click')
    await flushPromises()

    const alertBox = wrapper.find('[role="alert"]')
    expect(alertBox.exists()).toBe(true)
    expect(alertBox.text()).toBe('Delete failed')
  })

  it('returns early in deleteAccount if account is missing', async () => {
    const sellerStore = useSellerStore()
    sellerStore.account = null

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          ProductMainContainer: { template: '<div><slot/></div>' },
          AppCard: true,
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' }
        }
      }
    })

    await wrapper.vm.deleteAccount()
    expect(AccountService.deleteAccount).not.toHaveBeenCalled()
  })
})