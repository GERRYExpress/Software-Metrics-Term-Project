import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import AppNavItem from '@/components/AppNavItem.vue'

// Mock RouterLink as a global component
const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a class="router-link" :href="to"><slot /></a>',
  props: ['to']
}

const global = {
  stubs: {
    RouterLink: RouterLinkStub
  }
}

describe('AppNavItem', () => {
  const iconTests = [
    { name: 'home', routeName: 'home', iconName: 'Home' },
    { name: 'cart', routeName: 'cart', iconName: 'ShoppingCart' },
    { name: 'wishlist', routeName: 'wishlist', iconName: 'Heart' },
    { name: 'profile', routeName: 'profile', iconName: 'UserRound' }
  ]

  for (const { name, routeName, iconName } of iconTests) {
    it(`renders ${iconName} icon for name=${name}`, () => {
      const wrapper = mount(AppNavItem, { props: { name, routeName }, global })
      expect(wrapper.find('svg').exists()).toBe(true)
    })
  }

  it('links to correct route', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'cart', routeName: 'cart' }, global })
    const link = wrapper.findComponent(RouterLinkStub)
    expect(link.exists()).toBe(true)
    expect(link.props('to')).toEqual({ name: 'cart' })
  })

  it('shows count badge when count > 0', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'cart', routeName: 'cart', count: 3 }, global })
    const badge = wrapper.find('span.absolute')
    expect(badge.exists()).toBe(true)
    expect(badge.text()).toBe('3')
  })

  it('shows 9+ for count > 9', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'cart', routeName: 'cart', count: 15 }, global })
    const badge = wrapper.find('span.absolute')
    expect(badge.text()).toBe('9+')
  })

  it('hides badge when count is 0', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'cart', routeName: 'cart', count: 0 }, global })
    expect(wrapper.find('span.absolute').exists()).toBe(false)
  })

  it('hides badge when count is not provided', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'cart', routeName: 'cart' }, global })
    expect(wrapper.find('span.absolute').exists()).toBe(false)
  })

  it('applies base navigation classes', () => {
    const wrapper = mount(AppNavItem, { props: { name: 'home', routeName: 'home' }, global })
    const link = wrapper.findComponent(RouterLinkStub)
    expect(link.classes()).toContain('flex')
    expect(link.classes()).toContain('items-center')
    expect(link.classes()).toContain('justify-center')
    expect(link.classes()).toContain('bg-secondary')
    expect(link.classes()).toContain('rounded-full')
    expect(link.classes()).toContain('size-12')
    expect(link.classes()).toContain('shadow-md')
    expect(link.classes()).toContain('hover:scale-105')
    expect(link.classes()).toContain('hover:shadow-lg')
  })
})