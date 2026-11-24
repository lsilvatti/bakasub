import { describe, expect, it } from 'vitest'
import type { SubtitleTrackInfo } from '@shared/ipc-contract'
import {
  buildJobRequest,
  canContinue,
  initialWizard,
  pickTrack,
  stepsFor,
  wizardReducer,
  type WizardAction,
  type WizardState
} from './wizard-state'

const run = (...actions: WizardAction[]): WizardState =>
  actions.reduce(wizardReducer, initialWizard)

const track = (id: number, over: Partial<SubtitleTrackInfo> = {}): SubtitleTrackInfo => ({
  id,
  codecId: 'S_TEXT/ASS',
  format: 'ass',
  default: false,
  forced: false,
  ...over
})

describe('wizard', () => {
  it('a subtitle file skips the track step', () => {
    expect(stepsFor('subtitle')).toEqual(['file', 'preset', 'output'])
    const s = run({ type: 'file', path: '/a/ep01.en.srt' }, { type: 'next' })
    expect(s.step).toBe('preset')
    expect(s.file).toEqual({ path: '/a/ep01.en.srt', name: 'ep01.en.srt', kind: 'subtitle' })
  })

  it('a video goes through the track step and needs a track to continue', () => {
    let s = run({ type: 'file', path: 'C:\\v\\ep01.mkv' }, { type: 'next' })
    expect(s.step).toBe('track')
    expect(wizardReducer(s, { type: 'next' }).step).toBe('track')
    s = wizardReducer(s, { type: 'track', trackId: 3 })
    expect(wizardReducer(s, { type: 'next' }).step).toBe('preset')
  })

  it('unsupported files are ignored; a new file resets the track', () => {
    expect(run({ type: 'file', path: '/a/x.mp4' }).file).toBeNull()
    const s = run(
      { type: 'file', path: '/a/1.mkv' },
      { type: 'track', trackId: 2 },
      { type: 'file', path: '/a/2.mkv' }
    )
    expect(s.trackId).toBeNull()
  })

  it('preset step: model and language may come from the defaults; notes are capped', () => {
    const s = run({ type: 'file', path: '/a/x.srt' }, { type: 'next' })
    expect(canContinue(s)).toBe(false)
    expect(canContinue(s, { model: 'a/b', targetLang: 'pt-BR' })).toBe(true)
    const long = wizardReducer(s, { type: 'set', patch: { notes: 'x'.repeat(4001) } })
    expect(canContinue(long, { model: 'a/b', targetLang: 'pt-BR' })).toBe(false)
  })

  it('sends memory only when the user chose to ignore saved translations', () => {
    const s = run({ type: 'file', path: '/a/x.srt' })
    expect(initialWizard.memory).toBe('use')
    expect(buildJobRequest(s)).not.toHaveProperty('memory')
    const skip = wizardReducer(s, { type: 'set', patch: { memory: 'skip' } })
    expect(buildJobRequest(skip)).toMatchObject({ memory: 'skip' })
  })

  it('builds translate and pipeline requests, leaving defaults to the Main', () => {
    expect(
      buildJobRequest(
        run({ type: 'file', path: '/a/x.srt' }, { type: 'set', patch: { notes: '  hi ' } })
      )
    ).toEqual({ kind: 'translate', source: '/a/x.srt', notes: 'hi', confirmBeforeTranslate: true })
    expect(
      buildJobRequest(
        run(
          { type: 'file', path: '/a/x.mkv' },
          { type: 'track', trackId: 4 },
          { type: 'set', patch: { model: 'a/b', targetLang: 'es', outputDir: '/out' } },
          { type: 'mux', patch: { setDefault: false } }
        )
      )
    ).toEqual({
      kind: 'pipeline',
      source: '/a/x.mkv',
      trackId: 4,
      model: 'a/b',
      targetLang: 'es',
      outputDir: '/out',
      confirmBeforeTranslate: true,
      options: { outputMode: 'mux', setDefault: false, keepOriginalSubtitles: true }
    })
  })

  describe('output mode', () => {
    const video = run(
      { type: 'file', path: '/a/x.mkv' },
      { type: 'track', trackId: 4 },
      { type: 'mux', patch: { setDefault: false } }
    )
    const options = (s: typeof video, mode?: 'mux' | 'subtitle' | 'both'): unknown =>
      (buildJobRequest(s, null, mode) as { options?: unknown }).options

    it('follows the remembered default until the user picks', () => {
      expect(options(video)).toEqual({
        outputMode: 'mux',
        setDefault: false,
        keepOriginalSubtitles: true
      })
      expect(options(video, 'subtitle')).toEqual({ outputMode: 'subtitle' })
      const picked = wizardReducer(video, { type: 'set', patch: { outputMode: 'both' } })
      expect(options(picked, 'subtitle')).toEqual({
        outputMode: 'both',
        setDefault: false,
        keepOriginalSubtitles: true
      })
    })

    it('subtitle only sends no mux options (the MKV is never touched)', () => {
      const s = wizardReducer(video, { type: 'set', patch: { outputMode: 'subtitle' } })
      expect(options(s)).toEqual({ outputMode: 'subtitle' })
    })

    it('overwrite goes in the options only when on, for videos and subtitle files', () => {
      const on = wizardReducer(video, { type: 'set', patch: { overwrite: true } })
      expect(options(on)).toMatchObject({ outputMode: 'mux', overwrite: true })
      const srt = run(
        { type: 'file', path: '/a/x.srt' },
        { type: 'set', patch: { overwrite: true } }
      )
      expect(buildJobRequest(srt)).toMatchObject({
        kind: 'translate',
        options: { overwrite: true }
      })
      expect(buildJobRequest(run({ type: 'file', path: '/a/x.srt' }))).not.toHaveProperty('options')
    })
  })
})

describe('pickTrack', () => {
  it('prefers the preferred language, full over forced, and skips unsupported tracks', () => {
    const tracks = [
      track(1, { format: null, unsupported: 'image-based', language: 'en' }),
      track(2, { language: 'en', forced: true }),
      track(3, { language: 'en' }),
      track(4, { language: 'ja' })
    ]
    expect(pickTrack(tracks, ['en'])).toBe(3)
    expect(pickTrack(tracks, ['ja', 'en'])).toBe(4)
    expect(pickTrack(tracks, ['fr'])).toBe(3)
    expect(pickTrack([track(1, { format: null })], ['en'])).toBeNull()
  })
})

describe('next on the preset step', () => {
  it('advances with the defaults from Settings', () => {
    const s = run({ type: 'file', path: '/a/x.srt' }, { type: 'next' })
    expect(wizardReducer(s, { type: 'next' }).step).toBe('preset')
    expect(
      wizardReducer(s, { type: 'next', defaults: { model: 'a/b', targetLang: 'pt-BR' } }).step
    ).toBe('output')
  })
})

const mediaOf = (s: WizardState): unknown => (buildJobRequest(s) as { media?: unknown }).media

describe('media in the request', () => {
  const series = { provider: 'tmdb' as const, id: 7, kind: 'series' as const }
  it('sent only when included; the episode only for a series with both numbers', () => {
    const base = run({ type: 'file', path: '/a/x.srt' })
    const withMedia = wizardReducer(base, {
      type: 'media',
      patch: { ref: series, season: 1, episode: 5 }
    })
    expect(buildJobRequest(withMedia)).toMatchObject({
      media: { ref: series, season: 1, episode: 5 }
    })
    expect(mediaOf(wizardReducer(withMedia, { type: 'media', patch: { episode: null } }))).toEqual({
      ref: series
    })
    expect(
      buildJobRequest(wizardReducer(withMedia, { type: 'media', patch: { include: false } }))
    ).not.toHaveProperty('media')
    const movie = wizardReducer(withMedia, {
      type: 'media',
      patch: { ref: { ...series, kind: 'movie' } }
    })
    expect(mediaOf(movie)).toEqual({ ref: { ...series, kind: 'movie' } })
  })
  it('a new file clears the media choice', () => {
    const s = run(
      { type: 'file', path: '/a/x.srt' },
      { type: 'media', patch: { ref: series, seeded: true } },
      { type: 'file', path: '/a/y.srt' }
    )
    expect(s.media.ref).toBeNull()
    expect(s.media.seeded).toBe(false)
  })
})

describe('step bar navigation', () => {
  const d = { model: 'a/b', targetLang: 'pt-BR' }
  it('reached steps can be revisited forward and back without "Continue"', () => {
    let s = run(
      { type: 'file', path: '/a/x.srt' },
      { type: 'next' },
      { type: 'next', defaults: d },
      { type: 'go', step: 'file' }
    )
    expect(s.step).toBe('file')
    s = wizardReducer(s, { type: 'go', step: 'output', defaults: d })
    expect(s.step).toBe('output')
  })

  it('a step never reached cannot be jumped to', () => {
    const s = run({ type: 'file', path: '/a/x.srt' })
    expect(wizardReducer(s, { type: 'go', step: 'output', defaults: d }).step).toBe('file')
  })

  it('changing the file resets the reached steps and says so; the same file changes nothing', () => {
    let s = run(
      { type: 'file', path: '/a/x.mkv' },
      { type: 'track', trackId: 2 },
      { type: 'next' },
      { type: 'next' }
    )
    expect(s.furthest).toBe(2)
    expect(wizardReducer(s, { type: 'file', path: '/a/x.mkv' })).toBe(s)
    s = wizardReducer(wizardReducer(s, { type: 'go', step: 'file' }), {
      type: 'file',
      path: '/a/y.mkv'
    })
    expect(s).toMatchObject({ furthest: 0, notice: 'file-changed', trackId: null })
    expect(wizardReducer(s, { type: 'go', step: 'preset', defaults: d }).step).toBe('file')
  })
})

describe('external subtitle', () => {
  const video = { type: 'file', path: '/a/ep01.mkv' } as const
  const external = { path: '/subs/ep01.srt', name: 'ep01.srt', origin: 'file' } as const

  it('lets a video with no track continue, and excludes the track', () => {
    const s = run(video, { type: 'next' })
    expect(canContinue(s)).toBe(false)
    const withSub = wizardReducer(s, { type: 'subtitle', subtitle: external })
    expect(canContinue(withSub)).toBe(true)

    const swapped = run(
      video,
      { type: 'track', trackId: 2 },
      { type: 'subtitle', subtitle: external }
    )
    expect(swapped.trackId).toBeNull()
    expect(wizardReducer(swapped, { type: 'track', trackId: 3 })).toMatchObject({
      trackId: 3,
      subtitle: null
    })
  })

  it('removing it keeps the track step incomplete; a new video clears it', () => {
    const s = run(video, { type: 'subtitle', subtitle: external })
    expect(wizardReducer(s, { type: 'subtitle', subtitle: null }).subtitle).toBeNull()
    expect(wizardReducer(s, { type: 'file', path: '/a/ep02.mkv' }).subtitle).toBeNull()
  })

  it('builds a pipeline request with the subtitle instead of a track', () => {
    const s = run(video, {
      type: 'subtitle',
      subtitle: { ...external, origin: 'opensubtitles', language: 'en' }
    })
    const req = buildJobRequest(s)
    expect(req).toMatchObject({
      kind: 'pipeline',
      source: '/a/ep01.mkv',
      subtitle: { path: '/subs/ep01.srt', language: 'en' }
    })
    expect(req).not.toHaveProperty('trackId')
    expect(req).not.toHaveProperty('options.includeSourceSubtitle')
  })

  it('"also add the original" is sent only when an MKV is written', () => {
    const s = run(
      video,
      { type: 'subtitle', subtitle: external },
      { type: 'set', patch: { includeOriginal: true } }
    )
    expect(buildJobRequest(s)).toMatchObject({ options: { includeSourceSubtitle: true } })
    const subtitleOnly = wizardReducer(s, { type: 'set', patch: { outputMode: 'subtitle' } })
    expect(buildJobRequest(subtitleOnly)).not.toHaveProperty('options.includeSourceSubtitle')
  })
})
