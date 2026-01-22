# Screens Reference (Go/Btop Aesthetic)

## 1. Setup Wizard (Onboarding)

### Step 1: Provider & Access
╔══════════════════════════════════════════════════════════════════════════════╗
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
║ ▒   BAKASUB SETUP WIZARD                                       [STEP 1/3]  ▒ ║
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ┌── 1. AI PROVIDER (Select One) ───────────────────────────────────────┐   ║
║   │   (o) ♾️  OpenRouter (Recommended - Access to all models)            │   ║
║   │   ( ) 💎 Google Gemini API (Native)                                  │   ║
║   │   ( ) 🤖 OpenAI API (Native)                                         │   ║
║   │   ( ) 🏠 Local LLM (Ollama/LMStudio - Requires GPU)                  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 2. CREDENTIALS ────────────────────────────────────────────────────┐   ║
║   │   API KEY > ______________________________________________________   │   ║
║   │             [ VALIDATING KEY... OK ]                                 │   ║
║   │   *Enter the full key provided by the selected service.              │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 3. INTERFACE LANGUAGE ─────────────────────────────────────────────┐   ║
║   │   ( ) ENGLISH (Default)                                              │   ║
║   │   (o) PORTUGUÊS (Brasil)                                             │   ║
║   │   ( ) ESPAÑOL                                                        │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [Q] QUIT                                                   [NEXT STEP >]   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Step 2: System Dependencies
╔══════════════════════════════════════════════════════════════════════════════╗
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
║ ▒   BAKASUB SETUP WIZARD                                       [STEP 2/3]  ▒ ║
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   SYSTEM CHECK                                                               ║
║   Ensuring external tools (FFmpeg/MKVToolNix) are available in PATH.         ║
║                                                                              ║
║   ┌── FFmpeg (Media Processor) ──────────────────────────────────────────┐   ║
║   │   STATUS: [MISSING] -> [DOWNLOADING BINARY...]                       │   ║
║   │   PROGRESS:                                                          │   ║
║   │   [████████████████████████████░░░░░░░░░░░░] 72%                     │   ║
║   │   35MB / 48MB (4.2 MB/s)                                             │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── MKVToolNix (Container Engine) ─────────────────────────────────────┐   ║
║   │   STATUS: [FOUND]                                                    │   ║
║   │   PATH:   /usr/bin/mkvmerge                                          │   ║
║   │   VER:    v82.0 (Strapping Snorlax)                                  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [< BACK]                                    [WAITING FOR DOWNLOADS...]     ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Step 3: Core Defaults
╔══════════════════════════════════════════════════════════════════════════════╗
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
║ ▒   BAKASUB SETUP WIZARD                                       [STEP 3/3]  ▒ ║
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   1. TRANSLATION TARGET (Output Language)                                    ║
║      (o) PT-BR (Português)   ( ) EN-US (English)   ( ) ES (Español)          ║
║      ( ) JA-JP (Japanese)    ( ) FR-FR (Français)  ( ) DE (Deutsch)          ║
║      ( ) OTHER ISO CODE: [_______]  (e.g. 'it', 'ru', 'zh-cn')               ║
║                                                                              ║
║   2. PREFERENCES                                                             ║
║      [X] REMOVE HEARING IMPAIRED TAGS                                        ║
║      GLOBAL TEMPERATURE: [ 0.3   ]                                           ║
║                                                                              ║
║   ┌── [?] HELPER: WHAT IS TEMPERATURE? ──────────────────────────────────┐   ║
║   │  • Definition: Controls the "creativity" and randomness of the AI.   │   ║
║   │  • Low (0.0 - 0.2): Precise, literal. Best for Docs/Technical.       │   ║
║   │  • Mid (0.3 - 0.5): Balanced. Good for General Series.               │   ║
║   │  • High (0.6 - 1.0): Creative, diverse. Best for Anime/Drama/Slang.  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ℹ️  NOTE: AI Model selection available in Configuration > AI Models        ║
║            (FREE and ALL MODELS tabs with search)                            ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [< BACK]                                              [FINISH SETUP >]     ║
╚══════════════════════════════════════════════════════════════════════════════╝

## 2. Main Dashboard

╔══════════════════════════════════════════════════════════════════════════════╗
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
║ ▒  ____        _          ____        _      v1.0.0                        ▒ ║
║ ▒ | __ )  __ _| | ____ _ / ___| _   _| |__   [APP RUNNING]                 ▒ ║
║ ▒ |  _ \ / _` | |/ / _` |\___ \| | | | '_ \                                ▒ ║
║ ▒ | |_) | (_| |   < (_| | ___) | |_| | |_) | [API: ONLINE ] [CACHE: OK]    ▒ ║
║ ▒ |____/ \__,_|_|\_\__,_||____/ \__,_|_.__/  [!] UPDATE AVAILABLE (v1.1)   ▒ ║
║ ▒                                                                          ▒ ║
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ┌── 1. INPUT & MODE ───────────────────────────────────────────────────┐   ║
║   │   PATH > /srv/media/anime/one_piece/______________________________   │   ║
║   │                                                                      │   ║
║   │   (o) FULL PROCESS (Extract -> Translate -> Mux)                     │   ║
║   │       *Opens "Job Setup" screen for uninterrupted processing.        │   ║
║   │   ( ) WATCH MODE (Auto-process new files in folder)                  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 2. MODULES (Standalone) ─────────┐┌── 3. TOOLBOX (MKVToolNix) ─────┐   ║
║   │   [ 1 ] EXTRACT TRACKS             ││   [ 5 ] EDIT FLAGS / METADATA  │   ║
║   │   [ 2 ] TRANSLATE SUBTITLE         ││   [ 6 ] MANAGE ATTACHMENTS     │   ║
║   │   [ 3 ] MUX CONTAINER              ││   [ 7 ] ADD/REMOVE TRACKS      │   ║
║   │   [ 4 ] MANUAL REVIEW (Editor)     ││   [ 8 ] PROJECT GLOSSARY       │   ║
║   └────────────────────────────────────┘└────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 4. SYSTEM & AI ────────────────────────────────────────────────────┐   ║
║   │   MODEL: Gemini Flash 1.5 (Google)        [ m ] CHANGE MODEL         │   ║
║   │   TARGET: PT-BR  │  TEMP: 0.3             [ c ] CONFIGURATION        │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [k] KO-FI  │  DEPS: FFmpeg [OK] MKVToolNix [OK]  │  [q] QUIT                 ║
╚══════════════════════════════════════════════════════════════════════════════╝

## 3. Workflow: Full Process (Job Setup)

### Step A: Directory Detection
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║   ... (Dashboard Background dimmed/blurred using block chars) ...            ║
║   ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒   ║
║           ╔══════════════════════════════════════════════════════╗           ║
║           ║  DIRECTORY DETECTED                                  ║           ║
║           ╠══════════════════════════════════════════════════════╣           ║
║           ║                                                      ║           ║
║           ║  PATH: /srv/media/anime/one_piece/                   ║           ║
║           ║                                                      ║           ║
║           ║  ANALYSIS:                                           ║           ║
║           ║  • MKV Files Found: 24                               ║           ║
║           ║  • Subtitles Found: 00                               ║           ║
║           ║                                                      ║           ║
║           ║  HOW DO YOU WANT TO PROCEED?                         ║           ║
║           ║                                                      ║           ║
║           ║  [ b ] PROCESS BATCH (All 24 files)                  ║           ║
║           ║        *Apply same settings to all episodes.         ║           ║
║           ║                                                      ║           ║
║           ║  [ s ] SELECT SINGLE FILE                            ║           ║
║           ║        *Open file picker to choose one episode.      ║           ║
║           ║                                                      ║           ║
║           ║  [ ESC ] CANCEL                                      ║           ║
║           ║                                                      ║           ║
║           ╚══════════════════════════════════════════════════════╝           ║
║   ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Step B: Configuration & Pre-Flight (Standard State)
╔══════════════════════════════════════════════════════════════════════════════╗
║   JOB SETUP: /srv/media/anime/one_piece/ (Batch: 24 Files)                   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ┌── 1. EXTRACTION STRATEGY ────────────────────────────────────────────┐   ║
║   │   SUBTITLE SOURCE: [ Track 2 (eng) - ASS ]  [ AUTO-DETECT ]          │   ║
║   │   AUDIO REFERENCE: [ Track 1 (jpn) - FLAC]  (For context)            │   ║
║   │   [X] Extract Fonts (Attachments)                                    │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 2. TRANSLATION CONTEXT ────────────────────────────────────────────┐   ║
║   │   MEDIA TYPE:      [ < ANIME > ]  (Applies specialized prompt)       │   ║
║   │   TARGET LANG:     [ PT-BR ]                                         │   ║
║   │   GLOSSARY:        [ Auto-Inject (Series Name) ]                     │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 3. MUXING OUTPUT ──────────────────────────────────────────────────┐   ║
║   │   MODE:            [ REPLACE ORIGINAL (With .old backup) ]           │   ║
║   │   TRACK TITLE:     [ Português (BakaSub AI) ]                        │   ║
║   │   FLAGS:           [X] Set Default   [ ] Set Forced                  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── [i] COST ESTIMATION ───────────────────────────────────────────────┐   ║
║   │   Model: Gemini Flash 1.5  │  Est. Tokens: 1.2M  │  Est. Cost: $0.08 │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ESC] BACK      [ d ] SIMULATION (DRY RUN)     [ ENTER ] START JOB         ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Step B-2: Conflict State (Blocker)
╔══════════════════════════════════════════════════════════════════════════════╗
║   JOB SETUP: /srv/media/anime/one_piece/ (Batch: 24 Files)                   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ┌── 1. EXTRACTION STRATEGY ────────────────────────────────────────────┐   ║
║   │   SUBTITLE SOURCE: [!] MULTIPLE 'ENG' TRACKS FOUND  [ r ] RESOLVE    │   ║
║   │                    *Auto-detect failed. Please select source.        │   ║
║   │                                                                      │   ║
║   │   AUDIO REFERENCE: [ Track 1 (jpn) - FLAC]  (For context)            │   ║
║   │   [X] Extract Fonts (Attachments)                                    │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ... (Other panels remain same) ...                                         ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ESC] BACK                        [ START DISABLED - RESOLVE CONFLICTS ]   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Step C: Resolution Modal
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║   ... (Job Setup Background dimmed) ...                                      ║
║                                                                              ║
║           ╔══════════════════════════════════════════════════════╗           ║
║           ║  RESOLVE TRACK CONFLICT                              ║           ║
║           ╠══════════════════════════════════════════════════════╣           ║
║           ║                                                      ║           ║
║           ║  Multiple tracks detected with lang 'eng'.           ║           ║
║           ║  Select the Full Dialogue track:                     ║           ║
║           ║                                                      ║           ║
║           ║  ┌───┬──────┬──────┬─────────────────────────────┐   ║           ║
║           ║  │ # │ TYPE │ SIZE │ TRACK NAME (METADATA)       │   ║           ║
║           ║  ├───┼──────┼──────┼─────────────────────────────┤   ║           ║
║           ║  │ 2 │ ASS  │ 45KB │ Crunchyroll Subs            │   ║           ║
║           ║  │ 3 │ ASS  │ 5KB  │ Songs & Signs Only          │   ║           ║
║           ║  └───┴──────┴──────┴─────────────────────────────┘   ║           ║
║           ║                                                      ║           ║
║           ║  (o) Track 2  (Recommended - Larger Size)            ║           ║
║           ║  ( ) Track 3                                         ║           ║
║           ║                                                      ║           ║
║           ╚══════════════════════════════════════════════════════╝           ║
║               [ ESC ] CANCEL             [ ENTER ] CONFIRM                   ║
╚══════════════════════════════════════════════════════════════════════════════╝

## 4. Execution & Resilience

### Simulation Report (Dry Run)
*Triggered by `[ d ]` in Job Setup.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   DRY RUN REPORT (SIMULATION)                                                ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   ┌── JOB SUMMARY ───────────────────────────────────────────────────────┐   ║
║   │   FILES:       24 Episodes (.mkv)                                    │   ║
║   │   TOTAL SIZE:  32.4 GB                                               │   ║
║   │   PROVIDER:    Google Gemini 1.5 Flash                               │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── COST ANALYSIS (ESTIMATED) ─────────────────────────────────────────┐   ║
║   │   TOTAL CHARACTERS:    4,500,000 (approx)                            │   ║
║   │   INPUT TOKENS:        1.2M                                          │   ║
║   │   OUTPUT TOKENS:       900k                                          │   ║
║   │   ────────────────────────────────────────────────────────────────   │   ║
║   │   PRICE PER 1M:        $0.07                                         │   ║
║   │   ESTIMATED TOTAL:     $0.15 USD                                     │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── PRE-FLIGHT CHECKS ─────────────────────────────────────────────────┐   ║
║   │   [OK] Write Permissions (Output Folder)                             │   ║
║   │   [OK] FFmpeg/MKVMerge binaries found                                │   ║
║   │   [!!] WARNING: Episode 12 has broken timestamps (may cause desync)  │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] BACK TO SETUP               [ ENTER ] PROCEED TO PAYMENT/RUN       ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Job Running (The Logs)
*The main execution view. No interaction needed usually.*
╔══════════════════════════════════════════════════════════════════════════════╗
║  JOB RUNNING: One Piece E1080 (File 1/24)                   [ETA: 00:04:12]  ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > [10:04:22] Preflight check passed. Found 450 dialogue lines.              ║
║  > [10:04:23] Context: Anime Mode active. Temp: 0.7.                         ║
║  > [10:04:25] Batch 1 (Lines 1-50) sent to OpenRouter (Gemini Flash).        ║
║  > [10:04:28] Batch 1 received. Sanity Check: PASSED.                        ║
║  > [10:04:29] Batch 2 (Lines 51-100) sent...                                 ║
║  > [10:04:35] [WARN] Batch 2 Desync (Expected 50 IDs, got 48).               ║
║  > [10:04:35] └─ Engaging Anti-Desync Protocol (Split Strategy).             ║
║  > [10:04:36] └─ Split 2a (Lines 51-75) sent... OK.                          ║
║  > [10:04:41] └─ Split 2b (Lines 76-100) sent... OK.                         ║
║  > [10:04:42] Batch 3 (Lines 101-150) sent...                                ║
║  > [10:04:45] Batch 3 received.                                              ║
║  > [10:04:46] Processing Styles: Preserving 'MainDialogue' events.           ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  FILE PROGRESS:  [██████████████████████░░░░░░░░░░]  65%                     ║
║  BATCH TOTAL:    [█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]  1/24 Files              ║
║                                                                              ║
║  PERFORMANCE:    450 tok/s  │  ERRORS: 0 (1 Handled)                         ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Smart Resume Modal (Crash Recovery)
*Appears on startup if a .temp file is found.*
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║           ╔══════════════════════════════════════════════════════╗           ║
║           ║  UNFINISHED SESSION DETECTED                         ║           ║
║           ╠══════════════════════════════════════════════════════╣           ║
║           ║                                                      ║           ║
║           ║  BakaSub exited abnormally during the last run.      ║           ║
║           ║                                                      ║           ║
║           ║  FILE:   One Piece E1080.mkv                         ║           ║
║           ║  BATCH:  14 / 28 (50% Complete)                      ║           ║
║           ║  CACHE:  Translation saved to .bakasub.temp          ║           ║
║           ║                                                      ║           ║
║           ║  Do you want to resume from Batch 14?                ║           ║
║           ║                                                      ║           ║
║           ║  [ d ] DISCARD & RESTART       [ ENTER ] RESUME      ║           ║
║           ║                                                      ║           ║
║           ╚══════════════════════════════════════════════════════╝           ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Success Screen
╔══════════════════════════════════════════════════════════════════════════════╗
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
║ ▒   PROCESS COMPLETED SUCCESSFULLY                                         ▒ ║
║ ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒ ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   SUMMARY:                                                                   ║
║   ────────                                                                   ║
║   Files Processed:   24 / 24                                                 ║
║   Total Lines:       8,450                                                   ║
║   Time Elapsed:      00:45:12                                                ║
║                                                                              ║
║   AI STATISTICS:                                                             ║
║   ──────────────                                                             ║
║   Model Used:        Gemini Flash 1.5                                        ║
║   Input Tokens:      1.2M                                                    ║
║   Output Tokens:     850k                                                    ║
║   Total Cost:        $0.14 USD                                               ║
║                                                                              ║
║   OUTPUT:                                                                    ║
║   ───────                                                                    ║
║   Location:          /srv/media/anime/one_piece/                             ║
║   Status:            Muxed (Original Replaced, .old created)                 ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ o ] OPEN FOLDER           [ l ] VIEW LOGS           [ ESC ] DASHBOARD    ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Quality Gate (Errors Found)
*Triggered if Linter finds issues before Muxing.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   QUALITY GATE: ISSUES FOUND                                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   The automatic linter found potential issues in the translation.            ║
║   Please review before muxing.                                               ║
║                                                                              ║
║   ┌── DETECTED ISSUES ───────────────────────────────────────────────────┐   ║
║   │   ID    SEVERITY    ISSUE TYPE             CONTENT                   │   ║
║   │  ──────────────────────────────────────────────────────────────────  │   ║
║   │   12    [HIGH]      Broken Tags            {\an8}Olá mundo{\an}      │   ║
║   │   45    [MED]       English Residual       Hello, tudo bem?          │   ║
║   │   88    [LOW]       Glossary Mismatch      Traduzido "Nakama" como   │   ║
║   │                                            "Amigo" (Esperado: Compa) │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── ACTION (Select with Arrows) ───────────────────────────────────────┐   ║
║   │   (o) AUTO-FIX (Attempt to fix tags/terms via Regex)                 │   ║
║   │   ( ) MANUAL REVIEW (Open Editor)                                    │   ║
║   │   ( ) IGNORE AND CONTINUE                                            │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                       [ ENTER ] EXECUTE      ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Error Screen (BSOD)
*The global panic handler.*
╔══════════════════════════════════╗
║      CRITICAL SYSTEM ERROR       ║
╠══════════════════════════════════╣
║                                  ║
║  ConnectionResetError            ║
║  [Errno 104] Connection reset    ║
║                                  ║
║  Module: core/api.py             ║
║  Line: 404                       ║
║                                  ║
║  The application has crashed.    ║
║  A log has been saved.           ║
║                                  ║
║      [ g ] OPEN GITHUB ISSUE     ║
║          [ q ] EXIT APP          ║
║                                  ║
╚══════════════════════════════════╝

## 5. Toolbox (Standalone Tools)

### Project Glossary Editor
*Acessado via [ 8 ] no Dashboard.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   PROJECT GLOSSARY EDITOR                                                    ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   FILE: /srv/media/anime/one_piece/glossary.json                             ║
║                                                                              ║
║   ┌── TERMS MANAGEMENT ──────────────────────────────────────────────────┐   ║
║   │   [ a ] ADD NEW TERM    [ i ] IMPORT CSV    [ DEL ] REMOVE SELECTED  │   ║
║   │                                                                      │   ║
║   │   ORIGINAL TERM (Regex)        TRANSLATION              TYPE         │   ║
║   │  ──────────────────────────────────────────────────────────────────  │   ║
║   │   Nakama                       Companheiro              [Noun]       │   ║
║   │   Gomu Gomu no                 Gomu Gomu no             [Attack]     │   ║
║   │   Mugiwara                     Chapéu de Palha          [Name]       │   ║
║   │   Kaizoku-ou                   Rei dos Piratas          [Title]      │   ║
║   │                                                                      │   ║
║   │  [< PREV]   Page 1/1   [NEXT >]                                      │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── [?] HELPER: GLOSSARY ──────────────────────────────────────────────┐   ║
║   │  • These terms are injected into the System Prompt.                  │   ║
║   │  • Use this for persistent consistency across episodes/seasons.      │   ║
║   │  • "Regex" column supports simple patterns (e.g., ^Start).           │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CLOSE & SAVE                                                       ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Manual Review Editor (TUI)
*Acessado via [ 4 ] no Dashboard.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   MANUAL REVIEW MODE (TUI EDITOR)                                            ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   PROGRESS: Line 45 / 450  (10%)                                             ║
║                                                                              ║
║   ┌── CURRENT LINE (ID: 45) ─────────────────────────────────────────────┐   ║
║   │   TIMING: 00:04:12,500 --> 00:04:15,000                              │   ║
║   │                                                                      │   ║
║   │   [ORIGINAL]                                                         │   ║
║   │   Don't underestimate the sea, kid!                                  │   ║
║   │                                                                      │   ║
║   │   [TRANSLATION - EDITABLE]                                           │   ║
║   │   Não subestime o mar, garoto!_                                      │   ║
║   │                                                                      │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── CONTEXT (PREVIOUS LINES) ──────────────────────────────────────────┐   ║
║   │   43: Eu avisei para não ir lá.                                      │   ║
║   │   44: É muito perigoso nesta época do ano.                           │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── CONTROLS ──────────────────────────────────────────────────────────┐   ║
║   │   [ENTER] Commit & Next    [BACKSPACE] Undo    [CTRL+S] Save File    │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ LEFT ] PREV LINE        [ g ] GOTO LINE...        [ RIGHT ] NEXT LINE    ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Header Editor (MKVPropEdit)
*Acessado via [ 5 ] no Dashboard.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   TOOLBOX: HEADER EDITOR                                                     ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   FILE: Akira_1988.mkv                                                       ║
║                                                                              ║
║   ┌── GLOBAL METADATA ───────────────────────────────────────────────────┐   ║
║   │   MOVIE TITLE: [ Akira (1988) Remastered 4K ]                        │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── TRACK FLAGS (Select with Arrows, Toggle with Space) ───────────────┐   ║
║   │                                                                      │   ║
║   │   ID   TYPE   LANG   TRACK NAME                DEFAULT   FORCED      │   ║
║   │  ──────────────────────────────────────────────────────────────────  │   ║
║   │   1    Vid    jpn    Video HEVC                [ YES ]   [ NO  ]     │   ║
║   │   2    Aud    jpn    Original Audio            [ YES ]   [ NO  ]     │   ║
║   │   3    Aud    eng    English Dub               [ NO  ]   [ NO  ]     │   ║
║   │   4    Sub    por    BakaSub AI                [ YES ]   [ NO  ]     │   ║
║   │   5    Sub    eng    Signs & Songs             [ NO  ]   [ YES ]     │   ║
║   │                                                                      │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                      [ ENTER ] APPLY CHANGES║
╚══════════════════════════════════════════════════════════════════════════════╝

### Attachment Manager
*Acessado via [ 6 ] no Dashboard.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   TOOLBOX: ATTACHMENT MANAGER                                                ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   FILE: Akira_1988.mkv                                                       ║
║                                                                              ║
║   ┌── EXISTING ATTACHMENTS ──────────────────────────────────────────────┐   ║
║   │   [ d ] TOGGLE DELETE MODE                                           │   ║
║   │                                                                      │   ║
║   │   ID    FILENAME                     MIME-TYPE             SIZE      │   ║
║   │  ──────────────────────────────────────────────────────────────────  │   ║
║   │   1     Cover.jpg                    image/jpeg            1.2 MB    │   ║
║   │   2     OpenSans-Bold.ttf            application/x-font    140 KB    │   ║
║   │   3     Arial-Narrow.ttf             application/x-font    98 KB     │   ║
║   │                                                                      │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── ACTIONS ───────────────────────────────────────────────────────────┐   ║
║   │   [ a ] ADD NEW FILE...  (Opens Path Input)                          │   ║
║   │   [ e ] EXTRACT ALL...   (Dumps fonts to folder)                     │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                      [ ENTER ] EXECUTE      ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Quick Remuxer
*Acessado via [ 7 ] no Dashboard.*
╔══════════════════════════════════════════════════════════════════════════════╗
║   TOOLBOX: QUICK REMUXER                                                     ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║   BASE FILE: Akira_1988.mkv                                                  ║
║                                                                              ║
║   ┌── 1. SELECT TRACKS TO KEEP (Space to Toggle) ────────────────────────┐   ║
║   │   [X] 1: Video (HEVC)                                                │   ║
║   │   [X] 2: Audio (JPN)                                                 │   ║
║   │   [ ] 3: Audio (ENG) - [REMOVED]                                     │   ║
║   │   [ ] 4: Sub   (ENG) - [REMOVED]                                     │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 2. ADD EXTERNAL TRACKS ────────────────────────────────────────────┐   ║
║   │   [ a ] ADD FILE...                                                  │   ║
║   │                                                                      │   ║
║   │   PENDING ADDITION:                                                  │   ║
║   │   • Akira_Br_Dub_5.1.ac3  [ LANG: por ]                              │   ║
║   │   • Akira_Legenda.ass     [ LANG: por ]                              │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                      [ ENTER ] START MUXING ║
╚══════════════════════════════════════════════════════════════════════════════╝

## 6. Configuration Menu

### AI Providers Tab
╔══════════════════════════════════════════════════════════════════════════════╗
║   CONFIGURATION                                                              ║
╟──────────────────────────────────────────────────────────────────────────────╢
║   [ GENERAL ]   [ AI PROVIDERS ]   [ AI MODELS ]   [ PROMPTS ]               ║
║                                                                              ║
║   ┌── ACTIVE PROVIDER ───────────────────────────────────────────────────┐   ║
║   │   SELECT SERVICE:                                                    │   ║
║   │   (o) ♾️  OpenRouter (Recommended)                                   │   ║
║   │   ( ) 💎 Google Gemini API                                           │   ║
║   │   ( ) 🤖 OpenAI API                                                  │   ║
║   │   ( ) 🏠 Local LLM (Ollama/LMStudio)                                 │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── CONFIGURATION (Changes based on selection) ────────────────────────┐   ║
║   │   API KEY: sk-or-v1-********************************** [ SHOW ]      │   ║
║   │   BALANCE: $4.50 USD                                                 │   ║
║   │   STATUS:  [ CONNECTED ]                                             │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                            [ ENTER ] SAVE   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### General Settings
╔══════════════════════════════════════════════════════════════════════════════╗
║   CONFIGURATION                                                              ║
╟──────────────────────────────────────────────────────────────────────────────╢
║   [ GENERAL ]   [ AI PROVIDERS ]   [ AI MODELS ]   [ PROMPTS ]               ║
║                                                                              ║
║   ┌── INTERFACE LANGUAGE (UI) ───────────────────────────────────────────┐   ║
║   │   ( ) English (en)                                                   │   ║
║   │   (o) Português Brasil (pt-br)                                       │   ║
║   │   ( ) Español (es)                                                   │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── AUTOMATION BEHAVIOR ───────────────────────────────────────────────┐   ║
║   │   DEFAULT TARGET LANG:                                               │   ║
║   │   (o) PT-BR   ( ) EN-US   ( ) ES-LA   ( ) JA-JP   ( ) FR-FR          │   ║
║   │   ( ) OTHER:  [_______]                                              │   ║
║   │                                                                      │   ║
║   │   [X] TOUCHLESS MODE (Auto-start jobs)    [ c ] CONFIGURE RULES      │   ║
║   │       *Warning: Skips all confirmation screens.                      │   ║
║   │                                                                      │   ║
║   │   [X] REMOVE HEARING IMPAIRED TAGS (e.g., [Music], (Sigh))           │   ║
║   │   [X] AUTO-CHECK FOR UPDATES (GitHub API)                            │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ SAVE & EXIT ]                                           [ CANCEL ]       ║
╚══════════════════════════════════════════════════════════════════════════════╝
╚══════════════════════════════════════════════════════════════════════════════╝

### Touchless Configuration Modal
*Triggered by `[ c ]` in General Tab.*
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║           ╔══════════════════════════════════════════════════════╗           ║
║           ║  TOUCHLESS RULES (UNMANNED EXECUTION)                ║           ║
║           ╠══════════════════════════════════════════════════════╣           ║
║           ║                                                      ║           ║
║           ║  IF MULTIPLE SUBTITLES FOUND:                        ║           ║
║           ║  (o) AUTO-SELECT LARGEST (Likely full dialogue)      ║           ║
║           ║  ( ) AUTO-SELECT SMALLEST (Likely Signs/Songs)       ║           ║
║           ║  ( ) SKIP FILE (Log error and continue batch)        ║           ║
║           ║                                                      ║           ║
║           ║  DEFAULT CONTEXT PROFILE:                            ║           ║
║           ║  [ Anime (Factory)                    ▼ ]            ║           ║
║           ║                                                      ║           ║
║           ║  MUXING STRATEGY:                                    ║           ║
║           ║  (o) Replace Original (Backup .old)                  ║           ║
║           ║  ( ) Create New File                                 ║           ║
║           ║                                                      ║           ║
║           ╚══════════════════════════════════════════════════════╝           ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝

### AI Models Tab
╔══════════════════════════════════════════════════════════════════════════════╗
║   CONFIGURATION                                                              ║
╟──────────────────────────────────────────────────────────────────────────────╢
║   [ GENERAL ]   [ AI PROVIDERS ]   [ AI MODELS ]   [ PROMPTS ]               ║
║                                                                              ║
║   ┌── MODEL SELECTION ───────────────────────────────────────────────────┐   ║
║   │   SUB-TABS:  < FREE >   < ALL MODELS >   < SEARCH >                  │   ║
║   │   ────────────────────────────────────────────────────────────────   │   ║
║   │                                                                      │   ║
║   │   SEARCH > [Type to filter models...]                                │   ║
║   │                                                                      │   ║
║   │   NAME                          COST(1M)   CTX      TAGS             │   ║
║   │  ──────────────────────────────────────────────────────────────────  │   ║
║   │   ( ) Llama 3.3 70B              FREE       128k     [OpenSource]    │   ║
║   │   ( ) Qwen 2.5 72B               FREE       32k      [Multilingual]  │   ║
║   │   ( ) Mistral 7B Instruct        FREE       32k      [Fast]          │   ║
║   │   ( ) Gemma 2 9B                 FREE       8k       [Google]        │   ║
║   │   ( ) Phi-3 Medium               FREE       128k     [Microsoft]     │   ║
║   │                                                                      │   ║
║   │  [< PREV]   Page 1/5   [NEXT >]                                      │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── 💡 INFO ──────────────────────────────────────────────────────────┐   ║
║   │  • FREE: Zero-cost models (community/open-source)                    │   ║
║   │  • ALL MODELS: Complete catalog with pricing                         │   ║
║   │  • SEARCH: Filter by name, provider, or capability                   │   ║
║   │  • Use ↑↓ to select, [ENTER] to confirm                              │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                            [ ENTER ] SAVE   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Prompts Tab (Factory/Locked State)
╔══════════════════════════════════════════════════════════════════════════════╗
║   CONFIGURATION                                                              ║
╟──────────────────────────────────────────────────────────────────────────────╢
║   [ GENERAL ]   [ AI PROVIDERS ]   [ AI MODELS ]   [ PROMPTS ]               ║
║                                                                              ║
║   ┌── PROFILE MANAGER ───────────────────────────────────────────────────┐   ║
║   │   SELECT PROFILE:                                                    │   ║
║   │   [ 🔒 ANIME (Factory Default)         ▼ ] <--- Dropdown             │   ║
║   │   │  🔒 MOVIE (Factory Default)          │                           │   ║
║   │   │  🔒 SERIES (Factory Default)         │                           │   ║
║   │   │  🔒 DOCUMENTARY (Factory Default)    │                           │   ║
║   │   │  🔒 YOUTUBE (Factory Default)        │                           │   ║
║   │   │  👤 ANIME (Custom - Gírias)          │                           │   ║
║   │   └──────────────────────────────────────┘                           │   ║
║   │   STATUS: [LOCKED] Cannot be modified. Clone to customize.           │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── SYSTEM PROMPT PREVIEW (READ ONLY) ─────────────────────────────────┐   ║
║   │   You are a professional translator specializing in Japanese         │   ║
║   │   animation. Your goal is to translate the JSON payload while        │   ║
║   │   preserving honorifics (e.g., -san, -kun) and attack names.         │   ║
║   │   Use the glossary terms: {{glossary}}. Avoid localizing jokes.      │   ║
║   │   ... (content dimmed)                                               │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── OVERRIDES ─────────────────────────────────────────────────────────┐   ║
║   │   TEMPERATURE: [ 0.7 ] (Default)                                     │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ n ] CREATE NEW            [ c ] CLONE TO EDIT           [ ESC ] CANCEL   ║
╚══════════════════════════════════════════════════════════════════════════════╝

### Advanced Tab
╔══════════════════════════════════════════════════════════════════════════════╗
║   CONFIGURATION                                                              ║
╟──────────────────────────────────────────────────────────────────────────────╢
║   [ GENERAL ]   [ AI PROVIDERS ]   [ AI MODELS ]   [ PROMPTS ]               ║
║                                                                              ║
║   ┌── STORAGE & CACHE ───────────────────────────────────────────────────┐   ║
║   │   SEMANTIC CACHE:  [ 4.2 MB ]  (Approx. 15,000 cached terms)         │   ║
║   │   TEMP FILES:      [ 120 KB ]  (.bakasub.temp files)                 │   ║
║   │                                                                      │   ║
║   │   [ c ] CLEAR TRANSLATION CACHE (Forces re-translation next time)    │   ║
║   │   [ t ] CLEAR TEMP FILES        (Resets resume capability)           │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── DEBUGGING ─────────────────────────────────────────────────────────┐   ║
║   │   LOG LEVEL:       [ INFO ] (Default)                                │   ║
║   │                    [ DEBUG] (Verbose - Creates huge log files)       │   ║
║   │                                                                      │   ║
║   │   [ ] SAVE RAW JSON RESPONSES (For API debugging)                    │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
║                                                                              ║
║   ┌── SYSTEM INFO ───────────────────────────────────────────────────────┐   ║
║   │   VERSION: v1.0.0 (Build 20240520)                                   │   ║
║   │   GO VERSION: 1.22.2                                                 │   ║
║   │   PATH:    /usr/local/bin/bakasub                                    │   ║
║   └──────────────────────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════════════════════╣
║   [ ESC ] CANCEL                                            [ ENTER ] SAVE   ║
╚══════════════════════════════════════════════════════════════════════════════╝

      ┌──────────────────────────────────────────────────┐
      │  o   o                                    o   o  │ <--- Bobinas girando
      │ ┌──────────────────────────────────────────────┐ │
      │ │ ID: 12 | O: Olá mundo!   -> T: Hello world!  │ │ <--- Janela (Viewport)
      │ │ ID: 13 | O: Como vai?    -> T: How are you?  │ │      (Texto rolando)
      │ │ ID: 14 | O: Eu sou o...  -> T: I am the...   │ │
      │ └──────────────────────────────────────────────┘ │
      │  [░░░░░░░░░░░░░░░░░░░░░░░░░                  ]   │ <--- Progresso Físico
      └──────────────────────────────────────────────────┘