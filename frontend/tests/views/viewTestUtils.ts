import { createPinia, setActivePinia } from 'pinia'
import { vi } from 'vitest'
import { createRouter, createWebHistory, type RouteRecordRaw, Router } from 'vue-router'
import { mount, flushPromises } from '@vue/test-utils'

// Shared scaffolding for view tests: real-router factory, common component
// stubs, and Pinia reset. Replaces the copy-pasted setup blocks.
export function createViewTestRouter(routes: RouteRecordRaw[]) {
  return createRouter({ history: createWebHistory(), routes })
}

export const viewStubs = {
  AppCard: { template: '<div><slot/></div>' },
  BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
  ProductMainContainer: { template: '<div><slot/></div>' }
}

export function resetViewTestState() {
  setActivePinia(createPinia())
  vi.clearAllMocks()
}

// Shared layout stubs: session-expired modal exists once instead of
// copy-pasted in every layout suite. Navbar differs per layout.
export function layoutStubs(navbar: string) {
  return {
    [navbar]: true,
    AppFooter: true,
    SessionExpiredModal: {
      template: '<div class="session-modal" @click="$emit(\'confirm\')">Expired</div>',
      props: ['show']
    },
    NotificationToast: true,
    RouterView: true
  }
}

export async function mountView(component: any, router: Router, stubs: Record<string, any>, path: string) {
  router.push(path)
  await router.isReady()
  const wrapper = mount(component, { global: { plugins: [router], stubs } })
  await flushPromises()
  return wrapper
}
