// Package picker provides a file/directory browser for BakaSub.
// It implements a full-screen modal for selecting paths with Native Neon aesthetics.
package picker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// Neon color palette
var (
	neonPink = lipgloss.Color("#F700FF")
	cyan     = lipgloss.Color("#00FFFF")
	yellow   = lipgloss.Color("#FFFF00")
	gray     = lipgloss.Color("#808080")
	white    = lipgloss.Color("#FFFFFF")
)

// SelectionMode determines what can be selected
type SelectionMode int

const (
	// ModeDirectory allows selecting directories only
	ModeDirectory SelectionMode = iota
	// ModeFile allows selecting files only
	ModeFile
	// ModeBoth allows selecting both files and directories
	ModeBoth
)

// SelectedPathMsg is sent when a path is selected
type SelectedPathMsg struct {
	Path    string
	IsDir   bool
	Aborted bool
}

// Model represents the file picker state
type Model struct {
	filepicker    filepicker.Model
	selectedPath  string
	selectionMode SelectionMode
	width         int
	height        int
	err           error
	quitting      bool
	title         string
}

// New creates a new file picker model
func New(startDir string, mode SelectionMode) Model {
	// Default to home directory if not specified
	if startDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			startDir = "/"
		} else {
			startDir = home
		}
	}

	// Validate start directory exists
	if _, err := os.Stat(startDir); os.IsNotExist(err) {
		startDir, _ = os.UserHomeDir()
	}

	fp := filepicker.New()
	fp.CurrentDirectory = startDir
	fp.ShowPermissions = false
	fp.ShowSize = true
	fp.ShowHidden = false
	fp.Height = 15 // Will be adjusted on WindowSizeMsg

	// Configure file filtering based on mode
	if mode == ModeDirectory {
		fp.DirAllowed = true
		fp.FileAllowed = false
	} else if mode == ModeFile {
		fp.DirAllowed = false
		fp.FileAllowed = true
		// Filter for media files
		fp.AllowedTypes = []string{".mkv", ".mp4", ".avi", ".srt", ".ass", ".ssa", ".sub"}
	} else {
		fp.DirAllowed = true
		fp.FileAllowed = true
	}

	// Apply custom styles
	fp.Styles.Cursor = lipgloss.NewStyle().Foreground(neonPink).Bold(true)
	fp.Styles.Directory = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	fp.Styles.File = lipgloss.NewStyle().Foreground(white)
	fp.Styles.Selected = lipgloss.NewStyle().Foreground(neonPink).Bold(true)
	fp.Styles.Symlink = lipgloss.NewStyle().Foreground(yellow)

	title := locales.T("picker.title_directory")
	if mode == ModeFile {
		title = locales.T("picker.title_file")
	} else if mode == ModeBoth {
		title = locales.T("picker.title_both")
	}

	return Model{
		filepicker:    fp,
		selectionMode: mode,
		width:         0, // Will be set by WindowSizeMsg
		height:        0, // Will be set by WindowSizeMsg
		title:         title,
	}
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updatePickerHeight()
}

// updatePickerHeight recalculates the filepicker height based on current dimensions
func (m *Model) updatePickerHeight() {
	// Reserve space for:
	// Header section: 6 lines (bar + title + bar + empty + path + mode)
	// Empty line before picker: 1 line
	// Footer: 3-4 lines (empty + help OR empty + filetypes + help)
	// MainWindow borders: 2 lines
	// MainWindow padding: 0 lines (we only have horizontal padding)
	reserved := 6 + 1 + 4 + 2 // = 13 lines reserved
	m.filepicker.Height = m.height - reserved
	if m.filepicker.Height < 5 {
		m.filepicker.Height = 5
	}
}

// Init initializes the file picker
func (m Model) Init() tea.Cmd {
	// Request terminal size and initialize filepicker
	return tea.Batch(tea.WindowSize(), m.filepicker.Init())
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, func() tea.Msg {
				return SelectedPathMsg{Aborted: true}
			}

		case "tab":
			// Toggle selection mode
			switch m.selectionMode {
			case ModeDirectory:
				m.selectionMode = ModeFile
				m.filepicker.DirAllowed = false
				m.filepicker.FileAllowed = true
				m.filepicker.AllowedTypes = []string{".mkv", ".mp4", ".avi", ".srt", ".ass", ".ssa", ".sub"}
				m.title = locales.T("picker.title_file")
			case ModeFile:
				m.selectionMode = ModeBoth
				m.filepicker.DirAllowed = true
				m.filepicker.FileAllowed = true
				m.filepicker.AllowedTypes = nil
				m.title = locales.T("picker.title_both")
			default:
				m.selectionMode = ModeDirectory
				m.filepicker.DirAllowed = true
				m.filepicker.FileAllowed = false
				m.filepicker.AllowedTypes = nil
				m.title = locales.T("picker.title_directory")
			}
			return m, nil

		case " ", "enter":
			// Check if we have a valid selection
			if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
				info, err := os.Stat(path)
				if err == nil {
					return m, func() tea.Msg {
						return SelectedPathMsg{
							Path:  path,
							IsDir: info.IsDir(),
						}
					}
				}
			}
			if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
				// Handle disabled file (wrong type)
				m.err = nil
				_ = path
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePickerHeight()
	}

	// Update the filepicker
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Check for file selection from filepicker
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		info, err := os.Stat(path)
		if err == nil {
			return m, func() tea.Msg {
				return SelectedPathMsg{
					Path:  path,
					IsDir: info.IsDir(),
				}
			}
		}
	}

	return m, cmd
}

// View renders the file picker
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Check if terminal is too small or waiting for size
	if m.width == 0 || m.height == 0 {
		return locales.T("common.loading")
	}

	contentWidth := m.width - 4
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Header bar with title
	headerBarWidth := contentWidth - 4
	if headerBarWidth < 10 {
		headerBarWidth = 10
	}
	headerBar := strings.Repeat("▒", headerBarWidth)
	headerBarStyle := lipgloss.NewStyle().Foreground(gray)
	headerStyle := lipgloss.NewStyle().
		Foreground(neonPink).
		Bold(true)
	titleHeader := headerStyle.Render(m.title)

	// Current path display
	currentPath := m.filepicker.CurrentDirectory
	maxPathLen := contentWidth - 10
	if maxPathLen < 20 {
		maxPathLen = 20
	}
	if len(currentPath) > maxPathLen {
		currentPath = "..." + currentPath[len(currentPath)-(maxPathLen-3):]
	}
	pathStyle := lipgloss.NewStyle().
		Foreground(cyan).
		Bold(true)
	pathDisplay := "📁 " + pathStyle.Render(currentPath)

	// Mode indicator
	modeStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)
	modeText := locales.T("picker.mode_directory")
	if m.selectionMode == ModeFile {
		modeText = locales.T("picker.mode_file")
	} else if m.selectionMode == ModeBoth {
		modeText = locales.T("picker.mode_both")
	}
	modeDisplay := locales.T("picker.mode_label") + " " + modeStyle.Render(modeText)

	// File picker content
	pickerContent := m.filepicker.View()

	// Help text
	help := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderHotkey("↑/↓", locales.T("picker.hotkey_navigate")),
		"  ",
		renderHotkey("Enter", locales.T("picker.hotkey_select")),
		"  ",
		renderHotkey("Tab", locales.T("picker.hotkey_toggle")),
		"  ",
		renderHotkey("Esc", locales.T("picker.hotkey_cancel")),
	)

	// File type hint
	helpStyle := lipgloss.NewStyle().
		Foreground(gray)
	fileTypes := ""
	if m.selectionMode == ModeFile {
		fileTypes = helpStyle.Render(locales.T("picker.allowed_label") + " .mkv .mp4 .avi .srt .ass .ssa .sub")
	}

	// Build header section
	headerSection := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBarStyle.Render(headerBar),
		titleHeader,
		headerBarStyle.Render(headerBar),
		"",
		pathDisplay,
		modeDisplay,
	)

	// Footer with file types and help
	var footer string
	if fileTypes != "" {
		footer = lipgloss.JoinVertical(lipgloss.Left, "", fileTypes, help)
	} else {
		footer = lipgloss.JoinVertical(lipgloss.Left, "", help)
	}

	// Calculate available height for picker content
	// Header: 6 lines (bar + title + bar + empty + path + mode)
	// Footer: 2-3 lines (empty + help OR empty + filetypes + help)
	// Borders: 2 lines
	headerHeight := 6
	footerHeight := 3
	if fileTypes != "" {
		footerHeight = 4
	}
	borderHeight := 2

	availableHeight := m.height - headerHeight - footerHeight - borderHeight
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Style for fixed height picker content area
	pickerStyle := lipgloss.NewStyle().
		Height(availableHeight).
		MaxHeight(availableHeight)

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerSection,
		"",
		pickerStyle.Render(pickerContent),
		footer,
	)

	// Use styles.MainWindow like dashboard does
	// Don't set Height on MainWindow - let content determine height
	// and use lipgloss.Place to position in terminal
	mainContent := styles.MainWindow.
		Width(contentWidth).
		Render(content)

	// Place content at top of terminal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		mainContent,
	)
}

// renderHotkey creates a styled hotkey hint
func renderHotkey(key, text string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(neonPink).
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(white)

	return "[" + keyStyle.Render(key) + "] " + textStyle.Render(text)
}

// SelectedPath returns the currently selected path
func (m Model) SelectedPath() string {
	return m.selectedPath
}

// CurrentDirectory returns the current directory being browsed
func (m Model) CurrentDirectory() string {
	return m.filepicker.CurrentDirectory
}

// ScanDirectory scans a directory for MKV files
func ScanDirectory(path string) ([]string, error) {
	var files []string

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// Single file selected
		return []string{path}, nil
	}

	// Scan directory for MKV files
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".mkv" || ext == ".mp4" || ext == ".avi" {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	return files, nil
}
