# Jobs

`src/main/jobs/` runs every user action on files as a persisted **job**: extract tracks, translate a subtitle, mux subtitles into an MKV, or the whole pipeline. Pure Node apart from the IPC handlers; composed in `create-job-module.ts`.

```
jobs:create ──► resolveJobSpec (Zod + path guard + resolved preset snapshot) ──► jobs row (queued)
JobManager (FIFO, concurrency = settings 'jobs.concurrency', default 2)
  └─ JobRunner.run(spec, artifacts, workDir)       userData/cache/jobs/<id>/
        resume at the first step whose artifact does not verify
        probe → extract → preflight → translate → write → mux   (subset per kind)
  └─ status, step, progress, artifacts, result, error persisted at every transition
jobs:progress (throttled 250 ms; state transitions immediately)
```

During the translate step `progress.translation` carries `done/total/cached`, `chunk { completed, count, inFlight }` and the `timeline` map of the episode (cached, translated, failed and in-flight ranges, plus the lane cuts; see [`llm.md`](./llm.md#lanes-parallel-translation)). While every lane waits on OpenRouter (rate limit, lane cooldown after a 429, backoff), `progress.waiting { reason, untilMs }` is set; it is cleared when the request goes out and never kept in a terminal state. See [`rate-limits.md`](./rate-limits.md).

## Kinds

| `kind`      | Input                                           | Steps                                                     | Output                                                                                                                                                                                                            |
| ----------- | ----------------------------------------------- | --------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `extract`   | MKV + track ids                                 | probe → extract → write                                   | `<base>.<lang\|und>.<ext>` (`naming: 'player'`, default with one track; same language twice → `<base>.<lang>.<n>.<ext>`) or `<base>.<trackId>.<lang\|und>.<ext>` (`naming: 'track'`, default with several tracks) |
| `translate` | `.srt` / `.ass` / `.ssa` / `.vtt`               | translate → write                                         | `<base>.<targetLang>.<ext>`                                                                                                                                                                                       |
| `mux`       | MKV + 1..32 subtitles (language, name, default) | probe → validate → mux                                    | `<base>.<lang>.mkv`                                                                                                                                                                                               |
| `pipeline`  | MKV + track id, or MKV + external subtitle      | probe → extract → preflight → translate → [write] → [mux] | by `outputMode` (see below): `<base>.<targetLang>.mkv`, `<base>.<targetLang>.<ext>`, or both                                                                                                                      |

A `pipeline` whose source is an **external subtitle** (a file from disk or an OpenSubtitles download, for a video with no usable track) has no `extract` step and can also embed that subtitle (`options.includeSourceSubtitle`). See [`subtitle-sources.md`](./subtitle-sources.md).

### Pipeline output mode (`options.outputMode`)

| `outputMode`    | Steps after `translate` | Result                                                                                                                    |
| --------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `mux` (default) | `mux`                   | `<base>.<targetLang>.mkv`                                                                                                 |
| `subtitle`      | `write`                 | `<base>.<targetLang>.<ext>` next to the video (or in `outputDir`); the MKV is never read again, no remux free-space check |
| `both`          | `write` → `mux`         | both files                                                                                                                |

- The subtitle name is `<video base>.<BCP-47 language>.<ext>` (`Episode.mkv` → `Episode.pt-BR.ass`), the convention players use to auto-load external subtitles. `JobResult.subtitle` carries its path (also the first entry of `outputs`).
- **Compatibility**: `options.alsoWriteSubtitle` (old requests) is deprecated: when `outputMode` is absent, `true` → `both`, `false`/absent → `mux`; an explicit `outputMode` wins. New specs store `outputMode`; specs persisted before it (only `alsoWriteSubtitle`) keep working: `planFor` reads them through `pipelineOutputMode`, so queued/failed jobs resume with the same steps. The renderer's default stays `mux` (the last choice is remembered in Settings by the UI).
- `jobs:retry`/cancel need no special case: `write` is never reused from a previous attempt, and `subtitle` mode has no step that touches the MKV. Cancel removes the `.partial` and leaves the source and any finished output alone.
- **Name collisions** (`extract`, `translate`, `pipeline` subtitle): ` (2)`, ` (3)`… by default; `options.overwrite: true` replaces the existing file (atomic rename), except when the target is the job's own source file. The remuxed `.mkv` always gets a unique name.
- `extract` naming: `options.naming` (`player` | `track`); default `player` with a single track, `track` with several (batch: `options.extractNaming`, applied per item). Specs persisted without `naming` behave as `track`.

Steps (`steps/*.step.ts`) are shared by every kind (`plan.ts`):

- **probe**: `mkvmerge -J` → subtitle tracks (`media/mkv-probe.ts`). Supported: `S_TEXT/UTF8`, `S_TEXT/ASCII`, `S_TEXT/ASS`, `S_TEXT/SSA`, `S_TEXT/WEBVTT`. PGS/VobSub/DVB (`image-based`, would need OCR, out of scope) and other text codecs without a parser (`no-parser`, e.g. USF) are listed but rejected. Track language: `language_ietf`, else ISO 639-2 → BCP-47 table (`media/languages.ts`); `und` = unknown.
- **extract**: one `mkvextract` pass for all requested tracks into the work dir.
- **preflight**: token/cost estimate and coverage; may stop the job in `awaiting-confirmation` (see [`preflight.md`](./preflight.md)).
- **translate**: decode → parse → `translateDocument` (see [`llm.md`](./llm.md)) → `serializeRepaired` → UTF-8 (BOM if the source had one). `sourceLang`: job option, else filename suffix (`ep01.en.srt`), else the track language.
- **validate** (mux): each subtitle must parse; nothing is repaired. Failure → `INVALID_SUBTITLE` before touching the MKV.
- **write**: publishes work-dir files into the output dir (unique name, atomic copy).
- **mux**: `mkvmerge` keeps every original track and attachment (fonts); the new track gets `--language`, `--track-name` (`<Language> (BakaSub)` in the pipeline) and the default flag; with `setDefault` the original subtitle tracks lose theirs; `keepOriginalSubtitles: false` → `--no-subtitles`.

## Outputs never corrupt anything

- The source is never overwritten. Name collisions get ` (2)`, ` (3)`… (unless `overwrite` is requested for a subtitle output).
- Every final file is written as `.<name>.bakasub-partial.<ext>` in the destination directory and renamed only on success (atomic on the same filesystem); failures and cancellation delete it.
- Free space is checked before mux (`source + subtitles × 1.05`, `DISK_FULL`) and before extraction (100 MB in the work dir).

## mkvtoolnix (`media/mkvtoolnix*.ts`)

- Binaries come from the verified sidecars (`SidecarManager.resolve`, `SIDECAR_MISSING` if not installed) and run through `ProcessRunner` (abort → SIGTERM → SIGKILL after 3 s; **no timeout**, big files take long).
- `--gui-mode`: `#GUI#progress NN%` → step progress, `#GUI#warning` → job warnings, `#GUI#error` → error message. `ProcessRunner` gained `onStdoutLine` for this.
- Exit codes: `0` ok, **`1` = finished with warnings (success → `completed-with-warnings`)**, `2` error (`MEDIA_TOOL_FAILED`).
- Not yet validated against the real binaries (no mkvtoolnix on the dev machine at the time): the `mkvmerge -J` fields read (`codec_id`, `language`, `language_ietf`, `track_name`, `default_track`, `forced_track`, `num_index_entries`) and the `--gui-mode` position for `mkvextract` must be checked once the sidecars are installed.

## States

```
queued → running(step) → [awaiting-confirmation] → completed | completed-with-warnings
                        ↘ failed | cancelled | interrupted
```

- `awaiting-confirmation`: the preflight is waiting for the user ([`preflight.md`](./preflight.md)); not running, the slot is free, survives restarts; `jobs:confirm` re-queues it, `jobs:cancel` discards it.
- `completed-with-warnings`: translation left lines untranslated (`result.failedLines`: `cueId`, `reason`, and the source `excerpt` + `location` (line, start/end ms) built with `analyzeCoverage`, so the UI lists which lines stayed in the original language) or mkvtoolnix reported warnings (`result.warnings`). The term consistency report (`result.consistency`, see [`llm.md`](./llm.md#term-consistency-pass-qualityterm-consistencyts)) is informational and never changes the status; `jobs:dismiss-consistency` only sets `result.consistencyDismissed` on a terminal job that has one.
- `failed`: `error = { code, message }` (`AppError` code, or `INTERNAL` for bugs) and `step` = where it stopped. The work dir is kept.
- `cancelled` (user) deletes the work dir; `interrupted` (shutdown or crash) keeps it. On boot, `running` rows become `interrupted`.
- `result.usage` keeps the LLM cost even for cancelled/failed jobs.

## What the queue sees of a spec

The renderer never gets the spec; `toDto` derives small summaries from it:

- `Job.media` (`media-summary.ts`): title, artwork and episode of the frozen TMDB context.
- `Job.translation` (`translation-summary.ts`, `translate`/`pipeline` only): `model` (OpenRouter id), `preset` (`{ id, name }` as named at creation), `targetLang`, `sourceLang`, read from the `PresetSnapshot` and languages the spec froze in `resolveJobSpec`. Every field is optional and type-checked: a spec persisted without one (or with a malformed one) just omits it and the UI says "not recorded". No migration is needed.
- `Job.memory` (`jobMemoryMode`).

## Retry, cancel, remove

- `jobs:retry` (failed / cancelled / interrupted): back to `queued`, `attempt + 1`. The runner resumes at the **first step whose artifact does not verify** (`artifacts.ts`: probe info, extracted and translated files with size + sha256); everything from there on is redone. Outputs (`write`, `mux`) are always redone. Translation is cheap anyway thanks to the Job Memory. The retry keeps the job's `memory` choice; `jobs:retry { id, memory }` can switch it, which also drops the `preflight` and `translate` artifacts (see [`job-memory.md`](./job-memory.md)). Without a switch, an unfinished translation also re-runs the preflight, without asking again if the old estimate had been accepted, so the estimated vs actual comparison covers one attempt (see [`preflight.md`](./preflight.md#decision)).
- `jobs:cancel`: queued → `cancelled` directly; running → abort (LLM request and backoff, child processes); a translation already paid for is still written to the cache.
- `jobs:remove`: terminal jobs only; deletes the row and the work dir.
- `jobs:clear-completed`: the queue's "Clear finished". Deletes every `completed`, `completed-with-warnings` and `cancelled` job in one `DELETE … WHERE id IN (…)`, then their work dirs (a dir that fails to go is collected later as an orphan). Failed and interrupted jobs stay, since they still need a decision.
- Presets: the job stores a **resolved snapshot** of the preset (`spec.preset`, model and target language already decided, see [`presets.md`](./presets.md)), so editing the preset or updating the built-ins does not change a queued job or its retry; deleting a preset used by a non-terminal job → `PRESET_IN_USE`.

## Shutdown and housekeeping

- The `jobs` disposer is registered last (runs first): abort every job, wait up to 5 s, mark leftovers `interrupted`; then sidecars, Job Memory and the database close.
- Boot: work dirs without a job row, and dirs of terminal jobs finished more than 7 days ago, are deleted.
- `isBusy()` (any job running) blocks the Job Memory prune.

## Paths and security

Paths only come from the native dialogs (`dialog:open-files`, `dialog:choose-output-dir`) or drag & drop (`window.api.files.pathFor(file)` = `webUtils.getPathForFile`). Main revalidates every path in `jobs:create` / `media:probe` (`media/path-guard.ts`): absolute, resolved symlinks, regular file, extension allowlist (`.mkv` / `.srt .ass .ssa`) checked on both the given and the real path, readable, subtitles ≤ 50 MB; output dirs must be writable directories. Errors: `INVALID_PATH` with a stable reason (`relative`, `not-found`, `not-a-file`, `extension`, `too-large`, `not-writable`…).

## Batch mode (`src/main/batch/`)

A batch groups independent jobs created together; it executes nothing itself. Kinds: `extract`, `translate`, `pipeline`, `mux` (one `mux` job per MKV, with its paired subtitles). Every kind has a screen: Translate (`translate` / `pipeline`, "Several files or a folder"), Extract and Mux ("Select several files"); see [`renderer.md`](./renderer.md#translate-several-files-batch).

```
dialog:open-folder | dialog:open-files (multiple) | drag & drop
  └─ batches:plan   scan (folder or files) → probe MKVs (4 in parallel) → automatic selection
  └─ UI shows the plan; the user adjusts tracks / drops files
  └─ batches:create resolve + validate EVERY item, then one transaction: batch + all jobs (queued)
```

- **Scan** (`scan-files.ts`): folder from the native dialog, optional recursion, natural order (`ep2` before `ep10`), root files before subfolders, hidden entries skipped, symlinked directories not followed, at most 2,000 files and 10 levels (`truncated` in the plan). Every file goes through the same path guard as single jobs; loose files the guard rejects come back in `plan.rejected` with the reason.
- **Automatic track selection** (`track-selection.ts`):
  - `pipeline`: one track, by preferred source language (setting `tracks.preferredSourceLangs`, default `['en']`, list order, `en-US` matches `en`) → ASS/SSA before SRT → not forced → name not matching `sign|song|karaoke|lyric` → most events → lowest id.
  - `extract`: every supported text track.
  - Files without a usable track stay in the plan with `skipped` (`no-subtitle-track`, `no-supported-track`, `probe-failed` + `error`).
- **Create** (`batch-service.ts`): all items are resolved (`resolveJobSpec`, 4 in parallel) before anything is written; one bad item rejects the whole batch (`INVALID_INPUT`, message `item N: …`, always the first failing item in plan order). `translate`/`pipeline` take batch defaults `presetId`, `model`, `targetLang`, `notes` (each falls back to Settings: `translation.defaultPresetId` / `defaultModel` / `defaultTargetLang`; missing → `MODEL_NOT_SET` / `TARGET_LANG_NOT_SET`).
- **Works** (`groups`): the unit of context of a translate batch. `groups[]` (`BatchGroup`: `name?`, `media?` TMDB ref, `presetId?`, `model?`, `targetLang?`, `notes?`) and `items[].group` (index), with `items[].season`/`episode` for series. A group's settings replace the batch defaults for its files (notes too: an empty string means no notes); `media` comes only from the group, so files without a group, or in a group without `media`, carry no media context. The TMDB lookup goes through the same `MediaLookupSource` as single jobs (a lookup failure is logged and the job runs without it). Unknown group index or `groups` on an `extract` batch → `INVALID_INPUT`. The groups are kept in the batch's `options` JSON.
- **Mux batches** (`kind: 'mux'`, `mux-pairing.ts`): `batches:plan` scans MKVs **and** subtitles (`.srt .ass .ssa .vtt`) from a folder or loose files, without probing, and returns `BatchMuxPlan`: `items[]` = `{ video, subtitles[] }` (1..32, natural order, each with `language` and `matchedBy`), `unmatchedVideos`, `unmatchedSubtitles`, `rejected`, `truncated`. Pairing is deterministic (no edit-distance guessing):
  1. Same folder and same base name, case-insensitive, ignoring language/flag suffixes on the subtitle (`X.mkv` ↔ `X.pt-BR.srt`, `X.en.forced.ass`, `X.srt`). A `Subs/` or `Subtitles/` folder next to the video counts as the video's folder.
  2. Otherwise, same folder and same season/episode taken from the file names with the release-name parser of the media context (`S01E05`, `1x05`, or absolute anime numbering `Show - 28`).
  3. More than one candidate video (two files with the same base name or the same episode) is **ambiguous**: never resolved here, the subtitle goes to `unmatchedSubtitles` for the UI to place. A subtitle beyond the 32nd of a video also lands there.
  - Language comes from the filename suffix (`.pt-BR`, `.en`), `und` without one; the user edits it before `create`. `und` is a valid language for the mux job.
  - `batches:create` takes `items: { video, subtitles: { path, language, name?, default? }[] }[]` plus the batch-level `options` (`setDefault`, `keepOriginalSubtitles`; `alsoWriteSubtitle`, `outputMode` and `extractNaming` are ignored). Every item is resolved before the transaction (same `resolveJobSpec` as the single mux job); a bad item, a repeated video, an item of the wrong shape for the kind, or `groups` → `INVALID_INPUT` (`item N: …`). Subtitle parsing is validated in the job (`validate` step, `INVALID_SUBTITLE`), as in the single mux. Preset/model/language/notes/media settings do not apply.
  - Output `<base>.<lang>.mkv` (first subtitle's language) next to the source, mirrored under `outputDir` for a folder batch; collisions get ` (2)`. FIFO, cancel and remove behave as in the other kinds.
- **Output options** (`batches:create` `options`): `outputMode` (`pipeline`, default `mux`; `alsoWriteSubtitle` is the deprecated alias), `extractNaming` (`extract`), `overwrite` (`extract`/`translate`/`pipeline` subtitles); all forwarded to every item's job.
- **Output**: without `outputDir`, next to each source. With `outputDir` and a folder batch, the relative structure is **mirrored** (subfolders created on demand); items must be inside the batch folder. Loose-file batches with `outputDir` are flat; name collisions get ` (2)`.
- **Order**: jobs keep the plan order in the FIFO queue (staggered `created_at`).
- **State**: `Batch.counts` per job status, `total`, `state: active | finished`.
- **Cancel**: queued jobs first (so the queue cannot pick the next item of the same batch), then running ones.
- **Remove**: only when no job is active; removes the jobs (and their work dirs) and the batch.

## Pending

- Image-based subtitles (PGS/VobSub/DVB) are out of scope (OCR dropped).
