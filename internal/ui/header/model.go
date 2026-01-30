package header

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type Track struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Codec     string `json:"codec"`
	Language  string `json:"properties.language"`
	TrackName string `json:"properties.track_name"`
	Default   bool   `json:"properties.default_track"`
	Forced    bool   `json:"properties.forced_track"`
}

type MKVInfo struct {
	Tracks []Track `json:"tracks"`
}

// ClosedMsg is sent when the header editor should be closed
type ClosedMsg struct{}

type Model struct {
	table    table.Model
	tracks   []Track
	filePath string
	width    int
	height   int
	modified bool
}

func New(mkvPath string) (*Model, error) {
	tracks, err := analyzeMKV(mkvPath)
	if err != nil {
		return nil, err
	}

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "TYPE", Width: 8},
		{Title: "LANG", Width: 8},
		{Title: "TRACK NAME", Width: 30},
		{Title: "DEFAULT", Width: 10},
		{Title: "FORCED", Width: 10},
	}

	rows := make([]table.Row, len(tracks))
	for i, t := range tracks {
		defaultFlag := "[ NO  ]"
		if t.Default {
			defaultFlag = "[ YES ]"
		}
		forcedFlag := "[ NO  ]"
		if t.Forced {
			forcedFlag = "[ YES ]"
		}

		rows[i] = table.Row{
			strconv.Itoa(t.ID),
			t.Type,
			t.Language,
			t.TrackName,
			defaultFlag,
			forcedFlag,
		}
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
		table:    t,
		tracks:   tracks,
		filePath: mkvPath,
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
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return ClosedMsg{} }

		case "d":
			// Toggle Default flag for selected track
			selectedIdx := m.table.Cursor()
			if selectedIdx < len(m.tracks) {
				m.toggleDefault(selectedIdx)
				m.refreshTable()
				m.modified = true
			}

		case "f":
			// Toggle Forced flag for selected track
			selectedIdx := m.table.Cursor()
			if selectedIdx < len(m.tracks) {
				m.toggleForced(selectedIdx)
				m.refreshTable()
				m.modified = true
			}

		case "ctrl+s":
			// Apply changes to file
			if m.modified {
				m.applyChanges()
				m.modified = false
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

	header := styles.TitleStyle.Render(locales.T("header_editor.title")) + "\n"
	header += styles.SubtleStyle.Render(fmt.Sprintf("%s %s", locales.T("header_editor.file_label"), m.filePath)) + "\n\n"

	tableView := m.table.View()

	footer := "\n"
	if m.modified {
		footer += styles.ErrorStyle.Render(locales.T("header_editor.modified_warning")) + "\n"
	}
	footer += styles.KeyHintStyle.Render("[↑/↓]") + " " + locales.T("header_editor.navigate") + "  "
	footer += styles.KeyHintStyle.Render("[d]") + " " + locales.T("header_editor.toggle_default") + "  "
	footer += styles.KeyHintStyle.Render("[f]") + " " + locales.T("header_editor.toggle_forced") + "  "
	footer += styles.KeyHintStyle.Render("[Ctrl+S]") + " " + locales.T("header_editor.apply") + "  "
	footer += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("header_editor.exit")

	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}
	return styles.MainWindow.Width(contentWidth).Render(header + tableView + footer)
}

func (m *Model) resizeTable() {
	if m.width < 60 {
		return
	}
	// Fixed columns: ID=5, TYPE=8, LANG=8, DEFAULT=10, FORCED=10 = 41
	availableWidth := m.width - 51
	if availableWidth < 20 {
		availableWidth = 20
	}

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "TYPE", Width: 8},
		{Title: "LANG", Width: 8},
		{Title: "TRACK NAME", Width: availableWidth},
		{Title: "DEFAULT", Width: 10},
		{Title: "FORCED", Width: 10},
	}
	m.table.SetColumns(columns)
	m.table.SetHeight(m.height - 10)
}

func (m *Model) toggleDefault(idx int) {
	m.tracks[idx].Default = !m.tracks[idx].Default
}

func (m *Model) toggleForced(idx int) {
	m.tracks[idx].Forced = !m.tracks[idx].Forced
}

func (m *Model) refreshTable() {
	rows := make([]table.Row, len(m.tracks))
	for i, t := range m.tracks {
		defaultFlag := "[ NO  ]"
		if t.Default {
			defaultFlag = "[ YES ]"
		}
		forcedFlag := "[ NO  ]"
		if t.Forced {
			forcedFlag = "[ YES ]"
		}

		rows[i] = table.Row{
			strconv.Itoa(t.ID),
			t.Type,
			t.Language,
			t.TrackName,
			defaultFlag,
			forcedFlag,
		}
	}
	m.table.SetRows(rows)
}

func (m *Model) applyChanges() error {
	for _, track := range m.tracks {
		// Build mkvpropedit command
		args := []string{
			m.filePath,
			"--edit", fmt.Sprintf("track:%d", track.ID),
		}

		if track.Default {
			args = append(args, "--set", "flag-default=1")
		} else {
			args = append(args, "--set", "flag-default=0")
		}

		if track.Forced {
			args = append(args, "--set", "flag-forced=1")
		} else {
			args = append(args, "--set", "flag-forced=0")
		}

		cmd := exec.Command("mkvpropedit", args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func analyzeMKV(path string) ([]Track, error) {
	cmd := exec.Command("mkvmerge", "-J", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info struct {
		Tracks []struct {
			ID         int    `json:"id"`
			Type       string `json:"type"`
			Codec      string `json:"codec"`
			Properties struct {
				Language     string `json:"language"`
				TrackName    string `json:"track_name"`
				DefaultTrack bool   `json:"default_track"`
				ForcedTrack  bool   `json:"forced_track"`
			} `json:"properties"`
		} `json:"tracks"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	tracks := make([]Track, len(info.Tracks))
	for i, t := range info.Tracks {
		tracks[i] = Track{
			ID:        t.ID,
			Type:      t.Type,
			Codec:     t.Codec,
			Language:  t.Properties.Language,
			TrackName: t.Properties.TrackName,
			Default:   t.Properties.DefaultTrack,
			Forced:    t.Properties.ForcedTrack,
		}
	}

	return tracks, nil
}
