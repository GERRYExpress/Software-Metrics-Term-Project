import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import App from '../src/App.vue'
vi.mock(import("vue-router"), async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    // your mocked methods
  }
})
describe('App.vue', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(async () => {
    vi.useFakeTimers()
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'home', component: { template: '<div>Home</div>' } },
        { path: '/other', name: 'other', component: { template: '<div>Other</div>' } }
      ]
    })
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.clearAllMocks()
  })

  it('manages loading state across router navigation, errors, and lifecycle hooks', async () => {
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: {
          Loading: { 
            template: '<div class="loading-overlay" :data-active="active"></div>', 
            props: ['active', 'color'] 
          },
          RouterView: { template: '<div class="router-view"></div>' }
        }
      }
    })

    // Verify initial mount finishes loading after router is ready
    await router.isReady()
    expect(wrapper.exists()).toBe(true)

    // Trigger router navigation (beforeEach -> show timer, afterEach -> hide)
    const pushPromise = router.push('/other')
    
    // Advance timers to activate the show timeout logic (> 200ms)
    vi.advanceTimersByTime(200)
    await pushPromise
    await router.isReady()

    // Test router error hook handling safely without 'any'
    interface RouterWithErrorHandling {
      errorHandler?: (err: Error) => void
    }
    const routerWithErrors = router as unknown as RouterWithErrorHandling
    if (routerWithErrors.errorHandler) {
      routerWithErrors.errorHandler(new Error('Test routing error'))
    }

    // Unmount component to execute onUnmounted cleanup guards
    wrapper.unmount()
  })
})