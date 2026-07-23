import { describe, expect, it, vi } from 'vitest'
import { api, ApiClientError } from '../../src/api/client'

describe('api client', () => {
  it('unwraps success envelope', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => ({ json: async () => ({ success: true, data: { dbs: [] }, error: null }) })))
    await expect(api.listDbs()).resolves.toEqual({ dbs: [] })
    vi.unstubAllGlobals()
  })

  it('throws api errors', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => ({ json: async () => ({ success: false, data: null, error: { code: 'invalidRequest', message: 'bad', details: {} } }) })))
    await expect(api.listDbs()).rejects.toBeInstanceOf(ApiClientError)
    vi.unstubAllGlobals()
  })

  it('sends databaseType when adding mysql db', async () => {
    const fetchMock = vi.fn(async () => ({ json: async () => ({ success: true, data: { db: {} }, error: null }) }))
    vi.stubGlobal('fetch', fetchMock)
    await api.putDb('interview_db', 'mysql://root:secret@localhost:3306/interview_db', 'mysql')
    const calls = fetchMock.mock.calls as unknown as Array<[RequestInfo | URL, RequestInit?]>
    const requestInit = calls[0]?.[1]
    expect(requestInit).toBeDefined()
    expect(JSON.parse((requestInit?.body as string))).toEqual({
      url: 'mysql://root:secret@localhost:3306/interview_db',
      databaseType: 'mysql',
    })
    vi.unstubAllGlobals()
  })
})
