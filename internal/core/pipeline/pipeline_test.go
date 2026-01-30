package pipeline

import (
	"context"
	"testing"
	"time"

	"github.com/lsilvatti/bakasub/internal/core/ai"
	"github.com/lsilvatti/bakasub/internal/core/parser"
)

// MockProvider implements LLMProvider for testing
type MockProvider struct {
	Lines       []ai.Line
	Error       error
	CallCount   int
	LastPayload []ai.Line
	LastPrompt  string
}

func (m *MockProvider) SendBatch(ctx context.Context, payload []ai.Line, systemPrompt string) ([]ai.Line, error) {
	m.CallCount++
	m.LastPayload = payload
	m.LastPrompt = systemPrompt

	if m.Error != nil {
		return nil, m.Error
	}

	// Return translated lines
	result := make([]ai.Line, len(payload))
	for i, line := range payload {
		result[i] = ai.Line{
			ID:   line.ID,
			Text: "Translated: " + line.Text,
		}
	}
	return result, nil
}

func (m *MockProvider) ValidateKey(ctx context.Context) bool {
	return true
}

func (m *MockProvider) ListModels(ctx context.Context) ([]string, error) {
	return []string{"test-model"}, nil
}

// TestNew tests pipeline creation
func TestNew(t *testing.T) {
	provider := &MockProvider{}
	config := &PipelineConfig{
		BatchSize: 50,
	}

	pipeline := New(provider, nil, config)
	if pipeline == nil {
		t.Fatal("New returned nil")
	}

	if pipeline.Provider != provider {
		t.Error("provider not set correctly")
	}

	if pipeline.Config.BatchSize != 50 {
		t.Errorf("expected BatchSize 50, got %d", pipeline.Config.BatchSize)
	}

	if pipeline.Config.SlidingWindowSize != 3 {
		t.Errorf("expected default SlidingWindowSize 3, got %d", pipeline.Config.SlidingWindowSize)
	}
}

// TestNewPipelineDefaultBatchSize tests default batch size
func TestNewPipelineDefaultBatchSize(t *testing.T) {
	config := &PipelineConfig{
		BatchSize: 0, // Should default to 50
	}

	pipeline := New(nil, nil, config)
	if pipeline.Config.BatchSize != 50 {
		t.Errorf("expected default BatchSize 50, got %d", pipeline.Config.BatchSize)
	}
}

// TestNewPipelineDefaultSlidingWindow tests default sliding window size
func TestNewPipelineDefaultSlidingWindow(t *testing.T) {
	config := &PipelineConfig{
		SlidingWindowSize: 0, // Should default to 3
	}

	pipeline := New(nil, nil, config)
	if pipeline.Config.SlidingWindowSize != 3 {
		t.Errorf("expected default SlidingWindowSize 3, got %d", pipeline.Config.SlidingWindowSize)
	}
}

// TestPipelineConfigStruct tests PipelineConfig structure
func TestPipelineConfigStruct(t *testing.T) {
	config := PipelineConfig{
		InputPath:         "/input/file.mkv",
		OutputPath:        "/output/file.mkv",
		SourceLang:        "en",
		TargetLang:        "pt-br",
		Model:             "gpt-4o",
		Temperature:       0.5,
		BatchSize:         30,
		RemoveHI:          true,
		Glossary:          map[string]string{"Nakama": "Companheiro"},
		SystemPrompt:      "Translate this",
		SlidingWindowSize: 5,
	}

	if config.InputPath != "/input/file.mkv" {
		t.Errorf("unexpected InputPath: %q", config.InputPath)
	}

	if config.OutputPath != "/output/file.mkv" {
		t.Errorf("unexpected OutputPath: %q", config.OutputPath)
	}

	if config.SourceLang != "en" {
		t.Errorf("unexpected SourceLang: %q", config.SourceLang)
	}

	if config.TargetLang != "pt-br" {
		t.Errorf("unexpected TargetLang: %q", config.TargetLang)
	}

	if !config.RemoveHI {
		t.Error("RemoveHI should be true")
	}

	if config.Glossary["Nakama"] != "Companheiro" {
		t.Error("Glossary entry not set correctly")
	}

	if config.SlidingWindowSize != 5 {
		t.Errorf("unexpected SlidingWindowSize: %d", config.SlidingWindowSize)
	}
}

// TestResumeStateStruct tests ResumeState structure
func TestResumeStateStruct(t *testing.T) {
	state := ResumeState{
		FilePath:         "/path/to/file.mkv",
		CompletedBatches: 5,
		TotalBatches:     10,
		TranslatedLines: []parser.SubtitleLine{
			{Index: 0, Text: "Line 1"},
			{Index: 1, Text: "Line 2"},
		},
		Timestamp: time.Now(),
	}

	if state.FilePath != "/path/to/file.mkv" {
		t.Errorf("unexpected FilePath: %q", state.FilePath)
	}

	if state.CompletedBatches != 5 {
		t.Errorf("unexpected CompletedBatches: %d", state.CompletedBatches)
	}

	if state.TotalBatches != 10 {
		t.Errorf("unexpected TotalBatches: %d", state.TotalBatches)
	}

	if len(state.TranslatedLines) != 2 {
		t.Errorf("expected 2 translated lines, got %d", len(state.TranslatedLines))
	}
}

// TestTranslationBatchStruct tests TranslationBatch structure
func TestTranslationBatchStruct(t *testing.T) {
	batch := TranslationBatch{
		Lines: []parser.SubtitleLine{
			{Index: 0, Text: "Line 1"},
			{Index: 1, Text: "Line 2"},
		},
		ContextLines: []parser.SubtitleLine{
			{Index: -3, Text: "Context 1"},
		},
		BatchIndex:   5,
		TotalBatches: 10,
	}

	if len(batch.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(batch.Lines))
	}

	if len(batch.ContextLines) != 1 {
		t.Errorf("expected 1 context line, got %d", len(batch.ContextLines))
	}

	if batch.BatchIndex != 5 {
		t.Errorf("unexpected BatchIndex: %d", batch.BatchIndex)
	}

	if batch.TotalBatches != 10 {
		t.Errorf("unexpected TotalBatches: %d", batch.TotalBatches)
	}
}

// TestNewNERScanner tests NER scanner wrapper
func TestNewNERScanner(t *testing.T) {
	scanner := NewNERScanner()
	if scanner == nil {
		t.Fatal("NewNERScanner returned nil")
	}
}

// TestPipelineCallbacks tests log and progress callbacks
func TestPipelineCallbacks(t *testing.T) {
	config := &PipelineConfig{
		BatchSize: 10,
	}

	pipeline := New(nil, nil, config)
	logMessages := []string{}

	pipeline.LogCallback = func(msg string) {
		logMessages = append(logMessages, msg)
	}

	progressCalls := 0
	pipeline.ProgressCallback = func(current, total int) {
		progressCalls++
	}

	// Test log callback
	pipeline.log("Test message")
	if len(logMessages) != 1 || logMessages[0] != "Test message" {
		t.Error("LogCallback not working correctly")
	}

	// Test progress callback
	pipeline.progress(1, 10)
	if progressCalls != 1 {
		t.Error("ProgressCallback not working correctly")
	}
}

// TestPipelineLogNoCallback tests log with nil callback
func TestPipelineLogNoCallback(t *testing.T) {
	config := &PipelineConfig{}
	pipeline := New(nil, nil, config)

	// Should not panic with nil callback
	pipeline.log("Test message")
	pipeline.progress(1, 1)
}

// TestPipelineWithResume tests pipeline with resume state
func TestPipelineWithResume(t *testing.T) {
	config := &PipelineConfig{
		BatchSize: 10,
	}

	pipeline := New(nil, nil, config)
	pipeline.ResumeState = &ResumeState{
		FilePath:         "/test/file.mkv",
		CompletedBatches: 3,
		TotalBatches:     10,
		TranslatedLines: []parser.SubtitleLine{
			{Index: 0, Text: "Already translated"},
		},
	}

	if pipeline.ResumeState == nil {
		t.Error("ResumeState should be set")
	}

	if pipeline.ResumeState.CompletedBatches != 3 {
		t.Errorf("unexpected CompletedBatches: %d", pipeline.ResumeState.CompletedBatches)
	}
}

// TestBuildSystemPrompt tests system prompt building
func TestBuildSystemPrompt(t *testing.T) {
	config := &PipelineConfig{
		SystemPrompt: "You are a translator. Glossary: {{glossary}}",
		TargetLang:   "pt-br",
		Glossary: map[string]string{
			"Nakama": "Companheiro",
			"Sensei": "Mestre",
		},
	}

	pipeline := New(nil, nil, config)
	contextLines := []parser.SubtitleLine{
		{Index: 0, Text: "Previous line 1"},
		{Index: 1, Text: "Previous line 2"},
	}

	prompt := pipeline.buildSystemPrompt(contextLines)

	// Prompt should not be empty
	if prompt == "" {
		t.Error("buildSystemPrompt returned empty string")
	}
}
