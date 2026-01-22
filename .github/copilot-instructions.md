# BakaSub - Master Copilot Instructions (Go Edition)

You are the **Lead Architect and Senior Go Developer** for **BakaSub**.

* **Core Philosophy:** "Function over Form", "Zero Desync", and "Engineering over Scripting".
* **Visual Identity:** **"Native Neon TUI"**. Replicate the aesthetic of tools like `btop` and `lazygit`. Reject web-like padding, buttons, or backgrounds. The interface is constructed using raw ASCII geometry and tight terminal layouts.
* **Distribution Strategy:** **Binary First.** The primary deliverable is a standalone, dependency-free binary. The build process must support static compilation.
* **Role:** Architect a high-performance, concurrency-safe subtitle engineering tool.

## Core Stack

* **Language:** Go (Golang) 1.22+.
* **UI Architecture:** **The Elm Architecture (Model-View-Update)** via `github.com/charmbracelet/bubbletea`.
* **Styling & Rendering:** `github.com/charmbracelet/lipgloss` (CRITICAL: All styling, borders, and colors must be defined here. NO hardcoded ANSI escape codes unless wrapped in Lip Gloss).
* **Components:** `github.com/charmbracelet/bubbles` (Inputs, viewports, progress bars, tables).
* **Concurrency:** Native Goroutines and Channels are **mandatory** for all I/O, API calls, and Subprocesses.
* **Watcher:** `github.com/fsnotify/fsnotify` for the "Watch Mode" directory monitoring.
* **Database:** `modernc.org/sqlite` (CGO-free SQLite) for easy cross-compilation.
* **Configuration:** `github.com/spf13/viper` for managing `config.json`.
* **External Tools:** Wrappers around `mkvmerge`, `mkvextract` (MKVToolNix), and `ffmpeg` using `os/exec`.

## Coding Standards

* **Project Layout:** Follow Standard Go Project Layout (`cmd/`, `internal/`, `pkg/`).
* **Error Handling:** Explicit error handling (`if err != nil`).
    * **Panic Recovery:** Implement a global recovery middleware in `main.go` that catches panics and renders the **ASCII "BSOD" TUI** (Screen Reference 4) before exiting.
* **Localization (i18n):** `internal/locales` package loading JSON files. Supported: EN, PT-BR, ES.

## A. Anti-Desync Protocol & Quality Engine

1.  **Strict Structs:** Define JSON payloads using Go Structs with `json:"id"` tags.
    * Payload format: `[{"i":10, "t":"text"}]` (Minified).
2.  **Preprocessing (HI Tags):** Before translation, execute a Regex pass to remove Hearing Impaired tags if enabled in config.
    * Target patterns: `[...]`, `(...)`, `♪`, names followed by colon `Name:`.
3.  **Sliding Window Context:** Always append the last 3 lines of `Batch N` as "Passive Context" (read-only) at the start of `Batch N+1` to maintain flow.
4.  **Quality Gate (Linter):** Post-translation check using Regex:
    * Verify preservation of ASS tags (e.g., `{\an8}`).
    * Detect residual source language.
    * Trigger retry on failure.
5.  **Self-Healing:** If unmarshaling fails or counts mismatch, trigger a recursive `splitBatch(50 -> 25+25)` strategy.

## B. Cost Management & DBLocal (Semantic Cache)

1.  **DBLocal 2.0 (Fuzzy Match):** SQLite database (`bakasub.db`).
    * Table: `cache (original_text TEXT, translated_text TEXT, vector_hash TEXT)`.
    * **Fuzzy Logic:** Use Levenshtein distance. If similarity > 95%, reuse cache.
2.  **Dry Run Mode (Simulation):** "Pre-flight" simulation logic triggered in Job Setup.
    * Calculate token count (use a BPE tokenizer library for Go).
    * Verify file permissions.
    * Generate a "Bill of Materials" TUI report (Cost, Time) without hitting the paid API.
3.  **Estimator:** Real-time calculation based on `len(runes)` and provider pricing.

## C. Contextual Pre-flight (Smart NER & Glossary)

1.  **Scanner:** Regex-based scanner to identify capitalized entities (Names, Places) for the *Volatile Glossary*.
2.  **Project Glossary (`glossary.json`):**
    * **Behavior:** Check file directory for `glossary.json`. Load if exists; create template if not.
    * **Purpose:** Enforce consistency (e.g., "Nakama" -> "Companheiro") across a season.
3.  **Injection:** Inject merged glossary terms into the System Prompt template `{{glossary}}`.

## D. Visual Rules ("The Btop Aesthetic") - CRITICAL

1.  **No Mouse Dependency:** All interactions must be keyboard-driven (`[ k ]` key hints).
2.  **Layouts:** Use `lipgloss.JoinHorizontal` and `lipgloss.JoinVertical`. Handle `tea.WindowSizeMsg` to resize panes dynamically.
3.  **Borders:** Use `lipgloss.RoundedBorder()`.
    * **Main Window:** Double Border (`lipgloss.DoubleBorder()`).
    * **Panels:** Rounded or Normal Border.
4.  **Palette:**
    * **Accents:** Neon Pink (`#F700FF`), Cyan (`#00FFFF`), Yellow (`#FFFF00`).
    * **Backgrounds:** Transparent/Terminal Default. Avoid setting background colors on containers to maintain the "CLI feel".
5.  **Components:**
    * **Progress:** `bubbles/progress` styled with block characters (`█`, `░`).
    * **Tables:** `bubbles/table` styled densely (minimal padding).

## E. Modular AI Architecture (Provider Pattern)

1.  **Interface:** Define `type LLMProvider interface` in `internal/core/ai`.
    * `SendBatch(...)`, `ValidateKey()`, `ListModels()`.
2.  **Adapters:**
    * **OpenRouter:** Standard REST client.
    * **Gemini:** Native implementation.
    * **OpenAI:** Native implementation.
    * **Local LLM:** Client for Ollama/LMStudio endpoints.
3.  **Factory:** `ProviderFactory` instantiates the correct adapter based on `config.json`.

## F. Lifecycle & Resilience

1.  **Smart Resume:**
    * **State File:** Write a `.bakasub.temp` JSON file after every successful batch.
    * **Startup:** Check for this file. If found, trigger the **Resume Session Modal** (See Screens).
2.  **Watch Mode:**
    * Implement a specialized Goroutine using `fsnotify`.
    * Monitor input directory for new `.mkv` files.
    * Wait for file lock release -> trigger **Touchless Workflow**.
3.  **Update Checker:**
    * Async Goroutine on startup to query GitHub Releases.
    * Signal `MsgUpdateAvailable` to the Dashboard model.

## G. Workflows (Screen Logic)

### 1. Onboarding (Wizard)
* **Step 1:** Select Provider -> Input Key (or URL).
* **Step 2:** Check Dependencies (Auto-download static binaries if missing).
* **Step 3:** Defaults (Language List + Custom ISO input).

### 2. Job Setup (Full Mode)
* **Conflict Logic:** If multiple subtitle tracks match the target language (e.g., two 'eng' tracks), **DISABLE** the Start action.
    * Show `[ START DISABLED - RESOLVE CONFLICTS ]`.
    * Force user to open Resolution Modal.
* **Config Inheritance:** Pulls default Target Lang from config, but allows override (Simple cycler or input).

### 3. Configuration Menu
* **Providers:** Dynamic fields based on selection (Key vs URL).
* **General:**
    * Target Language: Radio List for top 5 + "Other" option with Text Input for ISO code.
    * Touchless Mode: Toggle + "Configure Rules" button.
* **Prompts (Profile System):**
    * **Factory Profiles:** Immutable (Locked). Can only be cloned.
    * **User Profiles:** Editable. Can be saved/deleted.

### 4. Toolbox
* **Manual Review:** A TUI-based text editor (using `bubbles/textarea`) to edit translation lines before muxing.
* **Header Editor:** In-place modification of MKV flags using `mkvpropedit`.

## UI Reference
Refer strictly to **screens.md** for the ASCII geometry and layout structure.

See: [UI Screens Reference](./screens.md)