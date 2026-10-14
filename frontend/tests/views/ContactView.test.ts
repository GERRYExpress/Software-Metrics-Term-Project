import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import ContactView from '../../src/views/ContactView.vue'
import { setActivePinia, createPinia } from 'pinia'

const ProductMainContainerStub = {
  name: 'ProductMainContainer',
  template: '<div class="product-main-container"><slot /></div>'
}

const AppCardStub = {
  name: 'AppCard',
  template: '<div class="app-card"><slot /></div>'
}

const global = {
  stubs: {
    ProductMainContainer: ProductMainContainerStub,
    AppCard: AppCardStub
  }
}

describe('ContactView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders contact page with title', () => {
    const wrapper = mount(ContactView, { global })
    expect(wrapper.find('h1').text()).toBe('Contacts')
  })

  it('shows address information', () => {
    const wrapper = mount(ContactView, { global })
    expect(wrapper.text()).toContain('College of Arts, Media, and Technology')
    expect(wrapper.text()).toContain('Chiang Mai University')
    expect(wrapper.text()).toContain('Huay Kaew Rd')
    expect(wrapper.text()).toContain('053-920299')
  })

  it('shows social/contact cards', () => {
    const wrapper = mount(ContactView, { global })
    expect(wrapper.text()).toContain('KaiDao5Fong2.0')
    expect(wrapper.text()).toContain('@KaiDao5Fong2.0')
    expect(wrapper.text()).toContain('admin@kd5f.com')
  })

  it('has Google Maps iframe', () => {
    const wrapper = mount(ContactView, { global })
    const iframe = wrapper.find('iframe')
    expect(iframe.exists()).toBe(true)
    expect(iframe.attributes('src')).toContain('google.com/maps')
  })

  it('renders two AppCard components', () => {
    const wrapper = mount(ContactView, { global })
    expect(wrapper.findAllComponents(AppCardStub)).toHaveLength(2)
  })
})