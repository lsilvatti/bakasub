import type { FreeQuota } from '../../llm/free-quota'
import { BrowserWindow, dialog, type OpenDialogOptions } from 'electron'
import { z } from 'zod'
import type { ModelList, OpenDialogKind } from '../../../shared/ipc-contract'
import type { ModelCatalog } from '../../llm/model-catalog'
import type { MkvToolnix } from '../../media/mkvtoolnix-service'
import {
  MAX_SUBTITLE_BYTES,
  MEDIA_EXTENSIONS,
  SUBTITLE_EXTENSIONS,
  validateInputFile,
  validateOutputDir
} from '../../media/path-guard'
import { createHandler } from '../create-handler'

const EXTENSIONS: Readonly<Record<OpenDialogKind, readonly string[]>> = {
  media: MEDIA_EXTENSIONS,
  subtitle: SUBTITLE_EXTENSIONS,
  any: [...MEDIA_EXTENSIONS, ...SUBTITLE_EXTENSIONS]
}

const FILTERS: Readonly<Record<OpenDialogKind, OpenDialogOptions['filters']>> = {
  media: [{ name: 'Matroska', extensions: ['mkv'] }],
  subtitle: [{ name: 'Subtitles', extensions: SUBTITLE_EXTENSIONS.map((e) => e.slice(1)) }],
  any: [
    {
      name: 'Video and subtitles',
      extensions: [...MEDIA_EXTENSIONS, ...SUBTITLE_EXTENSIONS].map((e) => e.slice(1))
    }
  ]
}

export function registerMediaHandlers(deps: {
  mkv: Pick<MkvToolnix, 'probe'>
  models: ModelCatalog
  freeQuota: Pick<FreeQuota, 'status'>
}): void {
  createHandler('media:probe', {
    input: z.object({ path: z.string().min(1).max(4096) }),
    handle: async ({ path }) => {
      const file = await validateInputFile(path, MEDIA_EXTENSIONS)
      const info = await deps.mkv.probe(file.path)
      return {
        subtitles: [...info.subtitles],
        attachments: info.attachments,
        ...(info.videoFps !== undefined ? { videoFps: info.videoFps } : {})
      }
    }
  })

  createHandler('models:free-quota', { handle: () => deps.freeQuota.status() })

  createHandler('models:list', {
    input: z.object({ force: z.boolean().optional() }),
    handle: async ({ force }): Promise<ModelList> => {
      const snap = await deps.models.get({ force })
      return {
        models: [...snap.models],
        fetchedAt: new Date(snap.fetchedAt).toISOString(),
        offline: snap.offline
      }
    }
  })

  createHandler('dialog:open-files', {
    input: z.object({
      kind: z.enum(['media', 'subtitle', 'any']),
      multiple: z.boolean().optional()
    }),
    handle: async ({ kind, multiple }, event) => {
      const win = BrowserWindow.fromWebContents(event.sender)
      const options: OpenDialogOptions = {
        properties: ['openFile', ...(multiple ? (['multiSelections'] as const) : [])],
        filters: FILTERS[kind]
      }
      const res = win
        ? await dialog.showOpenDialog(win, options)
        : await dialog.showOpenDialog(options)
      if (res.canceled) return { paths: [] }
      const paths: string[] = []
      for (const p of res.filePaths) {
        const isSub = SUBTITLE_EXTENSIONS.some((e) => p.toLowerCase().endsWith(e))
        const ok = await validateInputFile(
          p,
          EXTENSIONS[kind],
          isSub ? MAX_SUBTITLE_BYTES : undefined
        ).catch(() => null)
        if (ok) paths.push(ok.path)
      }
      return { paths }
    }
  })

  createHandler('dialog:choose-output-dir', {
    handle: async (_req, event) => {
      const win = BrowserWindow.fromWebContents(event.sender)
      const options: OpenDialogOptions = { properties: ['openDirectory', 'createDirectory'] }
      const res = win
        ? await dialog.showOpenDialog(win, options)
        : await dialog.showOpenDialog(options)
      if (res.canceled || !res.filePaths[0]) return { path: null }
      return { path: await validateOutputDir(res.filePaths[0]) }
    }
  })
}
