package job

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// ConflictModal represents the track selection modal
type ConflictModal struct {
	file          AnalyzedFile
	table         table.Model
	selectedTrack int
	width         int
	height        int
}

// NewConflictModal creates a new conflict resolution modal
func NewConflictModal(file AnalyzedFile) *ConflictModal {
	columns := []table.Column{
		{Title: "#", Width: 5},
		{Title: "TYPE", Width: 8},
		{Title: "SIZE", Width: 10},
		{Title: "TRACK NAME", Width: 40},
	}

	rows := []table.Row{}
	for _, track := range file.ConflictTracks {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", track.ID),
			strings.ToUpper(track.Codec),
			"~45KB",
			track.Name,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(styles.NeonPink).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(styles.NeonCyan).Background(lipgloss.Color("236")).Bold(true)
	t.SetStyles(s)

	return &ConflictModal{file: file, table: t, selectedTrack: 0}
}

// Update handles messages for the conflict modal
func (m ConflictModal) Update(msg tea.Msg) (ConflictModal, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			return m, nil
		case key.Matches(msg, keys.Enter):
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) > 0 {
				trackID := 0
				fmt.Sscanf(selectedRow[0], "%d", &trackID)
				return m, func() tea.Msg { return MsgConflictResolved{FileIndex: 0, TrackID: trackID} }
			}
			return m, nil
		case key.Matches(msg, keys.Up):
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		case key.Matches(msg, keys.Down):
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

// View renders the conflict modal
func (m ConflictModal) View() string {
	var s strings.Builder
	modalStyle := lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(styles.NeonPink).Padding(1, 2).Width(80)
	s.WriteString(styles.TitleStyle.Render("RESOLVE TRACK CONFLICT") + "\n\n")
	s.WriteString("Multiple tracks detected. Select the Full Dialogue track:\n\n")
	s.WriteString(m.table.View() + "\n\n")
	footer := styles.KeyHintStyle.Render("[ESC]") + " CANCEL      " + styles.KeyHintStyle.Render("[ENTER]") + " CONFIRM"
	s.WriteString(footer)
	return modalStyle.Render(s.String())
}

// DryRunReport represents the simulation report
type DryRunReport struct {
	config        JobConfig
	CanWrite      bool
	TokenCount    int
	EstimatedCost float64
	Warnings      []string
	width         int
	height        int
}

// NewDryRunReport creates a new dry run report
func NewDryRunReport(config JobConfig) *DryRunReport {
	return &DryRunReport{config: config, Warnings: []string{}}
}

// Update handles messages for the dry run report
func (m DryRunReport) Update(msg tea.Msg) (DryRunReport, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the dry run report
func (m DryRunReport) View() string {
	var s strings.Builder
	s.WriteString(styles.TitleStyle.Render("DRY RUN REPORT (SIMULATION)") + "\n\n")
	s.WriteString(styles.SectionStyle.Render("JOB SUMMARY") + "\n")
	s.WriteString(fmt.Sprintf("  FILES:       %d Episodes\n", len(m.config.Files)))
	s.WriteString(fmt.Sprintf("  PROVIDER:    %s\n\n", m.config.AIModel))
	s.WriteString(styles.SectionStyle.Render("COST ANALYSIS") + "\n")
	s.WriteString(fmt.Sprintf("  EST. TOKENS:       %dk\n", m.TokenCount/1000))
	s.WriteString(fmt.Sprintf("  ESTIMATED TOTAL:   $%.2f USD\n\n", m.EstimatedCost))
	s.WriteString(styles.SectionStyle.Render("PRE-FLIGHT CHECKS") + "\n")
	if m.CanWrite {
		s.WriteString(styles.SuccessStyle.Render("  [OK]") + " Write Permissions\n")
	} else {
		s.WriteString(styles.ErrorStyle.Render("  [!!]") + " No Write Permissions\n")
	}
	if len(m.Warnings) > 0 {
		for _, warning := range m.Warnings {
			s.WriteString(styles.WarningStyle.Render(fmt.Sprintf("  [!!] %s\n", warning)))
		}
	}
	s.WriteString("\n")
	footer := styles.KeyHintStyle.Render("[ESC]") + " BACK TO SETUP"
	s.WriteString(footer)
	return styles.AppStyle.Render(s.String())
}

// GlossaryEditor represents the glossary editing interface
type GlossaryEditor struct {
	path     string
	Terms    map[string]string
	table    table.Model
	addMode  bool
	inputs   []textinput.Model
	Closed   bool
	Modified bool
	width    int
	height   int
}

// GlossaryEntry represents a term in the glossary
type GlossaryEntry struct {
	Original    string `json:"original"`
	Translation string `json:"translation"`
	Type        string `json:"type"`
}

// NewGlossaryEditor creates a new glossary editor
func NewGlossaryEditor(path string, terms map[string]string) *GlossaryEditor {
	columns := []table.Column{
		{Title: "ORIGINAL", Width: 30},
		{Title: "TRANSLATION", Width: 30},
		{Title: "TYPE", Width: 15},
	}
	rows := []table.Row{}
	for original, translation := range terms {
		rows = append(rows, table.Row{original, translation, "[Noun]"})
	}
	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(true))
	inputs := make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Original term"
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Translation"
	return &GlossaryEditor{path: path, Terms: terms, table: t, inputs: inputs}
}

// Update handles messages for the glossary editor
func (m GlossaryEditor) Update(msg tea.Msg) (GlossaryEditor, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Escape) {
			if m.addMode {
				m.addMode = false
				return m, nil
			}
			m.Closed = true
			return m, nil
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the glossary editor
func (m GlossaryEditor) View() string {
	var s strings.Builder
	s.WriteString(styles.TitleStyle.Render("PROJECT GLOSSARY EDITOR") + "\n\n")
	s.WriteString(fmt.Sprintf("FILE: %s\n\n", m.path))
	s.WriteString(m.table.View() + "\n\n")
	s.WriteString(styles.KeyHintStyle.Render("[ESC]") + " CLOSE")
	return styles.AppStyle.Render(s.String())
}

// LoadGlossaryFromFile loads a glossary from a JSON file
func LoadGlossaryFromFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []GlossaryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	terms := make(map[string]string)
	for _, entry := range entries {
		terms[entry.Original] = entry.Translation
	}
	return terms, nil
}

// SaveGlossaryToFile saves a glossary to a JSON file
func SaveGlossaryToFile(path string, terms map[string]string) error {
	entries := []GlossaryEntry{}
	for original, translation := range terms {
		entries = append(entries, GlossaryEntry{Original: original, Translation: translation, Type: "noun"})
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
