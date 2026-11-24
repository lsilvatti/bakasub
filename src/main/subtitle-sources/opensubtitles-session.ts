import { z } from 'zod'
import { AppError } from '../errors'
import { cancelled, detach } from './abortable'
import type { OpenSubtitlesAuth, OpenSubtitlesClient } from './opensubtitles-client'

const Stored = z.object({ username: z.string().min(1), password: z.string().min(1) })
export type OpenSubtitlesLogin = z.infer<typeof Stored>

export interface LoginStore {
  get(): Promise<string | null>
  set(value: string): Promise<void>
  delete(): Promise<void>
}

/**
 * Holds the OpenSubtitles token for the stored account.
 *
 * `generation` is bumped by every logout and every explicit login (credential change). Each
 * login in flight captures it when it starts and drops its result if it changed meanwhile, so
 * a slow re-login can never resurrect a session the user already logged out of, nor overwrite
 * newer credentials.
 */
export class OpenSubtitlesSession {
  private auth: OpenSubtitlesAuth | null = null
  private pending: Promise<OpenSubtitlesAuth | null> | null = null
  private generation = 0

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
    const gen = this.invalidate()
    this.d.onSecret?.(password)
    const auth = await this.d.client.login(username, password, signal)
    this.d.onSecret?.(auth.token)
    if (gen !== this.generation) throw superseded()
    await this.d.store.set(JSON.stringify({ username, password } satisfies OpenSubtitlesLogin))
    if (gen !== this.generation) throw superseded()
    this.auth = auth
  }

  async logout(): Promise<void> {
    this.invalidate()
    await this.d.store.delete()
  }

  /**
   * `signal` only detaches this caller from a shared re-login (which keeps running, bounded by the
   * client timeout, for the other callers); the call itself gets the signal through its closure.
   */
  async run<T>(
    call: (auth: OpenSubtitlesAuth | undefined) => Promise<T>,
    signal?: AbortSignal
  ): Promise<T> {
    if (signal?.aborted) throw cancelled(RUN)
    const auth = await this.current(signal)
    try {
      return await call(auth ?? undefined)
    } catch (err) {
      if (!auth || !(err instanceof AppError) || err.code !== 'OPENSUBTITLES_AUTH') throw err
      if (this.auth === auth) this.auth = null
      if (signal?.aborted) throw cancelled(RUN)
      return call((await this.current(signal)) ?? undefined)
    }
  }

  private async current(signal?: AbortSignal): Promise<OpenSubtitlesAuth | null> {
    if (this.auth) return this.auth
    if (!this.pending) {
      const pending = this.relogin(this.generation).finally(() => {
        if (this.pending === pending) this.pending = null
      })
      this.pending = pending
    }
    return signal ? detach(this.pending, signal, RUN) : this.pending
  }

  /** Re-login from the stored account; a result from a stale generation is discarded. */
  private async relogin(gen: number): Promise<OpenSubtitlesAuth | null> {
    const login = await this.stored()
    if (gen !== this.generation) return this.auth
    if (!login) return null
    this.d.onSecret?.(login.password)
    const auth = await this.d.client.login(login.username, login.password)
    this.d.onSecret?.(auth.token)
    if (gen !== this.generation) return this.auth
    return (this.auth = auth)
  }

  /** Drops the token and any re-login in flight; returns the new generation. */
  private invalidate(): number {
    this.auth = null
    this.pending = null
    return ++this.generation
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

const RUN = 'OpenSubtitles request'

const superseded = (): AppError =>
  new AppError('CANCELLED', 'OpenSubtitles login superseded by a logout or another login')
