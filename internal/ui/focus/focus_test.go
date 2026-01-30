package focus

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

func TestModeConstants(t *testing.T) {
	if ModeNav != 0 {
		t.Errorf("ModeNav = %d, want 0", ModeNav)
	}

	if ModeInput != 1 {
		t.Errorf("ModeInput = %d, want 1", ModeInput)
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager(5)

	if m.mode != ModeNav {
		t.Errorf("mode = %d, want ModeNav", m.mode)
	}

	if m.activeField != -1 {
		t.Errorf("activeField = %d, want -1", m.activeField)
	}

	if m.totalFields != 5 {
		t.Errorf("totalFields = %d, want 5", m.totalFields)
	}
}

func TestManagerMode(t *testing.T) {
	m := NewManager(3)

	if m.Mode() != ModeNav {
		t.Errorf("Mode() = %d, want ModeNav", m.Mode())
	}

	m.EnterInput(0)

	if m.Mode() != ModeInput {
		t.Errorf("Mode() = %d, want ModeInput", m.Mode())
	}
}

func TestActiveField(t *testing.T) {
	m := NewManager(3)

	if m.ActiveField() != -1 {
		t.Errorf("ActiveField() = %d, want -1", m.ActiveField())
	}

	m.EnterInput(1)

	if m.ActiveField() != 1 {
		t.Errorf("ActiveField() = %d, want 1", m.ActiveField())
	}
}

func TestIsFieldActive(t *testing.T) {
	m := NewManager(3)
	m.EnterInput(1)

	if !m.IsFieldActive(1) {
		t.Error("IsFieldActive(1) should be true")
	}

	if m.IsFieldActive(0) {
		t.Error("IsFieldActive(0) should be false")
	}

	if m.IsFieldActive(2) {
		t.Error("IsFieldActive(2) should be false")
	}
}

func TestEnterInput(t *testing.T) {
	m := NewManager(3)

	m.EnterInput(2)

	if m.mode != ModeInput {
		t.Errorf("mode = %d, want ModeInput", m.mode)
	}

	if m.activeField != 2 {
		t.Errorf("activeField = %d, want 2", m.activeField)
	}
}

func TestEnterInputOutOfBounds(t *testing.T) {
	m := NewManager(3)

	m.EnterInput(-1)
	if m.mode != ModeNav {
		t.Error("mode should remain ModeNav for negative index")
	}

	m.EnterInput(5)
	if m.mode != ModeNav {
		t.Error("mode should remain ModeNav for out of bounds index")
	}
}

func TestExitInput(t *testing.T) {
	m := NewManager(3)
	m.EnterInput(1)

	m.ExitInput()

	if m.mode != ModeNav {
		t.Errorf("mode = %d, want ModeNav", m.mode)
	}

	if m.activeField != 1 {
		t.Errorf("activeField = %d, should be preserved as 1", m.activeField)
	}
}

func TestCycleNext(t *testing.T) {
	m := NewManager(3)
	m.activeField = 0

	m.CycleNext()
	if m.activeField != 1 {
		t.Errorf("activeField = %d, want 1", m.activeField)
	}

	m.CycleNext()
	if m.activeField != 2 {
		t.Errorf("activeField = %d, want 2", m.activeField)
	}

	m.CycleNext()
	if m.activeField != 0 {
		t.Errorf("activeField = %d, want 0 (wrapped)", m.activeField)
	}
}

func TestCycleNextInInputMode(t *testing.T) {
	m := NewManager(3)
	m.activeField = 0
	m.EnterInput(0)

	m.CycleNext()

	if m.activeField != 0 {
		t.Errorf("activeField should not change in ModeInput, got %d", m.activeField)
	}
}

func TestCyclePrev(t *testing.T) {
	m := NewManager(3)
	m.activeField = 1

	m.CyclePrev()
	if m.activeField != 0 {
		t.Errorf("activeField = %d, want 0", m.activeField)
	}

	m.CyclePrev()
	if m.activeField != 2 {
		t.Errorf("activeField = %d, want 2 (wrapped)", m.activeField)
	}
}

func TestCyclePrevInInputMode(t *testing.T) {
	m := NewManager(3)
	m.activeField = 1
	m.EnterInput(1)

	m.CyclePrev()

	if m.activeField != 1 {
		t.Errorf("activeField should not change in ModeInput, got %d", m.activeField)
	}
}

func TestFieldStyle(t *testing.T) {
	m := NewManager(3)

	style := m.FieldStyle(0, false)
	_ = style.Render("test")

	m.activeField = 0
	style = m.FieldStyle(0, true)
	_ = style.Render("test")

	m.EnterInput(0)
	style = m.FieldStyle(0, true)
	_ = style.Render("test")
}

func TestShouldBlinkCursor(t *testing.T) {
	m := NewManager(3)

	if m.ShouldBlinkCursor(0) {
		t.Error("cursor should not blink when not in input mode")
	}

	m.EnterInput(0)

	if !m.ShouldBlinkCursor(0) {
		t.Error("cursor should blink for active field in input mode")
	}

	if m.ShouldBlinkCursor(1) {
		t.Error("cursor should not blink for non-active field")
	}
}

func TestConfigureInput(t *testing.T) {
	m := NewManager(3)

	input := textinput.New()

	m.ConfigureInput(&input, 0)

	m.EnterInput(0)
	m.ConfigureInput(&input, 0)
}

func TestHandleTabCycle(t *testing.T) {
	m := NewManager(3)
	m.activeField = 0

	handled := m.HandleTabCycle("tab")
	if !handled {
		t.Error("tab should be handled in nav mode")
	}
	if m.activeField != 1 {
		t.Errorf("activeField = %d, want 1", m.activeField)
	}

	handled = m.HandleTabCycle("shift+tab")
	if !handled {
		t.Error("shift+tab should be handled in nav mode")
	}
	if m.activeField != 0 {
		t.Errorf("activeField = %d, want 0", m.activeField)
	}

	handled = m.HandleTabCycle("enter")
	if handled {
		t.Error("enter should not be handled by HandleTabCycle")
	}
}

func TestHandleTabCycleInInputMode(t *testing.T) {
	m := NewManager(3)
	m.EnterInput(0)

	handled := m.HandleTabCycle("tab")
	if handled {
		t.Error("tab should not be handled in input mode")
	}
}

func TestHandleEscape(t *testing.T) {
	m := NewManager(3)

	handled := m.HandleEscape()
	if handled {
		t.Error("escape should not be handled when not in input mode")
	}

	m.EnterInput(0)
	handled = m.HandleEscape()
	if !handled {
		t.Error("escape should be handled in input mode")
	}
	if m.mode != ModeNav {
		t.Errorf("mode = %d, want ModeNav", m.mode)
	}
}

func TestHandleEnter(t *testing.T) {
	m := NewManager(3)
	m.activeField = 1

	handled := m.HandleEnter()
	if !handled {
		t.Error("enter should be handled in nav mode with active field")
	}
	if m.mode != ModeInput {
		t.Errorf("mode = %d, want ModeInput", m.mode)
	}
}

func TestHandleEnterNoActiveField(t *testing.T) {
	m := NewManager(3)

	handled := m.HandleEnter()
	if handled {
		t.Error("enter should not be handled with no active field")
	}
}

func TestHandleEnterInInputMode(t *testing.T) {
	m := NewManager(3)
	m.EnterInput(0)

	handled := m.HandleEnter()
	if handled {
		t.Error("enter should not be handled when already in input mode")
	}
}
