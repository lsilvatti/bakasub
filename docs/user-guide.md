# User guide

How to use BakaSub day to day. For installation see the [README](../README.md); for problems see [`reporting-issues.md`](./reporting-issues.md).

## Concepts

- **Job**: one unit of work (extract, translate, mux, or the full pipeline for a video). Jobs run in the background and live in the **Queue**.
- **Batch**: many files at once. It creates one job per file.
- **Preset**: the _style_ of a translation (media type, temperature, style instructions, glossary). The model and the target language are chosen per job.
- **Sidecars**: ffmpeg and mkvtoolnix, the external tools that read and write MKV files. They are downloaded on demand.
- **Job Memory**: a local cache of translated lines.

## First run

The onboarding saves each step as you go:

1. **OpenRouter key** (required). Create one at [openrouter.ai/keys](https://openrouter.ai/keys) and add credits, or use free models. The key is checked before it is saved and stored encrypted with your OS keyring when there is one; if there is none, the app warns you that it is stored without OS protection. The key can be replaced but not removed, because the app needs it.
2. **Models** (required). Pick favorites and one default from the OpenRouter catalog. Each model shows its price per 1M tokens and badges: free, reliable JSON, moderated.
3. **Languages** (required). Choose your target languages and a default.
4. **Video tools** (optional, but needed for any MKV). The app asks for explicit consent before downloading ffmpeg and mkvtoolnix.
5. **TMDB key** (optional). Enables movie/series details (synopsis, episode, artwork) to improve translations.

Everything can be changed later in Settings, Models and Languages.

## Translate

**Translate** has two modes: one file, or several files / a folder.

### One file

Drop a file or pick one. A subtitle file (`.srt`, `.ass`, `.ssa`, `.vtt`) is translated directly. A video (`.mkv`) goes through the full pipeline: pick the source track, translate, and write the result.

Steps:

1. **File**, and for a video the **track**. Image-based tracks (PGS, VobSub, DVB) and unsupported formats are listed but disabled, with the reason. A track is preselected from your preferred source languages (Settings → General). If the video has **no subtitle that can be translated**, bring one in right there: choose a subtitle file you already have, or **search OpenSubtitles** (see below). The same two options are available, as "Use another subtitle", when the video does have tracks.
2. **Preset and context**. The preset can be suggested from the movie/series the app identifies from the file name (TMDB). Pick the model and target language (defaults are prefilled) and optionally write **notes** for the AI: synopsis, characters, tone (up to 4,000 characters).
3. **Output**. The folder (default: next to the source) and, for a video, what to write: the MKV with the new subtitle, the subtitle file only (for players that load external subtitles by name), or both. For an MKV you can set the new track as default and keep or drop the original subtitles; when the source was a subtitle from outside the video, you can also add that original subtitle to the MKV. A preview shows the final file names; by default existing files are kept and a ` (2)` suffix is added, unless you choose to replace them.
4. **Review cost** (see below), then confirm.
5. **Progress and result**: live progress, cost so far, and at the end the output files, final cost, warnings and any untranslated lines.

Your wizard state survives switching pages. Closing the app does not lose jobs: they are saved, and ones that were running are marked **interrupted** so you can retry them.

#### Searching OpenSubtitles

The search needs your own free OpenSubtitles API key (Settings → Keys; the in-app guide shows how to get one, and the search dialog asks for it the first time). It starts on its own from the file: a fingerprint of the video plus the title, season and episode read from the file name. You can change the title, numbers and languages and search again.

- **Made for this file** marks subtitles made for this exact video file: their timing matches. Other results match by title and are marked **Not verified for this file**: they may be out of sync; BakaSub does not adjust timing.
- When a result was timed for another frame rate than the video (for example 25 fps against 23.976 fps), it is flagged in orange: its lines will drift. The warning is repeated under the chosen subtitle, in the summary and in the review before translating. It never stops you from going on.
- **Machine translated** results are a poor source to translate from.
- Downloads are limited per day by OpenSubtitles (searching is not). Signing in to your OpenSubtitles account in Settings → Keys raises the limit. A subtitle already downloaded is kept and costs nothing to pick again.

Several files or a folder (below) do not search OpenSubtitles: videos without a usable subtitle are skipped there.

### Several files or a folder

Pick a folder (optionally including subfolders) or several files. The app scans them, picks tracks automatically from your preferred languages and groups files into **works** (series or movie), each with its own TMDB context, preset and settings. You review the plan, exclude files, change tracks, then create the batch. Creating a batch is all-or-nothing: if one item is invalid you are pointed to the file to fix or drop. Folder structure is mirrored in the output folder.

## Cost review and confirmation

Before anything is billed, each translation stops at **awaiting confirmation** and shows:

- lines in total, from cache, and to translate; number of requests and tokens;
- an estimated **cost range** (for a model with dynamic pricing only tokens are shown), the model and the date of the price catalog; a note when the model uses a different tokenizer (the estimate is approximate);
- a **coverage report**: what will _not_ be translated and why (lines with no letters such as `♪♪`, ASS drawings, karaoke lines, unreadable blocks, text that looks like an internal marker). Untranslated lines are kept as they are, never silently dropped;
- for subtitle files, the detected **encoding** and its confidence, and syntax problems found in the file (and which were fixed automatically on export).

You can set a threshold in Settings → General to continue without asking when the maximum estimated cost is below a value. Jobs awaiting confirmation survive restarts.

**Cancelling** stops the job immediately, but some providers keep generating (and billing) a request already in flight, so the final cost can include interrupted generations.

## Extract

Pick a video and one or more subtitle tracks, then an output folder. Files are named `<video>.<lang>.<ext>` so players find them, or with the track id if you prefer. Several tracks of the same language get a numeric suffix. A batch mode handles many videos at once.

## Mux

Pick a video and one or more subtitle files. For each subtitle set the language (guessed from the file name suffix, e.g. `ep01.pt-BR.ass`; `und` is allowed), the track name, and which one is default (at most one). Options: set the new subtitle as default, and keep the original subtitles. The batch mode pairs MKVs with subtitle files by base name or by season/episode and lets you move, drop or reassign pairs.

## Queue

Every job appears in the Queue with tabs **All / Active / Needs you / Done**. "Needs you" means awaiting confirmation, failed or interrupted. Select a job to see its progress timeline, cost, outputs, warnings and untranslated lines. Available actions depend on the state: confirm, cancel, retry (resumes from the step that failed), remove. Errors come with an action when possible (for example, no credits opens the OpenRouter credits page, missing tools offers to install them). The status bar shows running and queued counts, the session cost and a countdown when the app waits on rate limits.

## Presets and glossary

Four built-in presets (**Movie, Series, Anime, Documentary**) are read-only; use **Duplicate** to customize one. A custom preset has:

- a name and media type, a temperature (0 to 2) and style instructions. Write the style instructions in English whatever the target language;
- a **glossary** (termbase): source term, optional target and optional rule, injected as strict rules in the prompt. Useful to keep honorifics (`-san`, `-sensei`), names and jargon. After a translation, the result lists recurring terms that were rendered in different ways, so you can save your chosen rendering to the preset;
- export and import as JSON (with rename or skip when a name is taken).

A preset in use by a queued or running job cannot be deleted.

## Settings

- **General**: app language (system, English, Português do Brasil), default preset, jobs at the same time, auto-confirm threshold, preferred source languages, detailed logs, system notifications. Notifications appear only while BakaSub is in the background (never while its window is in front): a file finished, failed or waits for cost confirmation; a batch finished or waits for confirmation (one notification per batch); a new version is available. Each category can be turned off, or all of them at once. Clicking a notification brings the window back on that job, batch or the About panel.
- **Keys**: OpenRouter, TMDB and OpenSubtitles (API key, plus an optional account sign-in for a bigger daily download limit).
- **Translation cache**: size on disk, number of lines, last used date. You can prune entries not used for a number of days or clear everything. This is disabled while a job is running.
- **Video tools**: install or re-check ffmpeg and mkvtoolnix.
- **Advanced**: parallelism and rate-limit tuning, with the recommended range shown; values outside it warn but do not block. "Restore defaults" resets only these.
- **About**: version, platform, data folder, update panel (check, download, restart; every step is an explicit click).
  When a new version is found at startup, a dialog shows the version, date, size and release notes; "See details and update" opens this panel, "Later" hides it until the next start.

## The translation cache

Translated lines are saved locally, so reprocessing a file (or another file with the same lines) is nearly free. A cached line is reused only if the target language, model and glossary are the same. Short lines (three words or fewer) also take the previous line into account, so a "Yes." in two different scenes is not mixed up. Notes and TMDB context are not part of the key.

## Privacy: what leaves your computer

- **OpenRouter** receives the subtitle text to translate (with neighboring lines as context), your notes, the preset instructions and glossary, and TMDB synopsis if you chose to include it. Timestamps and formatting tags are not sent. The provider's own policy applies; some models are moderated and may refuse certain lines (they stay untranslated and are reported).
- **TMDB** (only if you add a key) receives the title parsed from the file name.
- **OpenSubtitles** (only if you add a key, and only when you open its search) receives the title, season and episode parsed from the file name, a fingerprint of the video file (its size and a checksum of its first and last 64 KiB), and, if you signed in, your account credentials. The video itself is never uploaded.
- **GitHub** is contacted to read the sidecar manifest and to check for app updates. The tools themselves are downloaded from their publishers' sites listed in the manifest (mkvtoolnix.download, ffmpeg.martin-riedl.de, gyan.dev), and verified against the SHA-256 in the manifest.
- There is no telemetry. Logs stay on your machine; the exported log file hides your API key and replaces your home folder with `~`.

## Where files are stored

Settings, cache, logs and downloaded tools live in your user data folder (Linux: `~/.config/bakasub`, Windows: `%APPDATA%\bakasub`). Uninstalling or updating the app does not delete it. Settings → About shows the exact path. Development builds use a separate `bakasub-dev` folder (see [`development.md`](./development.md#dev-data-folder)). Output files go where you chose in the Output step.

## Troubleshooting

| Symptom                                   | What to do                                                                                                                        |
| ----------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| "Video tools missing" / `SIDECAR_MISSING` | Install ffmpeg and mkvtoolnix in Settings → Video tools.                                                                          |
| Invalid key / no credits                  | Check the key in Settings → Keys; add credits at OpenRouter.                                                                      |
| Daily limit reached on a free model       | Free models have a daily quota (resets 00:00 UTC). Use another model or wait.                                                     |
| Many lines untranslated                   | Try a model with reliable JSON output; the lines already translated stay in the cache, so retrying only pays for what is missing. |
| Wrong characters (`Ã©`, `?`)              | The file encoding was probably misdetected. Re-save it as UTF-8 and try again.                                                    |
| Job is "interrupted"                      | The app closed while it ran. Open the job and retry; it resumes from the step that stopped.                                       |
| Linux: `chrome-sandbox` error             | Only for the unpacked build. Use the AppImage, or run with `--no-sandbox`.                                                        |
| Windows: SmartScreen warning              | The installer is not code-signed. Choose "More info" → "Run anyway" after checking `SHA256SUMS.txt`.                              |
| Auto-update does nothing on Linux         | Updates only work when running the AppImage.                                                                                      |

Still stuck? See [`reporting-issues.md`](./reporting-issues.md).
