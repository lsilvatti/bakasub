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
			return m, tea.Quit

		case "d":
			// Toggle delete mode
			if m.mode == ModeDelete {
				m.mode = ModeView
				m.message = ""
			} else {
				m.mode = ModeDelete
				m.message = "DELETE MODE: Press SPACE to mark, ENTER to execute"
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
	header := styles.TitleStyle.Render("TOOLBOX: ATTACHMENT MANAGER") + "\n"
	header += styles.SubtleStyle.Render(fmt.Sprintf("File: %s", filepath.Base(m.filePath))) + "\n\n"

	if m.mode == ModeAdd {
		return m.renderAddMode(header)
	}

	tableView := m.table.View()

	info := "\n"
	if m.message != "" {
		info += styles.WarningStyle.Render(m.message) + "\n"
	}

	footer := "\n"
	if m.mode == ModeDelete {
		footer += styles.KeyHintStyle.Render("[SPACE]") + " Mark  "
		footer += styles.KeyHintStyle.Render("[ENTER]") + " Execute Delete  "
	} else {
		footer += styles.KeyHintStyle.Render("[d]") + " Delete Mode  "
		footer += styles.KeyHintStyle.Render("[a]") + " Add File  "
		footer += styles.KeyHintStyle.Render("[e]") + " Extract All  "
	}
	footer += styles.KeyHintStyle.Render("[ESC]") + " Exit"

	return header + tableView + info + footer
}

func (m Model) renderAddMode(header string) string {
	content := header
	content += styles.SubtleStyle.Render("┌── ADD NEW ATTACHMENT ───────────────────────────────────────────────┐") + "\n"
	content += styles.SubtleStyle.Render("│                                                                     │") + "\n"
	content += styles.SubtleStyle.Render("│  Path: ") + m.addInput.View() + "\n"
	content += styles.SubtleStyle.Render("│                                                                     │") + "\n"
	content += styles.SubtleStyle.Render("│  *File will be embedded into the MKV container                      │") + "\n"
	content += styles.SubtleStyle.Render("└─────────────────────────────────────────────────────────────────────┘") + "\n\n"
	content += styles.KeyHintStyle.Render("[ENTER]") + " Confirm  "
	content += styles.KeyHintStyle.Render("[ESC]") + " Cancel"

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
		m.message = fmt.Sprintf("Error: %v", err)
		return
	}

	// Reload attachments
	attachments, err := listAttachments(m.filePath)
	if err != nil {
		m.message = fmt.Sprintf("Error reloading: %v", err)
		return
	}

	m.attachments = attachments
	m.message = "Attachment added successfully"
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
