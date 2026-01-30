package remuxer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type Track struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Codec    string `json:"codec"`
	Language string
	Name     string
	Selected bool
}

type State int

const (
	StateSelect State = iota
	StateConfirm
	StateAddExternal
	StateProcessing
	StateDone
)

// ExternalTrack represents an external file to add to the mux
type ExternalTrack struct {
	Path     string
	Type     string // audio, subtitle
	Language string
}

// ClosedMsg is sent when the remuxer should be closed
type ClosedMsg struct{}

type Model struct {
	tracks         []Track
	externalTracks []ExternalTrack
	filePath       string
	outputPath     string
	cursor         int
	state          State
	progress       string
	error          string
	width          int
	height         int

	// Add external track input
	addExternalPath string
	addExternalLang string
}

func New(mkvPath string) (*Model, error) {
	tracks, err := listTracks(mkvPath)
	if err != nil {
		return nil, err
	}

	// By default, all tracks are selected
	for i := range tracks {
		tracks[i].Selected = true
	}

	outputPath := generateOutputPath(mkvPath)

	return &Model{
		tracks:     tracks,
		filePath:   mkvPath,
		outputPath: outputPath,
		state:      StateSelect,
	}, nil
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Init() tea.Cmd {
	// Request current terminal size
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.state {
		case StateSelect:
			return m.updateSelect(msg)
		case StateAddExternal:
			return m.updateAddExternal(msg)
		case StateConfirm:
			return m.updateConfirm(msg)
		case StateDone:
			if msg.String() == "enter" || msg.String() == "q" {
				return m, func() tea.Msg { return ClosedMsg{} }
			}
		}
	}

	return m, nil
}

func (m Model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		return m, func() tea.Msg { return ClosedMsg{} }

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		totalItems := len(m.tracks) + len(m.externalTracks)
		if m.cursor < totalItems-1 {
			m.cursor++
		}

	case " ":
		// Toggle selection
		if m.cursor < len(m.tracks) {
			m.tracks[m.cursor].Selected = !m.tracks[m.cursor].Selected
		}

	case "a":
		// Add external track
		m.state = StateAddExternal
		return m, nil

	case "d", "delete", "backspace":
		// Remove external track if selected
		externalIdx := m.cursor - len(m.tracks)
		if externalIdx >= 0 && externalIdx < len(m.externalTracks) {
			m.externalTracks = append(m.externalTracks[:externalIdx], m.externalTracks[externalIdx+1:]...)
			if m.cursor > 0 && m.cursor >= len(m.tracks)+len(m.externalTracks) {
				m.cursor--
			}
		}

	case "enter":
		// Move to confirmation
		m.state = StateConfirm
	}

	return m, nil
}

func (m Model) updateAddExternal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateSelect
		m.addExternalPath = ""
		m.addExternalLang = ""
		return m, nil

	case "enter":
		// For now, this is a placeholder. In a full implementation,
		// this would integrate with a file picker.
		// For the TUI version, we'd need a text input for the path.
		m.state = StateSelect
		return m, nil

	case "1":
		// Placeholder: add a mock audio track
		m.externalTracks = append(m.externalTracks, ExternalTrack{
			Path:     "external_audio.ac3",
			Type:     "audio",
			Language: "por",
		})
		m.state = StateSelect

	case "2":
		// Placeholder: add a mock subtitle track
		m.externalTracks = append(m.externalTracks, ExternalTrack{
			Path:     "external_subs.ass",
			Type:     "subtitle",
			Language: "por",
		})
		m.state = StateSelect
	}

	return m, nil
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateSelect
		return m, nil

	case "y", "Y":
		// Execute remux
		m.state = StateProcessing
		return m, m.executeRemux()

	case "n", "N":
		m.state = StateSelect
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	// Check if terminal is too small or waiting for size
	if m.width == 0 || m.height == 0 {
		return locales.T("common.loading")
	}

	switch m.state {
	case StateSelect:
		return m.viewSelect()
	case StateAddExternal:
		return m.viewAddExternal()
	case StateConfirm:
		return m.viewConfirm()
	case StateProcessing:
		return m.viewProcessing()
	case StateDone:
		return m.viewDone()
	}
	return ""
}

func (m Model) viewSelect() string {
	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(locales.T("remuxer.title")))
	b.WriteString("\n")
	b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("%s %s", locales.T("remuxer.file_label"), filepath.Base(m.filePath))))
	b.WriteString("\n\n")

	// Track list section
	var trackList strings.Builder
	trackList.WriteString(styles.SectionStyle.Render(locales.T("remuxer.select_tracks")) + "\n\n")

	for i, track := range m.tracks {
		prefix := "[ ] "
		if track.Selected {
			prefix = "[X] "
		}

		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}

		line := fmt.Sprintf("%s%s%d: %s (%s) - %s - %s",
			cursor,
			prefix,
			track.ID,
			track.Type,
			track.Codec,
			track.Language,
			track.Name,
		)

		if i == m.cursor {
			line = styles.HighlightStyle.Render(line)
		}

		trackList.WriteString("   " + line + "\n")
	}
	b.WriteString(styles.Panel.Render(trackList.String()))
	b.WriteString("\n\n")

	// Add External Tracks section
	var externalSection strings.Builder
	externalSection.WriteString(styles.SectionStyle.Render(locales.T("remuxer.add_external")) + "\n\n")
	externalSection.WriteString("   " + styles.KeyHintStyle.Render("[ a ]") + " " + locales.T("remuxer.add_file") + "\n")

	if len(m.externalTracks) > 0 {
		externalSection.WriteString("\n   " + locales.T("remuxer.pending_addition") + "\n")
		for i, ext := range m.externalTracks {
			idx := len(m.tracks) + i
			cursor := "  "
			if idx == m.cursor {
				cursor = "> "
			}
			line := fmt.Sprintf("%s• %s  [ %s %s ]", cursor, ext.Path, locales.T("remuxer.lang_label"), ext.Language)
			if idx == m.cursor {
				line = styles.HighlightStyle.Render(line)
			}
			externalSection.WriteString("   " + line + "\n")
		}
	}
	b.WriteString(styles.Panel.Render(externalSection.String()))
	b.WriteString("\n\n")

	// Summary
	selected := 0
	for _, t := range m.tracks {
		if t.Selected {
			selected++
		}
	}
	b.WriteString(styles.SubtleStyle.Render(locales.Tf("remuxer.selected_summary", selected, len(m.tracks), len(m.externalTracks))))
	b.WriteString("\n\n")

	// Footer
	footer := styles.KeyHintStyle.Render("[↑/↓]") + " " + locales.T("remuxer.navigate") + "  "
	footer += styles.KeyHintStyle.Render("[SPACE]") + " " + locales.T("remuxer.toggle") + "  "
	footer += styles.KeyHintStyle.Render("[a]") + " " + locales.T("remuxer.add_external_short") + "  "
	if len(m.externalTracks) > 0 && m.cursor >= len(m.tracks) {
		footer += styles.KeyHintStyle.Render("[d]") + " " + locales.T("remuxer.remove") + "  "
	}
	footer += styles.KeyHintStyle.Render("[ENTER]") + " " + locales.T("remuxer.continue") + "  "
	footer += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("remuxer.cancel")
	b.WriteString(footer)

	return styles.MainWindow.Width(contentWidth).Render(b.String())
}

func (m Model) viewConfirm() string {
	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(locales.T("remuxer.confirm_title")))
	b.WriteString("\n\n")

	b.WriteString(styles.WarningStyle.Render("⚠ " + locales.T("remuxer.warning_new_file")))
	b.WriteString("\n\n")

	// Output section with Panel
	var outputSection strings.Builder
	outputSection.WriteString(styles.SectionStyle.Render(locales.T("remuxer.output_label")) + "\n\n")
	outputSection.WriteString("   " + locales.T("remuxer.new_file_label") + " " + m.outputPath)
	b.WriteString(styles.Panel.Render(outputSection.String()))
	b.WriteString("\n\n")

	// Show what will be included
	b.WriteString(styles.SubtleStyle.Render(locales.T("remuxer.tracks_to_include")))
	b.WriteString("\n")
	for _, track := range m.tracks {
		if track.Selected {
			line := fmt.Sprintf("  • Track %d: %s (%s)", track.ID, track.Type, track.Codec)
			b.WriteString(styles.SuccessStyle.Render(line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.KeyHintStyle.Render("[Y]") + " " + locales.T("remuxer.confirm") + "  ")
	b.WriteString(styles.KeyHintStyle.Render("[N]") + " " + locales.T("remuxer.cancel"))

	return styles.MainWindow.Width(contentWidth).Render(b.String())
}

func (m Model) viewAddExternal() string {
	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(locales.T("remuxer.add_external_track")))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtleStyle.Render(locales.T("remuxer.select_track_type")))
	b.WriteString("\n\n")

	// Placeholder options (in a real implementation, this would be a file picker)
	b.WriteString(styles.KeyHintStyle.Render("[1]") + " " + locales.T("remuxer.add_audio_track"))
	b.WriteString("\n")
	b.WriteString(styles.KeyHintStyle.Render("[2]") + " " + locales.T("remuxer.add_subtitle_track"))
	b.WriteString("\n\n")

	b.WriteString(styles.Dimmed.Render(locales.T("remuxer.future_file_browser")))
	b.WriteString("\n")
	b.WriteString(styles.Dimmed.Render(locales.T("remuxer.use_cli_note")))
	b.WriteString("\n\n")

	b.WriteString(styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("remuxer.cancel"))

	return styles.ModalStyle.Width(60).Render(b.String())
}

func (m Model) viewProcessing() string {
	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("REMUXING..."))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtleStyle.Render(m.progress))
	b.WriteString("\n\n")

	if m.error != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.error)))
	}

	return styles.MainWindow.Width(contentWidth).Render(b.String())
}

func (m Model) viewDone() string {
	contentWidth := m.width - 4
	if contentWidth < 60 {
		contentWidth = 60
	}

	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.ErrorStyle.Render(locales.T("remuxer.remux_failed")))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtleStyle.Render(m.error))
	} else {
		b.WriteString(styles.SuccessStyle.Render(locales.T("remuxer.remux_complete")))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("%s %s", locales.T("remuxer.created"), m.outputPath)))
	}

	b.WriteString("\n\n")
	b.WriteString(styles.KeyHintStyle.Render("[ENTER]") + " " + locales.T("common.close"))

	return styles.MainWindow.Width(contentWidth).Render(b.String())
}

func (m *Model) executeRemux() tea.Cmd {
	return func() tea.Msg {
		// Build mkvmerge command
		args := []string{
			"-o", m.outputPath,
		}

		// Add track selection
		for _, track := range m.tracks {
			if !track.Selected {
				args = append(args, "--no-"+track.Type+"s")
			}
		}

		// Add track IDs explicitly
		trackList := []string{}
		for _, track := range m.tracks {
			if track.Selected {
				trackList = append(trackList, strconv.Itoa(track.ID))
			}
		}
		if len(trackList) > 0 {
			args = append(args, "-d", strings.Join(trackList, ","))
			args = append(args, "-a", strings.Join(trackList, ","))
			args = append(args, "-s", strings.Join(trackList, ","))
		}

		args = append(args, m.filePath)

		m.progress = "Executing mkvmerge..."

		cmd := exec.Command("mkvmerge", args...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			m.error = fmt.Sprintf("%v\n%s", err, output)
		} else {
			m.error = ""
		}

		m.state = StateDone
		return nil
	}
}

func listTracks(path string) ([]Track, error) {
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
				Language  string `json:"language"`
				TrackName string `json:"track_name"`
			} `json:"properties"`
		} `json:"tracks"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	tracks := make([]Track, len(info.Tracks))
	for i, t := range info.Tracks {
		tracks[i] = Track{
			ID:       t.ID,
			Type:     t.Type,
			Codec:    t.Codec,
			Language: t.Properties.Language,
			Name:     t.Properties.TrackName,
			Selected: true,
		}
	}

	return tracks, nil
}

func generateOutputPath(originalPath string) string {
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(originalPath, ext)
	return base + "_remuxed" + ext
}
