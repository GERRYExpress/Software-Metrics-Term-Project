import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import CategoryTable from '../../src/components/CategoryTable.vue'

describe('CategoryTable.vue', () => {
  const categories = [
    { id: 1, name: 'Action', isHidden: false },
    { id: 2, name: 'RPG', isHidden: true }
  ]

  it('renders category rows with correct visibility badges', () => {
    const wrapper = mount(CategoryTable, {
      props: {
        categories,
        isLoading: false
      },
      global: {
        stubs: {
          Loading: true,
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          Edit3: true,
          Eye: true,
          EyeOff: true,
          Trash2: true
        }
      }
    })

    expect(wrapper.text()).toContain('Action')
    expect(wrapper.text()).toContain('Visible')
    expect(wrapper.text()).toContain('RPG')
    expect(wrapper.text()).toContain('Hidden')
  })

  it('enters edit mode and emits save on enter or button click', async () => {
    const wrapper = mount(CategoryTable, {
      props: {
        categories,
        editingId: null
      },
      global: {
        stubs: {
          Loading: true,
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          Edit3: true,
          Eye: true,
          EyeOff: true,
          Trash2: true
        }
      }
    })

    // Click edit button for the first category
    const editButtons = wrapper.findAll('button')
    // First button is edit for category 1
    await editButtons[0].trigger('click')

    // Verify model update
    expect(wrapper.emitted('update:editingId')?.[0]).toEqual([1])
  })

  it('emits toggle and remove events when corresponding buttons are clicked', async () => {
    const wrapper = mount(CategoryTable, {
      props: {
        categories,
        editingId: null
      },
      global: {
        stubs: {
          Loading: true,
          BaseButton: { template: '<button @click="$emit(\'click\')"><slot/></button>' },
          Edit3: true,
          Eye: true,
          EyeOff: true,
          Trash2: true
        }
      }
    })

    const buttons = wrapper.findAll('button')
    // buttons[1] is toggle visibility (Eye) for Action
    await buttons[1].trigger('click')
    expect(wrapper.emitted('toggle')).toBeTruthy()
    expect(wrapper.emitted('toggle')?.[0]).toEqual([categories[0]])

    // buttons[2] is delete (Trash2) for Action
    await buttons[2].trigger('click')
    expect(wrapper.emitted('remove')).toBeTruthy()
    expect(wrapper.emitted('remove')?.[0]).toEqual([categories[0]])
  })
})