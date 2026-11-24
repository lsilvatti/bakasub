import { describe, expect, it } from 'vitest'
import type { JobSpec } from './job-spec'
import { pipelineOutputMode } from './job-spec'
import { planFor } from './plan'

const pipeline = (over: object): JobSpec => ({ kind: 'pipeline', ...over }) as JobSpec

describe('planFor(pipeline)', () => {
  const head = ['probe', 'extract', 'preflight', 'translate']

  it('mux: termina no mux, sem write', () => {
    expect(planFor(pipeline({ outputMode: 'mux' }))).toEqual([...head, 'mux'])
  })
  it('subtitle: write e nenhum mux (o MKV não é tocado)', () => {
    expect(planFor(pipeline({ outputMode: 'subtitle' }))).toEqual([...head, 'write'])
  })
  it('both: write e depois mux', () => {
    expect(planFor(pipeline({ outputMode: 'both' }))).toEqual([...head, 'write', 'mux'])
  })
  it('jobs persistidos antes de outputMode: alsoWriteSubtitle continua valendo', () => {
    expect(planFor(pipeline({ alsoWriteSubtitle: true }))).toEqual([...head, 'write', 'mux'])
    expect(planFor(pipeline({ alsoWriteSubtitle: false }))).toEqual([...head, 'mux'])
    expect(planFor(pipeline({}))).toEqual([...head, 'mux'])
    expect(pipelineOutputMode({ alsoWriteSubtitle: true })).toBe('both')
    expect(pipelineOutputMode({ outputMode: 'subtitle', alsoWriteSubtitle: true })).toBe('subtitle')
  })
  it('external subtitle: nothing to extract, in all three output modes', () => {
    const ext = {
      externalSubtitle: { path: '/s.srt', size: 1, format: 'srt', includeInMux: false }
    }
    const start = ['probe', 'preflight', 'translate']
    expect(planFor(pipeline({ ...ext, outputMode: 'mux' }))).toEqual([...start, 'mux'])
    expect(planFor(pipeline({ ...ext, outputMode: 'subtitle' }))).toEqual([...start, 'write'])
    expect(planFor(pipeline({ ...ext, outputMode: 'both' }))).toEqual([...start, 'write', 'mux'])
  })
  it('outros tipos não mudam', () => {
    expect(planFor({ kind: 'extract' } as JobSpec)).toEqual(['probe', 'extract', 'write'])
    expect(planFor({ kind: 'translate' } as JobSpec)).toEqual(['preflight', 'translate', 'write'])
    expect(planFor({ kind: 'mux' } as JobSpec)).toEqual(['probe', 'validate', 'mux'])
  })
})
