import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import ConnectionPanel from '../../src/components/ConnectionPanel.vue'

describe('ConnectionPanel', () => {
  it('emits databaseType with add payload', async () => {
    const wrapper = mount(ConnectionPanel, {
      props: {
        dbs: [],
        selectedName: null,
        loading: false,
      },
      global: {
        stubs: {
          ElButton: { template: '<button><slot /></button>' },
          ElDialog: { template: '<div><slot /><slot name="footer" /></div>' },
          ElForm: { template: '<form><slot /></form>' },
          ElFormItem: { template: '<div><slot /></div>' },
          ElInput: {
            props: ['modelValue'],
            emits: ['update:modelValue'],
            template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
          },
          ElSelect: {
            props: ['modelValue'],
            emits: ['update:modelValue'],
            template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><slot /></select>',
          },
          ElOption: { template: '<option />' },
        },
      },
    })

    await wrapper.get('button').trigger('click')
    await wrapper.find('input').setValue('interview_db')
    const selects = wrapper.findAll('select')
    await selects[0].setValue('mysql')
    await wrapper.findAll('input')[1].setValue('mysql://root:secret@localhost:3306/interview_db')
    await wrapper.findAll('button').at(-1)?.trigger('click')

    expect(wrapper.emitted('add')?.[0][0]).toEqual({
      name: 'interview_db',
      url: 'mysql://root:secret@localhost:3306/interview_db',
      databaseType: 'mysql',
    })
  })
})
