import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import type { Job, JobMemoryMode } from '@shared/ipc-contract'
import { Badge } from '@renderer/components/ui/badge'
import { Button } from '@renderer/components/ui/button'
import { Tooltip } from '@renderer/components/ui/popover'
import {
  ACTIVE_STATUSES,
  isTerminal,
  RETRYABLE_STATUSES,
  useJob,
  useJobAction,
  usePreflight,
  useRemoveJob,
  useSetJobMemory,
  type JobAction
} from '@renderer/data/jobs'
import { errorMessage } from '@renderer/i18n/error-message'
import { episodeCode } from '@renderer/lib/episode'
import { JobProgress } from './job-progress'
import { JobResult } from './job-result'
import { JobSetup } from './job-setup'
import { MemoryOption } from './memory-option'
import { PreflightReview } from './preflight-review'
import { PromptSummary } from './prompt-summary'

/** The "do it again" label of the wizard that created the job. */
const ANOTHER_LABEL = {
  extract: 'extractPage.another',
  mux: 'muxPage.another',
  translate: 'translate.another',
  pipeline: 'translate.another'
} as const satisfies Record<Job['kind'], string>

const when = (iso: string, lang: string): string => {
  const d = new Date(iso)
  return Number.isNaN(d.getTime())
    ? ''
    : new Intl.DateTimeFormat(lang, { dateStyle: 'short', timeStyle: 'short' }).format(d)
}

/**
 * One job, whatever its state: preflight review (confirm/cancel), live progress, or the result.
 * Used by the wizard (after creating) and by the queue inspector.
 */
export function JobDetail({
  id,
  onAnother,
  onRemoved,
  showQueueLink = false,
  inspectPath,
  preflightNotice
}: {
  id: string
  /** Wizard: "Translate/Extract/Mux another…" once the job is over (label follows `job.kind`). */
  onAnother?: () => void
  onRemoved?: () => void
  showQueueLink?: boolean
  /** Wizard: the subtitle file behind the job, for the encoding and parser check. */
  inspectPath?: string
  /** Wizard: a non-blocking notice shown in the review (e.g. a subtitle that may be out of sync). */
  preflightNotice?: React.ReactNode
}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const query = useJob(id)
  const job = query.data
  const waiting = job?.status === 'awaiting-confirmation'
  const translates = job?.kind === 'translate' || job?.kind === 'pipeline'
  const active = !!job && ACTIVE_STATUSES.has(job.status)
  // Na revisão e no fim: o que foi (ou vai ser) enviado no System Prompt.
  const preflight = usePreflight(id, waiting || (translates && !active))
  const action = useJobAction()
  const remove = useRemoveJob()
  const setMemory = useSetJobMemory()
  const [retryChoice, setRetryMemory] = useState<JobMemoryMode | null>(null)
  const act = (a: JobAction, memory?: JobMemoryMode): void =>
    action.mutate(
      { id, action: a, memory },
      {
        onSuccess: () => {
          if (a === 'retry') setRetryMemory(null)
        },
        onError: (err) => toast.error(errorMessage(err))
      }
    )
  const changeMemory = (memory: JobMemoryMode): void =>
    setMemory.mutate({ id, memory }, { onError: (err) => toast.error(errorMessage(err)) })

  if (!job && query.isError) {
    const leave = onRemoved ?? onAnother
    return (
      <div className="grid gap-3" role="alert">
        <p className="text-sm text-muted-foreground">{errorMessage(query.error)}</p>
        {leave ? (
          <div>
            <Button variant="outline" onClick={leave}>
              {t('actions.back')}
            </Button>
          </div>
        ) : null}
      </div>
    )
  }
  if (!job) return <p className="text-sm text-muted-foreground">{t('translate.starting')}</p>
  const busy = action.isPending || remove.isPending || setMemory.isPending
  // Retry: starts from the job's original choice; the user may switch it before retrying.
  const retryable = translates && RETRYABLE_STATUSES.has(job.status)
  const retryMemory = retryChoice ?? job.memory ?? 'use'

  return (
    <div className="grid gap-4">
      <header className="grid gap-1">
        {job.media ? (
          <p className="text-sm text-violet">
            {job.media.title}
            {job.media.year ? ` (${job.media.year})` : ''}
            {job.media.episode
              ? ` · ${episodeCode(job.media.episode.season, job.media.episode.episode)}${job.media.episode.title ? `: ${job.media.episode.title}` : ''}`
              : ''}
          </p>
        ) : null}
        <h2 className="text-base font-semibold break-all">{job.source}</h2>
        <p className="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-muted-foreground">
          <Badge tone="violet">{t(`jobs.kind.${job.kind}`)}</Badge>
          {job.batchId ? <Badge>{t('jobs.batch')}</Badge> : null}
          {job.memory === 'skip' ? <Badge tone="primary">{t('jobs.memory.skipped')}</Badge> : null}
          <span>{t(`jobs.status.${job.status}`)}</span>
          {job.attempt > 1 ? <span>{t('jobs.attempt', { n: job.attempt })}</span> : null}
          <span>{t('jobs.created', { when: when(job.createdAt, i18n.language) })}</span>
          {job.finishedAt ? (
            <span>{t('jobs.finished', { when: when(job.finishedAt, i18n.language) })}</span>
          ) : null}
        </p>
        {translates ? (
          <div className="mt-2">
            <JobSetup setup={job.translation ?? {}} />
          </div>
        ) : null}
      </header>

      {waiting ? (
        preflight.data ? (
          <>
            <PreflightReview
              estimate={preflight.data.estimate}
              coverage={preflight.data.coverage}
              inspectPath={inspectPath}
              notice={preflightNotice}
              prompt={preflight.data.prompt ?? null}
            />
            <MemoryOption value={job.memory ?? 'use'} onChange={changeMemory} disabled={busy} />
          </>
        ) : (
          <p className="text-sm text-muted-foreground">{t('translate.starting')}</p>
        )
      ) : active ? (
        <JobProgress job={job} />
      ) : (
        <>
          <JobResult job={job} />
          {translates && preflight.data ? (
            <PromptSummary prompt={preflight.data.prompt ?? null} />
          ) : null}
        </>
      )}

      {retryable ? (
        <MemoryOption value={retryMemory} onChange={setRetryMemory} disabled={busy} />
      ) : null}

      <JobActions
        job={job}
        busy={busy}
        canConfirm={!!preflight.data}
        // Untouched = omitted: the Main keeps the job's original choice.
        onAct={(a) => act(a, a === 'retry' ? (retryChoice ?? undefined) : undefined)}
        onRemove={() =>
          remove.mutate(id, {
            onSuccess: onRemoved,
            onError: (err) => toast.error(errorMessage(err))
          })
        }
        onAnother={onAnother}
        showQueueLink={showQueueLink}
      />
    </div>
  )
}

function JobActions({
  job,
  busy,
  canConfirm,
  onAct,
  onRemove,
  onAnother,
  showQueueLink
}: {
  job: Job
  busy: boolean
  canConfirm: boolean
  onAct: (a: JobAction) => void
  onRemove: () => void
  onAnother?: () => void
  showQueueLink: boolean
}): React.JSX.Element {
  const { t } = useTranslation()
  const waiting = job.status === 'awaiting-confirmation'
  const active = ACTIVE_STATUSES.has(job.status)
  return (
    <div className="flex flex-wrap items-center justify-end gap-2 border-t pt-4">
      <Button variant="ghost" asChild className={showQueueLink ? '' : 'mr-auto'}>
        <Link to="/logs" search={{ job: job.id }}>
          {t('jobs.logs')}
        </Link>
      </Button>
      {showQueueLink ? (
        <Button variant="ghost" asChild className="mr-auto">
          <Link to="/queue" search={{ job: job.id }}>
            {t('jobs.openQueue')}
          </Link>
        </Button>
      ) : null}
      {waiting || active ? (
        <Tooltip content={t('jobs.cancelHint')}>
          <Button variant="destructive" disabled={busy} onClick={() => onAct('cancel')}>
            {t('jobs.cancel')}
          </Button>
        </Tooltip>
      ) : null}
      {isTerminal(job.status) && !onAnother ? (
        <Button variant="ghost" disabled={busy} onClick={onRemove}>
          {t('jobs.remove')}
        </Button>
      ) : null}
      {RETRYABLE_STATUSES.has(job.status) ? (
        <Button disabled={busy} onClick={() => onAct('retry')}>
          {t('jobs.retry')}
        </Button>
      ) : null}
      {waiting ? (
        <Button variant="default" disabled={busy || !canConfirm} onClick={() => onAct('confirm')}>
          {t('jobs.confirm')}
        </Button>
      ) : null}
      {onAnother && isTerminal(job.status) ? (
        <Button variant="default" onClick={onAnother}>
          {t(ANOTHER_LABEL[job.kind])}
        </Button>
      ) : null}
    </div>
  )
}
