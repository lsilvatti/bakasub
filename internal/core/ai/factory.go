package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/lsilvatti/bakasub/internal/config"
)

// ProviderFactory creates LLM provider instances based on configuration
type ProviderFactory struct {
	config *config.Config
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(cfg *config.Config) *ProviderFactory {
	return &ProviderFactory{
		config: cfg,
	}
}

// CreateProvider creates a provider instance based on the current configuration
func (f *ProviderFactory) CreateProvider(ctx context.Context) (LLMProvider, error) {
	if f.config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	// Normalize provider name
	providerName := strings.ToLower(strings.TrimSpace(f.config.AIProvider))

	// Get temperature (default to 0.3 if not set)
	temperature := f.config.Temperature
	if temperature == 0 {
		temperature = 0.3
	}

	// Get model
	model := f.config.Model
	if model == "" {
		return nil, fmt.Errorf("model not configured")
	}

	switch providerName {
	case "openrouter":
		if f.config.APIKey == "" {
			return nil, fmt.Errorf("API key not configured for OpenRouter")
		}
		return NewOpenRouterAdapter(f.config.APIKey, model, temperature), nil

	case "gemini", "google", "google-gemini":
		if f.config.APIKey == "" {
			return nil, fmt.Errorf("API key not configured for Gemini")
		}
		adapter, err := NewGeminiAdapter(ctx, f.config.APIKey, model, temperature)
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini adapter: %w", err)
		}
		return adapter, nil

	case "local", "ollama", "lmstudio":
		if f.config.LocalEndpoint == "" {
			return nil, fmt.Errorf("local endpoint not configured")
		}
		return NewLocalLLMAdapter(f.config.LocalEndpoint, model, temperature), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: openrouter, gemini, local)", providerName)
	}
}

// GetProviderInfo returns metadata about the currently configured provider
func (f *ProviderFactory) GetProviderInfo() (*ProviderInfo, error) {
	if f.config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	providerName := strings.ToLower(strings.TrimSpace(f.config.AIProvider))

	switch providerName {
	case "openrouter":
		return &ProviderInfo{
			Name:        "OpenRouter",
			Type:        "cloud",
			RequiresKey: true,
			Endpoint:    "https://openrouter.ai/api/v1",
		}, nil

	case "gemini", "google", "google-gemini":
		return &ProviderInfo{
			Name:        "Google Gemini",
			Type:        "cloud",
			RequiresKey: true,
			Endpoint:    "https://generativelanguage.googleapis.com",
		}, nil

	case "local", "ollama", "lmstudio":
		endpoint := f.config.LocalEndpoint
		if endpoint == "" {
			endpoint = "http://localhost:11434" // Default Ollama endpoint
		}
		return &ProviderInfo{
			Name:        "Local LLM",
			Type:        "local",
			RequiresKey: false,
			Endpoint:    endpoint,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// ValidateConfiguration checks if the current configuration is valid
func (f *ProviderFactory) ValidateConfiguration(ctx context.Context) error {
	// Check if provider is set
	if f.config.AIProvider == "" {
		return fmt.Errorf("AI provider not configured")
	}

	// Create provider to test
	provider, err := f.CreateProvider(ctx)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	// For Gemini, close the client after validation
	if gemini, ok := provider.(*GeminiAdapter); ok {
		defer gemini.Close()
	}

	// Validate the provider
	if !provider.ValidateKey(ctx) {
		return fmt.Errorf("provider validation failed (check API key/endpoint)")
	}

	return nil
}

// ListAvailableProviders returns a list of all supported provider names
func ListAvailableProviders() []string {
	return []string{
		"openrouter",
		"gemini",
		"local",
	}
}
