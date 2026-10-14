import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import SystemRequirement from '../../src/components/SystemRequirement.vue'

describe('SystemRequirement.vue', () => {
  it('renders fallback hyphens when hardware prop is empty', () => {
    const wrapper = mount(SystemRequirement)
    const values = wrapper.findAll('.font-bold.text-end')

    expect(values).toHaveLength(9)
    values.forEach(val => {
      expect(val.text()).toBe('-')
    })
  })

  it('renders hardware specs correctly when provided', () => {
    const wrapper = mount(SystemRequirement, {
      props: {
        hardware: {
          operatingSystem: 'Windows 11',
          processor: 'Intel Core i7',
          memory: '16 GB RAM',
          graphicsCard: 'NVIDIA RTX 3060',
          directXVersion: 'DX 12',
          network: 'Broadband',
          storage: '50 GB available space',
          soundcard: 'DirectX compatible',
          additionalNote: 'None'
        }
      }
    })

    expect(wrapper.text()).toContain('Windows 11')
    expect(wrapper.text()).toContain('Intel Core i7')
    expect(wrapper.text()).toContain('16 GB RAM')
    expect(wrapper.text()).toContain('NVIDIA RTX 3060')
  })
})