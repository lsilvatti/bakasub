package tape

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TranslationPair represents an original line and its translation
type TranslationPair struct {
	ID           int
	OriginalText string
	Translated   string
}

// Model represents the cassette tape visualization component
type Model struct {
	viewport   viewport.Model
	pairs      []TranslationPair
	width      int
	height     int
	progress   float64 // 0-100 for the progress bar
	spoolFrame int     // Animation frame for rotating spools
	autoScroll bool
	maxPairs   int // Maximum number of pairs to keep in memory

	// Styles
	neonPink lipgloss.Color
	cyan     lipgloss.Color
	yellow   lipgloss.Color
	dimmed   lipgloss.Color
}

// NewModel creates a new cassette tape model
func NewModel(width, height int) Model {
	vp := viewport.New(width-4, 5) // Fixed 5-line viewport
	vp.YPosition = 0

	return Model{
		viewport:   vp,
		pairs:      make([]TranslationPair, 0),
		width:      width,
		height:     height,
		progress:   0.0,
		spoolFrame: 0,
		autoScroll: true,
		maxPairs:   500, // Keep last 500 pairs in memory

		// Native Neon Colors
		neonPink: lipgloss.Color("#F700FF"),
		cyan:     lipgloss.Color("#00FFFF"),
		yellow:   lipgloss.Color("#FFFF00"),
		dimmed:   lipgloss.Color("#666666"),
	}
}

// AddPair adds a new translation pair and updates the viewport
func (m *Model) AddPair(pair TranslationPair) {
	m.pairs = append(m.pairs, pair)

	// Circular buffer: remove oldest if exceeds max
	if len(m.pairs) > m.maxPairs {
		m.pairs = m.pairs[1:]
	}

	// Advance spool animation frame
	m.spoolFrame = (m.spoolFrame + 1) % 4

	// Update viewport content
	m.updateViewportContent()

	// Auto-scroll to bottom if enabled
	if m.autoScroll {
		m.viewport.GotoBottom()
	}
}

// SetProgress updates the progress bar (0-100)
func (m *Model) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	m.progress = progress
}

// updateViewportContent rebuilds the viewport content from pairs
func (m *Model) updateViewportContent() {
	var sb strings.Builder

	// Calculate max width for truncation (accounting for viewport padding)
	maxWidth := m.viewport.Width - 2

	for _, pair := range m.pairs {
		// Format: [ID] Original >>> Translated
		idStr := lipgloss.NewStyle().
			Foreground(m.cyan).
			Bold(true).
			Render(fmt.Sprintf("[%d]", pair.ID))

		arrow := lipgloss.NewStyle().
			Foreground(m.neonPink).
			Bold(true).
			Render(" >>> ")

		// Truncate texts if needed
		originalText := truncate(pair.OriginalText, 30)
		translatedText := truncate(pair.Translated, 30)

		line := fmt.Sprintf("%s %s%s%s", idStr, originalText, arrow, translatedText)

		// Ensure line fits within viewport width
		if lipgloss.Width(line) > maxWidth {
			line = truncate(line, maxWidth)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	m.viewport.SetContent(sb.String())
}

// truncate shortens text with ellipsis if it exceeds maxLen
func truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return "..."
	}
	return text[:maxLen-3] + "..."
}

// Update handles viewport scrolling
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.autoScroll = false
			m.viewport, cmd = m.viewport.Update(msg)

		case "down", "j":
			m.viewport, cmd = m.viewport.Update(msg)
			// Re-enable auto-scroll if at bottom
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}

		case "G": // Jump to bottom (vim-style)
			m.viewport.GotoBottom()
			m.autoScroll = true
		}

	default:
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

// View renders the cassette tape visualization
func (m Model) View() string {
	// Calculate dimensions
	tapeWidth := m.width
	if tapeWidth < 52 {
		tapeWidth = 52 // Minimum width for tape design
	}

	// Determine spool characters based on animation frame
	var leftSpool, rightSpool string
	switch m.spoolFrame % 4 {
	case 0:
		leftSpool = "o"
		rightSpool = "o"
	case 1:
		leftSpool = "+"
		rightSpool = "+"
	case 2:
		leftSpool = "x"
		rightSpool = "x"
	case 3:
		leftSpool = "+"
		rightSpool = "+"
	}

	// Style for dimmed elements
	dimmedStyle := lipgloss.NewStyle().Foreground(m.dimmed)

	// === TOP BORDER ===
	topBorder := "┌" + strings.Repeat("─", tapeWidth-2) + "┐"

	// === SPOOL LINE ===
	spoolSpacing := tapeWidth - 14 // Space between spools
	spoolLine := fmt.Sprintf("│  %s   %s%s%s   %s  │",
		leftSpool,
		leftSpool,
		strings.Repeat(" ", spoolSpacing),
		rightSpool,
		rightSpool,
	)
	spoolLine = dimmedStyle.Render(spoolLine)

	// === VIEWPORT BORDER (TOP) ===
	vpTopBorder := "│ ┌" + strings.Repeat("─", tapeWidth-6) + "┐ │"

	// === VIEWPORT CONTENT (5 lines) ===
	viewportContent := m.viewport.View()
	vpLines := strings.Split(viewportContent, "\n")

	// Ensure exactly 5 lines
	for len(vpLines) < 5 {
		vpLines = append(vpLines, "")
	}
	if len(vpLines) > 5 {
		vpLines = vpLines[len(vpLines)-5:]
	}

	var vpRendered []string
	for _, line := range vpLines {
		// Pad line to fill viewport width
		lineWidth := lipgloss.Width(line)
		padding := m.viewport.Width - lineWidth
		if padding < 0 {
			padding = 0
		}
		paddedLine := line + strings.Repeat(" ", padding)
		vpRendered = append(vpRendered, fmt.Sprintf("│ │ %s │ │", paddedLine))
	}

	// === VIEWPORT BORDER (BOTTOM) ===
	vpBottomBorder := "│ └" + strings.Repeat("─", tapeWidth-6) + "┘ │"

	// === PROGRESS BAR ===
	progressBarWidth := tapeWidth - 6
	filled := int((m.progress / 100.0) * float64(progressBarWidth))
	if filled > progressBarWidth {
		filled = progressBarWidth
	}
	if filled < 0 {
		filled = 0
	}

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", progressBarWidth-filled)

	progressStyle := lipgloss.NewStyle().Foreground(m.yellow)
	progressBar := fmt.Sprintf("│  [%s%s]  │",
		progressStyle.Render(filledBar),
		dimmedStyle.Render(emptyBar),
	)

	// === BOTTOM BORDER ===
	bottomBorder := "└" + strings.Repeat("─", tapeWidth-2) + "┘"

	// === ASSEMBLE ===
	var lines []string
	lines = append(lines, topBorder)
	lines = append(lines, spoolLine)
	lines = append(lines, vpTopBorder)
	lines = append(lines, vpRendered...)
	lines = append(lines, vpBottomBorder)
	lines = append(lines, progressBar)
	lines = append(lines, bottomBorder)

	return strings.Join(lines, "\n")
}

// SetSize updates the component dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Update viewport width (fixed height of 5)
	if width > 4 {
		m.viewport.Width = width - 4
	}
}

// GetPairCount returns the number of translation pairs in memory
func (m Model) GetPairCount() int {
	return len(m.pairs)
}

// Clear removes all translation pairs
func (m *Model) Clear() {
	m.pairs = make([]TranslationPair, 0)
	m.viewport.SetContent("")
	m.progress = 0.0
	m.spoolFrame = 0
}
