import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import IdentificationView from '@/views/IdentificationView.vue'
import { setActivePinia, createPinia } from 'pinia'

const AuthMainContainerStub = {
  name: 'AuthMainContainer',
  template: '<div class="auth-main-container"><slot /></div>'
}

const AppCardStub = {
  name: 'AppCard',
  template: '<div class="app-card"><slot /></div>'
}

const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a class="router-link" :href="to"><slot /></a>',
  props: ['to']
}

const ShoppingBagStub = { name: 'ShoppingBag', template: '<svg />' }
const StoreStub = { name: 'Store', template: '<svg />' }
const ArrowLeftStub = { name: 'ArrowLeft', template: '<svg />' }

const global = {
  stubs: {
    AuthMainContainer: AuthMainContainerStub,
    AppCard: AppCardStub,
    RouterLink: RouterLinkStub,
    ShoppingBag: ShoppingBagStub,
    Store: StoreStub,
    ArrowLeft: ArrowLeftStub
  }
}

describe('IdentificationView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders identification page with title', () => {
    const wrapper = mount(IdentificationView, { global })
    expect(wrapper.find('h1').text()).toBe('Sign In')
    expect(wrapper.find('h2').text()).toBe('Why are you here?')
  })

  it('shows back to home link', () => {
    const wrapper = mount(IdentificationView, { global })
    const backLink = wrapper.findComponent(RouterLinkStub)
    expect(backLink.exists()).toBe(true)
    expect(backLink.props('to')).toEqual({ name: 'landing-page-view' })
  })

  it('shows two role selection buttons', () => {
    const wrapper = mount(IdentificationView, { global })
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(2)
    expect(buttons[0].text()).toContain("I'm here to buy a game!")
    expect(buttons[1].text()).toContain("I'm here to sell a game!")
  })

  it('has correct styling on buyer button', () => {
    const wrapper = mount(IdentificationView, { global })
    const buttons = wrapper.findAll('button')
    expect(buttons[0].classes()).toContain('bg-action')
    expect(buttons[0].classes()).toContain('text-tertiary')
  })

  it('has correct styling on seller button', () => {
    const wrapper = mount(IdentificationView, { global })
    const buttons = wrapper.findAll('button')
    expect(buttons[1].classes()).toContain('bg-secondary')
    expect(buttons[1].classes()).toContain('border-white/30')
  })

  })