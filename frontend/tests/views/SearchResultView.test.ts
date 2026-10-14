import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import SearchResultView from '../../src/views/SearchResultView.vue'
import GameService from '../../src/service/GameService'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Game } from '../../src/types'

vi.mock('@/service/GameService', () => ({
  default: {
    getGames: vi.fn(),
    getGamesByKeyword: vi.fn(),
    getGamesByCategory: vi.fn(),
    getGamesByKeywordAndCategory: vi.fn()
  }
}))
vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})

describe('SearchResultView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [{ path: '/search', name: 'search-result-view', component: SearchResultView }]
  })

  beforeEach(async () => {
    vi.clearAllMocks()
    router.push('/search?keyword=cyber&page=1')
    await router.isReady()
  })

  it('renders search results when games are found', async () => {
    vi.mocked(GameService.getGamesByKeyword).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Cyber Game', price: 50, posterUrl: 'url' }] as Game[],
      headers: { 'x-total-count': '1' }
    } as unknown as AxiosResponse<Game[]>)

    const wrapper = mount(SearchResultView, {
      props: { page: 1, category: '' },
      global: {
        plugins: [router],
        stubs: {
          GameCard: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          BaseInput: true
        }
      }
    })

    // Wait for async onMounted/watchEffect
    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('Search Results: cyber')
    expect(GameService.getGamesByKeyword).toHaveBeenCalledWith(12, 1, 'cyber')
  })

  it('renders empty state when no games are found', async () => {
    vi.mocked(GameService.getGamesByKeyword).mockResolvedValueOnce({
      data: [] as Game[],
      headers: { 'x-total-count': '0' }
    } as unknown as AxiosResponse<Game[]>)

    const wrapper = mount(SearchResultView, {
      props: { page: 1, category: '' },
      global: {
        plugins: [router],
        stubs: {
          GameCard: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          BaseInput: true
        }
      }
    })

    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('Congratulation! You found nothing.')
  })
})