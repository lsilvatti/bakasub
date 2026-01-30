package settings

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/config"
)

// TestTabConstants tests Tab constants
func TestTabConstants(t *testing.T) {
	if TabGeneral != 0 {
		t.Errorf("TabGeneral = %d, want 0", TabGeneral)
	}

	if TabProviders != 1 {
		t.Errorf("TabProviders = %d, want 1", TabProviders)
	}

	if TabModels != 2 {
		t.Errorf("TabModels = %d, want 2", TabModels)
	}

	if TabPrompts != 3 {
		t.Errorf("TabPrompts = %d, want 3", TabPrompts)
	}

	if TabAdvanced != 4 {
		t.Errorf("TabAdvanced = %d, want 4", TabAdvanced)
	}
}

// TestSavedMsgStruct tests SavedMsg structure
func TestSavedMsgStruct(t *testing.T) {
	cfg := config.Default()
	msg := SavedMsg{
		Config: cfg,
	}

	if msg.Config == nil {
		t.Error("Config should not be nil")
	}
}

// TestCancelledMsgStruct tests CancelledMsg structure
func TestCancelledMsgStruct(t *testing.T) {
	msg := CancelledMsg{}

	// Should exist
	_ = msg
}

// TestModelsLoadedMsgStruct tests modelsLoadedMsg structure
func TestModelsLoadedMsgStruct(t *testing.T) {
	msg := modelsLoadedMsg{
		models: nil,
		err:    nil,
	}

	if msg.models != nil {
		t.Error("models should be nil")
	}

	if msg.err != nil {
		t.Error("err should be nil")
	}
}

// TestNewWithNilConfig tests New with nil config
func TestNewWithNilConfig(t *testing.T) {
	model := New(nil)

	// Should not panic and should use default config
	if model.config == nil {
		t.Error("config should not be nil")
	}
}

// TestNewWithValidConfig tests New with valid config
func TestNewWithValidConfig(t *testing.T) {
	cfg := config.Default()
	cfg.APIKey = "test-key"

	model := New(cfg)

	if model.config == nil {
		t.Error("config should not be nil")
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		width:      120,
		height:     40,
		activeTab:  TabGeneral,
		saved:      false,
		hasError:   false,
		errMsg:     "",
		showAPIKey: false,
	}

	if model.width != 120 {
		t.Errorf("width = %d, want 120", model.width)
	}

	if model.height != 40 {
		t.Errorf("height = %d, want 40", model.height)
	}

	if model.activeTab != TabGeneral {
		t.Errorf("activeTab = %d, want TabGeneral", model.activeTab)
	}

	if model.saved {
		t.Error("saved should be false")
	}

	if model.hasError {
		t.Error("hasError should be false")
	}
}

// TestTabProgression tests tab progression values
func TestTabProgression(t *testing.T) {
	if TabProviders <= TabGeneral {
		t.Error("TabProviders should be greater than TabGeneral")
	}

	if TabModels <= TabProviders {
		t.Error("TabModels should be greater than TabProviders")
	}

	if TabPrompts <= TabModels {
		t.Error("TabPrompts should be greater than TabModels")
	}

	if TabAdvanced <= TabPrompts {
		t.Error("TabAdvanced should be greater than TabPrompts")
	}
}

// TestProviderSelection tests provider selection values
func TestProviderSelection(t *testing.T) {
	providers := []struct {
		idx  int
		name string
	}{
		{0, "openrouter"},
		{1, "gemini"},
		{2, "openai"},
		{3, "local"},
	}

	for _, p := range providers {
		model := Model{
			selectedProvider: p.idx,
		}

		if model.selectedProvider != p.idx {
			t.Errorf("selectedProvider = %d for %s, want %d", model.selectedProvider, p.name, p.idx)
		}
	}
}

// TestLogLevelSelection tests log level selection
func TestLogLevelSelection(t *testing.T) {
	logLevels := []struct {
		idx  int
		name string
	}{
		{0, "info"},
		{1, "debug"},
	}

	for _, ll := range logLevels {
		model := Model{
			selectedLogLevel: ll.idx,
		}

		if model.selectedLogLevel != ll.idx {
			t.Errorf("selectedLogLevel = %d for %s, want %d", model.selectedLogLevel, ll.name, ll.idx)
		}
	}
}

// TestTouchlessModalState tests touchless modal state
func TestTouchlessModalState(t *testing.T) {
	model := Model{
		showTouchlessModal: true,
		touchlessMultiSub:  0,
		touchlessProfile:   0,
		touchlessMuxMode:   0,
	}

	if !model.showTouchlessModal {
		t.Error("showTouchlessModal should be true")
	}

	// Test touchless options
	if model.touchlessMultiSub != 0 {
		t.Errorf("touchlessMultiSub = %d, want 0", model.touchlessMultiSub)
	}
}

// TestProfileManagement tests profile management fields
func TestProfileManagement(t *testing.T) {
	model := Model{
		profileKeys:        []string{"default", "custom"},
		selectedProfile:    0,
		editingPrompt:      false,
		editingProfileName: false,
		showProfileList:    true,
	}

	if len(model.profileKeys) != 2 {
		t.Errorf("len(profileKeys) = %d, want 2", len(model.profileKeys))
	}

	if model.selectedProfile != 0 {
		t.Errorf("selectedProfile = %d, want 0", model.selectedProfile)
	}

	if !model.showProfileList {
		t.Error("showProfileList should be true")
	}
}

// TestErrorState tests error state handling
func TestErrorState(t *testing.T) {
	model := Model{
		hasError: true,
		errMsg:   "test error message",
	}

	if !model.hasError {
		t.Error("hasError should be true")
	}

	if model.errMsg != "test error message" {
		t.Errorf("errMsg = %q, want test error message", model.errMsg)
	}
}
