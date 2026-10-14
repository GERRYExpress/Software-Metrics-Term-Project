import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import WishlistView from '../../src/views/WishlistView.vue'
import WishlistService from '../../src/service/WishlistService'
import { setActivePinia, createPinia } from 'pinia'
import type { AxiosResponse } from 'axios'

vi.mock('@/service/WishlistService')

const ProductMainContainerStub = {
  name: 'ProductMainContainer',
  template: '<div class="product-main-container"><slot /></div>'
}

const ProductListCardStub = {
  name: 'ProductListCard',
  template: '<div class="product-list-card" @click="$emit(\'delete-item\', 1, 1)"><slot /></div>',
  props: ['product'],
  emits: ['delete-item']
}

const RouterLinkStub = {
  name: 'RouterLink',
  template: '<a class="router-link" :href="to"><slot /></a>',
  props: ['to']
}

const globalStubs = {
  global: {
    stubs: {
      ProductMainContainer: ProductMainContainerStub,
      ProductListCard: ProductListCardStub,
      RouterLink: RouterLinkStub
    }
  }
}

describe('WishlistView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders wishlist page with title', () => {
    const wrapper = mount(WishlistView, globalStubs)
    expect(wrapper.find('h1').text()).toBe('Wishlists')
  })

  it('shows empty state when no wishlists', () => {
    vi.mocked(WishlistService.getWishlists).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<unknown[]>)

    const wrapper = mount(WishlistView, globalStubs)
    expect(wrapper.text()).toContain('Looks like your wishlist is empty')
  })

  it('shows link to search results in empty state', () => {
    vi.mocked(WishlistService.getWishlists).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<unknown[]>)

    const wrapper = mount(WishlistView, globalStubs)
    expect(wrapper.text()).toContain('Find more games')
    expect(wrapper.findComponent(RouterLinkStub).exists()).toBe(true)
  })
})