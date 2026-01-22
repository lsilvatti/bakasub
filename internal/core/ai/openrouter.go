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
		client:      &http.Client{Timeout: 120 * time.Second},
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

// ValidateKey checks if the API key is valid
func (o *OpenRouterAdapter) ValidateKey(ctx context.Context) bool {
	// Simple validation: try to list models
	models, err := o.ListModels(ctx)
	return err == nil && len(models) > 0
}

// ListModels returns available models from OpenRouter
func (o *OpenRouterAdapter) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse models response
	var modelsResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	models := make([]string, len(modelsResp.Data))
	for i, m := range modelsResp.Data {
		models[i] = m.ID
	}

	return models, nil
}
