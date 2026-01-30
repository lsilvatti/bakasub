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
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type State int

const (
	StateView State = iota
	StateAddTerm
)

type Entry struct {
	Original     string `json:"original"`
	Translation  string `json:"translation"`
	Type         string `json:"type"`
	AutoDetected bool   `json:"auto_detected"`
}

// ClosedMsg is sent when the glossary editor should be closed
type ClosedMsg struct{}

type Model struct {
	table    table.Model
	entries  []Entry
	filePath string
	width    int
	height   int
	state    State

	// Add term fields
	addOriginal    string
	addTranslation string
	addType        string
	addField       int // 0=original, 1=translation, 2=type

	// Pagination
	currentPage  int
	itemsPerPage int
}

func New(glossaryPath string) *Model {
	entries := loadOrCreate(glossaryPath)

	// Default columns (will be resized on WindowSizeMsg)
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
			source := "Auto"
			_ = source
		}
		rows[i] = table.Row{e.Original, e.Translation, e.Type, source}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Apply dense table styles for btop aesthetic
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(styles.Yellow).
		Bold(true).
		Padding(0, 1)
	s.Cell = s.Cell.
		Padding(0, 1)
	s.Selected = s.Selected.
		Foreground(styles.NeonPink).
		Background(styles.DarkGray).
		Bold(true).
		Padding(0, 1)
	t.SetStyles(s)

	return &Model{
		table:        t,
		entries:      entries,
		filePath:     glossaryPath,
		currentPage:  1,
		itemsPerPage: 15,
	}
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Resize table columns dynamically
	// Reserve: TYPE=15, SOURCE=10, borders/padding ~10
	availableWidth := width - 35
	halfWidth := availableWidth / 2
	if halfWidth < 15 {
		halfWidth = 15
	}

	columns := []table.Column{
		{Title: "ORIGINAL", Width: halfWidth},
		{Title: "TRANSLATION", Width: halfWidth},
		{Title: "TYPE", Width: 15},
		{Title: "SOURCE", Width: 10},
	}
	m.table.SetColumns(columns)
	m.table.SetHeight(height - 8)
}

func (m *Model) Init() tea.Cmd {
	// Request current terminal size
	return tea.WindowSize()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Handle add term mode
		if m.state == StateAddTerm {
			return m.handleAddTermInput(msg)
		}

		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return ClosedMsg{} }
		case "ctrl+s":
			m.save()
		case "a":
			// Add new term
			m.state = StateAddTerm
			m.addOriginal = ""
			m.addTranslation = ""
			m.addType = "Noun"
			m.addField = 0
			return m, nil
		case "d", "delete", "backspace":
			// Delete selected term
			if len(m.entries) > 0 {
				idx := m.table.Cursor()
				if idx >= 0 && idx < len(m.entries) {
					m.entries = append(m.entries[:idx], m.entries[idx+1:]...)
					m.refreshTable()
				}
			}
			return m, nil
		case "left", "h":
			// Previous page
			if m.currentPage > 1 {
				m.currentPage--
				m.refreshTable()
			}
			return m, nil
		case "right", "l":
			// Next page
			totalPages := m.getTotalPages()
			if m.currentPage < totalPages {
				m.currentPage++
				m.refreshTable()
			}
			return m, nil
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) handleAddTermInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.state = StateView
		return m, nil

	case "enter":
		if m.addField < 2 {
			m.addField++
		} else {
			// Save the new term
			if m.addOriginal != "" && m.addTranslation != "" {
				newEntry := Entry{
					Original:     m.addOriginal,
					Translation:  m.addTranslation,
					Type:         m.addType,
					AutoDetected: false,
				}
				m.entries = append(m.entries, newEntry)
				m.refreshTable()
			}
			m.state = StateView
		}
		return m, nil

	case "tab":
		m.addField = (m.addField + 1) % 3
		return m, nil

	case "backspace":
		switch m.addField {
		case 0:
			if len(m.addOriginal) > 0 {
				m.addOriginal = m.addOriginal[:len(m.addOriginal)-1]
			}
		case 1:
			if len(m.addTranslation) > 0 {
				m.addTranslation = m.addTranslation[:len(m.addTranslation)-1]
			}
		case 2:
			if len(m.addType) > 0 {
				m.addType = m.addType[:len(m.addType)-1]
			}
		}
		return m, nil

	default:
		// Handle character input
		if len(key) == 1 {
			switch m.addField {
			case 0:
				m.addOriginal += key
			case 1:
				m.addTranslation += key
			case 2:
				m.addType += key
			}
		}
	}

	return m, nil
}

func (m *Model) refreshTable() {
	// Get paginated entries
	start := (m.currentPage - 1) * m.itemsPerPage
	end := start + m.itemsPerPage
	if end > len(m.entries) {
		end = len(m.entries)
	}
	if start > len(m.entries) {
		start = 0
		m.currentPage = 1
	}

	paginatedEntries := m.entries[start:end]
	rows := make([]table.Row, len(paginatedEntries))
	for i, e := range paginatedEntries {
		source := "Manual"
		if e.AutoDetected {
			source = "Auto"
		}
		rows[i] = table.Row{e.Original, e.Translation, e.Type, source}
	}
	m.table.SetRows(rows)
}

func (m *Model) getTotalPages() int {
	if len(m.entries) == 0 {
		return 1
	}
	pages := len(m.entries) / m.itemsPerPage
	if len(m.entries)%m.itemsPerPage > 0 {
		pages++
	}
	return pages
}

func (m *Model) View() string {
	// Check if terminal is too small or waiting for size
	if m.width == 0 || m.height == 0 {
		return locales.T("common.loading")
	}

	if m.state == StateAddTerm {
		return m.viewAddTerm()
	}

	header := styles.TitleStyle.Render(locales.T("glossary_editor.title")) + "\n"
	header += styles.SubtleStyle.Render(fmt.Sprintf("%s %s", locales.T("glossary_editor.file_label"), m.filePath)) + "\n\n"

	// Terms management section
	controls := styles.SectionStyle.Render("┌── "+locales.T("glossary_editor.terms_management")+" ──────────────────────────────────────────────────┐") + "\n"
	controls += "│  " + styles.KeyHintStyle.Render("[ a ]") + " " + locales.T("glossary_editor.add_new_term") + "    "
	controls += styles.KeyHintStyle.Render("[ i ]") + " " + locales.T("glossary_editor.import_csv") + "    "
	controls += styles.KeyHintStyle.Render("[DEL]") + " " + locales.T("glossary_editor.remove_selected") + "\n"
	controls += styles.SectionStyle.Render("└──────────────────────────────────────────────────────────────────────┘") + "\n\n"

	tableView := m.table.View()

	// Pagination indicator
	totalPages := m.getTotalPages()
	pagination := fmt.Sprintf("\n  %s   "+locales.T("glossary_editor.page")+" %d/%d   %s\n",
		styles.KeyHintStyle.Render("[< "+locales.T("glossary_editor.prev_page")+"]"),
		m.currentPage,
		totalPages,
		styles.KeyHintStyle.Render("["+locales.T("glossary_editor.next_page")+" >]"),
	)

	footer := "\n" + styles.KeyHintStyle.Render("[↑/↓]") + " " + locales.T("glossary_editor.navigate") + "  "
	footer += styles.KeyHintStyle.Render("[←/→]") + " " + locales.T("glossary_editor.page") + "  "
	footer += styles.KeyHintStyle.Render("[a]") + " " + locales.T("glossary_editor.add") + "  "
	footer += styles.KeyHintStyle.Render("[d]") + " " + locales.T("glossary_editor.delete") + "  "
	footer += styles.KeyHintStyle.Render("[Ctrl+S]") + " " + locales.T("glossary_editor.save") + "  "
	footer += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("glossary_editor.exit")

	return styles.MainWindow.Width(m.width - 4).Render(header + controls + tableView + pagination + footer)
}

func (m *Model) viewAddTerm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(locales.T("glossary_editor.add_new_term")))
	b.WriteString("\n\n")

	// Original term
	origStyle := styles.SubtleStyle
	if m.addField == 0 {
		origStyle = styles.HighlightStyle
	}
	b.WriteString(origStyle.Render(locales.T("glossary_editor.original_term") + " "))
	b.WriteString(fmt.Sprintf("[%s]", m.addOriginal))
	if m.addField == 0 {
		b.WriteString("_")
	}
	b.WriteString("\n\n")

	// Translation
	transStyle := styles.SubtleStyle
	if m.addField == 1 {
		transStyle = styles.HighlightStyle
	}
	b.WriteString(transStyle.Render(locales.T("glossary_editor.translation_label") + " "))
	b.WriteString(fmt.Sprintf("[%s]", m.addTranslation))
	if m.addField == 1 {
		b.WriteString("_")
	}
	b.WriteString("\n\n")

	// Type
	typeStyle := styles.SubtleStyle
	if m.addField == 2 {
		typeStyle = styles.HighlightStyle
	}
	b.WriteString(typeStyle.Render(locales.T("glossary_editor.type_label") + " "))
	b.WriteString(fmt.Sprintf("[%s]", m.addType))
	if m.addField == 2 {
		b.WriteString("_")
	}
	b.WriteString("\n\n")

	b.WriteString(styles.KeyHintStyle.Render("[TAB]") + " " + locales.T("glossary_editor.next_field") + "  ")
	b.WriteString(styles.KeyHintStyle.Render("[ENTER]") + " " + locales.T("glossary_editor.confirm") + "  ")
	b.WriteString(styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("common.cancel"))

	return styles.ModalStyle.Width(60).Render(b.String())
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
