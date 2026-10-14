import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SellerProfileView from '../../src/views/SellerProfileView.vue'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})
describe('SellerProfileView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/seller/profile', name: 'seller-profile-view', component: SellerProfileView },
      { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
    ]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders seller profile details and handles logout', async () => {
    router.push('/seller/profile')
    await router.isReady()

    const { useSellerStore } = await import('../../src/stores/seller')
    const sellerStore = useSellerStore()
    vi.spyOn(sellerStore, 'loadSeller').mockResolvedValueOnce({
      id: 1,
      username: 'TestSeller',
      createdDate: '2026-01-01T00:00:00.000Z',
      totalGame: 3,
      totalSale: 150
    } as unknown as Awaited<ReturnType<typeof sellerStore.loadSeller>>)

    const wrapper = mount(SellerProfileView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          ProductMainContainer: { template: '<div><slot/></div>' }
        }
      }
    })

    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('ProfileLogoutDelete my account')
    expect(wrapper.text()).toContain('ProfileLogoutDelete my account')
    expect(wrapper.text()).toContain('ProfileLogoutDelete my account')

    const logoutButton = wrapper.findAll('button')[0]
    await logoutButton.trigger('click')

    expect(router.currentRoute.value.name).toBe('seller-profile-view')
  })
})