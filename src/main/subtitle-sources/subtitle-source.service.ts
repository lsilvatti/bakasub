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

export class SubtitleSourceService {
  private readonly fetch: typeof fetch

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

  async download(fileId: number, signal?: AbortSignal): Promise<SubtitleDownload> {
    const name = `${fileId}.srt`
    const path = join(this.d.cacheDir, name)
    const base = { path, name, format: 'srt' as const }
    if (await this.isValid(path)) return { ...base, cached: true }

    const { link, remaining, resetAt } = await this.d.session.run(
      (auth) => this.d.client.requestDownload(fileId, auth, signal),
      signal
    )
    if (!isAllowedDownloadLink(link))
      throw new AppError('NETWORK', 'OpenSubtitles returned a download link outside its domain')
    const bytes = await this.fetchFile(link, signal)

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

  private async fetchFile(link: string, signal?: AbortSignal): Promise<Buffer> {
    const timeout = AbortSignal.timeout(DOWNLOAD_TIMEOUT_MS)
    let res: Response
    try {
      res = await this.fetch(link, {
        signal: signal ? AbortSignal.any([signal, timeout]) : timeout
      })
    } catch (err) {
      if (signal?.aborted) throw new AppError('CANCELLED', 'Subtitle download cancelled')
      throw new AppError('NETWORK', 'OpenSubtitles download failed', { cause: err })
    }
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
      if (signal?.aborted) throw new AppError('CANCELLED', 'Subtitle download cancelled')
      throw new AppError('NETWORK', 'OpenSubtitles download interrupted', { cause: err })
    }
    if (size === 0) throw new AppError('INVALID_SUBTITLE', 'Downloaded subtitle is empty')
    return Buffer.concat(chunks)
  }
}
