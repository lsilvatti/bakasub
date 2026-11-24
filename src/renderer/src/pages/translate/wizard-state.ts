import type { JobCreateRequest, JobMemoryMode, OutputMode } from '@shared/ipc-contract'
import { noMedia, type MediaChoice } from '@renderer/components/media/media-choice'
import type { ExternalSubtitleValue } from '@renderer/components/media/external-subtitle'
import { sourceKindOf, type SourceKind } from '@shared/file-types'

export type WizardStep = 'file' | 'track' | 'preset' | 'output'

export type ExternalSubtitle = ExternalSubtitleValue

export interface WizardState {
  step: WizardStep
  file: { path: string; name: string; kind: SourceKind } | null
  trackId: number | null
  subtitle: ExternalSubtitle | null
  includeOriginal: boolean
  /** `null` = the TMDB suggestion for the work, else the default from Settings. */
  presetId: string | null
  model: string | null
  targetLang: string | null
  notes: string
  /** `null` = next to the source file. */
  outputDir: string | null
  /** `null` = the last choice remembered in Settings (`output.mode`), else `mux`. Videos only. */
  outputMode: OutputMode | null
  /** Replace a file that already has the output name instead of writing ` (2)`. */
  overwrite: boolean
  mux: { setDefault: boolean; keepOriginalSubtitles: boolean }
  /** `skip` = "Ignore saved translations" (off by default). */
  memory: JobMemoryMode
  media: MediaChoice
  /** Furthest step reached: completed steps can be revisited from the step bar. */
  furthest: number
  /** Something changed that made the later steps need another look. */
  notice: 'file-changed' | null
}

export const initialWizard: WizardState = {
  step: 'file',
  file: null,
  trackId: null,
  subtitle: null,
  includeOriginal: false,
  presetId: null,
  model: null,
  targetLang: null,
  notes: '',
  outputDir: null,
  outputMode: null,
  overwrite: false,
  mux: { setDefault: true, keepOriginalSubtitles: true },
  memory: 'use',
  media: noMedia,
  furthest: 0,
  notice: null
}

/** A subtitle file has no track to pick and nothing to mux. */
export const stepsFor = (kind: SourceKind | undefined): WizardStep[] =>
  kind === 'subtitle' ? ['file', 'preset', 'output'] : ['file', 'track', 'preset', 'output']

export const NOTES_MAX = 4000

export type WizardAction =
  | { type: 'file'; path: string }
  | { type: 'track'; trackId: number }
  | { type: 'subtitle'; subtitle: ExternalSubtitle | null }
  | {
      type: 'set'
      patch: Partial<
        Pick<
          WizardState,
          | 'presetId'
          | 'model'
          | 'targetLang'
          | 'notes'
          | 'outputDir'
          | 'outputMode'
          | 'overwrite'
          | 'memory'
          | 'includeOriginal'
        >
      >
    }
  | { type: 'mux'; patch: Partial<WizardState['mux']> }
  | { type: 'media'; patch: Partial<MediaChoice> }
  /** Jump from the step bar: any reached step whose previous steps are all valid. */
  | { type: 'go'; step: WizardStep; defaults?: Defaults }
  /** `defaults`: model/language from Settings, used when the user kept them. */
  | { type: 'next'; defaults?: Defaults }
  | { type: 'back' }
  | { type: 'reset' }

const baseName = (path: string): string => path.split(/[\\/]/).pop() ?? path

export function wizardReducer(s: WizardState, a: WizardAction): WizardState {
  const steps = stepsFor(s.file?.kind)
  const i = steps.indexOf(s.step)
  switch (a.type) {
    case 'file': {
      const kind = sourceKindOf(a.path)
      if (!kind) return s
      if (s.file?.path === a.path) return s
      // A new file invalidates the track and the media; preset, notes and output are kept.
      // Later steps must be walked again, and the user is told why.
      return {
        ...s,
        file: { path: a.path, name: baseName(a.path), kind },
        trackId: null,
        subtitle: null,
        media: noMedia,
        furthest: 0,
        notice: s.file ? 'file-changed' : null
      }
    }
    case 'track':
      return { ...s, trackId: a.trackId, subtitle: null }
    case 'subtitle':
      return { ...s, subtitle: a.subtitle, trackId: a.subtitle ? null : s.trackId }
    case 'set':
      return { ...s, ...a.patch }
    case 'mux':
      return { ...s, mux: { ...s.mux, ...a.patch } }
    case 'media':
      return { ...s, media: { ...s.media, ...a.patch } }
    case 'go':
      return canGo(s, a.step, a.defaults) ? { ...s, step: a.step } : s
    case 'next': {
      if (!canContinue(s, a.defaults) || i >= steps.length - 1) return s
      return { ...s, step: steps[i + 1], furthest: Math.max(s.furthest, i + 1), notice: null }
    }
    case 'back':
      return i > 0 ? { ...s, step: steps[i - 1] } : s
    case 'reset':
      return initialWizard
  }
}

/** What the current step needs before "Continue". Model/language may come from the defaults. */
export interface Defaults {
  model: string | null
  targetLang: string | null
}

export function canContinue(
  s: WizardState,
  defaults: Defaults = { model: null, targetLang: null }
): boolean {
  return stepValid(s, s.step, defaults)
}

/** Whether a step's own fields are complete. */
export function stepValid(
  s: WizardState,
  step: WizardStep,
  defaults: Defaults = { model: null, targetLang: null }
): boolean {
  switch (step) {
    case 'file':
      return s.file !== null
    case 'track':
      return s.trackId !== null || s.subtitle !== null
    case 'preset':
      return (
        !!(s.model ?? defaults.model) &&
        !!(s.targetLang ?? defaults.targetLang) &&
        s.notes.length <= NOTES_MAX
      )
    case 'output':
      return true
  }
}

/**
 * Step bar navigation: back to any step, forward to a step already reached, as long as every
 * step before it is still valid (no confirmation needed when nothing changed).
 */
export function canGo(s: WizardState, target: WizardStep, defaults?: Defaults): boolean {
  const steps = stepsFor(s.file?.kind)
  const t = steps.indexOf(target)
  if (t === -1) return false
  if (t <= steps.indexOf(s.step)) return true
  return t <= s.furthest && steps.slice(0, t).every((st) => stepValid(s, st, defaults))
}

/**
 * The `jobs:create` request. `suggestedPresetId`: the TMDB suggestion, used when the user kept
 * it. Unset model/language/preset are left to the Main's defaults. `defaultMode`: the output
 * mode remembered in Settings, used when the user did not pick one.
 */
export function buildJobRequest(
  s: WizardState,
  suggestedPresetId?: string | null,
  defaultMode: OutputMode = 'mux'
): JobCreateRequest {
  if (!s.file) throw new Error('No file selected')
  const presetId = s.presetId ?? suggestedPresetId
  const common = {
    source: s.file.path,
    ...(presetId ? { presetId } : {}),
    ...(s.model ? { model: s.model } : {}),
    ...(s.targetLang ? { targetLang: s.targetLang } : {}),
    ...(s.notes.trim() ? { notes: s.notes.trim() } : {}),
    ...(s.outputDir ? { outputDir: s.outputDir } : {}),
    ...mediaRequest(s.media),
    confirmBeforeTranslate: true,
    ...(s.memory === 'skip' ? { memory: 'skip' as const } : {})
  }
  const write = s.overwrite ? { overwrite: true } : {}
  if (s.file.kind === 'subtitle')
    return { kind: 'translate', ...common, ...(s.overwrite ? { options: write } : {}) }
  const outputMode = s.outputMode ?? defaultMode
  // Only a written video has mux options: the subtitle-only mode never touches the MKV.
  const muxed = outputMode !== 'subtitle'
  const options = { outputMode, ...(muxed ? { ...s.mux } : {}), ...write }
  if (s.subtitle)
    return {
      kind: 'pipeline',
      ...common,
      subtitle: {
        path: s.subtitle.path,
        ...(s.subtitle.language ? { language: s.subtitle.language } : {})
      },
      options: {
        ...options,
        ...(muxed && s.includeOriginal ? { includeSourceSubtitle: true } : {})
      }
    }
  if (s.trackId === null) throw new Error('No track selected')
  return { kind: 'pipeline', ...common, trackId: s.trackId, options }
}

/** `media` for the request: only when included; the episode only for a series with both numbers. */
function mediaRequest(
  m: MediaChoice
): Pick<Extract<JobCreateRequest, { kind: 'translate' }>, 'media'> {
  if (!m.include || !m.ref) return {}
  const episode =
    m.ref.kind === 'series' && m.season !== null && m.episode !== null
      ? { season: m.season, episode: m.episode }
      : {}
  return { media: { ref: m.ref, ...episode } }
}

export { pickTrack } from '@renderer/lib/tracks'
