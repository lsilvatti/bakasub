import { useTranslation } from 'react-i18next'
import { GuideDialog } from '@renderer/components/common/guide-dialog'
import { KeyForm } from '@renderer/components/forms/key-form'
import {
  useApiKeyStatus,
  useOpenSubtitlesStatus,
  useSetApiKey,
  useSetOpenSubtitlesKey,
  useSetTmdbKey,
  useTmdbKeyStatus
} from '@renderer/data/settings'
import { OpenSubtitlesAccount } from './opensubtitles-account'

/** The OpenRouter key cannot be removed here: the app needs it (onboarding would open again). */
export function KeysSection(): React.JSX.Element {
  const { t } = useTranslation()
  const or = useApiKeyStatus()
  const setOr = useSetApiKey()
  const tmdb = useTmdbKeyStatus()
  const setTmdb = useSetTmdbKey()
  const os = useOpenSubtitlesStatus()
  const setOs = useSetOpenSubtitlesKey()
  return (
    <div className="grid gap-6">
      <section className="grid gap-2">
        <div className="flex items-center justify-between gap-2">
          <h3 className="text-sm font-semibold">{t('settings.keys.openrouter')}</h3>
          <GuideDialog id="openrouter-key" />
        </div>
        <p className="text-xs text-muted-foreground">{t('settings.keys.openrouterHint')}</p>
        <KeyForm
          label={t('welcome.key.label')}
          placeholder="sk-or-v1-…"
          present={!!or.data?.present}
          noKeyring={or.data?.protection === 'none'}
          saving={setOr.isPending}
          error={setOr.error}
          onSave={(key) => setOr.mutateAsync(key)}
        />
      </section>
      <section className="grid gap-2">
        <div className="flex items-center justify-between gap-2">
          <h3 className="text-sm font-semibold">{t('settings.keys.tmdb')}</h3>
          <GuideDialog id="tmdb-key" />
        </div>
        <p className="text-xs text-muted-foreground">{t('settings.keys.tmdbHint')}</p>
        <KeyForm
          label={t('welcome.tmdb.label')}
          present={!!tmdb.data?.present}
          saving={setTmdb.isPending}
          error={setTmdb.error}
          onSave={(key) => setTmdb.mutateAsync(key)}
          onRemove={() => setTmdb.mutate(null)}
        />
      </section>
      <section className="grid gap-2">
        <div className="flex items-center justify-between gap-2">
          <h3 className="text-sm font-semibold">{t('settings.keys.opensubtitles')}</h3>
          <GuideDialog id="opensubtitles-key" />
        </div>
        <p className="text-xs text-muted-foreground">{t('settings.keys.opensubtitlesHint')}</p>
        <KeyForm
          label={t('settings.keys.opensubtitlesLabel')}
          present={!!os.data?.keyPresent}
          saving={setOs.isPending}
          error={setOs.error}
          onSave={(key) => setOs.mutateAsync(key)}
          onRemove={() => setOs.mutate(null)}
        />
        {os.data?.keyPresent ? <OpenSubtitlesAccount username={os.data.username} /> : null}
      </section>
    </div>
  )
}
