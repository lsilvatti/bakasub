package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOpenRouterAdapterStruct tests the OpenRouterAdapter structure
func TestOpenRouterAdapterStruct(t *testing.T) {
	adapter := NewOpenRouterAdapter("test-key", "gpt-4o", 0.7)
	if adapter == nil {
		t.Fatal("NewOpenRouterAdapter returned nil")
	}
}

// TestOpenRouterAdapterValidation tests ValidateKey with mock server
func TestOpenRouterAdapterValidation(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer valid-key" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Invalid API key",
					"code":    "invalid_key",
				},
			})
			return
		}

		// Return valid response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "ok",
					},
				},
			},
		})
	}))
	defer server.Close()

	// Test with invalid key
	adapter := &OpenRouterAdapter{
		apiKey:      "invalid-key",
		model:       "test-model",
		baseURL:     server.URL,
		client:      &http.Client{},
		temperature: 0.7,
	}

	ctx := context.Background()
	valid := adapter.ValidateKey(ctx)
	if valid {
		t.Error("Expected ValidateKey to return false for invalid key")
	}
}

// TestOpenRouterAdapterSendBatch tests the SendBatch method
func TestOpenRouterAdapterSendBatch(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return translated lines
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `[{"i":0,"t":"Olá mundo"}]`,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := &OpenRouterAdapter{
		apiKey:      "test-key",
		model:       "test-model",
		baseURL:     server.URL,
		client:      &http.Client{},
		temperature: 0.7,
	}

	ctx := context.Background()
	payload := []Line{{ID: 0, Text: "Hello world"}}
	result, err := adapter.SendBatch(ctx, payload, "Translate to Portuguese")

	if err != nil {
		t.Fatalf("SendBatch returned error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

// TestGeminiAdapterStruct tests the GeminiAdapter structure
func TestGeminiAdapterStruct(t *testing.T) {
	ctx := context.Background()
	adapter, err := NewGeminiAdapter(ctx, "test-key", "gemini-pro", 0.7)
	if err != nil {
		t.Fatalf("NewGeminiAdapter returned error: %v", err)
	}
	if adapter == nil {
		t.Fatal("NewGeminiAdapter returned nil")
	}
}

// TestGeminiAdapterSendBatch tests the Gemini SendBatch method
func TestGeminiAdapterSendBatch(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]interface{}{
							{"text": `[{"i":0,"t":"Olá mundo"}]`},
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := &GeminiAdapter{
		apiKey:      "test-key",
		model:       "gemini-pro",
		baseURL:     server.URL,
		client:      &http.Client{},
		temperature: 0.7,
	}

	ctx := context.Background()
	payload := []Line{{ID: 0, Text: "Hello world"}}
	result, err := adapter.SendBatch(ctx, payload, "Translate to Portuguese")

	if err != nil {
		t.Fatalf("SendBatch returned error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

// TestOpenAIAdapterStruct tests the OpenAIAdapter structure
func TestOpenAIAdapterStruct(t *testing.T) {
	adapter := NewOpenAIAdapter("test-key", "gpt-4o", 0.7)
	if adapter == nil {
		t.Fatal("NewOpenAIAdapter returned nil")
	}
}

// TestOpenAIAdapterSendBatch tests the OpenAI SendBatch method
func TestOpenAIAdapterSendBatch(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `[{"i":0,"t":"Olá mundo"}]`,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := &OpenAIAdapter{
		apiKey:      "test-key",
		model:       "gpt-4o",
		baseURL:     server.URL,
		client:      &http.Client{},
		temperature: 0.7,
	}

	ctx := context.Background()
	payload := []Line{{ID: 0, Text: "Hello world"}}
	result, err := adapter.SendBatch(ctx, payload, "Translate to Portuguese")

	if err != nil {
		t.Fatalf("SendBatch returned error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

// TestLocalLLMAdapterStruct tests the LocalLLMAdapter structure
func TestLocalLLMAdapterStruct(t *testing.T) {
	adapter := NewLocalLLMAdapter("http://localhost:11434", "llama2", 0.7)
	if adapter == nil {
		t.Fatal("NewLocalLLMAdapter returned nil")
	}
}

// TestLocalLLMAdapterSendBatch tests the Local LLM SendBatch method
func TestLocalLLMAdapterSendBatch(t *testing.T) {
	// Create mock server that simulates Ollama's /api/chat endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ollama uses /api/chat endpoint and returns response in a different format
		response := map[string]interface{}{
			"message": map[string]interface{}{
				"content": `[{"i":0,"t":"Olá mundo"}]`,
			},
			"done": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := &LocalLLMAdapter{
		endpoint:    server.URL,
		model:       "llama2",
		client:      &http.Client{},
		temperature: 0.7,
	}

	ctx := context.Background()
	payload := []Line{{ID: 0, Text: "Hello world"}}
	result, err := adapter.SendBatch(ctx, payload, "Translate to Portuguese")

	if err != nil {
		t.Fatalf("SendBatch returned error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
}

// TestProviderErrorStruct tests the ProviderError structure
func TestProviderErrorStruct(t *testing.T) {
	err := &ProviderError{
		Provider: "openrouter",
		Code:     "rate_limit",
		Message:  "Too many requests",
		Retry:    true,
	}

	if err.Provider != "openrouter" {
		t.Errorf("Expected Provider 'openrouter', got %q", err.Provider)
	}

	if err.Code != "rate_limit" {
		t.Errorf("Expected Code 'rate_limit', got %q", err.Code)
	}

	if !err.Retry {
		t.Error("Expected Retry to be true")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() returned empty string")
	}
}

// TestLineStruct tests the Line struct
func TestLineStruct(t *testing.T) {
	line := Line{
		ID:   1,
		Text: "Hello world",
	}

	if line.ID != 1 {
		t.Errorf("Expected ID 1, got %d", line.ID)
	}

	if line.Text != "Hello world" {
		t.Errorf("Expected Text 'Hello world', got %q", line.Text)
	}
}

// TestLineJSONSerialization tests Line JSON serialization
func TestLineJSONSerialization(t *testing.T) {
	line := Line{ID: 10, Text: "Test line"}
	data, err := json.Marshal(line)
	if err != nil {
		t.Fatalf("Failed to marshal Line: %v", err)
	}

	var unmarshaled Line
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Line: %v", err)
	}

	if unmarshaled.ID != line.ID || unmarshaled.Text != line.Text {
		t.Errorf("Unmarshaled Line doesn't match original")
	}
}

// TestAPIErrorHandling tests error handling for API errors
func TestAPIErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		errorCode   string
		expectRetry bool
	}{
		{"Rate Limit", 429, "rate_limit", true},
		{"Server Error", 500, "server_error", true},
		{"Invalid Key", 401, "invalid_key", false},
		{"Bad Request", 400, "bad_request", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"message": "Test error",
						"code":    tt.errorCode,
					},
				})
			}))
			defer server.Close()

			adapter := &OpenRouterAdapter{
				apiKey:      "test-key",
				model:       "test-model",
				baseURL:     server.URL,
				client:      &http.Client{},
				temperature: 0.7,
			}

			ctx := context.Background()
			_, err := adapter.SendBatch(ctx, []Line{{ID: 0, Text: "test"}}, "test")

			if err == nil {
				t.Error("Expected error but got nil")
			}
		})
	}
}
