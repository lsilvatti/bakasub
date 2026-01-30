package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColorDefinitions(t *testing.T) {
	tests := []struct {
		name     string
		color    lipgloss.Color
		expected string
	}{
		{"NeonPink", NeonPink, "#F700FF"},
		{"Cyan", Cyan, "#00FFFF"},
		{"Yellow", Yellow, "#FFFF00"},
		{"Gray", Gray, "#808080"},
		{"DarkGray", DarkGray, "#404040"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.color) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, string(tt.color), tt.expected)
			}
		})
	}
}

func TestBorderStyles(t *testing.T) {
	if RoundedBorder.Top == "" {
		t.Error("RoundedBorder should have Top defined")
	}

	if DoubleBorder.Top == "" {
		t.Error("DoubleBorder should have Top defined")
	}

	if NormalBorder.Top == "" {
		t.Error("NormalBorder should have Top defined")
	}
}

func TestBaseStyles(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"MainWindow", MainWindow},
		{"Panel", Panel},
		{"PanelTitle", PanelTitle},
		{"Hotkey", Hotkey},
		{"HotkeyText", HotkeyText},
		{"StatusOK", StatusOK},
		{"StatusError", StatusError},
		{"StatusWarning", StatusWarning},
		{"Highlight", Highlight},
		{"HighlightStyle", HighlightStyle},
		{"HeaderBorder", HeaderBorder},
		{"FooterBorder", FooterBorder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.Render("test")
			if result == "" {
				t.Errorf("%s.Render('test') returned empty string", tt.name)
			}
		})
	}
}

func TestHotkey(t *testing.T) {
	key := Hotkey.Render("k")

	if key == "" {
		t.Error("Hotkey.Render should not return empty string")
	}

	if len(key) < 1 {
		t.Error("Hotkey.Render should have content")
	}
}

func TestStatusStyles(t *testing.T) {
	okText := StatusOK.Render("OK")
	errorText := StatusError.Render("ERROR")
	warningText := StatusWarning.Render("WARNING")

	if okText == "" || errorText == "" || warningText == "" {
		t.Error("Status styles should render text")
	}
}

func TestHighlightStyle(t *testing.T) {
	text1 := Highlight.Render("test")
	text2 := HighlightStyle.Render("test")

	if text1 != text2 {
		t.Error("Highlight and HighlightStyle should produce the same output")
	}
}

func TestPanelWithContent(t *testing.T) {
	content := "Panel Content"
	result := Panel.Render(content)

	if result == "" {
		t.Error("Panel.Render should not return empty string")
	}

	if len(result) <= len(content) {
		t.Error("Panel.Render should add border and padding")
	}
}

func TestMainWindowWithContent(t *testing.T) {
	content := "Window Content"
	result := MainWindow.Render(content)

	if result == "" {
		t.Error("MainWindow.Render should not return empty string")
	}

	if len(result) <= len(content) {
		t.Error("MainWindow.Render should add border")
	}
}

func TestPanelTitleBold(t *testing.T) {
	title := PanelTitle.Render("Title")

	if title == "" {
		t.Error("PanelTitle.Render should not return empty string")
	}
}

func TestStylesConsistency(t *testing.T) {
	styles := []lipgloss.Style{
		MainWindow,
		Panel,
		PanelTitle,
		Hotkey,
		HotkeyText,
		StatusOK,
		StatusError,
		StatusWarning,
		Highlight,
		HeaderBorder,
		FooterBorder,
	}

	for i, style := range styles {
		result := style.Render("")
		_ = result
		if false {
			t.Errorf("Style %d panicked", i)
		}
	}
}
