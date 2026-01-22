package dashboard

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/picker"
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
	viewState   ViewState
	pickerModel picker.Model
}

// New creates a new dashboard model with default values
func New(cfg *config.Config) Model {
	if cfg == nil {
		cfg = config.Default()
	}

	return Model{
		width:           80,
		height:          24,
		selectedPath:    "", // No path selected initially
		analysis:        DirectoryAnalysis{},
		selectedMode:    0,
		apiOnline:       true,
		cacheOK:         true,
		updateAvailable: true,
		ffmpegOK:        true,
		mkvtoolnixOK:    true,
		currentModel:    cfg.Model,
		targetLang:      cfg.TargetLang,
		temperature:     fmt.Sprintf("%.1f", cfg.Temperature),
		config:          cfg,
		viewState:       ViewDashboard,
	}
}

// Init initializes the dashboard
func (m Model) Init() tea.Cmd {
	return nil
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
	// If we're in picker view, delegate to picker
	if m.viewState == ViewPicker {
		return m.updatePicker(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			m.viewState = ViewPicker
			return m, m.pickerModel.Init()

		case "k":
			// Open Ko-fi link in browser
			kofiURL := "https://ko-fi.com/bakasub"
			if m.config != nil && m.config.KofiUsername != "" {
				kofiURL = "https://ko-fi.com/" + m.config.KofiUsername
			}

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

// View renders the dashboard
func (m Model) View() string {
	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	// If we're in picker view, render the picker
	if m.viewState == ViewPicker {
		return m.pickerModel.View()
	}

	// Build all sections
	header := m.renderHeader()
	inputMode := m.renderInputMode()
	modulesToolbox := m.renderModulesToolbox()
	systemAI := m.renderSystemAI()
	footer := m.renderFooter()

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

	// Wrap in main window with double border (use safe width)
	contentWidth := layout.SafeWidth(m.width-4, 76)
	return styles.MainWindow.Width(contentWidth).Render(content)
}

// renderHeader creates the ASCII logo and status bar
func (m Model) renderHeader() string {
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
	apiStatus := styles.StatusOK.Render("[API: ONLINE]")
	if !m.apiOnline {
		apiStatus = styles.StatusError.Render("[API: OFFLINE]")
	}

	cacheStatus := styles.StatusOK.Render("[CACHE: OK]")
	if !m.cacheOK {
		cacheStatus = styles.StatusError.Render("[CACHE: ERROR]")
	}

	statusLine := apiStatus + " " + cacheStatus

	// Update notification
	updateNotice := ""
	if m.updateAvailable {
		updateNotice = styles.StatusWarning.Render("[!] UPDATE AVAILABLE (v1.1)")
	}

	// Status block for right side
	statusBlock := lipgloss.JoinVertical(
		lipgloss.Right,
		"v1.0.0",
		styles.StatusOK.Render("[APP RUNNING]"),
		"",
		statusLine,
		updateNotice,
	)

	// Create header with shaded background using block chars
	headerBar := strings.Repeat("▒", m.width-4)

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
func (m Model) renderInputMode() string {
	title := styles.PanelTitle.Render("1. INPUT & MODE")

	// Path display (read-only) with open hotkey
	pathLabel := "PATH > "
	var pathValue string
	if m.selectedPath == "" {
		// No path selected - show placeholder in dim gray
		pathValue = styles.Dimmed.Render("No file selected")
	} else if m.analysis.Scanning {
		// Currently scanning
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " + styles.StatusWarning.Render("[SCANNING...]")
	} else if m.analysis.Scanned && m.analysis.Error != nil {
		// Scan error
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " + styles.StatusError.Render("[ERROR]")
	} else if m.analysis.Scanned {
		// Scan complete - show file count
		pathValue = styles.CodeBlock.Render(m.selectedPath) + " " +
			styles.StatusOK.Render(fmt.Sprintf("[%d MKV]", m.analysis.MKVCount))
	} else {
		pathValue = styles.CodeBlock.Render(m.selectedPath)
	}

	pathLine := pathLabel + pathValue + "    " + styles.RenderHotkey("o", "OPEN PATH")

	// Mode selection
	fullModeIcon := "( )"
	watchModeIcon := "( )"
	if m.selectedMode == 0 {
		fullModeIcon = "(o)"
	} else {
		watchModeIcon = "(o)"
	}

	fullMode := styles.Highlight.Render(fullModeIcon) + " FULL PROCESS (Extract -> Translate -> Mux)"
	fullModeDesc := "    " + styles.Dimmed.Render("*Opens \"Job Setup\" screen for uninterrupted processing.")

	watchMode := styles.Highlight.Render(watchModeIcon) + " WATCH MODE (Auto-process new files in folder)"

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

	panelWidth := layout.SafeWidth(m.width-10, 60)
	return styles.Panel.Width(panelWidth).Render(panelContent)
}

// renderModulesToolbox creates the side-by-side Modules and Toolbox sections
func (m Model) renderModulesToolbox() string {
	// Modules Panel
	modulesTitle := styles.PanelTitle.Render("2. MODULES (Standalone)")
	modulesItems := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.RenderHotkey("1", "EXTRACT TRACKS"),
		styles.RenderHotkey("2", "TRANSLATE SUBTITLE"),
		styles.RenderHotkey("3", "MUX CONTAINER"),
		styles.RenderHotkey("4", "MANUAL REVIEW (Editor)"),
	)
	modulesContent := lipgloss.JoinVertical(lipgloss.Left, modulesTitle, modulesItems)
	halfWidth := layout.CalculateHalf(m.width, 6)
	halfWidth = layout.SafeWidth(halfWidth, 30)
	modulesPanel := styles.Panel.Width(halfWidth).Render(modulesContent)

	// Toolbox Panel
	toolboxTitle := styles.PanelTitle.Render("3. TOOLBOX (MKVToolNix)")
	toolboxItems := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.RenderHotkey("5", "EDIT FLAGS / METADATA"),
		styles.RenderHotkey("6", "MANAGE ATTACHMENTS"),
		styles.RenderHotkey("7", "ADD/REMOVE TRACKS"),
		styles.RenderHotkey("8", "PROJECT GLOSSARY"),
	)
	toolboxContent := lipgloss.JoinVertical(lipgloss.Left, toolboxTitle, toolboxItems)
	toolboxPanel := styles.Panel.Width(halfWidth).Render(toolboxContent)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, modulesPanel, toolboxPanel)
}

// renderSystemAI creates the System & AI info section
func (m Model) renderSystemAI() string {
	title := styles.PanelTitle.Render("4. SYSTEM & AI")

	// Model info
	modelLine := "MODEL: " + styles.CodeBlock.Render(m.currentModel) +
		"        " + styles.RenderHotkey("m", "CHANGE MODEL")

	// Settings info
	settingsLine := "TARGET: " + styles.Highlight.Render(m.targetLang) +
		"  │  TEMP: " + styles.Highlight.Render(m.temperature) +
		"             " + styles.RenderHotkey("c", "CONFIGURATION")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		modelLine,
		settingsLine,
	)

	panelWidth := layout.SafeWidth(m.width-10, 60)
	return styles.Panel.Width(panelWidth).Render(content)
}

// renderFooter creates the bottom status bar
func (m Model) renderFooter() string {
	// Left side: Ko-fi link with flash feedback
	kofiText := "KO-FI"
	if m.kofiFlash {
		// Flash cyan when activated
		kofiText = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Bold(true).Render("KO-FI")
		leftSide := styles.RenderHotkey("k", "") + kofiText
		return renderFooterContent(leftSide, m)
	}
	leftSide := styles.RenderHotkey("k", "KO-FI")

	return renderFooterContent(leftSide, m)
}

// renderFooterContent builds the footer content (extracted for reuse)
func renderFooterContent(leftSide string, m Model) string {
	// Middle: Dependency status
	ffmpegStatus := "FFmpeg [OK]"
	if !m.ffmpegOK {
		ffmpegStatus = "FFmpeg " + styles.StatusError.Render("[MISSING]")
	}

	mkvStatus := "MKVToolNix [OK]"
	if !m.mkvtoolnixOK {
		mkvStatus = "MKVToolNix " + styles.StatusError.Render("[MISSING]")
	}

	middle := fmt.Sprintf("DEPS: %s %s", ffmpegStatus, mkvStatus)

	// Right side: Quit
	rightSide := styles.RenderHotkey("q", "QUIT")

	// Join with separators
	footer := leftSide + "  │  " + middle + "  │  " + rightSide

	return styles.Footer.Render(footer)
}
