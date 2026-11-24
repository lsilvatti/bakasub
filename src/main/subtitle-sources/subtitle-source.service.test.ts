import { existsSync, mkdtempSync, readdirSync, realpathSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AppError } from '../errors'
import type { OpenSubtitlesAuth, OpenSubtitlesFile } from './opensubtitles-client'
import { isAllowedDownloadLink, SubtitleSourceService } from './subtitle-source.service'

const SRT = '1\n00:00:01,000 --> 00:00:02,000\nHello there.\n'
const LINK = 'https://www.opensubtitles.com/download/abc/file.srt'

const file = (fileId: number, over: Partial<OpenSubtitlesFile> = {}): OpenSubtitlesFile => ({
  fileId,
  release: `release-${fileId}`,
  language: 'en',
  downloads: 0,
  hashMatch: false,
  hearingImpaired: false,
  machineTranslated: false,
  trusted: false,
  ...over
})

let dir: string
let video: string
beforeEach(() => {
  dir = realpathSync.native(mkdtempSync(join(tmpdir(), 'bakasub-os-')))
  video = join(dir, '[Group] Some Show - S02E07 [1080p].mkv')
  writeFileSync(video, 'mkv')
})

function setup(o: {
  byHash?: OpenSubtitlesFile[]
  byQuery?: OpenSubtitlesFile[]
  hash?: string | null
  link?: string
  body?: string | Uint8Array
  status?: number
}): {
  service: SubtitleSourceService
  search: ReturnType<typeof vi.fn>
  requestDownload: ReturnType<typeof vi.fn>
  fetch: ReturnType<typeof vi.fn>
  cacheDir: string
} {
  const search = vi.fn(async (p: { moviehash?: string }) =>
    p.moviehash ? (o.byHash ?? []) : (o.byQuery ?? [])
  )
  const requestDownload = vi.fn(async () => ({ link: o.link ?? LINK, remaining: 9 }))
  const fetch = vi.fn(
    async () => new Response((o.body ?? SRT) as BodyInit, { status: o.status ?? 200 })
  )
  const cacheDir = join(dir, 'cache')
  return {
    service: new SubtitleSourceService({
      client: { search, requestDownload } as never,
      session: {
        run: (call: (a?: OpenSubtitlesAuth) => Promise<unknown>) => call(undefined)
      } as never,
      cacheDir,
      hash: async () => (o.hash === undefined ? 'abcdef0123456789' : o.hash),
      fetch: fetch as unknown as typeof globalThis.fetch
    }),
    search,
    requestDownload,
    fetch,
    cacheDir
  }
}

describe('search', () => {
  it('no query: title, season and episode come from the file name; searches by hash and by title', async () => {
    const { service, search } = setup({})
    const r = await service.search({ videoPath: video, languages: ['pt-BR', 'pt-br', 'en'] })
    expect(r).toMatchObject({ query: 'Some Show', season: 2, episode: 7, results: [] })
    expect(search).toHaveBeenCalledTimes(2)
    expect(search.mock.calls.map((c) => c[0])).toEqual([
      { moviehash: 'abcdef0123456789', languages: ['pt-br', 'en'] },
      { query: 'Some Show', languages: ['pt-br', 'en'], season: 2, episode: 7 }
    ])
  })

  it('with a user query, season and episode are only the requested ones', async () => {
    const { service, search } = setup({})
    const r = await service.search({ videoPath: video, query: ' Other Title ', languages: [] })
    expect(r.query).toBe('Other Title')
    expect(r).not.toHaveProperty('season')
    expect(search.mock.calls[1][0]).toEqual({
      query: 'Other Title',
      languages: [],
      season: undefined,
      episode: undefined
    })
  })

  it('no hash (small or unreadable file): searches by title only', async () => {
    const { service, search } = setup({ hash: null })
    await service.search({ videoPath: video, languages: [] })
    expect(search).toHaveBeenCalledTimes(1)
  })

  it('merges without duplicates, flags hash matches and ranks the most useful first', async () => {
    const { service } = setup({
      byHash: [file(1, { downloads: 3 })],
      byQuery: [
        file(1, { downloads: 3 }),
        file(2, { downloads: 900, machineTranslated: true }),
        file(3, { downloads: 10 }),
        file(4, { downloads: 5, trusted: true }),
        file(5, { downloads: 50, language: 'pt-br' })
      ]
    })
    const { results } = await service.search({ videoPath: video, languages: [] })
    expect(results.map((r) => r.fileId)).toEqual([1, 4, 5, 3, 2])
    expect(results[0].hashMatch).toBe(true)
    expect(results[1].hashMatch).toBe(false)
    expect(results.find((r) => r.fileId === 5)?.language).toBe('pt-BR')
  })

  it('the video goes through the path guard', async () => {
    const { service, search } = setup({})
    await expect(
      service.search({ videoPath: join(dir, 'nope.mkv'), languages: [] })
    ).rejects.toMatchObject({ code: 'INVALID_PATH' })
    await expect(
      service.search({ videoPath: join(dir, 'x.srt'), languages: [] })
    ).rejects.toMatchObject({ code: 'INVALID_PATH', message: 'extension' })
    expect(search).not.toHaveBeenCalled()
  })
})

describe('search cancellation', () => {
  /** Client and session that behave like the real ones: both honour the signal they are given. */
  function abortable(): {
    service: SubtitleSourceService
    search: ReturnType<typeof vi.fn>
    run: ReturnType<typeof vi.fn>
  } {
    const search = vi.fn(
      (_p: unknown, _a: unknown, signal?: AbortSignal) =>
        new Promise<OpenSubtitlesFile[]>((_resolve, reject) =>
          signal?.addEventListener('abort', () =>
            reject(new AppError('CANCELLED', 'OpenSubtitles request cancelled'))
          )
        )
    )
    const run = vi.fn((call: (a?: OpenSubtitlesAuth) => Promise<unknown>) => call(undefined))
    const service = new SubtitleSourceService({
      client: { search, requestDownload: vi.fn() } as never,
      session: { run } as never,
      cacheDir: join(dir, 'cache'),
      hash: async () => 'abcdef0123456789'
    })
    return { service, search, run }
  }

  it('aborting mid-flight rejects CANCELLED; the signal reaches both requests and both runs', async () => {
    const { service, search, run } = abortable()
    const ctrl = new AbortController()
    const pending = service.search({ videoPath: video, languages: ['en'] }, ctrl.signal)
    await vi.waitFor(() => expect(search).toHaveBeenCalledTimes(2))
    ctrl.abort()
    await expect(pending).rejects.toMatchObject({ code: 'CANCELLED' })
    expect(search.mock.calls.map((c) => c[2])).toEqual([ctrl.signal, ctrl.signal])
    expect(run.mock.calls.map((c) => c[1])).toEqual([ctrl.signal, ctrl.signal])
  })

  it('an already aborted signal sends nothing', async () => {
    const { service, search, run } = abortable()
    await expect(
      service.search({ videoPath: video, languages: ['en'] }, AbortSignal.abort())
    ).rejects.toMatchObject({ code: 'CANCELLED' })
    expect(run).not.toHaveBeenCalled()
    expect(search).not.toHaveBeenCalled()
  })
})

describe('download', () => {
  it('downloads, validates and stores under the id; the second request spends no quota', async () => {
    const { service, requestDownload, fetch, cacheDir } = setup({})
    const first = await service.download(77)
    expect(first).toEqual({
      path: join(cacheDir, '77.srt'),
      name: '77.srt',
      format: 'srt',
      cached: false,
      remaining: 9
    })
    expect(readdirSync(cacheDir)).toEqual(['77.srt'])

    const second = await service.download(77)
    expect(second).toEqual({ path: first.path, name: '77.srt', format: 'srt', cached: true })
    expect(requestDownload).toHaveBeenCalledTimes(1)
    expect(fetch).toHaveBeenCalledTimes(1)
  })

  it('a link outside the OpenSubtitles domain is not fetched', async () => {
    const { service, fetch } = setup({ link: 'https://evil.example/file.srt' })
    await expect(service.download(1)).rejects.toMatchObject({ code: 'NETWORK' })
    expect(fetch).not.toHaveBeenCalled()
  })

  it('content the parser rejects is not cached', async () => {
    const { service, cacheDir } = setup({ body: new Uint8Array([0, 1, 2, 3, 0, 0, 255]) })
    await expect(service.download(5)).rejects.toMatchObject({ code: 'INVALID_SUBTITLE' })
    expect(existsSync(cacheDir) ? readdirSync(cacheDir) : []).toEqual([])
  })

  it('an empty response or an HTTP error fails without leaving a file', async () => {
    const empty = setup({ body: '' })
    await expect(empty.service.download(5)).rejects.toMatchObject({ code: expect.any(String) })
    const gone = setup({ status: 410 })
    await expect(gone.service.download(5)).rejects.toMatchObject({ code: 'NETWORK' })
    expect(existsSync(gone.cacheDir)).toBe(false)
  })
})

describe('concurrent downloads', () => {
  /** Holds the fetch until `release()` so calls overlap for real. */
  function gated(): ReturnType<typeof setup> & { release: () => void } {
    const ctx = setup({})
    let release!: () => void
    const gate = new Promise<void>((r) => (release = r))
    ctx.fetch.mockImplementation(async () => {
      await gate
      return new Response(SRT, { status: 200 })
    })
    return { ...ctx, release }
  }

  it('two simultaneous calls for the same id make one request and get the same result', async () => {
    const { service, requestDownload, fetch, release, cacheDir } = gated()
    const a = service.download(77)
    const b = service.download(77)
    release()
    const [ra, rb] = await Promise.all([a, b])
    expect(ra).toEqual(rb)
    expect(ra).toMatchObject({ cached: false, remaining: 9 })
    expect(requestDownload).toHaveBeenCalledTimes(1)
    expect(fetch).toHaveBeenCalledTimes(1)
    expect(readdirSync(cacheDir)).toEqual(['77.srt'])
  })

  it('different ids are not joined', async () => {
    const { service, requestDownload, release } = gated()
    const both = Promise.all([service.download(1), service.download(2)])
    release()
    await both
    expect(requestDownload).toHaveBeenCalledTimes(2)
  })

  it('after a failure the entry is dropped and a new attempt downloads again', async () => {
    const { service, requestDownload, fetch, cacheDir } = setup({})
    fetch.mockResolvedValueOnce(new Response('', { status: 503 }))
    const failed = [service.download(5), service.download(5)]
    await expect(failed[0]).rejects.toMatchObject({ code: 'NETWORK' })
    await expect(failed[1]).rejects.toMatchObject({ code: 'NETWORK' })
    expect(requestDownload).toHaveBeenCalledTimes(1)

    await expect(service.download(5)).resolves.toMatchObject({ cached: false })
    expect(requestDownload).toHaveBeenCalledTimes(2)
    expect(readdirSync(cacheDir)).toEqual(['5.srt'])
  })

  it('a cancelled caller is detached; the shared download finishes for the others', async () => {
    const { service, requestDownload, release, cacheDir } = gated()
    const ctrl = new AbortController()
    const cancelledCall = service.download(77, ctrl.signal)
    const other = service.download(77)
    ctrl.abort()
    await expect(cancelledCall).rejects.toMatchObject({ code: 'CANCELLED' })
    release()
    await expect(other).resolves.toMatchObject({ cached: false })
    expect(requestDownload).toHaveBeenCalledTimes(1)
    expect(readdirSync(cacheDir)).toEqual(['77.srt'])
  })

  it('when the only caller cancels, the download still lands in the cache (quota not wasted)', async () => {
    const { service, requestDownload, release, cacheDir } = gated()
    const ctrl = new AbortController()
    const call = service.download(77, ctrl.signal)
    ctrl.abort()
    await expect(call).rejects.toMatchObject({ code: 'CANCELLED' })
    release()
    await vi.waitFor(() => expect(readdirSync(cacheDir)).toEqual(['77.srt']))
    await expect(service.download(77)).resolves.toMatchObject({ cached: true })
    expect(requestDownload).toHaveBeenCalledTimes(1)
  })

  it('an already aborted signal starts nothing', async () => {
    const { service, requestDownload } = setup({})
    await expect(service.download(77, AbortSignal.abort())).rejects.toMatchObject({
      code: 'CANCELLED'
    })
    expect(requestDownload).not.toHaveBeenCalled()
  })
})

describe('download redirects', () => {
  const redirect = (location: string, status = 302): Response =>
    new Response(null, { status, headers: { location } })

  function withFetch(responses: (url: string) => Response): ReturnType<typeof setup> {
    const ctx = setup({})
    ctx.fetch.mockImplementation(async (url: string) => responses(String(url)))
    return ctx
  }

  it('follows a redirect to an allowed host, checking every hop by hand', async () => {
    const { service, fetch, cacheDir } = withFetch((url) =>
      url === LINK
        ? redirect('https://dl.opensubtitles.org/file/77.srt')
        : new Response(SRT, { status: 200 })
    )
    await expect(service.download(77)).resolves.toMatchObject({ cached: false })
    expect(fetch).toHaveBeenCalledTimes(2)
    expect(fetch.mock.calls.map((c) => c[1]?.redirect)).toEqual(['manual', 'manual'])
    expect(fetch.mock.calls[1][0]).toBe('https://dl.opensubtitles.org/file/77.srt')
    expect(readdirSync(cacheDir)).toEqual(['77.srt'])
  })

  it('resolves a relative Location against the current hop', async () => {
    const { service, fetch } = withFetch((url) =>
      url === LINK ? redirect('/other/77.srt', 301) : new Response(SRT, { status: 200 })
    )
    await service.download(77)
    expect(fetch.mock.calls[1][0]).toBe('https://www.opensubtitles.com/other/77.srt')
  })

  it.each([
    ['a foreign host', 'https://evil.example/file.srt'],
    ['plain HTTP', 'http://dl.opensubtitles.org/file.srt'],
    ['a look-alike host', 'https://opensubtitles.com.evil.example/x']
  ])(
    'a redirect to %s is rejected before any request reaches it, nothing is written',
    async (_, target) => {
      const { service, fetch, cacheDir } = withFetch(() => redirect(target))
      await expect(service.download(9)).rejects.toMatchObject({ code: 'NETWORK' })
      expect(fetch).toHaveBeenCalledTimes(1)
      expect(existsSync(cacheDir)).toBe(false)
    }
  )

  it('a redirect without Location fails', async () => {
    const { service } = withFetch(() => new Response(null, { status: 307 }))
    await expect(service.download(9)).rejects.toMatchObject({ code: 'NETWORK' })
  })

  it('more than 3 redirects fail, even inside the allowed domain', async () => {
    let n = 0
    const { service, fetch, cacheDir } = withFetch(() =>
      redirect(`https://dl.opensubtitles.org/hop/${++n}`)
    )
    await expect(service.download(9)).rejects.toMatchObject({
      code: 'NETWORK',
      message: expect.stringContaining('too many redirects')
    })
    expect(fetch).toHaveBeenCalledTimes(4)
    expect(existsSync(cacheDir)).toBe(false)
  })

  it('3 redirects are still followed', async () => {
    let n = 0
    const { service, fetch } = withFetch(() =>
      n < 3 ? redirect(`https://dl.opensubtitles.org/hop/${++n}`) : new Response(SRT)
    )
    await expect(service.download(9)).resolves.toMatchObject({ cached: false })
    expect(fetch).toHaveBeenCalledTimes(4)
  })
})

describe('isAllowedDownloadLink', () => {
  it('HTTPS only, no credentials, on opensubtitles.com/.org and subdomains', () => {
    expect(isAllowedDownloadLink('https://www.opensubtitles.com/download/x')).toBe(true)
    expect(isAllowedDownloadLink('https://dl.opensubtitles.org/x')).toBe(true)
    expect(isAllowedDownloadLink('https://opensubtitles.com/x')).toBe(true)
    expect(isAllowedDownloadLink('http://www.opensubtitles.com/x')).toBe(false)
    expect(isAllowedDownloadLink('https://opensubtitles.com.evil.example/x')).toBe(false)
    expect(isAllowedDownloadLink('https://notopensubtitles.com/x')).toBe(false)
    expect(isAllowedDownloadLink('https://user:pw@www.opensubtitles.com/x')).toBe(false)
    expect(isAllowedDownloadLink('file:///etc/passwd')).toBe(false)
    expect(isAllowedDownloadLink('not a url')).toBe(false)
  })
})
