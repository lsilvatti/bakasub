import { describe, expect, it } from 'vitest'
import { FPS_TOLERANCE, fpsDiverges, subtitleSyncRisks } from './subtitle-sync'

describe('fpsDiverges', () => {
  it('fps iguais ou dentro da tolerância não acusam (23.976 ≈ 23.98)', () => {
    expect(fpsDiverges(23.976, 23.976)).toBe(false)
    expect(fpsDiverges(23.98, 23.976)).toBe(false)
    expect(fpsDiverges(23.976, 23.976 + FPS_TOLERANCE)).toBe(false)
    expect(fpsDiverges(1e9 / 41708333, 23.976)).toBe(false)
  })
  it('fps diferentes acusam (25 vs 23.976, 30 vs 29.97)', () => {
    expect(fpsDiverges(25, 23.976)).toBe(true)
    expect(fpsDiverges(30, 29.97)).toBe(true)
    expect(fpsDiverges(23.976, 23.976 + FPS_TOLERANCE * 2)).toBe(true)
  })
  it('fps ausente ou inválido em qualquer lado nunca é falso positivo', () => {
    expect(fpsDiverges(undefined, 23.976)).toBe(false)
    expect(fpsDiverges(25, undefined)).toBe(false)
    expect(fpsDiverges(null, null)).toBe(false)
    expect(fpsDiverges(0, 25)).toBe(false)
    expect(fpsDiverges(25, Number.NaN)).toBe(false)
    expect(fpsDiverges(-1, 25)).toBe(false)
  })
  it('tolerância configurável', () => {
    expect(fpsDiverges(25, 25.5, 1)).toBe(false)
    expect(fpsDiverges(25, 25.5, 0.1)).toBe(true)
  })
})

describe('subtitleSyncRisks', () => {
  it('hashMatch e fps iguais: nada a avisar', () => {
    expect(subtitleSyncRisks({ hashMatch: true, subtitleFps: 23.976, videoFps: 23.976 })).toEqual(
      []
    )
  })
  it('sem hashMatch: avisa mesmo sem fps', () => {
    expect(subtitleSyncRisks({ hashMatch: false })).toEqual(['no-hash-match'])
  })
  it('fps divergente vem antes do hash', () => {
    expect(subtitleSyncRisks({ hashMatch: false, subtitleFps: 25, videoFps: 23.976 })).toEqual([
      'fps-mismatch',
      'no-hash-match'
    ])
  })
  it('fps divergente acusa mesmo com hashMatch', () => {
    expect(subtitleSyncRisks({ hashMatch: true, subtitleFps: 25, videoFps: 23.976 })).toEqual([
      'fps-mismatch'
    ])
  })
  it('hashMatch sem fps do vídeo: nada a avisar', () => {
    expect(subtitleSyncRisks({ hashMatch: true, subtitleFps: 25 })).toEqual([])
  })
})
