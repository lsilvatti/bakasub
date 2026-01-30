# 🍜 BakaSub

> *"It's not like I made this subtitle tool for you or anything... B-Baka!"*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

**BakaSub** is an AI-powered subtitle translation tool for power users who demand **zero desync** and **native terminal aesthetics**. Born from frustration with clunky web interfaces and subtitle timing disasters.

Think `btop` meets `lazygit`, but for subtitles. No mouse, no bloat—just keyboard-driven efficiency.

---

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Dependencies](#-dependencies)
- [Quick Start](#-quick-start)
- [Usage Guide](#-usage-guide)
- [Configuration](#-configuration)
- [Troubleshooting](#-troubleshooting)
- [For Developers](#-for-developers)
- [Support](#-support)

---

## ✨ Features

| Feature | What it does |
|---------|-------------|
| 🤖 **AI Translation** | Supports OpenRouter, Google Gemini, OpenAI, and local LLMs (Ollama/LMStudio) |
| ⚡ **Zero Desync** | Sliding window context + quality gates keep your subs perfectly timed |
| 💾 **Smart Cache** | SQLite-backed fuzzy matching—why pay twice for the same line? |
| 🎨 **Neon TUI** | A terminal interface so pretty you'll forget GUIs exist |
| 📦 **Single Binary** | One file, no Python, no Node, no drama |
| 🔄 **Watch Mode** | Drop files in a folder, BakaSub handles the rest. Magic! ✨ |
| 🛠️ **MKV Toolbox** | Extract, mux, edit headers, manage fonts—all in one place |
| 🌍 **Trilingual UI** | English, Português (BR), Español |

---

## 🚀 Installation

### One-Line Install (Linux/macOS)

*"F-Fine, I'll make it easy for you... but only this once!"*

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Manual Download

Pick your platform, download, and you're done:

| Platform | Download Link |
|----------|---------------|
| 🐧 Linux (AMD64) | [bakasub-linux-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64) |
| 🪟 Windows (AMD64) | [bakasub-windows-amd64.exe](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe) |
| 🍎 macOS (Intel) | [bakasub-darwin-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64) |
| 🍎 macOS (Apple Silicon) | [bakasub-darwin-arm64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64) |

**Linux/macOS setup:**
```bash
chmod +x bakasub-*
sudo mv bakasub-* /usr/local/bin/bakasub
bakasub --version  # Verify it works!
```

**Windows:** Just put the `.exe` somewhere in your PATH or run it directly.

---

## 🔧 Dependencies

BakaSub needs two external tools. *"D-Don't look at me like that! You have to install them yourself... it's not like I can do everything for you!"*

**You MUST install these before running BakaSub:**

| Tool | What it does | Download |
|------|-------------|----------|
| **FFmpeg** | Media processing, stream extraction | [ffmpeg.org](https://ffmpeg.org/download.html) |
| **MKVToolNix** | MKV container manipulation | [mkvtoolnix.download](https://mkvtoolnix.download/downloads.html) |

### Quick Install Commands

**Ubuntu/Debian:**
```bash
sudo apt install ffmpeg mkvtoolnix
```

**Fedora:**
```bash
sudo dnf install ffmpeg mkvtoolnix
```

**Arch Linux:**
```bash
sudo pacman -S ffmpeg mkvtoolnix-cli
```

**macOS (Homebrew):**
```bash
brew install ffmpeg mkvtoolnix
```

**Windows:** Download installers from the links above, or use [Chocolatey](https://chocolatey.org/):
```powershell
choco install ffmpeg mkvtoolnix
```

---

## 🎬 Quick Start

### First Run

```bash
bakasub
```

On first launch, a setup wizard guides you through:

1. **AI Provider** — Choose your service and enter your API key
2. **Dependency Check** — Verifies FFmpeg and MKVToolNix are installed
3. **Defaults** — Set your target language and preferred model

*"I-I'm only helping because you clearly can't do it yourself!"*

### Basic Workflow

**Full Process Mode** — The most common use case:

1. Launch `bakasub`
2. Enter path to your MKV file or folder
3. Select **Full Process**
4. Press **Enter**
5. ☕ Grab coffee. You earned it.

**Watch Mode** — Set it and forget it:

1. Create a folder (e.g., `~/anime-incoming`)
2. Select **Watch Mode** in BakaSub
3. Point to your folder
4. Drop MKV files there anytime
5. BakaSub auto-processes new files as they appear

*Like a responsible download folder that actually cleans itself up.*

---

## 📖 Usage Guide

### Dashboard Keys

| Key | Action |
|-----|--------|
| `1` | Extract tracks from MKV |
| `2` | Translate subtitle file |
| `3` | Mux tracks into MKV |
| `4` | Manual review editor |
| `5` | Edit track flags/metadata |
| `6` | Manage attachments (fonts) |
| `7` | Quick remuxer |
| `8` | Project glossary |
| `m` | Change AI model |
| `c` | Open configuration |
| `q` | Quit |

### Job Setup Keys

| Key | Action |
|-----|--------|
| `Enter` | Start the job |
| `d` | Dry run (cost estimate without calling API) |
| `r` | Resolve track conflicts |
| `Esc` | Back to dashboard |

### Review Editor Keys

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate lines |
| `Enter` | Confirm edit, move to next |
| `Ctrl+S` | Save file |
| `g` | Go to line number |
| `Esc` | Exit editor |

### Toolbox Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | **Extract Tracks** | Rip subtitles or audio from MKV |
| 2 | **Translate Subtitle** | AI translation with your settings |
| 3 | **Mux Container** | Combine tracks into a new MKV |
| 4 | **Manual Review** | Split-view editor for corrections |
| 5 | **Header Editor** | Set default/forced track flags |
| 6 | **Attachments** | Add/remove fonts from MKV |
| 7 | **Remuxer** | Quick track add/remove |
| 8 | **Glossary** | Define terms for consistent translation across episodes |

---

## 🎭 Configuration

Config lives at `~/.config/bakasub/config.json`

```json
{
  "api_provider": "openrouter",
  "api_key": "sk-or-...",
  "target_lang": "pt-br",
  "remove_hi_tags": true,
  "global_temp": 0.3,
  "touchless_mode": false,
  "prompt_profile": "anime"
}
```

### Prompt Profiles

Different content needs different translation styles:

| Profile | Best for |
|---------|----------|
| **anime** | Preserves honorifics (-san, -kun), keeps attack names |
| **movie** | Formal tone, localized idioms |
| **series** | Balanced style for episodic content |
| **documentary** | Technical accuracy over creativity |
| **youtube** | Casual tone, internet slang-aware |

Clone factory profiles to customize them. *"I made defaults, but you can change them... if you think you know better!"*

### Interface Language

BakaSub supports: 🇬🇧 English (default) · 🇧🇷 Português · 🇪🇸 Español

Change in `Configuration > General > Interface Language`

---

## 🐛 Troubleshooting

### "API Error 401"

Your API key is invalid or expired.

→ Press `c` → AI Providers → Re-enter your key

### "Track Conflict Detected"

Multiple subtitle tracks match your language. BakaSub needs you to pick:

→ Press `r` in Job Setup  
→ Select the **full dialogue** track (usually larger file size)  
→ Signs/Songs tracks are typically smaller

### "FFmpeg Not Found"

Install FFmpeg using the commands in the [Dependencies](#-dependencies) section above.

*"I literally gave you the commands... just copy and paste them! Baka!"*

### Subtitles Are Desync'd

*"This should NEVER happen. My code is perfect!"* ...but if it does:

1. Check you selected the correct track (Signs/Songs ≠ Full Dialogue)
2. Verify source MKV isn't corrupt: `mkvmerge -i file.mkv`
3. [Open an issue](https://github.com/lsilvatti/bakasub/issues) with file info

---

## 👨‍💻 For Developers

*"Oh, you want to contribute? H-How bold of you..."*

### Building from Source

**Requirements:** Go 1.22+

```bash
git clone https://github.com/lsilvatti/bakasub.git
cd bakasub
go mod download
```

### Build Commands

```bash
make build-linux     # Linux AMD64
make build-windows   # Windows AMD64
make build-macos     # macOS Intel + ARM
make build-all       # All platforms
make install         # Build + install to /usr/local/bin
```

### Development

```bash
make dev    # Run without building
make test   # Run tests
make fmt    # Format code
make lint   # Run linter
```

### Contributing

1. Fork the repo
2. Create a feature branch: `git checkout -b cool-feature`
3. Commit your changes: `git commit -am 'Add cool feature'`
4. Push: `git push origin cool-feature`
5. Open a Pull Request

---

## 📜 License

MIT License — Do whatever you want, just don't blame me.

---

## 💖 Support

*"I-It's not like I need your support or anything... but if you insist..."*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

- ⭐ Star this repo
- 📢 Share with friends suffering from subtitle hell
- 🐛 Report bugs (but be nice about it!)

---

**Made with 💜 by someone who watched too much anime with terrible subtitles**

*"Omae wa mou... translated." — BakaSub, probably*
