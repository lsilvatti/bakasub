package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/spf13/viper"
)

func main() {
	// Initialize viper config for testing
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.bakasub")

	// Set default language
	viper.SetDefault("ui_language", "en")

	// Try to read config (will fail if not exists, but that's OK)
	_ = viper.ReadInConfig()

	fmt.Println("╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          BAKASUB i18n ENGINE DEMONSTRATION                        ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// List supported languages
	fmt.Println("📋 SUPPORTED LANGUAGES:")
	for code, name := range locales.GetSupportedLanguages() {
		fmt.Printf("   • %s: %s\n", code, name)
	}
	fmt.Println()

	// Test each language
	languages := []string{"en", "pt-br", "es"}

	for _, lang := range languages {
		fmt.Println(strings.Repeat("─", 70))
		fmt.Printf("🌐 TESTING LANGUAGE: %s (%s)\n", lang, locales.GetLanguageName(lang))
		fmt.Println(strings.Repeat("─", 70))

		// Load language
		if err := locales.Load(lang); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading language %s: %v\n", lang, err)
			continue
		}

		fmt.Printf("✓ Loaded: %s\n\n", locales.GetCurrentLocale())

		// Test basic translations
		fmt.Println("🔤 BASIC TRANSLATIONS:")
		fmt.Printf("   App Name:        %s\n", locales.T("app.name"))
		fmt.Printf("   App Tagline:     %s\n", locales.T("app.tagline"))
		fmt.Printf("   Quit:            %s\n", locales.T("common.quit"))
		fmt.Printf("   Loading:         %s\n", locales.T("common.loading"))
		fmt.Println()

		// Test wizard translations
		fmt.Println("🧙 WIZARD TRANSLATIONS:")
		fmt.Printf("   Title:           %s\n", locales.T("wizard.title"))
		fmt.Printf("   Step Indicator:  %s\n", locales.Tf("wizard.step_indicator", 1, 3))
		fmt.Printf("   Provider:        %s\n", locales.T("wizard.step1.provider_openrouter"))
		fmt.Printf("   API Key Label:   %s\n", locales.T("wizard.step1.api_key_label"))
		fmt.Println()

		// Test dashboard translations
		fmt.Println("📊 DASHBOARD TRANSLATIONS:")
		fmt.Printf("   Status:          %s\n", locales.T("dashboard.status.app_running"))
		fmt.Printf("   Extract Module:  %s\n", locales.T("dashboard.modules.extract"))
		fmt.Printf("   Glossary Tool:   %s\n", locales.T("dashboard.toolbox.glossary"))
		fmt.Println()

		// Test execution translations
		fmt.Println("⚙️  EXECUTION TRANSLATIONS:")
		fmt.Printf("   Running Status:  %s\n", locales.T("execution.status.running"))
		fmt.Printf("   Tape Title:      %s\n", locales.T("execution.tape.title"))
		fmt.Printf("   Pair Count:      %s\n", locales.Tf("execution.tape.pair_count", 42))
		fmt.Printf("   Pause Control:   %s\n", locales.T("execution.controls.pause"))
		fmt.Println()

		// Test error translations
		fmt.Println("❌ ERROR TRANSLATIONS:")
		fmt.Printf("   Terminal Small:  %s\n", locales.T("errors.terminal_too_small"))
		fmt.Printf("   Size Message:    %s\n", locales.Tf("errors.terminal_too_small_message", 80, 24))
		fmt.Printf("   BSOD Title:      %s\n", locales.T("errors.bsod_title"))
		fmt.Println()

		// Test missing key fallback
		fmt.Println("🔍 FALLBACK TEST (missing key):")
		fmt.Printf("   Missing Key:     %s\n", locales.T("this.key.does.not.exist"))
		fmt.Printf("   Expected:        this.key.does.not.exist\n")
		fmt.Println()
	}

	fmt.Println(strings.Repeat("═", 70))
	fmt.Println()

	// Test language switching
	fmt.Println("🔄 TESTING LANGUAGE SWITCHING:")
	fmt.Println()

	for i, lang := range []string{"en", "pt-br", "es", "en"} {
		fmt.Printf("%d. Switching to %s...\n", i+1, lang)
		if err := locales.Load(lang); err != nil {
			fmt.Fprintf(os.Stderr, "   ✗ Error: %v\n", err)
			continue
		}
		fmt.Printf("   ✓ Current locale: %s\n", locales.GetCurrentLocale())
		fmt.Printf("   ✓ Quit button: %s\n", locales.T("common.quit"))
		fmt.Println()
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     ALL TESTS COMPLETED                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════╝")
}
