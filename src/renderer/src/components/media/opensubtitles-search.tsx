import { useState, type FormEvent } from 'react'
import { Download, Search } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { SubtitleSearchRequest, SubtitleSearchResult } from '@shared/ipc-contract'
import { GuideDialog } from '@renderer/components/common/guide-dialog'
import { KeyForm } from '@renderer/components/forms/key-form'
import { Badge } from '@renderer/components/ui/badge'
import { Button } from '@renderer/components/ui/button'
import { Dialog, DialogContent } from '@renderer/components/ui/dialog'
import { Input } from '@renderer/components/ui/input'
import {
  useOpenSubtitlesStatus,
  useSetOpenSubtitlesKey,
  useSettings
} from '@renderer/data/settings'
import { useSubtitleDownload, useSubtitleSearch } from '@renderer/data/subtitles'
import { errorMessage } from '@renderer/i18n/error-message'
import { languageName } from '@renderer/lib/languages'
import { cn } from '@renderer/lib/utils'

export interface PickedDownload {
  path: string
  name: string
  language?: string
}

const toNumber = (s: string): number | undefined => {
  const n = Number(s)
  return s.trim() !== '' && Number.isInteger(n) && n >= 0 && n <= 9999 ? n : undefined
}

export function OpenSubtitlesSearch({
  videoPath,
  open,
  onOpenChange,
  onPick
}: {
  videoPath: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onPick: (subtitle: PickedDownload) => void
}): React.JSX.Element {
  const { t } = useTranslation()
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        title={t('opensubtitles.title')}
        description={t('opensubtitles.description')}
        className="w-[min(46rem,calc(100vw-2rem))]"
      >
        <KeyGate
          videoPath={videoPath}
          onPick={(s) => {
            onPick(s)
            onOpenChange(false)
          }}
        />
      </DialogContent>
    </Dialog>
  )
}

function KeyGate({
  videoPath,
  onPick
}: {
  videoPath: string
  onPick: (subtitle: PickedDownload) => void
}): React.JSX.Element | null {
  const { t } = useTranslation()
  const status = useOpenSubtitlesStatus()
  const setKey = useSetOpenSubtitlesKey()

  if (status.isPending) return null
  if (status.isError)
    return (
      <p role="alert" className="text-sm text-destructive">
        {errorMessage(status.error)}
      </p>
    )
  if (status.data.keyPresent)
    return (
      <SearchBody
        key={videoPath}
        videoPath={videoPath}
        loggedIn={!!status.data.username}
        onPick={onPick}
      />
    )
  return (
    <div className="grid gap-3">
      <p className="text-sm text-muted-foreground">{t('opensubtitles.needsKey')}</p>
      <KeyForm
        label={t('settings.keys.opensubtitlesLabel')}
        present={false}
        saving={setKey.isPending}
        error={setKey.error}
        onSave={(key) => setKey.mutateAsync(key)}
        guide={<GuideDialog id="opensubtitles-key" />}
      />
    </div>
  )
}

function SearchBody({
  videoPath,
  loggedIn,
  onPick
}: {
  videoPath: string
  loggedIn: boolean
  onPick: (subtitle: PickedDownload) => void
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const preferred = useSettings().data?.['tracks.preferredSourceLangs']
  const known = preferred?.length ? preferred : ['en']
  const [query, setQuery] = useState<string | null>(null)
  const [season, setSeason] = useState<string | null>(null)
  const [episode, setEpisode] = useState<string | null>(null)
  const [languages, setLanguages] = useState<string[] | null>(null)
  const langs = languages ?? known
  const [request, setRequest] = useState<SubtitleSearchRequest | null>(null)
  const search = useSubtitleSearch(request ?? { videoPath, languages: known })
  const download = useSubtitleDownload()

  const shown = {
    query: query ?? search.data?.query ?? '',
    season: season ?? (search.data?.season !== undefined ? String(search.data.season) : ''),
    episode: episode ?? (search.data?.episode !== undefined ? String(search.data.episode) : '')
  }

  const submit = (e: FormEvent): void => {
    e.preventDefault()
    const s = toNumber(shown.season)
    const ep = toNumber(shown.episode)
    setRequest({
      videoPath,
      ...(shown.query.trim() ? { query: shown.query.trim() } : {}),
      languages: langs,
      ...(s !== undefined ? { season: s } : {}),
      ...(ep !== undefined ? { episode: ep } : {})
    })
  }

  const toggle = (tag: string): void =>
    setLanguages(langs.includes(tag) ? langs.filter((l) => l !== tag) : [...langs, tag])

  const pick = (r: SubtitleSearchResult): void =>
    download.mutate(r.fileId, {
      onSuccess: (d) => {
        if (d.remaining !== undefined)
          toast.info(t('opensubtitles.remaining', { count: d.remaining }))
        onPick({ path: d.path, name: r.release, ...(r.language ? { language: r.language } : {}) })
      },
      onError: (err) => toast.error(errorMessage(err))
    })

  const results = search.data?.results ?? []

  return (
    <div className="grid gap-4">
      <form onSubmit={submit} className="grid gap-3">
        <div className="grid grid-cols-[minmax(0,1fr)_5rem_5rem] gap-2">
          <label className="grid gap-1 text-xs font-medium text-muted-foreground">
            {t('opensubtitles.query')}
            <Input
              value={shown.query}
              autoComplete="off"
              onChange={(e) => setQuery(e.target.value)}
            />
          </label>
          <label className="grid gap-1 text-xs font-medium text-muted-foreground">
            {t('opensubtitles.season')}
            <Input
              inputMode="numeric"
              value={shown.season}
              onChange={(e) => setSeason(e.target.value)}
            />
          </label>
          <label className="grid gap-1 text-xs font-medium text-muted-foreground">
            {t('opensubtitles.episode')}
            <Input
              inputMode="numeric"
              value={shown.episode}
              onChange={(e) => setEpisode(e.target.value)}
            />
          </label>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <div
            role="group"
            aria-label={t('opensubtitles.languages')}
            className="flex flex-wrap gap-1.5"
          >
            {known.map((tag) => {
              const on = langs.includes(tag)
              return (
                <button
                  key={tag}
                  type="button"
                  aria-pressed={on}
                  onClick={() => toggle(tag)}
                  className={cn(
                    'rounded-md border px-2 py-1 text-xs',
                    on ? 'border-primary bg-primary-soft/40' : 'text-muted-foreground'
                  )}
                >
                  {languageName(tag, i18n.language)}
                </button>
              )
            })}
          </div>
          {langs.length === 0 ? (
            <span className="text-xs text-muted-foreground">{t('opensubtitles.anyLanguage')}</span>
          ) : null}
          <Button type="submit" variant="default" className="ml-auto" disabled={search.isFetching}>
            <Search />
            {search.isFetching ? t('opensubtitles.searching') : t('opensubtitles.search')}
          </Button>
        </div>
      </form>

      {search.isError ? (
        <p role="alert" className="text-sm text-destructive">
          {errorMessage(search.error)}
        </p>
      ) : search.isPending ? (
        <p className="text-sm text-muted-foreground">{t('opensubtitles.searching')}</p>
      ) : results.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('opensubtitles.noResults')}</p>
      ) : (
        <ul className="grid max-h-[45vh] gap-1.5 overflow-y-auto pr-1">
          {results.map((r) => (
            <li
              key={r.fileId}
              className="flex items-center gap-3 rounded-lg border bg-background/40 px-3 py-2"
            >
              <div className="grid min-w-0 flex-1 gap-1">
                <span className="truncate text-sm font-medium" title={r.release}>
                  {r.release}
                </span>
                <span className="flex flex-wrap items-center gap-1.5 text-xs text-muted-foreground">
                  <span>{r.language ? languageName(r.language, i18n.language) : '—'}</span>
                  <span>{t('opensubtitles.downloads', { count: r.downloads })}</span>
                  {r.hashMatch ? (
                    <Badge tone="success">{t('opensubtitles.badge.hashMatch')}</Badge>
                  ) : null}
                  {r.trusted ? (
                    <Badge tone="primary">{t('opensubtitles.badge.trusted')}</Badge>
                  ) : null}
                  {r.hearingImpaired ? (
                    <Badge>{t('opensubtitles.badge.hearingImpaired')}</Badge>
                  ) : null}
                  {r.machineTranslated ? (
                    <Badge tone="warning">{t('opensubtitles.badge.machine')}</Badge>
                  ) : null}
                </span>
              </div>
              <Button
                size="sm"
                disabled={download.isPending}
                onClick={() => pick(r)}
                aria-label={t('opensubtitles.useNamed', { name: r.release })}
              >
                <Download />
                {download.isPending && download.variables === r.fileId
                  ? t('opensubtitles.downloading')
                  : t('opensubtitles.use')}
              </Button>
            </li>
          ))}
        </ul>
      )}

      <p className="text-xs text-muted-foreground">
        {t('opensubtitles.syncNote')}
        {loggedIn ? '' : ` ${t('opensubtitles.quotaNote')}`}
      </p>
    </div>
  )
}
