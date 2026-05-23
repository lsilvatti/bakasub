# BakaSub - The Full App 🌸✨

[Português   do Brasil](README-pt.md)

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/leosilvatto)

Hmph! So you found the real **BakaSub** repository, not just one of its nerdy internal organs? Good. This is the one end users should actually run on their own machine. It brings the frontend, backend, and PostgreSQL together through Docker Compose so you can translate subtitles without babysitting a Go server, a React app, and a database separately. B-baka!

> **Listen up!** If you just want to use Bakasub, this is the only repository you need. The backend and frontend repositories are developer-facing.

Need a precise split between user settings and infrastructure wiring? See [CONFIGURATION.md](CONFIGURATION.md).

## ✨ What Bakasub Actually Does

Bakasub is a local subtitle workflow app focused on AI-assisted translation and video subtitle handling.

It lets you:

- inspect subtitle tracks from your video files
- extract subtitle tracks from MKV and other supported containers
- translate subtitle files with OpenRouter models
- enrich the translation context with TMDB metadata and your own notes
- apply tuned presets for anime, movies, documentaries, and other contexts
- reuse cached translations to avoid paying twice for the same lines
- monitor jobs, costs, tokens, and logs in real time
- merge translated subtitles back into your video workflow

## 🌸 Features You Get

- **One-command stack**: frontend, backend, and PostgreSQL start together with Docker Compose.
- **Modern web UI**: no Electron shell, no separate desktop installer drama.
- **Live translation progress**: the interface streams job progress while the backend works.
- **Translation memory**: repeated lines can come straight from cache instead of burning more API credits.
- **Pre-flight estimates**: preview batches, token usage, and estimated cost before sending the job.
- **TMDB-aware context**: attach metadata from movies or series to improve translation quality.
- **Preset-driven workflows**: keep separate translation styles for anime, films, comedy, and more.
- **Logs and job history**: inspect what happened without digging through container internals like a gremlin.

## 🧰 What You Need Before You Start

You do **not** need to install Go, Node.js, or PostgreSQL on the host machine just to use the product.

You do need:

1. **Docker Engine** with the Compose plugin, or **Docker Desktop**.
2. A **released version** of this repository. Prefer a GitHub Release or a tagged version, not a random development commit.
3. A local **`library/` folder** inside this repository containing the videos and subtitles you want Bakasub to see.
4. An **OpenRouter API key**.
5. A **TMDB access token**.

Important:

- Translation stays locked until the OpenRouter key and TMDB token are configured and saved in the app Settings page.
- In the official Docker setup from this repository, FFmpeg and MKVToolNix already ship inside the backend image, so you should not need to install them on the host just to get started.

## 🚀 Quick Start

### 1. Get a real release of Bakasub

Use a tagged release of this repository.

Why? Because public releases pin the backend and frontend container images in `release.env`. A work-in-progress checkout may still contain placeholder image digests.

### 2. Put your media inside `library/`

Create the local library folder if it does not exist yet:

```bash
mkdir -p library
```

Then copy or symlink your media into it. Example:

```bash
ln -s /absolute/path/to/your/videos library/videos
```

The official product stack mounts `./library` into the backend container as `/videos`. The app-level browsing root is then saved in the database through **Settings**, not via environment variables.

### 3. Start the application

```bash
sh scripts/up.sh
```

This starts:

- the frontend on an automatically assigned local port
- the backend behind the frontend proxy at `/api/v1`
- the internal PostgreSQL database used by Bakasub

### 4. Open Bakasub in the browser

The script prints the exact local URL after startup.

### 5. Finish the first-time setup inside the app

Open the **Settings** page and configure:

- `OpenRouter API Key`
- `TMDB Access Token`
- `Library Root` if you want to limit the browser to a subfolder such as `/videos/videos`

Save both values. Until they are saved, the **Translate** page stays locked.

### 6. Start using the workflow

The normal flow is:

1. **Extract**: pick a video and inspect or extract subtitle tracks.
2. **Translate**: choose a subtitle file, a model, a preset, and a target language.
3. **Merge**: mux the translated result back into your video workflow when needed.
4. **Logs & Jobs**: review progress, errors, cost, and history.

## 🧭 A Good First Session

If you want the cleanest first run, do this:

1. Start the stack.
2. Open **Settings** and save your OpenRouter and TMDB credentials.
3. Confirm the **Library Root** in Settings. The default is `/videos`.
4. Go to **Models & Presets** and mark your favorite models.
5. Open **Extract** and inspect a video from the mounted library.
6. Open **Translate**, link TMDB metadata if helpful, run a pre-flight estimate, then submit the translation.
7. Watch the live progress dialog and confirm the result in **Logs & Jobs**.

## 🗂️ What Bakasub Can See

Bakasub only sees the folder mounted as `./library`, which appears inside the container as `/videos`.

That means:

- if a file is outside `library/`, it will not show up in the app
- if the library is empty, Bakasub will behave like your media collection is empty
- if the folder permission is blocked, the backend cannot inspect or process the files

If you do not want to copy files, use symlinks from `library/` to your actual media directories.

## 🔄 Updating Bakasub

When a new public release comes out:

1. Pull or download the new tagged version of this repository.
2. Keep your existing `library/` folder.
3. Start the updated stack again:

```bash
sh scripts/up.sh
```

If you want to be explicit, you can pull first:

```bash
docker compose --env-file release.env pull
sh scripts/up.sh
```

## 🛑 Stopping or Removing the Stack

Stop the containers without removing persistent data:

```bash
docker compose --env-file release.env down
```

Remove everything, including the PostgreSQL volume used by the app:

```bash
docker compose --env-file release.env down -v
```

Do the second one only if you really want to wipe the local app database.

## 🩹 Troubleshooting

### The frontend does not open

- Run `docker compose --env-file release.env port frontend 80` and open the printed URL.
- If the container failed to start, inspect `docker compose --env-file release.env logs`.

### The backend port is busy

- The backend is no longer exposed directly on the host in the product stack.
- Reach it through the frontend proxy using the URL printed by `sh scripts/up.sh`.

### Translation is locked

- Open **Settings**.
- Save a valid **OpenRouter API Key**.
- Save a valid **TMDB Access Token**.

The product expects both to be configured before translation is enabled.

### I cannot see my video files

- Confirm that the files really exist inside `library/` or inside a symlink created from `library/`.
- Confirm that Docker can read that folder.
- Confirm that the **Library Root** saved in Settings points to the mounted location you expect.

### The stack fails to start

- Make sure you are using a tagged or released version with real image digests in `release.env`.
- Run `docker compose --env-file release.env logs` to inspect what failed.

## 💖 Support The Project

If Bakasub saved your weekend release schedule, subtitle backlog, or general sanity, you can support the project here:

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/lsilvatti)