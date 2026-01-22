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
	fp.Height = 20

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

	title := "SELECT DIRECTORY"
	if mode == ModeFile {
		title = "SELECT FILE"
	} else if mode == ModeBoth {
		title = "SELECT FILE OR DIRECTORY"
	}

	return Model{
		filepicker:    fp,
		selectionMode: mode,
		width:         80,
		height:        24,
		title:         title,
	}
}

// Init initializes the file picker
func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
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
				m.title = "SELECT FILE"
			case ModeFile:
				m.selectionMode = ModeBoth
				m.filepicker.DirAllowed = true
				m.filepicker.FileAllowed = true
				m.filepicker.AllowedTypes = nil
				m.title = "SELECT FILE OR DIRECTORY"
			default:
				m.selectionMode = ModeDirectory
				m.filepicker.DirAllowed = true
				m.filepicker.FileAllowed = false
				m.filepicker.AllowedTypes = nil
				m.title = "SELECT DIRECTORY"
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
		m.filepicker.Height = m.height - 12
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

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(neonPink).
		Bold(true).
		Padding(0, 1)

	header := headerStyle.Render("╔══ " + m.title + " ══╗")

	// Current path display
	currentPath := m.filepicker.CurrentDirectory
	if len(currentPath) > m.width-20 {
		// Truncate long paths
		currentPath = "..." + currentPath[len(currentPath)-(m.width-23):]
	}
	pathStyle := lipgloss.NewStyle().
		Foreground(cyan).
		Bold(true)
	pathDisplay := "📁 " + pathStyle.Render(currentPath)

	// Mode indicator
	modeStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)
	modeText := "DIRECTORY"
	if m.selectionMode == ModeFile {
		modeText = "FILE"
	} else if m.selectionMode == ModeBoth {
		modeText = "FILE/DIR"
	}
	modeDisplay := "Mode: " + modeStyle.Render(modeText)

	// File picker content with border
	pickerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(cyan).
		Padding(0, 1).
		Width(m.width - 6)

	pickerContent := pickerStyle.Render(m.filepicker.View())

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(gray)

	help := lipgloss.JoinHorizontal(
		lipgloss.Left,
		renderHotkey("↑/↓", "Navigate"),
		"  ",
		renderHotkey("Enter", "Select/Open"),
		"  ",
		renderHotkey("Tab", "Toggle Mode"),
		"  ",
		renderHotkey("Esc", "Cancel"),
	)

	// File type hint
	fileTypes := ""
	if m.selectionMode == ModeFile {
		fileTypes = helpStyle.Render("Allowed: .mkv .mp4 .avi .srt .ass .ssa .sub")
	}

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		pathDisplay,
		modeDisplay,
		"",
		pickerContent,
		"",
		fileTypes,
		help,
	)

	// Wrap in main border
	mainStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(neonPink).
		Padding(1, 2).
		Width(m.width - 2).
		Height(m.height - 2)

	return mainStyle.Render(content)
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
