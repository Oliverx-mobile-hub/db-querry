import { describe, expect, it } from 'vitest'
import type { QueryResult } from '../../src/api/types'
import { buildExportFile, canExportResult, exportFilename, queryResultToCsv, safeDbName } from '../../src/utils/exportResults'

const result: QueryResult = {
  columns: [
    { name: 'id', dataType: 'int4' },
    { name: 'comma', dataType: 'text' },
    { name: 'quote', dataType: 'text' },
    { name: 'multiline', dataType: 'text' },
    { name: 'missing', dataType: 'text' },
    { name: 'amount', dataType: 'numeric' },
    { name: 'active', dataType: 'bool' },
    { name: 'created_at', dataType: 'timestamptz' },
  ],
  rows: [{
    id: 1,
    comma: 'Alice, Bob',
    quote: 'He said "Hello"',
    multiline: 'line1\nline2',
    missing: null,
    amount: 42,
    active: true,
    created_at: '2026-07-20T14:30:00Z',
  }],
  rowCount: 1,
  durationMs: 12,
  limitApplied: false,
  limit: null,
  empty: false,
}

describe('exportResults', () => {
  it('serializes CSV with stable column order and escaping', () => {
    expect(queryResultToCsv(result)).toBe([
      'id,comma,quote,multiline,missing,amount,active,created_at',
      '1,"Alice, Bob","He said ""Hello""","line1\nline2",,42,true,2026-07-20T14:30:00Z',
    ].join('\r\n'))
  })

  it('builds full JSON export file', () => {
    const file = buildExportFile(result, 'json', 'local', new Date(2026, 6, 20, 14, 30, 0))

    expect(file.filename).toBe('local-query-20260720-143000.json')
    expect(file.mimeType).toBe('application/json;charset=utf-8')
    expect(JSON.parse(file.content)).toEqual(result)
  })

  it('builds CSV export file names with safe db names', () => {
    expect(safeDbName('local/prod db')).toBe('local-prod-db')
    expect(exportFilename('local/prod db', 'csv', new Date(2026, 6, 20, 14, 30, 0))).toBe('local-prod-db-query-20260720-143000.csv')
  })

  it('detects export availability', () => {
    expect(canExportResult(result, false)).toBe(true)
    expect(canExportResult(null, false)).toBe(false)
    expect(canExportResult({ ...result, empty: true }, false)).toBe(false)
    expect(canExportResult(result, true)).toBe(false)
  })
})
