import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import LandingView from '../../src/views/LandingView.vue'
import GameService from '../../src/service/GameService'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import type { AxiosResponse } from 'axios'
import type { Game } from '../../src/types'

vi.mock('@/service/GameService')
vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-router')>()
  return {
    ...actual,
  }
})

describe('LandingView.vue', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', name: 'landing-view', component: LandingView },
      { path: '/search', name: 'search-result-view', component: { template: '<div>Search</div>' } },
      { path: '/login', name: 'login-view', component: { template: '<div>Login</div>' } },
      { path: '/register', name: 'register-view', component: { template: '<div>Register</div>' } }
    ]
  })

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renders carousel banners, featured games grid, and checks localStorage for welcome modal', async () => {
    vi.mocked(GameService.getGames).mockResolvedValueOnce({
      data: [{ id: 1, name: 'Featured Game', price: 60, posterUrl: 'url' }] as Game[],
      headers: {}
    } as unknown as AxiosResponse<Game[]>)

    router.push('/')
    await router.isReady()

    const wrapper = mount(LandingView, {
      global: {
        plugins: [router],
        stubs: {
          GameCard: true,
          ProductMainContainer: { template: '<div><slot/></div>' },
          WelcomeModal: {
            template: '<div class="welcome-modal" @click="$emit(\'close\')">Welcome</div>',
            props: ['show']
          },
          ChevronLeft: true,
          ChevronRight: true,
          img: true
        }
      }
    })

    await new Promise(resolve => setTimeout(resolve, 10))

    expect(wrapper.text()).toContain('HIGHLIGHTS')
    expect(wrapper.text()).toContain('FEATURED GAMES')
    const gameCard = wrapper.findComponent({ name: 'GameCard' })
    expect(gameCard.props('game').name).toBe('Featured Game')
    expect(wrapper.find('.welcome-modal').exists()).toBe(true)

    await wrapper.find('.welcome-modal').trigger('click')
    expect(localStorage.getItem('kd5f_welcome_seen')).toBe('true')
  })
})