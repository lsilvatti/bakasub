// @vitest-environment jsdom
import { QueryClient } from '@tanstack/react-query'
import { createMemoryHistory } from '@tanstack/react-router'
import { cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type {
  CoverageReport,
  Job,
  JobCreateRequest,
  MediaIdentification,
  PreflightEstimate,
  Preset,
  SubtitleSearchRequest,
  SubtitleSearchResult,
  SubtitleTrackInfo
} from '@shared/ipc-contract'
import { App } from '@renderer/app/app'
import { createAppRouter } from '@renderer/app/router'
import { installFakeApi, ok, type FakeApi } from '@renderer/test/fake-api'
import { onboardedSettings } from '@renderer/test/fixtures'

const preset: Preset = {
  id: 'builtin:series',
  builtin: true,
  name: 'Series',
  mediaType: 'series',
  temperature: 0.45,
  prompt: 'Natural dialogue.',
  glossary: [{ source: 'Frieren' }],
  copiedFrom: null,
  createdAt: null
}

const estimate: PreflightEstimate = {
  model: {
    id: 'a/b',
    pricing: { prompt: 1e-6, completion: 2e-6 },
    catalogFetchedAt: '2026-11-19T00:00:00.000Z'
  },
  approximate: false,
  lines: { total: 300, cached: 20, toTranslate: 280, copies: 0, chunks: 7 },
  tokens: { input: { low: 20000, high: 23000 }, output: { low: 8000, high: 9200 } },
  cost: { low: 0.036, high: 0.0414 },
  decision: 'needs-confirmation'
}

const coverage: CoverageReport = {
  totals: { cues: 301, blocks: 301, translatable: 300, skipped: 1, untranslated: 1 },
  byReason: { SKIPPED_NO_LETTERS: 1, UNTRANSLATED_KARAOKE: 1 },
  entries: [
    {
      ref: { kind: 'cue', cueId: 9 },
      status: 'untranslated',
      severity: 'warn',
      reason: 'UNTRANSLATED_KARAOKE',
      params: {},
      location: { line: 42, endLine: 42, startMs: 61_000 },
      excerpt: '{\\k20}La la la',
      outcome: 'kept-original',
      actions: []
    }
  ],
  truncated: false
}

const job = (over: Partial<Job>): Job => ({
  id: 'j1',
  kind: 'translate',
  batchId: null,
  status: 'awaiting-confirmation',
  step: 'preflight',
  source: 'ep01.en.srt',
  progress: { step: 'preflight' },
  result: null,
  error: null,
  attempt: 1,
  createdAt: '2026-11-19T10:00:00.000Z',
  updatedAt: '2026-11-19T10:00:00.000Z',
  startedAt: null,
  finishedAt: null,
  ...over
})

/** TMDB knows the show as Japanese animation (→ the Anime preset is suggested). */
const animeMedia = {
  identify: () =>
    ok({
      release: { title: 'Show' },
      status: 'matched' as const,
      match: {
        provider: 'tmdb' as const,
        id: 7,
        kind: 'series' as const,
        title: 'The Show',
        genres: ['Animation'],
        synopsis: 'Synopsis.'
      },
      episode: null
    } satisfies MediaIdentification),
  lookup: () =>
    ok({
      details: {
        provider: 'tmdb' as const,
        id: 7,
        kind: 'series' as const,
        title: 'The Show',
        genres: ['Animation'],
        genreIds: [16],
        originalLanguage: 'ja',
        synopsis: 'Synopsis.'
      },
      episode: null
    }),
  search: () => ok([])
}

function setup(
  path: string,
  tracks?: SubtitleTrackInfo[],
  media?: typeof animeMedia,
  opensubtitles?: { result?: Partial<SubtitleSearchResult>; videoFps?: number }
): {
  fake: FakeApi
  create: ReturnType<typeof vi.fn>
  confirm: ReturnType<typeof vi.fn>
  setMemory: ReturnType<typeof vi.fn>
  retry: ReturnType<typeof vi.fn>
  subtitleSearch: ReturnType<typeof vi.fn>
  subtitleDownload: ReturnType<typeof vi.fn>
  update: (over: Partial<Job>) => void
} {
  let current = job({})
  let preflightEstimate = estimate
  const create = vi.fn(async (req: JobCreateRequest) => {
    const memory = 'memory' in req && req.memory ? req.memory : 'use'
    current = job({ kind: req.kind, memory })
    preflightEstimate = { ...estimate, memory }
    return { ok: true as const, data: current }
  })
  const setMemory = vi.fn(async (_id: string, memory: 'use' | 'skip') => {
    current = { ...current, memory, status: 'queued', result: null }
    preflightEstimate = { ...estimate, memory }
    return { ok: true as const, data: current }
  })
  const retry = vi.fn(async (_id: string, memory?: 'use' | 'skip') => {
    current = { ...current, memory: memory ?? current.memory, status: 'queued' }
    return { ok: true as const, data: current }
  })
  const confirm = vi.fn(async () => {
    current = { ...current, status: 'queued', progress: { step: 'preflight' } }
    return { ok: true as const, data: current }
  })
  const subtitleSearch = vi.fn((req: SubtitleSearchRequest) =>
    ok({
      query: req.query ?? 'Show',
      season: 1,
      episode: 1,
      results: [
        {
          fileId: 77,
          release: 'Show.S01E01.WEB',
          language: 'en',
          downloads: 120,
          hashMatch: true,
          hearingImpaired: false,
          machineTranslated: false,
          trusted: false,
          ...opensubtitles?.result
        }
      ]
    })
  )
  const subtitleDownload = vi.fn((fileId: number) =>
    ok({
      path: `/cache/opensubtitles/${fileId}.srt`,
      name: `${fileId}.srt`,
      format: 'srt' as const,
      cached: false,
      remaining: 4
    })
  )
  const fake = installFakeApi({
    jobs: {
      list: () => ok([]),
      create,
      confirm,
      setMemory,
      retry,
      get: () => ok(current),
      preflight: () => ok({ estimate: preflightEstimate, coverage })
    },
    sidecars: {
      check: () =>
        ok({
          checkedAt: '',
          manifestUpdatedAt: '',
          offline: false,
          installDir: '',
          tools: [],
          ready: true
        })
    },
    settings: {
      apiKeyStatus: () => ok({ present: true, protection: 'os' as const }),
      openSubtitlesStatus: () => ok({ keyPresent: true, username: null }),
      get: () => ok({ ...onboardedSettings, 'translation.targetLangs': ['pt-BR', 'es'] })
    },
    presets: {
      list: () =>
        ok([preset, { ...preset, id: 'builtin:anime', name: 'Anime', mediaType: 'anime' as const }])
    },
    models: {
      list: () =>
        ok({
          fetchedAt: '',
          offline: false,
          models: [
            {
              id: 'a/b',
              name: 'Model AB',
              contextLength: null,
              maxCompletionTokens: null,
              isModerated: false,
              pricing: null,
              structuredOutputs: true,
              responseFormat: true
            }
          ]
        })
    },
    dialog: { openFiles: () => ok({ paths: [path] }) },
    subtitles: {
      search: subtitleSearch,
      cancelSearch: () => ok({ cancelled: false }),
      download: subtitleDownload
    },
    media: {
      ...media,
      probe: () =>
        ok({
          attachments: 2,
          ...(opensubtitles?.videoFps !== undefined ? { videoFps: opensubtitles.videoFps } : {}),
          subtitles: tracks ?? [
            {
              id: 1,
              codecId: 'S_HDMV/PGS',
              format: null,
              unsupported: 'image-based' as const,
              language: 'en',
              default: true,
              forced: false
            },
            {
              id: 2,
              codecId: 'S_TEXT/ASS',
              format: 'ass' as const,
              language: 'en',
              default: false,
              forced: false,
              eventCount: 310
            }
          ]
        })
    }
  })
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createAppRouter({
    queryClient,
    history: createMemoryHistory({ initialEntries: ['/translate'] })
  })
  render(<App queryClient={queryClient} router={router} />)
  return {
    fake,
    create,
    confirm,
    setMemory,
    retry,
    subtitleSearch,
    subtitleDownload,
    update: (over) => {
      current = { ...current, ...over }
    }
  }
}

afterEach(cleanup)

/** Translate opens straight on the single-file wizard. */
const openSingle = async (): Promise<void> => {
  await screen.findByRole('heading', { name: 'Translate a file' })
}

const next = async (): Promise<void> =>
  userEvent.click(await screen.findByRole('button', { name: 'Continue' }))

const summary = (): HTMLElement => screen.getByRole('complementary', { name: 'Summary' })

/** The summary's row labels, in order. */
const summaryRows = (): string[] =>
  within(summary())
    .queryAllByRole('button')
    .map((b) => b.firstElementChild?.textContent ?? '')

describe('Translate a file', () => {
  it('subtitle: no track step, defaults prefilled, preflight review, confirm, live progress', async () => {
    const { fake, create, confirm } = setup('/subs/ep01.en.srt')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    expect((await screen.findAllByText('ep01.en.srt')).length).toBeGreaterThan(0)
    const steps = screen.getByRole('list', { name: 'Steps' })
    expect(within(steps).queryByText('Track')).toBeNull()
    await next()

    expect(await screen.findByRole('heading', { name: 'How to translate' })).toBeTruthy()
    expect(await screen.findByText('Natural dialogue.')).toBeTruthy()
    expect(screen.getByText('1 glossary entry')).toBeTruthy()
    // In the selects and in the summary.
    expect(screen.getAllByText('Model AB').length).toBe(2)
    expect(screen.getAllByText('Brazilian Portuguese').length).toBe(2)
    await userEvent.type(screen.getByLabelText('Notes for the translator (optional)'), 'Elf mage')
    await next()

    expect(await screen.findByRole('heading', { name: 'Where to save' })).toBeTruthy()
    expect(screen.queryByText('Keep the original subtitles in the video')).toBeNull()
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith({
      kind: 'translate',
      source: '/subs/ep01.en.srt',
      notes: 'Elf mage',
      confirmBeforeTranslate: true
    })

    expect(await screen.findByText('280 to translate')).toBeTruthy()
    expect(screen.getByText('US$ 0.036 – 0.041')).toBeTruthy()
    expect(screen.getByText('1 line will stay in the original language')).toBeTruthy()
    expect(screen.getByText(/Karaoke line/)).toBeTruthy()

    await userEvent.click(screen.getByRole('button', { name: 'Translate' }))
    expect(confirm).toHaveBeenCalledWith('j1')
    fake.emit('jobs:progress', {
      id: 'j1',
      status: 'running',
      progress: {
        step: 'translate',
        translation: {
          done: 120,
          total: 280,
          cached: 20,
          chunk: { completed: 2, count: 7, inFlight: 1 }
        },
        cost: 0.012,
        waiting: { reason: 'cooldown', untilMs: Date.now() + 9000 }
      }
    })
    expect(await screen.findByText(/120 of 280 lines · 2 of 7 parts/)).toBeTruthy()
    expect(screen.getByText(/Paused after a rate limit · \d s/)).toBeTruthy()
  })

  it('ignore saved translations: wizard, preflight review (re-estimate) and retry', async () => {
    const { fake, create, setMemory, retry, update } = setup('/subs/ep01.en.srt')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await next()
    const option = await screen.findByRole('checkbox', { name: 'Ignore saved translations' })
    expect((option as HTMLButtonElement).getAttribute('aria-checked')).toBe('false')
    expect(screen.getByText(/reused at no cost/)).toBeTruthy()
    await userEvent.click(option)
    expect(screen.getByText(/new translations replace the saved ones/)).toBeTruthy()
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith(expect.objectContaining({ memory: 'skip' }))

    // Review: the estimate says the cache is ignored; switching back re-estimates.
    expect(await screen.findByText(/cache ignored/)).toBeTruthy()
    expect(screen.getByText('Saved translations ignored')).toBeTruthy()
    const review = screen.getByRole('checkbox', { name: 'Ignore saved translations' })
    expect(review.getAttribute('aria-checked')).toBe('true')
    await userEvent.click(review)
    expect(setMemory).toHaveBeenCalledWith('j1', 'use')

    // Failed: the retry starts from the job's choice and can switch it.
    update({ status: 'failed', memory: 'skip', error: { code: 'LLM_UNAVAILABLE', message: 'x' } })
    fake.emit('jobs:progress', { id: 'j1', status: 'failed', progress: { step: 'translate' } })
    await screen.findByRole('button', { name: 'Retry' })
    const again = screen.getByRole('checkbox', { name: 'Ignore saved translations' })
    expect(again.getAttribute('aria-checked')).toBe('true')
    await userEvent.click(again)
    await userEvent.click(screen.getByRole('button', { name: 'Retry' }))
    expect(retry).toHaveBeenCalledWith('j1', 'use')
  })

  it('video: subtitle-only output hides the mux options and previews the file name', async () => {
    const { create } = setup('/v/ep01.mkv')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await screen.findByRole('radiogroup')
    await next()
    await next()
    // Default: the MKV, with its mux options.
    const preview = await screen.findByRole('region', { name: 'Files that will be written' })
    expect(within(preview).getByText('ep01.pt-BR.mkv')).toBeTruthy()
    expect(screen.getByText('Keep the original subtitles in the video')).toBeTruthy()
    await userEvent.click(screen.getByRole('radio', { name: /Subtitle only/ }))
    expect(screen.queryByText('Keep the original subtitles in the video')).toBeNull()
    expect(screen.queryByText('Make the translated subtitle the default track')).toBeNull()
    expect(within(preview).queryByText('ep01.pt-BR.mkv')).toBeNull()
    expect(within(preview).getByText(/ep01\.pt-BR\.(ass|srt)/)).toBeTruthy()
    expect(screen.getByText(/same name as the video/)).toBeTruthy()
    await userEvent.click(
      screen.getByRole('checkbox', { name: /Replace files that already exist/ })
    )
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: 'pipeline',
        options: { outputMode: 'subtitle', overwrite: true }
      })
    )
  })

  it('video: track step preselects a translatable track; image tracks are disabled', async () => {
    const { create } = setup('/v/ep01.mkv')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    const group = await screen.findByRole('radiogroup')
    const radios = await within(group).findAllByRole('radio')
    expect((radios[0] as HTMLButtonElement).disabled).toBe(true)
    expect(within(radios[0]).getByText('Image subtitles are not supported')).toBeTruthy()
    expect(radios[1].getAttribute('aria-checked')).toBe('true')
    expect(screen.getByText(/2 fonts attached/)).toBeTruthy()
    await next()
    await next()
    await userEvent.click(await screen.findByText('Make the translated subtitle the default track'))
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith(
      expect.objectContaining({
        kind: 'pipeline',
        trackId: 2,
        options: { outputMode: 'mux', setDefault: false, keepOriginalSubtitles: true }
      })
    )
  })

  it('a video that already has a subtitle in the target language asks before translating', async () => {
    const { create } = setup('/v/ep01.mkv', [
      {
        id: 2,
        codecId: 'S_TEXT/ASS',
        format: 'ass',
        language: 'en',
        default: false,
        forced: false
      },
      {
        id: 5,
        codecId: 'S_TEXT/ASS',
        format: 'ass',
        language: 'por',
        name: 'Fansub',
        default: false,
        forced: false
      }
    ])
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await next()
    expect(await screen.findByText('Already has a Brazilian Portuguese subtitle')).toBeTruthy()
    await next()
    await userEvent.click(await screen.findByRole('button', { name: 'Estimate cost' }))
    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByText(/Brazilian Portuguese subtitle track \(Fansub\)/)).toBeTruthy()
    expect(create).not.toHaveBeenCalled()
    await userEvent.click(within(dialog).getByRole('button', { name: 'Translate anyway' }))
    expect(create).toHaveBeenCalled()
  })

  it('subtitle: the summary grows with the reached steps and never shows a placeholder', async () => {
    setup('/subs/ep01.en.srt')
    await openSingle()
    expect(within(summary()).getByText('Choose a file to get started.')).toBeTruthy()
    expect(summaryRows()).toEqual([])

    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await screen.findAllByText('ep01.en.srt')
    expect(summaryRows()).toEqual(['File'])
    await next()

    await screen.findByRole('heading', { name: 'How to translate' })
    // No TMDB work: no "Media" row.
    expect(summaryRows()).toEqual(['File', 'Preset', 'Model', 'Into'])
    await next()

    await screen.findByRole('heading', { name: 'Where to save' })
    expect(summaryRows()).toEqual(['File', 'Preset', 'Model', 'Into', 'Output'])
    expect(within(summary()).getByText('ep01.pt-BR.srt')).toBeTruthy()

    // Back to the first step: the reached rows stay, and still link to their step.
    await userEvent.click(screen.getByRole('button', { name: 'Back' }))
    await userEvent.click(screen.getByRole('button', { name: 'Back' }))
    await screen.findByRole('heading', { name: 'Choose what to translate' })
    expect(summaryRows()).toEqual(['File', 'Preset', 'Model', 'Into', 'Output'])
    expect(within(summary()).queryByText('—')).toBeNull()
    for (const b of within(summary()).getAllByRole('button'))
      expect(b.getAttribute('aria-label')).not.toContain('—')
    await userEvent.click(within(summary()).getByRole('button', { name: /^Output:/ }))
    expect(await screen.findByRole('heading', { name: 'Where to save' })).toBeTruthy()
  })

  it('video: the track row comes with the track step, the media row with a TMDB work', async () => {
    setup('/v/Show - 01.mkv', undefined, animeMedia)
    await openSingle()
    expect(summaryRows()).toEqual([])
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await screen.findAllByText('Show - 01.mkv')
    expect(summaryRows()).toEqual(['File'])
    await next()

    await screen.findByRole('radiogroup')
    expect(await within(summary()).findByText('English · ASS')).toBeTruthy()
    expect(summaryRows()).toEqual(['File', 'Track'])
    await next()

    expect(await within(summary()).findByText('The Show')).toBeTruthy()
    expect(summaryRows()).toEqual(['File', 'Track', 'Preset', 'Model', 'Into', 'Media'])
    await next()

    await screen.findByRole('heading', { name: 'Where to save' })
    expect(summaryRows()).toEqual(['File', 'Track', 'Preset', 'Model', 'Into', 'Media', 'Output'])
  })

  it('reached steps can be revisited from the step bar without Continue', async () => {
    setup('/subs/ep01.en.srt')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await next()
    expect(await screen.findByRole('heading', { name: 'Where to save' })).toBeTruthy()
    const steps = screen.getByRole('list', { name: 'Steps' })
    await userEvent.click(within(steps).getByRole('button', { name: /File/ }))
    expect(await screen.findByRole('heading', { name: 'Choose what to translate' })).toBeTruthy()
    await userEvent.click(within(steps).getByRole('button', { name: /Output/ }))
    expect(await screen.findByRole('heading', { name: 'Where to save' })).toBeTruthy()
  })

  it('suggests the preset from TMDB; picking another and back follows the suggestion again', async () => {
    const { create } = setup('/subs/Show - 01.en.srt', undefined, animeMedia)
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    expect(await screen.findByText('Suggested: Japanese animation on TMDB')).toBeTruthy()
    const presetSelect = screen.getAllByRole('combobox')[0]
    expect(within(presetSelect).getByText('Anime')).toBeTruthy()

    await userEvent.click(presetSelect)
    await userEvent.click(await screen.findByRole('option', { name: 'Series' }))
    await userEvent.click(await screen.findByRole('button', { name: /Back to the suggestion/ }))
    expect(await screen.findByText('Suggested: Japanese animation on TMDB')).toBeTruthy()

    await next()
    await userEvent.click(await screen.findByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith(expect.objectContaining({ presetId: 'builtin:anime' }))
  })

  it('keeps the wizard where it was after visiting another page', async () => {
    setup('/subs/ep01.en.srt')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await userEvent.type(screen.getByLabelText('Notes for the translator (optional)'), 'Elf mage')

    const nav = screen.getByRole('navigation')
    await userEvent.click(within(nav).getByRole('link', { name: 'Extract' }))
    expect(screen.queryByRole('heading', { name: 'How to translate' })).toBeNull()
    await userEvent.click(within(nav).getByRole('link', { name: 'Translate' }))

    expect(await screen.findByRole('heading', { name: 'How to translate' })).toBeTruthy()
    expect(
      (screen.getByLabelText('Notes for the translator (optional)') as HTMLTextAreaElement).value
    ).toBe('Elf mage')
  })
})

describe('a video with no subtitle to translate', () => {
  const pgsOnly: SubtitleTrackInfo[] = [
    {
      id: 1,
      codecId: 'S_HDMV/PGS',
      format: null,
      unsupported: 'image-based',
      language: 'en',
      default: true,
      forced: false
    }
  ]
  const toTrackStep = async (): Promise<void> => {
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    await screen.findByText('This video has no subtitle track that can be translated.')
  }

  it('cannot continue until a subtitle is brought in; a file from disk is enough', async () => {
    const { create } = setup('/videos/Show.S01E01.mkv', [])
    await toTrackStep()
    expect((screen.getByRole('button', { name: 'Continue' }) as HTMLButtonElement).disabled).toBe(
      true
    )
    await userEvent.click(screen.getByRole('button', { name: 'Choose a subtitle file' }))
    expect(await screen.findByText('Your file')).toBeTruthy()
    await next()
    await next()
    await screen.findByRole('heading', { name: 'Where to save' })
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith({
      kind: 'pipeline',
      source: '/videos/Show.S01E01.mkv',
      subtitle: { path: '/videos/Show.S01E01.mkv' },
      options: { outputMode: 'mux', setDefault: true, keepOriginalSubtitles: true },
      confirmBeforeTranslate: true
    })
  })

  it('image-only tracks: searches OpenSubtitles on its own, downloads the pick and can embed the original', async () => {
    const { create, subtitleSearch, subtitleDownload } = setup('/videos/Show.S01E01.mkv', pgsOnly)
    await toTrackStep()
    await userEvent.click(screen.getByRole('button', { name: 'Search OpenSubtitles' }))
    const dialog = await screen.findByRole('dialog')
    expect(await within(dialog).findByText('Made for this file')).toBeTruthy()
    expect(subtitleSearch).toHaveBeenCalledWith(
      {
        videoPath: '/videos/Show.S01E01.mkv',
        languages: onboardedSettings['tracks.preferredSourceLangs']
      },
      expect.any(String)
    )
    expect((within(dialog).getByLabelText('Title') as HTMLInputElement).value).toBe('Show')

    await userEvent.click(within(dialog).getByRole('button', { name: 'Use Show.S01E01.WEB' }))
    expect(subtitleDownload).toHaveBeenCalledWith(77)
    expect((await screen.findAllByText(/Show\.S01E01\.WEB/)).length).toBeGreaterThan(0)
    expect(screen.queryByRole('dialog')).toBeNull()

    await next()
    await next()
    await userEvent.click(
      await screen.findByRole('checkbox', { name: 'Also add the original subtitle to the video' })
    )
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalledWith({
      kind: 'pipeline',
      source: '/videos/Show.S01E01.mkv',
      subtitle: { path: '/cache/opensubtitles/77.srt', language: 'en' },
      options: {
        outputMode: 'mux',
        setDefault: true,
        keepOriginalSubtitles: true,
        includeSourceSubtitle: true
      },
      confirmBeforeTranslate: true
    })
  })

  it('warns before translating when the OpenSubtitles pick may be out of sync (never blocks)', async () => {
    const { create } = setup('/videos/Show.S01E01.mkv', pgsOnly, undefined, {
      result: { hashMatch: false, fps: 25 },
      videoFps: 23.976
    })
    await toTrackStep()
    await userEvent.click(screen.getByRole('button', { name: 'Search OpenSubtitles' }))
    const dialog = await screen.findByRole('dialog')
    expect(await within(dialog).findByText('25 fps ≠ video 23.976 fps')).toBeTruthy()
    expect(within(dialog).getByText('Not verified for this file')).toBeTruthy()
    expect(within(dialog).queryByText('Made for this file')).toBeNull()

    await userEvent.click(within(dialog).getByRole('button', { name: 'Use Show.S01E01.WEB' }))
    expect(
      await screen.findByText(/timed for 25 fps, but the video runs at 23\.976 fps/)
    ).toBeTruthy()
    expect(
      screen.getByText(/could not confirm this subtitle was made for this exact file/)
    ).toBeTruthy()
    // Summary row of the track step carries the short warning.
    expect(
      within(screen.getByRole('complementary', { name: 'Summary' })).getByText(
        '25 fps ≠ video 23.976 fps · Not verified for this file'
      )
    ).toBeTruthy()

    await next()
    await next()
    await userEvent.click(screen.getByRole('button', { name: 'Estimate cost' }))
    expect(create).toHaveBeenCalled()
    // Preflight review repeats the notice.
    expect(
      await screen.findByText(/timed for 25 fps, but the video runs at 23\.976 fps/)
    ).toBeTruthy()
  })

  it('hash match without a video frame rate: no sync warning', async () => {
    setup('/videos/Show.S01E01.mkv', pgsOnly, undefined, { result: { fps: 25 } })
    await toTrackStep()
    await userEvent.click(screen.getByRole('button', { name: 'Search OpenSubtitles' }))
    const dialog = await screen.findByRole('dialog')
    expect(await within(dialog).findByText('Made for this file')).toBeTruthy()
    expect(within(dialog).queryByText(/fps ≠ video/)).toBeNull()
    expect(within(dialog).queryByText('Not verified for this file')).toBeNull()
    await userEvent.click(within(dialog).getByRole('button', { name: 'Use Show.S01E01.WEB' }))
    await screen.findAllByText(/Show\.S01E01\.WEB/)
    expect(screen.queryByText(/out of sync/)).toBeNull()
  })

  it('a video with tracks still offers another subtitle, and picking one drops the track', async () => {
    setup('/videos/Show.S01E01.mkv')
    await openSingle()
    await userEvent.click(await screen.findByRole('button', { name: 'Browse files' }))
    await next()
    expect(await screen.findByText('Use another subtitle')).toBeTruthy()
    const radio = (): HTMLElement => screen.getByRole('radio', { checked: true })
    expect(radio()).toBeTruthy()
    await userEvent.click(screen.getByRole('button', { name: 'Choose a subtitle file' }))
    await screen.findByText('Your file')
    expect(screen.queryByRole('radio', { checked: true })).toBeNull()
    await userEvent.click(screen.getByRole('button', { name: 'Remove this subtitle' }))
    expect(await screen.findByRole('radio', { checked: true })).toBeTruthy()
  })
})
