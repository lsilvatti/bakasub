package execution

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/ui/components/tape"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogWarn
	LogError
	LogAI
	LogSuccess
)

// LogEntry represents a single log line with metadata
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

// LogBuffer is a circular buffer for log entries
type LogBuffer struct {
	entries []LogEntry
	maxSize int
	mu      sync.RWMutex
}

// NewLogBuffer creates a new circular log buffer
func NewLogBuffer(maxSize int) *LogBuffer {
	return &LogBuffer{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// AddLine adds a new log entry to the buffer
func (lb *LogBuffer) AddLine(level LogLevel, message string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	if len(lb.entries) >= lb.maxSize {
		// Remove oldest entry (circular buffer)
		lb.entries = lb.entries[1:]
	}

	lb.entries = append(lb.entries, entry)
}

// GetLines returns all log entries as formatted strings with syntax highlighting
func (lb *LogBuffer) GetLines() []string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	lines := make([]string, 0, len(lb.entries))
	for _, entry := range lb.entries {
		lines = append(lines, FormatLogEntry(entry))
	}
	return lines
}

// GetRawText returns all log entries as plain text (for viewport)
func (lb *LogBuffer) GetRawText() string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	var sb strings.Builder
	for _, entry := range lb.entries {
		sb.WriteString(FormatLogEntry(entry))
		sb.WriteString("\n")
	}
	return sb.String()
}

// Count returns the number of entries in the buffer
func (lb *LogBuffer) Count() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.entries)
}

// Clear removes all entries from the buffer
func (lb *LogBuffer) Clear() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.entries = make([]LogEntry, 0, lb.maxSize)
}

// FormatLogEntry formats a log entry with syntax highlighting
func FormatLogEntry(entry LogEntry) string {
	// Format timestamp [HH:MM:SS]
	timeStr := entry.Timestamp.Format("15:04:05")
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")).Faint(true)

	// Format level tag with appropriate color
	var levelStr string
	var levelStyle lipgloss.Style

	switch entry.Level {
	case LogInfo:
		levelStr = "[INFO]"
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")) // Cyan
	case LogWarn:
		levelStr = "[WARN]"
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")) // Yellow
	case LogError:
		levelStr = "[ERR]"
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true) // Red + Bold
	case LogAI:
		levelStr = "[AI]"
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F700FF")) // Neon Pink
	case LogSuccess:
		levelStr = "[OK]"
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // Green
	}

	// Format message (plain white)
	messageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	return fmt.Sprintf("%s %s %s",
		timeStyle.Render("["+timeStr+"]"),
		levelStyle.Render(levelStr),
		messageStyle.Render(entry.Message),
	)
}

// ParseLogLine attempts to detect log level from a plain text line
func ParseLogLine(line string) (LogLevel, string) {
	line = strings.TrimSpace(line)

	// Check for level indicators
	if strings.Contains(line, "[INFO]") || strings.Contains(line, "INFO:") {
		return LogInfo, strings.TrimPrefix(strings.TrimPrefix(line, "[INFO]"), "INFO:")
	}
	if strings.Contains(line, "[WARN]") || strings.Contains(line, "WARN:") || strings.Contains(line, "[WARNING]") {
		return LogWarn, strings.TrimPrefix(strings.TrimPrefix(line, "[WARN]"), "WARN:")
	}
	if strings.Contains(line, "[ERR]") || strings.Contains(line, "ERROR:") || strings.Contains(line, "[FAIL]") {
		return LogError, strings.TrimPrefix(strings.TrimPrefix(line, "[ERR]"), "ERROR:")
	}
	if strings.Contains(line, "[AI]") || strings.Contains(line, "AI:") {
		return LogAI, strings.TrimPrefix(strings.TrimPrefix(line, "[AI]"), "AI:")
	}
	if strings.Contains(line, "[OK]") || strings.Contains(line, "[SUCCESS]") {
		return LogSuccess, strings.TrimPrefix(strings.TrimPrefix(line, "[OK]"), "SUCCESS:")
	}

	// Default to info if no level detected
	return LogInfo, line
}

// KeyMap defines keyboard shortcuts for the execution screen
type KeyMap struct {
	Quit       key.Binding
	Pause      key.Binding
	Cancel     key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	Top        key.Binding
	Bottom     key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit:       key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit (when done)")),
		Pause:      key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "pause/resume")),
		Cancel:     key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("esc", "cancel job")),
		ScrollUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
		ScrollDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
		PageUp:     key.NewBinding(key.WithKeys("pgup"), key.WithHelp("pgup", "page up")),
		PageDown:   key.NewBinding(key.WithKeys("pgdown"), key.WithHelp("pgdn", "page down")),
		Top:        key.NewBinding(key.WithKeys("g", "home"), key.WithHelp("g", "jump to top")),
		Bottom:     key.NewBinding(key.WithKeys("G", "end"), key.WithHelp("G", "jump to bottom")),
	}
}

var keys = DefaultKeyMap()

// JobStatus represents the current state of the job
type JobStatus int

const (
	StatusRunning JobStatus = iota
	StatusPaused
	StatusCancelling
	StatusComplete
	StatusFailed
)

// Model represents the execution screen state
type Model struct {
	width  int
	height int

	// Job info
	jobName       string
	currentFile   string
	fileIndex     int
	totalFiles    int
	fileProgress  float64
	batchProgress float64
	status        JobStatus
	startTime     time.Time
	elapsedTime   time.Duration

	// Statistics
	linesProcessed int
	tokensUsed     int
	costSoFar      float64
	errors         int

	// Logging
	logBuffer    *LogBuffer
	viewport     viewport.Model
	autoScroll   bool
	lastLogCount int

	// Live Translation View (Cassette Tape)
	tapeView tape.Model

	// Control
	quitting bool
}

// New creates a new execution screen model
func New(jobName string, totalFiles int) Model {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFFF")).
		Padding(0, 1)

	tapeView := tape.NewModel(80, 12)

	return Model{
		width:        80,
		height:       24,
		jobName:      jobName,
		totalFiles:   totalFiles,
		fileIndex:    1,
		status:       StatusRunning,
		startTime:    time.Now(),
		logBuffer:    NewLogBuffer(1000), // Circular buffer: 1000 lines max
		viewport:     vp,
		autoScroll:   true,
		lastLogCount: 0,
		tapeView:     tapeView,
	}
}

// Init initializes the execution screen
func (m Model) Init() tea.Cmd {
	return nil
}

// LogMsg is sent when a new log line is added
type LogMsg struct {
	Level   LogLevel
	Message string
}

// ProgressMsg updates job progress
type ProgressMsg struct {
	FileProgress  float64
	BatchProgress float64
	CurrentFile   string
}

// StatsMsg updates statistics
type StatsMsg struct {
	LinesProcessed int
	TokensUsed     int
	CostSoFar      float64
	Errors         int
}

// StatusMsg updates job status
type StatusMsg struct {
	Status JobStatus
}

// TranslationMsg adds a new translation pair to the tape view
type TranslationMsg struct {
	ID           int
	OriginalText string
	Translated   string
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.status == StatusComplete || m.status == StatusFailed {
			// Allow quitting when done
			if key.Matches(msg, keys.Quit) {
				m.quitting = true
				return m, tea.Quit
			}
		}

		switch {
		case key.Matches(msg, keys.Cancel):
			if m.status == StatusRunning || m.status == StatusPaused {
				m.status = StatusCancelling
				m.logBuffer.AddLine(LogWarn, "Cancelling job...")
			}
		case key.Matches(msg, keys.Pause):
			if m.status == StatusRunning {
				m.status = StatusPaused
				m.logBuffer.AddLine(LogInfo, "Job paused by user")
			} else if m.status == StatusPaused {
				m.status = StatusRunning
				m.logBuffer.AddLine(LogInfo, "Job resumed")
			}
		case key.Matches(msg, keys.ScrollUp):
			m.viewport.LineUp(1)
			m.tapeView, _ = m.tapeView.Update(msg)
			m.autoScroll = false
		case key.Matches(msg, keys.ScrollDown):
			m.viewport.LineDown(1)
			m.tapeView, _ = m.tapeView.Update(msg)
			// Re-enable auto-scroll if at bottom
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}
		case key.Matches(msg, keys.PageUp):
			m.viewport.ViewUp()
			m.autoScroll = false
		case key.Matches(msg, keys.PageDown):
			m.viewport.ViewDown()
			if m.viewport.AtBottom() {
				m.autoScroll = true
			}
		case key.Matches(msg, keys.Top):
			m.viewport.GotoTop()
			m.autoScroll = false
		case key.Matches(msg, keys.Bottom):
			m.viewport.GotoBottom()
			m.tapeView, _ = m.tapeView.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			m.autoScroll = true
		}

	case LogMsg:
		m.logBuffer.AddLine(msg.Level, msg.Message)
		// Update viewport content
		m.viewport.SetContent(m.logBuffer.GetRawText())
		// Auto-scroll to bottom if enabled
		if m.autoScroll {
			m.viewport.GotoBottom()
		}

	case ProgressMsg:
		m.fileProgress = msg.FileProgress
		m.batchProgress = msg.BatchProgress
		m.currentFile = msg.CurrentFile

	case StatsMsg:
		m.linesProcessed = msg.LinesProcessed
		m.tokensUsed = msg.TokensUsed
		m.costSoFar = msg.CostSoFar
		m.errors = msg.Errors

	case StatusMsg:
		m.status = msg.Status

	case TranslationMsg:
		// Add translation pair to tape view
		m.tapeView.AddPair(tape.TranslationPair{
			ID:           msg.ID,
			OriginalText: msg.OriginalText,
			Translated:   msg.Translated,
		})
		// Update tape progress based on file progress
		m.tapeView.SetProgress(m.fileProgress)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Recalculate viewport dimensions
		headerHeight := 12 // Approximate height for header + progress
		footerHeight := 3
		availableHeight := m.height - headerHeight - footerHeight
		if availableHeight < 5 {
			availableHeight = 5
		}

		vpWidth := layout.SafeWidth(m.width-6, 70)
		m.viewport.Width = vpWidth
		m.viewport.Height = availableHeight

		// Resize tape view
		tapeWidth := layout.SafeWidth(m.width-4, 52)
		m.tapeView.SetSize(tapeWidth, 12)

		// Refresh content
		m.viewport.SetContent(m.logBuffer.GetRawText())
		if m.autoScroll {
			m.viewport.GotoBottom()
		}
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// SetSize updates the model dimensions and resizes components
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Resize viewport (logs)
	logHeight := 8
	if height > 24 {
		logHeight = height - 16
	}
	m.viewport.Width = width - 4
	m.viewport.Height = logHeight

	// Resize tape view
	m.tapeView.SetSize(width-4, 12)
}

// View renders the execution screen
func (m Model) View() string {
	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	header := m.renderHeader()
	progress := m.renderProgress()
	tapeView := m.renderTape()
	logs := m.renderLogs()
	footer := m.renderFooter()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		progress,
		"",
		tapeView,
		"",
		logs,
		"",
		footer,
	)

	contentWidth := layout.SafeWidth(m.width-4, 76)
	return styles.MainWindow.Width(contentWidth).Render(content)
}

func (m Model) renderHeader() string {
	// Status indicator
	var statusStr string
	var statusColor lipgloss.Color

	switch m.status {
	case StatusRunning:
		statusStr = "[▶ RUNNING]"
		statusColor = lipgloss.Color("#00FF00")
	case StatusPaused:
		statusStr = "[⏸ PAUSED]"
		statusColor = lipgloss.Color("#FFFF00")
	case StatusCancelling:
		statusStr = "[⏹ CANCELLING...]"
		statusColor = lipgloss.Color("#FF0000")
	case StatusComplete:
		statusStr = "[✓ COMPLETE]"
		statusColor = lipgloss.Color("#00FF00")
	case StatusFailed:
		statusStr = "[✗ FAILED]"
		statusColor = lipgloss.Color("#FF0000")
	}

	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)

	// Calculate elapsed time
	if m.status == StatusRunning || m.status == StatusPaused {
		m.elapsedTime = time.Since(m.startTime)
	}

	elapsed := fmt.Sprintf("%02d:%02d:%02d",
		int(m.elapsedTime.Hours()),
		int(m.elapsedTime.Minutes())%60,
		int(m.elapsedTime.Seconds())%60,
	)

	title := styles.TitleStyle.Render(fmt.Sprintf("JOB RUNNING: %s", m.jobName))
	status := statusStyle.Render(statusStr) + " [ETA: " + elapsed + "]"

	headerBar := strings.Repeat("▒", layout.SafeWidth(m.width-6, 70))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		title,
		status,
		headerBar,
	)
}

func (m Model) renderProgress() string {
	// File progress
	fileBar := layout.ProgressBar(int(m.fileProgress), 100, 50)
	fileInfo := fmt.Sprintf("FILE: %s (%d/%d)", m.currentFile, m.fileIndex, m.totalFiles)
	fileProgress := lipgloss.JoinVertical(
		lipgloss.Left,
		fileInfo,
		fileBar,
		fmt.Sprintf("%.0f%%", m.fileProgress),
	)

	// Batch progress
	batchBar := layout.ProgressBar(int(m.batchProgress), 100, 50)
	batchProgress := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("BATCH: %d/%d Files", m.fileIndex, m.totalFiles),
		batchBar,
		fmt.Sprintf("%.0f%%", m.batchProgress),
	)

	// Statistics
	stats := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("LINES: %d", m.linesProcessed),
		fmt.Sprintf("TOKENS: %s", formatNumber(m.tokensUsed)),
		fmt.Sprintf("COST: $%.4f", m.costSoFar),
		fmt.Sprintf("ERRORS: %d", m.errors),
	)

	panelWidth := layout.CalculateThird(m.width, 6)
	panelWidth = layout.SafeWidth(panelWidth, 20)

	filePanel := styles.Panel.Width(panelWidth).Render(fileProgress)
	batchPanel := styles.Panel.Width(panelWidth).Render(batchProgress)
	statsPanel := styles.Panel.Width(panelWidth).Render(stats)

	return lipgloss.JoinHorizontal(lipgloss.Top, filePanel, batchPanel, statsPanel)
}

func (m Model) renderTape() string {
	if m.tapeView.GetPairCount() == 0 {
		// Show empty state with instructions
		emptyMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Render("⏳ Waiting for translation pairs...")
		return lipgloss.JoinVertical(
			lipgloss.Left,
			styles.PanelTitle.Render("🎞️  LIVE TRANSLATION VIEW (CASSETTE TAPE)"),
			"",
			emptyMsg,
		)
	}

	title := styles.PanelTitle.Render(
		fmt.Sprintf("🎞️  LIVE TRANSLATION VIEW (%d pairs)", m.tapeView.GetPairCount()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.tapeView.View(),
	)
}

func (m Model) renderLogs() string {
	title := styles.PanelTitle.Render(
		fmt.Sprintf("EXECUTION LOGS (%d/%d lines)", m.logBuffer.Count(), 1000),
	)

	scrollIndicator := ""
	if !m.autoScroll {
		scrollIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Render(" [SCROLL MODE - Press G to jump to bottom]")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title+scrollIndicator,
		"",
		m.viewport.View(),
	)
}

func (m Model) renderFooter() string {
	if m.status == StatusComplete || m.status == StatusFailed {
		return styles.RenderHotkey("q", "QUIT")
	}

	var controls []string
	if m.status == StatusRunning {
		controls = append(controls, styles.RenderHotkey("p", "PAUSE"))
	} else if m.status == StatusPaused {
		controls = append(controls, styles.RenderHotkey("p", "RESUME"))
	}
	controls = append(controls, styles.RenderHotkey("esc", "CANCEL"))
	controls = append(controls, styles.RenderHotkey("↑↓", "SCROLL"))
	controls = append(controls, styles.RenderHotkey("G", "BOTTOM"))

	return strings.Join(controls, "  │  ")
}

// formatNumber formats large numbers with commas
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000)
}
