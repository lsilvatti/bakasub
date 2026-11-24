import type { JobMemoryMode, PromptSummary } from '../../../shared/ipc-contract'
import { AppError } from '../../errors'
import { buildMediaContext } from '../../media-context/media-block'
import { mediaScope } from '../../memory/memory-key'
import type { SubtitleDocument, SubtitleFormat } from '../../parsers'
import type { PresetSnapshot } from '../job-spec'
import type { StepContext } from '../step-context'
import { formatFromPath, readSubtitle, type TextFormatId } from './subtitle-io'

export interface TranslationInput {
  readonly doc: SubtitleDocument<unknown>
  readonly format: SubtitleFormat<unknown>
  readonly formatId: TextFormatId
  readonly preset: PresetSnapshot
  /** Bloco `<media>`: detalhes do TMDB congelados no job + notas do usuário. */
  readonly media?: string
  /** De onde veio o `<media>` (e por que o TMDB pedido ficou de fora), para o resumo do prompt. */
  readonly mediaSources: PromptSummary['media']
  /** Obra (TMDB) congelada no job, para a chave do Job Memory; `''` sem obra. */
  readonly mediaScope: string
  /** `skip`: ignorar as traduções salvas (tudo vai ao LLM; as novas sobrescrevem). */
  readonly memory: JobMemoryMode
  readonly targetLang: string
  readonly sourceLang?: string
}

/**
 * O documento a traduzir: o arquivo de origem (`translate`) ou a faixa extraída (`pipeline`).
 * Compartilhado por `preflight` e `translate`, para a estimativa valer exatamente para o que
 * será enviado.
 */
export async function translationInput(ctx: StepContext): Promise<TranslationInput> {
  const { spec } = ctx
  if (spec.kind !== 'translate' && spec.kind !== 'pipeline')
    throw new AppError('INVALID_INPUT', `translation step in a ${spec.kind} job`)
  let path: string
  let formatId: TextFormatId
  let sourceLang = spec.sourceLang
  if (spec.kind === 'translate') {
    path = spec.source
    formatId = formatFromPath(spec.source)
  } else if (spec.externalSubtitle) {
    path = spec.externalSubtitle.path
    formatId = spec.externalSubtitle.format
  } else {
    const f = ctx.artifacts.extract?.files.find((x) => x.trackId === spec.trackId)
    if (!f) throw new AppError('INVALID_INPUT', 'Extracted track is missing')
    path = f.path
    formatId = f.format
    sourceLang ??= f.language
  }
  const { doc, format } = await readSubtitle(path, formatId)
  return {
    doc,
    format,
    formatId,
    preset: spec.preset,
    ...withMedia(buildMediaContext(spec.media, spec.notes)),
    mediaSources: {
      tmdb: !!spec.media,
      episode: !!spec.media?.episode,
      notes: !!spec.notes?.trim(),
      ...(spec.mediaUnavailable ? { unavailable: spec.mediaUnavailable } : {})
    },
    mediaScope: mediaScope(spec.media?.details),
    memory: spec.memory === 'skip' ? 'skip' : 'use',
    targetLang: spec.targetLang,
    ...(sourceLang ? { sourceLang } : {})
  }
}

const withMedia = (media: string | undefined): { media?: string } => (media ? { media } : {})
