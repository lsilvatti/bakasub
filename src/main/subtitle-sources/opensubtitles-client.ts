import { randomUUID } from 'node:crypto'
import { z } from 'zod'
import { AppError } from '../errors'

const DEFAULT_HOST = 'api.opensubtitles.com'
const TIMEOUT_MS = 15_000
const API_HOST = /^[a-z0-9-]+\.opensubtitles\.com$/

export interface OpenSubtitlesAuth {
  readonly token: string
  readonly host?: string
}

export interface OpenSubtitlesSearchParams {
  readonly query?: string
  readonly moviehash?: string
  readonly languages?: readonly string[]
  readonly season?: number
  readonly episode?: number
}

export interface OpenSubtitlesFile {
  readonly fileId: number
  readonly release: string
  readonly language: string
  readonly downloads: number
  readonly hashMatch: boolean
  readonly hearingImpaired: boolean
  readonly machineTranslated: boolean
  readonly trusted: boolean
  readonly fps?: number
  readonly uploader?: string
}

export interface OpenSubtitlesDownloadLink {
  readonly link: string
  readonly remaining?: number
  readonly resetAt?: string
}

const bool = z.boolean().nullish()
const SearchBody = z.looseObject({
  data: z
    .array(
      z.looseObject({
        attributes: z.looseObject({
          language: z.string().nullish(),
          download_count: z.number().nullish(),
          hearing_impaired: bool,
          fps: z.number().nullish(),
          from_trusted: bool,
          ai_translated: bool,
          machine_translated: bool,
          moviehash_match: bool,
          release: z.string().nullish(),
          uploader: z.looseObject({ name: z.string().nullish() }).nullish(),
          files: z
            .array(z.looseObject({ file_id: z.number(), file_name: z.string().nullish() }))
            .nullish()
        })
      })
    )
    .default([])
})
const LoginBody = z.looseObject({ token: z.string().min(1), base_url: z.string().nullish() })
const DownloadBody = z.looseObject({
  link: z.string().min(1),
  remaining: z.number().nullish(),
  reset_time_utc: z.string().nullish()
})

export class OpenSubtitlesClient {
  private readonly fetch: typeof fetch

  constructor(
    private readonly o: {
      apiKey: () => Promise<string | null>
      userAgent: string
      fetch?: typeof fetch
    }
  ) {
    this.fetch = o.fetch ?? globalThis.fetch
  }

  async validate(apiKey: string, signal?: AbortSignal): Promise<void> {
    // The CDN caches `/subtitles` for 24 h regardless of the key: only a new query reaches the origin.
    const query = `bakasub-key-check-${randomUUID()}`
    await this.request('GET', '/subtitles', { params: { query }, apiKey, signal })
  }

  async login(
    username: string,
    password: string,
    signal?: AbortSignal
  ): Promise<OpenSubtitlesAuth> {
    const body = LoginBody.safeParse(
      await this.request('POST', '/login', { body: { username, password }, signal })
    )
    if (!body.success) throw new AppError('NETWORK', 'OpenSubtitles: unexpected login response')
    const host = body.data.base_url?.toLowerCase()
    return { token: body.data.token, ...(host && API_HOST.test(host) ? { host } : {}) }
  }

  async search(
    p: OpenSubtitlesSearchParams,
    auth?: OpenSubtitlesAuth,
    signal?: AbortSignal
  ): Promise<OpenSubtitlesFile[]> {
    const params: Record<string, string> = {}
    if (p.query) params.query = p.query.toLowerCase()
    if (p.moviehash) params.moviehash = p.moviehash.toLowerCase()
    if (p.languages?.length) params.languages = [...p.languages].sort().join(',').toLowerCase()
    if (p.season !== undefined) params.season_number = String(p.season)
    if (p.episode !== undefined) params.episode_number = String(p.episode)
    const body = SearchBody.safeParse(
      await this.request('GET', '/subtitles', { params, auth, signal })
    )
    if (!body.success) throw new AppError('NETWORK', 'OpenSubtitles: unexpected search response')
    return body.data.data.flatMap(({ attributes: a }) =>
      (a.files ?? []).map((f): OpenSubtitlesFile => ({
        fileId: f.file_id,
        release: a.release?.trim() || f.file_name?.trim() || String(f.file_id),
        language: a.language ?? '',
        downloads: a.download_count ?? 0,
        hashMatch: a.moviehash_match === true,
        hearingImpaired: a.hearing_impaired === true,
        machineTranslated: a.machine_translated === true || a.ai_translated === true,
        trusted: a.from_trusted === true,
        ...(a.fps ? { fps: a.fps } : {}),
        ...(a.uploader?.name ? { uploader: a.uploader.name } : {})
      }))
    )
  }

  async requestDownload(
    fileId: number,
    auth?: OpenSubtitlesAuth,
    signal?: AbortSignal
  ): Promise<OpenSubtitlesDownloadLink> {
    const body = DownloadBody.safeParse(
      await this.request('POST', '/download', {
        body: { file_id: fileId, sub_format: 'srt' },
        auth,
        signal
      })
    )
    if (!body.success) throw new AppError('NETWORK', 'OpenSubtitles: unexpected download response')
    const { link, remaining, reset_time_utc: resetAt } = body.data
    return {
      link,
      ...(typeof remaining === 'number' ? { remaining } : {}),
      ...(resetAt ? { resetAt } : {})
    }
  }

  private async request(
    method: 'GET' | 'POST',
    path: string,
    o: {
      params?: Record<string, string>
      body?: unknown
      apiKey?: string
      auth?: OpenSubtitlesAuth
      signal?: AbortSignal
    }
  ): Promise<unknown> {
    const apiKey = o.apiKey ?? (await this.o.apiKey())
    if (!apiKey)
      throw new AppError('OPENSUBTITLES_KEY_MISSING', 'No OpenSubtitles API key configured')
    const url = new URL(`https://${o.auth?.host ?? DEFAULT_HOST}/api/v1${path}`)
    for (const k of Object.keys(o.params ?? {}).sort()) url.searchParams.set(k, o.params![k])
    const headers: Record<string, string> = {
      Accept: 'application/json',
      'Api-Key': apiKey,
      'User-Agent': this.o.userAgent
    }
    if (o.auth) headers.Authorization = `Bearer ${o.auth.token}`
    if (o.body !== undefined) headers['Content-Type'] = 'application/json'

    const timeout = AbortSignal.timeout(TIMEOUT_MS)
    let res: Response
    try {
      res = await this.fetch(url, {
        method,
        headers,
        ...(o.body !== undefined ? { body: JSON.stringify(o.body) } : {}),
        signal: o.signal ? AbortSignal.any([o.signal, timeout]) : timeout
      })
    } catch (err) {
      if (o.signal?.aborted) throw new AppError('CANCELLED', 'OpenSubtitles request cancelled')
      throw new AppError('NETWORK', `OpenSubtitles unreachable: ${path}`, { cause: err })
    }
    if (res.status === 401 || res.status === 403 || (path === '/login' && res.status === 400))
      throw new AppError('OPENSUBTITLES_AUTH', `OpenSubtitles rejected the credentials (${path})`)
    if (res.status === 406)
      throw new AppError('OPENSUBTITLES_QUOTA', 'OpenSubtitles download quota exceeded')
    if (!res.ok) throw new AppError('NETWORK', `OpenSubtitles ${res.status} on ${path}`)
    try {
      return await res.json()
    } catch (err) {
      throw new AppError('NETWORK', `OpenSubtitles: invalid JSON on ${path}`, { cause: err })
    }
  }
}
