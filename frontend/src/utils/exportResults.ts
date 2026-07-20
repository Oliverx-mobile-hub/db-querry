import type { ExportFile, ExportFormat, QueryResult } from '../api/types'

export function canExportResult(result: QueryResult | null, loading: boolean): boolean {
  return !loading && result !== null && !result.empty
}

export function safeDbName(dbName: string | null): string {
  const value = (dbName || 'query').trim().replace(/[^A-Za-z0-9_-]+/g, '-').replace(/^-+|-+$/g, '')
  return value || 'query'
}

export function formatExportTimestamp(date: Date): string {
  const pad = (value: number) => String(value).padStart(2, '0')
  return [
    date.getFullYear(),
    pad(date.getMonth() + 1),
    pad(date.getDate()),
    '-',
    pad(date.getHours()),
    pad(date.getMinutes()),
    pad(date.getSeconds()),
  ].join('')
}

export function exportFilename(dbName: string | null, format: ExportFormat, now = new Date()): string {
  return `${safeDbName(dbName)}-query-${formatExportTimestamp(now)}.${format}`
}

export function queryResultToCsv(result: QueryResult): string {
  const headers = result.columns.map(column => column.name)
  const lines = [
    headers.map(serializeCsvField).join(','),
    ...result.rows.map(row => headers.map(header => serializeCsvField(row[header])).join(',')),
  ]
  return lines.join('\r\n')
}

export function buildExportFile(result: QueryResult, format: ExportFormat, dbName: string | null, now = new Date()): ExportFile {
  if (format === 'csv') {
    return {
      filename: exportFilename(dbName, 'csv', now),
      mimeType: 'text/csv;charset=utf-8',
      content: queryResultToCsv(result),
    }
  }

  return {
    filename: exportFilename(dbName, 'json', now),
    mimeType: 'application/json;charset=utf-8',
    content: JSON.stringify(result, null, 2),
  }
}

export function downloadExportFile(file: ExportFile): void {
  const blob = new Blob([file.content], { type: file.mimeType })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = file.filename
  link.style.display = 'none'
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

function serializeCsvField(value: unknown): string {
  const text = valueToCsvText(value)
  if (/[",\r\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

function valueToCsvText(value: unknown): string {
  if (value === null || value === undefined) return ''
  if (value instanceof Date) return value.toISOString()
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}
