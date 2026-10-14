import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import GalleryUploadCard from '../../src/components/GalleryUploadCard.vue'

describe('GalleryUploadCard.vue', () => {
  it('renders up to maxItems slots and pads with empty inputs', () => {
    const wrapper = mount(GalleryUploadCard, {
      props: {
        modelValue: ['https://example.com/img1.jpg'],
        maxItems: 3
      },
      global: {
        stubs: {
          Trash2: true,
          Plus: true
        }
      }
    })

    expect(wrapper.text()).toContain('Add up to 3 image or video URLs')
    expect(wrapper.text()).toContain('Media 1')
    expect(wrapper.text()).toContain('Media 2')
    expect(wrapper.text()).toContain('Media 3')
    expect(wrapper.findAll('input').length).toBe(3)
  })

  it('emits updated modelValue when typing in an input', async () => {
    const wrapper = mount(GalleryUploadCard, {
      props: {
        modelValue: [],
        maxItems: 2
      },
      global: {
        stubs: {
          Trash2: true,
          Plus: true
        }
      }
    })

    const input = wrapper.findAll('input')[0]
    await input.setValue('https://example.com/new.jpg')

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([['https://example.com/new.jpg']])
  })

  it('emits updated modelValue when removing a slot', async () => {
    const wrapper = mount(GalleryUploadCard, {
      props: {
        modelValue: ['https://example.com/img1.jpg', 'https://example.com/img2.jpg'],
        maxItems: 2
      },
      global: {
        stubs: {
          Trash2: { template: '<button @click="$emit(\'click\')">Remove</button>' },
          Plus: true
        }
      }
    })

    const removeButton = wrapper.find('button')
    await removeButton.trigger('click')

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
  })

  it('hides image element on error event', async () => {
    const wrapper = mount(GalleryUploadCard, {
      props: {
        modelValue: ['https://example.com/broken.jpg'],
        maxItems: 2
      },
      global: {
        stubs: {
          Trash2: true,
          Plus: true
        }
      }
    })

    const img = wrapper.find('img')
    expect(img.exists()).toBe(true)

    await img.trigger('error')
    expect((img.element as HTMLElement).style.display).toBe('none')
  })
})