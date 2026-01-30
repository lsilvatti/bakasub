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
	"github.com/lsilvatti/bakasub/internal/locales"
	"github.com/lsilvatti/bakasub/internal/ui/components"
	"github.com/lsilvatti/bakasub/internal/ui/layout"
	"github.com/lsilvatti/bakasub/internal/ui/styles"
)

type ViewState int

const (
	ViewMain              ViewState = iota
	ViewDirectoryDetected           // New: Modal asking batch or single file
	ViewConflictResolution
	ViewDryRunReport
	ViewGlossaryEditor
)

// Available media types for prompt selection
var mediaTypes = []string{"anime", "movie", "series", "documentary", "youtube"}

// Available mux modes
var muxModes = []string{"replace", "new-file"}

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

	// Interactive selection indices
	mediaTypeIdx int
	muxModeIdx   int

	// Directory detection state
	showDirModal   bool
	dirMKVCount    int
	dirIsDirectory bool

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

	// Find initial indices for media type and mux mode
	mediaTypeIdx := 0
	for i, mt := range mediaTypes {
		if mt == jobConfig.MediaType {
			mediaTypeIdx = i
			break
		}
	}

	muxModeIdx := 0
	for i, mm := range muxModes {
		if mm == jobConfig.MuxMode {
			muxModeIdx = i
			break
		}
	}

	return Model{
		cfg:             cfg,
		state:           ViewMain,
		jobConfig:       jobConfig,
		canStart:        false,
		mediaTypeIdx:    mediaTypeIdx,
		muxModeIdx:      muxModeIdx,
		analysisSpinner: components.NewNeonSpinner(),
	}
}

// SetSize updates the model dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Init() tea.Cmd {
	// Request terminal size and start spinner and analysis together
	return tea.Batch(tea.WindowSize(), m.analyzeDirectory, m.analysisSpinner.Start())
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

	case MsgDirectoryDetected:
		m.analyzing = false
		m.analysisSpinner.Stop()
		m.state = ViewDirectoryDetected
		m.showDirModal = true
		m.dirMKVCount = msg.MKVCount
		m.dirIsDirectory = msg.IsDir
		return m, nil

	case MsgBatchModeSelected:
		m.showDirModal = false
		m.state = ViewMain
		m.analyzing = true

		// Collect MKV files based on batch mode selection
		return m, func() tea.Msg {
			entries, err := os.ReadDir(m.jobConfig.InputPath)
			if err != nil {
				return MsgAnalysisComplete{Success: false, Error: err}
			}

			var mkvFiles []string
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".mkv") {
					mkvFiles = append(mkvFiles, filepath.Join(m.jobConfig.InputPath, entry.Name()))
				}
			}

			if msg.BatchMode {
				// Batch mode: analyze all files
				m.jobConfig.BatchMode = true
				return m.analyzeFiles(mkvFiles)
			}
			// Single mode: return to caller for file picker
			return MsgSelectSingleFile{Files: mkvFiles}
		}

	case MsgSingleFileSelected:
		m.jobConfig.BatchMode = false
		m.jobConfig.InputPath = msg.Path
		m.analyzing = true
		return m, func() tea.Msg {
			return m.analyzeFiles([]string{msg.Path})
		}

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
		// Send the StartJobMsg to parent (dashboard) with the job config
		return m, func() tea.Msg {
			return StartJobMsg{JobConfig: m.jobConfig}
		}

	case MsgCancelJob:
		// Send the CancelledMsg to parent (dashboard)
		return m, func() tea.Msg {
			return CancelledMsg{}
		}
	}

	return m.updateSubModels(msg)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewDirectoryDetected:
		switch msg.String() {
		case "b":
			// Batch mode: process all files
			return m, func() tea.Msg {
				return MsgBatchModeSelected{BatchMode: true}
			}
		case "s":
			// Single file mode: open file picker
			return m, func() tea.Msg {
				return MsgBatchModeSelected{BatchMode: false}
			}
		case "esc":
			return m, func() tea.Msg { return MsgCancelJob{} }
		}
		return m, nil

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
		case msg.String() == "m":
			// Cycle media type
			m.mediaTypeIdx = (m.mediaTypeIdx + 1) % len(mediaTypes)
			m.jobConfig.MediaType = mediaTypes[m.mediaTypeIdx]
			return m, nil
		case msg.String() == "M":
			// Cycle media type backwards
			m.mediaTypeIdx = (m.mediaTypeIdx - 1 + len(mediaTypes)) % len(mediaTypes)
			m.jobConfig.MediaType = mediaTypes[m.mediaTypeIdx]
			return m, nil
		case msg.String() == "x":
			// Cycle mux mode
			m.muxModeIdx = (m.muxModeIdx + 1) % len(muxModes)
			m.jobConfig.MuxMode = muxModes[m.muxModeIdx]
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
		if key.Matches(msg, keys.Enter) {
			// Proceed to execution from dry run report
			return m, func() tea.Msg { return MsgStartJob{} }
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

		// If multiple MKV files found, show directory detection modal
		if len(mkvFiles) > 1 {
			return MsgDirectoryDetected{
				Path:     m.jobConfig.InputPath,
				MKVCount: len(mkvFiles),
				IsDir:    true,
			}
		}

		m.jobConfig.BatchMode = false
	} else {
		mkvFiles = []string{m.jobConfig.InputPath}
		m.jobConfig.BatchMode = false
	}

	if len(mkvFiles) == 0 {
		return MsgAnalysisComplete{Success: false, Error: fmt.Errorf("no MKV files found")}
	}

	return m.analyzeFiles(mkvFiles)
}

func (m Model) analyzeFiles(mkvFiles []string) tea.Msg {
	var files []AnalyzedFile

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

	// Collect subtitle tracks
	var subTracks []media.Track
	for _, track := range fileInfo.Tracks {
		if track.Type == "subtitles" {
			subTracks = append(subTracks, track)
		}
	}

	// Check for conflicts: multiple tracks with the same language code
	// This happens when there are multiple subtitle options (e.g., "eng" full dialogue + "eng" signs only)
	langCounts := make(map[string]int)
	for _, track := range subTracks {
		langCounts[track.Language]++
	}

	// Find tracks that conflict (same language, multiple options)
	var conflictTracks []media.Track
	for _, track := range subTracks {
		if langCounts[track.Language] > 1 {
			conflictTracks = append(conflictTracks, track)
		}
	}

	analyzed.HasConflict = len(conflictTracks) > 1
	analyzed.ConflictTracks = conflictTracks

	// Auto-select if only one subtitle track exists, or no conflicts
	if len(subTracks) == 1 {
		analyzed.SelectedTrackID = subTracks[0].ID
		analyzed.HasConflict = false
	} else if len(conflictTracks) == 0 && len(subTracks) > 0 {
		// Multiple tracks but different languages - auto-select first one
		// Could be improved to select based on source language preference
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
	// Wait for terminal size
	if layout.IsWaitingForSize(m.width, m.height) {
		return locales.T("common.loading")
	}

	// Check if terminal is too small
	if layout.IsTooSmall(m.width, m.height) {
		return layout.RenderTooSmallWarning(m.width, m.height)
	}

	switch m.state {
	case ViewDirectoryDetected:
		return m.renderWithOverlay(m.renderDirectoryModal())
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
	contentWidth := m.width - 4

	if m.analyzing {
		spinnerContent := m.analysisSpinner.ViewWithCustomLabel(locales.T("common.loading"))
		return styles.MainWindow.Width(contentWidth).Render(spinnerContent)
	}

	if m.err != nil {
		errContent := styles.ErrorStyle.Render(fmt.Sprintf("%s: %v", locales.T("common.error"), m.err))
		return styles.MainWindow.Width(contentWidth).Render(errContent)
	}

	var s strings.Builder

	// Header bar
	headerBar := strings.Repeat("▒", contentWidth-4)
	s.WriteString(headerBar + "\n")

	title := styles.TitleStyle.Render(locales.T("job.title"))
	if m.jobConfig.BatchMode {
		title += styles.SubtleStyle.Render(fmt.Sprintf(" (%s)", locales.Tf("job.batch_label", len(m.jobConfig.Files))))
	}
	s.WriteString(title + "\n")
	s.WriteString(headerBar + "\n\n")

	// Section 1: Extraction Strategy
	s.WriteString(styles.SectionStyle.Render("1. "+locales.T("job.extraction.title")) + "\n")

	// Show subtitle source selection
	if m.hasConflicts {
		s.WriteString(styles.WarningStyle.Render("  "+locales.T("job.extraction.multiple_tracks_warning")) + " ")
		s.WriteString(styles.KeyHintStyle.Render("[ r ]") + " " + locales.T("job.extraction.resolve_button") + "\n")
		s.WriteString("      " + locales.T("job.extraction.resolve_help") + "\n")
	} else if len(m.jobConfig.Files) > 0 && len(m.jobConfig.Files[0].Tracks) > 0 {
		// Show selected track info
		selectedTrack := locales.T("job.extraction.auto_detect")
		for _, file := range m.jobConfig.Files {
			if file.SelectedTrackID >= 0 {
				for _, track := range file.Tracks {
					if track.ID == file.SelectedTrackID && track.Type == "subtitles" {
						selectedTrack = fmt.Sprintf("Track %d (%s) - %s", track.ID, track.Language, track.Codec)
						break
					}
				}
				break
			}
		}
		s.WriteString(fmt.Sprintf("  %s [ %s ]\n", locales.T("job.extraction.subtitle_source"), selectedTrack))
	} else {
		s.WriteString(fmt.Sprintf("  %s [ %s ]\n", locales.T("job.extraction.subtitle_source"), locales.T("job.extraction.auto_detect")))
	}

	// Show audio reference
	if len(m.jobConfig.Files) > 0 && len(m.jobConfig.Files[0].Tracks) > 0 {
		audioTrack := "None"
		for _, track := range m.jobConfig.Files[0].Tracks {
			if track.Type == "audio" {
				audioTrack = fmt.Sprintf("Track %d (%s) - %s", track.ID, track.Language, track.Codec)
				break
			}
		}
		s.WriteString(fmt.Sprintf("  %s [ %s ] (For context)\n", locales.T("job.extraction.audio_reference"), audioTrack))
	}

	// Show extract fonts option
	extractFonts := "✓"
	if !m.jobConfig.ExtractFonts {
		extractFonts = " "
	}
	s.WriteString(fmt.Sprintf("  [%s] %s\n", extractFonts, locales.T("job.extraction.extract_fonts")))
	s.WriteString("\n")

	// Section 2: Translation Context
	s.WriteString(styles.SectionStyle.Render("2. "+locales.T("job.translation.title")) + "\n")

	// Media Type - interactive selector
	mediaTypeDisplay := locales.T("job.media_types." + m.jobConfig.MediaType)
	s.WriteString(fmt.Sprintf("  %s %s  ", locales.T("job.translation.media_type"), styles.AccentStyle.Render("[ "+mediaTypeDisplay+" ]")))
	s.WriteString(styles.KeyHintStyle.Render("[ m ]") + " " + locales.T("common.next") + "\n")

	// Target Language
	s.WriteString(fmt.Sprintf("  %s %s\n", locales.T("job.translation.target_lang"), m.jobConfig.TargetLang))

	// Glossary
	glossaryTerms := locales.Tf("job.translation.glossary_terms", len(m.jobConfig.GlossaryTerms))
	s.WriteString(fmt.Sprintf("  %s %s\n", locales.T("job.translation.glossary"), glossaryTerms))
	s.WriteString("\n")

	// Section 3: Muxing Output
	s.WriteString(styles.SectionStyle.Render("3. "+locales.T("job.muxing.title")) + "\n")

	// Mux Mode - interactive selector
	muxModeDisplay := locales.T("job.mux_modes." + m.jobConfig.MuxMode)
	s.WriteString(fmt.Sprintf("  %s %s  ", locales.T("job.muxing.mode"), styles.AccentStyle.Render("[ "+muxModeDisplay+" ]")))
	s.WriteString(styles.KeyHintStyle.Render("[ x ]") + " " + locales.T("common.next") + "\n")
	s.WriteString("\n")

	// Cost Estimation Box
	s.WriteString(styles.InfoBoxStyle.Width(contentWidth-6).Render(fmt.Sprintf(
		"[i] %s\n"+
			"%s %s  |  %s %dk  |  %s $%.2f",
		locales.T("job.estimation.title"),
		locales.T("job.estimation.model"), m.jobConfig.AIModel,
		locales.T("job.estimation.tokens"), m.jobConfig.EstimatedTokens/1000,
		locales.T("job.estimation.cost"), m.jobConfig.EstimatedCost,
	)) + "\n\n")

	// Footer with keybindings
	footer := ""
	footer += styles.KeyHintStyle.Render("[ESC]") + " " + locales.T("job.footer.back") + "      "
	footer += styles.KeyHintStyle.Render("[ d ]") + " " + locales.T("job.footer.simulation") + "      "
	footer += styles.KeyHintStyle.Render("[ g ]") + " " + locales.T("job.translation.glossary") + "      "

	if m.canStart {
		footer += styles.SuccessStyle.Render("[ENTER] " + locales.T("job.footer.start"))
	} else {
		footer += styles.DisabledStyle.Render("[ " + locales.T("job.footer.start_disabled") + " ]")
	}

	s.WriteString(footer)

	return styles.MainWindow.Width(contentWidth).Render(s.String())
}

func (m Model) renderWithOverlay(overlay string) string {
	// Use dimmed background with modal overlay
	if m.width > 0 && m.height > 0 {
		return styles.RenderModalWithOverlay(overlay, m.width, m.height)
	}
	// Fallback to simple overlay
	base := m.renderMain()
	return base + "\n\n" + overlay
}

func (m Model) renderDirectoryModal() string {
	var s strings.Builder

	titleStyle := styles.TitleStyle.
		Width(54)

	contentStyle := styles.Panel.
		Width(54).
		Padding(1, 2)

	s.WriteString(titleStyle.Render(locales.T("job.directory_detected")))
	s.WriteString("\n")

	content := fmt.Sprintf(
		"%s %s\n\n"+
			"%s\n"+
			"• %s %d\n"+
			"• %s 00\n\n"+
			"%s\n\n"+
			"  %s "+locales.T("job.process_batch")+"\n"+
			"      "+locales.T("job.batch_note")+"\n\n"+
			"  %s "+locales.T("job.select_single_file")+"\n"+
			"      "+locales.T("job.single_file_note")+"\n\n"+
			"  %s "+locales.T("common.cancel"),
		locales.T("job.path_label"),
		filepath.Base(m.jobConfig.InputPath),
		locales.T("job.analysis"),
		locales.T("job.mkv_files_found"),
		m.dirMKVCount,
		locales.T("job.subtitles_found"),
		locales.T("job.how_to_proceed"),
		styles.KeyHintStyle.Render("[ b ]"),
		m.dirMKVCount,
		styles.KeyHintStyle.Render("[ s ]"),
		styles.KeyHintStyle.Render("[ESC]"),
	)

	s.WriteString(contentStyle.Render(content))

	return styles.ModalStyle.Width(60).Render(s.String())
}
