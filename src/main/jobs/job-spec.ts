import { dirname } from 'node:path'
import { z } from 'zod'
import type {
  ExtractNaming,
  JobCreateRequest,
  JobMemoryMode,
  JobMediaRequest,
  MediaContextLookup,
  MediaType,
  OutputMode,
  Preset,
  SubtitleTextFormat
} from '../../shared/ipc-contract'
import { formatFromPath } from './steps/subtitle-io'
import { AppError } from '../errors'
import { languageName } from '../llm/prompt-builder'
import { MAX_USER_NOTES } from '../presets/preset.schema'
import { canonicalTag } from '../media/languages'
import { langFromFilename } from '../media/naming'
import {
  MAX_SUBTITLE_BYTES,
  MEDIA_EXTENSIONS,
  SUBTITLE_EXTENSIONS,
  validateInputFile,
  validateOutputDir
} from '../media/path-guard'

/**
 * Preset congelado no job, já RESOLVIDO (modelo e idioma decididos na criação): editar o preset
 * ou atualizar os embutidos depois não muda um job enfileirado nem o seu retry.
 */
export interface PresetSnapshot {
  readonly id: string
  readonly name: string
  readonly mediaType: MediaType
  readonly model: string
  readonly glossary: Preset['glossary']
  readonly temperature: number
  readonly prompt: string
}

/**
 * Busca dos detalhes do TMDB para congelar no job. Falha = job sem os detalhes, mas com o motivo
 * congelado em `mediaUnavailable` (o job avisa) e reportado em `onError` (log).
 */
export interface MediaLookupSource {
  lookup(
    ref: JobMediaRequest['ref'],
    episode: { season?: number; episode?: number }
  ): Promise<MediaContextLookup>
  onError?: (err: unknown) => void
}

async function resolveMedia(
  req: JobMediaRequest | undefined,
  source: MediaLookupSource | undefined
): Promise<{ media?: MediaContextLookup; mediaUnavailable?: string }> {
  if (!req) return {}
  if (!source) return { mediaUnavailable: 'media context service is not available' }
  try {
    return { media: await source.lookup(req.ref, { season: req.season, episode: req.episode }) }
  } catch (err) {
    source.onError?.(err)
    return { mediaUnavailable: lookupFailure(err) }
  }
}

const MAX_REASON = 300
function lookupFailure(err: unknown): string {
  const code = err instanceof AppError ? `${err.code}: ` : ''
  const msg = err instanceof Error ? err.message : String(err)
  return `TMDB lookup failed (${code}${msg})`.slice(0, MAX_REASON)
}

/** Presets + padrões das configurações (`translation.*`). */
export interface PresetSource {
  get(id: string): Promise<Preset>
  defaults(): Promise<{ presetId: string; model: string | null; targetLang: string | null }>
}

interface Base {
  readonly source: string
  readonly sourceSize: number
  readonly outputDir: string
}
export interface MuxSpecOptions {
  readonly setDefault: boolean
  readonly keepOriginalSubtitles: boolean
}
export interface WriteSpecOptions {
  /** Ausente em jobs antigos = nome único (` (2)`). */
  readonly overwrite?: boolean
}

/** Modo de saída efetivo do pipeline, inclusive de specs persistidas antes de `outputMode`. */
export function pipelineOutputMode(spec: {
  readonly outputMode?: OutputMode
  readonly alsoWriteSubtitle?: boolean
}): OutputMode {
  return spec.outputMode ?? (spec.alsoWriteSubtitle ? 'both' : 'mux')
}

export interface MuxSubtitleSpec {
  readonly path: string
  readonly size: number
  readonly language: string
  readonly name: string
  readonly default: boolean
}

export interface ExternalSubtitleSpec {
  readonly path: string
  readonly size: number
  readonly format: SubtitleTextFormat
  readonly includeInMux: boolean
}

export type JobSpec =
  | (Base &
      WriteSpecOptions & {
        readonly kind: 'extract'
        readonly trackIds: readonly number[]
        /** Ausente em jobs antigos = `track`. */
        readonly naming?: ExtractNaming
      })
  | (Base &
      WriteSpecOptions & {
        readonly kind: 'translate'
        readonly preset: PresetSnapshot
        readonly targetLang: string
        readonly sourceLang?: string
        /** Contexto do usuário (vai no `<media>`). */
        readonly notes?: string
        /** Detalhes do TMDB congelados na criação (retry usa os mesmos, mesmo offline). */
        readonly media?: MediaContextLookup
        /** O usuário pediu o contexto do TMDB e ele não pôde ser obtido: o motivo (o job avisa). */
        readonly mediaUnavailable?: string
        readonly confirmBeforeTranslate: boolean
        /** Ausente em jobs antigos = `use`. Use `jobMemoryMode`. */
        readonly memory?: JobMemoryMode
      })
  | (Base &
      MuxSpecOptions & { readonly kind: 'mux'; readonly subtitles: readonly MuxSubtitleSpec[] })
  | (Base &
      MuxSpecOptions &
      WriteSpecOptions & {
        readonly kind: 'pipeline'
        readonly trackId?: number
        readonly externalSubtitle?: ExternalSubtitleSpec
        readonly preset: PresetSnapshot
        readonly targetLang: string
        readonly sourceLang?: string
        /** Jobs novos sempre trazem; jobs persistidos antigos só têm `alsoWriteSubtitle`. */
        readonly outputMode?: OutputMode
        /** @deprecated Legado: `true` = `both`. Use `pipelineOutputMode`. */
        readonly alsoWriteSubtitle?: boolean
        readonly notes?: string
        readonly media?: MediaContextLookup
        readonly mediaUnavailable?: string
        readonly confirmBeforeTranslate: boolean
        readonly memory?: JobMemoryMode
      })

/** `memory` efetivo de um spec persistido (`undefined` = não traduz). */
export function jobMemoryMode(spec: unknown): JobMemoryMode | undefined {
  const s = spec as { kind?: string; memory?: JobMemoryMode } | null
  if (s?.kind !== 'translate' && s?.kind !== 'pipeline') return undefined
  return s.memory === 'skip' ? 'skip' : 'use'
}

export const jobMemoryModeSchema = z.enum(['use', 'skip']) satisfies z.ZodType<JobMemoryMode>

const path = z.string().min(1).max(4096)
const lang = z.string().min(1).max(64)
const trackId = z.number().int().min(0).max(10_000)
const episodeNumber = z.number().int().min(0).max(9999)
const mediaRequest = z.object({
  ref: z.object({
    provider: z.literal('tmdb'),
    id: z.number().int().positive(),
    kind: z.enum(['movie', 'series'])
  }),
  season: episodeNumber.optional(),
  episode: episodeNumber.optional()
}) satisfies z.ZodType<JobMediaRequest>
const muxOptions = z
  .object({ setDefault: z.boolean().optional(), keepOriginalSubtitles: z.boolean().optional() })
  .optional()

export const JobCreateRequestSchema = z.discriminatedUnion('kind', [
  z.object({
    kind: z.literal('extract'),
    source: path,
    trackIds: z.array(trackId).min(1).max(64),
    outputDir: path.optional(),
    options: z
      .object({ overwrite: z.boolean().optional(), naming: z.enum(['player', 'track']).optional() })
      .optional()
  }),
  z.object({
    kind: z.literal('translate'),
    source: path,
    presetId: z.string().min(1).max(64).optional(),
    model: z.string().trim().min(1).max(200).optional(),
    notes: z.string().max(MAX_USER_NOTES).optional(),
    media: mediaRequest.optional(),
    targetLang: lang.optional(),
    sourceLang: lang.optional(),
    outputDir: path.optional(),
    options: z.object({ overwrite: z.boolean().optional() }).optional(),
    confirmBeforeTranslate: z.boolean().optional(),
    memory: jobMemoryModeSchema.optional()
  }),
  z.object({
    kind: z.literal('mux'),
    source: path,
    subtitles: z
      .array(
        z.object({
          path,
          language: lang,
          name: z.string().trim().min(1).max(200).optional(),
          default: z.boolean().optional()
        })
      )
      .min(1)
      .max(32),
    outputDir: path.optional(),
    options: muxOptions
  }),
  z.object({
    kind: z.literal('pipeline'),
    source: path,
    trackId: trackId.optional(),
    subtitle: z.object({ path, language: lang.optional() }).optional(),
    presetId: z.string().min(1).max(64).optional(),
    model: z.string().trim().min(1).max(200).optional(),
    notes: z.string().max(MAX_USER_NOTES).optional(),
    media: mediaRequest.optional(),
    targetLang: lang.optional(),
    sourceLang: lang.optional(),
    outputDir: path.optional(),
    options: z
      .object({
        setDefault: z.boolean().optional(),
        keepOriginalSubtitles: z.boolean().optional(),
        outputMode: z.enum(['mux', 'subtitle', 'both']).optional(),
        alsoWriteSubtitle: z.boolean().optional(),
        includeSourceSubtitle: z.boolean().optional(),
        overwrite: z.boolean().optional()
      })
      .optional(),
    confirmBeforeTranslate: z.boolean().optional(),
    memory: jobMemoryModeSchema.optional()
  })
]) satisfies z.ZodType<JobCreateRequest>

function requiredLang(tag: string): string {
  const t = canonicalTag(tag)
  if (!t) throw new AppError('INVALID_INPUT', `Invalid language tag: ${tag}`)
  return t
}
const optionalLang = (tag: string | undefined): string | undefined =>
  tag === undefined ? undefined : requiredLang(tag)

/** Preset (estilo) + modelo e idioma: o pedido manda; senão, os padrões das configurações. */
async function resolvePreset(
  req: { presetId?: string; model?: string; targetLang?: string },
  presets: PresetSource
): Promise<{ preset: PresetSnapshot; targetLang: string }> {
  const defaults = await presets.defaults()
  const p = await presets.get(req.presetId ?? defaults.presetId)
  const model = req.model ?? defaults.model
  if (!model) throw new AppError('MODEL_NOT_SET', 'No model in the request and no default model')
  const lang = req.targetLang ?? defaults.targetLang
  if (!lang)
    throw new AppError('TARGET_LANG_NOT_SET', 'No target language in the request and no default')
  return {
    preset: {
      id: p.id,
      name: p.name,
      mediaType: p.mediaType,
      model,
      glossary: p.glossary,
      temperature: p.temperature,
      prompt: p.prompt
    },
    targetLang: requiredLang(lang)
  }
}

const withNotes = (notes: string | undefined): { notes?: string } =>
  notes?.trim() ? { notes: notes.trim() } : {}

/**
 * Pedido do renderer -> spec congelado. Revalida todos os caminhos no Main (o renderer nunca é
 * fonte confiável), copia o preset e aplica os padrões. Lança `INVALID_INPUT`/`INVALID_PATH`.
 */
export async function resolveJobSpec(
  raw: JobCreateRequest,
  presets: PresetSource,
  media?: MediaLookupSource
): Promise<JobSpec> {
  const parsed = JobCreateRequestSchema.safeParse(raw)
  if (!parsed.success) throw new AppError('INVALID_INPUT', parsed.error.message)
  const req = parsed.data

  const sourceExts = req.kind === 'translate' ? SUBTITLE_EXTENSIONS : MEDIA_EXTENSIONS
  const src = await validateInputFile(
    req.source,
    sourceExts,
    req.kind === 'translate' ? MAX_SUBTITLE_BYTES : undefined
  )
  const outputDir = await validateOutputDir(req.outputDir ?? dirname(src.path))
  const base = { source: src.path, sourceSize: src.size, outputDir }

  switch (req.kind) {
    case 'extract': {
      const trackIds = [...new Set(req.trackIds)].sort((a, b) => a - b)
      return {
        ...base,
        kind: 'extract',
        trackIds,
        naming: req.options?.naming ?? (trackIds.length === 1 ? 'player' : 'track'),
        ...overwrite(req.options?.overwrite)
      }
    }
    case 'translate': {
      const { preset, targetLang } = await resolvePreset(req, presets)
      return {
        ...base,
        kind: 'translate',
        preset,
        targetLang,
        ...withNotes(req.notes),
        ...(await resolveMedia(req.media, media)),
        ...withLang(optionalLang(req.sourceLang) ?? langFromFilename(src.path)),
        ...overwrite(req.options?.overwrite),
        confirmBeforeTranslate: req.confirmBeforeTranslate ?? true,
        memory: req.memory ?? 'use'
      }
    }
    case 'mux': {
      const subtitles: MuxSubtitleSpec[] = []
      for (const s of req.subtitles) {
        const f = await validateInputFile(s.path, SUBTITLE_EXTENSIONS, MAX_SUBTITLE_BYTES)
        // `und` (idioma desconhecido) é válido numa faixa de legenda; o padrão do batch mux.
        const language = s.language.toLowerCase() === 'und' ? 'und' : requiredLang(s.language)
        subtitles.push({
          path: f.path,
          size: f.size,
          language,
          name: s.name ?? languageName(language),
          default: s.default ?? ((req.options?.setDefault ?? true) && subtitles.length === 0)
        })
      }
      return {
        ...base,
        kind: 'mux',
        subtitles,
        setDefault: req.options?.setDefault ?? true,
        keepOriginalSubtitles: req.options?.keepOriginalSubtitles ?? true
      }
    }
    case 'pipeline': {
      const origin = await pipelineOrigin(req)
      const { preset, targetLang } = await resolvePreset(req, presets)
      return {
        ...base,
        kind: 'pipeline',
        ...origin.spec,
        preset,
        targetLang,
        ...withNotes(req.notes),
        ...(await resolveMedia(req.media, media)),
        ...withLang(optionalLang(req.sourceLang) ?? origin.sourceLang),
        setDefault: req.options?.setDefault ?? true,
        keepOriginalSubtitles: req.options?.keepOriginalSubtitles ?? true,
        // `alsoWriteSubtitle` (legado) só vale sem `outputMode`: `true` = `both`.
        outputMode: req.options?.outputMode ?? (req.options?.alsoWriteSubtitle ? 'both' : 'mux'),
        ...overwrite(req.options?.overwrite),
        confirmBeforeTranslate: req.confirmBeforeTranslate ?? true,
        memory: req.memory ?? 'use'
      }
    }
  }
}

async function pipelineOrigin(
  req: Extract<z.infer<typeof JobCreateRequestSchema>, { kind: 'pipeline' }>
): Promise<{
  spec: { trackId: number } | { externalSubtitle: ExternalSubtitleSpec }
  sourceLang?: string
}> {
  // Exatamente uma origem: o retorno é a união discriminada, então o invariante
  // é garantido pelo narrowing aqui (sem `!`) e checado antes de tocar no disco.
  const { trackId, subtitle } = req
  if (subtitle === undefined) {
    if (trackId === undefined) throw exactlyOneOrigin()
    return { spec: { trackId } }
  }
  if (trackId !== undefined) throw exactlyOneOrigin()
  const f = await validateInputFile(subtitle.path, SUBTITLE_EXTENSIONS, MAX_SUBTITLE_BYTES)
  const mode = req.options?.outputMode ?? (req.options?.alsoWriteSubtitle ? 'both' : 'mux')
  return {
    spec: {
      externalSubtitle: {
        path: f.path,
        size: f.size,
        format: formatFromPath(f.path),
        includeInMux: mode !== 'subtitle' && (req.options?.includeSourceSubtitle ?? false)
      }
    },
    sourceLang: optionalLang(subtitle.language) ?? langFromFilename(f.path)
  }
}

const exactlyOneOrigin = (): AppError =>
  new AppError('INVALID_INPUT', 'Exactly one of trackId and subtitle is required')

const overwrite = (v: boolean | undefined): { overwrite?: true } => (v ? { overwrite: true } : {})

const withLang = (sourceLang: string | undefined): { sourceLang?: string } =>
  sourceLang ? { sourceLang } : {}
