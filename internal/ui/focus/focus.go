package focus

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

// Mode represents the current input focus state
type Mode int

const (
	// ModeNav means navigation/global hotkeys are active
	ModeNav Mode = iota
	// ModeInput means a text input has focus and is capturing keystrokes
	ModeInput
)

// Manager handles focus state and visual feedback for text inputs
type Manager struct {
	mode        Mode
	activeField int // Index of the currently focused field (-1 = none)
	totalFields int // Total number of focusable fields
}

// NewManager creates a new focus manager
func NewManager(totalFields int) *Manager {
	return &Manager{
		mode:        ModeNav,
		activeField: -1,
		totalFields: totalFields,
	}
}

// Mode returns the current focus mode
func (m *Manager) Mode() Mode {
	return m.mode
}

// ActiveField returns the index of the currently focused field
func (m *Manager) ActiveField() int {
	return m.activeField
}

// IsFieldActive checks if a specific field is active
func (m *Manager) IsFieldActive(index int) bool {
	return m.activeField == index
}

// EnterInput activates a text input field
func (m *Manager) EnterInput(fieldIndex int) {
	if fieldIndex >= 0 && fieldIndex < m.totalFields {
		m.mode = ModeInput
		m.activeField = fieldIndex
	}
}

// ExitInput returns to navigation mode
func (m *Manager) ExitInput() {
	m.mode = ModeNav
	// Keep activeField so we know which field was selected
}

// CycleNext moves to the next field in navigation mode
func (m *Manager) CycleNext() {
	if m.mode == ModeNav {
		m.activeField = (m.activeField + 1) % m.totalFields
	}
}

// CyclePrev moves to the previous field in navigation mode
func (m *Manager) CyclePrev() {
	if m.mode == ModeNav {
		m.activeField--
		if m.activeField < 0 {
			m.activeField = m.totalFields - 1
		}
	}
}

// FieldStyle returns the appropriate style for a text input field
// based on its focus state. Does NOT add border (textinput has its own rendering)
func (m *Manager) FieldStyle(fieldIndex int, isActive bool) lipgloss.Style {
	baseStyle := lipgloss.NewStyle()

	if !isActive {
		// Non-selected field - dim gray
		return baseStyle.Foreground(lipgloss.Color("#808080"))
	}

	// This field is selected/active
	if m.mode == ModeInput && m.activeField == fieldIndex {
		// Actively editing - Neon Pink
		return baseStyle.Foreground(styles.NeonPink)
	}

	// Selected but in nav mode - Cyan
	return baseStyle.Foreground(styles.Cyan)
}

// ShouldBlinkCursor returns whether the cursor should blink for a field
func (m *Manager) ShouldBlinkCursor(fieldIndex int) bool {
	return m.mode == ModeInput && m.activeField == fieldIndex
}

// ConfigureInput sets up a textinput.Model with the appropriate focus state
func (m *Manager) ConfigureInput(input *textinput.Model, fieldIndex int) {
	if m.mode == ModeInput && m.activeField == fieldIndex {
		input.Focus()
	} else {
		input.Blur()
	}
}

// HandleTabCycle handles Tab/Shift+Tab for field cycling
// Returns true if the key was handled
func (m *Manager) HandleTabCycle(key string) bool {
	if m.mode == ModeNav {
		switch key {
		case "tab":
			m.CycleNext()
			return true
		case "shift+tab":
			m.CyclePrev()
			return true
		}
	}
	return false
}

// HandleEscape handles ESC key - exits input mode
// Returns true if the key was handled
func (m *Manager) HandleEscape() bool {
	if m.mode == ModeInput {
		m.ExitInput()
		return true
	}
	return false
}

// HandleEnter handles ENTER key - enters input mode for selected field
// Returns true if the key was handled
func (m *Manager) HandleEnter() bool {
	if m.mode == ModeNav && m.activeField >= 0 {
		m.EnterInput(m.activeField)
		return true
	}
	return false
}
