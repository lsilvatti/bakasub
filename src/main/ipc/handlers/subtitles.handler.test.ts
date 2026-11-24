import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AppError } from '../../errors'

type Opts = {
  input?: { safeParse(v: unknown): { success: boolean; data?: unknown } }
  handle(p: unknown, event: unknown): unknown
}
const registered = new Map<string, Opts>()
vi.mock('../create-handler', () => ({
  createHandler: (channel: string, opts: Opts) => registered.set(channel, opts)
}))
vi.mock('../../logging/logger', () => ({ redactor: { add: vi.fn() } }))

import { registerSubtitleHandlers } from './subtitles.handler'

const ID = '0b7f7e43-5b55-4a1c-9d0e-6a0f1f4e2f11'
const sender = (id: number): { sender: { id: number } } => ({ sender: { id } })
const call = (channel: string, payload: unknown, event = sender(1)): unknown => {
  const opts = registered.get(channel)!
  const parsed = opts.input!.safeParse(payload)
  if (!parsed.success) throw new Error('invalid input')
  return opts.handle(parsed.data, event)
}

describe('subtitles:search cancellation', () => {
  const search = vi.fn()
  beforeEach(() => {
    registered.clear()
    search.mockReset()
    registerSubtitleHandlers({
      vault: { has: vi.fn() } as never,
      subtitleSources: {
        client: { validate: vi.fn() },
        session: { login: vi.fn(), logout: vi.fn(), username: vi.fn() },
        service: { search, download: vi.fn() } as never
      }
    })
    search.mockImplementation(
      (_req: unknown, signal?: AbortSignal) =>
        new Promise((resolve, reject) => {
          if (!signal) return resolve({ query: 'q', results: [] })
          signal.addEventListener('abort', () => reject(new AppError('CANCELLED', 'cancelled')))
        })
    )
  })

  it('a tagged search gets a signal the cancel channel aborts; searchId never reaches the service', async () => {
    const running = call('subtitles:search', { videoPath: '/v.mkv', languages: [], searchId: ID })
    expect(search.mock.calls[0]![0]).toEqual({ videoPath: '/v.mkv', languages: [] })
    expect(await call('subtitles:cancel-search', { searchId: ID })).toEqual({ cancelled: true })
    await expect(running).rejects.toMatchObject({ code: 'CANCELLED' })
    expect(await call('subtitles:cancel-search', { searchId: ID })).toEqual({ cancelled: false })
  })

  it('only the sender that started it can cancel it', async () => {
    const running = call('subtitles:search', { videoPath: '/v.mkv', languages: [], searchId: ID })
    expect(await call('subtitles:cancel-search', { searchId: ID }, sender(2))).toEqual({
      cancelled: false
    })
    await call('subtitles:cancel-search', { searchId: ID })
    await expect(running).rejects.toMatchObject({ code: 'CANCELLED' })
  })

  it('an untagged search still works (no signal)', async () => {
    await expect(
      call('subtitles:search', { videoPath: '/v.mkv', languages: [] })
    ).resolves.toMatchObject({ query: 'q' })
    expect(search.mock.calls[0]![1]).toBeUndefined()
  })

  it('searchId must be a UUID', () => {
    expect(() => call('subtitles:cancel-search', { searchId: 'nope' })).toThrow('invalid input')
    expect(() =>
      call('subtitles:search', { videoPath: '/v.mkv', languages: [], searchId: '../x' })
    ).toThrow('invalid input')
  })
})
