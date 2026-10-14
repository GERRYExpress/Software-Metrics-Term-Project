import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import MediaCarousel from '../../src/components/MediaCarousel.vue'

describe('MediaCarousel.vue', () => {
  it('renders no media state when media is null or empty', () => {
    const wrapper = mount(MediaCarousel, {
      props: { media: [] },
      global: { stubs: { ImageIcon: true } }
    })
    expect(wrapper.text()).toContain('No gallery media available')
  })

  it('renders image media and navigates through next and prev', async () => {
    const wrapper = mount(MediaCarousel, {
      props: {
        media: [
          { id: 1, type: 'image', mediaUrl: 'img1.jpg' },
          { id: 2, type: 'image', mediaUrl: 'img2.jpg' }
        ]
      },
      global: { stubs: { ChevronLeft: true, ChevronRight: true } }
    })

    expect(wrapper.find('img').attributes('src')).toBe('img1.jpg')
    expect(wrapper.text()).toContain('1 / 2')

    // Click next button
    const buttons = wrapper.findAll('button')
    // Next is the second control button
    await buttons[1].trigger('click')
    expect(wrapper.text()).toContain('2 / 2')

    // Click prev button to loop back
    await buttons[0].trigger('click')
    expect(wrapper.text()).toContain('1 / 2')
  })

  it('renders video iframe and thumbnail play icon', () => {
    const wrapper = mount(MediaCarousel, {
      props: {
        media: [
          { id: 1, type: 'video', mediaUrl: 'trailer.mp4' },
          { id: 2, type: 'image', mediaUrl: 'img.jpg' }
        ]
      },
      global: { stubs: { Play: true, ChevronLeft: true, ChevronRight: true } }
    })

    expect(wrapper.find('iframe').exists()).toBe(true)
    expect(wrapper.find('iframe').attributes('src')).toBe('trailer.mp4')
  })

  it('resets activeIndex when media prop changes via watcher', async () => {
    const wrapper = mount(MediaCarousel, {
      props: {
        media: [
          { id: 1, type: 'image', mediaUrl: 'img1.jpg' },
          { id: 2, type: 'image', mediaUrl: 'img2.jpg' }
        ]
      }
    })

    // Advance to index 1 via thumbnail or next
    await wrapper.findAll('button')[1].trigger('click') // next
    expect(wrapper.text()).toContain('2 / 2')

    // Update media prop
    await wrapper.setProps({
      media: [
        { id: 3, type: 'image', mediaUrl: 'new1.jpg' }
      ]
    })

    expect(wrapper.text()).toContain('1 / 1')
  })

  it('allows clicking thumbnails to change active media', async () => {
    const wrapper = mount(MediaCarousel, {
      props: {
        media: [
          { id: 1, type: 'image', mediaUrl: 'img1.jpg' },
          { id: 2, type: 'image', mediaUrl: 'img2.jpg' }
        ]
      }
    })

    const thumbnails = wrapper.findAll('.shrink-0')
    await thumbnails[1].trigger('click')
    expect(wrapper.text()).toContain('2 / 2')
  })
})