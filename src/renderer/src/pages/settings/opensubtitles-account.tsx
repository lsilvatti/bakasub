import { useState, type FormEvent } from 'react'
import { CheckCircle2, LogIn, ShieldAlert } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@renderer/components/ui/button'
import { Field } from '@renderer/components/ui/field'
import { Input } from '@renderer/components/ui/input'
import { useOpenSubtitlesLogin } from '@renderer/data/settings'
import { errorMessage } from '@renderer/i18n/error-message'

/**
 * `noKeyring`: the secret vault has no OS keyring (`protection: 'none'`), so the password is kept
 * without system protection. Shown before sign-in too, so the user decides before typing it.
 */
export function OpenSubtitlesAccount({
  username,
  noKeyring = false
}: {
  username: string | null
  noKeyring?: boolean
}): React.JSX.Element {
  const { t } = useTranslation()
  const login = useOpenSubtitlesLogin()
  const [user, setUser] = useState('')
  const [password, setPassword] = useState('')

  const submit = (e: FormEvent): void => {
    e.preventDefault()
    if (!user.trim() || !password) return
    login.mutate(
      { username: user.trim(), password },
      {
        onSuccess: () => {
          setUser('')
          setPassword('')
        }
      }
    )
  }

  return (
    <div className="grid gap-2">
      <h4 className="text-sm font-medium">{t('settings.keys.opensubtitlesAccount')}</h4>
      <p className="text-xs text-muted-foreground">
        {t(
          noKeyring
            ? 'settings.keys.opensubtitlesAccountHintNoKeyring'
            : 'settings.keys.opensubtitlesAccountHint'
        )}
      </p>
      {noKeyring ? (
        <p role="alert" className="flex gap-2 text-xs text-warning">
          <ShieldAlert className="size-4 shrink-0" />
          {t('settings.keys.opensubtitlesNoKeyring')}
        </p>
      ) : null}
      {username ? (
        <div className="flex flex-wrap items-center gap-3 rounded-md border bg-background/40 px-3 py-2">
          <CheckCircle2 className="size-4 text-success" />
          <span className="text-sm">{t('settings.keys.connectedAs', { username })}</span>
          <Button
            size="sm"
            variant="destructive"
            className="ml-auto"
            disabled={login.isPending}
            onClick={() => login.mutate(null)}
          >
            {t('settings.keys.disconnect')}
          </Button>
        </div>
      ) : (
        <form onSubmit={submit} className="grid gap-2">
          <div className="grid gap-2 sm:grid-cols-2">
            <Field label={t('settings.keys.username')} hint={t('settings.keys.usernameHint')}>
              <Input
                autoComplete="off"
                spellCheck={false}
                value={user}
                onChange={(e) => setUser(e.target.value)}
              />
            </Field>
            <Field
              label={t('settings.keys.password')}
              error={login.error ? errorMessage(login.error) : undefined}
            >
              <Input
                type="password"
                autoComplete="off"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Field>
          </div>
          <div>
            <Button type="submit" disabled={login.isPending || !user.trim() || !password}>
              <LogIn />
              {login.isPending ? t('settings.keys.connecting') : t('settings.keys.connect')}
            </Button>
          </div>
        </form>
      )}
    </div>
  )
}
