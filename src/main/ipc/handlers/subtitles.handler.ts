import { z } from 'zod'
import type { OpenSubtitlesStatus } from '../../../shared/ipc-contract'
import { redactor } from '../../logging/logger'
import type { SecretVault } from '../../settings/secret-vault'
import type { OpenSubtitlesClient } from '../../subtitle-sources/opensubtitles-client'
import type { OpenSubtitlesSession } from '../../subtitle-sources/opensubtitles-session'
import type { SubtitleSourceService } from '../../subtitle-sources/subtitle-source.service'
import { createHandler } from '../create-handler'

const secret = (min: number, max: number): z.ZodString =>
  z
    .string()
    .min(min)
    .max(max)
    .regex(/^[\x21-\x7e]+$/)
const Num = z.number().int().min(0).max(9999)

export function registerSubtitleHandlers(deps: {
  vault: SecretVault
  subtitleSources: {
    client: Pick<OpenSubtitlesClient, 'validate'>
    session: Pick<OpenSubtitlesSession, 'login' | 'logout' | 'username'>
    service: Pick<SubtitleSourceService, 'search' | 'download'>
  }
}): void {
  const { client, session, service } = deps.subtitleSources
  const status = async (): Promise<OpenSubtitlesStatus> => ({
    keyPresent: await deps.vault.has('opensubtitles.apiKey'),
    username: await session.username()
  })

  createHandler('settings:set-opensubtitles-key', {
    input: z.object({ key: secret(8, 256) }),
    handle: async ({ key }) => {
      redactor.add(key)
      await client.validate(key)
      await deps.vault.set('opensubtitles.apiKey', key)
      return status()
    }
  })
  createHandler('settings:clear-opensubtitles-key', {
    handle: async () => {
      await session.logout()
      await deps.vault.delete('opensubtitles.apiKey')
      return status()
    }
  })
  createHandler('settings:opensubtitles-status', { handle: status })
  createHandler('settings:opensubtitles-login', {
    input: z.object({
      username: z.string().trim().min(1).max(100),
      password: z.string().min(1).max(256)
    }),
    handle: async ({ username, password }) => {
      redactor.add(password)
      await session.login(username, password)
      return status()
    }
  })
  createHandler('settings:opensubtitles-logout', {
    handle: async () => {
      await session.logout()
      return status()
    }
  })

  createHandler('subtitles:search', {
    input: z.object({
      videoPath: z.string().min(1).max(4096),
      query: z.string().trim().max(200).optional(),
      languages: z
        .array(
          z
            .string()
            .min(2)
            .max(16)
            .regex(/^[A-Za-z0-9-]+$/)
        )
        .max(10),
      season: Num.optional(),
      episode: Num.optional()
    }),
    handle: (req) => service.search(req)
  })
  createHandler('subtitles:download', {
    input: z.object({ fileId: z.number().int().positive().max(Number.MAX_SAFE_INTEGER) }),
    handle: ({ fileId }) => service.download(fileId)
  })
}
