import { mount } from '@vue/test-utils'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import SellerNavbar from '../../src/components/SellerNavbar.vue'
import { useSellerStore } from '@/stores/seller'

vi.mock('@/stores/seller', () => ({
  useSellerStore: vi.fn()
}))

describe('SellerNavbar.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('loads seller data on mount and displays profile fallback', async () => {
    const mockLoadSeller = vi.fn().mockResolvedValue({})
    vi.mocked(useSellerStore).mockReturnValue({
      loadSeller: mockLoadSeller,
      account: null
    } as unknown as ReturnType<typeof useSellerStore>)

    const wrapper = mount(SellerNavbar, {
      global: {
        stubs: {
          Store: true,
          LayoutDashboard: true,
          PackagePlus: true,
          UserRound: true,
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    expect(mockLoadSeller).toHaveBeenCalled()
    expect(wrapper.text()).toContain('Profile')
    expect(wrapper.text()).toContain('KD5F Seller')
  })

  it('displays seller username when account is loaded', () => {
    const mockLoadSeller = vi.fn().mockResolvedValue({})
    vi.mocked(useSellerStore).mockReturnValue({
      loadSeller: mockLoadSeller,
      account: { username: 'CoolSeller' }
    } as unknown as ReturnType<typeof useSellerStore>)

    const wrapper = mount(SellerNavbar, {
      global: {
        stubs: {
          Store: true,
          LayoutDashboard: true,
          PackagePlus: true,
          UserRound: true,
          RouterLink: {
            template: '<a><slot /></a>'
          }
        }
      }
    })

    expect(wrapper.text()).toContain('CoolSeller')
  })
})