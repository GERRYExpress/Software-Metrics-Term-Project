import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import AppNavbar from '../../src/components/AppNavbar.vue'
import CategoryService from '@/service/CategoryService.ts'

vi.mock('@/service/CategoryService.ts', () => ({
  default: {
    getCategories: vi.fn().mockResolvedValue({
      data: [
        { id: 1, name: 'Action' },
        { id: 2, name: 'RPG' }
      ]
    })
  }
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  }),
  useRoute: () => ({
    query: {}
  })
}))

describe('AppNavbar.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('fetches categories on mount', async () => {
    mount(AppNavbar, {
      global: {
        stubs: {
          Search: true,
          BaseInput: true,
          BaseSelect: true,
          AppNavItem: true,
          RouterLink: true
        }
      }
    })

    // Allow microtasks from onMounted to flush
    await vi.dynamicImportSettled?.() || new Promise(resolve => setTimeout(resolve, 10))

    expect(CategoryService.getCategories).toHaveBeenCalled()
  })
})