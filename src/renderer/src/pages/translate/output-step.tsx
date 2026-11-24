import { useTranslation } from 'react-i18next'
import type { OutputMode } from '@shared/ipc-contract'
import { Checkbox } from '@renderer/components/ui/checkbox'
import { OutputFolder } from '@renderer/components/files/output-folder'
import { NamePreview, OverwriteOption, RadioChoice } from '@renderer/components/output/output-parts'

export interface MuxChoice {
  setDefault: boolean
  keepOriginalSubtitles: boolean
}

export interface OutputPreview {
  names: string[]
  folder: string
  /** A batch previews its first file only. */
  example?: boolean
}

/**
 * Output folder; for videos also what to write (video, subtitle only or both) with the mux
 * options that only apply when a video is written. Hidden for subtitle files.
 */
export function OutputStep({
  video,
  outputDir,
  mode,
  mux,
  overwrite,
  includeOriginal,
  preview,
  onFolder,
  onMode,
  onMux,
  onOverwrite,
  onIncludeOriginal,
  summary
}: {
  video: boolean
  outputDir: string | null
  mode: OutputMode
  mux: MuxChoice
  overwrite: boolean
  includeOriginal?: boolean
  preview?: OutputPreview | null
  onFolder: (dir: string | null) => void
  onMode: (mode: OutputMode) => void
  onMux: (patch: Partial<MuxChoice>) => void
  onOverwrite: (v: boolean) => void
  onIncludeOriginal?: (v: boolean) => void
  /** Replaces the closing note (the batch has its own). */
  summary?: React.ReactNode
}): React.JSX.Element {
  const { t } = useTranslation()
  const writesVideo = video && mode !== 'subtitle'
  const writesSubtitle = !video || mode !== 'mux'
  const options: { key: keyof MuxChoice; label: string }[] = [
    { key: 'setDefault', label: t('translate.output.setDefault') },
    { key: 'keepOriginalSubtitles', label: t('translate.output.keepOriginal') }
  ]
  return (
    <div className="grid gap-4">
      {video ? (
        <>
          <RadioChoice
            legend={t('translate.output.mode.label')}
            name="output-mode"
            value={mode}
            onChange={onMode}
            options={(['mux', 'subtitle', 'both'] as const).map((m) => ({
              value: m,
              label: t(`translate.output.mode.${m}`),
              hint: t(`translate.output.mode.${m}Hint`)
            }))}
          />
          {writesSubtitle ? (
            <p className="-mt-2 text-xs text-muted-foreground">
              {t('translate.output.playersNote')}
            </p>
          ) : null}
        </>
      ) : null}
      <div className="grid gap-1.5">
        <span className="text-sm font-medium">{t('translate.output.folder')}</span>
        <OutputFolder value={outputDir} onChange={onFolder} />
      </div>
      {writesVideo ? (
        <div className="grid gap-2">
          {options.map((o) => (
            <label key={o.key} className="flex items-center gap-2.5 text-sm">
              <Checkbox
                checked={mux[o.key]}
                onCheckedChange={(v) => onMux({ [o.key]: v === true })}
              />
              {o.label}
            </label>
          ))}
          {includeOriginal !== undefined ? (
            <label className="flex items-center gap-2.5 text-sm">
              <Checkbox
                checked={includeOriginal}
                onCheckedChange={(v) => onIncludeOriginal?.(v === true)}
              />
              {t('translate.output.includeOriginal')}
            </label>
          ) : null}
        </div>
      ) : null}
      {preview ? <NamePreview {...preview} /> : null}
      <OverwriteOption value={overwrite} onChange={onOverwrite} />
      {summary ?? <p className="text-xs text-muted-foreground">{t('translate.output.summary')}</p>}
    </div>
  )
}
