import { join } from 'node:path'
import { AppError } from '../../errors'
import { ensureFreeSpace } from '../../media/path-guard'
import { fileArtifact } from '../artifacts'
import type { Step } from '../step-context'

const EXTRACT_FREE_BYTES = 100 * 1024 * 1024

/** Uma única passada de `mkvextract` para todas as faixas, dentro do diretório de trabalho. */
export const extractStep: Step = async (ctx) => {
  const { spec } = ctx
  const ids =
    spec.kind === 'extract'
      ? spec.trackIds
      : spec.kind === 'pipeline' && spec.trackId !== undefined
        ? [spec.trackId]
        : []
  const tracks = ids.map((id) => {
    const t = ctx.artifacts.probe?.tracks.find((x) => x.id === id)
    if (!t?.format) throw new AppError('INVALID_INPUT', `Track ${id} is not available`)
    return {
      id,
      format: t.format,
      language: t.language,
      output: join(ctx.workDir, `track-${id}.${t.format}`)
    }
  })
  await ensureFreeSpace(ctx.workDir, EXTRACT_FREE_BYTES)
  const { warnings } = await ctx.deps.mkv.extract(
    spec.source,
    tracks.map((t) => ({ trackId: t.id, output: t.output })),
    { signal: ctx.signal, onProgress: (f) => ctx.progress({ stepFraction: f }) }
  )
  warnings.forEach((w) => ctx.warn(w))
  const files: NonNullable<typeof ctx.artifacts.extract>['files'] = []
  for (const t of tracks) {
    const a = await fileArtifact(t.output).catch(() => {
      throw new AppError('MEDIA_TOOL_FAILED', `mkvextract did not produce track ${t.id}`)
    })
    files.push({
      ...a,
      trackId: t.id,
      format: t.format,
      ...(t.language ? { language: t.language } : {})
    })
  }
  ctx.artifacts.extract = { files }
}
