import { rename, unlink } from 'node:fs/promises'
import { AppError } from '../../errors'
import { languageName } from '../../llm/prompt-builder'
import { outputNames, uniquePath } from '../../media/naming'
import { ensureFreeSpace } from '../../media/path-guard'
import type { MuxSubtitle } from '../../media/mkvtoolnix'
import type { Step } from '../step-context'

const SPACE_MARGIN = 1.05

/** MKV novo ao lado (nunca sobrescreve a origem), gravado em `.partial` e renomeado no fim. */
export const muxStep: Step = async (ctx) => {
  const { spec } = ctx
  let subtitles: MuxSubtitle[]
  let subtitleBytes: number
  let lang: string
  if (spec.kind === 'pipeline') {
    const t = ctx.artifacts.translate
    if (!t) throw new AppError('INVALID_INPUT', 'Translated file is missing')
    lang = spec.targetLang
    subtitles = [
      {
        path: t.file.path,
        language: lang,
        name: `${languageName(lang)} (BakaSub)`,
        default: spec.setDefault
      }
    ]
    subtitleBytes = t.file.size
    const ext = spec.externalSubtitle
    if (ext?.includeInMux) {
      const srcLang = spec.sourceLang ?? 'und'
      subtitles.push({
        path: ext.path,
        language: srcLang,
        name: srcLang === 'und' ? 'Original' : languageName(srcLang),
        default: false
      })
      subtitleBytes += ext.size
    }
  } else if (spec.kind === 'mux') {
    lang = spec.subtitles[0].language
    subtitles = spec.subtitles.map((s) => ({
      path: s.path,
      language: s.language,
      name: s.name,
      default: s.default
    }))
    subtitleBytes = spec.subtitles.reduce((n, s) => n + s.size, 0)
  } else throw new AppError('INVALID_INPUT', `mux step in a ${spec.kind} job`)

  const final = await uniquePath(spec.outputDir, outputNames.muxed(spec.source, lang))
  const partial = outputNames.partial(final)
  await ensureFreeSpace(spec.outputDir, Math.ceil((spec.sourceSize + subtitleBytes) * SPACE_MARGIN))

  try {
    const { warnings } = await ctx.deps.mkv.mux(
      {
        source: spec.source,
        output: partial,
        subtitles,
        sourceSubtitleTrackIds: (ctx.artifacts.probe?.tracks ?? []).map((t) => t.id),
        setDefault: spec.setDefault,
        keepOriginalSubtitles: spec.keepOriginalSubtitles
      },
      { signal: ctx.signal, onProgress: (f) => ctx.progress({ stepFraction: f }) }
    )
    await rename(partial, final)
    warnings.forEach((w) => ctx.warn(w))
    ctx.artifacts.mux = { output: final, warnings }
  } catch (err) {
    await unlink(partial).catch(() => undefined)
    throw err
  }
}
