package tokenizer

import (
	"testing"
)

func TestNewEstimator(t *testing.T) {
	estimator := NewEstimator()

	if estimator == nil {
		t.Fatal("NewEstimator returned nil")
	}

	if estimator.charsPerToken != 4.0 {
		t.Errorf("expected charsPerToken 4.0, got %f", estimator.charsPerToken)
	}
}

func TestEstimateTokens(t *testing.T) {
	estimator := NewEstimator()

	tests := []struct {
		name     string
		text     string
		minToken int
		maxToken int
	}{
		{"empty string", "", 0, 0},
		{"single word", "hello", 1, 5},
		{"sentence", "Hello, how are you today?", 3, 15},
		{"long text", "This is a longer piece of text that contains multiple sentences. It should produce more tokens.", 15, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := estimator.EstimateTokens(tt.text)
			if tokens < tt.minToken || tokens > tt.maxToken {
				t.Errorf("token count %d not in expected range [%d, %d]", tokens, tt.minToken, tt.maxToken)
			}
		})
	}
}

func TestEstimateByChars(t *testing.T) {
	estimator := NewEstimator()

	// 16 chars / 4 chars per token = 4 tokens
	text := "1234567890123456"
	tokens := estimator.estimateByChars(text)

	if tokens != 4 {
		t.Errorf("expected 4 tokens, got %d", tokens)
	}
}

func TestEstimateByWords(t *testing.T) {
	estimator := NewEstimator()

	text := "one two three four five"
	tokens := estimator.estimateByWords(text)

	// 5 words * 1.4 = 7
	if tokens != 7 {
		t.Errorf("expected 7 tokens, got %d", tokens)
	}
}

func TestEstimateByRunes(t *testing.T) {
	estimator := NewEstimator()

	text := "Hello World!"
	tokens := estimator.estimateByRunes(text)

	// Should count word segments and punctuation
	if tokens <= 0 {
		t.Error("should produce positive token count")
	}
}

func TestEstimateBatch(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{
		"Hello world",
		"How are you",
		"This is a test",
	}

	total := estimator.EstimateBatch(lines)

	if total <= 0 {
		t.Error("batch should produce positive token count")
	}

	// Should be sum of individual estimates
	individual := 0
	for _, line := range lines {
		individual += estimator.EstimateTokens(line)
	}

	if total != individual {
		t.Errorf("batch total %d should equal sum of individual %d", total, individual)
	}
}

func TestEstimateBatchEmpty(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{}
	total := estimator.EstimateBatch(lines)

	if total != 0 {
		t.Errorf("expected 0 tokens for empty batch, got %d", total)
	}
}

func TestEstimateCost(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{
		"Line one for translation",
		"Line two for translation",
		"Line three for translation",
	}

	estimate := estimator.EstimateCost(lines, "gpt-4o")

	if estimate.InputTokens <= 0 {
		t.Error("InputTokens should be positive")
	}

	if estimate.OutputTokens <= 0 {
		t.Error("OutputTokens should be positive")
	}

	if estimate.TotalTokens != estimate.InputTokens+estimate.OutputTokens {
		t.Error("TotalTokens should equal InputTokens + OutputTokens")
	}

	if estimate.CostUSD < 0 {
		t.Error("CostUSD should not be negative")
	}

	if estimate.FormattedCost == "" {
		t.Error("FormattedCost should not be empty")
	}
}

func TestEstimateCostFreeModel(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{"Test line"}
	estimate := estimator.EstimateCost(lines, "llama-3.3-70b")

	if estimate.CostUSD != 0 {
		t.Errorf("expected 0 cost for free model, got %f", estimate.CostUSD)
	}
}

func TestEstimateCostUnknownModel(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{"Test line"}
	estimate := estimator.EstimateCost(lines, "unknown-model-xyz")

	// Should fall back to default pricing
	if estimate.CostUSD < 0 {
		t.Error("should use default pricing for unknown model")
	}
}

func TestNormalizeModelName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gpt-4o", "gpt-4o"},
		{"GPT-4O", "gpt-4o"},
		{"gpt-4o-mini", "gpt-4o-mini"},
		{"gemini-1.5-flash", "gemini-1.5-flash"},
		{"claude-3.5-sonnet", "claude-3.5-sonnet"},
		{"some-model:free", "free"},
		{"unknown", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeModelName(tt.input)
			// Check that it returns a valid key
			if _, ok := ModelPricing[result]; !ok {
				t.Errorf("normalized model %q not found in ModelPricing", result)
			}
		})
	}
}

func TestModelPricing(t *testing.T) {
	// Verify all expected models are in pricing table
	expectedModels := []string{
		"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-3.5-turbo",
		"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash",
		"claude-3-opus", "claude-3-sonnet", "claude-3-haiku", "claude-3.5-sonnet",
		"free", "llama-3.3-70b", "qwen-2.5-72b", "default",
	}

	for _, model := range expectedModels {
		if _, ok := ModelPricing[model]; !ok {
			t.Errorf("model %q missing from ModelPricing", model)
		}
	}
}

func TestPricingStruct(t *testing.T) {
	pricing := Pricing{
		InputPer1M:  2.50,
		OutputPer1M: 10.00,
	}

	if pricing.InputPer1M != 2.50 {
		t.Errorf("unexpected InputPer1M: %f", pricing.InputPer1M)
	}

	if pricing.OutputPer1M != 10.00 {
		t.Errorf("unexpected OutputPer1M: %f", pricing.OutputPer1M)
	}
}

func TestCostEstimateStruct(t *testing.T) {
	estimate := CostEstimate{
		InputTokens:   1000,
		OutputTokens:  800,
		TotalTokens:   1800,
		CostUSD:       0.05,
		FormattedCost: "$0.05",
	}

	if estimate.InputTokens != 1000 {
		t.Errorf("unexpected InputTokens: %d", estimate.InputTokens)
	}

	if estimate.OutputTokens != 800 {
		t.Errorf("unexpected OutputTokens: %d", estimate.OutputTokens)
	}

	if estimate.TotalTokens != 1800 {
		t.Errorf("unexpected TotalTokens: %d", estimate.TotalTokens)
	}

	if estimate.CostUSD != 0.05 {
		t.Errorf("unexpected CostUSD: %f", estimate.CostUSD)
	}

	if estimate.FormattedCost != "$0.05" {
		t.Errorf("unexpected FormattedCost: %q", estimate.FormattedCost)
	}
}

func TestEstimateWithASSTags(t *testing.T) {
	estimator := NewEstimator()

	textWithTags := `{\an8}Hello World{\b1}`
	textWithoutTags := "Hello World"

	tokensWithTags := estimator.EstimateTokens(textWithTags)
	tokensWithoutTags := estimator.EstimateTokens(textWithoutTags)

	// Text with ASS tags should produce more tokens
	if tokensWithTags <= tokensWithoutTags {
		t.Error("text with ASS tags should produce more tokens")
	}
}

func TestEstimateUnicode(t *testing.T) {
	estimator := NewEstimator()

	// Japanese text
	japaneseText := "こんにちは世界"
	tokens := estimator.EstimateTokens(japaneseText)

	if tokens <= 0 {
		t.Error("should estimate tokens for unicode text")
	}

	// Portuguese with accents
	portugueseText := "Olá, como você está? Tudo bem?"
	ptTokens := estimator.EstimateTokens(portugueseText)

	if ptTokens <= 0 {
		t.Error("should estimate tokens for Portuguese text")
	}
}

func TestEstimateSystemPromptOverhead(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{"Test"}
	estimate := estimator.EstimateCost(lines, "gpt-4o")

	// Should include ~500 tokens for system prompt
	if estimate.InputTokens < 500 {
		t.Error("should include system prompt overhead in input tokens")
	}
}

func TestOutputTokenRatio(t *testing.T) {
	estimator := NewEstimator()

	lines := []string{"This is a test line for translation purposes"}
	estimate := estimator.EstimateCost(lines, "gpt-4o")

	// Output should be ~80% of input (minus system prompt overhead)
	inputWithoutOverhead := estimate.InputTokens - 500
	expectedOutput := int(float64(inputWithoutOverhead) * 0.8)

	// Allow some tolerance
	if estimate.OutputTokens < expectedOutput-5 || estimate.OutputTokens > expectedOutput+5 {
		t.Errorf("output tokens %d not close to expected %d", estimate.OutputTokens, expectedOutput)
	}
}
