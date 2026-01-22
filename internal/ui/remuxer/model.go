package remuxer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	StateProcessing
	StateDone
)

type Model struct {
	tracks     []Track
	filePath   string
	outputPath string
	cursor     int
	state      State
	progress   string
	error      string
	width      int
	height     int
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

func (m Model) Init() tea.Cmd {
	return nil
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
		case StateConfirm:
			return m.updateConfirm(msg)
		case StateDone:
			if msg.String() == "enter" || msg.String() == "q" {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) updateSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.tracks)-1 {
			m.cursor++
		}

	case " ":
		// Toggle selection
		m.tracks[m.cursor].Selected = !m.tracks[m.cursor].Selected

	case "enter":
		// Move to confirmation
		m.state = StateConfirm
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
	switch m.state {
	case StateSelect:
		return m.viewSelect()
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
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("TOOLBOX: QUICK REMUXER"))
	b.WriteString("\n")
	b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("File: %s", filepath.Base(m.filePath))))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtleStyle.Render("┌── SELECT TRACKS TO KEEP ────────────────────────────────────────────┐"))
	b.WriteString("\n")

	// Track list
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
		} else {
			line = styles.SubtleStyle.Render("│ ") + line
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString(styles.SubtleStyle.Render("└─────────────────────────────────────────────────────────────────────┘"))
	b.WriteString("\n\n")

	// Summary
	selected := 0
	for _, t := range m.tracks {
		if t.Selected {
			selected++
		}
	}
	b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("Selected: %d / %d tracks", selected, len(m.tracks))))
	b.WriteString("\n\n")

	// Footer
	footer := styles.KeyHintStyle.Render("[↑/↓]") + " Navigate  "
	footer += styles.KeyHintStyle.Render("[SPACE]") + " Toggle  "
	footer += styles.KeyHintStyle.Render("[ENTER]") + " Continue  "
	footer += styles.KeyHintStyle.Render("[ESC]") + " Cancel"
	b.WriteString(footer)

	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("CONFIRM REMUX"))
	b.WriteString("\n\n")

	b.WriteString(styles.WarningStyle.Render("⚠ WARNING: This will create a new file"))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtleStyle.Render("┌── OUTPUT ───────────────────────────────────────────────────────────┐"))
	b.WriteString("\n")
	b.WriteString(styles.SubtleStyle.Render("│  New File: ") + m.outputPath + "\n")
	b.WriteString(styles.SubtleStyle.Render("└─────────────────────────────────────────────────────────────────────┘"))
	b.WriteString("\n\n")

	// Show what will be included
	b.WriteString(styles.SubtleStyle.Render("Tracks to include:"))
	b.WriteString("\n")
	for _, track := range m.tracks {
		if track.Selected {
			line := fmt.Sprintf("  • Track %d: %s (%s)", track.ID, track.Type, track.Codec)
			b.WriteString(styles.SuccessStyle.Render(line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.KeyHintStyle.Render("[Y]") + " Confirm  ")
	b.WriteString(styles.KeyHintStyle.Render("[N]") + " Cancel")

	return b.String()
}

func (m Model) viewProcessing() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("REMUXING..."))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtleStyle.Render(m.progress))
	b.WriteString("\n\n")

	if m.error != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.error)))
	}

	return b.String()
}

func (m Model) viewDone() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.ErrorStyle.Render("REMUX FAILED"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtleStyle.Render(m.error))
	} else {
		b.WriteString(styles.SuccessStyle.Render("REMUX COMPLETE"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtleStyle.Render(fmt.Sprintf("Created: %s", m.outputPath)))
	}

	b.WriteString("\n\n")
	b.WriteString(styles.KeyHintStyle.Render("[ENTER]") + " Exit")

	return b.String()
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
