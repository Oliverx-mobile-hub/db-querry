import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import MetadataExplorer from '../../src/components/MetadataExplorer.vue'

describe('MetadataExplorer', () => {
  it('renders empty metadata', () => {
    const wrapper = mount(MetadataExplorer, { props: { metadata: { databaseType: 'postgres', schemas: [] } } })
    expect(wrapper.text()).toContain('0 objects')
  })
})

