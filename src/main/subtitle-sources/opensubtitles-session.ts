import { z } from 'zod'
import { AppError } from '../errors'
import type { OpenSubtitlesAuth, OpenSubtitlesClient } from './opensubtitles-client'

const Stored = z.object({ username: z.string().min(1), password: z.string().min(1) })
export type OpenSubtitlesLogin = z.infer<typeof Stored>

export interface LoginStore {
  get(): Promise<string | null>
  set(value: string): Promise<void>
  delete(): Promise<void>
}

export class OpenSubtitlesSession {
  private auth: OpenSubtitlesAuth | null = null
  private pending: Promise<OpenSubtitlesAuth | null> | null = null

  constructor(
    private readonly d: {
      client: Pick<OpenSubtitlesClient, 'login'>
      store: LoginStore
      onSecret?: (secret: string) => void
    }
  ) {}

  async username(): Promise<string | null> {
    return (await this.stored())?.username ?? null
  }

  async login(username: string, password: string, signal?: AbortSignal): Promise<void> {
    this.d.onSecret?.(password)
    const auth = await this.d.client.login(username, password, signal)
    this.d.onSecret?.(auth.token)
    await this.d.store.set(JSON.stringify({ username, password } satisfies OpenSubtitlesLogin))
    this.auth = auth
  }

  async logout(): Promise<void> {
    this.auth = null
    this.pending = null
    await this.d.store.delete()
  }
  async run<T>(
    call: (auth: OpenSubtitlesAuth | undefined) => Promise<T>,
    signal?: AbortSignal
  ): Promise<T> {
    const auth = await this.current(signal)
    try {
      return await call(auth ?? undefined)
    } catch (err) {
      if (!auth || !(err instanceof AppError) || err.code !== 'OPENSUBTITLES_AUTH') throw err
      if (this.auth === auth) this.auth = null
      return call((await this.current(signal)) ?? undefined)
    }
  }

  private async current(signal?: AbortSignal): Promise<OpenSubtitlesAuth | null> {
    if (this.auth) return this.auth
    this.pending ??= this.relogin(signal).finally(() => (this.pending = null))
    return this.pending
  }

  private async relogin(signal?: AbortSignal): Promise<OpenSubtitlesAuth | null> {
    const login = await this.stored()
    if (!login) return null
    this.d.onSecret?.(login.password)
    const auth = await this.d.client.login(login.username, login.password, signal)
    this.d.onSecret?.(auth.token)
    return (this.auth = auth)
  }

  private async stored(): Promise<OpenSubtitlesLogin | null> {
    const raw = await this.d.store.get()
    if (!raw) return null
    try {
      const parsed = Stored.safeParse(JSON.parse(raw))
      return parsed.success ? parsed.data : null
    } catch {
      return null
    }
  }
}
