import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import CartView from '../../src/views/CartView.vue'
import CartService from '../../src/service/CartService'
import { setActivePinia, createPinia } from 'pinia'
import type { AxiosResponse } from 'axios'

vi.mock('@/service/CartService')

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

const PriceSummaryCardStub = {
  name: 'PriceSummaryCard',
  template: '<div class="price-summary-card" @click="$emit(\'checkout\')"><slot /></div>',
  props: ['games'],
  emits: ['checkout']
}

const PaymentModalStub = {
  name: 'PaymentModal',
  template: '<div class="payment-modal" v-if="show"><slot /></div>',
  props: ['show', 'totalAmount']
}

const PaymentSuccessModalStub = {
  name: 'PaymentSuccessModal',
  template: '<div class="payment-success-modal" v-if="show"><slot /></div>',
  props: ['show']
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
      PriceSummaryCard: PriceSummaryCardStub,
      PaymentModal: PaymentModalStub,
      PaymentSuccessModal: PaymentSuccessModalStub,
      RouterLink: RouterLinkStub
    }
  }
}

describe('CartView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders cart page with title', () => {
    const wrapper = mount(CartView, globalStubs)
    expect(wrapper.find('h1').text()).toBe('Carts')
  })

  it('shows empty state when no carts', () => {
    // FIX: Use unknown[] instead of any
    vi.mocked(CartService.getCarts).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<unknown[]>)

    const wrapper = mount(CartView, globalStubs)
    expect(wrapper.text()).toContain('Your cart is currently empty')
  })

  it('shows link to search results in empty state', () => {
    // FIX: Use unknown[] instead of any
    vi.mocked(CartService.getCarts).mockResolvedValue({
      data: []
    } as unknown as AxiosResponse<unknown[]>)

    const wrapper = mount(CartView, globalStubs)
    expect(wrapper.text()).toContain('Find more games')
    expect(wrapper.findComponent(RouterLinkStub).exists()).toBe(true)
  })

  it('has PriceSummaryCard component', () => {
    const wrapper = mount(CartView, globalStubs)
    expect(wrapper.findComponent(PriceSummaryCardStub).exists()).toBe(true)
  })

  it('has PaymentModal and PaymentSuccessModal components', () => {
    const wrapper = mount(CartView, globalStubs)
    expect(wrapper.findComponent(PaymentModalStub).exists()).toBe(true)
    expect(wrapper.findComponent(PaymentSuccessModalStub).exists()).toBe(true)
  })
})