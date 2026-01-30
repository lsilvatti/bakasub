package job

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/config"
)

// TestViewStateConstants tests ViewState constants
func TestViewStateConstants(t *testing.T) {
	if ViewMain != 0 {
		t.Errorf("ViewMain = %d, want 0", ViewMain)
	}

	if ViewDirectoryDetected != 1 {
		t.Errorf("ViewDirectoryDetected = %d, want 1", ViewDirectoryDetected)
	}

	if ViewConflictResolution != 2 {
		t.Errorf("ViewConflictResolution = %d, want 2", ViewConflictResolution)
	}

	if ViewDryRunReport != 3 {
		t.Errorf("ViewDryRunReport = %d, want 3", ViewDryRunReport)
	}

	if ViewGlossaryEditor != 4 {
		t.Errorf("ViewGlossaryEditor = %d, want 4", ViewGlossaryEditor)
	}
}

// TestDefaultKeyMap tests DefaultKeyMap function
func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	// Verify key bindings are set
	if len(km.Enter.Keys()) == 0 {
		t.Error("Enter key binding should have keys")
	}

	if len(km.Escape.Keys()) == 0 {
		t.Error("Escape key binding should have keys")
	}

	if len(km.DryRun.Keys()) == 0 {
		t.Error("DryRun key binding should have keys")
	}
}

// TestNew tests New function
func TestNew(t *testing.T) {
	cfg := config.Default()
	model := New(cfg, "/test/video.mkv")

	if model.cfg == nil {
		t.Error("cfg should not be nil")
	}

	if model.state != ViewMain {
		t.Errorf("state = %d, want ViewMain", model.state)
	}

	if model.jobConfig.InputPath != "/test/video.mkv" {
		t.Errorf("InputPath = %q, want /test/video.mkv", model.jobConfig.InputPath)
	}
}

// TestJobConfigDefaults tests JobConfig default values
func TestJobConfigDefaults(t *testing.T) {
	cfg := config.Default()
	model := New(cfg, "/test/video.mkv")

	jc := model.jobConfig

	if jc.RemoveHITags != true {
		t.Error("RemoveHITags should default to true")
	}

	if jc.SetDefault != true {
		t.Error("SetDefault should default to true")
	}

	if jc.BackupOriginal != true {
		t.Error("BackupOriginal should default to true")
	}

	if jc.ExtractFonts != true {
		t.Error("ExtractFonts should default to true")
	}

	if jc.AutoDetectTrack != true {
		t.Error("AutoDetectTrack should default to true")
	}

	if jc.GlossaryTerms == nil {
		t.Error("GlossaryTerms should not be nil")
	}
}

// TestMsgDirectoryDetected tests MsgDirectoryDetected structure
func TestMsgDirectoryDetected(t *testing.T) {
	msg := MsgDirectoryDetected{
		Path:     "/test/dir",
		MKVCount: 5,
		IsDir:    true,
	}

	if msg.Path != "/test/dir" {
		t.Errorf("Path = %q, want /test/dir", msg.Path)
	}

	if msg.MKVCount != 5 {
		t.Errorf("MKVCount = %d, want 5", msg.MKVCount)
	}

	if !msg.IsDir {
		t.Error("IsDir should be true")
	}
}

// TestMsgBatchModeSelected tests MsgBatchModeSelected structure
func TestMsgBatchModeSelected(t *testing.T) {
	msg := MsgBatchModeSelected{BatchMode: true}

	if !msg.BatchMode {
		t.Error("BatchMode should be true")
	}
}

// TestMsgAnalysisComplete tests MsgAnalysisComplete structure
func TestMsgAnalysisComplete(t *testing.T) {
	msg := MsgAnalysisComplete{
		Files:   []AnalyzedFile{},
		Success: true,
		Error:   nil,
	}

	if !msg.Success {
		t.Error("Success should be true")
	}

	if msg.Error != nil {
		t.Error("Error should be nil")
	}
}

// TestMsgConflictResolved tests MsgConflictResolved structure
func TestMsgConflictResolved(t *testing.T) {
	msg := MsgConflictResolved{
		FileIndex: 2,
		TrackID:   5,
	}

	if msg.FileIndex != 2 {
		t.Errorf("FileIndex = %d, want 2", msg.FileIndex)
	}

	if msg.TrackID != 5 {
		t.Errorf("TrackID = %d, want 5", msg.TrackID)
	}
}

// TestMsgCostEstimated tests MsgCostEstimated structure
func TestMsgCostEstimated(t *testing.T) {
	msg := MsgCostEstimated{
		TotalChars:    10000,
		EstimatedCost: 0.05,
		TokenCount:    2500,
	}

	if msg.TotalChars != 10000 {
		t.Errorf("TotalChars = %d, want 10000", msg.TotalChars)
	}

	if msg.EstimatedCost != 0.05 {
		t.Errorf("EstimatedCost = %f, want 0.05", msg.EstimatedCost)
	}
}

// TestMsgDryRunComplete tests MsgDryRunComplete structure
func TestMsgDryRunComplete(t *testing.T) {
	msg := MsgDryRunComplete{
		CanWrite:      true,
		TokenCount:    5000,
		EstimatedCost: 0.10,
		Warnings:      []string{"Warning 1", "Warning 2"},
	}

	if !msg.CanWrite {
		t.Error("CanWrite should be true")
	}

	if len(msg.Warnings) != 2 {
		t.Errorf("len(Warnings) = %d, want 2", len(msg.Warnings))
	}
}

// TestStartJobMsg tests StartJobMsg structure
func TestStartJobMsg(t *testing.T) {
	msg := StartJobMsg{
		JobConfig: JobConfig{
			InputPath:  "/test/video.mkv",
			TargetLang: "pt-br",
		},
	}

	if msg.JobConfig.InputPath != "/test/video.mkv" {
		t.Errorf("InputPath = %q, want /test/video.mkv", msg.JobConfig.InputPath)
	}
}

// TestModelStateTransitions tests state transitions
func TestModelStateTransitions(t *testing.T) {
	cfg := config.Default()
	model := New(cfg, "/test/video.mkv")

	// Initial state
	if model.state != ViewMain {
		t.Errorf("Initial state = %d, want ViewMain", model.state)
	}

	// Simulate state change
	model.state = ViewConflictResolution
	if model.state != ViewConflictResolution {
		t.Error("State should be ViewConflictResolution")
	}
}
