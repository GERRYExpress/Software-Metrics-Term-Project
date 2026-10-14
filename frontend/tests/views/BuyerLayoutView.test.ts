import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import BuyerLayoutView from '../../src/views/BuyerLayoutView.vue'
import { createViewTestRouter, resetViewTestState, layoutStubs } from './viewTestUtils'

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})

describe('BuyerLayoutView.vue', () => {
  const router = createViewTestRouter([
    { path: '/', name: 'landing-view', component: { template: '<div>Landing</div>' } },
    { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } }
  ])

  beforeEach(() => {
    resetViewTestState()
  })

  it('renders layout components, calls fetchInitialData on mount, and handles session expiration modal', async () => {
    router.push('/?expired=true')
    await router.isReady()

    const { useUserActionStore } = await import('../../src/stores/userAction')
    const userActionStore = useUserActionStore()
    const fetchInitialDataSpy = vi.spyOn(userActionStore, 'fetchInitialData')

    const wrapper = mount(BuyerLayoutView, {
      global: {
        plugins: [router],
        stubs: layoutStubs('AppNavbar')
      }
    })

    expect(fetchInitialDataSpy).toHaveBeenCalled()
    await wrapper.find('.session-modal').trigger('click')
    await router.isReady() // Ensure router finishes navigation
    expect(router.currentRoute.value.name).toBe('landing-view')
  })
})
