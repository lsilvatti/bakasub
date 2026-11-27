import { useCallback } from 'react'
import { Info } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { FileDrop } from '@renderer/components/files/file-drop'
import { JobDetail } from '@renderer/components/jobs/job-detail'
import { Backdrop } from '@renderer/components/media/backdrop'
import type { MediaChoice } from '@renderer/components/media/media-choice'
import {
  ExternalSubtitle,
  type ExternalSubtitleValue
} from '@renderer/components/media/external-subtitle'
import { SyncWarning } from '@renderer/components/media/sync-warning'
import { MediaPanel } from '@renderer/components/media/media-panel'
import { TrackPicker } from '@renderer/components/media/track-picker'
import { Button } from '@renderer/components/ui/button'
import { Dialog, DialogClose, DialogContent, DialogFooter } from '@renderer/components/ui/dialog'
import { useCreateJob, useJob } from '@renderer/data/jobs'
import { useMediaLookup, useProbe } from '@renderer/data/media'
import { usePresets } from '@renderer/data/presets'
import { useOutputPrefs } from '@renderer/data/output-prefs'
import { useSettings } from '@renderer/data/settings'
import { errorMessage } from '@renderer/i18n/error-message'
import { existingTarget } from '@renderer/lib/existing-subtitle'
import { languageName } from '@renderer/lib/languages'
import { outputFolder, translateOutputNames } from '@renderer/lib/output-names'
import { usePageReducer, usePageState } from '@renderer/lib/page-state'
import { presetLabel } from '@renderer/lib/preset-labels'
import { suggestPreset } from '@renderer/lib/preset-suggestion'
import { syncCheck, useSyncMessages } from '@renderer/lib/sync-check'
import { MemoryOption } from '@renderer/components/jobs/memory-option'
import { OutputStep } from './output-step'
import { PresetStep } from './preset-step'
import { suggestionHint } from './suggestion-hint'
import { StepBar } from './step-bar'
import { WizardSummary } from './wizard-summary'
import {
  buildJobRequest,
  canContinue,
  canGo,
  initialWizard,
  stepsFor,
  stepValid,
  wizardReducer,
  type WizardStep
} from './wizard-state'

/**
 * "Translate a file": file → track (video) → preset & context → output → job (preflight, progress).
 * `onSwitchMode`: from the first step, switch to the several-files (batch) mode.
 */
export function TranslatePage({
  onSwitchMode
}: { onSwitchMode?: () => void } = {}): React.JSX.Element {
  const { t, i18n } = useTranslation()
  const [state, dispatch] = usePageReducer('translate.single', wizardReducer, initialWizard)
  const [jobId, setJobId] = usePageState<string | null>('translate.single.job', null)
  const [confirmExisting, setConfirmExisting] = usePageState('translate.single.confirm', false)
  const settings = useSettings().data
  const presets = usePresets().data ?? []
  const create = useCreateJob()
  const prefs = useOutputPrefs()
  const onTrack = useCallback(
    (ids: number[]) => ids[0] !== undefined && dispatch({ type: 'track', trackId: ids[0] }),
    [dispatch]
  )
  const onSubtitle = useCallback(
    (subtitle: ExternalSubtitleValue | null) => dispatch({ type: 'subtitle', subtitle }),
    [dispatch]
  )
  const onMedia = useCallback(
    (patch: Partial<MediaChoice>) => dispatch({ type: 'media', patch }),
    [dispatch]
  )

  const defaults = {
    model: settings?.['translation.defaultModel'] ?? null,
    targetLang: settings?.['translation.defaultTargetLang'] ?? null
  }
  const model = state.model ?? defaults.model
  const targetLang = state.targetLang ?? defaults.targetLang

  // Cached queries: the same keys the steps already use.
  const video = state.file?.kind === 'video' ? state.file.path : null
  const probe = useProbe(video).data
  const track = probe?.subtitles.find((tr) => tr.id === state.trackId)
  const noUsableTrack = !!probe && !probe.subtitles.some((tr) => tr.format !== null)
  // OpenSubtitles pick: possible desync, flagged before the paid translation (never blocks).
  const sync = syncCheck(state.subtitle?.match, probe?.videoFps)
  const syncMessages = useSyncMessages()
  const lookup = useMediaLookup(state.media.ref, null).data
  const work = lookup && state.media.ref && lookup.details.id === state.media.ref.id ? lookup : null
  // The details only reach the AI (and the backdrop) when included; the preset suggestion
  // (like the batch's) follows the work either way.
  const chosen = state.media.include ? (work ?? undefined) : undefined
  const suggestion = work ? suggestPreset(work.details, presets) : null
  // `presetId: null` = follow the TMDB suggestion, else the default from Settings.
  const presetId =
    state.presetId ?? suggestion?.presetId ?? settings?.['translation.defaultPresetId'] ?? null
  const preset = presets.find((p) => p.id === presetId)
  const presetName = (id: string): string => {
    const p = presets.find((x) => x.id === id)
    return p ? presetLabel(p, t) : id
  }
  const created = useJob(jobId).data
  // Artwork: the chosen title while editing, the job's frozen media afterwards.
  const backdrop = jobId ? created?.media?.backdropUrl : chosen?.details.backdropUrl
  const existing = state.file
    ? existingTarget(targetLang, {
        path: state.file.path,
        kind: state.file.kind,
        tracks: probe?.subtitles
      })
    : null
  const targetName = targetLang ? languageName(targetLang, i18n.language) : ''

  const steps = stepsFor(state.file?.kind)
  const index = steps.indexOf(state.step)
  const last = index === steps.length - 1
  const ready = canContinue(state, defaults)
  const reachable = (s: WizardStep): boolean => canGo(state, s, defaults)
  const go = (s: WizardStep): void => dispatch({ type: 'go', step: s, defaults })

  const isVideo = state.file?.kind === 'video'
  const outputMode = state.outputMode ?? prefs.mode
  const extOf = (path: string | undefined): string => path?.split('.').pop()?.toLowerCase() ?? 'srt'
  const subtitleExt = !isVideo
    ? extOf(state.file?.path)
    : state.subtitle
      ? extOf(state.subtitle.path)
      : (track?.format ?? 'srt')
  const outputNames =
    state.file && targetLang
      ? translateOutputNames(
          state.file.path,
          targetLang,
          subtitleExt,
          isVideo ? outputMode : 'subtitle'
        )
      : []
  const preview = state.file
    ? { names: outputNames, folder: outputFolder(state.file.path, state.outputDir) }
    : null

  const createJob = (): void =>
    create.mutate(buildJobRequest(state, suggestion?.presetId, prefs.mode), {
      onSuccess: (job) => {
        if (isVideo) prefs.remember({ mode: outputMode })
        setConfirmExisting(false)
        setJobId(job.id)
      },
      onError: (err) => toast.error(errorMessage(err))
    })
  // Already has a subtitle in the target language: ask first (never silently).
  const start = (): void => (existing ? setConfirmExisting(true) : createJob())

  return (
    <div className="relative min-h-full">
      <Backdrop url={backdrop} />
      <div className="relative mx-auto grid max-w-6xl gap-5 px-6 py-7">
        <h1 className="text-xl font-semibold tracking-tight">{t('pages.translate.title')}</h1>

        {jobId ? (
          <section className="mx-auto w-full max-w-3xl rounded-xl border bg-card/92 p-5">
            <JobDetail
              id={jobId}
              preflightNotice={sync ? <SyncWarning check={sync} /> : undefined}
              showQueueLink
              inspectPath={state.file?.kind === 'subtitle' ? state.file.path : undefined}
              onAnother={() => {
                setJobId(null)
                dispatch({ type: 'reset' })
              }}
            />
          </section>
        ) : (
          <>
            <StepBar
              steps={steps}
              current={state.step}
              label={(s) => t(`translate.steps.${s}`)}
              done={(s) => steps.indexOf(s) <= state.furthest && stepValid(state, s, defaults)}
              reachable={reachable}
              onGo={go}
            />
            {state.notice === 'file-changed' ? (
              <p
                role="status"
                className="mx-auto flex items-center gap-2 rounded-md bg-warning/10 px-3 py-2 text-sm text-warning"
              >
                <Info className="size-4 shrink-0" />
                {t('translate.fileChanged')}
              </p>
            ) : null}

            <div className="grid items-start gap-5 lg:grid-cols-[minmax(0,1fr)_17rem]">
              <div className="grid min-w-0 gap-5">
                {state.step === 'preset' && state.file ? (
                  <MediaPanel fileName={state.file.name} value={state.media} onChange={onMedia} />
                ) : null}

                <section
                  className="grid gap-4 rounded-xl border bg-card/92 p-5"
                  aria-labelledby="translate-step"
                >
                  <h2 id="translate-step" className="text-base font-semibold">
                    {t(`translate.${state.step}.heading`)}
                  </h2>
                  {state.step === 'file' ? (
                    <>
                      <p className="-mt-2 text-sm text-muted-foreground">
                        {t('translate.file.body')}
                      </p>
                      <FileDrop
                        path={state.file?.path ?? null}
                        onFile={(path) => dispatch({ type: 'file', path })}
                      />
                    </>
                  ) : null}
                  {state.step === 'track' && state.file ? (
                    <>
                      {noUsableTrack ? null : (
                        <TrackPicker
                          path={state.file.path}
                          mode="single"
                          value={state.trackId === null ? [] : [state.trackId]}
                          onChange={onTrack}
                          autoSelect={state.subtitle === null}
                        />
                      )}
                      {probe ? (
                        <ExternalSubtitle
                          videoPath={state.file.path}
                          videoFps={probe.videoFps}
                          value={state.subtitle}
                          onChange={onSubtitle}
                          prominent={noUsableTrack}
                        />
                      ) : null}
                    </>
                  ) : null}
                  {state.step === 'preset' ? (
                    <PresetStep
                      value={{ ...state, presetId }}
                      onChange={(patch) =>
                        dispatch({
                          type: 'set',
                          // Picking the suggestion again goes back to following it.
                          patch:
                            suggestion && patch.presetId === suggestion.presetId
                              ? { ...patch, presetId: null }
                              : patch
                        })
                      }
                      presetHint={suggestionHint({
                        t,
                        suggestion,
                        presetId,
                        hasWork: !!state.media.ref,
                        presetName,
                        onUseSuggestion: () => dispatch({ type: 'set', patch: { presetId: null } })
                      })}
                    />
                  ) : null}
                  {state.step === 'output' ? (
                    <OutputStep
                      video={isVideo}
                      outputDir={state.outputDir}
                      mode={outputMode}
                      mux={state.mux}
                      overwrite={state.overwrite}
                      includeOriginal={state.subtitle ? state.includeOriginal : undefined}
                      onIncludeOriginal={(includeOriginal) =>
                        dispatch({ type: 'set', patch: { includeOriginal } })
                      }
                      preview={preview}
                      onFolder={(outputDir) => dispatch({ type: 'set', patch: { outputDir } })}
                      onMode={(m) => dispatch({ type: 'set', patch: { outputMode: m } })}
                      onMux={(patch) => dispatch({ type: 'mux', patch })}
                      onOverwrite={(overwrite) => dispatch({ type: 'set', patch: { overwrite } })}
                    />
                  ) : null}
                  {state.step === 'output' ? (
                    <MemoryOption
                      value={state.memory}
                      onChange={(memory) => dispatch({ type: 'set', patch: { memory } })}
                    />
                  ) : null}

                  <div className="flex items-center justify-end gap-2 border-t pt-4">
                    {index > 0 || onSwitchMode ? (
                      <Button
                        variant="ghost"
                        className="mr-auto"
                        onClick={() => (index > 0 ? dispatch({ type: 'back' }) : onSwitchMode?.())}
                      >
                        {index > 0 ? t('translate.back') : t('choose.toMany')}
                      </Button>
                    ) : null}
                    {last ? (
                      <Button variant="default" disabled={create.isPending} onClick={start}>
                        {create.isPending ? t('translate.starting') : t('translate.start')}
                      </Button>
                    ) : (
                      <Button
                        variant="default"
                        disabled={!ready}
                        onClick={() => dispatch({ type: 'next', defaults })}
                      >
                        {t('translate.next')}
                      </Button>
                    )}
                  </div>
                </section>
              </div>

              <WizardSummary
                state={state}
                preset={preset}
                model={model}
                targetLang={targetLang}
                track={track}
                syncWarn={sync ? syncMessages(sync, 'short').join(' · ') : undefined}
                media={chosen}
                existing={existing}
                outputNames={outputNames}
                canGo={reachable}
                onGo={go}
              />
            </div>
          </>
        )}
      </div>

      <Dialog open={confirmExisting} onOpenChange={setConfirmExisting}>
        <DialogContent
          title={t('translate.existing.title', { language: targetName })}
          description={
            existing?.where === 'track'
              ? t('translate.existing.track', {
                  language: targetName,
                  name: existing.track.name ?? `#${existing.track.id}`
                })
              : t('translate.existing.file', { language: targetName })
          }
        >
          <DialogFooter>
            <DialogClose asChild>
              <Button>{t('actions.cancel')}</Button>
            </DialogClose>
            <Button variant="default" disabled={create.isPending} onClick={createJob}>
              {t('translate.existing.go')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
