/**
 * Risco de dessincronia de uma legenda externa em relação ao vídeo, avaliado ANTES da tradução
 * (que é paga). Só aponta riscos: nunca bloqueia e nunca corrige o timing.
 */

/** Diferença máxima, em quadros por segundo, para considerar dois fps iguais (23.976 ≈ 23.98). */
export const FPS_TOLERANCE = 0.01

/**
 * - `fps-mismatch`: os dois lados informam fps e eles divergem além da tolerância (deriva certa).
 * - `no-hash-match`: o OpenSubtitles não confirmou pelo hash que a legenda é deste arquivo exato.
 */
export type SyncRisk = 'fps-mismatch' | 'no-hash-match'

export interface SyncInput {
  /** `moviehash_match` do OpenSubtitles. */
  readonly hashMatch: boolean
  /** fps declarado no resultado do OpenSubtitles. */
  readonly subtitleFps?: number | null
  /** fps da primeira faixa de vídeo (probe). */
  readonly videoFps?: number | null
}

const known = (fps: number | null | undefined): fps is number =>
  typeof fps === 'number' && Number.isFinite(fps) && fps > 0

/** `true` só quando os DOIS fps são conhecidos e divergem; fps ausente/inválido nunca acusa. */
export function fpsDiverges(
  subtitleFps: number | null | undefined,
  videoFps: number | null | undefined,
  tolerance = FPS_TOLERANCE
): boolean {
  if (!known(subtitleFps) || !known(videoFps)) return false
  // Arredonda para não acusar ruído de ponto flutuante (1e9 / 41708333 = 23.97602…).
  return Math.round(Math.abs(subtitleFps - videoFps) * 1e6) / 1e6 > tolerance
}

/** Riscos em ordem de gravidade (fps primeiro). Lista vazia = nada a avisar. */
export function subtitleSyncRisks(input: SyncInput): SyncRisk[] {
  const risks: SyncRisk[] = []
  if (fpsDiverges(input.subtitleFps, input.videoFps)) risks.push('fps-mismatch')
  if (!input.hashMatch) risks.push('no-hash-match')
  return risks
}
