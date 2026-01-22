package ai

import (
	"context"
	"fmt"
)

// Line represents a single subtitle line for translation
type Line struct {
	ID   int    `json:"i"` // Line ID (minified JSON key)
	Text string `json:"t"` // Text content (minified JSON key)
}

// LLMProvider defines the interface for AI translation providers
type LLMProvider interface {
	// SendBatch sends a batch of lines to the AI for translation
	// Returns translated lines in the same order as input
	SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error)

	// ValidateKey checks if the API key/endpoint is valid
	ValidateKey(ctx context.Context) bool

	// ListModels returns available models for this provider
	ListModels(ctx context.Context) ([]string, error)
}

// ProviderInfo contains metadata about a provider
type ProviderInfo struct {
	Name        string // Provider name (openrouter, gemini, openai, local)
	Type        string // cloud or local
	RequiresKey bool   // Whether API key is required
	Endpoint    string // Base URL for API
}

// ProviderError represents an error from a provider
type ProviderError struct {
	Provider string // Provider name
	Code     string // Error code (rate_limit, invalid_key, etc.)
	Message  string // Human-readable message
	Retry    bool   // Whether the request can be retried
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Provider, e.Code, e.Message)
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	if provErr, ok := err.(*ProviderError); ok {
		return provErr.Code == "rate_limit"
	}
	return false
}

// IsAuthError checks if the error is an authentication error
func IsAuthError(err error) bool {
	if provErr, ok := err.(*ProviderError); ok {
		return provErr.Code == "invalid_key" || provErr.Code == "unauthorized"
	}
	return false
}
