import { mkdtempSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { beforeAll, describe, expect, it } from 'vitest'
import { movieHash } from './movie-hash'

const CHUNK = 64 * 1024
let dir: string
beforeAll(() => {
  dir = mkdtempSync(join(tmpdir(), 'bakasub-hash-'))
})

const file = (name: string, content: Buffer): string => {
  const path = join(dir, name)
  writeFileSync(path, content)
  return path
}

describe('movieHash (OSDb)', () => {
  it('an all-zero file hashes to its own size', async () => {
    expect(await movieHash(file('zeros.mkv', Buffer.alloc(3 * CHUNK)))).toBe('0000000000030000')
  })

  it('sums the little-endian uint64s of the first and last chunk; the middle is ignored', async () => {
    const a = Buffer.alloc(3 * CHUNK)
    a.writeBigUInt64LE(5n, 0)
    a.writeBigUInt64LE(7n, 3 * CHUNK - 8)
    const b = Buffer.from(a)
    b.fill(0xff, CHUNK, 2 * CHUNK)
    expect(await movieHash(file('a.mkv', a))).toBe('000000000003000c')
    expect(await movieHash(file('b.mkv', b))).toBe('000000000003000c')
  })

  it('the sum wraps at 64 bits', async () => {
    const buf = Buffer.alloc(2 * CHUNK)
    buf.writeBigUInt64LE(0xffffffffffffffffn, 0)
    expect(await movieHash(file('wrap.mkv', buf))).toBe('000000000001ffff')
  })

  it('between 64 and 128 KiB the chunks overlap; under 64 KiB there is no hash', async () => {
    const buf = Buffer.alloc(CHUNK)
    buf.writeBigUInt64LE(1n, 0)
    expect(await movieHash(file('one.mkv', buf))).toBe('0000000000010002')
    expect(await movieHash(file('tiny.mkv', Buffer.alloc(CHUNK - 1)))).toBeNull()
  })
})
