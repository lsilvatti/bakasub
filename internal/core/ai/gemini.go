package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiAdapter implements LLMProvider for Google Gemini API using REST API
type GeminiAdapter struct {
	apiKey      string
	model       string
	baseURL     string
	client      *http.Client
	temperature float64
}

// NewGeminiAdapter creates a new Gemini adapter
func NewGeminiAdapter(ctx context.Context, apiKey, model string, temperature float64) (*GeminiAdapter, error) {
	return &GeminiAdapter{
		apiKey:      apiKey,
		model:       model,
		baseURL:     "https://generativelanguage.googleapis.com/v1beta",
		client:      &http.Client{Timeout: 120 * time.Second},
		temperature: temperature,
	}, nil
}

// geminiRequest represents the API request structure
type geminiRequest struct {
	Contents         []geminiContent `json:"contents"`
	GenerationConfig geminiGenConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenConfig struct {
	Temperature float64 `json:"temperature,omitempty"`
}

// geminiResponse represents the API response structure
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

// SendBatch sends a batch of lines for translation
func (g *GeminiAdapter) SendBatch(ctx context.Context, payload []Line, systemPrompt string) ([]Line, error) {
	// Convert payload to minified JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build contents (combining system prompt and user payload)
	fullPrompt := systemPrompt + "\n\n" + string(payloadJSON)
	contents := []geminiContent{
		{
			Role: "user",
			Parts: []geminiPart{
				{Text: fullPrompt},
			},
		},
	}

	// Build request
	reqBody := geminiRequest{
		Contents: contents,
		GenerationConfig: geminiGenConfig{
			Temperature: g.temperature,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", g.baseURL, g.model, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "gemini",
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
	var apiResp geminiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.Error != nil {
		code := "unknown"
		if apiResp.Error.Code == 429 {
			code = "rate_limit"
		} else if apiResp.Error.Code == 401 || apiResp.Error.Code == 403 {
			code = "invalid_key"
		}

		retry := apiResp.Error.Code == 429 || apiResp.Error.Code >= 500

		return nil, &ProviderError{
			Provider: "gemini",
			Code:     code,
			Message:  apiResp.Error.Message,
			Retry:    retry,
		}
	}

	// Check for valid response
	if len(apiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	// Extract text content
	var content string
	for _, part := range apiResp.Candidates[0].Content.Parts {
		content += part.Text
	}

	if content == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	// Parse translated lines from response
	var translatedLines []Line
	if err := json.Unmarshal([]byte(content), &translatedLines); err != nil {
		return nil, fmt.Errorf("failed to parse translated lines: %w", err)
	}

	return translatedLines, nil
}

// ValidateKey checks if the API key is valid
func (g *GeminiAdapter) ValidateKey(ctx context.Context) bool {
	// Simple validation: try to list models
	models, err := g.ListModels(ctx)
	return err == nil && len(models) > 0
}

// ListModels returns available Gemini models
func (g *GeminiAdapter) ListModels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/models?key=%s", g.baseURL, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Provider: "gemini",
			Code:     "network_error",
			Message:  err.Error(),
			Retry:    true,
		}
	}
	defer resp.Body.Close()

	// Check for authentication errors
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Provider: "gemini",
			Code:     "invalid_key",
			Message:  fmt.Sprintf("Invalid API key: %s", string(body)),
			Retry:    false,
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ProviderError{
			Provider: "gemini",
			Code:     "http_error",
			Message:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
			Retry:    resp.StatusCode >= 500,
		}
	}

	// Parse models response
	var modelsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse models: %w", err)
	}

	// Filter for generateContent-capable models
	var models []string
	for _, m := range modelsResp.Models {
		// Only include models that support generateContent (typically start with "models/gemini")
		if strings.Contains(m.Name, "gemini") {
			// Remove "models/" prefix if present
			name := strings.TrimPrefix(m.Name, "models/")
			models = append(models, name)
		}
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no compatible models found")
	}

	return models, nil
}

// Close is a no-op for HTTP-based implementation
func (g *GeminiAdapter) Close() error {
	return nil
}
