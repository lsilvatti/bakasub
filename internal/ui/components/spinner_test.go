package components

import (
	"strings"
	"testing"
)

func TestNewNeonSpinner(t *testing.T) {
	spinner := NewNeonSpinner()

	if spinner.active {
		t.Error("new spinner should not be active")
	}

	if spinner.label != "" {
		t.Errorf("new spinner should have empty label, got %q", spinner.label)
	}
}

func TestNewNeonSpinnerWithLabel(t *testing.T) {
	label := "Loading..."
	spinner := NewNeonSpinnerWithLabel(label)

	if spinner.label != label {
		t.Errorf("spinner label = %q, want %q", spinner.label, label)
	}

	if spinner.active {
		t.Error("new spinner should not be active")
	}
}

func TestNeonSpinnerStart(t *testing.T) {
	spinner := NewNeonSpinner()

	cmd := spinner.Start()

	if !spinner.active {
		t.Error("spinner should be active after Start()")
	}

	if cmd == nil {
		t.Error("Start() should return a tea.Cmd")
	}
}

func TestNeonSpinnerStop(t *testing.T) {
	spinner := NewNeonSpinner()
	_ = spinner.Start()

	spinner.Stop()

	if spinner.active {
		t.Error("spinner should not be active after Stop()")
	}
}

func TestNeonSpinnerSetLabel(t *testing.T) {
	spinner := NewNeonSpinner()

	spinner.SetLabel("New Label")

	if spinner.label != "New Label" {
		t.Errorf("label = %q, want %q", spinner.label, "New Label")
	}
}

func TestNeonSpinnerIsActive(t *testing.T) {
	spinner := NewNeonSpinner()

	if spinner.IsActive() {
		t.Error("new spinner should not be active")
	}

	_ = spinner.Start()

	if !spinner.IsActive() {
		t.Error("spinner should be active after Start()")
	}

	spinner.Stop()

	if spinner.IsActive() {
		t.Error("spinner should not be active after Stop()")
	}
}

func TestNeonSpinnerViewInactive(t *testing.T) {
	spinner := NewNeonSpinner()

	view := spinner.View()

	if view != "" {
		t.Errorf("inactive spinner View() = %q, want empty string", view)
	}
}

func TestNeonSpinnerViewActive(t *testing.T) {
	spinner := NewNeonSpinner()
	_ = spinner.Start()

	view := spinner.View()

	// Active spinner should render something
	// (The actual content depends on the spinner state)
	if view == "" {
		t.Error("active spinner should render something")
	}
}

func TestNeonSpinnerViewWithLabel(t *testing.T) {
	spinner := NewNeonSpinnerWithLabel("Processing")
	_ = spinner.Start()

	view := spinner.View()

	if !strings.Contains(view, "Processing") {
		t.Errorf("view should contain label, got %q", view)
	}
}

func TestNeonSpinnerViewWithCustomLabel(t *testing.T) {
	spinner := NewNeonSpinner()
	_ = spinner.Start()

	view := spinner.ViewWithCustomLabel("Custom")

	if !strings.Contains(view, "Custom") {
		t.Errorf("view should contain custom label, got %q", view)
	}
}

func TestNeonSpinnerViewWithCustomLabelInactive(t *testing.T) {
	spinner := NewNeonSpinner()

	view := spinner.ViewWithCustomLabel("Custom")

	if view != "" {
		t.Errorf("inactive spinner ViewWithCustomLabel() = %q, want empty string", view)
	}
}

func TestNewNeonSpinnerLine(t *testing.T) {
	spinner := NewNeonSpinnerLine()

	if spinner.active {
		t.Error("new spinner should not be active")
	}
}

func TestNeonSpinnerUpdate(t *testing.T) {
	spinner := NewNeonSpinner()
	_ = spinner.Start()

	// Update should not panic with nil message
	updatedSpinner, _ := spinner.Update(nil)

	if !updatedSpinner.active {
		t.Error("spinner should still be active after Update")
	}
}

func TestNeonSpinnerUpdateInactive(t *testing.T) {
	spinner := NewNeonSpinner()

	// Update on inactive spinner should return nil cmd
	updatedSpinner, cmd := spinner.Update(nil)

	if cmd != nil {
		t.Error("Update on inactive spinner should return nil cmd")
	}

	if updatedSpinner.active {
		t.Error("spinner should remain inactive")
	}
}
