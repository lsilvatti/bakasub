# BakaSub - The Full App 🌸✨

[Português do Brasil](README-pt.md)

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/leosilvatto)

Hmph! So you found the real **BakaSub** repository, not just one of its nerdy internal organs? Good. This is the one end users should actually run on their own machine. It brings the frontend, backend, and PostgreSQL together through Docker Compose so you can translate subtitles without babysitting a Go server, a React app, and a database separately. B-baka!

> **Listen up!** If you just want to use Bakasub, this is the only repository you need. The backend and frontend repositories are developer-facing.

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
3. An **absolute path** on your machine containing the videos and subtitles you want Bakasub to see.
4. An **OpenRouter API key**.
5. A **TMDB access token**.

Important:

- Translation stays locked until the OpenRouter key and TMDB token are configured and saved in the app Settings page.
- In the official Docker setup from this repository, FFmpeg and MKVToolNix already ship inside the backend image, so you should not need to install them on the host just to get started.

## 🚀 Quick Start

### 1. Get a real release of Bakasub

Use a tagged release of this repository.

Why? Because public releases pin the backend and frontend container images in `release.env`. A work-in-progress checkout may still contain placeholder image digests.

### 2. Create your local runtime config

Copy the example file:

```bash
cp .env.example .env
```

Then edit `.env` and set at least this value:

```env
VIDEO_DIR=/absolute/path/to/your/videos
```

Optional settings you can change there too:

- `FRONTEND_HOST_PORT` if `3000` is busy
- `BACKEND_HOST_PORT` if `8080` is busy
- `POSTGRES_DB`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` if you want custom local database credentials for the composed stack

### 3. Start the application

```bash
docker compose --env-file release.env --env-file .env up -d
```

This starts:

- the frontend at `http://localhost:3000`
- the backend API at `http://localhost:8080/api/v1`
- the internal PostgreSQL database used by Bakasub

### 4. Open Bakasub in the browser

Go to:

```text
http://localhost:3000
```

### 5. Finish the first-time setup inside the app

Open the **Settings** page and configure:

- `OpenRouter API Key`
- `TMDB Access Token`

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
3. Go to **Models & Presets** and mark your favorite models.
4. Open **Extract** and inspect a video from your mounted `VIDEO_DIR`.
5. Open **Translate**, link TMDB metadata if helpful, run a pre-flight estimate, then submit the translation.
6. Watch the live progress dialog and confirm the result in **Logs & Jobs**.

## 🗂️ What Bakasub Can See

Bakasub only sees the folder mounted as `VIDEO_DIR`.

That means:

- if a file is outside that directory, it will not show up in the app
- if the path is wrong, Bakasub will behave like your library is empty
- if the folder permission is blocked, the backend cannot inspect or process the files

Use an absolute path. Relative paths are not enough here.

## 🔄 Updating Bakasub

When a new public release comes out:

1. Pull or download the new tagged version of this repository.
2. Keep your existing `.env` file.
3. Start the updated stack again:

```bash
docker compose --env-file release.env --env-file .env up -d
```

If you want to be explicit, you can pull first:

```bash
docker compose --env-file release.env --env-file .env pull
docker compose --env-file release.env --env-file .env up -d
```

## 🛑 Stopping or Removing the Stack

Stop the containers without removing persistent data:

```bash
docker compose --env-file release.env --env-file .env down
```

Remove everything, including the PostgreSQL volume used by the app:

```bash
docker compose --env-file release.env --env-file .env down -v
```

Do the second one only if you really want to wipe the local app database.

## 🩹 Troubleshooting

### The frontend does not open

- Check whether port `3000` is already in use.
- If needed, change `FRONTEND_HOST_PORT` in `.env` and run the stack again.

### The backend port is busy

- Change `BACKEND_HOST_PORT` in `.env`.
- Restart with the same `docker compose --env-file release.env --env-file .env up -d` command.

### Translation is locked

- Open **Settings**.
- Save a valid **OpenRouter API Key**.
- Save a valid **TMDB Access Token**.

The product expects both to be configured before translation is enabled.

### I cannot see my video files

- Confirm that `VIDEO_DIR` points to the correct absolute path.
- Confirm that Docker can read that directory.
- Confirm that the files really live inside that mounted folder.

### The stack fails to start

- Make sure you are using a tagged or released version with real image digests in `release.env`.
- Run `docker compose --env-file release.env --env-file .env logs` to inspect what failed.

## 💖 Support The Project

If Bakasub saved your weekend release schedule, subtitle backlog, or general sanity, you can support the project here:

[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=000000)](https://buymeacoffee.com/lsilvatti)