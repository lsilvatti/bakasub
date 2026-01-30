package layout

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCalculateHalf(t *testing.T) {
	tests := []struct {
		width    int
		padding  int
		expected int
	}{
		{100, 2, 48},
		{80, 4, 36},
		{60, 0, 30},
		{0, 2, -2},
	}

	for _, tt := range tests {
		result := CalculateHalf(tt.width, tt.padding)
		if result != tt.expected {
			t.Errorf("CalculateHalf(%d, %d) = %d, want %d", tt.width, tt.padding, result, tt.expected)
		}
	}
}

func TestCalculateThird(t *testing.T) {
	tests := []struct {
		width    int
		padding  int
		expected int
	}{
		{120, 2, 38},
		{90, 3, 27},
		{60, 0, 20},
	}

	for _, tt := range tests {
		result := CalculateThird(tt.width, tt.padding)
		if result != tt.expected {
			t.Errorf("CalculateThird(%d, %d) = %d, want %d", tt.width, tt.padding, result, tt.expected)
		}
	}
}

func TestCalculateTwoThirds(t *testing.T) {
	tests := []struct {
		width    int
		padding  int
		expected int
	}{
		{120, 2, 78},
		{90, 3, 57},
		{60, 0, 40},
	}

	for _, tt := range tests {
		result := CalculateTwoThirds(tt.width, tt.padding)
		if result != tt.expected {
			t.Errorf("CalculateTwoThirds(%d, %d) = %d, want %d", tt.width, tt.padding, result, tt.expected)
		}
	}
}

func TestCalculateQuarter(t *testing.T) {
	tests := []struct {
		width    int
		padding  int
		expected int
	}{
		{100, 2, 23},
		{80, 4, 16},
		{40, 0, 10},
	}

	for _, tt := range tests {
		result := CalculateQuarter(tt.width, tt.padding)
		if result != tt.expected {
			t.Errorf("CalculateQuarter(%d, %d) = %d, want %d", tt.width, tt.padding, result, tt.expected)
		}
	}
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"Hello", 10, "Hello"},
		{"Hello World", 5, "He..."},
		{"Hello World", 11, "Hello World"},
		{"Hello World", 8, "Hello..."},
		{"Hi", 2, "Hi"},
		{"Hi", 1, "H"},
		{"", 5, ""},
		{"Test", 0, ""},
		{"Test", -1, ""},
		{"ABC", 3, "ABC"},
	}

	for _, tt := range tests {
		result := TruncateToWidth(tt.input, tt.width)
		if result != tt.expected {
			t.Errorf("TruncateToWidth(%q, %d) = %q, want %q", tt.input, tt.width, result, tt.expected)
		}
	}
}

func TestTruncateLines(t *testing.T) {
	tests := []struct {
		input    string
		maxLines int
		expected string
	}{
		{"line1\nline2\nline3", 2, "line1\nline2"},
		{"line1\nline2", 5, "line1\nline2"},
		{"single line", 1, "single line"},
		{"line1\nline2\nline3\nline4", 3, "line1\nline2\nline3"},
		{"", 2, ""},
		{"test", 0, ""},
		{"test", -1, ""},
	}

	for _, tt := range tests {
		result := TruncateLines(tt.input, tt.maxLines)
		if result != tt.expected {
			t.Errorf("TruncateLines(%q, %d) = %q, want %q", tt.input, tt.maxLines, result, tt.expected)
		}
	}
}

func TestIsTooSmall(t *testing.T) {
	tests := []struct {
		width    int
		height   int
		expected bool
	}{
		{80, 24, false},
		{100, 30, false},
		{79, 24, true},
		{80, 23, true},
		{70, 20, true},
		{0, 0, false},
		{0, 24, false},
		{80, 0, false},
		{200, 50, false},
	}

	for _, tt := range tests {
		result := IsTooSmall(tt.width, tt.height)
		if result != tt.expected {
			t.Errorf("IsTooSmall(%d, %d) = %v, want %v", tt.width, tt.height, result, tt.expected)
		}
	}
}

func TestIsWaitingForSize(t *testing.T) {
	tests := []struct {
		width    int
		height   int
		expected bool
	}{
		{0, 0, true},
		{0, 24, true},
		{80, 0, true},
		{80, 24, false},
		{100, 50, false},
	}

	for _, tt := range tests {
		result := IsWaitingForSize(tt.width, tt.height)
		if result != tt.expected {
			t.Errorf("IsWaitingForSize(%d, %d) = %v, want %v", tt.width, tt.height, result, tt.expected)
		}
	}
}

func TestMinDimensions(t *testing.T) {
	if MinWidth != 80 {
		t.Errorf("MinWidth = %d, want 80", MinWidth)
	}

	if MinHeight != 24 {
		t.Errorf("MinHeight = %d, want 24", MinHeight)
	}
}

func TestPlaceInPanel(t *testing.T) {
	content := "Hello"
	result := PlaceInPanel(content, 20, 5, lipgloss.Center, lipgloss.Center)

	if result == "" {
		t.Error("PlaceInPanel should return non-empty string")
	}

	if !strings.Contains(result, "Hello") {
		t.Error("PlaceInPanel result should contain the content")
	}
}

func TestRenderTooSmallWarning(t *testing.T) {
	result := RenderTooSmallWarning(70, 20)

	if result == "" {
		t.Error("RenderTooSmallWarning should return non-empty string")
	}

	if !strings.Contains(result, "70") || !strings.Contains(result, "20") {
		t.Error("warning should contain current dimensions")
	}
}
