package review

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/core/parser"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/focus"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// ClosedMsg is sent when the review editor should be closed
type ClosedMsg struct{}

type Model struct {
	originalLines   []parser.SubtitleLine
	translatedLines []parser.SubtitleLine
	currentIndex    int
	editor          textarea.Model
	focusManager    *focus.Manager
	filePath        string
	saved           bool
	width           int
	height          int
}

func New(originalPath, translatedPath string) (*Model, error) {
	origFile, err := parser.ParseFile(originalPath)
	if err != nil {
		return nil, err
	}

	// If translatedPath is empty, use originalPath (single file editing mode)
	var transFile *parser.SubtitleFile
	if translatedPath == "" {
		// Clone the original file for editing
		transFile = origFile
		translatedPath = originalPath
	} else {
		transFile, err = parser.ParseFile(translatedPath)
		if err != nil {
			return nil, err
		}
	}

	ta := textarea.New()
	ta.Focus()
	ta.SetHeight(5)

	m := &Model{
		originalLines:   origFile.Lines,
		translatedLines: transFile.Lines,
		currentIndex:    0,
		editor:          ta,
		focusManager:    focus.NewManager(1), // 1 text area field
		filePath:        translatedPath,
		saved:           true,
	}

	// Start in input mode since the editor should be active by default
	m.focusManager.EnterInput(0)

	if len(m.translatedLines) > 0 {
		m.editor.SetValue(m.translatedLines[0].Text)
	}

	return m, nil
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.editor.SetWidth(width/2 - 4)
}

func (m Model) Init() tea.Cmd {
	// Request terminal size and start editor blink
	return tea.Batch(tea.WindowSize(), textarea.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.editor.SetWidth(msg.Width/2 - 4)

	case tea.KeyMsg:
		// GATEKEEPER LOGIC for Review screen
		// The review screen is special - it starts in input mode by default
		// since the main purpose is editing

		if m.focusManager.Mode() == focus.ModeInput {
			// In input mode - handle special keys that should work in edit mode
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
				// Save while staying in edit mode
				return m, m.saveFile()

			case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
				// Tab moves to next line (saves current, loads next)
				m.saved = false
				m.nextLine()
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
				// Shift+Tab moves to previous line
				m.saved = false
				m.prevLine()
				return m, nil

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				// ESC exits the editor completely
				return m, func() tea.Msg { return ClosedMsg{} }
			}

			// All other keys go to the editor
			m.saved = false
		}
	}

	m.editor, cmd = m.editor.Update(msg)
	if m.currentIndex < len(m.translatedLines) {
		m.translatedLines[m.currentIndex].Text = m.editor.Value()
	}

	return m, cmd
}

func (m *Model) nextLine() {
	if m.currentIndex < len(m.translatedLines)-1 {
		m.currentIndex++
		m.editor.SetValue(m.translatedLines[m.currentIndex].Text)
	}
}

func (m *Model) prevLine() {
	if m.currentIndex > 0 {
		m.currentIndex--
		m.editor.SetValue(m.translatedLines[m.currentIndex].Text)
	}
}

func (m Model) saveFile() tea.Cmd {
	return func() tea.Msg {
		content := parser.ReassembleSRT(m.translatedLines)
		if err := os.WriteFile(m.filePath, []byte(content), 0644); err != nil {
			return err
		}
		return nil
	}
}

func (m Model) View() string {
	// Wait for terminal size
	if layout.IsWaitingForSize(m.width, m.height) {
		return locales.T("common.loading")
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	if m.width == 0 {
		return locales.T("common.loading")
	}

	leftWidth := layout.CalculateHalf(m.width, 2)
	leftWidth = layout.SafeWidth(leftWidth, 30)
	rightWidth := leftWidth

	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height - 10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Cyan).
		Padding(1)

	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height - 10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.NeonPink).
		Padding(1)

	var originalText string
	if m.currentIndex < len(m.originalLines) {
		orig := m.originalLines[m.currentIndex]
		originalText = fmt.Sprintf("[%d] %s\n\n%s",
			m.currentIndex+1,
			orig.StartTime,
			orig.Text)
	}

	leftPanel := leftStyle.Render(
		styles.TitleStyle.Render(locales.T("review.original_readonly")) + "\n\n" +
			originalText,
	)

	rightPanel := rightStyle.Render(
		styles.TitleStyle.Render(locales.T("review.translated_editable")) + "\n\n" +
			m.editor.View(),
	)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	footer := styles.KeyHintStyle.Render("[Tab/Shift+Tab]") + " " + locales.T("review.navigate") + "  " +
		styles.KeyHintStyle.Render("[Ctrl+S]") + " " + locales.T("review.save") + "  " +
		styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("review.exit")

	if !m.saved {
		footer += " " + styles.WarningStyle.Render(locales.T("review.modified"))
	}

	progress := locales.Tf("review.progress_simple", m.currentIndex+1, len(m.translatedLines))

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.TitleStyle.Render(locales.T("review.editor_title")),
		progress,
		"",
		content,
		"",
		footer,
	)
}
