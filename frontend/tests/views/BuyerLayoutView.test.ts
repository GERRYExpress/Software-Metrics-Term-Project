import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import BuyerLayoutView from '../../src/views/BuyerLayoutView.vue'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})

describe('BuyerLayoutView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'landing-view', component: { template: '<div>Landing</div>' } },
      { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
    ]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders layout components, calls fetchInitialData on mount, and handles session expiration modal', async () => {
    router.push('/?expired=true')
    await router.isReady()

    const { useUserActionStore } = await import('@/stores/userAction')
    const userActionStore = useUserActionStore()
    const fetchInitialDataSpy = vi.spyOn(userActionStore, 'fetchInitialData')

    const wrapper = mount(BuyerLayoutView, {
      global: {
        plugins: [router],
        stubs: {
          AppNavbar: true,
          AppFooter: true,
          SessionExpiredModal: {
            template: '<div class="session-modal" @click="$emit(\'confirm\')">Expired</div>',
            props: ['show']
          },
          NotificationToast: true,
          RouterView: true
        }
      }
    })

    expect(fetchInitialDataSpy).toHaveBeenCalled()

    await wrapper.find('.session-modal').trigger('click')
    expect(router.currentRoute.value.name).toBe('login-view')
  })
})