import { mkdirSync, mkdtempSync, realpathSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { beforeAll, describe, expect, it, vi } from 'vitest'
import { TEST_PRESET, testPresetSource } from '../presets/test-presets'
import { AppError } from '../errors'
import { jobMemoryMode, resolveJobSpec } from './job-spec'

let dir: string
const presets = testPresetSource()
beforeAll(() => {
  dir = realpathSync.native(mkdtempSync(join(tmpdir(), 'bakasub-spec-')))
  writeFileSync(join(dir, 'ep01.mkv'), 'mkv')
  writeFileSync(join(dir, 'ep01.en.srt'), 'srt')
  writeFileSync(join(dir, 'ep01.pt-BR.ass'), 'ass')
  mkdirSync(join(dir, 'out'))
})

describe('resolveJobSpec', () => {
  it('translate: congela o preset resolvido, idioma padrão do perfil, sourceLang do nome, saída ao lado', async () => {
    const s = await resolveJobSpec(
      { kind: 'translate', source: join(dir, 'ep01.en.srt'), presetId: 'p-test' },
      presets
    )
    expect(s).toEqual({
      kind: 'translate',
      source: join(dir, 'ep01.en.srt'),
      sourceSize: 3,
      preset: {
        id: 'p-test',
        name: 'P',
        mediaType: 'anime',
        model: 'a/b',
        glossary: TEST_PRESET.glossary,
        temperature: 0.3,
        prompt: 'Keep honorifics.'
      },
      targetLang: 'pt-BR',
      sourceLang: 'en',
      outputDir: dir,
      confirmBeforeTranslate: true,
      memory: 'use'
    })
  })

  it("memory: 'skip' congelado no spec; valor desconhecido é INVALID_INPUT", async () => {
    const base = {
      kind: 'translate' as const,
      source: join(dir, 'ep01.en.srt'),
      presetId: 'p-test'
    }
    expect(await resolveJobSpec({ ...base, memory: 'skip' }, presets)).toMatchObject({
      memory: 'skip'
    })
    await expect(
      resolveJobSpec({ ...base, memory: 'clear' } as never, presets)
    ).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })

  it('jobMemoryMode: jobs antigos sem o campo = use; jobs que não traduzem = undefined', () => {
    expect(jobMemoryMode({ kind: 'translate' })).toBe('use')
    expect(jobMemoryMode({ kind: 'pipeline', memory: 'skip' })).toBe('skip')
    expect(jobMemoryMode({ kind: 'extract' })).toBeUndefined()
    expect(jobMemoryMode(null)).toBeUndefined()
  })

  it('pipeline: opções padrão; sourceLang explícito canonizado', async () => {
    const s = await resolveJobSpec(
      {
        kind: 'pipeline',
        source: join(dir, 'ep01.mkv'),
        trackId: 2,
        presetId: 'p-test',
        targetLang: 'es-419',
        sourceLang: 'EN',
        outputDir: join(dir, 'out')
      },
      presets
    )
    expect(s).toMatchObject({
      kind: 'pipeline',
      trackId: 2,
      targetLang: 'es-419',
      sourceLang: 'en',
      outputDir: join(dir, 'out'),
      setDefault: true,
      keepOriginalSubtitles: true,
      outputMode: 'mux',
      confirmBeforeTranslate: true,
      memory: 'use'
    })
  })

  it('pipeline: outputMode, migração de alsoWriteSubtitle e overwrite', async () => {
    const req = {
      kind: 'pipeline' as const,
      source: join(dir, 'ep01.mkv'),
      trackId: 2,
      presetId: 'p-test',
      targetLang: 'pt-BR'
    }
    const mode = async (options?: object): Promise<unknown> =>
      (await resolveJobSpec({ ...req, options }, presets)) as unknown
    expect(await mode()).toMatchObject({ outputMode: 'mux' })
    expect(await mode({ alsoWriteSubtitle: true })).toMatchObject({ outputMode: 'both' })
    expect(await mode({ alsoWriteSubtitle: false })).toMatchObject({ outputMode: 'mux' })
    expect(await mode({ outputMode: 'subtitle' })).toMatchObject({ outputMode: 'subtitle' })
    // outputMode explícito vence o legado
    expect(await mode({ outputMode: 'subtitle', alsoWriteSubtitle: true })).toMatchObject({
      outputMode: 'subtitle'
    })
    expect(await mode({ overwrite: true })).toMatchObject({ overwrite: true })
    expect(await mode()).not.toHaveProperty('overwrite')
    await expect(mode({ outputMode: 'nope' })).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })

  it('extract: nomeação padrão (player com uma faixa, track com várias) e explícita', async () => {
    const ex = (trackIds: number[], options?: object): Promise<unknown> =>
      resolveJobSpec({ kind: 'extract', source: join(dir, 'ep01.mkv'), trackIds, options }, presets)
    expect(await ex([2])).toMatchObject({ naming: 'player' })
    expect(await ex([2, 3])).toMatchObject({ naming: 'track' })
    expect(await ex([2, 3], { naming: 'player', overwrite: true })).toMatchObject({
      naming: 'player',
      overwrite: true
    })
    expect(await ex([2], { naming: 'track' })).toMatchObject({ naming: 'track' })
  })

  it('mux: idioma pré-preenchido pelo nome não é necessário quando enviado; nome padrão do idioma', async () => {
    const s = await resolveJobSpec(
      {
        kind: 'mux',
        source: join(dir, 'ep01.mkv'),
        subtitles: [{ path: join(dir, 'ep01.pt-BR.ass'), language: 'pt-br' }]
      },
      presets
    )
    expect(s).toMatchObject({
      kind: 'mux',
      subtitles: [
        {
          path: join(dir, 'ep01.pt-BR.ass'),
          language: 'pt-BR',
          name: 'Brazilian Portuguese',
          default: true
        }
      ],
      setDefault: true
    })
  })

  it('extract: ids únicos e não vazios', async () => {
    const s = await resolveJobSpec(
      { kind: 'extract', source: join(dir, 'ep01.mkv'), trackIds: [3, 2, 3] },
      presets
    )
    expect(s).toMatchObject({ kind: 'extract', trackIds: [2, 3] })
    await expect(
      resolveJobSpec({ kind: 'extract', source: join(dir, 'ep01.mkv'), trackIds: [] }, presets)
    ).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })

  it('caminhos inválidos: extensão errada para o tipo, relativo, pasta de saída inexistente', async () => {
    await expect(
      resolveJobSpec(
        { kind: 'translate', source: join(dir, 'ep01.mkv'), presetId: 'p-test' },
        presets
      )
    ).rejects.toMatchObject({ code: 'INVALID_PATH' })
    await expect(
      resolveJobSpec({ kind: 'extract', source: 'ep01.mkv', trackIds: [1] }, presets)
    ).rejects.toMatchObject({ code: 'INVALID_PATH' })
    await expect(
      resolveJobSpec(
        {
          kind: 'extract',
          source: join(dir, 'ep01.mkv'),
          trackIds: [1],
          outputDir: join(dir, 'nope')
        },
        presets
      )
    ).rejects.toMatchObject({ code: 'INVALID_PATH' })
  })

  it('forma inválida e idioma inválido => INVALID_INPUT', async () => {
    await expect(resolveJobSpec({ kind: 'nope' } as never, presets)).rejects.toMatchObject({
      code: 'INVALID_INPUT'
    })
    await expect(
      resolveJobSpec(
        {
          kind: 'translate',
          source: join(dir, 'ep01.en.srt'),
          presetId: 'p-test',
          targetLang: 'xx yy'
        },
        presets
      )
    ).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })

  it('preset padrão; modelo e idioma: pedido > padrão das configurações', async () => {
    const builtin = { ...TEST_PRESET, id: 'builtin:anime' }
    const src = testPresetSource([builtin], { model: 'm/default', targetLang: 'es' })
    const s = await resolveJobSpec({ kind: 'translate', source: join(dir, 'ep01.en.srt') }, src)
    expect(s).toMatchObject({
      preset: { id: 'builtin:anime', model: 'm/default' },
      targetLang: 'es'
    })
    const s2 = await resolveJobSpec(
      {
        kind: 'translate',
        source: join(dir, 'ep01.en.srt'),
        targetLang: 'fr',
        model: 'm/picked'
      },
      src
    )
    expect(s2).toMatchObject({ targetLang: 'fr', preset: { model: 'm/picked' } })
  })

  it('sem modelo no pedido nem padrão: MODEL_NOT_SET', async () => {
    const src = testPresetSource([TEST_PRESET], { model: null })
    await expect(
      resolveJobSpec({ kind: 'translate', source: join(dir, 'ep01.en.srt') }, src)
    ).rejects.toMatchObject({ code: 'MODEL_NOT_SET' })
  })

  it('sem idioma no pedido nem padrão: TARGET_LANG_NOT_SET', async () => {
    const src = testPresetSource([TEST_PRESET], { targetLang: null })
    await expect(
      resolveJobSpec({ kind: 'translate', source: join(dir, 'ep01.en.srt') }, src)
    ).rejects.toMatchObject({ code: 'TARGET_LANG_NOT_SET' })
    const s = await resolveJobSpec(
      { kind: 'translate', source: join(dir, 'ep01.en.srt'), targetLang: 'ja' },
      src
    )
    expect(s).toMatchObject({ targetLang: 'ja' })
  })

  it('notas do usuário: trim; vazias somem; limite de tamanho', async () => {
    const s = await resolveJobSpec(
      { kind: 'translate', source: join(dir, 'ep01.en.srt'), notes: '  Naruto, ep. 3  ' },
      presets
    )
    expect(s).toMatchObject({ notes: 'Naruto, ep. 3' })
    const s2 = await resolveJobSpec(
      { kind: 'translate', source: join(dir, 'ep01.en.srt'), notes: '  ' },
      presets
    )
    expect('notes' in s2).toBe(false)
    await expect(
      resolveJobSpec(
        { kind: 'translate', source: join(dir, 'ep01.en.srt'), notes: 'x'.repeat(4001) },
        presets
      )
    ).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })

  it('media context: the Main looks the details up and freezes them in the spec', async () => {
    const lookup = vi.fn(async () => ({
      details: {
        provider: 'tmdb' as const,
        id: 7,
        kind: 'series' as const,
        title: 'Frieren',
        genres: [],
        synopsis: 'syn'
      },
      episode: { season: 1, episode: 5, title: 'E5' }
    }))
    const s = await resolveJobSpec(
      {
        kind: 'translate',
        source: join(dir, 'ep01.en.srt'),
        media: { ref: { provider: 'tmdb', id: 7, kind: 'series' }, season: 1, episode: 5 }
      },
      presets,
      { lookup }
    )
    expect(lookup).toHaveBeenCalledWith(
      { provider: 'tmdb', id: 7, kind: 'series' },
      { season: 1, episode: 5 }
    )
    expect(s).toMatchObject({
      media: { details: { id: 7, title: 'Frieren' }, episode: { episode: 5 } }
    })
  })

  it('media lookup failure: no details, the reason is frozen and the error is reported', async () => {
    const onError = vi.fn()
    const s = await resolveJobSpec(
      {
        kind: 'translate',
        source: join(dir, 'ep01.en.srt'),
        media: { ref: { provider: 'tmdb', id: 7, kind: 'movie' } }
      },
      presets,
      {
        lookup: async () => Promise.reject(new AppError('TMDB_KEY_MISSING', 'no key')),
        onError
      }
    )
    expect(s).not.toHaveProperty('media')
    expect(s).toMatchObject({ mediaUnavailable: 'TMDB lookup failed (TMDB_KEY_MISSING: no key)' })
    expect(onError).toHaveBeenCalled()
  })

  it('media requested without a lookup service: the reason is frozen too', async () => {
    const s = await resolveJobSpec(
      {
        kind: 'pipeline',
        source: join(dir, 'ep01.mkv'),
        trackId: 2,
        media: { ref: { provider: 'tmdb', id: 7, kind: 'movie' } }
      },
      presets
    )
    expect(s).toMatchObject({ mediaUnavailable: expect.stringContaining('not available') })
  })

  it('no media requested: no details and no reason', async () => {
    const s = await resolveJobSpec({ kind: 'translate', source: join(dir, 'ep01.en.srt') }, presets)
    expect(s).not.toHaveProperty('media')
    expect(s).not.toHaveProperty('mediaUnavailable')
  })

  it('media ref is validated', async () => {
    await expect(
      resolveJobSpec(
        {
          kind: 'translate',
          source: join(dir, 'ep01.en.srt'),
          media: { ref: { provider: 'imdb', id: -1, kind: 'movie' } } as never
        },
        presets
      )
    ).rejects.toMatchObject({ code: 'INVALID_INPUT' })
  })
})

describe('resolveJobSpec: pipeline with an external subtitle', () => {
  const base = { kind: 'pipeline' as const, presetId: 'p-test' }

  it('freezes path, size and format; sourceLang from the file name; no trackId', async () => {
    const s = await resolveJobSpec(
      { ...base, source: join(dir, 'ep01.mkv'), subtitle: { path: join(dir, 'ep01.en.srt') } },
      presets
    )
    expect(s).toMatchObject({
      kind: 'pipeline',
      source: join(dir, 'ep01.mkv'),
      externalSubtitle: {
        path: join(dir, 'ep01.en.srt'),
        size: 3,
        format: 'srt',
        includeInMux: false
      },
      sourceLang: 'en'
    })
    expect(s).not.toHaveProperty('trackId')
  })

  it('the given language beats the file name; the request sourceLang beats both', async () => {
    const source = join(dir, 'ep01.mkv')
    const path = join(dir, 'ep01.en.srt')
    expect(
      await resolveJobSpec({ ...base, source, subtitle: { path, language: 'ja' } }, presets)
    ).toMatchObject({ sourceLang: 'ja' })
    expect(
      await resolveJobSpec(
        { ...base, source, sourceLang: 'ko', subtitle: { path, language: 'ja' } },
        presets
      )
    ).toMatchObject({ sourceLang: 'ko' })
  })

  it('includeSourceSubtitle only applies when an MKV is written', async () => {
    const req = {
      ...base,
      source: join(dir, 'ep01.mkv'),
      subtitle: { path: join(dir, 'ep01.pt-BR.ass') }
    }
    expect(
      await resolveJobSpec({ ...req, options: { includeSourceSubtitle: true } }, presets)
    ).toMatchObject({ externalSubtitle: { format: 'ass', includeInMux: true } })
    expect(
      await resolveJobSpec(
        { ...req, options: { includeSourceSubtitle: true, outputMode: 'subtitle' } },
        presets
      )
    ).toMatchObject({ externalSubtitle: { includeInMux: false } })
  })

  it('requires exactly one of trackId and subtitle', async () => {
    const source = join(dir, 'ep01.mkv')
    const error = {
      code: 'INVALID_INPUT',
      message: 'Exactly one of trackId and subtitle is required'
    }
    // Zero origens.
    await expect(resolveJobSpec({ ...base, source }, presets)).rejects.toMatchObject(error)
    // Duas origens.
    await expect(
      resolveJobSpec(
        { ...base, source, trackId: 2, subtitle: { path: join(dir, 'ep01.en.srt') } },
        presets
      )
    ).rejects.toMatchObject(error)
    // Duas origens com legenda inexistente: a checagem vem antes de tocar no disco.
    await expect(
      resolveJobSpec(
        { ...base, source, trackId: 2, subtitle: { path: join(dir, 'nope.srt') } },
        presets
      )
    ).rejects.toMatchObject(error)
  })

  it('trackId 0 is a valid single origin (no truthiness check)', async () => {
    const spec = await resolveJobSpec(
      { ...base, source: join(dir, 'ep01.mkv'), trackId: 0 },
      presets
    )
    expect(spec).toMatchObject({ kind: 'pipeline', trackId: 0 })
    expect(spec).not.toHaveProperty('externalSubtitle')
  })

  it('the external subtitle goes through the same path guard (extension, existence)', async () => {
    const source = join(dir, 'ep01.mkv')
    await expect(
      resolveJobSpec({ ...base, source, subtitle: { path: source } }, presets)
    ).rejects.toMatchObject({ code: 'INVALID_PATH', message: 'extension' })
    await expect(
      resolveJobSpec({ ...base, source, subtitle: { path: join(dir, 'nope.srt') } }, presets)
    ).rejects.toMatchObject({ code: 'INVALID_PATH', message: 'not-found' })
  })
})
