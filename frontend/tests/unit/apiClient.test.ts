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
})
