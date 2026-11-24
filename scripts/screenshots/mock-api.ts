/**
 * In-memory `window.api` for the marketing screenshots: the renderer runs in a plain browser
 * (no Electron, no Main) against a fictional library. Typed against the IPC contract, so a
 * contract change that the mock no longer satisfies fails `tsc`.
 *
 * Loaded before `main.tsx` by `vite.config.mts`. Nothing here ships with the app.
 */
import type {
  BakaSubApi,
  CoverageReport,
  EpisodeDetails,
  IpcResult,
  Job,
  MediaDetails,
  ModelSummary,
  PreflightEstimate,
  Preset,
  PromptSummary,
  PublicSettings
} from '@shared/ipc-contract'
import { BUILTIN_PRESETS } from '@main/presets/builtin-presets'
import { backdrop, poster } from './art'

/** Like the real IPC, answers are copies: the renderer's cache never aliases the mock's state. */
const ok = <T>(data: T): Promise<IpcResult<T>> =>
  Promise.resolve({ ok: true, data: structuredClone(data) })
const NOW = Date.now()
const ago = (min: number): string => new Date(NOW - min * 60_000).toISOString()

const LIBRARY = '/home/user/Videos/Starlight Courier'
const VIDEO = `${LIBRARY}/Starlight.Courier.S01E07.1080p.WEB-DL.mkv`
const ART = { backdropUrl: backdrop(), posterUrl: poster() }

const show: MediaDetails = {
  provider: 'tmdb',
  id: 424242,
  kind: 'series',
  title: 'Starlight Courier',
  originalTitle: 'Hoshiakari no Haitatsunin',
  year: 2024,
  ...ART,
  genres: ['Animation', 'Sci-Fi & Fantasy', 'Drama'],
  genreIds: [16, 10765, 18],
  originalLanguage: 'ja',
  originCountries: ['JP'],
  runtime: 24,
  synopsis:
    'Mika, a rookie courier on the last night train between the colonies of Orion, delivers letters nobody else will carry. Each stop hides a stranger with something left unsaid.'
}
const episode: EpisodeDetails = {
  season: 1,
  episode: 7,
  title: 'The Last Train to Orion',
  synopsis:
    'A storm strands the train at a ghost station, and Mika must decide whether to open a letter addressed to herself.'
}
const mediaSummary = (ep = episode): Job['media'] => ({
  title: show.title,
  kind: 'series',
  year: show.year,
  ...ART,
  episode: { season: ep.season, episode: ep.episode, title: ep.title }
})

const settings: PublicSettings = {
  'translation.defaultModel': 'google/gemini-2.5-flash',
  'translation.defaultTargetLang': 'pt-BR',
  'translation.defaultPresetId': 'builtin:series',
  'translation.targetLangs': ['pt-BR', 'es', 'fr'],
  'models.favorites': [
    'google/gemini-2.5-flash',
    'anthropic/claude-sonnet-4.5',
    'deepseek/deepseek-chat-v3.1'
  ],
  'jobs.concurrency': 2,
  'preflight.autoConfirmBelowUsd': null,
  'tracks.preferredSourceLangs': ['en'],
  'logs.level': 'info',
  'updates.autoCheck': true,
  'notifications.enabled': true,
  'notifications.job': true,
  'notifications.batch': true,
  'notifications.update': true,
  'ui.language': 'en',
  'output.mode': 'mux',
  'output.extractNaming': 'player',
  'llm.rateLimit.freeRpm': 20,
  'llm.rateLimit.paidRpm': null,
  'llm.rateLimit.freeTier': 'auto',
  'llm.retry.maxThrottled': 5,
  'llm.retry.maxTransient': 3,
  'llm.retry.baseDelayMs': 1000,
  'llm.retry.maxDelayMs': 30000,
  'llm.retry.chunkRounds': 2,
  'llm.concurrency.perJob': 4,
  'llm.concurrency.paidMaxInFlight': 8,
  'llm.concurrency.freeMaxInFlight': 3
}

const model = (
  id: string,
  name: string,
  prompt: number,
  completion: number,
  tokenizer = 'Other'
): ModelSummary => ({
  id,
  name,
  contextLength: 1_000_000,
  maxCompletionTokens: 65_536,
  isModerated: false,
  pricing: { prompt: prompt / 1e6, completion: completion / 1e6 },
  structuredOutputs: true,
  responseFormat: true,
  tokenizer
})
const models: ModelSummary[] = [
  model('google/gemini-2.5-flash', 'Google: Gemini 2.5 Flash', 0.3, 2.5, 'Gemini'),
  model('anthropic/claude-sonnet-4.5', 'Anthropic: Claude Sonnet 4.5', 3, 15, 'Claude'),
  model('deepseek/deepseek-chat-v3.1', 'DeepSeek: DeepSeek V3.1', 0.2, 0.8),
  model('openai/gpt-4.1-mini', 'OpenAI: GPT-4.1 Mini', 0.4, 1.6, 'GPT'),
  model('meta-llama/llama-3.3-70b-instruct:free', 'Meta: Llama 3.3 70B (free)', 0, 0, 'Llama')
]

const customPreset: Preset = {
  id: 'p-starlight',
  builtin: false,
  name: 'Starlight Courier',
  mediaType: 'Anime',
  temperature: 0.4,
  prompt: [
    BUILTIN_PRESETS.find((p) => p.id === 'builtin:anime')!.prompt,
    '- Mika speaks casually; the Conductor always speaks formally and archaically.'
  ].join('\n'),
  glossary: [
    { source: 'Mika-chan', rule: 'Keep the honorific.' },
    { source: 'Orion Line', target: 'Linha Órion' },
    { source: 'Conductor', target: 'Maquinista', rule: 'Title, always capitalized.' },
    { source: 'Lumen Post', target: 'Correio Lúmen' },
    { source: 'senpai' },
    { source: 'ghost station', target: 'estação fantasma' }
  ],
  copiedFrom: 'builtin:anime',
  createdAt: ago(60 * 24 * 9)
}
const presets: Preset[] = [
  ...BUILTIN_PRESETS.map((p) => ({ ...p, glossary: [...p.glossary] })),
  customPreset
]

const timeline = (
  total: number,
  done: number
): NonNullable<Job['progress']['translation']>['timeline'] => ({
  total,
  cached: [[0, 64]],
  translated: [
    [64, done],
    [Math.round(total * 0.5), Math.round(total * 0.5) + done - 120]
  ],
  failed: [],
  active: [
    [done, done + 24],
    [Math.round(total * 0.5) + done - 120, Math.round(total * 0.5) + done - 96]
  ],
  lanes: [
    [0, Math.round(total * 0.5)],
    [Math.round(total * 0.5), total]
  ]
})

const usage = (cost: number, requests: number): NonNullable<Job['result']>['usage'] => ({
  requests,
  promptTokens: requests * 5200,
  completionTokens: requests * 1900,
  cachedTokens: requests * 900,
  cost,
  costIncomplete: false
})

const episodeN = (n: number, title: string): EpisodeDetails => ({ season: 1, episode: n, title })
const baseJob = (over: Partial<Job> & Pick<Job, 'id' | 'status' | 'source'>): Job => ({
  kind: 'pipeline',
  batchId: null,
  step: null,
  progress: { step: null },
  result: null,
  error: null,
  attempt: 1,
  createdAt: ago(40),
  updatedAt: ago(1),
  startedAt: ago(38),
  finishedAt: null,
  ...over
})

const jobs: Job[] = [
  baseJob({
    id: 'j-08',
    status: 'running',
    step: 'translate',
    source: 'Starlight.Courier.S01E08.1080p.WEB-DL.mkv',
    media: mediaSummary(episodeN(8, 'Letters from the Ninth Moon')),
    progress: {
      step: 'translate',
      translation: {
        done: 238,
        total: 402,
        cached: 64,
        chunk: { completed: 7, count: 12, inFlight: 2 },
        timeline: timeline(402, 238)
      },
      cost: 0.0142
    }
  }),
  baseJob({
    id: 'j-09',
    status: 'queued',
    source: 'Starlight.Courier.S01E09.1080p.WEB-DL.mkv',
    media: mediaSummary(episodeN(9, 'Return to Sender')),
    startedAt: null
  }),
  baseJob({
    id: 'j-10',
    status: 'queued',
    source: 'Starlight.Courier.S01E10.1080p.WEB-DL.mkv',
    media: mediaSummary(episodeN(10, 'Postmarked Tomorrow')),
    startedAt: null
  }),
  baseJob({
    id: 'j-06',
    status: 'completed',
    step: 'mux',
    source: 'Starlight.Courier.S01E06.1080p.WEB-DL.mkv',
    media: mediaSummary(episodeN(6, 'Stamps and Stardust')),
    progress: { step: 'mux', stepFraction: 1 },
    createdAt: ago(95),
    startedAt: ago(94),
    finishedAt: ago(82),
    result: {
      outputs: [`${LIBRARY}/Starlight.Courier.S01E06.1080p.WEB-DL.pt-BR.mkv`],
      warnings: [],
      usage: usage(0.0231, 12),
      cache: { hits: 71, misses: 333 }
    }
  }),
  baseJob({
    id: 'j-05',
    status: 'completed-with-warnings',
    step: 'mux',
    source: 'Starlight.Courier.S01E05.1080p.WEB-DL.mkv',
    media: mediaSummary(episodeN(5, 'A Parcel for Nobody')),
    progress: { step: 'mux', stepFraction: 1 },
    createdAt: ago(130),
    startedAt: ago(129),
    finishedAt: ago(117),
    result: {
      outputs: [`${LIBRARY}/Starlight.Courier.S01E05.1080p.WEB-DL.pt-BR.mkv`],
      warnings: [],
      usage: usage(0.0198, 11),
      cache: { hits: 40, misses: 358 },
      failedLines: [
        {
          cueId: 211,
          reason: 'UNTRANSLATED_VALIDATION_FAILED',
          excerpt: '{\\i1}Next time on Starlight Courier...{\\i0}',
          location: { line: 455, endLine: 455, startMs: 1_402_300, endMs: 1_405_100 }
        }
      ]
    }
  }),
  baseJob({
    id: 'j-mv',
    kind: 'translate',
    status: 'completed',
    source: 'Paper.Lanterns.2019.720p.BluRay.en.srt',
    progress: { step: 'write', stepFraction: 1 },
    createdAt: ago(60 * 26),
    startedAt: ago(60 * 26),
    finishedAt: ago(60 * 26 - 4),
    result: {
      outputs: ['/home/user/Videos/Paper.Lanterns.2019.720p.BluRay.pt-BR.srt'],
      subtitle: '/home/user/Videos/Paper.Lanterns.2019.720p.BluRay.pt-BR.srt',
      warnings: [],
      usage: usage(0.0412, 18),
      cache: { hits: 0, misses: 1204 }
    }
  })
]

const estimate: PreflightEstimate = {
  model: {
    id: 'google/gemini-2.5-flash',
    pricing: models[0].pricing,
    catalogFetchedAt: ago(12)
  },
  approximate: true,
  lines: { total: 412, cached: 58, toTranslate: 331, copies: 23, chunks: 11 },
  tokens: { input: { low: 41_800, high: 62_700 }, output: { low: 10_900, high: 16_400 } },
  cost: { low: 0.0398, high: 0.0598 },
  decision: 'needs-confirmation'
}

const coverage: CoverageReport = {
  totals: { cues: 412, blocks: 0, translatable: 389, skipped: 22, untranslated: 1 },
  byReason: { SKIPPED_DRAWING: 14, SKIPPED_NO_LETTERS: 8, UNTRANSLATED_KARAOKE: 1 },
  entries: [
    {
      ref: { kind: 'cue', cueId: 3 },
      status: 'untranslated',
      severity: 'warn',
      reason: 'UNTRANSLATED_KARAOKE',
      params: {},
      location: { line: 31, endLine: 31, startMs: 62_000, endMs: 71_500 },
      excerpt: '{\\k32}Ho{\\k28}shi{\\k40}a{\\k35}ka{\\k30}ri',
      outcome: 'kept-original',
      actions: []
    }
  ],
  truncated: false
}

const prompt: PromptSummary = {
  system: 'You are a professional subtitle translator…',
  systemTokens: 1184,
  glossaryEntries: customPreset.glossary.length,
  style: true,
  media: { tmdb: true, episode: true, notes: false }
}

const listeners = new Map<string, Set<(p: unknown) => void>>()
let seq = 0

const api: BakaSubApi = {
  app: {
    getInfo: () =>
      ok({
        name: 'BakaSub',
        version: '1.0.7',
        platform: 'linux',
        arch: 'x64',
        electron: '39.0.0',
        userDataPath: '~/.config/bakasub'
      }),
    checkUpdate: () =>
      ok({
        current: '1.0.7',
        latest: '1.0.7',
        updateAvailable: false,
        publishedAt: '',
        offline: false
      }),
    getUpdateState: () => ok({ state: 'not-available', current: '1.0.7' }),
    downloadUpdate: () => ok({ state: 'not-available', current: '1.0.7' }),
    cancelUpdateDownload: () => ok({ cancelled: false }),
    installUpdate: () => ok({ installing: true }),
    openExternal: () => ok({ opened: true })
  },
  sidecars: {
    check: () =>
      ok({
        checkedAt: ago(1),
        manifestUpdatedAt: ago(60 * 24),
        offline: false,
        installDir: '~/.config/bakasub/bin',
        ready: true,
        tools: ['ffmpeg', 'mkvextract', 'mkvmerge'].map((name) => ({
          name,
          requiredVersion: name === 'ffmpeg' ? '7.1' : '92.0',
          installedVersion: name === 'ffmpeg' ? '7.1' : '92.0',
          state: 'ok' as const,
          path: `~/.config/bakasub/bin/${name}`
        }))
      }),
    install: () => ok({ outcome: 'declined' }),
    cancel: () => ok({ cancelled: false })
  },
  memory: {
    stats: () => ok({ rows: 18_342, bytes: 6_815_744, oldestUsedAt: ago(60 * 24 * 40) }),
    prune: () =>
      ok({ deleted: 0, stats: { rows: 18_342, bytes: 6_815_744, oldestUsedAt: ago(60 * 24 * 40) } })
  },
  settings: {
    setApiKey: () => ok({ present: true, protection: 'os' }),
    clearApiKey: () => ok({ present: false, protection: 'os' }),
    apiKeyStatus: () => ok({ present: true, protection: 'os' }),
    setTmdbKey: () => ok({ present: true }),
    clearTmdbKey: () => ok({ present: false }),
    tmdbKeyStatus: () => ok({ present: true }),
    get: () => ok(settings),
    update: (patch) => ok(Object.assign(settings, patch)),
    reset: () => ok(settings)
  },
  logs: {
    tail: () => ok([]),
    watch: (on) => ok({ watching: on }),
    openFolder: () => ok({ opened: true }),
    export: () => ok({ path: null }),
    report: () => ok({ accepted: true })
  },
  presets: {
    list: () => ok(presets),
    create: (input) =>
      ok({ ...input, id: `p-${++seq}`, builtin: false, copiedFrom: null, createdAt: ago(0) }),
    update: (id, input) => ok({ ...presets.find((p) => p.id === id)!, ...input }),
    duplicate: (id) =>
      ok({ ...presets.find((p) => p.id === id)!, id: `p-${++seq}`, builtin: false }),
    delete: () => ok({ deleted: true }),
    export: () => ok({ path: null, count: 0 }),
    import: () => ok({ path: null, imported: [], skipped: [] })
  },
  models: {
    list: () => ok({ models, fetchedAt: ago(12), offline: false }),
    freeQuota: () => ok({ tier: 'over-10', source: 'auto', limit: 1000, used: 37, remaining: 963 })
  },
  media: {
    probe: () =>
      ok({
        attachments: 6,
        subtitles: [
          {
            id: 2,
            codecId: 'S_TEXT/ASS',
            format: 'ass',
            language: 'en',
            name: 'Full Subtitles',
            default: true,
            forced: false,
            eventCount: 412
          },
          {
            id: 3,
            codecId: 'S_TEXT/ASS',
            format: 'ass',
            language: 'en',
            name: 'Signs & Songs',
            default: false,
            forced: true,
            eventCount: 46
          },
          {
            id: 4,
            codecId: 'S_TEXT/UTF8',
            format: 'srt',
            language: 'es',
            name: 'Español (Latinoamérica)',
            default: false,
            forced: false,
            eventCount: 398
          },
          {
            id: 5,
            codecId: 'S_HDMV/PGS',
            format: null,
            unsupported: 'image-based',
            language: 'ja',
            name: 'Japanese',
            default: false,
            forced: false
          }
        ]
      }),
    identify: () =>
      ok({
        release: { title: 'Starlight Courier', season: 1, episode: 7 },
        status: 'matched',
        match: show,
        episode
      }),
    search: () => ok([show]),
    lookup: () => ok({ details: show, episode })
  },
  dialog: {
    openFiles: () => ok({ paths: [VIDEO] }),
    chooseOutputDir: () => ok({ path: null }),
    openFolder: () => ok({ path: LIBRARY })
  },
  batches: {
    plan: (() =>
      Promise.reject(new Error('not in the screenshots'))) as BakaSubApi['batches']['plan'],
    create: () => Promise.reject(new Error('not in the screenshots')),
    list: () => ok([]),
    get: () => Promise.reject(new Error('not in the screenshots')),
    jobs: () => ok([]),
    cancel: () => Promise.reject(new Error('not in the screenshots')),
    remove: () => ok({ removed: true }),
    confirm: () => Promise.reject(new Error('not in the screenshots'))
  },
  jobs: {
    create: (req) => {
      const job = baseJob({
        id: `j-new-${++seq}`,
        kind: req.kind,
        status: 'awaiting-confirmation',
        step: 'preflight',
        source: req.source.split('/').pop()!,
        media: mediaSummary(),
        createdAt: ago(0),
        startedAt: ago(0),
        result: { outputs: [], warnings: [], preflight: estimate }
      })
      jobs.unshift(job)
      return ok(job)
    },
    list: () => ok(jobs),
    get: (id) => {
      const job = jobs.find((j) => j.id === id)
      return job
        ? ok(job)
        : Promise.resolve({ ok: false, error: { code: 'NOT_FOUND', message: 'Job not found' } })
    },
    cancel: (id) => ok(jobs.find((j) => j.id === id)!),
    retry: (id) => ok(jobs.find((j) => j.id === id)!),
    remove: () => ok({ removed: true }),
    clearCompleted: () => ok({ removed: [] }),
    confirm: (id) => ok(jobs.find((j) => j.id === id)!),
    dismissConsistency: (id) => ok(jobs.find((j) => j.id === id)!),
    preflight: () => ok({ estimate, coverage, prompt })
  },
  files: {
    inspect: () =>
      ok({
        file: { name: 'subtitle.ass', sizeBytes: 48_210 },
        format: 'ass',
        encoding: { name: 'utf-8', confidence: 'certain', score: null, bom: false, eol: 'lf' },
        diagnostics: { items: [], totals: { info: 0, warn: 0, error: 0 }, truncated: false },
        coverage
      }),
    pathFor: (file) => `/home/user/Videos/${file.name}`
  },
  events: {
    on(channel, listener) {
      const set = listeners.get(channel) ?? new Set()
      set.add(listener as (p: unknown) => void)
      listeners.set(channel, set)
      return () => set.delete(listener as (p: unknown) => void)
    }
  }
}

Object.defineProperty(window, 'api', { value: api, configurable: true })
