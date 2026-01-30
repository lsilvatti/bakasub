package ai

import (
	"context"
	"errors"
	"testing"
)

// TestProviderErrorError tests ProviderError.Error() method
func TestProviderErrorError(t *testing.T) {
	err := &ProviderError{
		Provider: "openrouter",
		Code:     "rate_limit",
		Message:  "Too many requests",
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() should not return empty string")
	}

	// Should contain provider, code, and message
	if !containsStr(errStr, "openrouter") {
		t.Errorf("Error() should contain provider: %q", errStr)
	}

	if !containsStr(errStr, "rate_limit") {
		t.Errorf("Error() should contain code: %q", errStr)
	}

	if !containsStr(errStr, "Too many requests") {
		t.Errorf("Error() should contain message: %q", errStr)
	}
}

// TestIsRateLimitError tests IsRateLimitError function
func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "rate_limit error",
			err: &ProviderError{
				Code: "rate_limit",
			},
			want: true,
		},
		{
			name: "other error",
			err: &ProviderError{
				Code: "invalid_key",
			},
			want: false,
		},
		{
			name: "generic error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRateLimitError(tt.err)
			if got != tt.want {
				t.Errorf("IsRateLimitError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsAuthError tests IsAuthError function
func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "invalid_key error",
			err: &ProviderError{
				Code: "invalid_key",
			},
			want: true,
		},
		{
			name: "unauthorized error",
			err: &ProviderError{
				Code: "unauthorized",
			},
			want: true,
		},
		{
			name: "other error",
			err: &ProviderError{
				Code: "rate_limit",
			},
			want: false,
		},
		{
			name: "generic error",
			err:  errors.New("generic error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAuthError(tt.err)
			if got != tt.want {
				t.Errorf("IsAuthError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLLMProviderInterface verifies the interface signature
func TestLLMProviderInterface(t *testing.T) {
	// Create a mock that implements the interface
	var _ LLMProvider = &mockProvider{}
}

// mockProvider is a mock implementation of LLMProvider for testing
type mockProvider struct {
	sendBatchFunc   func(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error)
	validateKeyFunc func(ctx context.Context) bool
	listModelsFunc  func(ctx context.Context) ([]string, error)
}

func (m *mockProvider) SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
	if m.sendBatchFunc != nil {
		return m.sendBatchFunc(ctx, payload, systemPrompt)
	}
	return payload, nil
}

func (m *mockProvider) ValidateKey(ctx context.Context) bool {
	if m.validateKeyFunc != nil {
		return m.validateKeyFunc(ctx)
	}
	return true
}

func (m *mockProvider) ListModels(ctx context.Context) ([]string, error) {
	if m.listModelsFunc != nil {
		return m.listModelsFunc(ctx)
	}
	return []string{"test-model"}, nil
}

// TestMockProviderSendBatch tests the mock provider
func TestMockProviderSendBatch(t *testing.T) {
	mock := &mockProvider{
		sendBatchFunc: func(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
			result := make([]Line, len(payload))
			for i, l := range payload {
				result[i] = Line{ID: l.ID, Text: "Translated: " + l.Text}
			}
			return result, nil
		},
	}

	ctx := context.Background()
	payload := []Line{{ID: 0, Text: "Hello"}}
	result, err := mock.SendBatch(ctx, payload, "test prompt")

	if err != nil {
		t.Fatalf("SendBatch failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("len(result) = %d, want 1", len(result))
	}

	if result[0].Text != "Translated: Hello" {
		t.Errorf("result[0].Text = %q, want Translated: Hello", result[0].Text)
	}
}

// TestMockProviderValidateKey tests mock ValidateKey
func TestMockProviderValidateKey(t *testing.T) {
	mock := &mockProvider{
		validateKeyFunc: func(ctx context.Context) bool {
			return true
		},
	}

	ctx := context.Background()
	valid := mock.ValidateKey(ctx)

	if !valid {
		t.Error("ValidateKey should return true")
	}
}

// TestMockProviderListModels tests mock ListModels
func TestMockProviderListModels(t *testing.T) {
	mock := &mockProvider{
		listModelsFunc: func(ctx context.Context) ([]string, error) {
			return []string{"model1", "model2"}, nil
		},
	}

	ctx := context.Background()
	models, err := mock.ListModels(ctx)

	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("len(models) = %d, want 2", len(models))
	}
}

// helper function
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
