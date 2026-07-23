import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import type { MetadataDocument } from '../../src/api/types'
import MetadataExplorer from '../../src/components/MetadataExplorer.vue'

describe('MetadataExplorer', () => {
  it('renders empty metadata', () => {
    const wrapper = mount(MetadataExplorer, {
      props: { dbName: 'local', metadata: { databaseType: 'postgres', schemas: [] } },
      global: {
        stubs: {
          ElEmpty: { template: '<div />' },
        },
      },
    })
    expect(wrapper.text()).toContain('0 objects')
  })

  it('renders database, table, and wrapping column details after expanding a table', async () => {
    const metadata: MetadataDocument = {
      databaseType: 'mysql',
      schemas: [{
        name: 'interview_db',
        objects: [{
          schema: 'interview_db',
          name: 'candidates',
          type: 'table',
          comment: '',
          columns: [{
            name: 'gender',
            dataType: "enum('female','male','non_binary','undisclosed')",
            nullable: false,
            primaryKey: false,
            ordinal: 1,
            comment: '',
          }],
        }],
      }],
    }
    const wrapper = mount(MetadataExplorer, {
      props: { dbName: 'interview_db', metadata },
      global: {
        stubs: {
          ElButton: { template: '<button><slot /></button>' },
          ElEmpty: { template: '<div />' },
          ElInput: { template: '<input />' },
        },
      },
    })

    expect(wrapper.text()).toContain('interview_db')
    expect(wrapper.text()).toContain('candidates')
    await wrapper.get('.object-row').trigger('click')
    expect(wrapper.text()).toContain('gender')
    expect(wrapper.text()).toContain("enum('female','male','non_binary','undisclosed')")
    expect(wrapper.text()).toContain('NOT NULL')
  })
})
