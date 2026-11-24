# Architecture overview

Entry point for developers. Each area links to the document that holds the details.

## Stack

| Layer             | Technology                                                                                             |
| ----------------- | ------------------------------------------------------------------------------------------------------ |
| Wrapper / backend | Electron (Main Process as a Node.js + TypeScript backend)                                              |
| Frontend          | React 19 (Vite) + Tailwind CSS v4 + TanStack Query + TanStack Router                                   |
| Communication     | Strict native IPC: `contextIsolation: true`, `nodeIntegration: false`, `sandbox: true`, preload bridge |
| Database          | SQLite (WAL) via Drizzle ORM + LibSQL driver                                                           |
| Validation        | Zod (IPC payloads and remote manifests)                                                                |
| LLM               | OpenRouter over SSE (native `fetch`)                                                                   |
| Sidecars          | mkvtoolnix (`mkvextract`/`mkvmerge`) and ffmpeg via `child_process`                                    |
| Logging           | `electron-log`, rotating file in `userData/logs`                                                       |

Tooling runs on Node **26** (`.nvmrc`). The Main Process runs on Electron 44's embedded Node (24.x), so `@types/node` stays at `^24`.

Supported platforms: **Linux x64** and **Windows x64**. There is no macOS build.

## Layout

```
src/
  main/
    config/        paths (userData), remote URLs
    database/      client (WAL, migrate on boot), schema/, repositories/
    ipc/           create-handler (sender check + Zod), register, per-domain handlers/
    logging/       electron-log (scopes: app, db, ipc, sidecar, llm, jobs, renderer)
    services/      shutdown-coordinator, app-info, app-update, native consent dialog
    sidecars/      sidecar module (no `electron` dependency, testable in plain Node)
    parsers/       subtitle parsing: decoder, common AST, SRT, ASS/SSA, WebVTT, segmentation, coverage (pure, testable)
    llm/           OpenRouter client (SSE), model catalog, chunking, prompts, retries, translateDocument
    memory/        Job Memory: context hash and cache service
    media-context/ TMDB media context: release-name parser, client, cache, artwork protocol
    subtitle-sources/ OpenSubtitles: client, session, movie hash, search and cached download
    media/         mkvtoolnix facade (probe/extract/mux), languages, output naming, path guard
    jobs/          job queue, runner, steps (probe/extract/translate/validate/write/mux), job specs
    batch/         batch planning (scan, probe, track selection) and batch service
    preflight/     token counting (js-tiktoken) and cost estimate
    presets/       built-in (immutable) and custom translation presets
    settings/      typed settings registry and the encrypted secret vault
    window/        BrowserWindow factory
  preload/         the only surface exposed to the renderer (window.api), allowlisted
  renderer/        React + Tailwind + TanStack Query + TanStack Router + i18next
    src/guides/    in-app user guides, one folder per language (en, pt-BR)
  shared/          typed IPC contract and pure utilities (no zod/electron)
drizzle/           versioned migrations
scripts/           cross-build (Docker/wine), icons and native binaries
docs/              detailed system flows
```

## Main flows

- **IPC**: `src/shared/ipc-contract.ts` is the single source of truth for channels, payloads and `window.api`. The preload only exposes channels listed in `INVOKE_CHANNELS`/`EVENT_CHANNELS` and never the raw `ipcRenderer`. Every handler goes through `createHandler`, which rejects untrusted `senderFrame`s, validates input with Zod and returns a `{ ok, data } | { ok: false, error }` envelope. See [`ipc.md`](./ipc.md).
- **Graceful shutdown**: `ShutdownCoordinator` runs disposers in LIFO order with a global timeout on `before-quit` and on `SIGINT`/`SIGTERM`. It aborts running jobs first (marking them `interrupted`), then releases sidecar child processes (SIGTERM → SIGKILL), the Job Memory writes and the database (`wal_checkpoint(TRUNCATE)`).
- **Subtitle parsing**: lossless AST per format; only text nodes go to the LLM. See [`parsing.md`](./parsing.md).
- **Translation**: chunking, structured JSON, validation, retries, termbase. See [`llm.md`](./llm.md), [`rate-limits.md`](./rate-limits.md) and [`presets.md`](./presets.md).
- **Cost and cache**: [`preflight.md`](./preflight.md) and [`job-memory.md`](./job-memory.md).
- **Jobs and batches**: queue, runner, steps, cancellation. See [`jobs.md`](./jobs.md).
- **Media context (TMDB)**: [`media-context.md`](./media-context.md).
- **External subtitles** (file from disk or OpenSubtitles, for videos with no usable track): [`subtitle-sources.md`](./subtitle-sources.md).
- **Sidecars**: remote manifest, verified download, pruning. See [`sidecars.md`](./sidecars.md).
- **Database**: SQLite in WAL mode (`synchronous=NORMAL`, `foreign_keys=ON`, `busy_timeout=5000`), Drizzle migrations applied on boot (`drizzle/` ships in `resources/`). See [`database.md`](./database.md).
- **Native notifications**: emitted by the Main, only in the background, per-category settings. See [`notifications.md`](./notifications.md).
- **Renderer**: [`renderer.md`](./renderer.md). **Brand tokens** (shared with the site): [`brand-tokens.md`](./brand-tokens.md). **Website** (`site/`): [`site.md`](./site.md). **Logging**: [`logging.md`](./logging.md). **Build and release**: [`release.md`](./release.md).

## User data

Everything lives in `app.getPath('userData')`, never in the install directory:

```
bakasub.db            SQLite (WAL)
logs/main.log         JSON lines, rotates at 5 MB (5 files kept)
cache/manifests/      last valid remote manifest (offline fallback)
cache/artwork/        TMDB artwork (LRU, 200 MB)
cache/opensubtitles/  subtitles downloaded from OpenSubtitles (<fileId>.srt)
cache/jobs/<id>/      per-job work directory
bin/                  sidecars (+ installed.json)
```

The folder depends on how the app runs. `src/main/bootstrap/user-data.ts` is the **first import** of `src/main/index.ts` and calls `app.setPath('userData', …)` before the logger, the single-instance lock, the database or any `app.getPath('userData')`; the rule itself is the pure `resolveUserDataDir` (`src/main/config/user-data-dir.ts`, unit-tested):

1. `BAKASUB_USER_DATA` set (non-empty): that directory, packaged or not (relative paths resolve from the working directory).
2. Unpackaged (`!app.isPackaged`, i.e. `electron-vite dev` or `electron .`): `<appData>/bakasub-dev`.
3. Installed: Electron's default, `<appData>/bakasub`.

|                      | Installed app                                  | Development (`npm run dev`) |
| -------------------- | ---------------------------------------------- | --------------------------- |
| userData (Linux)     | `~/.config/bakasub`                            | `~/.config/bakasub-dev`     |
| userData (Windows)   | `%APPDATA%\bakasub`                            | `%APPDATA%\bakasub-dev`     |
| Database             | `<userData>/bakasub.db` (+ `-wal`, `-shm`)     | same, in the dev folder     |
| Logs                 | `<userData>/logs/main.log`                     | same, in the dev folder     |
| Sidecars             | `<userData>/bin/` (`installed.json`)           | same, in the dev folder     |
| Cache                | `<userData>/cache/` (manifests, artwork, jobs) | same, in the dev folder     |
| Single-instance lock | keyed on `bakasub`                             | keyed on `bakasub-dev`      |

So a dev build and the installed app can run at the same time without sharing the lock, the WAL database, logs, sidecars or cache, and a new migration applied in dev never touches the installed database. Settings → About shows the real `userDataPath`.

## Sidecars in short

ffmpeg and mkvtoolnix are **not bundled**. The app reads [`manifests/sidecar.json`](../manifests/sidecar.json) from this repository (via `raw.githubusercontent.com`): URL, SHA-256, archive type and binary path per `tool × platform × arch`.

- Never downloads automatically: `sidecars:install` requires a native dialog in the Main Process (default button: Cancel), which a compromised renderer cannot approve.
- Streaming download with on-the-fly SHA-256, HTTP Range resume and exponential-backoff retry; a hash mismatch discards the file.
- Every binary is executed (`--version`) before being recorded in `bin/installed.json`.
- All external CLI execution goes through `ProcessRunner` (timeout, SIGTERM→SIGKILL, stderr captured to the log, tracked for shutdown).

```ts
const { path, env } = await sidecars.resolve('mkvmerge')
// `env` is required: on Linux it injects LD_LIBRARY_PATH for the extracted AppImage libs
```

`mkvextract` writes `.srt` files with a UTF-8 BOM; the lexer strips it before reading the first cue.

## Feature status (1.0.0)

- [x] Strict IPC, logging, database with migrations, graceful shutdown
- [x] Sidecars: manifests, verifiable download, resume, pruning, consent
- [x] App version check and auto-update (GitHub Releases)
- [x] Build and release (Linux + Windows)
- [x] **SRT, ASS/SSA and WebVTT parsers**: lossless AST, tolerant scanner with diagnostics, safe repair mode, XLIFF-style segmentation, coverage report. Image-based subtitles (PGS/VobSub) are out of scope
- [x] **Job Manager**: persisted jobs (`extract`, `translate`, `mux`, `pipeline`), FIFO queue with bounded concurrency, per-job `AbortController`, retry from the failed step, atomic outputs, presets (4 built-ins + custom, JSON export/import)
- [x] **Batch mode** for Translate, Extract and Mux, with MKV ↔ subtitle pairing for mux
- [x] **LLM layer**: `<context>`/`<target>` chunks, parallel lanes, AIMD in-flight cap, structured JSON validated with Zod, per-line repair, termbase, SSE with cancellation and cost accounting, encrypted API key
- [x] **Term consistency pass** (local, no LLM)
- [x] **Job Memory**: validated cache keyed by context hash + language + model + glossary
- [x] **Preflight**: local token estimate × cache hit rate × live OpenRouter pricing, coverage report, confirmation before translating
- [x] **Media context** (TMDB): release-name parsing, conservative auto-match, local artwork cache
- [x] **External subtitles** for videos with no usable track: a file from disk or an OpenSubtitles search (movie hash + title), single-file wizard only
- [x] **Rate limits and retries**: per-lane limiter, 429 cooldown, free-model daily quota
- [x] **Frontend**: onboarding, Translate wizard, queue, Extract, Mux, Presets, Models, Languages, Logs, Settings; English and Brazilian Portuguese
- [x] **Parser diagnostics and encoding confidence** for subtitle files in the Translate wizard preflight

Known gaps: no batch view after leaving the page that created it and no "remove batch" (`batches:list` is unused); no encoding override in `jobs:create`; videos and the queue detail do not show parser diagnostics.
