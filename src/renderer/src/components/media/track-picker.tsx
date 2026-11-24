import { useEffect } from 'react'
import { Check } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@renderer/components/ui/badge'
import { useProbe } from '@renderer/data/media'
import { useSettings } from '@renderer/data/settings'
import { errorMessage } from '@renderer/i18n/error-message'
import { languageName } from '@renderer/lib/languages'
import { pickTrack } from '@renderer/lib/tracks'
import { cn } from '@renderer/lib/utils'

/**
 * Subtitle tracks of a video (`media:probe`). `single` (translate) or `multi` (extract).
 * Unsupported tracks are listed but disabled, with the reason. One track is preselected.
 */
export function TrackPicker({
  path,
  mode,
  value,
  onChange,
  autoSelect = true
}: {
  path: string
  mode: 'single' | 'multi'
  value: readonly number[]
  onChange: (ids: number[]) => void
  autoSelect?: boolean
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const probe = useProbe(path)
  const preferred = useSettings().data?.['tracks.preferredSourceLangs']

  useEffect(() => {
    if (!autoSelect || value.length > 0 || !probe.data) return
    const id = pickTrack(probe.data.subtitles, preferred ?? ['en'])
    if (id !== null) onChange([id])
  }, [autoSelect, probe.data, preferred, value.length, onChange])

  if (probe.isPending)
    return <p className="text-sm text-muted-foreground">{t('translate.track.loading')}</p>
  if (probe.isError)
    return (
      <p role="alert" className="text-sm text-destructive">
        {errorMessage(probe.error)}
      </p>
    )
  const tracks = probe.data.subtitles
  if (!tracks.some((tr) => tr.format !== null))
    return <p className="text-sm text-warning">{t('translate.track.none')}</p>

  const toggle = (id: number): void =>
    onChange(
      mode === 'single' ? [id] : value.includes(id) ? value.filter((x) => x !== id) : [...value, id]
    )

  return (
    <div className="grid gap-2">
      <div
        role={mode === 'single' ? 'radiogroup' : 'group'}
        aria-label={t('translate.track.heading')}
        className="grid gap-1.5"
      >
        {tracks.map((tr) => {
          const disabled = tr.format === null
          const selected = value.includes(tr.id)
          return (
            <button
              key={tr.id}
              type="button"
              role={mode === 'single' ? 'radio' : 'checkbox'}
              aria-checked={selected}
              disabled={disabled}
              onClick={() => toggle(tr.id)}
              className={cn(
                'flex flex-wrap items-center gap-2 rounded-lg border px-3 py-2 text-left text-sm disabled:opacity-45',
                selected ? 'border-primary bg-primary-soft/40' : 'bg-background/40 hover:bg-accent'
              )}
            >
              {mode === 'multi' ? (
                <span
                  aria-hidden
                  className={cn(
                    'grid size-4 place-items-center rounded-[4px] border',
                    selected && 'border-primary bg-primary text-primary-foreground'
                  )}
                >
                  {selected ? <Check className="size-3" strokeWidth={3} /> : null}
                </span>
              ) : null}
              <span className="font-medium">
                {tr.language
                  ? languageName(tr.language, i18n.language)
                  : t('translate.track.unknownLang')}
              </span>
              {tr.name ? <span className="text-muted-foreground">{tr.name}</span> : null}
              <span className="text-xs text-muted-foreground uppercase">
                {tr.format ?? tr.codecId}
              </span>
              {tr.eventCount !== undefined ? (
                <span className="text-xs text-muted-foreground">
                  {t('translate.track.events', { count: tr.eventCount })}
                </span>
              ) : null}
              <span className="ml-auto flex gap-1">
                {tr.default ? <Badge>{t('translate.track.default')}</Badge> : null}
                {tr.forced ? <Badge tone="violet">{t('translate.track.forced')}</Badge> : null}
                {tr.unsupported ? (
                  <Badge tone="warning">
                    {tr.unsupported === 'image-based'
                      ? t('translate.track.imageBased')
                      : t('translate.track.noParser')}
                  </Badge>
                ) : null}
              </span>
            </button>
          )
        })}
      </div>
      {probe.data.attachments > 0 ? (
        <p className="text-xs text-muted-foreground">
          {t('translate.track.fonts', { count: probe.data.attachments })}
        </p>
      ) : null}
    </div>
  )
}
