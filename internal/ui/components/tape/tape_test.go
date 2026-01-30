package tape

import (
	"testing"
)

func TestTranslationPairStruct(t *testing.T) {
	pair := TranslationPair{
		ID:           1,
		OriginalText: "Hello",
		Translated:   "Olá",
	}

	if pair.ID != 1 {
		t.Errorf("ID = %d, want 1", pair.ID)
	}

	if pair.OriginalText != "Hello" {
		t.Errorf("OriginalText = %q, want %q", pair.OriginalText, "Hello")
	}

	if pair.Translated != "Olá" {
		t.Errorf("Translated = %q, want %q", pair.Translated, "Olá")
	}
}

func TestNewModel(t *testing.T) {
	m := NewModel(80, 24)

	if m.width != 80 {
		t.Errorf("width = %d, want 80", m.width)
	}

	if m.height != 24 {
		t.Errorf("height = %d, want 24", m.height)
	}

	if m.progress != 0.0 {
		t.Errorf("progress = %f, want 0.0", m.progress)
	}

	if m.spoolFrame != 0 {
		t.Errorf("spoolFrame = %d, want 0", m.spoolFrame)
	}

	if !m.autoScroll {
		t.Error("autoScroll should be true by default")
	}

	if m.maxPairs != 500 {
		t.Errorf("maxPairs = %d, want 500", m.maxPairs)
	}

	if len(m.pairs) != 0 {
		t.Errorf("pairs should be empty, got %d", len(m.pairs))
	}
}

func TestModelColors(t *testing.T) {
	m := NewModel(80, 24)

	if string(m.neonPink) != "#F700FF" {
		t.Errorf("neonPink = %q, want %q", string(m.neonPink), "#F700FF")
	}

	if string(m.cyan) != "#00FFFF" {
		t.Errorf("cyan = %q, want %q", string(m.cyan), "#00FFFF")
	}

	if string(m.yellow) != "#FFFF00" {
		t.Errorf("yellow = %q, want %q", string(m.yellow), "#FFFF00")
	}

	if string(m.dimmed) != "#666666" {
		t.Errorf("dimmed = %q, want %q", string(m.dimmed), "#666666")
	}
}

func TestAddPair(t *testing.T) {
	m := NewModel(80, 24)

	pair := TranslationPair{
		ID:           1,
		OriginalText: "Hello",
		Translated:   "Olá",
	}

	m.AddPair(pair)

	if len(m.pairs) != 1 {
		t.Errorf("pairs length = %d, want 1", len(m.pairs))
	}

	if m.pairs[0].ID != 1 {
		t.Errorf("pairs[0].ID = %d, want 1", m.pairs[0].ID)
	}
}

func TestAddPairAdvancesFrame(t *testing.T) {
	m := NewModel(80, 24)

	initialFrame := m.spoolFrame

	m.AddPair(TranslationPair{ID: 1, OriginalText: "Test", Translated: "Teste"})

	expectedFrame := (initialFrame + 1) % 4
	if m.spoolFrame != expectedFrame {
		t.Errorf("spoolFrame = %d, want %d", m.spoolFrame, expectedFrame)
	}
}

func TestAddPairCircularBuffer(t *testing.T) {
	m := NewModel(80, 24)
	m.maxPairs = 5

	for i := 1; i <= 7; i++ {
		m.AddPair(TranslationPair{
			ID:           i,
			OriginalText: "Original",
			Translated:   "Translated",
		})
	}

	if len(m.pairs) != 5 {
		t.Errorf("pairs length = %d, want 5", len(m.pairs))
	}

	if m.pairs[0].ID != 3 {
		t.Errorf("pairs[0].ID = %d, want 3", m.pairs[0].ID)
	}
}

func TestSetProgress(t *testing.T) {
	m := NewModel(80, 24)

	tests := []struct {
		input    float64
		expected float64
	}{
		{50.0, 50.0},
		{0.0, 0.0},
		{100.0, 100.0},
		{-10.0, 0.0},
		{150.0, 100.0},
	}

	for _, tt := range tests {
		m.SetProgress(tt.input)
		if m.progress != tt.expected {
			t.Errorf("SetProgress(%f) -> progress = %f, want %f", tt.input, m.progress, tt.expected)
		}
	}
}

func TestSetSize(t *testing.T) {
	m := NewModel(80, 24)

	m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}

	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

func TestAutoScrollBehavior(t *testing.T) {
	m := NewModel(80, 24)

	if !m.autoScroll {
		t.Error("autoScroll should be true by default")
	}

	// Manually set autoScroll to false to test behavior
	m.autoScroll = false

	if m.autoScroll {
		t.Error("autoScroll should be false after setting to false")
	}

	m.autoScroll = true

	if !m.autoScroll {
		t.Error("autoScroll should be true after setting to true")
	}
}

func TestView(t *testing.T) {
	m := NewModel(80, 24)

	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}
}

func TestViewWithPairs(t *testing.T) {
	m := NewModel(80, 24)

	m.AddPair(TranslationPair{ID: 1, OriginalText: "Hello", Translated: "Olá"})
	m.AddPair(TranslationPair{ID: 2, OriginalText: "World", Translated: "Mundo"})

	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string with pairs")
	}
}

func TestViewWithProgress(t *testing.T) {
	m := NewModel(80, 24)
	m.SetProgress(75.0)

	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string with progress")
	}
}

func TestClear(t *testing.T) {
	m := NewModel(80, 24)

	m.AddPair(TranslationPair{ID: 1, OriginalText: "Test", Translated: "Teste"})
	m.SetProgress(50.0)

	m.Clear()

	if len(m.pairs) != 0 {
		t.Errorf("pairs should be empty after Clear(), got %d", len(m.pairs))
	}

	if m.progress != 0.0 {
		t.Errorf("progress should be 0 after Clear(), got %f", m.progress)
	}
}

func TestGetPairCount(t *testing.T) {
	m := NewModel(80, 24)

	m.AddPair(TranslationPair{ID: 1, OriginalText: "Hello", Translated: "Olá"})
	m.AddPair(TranslationPair{ID: 2, OriginalText: "World", Translated: "Mundo"})

	count := m.GetPairCount()

	if count != 2 {
		t.Errorf("GetPairCount() = %d, want 2", count)
	}
}

func TestSpoolFrameAnimation(t *testing.T) {
	m := NewModel(80, 24)

	for i := 0; i < 8; i++ {
		expectedFrame := i % 4
		if m.spoolFrame != expectedFrame {
			t.Errorf("iteration %d: spoolFrame = %d, want %d", i, m.spoolFrame, expectedFrame)
		}
		m.AddPair(TranslationPair{ID: i, OriginalText: "Test", Translated: "Teste"})
	}
}
