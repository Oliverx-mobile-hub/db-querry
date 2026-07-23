import type { ApiResponse, DatabaseType, DbMetadataResponse, DbSummary, GeneratedSqlDraft, QueryResult } from './types'

const baseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/$/, '')

export class ApiClientError extends Error {
  constructor(public code: string, message: string, public details: Record<string, unknown> = {}) {
    super(message)
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers ?? {}),
    },
  })
  const body = (await response.json()) as ApiResponse<T>
  if (!body.success) {
    throw new ApiClientError(body.error.code, body.error.message, body.error.details)
  }
  return body.data
}

export const api = {
  listDbs: () => request<{ dbs: DbSummary[] }>('/api/v1/dbs'),
  putDb: (name: string, url: string, databaseType: DatabaseType) => request<{ db: DbSummary }>(`/api/v1/dbs/${encodeURIComponent(name)}`, { method: 'PUT', body: JSON.stringify({ url, databaseType }) }),
  deleteDb: (name: string) => request<{ deleted: boolean; name: string }>(`/api/v1/dbs/${encodeURIComponent(name)}`, { method: 'DELETE' }),
  getMetadata: (name: string) => request<DbMetadataResponse>(`/api/v1/dbs/${encodeURIComponent(name)}`),
  query: (name: string, sql: string) => request<QueryResult>(`/api/v1/dbs/${encodeURIComponent(name)}/query`, { method: 'POST', body: JSON.stringify({ sql }) }),
  generateSql: (name: string, prompt: string) => request<GeneratedSqlDraft>(`/api/v1/dbs/${encodeURIComponent(name)}/query/natural`, { method: 'POST', body: JSON.stringify({ prompt }) }),
}
