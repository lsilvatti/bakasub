package ai

import (
	"context"
	"testing"

	"github.com/lsilvatti/bakasub/internal/config"
)

func TestNewProviderFactory(t *testing.T) {
	cfg := config.Default()
	factory := NewProviderFactory(cfg)

	if factory == nil {
		t.Fatal("NewProviderFactory returned nil")
	}

	if factory.config != cfg {
		t.Error("factory config not set correctly")
	}
}

func TestProviderFactoryNilConfig(t *testing.T) {
	factory := &ProviderFactory{
		config: nil,
	}

	ctx := context.Background()
	_, err := factory.CreateProvider(ctx)

	if err == nil {
		t.Error("Expected error for nil config")
	}
}

func TestCreateProviderOpenRouter(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "openrouter"
	cfg.APIKey = "test-key"
	cfg.Model = "gpt-4o"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestCreateProviderOpenAI(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "openai"
	cfg.APIKey = "test-key"
	cfg.Model = "gpt-4o"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestCreateProviderGemini(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "gemini"
	cfg.APIKey = "test-key"
	cfg.Model = "gemini-pro"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestCreateProviderLocal(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "local"
	cfg.LocalEndpoint = "http://localhost:11434"
	cfg.Model = "llama2"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}

func TestCreateProviderMissingAPIKey(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "openrouter"
	cfg.APIKey = ""
	cfg.Model = "gpt-4o"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	_, err := factory.CreateProvider(ctx)
	if err == nil {
		t.Error("Expected error for missing API key")
	}
}

func TestCreateProviderMissingModel(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "openrouter"
	cfg.APIKey = "test-key"
	cfg.Model = ""

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	_, err := factory.CreateProvider(ctx)
	if err == nil {
		t.Error("Expected error for missing model")
	}
}

func TestCreateProviderMissingLocalEndpoint(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "local"
	cfg.LocalEndpoint = ""
	cfg.Model = "llama2"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	_, err := factory.CreateProvider(ctx)
	if err == nil {
		t.Error("Expected error for missing local endpoint")
	}
}

func TestCreateProviderUnsupported(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "unsupported-provider"
	cfg.Model = "test-model"

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	_, err := factory.CreateProvider(ctx)
	if err == nil {
		t.Error("Expected error for unsupported provider")
	}
}

func TestCreateProviderNameNormalization(t *testing.T) {
	tests := []struct {
		name        string
		providerStr string
		apiKey      string
		endpoint    string
		model       string
		wantErr     bool
	}{
		{"openrouter lowercase", "openrouter", "key", "", "model", false},
		{"openrouter uppercase", "OPENROUTER", "key", "", "model", false},
		{"openrouter mixed", "OpenRouter", "key", "", "model", false},
		{"openrouter trimmed", "  openrouter  ", "key", "", "model", false},
		{"gemini", "gemini", "key", "", "model", false},
		{"google", "google", "key", "", "model", false},
		{"google-gemini", "google-gemini", "key", "", "model", false},
		{"local", "local", "", "http://localhost:11434", "model", false},
		{"ollama", "ollama", "", "http://localhost:11434", "model", false},
		{"lmstudio", "lmstudio", "", "http://localhost:11434", "model", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Default()
			cfg.AIProvider = tt.providerStr
			cfg.APIKey = tt.apiKey
			cfg.LocalEndpoint = tt.endpoint
			cfg.Model = tt.model

			factory := NewProviderFactory(cfg)
			ctx := context.Background()

			_, err := factory.CreateProvider(ctx)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetProviderInfo(t *testing.T) {
	tests := []struct {
		provider string
		wantName string
		wantType string
		wantKey  bool
	}{
		{"openrouter", "OpenRouter", "cloud", true},
		{"openai", "OpenAI", "cloud", true},
		{"gemini", "Google Gemini", "cloud", true},
		{"local", "Local LLM", "local", false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			cfg := config.Default()
			cfg.AIProvider = tt.provider

			factory := NewProviderFactory(cfg)
			info, err := factory.GetProviderInfo()
			if err != nil {
				t.Fatalf("GetProviderInfo failed: %v", err)
			}

			if info.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", info.Name, tt.wantName)
			}

			if info.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", info.Type, tt.wantType)
			}

			if info.RequiresKey != tt.wantKey {
				t.Errorf("RequiresKey = %v, want %v", info.RequiresKey, tt.wantKey)
			}
		})
	}
}

func TestGetProviderInfoNilConfig(t *testing.T) {
	factory := &ProviderFactory{config: nil}

	_, err := factory.GetProviderInfo()
	if err == nil {
		t.Error("Expected error for nil config")
	}
}

func TestGetProviderInfoUnsupported(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "unsupported"

	factory := NewProviderFactory(cfg)
	_, err := factory.GetProviderInfo()
	if err == nil {
		t.Error("Expected error for unsupported provider")
	}
}

func TestProviderInfoStruct(t *testing.T) {
	info := ProviderInfo{
		Name:        "Test Provider",
		Type:        "cloud",
		RequiresKey: true,
		Endpoint:    "https://api.example.com",
	}

	if info.Name != "Test Provider" {
		t.Errorf("Name = %q, want Test Provider", info.Name)
	}

	if info.Type != "cloud" {
		t.Errorf("Type = %q, want cloud", info.Type)
	}

	if !info.RequiresKey {
		t.Error("RequiresKey should be true")
	}

	if info.Endpoint != "https://api.example.com" {
		t.Errorf("Endpoint = %q, want https://api.example.com", info.Endpoint)
	}
}

func TestDefaultTemperature(t *testing.T) {
	cfg := config.Default()
	cfg.AIProvider = "openrouter"
	cfg.APIKey = "test-key"
	cfg.Model = "test-model"
	cfg.Temperature = 0 // Should default to 0.3

	factory := NewProviderFactory(cfg)
	ctx := context.Background()

	provider, err := factory.CreateProvider(ctx)
	if err != nil {
		t.Fatalf("CreateProvider failed: %v", err)
	}

	if provider == nil {
		t.Fatal("Provider is nil")
	}
}
