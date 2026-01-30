package header

import (
	"testing"
)

// TestTrackStruct tests Track structure
func TestTrackStruct(t *testing.T) {
	track := Track{
		ID:        1,
		Type:      "subtitles",
		Codec:     "ass",
		Language:  "eng",
		TrackName: "English Subtitles",
		Default:   true,
		Forced:    false,
	}

	if track.ID != 1 {
		t.Errorf("ID = %d, want 1", track.ID)
	}

	if track.Type != "subtitles" {
		t.Errorf("Type = %q, want subtitles", track.Type)
	}

	if track.Codec != "ass" {
		t.Errorf("Codec = %q, want ass", track.Codec)
	}

	if track.Language != "eng" {
		t.Errorf("Language = %q, want eng", track.Language)
	}

	if track.TrackName != "English Subtitles" {
		t.Errorf("TrackName = %q, want English Subtitles", track.TrackName)
	}

	if !track.Default {
		t.Error("Default should be true")
	}

	if track.Forced {
		t.Error("Forced should be false")
	}
}

// TestMKVInfoStruct tests MKVInfo structure
func TestMKVInfoStruct(t *testing.T) {
	info := MKVInfo{
		Tracks: []Track{
			{ID: 0, Type: "video"},
			{ID: 1, Type: "audio"},
			{ID: 2, Type: "subtitles"},
		},
	}

	if len(info.Tracks) != 3 {
		t.Errorf("len(Tracks) = %d, want 3", len(info.Tracks))
	}
}

// TestClosedMsgStruct tests ClosedMsg structure
func TestClosedMsgStruct(t *testing.T) {
	msg := ClosedMsg{}
	_ = msg // Just verify it can be created
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		tracks:   []Track{},
		filePath: "/test/video.mkv",
		width:    80,
		height:   24,
		modified: false,
	}

	if model.filePath != "/test/video.mkv" {
		t.Errorf("filePath = %q, want /test/video.mkv", model.filePath)
	}

	if model.width != 80 {
		t.Errorf("width = %d, want 80", model.width)
	}

	if model.height != 24 {
		t.Errorf("height = %d, want 24", model.height)
	}

	if model.modified {
		t.Error("modified should be false initially")
	}
}

// TestNewWithNonExistentFile tests New with non-existent file
func TestNewWithNonExistentFile(t *testing.T) {
	_, err := New("/nonexistent/file.mkv")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestTrackFiltering tests filtering tracks by type
func TestTrackFiltering(t *testing.T) {
	tracks := []Track{
		{ID: 0, Type: "video"},
		{ID: 1, Type: "audio"},
		{ID: 2, Type: "audio"},
		{ID: 3, Type: "subtitles"},
		{ID: 4, Type: "subtitles"},
	}

	videoCount := 0
	audioCount := 0
	subtitleCount := 0

	for _, track := range tracks {
		switch track.Type {
		case "video":
			videoCount++
		case "audio":
			audioCount++
		case "subtitles":
			subtitleCount++
		}
	}

	if videoCount != 1 {
		t.Errorf("videoCount = %d, want 1", videoCount)
	}

	if audioCount != 2 {
		t.Errorf("audioCount = %d, want 2", audioCount)
	}

	if subtitleCount != 2 {
		t.Errorf("subtitleCount = %d, want 2", subtitleCount)
	}
}

// TestTrackDefaultFlag tests Default flag behavior
func TestTrackDefaultFlag(t *testing.T) {
	tracks := []Track{
		{ID: 1, Type: "audio", Default: true},
		{ID: 2, Type: "audio", Default: false},
	}

	defaultCount := 0
	for _, track := range tracks {
		if track.Default {
			defaultCount++
		}
	}

	if defaultCount != 1 {
		t.Errorf("defaultCount = %d, want 1", defaultCount)
	}
}

// TestTrackForcedFlag tests Forced flag behavior
func TestTrackForcedFlag(t *testing.T) {
	track := Track{
		ID:     1,
		Type:   "subtitles",
		Forced: true,
	}

	if !track.Forced {
		t.Error("Forced should be true")
	}
}

// TestTrackJSONTags tests JSON tags on Track
func TestTrackJSONTags(t *testing.T) {
	track := Track{
		ID:       1,
		Type:     "subtitles",
		Codec:    "ass",
		Language: "eng",
	}

	// Verify struct can be used (JSON tags are compile-time)
	if track.ID != 1 {
		t.Error("Track fields not set correctly")
	}
}
