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
	"github.com/lsilvatti/bakasub/internal/core/linter"
	"github.com/lsilvatti/bakasub/internal/core/media"
	"github.com/lsilvatti/bakasub/internal/core/ner"
	"github.com/lsilvatti/bakasub/internal/core/parser"
)

// NewNERScanner creates a new NER scanner (wrapper for ner package)
func NewNERScanner() *ner.Scanner {
	return ner.NewScanner()
}

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
	SlidingWindowSize int    // Number of lines for context
	TrackID           int    // Subtitle track ID to extract (-1 for auto-detect)
	MuxMode           string // "replace" or "new-file"
	BackupOriginal    bool   // Create backup before replace
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

	// Determine track ID to use
	trackID := p.Config.TrackID
	if trackID < 0 {
		// Auto-detect: find first subtitle track
		p.log("Auto-detecting subtitle track...")
		fileInfo, err := media.Analyze(p.Config.InputPath)
		if err != nil {
			return fmt.Errorf("failed to analyze file: %w", err)
		}
		subTracks := media.GetSubtitleTracks(fileInfo)
		if len(subTracks) == 0 {
			return fmt.Errorf("no subtitle tracks found in file")
		}
		trackID = subTracks[0].ID
		p.log(fmt.Sprintf("Using subtitle track %d (%s)", trackID, subTracks[0].Language))
	}

	// Step 1: Extract subtitle track
	p.log("Extracting subtitle track...")
	tempSubPath := filepath.Join(os.TempDir(), "bakasub_temp.ass")
	defer os.Remove(tempSubPath)

	if err := media.ExtractSubtitleTrack(p.Config.InputPath, trackID, tempSubPath); err != nil {
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

	// Step 3.5: NER Scan for Volatile Glossary (if no project glossary provided)
	if len(p.Config.Glossary) == 0 {
		p.log("Scanning for named entities (Volatile Glossary)...")
		nerScanner := NewNERScanner()
		entities := nerScanner.ScanLines(subFile.Lines)
		if len(entities) > 0 {
			p.log(fmt.Sprintf("Detected %d potential entities", len(entities)))
			// Create volatile glossary from detected entities
			p.Config.Glossary = make(map[string]string)
			for _, e := range entities {
				// Only add high-confidence entities (preserve original form)
				if e.Confidence >= 0.7 {
					p.Config.Glossary[e.Text] = e.Text
				}
			}
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

	// Handle replace mode - use temp file to avoid overwriting source
	outputPath := p.Config.OutputPath
	isReplaceMode := p.Config.MuxMode == "replace" || outputPath == p.Config.InputPath
	var tempOutputPath string

	if isReplaceMode {
		// Create backup if enabled
		if p.Config.BackupOriginal {
			backupPath := p.Config.InputPath + ".bak"
			p.log(fmt.Sprintf("Creating backup: %s", filepath.Base(backupPath)))
			input, err := os.ReadFile(p.Config.InputPath)
			if err != nil {
				return fmt.Errorf("failed to read input for backup: %w", err)
			}
			if err := os.WriteFile(backupPath, input, 0644); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		}

		// Use temp file for output to avoid same input/output error
		tempOutputPath = filepath.Join(os.TempDir(), "bakasub_mux_temp.mkv")
		outputPath = tempOutputPath
		defer os.Remove(tempOutputPath)
	}

	if err := media.MuxSubtitle(p.Config.InputPath, translatedPath, outputPath); err != nil {
		return fmt.Errorf("mux failed: %w", err)
	}

	// In replace mode, move temp file to original location
	if isReplaceMode && tempOutputPath != "" {
		p.log("Replacing original file...")
		// Read the temp file
		tempData, err := os.ReadFile(tempOutputPath)
		if err != nil {
			return fmt.Errorf("failed to read temp output: %w", err)
		}
		// Write to original location
		if err := os.WriteFile(p.Config.InputPath, tempData, 0644); err != nil {
			return fmt.Errorf("failed to replace original file: %w", err)
		}
	}

	// Clean up resume state
	p.clearResumeState()

	p.log("Translation complete!")
	return nil
}

// translateBatch translates a single batch with anti-desync protocol
func (p *Pipeline) translateBatch(ctx context.Context, batch TranslationBatch) ([]parser.SubtitleLine, error) {
	return p.translateBatchWithRetry(ctx, batch, 0)
}

// translateBatchWithRetry implements self-healing split strategy
// maxDepth prevents infinite recursion (max 3 levels: 50 -> 25 -> 12 -> 6)
func (p *Pipeline) translateBatchWithRetry(ctx context.Context, batch TranslationBatch, depth int) ([]parser.SubtitleLine, error) {
	const maxRetryDepth = 3
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

	// Handle errors or desync with self-healing split strategy
	if err != nil || len(response) != len(needsTranslation) {
		if err != nil {
			p.log(fmt.Sprintf("  AI ERROR: %v", err))
		} else {
			p.log(fmt.Sprintf("  DESYNC DETECTED: Expected %d, got %d", len(needsTranslation), len(response)))
		}

		// If we can still split, try self-healing
		if depth < maxRetryDepth && len(batch.Lines) > 1 {
			p.log(fmt.Sprintf("  └─ Engaging Self-Healing Protocol (Split Strategy, depth=%d)", depth+1))

			// Split batch in half
			mid := len(batch.Lines) / 2

			batchA := TranslationBatch{
				Lines:        batch.Lines[:mid],
				ContextLines: batch.ContextLines,
				BatchIndex:   batch.BatchIndex,
				TotalBatches: batch.TotalBatches,
			}

			batchB := TranslationBatch{
				Lines:        batch.Lines[mid:],
				ContextLines: batch.Lines[max(0, mid-p.Config.SlidingWindowSize):mid], // Use end of A as context for B
				BatchIndex:   batch.BatchIndex,
				TotalBatches: batch.TotalBatches,
			}

			p.log(fmt.Sprintf("  └─ Split %da (Lines 1-%d) processing...", batch.BatchIndex+1, mid))
			resultA, errA := p.translateBatchWithRetry(ctx, batchA, depth+1)
			if errA != nil {
				return nil, fmt.Errorf("split A failed: %w", errA)
			}

			p.log(fmt.Sprintf("  └─ Split %db (Lines %d-%d) processing...", batch.BatchIndex+1, mid+1, len(batch.Lines)))
			resultB, errB := p.translateBatchWithRetry(ctx, batchB, depth+1)
			if errB != nil {
				return nil, fmt.Errorf("split B failed: %w", errB)
			}

			// Merge results
			return append(resultA, resultB...), nil
		}

		// Max depth reached, fail
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("desync: count mismatch after %d splits", depth)
	}

	// Apply translations
	for _, resp := range response {
		if resp.ID >= 0 && resp.ID < len(translatedLines) {
			translatedLines[resp.ID].Text = resp.Text
			// Cache the translation
			p.Cache.SaveTranslation(batch.Lines[resp.ID].Text, resp.Text, langPair)
		}
	}

	// Quality Gate: Run linter on translated lines
	if depth == 0 { // Only lint at top level to avoid retry loops
		lintResult := p.lintTranslation(translatedLines)
		if !lintResult.PassedAll && len(lintResult.Issues) > 0 {
			highSeverityCount := 0
			for _, issue := range lintResult.Issues {
				if issue.Severity == linter.SeverityHigh {
					highSeverityCount++
				}
			}

			// Only retry if high severity issues found
			if highSeverityCount > 0 && depth < maxRetryDepth {
				p.log(fmt.Sprintf("  Quality Gate: %d HIGH severity issues, retrying...", highSeverityCount))
				return p.translateBatchWithRetry(ctx, batch, depth+1)
			}
		}
	}

	return translatedLines, nil
}

// buildSystemPrompt creates system prompt with sliding window context and glossary
func (p *Pipeline) buildSystemPrompt(contextLines []parser.SubtitleLine) string {
	prompt := p.Config.SystemPrompt

	// Inject glossary
	if len(p.Config.Glossary) > 0 {
		glossaryText := "\n\nGlossary (preserve these terms exactly as specified):\n"
		for orig, trans := range p.Config.Glossary {
			glossaryText += fmt.Sprintf("- \"%s\" -> \"%s\"\n", orig, trans)
		}
		prompt = strings.Replace(prompt, "{{glossary}}", glossaryText, 1)
	} else {
		prompt = strings.Replace(prompt, "{{glossary}}", "", 1)
	}

	// Add sliding window context (passive context from previous batch)
	// Per spec: last 3 lines of Batch N appended as read-only context at start of Batch N+1
	if len(contextLines) > 0 {
		contextText := "\n\n---\nPASSIVE CONTEXT (Previous lines for reference - DO NOT translate these):\n"
		for i, line := range contextLines {
			contextText += fmt.Sprintf("%d. %s\n", i+1, line.Text)
		}
		contextText += "---\n"
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

// LoadResumeState loads .bakasub.temp if it exists in the video's directory
func LoadResumeState(path string) (*ResumeState, error) {
	// If path points directly to a .temp file, use it
	// Otherwise, assume it's a video path and look for the temp file
	tempPath := path
	if !strings.HasSuffix(path, ".bakasub.temp") {
		tempPath = filepath.Join(filepath.Dir(path), ".bakasub.temp")
	}

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

// lintTranslation runs quality checks on translated lines
func (p *Pipeline) lintTranslation(lines []parser.SubtitleLine) linter.Result {
	// Extract text from lines
	texts := make([]string, len(lines))
	for i, line := range lines {
		texts[i] = line.Text
	}

	// Build lint options
	opts := linter.CheckOptions{
		SourceLang: p.Config.SourceLang,
		TargetLang: p.Config.TargetLang,
		Glossary:   p.Config.Glossary,
	}

	return linter.Check(texts, opts)
}
