import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import RegisterFormView from '@/views/RegisterFormView.vue'
import AccountService from '@/service/AccountService'
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

const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a class="router-link" :href="to"><slot /></a>',
  props: ['to']
}

const UserRoundPlusStub = { name: 'UserRoundPlus', template: '<svg />' }

const global = {
  stubs: {
    AuthMainContainer: AuthMainContainerStub,
    AppCard: AppCardStub,
    BaseButton: BaseButtonStub,
    BaseInput: BaseInputStub,
    RouterLink: RouterLinkStub,
    UserRoundPlus: UserRoundPlusStub
  }
}

describe('RegisterFormView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders register form with username, password, and confirm password inputs', () => {
    const wrapper = mount(RegisterFormView, { global })
    expect(wrapper.find('h1').text()).toBe('Create new account')
    expect(wrapper.findAllComponents(BaseInputStub)).toHaveLength(3)
    expect(wrapper.findComponent(BaseButtonStub).exists()).toBe(true)
  })

  it('shows link to login view', () => {
    const wrapper = mount(RegisterFormView, { global })
    expect(wrapper.findComponent(RouterLinkStub).exists()).toBe(true)
    expect(wrapper.text()).toContain('Sign in here')
  })

  it('has Register button with correct text', () => {
    const wrapper = mount(RegisterFormView, { global })
    const button = wrapper.findComponent(BaseButtonStub)
    expect(button.text()).toBe('Register')
  })

  it('has password input with type password', () => {
    const wrapper = mount(RegisterFormView, { global })
    const inputs = wrapper.findAllComponents(BaseInputStub)
    expect(inputs[1].props('type')).toBe('password')
    expect(inputs[2].props('type')).toBe('password')
  })

  it('shows back to home link', () => {
    const wrapper = mount(RegisterFormView, { global })
    const backLink = wrapper.findComponent(RouterLinkStub)
    expect(backLink.props('to')).toEqual({ name: 'landing-page-view' })
  })

  it('has correct labels on inputs', () => {
    const wrapper = mount(RegisterFormView, { global })
    const inputs = wrapper.findAllComponents(BaseInputStub)
    expect(inputs[0].props('label')).toBe('Username')
    expect(inputs[1].props('label')).toBe('Password')
    expect(inputs[2].props('label')).toBe('Confirm Password')
  })
})