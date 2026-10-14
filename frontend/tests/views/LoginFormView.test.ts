import { flushPromises } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import LoginFormView from '@/views/LoginFormView.vue'
import AccountService from '@/service/AccountService'
import type { Router } from 'vue-router'
import { createViewTestRouter, resetViewTestState, mountView } from './viewTestUtils'
import type { AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/auth'

vi.mock('@/service/AccountService')
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('LoginFormView.vue', () => {
  let router: Router

  const loginErrorStubs = {
    AuthMainContainer: { template: '<div><slot/></div>' },
    AppCard: { template: '<div><slot/></div>' },
    BaseInput: true,
    BaseToggle: true,
    BaseButton: { template: '<button type="submit"><slot/></button>' },
    UserRound: true
  }

  beforeEach(() => {
    resetViewTestState()

    router = createViewTestRouter([
      { path: '/', name: 'landing-page-view', component: { template: '<div>Home</div>' } },
      { path: '/login', name: 'login-view', component: LoginFormView },
      { path: '/register', name: 'register-view', component: { template: '<div>Register</div>' } },
      { path: '/seller', name: 'seller-dashboard-view', component: { template: '<div>Seller Dashboard</div>' } }
    ])
  })

  it('renders login form elements correctly', async () => {
    const wrapper = await mountView(LoginFormView, router, {
      AuthMainContainer: { template: '<div><slot/></div>' },
      AppCard: { template: '<div><slot/></div>' },
      BaseInput: { 
        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" :placeholder="placeholder" />',
        props: ['modelValue', 'placeholder']
      },
      BaseToggle: { template: '<button @click="$emit(\'toggleRole\')">Toggle</button>' },
      BaseButton: { template: '<button type="submit"><slot/></button>' },
      UserRound: true
    }, '/login')

    expect(wrapper.text()).toContain('Welcome Back!')
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.findAllComponents({ name: 'RouterLink' })).toHaveLength(2)
  })

  it('toggles role correctly and updates router query', async () => {
    const wrapper = await mountView(LoginFormView, router, {
      AuthMainContainer: { template: '<div><slot/></div>' },
      AppCard: { template: '<div><slot/></div>' },
      BaseInput: true,
      BaseToggle: { 
        template: '<button class="toggle-btn" @click="$emit(\'toggleRole\')">Toggle</button>' 
      },
      BaseButton: true,
      UserRound: true
    }, '/login')

    await wrapper.find('.toggle-btn').trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.query.role).toBe('true')
  })

  it('handles successful buyer login and navigates to landing page', async () => {
    const authStore = useAuthStore()
    const setSessionSpy = vi.spyOn(authStore, 'setSession')

    vi.mocked(AccountService.authenticate).mockResolvedValueOnce({
      data: {
        account: { id: 1, username: 'buyer1' },
        token: 'fake-jwt-token'
      }
    } as unknown as AxiosResponse)

    const wrapper = await mountView(LoginFormView, router, {
      AuthMainContainer: { template: '<div><slot/></div>' },
      AppCard: { template: '<div><slot/></div>' },
      BaseInput: { 
        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
        props: ['modelValue']
      },
      BaseToggle: true,
      BaseButton: { template: '<button type="submit"><slot/></button>' },
      UserRound: true
    }, '/login')

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('buyer1')
    await inputs[1].setValue('password123')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(AccountService.authenticate).toHaveBeenCalledWith({
      username: 'buyer1',
      password: 'password123',
      isSeller: false
    })
    expect(setSessionSpy).toHaveBeenCalledWith({ id: 1, username: 'buyer1' }, 'fake-jwt-token')
    expect(router.currentRoute.value.name).toBe('landing-page-view')
  })

  it('handles successful seller login and navigates to seller dashboard', async () => {
    vi.mocked(AccountService.authenticate).mockResolvedValueOnce({
      data: {
        account: { id: 2, username: 'seller1', isSeller: true },
        token: 'fake-jwt-token'
      }
    } as unknown as AxiosResponse)

    const wrapper = await mountView(LoginFormView, router, {
      AuthMainContainer: { template: '<div><slot/></div>' },
      AppCard: { template: '<div><slot/></div>' },
      BaseInput: { 
        template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
        props: ['modelValue']
      },
      BaseToggle: { 
        template: '<button class="toggle-btn" @click="$emit(\'toggleRole\')">Toggle</button>' 
      },
      BaseButton: { template: '<button type="submit"><slot/></button>' },
      UserRound: true
    }, '/login')

    await wrapper.find('.toggle-btn').trigger('click')

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('seller1')
    await inputs[1].setValue('sellerpass')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(router.currentRoute.value.name).toBe('seller-dashboard-view')
  })

  it('displays custom error message from API response on login failure', async () => {
    const apiError = {
      response: {
        data: {
          error: 'Invalid username or password.'
        }
      }
    }
    vi.mocked(AccountService.authenticate).mockRejectedValueOnce(apiError)

    const wrapper = await mountView(LoginFormView, router, loginErrorStubs, '/login')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    const alertBox = wrapper.find('[role="alert"]')
    expect(alertBox.exists()).toBe(true)
    expect(alertBox.text()).toBe('Invalid username or password.')
  })

  it('displays default fallback error message when API error has no response data', async () => {
    vi.mocked(AccountService.authenticate).mockRejectedValueOnce(new Error('Network Error'))

    const wrapper = await mountView(LoginFormView, router, loginErrorStubs, '/login')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    const alertBox = wrapper.find('[role="alert"]')
    expect(alertBox.exists()).toBe(true)
    expect(alertBox.text()).toBe('Unable to sign in. Check that the backend is running and try again.')
  })
})