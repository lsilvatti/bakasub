// Package tokenizer provides BPE token estimation for AI cost calculations.
// This is a simplified estimator based on common BPE tokenization patterns.
package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Estimator provides token count estimation for various AI models
type Estimator struct {
	// Average characters per token varies by language
	charsPerToken float64
}

// NewEstimator creates a new token estimator
func NewEstimator() *Estimator {
	return &Estimator{
		// Average for English text (GPT models)
		// Most languages average 3-4 chars per token
		charsPerToken: 4.0,
	}
}

// EstimateTokens returns an estimated token count for the given text
func (e *Estimator) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Count using multiple heuristics and average them
	estimates := []int{
		e.estimateByChars(text),
		e.estimateByWords(text),
		e.estimateByRunes(text),
	}

	total := 0
	for _, est := range estimates {
		total += est
	}

	return total / len(estimates)
}

// estimateByChars uses character count divided by average chars per token
func (e *Estimator) estimateByChars(text string) int {
	return int(float64(len(text)) / e.charsPerToken)
}

// estimateByWords uses word count as a rough approximation
// On average, 1 word ≈ 1.3-1.5 tokens for English
func (e *Estimator) estimateByWords(text string) int {
	words := strings.Fields(text)
	return int(float64(len(words)) * 1.4)
}

// estimateByRunes handles unicode better for non-ASCII languages
func (e *Estimator) estimateByRunes(text string) int {
	runes := []rune(text)

	// Count token-like segments
	tokenCount := 0
	inWord := false

	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				inWord = true
				tokenCount++
			}
		} else {
			inWord = false
			// Punctuation and special characters often become separate tokens
			if !unicode.IsSpace(r) {
				tokenCount++
			}
		}
	}

	// Add tokens for special patterns
	// ASS tags are typically 1-2 tokens each
	assTagPattern := regexp.MustCompile(`\{[^}]*\}`)
	assTags := assTagPattern.FindAllString(text, -1)
	tokenCount += len(assTags)

	return tokenCount
}

// EstimateBatch estimates tokens for multiple text lines
func (e *Estimator) EstimateBatch(lines []string) int {
	total := 0
	for _, line := range lines {
		total += e.EstimateTokens(line)
	}
	return total
}

// CostEstimate contains cost estimation results
type CostEstimate struct {
	InputTokens   int
	OutputTokens  int // Estimated as ~80% of input for translation
	TotalTokens   int
	CostUSD       float64
	FormattedCost string
}

// Pricing contains model pricing information
type Pricing struct {
	InputPer1M  float64 // USD per 1M input tokens
	OutputPer1M float64 // USD per 1M output tokens
}

// Common model pricing (as of 2024)
var ModelPricing = map[string]Pricing{
	// OpenAI
	"gpt-4o":        {InputPer1M: 2.50, OutputPer1M: 10.00},
	"gpt-4o-mini":   {InputPer1M: 0.15, OutputPer1M: 0.60},
	"gpt-4-turbo":   {InputPer1M: 10.00, OutputPer1M: 30.00},
	"gpt-3.5-turbo": {InputPer1M: 0.50, OutputPer1M: 1.50},

	// Google
	"gemini-1.5-flash": {InputPer1M: 0.075, OutputPer1M: 0.30},
	"gemini-1.5-pro":   {InputPer1M: 1.25, OutputPer1M: 5.00},
	"gemini-2.0-flash": {InputPer1M: 0.10, OutputPer1M: 0.40},

	// Anthropic (via OpenRouter)
	"claude-3-opus":     {InputPer1M: 15.00, OutputPer1M: 75.00},
	"claude-3-sonnet":   {InputPer1M: 3.00, OutputPer1M: 15.00},
	"claude-3-haiku":    {InputPer1M: 0.25, OutputPer1M: 1.25},
	"claude-3.5-sonnet": {InputPer1M: 3.00, OutputPer1M: 15.00},

	// Free/Open Source (OpenRouter free tier)
	"free":          {InputPer1M: 0.00, OutputPer1M: 0.00},
	"llama-3.3-70b": {InputPer1M: 0.00, OutputPer1M: 0.00},
	"qwen-2.5-72b":  {InputPer1M: 0.00, OutputPer1M: 0.00},

	// Default fallback
	"default": {InputPer1M: 0.10, OutputPer1M: 0.40},
}

// EstimateCost calculates the estimated cost for translating the given text
func (e *Estimator) EstimateCost(lines []string, model string) CostEstimate {
	inputTokens := e.EstimateBatch(lines)

	// Output tokens are typically ~80% of input for translation
	// (translated text is usually shorter or similar length)
	outputTokens := int(float64(inputTokens) * 0.8)

	// Add system prompt overhead (~500 tokens typically)
	inputTokens += 500

	pricing, ok := ModelPricing[normalizeModelName(model)]
	if !ok {
		pricing = ModelPricing["default"]
	}

	costUSD := (float64(inputTokens) * pricing.InputPer1M / 1000000) +
		(float64(outputTokens) * pricing.OutputPer1M / 1000000)

	return CostEstimate{
		InputTokens:   inputTokens,
		OutputTokens:  outputTokens,
		TotalTokens:   inputTokens + outputTokens,
		CostUSD:       costUSD,
		FormattedCost: formatCost(costUSD),
	}
}

// normalizeModelName extracts the base model name for pricing lookup
func normalizeModelName(model string) string {
	model = strings.ToLower(model)

	// Check for free models
	if strings.Contains(model, "free") || strings.Contains(model, ":free") {
		return "free"
	}

	// Try to match known model patterns
	patterns := []struct {
		pattern string
		key     string
	}{
		{"gpt-4o-mini", "gpt-4o-mini"},
		{"gpt-4o", "gpt-4o"},
		{"gpt-4-turbo", "gpt-4-turbo"},
		{"gpt-3.5", "gpt-3.5-turbo"},
		{"gemini-1.5-flash", "gemini-1.5-flash"},
		{"gemini-1.5-pro", "gemini-1.5-pro"},
		{"gemini-2.0-flash", "gemini-2.0-flash"},
		{"claude-3-opus", "claude-3-opus"},
		{"claude-3.5-sonnet", "claude-3.5-sonnet"},
		{"claude-3-sonnet", "claude-3-sonnet"},
		{"claude-3-haiku", "claude-3-haiku"},
		{"llama-3.3", "llama-3.3-70b"},
		{"llama-3-70b", "llama-3.3-70b"},
		{"qwen-2.5", "qwen-2.5-72b"},
	}

	for _, p := range patterns {
		if strings.Contains(model, p.pattern) {
			return p.key
		}
	}

	return "default"
}

// formatCost formats a USD cost for display
func formatCost(cost float64) string {
	if cost < 0.01 {
		return "$0.00"
	}
	return fmt.Sprintf("$%.2f", cost)
}
