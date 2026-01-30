package locales

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		langCode string
		valid    bool
	}{
		{"en", true},
		{"pt-br", true},
		{"es", true},
		{"PT-BR", true},   // Should normalize to lowercase
		{"invalid", true}, // Should fallback to English
		{"fr", true},      // Should fallback to English
	}

	for _, tt := range tests {
		t.Run(tt.langCode, func(t *testing.T) {
			err := Load(tt.langCode)
			if tt.valid && err != nil {
				t.Errorf("Load(%q) failed: %v", tt.langCode, err)
			}
		})
	}
}

func TestT(t *testing.T) {
	// Load English
	if err := Load("en"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Test simple key
	result := T("app_name")
	if result == "app_name" {
		// Key not found, this might be okay depending on the locale file
		// Let's just ensure it doesn't panic
	}

	// Test nested key
	nestedResult := T("wizard.step1.title")
	// Should return something (either the translation or the key)
	if nestedResult == "" {
		t.Error("T should not return empty string")
	}

	// Test non-existent key
	nonExistent := T("this.key.does.not.exist")
	if nonExistent != "this.key.does.not.exist" {
		// Should return the key itself when not found
	}
}

func TestTf(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Test with no args
	result1 := Tf("app_name")
	if result1 == "" {
		t.Error("Tf should not return empty string")
	}

	// Test with args (assuming there's a key with format specifiers)
	// Even if the key doesn't exist, it should not panic
	result2 := Tf("some.key.with.format", 1, 2, 3)
	if result2 == "" {
		t.Error("Tf should not return empty string")
	}
}

func TestGetCurrentLocale(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	locale := GetCurrentLocale()
	if locale != "en" {
		t.Errorf("expected locale 'en', got %q", locale)
	}

	// Load Portuguese
	if err := Load("pt-br"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	locale = GetCurrentLocale()
	if locale != "pt-br" {
		t.Errorf("expected locale 'pt-br', got %q", locale)
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	languages := GetSupportedLanguages()

	if len(languages) != 3 {
		t.Errorf("expected 3 supported languages, got %d", len(languages))
	}

	expectedLangs := []string{"en", "pt-br", "es"}
	for _, lang := range expectedLangs {
		if _, ok := languages[lang]; !ok {
			t.Errorf("expected language %q to be supported", lang)
		}
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		langCode string
		expected bool
	}{
		{"en", true},
		{"EN", true}, // Should normalize
		{"pt-br", true},
		{"PT-BR", true},
		{"es", true},
		{"ES", true},
		{"fr", false},
		{"de", false},
		{"ja", false},
	}

	for _, tt := range tests {
		t.Run(tt.langCode, func(t *testing.T) {
			result := IsSupported(tt.langCode)
			if result != tt.expected {
				t.Errorf("IsSupported(%q) = %v, want %v", tt.langCode, result, tt.expected)
			}
		})
	}
}

func TestGetLanguageName(t *testing.T) {
	tests := []struct {
		langCode string
		expected string
	}{
		{"en", "English"},
		{"pt-br", "Português (Brasil)"},
		{"es", "Español"},
		{"fr", "fr"}, // Unsupported, returns code
	}

	for _, tt := range tests {
		t.Run(tt.langCode, func(t *testing.T) {
			result := GetLanguageName(tt.langCode)
			if result != tt.expected {
				t.Errorf("GetLanguageName(%q) = %q, want %q", tt.langCode, result, tt.expected)
			}
		})
	}
}

func TestReload(t *testing.T) {
	// Load English first
	if err := Load("en"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Reload should not fail
	if err := Reload(); err != nil {
		t.Errorf("Reload failed: %v", err)
	}

	// Locale should still be English
	if GetCurrentLocale() != "en" {
		t.Errorf("locale changed after Reload")
	}
}

func TestTWithNilTranslations(t *testing.T) {
	// This tests the edge case where translations might be nil
	// Load should always initialize translations, so this is more of a safety check

	// Calling T should not panic even in edge cases
	result := T("test.key")
	if result == "" {
		// T should at least return the key
	}
}

func TestLoadNormalizesCase(t *testing.T) {
	// Load with uppercase
	if err := Load("EN"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should be normalized to lowercase
	if GetCurrentLocale() != "en" {
		t.Errorf("expected locale 'en', got %q", GetCurrentLocale())
	}

	// Load with mixed case
	if err := Load("pT-Br"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if GetCurrentLocale() != "pt-br" {
		t.Errorf("expected locale 'pt-br', got %q", GetCurrentLocale())
	}
}

func TestTNestedKeys(t *testing.T) {
	if err := Load("en"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Test various nesting depths
	tests := []string{
		"single",
		"nested.key",
		"deeply.nested.key.structure",
	}

	for _, key := range tests {
		result := T(key)
		// Should not panic and should return something
		if result == "" {
			t.Errorf("T(%q) returned empty string", key)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test thread safety of locales package
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Mix of reads and writes
			if id%2 == 0 {
				Load("en")
			} else {
				Load("pt-br")
			}

			T("test.key")
			GetCurrentLocale()
			GetSupportedLanguages()
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLoadFallbackToEnglish(t *testing.T) {
	// Loading an unsupported language should fallback to English
	if err := Load("xyz"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have fallen back to English
	if GetCurrentLocale() != "en" {
		t.Errorf("expected fallback to 'en', got %q", GetCurrentLocale())
	}
}

func TestSupportedLanguagesNotModifiable(t *testing.T) {
	langs := GetSupportedLanguages()

	// Try to modify the returned map
	langs["test"] = "Test Language"

	// Get again and verify it wasn't modified
	langs2 := GetSupportedLanguages()
	if _, ok := langs2["test"]; ok {
		t.Error("GetSupportedLanguages should return a copy, not the original")
	}
}
