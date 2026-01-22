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

// LocalLLMAdapter implements LLMProvider for local LLM servers (Ollama, LMStudio)
type LocalLLMAdapter struct {
	endpoint    string
	model       string
	client      *http.Client
	temperature float64
}

// NewLocalLLMAdapter creates a new local LLM adapter
func NewLocalLLMAdapter(endpoint, model string, temperature float64) *LocalLLMAdapter {
	return &LocalLLMAdapter{
		endpoint:    endpoint,
		model:       model,
		client:      &http.Client{Timeout: 300 * time.Second}, // Longer timeout for local inference
		temperature: temperature,
	}
}

// localLLMRequest represents the Ollama API request structure
type localLLMRequest struct {
	Model       string            `json:"model"`
	Messages    []localLLMMessage `json:"messages"`
	Stream      bool              `json:"stream"`
	Temperature float64           `json:"temperature"`
}

type localLLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// localLLMResponse represents the Ollama API response structure
type localLLMResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

// SendBatch sends a batch of lines for translation
func (l *LocalLLMAdapter) SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
	// Convert payload to minified JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build messages (Ollama format)
	messages := []localLLMMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: string(payloadJSON)},
	}

	// Build request
	reqBody := localLLMRequest{
		Model:       l.model,
		Messages:    messages,
		Stream:      false,
		Temperature: l.temperature,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request (Ollama uses /api/chat endpoint)
	url := l.endpoint + "/api/chat"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "local",
			Code:     "network_error",
			Message:  fmt.Sprintf("failed to connect to %s: %v", l.endpoint, err),
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
	var apiResp localLLMResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.Error != "" {
		return nil, &ProviderError{
			Provider: "local",
			Code:     "inference_error",
			Message:  apiResp.Error,
			Retry:    false,
		}
	}

	// Parse translated lines from response
	content := apiResp.Message.Content

	var translatedLines []Line
	if err := json.Unmarshal([]byte(content), &translatedLines); err != nil {
		return nil, fmt.Errorf("failed to parse translated lines: %w", err)
	}

	return translatedLines, nil
}

// ValidateKey checks if the local server is accessible
func (l *LocalLLMAdapter) ValidateKey(ctx context.Context) bool {
	// For local LLM, just check if the server is reachable
	req, err := http.NewRequestWithContext(ctx, "GET", l.endpoint+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// ListModels returns available models from the local server
func (l *LocalLLMAdapter) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", l.endpoint+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "local",
			Code:     "network_error",
			Message:  err.Error(),
			Retry:    true,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse Ollama tags response
	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	models := make([]string, len(tagsResp.Models))
	for i, m := range tagsResp.Models {
		models[i] = m.Name
	}

	return models, nil
}
