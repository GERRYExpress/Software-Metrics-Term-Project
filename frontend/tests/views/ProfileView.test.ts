import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import ProfileView from '../../src/views/ProfileView.vue'
import AccountService from '@/service/AccountService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'

vi.mock('@/service/AccountService')

describe('ProfileView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/profile', name: 'profile-view', component: ProfileView },
      { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
    ]
  })

  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    router.push('/profile')
    await router.isReady()
  })

  it('renders buyer profile details and handles logout and account deletion', async () => {
    const { useAuthStore } = await import('@/stores/auth')
    const authStore = useAuthStore()
    authStore.account = {
      id: 7,
      username: 'TestBuyer',
      createdDate: '2026-01-01T00:00:00.000Z',
      totalSpend: 120.5,
      totalGame: 4
    }

    vi.mocked(AccountService.deleteAccount).mockResolvedValueOnce({} as unknown as AxiosResponse<void>)

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router],
        stubs: {
          AppCard: { template: '<div><slot/></div>' },
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          ProductMainContainer: { template: '<div><slot/></div>' }
        }
      }
    })

    expect(wrapper.text()).toContain('TestBuyer')
    expect(wrapper.text()).toContain('$120.50')
    expect(wrapper.text()).toContain('4')

    const buttons = wrapper.findAll('button')
    // Click Logout button
    await buttons[0].trigger('click')
    expect(router.currentRoute.value.name).toBe('login-view')

    // Reset router navigation back to profile for delete test
    router.push('/profile')
    await router.isReady()

    // Click Delete My Account button to open modal
    await buttons[1].trigger('click')
    expect(wrapper.text()).toContain('Are you sure?')

    // Confirm deletion inside modal
    const modalConfirmButton = wrapper.findAll('button')[wrapper.findAll('button').length - 1]
    await modalConfirmButton.trigger('click')

    expect(AccountService.deleteAccount).toHaveBeenCalledWith(7)
  })
})