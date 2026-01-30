package modelselect

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/ui/focus"
)

func TestModelInfoStruct(t *testing.T) {
	info := ModelInfo{
		ID:          "gpt-4",
		Name:        "GPT-4",
		ContextSize: "128K",
		PricePerM:   "$30.00",
		IsFree:      false,
	}

	if info.ID != "gpt-4" {
		t.Errorf("ID = %q, want %q", info.ID, "gpt-4")
	}

	if info.Name != "GPT-4" {
		t.Errorf("Name = %q, want %q", info.Name, "GPT-4")
	}

	if info.ContextSize != "128K" {
		t.Errorf("ContextSize = %q, want %q", info.ContextSize, "128K")
	}

	if info.PricePerM != "$30.00" {
		t.Errorf("PricePerM = %q, want %q", info.PricePerM, "$30.00")
	}

	if info.IsFree {
		t.Error("IsFree should be false")
	}
}

func TestModelInfoFreeModel(t *testing.T) {
	info := ModelInfo{
		ID:          "llama-free",
		Name:        "Llama Free",
		ContextSize: "8K",
		PricePerM:   "FREE",
		IsFree:      true,
	}

	if !info.IsFree {
		t.Error("IsFree should be true")
	}
}

func TestNew(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	if m.currentTab != 0 {
		t.Errorf("currentTab = %d, want 0 (FREE)", m.currentTab)
	}

	if m.selectedIndex != 0 {
		t.Errorf("selectedIndex = %d, want 0", m.selectedIndex)
	}

	if m.scrollOffset != 0 {
		t.Errorf("scrollOffset = %d, want 0", m.scrollOffset)
	}

	if m.visibleRows != 7 {
		t.Errorf("visibleRows = %d, want 7", m.visibleRows)
	}

	if len(m.availableModels) != 0 {
		t.Errorf("availableModels should be empty, got %d", len(m.availableModels))
	}
}

func TestSetWidth(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	m.SetWidth(100)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}

	// Content width should be width - 6
	expectedContentWidth := 94
	if m.contentWidth != expectedContentWidth {
		t.Errorf("contentWidth = %d, want %d", m.contentWidth, expectedContentWidth)
	}
}

func TestSetWidthMinimum(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	// Set very small width
	m.SetWidth(30)

	// Content width should be at least 50
	if m.contentWidth < 50 {
		t.Errorf("contentWidth = %d, should be at least 50", m.contentWidth)
	}
}

func TestSetActive(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	if m.isActive {
		t.Error("model should not be active initially")
	}

	m.SetActive(true)

	if !m.isActive {
		t.Error("model should be active after SetActive(true)")
	}

	if !m.IsActive() {
		t.Error("IsActive() should return true")
	}

	m.SetActive(false)

	if m.isActive {
		t.Error("model should not be active after SetActive(false)")
	}
}

func TestSetModels(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	models := []ModelInfo{
		{ID: "model-a", Name: "Model A", ContextSize: "8K", PricePerM: "$1.00", IsFree: false},
		{ID: "model-b", Name: "Model B", ContextSize: "32K", PricePerM: "FREE", IsFree: true},
		{ID: "model-c", Name: "Model C", ContextSize: "128K", PricePerM: "$5.00", IsFree: false},
	}

	m.SetModels(models)

	if len(m.availableModels) != 3 {
		t.Errorf("availableModels length = %d, want 3", len(m.availableModels))
	}
}

func TestGetSelectedModel(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	models := []ModelInfo{
		{ID: "model-a", Name: "Model A", ContextSize: "8K", PricePerM: "$1.00", IsFree: false},
		{ID: "model-b", Name: "Model B", ContextSize: "32K", PricePerM: "FREE", IsFree: true},
	}

	m.SetModels(models)

	// Initially should return something (or empty if no filtered models)
	selected := m.GetSelectedModel()
	_ = selected // Just verify no panic
}

func TestView(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)
	m.SetWidth(80)

	models := []ModelInfo{
		{ID: "model-a", Name: "Model A", ContextSize: "8K", PricePerM: "$1.00", IsFree: false},
		{ID: "model-b", Name: "Model B", ContextSize: "32K", PricePerM: "FREE", IsFree: true},
	}

	m.SetModels(models)

	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}
}

func TestViewWithActiveState(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)
	m.SetWidth(80)
	m.SetActive(true)

	models := []ModelInfo{
		{ID: "model-a", Name: "Model A", ContextSize: "8K", PricePerM: "$1.00", IsFree: false},
	}

	m.SetModels(models)

	view := m.View()

	if view == "" {
		t.Error("View() should not return empty string when active")
	}
}

func TestEmptyModels(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)
	m.SetWidth(80)

	// No models set
	view := m.View()

	// Should still render without panic
	if view == "" {
		t.Error("View() should not return empty string even with no models")
	}
}

func TestTabSwitching(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	// Default is FREE tab (0)
	if m.currentTab != 0 {
		t.Errorf("currentTab = %d, want 0", m.currentTab)
	}
}

func TestSearchInput(t *testing.T) {
	fm := focus.NewManager(3)
	m := New(fm)

	// Search input should be initialized
	if m.searchInput.Placeholder == "" {
		t.Error("searchInput should have a placeholder")
	}
}
