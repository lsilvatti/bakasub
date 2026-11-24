import type { JobStep } from '../../shared/ipc-contract'
import { pipelineOutputMode, type JobSpec } from './job-spec'
import type { Step } from './step-context'
import { extractStep } from './steps/extract.step'
import { muxStep } from './steps/mux.step'
import { preflightStep } from './steps/preflight.step'
import { probeStep } from './steps/probe.step'
import { translateStep } from './steps/translate.step'
import { validateStep } from './steps/validate.step'
import { writeStep } from './steps/write.step'

export const STEPS: Readonly<Record<JobStep, Step>> = {
  probe: probeStep,
  extract: extractStep,
  preflight: preflightStep,
  translate: translateStep,
  validate: validateStep,
  write: writeStep,
  mux: muxStep
}

/** Etapas de cada tipo de job, em ordem. */
export function planFor(spec: JobSpec): JobStep[] {
  switch (spec.kind) {
    case 'extract':
      return ['probe', 'extract', 'write']
    case 'translate':
      return ['preflight', 'translate', 'write']
    case 'mux':
      return ['probe', 'validate', 'mux']
    case 'pipeline': {
      const mode = pipelineOutputMode(spec)
      return [
        'probe',
        ...(spec.externalSubtitle ? [] : (['extract'] as const)),
        'preflight',
        'translate',
        ...(mode === 'mux' ? [] : (['write'] as const)),
        ...(mode === 'subtitle' ? [] : (['mux'] as const))
      ]
    }
  }
}
