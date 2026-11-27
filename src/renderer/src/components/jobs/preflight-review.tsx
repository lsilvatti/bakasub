import { AlertTriangle, Coins, Info } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type {
  CoverageReport,
  PreflightEstimate,
  PromptSummary as Prompt
} from '@shared/ipc-contract'
import { formatDate, formatInt, formatUsd } from '@renderer/lib/format'
import { CoverageList } from './coverage-list'
import { CoverageSummary, FileCheck } from './file-check'
import { PromptSummary } from './prompt-summary'

/** The numbers the user confirms: lines, requests, tokens, cost, free quota and coverage. */
export function PreflightReview({
  estimate: e,
  coverage,
  inspectPath,
  prompt,
  notice
}: {
  estimate: PreflightEstimate
  coverage: CoverageReport
  /** What goes into the system prompt; `null` for jobs created before it was recorded. */
  prompt?: Prompt | null
  /** Subtitle file to check for encoding and parser warnings (not available for videos). */
  inspectPath?: string
  /** Non-blocking notice from the wizard, shown above the numbers. */
  notice?: React.ReactNode
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const n = (v: number): string => formatInt(v, i18n.language)
  const free =
    e.model.pricing !== null && e.model.pricing.prompt === 0 && e.model.pricing.completion === 0
  const q = e.freeQuota

  return (
    <div className="grid gap-4">
      {notice}
      {e.lines.toTranslate === 0 ? (
        <p className="flex items-center gap-2 text-sm text-success">
          <Info className="size-4" />
          {t('preflight.nothing')}
        </p>
      ) : (
        <dl className="grid gap-3 sm:grid-cols-3">
          <div className="grid gap-0.5 rounded-lg border bg-background/40 p-3">
            <dt className="text-xs text-muted-foreground">{t('preflight.lines')}</dt>
            <dd className="text-lg font-semibold">
              {t('preflight.toTranslate', { count: e.lines.toTranslate })}
            </dd>
            <dd className="text-xs text-muted-foreground">
              {t('preflight.total', { count: e.lines.total })}
              {e.lines.cached > 0
                ? ` · ${t('preflight.cachedLines', { count: e.lines.cached })}`
                : ''}
              {e.lines.copies > 0 ? ` · ${t('preflight.copies', { count: e.lines.copies })}` : ''}
              {e.memory === 'skip' ? ` · ${t('preflight.cacheSkipped')}` : ''}
            </dd>
          </div>
          <div className="grid gap-0.5 rounded-lg border bg-background/40 p-3">
            <dt className="text-xs text-muted-foreground">{t('preflight.tokens')}</dt>
            <dd className="text-sm">
              {t('preflight.tokensIn', {
                low: n(e.tokens.input.low),
                high: n(e.tokens.input.high)
              })}
            </dd>
            <dd className="text-sm">
              {t('preflight.tokensOut', {
                low: n(e.tokens.output.low),
                high: n(e.tokens.output.high)
              })}
            </dd>
            <dd className="text-xs text-muted-foreground">
              {t('preflight.parts', { count: e.lines.chunks })}
            </dd>
          </div>
          <div className="grid gap-0.5 rounded-lg border border-primary/40 bg-primary-soft/40 p-3">
            <dt className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <Coins className="size-3.5" />
              {t('preflight.cost')}
            </dt>
            <dd className="text-lg font-semibold text-primary">
              {free
                ? t('preflight.free')
                : e.cost
                  ? t('preflight.costRange', {
                      low: formatUsd(e.cost.low),
                      high: formatUsd(e.cost.high)
                    })
                  : '—'}
            </dd>
            {!e.cost && !free ? (
              <dd className="text-xs text-muted-foreground">{t('preflight.dynamic')}</dd>
            ) : null}
          </div>
        </dl>
      )}

      {q ? (
        q.enough === false ? (
          <p
            role="alert"
            className="flex gap-2 rounded-md border border-warning/40 bg-warning/10 p-3 text-sm text-warning"
          >
            <AlertTriangle className="size-4 shrink-0" />
            {t('preflight.quotaShort', { requests: q.requests, remaining: q.remaining })}
          </p>
        ) : q.enough === null ? (
          <p className="text-xs text-muted-foreground">{t('preflight.quotaUnknown')}</p>
        ) : (
          <p className="text-xs text-muted-foreground">
            {t('preflight.quotaOk', { remaining: q.remaining, limit: q.limit })}
          </p>
        )
      ) : null}

      {prompt !== undefined ? <PromptSummary prompt={prompt} /> : null}
      {inspectPath ? <FileCheck path={inspectPath} /> : null}
      <CoverageSummary report={coverage} />
      <CoverageList report={coverage} />

      <p className="text-xs text-muted-foreground">
        {t('preflight.model', {
          model: e.model.id,
          date: e.model.catalogFetchedAt ? formatDate(e.model.catalogFetchedAt, i18n.language) : '—'
        })}
        {e.approximate ? ` · ${t('preflight.approximate')}` : ''}
      </p>
    </div>
  )
}
