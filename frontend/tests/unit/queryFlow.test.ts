import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import ResultTable from '../../src/components/ResultTable.vue'

describe('ResultTable', () => {
  it('shows row count', () => {
    const wrapper = mount(ResultTable, {
      props: {
        result: {
          columns: [{ name: 'id', dataType: 'integer' }],
          rows: [{ id: 1 }],
          rowCount: 1,
          durationMs: 3,
          limitApplied: false,
          limit: null,
          empty: false,
        },
      },
    })
    expect(wrapper.text()).toContain('1 rows')
  })
})
