import { useState } from 'react'
import { FileText, Globe, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Badge } from '@renderer/components/ui/badge'
import { Button } from '@renderer/components/ui/button'
import { pickFile } from '@renderer/data/media'
import { errorMessage } from '@renderer/i18n/error-message'
import { languageName } from '@renderer/lib/languages'
import { cn } from '@renderer/lib/utils'
import { OpenSubtitlesSearch } from './opensubtitles-search'

export interface ExternalSubtitleValue {
  path: string
  name: string
  origin: 'file' | 'opensubtitles'
  language?: string
}

const baseName = (path: string): string => path.split(/[\\/]/).pop() ?? path

export function ExternalSubtitle({
  videoPath,
  value,
  onChange,
  prominent
}: {
  videoPath: string
  value: ExternalSubtitleValue | null
  onChange: (value: ExternalSubtitleValue | null) => void
  prominent: boolean
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const [searching, setSearching] = useState(false)

  const browse = (): void => {
    pickFile('subtitle').then(
      (path) => path && onChange({ path, name: baseName(path), origin: 'file' }),
      (err: unknown) => toast.error(errorMessage(err))
    )
  }

  return (
    <div
      className={cn(
        'grid gap-3',
        prominent ? 'rounded-lg border border-warning/40 bg-warning/5 p-4' : 'border-t pt-4'
      )}
    >
      {value ? (
        <div className="flex flex-wrap items-center gap-2 rounded-lg border border-primary bg-primary-soft/40 px-3 py-2 text-sm">
          {value.origin === 'file' ? (
            <FileText className="size-4 shrink-0" />
          ) : (
            <Globe className="size-4 shrink-0" />
          )}
          <span className="min-w-0 flex-1 truncate font-medium" title={value.name}>
            {value.name}
          </span>
          {value.language ? <Badge>{languageName(value.language, i18n.language)}</Badge> : null}
          <Badge tone="primary">{t(`translate.external.origin.${value.origin}`)}</Badge>
          <Button
            size="icon"
            variant="ghost"
            aria-label={t('translate.external.remove')}
            onClick={() => onChange(null)}
          >
            <X />
          </Button>
        </div>
      ) : (
        <div className="grid gap-1">
          <p className={cn('text-sm font-medium', prominent && 'text-warning')}>
            {prominent ? t('translate.track.none') : t('translate.external.other')}
          </p>
          <p className="text-xs text-muted-foreground">
            {prominent ? t('translate.external.bodyNone') : t('translate.external.body')}
          </p>
        </div>
      )}
      <div className="flex flex-wrap gap-2">
        <Button size="sm" onClick={browse}>
          <FileText />
          {t('translate.external.file')}
        </Button>
        <Button size="sm" onClick={() => setSearching(true)}>
          <Globe />
          {t('translate.external.search')}
        </Button>
      </div>
      <OpenSubtitlesSearch
        videoPath={videoPath}
        open={searching}
        onOpenChange={setSearching}
        onPick={(s) => onChange({ ...s, origin: 'opensubtitles' })}
      />
    </div>
  )
}
