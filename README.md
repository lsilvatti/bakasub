# 🍜 BakaSub

> *"It's not like I made this subtitle tool for you or anything... B-Baka!"*

**BakaSub** is a blazing-fast, AI-powered subtitle translation tool built for power users who demand **zero desync** and **native terminal aesthetics**. Born from frustration with clunky web interfaces and subtitle timing disasters, BakaSub brings professional-grade translation automation to your terminal.

Think of it as `btop` meets `lazygit`, but for subtitles. No mouse required, no bloated GUI, just pure keyboard-driven efficiency.

## ✨ Features

- **🤖 AI-Powered Translation**: OpenRouter, Google Gemini, OpenAI, or local LLM support
- **⚡ Zero Desync Protocol**: Sliding window context + quality gates ensure perfect timing
- **💾 Smart Caching**: SQLite-backed fuzzy matching saves you money on repeated translations
- **🎨 Native Neon TUI**: btop-inspired interface that looks *chef's kiss* in your terminal
- **📦 Binary First**: Single executable, no dependencies (except FFmpeg/MKVToolNix)
- **🔄 Watch Mode**: Drop files in a folder, walk away, let BakaSub handle it
- **🛠️ MKV Toolbox**: Extract, mux, edit headers, manage fonts - all in one place
- **🌍 Trilingual**: Interface available in English, Português (BR), and Español

### Why BakaSub?

| 💀 Old Way | ✨ BakaSub Way |
|-----------|---------------|
| Export subtitles manually | Auto-extracts from MKV |
| Copy-paste into web translator | Batch API calls with context |
| Fix timing desync for 2 hours | Anti-desync protocol built-in |
| Manually remux into video | One-step muxing with backups |
| Hope you didn't mess up | Quality gate catches errors |

## 🚀 Installation

### Quick Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Manual Installation

1. **Download** the latest release for your platform:
   - [Linux (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64)
   - [Windows (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe)
   - [macOS (Intel)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64)
   - [macOS (Apple Silicon)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64)

2. **Make executable** (Linux/macOS):
   ```bash
   chmod +x bakasub-*
   sudo mv bakasub-* /usr/local/bin/bakasub
   ```

3. **Verify installation**:
   ```bash
   bakasub --version
   ```

### Dependencies

BakaSub needs these external tools (the wizard will offer to download them):

- **FFmpeg**: Media processing
- **MKVToolNix**: Container manipulation

## 🎬 Quick Start

### First Run (Setup Wizard)

On first launch, BakaSub walks you through:

1. **AI Provider Setup**: Choose your service (OpenRouter recommended) and enter API key
2. **Dependency Check**: Auto-downloads FFmpeg/MKVToolNix if missing
3. **Defaults**: Set your target language and preferred model

```bash
bakasub
```

### Basic Workflow: Full Process Mode

The most common use case - translate everything in one shot:

1. Launch BakaSub
2. Enter the path to your MKV file or folder
3. Select **"Full Process"** mode
4. Press **Enter** to start
5. Grab coffee while BakaSub does its magic ☕

### Watch Mode (Set It and Forget It)

Perfect for automation or batch processing:

1. Create a folder (e.g., `~/anime-incoming`)
2. In BakaSub, select **"Watch Mode"**
3. Point it to your folder
4. Drop files into the folder
5. BakaSub auto-processes new files as they appear

*Like a responsible adult's download folder, but it actually cleans itself up.*

## ⌨️ Key Bindings

### Dashboard

| Key | Action |
|-----|--------|
| `1-4` | Launch modules (Extract, Translate, Mux, Review) |
| `5-8` | Open toolbox (Header Editor, Glossary, etc.) |
| `m` | Change AI model |
| `c` | Open configuration |
| `q` | Quit |

### Job Setup

| Key | Action |
|-----|--------|
| `Enter` | Start job |
| `d` | Dry run (cost estimate) |
| `r` | Resolve track conflicts |
| `Esc` | Back to dashboard |

### Manual Review Editor

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate lines |
| `Enter` | Commit edit & next |
| `Ctrl+S` | Save file |
| `g` | Go to line number |
| `Esc` | Exit editor |

### Header Editor

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate tracks |
| `Space` | Toggle flags (Default/Forced) |
| `Enter` | Apply changes |
| `Esc` | Cancel |

## 🎭 Configuration

Config lives at `~/.config/bakasub/config.json`. Key settings:

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

BakaSub ships with specialized prompts for different content types:

- **Anime**: Preserves honorifics (-san, -kun), keeps attack names
- **Movie**: Formal tone, localized idioms
- **Series**: Balanced style for episodic content
- **Documentary**: Technical accuracy over creativity
- **YouTube**: Casual tone, internet slang-aware

You can clone factory profiles and customize them.

## 🛠️ Toolbox Modules

### Standalone Operations

1. **Extract Tracks**: Rip subtitles/audio from MKV
2. **Translate Subtitle**: AI translation with your settings
3. **Mux Container**: Combine tracks into MKV
4. **Manual Review**: Split-view editor for corrections

### MKVToolNix Tools

5. **Edit Flags/Metadata**: Set default tracks, forced subs
6. **Manage Attachments**: Add/remove fonts from MKV
7. **Add/Remove Tracks**: Quick remuxer with track selection
8. **Project Glossary**: Define terms for consistent translation

## 🌍 Localization

BakaSub's interface supports:

- 🇬🇧 **English** (default)
- 🇧🇷 **Português (Brasil)**
- 🇪🇸 **Español**

Switch in `Configuration > General > Interface Language`.

## 🐛 Troubleshooting

### "API Error 401"

Your API key is invalid or expired. Run `bakasub` → `c` (config) → AI Providers → re-enter key.

### "Track Conflict Detected"

Multiple subtitle tracks match your target language. BakaSub needs you to pick:
- Press `r` in Job Setup
- Select the **full dialogue** track (usually the larger file size)
- Signs/Songs tracks are typically smaller

### "FFmpeg Not Found"

Install FFmpeg:
- **Ubuntu/Debian**: `sudo apt install ffmpeg`
- **macOS**: `brew install ffmpeg`
- **Windows**: Download from [ffmpeg.org](https://ffmpeg.org)

Or let the Setup Wizard download it for you.

### Subtitles Are Desync'd

This should NEVER happen thanks to our anti-desync protocol. If it does:
1. Check that you selected the correct subtitle track (Signs/Songs ≠ Full Dialogue)
2. Verify the source MKV isn't already corrupt (`mkvmerge -i file.mkv`)
3. Open a GitHub issue with the file info

## 🤝 Contributing

Found a bug? Want a feature? Contributions are welcome!

1. Fork the repo
2. Create a feature branch (`git checkout -b cool-feature`)
3. Commit your changes (`git commit -am 'Add cool feature'`)
4. Push to the branch (`git push origin cool-feature`)
5. Open a Pull Request

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

## 💖 Support

Like BakaSub? Consider supporting development:

- ⭐ Star the repo
- ☕ [Buy me a coffee](https://ko-fi.com/lsilvatti) *(we accept headpats too)*
- 📢 Share with friends who suffer from subtitle hell

---

**Made with 💜 by someone who watched too much anime with terrible subtitles**

*"Omae wa mou... translated." - BakaSub, probably*