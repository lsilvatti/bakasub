package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// NeonSpinner is a standardized Neon Pink spinner component for BakaSub
type NeonSpinner struct {
	spinner spinner.Model
	active  bool
	label   string
}

// NewNeonSpinner creates a new Neon Pink spinner
func NewNeonSpinner() NeonSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.NeonPink)

	return NeonSpinner{
		spinner: s,
		active:  false,
		label:   "",
	}
}

// NewNeonSpinnerWithLabel creates a spinner with a label
func NewNeonSpinnerWithLabel(label string) NeonSpinner {
	ns := NewNeonSpinner()
	ns.label = label
	return ns
}

// Start activates the spinner
func (ns *NeonSpinner) Start() tea.Cmd {
	ns.active = true
	return ns.spinner.Tick
}

// Stop deactivates the spinner
func (ns *NeonSpinner) Stop() {
	ns.active = false
}

// SetLabel updates the spinner label
func (ns *NeonSpinner) SetLabel(label string) {
	ns.label = label
}

// IsActive returns whether the spinner is active
func (ns *NeonSpinner) IsActive() bool {
	return ns.active
}

// Update handles spinner updates
func (ns NeonSpinner) Update(msg tea.Msg) (NeonSpinner, tea.Cmd) {
	if !ns.active {
		return ns, nil
	}

	var cmd tea.Cmd
	ns.spinner, cmd = ns.spinner.Update(msg)
	return ns, cmd
}

// View renders the spinner
func (ns NeonSpinner) View() string {
	if !ns.active {
		return ""
	}

	if ns.label == "" {
		return ns.spinner.View()
	}

	return ns.spinner.View() + " " + ns.label
}

// ViewWithCustomLabel renders the spinner with a custom inline label
func (ns NeonSpinner) ViewWithCustomLabel(label string) string {
	if !ns.active {
		return ""
	}

	return ns.spinner.View() + " " + label
}

// Alternative spinner styles for variety

// NewNeonSpinnerLine creates a spinner with the Line animation
func NewNeonSpinnerLine() NeonSpinner {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(styles.NeonPink)

	return NeonSpinner{
		spinner: s,
		active:  false,
	}
}

// NewNeonSpinnerPoints creates a spinner with the Points animation
func NewNeonSpinnerPoints() NeonSpinner {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(styles.NeonPink)

	return NeonSpinner{
		spinner: s,
		active:  false,
	}
}

// NewNeonSpinnerMiniDot creates a spinner with the MiniDot animation (compact)
func NewNeonSpinnerMiniDot() NeonSpinner {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(styles.NeonPink)

	return NeonSpinner{
		spinner: s,
		active:  false,
	}
}
