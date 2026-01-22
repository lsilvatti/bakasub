package glossary

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lsilvatti/bakasub/internal/core/parser"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type Entry struct {
	Original     string `json:"original"`
	Translation  string `json:"translation"`
	Type         string `json:"type"`
	AutoDetected bool   `json:"auto_detected"`
}

type Model struct {
	table    table.Model
	entries  []Entry
	filePath string
	width    int
	height   int
}

func New(glossaryPath string) *Model {
	entries := loadOrCreate(glossaryPath)

	columns := []table.Column{
		{Title: "ORIGINAL", Width: 30},
		{Title: "TRANSLATION", Width: 30},
		{Title: "TYPE", Width: 15},
		{Title: "SOURCE", Width: 10},
	}

	rows := make([]table.Row, len(entries))
	for i, e := range entries {
		source := "Manual"
		if e.AutoDetected {
			source = "Auto"
		}
		rows[i] = table.Row{e.Original, e.Translation, e.Type, source}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return &Model{
		table:    t,
		entries:  entries,
		filePath: glossaryPath,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		case "ctrl+s":
			m.save()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := styles.TitleStyle.Render("PROJECT GLOSSARY") + "\n"
	header += styles.SubtleStyle.Render(fmt.Sprintf("File: %s", m.filePath)) + "\n\n"

	tableView := m.table.View()

	footer := "\n" + styles.KeyHintStyle.Render("[↑/↓]") + " Navigate  "
	footer += styles.KeyHintStyle.Render("[Ctrl+S]") + " Save  "
	footer += styles.KeyHintStyle.Render("[ESC]") + " Exit"

	return header + tableView + footer
}

func (m Model) save() error {
	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

func loadOrCreate(path string) []Entry {
	data, err := os.ReadFile(path)
	if err != nil {
		return []Entry{}
	}

	var entries []Entry
	json.Unmarshal(data, &entries)
	return entries
}

func AutoDetectTerms(subtitlePath string) ([]Entry, error) {
	subFile, err := parser.ParseFile(subtitlePath)
	if err != nil {
		return nil, err
	}

	termMap := make(map[string]bool)
	entries := []Entry{}

	capitalized := regexp.MustCompile(`\b([A-Z][a-z]{2,})\b`)

	for _, line := range subFile.Lines {
		matches := capitalized.FindAllString(line.Text, -1)
		for _, match := range matches {
			if !termMap[match] && !isCommonWord(match) {
				termMap[match] = true
				entries = append(entries, Entry{
					Original:     match,
					Translation:  match,
					Type:         "name",
					AutoDetected: true,
				})
			}
		}
	}

	return entries, nil
}

func isCommonWord(word string) bool {
	common := []string{"The", "This", "That", "With", "From", "Have", "What", "When", "Where", "Which", "While"}
	word = strings.ToLower(word)
	for _, c := range common {
		if strings.ToLower(c) == word {
			return true
		}
	}
	return false
}

func MergeGlossaries(auto, manual []Entry) []Entry {
	merged := make(map[string]Entry)

	for _, e := range auto {
		merged[e.Original] = e
	}

	for _, e := range manual {
		merged[e.Original] = e
	}

	result := []Entry{}
	for _, e := range merged {
		result = append(result, e)
	}

	return result
}
