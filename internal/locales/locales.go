package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

//go:embed all:*.json
var localeFiles embed.FS

// currentLocale holds the currently loaded language code
var currentLocale string = "en"

// translations holds the current translation map
var translations map[string]interface{}

// mutex protects concurrent access to translations
var mu sync.RWMutex

// supportedLanguages maps language codes to their display names
var supportedLanguages = map[string]string{
	"en":    "English",
	"pt-br": "Português (Brasil)",
	"es":    "Español",
}

// init loads English by default when the package is imported
func init() {
	// Load English as the default language
	// This ensures translations are always available
	if err := Load("en"); err != nil {
		// This should never happen as English is always embedded
		panic(fmt.Sprintf("Failed to load default English locale: %v", err))
	}
}

// Init initializes the i18n system by loading the configured language
// If no language is configured, it defaults to "en"
func Init() error {
	langCode := viper.GetString("ui_language")
	if langCode == "" {
		langCode = "en"
	}
	err := Load(langCode)
	if err != nil {
		// If loading fails, fallback to English
		return Load("en")
	}
	return nil
}

// Load loads the translation file for the specified language code
// Returns error if the language is not supported or file cannot be parsed
func Load(langCode string) error {
	// Normalize language code to lowercase
	langCode = strings.ToLower(langCode)

	// Check if language is supported
	if _, ok := supportedLanguages[langCode]; !ok {
		// Fallback to English if unsupported
		langCode = "en"
	}

	// Construct filename
	filename := fmt.Sprintf("%s.json", langCode)

	// Read embedded file
	data, err := localeFiles.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read locale file %s: %w", filename, err)
	}

	// Parse JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("failed to parse locale file %s: %w", filename, err)
	}

	// Update global translations
	mu.Lock()
	translations = parsed
	currentLocale = langCode
	mu.Unlock()

	return nil
}

// T retrieves a translation for the given key
// Supports nested keys using dot notation (e.g., "wizard.step1.title")
// Falls back to the key itself if not found
func T(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if translations == nil {
		return key
	}

	// Split key by dots for nested access
	parts := strings.Split(key, ".")

	// Navigate through the nested map
	var current interface{} = translations
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			if next, ok := v[part]; ok {
				current = next
			} else {
				// Key not found, return original key
				return key
			}
		default:
			// Not a map, cannot navigate further
			return key
		}
	}

	// Check if final value is a string
	if str, ok := current.(string); ok {
		return str
	}

	// Final value is not a string, return key
	return key
}

// Tf is a formatted translation helper that uses fmt.Sprintf
// Example: Tf("wizard.step_indicator", 1, 3) -> "STEP 1/3"
func Tf(key string, args ...interface{}) string {
	template := T(key)
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}

// GetCurrentLocale returns the currently loaded language code
func GetCurrentLocale() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLocale
}

// GetSupportedLanguages returns a map of supported language codes and names
func GetSupportedLanguages() map[string]string {
	result := make(map[string]string)
	for k, v := range supportedLanguages {
		result[k] = v
	}
	return result
}

// IsSupported checks if a language code is supported
func IsSupported(langCode string) bool {
	langCode = strings.ToLower(langCode)
	_, ok := supportedLanguages[langCode]
	return ok
}

// GetLanguageName returns the display name for a language code
// Returns the code itself if not found
func GetLanguageName(langCode string) string {
	langCode = strings.ToLower(langCode)
	if name, ok := supportedLanguages[langCode]; ok {
		return name
	}
	return langCode
}

// Reload reloads the current language (useful after config changes)
func Reload() error {
	return Load(currentLocale)
}

// SwitchLanguage switches to a different language and updates the config
func SwitchLanguage(langCode string) error {
	if err := Load(langCode); err != nil {
		return err
	}

	// Update viper config
	viper.Set("ui_language", langCode)
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, try to create it
		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("failed to save language preference: %w", err)
		}
	}

	return nil
}
