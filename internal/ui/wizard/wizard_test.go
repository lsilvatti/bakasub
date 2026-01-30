package wizard

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/config"
)

// TestStepConstants tests Step constants
func TestStepConstants(t *testing.T) {
	if StepLanguageDeps != 0 {
		t.Errorf("StepLanguageDeps = %d, want 0", StepLanguageDeps)
	}

	if StepProvider != 1 {
		t.Errorf("StepProvider = %d, want 1", StepProvider)
	}

	if StepDefaults != 2 {
		t.Errorf("StepDefaults = %d, want 2", StepDefaults)
	}
}

// TestIsoLanguageNames tests isoLanguageNames map
func TestIsoLanguageNames(t *testing.T) {
	testCases := []struct {
		code     string
		expected string
	}{
		{"en", "English"},
		{"pt-br", "Português (Brasil)"},
		{"es", "Español"},
		{"ja", "日本語 (Japanese)"},
		{"fr", "Français"},
		{"de", "Deutsch"},
	}

	for _, tc := range testCases {
		t.Run(tc.code, func(t *testing.T) {
			name, ok := isoLanguageNames[tc.code]
			if !ok {
				t.Errorf("isoLanguageNames missing code %q", tc.code)
				return
			}
			if name != tc.expected {
				t.Errorf("isoLanguageNames[%q] = %q, want %q", tc.code, name, tc.expected)
			}
		})
	}
}

// TestIsoLanguageNamesComplete tests isoLanguageNames has expected entries
func TestIsoLanguageNamesComplete(t *testing.T) {
	requiredCodes := []string{
		"en", "pt-br", "es", "fr", "de", "it", "ja", "ko", "zh", "ru",
	}

	for _, code := range requiredCodes {
		if _, ok := isoLanguageNames[code]; !ok {
			t.Errorf("isoLanguageNames missing required code %q", code)
		}
	}
}

// TestModelInfoStruct tests ModelInfo structure
func TestModelInfoStruct(t *testing.T) {
	info := ModelInfo{
		ID:          "anthropic/claude-3-opus",
		Name:        "Claude 3 Opus",
		ContextSize: "200K",
		PricePerM:   "$15.00",
		IsFree:      false,
	}

	if info.ID != "anthropic/claude-3-opus" {
		t.Errorf("ID = %q, want anthropic/claude-3-opus", info.ID)
	}

	if info.Name != "Claude 3 Opus" {
		t.Errorf("Name = %q, want Claude 3 Opus", info.Name)
	}

	if info.ContextSize != "200K" {
		t.Errorf("ContextSize = %q, want 200K", info.ContextSize)
	}

	if info.IsFree {
		t.Error("IsFree should be false")
	}
}

// TestModelInfoFree tests free ModelInfo
func TestModelInfoFree(t *testing.T) {
	info := ModelInfo{
		ID:          "google/gemini-pro",
		Name:        "Gemini Pro",
		ContextSize: "128K",
		PricePerM:   "Free",
		IsFree:      true,
	}

	if !info.IsFree {
		t.Error("IsFree should be true")
	}
}

// TestNewWithNilConfig tests New with nil config
func TestNewWithNilConfig(t *testing.T) {
	model := New(nil)

	// Should not panic
	if model.step != StepLanguageDeps {
		t.Errorf("step = %d, want StepLanguageDeps", model.step)
	}
}

// TestNewWithValidConfig tests New with valid config
func TestNewWithValidConfig(t *testing.T) {
	cfg := config.Default()
	model := New(cfg)

	if model.config == nil {
		t.Error("config should not be nil")
	}

	if model.step != StepLanguageDeps {
		t.Errorf("step = %d, want StepLanguageDeps", model.step)
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		step:     StepProvider,
		width:    120,
		height:   40,
		quitting: false,
		finished: false,
	}

	if model.step != StepProvider {
		t.Errorf("step = %d, want StepProvider", model.step)
	}

	if model.width != 120 {
		t.Errorf("width = %d, want 120", model.width)
	}

	if model.height != 40 {
		t.Errorf("height = %d, want 40", model.height)
	}

	if model.quitting {
		t.Error("quitting should be false")
	}

	if model.finished {
		t.Error("finished should be false")
	}
}

// TestUILanguageSelection tests UI language selection values
func TestUILanguageSelection(t *testing.T) {
	languages := []struct {
		idx  int
		name string
	}{
		{0, "English"},
		{1, "PT-BR"},
		{2, "Español"},
	}

	for _, l := range languages {
		model := Model{
			uiLanguageSelection: l.idx,
		}

		if model.uiLanguageSelection != l.idx {
			t.Errorf("uiLanguageSelection = %d for %s, want %d", model.uiLanguageSelection, l.name, l.idx)
		}
	}
}

// TestProviderSelection tests provider selection values
func TestProviderSelection(t *testing.T) {
	providers := []int{0, 1, 2, 3}

	for _, idx := range providers {
		model := Model{
			providerSelection: idx,
		}

		if model.providerSelection != idx {
			t.Errorf("providerSelection = %d, want %d", model.providerSelection, idx)
		}
	}
}

// TestDependencyState tests dependency check state
func TestDependencyState(t *testing.T) {
	model := Model{
		depStatus: map[string]bool{
			"ffmpeg":     true,
			"mkvmerge":   true,
			"mkvextract": false,
		},
		depDownloading:   false,
		downloadProgress: 0.0,
		checkComplete:    false,
	}

	if !model.depStatus["ffmpeg"] {
		t.Error("ffmpeg should be true")
	}

	if !model.depStatus["mkvmerge"] {
		t.Error("mkvmerge should be true")
	}

	if model.depStatus["mkvextract"] {
		t.Error("mkvextract should be false")
	}
}

// TestKeyValidation tests key validation state
func TestKeyValidation(t *testing.T) {
	model := Model{
		validatingKey:    true,
		keyValidated:     false,
		keyValidationErr: "",
	}

	if !model.validatingKey {
		t.Error("validatingKey should be true")
	}

	if model.keyValidated {
		t.Error("keyValidated should be false")
	}
}

// TestKeyValidationError tests key validation error state
func TestKeyValidationError(t *testing.T) {
	model := Model{
		validatingKey:    false,
		keyValidated:     false,
		keyValidationErr: "Invalid API key",
	}

	if model.keyValidationErr != "Invalid API key" {
		t.Errorf("keyValidationErr = %q, want Invalid API key", model.keyValidationErr)
	}
}

// TestStep3State tests step 3 state
func TestStep3State(t *testing.T) {
	model := Model{
		activeStep3Section: 0,
		languageSelection:  0,
		targetLangOther:    false,
		hiTagsRemoval:      true,
		tempValue:          0.7,
	}

	if model.targetLangOther {
		t.Error("targetLangOther should be false")
	}

	if !model.hiTagsRemoval {
		t.Error("hiTagsRemoval should be true")
	}

	if model.tempValue != 0.7 {
		t.Errorf("tempValue = %f, want 0.7", model.tempValue)
	}
}

// TestStepProgression tests step progression
func TestStepProgression(t *testing.T) {
	if StepProvider <= StepLanguageDeps {
		t.Error("StepProvider should be greater than StepLanguageDeps")
	}

	if StepDefaults <= StepProvider {
		t.Error("StepDefaults should be greater than StepProvider")
	}
}
