import { mkdir, rename, stat, unlink, writeFile } from 'node:fs/promises'
import { join } from 'node:path'
import type {
  SubtitleDownload,
  SubtitleSearchRequest,
  SubtitleSearchResponse,
  SubtitleSearchResult
} from '../../shared/ipc-contract'
import { AppError } from '../errors'
import { readSubtitle } from '../jobs/steps/subtitle-io'
import { canonicalTag } from '../media/languages'
import { MAX_SUBTITLE_BYTES, MEDIA_EXTENSIONS, validateInputFile } from '../media/path-guard'
import { parseReleaseName } from '../media-context/release-name'
import type { OpenSubtitlesClient, OpenSubtitlesFile } from './opensubtitles-client'
import type { OpenSubtitlesSession } from './opensubtitles-session'

const DOWNLOAD_TIMEOUT_MS = 30_000
const MAX_REDIRECTS = 3
const REDIRECT_STATUS = new Set([301, 302, 303, 307, 308])
const MAX_RESULTS = 60
const LINK_HOST = /(^|\.)opensubtitles\.(com|org)$/

export function isAllowedDownloadLink(raw: string): boolean {
  try {
    const url = new URL(raw)
    return (
      url.protocol === 'https:' &&
      url.username === '' &&
      url.password === '' &&
      LINK_HOST.test(url.hostname.toLowerCase())
    )
  } catch {
    return false
  }
}

function resolveLocation(location: string | null, base: string): string | null {
  if (!location) return null
  try {
    return new URL(location, base).toString()
  } catch {
    return null
  }
}

function rank(a: SubtitleSearchResult, b: SubtitleSearchResult): number {
  return (
    Number(b.hashMatch) - Number(a.hashMatch) ||
    Number(a.machineTranslated) - Number(b.machineTranslated) ||
    Number(b.trusted) - Number(a.trusted) ||
    b.downloads - a.downloads ||
    a.fileId - b.fileId
  )
}

function toResult(f: OpenSubtitlesFile, hashMatch: boolean): SubtitleSearchResult {
  return {
    fileId: f.fileId,
    release: f.release,
    language: canonicalTag(f.language) ?? f.language,
    downloads: f.downloads,
    hashMatch: hashMatch || f.hashMatch,
    hearingImpaired: f.hearingImpaired,
    machineTranslated: f.machineTranslated,
    trusted: f.trusted,
    ...(f.fps ? { fps: f.fps } : {}),
    ...(f.uploader ? { uploader: f.uploader } : {})
  }
}

const cancelled = (): AppError => new AppError('CANCELLED', 'Subtitle download cancelled')

/** Settles with `shared`, or rejects with `CANCELLED` as soon as `signal` aborts, whichever first. */
function detach<T>(shared: Promise<T>, signal: AbortSignal): Promise<T> {
  return new Promise<T>((resolve, reject) => {
    const onAbort = (): void => reject(cancelled())
    signal.addEventListener('abort', onAbort, { once: true })
    shared.then(resolve, reject).finally(() => signal.removeEventListener('abort', onAbort))
  })
}

export class SubtitleSourceService {
  private readonly fetch: typeof fetch
  private readonly inflight = new Map<number, Promise<SubtitleDownload>>()

  constructor(
    private readonly d: {
      client: Pick<OpenSubtitlesClient, 'search' | 'requestDownload'>
      session: Pick<OpenSubtitlesSession, 'run'>
      cacheDir: string
      hash: (path: string) => Promise<string | null>
      fetch?: typeof fetch
    }
  ) {
    this.fetch = d.fetch ?? globalThis.fetch
  }

  async search(req: SubtitleSearchRequest, signal?: AbortSignal): Promise<SubtitleSearchResponse> {
    const video = await validateInputFile(req.videoPath, MEDIA_EXTENSIONS)
    const typed = req.query?.trim()
    const release = parseReleaseName(video.path)
    const query = typed || release.title
    const season = typed ? req.season : (req.season ?? release.season)
    const episode = typed ? req.episode : (req.episode ?? release.episode)
    const languages = [...new Set(req.languages.map((l) => l.toLowerCase()))]
    const moviehash = await this.d.hash(video.path).catch(() => null)

    const [byHash, byQuery] = await Promise.all([
      moviehash
        ? this.d.session.run((auth) => this.d.client.search({ moviehash, languages }, auth, signal))
        : Promise.resolve([]),
      this.d.session.run((auth) =>
        this.d.client.search({ query, languages, season, episode }, auth, signal)
      )
    ])

    const merged = new Map<number, SubtitleSearchResult>()
    for (const f of byHash) merged.set(f.fileId, toResult(f, true))
    for (const f of byQuery) if (!merged.has(f.fileId)) merged.set(f.fileId, toResult(f, false))
    return {
      query,
      ...(season !== undefined ? { season } : {}),
      ...(episode !== undefined ? { episode } : {}),
      results: [...merged.values()].sort(rank).slice(0, MAX_RESULTS)
    }
  }

  /**
   * Concurrent calls for the same id share one download: one quota spend, one `.partial` writer.
   * A caller's signal only detaches that caller (it rejects with `CANCELLED`); the shared download
   * keeps running, bounded by its own timeout, so the others still get the file and the quota
   * already spent ends up in the cache. The entry is dropped once it settles, success or error,
   * so a later call after a failure starts a fresh attempt.
   */
  download(fileId: number, signal?: AbortSignal): Promise<SubtitleDownload> {
    if (signal?.aborted) return Promise.reject(cancelled())
    let shared = this.inflight.get(fileId)
    if (!shared) {
      const started = this.downloadToCache(fileId).finally(() => {
        if (this.inflight.get(fileId) === started) this.inflight.delete(fileId)
      })
      this.inflight.set(fileId, started)
      shared = started
    }
    return signal ? detach(shared, signal) : shared
  }

  private async downloadToCache(fileId: number): Promise<SubtitleDownload> {
    const name = `${fileId}.srt`
    const path = join(this.d.cacheDir, name)
    const base = { path, name, format: 'srt' as const }
    if (await this.isValid(path)) return { ...base, cached: true }

    const { link, remaining, resetAt } = await this.d.session.run((auth) =>
      this.d.client.requestDownload(fileId, auth)
    )
    if (!isAllowedDownloadLink(link))
      throw new AppError('NETWORK', 'OpenSubtitles returned a download link outside its domain')
    const bytes = await this.fetchFile(link)

    await mkdir(this.d.cacheDir, { recursive: true })
    const partial = `${path}.partial`
    try {
      await writeFile(partial, bytes)
      await readSubtitle(partial, 'srt')
      await rename(partial, path)
    } catch (err) {
      await unlink(partial).catch(() => undefined)
      throw err
    }
    return {
      ...base,
      cached: false,
      ...(remaining !== undefined ? { remaining } : {}),
      ...(resetAt ? { resetAt } : {})
    }
  }

  private async isValid(path: string): Promise<boolean> {
    try {
      if (!(await stat(path)).isFile()) return false
      await readSubtitle(path, 'srt')
      return true
    } catch {
      return false
    }
  }

  /**
   * Follows redirects by hand so every hop is checked against the link allowlist before any
   * request reaches it; `redirect: 'follow'` would contact the foreign host first.
   */
  private async followRedirects(link: string, signal: AbortSignal): Promise<Response> {
    let url = link
    for (let hop = 0; ; hop++) {
      let res: Response
      try {
        res = await this.fetch(url, { redirect: 'manual', signal })
      } catch (err) {
        throw new AppError('NETWORK', 'OpenSubtitles download failed', { cause: err })
      }
      // A browser-style fetch hides the Location of a manual redirect: nothing safe to follow.
      if (res.type === 'opaqueredirect')
        throw new AppError('NETWORK', 'OpenSubtitles download redirect could not be verified')
      if (!REDIRECT_STATUS.has(res.status)) return res

      await res.body?.cancel().catch(() => undefined)
      if (hop >= MAX_REDIRECTS)
        throw new AppError('NETWORK', 'OpenSubtitles download: too many redirects')
      const next = resolveLocation(res.headers.get('location'), url)
      if (!next || !isAllowedDownloadLink(next))
        throw new AppError('NETWORK', 'OpenSubtitles download redirected outside its domain')
      url = next
    }
  }

  /** No caller signal on purpose: the fetch is shared (see `download`), only the timeout ends it. */
  private async fetchFile(link: string): Promise<Buffer> {
    const res = await this.followRedirects(link, AbortSignal.timeout(DOWNLOAD_TIMEOUT_MS))
    if (!res.ok || !res.body) throw new AppError('NETWORK', `OpenSubtitles download ${res.status}`)
    const chunks: Uint8Array[] = []
    let size = 0
    try {
      for await (const chunk of res.body as unknown as AsyncIterable<Uint8Array>) {
        size += chunk.byteLength
        if (size > MAX_SUBTITLE_BYTES) {
          await res.body.cancel().catch(() => undefined)
          throw new AppError('INVALID_SUBTITLE', 'Downloaded subtitle is too large')
        }
        chunks.push(chunk)
      }
    } catch (err) {
      if (err instanceof AppError) throw err
      throw new AppError('NETWORK', 'OpenSubtitles download interrupted', { cause: err })
    }
    if (size === 0) throw new AppError('INVALID_SUBTITLE', 'Downloaded subtitle is empty')
    return Buffer.concat(chunks)
  }
}
