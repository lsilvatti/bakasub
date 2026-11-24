import { describe, expect, it, vi } from 'vitest'
import { OpenSubtitlesClient } from './opensubtitles-client'

type Call = { url: URL; method: string; headers: Record<string, string>; body: unknown }
type Route = { status?: number; body?: unknown }

function setup(
  routes: Record<string, Route>,
  apiKey: string | null = 'key-1234'
): { client: OpenSubtitlesClient; calls: Call[] } {
  const calls: Call[] = []
  const fetch = vi.fn(async (input: string | URL, init?: RequestInit) => {
    const url = new URL(String(input))
    calls.push({
      url,
      method: init?.method ?? 'GET',
      headers: (init?.headers ?? {}) as Record<string, string>,
      body: init?.body ? JSON.parse(String(init.body)) : undefined
    })
    const r = routes[url.pathname]
    if (!r) return new Response('{}', { status: 404 })
    return new Response(JSON.stringify(r.body ?? {}), { status: r.status ?? 200 })
  })
  return {
    client: new OpenSubtitlesClient({
      apiKey: async () => apiKey,
      userAgent: 'BakaSub v9.9.9',
      fetch: fetch as unknown as typeof globalThis.fetch
    }),
    calls
  }
}

const SEARCH = '/api/v1/subtitles'

describe('headers and credentials', () => {
  it('sends Api-Key and User-Agent; Bearer and the login host only with a session', async () => {
    const { client, calls } = setup({ [SEARCH]: { body: { data: [] } } })
    await client.search({ query: 'x' })
    expect(calls[0].url.host).toBe('api.opensubtitles.com')
    expect(calls[0].headers).toMatchObject({
      'Api-Key': 'key-1234',
      'User-Agent': 'BakaSub v9.9.9'
    })
    expect(calls[0].headers.Authorization).toBeUndefined()

    await client.search({ query: 'x' }, { token: 'jwt', host: 'vip-api.opensubtitles.com' })
    expect(calls[1].url.host).toBe('vip-api.opensubtitles.com')
    expect(calls[1].headers.Authorization).toBe('Bearer jwt')
  })

  it('without a key nothing reaches the network: OPENSUBTITLES_KEY_MISSING', async () => {
    const { client, calls } = setup({}, null)
    await expect(client.search({ query: 'x' })).rejects.toMatchObject({
      code: 'OPENSUBTITLES_KEY_MISSING'
    })
    expect(calls).toHaveLength(0)
  })

  it('validate uses the candidate key, not the stored one', async () => {
    const { client, calls } = setup({ [SEARCH]: { body: { data: [] } } }, null)
    await client.validate('candidate-key')
    expect(calls[0].headers['Api-Key']).toBe('candidate-key')
  })

  it('validate never repeats a query: the CDN would answer any key from its cache', async () => {
    const { client, calls } = setup({ [SEARCH]: { body: { data: [] } } })
    await client.validate('candidate-key')
    await client.validate('candidate-key')
    const [a, b] = calls.map((c) => c.url.searchParams.get('query'))
    expect(a).toMatch(/^bakasub-key-check-[0-9a-f-]{36}$/)
    expect(b).not.toBe(a)
  })

  it('401/403 = AUTH, 406 = QUOTA, anything else = NETWORK', async () => {
    for (const [status, code] of [
      [401, 'OPENSUBTITLES_AUTH'],
      [403, 'OPENSUBTITLES_AUTH'],
      [406, 'OPENSUBTITLES_QUOTA'],
      [429, 'NETWORK'],
      [500, 'NETWORK']
    ] as const) {
      const { client } = setup({ [SEARCH]: { status } })
      await expect(client.search({ query: 'x' })).rejects.toMatchObject({ code })
    }
  })
})

describe('login', () => {
  it('returns the token and the given host, when it belongs to the service', async () => {
    const ok = setup({
      '/api/v1/login': { body: { token: 'jwt', base_url: 'vip-api.opensubtitles.com' } }
    })
    expect(await ok.client.login('u', 'p')).toEqual({
      token: 'jwt',
      host: 'vip-api.opensubtitles.com'
    })
    expect(ok.calls[0]).toMatchObject({ method: 'POST', body: { username: 'u', password: 'p' } })

    const evil = setup({ '/api/v1/login': { body: { token: 'jwt', base_url: 'evil.example' } } })
    expect(await evil.client.login('u', 'p')).toEqual({ token: 'jwt' })
  })

  it('rejected login credentials (400 or 401) are AUTH; a 400 elsewhere is NETWORK', async () => {
    for (const status of [400, 401]) {
      const { client } = setup({ '/api/v1/login': { status } })
      await expect(client.login('me@example.com', 'p')).rejects.toMatchObject({
        code: 'OPENSUBTITLES_AUTH'
      })
    }
    const { client } = setup({ [SEARCH]: { status: 400 } })
    await expect(client.search({ query: 'x' })).rejects.toMatchObject({ code: 'NETWORK' })
  })

  it('a response without a token is NETWORK', async () => {
    const { client } = setup({ '/api/v1/login': { body: {} } })
    await expect(client.login('u', 'p')).rejects.toMatchObject({ code: 'NETWORK' })
  })
})

describe('search', () => {
  it('parameters sorted and lower-cased (the API redirects otherwise)', async () => {
    const { client, calls } = setup({ [SEARCH]: { body: { data: [] } } })
    await client.search({
      query: 'Sousou No Frieren',
      moviehash: 'ABCDEF0123456789',
      languages: ['pt-BR', 'en'],
      season: 1,
      episode: 5
    })
    expect(calls[0].url.search).toBe(
      '?episode_number=5&languages=en%2Cpt-br&moviehash=abcdef0123456789' +
        '&query=sousou+no+frieren&season_number=1'
    )
  })

  it('one result per file; missing or null fields are tolerated', async () => {
    const { client } = setup({
      [SEARCH]: {
        body: {
          data: [
            {
              attributes: {
                language: 'pt-BR',
                download_count: 42,
                hearing_impaired: true,
                fps: 23.976,
                from_trusted: true,
                ai_translated: false,
                machine_translated: true,
                moviehash_match: true,
                release: ' Show.S01E05.1080p ',
                uploader: { name: 'someone' },
                files: [
                  { file_id: 11, file_name: 'a.srt' },
                  { file_id: 12, file_name: 'b.srt' }
                ]
              }
            },
            { attributes: { language: null, release: null, files: [{ file_id: 13 }] } },
            { attributes: { files: null } }
          ]
        }
      }
    })
    const r = await client.search({ query: 'show' })
    expect(r.map((f) => f.fileId)).toEqual([11, 12, 13])
    expect(r[0]).toEqual({
      fileId: 11,
      release: 'Show.S01E05.1080p',
      language: 'pt-BR',
      downloads: 42,
      hashMatch: true,
      hearingImpaired: true,
      machineTranslated: true,
      trusted: true,
      fps: 23.976,
      uploader: 'someone'
    })
    expect(r[2]).toEqual({
      fileId: 13,
      release: '13',
      language: '',
      downloads: 0,
      hashMatch: false,
      hearingImpaired: false,
      machineTranslated: false,
      trusted: false
    })
  })

  it('an unexpected body is NETWORK', async () => {
    const { client } = setup({ [SEARCH]: { body: { data: 'nope' } } })
    await expect(client.search({ query: 'x' })).rejects.toMatchObject({ code: 'NETWORK' })
  })
})

describe('requestDownload', () => {
  it('asks for SRT and returns the link and the quota', async () => {
    const { client, calls } = setup({
      '/api/v1/download': {
        body: {
          link: 'https://www.opensubtitles.com/download/abc/file.srt',
          remaining: 4,
          reset_time_utc: '2026-11-25T00:00:00.000Z'
        }
      }
    })
    expect(await client.requestDownload(11)).toEqual({
      link: 'https://www.opensubtitles.com/download/abc/file.srt',
      remaining: 4,
      resetAt: '2026-11-25T00:00:00.000Z'
    })
    expect(calls[0]).toMatchObject({ method: 'POST', body: { file_id: 11, sub_format: 'srt' } })
  })
})
