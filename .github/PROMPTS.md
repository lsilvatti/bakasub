# PROMPTS.md - BakaSub Execution Plan

## Phase 1: Foundation, UI Shell & Configuration

### Prompt 1.1: Project Skeleton & Stack Initialization
**Context:** We are starting the development of "BakaSub", a high-fidelity TUI for subtitle engineering using Go.
**Reference:** `copilot-instructions.md` (Core Stack, Coding Standards).

**Task:**
Initialize the Go project and set up the directory structure following the Standard Go Project Layout.
1.  Initialize module `github.com/lsilvatti/bakasub`.
2.  Create the following folder structure:
    * `cmd/bakasub/`: Entry point (`main.go`).
    * `internal/ui/`: Bubble Tea models and components.
    * `internal/core/`: Business logic (AI, Media, DB).
    * `internal/config/`: Configuration management (Viper).
    * `pkg/utils/`: Helper functions.
3.  Install core dependencies:
    * `go get github.com/charmbracelet/bubbletea`
    * `go get github.com/charmbracelet/lipgloss`
    * `go get github.com/charmbracelet/bubbles`
    * `go get github.com/spf13/viper`
    * `go get modernc.org/sqlite`
4.  Create a `internal/ui/styles/styles.go` file to define the **"Native Neon"** design system:
    * **Colors:** Define `NeonPink` (#F700FF), `Cyan` (#00FFFF), `Yellow` (#FFFF00).
    * **Borders:** Create a standard `lipgloss.RoundedBorder()` style.
    * **Backgrounds:** Ensure all panels have transparent or terminal-default backgrounds (no solid blocks).
5.  Create a basic `main.go` that initializes a simple Bubble Tea program printing "BakaSub Core Loaded" with the defined styles to verify the setup.

### Prompt 1.2: The "Btop" Dashboard Layout
**Context:** The visual identity is critical. It must mimic the density and layout of `btop`.
**Reference:** `screens.md` (Section 2: Main Dashboard).

**Task:**
Implement the Main Dashboard UI in `internal/ui/dashboard`.
1.  Use `lipgloss.JoinHorizontal` and `lipgloss.JoinVertical` to strictly replicate the grid layout from `screens.md`.
2.  **Constraint:** Do NOT use standard button widgets. Create interactive text elements that display hotkeys like `[ k ] KEY`.
3.  Implement the 4 main dashboard sections as `lipgloss` panels:
    * **Input & Mode:** Shows current path and mode toggle.
    * **Modules:** List of actions (Extract, Translate, Mux, Review).
    * **Toolbox:** List of tools (Header Editor, Glossary, etc.).
    * **System & AI:** Placeholder for Model info, Temp, and Config status.
4.  Implement the Header (ASCII Art Logo placeholder) and Footer (Dependencies status).
5.  Ensure the layout handles `tea.WindowSizeMsg` to resize panels dynamically.

### Prompt 1.3: Configuration Logic & UI
**Context:** We need to persist user settings and manage the prompt profiles.
**Reference:** `copilot-instructions.md` (Config), `screens.md` (Section 6).

**Task:**
Implement the Configuration logic (`internal/config`) and the Settings UI (`internal/ui/settings`).
1.  **Backend:** Create `config.go` using `viper`.
    * Define struct `Config` with fields: `APIProvider`, `APIKey`, `TargetLang`, `TouchlessRules` (struct), `PromptProfiles` (map).
    * Implement `Load()` and `Save()` methods targeting `config.json`.
2.  **UI - Tabbed Layout:** Create a view with tabs: General, Providers, Models, Prompts, Advanced.
3.  **UI - General Tab:**
    * Implement the "Target Language" selector as a Radio List for top 5 languages + a Text Input for "Custom ISO" if "Other" is selected.
4.  **UI - Prompts Tab:**
    * Implement the **Profile Manager**: A list showing "Factory" (Locked 🔒) and "User" (Editable 👤) profiles.
    * Factory profiles cannot be edited, only cloned. User profiles can be edited/deleted.

### Prompt 1.4: The Onboarding Wizard
**Context:** First-run experience when `config.json` is missing.
**Reference:** `screens.md` (Section 1: Setup Wizard).

**Task:**
Implement the Onboarding Wizard in `internal/ui/wizard`.
1.  **Step 1 (Access):** Allow user to select Provider (OpenRouter, Gemini, etc.) and input the API Key (or URL for local). Mask the key input.
2.  **Step 2 (Dependencies):** Check for `ffmpeg` and `mkvmerge` in PATH.
    * *Mock Logic:* If missing, simulate a download with a progress bar (we will implement the real downloader later).
3.  **Step 3 (Defaults):** Select Default Model and Target Language using the same components from the Settings screen.
4.  On completion, save `config.json` and transition to the Dashboard.

## Phase 2: Core Logic, Backend & Orchestration

### Prompt 2.1: MKVToolNix Wrappers (The Engine)
**Context:** The application needs to analyze and manipulate MKV files.
**Reference:** `copilot-instructions.md` (External Tools).

**Task:**
Create `internal/core/media` package.
1.  Implement `Analyze(path string)`: Wraps `mkvmerge -J {file}` to parse JSON metadata.
    * Return a struct containing Tracks (ID, Type, Lang, Codec) and Attachments.
2.  Implement `ExtractTrack(trackID, outputPath)`: Wraps `mkvextract`.
3.  Implement `Mux(sources, output)`: Wraps `mkvmerge`.
4.  **Conflict Detection:** Create a helper function that takes a parsed file and a Target Lang, and returns `true` if multiple subtitle tracks match that language.

### Prompt 2.2: AI Modular Architecture
**Context:** Supporting multiple AI providers with a unified interface.
**Reference:** `copilot-instructions.md` (Section E).

**Task:**
Create `internal/core/ai` package.
1.  Define the Interface:
    ```go
    type LLMProvider interface {
        SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error)
        ValidateKey() bool
        ListModels() ([]string, error)
    }
    ```
2.  Implement Adapters:
    * **OpenRouter:** Standard HTTP client with `Authorization: Bearer` header.
    * **Gemini:** Use `google.golang.org/genai` (native SDK).
    * **LocalLLM:** HTTP client pointing to user-defined URL (Ollama format).
3.  Implement a `ProviderFactory` that reads `config.json` and instantiates the correct adapter.

### Prompt 2.3: DBLocal (Semantic Cache)
**Context:** Local database to save costs by caching translations.
**Reference:** `copilot-instructions.md` (Section B), `Bakasub.pdf`.

**Task:**
Implement `internal/core/db`.
1.  Initialize SQLite connection to `bakasub.db` using `modernc.org/sqlite`.
2.  Create table `cache` with columns: `original_hash` (index), `original_text`, `translated_text`, `lang_pair`.
3.  Implement `GetFuzzyMatch(text string, threshold float64)`:
    * Calculate Levenshtein distance (use a library like `github.com/agnivade/levenshtein`).
    * Return cached translation if similarity > 95%.
4.  Implement `SaveTranslation(original, translated, lang)`.
5. DB Must be thread-safe for concurrent reads/writes.
6. Must persist trough jobs, the idea is to save money by reusing previous translations.

### Prompt 2.4: The Job Setup Flow (Pre-Flight)
**Context:** The critical preparation screen before execution.
**Reference:** `screens.md` (Section 3: Workflow).

**Task:**
Implement the `JobSetup` model in `internal/ui/job`.
1.  **Directory Analysis:** When a folder is selected from Dashboard, scan for MKV files.
2.  **Conflict UI:** Use the logic from Prompt 2.1. If a conflict is detected:
    * Change the "Start" button style/text to `[ START DISABLED ]`.
    * Implement the **Resolution Modal** overlay to let user pick the correct track.
3.  **Cost Estimator:** Calculate (Total Characters in subtitle tracks) * (Model Price per 1M) and update the UI label dynamically.
4.  **Dry Run:** Triggered by `[ d ]`. Check write permissions, calculate tokens, and display the "Simulation Report" TUI.
5. **Glossary**: Save common terms and their translations to be used in the AI prompt. Save names and locations in `glossary.json` in the video folder to be used on the translation to keep context.

### Prompt 2.5: Execution Engine & Automation
**Context:** Running the actual translation job and handling "Watch Mode".
**Reference:** `screens.md` (Section 4), `copilot-instructions.md` (Section F).

**Task:**
Implement the execution loop.
1.  **The Worker:** Create a Bubble Tea `Cmd` that runs the pipeline:
    * Extract -> Parse (ASS/SRT) -> Batching -> DBLocal Check -> AI Request -> Reassemble -> Mux.
2.  **UI Updates:** Stream logs to a `viewport` and update the Block Progress Bar (`[████░░]`).
3.  **Smart Resume:**
    * Write a `.bakasub.temp` JSON file after every successful batch.
    * On startup check, if `.temp` exists, trigger the **Resume Modal**.
4.  **Watch Mode:** Implement `internal/core/watcher` using `fsnotify`. Monitor input dir, debounce file events, and trigger the pipeline using "Touchless Rules" from config.
5.  **Sliding Window Context**: In the pipeline, ensure that the last 3 lines of Batch N are appended to the system prompt of Batch N+1 as passive context (read-only) to maintain dialogue fluidity, as specified in copilot-instructions.md

## Phase 3: Toolbox, Resilience & Distribution

### Prompt 3.1: Manual Review Editor & Glossary
**Context:** Power tools for manual corrections.
**Reference:** `screens.md` (Section 5: Toolbox).

**Task:**
1.  **Manual Review:** Implement `internal/ui/review`.
    * Create a split-view: Original (Left, Read-only) vs Translated (Right, Editable `textarea`).
    * Implement navigation (`Tab`/Arrows) to switch lines.
    * Implement `Ctrl+S` to save the file to disk.
2.  **Project Glossary:** Implement `internal/ui/glossary`.
    * A table view allowing Add/Edit/Delete of Regex/Translation pairs.
    * Save logic: Serialize to `glossary.json` in the video folder.
    * **Integration:** Update `core/ai` to inject this file's content into the System Prompt `{{glossary}}` variable.

### Prompt 3.2: Header Editor & Quality Gate
**Context:** Metadata editing and error checking.
**Reference:** `screens.md` (Section 5 & 4).

**Task:**
1.  **Header Editor:**
    * UI: A table listing tracks from `mkvmerge -J`.
    * Logic: Use `mkvpropedit` to toggle "Default" and "Forced" flags instantly.
2.  **Quality Gate (Linter):** Implement `internal/core/linter`.
    * Run regex checks after translation (e.g., mismatched brackets `{}`, `[]`).
    * If errors found, show the "Issues Found" screen before muxing.
    * Allow "Auto-Fix" (Regex replace) or "Manual Review" (jump to Editor).

### Prompt 3.3: Toolbox Expansion (Attachments & Remuxer)
**Context:** Remaining MKV manipulation tools defined in the design.
**Reference:** `screens.md` (Section 5: Toolbox).

**Task:**
1.  **Attachment Manager:** Implement `internal/ui/attachments`.
    * UI: List existing attachments (fonts, images) with Size/MIME.
    * Actions: 
        * `[ d ]`: Toggle delete mode (mark for removal).
        * `[ a ]`: Open Path Input Modal to add a new file.
        * `[ e ]`: Extract all attachments to a folder.
    * Backend: Use `mkvpropedit --add-attachment` and `--delete-attachment`.
2.  **Quick Remuxer:** Implement `internal/ui/remuxer`.
    * UI: A Checkbox List showing all video/audio/subtitle tracks.
    * Logic: Allow user to uncheck tracks they want to remove (e.g., dubs).
    * Action: Generate a `mkvmerge -o new_file.mkv` command that only includes the selected track IDs. This is a destructive/creation action, so show a confirmation or progress spinner.

### Prompt 3.4: Panic Handling & Update Checker
**Context:** Production-readiness features.
**Reference:** `copilot-instructions.md` (Coding Standards).

**Task:**
1.  **BSOD:** Implement a global `recover()` middleware in `main.go`.
    * If a panic occurs, clear screen and render the "Critical System Error" TUI with the stack trace.
    * Provide a link to GitHub Issues.
2.  **Update Checker:**
    * Create an async goroutine that queries `https://api.github.com/repos/lsilvatti/bakasub/releases/latest`.
    * If version > current, send a `MsgUpdateAvailable` to the Dashboard model to show the `[!]` indicator.

### Prompt 3.5: Build & Distribution Strategy
**Context:** "Binary First" strategy.
**Reference:** `Bakasub.pdf` (Distribution).

**Task:**
Prepare the project for release.
1.  Create a `Makefile` with targets:
    * `build-win`: `GOOS=windows go build -o bin/bakasub.exe -ldflags="-s -w" ./cmd/bakasub`
    * `build-linux`: `GOOS=linux go build -o bin/bakasub -ldflags="-s -w" ./cmd/bakasub`
2.  Create a `.github/workflows/release.yml` Action:
    * Trigger on tag creation.
    * Build for Windows/Linux.
    * Upload artifacts to GitHub Releases.
3.  **Documentation:** Generate templates for `README.md` (English), `README-pt.md` (Portuguese), and `README-es.md` (Spanish) detailing installation and key bindingsm, use a friendly tone with some jokes about animes, series and tv shows, nothing much specific, more generic/steriotype jokes like “It’s not like I made this subtitle tool for you or anything... B-Baka!” 