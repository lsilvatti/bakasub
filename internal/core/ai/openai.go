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

// OpenAIAdapter implements LLMProvider for OpenAI API
type OpenAIAdapter struct {
	apiKey      string
	model       string
	baseURL     string
	client      *http.Client
	temperature float64
}

// NewOpenAIAdapter creates a new OpenAI adapter
func NewOpenAIAdapter(apiKey, model string, temperature float64) *OpenAIAdapter {
	return &OpenAIAdapter{
		apiKey:      apiKey,
		model:       model,
		baseURL:     "https://api.openai.com/v1",
		client:      &http.Client{Timeout: 120 * time.Second},
		temperature: temperature,
	}
}

// openAIRequest represents the API request structure
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse represents the API response structure
type openAIResponse struct {
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
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// SendBatch sends a batch of lines for translation
func (o *OpenAIAdapter) SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
	// Convert payload to minified JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build messages
	messages := []openAIMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: string(payloadJSON)},
	}

	// Build request
	reqBody := openAIRequest{
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

	// Send request
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "openai",
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
	var apiResp openAIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.Error != nil {
		code := "unknown"
		retry := false
		if apiResp.Error.Type == "insufficient_quota" || apiResp.Error.Code == "rate_limit_exceeded" {
			code = "rate_limit"
			retry = true
		} else if apiResp.Error.Type == "invalid_request_error" && apiResp.Error.Code == "invalid_api_key" {
			code = "invalid_key"
		}

		return nil, &ProviderError{
			Provider: "openai",
			Code:     code,
			Message:  apiResp.Error.Message,
			Retry:    retry,
		}
	}

	// Check for valid response
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse translated lines from response
	content := apiResp.Choices[0].Message.Content

	var translatedLines []Line
	if err := json.Unmarshal([]byte(content), &translatedLines); err != nil {
		return nil, fmt.Errorf("failed to parse translated lines: %w", err)
	}

	return translatedLines, nil
}

// ValidateKey checks if the API key is valid by making a simple API request
func (o *OpenAIAdapter) ValidateKey(ctx context.Context) bool {
	models, err := o.ListModels(ctx)
	return err == nil && len(models) > 0
}

// ListModels returns available models from OpenAI
func (o *OpenAIAdapter) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "openai",
			Code:     "network_error",
			Message:  err.Error(),
			Retry:    true,
		}
	}
	defer resp.Body.Close()

	// Check for authentication errors
	if resp.StatusCode == http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Provider: "openai",
			Code:     "invalid_key",
			Message:  fmt.Sprintf("Invalid API key: %s", string(body)),
			Retry:    false,
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Provider: "openai",
			Code:     "http_error",
			Message:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
			Retry:    resp.StatusCode >= 500,
		}
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

	// Filter for GPT models
	var models []string
	for _, m := range modelsResp.Data {
		// Include gpt models for chat
		if len(m.ID) >= 3 && m.ID[:3] == "gpt" {
			models = append(models, m.ID)
		}
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no compatible GPT models found")
	}

	return models, nil
}

// Close is a no-op for HTTP-based implementation
func (o *OpenAIAdapter) Close() error {
	return nil
}
