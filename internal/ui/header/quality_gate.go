package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/core/linter"
	"github.com/lsilvatti/bakasub/internal/locales"
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
	b.WriteString(styles.TitleStyle.Render(locales.T("quality_gate.title")))
	b.WriteString("\n\n")

	// Description
	b.WriteString(styles.SubtleStyle.Render(locales.T("quality_gate.description")))
	b.WriteString("\n")
	b.WriteString(styles.SubtleStyle.Render(locales.T("quality_gate.review_before_mux")))
	b.WriteString("\n\n")

	// Issues table
	b.WriteString(renderIssuesTable(m.result.Issues))
	b.WriteString("\n\n")

	// Action selection with Panel
	var actionContent strings.Builder
	actionContent.WriteString(styles.SectionStyle.Render("ACTION (Select with Arrows)") + "\n\n")

	actions := []struct {
		name        string
		description string
	}{
		{locales.T("quality_gate.auto_fix"), ""},
		{locales.T("quality_gate.manual_review"), ""},
		{locales.T("quality_gate.ignore_continue"), ""},
	}

	for i, action := range actions {
		prefix := "( ) "
		if i == m.selectedIdx {
			prefix = "(o) "
		}
		line := "   " + prefix + action.name
		if i == m.selectedIdx {
			line = styles.HighlightStyle.Render(line)
		}
		actionContent.WriteString(line + "\n")
	}
	b.WriteString(styles.Panel.Render(actionContent.String()))
	b.WriteString("\n\n")

	// Footer
	footer := styles.KeyHintStyle.Render("[↑/↓]") + " Select  "
	footer += styles.KeyHintStyle.Render("[ENTER]") + " Execute"
	b.WriteString(footer)

	return b.String()
}

func renderIssuesTable(issues []linter.Issue) string {
	var issuesContent strings.Builder
	issuesContent.WriteString(styles.SectionStyle.Render("DETECTED ISSUES") + "\n\n")

	// Header
	issuesContent.WriteString(fmt.Sprintf("   %-6s  %-10s  %-20s  %-30s\n", "ID", "SEVERITY", "ISSUE TYPE", "CONTENT"))
	issuesContent.WriteString("   " + strings.Repeat("─", 70) + "\n")

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

		issuesContent.WriteString(fmt.Sprintf("   %-6d  ", issue.LineID))
		issuesContent.WriteString(severityColor.Render(fmt.Sprintf("[%-8s]", issue.Severity)))
		issuesContent.WriteString(fmt.Sprintf("  %-20s  %-28s\n", issue.IssueType, truncateForTable(issue.Content, 28)))
	}

	if len(issues) > 10 {
		issuesContent.WriteString(fmt.Sprintf("   ... and %d more issues\n", len(issues)-10))
	}

	return styles.Panel.Render(issuesContent.String())
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
