import { describe, expect, it } from 'vitest'
import { buildExtractArgs, buildMuxArgs, interpretExit, parseGuiLine } from './mkvtoolnix'
import { parseIdentify } from './mkv-probe'
import { toBcp47 } from './languages'

/** Forma de `mkvmerge -J` (campos usados + ruído). Validar contra o binário real. */
const IDENTIFY = {
  container: { recognized: true, supported: true, type: 'Matroska' },
  tracks: [
    {
      id: 0,
      type: 'video',
      codec: 'AVC',
      properties: { codec_id: 'V_MPEG4/ISO/AVC', default_duration: 41708333 }
    },
    {
      id: 1,
      type: 'audio',
      codec: 'AAC',
      properties: { codec_id: 'A_AAC', language: 'jpn', language_ietf: 'ja' }
    },
    {
      id: 2,
      type: 'subtitles',
      codec: 'SubStationAlpha',
      properties: {
        codec_id: 'S_TEXT/ASS',
        language: 'eng',
        track_name: 'Full Subs',
        default_track: true,
        forced_track: false,
        num_index_entries: 412
      }
    },
    {
      id: 3,
      type: 'subtitles',
      codec: 'SubRip/SRT',
      properties: {
        codec_id: 'S_TEXT/UTF8',
        language: 'por',
        language_ietf: 'pt-BR',
        forced_track: true
      }
    },
    {
      id: 4,
      type: 'subtitles',
      codec: 'HDMV PGS',
      properties: { codec_id: 'S_HDMV/PGS', language: 'und' }
    },
    { id: 5, type: 'subtitles', codec: 'WebVTT', properties: { codec_id: 'S_TEXT/WEBVTT' } },
    {
      id: 6,
      type: 'subtitles',
      codec: 'SSA',
      properties: { codec_id: 'S_TEXT/SSA', language: 'xyz' }
    },
    { id: 7, type: 'subtitles', codec: 'USF', properties: { codec_id: 'S_TEXT/USF' } }
  ],
  attachments: [{ id: 1, file_name: 'font.ttf' }]
}

describe('toBcp47', () => {
  it('ISO 639-2 (B e T) para BCP-47; und/desconhecido => undefined', () => {
    expect(toBcp47('eng')).toBe('en')
    expect(toBcp47('jpn')).toBe('ja')
    expect(toBcp47('por')).toBe('pt')
    expect(toBcp47('fre')).toBe('fr')
    expect(toBcp47('fra')).toBe('fr')
    expect(toBcp47('und')).toBeUndefined()
    expect(toBcp47('xyz')).toBeUndefined()
  })
})

describe('parseIdentify', () => {
  const info = parseIdentify(IDENTIFY)
  it('só faixas de legenda, com formato, idioma, flags e eventos', () => {
    expect(info.subtitles.map((t) => t.id)).toEqual([2, 3, 4, 5, 6, 7])
    expect(info.subtitles[0]).toEqual({
      id: 2,
      codecId: 'S_TEXT/ASS',
      format: 'ass',
      language: 'en',
      name: 'Full Subs',
      default: true,
      forced: false,
      eventCount: 412
    })
    expect(info.subtitles[1]).toMatchObject({ format: 'srt', language: 'pt-BR', forced: true })
  })
  it('WebVTT é suportado; PGS e texto sem parser ficam desabilitados com motivo', () => {
    expect(info.subtitles[2]).toMatchObject({ format: null, unsupported: 'image-based' })
    expect(info.subtitles[3]).toMatchObject({ format: 'vtt', codecId: 'S_TEXT/WEBVTT' })
    expect(info.subtitles[3]).not.toHaveProperty('unsupported')
    expect(info.subtitles[4]).toMatchObject({ format: 'ssa', language: undefined })
    expect(info.subtitles[5]).toMatchObject({ format: null, unsupported: 'no-parser' })
  })
  it('anexos contados (fontes do ASS)', () => {
    expect(info.attachments).toBe(1)
  })
  it('fps do vídeo vem do default_duration da primeira faixa de vídeo', () => {
    expect(info.videoFps).toBe(23.976)
    const pal = { tracks: [{ id: 0, type: 'video', properties: { default_duration: 40000000 } }] }
    expect(parseIdentify(pal).videoFps).toBe(25)
  })
  it('sem default_duration (ou inválido) não inventa fps', () => {
    expect(parseIdentify({ tracks: [{ id: 0, type: 'video' }] })).not.toHaveProperty('videoFps')
    const zero = { tracks: [{ id: 0, type: 'video', properties: { default_duration: 0 } }] }
    expect(parseIdentify(zero)).not.toHaveProperty('videoFps')
  })
  it('container não reconhecido lança', () => {
    expect(() => parseIdentify({ container: { recognized: false }, tracks: [] })).toThrow()
  })
})

describe('parseGuiLine', () => {
  it('progresso, aviso, erro; o resto é null', () => {
    expect(parseGuiLine('#GUI#progress 42%')).toEqual({ type: 'progress', percent: 42 })
    expect(parseGuiLine('#GUI#warning Something odd')).toEqual({
      type: 'warning',
      message: 'Something odd'
    })
    expect(parseGuiLine('#GUI#error Bad thing')).toEqual({ type: 'error', message: 'Bad thing' })
    expect(parseGuiLine('Multiplexing took 3 seconds.')).toBeNull()
  })
})

describe('interpretExit', () => {
  it('0 = ok, 1 = avisos (sucesso), outros = erro', () => {
    expect(interpretExit(0)).toBe('ok')
    expect(interpretExit(1)).toBe('warnings')
    expect(interpretExit(2)).toBe('error')
    expect(interpretExit(null)).toBe('error')
  })
})

describe('buildMuxArgs', () => {
  const base = {
    source: '/v/ep.mkv',
    output: '/out/.ep.pt-BR.bakasub-partial.mkv',
    subtitles: [
      {
        path: '/w/sub.ass',
        language: 'pt-BR',
        name: 'Brazilian Portuguese (BakaSub)',
        default: true
      }
    ],
    sourceSubtitleTrackIds: [2, 3]
  }
  it('mantém originais tirando o default delas; nova faixa default com idioma e nome', () => {
    expect(buildMuxArgs({ ...base, setDefault: true, keepOriginalSubtitles: true })).toEqual([
      '--gui-mode',
      '-o',
      '/out/.ep.pt-BR.bakasub-partial.mkv',
      '--default-track-flag',
      '2:no',
      '--default-track-flag',
      '3:no',
      '/v/ep.mkv',
      '--language',
      '0:pt-BR',
      '--track-name',
      '0:Brazilian Portuguese (BakaSub)',
      '--default-track-flag',
      '0:yes',
      '/w/sub.ass'
    ])
  })
  it('sem setDefault não mexe nas originais; sem keep remove as legendas originais', () => {
    const a = buildMuxArgs({
      ...base,
      subtitles: [{ ...base.subtitles[0], default: false }],
      setDefault: false,
      keepOriginalSubtitles: false
    })
    expect(a).toContain('--no-subtitles')
    expect(a).not.toContain('2:no')
    expect(a).toContain('0:no')
  })
})

describe('buildExtractArgs', () => {
  it('uma chamada com todas as faixas', () => {
    expect(
      buildExtractArgs('/v/ep.mkv', [
        { trackId: 2, output: '/w/2.ass' },
        { trackId: 3, output: '/w/3.srt' }
      ])
    ).toEqual(['--gui-mode', '/v/ep.mkv', 'tracks', '2:/w/2.ass', '3:/w/3.srt'])
  })
})
