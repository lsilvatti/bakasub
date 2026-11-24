import { describe, expect, it, vi } from 'vitest'
import { AppError } from '../errors'
import type { OpenSubtitlesAuth } from './opensubtitles-client'
import { OpenSubtitlesSession, type LoginStore } from './opensubtitles-session'

function setup(stored: string | null = null): {
  session: OpenSubtitlesSession
  login: ReturnType<typeof vi.fn>
  store: { value: string | null }
  secrets: string[]
} {
  const store = { value: stored }
  const secrets: string[] = []
  let n = 0
  const login = vi.fn(async (): Promise<OpenSubtitlesAuth> => ({ token: `token-${++n}` }))
  const loginStore: LoginStore = {
    get: async () => store.value,
    set: async (v) => void (store.value = v),
    delete: async () => void (store.value = null)
  }
  return {
    session: new OpenSubtitlesSession({
      client: { login },
      store: loginStore,
      onSecret: (s) => secrets.push(s)
    }),
    login,
    store,
    secrets
  }
}

const CREDS = JSON.stringify({ username: 'u', password: 'hunter2-long' })
const authError = (): AppError => new AppError('OPENSUBTITLES_AUTH', 'rejected')

describe('OpenSubtitlesSession', () => {
  it('no account: calls anonymously and never logs in', async () => {
    const { session, login } = setup()
    const call = vi.fn(async (auth?: OpenSubtitlesAuth) => auth ?? 'anonymous')
    expect(await session.run(call)).toBe('anonymous')
    expect(login).not.toHaveBeenCalled()
    expect(await session.username()).toBeNull()
  })

  it('login stores the credentials only when accepted; password and token go to the redactor', async () => {
    const ok = setup()
    await ok.session.login('u', 'hunter2-long')
    expect(JSON.parse(ok.store.value!)).toEqual({ username: 'u', password: 'hunter2-long' })
    expect(ok.secrets).toEqual(['hunter2-long', 'token-1'])
    expect(await ok.session.username()).toBe('u')

    const bad = setup()
    bad.login.mockRejectedValueOnce(authError())
    await expect(bad.session.login('u', 'wrong-password')).rejects.toMatchObject({
      code: 'OPENSUBTITLES_AUTH'
    })
    expect(bad.store.value).toBeNull()
  })

  it('stored account: logs in once, even with concurrent calls, and reuses the token', async () => {
    const { session, login } = setup(CREDS)
    const seen = await Promise.all([
      session.run(async (a) => a?.token),
      session.run(async (a) => a?.token)
    ])
    expect(seen).toEqual(['token-1', 'token-1'])
    expect(await session.run(async (a) => a?.token)).toBe('token-1')
    expect(login).toHaveBeenCalledTimes(1)
  })

  it('rejected token: logs in again and retries once', async () => {
    const { session, login } = setup(CREDS)
    const call = vi
      .fn<(a?: OpenSubtitlesAuth) => Promise<string | undefined>>()
      .mockRejectedValueOnce(authError())
      .mockImplementation(async (a) => a?.token)
    expect(await session.run(call)).toBe('token-2')
    expect(login).toHaveBeenCalledTimes(2)
  })

  it('a second rejection propagates; a non-auth error does not log in again; anonymous does not retry', async () => {
    const a = setup(CREDS)
    await expect(a.session.run(async () => Promise.reject(authError()))).rejects.toMatchObject({
      code: 'OPENSUBTITLES_AUTH'
    })
    expect(a.login).toHaveBeenCalledTimes(2)

    const b = setup(CREDS)
    await expect(
      b.session.run(async () => Promise.reject(new AppError('NETWORK', 'down')))
    ).rejects.toMatchObject({ code: 'NETWORK' })
    expect(b.login).toHaveBeenCalledTimes(1)

    const c = setup()
    const call = vi.fn(async () => Promise.reject(authError()))
    await expect(c.session.run(call)).rejects.toMatchObject({ code: 'OPENSUBTITLES_AUTH' })
    expect(call).toHaveBeenCalledTimes(1)
  })

  it('logout clears credentials and token; corrupted credentials count as no account', async () => {
    const { session, store } = setup(CREDS)
    await session.run(async () => undefined)
    await session.logout()
    expect(store.value).toBeNull()
    expect(await session.run(async (a) => a ?? 'anonymous')).toBe('anonymous')

    const broken = setup('{not json')
    expect(await broken.session.username()).toBeNull()
    expect(await broken.session.run(async (a) => a ?? 'anonymous')).toBe('anonymous')
  })
})
