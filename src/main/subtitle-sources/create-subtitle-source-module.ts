import { join } from 'node:path'
import type { SecretVault } from '../settings/secret-vault'
import { movieHash } from './movie-hash'
import { OpenSubtitlesClient } from './opensubtitles-client'
import { OpenSubtitlesSession } from './opensubtitles-session'
import { SubtitleSourceService } from './subtitle-source.service'

export function createSubtitleSourceModule(o: {
  vault: SecretVault
  cacheDir: string
  appVersion: string
  onSecret?: (secret: string) => void
}): {
  client: OpenSubtitlesClient
  session: OpenSubtitlesSession
  service: SubtitleSourceService
} {
  const client = new OpenSubtitlesClient({
    apiKey: () => o.vault.get('opensubtitles.apiKey'),
    userAgent: `BakaSub v${o.appVersion}`
  })
  const session = new OpenSubtitlesSession({
    client,
    store: {
      get: () => o.vault.get('opensubtitles.login'),
      set: (value) => o.vault.set('opensubtitles.login', value),
      delete: () => o.vault.delete('opensubtitles.login')
    },
    onSecret: o.onSecret
  })
  const service = new SubtitleSourceService({
    client,
    session,
    cacheDir: join(o.cacheDir, 'opensubtitles'),
    hash: movieHash
  })
  return { client, session, service }
}
