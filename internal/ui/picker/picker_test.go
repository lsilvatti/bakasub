package picker

import (
	"testing"
)

// TestSelectionModeConstants tests SelectionMode constants
func TestSelectionModeConstants(t *testing.T) {
	if ModeDirectory != 0 {
		t.Errorf("ModeDirectory = %d, want 0", ModeDirectory)
	}

	if ModeFile != 1 {
		t.Errorf("ModeFile = %d, want 1", ModeFile)
	}

	if ModeBoth != 2 {
		t.Errorf("ModeBoth = %d, want 2", ModeBoth)
	}
}

// TestSelectedPathMsgStruct tests SelectedPathMsg structure
func TestSelectedPathMsgStruct(t *testing.T) {
	msg := SelectedPathMsg{
		Path:    "/test/path",
		IsDir:   true,
		Aborted: false,
	}

	if msg.Path != "/test/path" {
		t.Errorf("Path = %q, want /test/path", msg.Path)
	}

	if !msg.IsDir {
		t.Error("IsDir should be true")
	}

	if msg.Aborted {
		t.Error("Aborted should be false")
	}
}

// TestSelectedPathMsgAborted tests aborted selection
func TestSelectedPathMsgAborted(t *testing.T) {
	msg := SelectedPathMsg{
		Path:    "",
		IsDir:   false,
		Aborted: true,
	}

	if !msg.Aborted {
		t.Error("Aborted should be true")
	}
}

// TestNew tests New function with default directory
func TestNew(t *testing.T) {
	model := New("", ModeFile)

	// Should not panic and should have valid state
	if model.selectionMode != ModeFile {
		t.Errorf("selectionMode = %d, want ModeFile", model.selectionMode)
	}
}

// TestNewWithNonExistentDir tests New with non-existent directory
func TestNewWithNonExistentDir(t *testing.T) {
	model := New("/nonexistent/path/that/wont/exist", ModeDirectory)

	// Should fall back to home directory
	if model.selectionMode != ModeDirectory {
		t.Errorf("selectionMode = %d, want ModeDirectory", model.selectionMode)
	}
}

// TestNewWithModeDirectory tests New with directory mode
func TestNewWithModeDirectory(t *testing.T) {
	model := New("", ModeDirectory)

	if model.selectionMode != ModeDirectory {
		t.Errorf("selectionMode = %d, want ModeDirectory", model.selectionMode)
	}
}

// TestNewWithModeBoth tests New with both mode
func TestNewWithModeBoth(t *testing.T) {
	model := New("", ModeBoth)

	if model.selectionMode != ModeBoth {
		t.Errorf("selectionMode = %d, want ModeBoth", model.selectionMode)
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		selectedPath:  "/test/path",
		selectionMode: ModeFile,
		width:         80,
		height:        24,
		quitting:      false,
		title:         "Select a file",
	}

	if model.selectedPath != "/test/path" {
		t.Errorf("selectedPath = %q, want /test/path", model.selectedPath)
	}

	if model.selectionMode != ModeFile {
		t.Errorf("selectionMode = %d, want ModeFile", model.selectionMode)
	}

	if model.width != 80 {
		t.Errorf("width = %d, want 80", model.width)
	}

	if model.height != 24 {
		t.Errorf("height = %d, want 24", model.height)
	}

	if model.quitting {
		t.Error("quitting should be false")
	}
}

// TestNeonColors tests that color variables are defined
func TestNeonColors(t *testing.T) {
	// Verify colors are non-empty
	if neonPink == "" {
		t.Error("neonPink should not be empty")
	}

	if cyan == "" {
		t.Error("cyan should not be empty")
	}

	if yellow == "" {
		t.Error("yellow should not be empty")
	}
}
