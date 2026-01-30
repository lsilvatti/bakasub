package dashboard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/pipeline"
	"github.com/lsilvatti/bakasub/internal/core/watcher"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/attachments"
	"github.com/lsilvatti/bakasub/internal/ui/execution"
	"github.com/lsilvatti/bakasub/internal/ui/glossary"
	"github.com/lsilvatti/bakasub/internal/ui/header"
	"github.com/lsilvatti/bakasub/internal/ui/job"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/picker"
	"github.com/lsilvatti/bakasub/internal/ui/remuxer"
	"github.com/lsilvatti/bakasub/internal/ui/review"
	"github.com/lsilvatti/bakasub/internal/ui/settings"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
	"github.com/lsilvatti/bakasub/pkg/utils"
)

// ViewState represents the current view being displayed
type ViewState int

const (
	// ViewDashboard shows the main dashboard
	ViewDashboard ViewState = iota
	// ViewPicker shows the file/directory picker
	ViewPicker
	// ViewSettings shows configuration screen
	ViewSettings
	// ViewJob shows job setup screen
	ViewJob
	// ViewExecution shows job execution screen
	ViewExecution
	// ViewHeader shows the header/metadata editor
	ViewHeader
	// ViewAttachments shows the attachment manager
	ViewAttachments
	// ViewRemuxer shows the quick remuxer
	ViewRemuxer
	// ViewGlossary shows the project glossary editor
	ViewGlossary
	// ViewReview shows the manual review editor
	ViewReview
)

// DirectoryAnalysis holds results of scanning a directory
type DirectoryAnalysis struct {
	Path     string
	MKVCount int
	SubCount int
	IsDir    bool
	Scanned  bool
	Scanning bool
	Error    error
}

// Model represents the dashboard state
type Model struct {
	width           int
	height          int
	selectedPath    string            // Currently selected path (from picker)
	analysis        DirectoryAnalysis // Results of directory scan
	selectedMode    int               // 0 = Full Process, 1 = Watch Mode
	apiOnline       bool
	cacheOK         bool
	updateAvailable bool
	ffmpegOK        bool
	mkvtoolnixOK    bool
	currentModel    string
	targetLang      string
	temperature     string
	kofiFlash       bool // Visual feedback when Ko-fi link is activated
	config          *config.Config

	// View management
	viewState      ViewState
	pickerModel    picker.Model
	settingsModel  settings.Model
	jobModel       job.Model
	executionModel execution.Model

	// Standalone module models
	headerModel      *header.Model
	attachmentsModel *attachments.Model
	remuxerModel     *remuxer.Model
	glossaryModel    *glossary.Model
	reviewModel      *review.Model

	// Module error (for displaying errors when opening modules)
	moduleError error

	// Update checker state
	latestVersion string
	releaseURL    string

	// Resume state (Smart Resume feature)
	resumeState     *pipeline.ResumeState
	showResumeModal bool

	// Watch Mode state
	watchModeActive  bool
	watchModePath    string
	watchModeWatcher *watcher.Watcher
	watchFileChan    chan string
}

// New creates a new dashboard model with default values
func New(cfg *config.Config) Model {
	if cfg == nil {
		cfg = config.Default()
	}

	// Apply interface language from config
	if cfg.InterfaceLang != "" {
		locales.Load(cfg.InterfaceLang)
	}

	return Model{
		width:           0,  // Will be set by WindowSizeMsg
		height:          0,  // Will be set by WindowSizeMsg
		selectedPath:    "", // No path selected initially
		analysis:        DirectoryAnalysis{},
		selectedMode:    0,
		apiOnline:       true,
		cacheOK:         true,
		updateAvailable: false, // Will be set by async update check
		ffmpegOK:        true,
		mkvtoolnixOK:    true,
		currentModel:    cfg.Model,
		targetLang:      cfg.TargetLang,
		temperature:     fmt.Sprintf("%.1f", cfg.Temperature),
		config:          cfg,
		viewState:       ViewDashboard,
		latestVersion:   "",
		releaseURL:      "",
	}
}

// Init initializes the dashboard
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Request current terminal size immediately
	cmds = append(cmds, tea.WindowSize())

	// Start async update check if enabled
	if m.config != nil && m.config.AutoCheckUpdates {
		cmds = append(cmds, utils.CheckForUpdates(utils.Version))
	}

	// Check for resume state files (.bakasub.temp)
	cmds = append(cmds, checkResumeState())

	return tea.Batch(cmds...)
}

// resumeStateFoundMsg is sent when a .bakasub.temp file is found
type resumeStateFoundMsg struct {
	state *pipeline.ResumeState
}

// checkResumeState looks for .bakasub.temp files in common locations
func checkResumeState() tea.Cmd {
	return func() tea.Msg {
		// Check home directory and current directory for temp files
		locations := []string{"."}
		if home, err := os.UserHomeDir(); err == nil {
			locations = append(locations, home)
		}

		for _, loc := range locations {
			tempPath := filepath.Join(loc, ".bakasub.temp")
			if _, err := os.Stat(tempPath); err == nil {
				// Found a temp file, try to load it
				state, err := pipeline.LoadResumeState(tempPath)
				if err == nil && state != nil {
					return resumeStateFoundMsg{state: state}
				}
			}
		}

		return nil
	}
}

// resumeDiscardedMsg is sent when user discards resume state
type resumeDiscardedMsg struct{}

// discardResumeState removes the .bakasub.temp file
func discardResumeState(state *pipeline.ResumeState) tea.Cmd {
	return func() tea.Msg {
		if state != nil && state.FilePath != "" {
			tempPath := filepath.Join(filepath.Dir(state.FilePath), ".bakasub.temp")
			os.Remove(tempPath)
		}
		return resumeDiscardedMsg{}
	}
}

// scanDirectoryMsg is sent when directory scanning is complete
type scanDirectoryMsg struct {
	path     string
	mkvCount int
	subCount int
	isDir    bool
	err      error
}

// kofiFlashMsg is sent to reset the Ko-fi flash visual feedback
type kofiFlashMsg struct{}

// scanDirectory scans the selected path for media files
func scanDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		files, err := picker.ScanDirectory(path)
		if err != nil {
			return scanDirectoryMsg{path: path, err: err}
		}

		// Count subtitle files separately
		subCount := 0
		mkvCount := len(files)

		return scanDirectoryMsg{
			path:     path,
			mkvCount: mkvCount,
			subCount: subCount,
			isDir:    true,
			err:      nil,
		}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate to sub-views if active
	switch m.viewState {
	case ViewPicker:
		return m.updatePicker(msg)
	case ViewSettings:
		return m.updateSettings(msg)
	case ViewJob:
		return m.updateJob(msg)
	case ViewExecution:
		return m.updateExecution(msg)
	case ViewHeader:
		return m.updateHeader(msg)
	case ViewAttachments:
		return m.updateAttachments(msg)
	case ViewRemuxer:
		return m.updateRemuxer(msg)
	case ViewGlossary:
		return m.updateGlossary(msg)
	case ViewReview:
		return m.updateReview(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle resume modal interactions first
		if m.showResumeModal {
			switch msg.String() {
			case "enter":
				// Resume from saved state
				if m.resumeState != nil {
					// Load job with resume state
					m.selectedPath = m.resumeState.FilePath
					m.jobModel = job.New(m.config, m.selectedPath)
					m.jobModel.SetSize(m.width, m.height)
					m.viewState = ViewJob
					m.showResumeModal = false
					return m, m.jobModel.Init()
				}
			case "d":
				// Discard resume state
				m.showResumeModal = false
				return m, discardResumeState(m.resumeState)
			case "esc":
				// Just close modal without action
				m.showResumeModal = false
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "o", " ":
			// Open file picker
			startDir := m.selectedPath
			if startDir == "" {
				startDir, _ = utils.GetHomeDir()
			}
			m.pickerModel = picker.New(startDir, picker.ModeDirectory)
			m.pickerModel.SetSize(m.width, m.height)
			m.viewState = ViewPicker
			return m, m.pickerModel.Init()

		case "c":
			// Open configuration/settings
			m.settingsModel = settings.New(m.config)
			m.settingsModel.SetSize(m.width, m.height)
			m.viewState = ViewSettings
			return m, m.settingsModel.Init()

		case "enter":
			// Start Full Process or Watch Mode based on selection
			if m.selectedPath == "" {
				// No path selected, do nothing
				return m, nil
			}
			if m.selectedMode == 0 {
				// Full Process - open Job Setup
				m.jobModel = job.New(m.config, m.selectedPath)
				m.jobModel.SetSize(m.width, m.height)
				m.viewState = ViewJob
				return m, m.jobModel.Init()
			}
			// Watch Mode - start file watcher with touchless processing
			if m.selectedPath != "" {
				return m, m.startWatchMode()
			}
			return m, nil

		case "1":
			// Extract Tracks (standalone module) - opens Job Setup in extraction mode
			if m.selectedPath != "" {
				m.jobModel = job.New(m.config, m.selectedPath)
				m.jobModel.SetSize(m.width, m.height)
				m.viewState = ViewJob
				return m, m.jobModel.Init()
			}
			return m, nil

		case "2":
			// Translate Subtitle (standalone module) - opens Job Setup
			if m.selectedPath != "" {
				m.jobModel = job.New(m.config, m.selectedPath)
				m.jobModel.SetSize(m.width, m.height)
				m.viewState = ViewJob
				return m, m.jobModel.Init()
			}
			return m, nil

		case "3":
			// Mux Container (standalone module) - opens Remuxer
			if m.selectedPath != "" {
				model, err := remuxer.New(m.findFirstMKV())
				if err != nil {
					m.moduleError = err
					return m, nil
				}
				m.remuxerModel = model
				m.remuxerModel.SetSize(m.width, m.height)
				m.viewState = ViewRemuxer
				return m, m.remuxerModel.Init()
			}
			return m, nil

		case "4":
			// Manual Review (Editor)
			if m.selectedPath != "" {
				// Try to find subtitle files in the directory
				subFiles := m.findSubtitleFiles()
				if len(subFiles) == 0 {
					m.moduleError = fmt.Errorf("%s", locales.T("review.no_subtitle_found"))
					return m, nil
				}
				// Open review with the first subtitle file found
				reviewModel, err := review.New(subFiles[0], "")
				if err != nil {
					m.moduleError = err
					return m, nil
				}
				m.reviewModel = reviewModel
				m.reviewModel.SetSize(m.width, m.height)
				m.viewState = ViewReview
				return m, m.reviewModel.Init()
			}
			return m, nil

		case "5":
			// Edit Flags / Metadata (Header Editor)
			if m.selectedPath != "" {
				model, err := header.New(m.findFirstMKV())
				if err != nil {
					m.moduleError = err
					return m, nil
				}
				m.headerModel = model
				m.headerModel.SetSize(m.width, m.height)
				m.viewState = ViewHeader
				return m, m.headerModel.Init()
			}
			return m, nil

		case "6":
			// Manage Attachments
			if m.selectedPath != "" {
				model, err := attachments.New(m.findFirstMKV())
				if err != nil {
					m.moduleError = err
					return m, nil
				}
				m.attachmentsModel = model
				m.attachmentsModel.SetSize(m.width, m.height)
				m.viewState = ViewAttachments
				return m, m.attachmentsModel.Init()
			}
			return m, nil

		case "7":
			// Add/Remove Tracks (Remuxer)
			if m.selectedPath != "" {
				model, err := remuxer.New(m.findFirstMKV())
				if err != nil {
					m.moduleError = err
					return m, nil
				}
				m.remuxerModel = model
				m.remuxerModel.SetSize(m.width, m.height)
				m.viewState = ViewRemuxer
				return m, m.remuxerModel.Init()
			}
			return m, nil

		case "8":
			// Project Glossary
			if m.selectedPath != "" {
				glossaryPath := filepath.Join(filepath.Dir(m.findFirstMKV()), "glossary.json")
				m.glossaryModel = glossary.New(glossaryPath)
				m.glossaryModel.SetSize(m.width, m.height)
				m.viewState = ViewGlossary
				return m, m.glossaryModel.Init()
			}
			return m, nil

		case "m":
			// Change Model - open settings at Models tab
			m.settingsModel = settings.New(m.config)
			m.settingsModel.SetSize(m.width, m.height)
			m.viewState = ViewSettings
			return m, m.settingsModel.Init()

		case "k":
			// Open Ko-fi link in browser
			kofiURL := "https://ko-fi.com/lsilvatti"

			// Open browser (non-blocking)
			go utils.OpenURL(kofiURL)

			// Enable flash visual feedback
			m.kofiFlash = true

			// Schedule flash reset after 300ms
			return m, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
				return kofiFlashMsg{}
			})

		case "up", "down":
			// Toggle mode selection
			if m.selectedMode == 0 {
				m.selectedMode = 1
			} else {
				m.selectedMode = 0
			}
		}

	case utils.MsgUpdateAvailable:
		// Update available from async check
		m.updateAvailable = true
		m.latestVersion = msg.LatestVersion
		m.releaseURL = msg.ReleaseURL

	case utils.MsgUpdateCheckFailed:
		// Update check failed, silently ignore
		m.updateAvailable = false

	case resumeStateFoundMsg:
		// Show resume modal if valid state found
		if msg.state != nil {
			m.resumeState = msg.state
			m.showResumeModal = true
		}

	case resumeDiscardedMsg:
		// User discarded the resume state
		m.resumeState = nil
		m.showResumeModal = false

	case watchModeStartedMsg:
		// Watch mode has started - store watcher and start listening for events
		m.watchModeActive = true
		m.watchModePath = msg.path
		m.watchModeWatcher = msg.watcher
		m.watchFileChan = msg.fileChan
		m.moduleError = nil
		// Start listening for watch mode events
		return m, listenWatchModeFiles(msg.fileChan)

	case watchModeFileMsg:
		// New file detected in watch mode - auto-start job
		m.jobModel = job.New(m.config, msg.path)
		m.jobModel.SetSize(m.width, m.height)
		m.viewState = ViewJob
		// Continue listening for more files
		return m, tea.Batch(m.jobModel.Init(), listenWatchModeFiles(m.watchFileChan))

	case watchModeErrorMsg:
		// Watch mode error
		m.watchModeActive = false
		if m.watchModeWatcher != nil {
			m.watchModeWatcher.Stop()
			m.watchModeWatcher = nil
		}
		m.moduleError = msg.err

	case scanDirectoryMsg:
		// Update analysis results
		m.analysis = DirectoryAnalysis{
			Path:     msg.path,
			MKVCount: msg.mkvCount,
			SubCount: msg.subCount,
			IsDir:    msg.isDir,
			Scanned:  true,
			Scanning: false,
			Error:    msg.err,
		}

	case kofiFlashMsg:
		// Reset Ko-fi flash
		m.kofiFlash = false

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// updatePicker handles updates when in picker view
func (m Model) updatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case picker.SelectedPathMsg:
		if msg.Aborted {
			// User cancelled, return to dashboard
			m.viewState = ViewDashboard
			return m, nil
		}
		// Path selected, update and scan
		m.selectedPath = msg.Path
		m.analysis.Scanning = true
		m.viewState = ViewDashboard
		return m, scanDirectory(msg.Path)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to picker
		var cmd tea.Cmd
		pickerModel, cmd := m.pickerModel.Update(msg)
		m.pickerModel = pickerModel.(picker.Model)
		return m, cmd
	}

	// Delegate to picker
	var cmd tea.Cmd
	pickerModel, cmd := m.pickerModel.Update(msg)
	m.pickerModel = pickerModel.(picker.Model)
	return m, cmd
}

// updateSettings handles updates when in settings view
func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case settings.SavedMsg:
		// Settings were saved, update our config reference
		if msg.Config != nil {
			m.config = msg.Config
			m.currentModel = msg.Config.Model
			m.targetLang = msg.Config.TargetLang
			m.temperature = fmt.Sprintf("%.1f", msg.Config.Temperature)
			// Reload locales if interface language changed
			if msg.Config.InterfaceLang != "" {
				locales.Load(msg.Config.InterfaceLang)
			}
		}
		m.viewState = ViewDashboard
		return m, nil

	case settings.CancelledMsg:
		// Settings were cancelled
		m.viewState = ViewDashboard
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to settings
		m.settingsModel.SetSize(msg.Width, msg.Height)
	}

	// Delegate all key handling to settings (including ESC)
	var cmd tea.Cmd
	settingsModel, cmd := m.settingsModel.Update(msg)
	m.settingsModel = settingsModel.(settings.Model)
	return m, cmd
}

// updateJob handles updates when in job setup view
func (m Model) updateJob(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			// Return to dashboard
			m.viewState = ViewDashboard
			return m, nil
		}

	case job.StartJobMsg:
		// Job started, transition to execution view with full job config
		execConfig := execution.JobConfig{
			InputPath:       msg.JobConfig.InputPath,
			BatchMode:       msg.JobConfig.BatchMode,
			SourceLang:      msg.JobConfig.SourceLang,
			TargetLang:      msg.JobConfig.TargetLang,
			MediaType:       msg.JobConfig.MediaType,
			AIModel:         msg.JobConfig.AIModel,
			Temperature:     msg.JobConfig.Temperature,
			GlossaryPath:    msg.JobConfig.GlossaryPath,
			GlossaryTerms:   msg.JobConfig.GlossaryTerms,
			RemoveHITags:    msg.JobConfig.RemoveHITags,
			MuxMode:         msg.JobConfig.MuxMode,
			SetDefault:      msg.JobConfig.SetDefault,
			BackupOriginal:  msg.JobConfig.BackupOriginal,
			ExtractFonts:    msg.JobConfig.ExtractFonts,
			AutoDetectTrack: msg.JobConfig.AutoDetectTrack,
		}

		// Convert analyzed files
		for _, f := range msg.JobConfig.Files {
			execConfig.Files = append(execConfig.Files, execution.AnalyzedFile{
				Path:            f.Path,
				Filename:        f.Filename,
				SelectedTrackID: f.SelectedTrackID,
			})
		}

		m.executionModel = execution.New(m.config, execConfig)
		m.executionModel.SetSize(m.width, m.height)
		m.viewState = ViewExecution
		return m, m.executionModel.Init()

	case job.CancelledMsg:
		// Job was cancelled
		m.viewState = ViewDashboard
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to job
		m.jobModel.SetSize(msg.Width, msg.Height)
	}

	// Delegate to job
	var cmd tea.Cmd
	jobModel, cmd := m.jobModel.Update(msg)
	m.jobModel = jobModel.(job.Model)
	return m, cmd
}

// updateExecution handles updates when in execution view
func (m Model) updateExecution(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Only allow ESC to return to dashboard if job is complete or failed
		// The execution model handles its own ESC for cancellation
		if msg.String() == "q" {
			// Let execution model handle the quit key
		}

	case execution.CompletedMsg:
		// Execution completed, return to dashboard
		m.viewState = ViewDashboard
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to execution
		m.executionModel.SetSize(msg.Width, msg.Height)
	}

	// Delegate to execution
	var cmd tea.Cmd
	executionModel, cmd := m.executionModel.Update(msg)
	m.executionModel = executionModel.(execution.Model)
	return m, cmd
}

// updateHeader handles updates when in header editor view
func (m Model) updateHeader(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.headerModel == nil {
		m.viewState = ViewDashboard
		return m, nil
	}

	switch msg := msg.(type) {
	case header.ClosedMsg:
		m.viewState = ViewDashboard
		m.headerModel = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to header
		if m.headerModel != nil {
			m.headerModel.SetSize(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	headerModel, cmd := m.headerModel.Update(msg)
	if hm, ok := headerModel.(header.Model); ok {
		m.headerModel = &hm
	}
	return m, cmd
}

// updateAttachments handles updates when in attachments manager view
func (m Model) updateAttachments(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.attachmentsModel == nil {
		m.viewState = ViewDashboard
		return m, nil
	}

	switch msg := msg.(type) {
	case attachments.ClosedMsg:
		m.viewState = ViewDashboard
		m.attachmentsModel = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to attachments
		if m.attachmentsModel != nil {
			m.attachmentsModel.SetSize(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	attachmentsModel, cmd := m.attachmentsModel.Update(msg)
	if am, ok := attachmentsModel.(attachments.Model); ok {
		m.attachmentsModel = &am
	}
	return m, cmd
}

// updateRemuxer handles updates when in remuxer view
func (m Model) updateRemuxer(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.remuxerModel == nil {
		m.viewState = ViewDashboard
		return m, nil
	}

	switch msg := msg.(type) {
	case remuxer.ClosedMsg:
		m.viewState = ViewDashboard
		m.remuxerModel = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to remuxer
		if m.remuxerModel != nil {
			m.remuxerModel.SetSize(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	remuxerModel, cmd := m.remuxerModel.Update(msg)
	if rm, ok := remuxerModel.(remuxer.Model); ok {
		m.remuxerModel = &rm
	}
	return m, cmd
}

// updateGlossary handles updates when in glossary editor view
func (m Model) updateGlossary(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.glossaryModel == nil {
		m.viewState = ViewDashboard
		return m, nil
	}

	switch msg := msg.(type) {
	case glossary.ClosedMsg:
		m.viewState = ViewDashboard
		m.glossaryModel = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.glossaryModel.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	glossaryModel, cmd := m.glossaryModel.Update(msg)
	// Glossary uses pointer receiver, so the type assertion is correct
	if gm, ok := glossaryModel.(*glossary.Model); ok {
		m.glossaryModel = gm
	}
	return m, cmd
}

// updateReview handles updates when in review editor view
func (m Model) updateReview(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.reviewModel == nil {
		m.viewState = ViewDashboard
		return m, nil
	}

	switch msg := msg.(type) {
	case review.ClosedMsg:
		m.viewState = ViewDashboard
		m.reviewModel = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Propagate size to review
		if m.reviewModel != nil {
			m.reviewModel.SetSize(msg.Width, msg.Height)
		}
	}

	var cmd tea.Cmd
	reviewModel, cmd := m.reviewModel.Update(msg)
	// Review uses value receiver, so we need to convert
	if rm, ok := reviewModel.(review.Model); ok {
		m.reviewModel = &rm
	}
	return m, cmd
}

// findFirstMKV returns the path to the first MKV file in the selected path
func (m Model) findFirstMKV() string {
	if m.selectedPath == "" {
		return ""
	}

	// If the selected path is already an MKV file
	if strings.HasSuffix(strings.ToLower(m.selectedPath), ".mkv") {
		return m.selectedPath
	}

	// If it's a directory, find the first MKV file
	files, err := picker.ScanDirectory(m.selectedPath)
	if err != nil || len(files) == 0 {
		return m.selectedPath // Return as-is and let the module handle the error
	}

	return files[0]
}

// findSubtitleFiles returns subtitle files in the selected directory
func (m Model) findSubtitleFiles() []string {
	if m.selectedPath == "" {
		return nil
	}

	dir := m.selectedPath
	info, err := os.Stat(dir)
	if err != nil {
		return nil
	}

	if !info.IsDir() {
		dir = filepath.Dir(dir)
	}

	var subtitles []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, ".ass") || strings.HasSuffix(name, ".srt") ||
			strings.HasSuffix(name, ".ssa") || strings.HasSuffix(name, ".sub") {
			subtitles = append(subtitles, filepath.Join(dir, entry.Name()))
		}
	}

	return subtitles
}

// watchModeStartedMsg is sent when watch mode starts
type watchModeStartedMsg struct {
	path     string
	watcher  *watcher.Watcher
	fileChan chan string
}

// watchModeFileMsg is sent when a new file is detected
type watchModeFileMsg struct {
	path string
}

// watchModeErrorMsg is sent when an error occurs in watch mode
type watchModeErrorMsg struct {
	err error
}

// startWatchMode initiates the directory watcher
func (m Model) startWatchMode() tea.Cmd {
	selectedPath := m.selectedPath

	return func() tea.Msg {
		// Create watcher
		w, err := watcher.New(selectedPath)
		if err != nil {
			return watchModeErrorMsg{err: err}
		}

		// Create channel for file events
		fileChan := make(chan string, 10)

		// Set callback for new files
		w.OnNewFile = func(filePath string) {
			select {
			case fileChan <- filePath:
			default:
				// Channel full, skip
			}
		}

		// Set callback for errors
		w.OnError = func(err error) {
			// Errors are logged but don't stop the watcher
		}

		// Start watching
		if err := w.Start(); err != nil {
			return watchModeErrorMsg{err: err}
		}

		// Return started message with channel reference
		return watchModeStartedMsg{
			path:     selectedPath,
			watcher:  w,
			fileChan: fileChan,
		}
	}
}

// listenWatchMode creates a command that listens for watch mode events
func listenWatchModeFiles(fileChan chan string) tea.Cmd {
	if fileChan == nil {
		return nil
	}
	return func() tea.Msg {
		filePath, ok := <-fileChan
		if !ok {
			return nil
		}
		return watchModeFileMsg{path: filePath}
	}
}

// View renders the dashboard
func (m Model) View() string {
	// Wait for terminal size
	if layout.IsWaitingForSize(m.width, m.height) {
		return locales.T("common.loading")
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	// Render the appropriate view based on state
	switch m.viewState {
	case ViewPicker:
		return m.pickerModel.View()
	case ViewSettings:
		return m.settingsModel.View()
	case ViewJob:
		return m.jobModel.View()
	case ViewExecution:
		return m.executionModel.View()
	case ViewHeader:
		if m.headerModel != nil {
			return m.headerModel.View()
		}
	case ViewAttachments:
		if m.attachmentsModel != nil {
			return m.attachmentsModel.View()
		}
	case ViewRemuxer:
		if m.remuxerModel != nil {
			return m.remuxerModel.View()
		}
	case ViewGlossary:
		if m.glossaryModel != nil {
			return m.glossaryModel.View()
		}
	case ViewReview:
		if m.reviewModel != nil {
			return m.reviewModel.View()
		}
	}

	// Calculate content width (full width minus border/padding)
	contentWidth := m.width - 4 // 2 for border, 2 for padding

	// Build all sections for main dashboard
	header := m.renderHeader(contentWidth)
	inputMode := m.renderInputMode(contentWidth)
	modulesToolbox := m.renderModulesToolbox(contentWidth)
	systemAI := m.renderSystemAI(contentWidth)
	footer := m.renderFooter(contentWidth)

	// Compose the full dashboard
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		inputMode,
		"",
		modulesToolbox,
		"",
		systemAI,
		"",
		footer,
	)

	// Wrap in main window with double border - full width
	dashboard := styles.MainWindow.Width(contentWidth).Render(content)

	// If resume modal is active, overlay it on top of the dashboard
	if m.showResumeModal && m.resumeState != nil {
		return m.renderResumeModal(dashboard, contentWidth)
	}

	return dashboard
}

// renderHeader creates the ASCII logo and status bar
func (m Model) renderHeader(contentWidth int) string {
	// ASCII Art Logo
	logo := []string{
		" ____        _          ____        _     ",
		"| __ )  __ _| | ____ _ / ___| _   _| |__  ",
		"|  _ \\ / _` | |/ / _` |\\___ \\| | | | '_ \\ ",
		"| |_) | (_| |   < (_| | ___) | |_| | |_) |",
		"|____/ \\__,_|_|\\_\\__,_||____/ \\__,_|_.__/ ",
	}

	// Style the logo
	styledLogo := ""
	for _, line := range logo {
		styledLogo += styles.Logo.Render(line) + "\n"
	}

	// Status indicators
	apiStatus := styles.StatusOK.Render("[" + locales.T("dashboard.status.api_online") + "]")
	if !m.apiOnline {
		apiStatus = styles.StatusError.Render("[" + locales.T("dashboard.status.api_offline") + "]")
	}

	cacheStatus := styles.StatusOK.Render("[" + locales.T("dashboard.status.cache_ok") + "]")
	if !m.cacheOK {
		cacheStatus = styles.StatusError.Render("[" + locales.T("dashboard.status.cache_error") + "]")
	}

	statusLine := apiStatus + " " + cacheStatus

	// Update notification
	updateNotice := ""
	if m.updateAvailable && m.latestVersion != "" {
		updateNotice = styles.StatusWarning.Render(fmt.Sprintf("[!] %s (%s)", locales.T("dashboard.status.update_available"), m.latestVersion))
	} else if m.updateAvailable {
		updateNotice = styles.StatusWarning.Render("[!] " + locales.T("dashboard.status.update_available"))
	}

	// Status block for right side
	statusBlock := lipgloss.JoinVertical(
		lipgloss.Right,
		utils.Version,
		styles.StatusOK.Render("["+locales.T("dashboard.status.app_running")+"]"),
		"",
		statusLine,
		updateNotice,
	)

	// Create header with shaded background using block chars
	headerBar := strings.Repeat("▒", contentWidth-4)

	// Build title line with logo on left and status on right
	titleArea := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styledLogo,
		strings.Repeat(" ", 10),
		statusBlock,
	)

	headerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		titleArea,
		headerBar,
	)

	return headerContent
}

// renderInputMode creates the Input & Mode section
func (m Model) renderInputMode(contentWidth int) string {
	title := styles.PanelTitle.Render(locales.T("dashboard.input.title_numbered"))

	// Path display (read-only) with open hotkey
	pathLabel := locales.T("dashboard.input.path_label") + " > "
	var pathValue string
	if m.selectedPath == "" {
		// No path selected - show placeholder in dim gray
		pathValue = styles.Dimmed.Render(locales.T("dashboard.input.no_file_selected"))
	} else if m.analysis.Scanning {
		// Currently scanning
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " + styles.StatusWarning.Render("["+locales.T("dashboard.input.scanning")+"]")
	} else if m.analysis.Scanned && m.analysis.Error != nil {
		// Scan error
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " + styles.StatusError.Render("["+locales.T("dashboard.input.error")+"]")
	} else if m.analysis.Scanned {
		// Scan complete - show file count
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " +
			styles.StatusOK.Render(fmt.Sprintf("["+locales.T("dashboard.input.mkv_count")+"]", m.analysis.MKVCount))
	} else {
		pathValue = styles.CodeBlock.Render(m.selectedPath)
	}

	pathLine := pathLabel + pathValue + "    " + styles.RenderHotkey("o", locales.T("dashboard.input.open_path"))

	// Mode selection
	fullModeIcon := "( )"
	watchModeIcon := "( )"
	if m.selectedMode == 0 {
		fullModeIcon = "(o)"
	} else {
		watchModeIcon = "(o)"
	}

	fullMode := styles.Highlight.Render(fullModeIcon) + " " + locales.T("dashboard.input.mode_full")
	fullModeDesc := "    " + styles.Dimmed.Render(locales.T("dashboard.input.mode_full_help"))

	watchMode := styles.Highlight.Render(watchModeIcon) + " " + locales.T("dashboard.input.mode_watch")

	modeContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		pathLine,
		"",
		fullMode,
		fullModeDesc,
		watchMode,
	)

	panelContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		modeContent,
	)

	return styles.Panel.Width(contentWidth - 6).Render(panelContent)
}

// renderModulesToolbox creates the side-by-side Modules and Toolbox sections
func (m Model) renderModulesToolbox(contentWidth int) string {
	// Modules Panel
	modulesTitle := styles.PanelTitle.Render(locales.T("dashboard.modules.title_numbered"))
	modulesItems := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.RenderHotkey("1", locales.T("dashboard.modules.extract")),
		styles.RenderHotkey("2", locales.T("dashboard.modules.translate")),
		styles.RenderHotkey("3", locales.T("dashboard.modules.mux")),
		styles.RenderHotkey("4", locales.T("dashboard.modules.review")),
	)
	modulesContent := lipgloss.JoinVertical(lipgloss.Left, modulesTitle, modulesItems)
	halfWidth := (contentWidth - 10) / 2
	modulesPanel := styles.Panel.Width(halfWidth).Render(modulesContent)

	// Toolbox Panel
	toolboxTitle := styles.PanelTitle.Render(locales.T("dashboard.toolbox.title_numbered"))
	toolboxItems := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.RenderHotkey("5", locales.T("dashboard.toolbox.edit_flags")),
		styles.RenderHotkey("6", locales.T("dashboard.toolbox.manage_attachments")),
		styles.RenderHotkey("7", locales.T("dashboard.toolbox.track_manager")),
		styles.RenderHotkey("8", locales.T("dashboard.toolbox.glossary")),
	)
	toolboxContent := lipgloss.JoinVertical(lipgloss.Left, toolboxTitle, toolboxItems)
	toolboxPanel := styles.Panel.Width(halfWidth).Render(toolboxContent)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, modulesPanel, toolboxPanel)
}

// renderSystemAI creates the System & AI info section
func (m Model) renderSystemAI(contentWidth int) string {
	title := styles.PanelTitle.Render(locales.T("dashboard.system.title_numbered"))

	// Model info
	modelLine := locales.T("dashboard.system.model") + " " + styles.CodeBlock.Render(m.currentModel) +
		"        " + styles.RenderHotkey("m", locales.T("dashboard.system.change_model"))

	// Settings info
	settingsLine := locales.T("dashboard.system.target") + " " + styles.Highlight.Render(m.targetLang) +
		"  │  " + locales.T("dashboard.system.temp") + " " + styles.Highlight.Render(m.temperature) +
		"             " + styles.RenderHotkey("c", locales.T("dashboard.system.configuration"))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		modelLine,
		settingsLine,
	)

	return styles.Panel.Width(contentWidth - 6).Render(content)
}

// renderFooter creates the bottom status bar
func (m Model) renderFooter(contentWidth int) string {
	// Left side: Ko-fi link with flash feedback
	kofiText := locales.T("dashboard.footer.kofi")
	if m.kofiFlash {
		// Flash cyan when activated
		kofiText = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true).Render(locales.T("dashboard.footer.kofi"))
		leftSide := styles.RenderHotkey("k", "") + kofiText
		return renderFooterContent(leftSide, m)
	}
	leftSide := styles.RenderHotkey("k", locales.T("dashboard.footer.kofi"))

	return renderFooterContent(leftSide, m)
}

// renderFooterContent builds the footer content (extracted for reuse)
func renderFooterContent(leftSide string, m Model) string {
	// Middle: Dependency status
	ffmpegStatus := locales.T("dashboard.footer.ffmpeg_ok")
	if !m.ffmpegOK {
		ffmpegStatus = "FFmpeg " + styles.StatusError.Render("["+locales.T("status_indicators.missing")+"]")
	}

	mkvStatus := locales.T("dashboard.footer.mkvtoolnix_ok")
	if !m.mkvtoolnixOK {
		mkvStatus = "MKVToolNix " + styles.StatusError.Render("["+locales.T("status_indicators.missing")+"]")
	}

	middle := fmt.Sprintf("%s %s %s", locales.T("dashboard.footer.deps"), ffmpegStatus, mkvStatus)

	// Right side: Quit
	rightSide := styles.RenderHotkey("q", locales.T("dashboard.footer.quit"))

	// Join with separators
	footer := leftSide + "  │  " + middle + "  │  " + rightSide

	return styles.Footer.Render(footer)
}

// renderResumeModal renders the smart resume modal overlay
func (m Model) renderResumeModal(background string, contentWidth int) string {
	if m.resumeState == nil {
		return background
	}

	// Calculate modal dimensions
	modalWidth := 56
	if contentWidth < modalWidth+4 {
		modalWidth = contentWidth - 4
	}

	// Extract filename from path
	filename := filepath.Base(m.resumeState.FilePath)
	if len(filename) > modalWidth-10 {
		filename = filename[:modalWidth-13] + "..."
	}

	// Calculate percentage complete
	percent := 0
	if m.resumeState.TotalBatches > 0 {
		percent = (m.resumeState.CompletedBatches * 100) / m.resumeState.TotalBatches
	}

	// Build modal content
	titleStyle := lipgloss.NewStyle().Foreground(styles.Cyan).Bold(true)
	title := titleStyle.Render(locales.T("resume_modal.title"))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		locales.T("resume_modal.message"),
		"",
		fmt.Sprintf("%s %s", locales.T("resume_modal.file_label"), filename),
		fmt.Sprintf("%s %d / %d (%d%% %s)", locales.T("resume_modal.batch_label"), m.resumeState.CompletedBatches, m.resumeState.TotalBatches, percent, locales.T("resume_modal.complete")),
		locales.T("resume_modal.cache_label")+" "+locales.T("resume_modal.cache_saved"),
		"",
		locales.Tf("resume_modal.resume_question", m.resumeState.CompletedBatches),
		"",
	)

	// Controls
	controls := lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.RenderHotkey("d", locales.T("resume_modal.discard_restart")),
		"       ",
		styles.RenderHotkey("ENTER", locales.T("resume_modal.resume")),
	)

	// Compose modal
	modalContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		content,
		controls,
	)

	// Style modal with border
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.Cyan).
		Padding(1, 2).
		Width(modalWidth)

	modal := modalStyle.Render(modalContent)

	// Center modal on screen
	modalHeight := lipgloss.Height(modal)
	modalLines := strings.Split(modal, "\n")

	// Get background lines
	bgLines := strings.Split(background, "\n")

	// Calculate vertical offset to center
	startY := (len(bgLines) - modalHeight) / 2
	if startY < 0 {
		startY = 0
	}

	// Calculate horizontal offset to center
	startX := (m.width - modalWidth - 4) / 2
	if startX < 0 {
		startX = 0
	}

	// Overlay modal on background (simple version - just place it in center)
	result := make([]string, len(bgLines))
	for i, line := range bgLines {
		if i >= startY && i < startY+modalHeight {
			modalLineIdx := i - startY
			if modalLineIdx < len(modalLines) {
				// Pad modal line to position it correctly
				paddedModal := strings.Repeat(" ", startX) + modalLines[modalLineIdx]
				// Use dimmed background and overlay modal
				result[i] = paddedModal
			} else {
				result[i] = line
			}
		} else {
			result[i] = line
		}
	}

	return strings.Join(result, "\n")
}
