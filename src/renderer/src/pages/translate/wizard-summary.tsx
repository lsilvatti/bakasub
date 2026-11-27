import { useTranslation } from 'react-i18next'
import type { MediaContextLookup, Preset, SubtitleTrackInfo } from '@shared/ipc-contract'
import { useModels } from '@renderer/data/models'
import { episodeCode } from '@renderer/lib/episode'
import type { ExistingTarget } from '@renderer/lib/existing-subtitle'
import { languageName } from '@renderer/lib/languages'
import { presetLabel } from '@renderer/lib/preset-labels'
import { SummaryList } from './summary-list'
import type { SummaryRow } from './summary-rows'
import { stepsFor, type WizardState, type WizardStep } from './wizard-state'

/** What was chosen so far: a row per item of a reached step; a row takes the user to its step. */
export function WizardSummary({
  state,
  preset,
  model,
  targetLang,
  track,
  syncWarn,
  media,
  existing,
  outputNames,
  canGo,
  onGo
}: {
  state: WizardState
  preset: Preset | undefined
  model: string | null
  targetLang: string | null
  track: SubtitleTrackInfo | undefined
  /** External subtitle that may be out of sync with the video. */
  syncWarn?: string
  media: MediaContextLookup | undefined
  existing: ExistingTarget | null
  /** Files the job will write (the subtitle path, with the output folder). */
  outputNames: readonly string[]
  canGo: (s: WizardStep) => boolean
  onGo: (s: WizardStep) => void
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const catalog = useModels().data?.models
  const lang = (tag: string): string => languageName(tag, i18n.language)

  const rows: SummaryRow<WizardStep>[] = [
    { step: 'file', label: t('translate.summary.file'), value: state.file?.name ?? null },
    {
      step: 'track',
      label: t('translate.summary.track'),
      value: state.subtitle
        ? `${state.subtitle.name} · ${t(`translate.external.origin.${state.subtitle.origin}`)}`
        : track
          ? [
              track.language ? lang(track.language) : t('translate.track.unknownLang'),
              track.format?.toUpperCase()
            ]
              .filter(Boolean)
              .join(' · ')
          : null,
      ...(state.subtitle && syncWarn ? { warn: syncWarn } : {})
    },
    {
      step: 'preset',
      label: t('translate.summary.preset'),
      value: preset ? presetLabel(preset, t) : null
    },
    {
      step: 'preset',
      label: t('translate.summary.model'),
      value: model ? (catalog?.find((m) => m.id === model)?.name ?? model) : null
    },
    {
      step: 'preset',
      label: t('translate.summary.language'),
      value: targetLang ? lang(targetLang) : null,
      ...(existing && targetLang
        ? { warn: t('translate.existing.inline', { language: lang(targetLang) }) }
        : {})
    },
    {
      step: 'preset',
      label: t('translate.summary.media'),
      // Only a work the user has (or turned off); nothing while the details load.
      value: !state.media.ref
        ? null
        : !state.media.include
          ? t('translate.summary.mediaOff')
          : media
            ? `${media.details.title}${media.episode ? ` · ${episodeCode(media.episode.season, media.episode.episode)}` : ''}`
            : null
    },
    {
      step: 'output',
      label: t('translate.summary.output'),
      value: state.outputDir ?? t('translate.output.nextToSource'),
      detail: outputNames.join(' · ')
    }
  ]

  return (
    <SummaryList
      id="wizard-summary"
      title={t('translate.summary.title')}
      empty={t('translate.summary.empty')}
      rows={rows}
      steps={stepsFor(state.file?.kind)}
      furthest={state.furthest}
      current={state.step}
      stepLabel={(s) => t(`translate.steps.${s}`)}
      canGo={canGo}
      onGo={onGo}
    />
  )
}
