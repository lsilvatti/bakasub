import { useTranslation } from 'react-i18next'
import { subtitleSyncRisks, type SyncRisk } from '@shared/subtitle-sync'

/** What OpenSubtitles said about a result (absent for a file the user picked themselves). */
export interface SubtitleMatch {
  hashMatch: boolean
  fps?: number
}

export interface SyncCheck {
  risks: SyncRisk[]
  subtitleFps?: number
  videoFps?: number
}

/** Risks of an OpenSubtitles subtitle against the video; `null` = nothing to say. */
export function syncCheck(
  match: SubtitleMatch | undefined,
  videoFps: number | undefined
): SyncCheck | null {
  if (!match) return null
  const risks = subtitleSyncRisks({ hashMatch: match.hashMatch, subtitleFps: match.fps, videoFps })
  return risks.length
    ? {
        risks,
        ...(match.fps !== undefined ? { subtitleFps: match.fps } : {}),
        ...(videoFps !== undefined ? { videoFps } : {})
      }
    : null
}

export const formatFps = (fps: number | undefined, lang: string): string =>
  fps === undefined ? '?' : new Intl.NumberFormat(lang, { maximumFractionDigits: 3 }).format(fps)

/** One short line per risk: `sync.<risk>` with the two frame rates. */
export function useSyncMessages(): (check: SyncCheck, form: 'short' | 'long') => string[] {
  const { t, i18n } = useTranslation()
  return (check, form) =>
    check.risks.map((risk) =>
      t(`opensubtitles.sync.${form}.${risk}`, {
        subtitle: formatFps(check.subtitleFps, i18n.language),
        video: formatFps(check.videoFps, i18n.language)
      })
    )
}
