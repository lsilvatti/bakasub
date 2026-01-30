package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSRT(t *testing.T) {
	srtContent := `1
00:00:01,000 --> 00:00:04,000
Hello, world!

2
00:00:05,000 --> 00:00:08,000
How are you?

3
00:00:10,000 --> 00:00:15,000
This is a test
with multiple lines.
`

	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")
	if err := os.WriteFile(srtPath, []byte(srtContent), 0644); err != nil {
		t.Fatal(err)
	}

	sf, err := ParseFile(srtPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if sf.Format != "srt" {
		t.Errorf("expected format 'srt', got %q", sf.Format)
	}

	if sf.LineCount != 3 {
		t.Errorf("expected 3 lines, got %d", sf.LineCount)
	}

	if sf.Lines[0].Text != "Hello, world!" {
		t.Errorf("unexpected first line text: %q", sf.Lines[0].Text)
	}

	if sf.Lines[0].StartTime != "00:00:01,000" {
		t.Errorf("unexpected start time: %q", sf.Lines[0].StartTime)
	}

	if sf.Lines[0].EndTime != "00:00:04,000" {
		t.Errorf("unexpected end time: %q", sf.Lines[0].EndTime)
	}

	// Test multi-line subtitle
	expectedMultiline := "This is a test\nwith multiple lines."
	if sf.Lines[2].Text != expectedMultiline {
		t.Errorf("expected multiline %q, got %q", expectedMultiline, sf.Lines[2].Text)
	}
}

func TestParseASS(t *testing.T) {
	assContent := `[Script Info]
Title: Test Subtitle
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,Arial,48,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,2,2,10,10,10,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:01.00,0:00:04.00,Default,,0000,0000,0000,,Hello, world!
Dialogue: 0,0:00:05.00,0:00:08.00,Default,,0000,0000,0000,,{\an8}Top positioned text
`

	tmpDir := t.TempDir()
	assPath := filepath.Join(tmpDir, "test.ass")
	if err := os.WriteFile(assPath, []byte(assContent), 0644); err != nil {
		t.Fatal(err)
	}

	sf, err := ParseFile(assPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if sf.Format != "ass" {
		t.Errorf("expected format 'ass', got %q", sf.Format)
	}

	if sf.LineCount != 2 {
		t.Errorf("expected 2 lines, got %d", sf.LineCount)
	}

	if sf.Lines[0].Text != "Hello, world!" {
		t.Errorf("unexpected first line text: %q", sf.Lines[0].Text)
	}

	if sf.Lines[0].Style != "Default" {
		t.Errorf("unexpected style: %q", sf.Lines[0].Style)
	}

	if sf.Lines[0].StartTime != "0:00:01.00" {
		t.Errorf("unexpected start time: %q", sf.Lines[0].StartTime)
	}

	// Test ASS tags preserved
	if !strings.Contains(sf.Lines[1].Text, "{\\an8}") {
		t.Error("ASS tags should be preserved")
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.srt")
	if err == nil {
		t.Error("ParseFile should fail for non-existent file")
	}
}

func TestRemoveHearingImpairedTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "remove brackets",
			input:    "Hello [Music] world",
			expected: "Hello world",
		},
		{
			name:     "remove parentheses",
			input:    "Hello (sighs) world",
			expected: "Hello world",
		},
		{
			name:     "remove music symbols",
			input:    "♪ La la la ♪",
			expected: "La la la",
		},
		{
			name:     "remove speaker labels",
			input:    "JOHN: Hello there",
			expected: "Hello there",
		},
		{
			name:     "remove speaker with dash",
			input:    "- NARRATOR: Once upon a time",
			expected: "Once upon a time",
		},
		{
			name:     "complex speaker label",
			input:    "Dr. Smith: How are you?",
			expected: "How are you?",
		},
		{
			name:     "no HI tags",
			input:    "Normal subtitle text",
			expected: "Normal subtitle text",
		},
		{
			name:     "multiple HI patterns",
			input:    "[Music] JOHN: Hello (laughs) ♪",
			expected: "Hello",
		},
		{
			name:     "remove music note alt",
			input:    "♫ Song lyrics ♫",
			expected: "Song lyrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveHearingImpairedTags(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBatchLines(t *testing.T) {
	lines := []SubtitleLine{
		{Index: 0, Text: "Line 1"},
		{Index: 1, Text: "Line 2"},
		{Index: 2, Text: "Line 3"},
		{Index: 3, Text: "Line 4"},
		{Index: 4, Text: "Line 5"},
	}

	// Batch size of 2
	batches := BatchLines(lines, 2)
	if len(batches) != 3 {
		t.Errorf("expected 3 batches, got %d", len(batches))
	}

	if len(batches[0]) != 2 {
		t.Errorf("expected first batch size 2, got %d", len(batches[0]))
	}

	if len(batches[1]) != 2 {
		t.Errorf("expected second batch size 2, got %d", len(batches[1]))
	}

	if len(batches[2]) != 1 {
		t.Errorf("expected third batch size 1, got %d", len(batches[2]))
	}
}

func TestBatchLinesEmpty(t *testing.T) {
	lines := []SubtitleLine{}
	batches := BatchLines(lines, 10)

	if len(batches) != 0 {
		t.Errorf("expected 0 batches for empty input, got %d", len(batches))
	}
}

func TestBatchLinesSingleBatch(t *testing.T) {
	lines := []SubtitleLine{
		{Index: 0, Text: "Line 1"},
		{Index: 1, Text: "Line 2"},
	}

	batches := BatchLines(lines, 10)
	if len(batches) != 1 {
		t.Errorf("expected 1 batch, got %d", len(batches))
	}

	if len(batches[0]) != 2 {
		t.Errorf("expected batch size 2, got %d", len(batches[0]))
	}
}

func TestReassembleSRT(t *testing.T) {
	lines := []SubtitleLine{
		{Index: 1, StartTime: "00:00:01,000", EndTime: "00:00:04,000", Text: "Hello"},
		{Index: 2, StartTime: "00:00:05,000", EndTime: "00:00:08,000", Text: "World"},
	}

	result := ReassembleSRT(lines)

	if !strings.Contains(result, "1\n") {
		t.Error("should contain index 1")
	}

	if !strings.Contains(result, "00:00:01,000 --> 00:00:04,000") {
		t.Error("should contain timing")
	}

	if !strings.Contains(result, "Hello") {
		t.Error("should contain text")
	}
}

func TestReassembleASS(t *testing.T) {
	header := `[Script Info]
Title: Test

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`

	lines := []SubtitleLine{
		{
			Layer:     0,
			StartTime: "0:00:01.00",
			EndTime:   "0:00:04.00",
			Style:     "Default",
			MarginL:   0,
			MarginR:   0,
			MarginV:   0,
			Effect:    "",
			Text:      "Hello, world!",
		},
	}

	result := ReassembleASS(header, lines)

	if !strings.Contains(result, "Dialogue:") {
		t.Error("should contain Dialogue line")
	}

	if !strings.Contains(result, "Hello, world!") {
		t.Error("should contain text")
	}

	if !strings.Contains(result, "0:00:01.00") {
		t.Error("should contain start time")
	}
}

func TestSubtitleLineStruct(t *testing.T) {
	line := SubtitleLine{
		Index:      1,
		StartTime:  "00:00:01,000",
		EndTime:    "00:00:05,000",
		Text:       "Test text",
		Style:      "Default",
		OriginalID: 1,
		Layer:      0,
		MarginL:    10,
		MarginR:    10,
		MarginV:    20,
		Effect:     "",
		RawEvent:   "Dialogue: 0,...",
	}

	if line.Index != 1 {
		t.Errorf("unexpected Index: %d", line.Index)
	}

	if line.Text != "Test text" {
		t.Errorf("unexpected Text: %q", line.Text)
	}

	if line.Style != "Default" {
		t.Errorf("unexpected Style: %q", line.Style)
	}
}

func TestSubtitleFileStruct(t *testing.T) {
	sf := SubtitleFile{
		Format:       "srt",
		Header:       "",
		Lines:        []SubtitleLine{{Text: "Test"}},
		LineCount:    1,
		EventsHeader: "",
	}

	if sf.Format != "srt" {
		t.Errorf("unexpected Format: %q", sf.Format)
	}

	if sf.LineCount != 1 {
		t.Errorf("unexpected LineCount: %d", sf.LineCount)
	}
}

func TestParseSRTNoTrailingNewline(t *testing.T) {
	// Test SRT file without trailing newline
	srtContent := `1
00:00:01,000 --> 00:00:04,000
Hello, world!`

	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")
	if err := os.WriteFile(srtPath, []byte(srtContent), 0644); err != nil {
		t.Fatal(err)
	}

	sf, err := ParseFile(srtPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if sf.LineCount != 1 {
		t.Errorf("expected 1 line, got %d", sf.LineCount)
	}
}

func TestParseSRTWithDotTimeSeparator(t *testing.T) {
	// Some SRT files use . instead of , for milliseconds
	srtContent := `1
00:00:01.000 --> 00:00:04.000
Hello, world!

`

	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")
	if err := os.WriteFile(srtPath, []byte(srtContent), 0644); err != nil {
		t.Fatal(err)
	}

	sf, err := ParseFile(srtPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if sf.LineCount != 1 {
		t.Errorf("expected 1 line, got %d", sf.LineCount)
	}
}

func TestReassembleSRTRoundTrip(t *testing.T) {
	// Create original lines
	original := []SubtitleLine{
		{StartTime: "00:00:01,000", EndTime: "00:00:04,000", Text: "Line one"},
		{StartTime: "00:00:05,000", EndTime: "00:00:08,000", Text: "Line two"},
	}

	// Reassemble to SRT
	srtContent := ReassembleSRT(original)

	// Write to temp file
	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "roundtrip.srt")
	if err := os.WriteFile(srtPath, []byte(srtContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Parse back
	parsed, err := ParseFile(srtPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(parsed.Lines) != len(original) {
		t.Fatalf("expected %d lines, got %d", len(original), len(parsed.Lines))
	}

	for i, line := range parsed.Lines {
		if line.Text != original[i].Text {
			t.Errorf("line %d: expected text %q, got %q", i, original[i].Text, line.Text)
		}
	}
}
