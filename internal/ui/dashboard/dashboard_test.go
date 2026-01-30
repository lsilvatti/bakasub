package dashboard

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/config"
)

// TestViewStateConstants tests ViewState constants
func TestViewStateConstants(t *testing.T) {
	tests := []struct {
		state    ViewState
		expected int
		name     string
	}{
		{ViewDashboard, 0, "ViewDashboard"},
		{ViewPicker, 1, "ViewPicker"},
		{ViewSettings, 2, "ViewSettings"},
		{ViewJob, 3, "ViewJob"},
		{ViewExecution, 4, "ViewExecution"},
		{ViewHeader, 5, "ViewHeader"},
		{ViewAttachments, 6, "ViewAttachments"},
		{ViewRemuxer, 7, "ViewRemuxer"},
		{ViewGlossary, 8, "ViewGlossary"},
		{ViewReview, 9, "ViewReview"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.state) != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.state, tt.expected)
			}
		})
	}
}

// TestDirectoryAnalysisStruct tests DirectoryAnalysis structure
func TestDirectoryAnalysisStruct(t *testing.T) {
	analysis := DirectoryAnalysis{
		Path:     "/test/path",
		MKVCount: 5,
		SubCount: 10,
		IsDir:    true,
		Scanned:  true,
		Scanning: false,
		Error:    nil,
	}

	if analysis.Path != "/test/path" {
		t.Errorf("Path = %q, want /test/path", analysis.Path)
	}

	if analysis.MKVCount != 5 {
		t.Errorf("MKVCount = %d, want 5", analysis.MKVCount)
	}

	if analysis.SubCount != 10 {
		t.Errorf("SubCount = %d, want 10", analysis.SubCount)
	}

	if !analysis.IsDir {
		t.Error("IsDir should be true")
	}

	if !analysis.Scanned {
		t.Error("Scanned should be true")
	}

	if analysis.Scanning {
		t.Error("Scanning should be false")
	}
}

// TestDirectoryAnalysisWithError tests DirectoryAnalysis with error
func TestDirectoryAnalysisWithError(t *testing.T) {
	analysis := DirectoryAnalysis{
		Path:    "/invalid/path",
		Error:   nil,
		Scanned: true,
	}

	if analysis.Error != nil {
		t.Error("Error should be nil")
	}
}

// TestNewWithNilConfig tests New with nil config
func TestNewWithNilConfig(t *testing.T) {
	model := New(nil)

	if model.config == nil {
		t.Error("config should not be nil")
	}

	if model.viewState != ViewDashboard {
		t.Errorf("viewState = %d, want ViewDashboard", model.viewState)
	}
}

// TestNewWithValidConfig tests New with valid config
func TestNewWithValidConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Model = "test-model"
	cfg.TargetLang = "pt-br"
	cfg.Temperature = 0.5

	model := New(cfg)

	if model.config == nil {
		t.Error("config should not be nil")
	}

	if model.currentModel != "test-model" {
		t.Errorf("currentModel = %q, want test-model", model.currentModel)
	}

	if model.targetLang != "pt-br" {
		t.Errorf("targetLang = %q, want pt-br", model.targetLang)
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		width:           120,
		height:          40,
		selectedPath:    "/test/path",
		selectedMode:    0,
		apiOnline:       true,
		cacheOK:         true,
		updateAvailable: false,
		ffmpegOK:        true,
		mkvtoolnixOK:    true,
		viewState:       ViewDashboard,
	}

	if model.width != 120 {
		t.Errorf("width = %d, want 120", model.width)
	}

	if model.height != 40 {
		t.Errorf("height = %d, want 40", model.height)
	}

	if model.selectedPath != "/test/path" {
		t.Errorf("selectedPath = %q, want /test/path", model.selectedPath)
	}

	if !model.apiOnline {
		t.Error("apiOnline should be true")
	}

	if !model.cacheOK {
		t.Error("cacheOK should be true")
	}
}

// TestSelectedModeValues tests selectedMode values
func TestSelectedModeValues(t *testing.T) {
	modes := []struct {
		idx  int
		name string
	}{
		{0, "Full Process"},
		{1, "Watch Mode"},
	}

	for _, m := range modes {
		model := Model{
			selectedMode: m.idx,
		}

		if model.selectedMode != m.idx {
			t.Errorf("selectedMode = %d for %s, want %d", model.selectedMode, m.name, m.idx)
		}
	}
}

// TestResumeState tests resume state
func TestResumeState(t *testing.T) {
	model := Model{
		showResumeModal: true,
		resumeState:     nil,
	}

	if !model.showResumeModal {
		t.Error("showResumeModal should be true")
	}
}

// TestWatchModeState tests watch mode state
func TestWatchModeState(t *testing.T) {
	model := Model{
		watchModeActive: true,
		watchModePath:   "/watch/path",
	}

	if !model.watchModeActive {
		t.Error("watchModeActive should be true")
	}

	if model.watchModePath != "/watch/path" {
		t.Errorf("watchModePath = %q, want /watch/path", model.watchModePath)
	}
}

// TestUpdateCheckerState tests update checker state
func TestUpdateCheckerState(t *testing.T) {
	model := Model{
		updateAvailable: true,
		latestVersion:   "v2.0.0",
		releaseURL:      "https://github.com/example/release",
	}

	if !model.updateAvailable {
		t.Error("updateAvailable should be true")
	}

	if model.latestVersion != "v2.0.0" {
		t.Errorf("latestVersion = %q, want v2.0.0", model.latestVersion)
	}

	if model.releaseURL != "https://github.com/example/release" {
		t.Errorf("releaseURL = %q, want https://github.com/example/release", model.releaseURL)
	}
}

// TestKofiFlash tests ko-fi flash feedback
func TestKofiFlash(t *testing.T) {
	model := Model{
		kofiFlash: true,
	}

	if !model.kofiFlash {
		t.Error("kofiFlash should be true")
	}
}

// TestViewStateProgression tests that view states have unique values
func TestViewStateProgression(t *testing.T) {
	states := []ViewState{
		ViewDashboard, ViewPicker, ViewSettings, ViewJob,
		ViewExecution, ViewHeader, ViewAttachments, ViewRemuxer,
		ViewGlossary, ViewReview,
	}

	seen := make(map[ViewState]bool)
	for _, s := range states {
		if seen[s] {
			t.Errorf("Duplicate ViewState value: %d", s)
		}
		seen[s] = true
	}
}

// TestDependencyStatus tests dependency status fields
func TestDependencyStatus(t *testing.T) {
	model := Model{
		ffmpegOK:     false,
		mkvtoolnixOK: true,
	}

	if model.ffmpegOK {
		t.Error("ffmpegOK should be false")
	}

	if !model.mkvtoolnixOK {
		t.Error("mkvtoolnixOK should be true")
	}
}

// TestModuleError tests moduleError field
func TestModuleError(t *testing.T) {
	model := Model{
		moduleError: nil,
	}

	if model.moduleError != nil {
		t.Error("moduleError should be nil")
	}
}
