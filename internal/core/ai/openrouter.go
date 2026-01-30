package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterAdapter implements LLMProvider for OpenRouter API
type OpenRouterAdapter struct {
	apiKey      string
	model       string
	baseURL     string
	client      *http.Client
	temperature float64
}

// NewOpenRouterAdapter creates a new OpenRouter adapter
func NewOpenRouterAdapter(apiKey, model string, temperature float64) *OpenRouterAdapter {
	return &OpenRouterAdapter{
		apiKey:      apiKey,
		model:       model,
		baseURL:     "https://openrouter.ai/api/v1",
		client:      &http.Client{Timeout: 30 * time.Second},
		temperature: temperature,
	}
}

// openRouterRequest represents the API request structure
type openRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []openRouterMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openRouterResponse represents the API response structure
type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// SendBatch sends a batch of lines for translation
func (o *OpenRouterAdapter) SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
	// Convert payload to minified JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build messages
	messages := []openRouterMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: string(payloadJSON)},
	}

	// Build request
	reqBody := openRouterRequest{
		Model:       o.model,
		Messages:    messages,
		Temperature: o.temperature,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/lsilvatti/bakasub")
	req.Header.Set("X-Title", "BakaSub")

	// Send request
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "openrouter",
			Code:     "network_error",
			Message:  err.Error(),
			Retry:    true,
		}
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var apiResp openRouterResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.Error != nil {
		code := apiResp.Error.Code
		if code == "" {
			code = "unknown"
		}

		retry := code == "rate_limit" || code == "timeout" || resp.StatusCode >= 500

		return nil, &ProviderError{
			Provider: "openrouter",
			Code:     code,
			Message:  apiResp.Error.Message,
			Retry:    retry,
		}
	}

	// Check for valid response
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Parse translated lines from response
	content := apiResp.Choices[0].Message.Content

	var translatedLines []Line
	if err := json.Unmarshal([]byte(content), &translatedLines); err != nil {
		return nil, fmt.Errorf("failed to parse translated lines: %w", err)
	}

	return translatedLines, nil
}

// ValidateKey checks if the API key is valid by making a minimal chat request
func (o *OpenRouterAdapter) ValidateKey(ctx context.Context) bool {
	// Make a minimal chat completion request to validate the key
	// OpenRouter's /models endpoint doesn't require auth, so we need to test with a real request
	reqBody := openRouterRequest{
		Model: "meta-llama/llama-3.3-70b-instruct:free", // Use a free model that's reliably available
		Messages: []openRouterMessage{
			{Role: "user", Content: "test"},
		},
		Temperature: 0.0,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(reqJSON))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/lsilvatti/bakasub")
	req.Header.Set("X-Title", "BakaSub")

	resp, err := o.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Read response body for debugging
	body, _ := io.ReadAll(resp.Body)

	// Check for authentication errors
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		// Try to parse error message
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
			} `json:"error"`
		}
		json.Unmarshal(body, &errResp)
		return false
	}

	// Any 2xx status means the key is valid
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// ListModels returns available models from OpenRouter
func (o *OpenRouterAdapter) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/lsilvatti/bakasub")
	req.Header.Set("X-Title", "BakaSub")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "openrouter",
			Code:     "network_error",
			Message:  err.Error(),
			Retry:    true,
		}
	}
	defer resp.Body.Close()

	// Check for authentication/authorization errors
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		// Try to parse error from JSON
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return nil, &ProviderError{
				Provider: "openrouter",
				Code:     "invalid_key",
				Message:  errResp.Error.Message,
				Retry:    false,
			}
		}
		return nil, &ProviderError{
			Provider: "openrouter",
			Code:     "invalid_key",
			Message:  fmt.Sprintf("Authentication failed (HTTP %d): %s", resp.StatusCode, string(body)),
			Retry:    false,
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Provider: "openrouter",
			Code:     "http_error",
			Message:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
			Retry:    resp.StatusCode >= 500,
		}
	}

	// Parse models response with pricing
	var modelsResp struct {
		Data []OpenRouterModel `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	models := make([]string, len(modelsResp.Data))
	for i, m := range modelsResp.Data {
		// Format: id|prompt_price|context_length
		// Example: openai/gpt-4|0.00003|8192
		models[i] = fmt.Sprintf("%s|%s|%d", m.ID, m.Pricing.Prompt, m.ContextLength)
	}

	return models, nil
}

// OpenRouterModel represents a model from the API
type OpenRouterModel struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ContextLength int    `json:"context_length"`
	Pricing       struct {
		Prompt     string `json:"prompt"`
		Completion string `json:"completion"`
	} `json:"pricing"`
}
