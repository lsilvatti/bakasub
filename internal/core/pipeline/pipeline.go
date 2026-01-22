package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lsilvatti/bakasub/internal/core/ai"
	"github.com/lsilvatti/bakasub/internal/core/db"
	"github.com/lsilvatti/bakasub/internal/core/media"
	"github.com/lsilvatti/bakasub/internal/core/parser"
)

// Pipeline orchestrates the translation workflow
type Pipeline struct {
	Provider         ai.LLMProvider
	Cache            *db.Cache
	Config           *PipelineConfig
	ResumeState      *ResumeState
	LogCallback      func(string)
	ProgressCallback func(current, total int)
}

// PipelineConfig holds pipeline configuration
type PipelineConfig struct {
	InputPath         string
	OutputPath        string
	SourceLang        string
	TargetLang        string
	Model             string
	Temperature       float64
	BatchSize         int
	RemoveHI          bool
	Glossary          map[string]string
	SystemPrompt      string
	SlidingWindowSize int // Number of lines for context
}

// ResumeState holds state for smart resume
type ResumeState struct {
	FilePath         string
	CompletedBatches int
	TotalBatches     int
	TranslatedLines  []parser.SubtitleLine
	Timestamp        time.Time
}

// TranslationBatch represents a batch of lines to translate
type TranslationBatch struct {
	Lines        []parser.SubtitleLine
	ContextLines []parser.SubtitleLine // Sliding window context
	BatchIndex   int
	TotalBatches int
}

// New creates a new pipeline instance
func New(provider ai.LLMProvider, cache *db.Cache, config *PipelineConfig) *Pipeline {
	if config.SlidingWindowSize == 0 {
		config.SlidingWindowSize = 3
	}
	if config.BatchSize == 0 {
		config.BatchSize = 50
	}

	return &Pipeline{
		Provider: provider,
		Cache:    cache,
		Config:   config,
	}
}

// Execute runs the full translation pipeline
func (p *Pipeline) Execute(ctx context.Context) error {
	p.log("Starting translation pipeline...")

	// Step 1: Extract subtitle track
	p.log("Extracting subtitle track...")
	tempSubPath := filepath.Join(os.TempDir(), "bakasub_temp.ass")
	defer os.Remove(tempSubPath)

	if err := media.ExtractSubtitleTrack(p.Config.InputPath, 2, tempSubPath); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Step 2: Parse subtitles
	p.log("Parsing subtitle file...")
	subFile, err := parser.ParseFile(tempSubPath)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}
	p.log(fmt.Sprintf("Found %d lines", subFile.LineCount))

	// Step 3: Preprocessing (remove HI tags if enabled)
	if p.Config.RemoveHI {
		p.log("Removing hearing impaired tags...")
		for i := range subFile.Lines {
			subFile.Lines[i].Text = parser.RemoveHearingImpairedTags(subFile.Lines[i].Text)
		}
	}

	// Step 4: Batch and translate
	batches := parser.BatchLines(subFile.Lines, p.Config.BatchSize)
	p.log(fmt.Sprintf("Split into %d batches", len(batches)))

	translatedLines := []parser.SubtitleLine{}
	var contextWindow []parser.SubtitleLine

	startBatch := 0
	if p.ResumeState != nil {
		startBatch = p.ResumeState.CompletedBatches
		translatedLines = p.ResumeState.TranslatedLines
		p.log(fmt.Sprintf("Resuming from batch %d/%d", startBatch+1, len(batches)))
	}

	for i := startBatch; i < len(batches); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		p.log(fmt.Sprintf("Processing batch %d/%d...", i+1, len(batches)))
		p.progress(i+1, len(batches))

		batch := TranslationBatch{
			Lines:        batches[i],
			ContextLines: contextWindow,
			BatchIndex:   i,
			TotalBatches: len(batches),
		}

		translated, err := p.translateBatch(ctx, batch)
		if err != nil {
			return fmt.Errorf("batch %d failed: %w", i+1, err)
		}

		translatedLines = append(translatedLines, translated...)

		// Update sliding window with last N lines
		if len(translated) >= p.Config.SlidingWindowSize {
			contextWindow = translated[len(translated)-p.Config.SlidingWindowSize:]
		} else {
			contextWindow = translated
		}

		// Save resume state after each batch
		if err := p.saveResumeState(i+1, len(batches), translatedLines); err != nil {
			p.log(fmt.Sprintf("Warning: Failed to save resume state: %v", err))
		}
	}

	// Step 5: Reassemble subtitle file
	p.log("Reassembling subtitle file...")
	var content string
	if subFile.Format == "ass" {
		content = parser.ReassembleASS(subFile.Header, translatedLines)
	} else {
		content = parser.ReassembleSRT(translatedLines)
	}

	translatedPath := filepath.Join(os.TempDir(), "bakasub_translated"+filepath.Ext(tempSubPath))
	if err := os.WriteFile(translatedPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	defer os.Remove(translatedPath)

	// Step 6: Mux back into video
	p.log("Muxing translated subtitle...")
	if err := media.MuxSubtitle(p.Config.InputPath, translatedPath, p.Config.OutputPath); err != nil {
		return fmt.Errorf("mux failed: %w", err)
	}

	// Clean up resume state
	p.clearResumeState()

	p.log("Translation complete!")
	return nil
}

// translateBatch translates a single batch with anti-desync protocol
func (p *Pipeline) translateBatch(ctx context.Context, batch TranslationBatch) ([]parser.SubtitleLine, error) {
	langPair := fmt.Sprintf("%s->%s", p.Config.SourceLang, p.Config.TargetLang)

	// Check cache for each line
	cachedCount := 0
	translatedLines := make([]parser.SubtitleLine, len(batch.Lines))
	needsTranslation := []int{}

	for i, line := range batch.Lines {
		if cached, found := p.Cache.GetExactMatch(line.Text, langPair); found {
			translatedLines[i] = line
			translatedLines[i].Text = cached
			cachedCount++
		} else if cached, found := p.Cache.GetFuzzyMatch(line.Text, langPair, 0.95); found {
			translatedLines[i] = line
			translatedLines[i].Text = cached.TranslatedText
			cachedCount++
		} else {
			needsTranslation = append(needsTranslation, i)
		}
	}

	if cachedCount > 0 {
		p.log(fmt.Sprintf("  Cache hit: %d/%d lines", cachedCount, len(batch.Lines)))
	}

	if len(needsTranslation) == 0 {
		return translatedLines, nil
	}

	// Build system prompt with context and glossary
	systemPrompt := p.buildSystemPrompt(batch.ContextLines)

	// Prepare payload for AI
	payload := []ai.Line{}
	for _, idx := range needsTranslation {
		payload = append(payload, ai.Line{ID: idx, Text: batch.Lines[idx].Text})
	}

	// Send to AI provider
	response, err := p.Provider.SendBatch(ctx, payload, systemPrompt)
	if err != nil {
		return nil, err
	}

	// Anti-desync check: verify count matches
	if len(response) != len(needsTranslation) {
		p.log(fmt.Sprintf("  DESYNC DETECTED: Expected %d, got %d", len(needsTranslation), len(response)))
		// Trigger split strategy (not implemented in this demo)
		return nil, fmt.Errorf("desync: count mismatch")
	}

	// Apply translations
	for _, resp := range response {
		if resp.ID >= 0 && resp.ID < len(translatedLines) {
			translatedLines[resp.ID].Text = resp.Text
			// Cache the translation
			p.Cache.SaveTranslation(batch.Lines[resp.ID].Text, resp.Text, langPair)
		}
	}

	return translatedLines, nil
}

// buildSystemPrompt creates system prompt with sliding window context and glossary
func (p *Pipeline) buildSystemPrompt(contextLines []parser.SubtitleLine) string {
	prompt := p.Config.SystemPrompt

	// Inject glossary
	if len(p.Config.Glossary) > 0 {
		glossaryText := "\n\nGlossary (preserve these terms):\n"
		for orig, trans := range p.Config.Glossary {
			glossaryText += fmt.Sprintf("- %s -> %s\n", orig, trans)
		}
		prompt = strings.Replace(prompt, "{{glossary}}", glossaryText, 1)
	} else {
		prompt = strings.Replace(prompt, "{{glossary}}", "", 1)
	}

	// Add sliding window context
	if len(contextLines) > 0 {
		contextText := "\n\nPrevious lines (for context only, do not translate):\n"
		for _, line := range contextLines {
			contextText += fmt.Sprintf("- %s\n", line.Text)
		}
		prompt += contextText
	}

	return prompt
}

func (p *Pipeline) log(msg string) {
	if p.LogCallback != nil {
		p.LogCallback(msg)
	}
}

func (p *Pipeline) progress(current, total int) {
	if p.ProgressCallback != nil {
		p.ProgressCallback(current, total)
	}
}

// saveResumeState writes .bakasub.temp file
func (p *Pipeline) saveResumeState(completed, total int, lines []parser.SubtitleLine) error {
	state := ResumeState{
		FilePath:         p.Config.InputPath,
		CompletedBatches: completed,
		TotalBatches:     total,
		TranslatedLines:  lines,
		Timestamp:        time.Now(),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	tempPath := filepath.Join(filepath.Dir(p.Config.InputPath), ".bakasub.temp")
	return os.WriteFile(tempPath, data, 0644)
}

func (p *Pipeline) clearResumeState() {
	tempPath := filepath.Join(filepath.Dir(p.Config.InputPath), ".bakasub.temp")
	os.Remove(tempPath)
}

// LoadResumeState loads .bakasub.temp if it exists
func LoadResumeState(videoPath string) (*ResumeState, error) {
	tempPath := filepath.Join(filepath.Dir(videoPath), ".bakasub.temp")
	data, err := os.ReadFile(tempPath)
	if err != nil {
		return nil, err
	}

	var state ResumeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
