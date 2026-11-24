<h1 align="center">
  <img src="./docs/bakasub-logo-en.svg" alt="BakaSub" width="560" />
</h1>

🇬🇧 **English** · 🇧🇷 [Português](./README-pt.md)

**Extract, translate with AI and re-mux the subtitles of your `.mkv` videos**, from a desktop app, on your own machine.

Windows x64 and Linux x64 · [GPL-3.0](./LICENSE) · [Website](https://bakasub.silvatti.com.br/)

<p align="center">
  <a href="https://buymeacoffee.com/leosilvatto"><img src="https://img.shields.io/badge/Buy%20me%20a%20coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Buy me a coffee" /></a>
</p>

## What it does

- **Translates** `.srt`, `.ass`, `.ssa` and `.vtt` with any model from [OpenRouter](https://openrouter.ai). Only the text goes to the model: timestamps, styles and tags come back exactly where they were.
- **Handles the whole MKV workflow**: extract a subtitle track, translate it, and mux it back into the video (or just save the subtitle file). Use each step on its own or run them all in one go.
- **Shows the cost before you pay**: a token count and a price estimate from live OpenRouter pricing, plus a list of what will _not_ be translated.
- **Remembers what it translated**: a local cache makes reprocessing a file nearly free.
- **Works when the video has no subtitle**: bring your own subtitle file or search [OpenSubtitles](https://www.opensubtitles.com/) from the app, then translate it like any other.
- **Knows what you are watching**: it reads the release name and can pull the synopsis and episode details from TMDB to improve the translation.
- **Keeps your terms consistent** with presets (Movie, Series, Anime, Documentary, or your own) and a glossary for names, honorifics and jargon.
- **Works in bulk**: translate, extract or mux a whole folder or season at once.
- English and Brazilian Portuguese interface. No account, no telemetry; your key and your files stay on your computer.

## Install

Download the latest version from the [Releases page](https://github.com/lsilvatti/bakasub/releases):

| System  | File                    | Notes                                                                                                               |
| ------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------- |
| Windows | `bakasub-<v>-setup.exe` | Windows may show a SmartScreen warning because the installer is not code-signed. Choose "More info" → "Run anyway". |
| Linux   | `bakasub-<v>.AppImage`  | `chmod +x bakasub-<v>.AppImage`, then run it. Automatic updates only work with the AppImage.                        |

You can check a download against `SHA256SUMS.txt` from the same release.

## Getting started

You need an **OpenRouter API key**: create one at [openrouter.ai/keys](https://openrouter.ai/keys). Some models are free (with a daily limit); others are paid per use.

1. Open BakaSub. The first-run guide asks for your OpenRouter key, a favorite and default model, and your target languages.
2. Let it download **ffmpeg and mkvtoolnix** when asked (needed for `.mkv`; it only downloads with your consent). Subtitle files alone do not need them.
3. Optionally add a free [TMDB](https://www.themoviedb.org/settings/api) key to enrich translations with synopses.
4. Go to **Translate**, drop a video or a subtitle file, pick a preset, review the cost and confirm.

Progress, cost and results are in the **Queue**. A full walkthrough is in the [user guide](./docs/user-guide.md).

## Good to know

- **Cost and privacy.** You pay OpenRouter directly for what you translate. The subtitle text, your notes and the glossary are sent to OpenRouter; the file name is sent to TMDB if you add a key, and to OpenSubtitles (with a fingerprint of the video file) if you add its key and search there. Nothing is sent without an action of yours, and nothing goes to the author. Details in the [user guide](./docs/user-guide.md#privacy-what-leaves-your-computer).
- **Your data** stays in your user data folder. Updating or uninstalling the app does not delete it. Development builds (`npm run dev`) use a separate folder, so they never touch the installed app's data:

  |         | Installed           | Development             |
  | ------- | ------------------- | ----------------------- |
  | Linux   | `~/.config/bakasub` | `~/.config/bakasub-dev` |
  | Windows | `%APPDATA%\bakasub` | `%APPDATA%\bakasub-dev` |

  Database, logs, downloaded tools and cache are all inside it. `BAKASUB_USER_DATA=<dir>` overrides the folder (tests, scripts). Details in [Development](./docs/development.md#dev-data-folder).

- **Not supported**: image-based subtitles (PGS, VobSub, DVB) and macOS.
- **Known limitations**: no overview screen for a batch after you leave its page (its jobs show one by one in the Queue); no manual encoding override, so a badly detected file should be re-saved as UTF-8.

## Documentation

- **Using the app**: [User guide](./docs/user-guide.md) · [Reporting a problem](./docs/reporting-issues.md)
- **Contributing**: [CONTRIBUTING.md](./CONTRIBUTING.md) · [Development](./docs/development.md) · [Architecture](./docs/architecture.md)
- **Internals** (English): [IPC](./docs/ipc.md) · [Sidecars](./docs/sidecars.md) · [Database](./docs/database.md) · [Subtitle parsing](./docs/parsing.md) · [Job Memory](./docs/job-memory.md) · [LLM translation](./docs/llm.md) · [Jobs](./docs/jobs.md) · [Preflight](./docs/preflight.md) · [Rate limits](./docs/rate-limits.md) · [Presets](./docs/presets.md) · [Media context](./docs/media-context.md) · [External subtitles](./docs/subtitle-sources.md) · [Logging](./docs/logging.md) · [Renderer](./docs/renderer.md) · [Brand tokens](./docs/brand-tokens.md) · [Website](./docs/site.md) · [Build and release](./docs/release.md)

## Support and contributing

Found a bug or have an idea? [Open an issue](https://github.com/lsilvatti/bakasub/issues/new/choose); the [guide](./docs/reporting-issues.md) says what to include. Pull requests are welcome after a conversation in an issue: see [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

BakaSub is free software under the [GNU GPL v3](./LICENSE). ffmpeg and mkvtoolnix are separate programs, downloaded by you and run as external processes; they keep their own licenses. Movie and series data comes from TMDB (attribution is shown in the app); this product uses the TMDB API but is not endorsed or certified by TMDB.
