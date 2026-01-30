package langselector

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/ui/focus"
)

func TestLanguageStruct(t *testing.T) {
	lang := Language{
		Code:    "en",
		Display: "English",
	}

	if lang.Code != "en" {
		t.Errorf("Code = %q, want %q", lang.Code, "en")
	}

	if lang.Display != "English" {
		t.Errorf("Display = %q, want %q", lang.Display, "English")
	}
}

func TestModeConstants(t *testing.T) {
	if ModeUILanguage != 0 {
		t.Errorf("ModeUILanguage = %d, want 0", ModeUILanguage)
	}

	if ModeTargetLanguage != 1 {
		t.Errorf("ModeTargetLanguage = %d, want 1", ModeTargetLanguage)
	}
}

func TestISOLanguageNames(t *testing.T) {
	expectedMappings := map[string]string{
		"en":    "English",
		"pt-br": "Português (Brasil)",
		"es":    "Español",
		"fr":    "Français",
		"de":    "Deutsch",
		"ja":    "日本語 (Japanese)",
	}

	for code, expectedName := range expectedMappings {
		if name, ok := isoLanguageNames[code]; !ok {
			t.Errorf("ISO code %q should be mapped", code)
		} else if name != expectedName {
			t.Errorf("isoLanguageNames[%q] = %q, want %q", code, name, expectedName)
		}
	}
}

func TestNewUILanguageSelector(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)

	if selector.mode != ModeUILanguage {
		t.Errorf("mode = %d, want %d", selector.mode, ModeUILanguage)
	}

	if selector.selectedIndex != 0 {
		t.Errorf("selectedIndex = %d, want 0", selector.selectedIndex)
	}

	if len(selector.languages) != 3 {
		t.Errorf("expected 3 languages, got %d", len(selector.languages))
	}

	expectedCodes := []string{"en", "pt-br", "es"}
	for i, code := range expectedCodes {
		if selector.languages[i].Code != code {
			t.Errorf("languages[%d].Code = %q, want %q", i, selector.languages[i].Code, code)
		}
	}

	if selector.hasCustomOption {
		t.Error("hasCustomOption should be false for UI language selector")
	}
}

func TestNewTargetLanguageSelector(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewTargetLanguageSelector(fm)

	if selector.mode != ModeTargetLanguage {
		t.Errorf("mode = %d, want %d", selector.mode, ModeTargetLanguage)
	}

	if !selector.hasCustomOption {
		t.Error("hasCustomOption should be true for target language selector")
	}

	if len(selector.languages) < 5 {
		t.Errorf("expected at least 5 languages for target selector, got %d", len(selector.languages))
	}
}

func TestModelSetWidth(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)

	selector.SetWidth(100)

	if selector.width != 100 {
		t.Errorf("width = %d, want 100", selector.width)
	}
}

func TestModelGetSelectedCode(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)

	selected := selector.GetSelectedCode()

	if selected != "en" {
		t.Errorf("GetSelectedCode() = %q, want %q", selected, "en")
	}
}

func TestModelSetSelectedByCode(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)

	selector.SetSelectedByCode("pt-br")

	selected := selector.GetSelectedCode()
	if selected != "pt-br" {
		t.Errorf("after SetSelectedByCode('pt-br'), GetSelectedCode() = %q", selected)
	}
}

func TestModelSetActive(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)

	if selector.isActive {
		t.Error("selector should not be active initially")
	}

	selector.SetActive(true)

	if !selector.isActive {
		t.Error("selector should be active after SetActive(true)")
	}

	selector.SetActive(false)

	if selector.isActive {
		t.Error("selector should not be active after SetActive(false)")
	}
}

func TestModelView(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)
	selector.SetWidth(70)

	view := selector.View()

	if view == "" {
		t.Error("View() should not return empty string")
	}
}

func TestModelViewActive(t *testing.T) {
	fm := focus.NewManager(3)
	selector := NewUILanguageSelector(fm)
	selector.SetWidth(70)
	selector.SetActive(true)

	view := selector.View()

	if view == "" {
		t.Error("View() should not return empty string when active")
	}
}

func TestISOCodeValidation(t *testing.T) {
	commonCodes := []string{
		"en", "en-us", "en-gb",
		"pt", "pt-br",
		"es", "es-la", "es-es",
		"fr", "de", "it", "ja", "ko", "zh",
		"ru", "ar", "hi", "th", "vi",
	}

	for _, code := range commonCodes {
		if _, ok := isoLanguageNames[code]; !ok {
			t.Errorf("ISO code %q should be in isoLanguageNames", code)
		}
	}
}
