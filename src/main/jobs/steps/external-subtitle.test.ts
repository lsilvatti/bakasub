import { mkdtempSync, realpathSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { beforeAll, describe, expect, it, vi } from 'vitest'
import type { Artifacts } from '../artifacts'
import type { JobSpec } from '../job-spec'
import type { StepContext } from '../step-context'
import { extractStep } from './extract.step'
import { muxStep } from './mux.step'
import { probeStep } from './probe.step'
import { translationInput } from './translation-input'

const SRT = '1\n00:00:01,000 --> 00:00:02,000\nHello there.\n'

let dir: string
let external: string
let translated: string
beforeAll(() => {
  dir = realpathSync.native(mkdtempSync(join(tmpdir(), 'bakasub-ext-')))
  external = join(dir, 'downloaded.srt')
  translated = join(dir, 'translated.srt')
  writeFileSync(external, SRT)
  writeFileSync(translated, SRT)
  writeFileSync(join(dir, 'ep01.mkv'), 'mkv')
})

const spec = (over: object = {}): JobSpec =>
  ({
    kind: 'pipeline',
    source: join(dir, 'ep01.mkv'),
    sourceSize: 3,
    outputDir: dir,
    externalSubtitle: { path: external, size: SRT.length, format: 'srt', includeInMux: false },
    preset: {
      id: 'p',
      name: 'P',
      mediaType: 'anime',
      model: 'a/b',
      glossary: [],
      temperature: 0.3,
      prompt: ''
    },
    targetLang: 'pt-BR',
    sourceLang: 'en',
    setDefault: true,
    keepOriginalSubtitles: true,
    outputMode: 'mux',
    confirmBeforeTranslate: true,
    memory: 'use',
    ...over
  }) as JobSpec

function ctx(
  s: JobSpec,
  artifacts: Artifacts = {}
): { ctx: StepContext; mkv: Record<'probe' | 'extract' | 'mux', ReturnType<typeof vi.fn>> } {
  const mkv = {
    probe: vi.fn(async () => ({ subtitles: [], attachments: 0 })),
    extract: vi.fn(async () => ({ warnings: [] })),
    mux: vi.fn(async (plan: { output: string }) => {
      writeFileSync(plan.output, 'muxed')
      return { warnings: [] }
    })
  }
  return {
    mkv,
    ctx: {
      spec: s,
      workDir: dir,
      artifacts,
      signal: new AbortController().signal,
      deps: { mkv } as never,
      progress: () => undefined,
      warn: () => undefined,
      recordUsage: () => undefined
    }
  }
}

describe('pipeline with an external subtitle', () => {
  it('probe accepts an MKV with no subtitle track', async () => {
    const c = ctx(spec())
    await probeStep(c.ctx)
    expect(c.ctx.artifacts.probe).toEqual({ tracks: [] })
  })

  it('extract, if called, asks for no track', async () => {
    const c = ctx(spec(), { probe: { tracks: [] } })
    await extractStep(c.ctx)
    expect(c.mkv.extract).toHaveBeenCalledWith(expect.any(String), [], expect.anything())
    expect(c.ctx.artifacts.extract).toEqual({ files: [] })
  })

  it('translation reads the external file directly, with no extract artifact', async () => {
    const input = await translationInput(ctx(spec()).ctx)
    expect(input.formatId).toBe('srt')
    expect(input.sourceLang).toBe('en')
    expect(input.targetLang).toBe('pt-BR')
    expect(input.doc.items.filter((i) => i.kind === 'cue')).toHaveLength(1)
  })

  it('missing external file: fails on read, not on extract', async () => {
    const gone = spec({
      externalSubtitle: { path: join(dir, 'gone.srt'), size: 1, format: 'srt', includeInMux: false }
    })
    await expect(translationInput(ctx(gone).ctx)).rejects.toMatchObject({ code: 'ENOENT' })
  })

  const translate = {
    file: { path: '', size: SRT.length, sha256: 'x', format: 'srt' as const },
    failedLines: [],
    usage: {} as never,
    cache: { hits: 0, misses: 0 }
  }

  it('mux: by default only the translation goes into the MKV', async () => {
    const c = ctx(spec(), {
      translate: { ...translate, file: { ...translate.file, path: translated } }
    })
    await muxStep(c.ctx)
    const plan = c.mkv.mux.mock.calls[0][0]
    expect(plan.subtitles).toEqual([
      {
        path: translated,
        language: 'pt-BR',
        name: expect.stringContaining('(BakaSub)'),
        default: true
      }
    ])
    expect(plan.sourceSubtitleTrackIds).toEqual([])
  })

  it('mux with includeInMux: the original goes after it, in the source language, never as default', async () => {
    const s = spec({
      externalSubtitle: { path: external, size: SRT.length, format: 'srt', includeInMux: true }
    })
    const c = ctx(s, { translate: { ...translate, file: { ...translate.file, path: translated } } })
    await muxStep(c.ctx)
    const { subtitles } = c.mkv.mux.mock.calls[0][0]
    expect(subtitles).toHaveLength(2)
    expect(subtitles[0]).toMatchObject({ path: translated, default: true })
    expect(subtitles[1]).toMatchObject({ path: external, language: 'en', default: false })
  })

  it('mux with includeInMux and an unknown source language: an `und` track', async () => {
    const s = spec({
      sourceLang: undefined,
      externalSubtitle: { path: external, size: SRT.length, format: 'srt', includeInMux: true }
    })
    const c = ctx(s, { translate: { ...translate, file: { ...translate.file, path: translated } } })
    await muxStep(c.ctx)
    expect(c.mkv.mux.mock.calls[0][0].subtitles[1]).toMatchObject({
      language: 'und',
      name: 'Original'
    })
  })
})
