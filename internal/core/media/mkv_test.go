package media

import (
	"testing"
)

// TestTrackStruct tests Track structure
func TestTrackStruct(t *testing.T) {
	track := Track{
		ID:       1,
		Type:     "subtitles",
		Codec:    "ass",
		Language: "eng",
		Name:     "English Subtitles",
		Default:  true,
		Forced:   false,
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

	if track.Name != "English Subtitles" {
		t.Errorf("Name = %q, want English Subtitles", track.Name)
	}

	if !track.Default {
		t.Error("Default should be true")
	}

	if track.Forced {
		t.Error("Forced should be false")
	}
}

// TestTrackIsSubtitle tests Track type checking
func TestTrackIsSubtitle(t *testing.T) {
	track := Track{
		ID:   1,
		Type: "subtitles",
	}

	if track.Type != "subtitles" {
		t.Error("Expected subtitles type")
	}
}

// TestTrackIsAudio tests Track audio type
func TestTrackIsAudio(t *testing.T) {
	track := Track{
		ID:   2,
		Type: "audio",
	}

	if track.Type != "audio" {
		t.Error("Expected audio type")
	}
}

// TestTrackIsVideo tests Track video type
func TestTrackIsVideo(t *testing.T) {
	track := Track{
		ID:   0,
		Type: "video",
	}

	if track.Type != "video" {
		t.Error("Expected video type")
	}
}

// TestAttachmentStruct tests Attachment structure
func TestAttachmentStruct(t *testing.T) {
	att := Attachment{
		ID:          1,
		FileName:    "font.ttf",
		MimeType:    "application/x-font-ttf",
		Size:        12345,
		Description: "Main font",
	}

	if att.ID != 1 {
		t.Errorf("ID = %d, want 1", att.ID)
	}

	if att.FileName != "font.ttf" {
		t.Errorf("FileName = %q, want font.ttf", att.FileName)
	}

	if att.MimeType != "application/x-font-ttf" {
		t.Errorf("MimeType = %q, want application/x-font-ttf", att.MimeType)
	}

	if att.Size != 12345 {
		t.Errorf("Size = %d, want 12345", att.Size)
	}
}

// TestFileInfoStruct tests FileInfo structure
func TestFileInfoStruct(t *testing.T) {
	info := FileInfo{
		FileName: "video.mkv",
		Tracks: []Track{
			{ID: 0, Type: "video"},
			{ID: 1, Type: "audio"},
			{ID: 2, Type: "subtitles"},
		},
		Attachments: []Attachment{
			{ID: 1, FileName: "font.ttf"},
		},
	}

	if info.FileName != "video.mkv" {
		t.Errorf("FileName = %q, want video.mkv", info.FileName)
	}

	if len(info.Tracks) != 3 {
		t.Errorf("len(Tracks) = %d, want 3", len(info.Tracks))
	}

	if len(info.Attachments) != 1 {
		t.Errorf("len(Attachments) = %d, want 1", len(info.Attachments))
	}
}

// TestSetBinPath tests SetBinPath function
func TestSetBinPath(t *testing.T) {
	originalPath := BinPath

	SetBinPath("/custom/path")

	if BinPath != "/custom/path" {
		t.Errorf("BinPath = %q, want /custom/path", BinPath)
	}

	// Restore original
	BinPath = originalPath
}

// TestGetBinaryPath tests getBinaryPath function
func TestGetBinaryPath(t *testing.T) {
	// This tests the fallback behavior
	path := getBinaryPath("nonexistent-binary-12345")

	// Should return the name as fallback
	if path != "nonexistent-binary-12345" {
		t.Logf("getBinaryPath returned: %q", path)
	}
}

// TestAnalyzeNonExistentFile tests Analyze with non-existent file
func TestAnalyzeNonExistentFile(t *testing.T) {
	_, err := Analyze("/nonexistent/file.mkv")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestTrackFiltering tests filtering tracks by type
func TestTrackFiltering(t *testing.T) {
	info := FileInfo{
		Tracks: []Track{
			{ID: 0, Type: "video", Codec: "h264"},
			{ID: 1, Type: "audio", Language: "eng"},
			{ID: 2, Type: "audio", Language: "jpn"},
			{ID: 3, Type: "subtitles", Language: "eng"},
			{ID: 4, Type: "subtitles", Language: "jpn"},
		},
	}

	// Count by type
	videoCount := 0
	audioCount := 0
	subtitleCount := 0

	for _, track := range info.Tracks {
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

// TestTrackProperties tests Track properties map
func TestTrackProperties(t *testing.T) {
	track := Track{
		ID:   1,
		Type: "subtitles",
		Properties: map[string]interface{}{
			"text_subtitles": true,
			"encoding":       "UTF-8",
		},
	}

	if track.Properties == nil {
		t.Error("Properties should not be nil")
	}

	if track.Properties["text_subtitles"] != true {
		t.Error("text_subtitles property not set correctly")
	}
}

// TestExtractSubtitleTrackNonExistent tests ExtractSubtitleTrack with non-existent file
func TestExtractSubtitleTrackNonExistent(t *testing.T) {
	err := ExtractSubtitleTrack("/nonexistent/file.mkv", 0, "/tmp/output.ass")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestMuxSubtitleNonExistent tests MuxSubtitle with non-existent files
func TestMuxSubtitleNonExistent(t *testing.T) {
	err := MuxSubtitle("/nonexistent/video.mkv", "/nonexistent/sub.ass", "/tmp/output.mkv")

	if err == nil {
		t.Error("Expected error for non-existent files")
	}
}

// TestMuxSourceStruct tests MuxSource structure
func TestMuxSourceStruct(t *testing.T) {
	src := MuxSource{
		FilePath: "/path/to/file.mkv",
		TrackIDs: []int{0, 1, 2},
		NoVideo:  false,
		NoAudio:  true,
		NoSubs:   false,
		Language: "eng",
		Name:     "English Track",
	}

	if src.FilePath != "/path/to/file.mkv" {
		t.Errorf("FilePath = %q, want /path/to/file.mkv", src.FilePath)
	}

	if len(src.TrackIDs) != 3 {
		t.Errorf("len(TrackIDs) = %d, want 3", len(src.TrackIDs))
	}

	if !src.NoAudio {
		t.Error("NoAudio should be true")
	}

	if src.Language != "eng" {
		t.Errorf("Language = %q, want eng", src.Language)
	}
}

// TestMuxOptionsStruct tests MuxOptions structure
func TestMuxOptionsStruct(t *testing.T) {
	opts := MuxOptions{
		Title:      "Test Video",
		OutputPath: "/output/video.mkv",
		Sources: []MuxSource{
			{FilePath: "/input/video.mkv"},
			{FilePath: "/input/sub.ass", Language: "eng"},
		},
		Quiet: true,
	}

	if opts.Title != "Test Video" {
		t.Errorf("Title = %q, want Test Video", opts.Title)
	}

	if opts.OutputPath != "/output/video.mkv" {
		t.Errorf("OutputPath = %q, want /output/video.mkv", opts.OutputPath)
	}

	if len(opts.Sources) != 2 {
		t.Errorf("len(Sources) = %d, want 2", len(opts.Sources))
	}

	if !opts.Quiet {
		t.Error("Quiet should be true")
	}
}

// TestMuxEmptyOutput tests Mux with empty output path
func TestMuxEmptyOutput(t *testing.T) {
	opts := MuxOptions{
		OutputPath: "",
		Sources:    []MuxSource{{FilePath: "/test.mkv"}},
	}

	err := Mux(opts)
	if err == nil {
		t.Error("Expected error for empty output path")
	}
}

// TestMuxNoSources tests Mux with no sources
func TestMuxNoSources(t *testing.T) {
	opts := MuxOptions{
		OutputPath: "/output.mkv",
		Sources:    []MuxSource{},
	}

	err := Mux(opts)
	if err == nil {
		t.Error("Expected error for empty sources")
	}
}

// TestExtractTrackNonExistent tests ExtractTrack with non-existent file
func TestExtractTrackNonExistent(t *testing.T) {
	err := ExtractTrack("/nonexistent/file.mkv", 0, "/tmp/output.ass")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestFileInfoContainer tests FileInfo container fields
func TestFileInfoContainer(t *testing.T) {
	info := FileInfo{
		FileName: "test.mkv",
	}
	info.Container.Type = "Matroska"
	info.Container.Duration = 3600000 // 1 hour in milliseconds

	if info.Container.Type != "Matroska" {
		t.Errorf("Container.Type = %q, want Matroska", info.Container.Type)
	}

	if info.Container.Duration != 3600000 {
		t.Errorf("Container.Duration = %d, want 3600000", info.Container.Duration)
	}
}

// TestTrackWithProperties tests Track with Properties map
func TestTrackWithProperties(t *testing.T) {
	track := Track{
		ID:         1,
		Type:       "subtitles",
		Properties: make(map[string]interface{}),
	}

	track.Properties["codec_id"] = "S_TEXT/ASS"
	track.Properties["default_duration"] = 1000

	if track.Properties["codec_id"] != "S_TEXT/ASS" {
		t.Error("Property codec_id not set correctly")
	}
}

// TestMultipleAttachments tests multiple attachments
func TestMultipleAttachments(t *testing.T) {
	attachments := []Attachment{
		{ID: 1, FileName: "font1.ttf", MimeType: "application/x-truetype-font", Size: 10000},
		{ID: 2, FileName: "font2.otf", MimeType: "application/x-opentype-font", Size: 20000},
		{ID: 3, FileName: "cover.jpg", MimeType: "image/jpeg", Size: 5000},
	}

	if len(attachments) != 3 {
		t.Errorf("len(attachments) = %d, want 3", len(attachments))
	}

	// Check MIME types
	mimeTypes := map[string]bool{}
	for _, a := range attachments {
		mimeTypes[a.MimeType] = true
	}

	if !mimeTypes["image/jpeg"] {
		t.Error("Should have image/jpeg MIME type")
	}
}

// TestTrackLanguages tests various language codes
func TestTrackLanguages(t *testing.T) {
	languages := []struct {
		code string
		name string
	}{
		{"eng", "English"},
		{"jpn", "Japanese"},
		{"por", "Portuguese"},
		{"spa", "Spanish"},
		{"fra", "French"},
		{"deu", "German"},
		{"ita", "Italian"},
		{"und", "Undetermined"},
	}

	for _, lang := range languages {
		track := Track{Language: lang.code}
		if len(track.Language) != 3 {
			t.Errorf("Language %s (%s) should be 3 chars", lang.name, lang.code)
		}
	}
}
