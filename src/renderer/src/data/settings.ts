import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseMutationResult,
  type UseQueryResult
} from '@tanstack/react-query'
import type {
  ApiKeyStatus,
  OpenSubtitlesStatus,
  PublicSettings,
  TmdbKeyStatus
} from '@shared/ipc-contract'
import { isOnboarded as onboardingRule } from '@shared/onboarding'
import { unwrap } from '@renderer/lib/ipc'
import { qk } from './query-keys'

/** Shared with the router guard (`ensureQueryData`). */
export const settingsQuery = {
  queryKey: qk.settings.all(),
  queryFn: (): Promise<PublicSettings> => unwrap(window.api.settings.get())
}
export const apiKeyStatusQuery = {
  queryKey: qk.settings.apiKey(),
  queryFn: (): Promise<ApiKeyStatus> => unwrap(window.api.settings.apiKeyStatus())
}

export const useSettings = (): UseQueryResult<PublicSettings> => useQuery(settingsQuery)

/** What the app needs before translating (rule shared with the Main's auto update gate). */
export const isOnboarded = (
  key: ApiKeyStatus | undefined,
  s: PublicSettings | undefined
): boolean =>
  onboardingRule({
    apiKeyPresent: !!key?.present,
    defaultModel: s?.['translation.defaultModel'],
    defaultTargetLang: s?.['translation.defaultTargetLang']
  })

/** The Main answers with the full snapshot (rules applied): it replaces the cache as is. */
export function useUpdateSettings(): UseMutationResult<
  PublicSettings,
  Error,
  Partial<PublicSettings>
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (patch) => unwrap(window.api.settings.update(patch)),
    onSuccess: (s) => {
      qc.setQueryData(qk.settings.all(), s)
      // Tier override changes the quota the Main reports.
      void qc.invalidateQueries({ queryKey: qk.models.freeQuota() })
    }
  })
}

export function useResetSettings(): UseMutationResult<
  PublicSettings,
  Error,
  (keyof PublicSettings)[]
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (keys) => unwrap(window.api.settings.reset(keys)),
    onSuccess: (s) => {
      qc.setQueryData(qk.settings.all(), s)
      void qc.invalidateQueries({ queryKey: qk.models.freeQuota() })
    }
  })
}

export const useApiKeyStatus = (): UseQueryResult<ApiKeyStatus> => useQuery(apiKeyStatusQuery)

/** Write only: the key never comes back. A new key changes the catalog auth and the free quota. */
export function useSetApiKey(): UseMutationResult<ApiKeyStatus, Error, string | null> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (key) =>
      unwrap(key === null ? window.api.settings.clearApiKey() : window.api.settings.setApiKey(key)),
    onSuccess: (s) => {
      qc.setQueryData(qk.settings.apiKey(), s)
      void qc.invalidateQueries({ queryKey: qk.models.all })
    }
  })
}

export const useTmdbKeyStatus = (): UseQueryResult<TmdbKeyStatus> =>
  useQuery({
    queryKey: qk.settings.tmdbKey(),
    queryFn: () => unwrap(window.api.settings.tmdbKeyStatus())
  })

/** Validated against TMDB before saving (`TMDB_AUTH` when rejected). `null` removes it. */
export function useSetTmdbKey(): UseMutationResult<TmdbKeyStatus, Error, string | null> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (key) =>
      unwrap(
        key === null ? window.api.settings.clearTmdbKey() : window.api.settings.setTmdbKey(key)
      ),
    onSuccess: (s) => {
      qc.setQueryData(qk.settings.tmdbKey(), s)
      // Identification said `no-key` before: ask again.
      void qc.invalidateQueries({ queryKey: ['media'] })
    }
  })
}

export const useOpenSubtitlesStatus = (): UseQueryResult<OpenSubtitlesStatus> =>
  useQuery({
    queryKey: qk.settings.openSubtitles(),
    queryFn: () => unwrap(window.api.settings.openSubtitlesStatus())
  })

export function useSetOpenSubtitlesKey(): UseMutationResult<
  OpenSubtitlesStatus,
  Error,
  string | null
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (key) =>
      unwrap(
        key === null
          ? window.api.settings.clearOpenSubtitlesKey()
          : window.api.settings.setOpenSubtitlesKey(key)
      ),
    onSuccess: (s) => {
      qc.setQueryData(qk.settings.openSubtitles(), s)
      void qc.invalidateQueries({ queryKey: ['subtitles'] })
    }
  })
}

export function useOpenSubtitlesLogin(): UseMutationResult<
  OpenSubtitlesStatus,
  Error,
  { username: string; password: string } | null
> {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (login) =>
      unwrap(
        login === null
          ? window.api.settings.openSubtitlesLogout()
          : window.api.settings.openSubtitlesLogin(login.username, login.password)
      ),
    onSuccess: (s) => qc.setQueryData(qk.settings.openSubtitles(), s)
  })
}
