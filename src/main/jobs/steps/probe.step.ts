import { AppError } from '../../errors'
import type { Step } from '../step-context'

/** `mkvmerge -J`; confere que as faixas pedidas existem e são traduzíveis/extraíveis. */
export const probeStep: Step = async (ctx) => {
  const info = await ctx.deps.mkv.probe(ctx.spec.source, ctx.signal)
  const wanted =
    ctx.spec.kind === 'extract'
      ? ctx.spec.trackIds
      : ctx.spec.kind === 'pipeline' && ctx.spec.trackId !== undefined
        ? [ctx.spec.trackId]
        : []
  for (const id of wanted) {
    const t = info.subtitles.find((s) => s.id === id)
    if (!t) throw new AppError('INVALID_INPUT', `No subtitle track with id ${id}`)
    if (!t.format)
      throw new AppError('INVALID_INPUT', `Track ${id} is not supported (${t.unsupported})`)
  }
  ctx.artifacts.probe = { tracks: [...info.subtitles] }
}
