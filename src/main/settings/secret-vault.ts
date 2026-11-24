import { createCipheriv, createDecipheriv, randomBytes } from 'node:crypto'
import type { SecretsRepository } from '../database/repositories/secrets.repository'

/** Subconjunto do `safeStorage` do Electron (injetado: o cofre não importa `electron`). */
export interface SafeStorageLike {
  isEncryptionAvailable(): boolean
  encryptString(plain: string): Buffer
  decryptString(encrypted: Buffer): string
  /** Só Linux. `basic_text` = sem keyring real (ofuscação). */
  getSelectedStorageBackend?(): string
}

/** Segredos conhecidos. Novos segredos entram aqui. */
export type SecretKey =
  'openrouter.apiKey' | 'tmdb.credential' | 'opensubtitles.apiKey' | 'opensubtitles.login'

/** `os`: chave mestra protegida pelo cofre do SO. `none`: guardada crua no banco. */
export type VaultProtection = 'os' | 'none'

const KEY_VERSION = 1
const ALGO = 'aes-256-gcm'

function osAvailable(safe: SafeStorageLike | null): safe is SafeStorageLike {
  if (!safe?.isEncryptionAvailable()) return false
  return safe.getSelectedStorageBackend?.() !== 'basic_text'
}

/**
 * Envelope encryption: uma chave mestra aleatória (32 bytes) cifra cada segredo com
 * AES-256-GCM (IV por registro; o nome do segredo é AAD, então copiar ciphertext entre linhas
 * é detectado). A chave mestra fica em `vault_keys`, embrulhada pelo `safeStorage` quando o SO
 * oferece um cofre real.
 */
export class SecretVault {
  private constructor(
    private readonly repo: SecretsRepository,
    private readonly master: Buffer,
    readonly protection: VaultProtection,
    /** `true` quando a chave mestra anterior era irrecuperável e o cofre foi recriado. */
    readonly recovered: boolean
  ) {}

  static async open(repo: SecretsRepository, safe: SafeStorageLike | null): Promise<SecretVault> {
    const os = osAvailable(safe)
    const stored = await repo.latestKey()

    if (stored) {
      const master = unwrap(stored.wrapped, stored.protection, safe)
      if (master) {
        // Promoção: chave guardada crua passa a ser protegida quando o SO oferece cofre.
        if (stored.protection === 'none' && os) {
          await repo.putKey({ ...stored, wrapped: wrap(master, safe), protection: 'os' })
          return new SecretVault(repo, master, 'os', false)
        }
        return new SecretVault(repo, master, stored.protection, false)
      }
      // Keyring trocado/perdido: sem a chave mestra os segredos são lixo. Recomeça.
      await repo.reset()
    }

    const master = randomBytes(32)
    const protection: VaultProtection = os ? 'os' : 'none'
    await repo.putKey({
      version: KEY_VERSION,
      wrapped: os ? wrap(master, safe) : master,
      protection
    })
    return new SecretVault(repo, master, protection, stored !== null)
  }

  async set(key: SecretKey, value: string): Promise<void> {
    const iv = randomBytes(12)
    const cipher = createCipheriv(ALGO, this.master, iv)
    cipher.setAAD(Buffer.from(key, 'utf8'))
    const ciphertext = Buffer.concat([cipher.update(value, 'utf8'), cipher.final()])
    await this.repo.put(key, { ciphertext, iv, tag: cipher.getAuthTag(), keyVersion: KEY_VERSION })
  }

  /** `null` quando ausente ou quando não autentica (adulterado/movido): nunca lança por conteúdo. */
  async get(key: SecretKey): Promise<string | null> {
    const row = await this.repo.get(key)
    if (!row) return null
    try {
      const decipher = createDecipheriv(ALGO, this.master, row.iv)
      decipher.setAAD(Buffer.from(key, 'utf8'))
      decipher.setAuthTag(row.tag)
      return Buffer.concat([decipher.update(row.ciphertext), decipher.final()]).toString('utf8')
    } catch {
      return null
    }
  }

  async has(key: SecretKey): Promise<boolean> {
    return (await this.get(key)) !== null
  }

  delete(key: SecretKey): Promise<void> {
    return this.repo.delete(key)
  }
}

function wrap(master: Buffer, safe: SafeStorageLike | null): Buffer {
  return safe!.encryptString(master.toString('base64'))
}

function unwrap(
  wrapped: Buffer,
  protection: VaultProtection,
  safe: SafeStorageLike | null
): Buffer | null {
  if (protection === 'none') return wrapped.length === 32 ? Buffer.from(wrapped) : null
  if (!safe?.isEncryptionAvailable()) return null
  try {
    const key = Buffer.from(safe.decryptString(Buffer.from(wrapped)), 'base64')
    return key.length === 32 ? key : null
  } catch {
    return null
  }
}
