import { describe, expect, it } from 'vitest'
import { AbortRegistry } from './abort-registry'

const hang = (signal: AbortSignal): Promise<never> =>
  new Promise((_r, reject) => signal.addEventListener('abort', () => reject(new Error('aborted'))))

describe('AbortRegistry', () => {
  it('cancel aborts the matching request and drops it', async () => {
    const reg = new AbortRegistry()
    const running = reg.run(1, 'a', hang)
    expect(reg.cancel(1, 'a')).toBe(true)
    await expect(running).rejects.toThrow('aborted')
    expect(reg.size).toBe(0)
    expect(reg.cancel(1, 'a')).toBe(false)
  })

  it('another sender cannot cancel it', async () => {
    const reg = new AbortRegistry()
    const running = reg.run(1, 'a', hang)
    expect(reg.cancel(2, 'a')).toBe(false)
    expect(reg.size).toBe(1)
    reg.cancel(1, 'a')
    await expect(running).rejects.toThrow()
  })

  it('a settled request leaves nothing behind', async () => {
    const reg = new AbortRegistry()
    await expect(reg.run(1, 'a', async () => 42)).resolves.toBe(42)
    await expect(
      reg.run(1, 'b', async () => {
        throw new Error('boom')
      })
    ).rejects.toThrow('boom')
    expect(reg.size).toBe(0)
  })

  it('reusing an id in flight aborts the older request; the newer one keeps its entry', async () => {
    const reg = new AbortRegistry()
    const first = reg.run(1, 'a', hang)
    const second = reg.run(1, 'a', hang)
    await expect(first).rejects.toThrow('aborted')
    expect(reg.size).toBe(1)
    expect(reg.cancel(1, 'a')).toBe(true)
    await expect(second).rejects.toThrow('aborted')
  })
})
