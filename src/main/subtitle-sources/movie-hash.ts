import { open } from 'node:fs/promises'

const CHUNK = 64 * 1024
const MASK = (1n << 64n) - 1n

/** OpenSubtitles "moviehash" (OSDb). `null` for files under 64 KiB. */
export async function movieHash(path: string): Promise<string | null> {
  const file = await open(path, 'r')
  try {
    const { size } = await file.stat()
    if (size < CHUNK) return null
    let hash = BigInt(size)
    const buf = Buffer.alloc(CHUNK)
    for (const position of [0, size - CHUNK]) {
      const { bytesRead } = await file.read(buf, 0, CHUNK, position)
      if (bytesRead !== CHUNK) return null
      for (let i = 0; i < CHUNK; i += 8) hash = (hash + buf.readBigUInt64LE(i)) & MASK
    }
    return hash.toString(16).padStart(16, '0')
  } finally {
    await file.close()
  }
}
