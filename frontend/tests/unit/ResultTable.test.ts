import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import type { QueryResult } from '../../src/api/types'
import ResultTable from '../../src/components/ResultTable.vue'

vi.mock('../../src/utils/exportResults', async importOriginal => {
  const actual = await importOriginal<typeof import('../../src/utils/exportResults')>()
  return {
    ...actual,
    downloadExportFile: vi.fn(),
  }
})

const result: QueryResult = {
  columns: [{ name: 'id', dataType: 'int4' }],
  rows: [{ id: 1 }],
  rowCount: 1,
  durationMs: 8,
  limitApplied: false,
  limit: null,
  empty: false,
}

function mountResultTable(props: { result: QueryResult | null; loading: boolean; dbName?: string | null }) {
  return mount(ResultTable, {
    props: { dbName: 'local', ...props },
    global: {
      stubs: {
        ElButton: {
          props: ['disabled'],
          emits: ['click'],
          template: '<button :disabled="disabled" @click="$emit(\'click\', $event)"><slot /></button>',
        },
        ElEmpty: { template: '<div />' },
        ElTable: { template: '<table><slot /></table>' },
        ElTableColumn: { template: '<col />' },
      },
    },
  })
}

describe('ResultTable', () => {
  it('enables export buttons for non-empty results', () => {
    const wrapper = mountResultTable({ result, loading: false })

    expect(wrapper.get('[data-testid="export-csv"]').attributes('disabled')).toBeUndefined()
    expect(wrapper.get('[data-testid="export-json"]').attributes('disabled')).toBeUndefined()
  })

  it('disables export buttons without results', () => {
    const wrapper = mountResultTable({ result: null, loading: false })

    expect(wrapper.get('[data-testid="export-csv"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="export-json"]').attributes('disabled')).toBeDefined()
  })

  it('disables export buttons for empty results', () => {
    const wrapper = mountResultTable({ result: { ...result, empty: true, rows: [], rowCount: 0 }, loading: false })

    expect(wrapper.get('[data-testid="export-csv"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="export-json"]').attributes('disabled')).toBeDefined()
  })

  it('disables export buttons while querying', () => {
    const wrapper = mountResultTable({ result, loading: true })

    expect(wrapper.get('[data-testid="export-csv"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[data-testid="export-json"]').attributes('disabled')).toBeDefined()
  })
})
