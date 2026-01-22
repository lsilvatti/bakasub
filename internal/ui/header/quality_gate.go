package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/core/linter"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type ActionType int

const (
	ActionAutoFix ActionType = iota
	ActionManualReview
	ActionIgnore
)

type QualityGateModel struct {
	result       linter.Result
	selectedIdx  int
	actionChoice ActionType
	width        int
	height       int
}

func NewQualityGate(result linter.Result) *QualityGateModel {
	return &QualityGateModel{
		result:       result,
		selectedIdx:  0,
		actionChoice: ActionAutoFix,
	}
}

func (m QualityGateModel) Init() tea.Cmd {
	return nil
}

func (m QualityGateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "j":
			if m.selectedIdx < 2 {
				m.selectedIdx++
			}

		case "enter":
			// Execute selected action
			return m, tea.Quit
		}
	}

	// Update action choice based on selection
	m.actionChoice = ActionType(m.selectedIdx)

	return m, nil
}

func (m QualityGateModel) View() string {
	var b strings.Builder

	// Header
	b.WriteString(styles.TitleStyle.Render("QUALITY GATE: ISSUES FOUND"))
	b.WriteString("\n\n")

	// Description
	b.WriteString(styles.SubtleStyle.Render("The automatic linter found potential issues in the translation."))
	b.WriteString("\n")
	b.WriteString(styles.SubtleStyle.Render("Please review before muxing."))
	b.WriteString("\n\n")

	// Issues table
	b.WriteString(renderIssuesTable(m.result.Issues))
	b.WriteString("\n\n")

	// Action selection
	b.WriteString(styles.SubtleStyle.Render("┌── ACTION (Select with Arrows) ───────────────────────────────────────┐"))
	b.WriteString("\n")

	actions := []struct {
		name        string
		description string
	}{
		{"AUTO-FIX (Attempt to fix tags/terms via Regex)", ""},
		{"MANUAL REVIEW (Open Editor)", ""},
		{"IGNORE AND CONTINUE", ""},
	}

	for i, action := range actions {
		prefix := "( ) "
		if i == m.selectedIdx {
			prefix = "(o) "
		}
		line := styles.SubtleStyle.Render("│   ") + prefix + action.name
		if i == m.selectedIdx {
			line = styles.HighlightStyle.Render(line)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString(styles.SubtleStyle.Render("└──────────────────────────────────────────────────────────────────────┘"))
	b.WriteString("\n\n")

	// Footer
	footer := styles.KeyHintStyle.Render("[↑/↓]") + " Select  "
	footer += styles.KeyHintStyle.Render("[ENTER]") + " Execute"
	b.WriteString(footer)

	return b.String()
}

func renderIssuesTable(issues []linter.Issue) string {
	var b strings.Builder

	b.WriteString(styles.SubtleStyle.Render("┌── DETECTED ISSUES ───────────────────────────────────────────────────┐"))
	b.WriteString("\n")

	// Header
	header := fmt.Sprintf("│   %-6s  %-10s  %-20s  %-30s │", "ID", "SEVERITY", "ISSUE TYPE", "CONTENT")
	b.WriteString(styles.SubtleStyle.Render(header))
	b.WriteString("\n")

	separator := "│  " + strings.Repeat("─", 70) + " │"
	b.WriteString(styles.SubtleStyle.Render(separator))
	b.WriteString("\n")

	// Rows (limit to 10 for display)
	displayCount := len(issues)
	if displayCount > 10 {
		displayCount = 10
	}

	for i := 0; i < displayCount; i++ {
		issue := issues[i]
		severityColor := styles.SubtleStyle
		switch issue.Severity {
		case linter.SeverityHigh:
			severityColor = styles.ErrorStyle
		case linter.SeverityMedium:
			severityColor = styles.WarningStyle
		case linter.SeverityLow:
			severityColor = styles.SubtleStyle
		}

		row := fmt.Sprintf("│   %-6d  ", issue.LineID)
		b.WriteString(styles.SubtleStyle.Render(row))
		b.WriteString(severityColor.Render(fmt.Sprintf("[%-8s]", issue.Severity)))
		b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("  %-20s  %-28s │", issue.IssueType, truncateForTable(issue.Content, 28))))
		b.WriteString("\n")
	}

	if len(issues) > 10 {
		more := fmt.Sprintf("│   ... and %d more issues", len(issues)-10)
		b.WriteString(styles.SubtleStyle.Render(more))
		b.WriteString(strings.Repeat(" ", 70-len(more)+1))
		b.WriteString(styles.SubtleStyle.Render("│"))
		b.WriteString("\n")
	}

	b.WriteString(styles.SubtleStyle.Render("└──────────────────────────────────────────────────────────────────────┘"))

	return b.String()
}

func truncateForTable(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetSelectedAction returns the currently selected action
func (m QualityGateModel) GetSelectedAction() ActionType {
	return m.actionChoice
}
