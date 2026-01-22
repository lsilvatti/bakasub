package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/ai"
)

func main() {
	// Parse command-line flags
	validateOnly := flag.Bool("validate", false, "Only validate provider configuration")
	listModels := flag.Bool("list", false, "List available models")
	testTranslation := flag.Bool("test", false, "Test translation with sample data")
	providerFlag := flag.String("provider", "", "Override provider (openrouter, gemini, local)")
	model := flag.String("model", "", "Override model")
	apiKey := flag.String("key", "", "Override API key")
	endpoint := flag.String("endpoint", "", "Override endpoint (for local LLM)")
	temperature := flag.Float64("temp", 0.3, "Temperature (0.0-1.0)")

	flag.Parse()

	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  BakaSub AI Package Demo                                                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("⚠ Config not found, using defaults or CLI args\n")
		cfg = &config.Config{
			AIProvider:  "openrouter",
			Model:       "google/gemini-flash-1.5",
			Temperature: 0.3,
		}
	}

	// Override with CLI flags if provided
	if *providerFlag != "" {
		cfg.AIProvider = *providerFlag
	}
	if *model != "" {
		cfg.Model = *model
	}
	if *apiKey != "" {
		cfg.APIKey = *apiKey
	}
	if *endpoint != "" {
		cfg.LocalEndpoint = *endpoint
	}
	if *temperature != 0.3 {
		cfg.Temperature = *temperature
	}

	// Display current configuration
	fmt.Println("┌── CONFIGURATION ────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ Provider:    %-60s │\n", cfg.AIProvider)
	fmt.Printf("│ Model:       %-60s │\n", cfg.Model)
	if cfg.AIProvider == "local" || cfg.AIProvider == "ollama" {
		fmt.Printf("│ Endpoint:    %-60s │\n", cfg.LocalEndpoint)
	} else {
		keyDisplay := maskAPIKey(cfg.APIKey)
		fmt.Printf("│ API Key:     %-60s │\n", keyDisplay)
	}
	fmt.Printf("│ Temperature: %.2f%-58s │\n", cfg.Temperature, "")
	fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Create provider factory
	factory := ai.NewProviderFactory(cfg)

	// Get provider info
	info, err := factory.GetProviderInfo()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("┌── PROVIDER INFO ────────────────────────────────────────────────────────────┐")
	fmt.Printf("│ Name:        %-60s │\n", info.Name)
	fmt.Printf("│ Type:        %-60s │\n", info.Type)
	fmt.Printf("│ Endpoint:    %-60s │\n", info.Endpoint)
	fmt.Printf("│ Requires Key: %-59v │\n", info.RequiresKey)
	fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create provider
	fmt.Println("⏳ Creating provider instance...")
	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to create provider: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Provider created successfully")
	fmt.Println()

	// Close Gemini client if applicable
	if gemini, ok := provider.(*ai.GeminiAdapter); ok {
		defer gemini.Close()
	}

	// Validate configuration
	if *validateOnly || *listModels || *testTranslation {
		fmt.Println("⏳ Validating provider configuration...")
		if provider.ValidateKey(ctx) {
			fmt.Println("✓ Provider validation successful")
		} else {
			fmt.Println("❌ Provider validation failed (check API key/endpoint)")
			os.Exit(1)
		}
		fmt.Println()
	}

	// List models
	if *listModels {
		fmt.Println("⏳ Fetching available models...")
		models, err := provider.ListModels(ctx)
		if err != nil {
			fmt.Printf("❌ Failed to list models: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("┌── AVAILABLE MODELS ─────────────────────────────────────────────────────────┐")
		for i, model := range models {
			if i >= 20 {
				fmt.Printf("│ ... and %d more models%-50s │\n", len(models)-20, "")
				break
			}
			fmt.Printf("│ • %-72s │\n", model)
		}
		fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
		fmt.Println()
	}

	// Test translation
	if *testTranslation {
		fmt.Println("⏳ Testing translation with sample data...")

		// Sample subtitle lines (English -> Portuguese)
		sampleLines := []ai.Line{
			{ID: 1, Text: "Hello, how are you?"},
			{ID: 2, Text: "I'm fine, thank you!"},
			{ID: 3, Text: "What a beautiful day!"},
		}

		systemPrompt := `You are a professional translator. Translate the following JSON array from English to Portuguese (Brazil).
Maintain the same structure: [{"i":1, "t":"translated text"}].
Be natural and conversational.`

		fmt.Println()
		fmt.Println("┌── INPUT (SAMPLE DATA) ──────────────────────────────────────────────────────┐")
		for _, line := range sampleLines {
			fmt.Printf("│ [%d] %-69s │\n", line.ID, line.Text)
		}
		fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
		fmt.Println()

		startTime := time.Now()
		translatedLines, err := provider.SendBatch(ctx, sampleLines, systemPrompt)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Translation failed: %v\n", err)

			// Check error type
			if ai.IsRateLimitError(err) {
				fmt.Println("   → Rate limit hit. Try again later.")
			} else if ai.IsAuthError(err) {
				fmt.Println("   → Authentication failed. Check your API key.")
			}

			os.Exit(1)
		}

		fmt.Println("┌── OUTPUT (TRANSLATED) ──────────────────────────────────────────────────────┐")
		for _, line := range translatedLines {
			fmt.Printf("│ [%d] %-69s │\n", line.ID, line.Text)
		}
		fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
		fmt.Println()

		fmt.Printf("✓ Translation completed in %.2fs\n", duration.Seconds())
		fmt.Println()

		// Display raw JSON (for debugging)
		fmt.Println("┌── RAW JSON OUTPUT ──────────────────────────────────────────────────────────┐")
		jsonOutput, _ := json.MarshalIndent(translatedLines, "", "  ")
		fmt.Printf("│ %s\n", jsonOutput)
		fmt.Println("└─────────────────────────────────────────────────────────────────────────────┘")
		fmt.Println()
	}

	if !*validateOnly && !*listModels && !*testTranslation {
		fmt.Println("ℹ No action specified. Use:")
		fmt.Println("  --validate    Validate provider configuration")
		fmt.Println("  --list        List available models")
		fmt.Println("  --test        Test translation with sample data")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run cmd/ai-demo/main.go --validate")
		fmt.Println("  go run cmd/ai-demo/main.go --list --provider gemini")
		fmt.Println("  go run cmd/ai-demo/main.go --test --provider openrouter --model google/gemini-flash-1.5")
		fmt.Println()
	}

	fmt.Println("✓ Demo completed")
}

func maskAPIKey(key string) string {
	if key == "" {
		return "[NOT SET]"
	}
	if len(key) <= 8 {
		return "********"
	}
	return key[:4] + "********************************" + key[len(key)-4:]
}
