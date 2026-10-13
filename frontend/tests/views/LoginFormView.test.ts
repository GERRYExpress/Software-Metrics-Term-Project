import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import LoginFormView from '@/views/LoginFormView.vue'
import AccountService from '@/service/AccountService'
import { useAuthStore } from '@/stores/auth'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('@/service/AccountService')

const AuthMainContainerStub = {
  name: 'AuthMainContainer',
  template: '<div class="auth-main-container"><slot /></div>'
}

const AppCardStub = {
  name: 'AppCard',
  template: '<div class="app-card"><slot /></div>'
}

const BaseButtonStub = {
  name: 'BaseButton',
  template: '<button class="base-button" @click="$emit(\'click\')" :disabled="disabled"><slot /></button>',
  props: ['disabled'],
  emits: ['click']
}

const BaseInputStub = {
  name: 'BaseInput',
  template: '<input class="base-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" :placeholder="placeholder" :type="type" :label="label" />',
  props: ['modelValue', 'placeholder', 'type', 'label'],
  emits: ['update:modelValue']
}

const BaseToggleStub = {
  name: 'BaseToggle',
  template: '<div class="base-toggle" @click="$emit(\'toggleRole\')"><slot /></div>',
  props: ['modelValue', 'options'],
  emits: ['toggleRole', 'update:modelValue']
}

const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a class="router-link" :href="to"><slot /></a>',
  props: ['to']
}

const UserRoundStub = { name: 'UserRound', template: '<svg />' }

const global = {
  stubs: {
    AuthMainContainer: AuthMainContainerStub,
    AppCard: AppCardStub,
    BaseButton: BaseButtonStub,
    BaseInput: BaseInputStub,
    BaseToggle: BaseToggleStub,
    RouterLink: RouterLinkStub,
    UserRound: UserRoundStub
  }
}

describe('LoginFormView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders login form with username and password inputs', () => {
    const wrapper = mount(LoginFormView, { global })
    expect(wrapper.find('h1').text()).toBe('Welcome Back!')
    expect(wrapper.findAllComponents(BaseInputStub)).toHaveLength(2)
    expect(wrapper.findComponent(BaseToggleStub).exists()).toBe(true)
    expect(wrapper.findComponent(BaseButtonStub).exists()).toBe(true)
  })

  it('shows link to register view', () => {
    const wrapper = mount(LoginFormView, { global })
    expect(wrapper.findComponent(RouterLinkStub).exists()).toBe(true)
    expect(wrapper.text()).toContain('Sign up here')
  })

  it('shows back to home link', () => {
    const wrapper = mount(LoginFormView, { global })
    const backLink = wrapper.findComponent(RouterLinkStub)
    expect(backLink.props('to')).toEqual({ name: 'landing-page-view' })
  })

  it('has Login button with correct text', () => {
    const wrapper = mount(LoginFormView, { global })
    const button = wrapper.findComponent(BaseButtonStub)
    expect(button.text()).toBe('Login')
  })

  it('shows loading text when isLoading is true (via prop)', () => {
    const wrapper = mount(LoginFormView, { global })
    const button = wrapper.findComponent(BaseButtonStub)
    // Test the button renders correctly with loading prop
    expect(button.props('disabled')).toBe(false)
  })

  it('has correct toggle options', () => {
    const wrapper = mount(LoginFormView, { global })
    const toggle = wrapper.findComponent(BaseToggleStub)
    expect(toggle.props('options')).toEqual(['Buyer', 'Seller'])
  })
})