import { z } from 'zod'
import type { SubtitleTextFormat } from '../../shared/ipc-contract'
import { canonicalTag, toBcp47 } from './languages'

export type SubtitleFormatId = SubtitleTextFormat

export interface SubtitleTrack {
  readonly id: number
  readonly codecId: string
  /** `null` quando a faixa não pode ser traduzida (ver `unsupported`). */
  readonly format: SubtitleFormatId | null
  /** `image-based`: PGS/VobSub/DVB (exigiria OCR, fora do escopo). `no-parser`: texto sem parser. */
  readonly unsupported?: 'image-based' | 'no-parser'
  readonly language?: string
  readonly name?: string
  readonly default: boolean
  readonly forced: boolean
  readonly eventCount?: number
}

export interface MkvInfo {
  readonly subtitles: readonly SubtitleTrack[]
  readonly attachments: number
  /** fps da primeira faixa de vídeo (`default_duration`); ausente quando o container não informa. */
  readonly videoFps?: number
}

const CODECS: Readonly<Record<string, SubtitleFormatId>> = {
  'S_TEXT/UTF8': 'srt',
  'S_TEXT/ASCII': 'srt',
  'S_TEXT/ASS': 'ass',
  'S_TEXT/SSA': 'ssa',
  'S_TEXT/WEBVTT': 'vtt'
}
const IMAGE_CODECS = new Set(['S_HDMV/PGS', 'S_VOBSUB', 'S_DVBSUB', 'S_IMAGE/BMP'])

const Identify = z.looseObject({
  container: z.looseObject({ recognized: z.boolean().optional() }).optional(),
  tracks: z
    .array(
      z.looseObject({
        id: z.number(),
        type: z.string(),
        properties: z
          .looseObject({
            codec_id: z.string().optional(),
            language: z.string().optional(),
            language_ietf: z.string().optional(),
            track_name: z.string().optional(),
            default_track: z.boolean().optional(),
            forced_track: z.boolean().optional(),
            num_index_entries: z.number().optional(),
            /** Duração de um quadro, em nanossegundos (faixas de vídeo). */
            default_duration: z.number().optional()
          })
          .optional()
      })
    )
    .default([]),
  attachments: z.array(z.unknown()).optional()
})

/** `mkvmerge -J` -> faixas de legenda. Lança para JSON fora da forma ou container não reconhecido. */
export function parseIdentify(json: unknown): MkvInfo {
  const id = Identify.parse(json)
  if (id.container?.recognized === false) throw new Error('Container not recognized')
  const subtitles = id.tracks
    .filter((t) => t.type === 'subtitles')
    .map((t): SubtitleTrack => {
      const p = t.properties ?? {}
      const codecId = p.codec_id ?? ''
      const format = CODECS[codecId] ?? null
      const language = canonicalTag(p.language_ietf) ?? toBcp47(p.language)
      return {
        id: t.id,
        codecId,
        format,
        ...(format === null
          ? { unsupported: IMAGE_CODECS.has(codecId) ? 'image-based' : 'no-parser' }
          : {}),
        language,
        ...(p.track_name ? { name: p.track_name } : {}),
        default: p.default_track ?? false,
        forced: p.forced_track ?? false,
        ...(p.num_index_entries !== undefined ? { eventCount: p.num_index_entries } : {})
      } as SubtitleTrack
    })
  const videoFps = frameRate(
    id.tracks.find((t) => t.type === 'video')?.properties?.default_duration
  )
  return {
    subtitles,
    attachments: id.attachments?.length ?? 0,
    ...(videoFps !== undefined ? { videoFps } : {})
  }
}

/** ns por quadro -> fps com 3 casas (41708333 -> 23.976); `undefined` se ausente ou absurdo. */
function frameRate(defaultDuration: number | undefined): number | undefined {
  if (defaultDuration === undefined || !(defaultDuration > 0)) return undefined
  const fps = Math.round((1e9 / defaultDuration) * 1000) / 1000
  return fps > 0 && fps <= 1000 ? fps : undefined
}
