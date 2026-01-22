package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/media"
	"github.com/lsilvatti/bakasub/internal/ui/components"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type ViewState int

const (
	ViewMain ViewState = iota
	ViewConflictResolution
	ViewDryRunReport
	ViewGlossaryEditor
)

type Model struct {
	cfg    *config.Config
	width  int
	height int

	state           ViewState
	jobConfig       JobConfig
	analyzing       bool
	analysisSpinner components.NeonSpinner
	hasConflicts    bool
	canStart        bool
	selectedFileIdx int

	conflictModal  *ConflictModal
	dryRunReport   *DryRunReport
	glossaryEditor *GlossaryEditor

	err error
}

type KeyMap struct {
	Enter    key.Binding
	Escape   key.Binding
	DryRun   key.Binding
	Glossary key.Binding
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	Space    key.Binding
	Resolve  key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start job")),
		Escape:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back/cancel")),
		DryRun:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "dry run simulation")),
		Glossary: key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "edit glossary")),
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		Space:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
		Resolve:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "resolve conflict")),
	}
}

var keys = DefaultKeyMap()

func New(cfg *config.Config, inputPath string) Model {
	jobConfig := JobConfig{
		InputPath:       inputPath,
		TargetLang:      cfg.TargetLang,
		MediaType:       "anime",
		AIModel:         cfg.Model,
		Temperature:     cfg.Temperature,
		RemoveHITags:    true,
		MuxMode:         "replace",
		SetDefault:      true,
		BackupOriginal:  true,
		ExtractFonts:    true,
		AutoDetectTrack: true,
		GlossaryTerms:   make(map[string]string),
	}

	return Model{
		cfg:             cfg,
		state:           ViewMain,
		jobConfig:       jobConfig,
		canStart:        false,
		analysisSpinner: components.NewNeonSpinner(),
	}
}

func (m Model) Init() tea.Cmd {
	// Start spinner and analysis together
	return tea.Batch(m.analyzeDirectory, m.analysisSpinner.Start())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update spinner
	var spinnerCmd tea.Cmd
	m.analysisSpinner, spinnerCmd = m.analysisSpinner.Update(msg)
	if spinnerCmd != nil {
		cmds = append(cmds, spinnerCmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		model, cmd := m.handleKeyPress(msg)
		m = model.(Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case MsgAnalysisComplete:
		m.analyzing = false
		m.analysisSpinner.Stop() // Stop spinner when analysis completes
		if msg.Success {
			m.jobConfig.Files = msg.Files
			m.hasConflicts = m.checkConflicts()
			m.canStart = !m.hasConflicts
			cmds = append(cmds, m.loadGlossary)
			return m, tea.Batch(cmds...)
		}
		m.err = msg.Error
		return m, tea.Batch(cmds...)

	case MsgConflictResolved:
		if msg.FileIndex >= 0 && msg.FileIndex < len(m.jobConfig.Files) {
			m.jobConfig.Files[msg.FileIndex].SelectedTrackID = msg.TrackID
			m.jobConfig.Files[msg.FileIndex].HasConflict = false
			m.hasConflicts = m.checkConflicts()
			m.canStart = !m.hasConflicts
			m.state = ViewMain
		}
		return m, nil

	case MsgCostEstimated:
		m.jobConfig.EstimatedChars = msg.TotalChars
		m.jobConfig.EstimatedTokens = msg.TokenCount
		m.jobConfig.EstimatedCost = msg.EstimatedCost
		return m, nil

	case MsgDryRunComplete:
		if m.dryRunReport != nil {
			m.dryRunReport.CanWrite = msg.CanWrite
			m.dryRunReport.TokenCount = msg.TokenCount
			m.dryRunReport.EstimatedCost = msg.EstimatedCost
			m.dryRunReport.Warnings = msg.Warnings
		}
		return m, nil

	case MsgGlossaryLoaded:
		m.jobConfig.GlossaryTerms = msg.Terms
		m.jobConfig.GlossaryPath = msg.Path
		return m, m.estimateCost

	case MsgStartJob:
		return m, nil

	case MsgCancelJob:
		return m, nil
	}

	return m.updateSubModels(msg)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewMain:
		switch {
		case key.Matches(msg, keys.Escape):
			return m, func() tea.Msg { return MsgCancelJob{} }
		case key.Matches(msg, keys.Enter):
			if m.canStart {
				return m, func() tea.Msg { return MsgStartJob{} }
			}
			return m, nil
		case key.Matches(msg, keys.DryRun):
			m.state = ViewDryRunReport
			m.dryRunReport = NewDryRunReport(m.jobConfig)
			return m, m.runDryRun
		case key.Matches(msg, keys.Glossary):
			m.state = ViewGlossaryEditor
			m.glossaryEditor = NewGlossaryEditor(m.jobConfig.GlossaryPath, m.jobConfig.GlossaryTerms)
			return m, nil
		case key.Matches(msg, keys.Resolve):
			if m.hasConflicts {
				for i, file := range m.jobConfig.Files {
					if file.HasConflict {
						m.selectedFileIdx = i
						m.state = ViewConflictResolution
						m.conflictModal = NewConflictModal(file)
						return m, nil
					}
				}
			}
			return m, nil
		}
	case ViewConflictResolution:
		if m.conflictModal != nil {
			updated, cmd := m.conflictModal.Update(msg)
			m.conflictModal = &updated
			return m, cmd
		}
	case ViewDryRunReport:
		if key.Matches(msg, keys.Escape) {
			m.state = ViewMain
			return m, nil
		}
	case ViewGlossaryEditor:
		if m.glossaryEditor != nil {
			updated, cmd := m.glossaryEditor.Update(msg)
			m.glossaryEditor = &updated
			if updated.Closed {
				m.state = ViewMain
				if updated.Modified {
					m.jobConfig.GlossaryTerms = updated.Terms
					return m, m.saveGlossary
				}
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) updateSubModels(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewConflictResolution:
		if m.conflictModal != nil {
			updated, cmd := m.conflictModal.Update(msg)
			m.conflictModal = &updated
			return m, cmd
		}
	case ViewDryRunReport:
		if m.dryRunReport != nil {
			updated, cmd := m.dryRunReport.Update(msg)
			m.dryRunReport = &updated
			return m, cmd
		}
	case ViewGlossaryEditor:
		if m.glossaryEditor != nil {
			updated, cmd := m.glossaryEditor.Update(msg)
			m.glossaryEditor = &updated
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) checkConflicts() bool {
	for _, file := range m.jobConfig.Files {
		if file.HasConflict {
			return true
		}
	}
	return false
}

func (m Model) analyzeDirectory() tea.Msg {
	var files []AnalyzedFile

	info, err := os.Stat(m.jobConfig.InputPath)
	if err != nil {
		return MsgAnalysisComplete{Success: false, Error: err}
	}

	var mkvFiles []string
	if info.IsDir() {
		entries, err := os.ReadDir(m.jobConfig.InputPath)
		if err != nil {
			return MsgAnalysisComplete{Success: false, Error: err}
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".mkv") {
				mkvFiles = append(mkvFiles, filepath.Join(m.jobConfig.InputPath, entry.Name()))
			}
		}
		m.jobConfig.BatchMode = len(mkvFiles) > 1
	} else {
		mkvFiles = []string{m.jobConfig.InputPath}
		m.jobConfig.BatchMode = false
	}

	if len(mkvFiles) == 0 {
		return MsgAnalysisComplete{Success: false, Error: fmt.Errorf("no MKV files found")}
	}

	for _, path := range mkvFiles {
		analyzed, err := m.analyzeFile(path)
		if err != nil {
			continue
		}
		files = append(files, analyzed)
	}

	return MsgAnalysisComplete{Files: files, Success: len(files) > 0, Error: nil}
}

func (m Model) analyzeFile(path string) (AnalyzedFile, error) {
	fileInfo, err := media.Analyze(path)
	if err != nil {
		return AnalyzedFile{}, err
	}

	analyzed := AnalyzedFile{
		Path:            path,
		Filename:        filepath.Base(path),
		Tracks:          fileInfo.Tracks,
		Attachments:     fileInfo.Attachments,
		SelectedTrackID: -1,
	}

	var subTracks []media.Track
	for _, track := range fileInfo.Tracks {
		if track.Type == "subtitles" {
			subTracks = append(subTracks, track)
		}
	}

	analyzed.HasConflict = len(subTracks) > 1
	analyzed.ConflictTracks = subTracks

	if len(subTracks) == 1 {
		analyzed.SelectedTrackID = subTracks[0].ID
		analyzed.HasConflict = false
	}

	return analyzed, nil
}

func (m Model) estimateCost() tea.Msg {
	totalChars := len(m.jobConfig.Files) * 10000
	tokenCount := totalChars / 4
	pricePerM := 0.15
	if m.jobConfig.ModelPricePerM > 0 {
		pricePerM = m.jobConfig.ModelPricePerM
	}
	estimatedCost := (float64(tokenCount) / 1000000) * pricePerM

	return MsgCostEstimated{
		TotalChars:    totalChars,
		EstimatedCost: estimatedCost,
		TokenCount:    tokenCount,
	}
}

func (m Model) loadGlossary() tea.Msg {
	dir := m.jobConfig.InputPath
	if !filepath.IsAbs(dir) {
		dir, _ = filepath.Abs(dir)
	}

	info, err := os.Stat(dir)
	if err == nil && !info.IsDir() {
		dir = filepath.Dir(dir)
	}

	glossaryPath := filepath.Join(dir, "glossary.json")
	terms, err := LoadGlossaryFromFile(glossaryPath)
	if err != nil {
		return MsgGlossaryLoaded{Terms: make(map[string]string), Path: glossaryPath}
	}

	return MsgGlossaryLoaded{Terms: terms, Path: glossaryPath}
}

func (m Model) saveGlossary() tea.Msg {
	err := SaveGlossaryToFile(m.jobConfig.GlossaryPath, m.jobConfig.GlossaryTerms)
	return MsgGlossarySaved{Success: err == nil, Error: err}
}

func (m Model) runDryRun() tea.Msg {
	warnings := []string{}

	dir := m.jobConfig.InputPath
	info, err := os.Stat(dir)
	if err == nil && !info.IsDir() {
		dir = filepath.Dir(dir)
	}

	canWrite := true
	testFile := filepath.Join(dir, ".bakasub_writetest")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		canWrite = false
		warnings = append(warnings, "No write permission in output directory")
	} else {
		os.Remove(testFile)
	}

	for _, file := range m.jobConfig.Files {
		if strings.Contains(strings.ToLower(file.Filename), "broken") {
			warnings = append(warnings, fmt.Sprintf("File '%s' may have timestamp issues", file.Filename))
		}
	}

	return MsgDryRunComplete{
		CanWrite:      canWrite,
		TokenCount:    m.jobConfig.EstimatedTokens,
		EstimatedCost: m.jobConfig.EstimatedCost,
		Warnings:      warnings,
	}
}

func (m Model) View() string {
	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	switch m.state {
	case ViewConflictResolution:
		if m.conflictModal != nil {
			return m.renderWithOverlay(m.conflictModal.View())
		}
	case ViewDryRunReport:
		if m.dryRunReport != nil {
			return m.dryRunReport.View()
		}
	case ViewGlossaryEditor:
		if m.glossaryEditor != nil {
			return m.glossaryEditor.View()
		}
	}
	return m.renderMain()
}

func (m Model) renderMain() string {
	if m.analyzing {
		// Show spinner while analyzing
		return styles.AppStyle.Render(m.analysisSpinner.ViewWithCustomLabel("Analyzing directory..."))
	}

	if m.err != nil {
		return styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var s strings.Builder

	title := styles.TitleStyle.Render("JOB SETUP")
	if m.jobConfig.BatchMode {
		title += styles.SubtleStyle.Render(fmt.Sprintf(" (Batch: %d Files)", len(m.jobConfig.Files)))
	}
	s.WriteString(title + "\n\n")

	s.WriteString(styles.SectionStyle.Render("1. EXTRACTION STRATEGY") + "\n")

	if m.hasConflicts {
		s.WriteString(styles.WarningStyle.Render("  [!] MULTIPLE TRACKS FOUND") + " ")
		s.WriteString(styles.KeyHintStyle.Render("[ r ]") + " RESOLVE\n")
	} else {
		s.WriteString("  ✓ Tracks auto-detected\n")
	}
	s.WriteString("\n")

	s.WriteString(styles.SectionStyle.Render("2. TRANSLATION CONTEXT") + "\n")
	s.WriteString(fmt.Sprintf("  Media Type: %s\n", m.jobConfig.MediaType))
	s.WriteString(fmt.Sprintf("  Target Lang: %s\n", m.jobConfig.TargetLang))
	s.WriteString(fmt.Sprintf("  Glossary: %d terms\n", len(m.jobConfig.GlossaryTerms)))
	s.WriteString("\n")

	s.WriteString(styles.SectionStyle.Render("3. MUXING OUTPUT") + "\n")
	s.WriteString(fmt.Sprintf("  Mode: %s\n", m.jobConfig.MuxMode))
	s.WriteString("\n")

	s.WriteString(styles.InfoBoxStyle.Render(fmt.Sprintf(
		"[i] COST ESTIMATION\n"+
			"Model: %s  |  Est. Tokens: %dk  |  Est. Cost: $%.2f",
		m.jobConfig.AIModel,
		m.jobConfig.EstimatedTokens/1000,
		m.jobConfig.EstimatedCost,
	)) + "\n\n")

	footer := ""
	footer += styles.KeyHintStyle.Render("[ESC]") + " BACK      "
	footer += styles.KeyHintStyle.Render("[ d ]") + " SIMULATION (DRY RUN)      "
	footer += styles.KeyHintStyle.Render("[ g ]") + " GLOSSARY      "

	if m.canStart {
		footer += styles.SuccessStyle.Render("[ENTER] START JOB")
	} else {
		footer += styles.DisabledStyle.Render("[ START DISABLED - RESOLVE CONFLICTS ]")
	}

	s.WriteString(footer)

	return styles.AppStyle.Render(s.String())
}

func (m Model) renderWithOverlay(overlay string) string {
	base := m.renderMain()
	return base + "\n\n" + overlay
}
