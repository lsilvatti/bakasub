package remuxer

import (
	"testing"
)

// TestStateConstants tests State constants
func TestStateConstants(t *testing.T) {
	if StateSelect != 0 {
		t.Errorf("StateSelect = %d, want 0", StateSelect)
	}

	if StateConfirm != 1 {
		t.Errorf("StateConfirm = %d, want 1", StateConfirm)
	}

	if StateAddExternal != 2 {
		t.Errorf("StateAddExternal = %d, want 2", StateAddExternal)
	}

	if StateProcessing != 3 {
		t.Errorf("StateProcessing = %d, want 3", StateProcessing)
	}

	if StateDone != 4 {
		t.Errorf("StateDone = %d, want 4", StateDone)
	}
}

// TestTrackStruct tests Track structure
func TestTrackStruct(t *testing.T) {
	track := Track{
		ID:       1,
		Type:     "subtitles",
		Codec:    "SubRip/SRT",
		Language: "eng",
		Name:     "English",
		Selected: true,
	}

	if track.ID != 1 {
		t.Errorf("ID = %d, want 1", track.ID)
	}

	if track.Type != "subtitles" {
		t.Errorf("Type = %q, want subtitles", track.Type)
	}

	if track.Codec != "SubRip/SRT" {
		t.Errorf("Codec = %q, want SubRip/SRT", track.Codec)
	}

	if track.Language != "eng" {
		t.Errorf("Language = %q, want eng", track.Language)
	}

	if track.Name != "English" {
		t.Errorf("Name = %q, want English", track.Name)
	}

	if !track.Selected {
		t.Error("Selected should be true")
	}
}

// TestExternalTrackStruct tests ExternalTrack structure
func TestExternalTrackStruct(t *testing.T) {
	ext := ExternalTrack{
		Path:     "/path/to/sub.srt",
		Type:     "subtitle",
		Language: "por",
	}

	if ext.Path != "/path/to/sub.srt" {
		t.Errorf("Path = %q, want /path/to/sub.srt", ext.Path)
	}

	if ext.Type != "subtitle" {
		t.Errorf("Type = %q, want subtitle", ext.Type)
	}

	if ext.Language != "por" {
		t.Errorf("Language = %q, want por", ext.Language)
	}
}

// TestClosedMsgStruct tests ClosedMsg structure
func TestClosedMsgStruct(t *testing.T) {
	msg := ClosedMsg{}

	// ClosedMsg is an empty struct used for signaling
	_ = msg
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		state:      StateSelect,
		filePath:   "/input.mkv",
		outputPath: "/output.mkv",
		cursor:     0,
	}

	if model.state != StateSelect {
		t.Errorf("state = %d, want StateSelect", model.state)
	}

	if model.filePath != "/input.mkv" {
		t.Errorf("filePath = %q, want /input.mkv", model.filePath)
	}

	if model.outputPath != "/output.mkv" {
		t.Errorf("outputPath = %q, want /output.mkv", model.outputPath)
	}

	if model.cursor != 0 {
		t.Errorf("cursor = %d, want 0", model.cursor)
	}
}

// TestNewWithEmptyPath tests New with empty path
func TestNewWithEmptyPath(t *testing.T) {
	_, err := New("")

	if err == nil {
		t.Error("expected error for empty path")
	}
}

// TestNewWithNonExistentFile tests New with non-existent file
func TestNewWithNonExistentFile(t *testing.T) {
	_, err := New("/nonexistent/file.mkv")

	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

// TestTrackTypeValues tests common track type values
func TestTrackTypeValues(t *testing.T) {
	trackTypes := []string{"video", "audio", "subtitles"}

	for _, tt := range trackTypes {
		track := Track{Type: tt}
		if track.Type == "" {
			t.Errorf("Track type should not be empty for %q", tt)
		}
	}
}

// TestTrackLanguageCodes tests ISO language codes
func TestTrackLanguageCodes(t *testing.T) {
	langCodes := []string{"eng", "por", "spa", "jpn", "und"}

	for _, code := range langCodes {
		track := Track{Language: code}
		if len(track.Language) != 3 {
			t.Errorf("Language code %q should be 3 characters", code)
		}
	}
}

// TestExternalTrackEmpty tests empty ExternalTrack
func TestExternalTrackEmpty(t *testing.T) {
	ext := ExternalTrack{}

	if ext.Path != "" {
		t.Error("Path should be empty")
	}

	if ext.Type != "" {
		t.Error("Type should be empty")
	}

	if ext.Language != "" {
		t.Error("Language should be empty")
	}
}

// TestMultipleTracks tests slice of tracks
func TestMultipleTracks(t *testing.T) {
	tracks := []Track{
		{ID: 0, Type: "video", Codec: "HEVC", Selected: true},
		{ID: 1, Type: "audio", Codec: "AAC", Language: "jpn", Selected: true},
		{ID: 2, Type: "subtitles", Codec: "ASS", Language: "eng", Selected: false},
	}

	if len(tracks) != 3 {
		t.Errorf("len(tracks) = %d, want 3", len(tracks))
	}

	// Video track should be first
	if tracks[0].Type != "video" {
		t.Error("First track should be video")
	}

	// Only 2 should be selected
	selected := 0
	for _, tr := range tracks {
		if tr.Selected {
			selected++
		}
	}
	if selected != 2 {
		t.Errorf("selected = %d, want 2", selected)
	}
}

// TestStateProgression tests state progression values
func TestStateProgression(t *testing.T) {
	// States should increase sequentially
	if StateConfirm <= StateSelect {
		t.Error("StateConfirm should be greater than StateSelect")
	}

	if StateAddExternal <= StateConfirm {
		t.Error("StateAddExternal should be greater than StateConfirm")
	}

	if StateProcessing <= StateAddExternal {
		t.Error("StateProcessing should be greater than StateAddExternal")
	}

	if StateDone <= StateProcessing {
		t.Error("StateDone should be greater than StateProcessing")
	}
}
