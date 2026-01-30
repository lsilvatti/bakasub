package attachments

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type Attachment struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	MIMEType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type MKVAttachments struct {
	Attachments []Attachment `json:"attachments"`
}

type Mode int

const (
	ModeView Mode = iota
	ModeDelete
	ModeAdd
)

// ClosedMsg is sent when the attachment manager should be closed
type ClosedMsg struct{}

type Model struct {
	table        table.Model
	attachments  []Attachment
	filePath     string
	mode         Mode
	deleteMarked map[int]bool
	addInput     textinput.Model
	width        int
	height       int
	message      string
}

func New(mkvPath string) (*Model, error) {
	attachments, err := listAttachments(mkvPath)
	if err != nil {
		return nil, err
	}

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "FILENAME", Width: 30},
		{Title: "MIME-TYPE", Width: 25},
		{Title: "SIZE", Width: 12},
		{Title: "STATUS", Width: 10},
	}

	rows := buildRows(attachments, make(map[int]bool))

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

	input := textinput.New()
	input.Placeholder = "Path to file..."
	input.Width = 60

	return &Model{
		table:        t,
		attachments:  attachments,
		filePath:     mkvPath,
		mode:         ModeView,
		deleteMarked: make(map[int]bool),
		addInput:     input,
	}, nil
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.resizeTable()
}

func (m Model) Init() tea.Cmd {
	// Request current terminal size
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeTable()

	case tea.KeyMsg:
		// Handle add mode separately
		if m.mode == ModeAdd {
			switch msg.String() {
			case "esc":
				m.mode = ModeView
				m.addInput.SetValue("")
				m.addInput.Blur()
				return m, nil
			case "enter":
				if m.addInput.Value() != "" {
					m.addAttachment(m.addInput.Value())
					m.mode = ModeView
					m.addInput.SetValue("")
					m.addInput.Blur()
				}
				return m, nil
			}
			m.addInput, cmd = m.addInput.Update(msg)
			return m, cmd
		}

		// Normal mode key handling
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return ClosedMsg{} }

		case "d":
			// Toggle delete mode
			if m.mode == ModeDelete {
				m.mode = ModeView
				m.message = ""
			} else {
				m.mode = ModeDelete
				m.message = locales.T("attachments.delete_mode_message")
			}
			m.refreshTable()

		case "a":
			// Enter add mode
			m.mode = ModeAdd
			m.addInput.Focus()
			return m, textinput.Blink

		case "e":
			// Extract all attachments
			m.extractAll()

		case " ":
			// Toggle mark for deletion (only in delete mode)
			if m.mode == ModeDelete {
				selectedIdx := m.table.Cursor()
				if selectedIdx < len(m.attachments) {
					id := m.attachments[selectedIdx].ID
					m.deleteMarked[id] = !m.deleteMarked[id]
					m.refreshTable()
				}
			}

		case "enter":
			// Execute deletion (only in delete mode)
			if m.mode == ModeDelete && len(m.deleteMarked) > 0 {
				m.executeDelete()
				m.mode = ModeView
				m.message = ""
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// Check if terminal is too small or waiting for size
	if m.width == 0 || m.height == 0 {
		return locales.T("common.loading")
	}

	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	header := styles.TitleStyle.Render(locales.T("attachments.title")) + "\n"
	header += styles.SubtleStyle.Render(fmt.Sprintf("%s %s", locales.T("attachments.file_label"), filepath.Base(m.filePath))) + "\n\n"

	if m.mode == ModeAdd {
		return styles.MainWindow.Width(contentWidth).Render(m.renderAddMode(header))
	}

	tableView := m.table.View()

	info := "\n"
	if m.message != "" {
		info += styles.WarningStyle.Render(m.message) + "\n"
	}

	footer := "\n"
	if m.mode == ModeDelete {
		footer += styles.KeyHintStyle.Render("[SPACE]") + " " + locales.T("attachments.mark") + "  "
		footer += styles.KeyHintStyle.Render("[ENTER]") + " " + locales.T("attachments.execute_delete") + "  "
	} else {
		footer += styles.KeyHintStyle.Render("[d]") + " " + locales.T("attachments.delete_mode") + "  "
		footer += styles.KeyHintStyle.Render("[a]") + " " + locales.T("attachments.add_file") + "  "
		footer += styles.KeyHintStyle.Render("[e]") + " " + locales.T("attachments.extract_all") + "  "
	}
	footer += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("attachments.exit")

	return styles.MainWindow.Width(contentWidth).Render(header + tableView + info + footer)
}

func (m *Model) resizeTable() {
	if m.width < 60 {
		return
	}
	// Resize table columns: ID=5, MIME=25, SIZE=12, STATUS=10 = 52 fixed
	availableWidth := m.width - 62
	if availableWidth < 20 {
		availableWidth = 20
	}

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "FILENAME", Width: availableWidth},
		{Title: "MIME-TYPE", Width: 25},
		{Title: "SIZE", Width: 12},
		{Title: "STATUS", Width: 10},
	}
	m.table.SetColumns(columns)
	m.table.SetHeight(m.height - 10)
}

func (m Model) renderAddMode(header string) string {
	content := header
	content += styles.SubtleStyle.Render("┌── "+locales.T("attachments.add_attachment")+" ───────────────────────────────────────────────┐") + "\n"
	content += styles.SubtleStyle.Render("│                                                                     │") + "\n"
	content += styles.SubtleStyle.Render("│  "+locales.T("attachments.path_label")+" ") + m.addInput.View() + "\n"
	content += styles.SubtleStyle.Render("│                                                                     │") + "\n"
	content += styles.SubtleStyle.Render("│  "+locales.T("attachments.embed_note")+"                      │") + "\n"
	content += styles.SubtleStyle.Render("└─────────────────────────────────────────────────────────────────────┘") + "\n\n"
	content += styles.KeyHintStyle.Render("[ENTER]") + " " + locales.T("attachments.confirm") + "  "
	content += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("attachments.cancel")

	return content
}

func (m *Model) refreshTable() {
	m.table.SetRows(buildRows(m.attachments, m.deleteMarked))
}

func (m *Model) addAttachment(filePath string) {
	args := []string{
		m.filePath,
		"--add-attachment", filePath,
	}

	cmd := exec.Command("mkvpropedit", args...)
	if err := cmd.Run(); err != nil {
		m.message = fmt.Sprintf("%s %v", locales.T("attachments.error_label"), err)
		return
	}

	// Reload attachments
	attachments, err := listAttachments(m.filePath)
	if err != nil {
		m.message = fmt.Sprintf("%s %v", locales.T("attachments.error_label"), err)
		return
	}

	m.attachments = attachments
	m.message = locales.T("attachments.attachment_added")
	m.refreshTable()
}

func (m *Model) executeDelete() {
	for id := range m.deleteMarked {
		args := []string{
			m.filePath,
			"--delete-attachment", strconv.Itoa(id),
		}

		cmd := exec.Command("mkvpropedit", args...)
		if err := cmd.Run(); err != nil {
			m.message = fmt.Sprintf("Error deleting %d: %v", id, err)
			continue
		}
	}

	// Reload attachments
	attachments, err := listAttachments(m.filePath)
	if err != nil {
		m.message = fmt.Sprintf("Error reloading: %v", err)
		return
	}

	m.attachments = attachments
	m.deleteMarked = make(map[int]bool)
	m.message = "Attachments deleted successfully"
	m.refreshTable()
}

func (m *Model) extractAll() {
	extractDir := filepath.Join(filepath.Dir(m.filePath), "attachments_extracted")
	os.MkdirAll(extractDir, 0755)

	args := []string{
		"attachments",
		m.filePath,
	}

	for _, att := range m.attachments {
		args = append(args, fmt.Sprintf("%d:%s", att.ID, filepath.Join(extractDir, att.FileName)))
	}

	cmd := exec.Command("mkvextract", args...)
	if err := cmd.Run(); err != nil {
		m.message = fmt.Sprintf("Error extracting: %v", err)
		return
	}

	m.message = fmt.Sprintf("Extracted to: %s", extractDir)
}

func buildRows(attachments []Attachment, deleteMarked map[int]bool) []table.Row {
	rows := make([]table.Row, len(attachments))
	for i, att := range attachments {
		status := ""
		if deleteMarked[att.ID] {
			status = "[DELETE]"
		}

		sizeStr := formatSize(att.Size)

		rows[i] = table.Row{
			strconv.Itoa(att.ID),
			att.FileName,
			att.MIMEType,
			sizeStr,
			status,
		}
	}
	return rows
}

func listAttachments(path string) ([]Attachment, error) {
	cmd := exec.Command("mkvmerge", "-J", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info struct {
		Attachments []struct {
			ID       int    `json:"id"`
			FileName string `json:"file_name"`
			MIMEType string `json:"content_type"`
			Size     int64  `json:"size"`
		} `json:"attachments"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	attachments := make([]Attachment, len(info.Attachments))
	for i, a := range info.Attachments {
		attachments[i] = Attachment{
			ID:       a.ID,
			FileName: a.FileName,
			MIMEType: a.MIMEType,
			Size:     a.Size,
		}
	}

	return attachments, nil
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)

	if bytes < KB {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < MB {
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
}
