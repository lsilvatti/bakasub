# External subtitle sources

Some MKVs have no subtitle track that can be translated: they have none, or only image-based ones. For those, the single-file Translate wizard accepts a subtitle from outside the video, either a file from disk or a download from [OpenSubtitles](https://www.opensubtitles.com/). The same option appears under the track list when the video does have tracks.

Batch mode does not use external subtitles. Items without a usable track are still skipped (`no-subtitle-track` / `no-supported-track`).

## How the subtitle reaches the job

OpenSubtitles never runs inside a job. The wizard searches and downloads, and the download ends as a local file. The job receives a path to a subtitle file and nothing else.

```
Track step ─┬─ track of the MKV ────────────► pipeline { trackId }
            ├─ file from disk ──────────────► pipeline { subtitle: { path } }
            └─ OpenSubtitles: search → download → userData/cache/opensubtitles/<fileId>.srt
                                              └► pipeline { subtitle: { path, language } }
```

Because of this, a job works offline and a retry never spends download quota. The LLM is still the only network call a job makes.

## The pipeline with an external subtitle

A `pipeline` request takes exactly one of `trackId` and `subtitle: { path, language? }`. Sending both or neither is `INVALID_INPUT`. `resolveJobSpec` checks this; the Zod union does not. The path goes through `validateInputFile`, like every path from the renderer: absolute, a real file, a subtitle extension, at most 50 MB.

The frozen spec stores `externalSubtitle: { path, size, format, includeInMux }` and no `trackId`. Specs saved before this feature always have `trackId` and keep working.

| Step                     | With an external subtitle                                                                               |
| ------------------------ | ------------------------------------------------------------------------------------------------------- |
| `probe`                  | Still runs. It validates the container and lists the original subtitle tracks for the mux.              |
| `extract`                | Left out of the plan by `planFor`.                                                                      |
| `preflight`, `translate` | `translationInput` reads `externalSubtitle.path` directly, as the `translate` kind reads its source.    |
| `write`                  | Unchanged. The output name comes from the video and the target language.                                |
| `mux`                    | Adds the translated track as usual. With `includeInMux` it also adds the external subtitle (see below). |

The source language is the request's `sourceLang`. Without it, the job uses `subtitle.language`, which a download always has. Without that, it uses the language tag in the file name (`ep01.en.srt`).

`options.includeSourceSubtitle` adds the external subtitle to the MKV. It is off by default and shows in the Output step as "Also add the original subtitle to the video". It needs an MKV to be written, so `outputMode: 'subtitle'` turns it off. The track goes after the translated one and is never the default. It is tagged with the source language, or `und` and the name `Original` when the language is unknown. `mkvmerge` receives the file as is, like the subtitles of a `mux` job.

The job does not copy the subtitle into its work directory. If the file is moved or the cache is cleared before the job runs, the job fails when it tries to read it.

## OpenSubtitles (`src/main/subtitle-sources/`)

| File                               | Role                                                                                              |
| ---------------------------------- | ------------------------------------------------------------------------------------------------- |
| `opensubtitles-client.ts`          | Stateless client of the REST API (`https://api.opensubtitles.com/api/v1`), with injected `fetch`. |
| `opensubtitles-session.ts`         | Optional user session. Keeps the token in memory and logs in again once if it is rejected.        |
| `movie-hash.ts`                    | OSDb hash of the video file.                                                                      |
| `subtitle-source.service.ts`       | Search, download, validation and cache.                                                           |
| `create-subtitle-source-module.ts` | Composition root (vault, cache dir, `User-Agent`).                                                |

### Credentials

The API key is required. It is the user's own consumer key, stored as `opensubtitles.apiKey` in the `SecretVault` and sent as `Api-Key`. The app ships no key.

Before saving a key, the client checks it with a search for a random query that is never repeated. A fixed query would not work: the API's CDN caches `/subtitles` responses for 24 hours and ignores the key, so a query seen before answers 200 to any key. A bad key is refused with 403 only when the request reaches the origin.

The account is optional and raises the daily download quota. It is stored as `opensubtitles.login`, a JSON `{ username, password }` in the vault, and saved only after OpenSubtitles accepts the login. The password has to be kept because the session token lasts about a day. The token itself is never written to disk. The API takes the user name, and an email is refused.

The API key, the password and every token are given to the log redactor before the first request that carries them.

Every request sends `User-Agent: BakaSub v<version>`, which the API requires. It contains no user data.

### Search

`subtitles:search` takes the video path, the languages, and optionally `query`, `season` and `episode`. It sends two requests in parallel:

1. One by movie hash, which is the file size plus a 64-bit sum of the first and last 64 KiB. These results were made for this exact release, so their timing matches.
2. One by title. It uses the user's query when there is one. Otherwise it uses the title, season and episode parsed from the file name (`media-context/release-name.ts`).

The service merges the results by file id and flags the hash matches. It orders them by hash match, then human over machine translation, then trusted uploader, then download count. The response includes what was searched, so the form can show it and the user can refine it.

Query parameters are sent lower-cased and sorted, because the API redirects otherwise. Searching spends no download quota. A file smaller than 64 KiB has no hash and is searched by title only.

### Download

`subtitles:download` takes a file id and does this:

1. If `cache/opensubtitles/<fileId>.srt` exists and still parses, it returns that file. No request is made and no quota is spent.
2. It asks the API for a temporary link, always converted to SRT.
3. It fetches the link only if it is HTTPS on `opensubtitles.com` or `opensubtitles.org`, or one of their subdomains. The 50 MB cap applies while the body streams in.
4. It writes the body to `.partial` and parses it with the SRT parser. The file is renamed into the cache only if it parses.

The local file name is built from the id. The remote file name never reaches the file system.

### Errors

| Code                        | When                                                                                         |
| --------------------------- | -------------------------------------------------------------------------------------------- |
| `OPENSUBTITLES_KEY_MISSING` | The vault has no API key. No request is made.                                                |
| `OPENSUBTITLES_AUTH`        | 401 or 403, after one re-login. Also 400 on `/login`, for a wrong user name or password.     |
| `OPENSUBTITLES_QUOTA`       | 406 on `/download`. The daily download quota is used up.                                     |
| `NETWORK`                   | Unreachable, 429 or 5xx, an unexpected response shape, or a download link on another domain. |
| `INVALID_SUBTITLE`          | The downloaded file is empty, too large or does not parse.                                   |

## Known limits

BakaSub does not adjust timing. A result without a hash match may be out of sync with the video, and the UI warns about it.

Batch mode does not use external sources.

The automated tests use a fake `fetch`, and nothing in CI talks to the service. If OpenSubtitles changes its API, the app reports `NETWORK` for the unexpected response. The contract was checked by hand against the live API on 2026-11-24:

- headers, parameter ordering and response fields;
- 403 for a bad key, 401 for a bad login, 400 for an email used as user name;
- re-login after a rejected token;
- download with the API key only, limited to 5 a day, or 20 a day with a free account;
- download link on `www.opensubtitles.com`.
